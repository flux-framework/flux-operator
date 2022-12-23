#!/bin/bash

NAME=${1:test}
NAMESPACE=${2:-flux-operator}

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT=$(dirname ${HERE})
TEST_DIR="${ROOT}/examples/tests/${NAME}"

echo "Namespace: ${NAMESPACE}"
kubectl get -n ${NAMESPACE} pod
pods=$(kubectl get -n ${NAMESPACE} pod --output=jsonpath={.items..metadata.name}); 
echo "Pods: ${pods}"
pod="${pods%% *}"
echo "Pod: ${pod}"

# Prepare actual and tested comparison
expected=${TEST_DIR}/test.out.correct
actual=${TEST_DIR}/test.out
kubectl logs -n ${NAMESPACE} ${pod} -f > ${actual} 2>&1

echo "Actual:"
cat ${actual}

# Only run if we've derived this file
if [[ -e "${expected}" ]]; then
    echo "Expected:"
    cat ${expected}
    diff ${expected} ${actual}
fi

# Ensure all containers exit code 0
for exitcode in $(kubectl get -n flux-operator pod --output=jsonpath={.items...containerStatuses..state.terminated.exitCode}); do
   if [[ "${exitcode}" != "0" ]]; then
       echo "Container in ${NAME} had nonzero exit code"
       exit 1
    fi 
done