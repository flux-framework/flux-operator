#!/bin/bash

# This is a test script that, instead of applying a yaml file,
# brings up the MiniCluster and runs Python tests using it

# Usage: /bin/bash script/test.sh $name 30
HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT=$(dirname ${HERE})
cd ${ROOT}

set -eEu -o pipefail

cmd="$@"

echo "Testing Command: ${cmd}"

# Quick helper script to run a test
make clean >> /dev/null
make run > /dev/null 2> /dev/null &
pid=$!
echo "PID for running cluster is ${pid}"

# If there is a pre-run script
${cmd}

kill ${pid} || true
kill $(lsof -t -i:8080) || true
