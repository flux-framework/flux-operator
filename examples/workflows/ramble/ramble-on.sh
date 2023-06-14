#!/bin/bash

# We use a bash script wrapper since we need to source environments in
# the job. If you use spack with ramble, you'd need to source it too.

. /home/flux/ramble/share/ramble/setup-env.sh
ramble workspace activate /home/flux/test_workspace
ramble workspace concretize
ramble workspace setup
ramble on