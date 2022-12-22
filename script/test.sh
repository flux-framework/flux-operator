#!/bin/bash

# Usage: /bin/bash script/test.sh $name 30
HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT=$(dirname ${HERE})
cd ${ROOT}

set -eEu -o pipefail

name=${1}
jobtime=${2:-30}

echo "   Name: ${name}"
echo "Jobtime: ${jobtime}"

# Output and error files
out="./examples/tests/${name}/${name}-log.out"
err="./examples/tests/${name}/${name}-log.err"

# Quick helper script to run a test
make clean >> /dev/null
make run > ${out} 2> ${err} &
pid=$!
echo "PID for running cluster is ${pid}"
kubectl apply -f examples/tests/${name}/minicluster-${name}.yaml
make list
sleep ${jobtime}
/bin/bash examples/tests/${name}/test.sh ${name} || (
    echo "Tests for ${name} were not successful"
    cat ${out}
    cat ${err}
    kill ${pid} || echo "I am already dead 😭️"
    kill $(lsof -t -i:8080) || true
    exit 1;
)
kill ${pid} || true
kill $(lsof -t -i:8080) || true