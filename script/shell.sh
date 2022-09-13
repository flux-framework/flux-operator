#!/bin/bash
# Usage ./shell.sh <podname>
kubectl exec --stdin --tty -n flux-operator ${1} -- /bin/bash