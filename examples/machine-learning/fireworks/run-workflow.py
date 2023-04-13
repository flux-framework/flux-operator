#!/usr/bin/env python3
#
#-  run-workflow.py ~~
#
#-  Do not run this script manually for this demo. This script can run on a
#   login node, but the demo will fail on the second step of the workflow
#   because of MPI. Instead, this demo runs this script for you as part of
#   LSF batch script:
#       $ bsub batch-runner.lsf
#

import os
from fireworks import LaunchPad
from fireworks.core.rocket_launcher import rapidfire

# Set up and reset the LaunchPad using MongoDB URI string.
launchpad = LaunchPad(host = os.getenv("MONGODB_URI"), uri_mode = True)

# Launch workflow "locally" -- this will run on a Summit batch node.
rapidfire(launchpad)
print("Done.")

#-  vim:set syntax=python: