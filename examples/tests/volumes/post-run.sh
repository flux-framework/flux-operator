#!/bin/bash

echo "Cleaning up /tmp/data in minikube"
minikube ssh -- sudo rm -rf /tmp/data