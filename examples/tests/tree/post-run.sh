#!/bin/bash

echo "Checking log in minikube"
minikube ssh -- cat /tmp/data/tree.out
echo "Cleaning up /tmp/data in minikube"
minikube ssh -- sudo rm -rf /tmp/data

