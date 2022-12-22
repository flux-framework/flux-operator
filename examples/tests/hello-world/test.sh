#!/bin/bash

NAME=${1:test}
NAMESPACE=${2:-flux-operator}

TEST_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "Namespace: ${NAMESPACE}"
kubectl get -n ${NAMESPACE} pods
pods=$(kubectl get -n ${NAMESPACE} pod --output=jsonpath={.items..metadata.name}); 
echo "Pods: ${pods}"
pod="${pods%% *}"
echo "Pod: ${pod}"

# Prepare actual and tested comparison
expected=${TEST_DIR}/test.out
actual=${TEST_DIR}/${NAME}-test.out
kubectl logs -n ${NAMESPACE} ${pod} -f > ${actual} 2>&1

echo "Expected:"
cat ${expected}

echo "Actual:"
cat ${actual}

diff ${expected} ${actual}