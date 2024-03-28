#!/bin/bash

docker build -t gozmq .
kind load docker-image gozmq
kubectl delete -f minicluster.yaml
kubectl apply -f minicluster.yaml