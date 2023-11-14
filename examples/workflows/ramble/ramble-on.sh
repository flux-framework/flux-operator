#!/bin/bash

# We use a bash script wrapper since we need to source environments in
# the job. If you use spack with ramble, you'd need to source it too.

. /opt/ramble/share/ramble/setup-env.sh
ramble workspace activate /opt/test_workspace
ramble workspace concretize || true
ramble workspace setup
ramble on