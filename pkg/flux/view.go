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

		hostTemplate = `hosts = [{host="%s", bind="tcp://eth0:%s", connect="tcp://%s:%s"},
		 {host="%s"}]`

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
curve_cert = "%s/view/etc/curve/curve.cert"
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
func GenerateFluxEntrypoint(cluster *api.MiniCluster) (string, error) {

	// fluxRoot for the view is in /opt/view/lib
	// This must be consistent between the flux-view containers
	// github.com:converged-computing/flux-views.git
	fluxRoot := "/opt/view"

	mainHost := fmt.Sprintf("%s-0", cluster.Name)

	// Generate the curve certificate
	curveCert, err := GetCurveCert(cluster)
	if err != nil {
		return "", err
	}

	// Generate hostlists, this is the lead broker
	hosts := generateHostlist(cluster, cluster.Spec.MaxSize)
	brokerConfig := generateBrokerConfig(cluster, hosts)

	setup := `#!/bin/sh
fluxroot=%s
mainHost=%s
echo "Hello I am hostname $(hostname) running setup."

# Always use verbose, no reason to not here
echo "Flux install root: ${fluxroot}"
export fluxroot

# Add flux to the path
export PATH=/opt/view/bin:$PATH

# Cron directory
mkdir -p $fluxroot/etc/flux/system/cron.d
mkdir -p $fluxroot/var/lib/flux

# These actions need to happen on all hosts
mkdir -p $fluxroot/etc/flux/system
hosts="%s"
echo "flux R encode --hosts=${hosts} --local"
flux R encode --hosts=${hosts} --local > ${fluxroot}/etc/flux/system/R

echo
echo "📦 Resources"
cat ${fluxroot}/etc/flux/system/R

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

# Generate the certificate (ONLY if the lead broker)
if [[ "$(hostname)" == "${mainHost}" ]]; then
echo "Generating curve certificate at main host..."
cat <<EOT >> ${fluxroot}/etc/curve/curve.cert
%s
EOT
echo
echo "🌟️ Curve Certificate"
cat ${fluxroot}/etc/curve/curve.cert
fi

# Now prepare to copy finished spack view over
echo "Moving content from /opt/view to be in shared volume at %s"
view=$(ls /opt/views/._view/)
view="/opt/views/._view/${view}"

# Give a little extra wait time
sleep 10

viewroot="%s"
mkdir -p $viewroot/view
# We have to move both of these paths, *sigh*
cp -R ${view}/* $viewroot/view
cp -R /opt/software $viewroot/

# This is a marker to indicate the copy is done
touch $viewroot/flux-operator-done.txt

# Sleep forever, the application needs to run and end
echo "Sleeping forever so %s can be shared and used for application containers."
sleep infinity
`

	return fmt.Sprintf(
		setup,
		fluxRoot,
		mainHost,
		hosts,
		brokerConfig,
		curveCert,
		cluster.Spec.Flux.Container.MountPath,
		cluster.Spec.Flux.Container.MountPath,
		cluster.Spec.Flux.Container.Name,
	), nil
}
