#!/bin/bash

MERLIN_INFO=$1
cd /workflow/cavity

echo "***** Setting Up Mesh *****"
python $MERLIN_INFO/scripts/mesh_param_script.py -scripts_dir $MERLIN_INFO/scripts/
mv blockMeshDict.txt system/blockMeshDict
if [ -e system/blockMeshDict ]; then
  echo "... blockMeshDict.txt complete"
fi

if [ -e system/controlDict ]; then
  CONTROL="system/controlDict"
  echo "***** Setting Control Dictionary *****"
  sed -i "30s/.*/writeControl    runTime;/" ${CONTROL}
  sed -i "26s/.*/endTime         1;/"  ${CONTROL}
  sed -i "32s/.*/writeInterval   .1;/" ${CONTROL}
  echo "... system/controlDict edits complete"
else
  echo "Can't find system/controlDict"
fi
