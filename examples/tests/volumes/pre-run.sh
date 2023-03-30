#!/bin/bash

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
TESTS=$(dirname ${HERE})

echo "Copying local volume to /tmp/data mounted in kind"

# The "data" volume will be mounted at /mnt/data
mkdir -p /tmp/data
cp ${TESTS}/data/pancakes.txt /tmp/data/pancakes.txt
ls /tmp/data
