#!/bin/bash

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
TESTS=$(dirname ${HERE})

# The "data" volume will be mounted at /mnt/data
minikube ssh -- mkdir -p /tmp/data