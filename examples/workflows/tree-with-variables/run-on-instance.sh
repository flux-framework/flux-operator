#!/bin/bash

# It's hard to see this for a quick job, so let's write to a file!
outdir=${1}
outfile="${outdir}/${FLUX_TREE_ID}-output.txt"

# ID string uniquely identifying the hierarchical path of the Flux instance on which Jobscript is being executed
echo "FLUX_TREE_ID ${FLUX_TREE_ID}" > "${outfile}"

# the integer ID of each jobscript invocation local to the Flux instance. It starts from 1 and sequentially increases.
echo "FLUX_TREE_JOBSCRIPT_INDEX ${FLUX_TREE_JOBSCRIPT_INDEX}" >> "${outfile}"

# the number nodes assigned to the instance
echo "FLUX_TREE_NNODES ${FLUX_TREE_NNODES}" >> "${outfile}"

# the number of cores per node assigned to the instance
echo "FLUX_TREE_NCORES_PER_NODE ${FLUX_TREE_NCORES_PER_NODE}" >> "${outfile}"

# the number of GPUs per node assigned to the instance.
echo "FLUX_TREE_NGPUS_PER_NODE ${FLUX_TREE_NGPUS_PER_NODE}" >> "${outfile}"