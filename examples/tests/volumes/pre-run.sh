#!/bin/bash

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "Copying local volume to /tmp/data-volumes in minikube"

# The "data" volume will be mounted at /mnt/data
minikube ssh -- mkdir -p /tmp/data
minikube cp ${HERE}/data/pancakes.txt /tmp/data/pancakes.txt
minikube ssh ls /tmp/data
