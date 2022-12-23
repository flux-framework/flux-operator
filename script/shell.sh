#!/bin/bash
# Usage ./shell.sh <podname>
kubectl exec --stdin --tty -n flux-operator ${@} -- /bin/bash