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

		hostTemplate = `hosts = [{host="%s", bind="tcp://eth0:%s", connect="tcp://%s:%d"},
		 {host="%d"}]`

		hostBlock = fmt.Sprintf(
			hostTemplate,
			cluster.Spec.Flux.Bursting.LeadBroker.Address,
			cluster.Spec.Flux.Bursting.LeadBroker.Port,
			cluster.Spec.Flux.Bursting.LeadBroker.Address,
			cluster.Spec.Flux.Bursting.LeadBroker.Port,
		)
	}
	return hostBlock
}

func generateBrokerConfig(
	cluster *api.MiniCluster,
	hosts string,
	containerIndex int,
) string {

	if cluster.Spec.Flux.BrokerConfig != "" {
		return cluster.Spec.Flux.BrokerConfig
	}

	// Port assembled based on index. Right now this only supports up
	defaultPort := 8050 + containerIndex
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
path = "%s/view/etc/flux/system/R-%d"

[bootstrap]
curve_cert = "%s/view/curve/curve.cert"
default_port = %d
default_bind = "%s"
default_connect = "%s"
%s

[archive]
dbpath = "%s/view/var/lib/flux/job-archive-%d.sqlite"
period = "1m"
busytimeout = "50s"

[sched-fluxion-qmanager]
queue-policy = "%s"
`
	return fmt.Sprintf(
		template,
		cluster.Spec.Flux.Container.MountPath,
		containerIndex,
		cluster.Spec.Flux.Container.MountPath,
		defaultPort,
		defaultBind,
		defaultConnect,
		hostBlock,
		cluster.Spec.Flux.Container.MountPath,
		containerIndex,
		cluster.Spec.Flux.Scheduler.QueuePolicy,
	)

}

// generateFluxEntrypoint generates the flux entrypoint to prepare flux
// This is run inside of the flux container that will be copied to the empty volume
// If the flux container is disabled, we still add an init container with
// the broker config, etc., but we don't expect a flux view there.
func GenerateFluxEntrypoint(
	cluster *api.MiniCluster,
) (string, error) {

	// fluxRoot for the view is in /opt/view/lib
	// This must be consistent between the flux-view containers
	// github.com:converged-computing/flux-views.git
	fluxRoot := "/opt/view"

	// If we are disabling the view, it won't have flux (or extra spack copies)
	// We copy our faux flux config directory (not a symlink) to the mount path
	spackView := fmt.Sprintf(`mkdir -p $viewroot/software
  cp -R /opt/view/* %s/view`,
		cluster.Spec.Flux.Container.MountPath,
	)

	generateHosts := `echo '📦 Flux view disabled, not generating resources here.'
  mkdir -p ${fluxroot}/etc/flux/system
  `

	// Create a different broker.toml for each runFlux container
	if !cluster.Spec.Flux.Container.Disable {

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

	// Generate a broker config for each potential running flux container
	brokerConfigs := ""
	for i, container := range cluster.Spec.Containers {
		if !container.RunFlux {
			continue
		}

		// Generate hostlists, this is the lead broker
		hosts := generateHostlist(cluster, container, cluster.Spec.MaxSize)

		// Create a different broker.toml for each runFlux container
		if !cluster.Spec.Flux.Container.Disable {
			generateHosts = fmt.Sprintf(`
echo "flux R encode --hosts=${hosts} --local"
flux R encode --hosts=${hosts} --local > ${fluxroot}/etc/flux/system/R-%d
  
echo
echo "📦 Resources"
cat ${fluxroot}/etc/flux/system/R-%d`, i, i)
		}

		brokerConfig := generateBrokerConfig(cluster, hosts, i)
		brokerConfigs += fmt.Sprintf(`
# Write the broker configuration
mkdir -p ${fluxroot}/etc/flux/config-%d

cat <<EOT >> ${fluxroot}/etc/flux/config-%d/broker.toml
%s
EOT

# These actions need to happen on all hosts
mkdir -p $fluxroot/etc/flux/system
hosts="%s"

# Echo hosts here in case the main container needs to generate
echo "${hosts}" > ${fluxroot}/etc/flux/system/hostlist-%d
%s

# Cron directory
mkdir -p $fluxroot/etc/flux/system/cron-%d.d
mkdir -p $fluxroot/var/lib/flux

# The rundir needs to be created first, and owned by user flux
# Along with the state directory and curve certificate
mkdir -p ${fluxroot}/run/flux ${fluxroot}/etc/curve

`, i, i, brokerConfig, hosts, i, generateHosts, i)
	}

	setup := `#!/bin/sh
fluxroot=%s
echo "Hello I am hostname $(hostname) running setup."

# Always use verbose, no reason to not here
echo "Flux install root: ${fluxroot}"
export fluxroot

# Add flux to the path (if using view)
export PATH=/opt/view/bin:$PATH

# If the view doesn't exist, ensure basic paths do
mkdir -p $fluxroot/bin

%s

echo
echo "🐸 Broker Configuration"
for filename in $(find ${fluxroot}/etc/flux -name broker.toml)
  do
  cat $filename
done

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
		brokerConfigs,
		cluster.Spec.Flux.Container.MountPath,
		spackView,
	), nil
}
