#!/bin/bash

job_name="${1:-flux-sample}"
job_port="${2:-8080}"
nodes="${3:-2}"
echo "I am hostname $(hostname) and rank ${FLUX_TASK_RANK} of ${nodes} nodes. The job is ${job_name} and master is on port ${job_port}"

# This will be parsed by the main.py to get the rank
export LOCAL_RANK=${FLUX_TASK_RANK}

# Not ideal, but it kind of works
torchrun --node_rank ${LOCAL_RANK} --nnodes ${nodes} --nproc_per_node 2 --master_addr ${job_name}-0.flux-service.default.svc.cluster.local --master_port ${job_port} ./main.py