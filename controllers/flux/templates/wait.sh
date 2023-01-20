#!/bin/sh

# This waiting script is run continuously and only updates
# hosts and runs the job (breaking from the while) when 
# update_hosts.sh has been populated. This means the pod usually
# needs to be updated with the config map that has ips!

# Always run flux commands (and the broker) as flux user
asFlux="sudo -u flux -E"

# If any preCommand logic is defined
{{ .PreCommand}}

{{ if not .TestMode }}# Show asFlux directive once
printf "\nAs Flux prefix for flux commands: ${asFlux}\n"{{ end }}

# We currently require sudo and an ubuntu base
which sudo > /dev/null 2>&1 || (echo "sudo is required to be installed" && exit 1);
which flux > /dev/null 2>&1 || (echo "flux is required to be installed" && exit 1);

# Broker Options: important!
# The local-uri setting places the unix domain socket in rundir 
#   if FLUX_URI is not set, tools know where to connect.
#   -Slog-stderr-level= can be set to 7 for larger debug level
#   or exposed as a variable
brokerOptions="-Scron.directory=/etc/flux/system/cron.d \
  -Stbon.fanout=256 \
  -Srundir=/run/flux \
  -Sstatedir=${STATE_DIRECTORY:-/var/lib/flux} \
  -Slocal-uri=local:///run/flux/local \
{{ if not .TestMode }}  -Slog-stderr-level={{or .FluxLogLevel 6}} {{ else }} -Slog-stderr-level=0 {{ end }} \
  -Slog-stderr-mode=local"

# quorum settings influence how the instance treats missing ranks
#   by default all ranks must be online before work is run, but
#   we want it to be OK to run when a few are down
# These are currently removed because we want the main rank to
# wait for all the others, and then they clean up nicely
#  -Sbroker.quorum=0 \
#  -Sbroker.quorum-timeout=none \

# This should be added to keep running as a service
#  -Sbroker.rc2_none \

# Run diagnostics instead of a command
run_diagnostics() {
    printf "\n🐸 ${asFlux} flux start -o --config /etc/flux/config ${brokerOptions} flux overlay status\n"
    ${asFlux} flux start -o --config /etc/flux/config ${brokerOptions} flux overlay status
    printf "\n🐸 ${asFlux} flux start -o --config /etc/flux/config ${brokerOptions} flux lsattr -v\n"
    ${asFlux} flux start -o --config /etc/flux/config ${brokerOptions} flux lsattr -v
    printf "\n🐸 ${asFlux} flux start -o --config /etc/flux/config ${brokerOptions} flux dmesg\n"
    ${asFlux} flux start -o --config /etc/flux/config ${brokerOptions} flux dmesg
    printf "\n💤 sleep infinity\n"
    sleep infinity
}

# The statedir similarly should exist and have plenty of available space.
# If there are differences in containers / volumes this could eventually be
# exposed as STATEDIR variable
export STATE_DIR=/var/lib/flux
mkdir -p ${STATE_DIR}

# Cron directory
mkdir -p /etc/flux/system/cron.d

# uuid for flux token (auth)
FLUX_TOKEN="{{ .FluxToken}}"

# Main host <name>-0
mainHost="{{ .MainHost}}"

# The working directory should be set by the CRD or the container
workdir=${PWD}

{{ if not .TestMode }}
printf "\n👋 Hello, I'm $(hostname)\n"
printf "The main host is ${mainHost}\n"
printf "The working directory is ${workdir}\n"
ls ${workdir}{{ end }}

# These actions need to happen on all hosts
# Configure resources
mkdir -p /etc/flux/system

# --cores=IDS Assign cores with IDS to each rank in R, so we  assign 0-(N-1) to each host
{{ if not .TestMode }}echo "flux R encode --hosts={{ .Hosts}} {{if .Cores}}--cores=0-{{.Cores}}{{ end }}"{{ end }}
flux R encode --hosts={{ .Hosts}} {{if .Cores}}--cores=0-{{.Cores}}{{ end }} > /etc/flux/system/R
{{ if not .TestMode }}printf "\n📦 Resources\n"
cat /etc/flux/system/R{{ end }}

# Do we want to run diagnostics instead of regular entrypoint?
diagnostics="{{ .Diagnostics}}"
{{ if not .TestMode }}printf "\n🐸 Diagnostics: ${diagnostics}\n"{{ end }}

# Flux option flags
option_flags="{{ .FluxOptionFlags}}"
if [ "${option_flags}" != "" ]; then
    # Make sure we don't get rid of any already defined flags
    existing_flags="${FLUX_OPTION_FLAGS:-}"

    # provide them first so they are replaced by new ones here
    if [ "${existing_flags}" != "" ]; then
        export FLUX_OPTION_FLAGS="${existing_flags} ${option_flags}"
    else 
        export FLUX_OPTION_FLAGS="${option_flags}"
    fi
{{ if not .TestMode }}    echo "🚩️ Flux Option Flags defined"{{ end }}
fi

mkdir -p /etc/flux/imp/conf.d/
cat <<EOT >> /etc/flux/imp/conf.d/imp.toml
[exec]
allowed-users = [ "flux", "root" ]
allowed-shells = [ "/usr/libexec/flux/flux-shell" ]	
EOT

{{ if not .TestMode }}printf "\n🦊 Independent Minister of Privilege\n"
cat /etc/flux/imp/conf.d/imp.toml

printf "\n🐸 Broker Configuration\n"
cat /etc/flux/config/broker.toml{{ end }}

# Add a flux user (required)
sudo adduser --disabled-password --uid 1000 --gecos "" flux > /dev/null 2>&1 || {{ if not .TestMode }} printf "flux user is already added.\n"{{ else }}true{{ end }}

# The rundir needs to be created first, and owned by user flux
# Along with the state directory and curve certificate
mkdir -p /run/flux /etc/curve

# Show generated curve certificate - the munge.key should already be equivalent (and exist)
{{ if not .TestMode }}cat /mnt/curve/curve.cert{{ end }}
cp /mnt/curve/curve.cert /etc/curve/curve.cert

# Remove group and other reead
chmod o-r /etc/curve/curve.cert
chmod g-r /etc/curve/curve.cert

# We must get the correct flux user id - this user needs to own
# the run directory and these others
fluxuid=$(id -u flux)
chown -R ${fluxuid} /run/flux ${STATE_DIR} /etc/curve/curve.cert ${workdir}

{{ if not .TestMode }}
printf "\n✨ Curve certificate generated by helper pod\n"
cat /etc/curve/curve.cert{{ end }}

# Are we running diagnostics or the start command?
if [ "${diagnostics}" == "true" ]; then
    run_diagnostics
else

    # Start flux with the original entrypoint
    if [ $(hostname) == "${mainHost}" ]; then

        # No command - use default to start server
{{ if not .TestMode }}        echo "Extra arguments are: $@"{{ end }}
        if [ "$@" == "" ]; then

            # Start restful API server
            startServer="uvicorn app.main:app --host=0.0.0.0 --port={{or .FluxRestfulPort 5000}} {{if .Size }}--workers {{.Size}}{{ end }}"
            git clone -b {{or .FluxRestfulBranch "main"}} --depth 1 https://github.com/flux-framework/flux-restful-api /flux-restful-api > /dev/null 2>&1
            cd /flux-restful-api

            # Install python requirements, with preference for python3
            python3 -m pip install -r requirements.txt > /dev/null 2>&1 || python -m pip install -r requirements.txt > /dev/null 2>&1

            # Generate a random flux token
            FLUX_USER=flux 
            FLUX_REQUIRE_AUTH=true
            FLUX_NUMBER_NODES={{ .ClusterSize}}
            export FLUX_TOKEN FLUX_USER FLUX_REQUIRE_AUTH FLUX_NUMBER_NODES

{{ if not .TestMode }}
            printf "\n 🔑 Your Credentials! These will allow you to control your MiniCluster with flux-framework/flux-restful-api\n"
            printf "export FLUX_TOKEN=${FLUX_TOKEN}\n"
            printf "export FLUX_USER=${FLUX_USER}\n"

            # -o is an "option" for the broker
            # -S corresponds to a shortened --setattr=ATTR=VAL
            printf "\n🌀 flux start -o --config /etc/flux/config ${brokerOptions} ${startServer}\n"{{ end }}
            ${asFlux} flux start -o --config /etc/flux/config ${brokerOptions} ${startServer}

        # Case 2: Fall back to provided command
        else
{{ if not .TestMode }} 
            printf "\n🌀 flux start -o --config /etc/flux/config ${brokerOptions} flux mini submit {{ if gt .Tasks .Size }} -N {{.Size}}{{ end }} -n {{.Tasks}} {{ if .FluxOptionFlags }}{{ .FluxOptionFlags}}{{ end }} --watch $@\n"{{ end }}
            ${asFlux} flux start -o --config /etc/flux/config ${brokerOptions} flux mini submit {{ if gt .Tasks .Size }} -N {{.Size}}{{ end }} -n {{.Tasks}} {{ if .FluxOptionFlags }}{{ .FluxOptionFlags}}{{ end }} --watch $@
        fi
    else
        # Sleep until the broker is ready
{{ if not .TestMode }}
        printf "\n🌀 flux start -o --config /etc/flux/config ${brokerOptions}\n"{{ end }}
        while true
        do
            ${asFlux} flux start -o --config /etc/flux/config ${brokerOptions}
            {{ if not .TestMode }}printf "\n😪 Sleeping 15s until broker is ready..."{{ end }}
            sleep 15
        done
    fi
fi
