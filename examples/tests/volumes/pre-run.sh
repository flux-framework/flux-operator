#!/bin/bash

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "Copying local volume to /tmp/data-volumes in minikube"

# The "data" volume will be mounted at /mnt/data
minikube ssh -- mkdir -p /mnt/data
minikube cp ${HERE}/data/pancakes.txt /mnt/data/pancakes.txt
minikube ssh ls /mnt/data
