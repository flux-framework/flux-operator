#!/bin/bash

echo "Cleaning up /data in minikube"
minikube ssh -- sudo rm -rf /data
