#!/bin/bash
# Usage ./log.sh <podname>
kubectl logs -n flux-operator ${@}