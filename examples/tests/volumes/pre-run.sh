#!/bin/bash

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
TESTS=$(dirname ${HERE})

echo "Copying local volume to /tmp/data-volumes in minikube"

# The "data" volume will be mounted at /mnt/data
minikube ssh -- mkdir -p /data
minikube cp ${TESTS}/data/pancakes.txt /data/pancakes.txt
minikube ssh ls /data
