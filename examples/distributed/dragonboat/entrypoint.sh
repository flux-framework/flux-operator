#!/bin/bash

# Can't be zero
replica=$((FLUX_TASK_RANK+1))
/opt/dragon/example-helloworld -replicaid $replica -addr flux-sample-$FLUX_TASK_RANK.flux-service.default.svc.cluster.local:6300${replica}