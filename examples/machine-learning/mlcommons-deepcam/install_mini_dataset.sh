#!/bin/bash

# This script will take the downloaded small batch of data files,
# make the required number of duplicates, and install in the specified
# directory in train/ and validation/ subfolders.

if [ $# -lt 2 ]; then
    echo "Usage:"
    echo "  $0 DOWNLOADED_DATA_DIR INSTALLATION_TARGET_DIR [NUM_COPIES]"
    exit 1
fi

sourceDir=$1
targetDir=$2
numCopies=1
if [ $# -ge 3 ]; then
    numCopies=$3
fi

# First, we prepare the train directory by duplicating every file numCopies times
mkdir -p $targetDir/train
for f in $(ls $sourceDir | grep "data-.*.h5"); do
    echo $f
    for (( i=0; i<$numCopies; i++ )); do
        outFile=$targetDir/train/${f/.h5/-$i.h5}
        echo "  $outFile"
        cp $sourceDir/$f $outFile
    done
done

# Copy in the stats file
cp $sourceDir/stats.h5 $targetDir/

# Now copy the training directory to the validation directory
cp -r $targetDir/train $targetDir/validation
