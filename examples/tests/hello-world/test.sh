#!/bin/bash

NAMESPACE=${1:-flux-operator}
echo "Namespace: ${NAMESPACE}"
kubectl get -n ${NAMESPACE} pods
pods=$(kubectl get -n ${NAMESPACE} pod --output=jsonpath={.items..metadata.name}); 
echo "Pods: ${pods}"
pod="${pods%% *}"
echo "Pod: ${pod}"

kubectl logs -n ${NAMESPACE} ${pod}