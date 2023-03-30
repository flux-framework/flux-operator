#!/bin/bash

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "Cloning snakemake run data to /tmp/workflow"

git clone --depth 1 https://github.com/snakemake/snakemake-tutorial-data /tmp/data

wget -O /tmp/data/Snakefile https://raw.githubusercontent.com/rse-ops/flux-hpc/main/snakemake/atacseq/Snakefile
mkdir -p /tmp/data/scripts
wget -O /tmp/data/scripts/plot-quals.py https://raw.githubusercontent.com/rse-ops/flux-hpc/main/snakemake/atacseq/scripts/plot-quals.py