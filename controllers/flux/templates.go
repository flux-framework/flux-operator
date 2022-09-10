/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

var (
	brokerConfigTemplate = `
[access]
allow-guest-user = true
allow-root-owner = true

[bootstrap]
curve_cert = "/mnt/curve/curve.cert"
default_port = 8050
default_bind = "tcp://eth0:%%p"
default_connect = "tcp://%%h:%%p"
hosts = [
	{ host="%s-%s"},
]
`
	waitToStartTemplate = `#!/bin/sh

# This waiting script is run continuously and only updates
# hosts and runs the job (breaking from the while) when 
# update_hosts.sh has been populated. This means the pod usually
# needs to be updated with the config map that has ips!

# We determine the update_hosts.sh is ready when it has content
count_lines() {
	lines=$(cat /flux_operator/update_hosts.sh | wc -l)
	echo $lines
}

while [ $(count_lines) -lt 2 ];
do
    echo "Host updating script not available yet, waiting..."
    sleep 5s
done             

# Run to discover hosts
/bin/sh /flux_operator/update_hosts.sh

# Show host updates
cat /etc/hosts

# Start flux with the original entrypoint
printf "/bin/sh /flux_operator/start.sh $@\n"
/bin/sh /flux_operator/start.sh $@
`

	startFluxTemplate = `#/bin/sh
printf "flux start -o --config /etc/flux/config $@\n"
flux start -o --config /etc/flux/config $@
`
)
