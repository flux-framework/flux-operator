#!/bin/bash

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "Copying local volume to /tmp/data-volumes in minikube"

# We don't care if this works or not - mkdir -p seems to bork
minikube ssh -- mkdir -p /tmp/data-volumes
minikube cp ${HERE}/data/pancakes.txt /tmp/data-volumes/pancakes.txt
minikube ssh ls /tmp/data-volumes