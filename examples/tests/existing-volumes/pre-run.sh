#!/bin/bash

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
TESTS=$(dirname ${HERE})

echo "Copying local volume to /tmp/data-volumes in minikube"

# The "data" volume will be mounted at /mnt/data
minikube ssh -- mkdir -p /tmp/data
minikube cp ${TESTS}/data/pancakes.txt /tmp/data/pancakes.txt
minikube ssh ls /tmp/data
kubectl apply -f ${HERE}/pv.yaml 
kubectl apply -f ${HERE}/pvc.yaml