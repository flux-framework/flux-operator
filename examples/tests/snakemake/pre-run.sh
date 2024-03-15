#!/bin/bash

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "Cloning snakemake run data to /tmp/workflow"
git clone --depth 1 https://github.com/snakemake/snakemake-tutorial-data /tmp/workflow

wget -O /tmp/workflow/Snakefile https://raw.githubusercontent.com/rse-ops/flux-hpc/main/snakemake/atacseq/Snakefile
mkdir -p /tmp/workflow/scripts
wget -O /tmp/workflow/scripts/plot-quals.py https://raw.githubusercontent.com/rse-ops/flux-hpc/main/snakemake/atacseq/scripts/plot-quals.py

echo "Preparing to mount into MiniKube"
minikube ssh -- mkdir -p /data
minikube mount /tmp/workflow:/data &
