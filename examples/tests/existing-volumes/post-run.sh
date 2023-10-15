#!/bin/bash

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "Cleaning up /tmp/data in minikube"

# pods usually need to be deleted before the pvc/pv
kubectl delete -f ${HERE}/minicluster.yaml
kubectl delete pods --all --grace-period=0 --force
kubectl delete pvc --all --grace-period=0 --force
kubectl delete pv --all --grace-period=0 --force
minikube ssh -- sudo rm -rf /tmp/data
