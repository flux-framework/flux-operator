#!/bin/bash
# Usage ./describe.sh <podname>
kubectl describe pod -n flux-operator ${1}
