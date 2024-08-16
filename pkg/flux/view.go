/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package flux

import (
	"fmt"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"
)

// generateHostBlock generates the host block for the flux config
func generateHostBlock(hosts string, cluster *api.MiniCluster) string {

	// Default hostBlock is simply the provided names
	hostTemplate := `hosts = [
{ host="%s"},
]		
`
	hostBlock := fmt.Sprintf(hostTemplate, hosts)

	// Unless we have a bursting broker address
	if cluster.Spec.Flux.Bursting.LeadBroker.Address != "" {

		hostTemplate = `hosts = [{host="%s", bind="tcp://eth0:%d", connect="tcp://%s:%d"},
		 {host="%s"}]`

		hostBlock = fmt.Sprintf(
			hostTemplate,
			cluster.Spec.Flux.Bursting.LeadBroker.Address,
			cluster.Spec.Flux.Bursting.LeadBroker.Port,
			cluster.Spec.Flux.Bursting.LeadBroker.Address,
			cluster.Spec.Flux.Bursting.LeadBroker.Port,
			hosts,
		)
	}
	return hostBlock
}

func generateBrokerConfig(cluster *api.MiniCluster, hosts string) string {

	if cluster.Spec.Flux.BrokerConfig != "" {
		return cluster.Spec.Flux.BrokerConfig
	}

	hostBlock := generateHostBlock(hosts, cluster)
	fqdn := fmt.Sprintf("%s.%s.svc.cluster.local", cluster.Spec.Network.HeadlessName, cluster.Namespace)

	// These shouldn't be formatted in block
	defaultBind := "tcp://eth0:%p"
	defaultConnect := "tcp://%h" + fmt.Sprintf(".%s:", fqdn) + "%p"

	template := `[access]
allow-guest-user = true
allow-root-owner = true

# Point to resource definition generated with flux-R(1).
[resource]
path = "%s/view/etc/flux/system/R"

[bootstrap]
curve_cert = "%s/view/curve/curve.cert"
default_port = 8050
default_bind = "%s"
default_connect = "%s"
%s

[archive]
dbpath = "%s/view/var/lib/flux/job-archive.sqlite"
period = "1m"
busytimeout = "50s"

[sched-fluxion-qmanager]
queue-policy = "%s"
`
	return fmt.Sprintf(
		template,
		cluster.Spec.Flux.Container.MountPath,
		cluster.Spec.Flux.Container.MountPath,
		defaultBind,
		defaultConnect,
		hostBlock,
		cluster.Spec.Flux.Container.MountPath,
		cluster.Spec.Flux.Scheduler.QueuePolicy,
	)

}

// generateFluxEntrypoint generates the flux entrypoint to prepare flux
// This is run inside of the flux container that will be copied to the empty volume
// If the flux container is disabled, we still add an init container with
// the broker config, etc., but we don't expect a flux view there.
func GenerateFluxEntrypoint(cluster *api.MiniCluster) (string, error) {

	// fluxRoot for the view is in /opt/view/lib
	// This must be consistent between the flux-view containers
	// github.com:converged-computing/flux-views.git
	fluxRoot := "/opt/view"

	// Get the main host, either cluster name or custom
	mainHost := cluster.MainHost()

	// Generate hostlists, this is the lead broker
	hosts := generateHostlist(cluster, cluster.Spec.MaxSize)
	brokerConfig := generateBrokerConfig(cluster, hosts)

	// If we are disabling the view, it won't have flux (or extra spack copies)
	// We copy our faux flux config directory (not a symlink) to the mount path
	spackView := fmt.Sprintf(`mkdir -p $viewroot/software
cp -R /opt/view/* %s/view`,
		cluster.Spec.Flux.Container.MountPath,
	)

	generateHosts := `echo '📦 Flux view disabled, not generating resources here.'
mkdir -p ${fluxroot}/etc/flux/system
`
	if !cluster.Spec.Flux.Container.Disable {
		generateHosts = `
echo "flux R encode --hosts=${hosts} --local"
flux R encode --hosts=${hosts} --local > ${fluxroot}/etc/flux/system/R

echo
echo "📦 Resources"
cat ${fluxroot}/etc/flux/system/R`

		spackView = `# Now prepare to copy finished spack view over
echo "Moving content from /opt/view to be in shared volume at %s"
# Note that /opt/view is a symlink to here!
view=$(ls /opt/views/._view/)
view="/opt/views/._view/${view}"

# Give a little extra wait time
sleep 10

# We have to move both of these paths, *sigh*
cp -R ${view}/* $viewroot/view
cp -R /opt/software $viewroot/
`
	}

	setup := `#!/bin/sh
fluxroot=%s
mainHost=%s
echo "Hello I am hostname $(hostname) running setup."

# Always use verbose, no reason to not here
echo "Flux install root: ${fluxroot}"
export fluxroot

# Add flux to the path (if using view)
export PATH=/opt/view/bin:$PATH

# If the view doesn't exist, ensure basic paths do
mkdir -p $fluxroot/bin

# Cron directory
mkdir -p $fluxroot/etc/flux/system/cron.d
mkdir -p $fluxroot/var/lib/flux

# These actions need to happen on all hosts
mkdir -p $fluxroot/etc/flux/system
hosts="%s"

# Echo hosts here in case the main container needs to generate
echo "${hosts}" > ${fluxroot}/etc/flux/system/hostlist
%s

# Write the broker configuration
mkdir -p ${fluxroot}/etc/flux/config
cat <<EOT >> ${fluxroot}/etc/flux/config/broker.toml
%s
EOT

echo
echo "🐸 Broker Configuration"
cat ${fluxroot}/etc/flux/config/broker.toml

# The rundir needs to be created first, and owned by user flux
# Along with the state directory and curve certificate
mkdir -p ${fluxroot}/run/flux ${fluxroot}/etc/curve

# View the curve certificate
echo "🌟️ Curve Certificate"
cat /flux_operator/curve.cert

viewroot="%s"
mkdir -p $viewroot/view

%s

# This is a marker to indicate the copy is done
touch $viewroot/flux-operator-done.txt
echo "Application is done."
`

	return fmt.Sprintf(
		setup,
		fluxRoot,
		mainHost,
		hosts,
		generateHosts,
		brokerConfig,
		cluster.Spec.Flux.Container.MountPath,
		spackView,
	), nil
}
