#!/bin/bash

echo "Cleaning up /tmp/data-volumes in minikube"
minikube ssh -- sudo rm -rf /tmp/data-volumes