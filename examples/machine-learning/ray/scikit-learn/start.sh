#!/bin/bash
export LC_ALL=C.UTF-8
export LANG=C.UTF-8

name=$(hostname)
echo "Hello I am ${name}"

# We only care about the head node address
address="flux-sample-0.flux-service.flux-operator.svc.cluster.local"
port="6379"
password="austinpowersyeahbaby"
echo "Head node address is ${address}:${port}"

# This is the "head noded"
if [[ "${name}" == "flux-sample-0" ]]; then
    echo "ray start --head --node-ip-address=${address} --port=6379 --redis-password=${password} --temp-dir=/tmp/workflow/tmp --disable-usage-stats"
    ray start --head --node-ip-address=${address} --port=${port} --redis-password=${password} --temp-dir=/tmp/workflow/tmp --disable-usage-stats
    sleep infinity
else
    # This triggers an error when I add port to the string, so we rely on using the default
    echo "ray start --address ${address}:${port} --redis-password=${password} --disable-usage-stats"
    ray start --address ${address}:${port} --redis-password=${password} --disable-usage-stats
    sleep infinity
fi