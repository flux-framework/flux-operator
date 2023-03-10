#!/bin/sh

# This waiting script is run continuously and only updates
# hosts and runs the job (breaking from the while) when
# update_hosts.sh has been populated. This means the pod usually
# needs to be updated with the config map that has ips!

# If we are not in strict, don't set strict mode
{{ if not .Logging.Strict }}set -eEu -o pipefail{{ end }}

# Set the flux user from the getgo
fluxuser={{ if .Container.FluxUser.Name}}{{ .Container.FluxUser.Name }}{{ else }}flux{{ end }}
fluxuid={{ if .Container.FluxUser.Uid}}{{ .Container.FluxUser.Uid }}{{ else }}1000{{ end }}

{{ if not .Logging.Quiet }}# Show asFlux directive once
printf "\nFlux username: ${fluxuser}\n"{{ end }}

# Ensure pythonpath is set to something
if [ -z ${PYTHONPATH+x} ]; then
  PYTHONPATH=""
fi

# commands to be run as root
asSudo="sudo -E PYTHONPATH=$PYTHONPATH -E PATH=$PATH"

# Always run flux commands (and the broker) as flux user, unless requested otherwise (e.g., for storage)
{{ if .Container.Commands.RunFluxAsRoot }}
# Storage won't have write if we are the flux user
asFlux="${asSudo} -E HOME=/home/${fluxuser}"{{ else }}
# and ensure the home is targeted to be there too.
asFlux="sudo -u ${fluxuser} -E PYTHONPATH=$PYTHONPATH -E PATH=$PATH -E HOME=/home/${fluxuser}"
{{ end }}

# We currently require sudo and an ubuntu base
which sudo > /dev/null 2>&1 || (echo "sudo is required to be installed" && exit 1);
which flux > /dev/null 2>&1 || (echo "flux is required to be installed" && exit 1);

# Add fluxuser to sudoers
echo "${fluxuser} ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

# Add a flux user (required) that should exist before pre-command
sudo adduser --disabled-password --uid ${fluxuid} --gecos "" ${fluxuser} > /dev/null 2>&1 || {{ if not .Logging.Quiet }} printf "${fluxuser} user is already added.\n"{{ else }}true{{ end }}

{{ if .Users }}{{range $username := .Users}}# Add additional users
printf "Adding '{{.Name}}' with password '{{ .Password}}'\n"
sudo useradd -m -p $(openssl passwd '{{ .Password }}') {{.Name}}
{{ end }}{{ end }}

# Show user permissions / ids
{{ if not .Logging.Quiet }}printf "${fluxuser} user identifiers:\n$(id ${fluxuser})\n"{{ end }}

# If any preCommand logic is defined
{{ .Container.PreCommand}}

# And pre command logic that isn't passed to the certificate generator
{{ .Container.Commands.Pre}}

{{ if not .Logging.Quiet }}# Show asFlux directive once
printf "\nAs Flux prefix for flux commands: ${asFlux}\n"{{ end }}


# We use the actual time command and not the wrapper, otherwise we get there is no argument -f
{{ if .Logging.Timed }}which /usr/bin/time > /dev/null 2>&1 || (echo "/usr/bin/time is required to use logging.timed true" && exit 1);{{ end }}

# Broker Options: important!
# The local-uri setting places the unix domain socket in rundir
#   if FLUX_URI is not set, tools know where to connect.
#   -Slog-stderr-level= can be set to 7 for larger debug level
#   or exposed as a variable
brokerOptions="-Scron.directory=/etc/flux/system/cron.d \
  -Stbon.fanout=256 \
  -Srundir=/run/flux {{ if .Interactive }}-Sbroker.rc2_none {{ end }} \
  -Sstatedir=${STATE_DIRECTORY:-/var/lib/flux} \
  -Slocal-uri=local:///run/flux/local \
{{ if not .Logging.Quiet }}  -Slog-stderr-level={{or .Container.FluxLogLevel 6}} {{ else }} -Slog-stderr-level=0 {{ end }} \
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

# Main host <name>-0
mainHost="{{ .MainHost}}"

# The working directory should be set by the CRD or the container
workdir=${PWD}

# And if we are using fusefs / object storage, ensure we can see contents
mkdir -p ${workdir}

{{ if not .Logging.Quiet }}
printf "\n👋 Hello, I'm $(hostname)\n"
printf "The main host is ${mainHost}\n"
printf "The working directory is ${workdir}, contents include:\n"
ls ${workdir}
printf "End of file listing, if you see nothing above there are no files.\n"{{ end }}

# These actions need to happen on all hosts
# Configure resources
mkdir -p /etc/flux/system

# --cores=IDS Assign cores with IDS to each rank in R, so we  assign 0-(N-1) to each host
{{ if not .Logging.Quiet }}echo "flux R encode --hosts={{ .Hosts}} {{if .Cores}}--cores=0-{{.Cores}}{{ end }}"{{ end }}
flux R encode --hosts={{ .Hosts}} {{if .Cores}}--cores=0-{{.Cores}}{{ end }} > /etc/flux/system/R
{{ if not .Logging.Quiet }}printf "\n📦 Resources\n"
cat /etc/flux/system/R{{ end }}

# Do we want to run diagnostics instead of regular entrypoint?
diagnostics="{{ .Container.Diagnostics}}"
{{ if not .Logging.Quiet }}printf "\n🐸 Diagnostics: ${diagnostics}\n"{{ end }}

# Flux option flags
option_flags="{{ .Container.FluxOptionFlags}}"
if [ "${option_flags}" != "" ]; then
    # Make sure we don't get rid of any already defined flags
    existing_flags="${FLUX_OPTION_FLAGS:-}"

    # provide them first so they are replaced by new ones here
    if [ "${existing_flags}" != "" ]; then
        export FLUX_OPTION_FLAGS="${existing_flags} ${option_flags}"
    else
        export FLUX_OPTION_FLAGS="${option_flags}"
    fi
{{ if not .Logging.Quiet }}    echo "🚩️ Flux Option Flags defined"{{ end }}
fi

mkdir -p /etc/flux/imp/conf.d/
cat <<EOT >> /etc/flux/imp/conf.d/imp.toml
[exec]
allowed-users = [ "${fluxuser}", "root" ]
allowed-shells = [ "/usr/libexec/flux/flux-shell" ]
EOT

{{ if not .Logging.Quiet }}printf "\n🦊 Independent Minister of Privilege\n"
cat /etc/flux/imp/conf.d/imp.toml

printf "\n🐸 Broker Configuration\n"
cat /etc/flux/config/broker.toml{{ end }}

# If we are communicating via the flux uri this service needs to be started
chmod u+s /usr/libexec/flux/flux-imp
chmod 4755 /usr/libexec/flux/flux-imp
chmod 0644 /etc/flux/imp/conf.d/imp.toml
sudo service munge start > /dev/null 2>&1

# The rundir needs to be created first, and owned by user flux
# Along with the state directory and curve certificate
mkdir -p /run/flux /etc/curve

# Show generated curve certificate - the munge.key should already be equivalent (and exist)
cp /mnt/curve/curve.cert /etc/curve/curve.cert

# Remove group and other read
chmod o-r /etc/curve/curve.cert
chmod g-r /etc/curve/curve.cert

# We must get the correct flux user id - this user needs to own
# the run directory and these others
fluxuid=$(id -u ${fluxuser})

{{ if not .Container.Commands.RunFluxAsRoot }}chown -R ${fluxuid} /run/flux ${STATE_DIR} /etc/curve/curve.cert ${workdir}{{ end }}

# Make directory world read/writable
chmod -R 0777 ${workdir}

{{ if not .Logging.Quiet }}
printf "\n🔒️ Working directory permissions:\n$(ls -l ${workdir})\n\n"
printf "\n✨ Curve certificate generated by helper pod\n"
cat /etc/curve/curve.cert{{ end }}

function run_flux_restful() {

    # Start restful API server
    branch={{if .FluxRestful.Branch}}{{.FluxRestful.Branch}}{{else}}main{{end}}
    startServer="uvicorn app.main:app --host=0.0.0.0 --port={{or .FluxRestful.Port 5000}}"
    printf "Cloning flux-framework/flux-restful-api branch ${branch}\n"
    git clone -b ${branch} --depth 1 https://github.com/flux-framework/flux-restful-api /flux-restful-api > /dev/null 2>&1
    cd /flux-restful-api
            
    # Export the main flux user and token "superuser"
    export FLUX_USER={{ .FluxUser}}
    export FLUX_TOKEN={{ .FluxToken}}

    # Install python requirements, with preference for python3
    python3 -m pip install -r requirements.txt > /dev/null 2>&1 || python -m pip install -r requirements.txt > /dev/null 2>&1

    # Prepare databases!
    alembic revision --autogenerate -m "Create intital tables"
    alembic upgrade head
    python3 app/db/init_db.py init || python app/db/init_db.py init

    {{ if .Users }}{{range $username := .Users}}# Add additional users
    printf "Adding '{{.Name}}' with password '{{ .Password}}'\n"
    python3 ./app/db/init_db.py add-user "{{.Name}}" "{{.Password}}" || python ./app/db/init_db.py add-user "{{.Name}}" "{{.Password}}"
    {{ end }}{{ end }}

    # Shared envars across user modes
    export FLUX_REQUIRE_AUTH=true
    export FLUX_SECRET_KEY={{ .FluxRestful.SecretKey}}
    export FLUX_NUMBER_NODES={{ .Size}}

    printf "\n 🔑 Use your Flux user and token credentials to authenticate with the MiniCluster with flux-framework/flux-restful-api\n"

    # -o is an "option" for the broker
    # -S corresponds to a shortened --setattr=ATTR=VAL 
    printf "\n🌀 ${asFlux} flux broker --config-path /etc/flux/config ${brokerOptions} ${startServer}\n"
    ${asFlux} flux broker --config-path /etc/flux/config ${brokerOptions} ${startServer}
}

# Run an interactive cluster, giving no command to flux start
function run_interactive_cluster() {
    printf "\n🌀 ${asFlux} flux broker --config-path /etc/flux/config ${brokerOptions}\n"
    ${asFlux} flux broker --config-path /etc/flux/config ${brokerOptions}
}

# Are we running diagnostics or the start command?
if [ "${diagnostics}" == "true" ]; then
    run_diagnostics
else

    # Start flux with the original entrypoint
    if [ $(hostname) == "${mainHost}" ]; then

        # No command - use default to start server
{{ if not .Logging.Quiet }}        echo "Extra command arguments are: $@"{{ end }}
        if [ "$@" == "" ]; then

            # An interactive job also doesn't require a command
            {{ if .Interactive }}run_interactive_cluster
            {{ else }}run_flux_restful{{ end }}

        # Case 2: Fall back to provided command
        else
{{ if not .Logging.Quiet }} # if tasks >= size
            printf "\n🌀 flux start -o --config /etc/flux/config ${brokerOptions} flux submit {{ if ge .Tasks .Size }} -N {{.Size}}{{ end }} -n {{.Tasks}} --quiet {{ if .Container.FluxOptionFlags }}{{ .Container.FluxOptionFlags}}{{ end }} --watch{{ if .Logging.Debug }} -vvv{{ end }} $@\n"{{ end }}
            {{ if .Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxstart wall time %E" {{ end }}${asFlux} flux start -o --config /etc/flux/config ${brokerOptions} {{ if .Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxsubmit wall time %E" {{ end }} flux submit {{ if ge .Tasks .Size }} -N {{.Size}}{{ end }} -n {{.Tasks}} --quiet {{ if .Container.FluxOptionFlags }}{{ .Container.FluxOptionFlags}}{{ end }} --watch{{ if .Logging.Debug }} -vvv{{ end }} $@
        fi
    else
        # Sleep until the broker is ready
{{ if not .Logging.Quiet }}
        printf "\n🌀 flux start -o --config /etc/flux/config ${brokerOptions}\n"{{ end }}
        while true
        do
            {{ if .Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxstart wall time %E" {{ end }}${asFlux} flux start -o --config /etc/flux/config ${brokerOptions}
            {{ if not .Logging.Quiet }}printf "\n😪 Sleeping 15s until broker is ready..."{{ end }}
            sleep 15
        done
    fi
fi
