#!/bin/bash

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
TESTS=$(dirname ${HERE})

echo "Copying local volume to /tmp/data for kind"
mkdir -p /tmp/data

# The "data" volume will be mounted at /mnt/data
cp ${TESTS}/data/pancakes.txt /tmp/data/pancakes.txt
ls /tmp/data
kubectl apply -f ${HERE}/pv.yaml 
kubectl apply -f ${HERE}/pvc.yaml