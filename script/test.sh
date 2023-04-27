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
out="${ROOT}/examples/tests/${name}/${name}-log.out"
err="${ROOT}/examples/tests/${name}/${name}-log.err"

# Quick helper script to run a test
make clean >> /dev/null 2>&1
make run > ${out} 2> ${err} &
pid=$!
echo "PID for running cluster is ${pid}"

# If there is a pre-run script
/bin/bash examples/tests/${name}/pre-run.sh || true
kubectl apply -f examples/tests/${name}/minicluster.yaml
echo "Sleeping for ${jobtime} seconds to allow job to complete 😴️."
sleep ${jobtime}
/bin/bash ${HERE}/check-output.sh ${name} || (
    echo "Tests for ${name} were not successful"
    kill ${pid} || echo "I am already dead 😭️"
    echo "$out"
    echo "$err"
    kill $(lsof -t -i:8080) || true
    /bin/bash examples/tests/${name}/post-run.sh || true
    exit 1;
)
kill ${pid} || true
kill $(lsof -t -i:8080) || true
/bin/bash examples/tests/${name}/post-run.sh || true
