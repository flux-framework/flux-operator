#!/bin/sh

# This script exists just to generate and print the curve.cert for adding
# as a config map to the actual Flux MiniCluster

# Always run flux commands (and the broker) as flux user
asFlux="sudo -u flux -E"

# If any preCommand logic is defined
{{ .PreCommand}}

# We currently require sudo and an ubuntu base
which sudo > /dev/null 2>&1 || (echo "sudo is required to be installed" && exit 1);
which flux > /dev/null 2>&1 || (echo "flux is required to be installed" && exit 1);

# Add a flux user (required)
sudo adduser --disabled-password --uid 1000 --gecos "" flux > /dev/null 2>&1 || true

# The entire purpose of this script is to generate the certificate and print to logs
mkdir -p /mnt/curve
flux keygen /mnt/curve/curve.cert
cat /mnt/curve/curve.cert
