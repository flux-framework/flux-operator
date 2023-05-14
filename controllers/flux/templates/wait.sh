#!/bin/sh

# This waiting script is run continuously and only updates
# hosts and runs the job (breaking from the while) when
# update_hosts.sh has been populated. This means the pod usually
# needs to be updated with the config map that has ips!

# If any initCommand logic is defined
{{ .Container.Commands.Init}} {{ if .Spec.Logging.Quiet }}> /dev/null{{ end }}

# If we are not in strict, don't set strict mode
{{ if .Spec.Logging.Strict }}set -eEu -o pipefail{{ end }}

# Set the flux user from the getgo
fluxuser={{ if .Container.FluxUser.Name}}{{ .Container.FluxUser.Name }}{{ else }}flux{{ end }}
fluxuid={{ if .Container.FluxUser.Uid}}{{ .Container.FluxUser.Uid }}{{ else }}1000{{ end }}
fluxroot={{ if .Spec.Flux.InstallRoot }}{{ .Spec.Flux.InstallRoot }}{{ else }}/usr{{ end }}
if [ "${fluxroot}" == "" ]; then
  fluxroot="/usr"
fi

{{ if not .Spec.Logging.Quiet }}# Show asFlux directive once
printf "\nFlux username: ${fluxuser}\n"{{ end }}

# Ensure pythonpath is set to something
if [ -z ${PYTHONPATH+x} ]; then
  PYTHONPATH=""
fi
if [ -z ${LD_LIBRARY_PATH+x} ]; then
  LD_LIBRARY_PATH=""
fi

# Set the flux root
{{ if not .Spec.Logging.Quiet }}printf "\nFlux install root: ${fluxroot}\n"{{ end }}
export fluxroot

# commands to be run as root
asSudo="sudo -E PYTHONPATH=$PYTHONPATH -E PATH=$PATH -E LD_LIBRARY_PATH=${LD_LIBRARY_PATH}"

# Always run flux commands (and the broker) as flux user, unless requested otherwise (e.g., for storage)
{{ if .Container.Commands.RunFluxAsRoot }}
# Storage won't have write if we are the flux user
asFlux="${asSudo} -E HOME=/home/${fluxuser}"{{ else }}
# and ensure the home is targeted to be there too.
asFlux="sudo -u ${fluxuser} -E PYTHONPATH=$PYTHONPATH -E PATH=$PATH -E LD_LIBRARY_PATH=${LD_LIBRARY_PATH} -E HOME=/home/${fluxuser}"
{{ end }}

# If any preCommand logic is defined
{{ .Container.PreCommand}} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

# And pre command logic that isn't passed to the certificate generator
{{ .Container.Commands.Pre}} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

# We currently require sudo and an ubuntu base
which sudo > /dev/null 2>&1 || (echo "sudo is required to be installed" && exit 1);
which flux > /dev/null 2>&1 || (echo "flux is required to be installed" && exit 1);

# Add fluxuser to sudoers, only if not running as root
{{ if not .Container.Commands.RunFluxAsRoot }}
echo "${fluxuser} ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

# Add a flux user (required) that should exist before pre-command
sudo adduser --disabled-password --uid ${fluxuid} --gecos "" ${fluxuser} > /dev/null 2>&1 || {{ if not .Spec.Logging.Quiet }} printf "${fluxuser} user is already added.\n"{{ else }}true{{ end }}

# Show user permissions / ids
{{ if not .Spec.Logging.Quiet }}printf "${fluxuser} user identifiers:\n$(id ${fluxuser})\n"{{ end }}
{{ end }}

{{ if .Spec.Users }}{{range $username := .Spec.Users}}# Add additional users
printf "Adding '{{.Name}}' with password '{{ .Password}}'\n"
sudo useradd -m -p $(openssl passwd '{{ .Password }}') {{.Name}}
{{ end }}{{ end }}

{{ if not .Spec.Logging.Quiet }}# Show asFlux directive once
printf "\nAs Flux prefix for flux commands: ${asFlux}\n"{{ end }}

# We use the actual time command and not the wrapper, otherwise we get there is no argument -f
{{ if .Spec.Logging.Timed }}which /usr/bin/time > /dev/null 2>&1 || (echo "/usr/bin/time is required to use logging.timed true" && exit 1);{{ end }}

# If the user wants to save/load archives, we set the state directory to that
# The statedir similarly should exist and have plenty of available space.
export STATE_DIR=/var/lib/flux
export FLUX_OUTPUT_DIR={{ if .Container.Logs }}{{.Container.Logs}}{{ else }}/tmp/fluxout{{ end }}
mkdir -p ${STATE_DIR} ${FLUX_OUTPUT_DIR}

# Broker Options: important!
# The local-uri setting places the unix domain socket in rundir
#   if FLUX_URI is not set, tools know where to connect.
brokerOptions="-Scron.directory=/etc/flux/system/cron.d \
  -Stbon.fanout=256 \
  -Srundir=/run/flux {{ if .Spec.Interactive }}-Sbroker.rc2_none {{ end }} \
  -Sstatedir=${STATE_DIR} \
  -Slocal-uri=local:///run/flux/local \
{{ if .Spec.Flux.ConnectTimeout }}-Stbon.connect_timeout={{ .Spec.Flux.ConnectTimeout }}{{ end }} \
{{ if .RequiredRanks }}-Sbroker.quorum={{ .RequiredRanks }}{{ end }} \
{{ if .Spec.Logging.Zeromq }}-Stbon.zmqdebug=1{{ end }} \
{{ if not .Spec.Logging.Quiet }} -Slog-stderr-level={{or .Spec.Flux.LogLevel 6}} {{ else }} -Slog-stderr-level=0 {{ end }} \
  -Slog-stderr-mode=local"

# if we are given an archive to use, load first, not required to exist
# Note that we ask the user to dump in interactive mode - I am not
# sure that doing it with a hook ensures the dump will be successful.
{{if .Spec.Archive.Path }}
if [[ -e "{{ .Spec.Archive.Path}}" ]]; then
{{ if not .Spec.Logging.Quiet }}printf "🧊️ Found existing archive at {{ .Spec.Archive.Path}} loading into state directory\nBefore:\n"{{ end }}
brokerOptions="${brokerOptions} -Scontent.restore={{ .Spec.Archive.Path}}"
fi{{ end }}

# quorum settings influence how the instance treats missing ranks
#   by default all ranks must be online before work is run, but
#   we want it to be OK to run when a few are down
# These are currently removed because we want the main rank to
# wait for all the others, and then they clean up nicely
#  -Sbroker.quorum=0 \
#  -Sbroker.quorum-timeout=none \

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

# Cron directory
mkdir -p /etc/flux/system/cron.d

# Main host <name>-0 and the fully qualified domain name
mainHost="{{ .MainHost}}"

# The working directory should be set by the CRD or the container
workdir=${PWD}

# And if we are using fusefs / object storage, ensure we can see contents
mkdir -p ${workdir}

# We always do a listing in case it's needed to "expose" object storage
# Often this isn't enough and a list of the paths needed should be
# added to containers[].commands.pre
{{ if not .Spec.Logging.Quiet }}
printf "\n👋 Hello, I'm $(hostname)\n"
printf "The main host is ${mainHost}\n"
printf "The working directory is ${workdir}, contents include:\n"
ls -R ${workdir}
printf "End of file listing, if you see nothing above there are no files.\n"{{ else }}
ls -R ${workdir} > /dev/null 2>&1
{{ end }}

# These actions need to happen on all hosts
# Configure resources
mkdir -p /etc/flux/system

# --cores=IDS Assign cores with IDS to each rank in R, so we  assign 0-(N-1) to each host
{{ if not .Spec.Logging.Quiet }}echo "flux R encode --hosts={{ .Hosts}} {{if .Cores}}--cores=0-{{.Cores}}{{ end }}"{{ end }}
flux R encode --hosts={{ .Hosts}} {{if .Cores}}--cores=0-{{.Cores}}{{ else }}--local{{ end }} > /etc/flux/system/R
{{ if not .Spec.Logging.Quiet }}printf "\n📦 Resources\n"
cat /etc/flux/system/R{{ end }}

# Do we want to run diagnostics instead of regular entrypoint?
diagnostics="{{ .Container.Diagnostics}}"
{{ if not .Spec.Logging.Quiet }}printf "\n🐸 Diagnostics: ${diagnostics}\n"{{ end }}

# Flux option flags
option_flags="{{ .Spec.Flux.OptionFlags}}"
if [ "${option_flags}" != "" ]; then
    # Make sure we don't get rid of any already defined flags
    existing_flags="${FLUX_OPTION_FLAGS:-}"

    # provide them first so they are replaced by new ones here
    if [ "${existing_flags}" != "" ]; then
        export FLUX_OPTION_FLAGS="${existing_flags} ${option_flags}"
    else
        export FLUX_OPTION_FLAGS="${option_flags}"
    fi
{{ if not .Spec.Logging.Quiet }}    echo "🚩️ Flux Option Flags defined"{{ end }}
fi

mkdir -p /etc/flux/imp/conf.d/

{{ if .Container.Commands.RunFluxAsRoot }}
cat <<EOT >> /etc/flux/imp/conf.d/imp.toml
[exec]
allowed-users = [ "root" ]
allowed-shells = [ "${fluxroot}/libexec/flux/flux-shell" ]
EOT
{{ else }}
cat <<EOT >> /etc/flux/imp/conf.d/imp.toml
[exec]
allowed-users = [ "${fluxuser}", "root" ]
allowed-shells = [ "${fluxroot}/libexec/flux/flux-shell" ]
EOT
{{ end }}

{{ if not .Spec.Logging.Quiet }}printf "\n🦊 Independent Minister of Privilege\n"
cat /etc/flux/imp/conf.d/imp.toml

printf "\n🐸 Broker Configuration\n"
cat /etc/flux/config/broker.toml{{ end }}

# If we are communicating via the flux uri this service needs to be started
chmod u+s ${fluxroot}/libexec/flux/flux-imp
chmod 4755 ${fluxroot}/libexec/flux/flux-imp
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

# Either the flux user owns the instance, or root
{{ if not .Container.Commands.RunFluxAsRoot }}

# We must get the correct flux user id - this user needs to own
# the run directory and these others
fluxuid=$(id -u ${fluxuser})

chown -R ${fluxuid} /run/flux ${STATE_DIR} /etc/curve/curve.cert ${workdir} ${FLUX_OUTPUT_DIR}{{ end }}

# Make directory world read/writable
chmod -R 0777 ${workdir}

{{ if not .Spec.Logging.Quiet }}
printf "\n🧊️ State Directory:\n$(ls -l ${STATE_DIR})\n\n"
printf "\n🔒️ Working directory permissions:\n$(ls -l ${workdir})\n\n"
printf "\n✨ Curve certificate generated by helper pod\n"
cat /etc/curve/curve.cert{{ end }}

function run_flux_restful() {

    # Start restful API server
    branch={{if .Spec.FluxRestful.Branch}}{{.Spec.FluxRestful.Branch}}{{else}}main{{end}}
    startServer="uvicorn app.main:app --host=0.0.0.0 --port={{or .Spec.FluxRestful.Port 5000}}"
    printf "Cloning flux-framework/flux-restful-api branch ${branch}\n"
    git clone -b ${branch} --depth 1 https://github.com/flux-framework/flux-restful-api /flux-restful-api > /dev/null 2>&1
    cd /flux-restful-api
            
    # Export the main flux user and token "superuser"
    export FLUX_USER={{ .FluxUser}}
    export FLUX_TOKEN={{ .FluxToken}}
    printf "🔒️ Credentials, my friend!\n    FLUX_USER: ${FLUX_USER}\n    FLUX_TOKEN: ${FLUX_TOKEN}\n\n"

    # Install python requirements, with preference for python3
    python3 -m pip install -r requirements.txt > /dev/null 2>&1 || python -m pip install -r requirements.txt > /dev/null 2>&1

    # Prepare databases!
    alembic revision --autogenerate -m "Create intital tables"
    alembic upgrade head
    python3 app/db/init_db.py init || python app/db/init_db.py init

    {{ if .Spec.Users }}{{range $username := .Spec.Users}}# Add additional users
    printf "Adding '{{.Name}}' with password '{{ .Password}}'\n"
    python3 ./app/db/init_db.py add-user "{{.Name}}" "{{.Password}}" || python ./app/db/init_db.py add-user "{{.Name}}" "{{.Password}}"
    {{ end }}{{ end }}

    # Shared envars across user modes
    # For the RestFul API, we can't easily scale this up so MaxSize is largely ignored
    export FLUX_REQUIRE_AUTH=true
    export FLUX_SECRET_KEY={{ .Spec.FluxRestful.SecretKey}}
    export FLUX_NUMBER_NODES={{ .Spec.Size}}

    printf "\n 🔑 Use your Flux user and token credentials to authenticate with the MiniCluster with flux-framework/flux-restful-api\n"

    # -o is an "option" for the broker
    # -S corresponds to a shortened --setattr=ATTR=VAL 
    printf "\n🌀 ${asFlux} {{.Container.Commands.Prefix}} flux broker --config-path /etc/flux/config ${brokerOptions} ${startServer}\n"
    ${asFlux} {{.Container.Commands.Prefix}} flux broker --config-path /etc/flux/config ${brokerOptions} ${startServer}
}

# Run an interactive cluster, giving no command to flux start
function run_interactive_cluster() {
    printf "\n🌀 ${asFlux} {{.Container.Commands.Prefix}} flux broker --config-path /etc/flux/config ${brokerOptions}\n"
    ${asFlux} {{.Container.Commands.Prefix}} flux broker --config-path /etc/flux/config ${brokerOptions}
}

# Are we running diagnostics or the start command?
if [ "${diagnostics}" == "true" ]; then
    run_diagnostics
else

    # Start flux with the original entrypoint
    if [ $(hostname) == "${mainHost}" ]; then

        # If it's a batch job, we write the script for the broker to run
        {{ if .Container.Batch }}rm -rf flux-job.batch
        echo "#!/bin/bash
{{ if .Container.BatchRaw }}{{range $index, $line := .Batch}}{{ if $line }}{{$line}}{{ end }}
{{ end }}
{{ else }}{{range $index, $line := .Batch}}{{ if $line }}flux submit --flags waitable --error=${FLUX_OUTPUT_DIR}/job-{{$index}}.err --output=${FLUX_OUTPUT_DIR}/job-{{$index}}.out {{$line}}{{ end }}
{{ end }}
flux queue idle
flux jobs -a{{ end }}
" >> flux-job.batch
        chmod +x flux-job.batch
        {{ if not .Container.Commands.RunFluxAsRoot }}chown -R ${fluxuid} flux-job.batch{{ end }}
        {{ end }} # end if container batch

        # Commands only run by the broker
        {{ .Container.Commands.BrokerPre}} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

        # No command - use default to start server
{{ if not .Spec.Logging.Quiet }}        echo "Extra command arguments are: $@"{{ end }}
        if [ "$@" == "" ]; then

            # An interactive job also doesn't require a command
            {{ if .Spec.Interactive }}run_interactive_cluster
            {{ else }}run_flux_restful{{ end }}

        # Case 2: Fall back to provided command
        else

            # If we are running a batch job, no launcher mode
            {{ if .Container.Batch }}
            {{ if not .Spec.Logging.Quiet }}printf "✨️ Prepared Batch Job:\n"
            cat flux-job.batch
            {{ end }}

            flags="{{ if ge .Spec.Tasks .Spec.Size }} -N {{.Spec.Size}}{{ end }} -n {{.Spec.Tasks}} {{ if .Spec.Flux.OptionFlags }}{{ .Spec.Flux.OptionFlags}}{{ end }} {{ if .Spec.Logging.Debug }} -vvv{{ end }}"
            {{ if not .Spec.Logging.Quiet }}          
            printf "\n🌀 Batch Mode: flux start -o --config /etc/flux/config ${brokerOptions} {{.Container.Commands.Prefix}} sh -c 'flux batch ${flags} --flags waitable ./flux-job.batch && flux job wait --all'\n"
            {{ end }}
            {{ if .Spec.Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxstart wall time %E" {{ end }}${asFlux} flux start -o --config /etc/flux/config ${brokerOptions} {{ if .Spec.Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxsubmit wall time %E" {{ end }} {{.Container.Commands.Prefix}} sh -c "flux batch ${flags} --flags waitable ./flux-job.batch && flux job wait --all" {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

            {{ else }} # else for if container.batch
            {{ if not .Spec.Logging.Quiet }} # if tasks >= size
            # Container launchers are snakemake, nextflow, that will launch their own jobs
            {{ if .Container.Launcher }}
            printf "\n🌀 Launcher Mode: flux start -o --config /etc/flux/config ${brokerOptions} {{.Container.Commands.Prefix}} $@\n"
            {{ else }}
            printf "\n🌀 Submit Mode: flux start -o --config /etc/flux/config ${brokerOptions} {{.Container.Commands.Prefix}} flux submit {{ if ge .Spec.Tasks .Spec.Size }} -N {{.Spec.Size}}{{ end }} -n {{.Spec.Tasks}} --quiet {{ if .Spec.Flux.OptionFlags }}{{ .Spec.Flux.OptionFlags}}{{ end }} --watch{{ if .Spec.Logging.Debug }} -vvv{{ end }} $@\n"
            {{ end }}
{{ end }}
            {{ if .Container.Launcher }}
            {{ if .Spec.Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxstart wall time %E" {{ end }}${asFlux} flux start -o --config /etc/flux/config ${brokerOptions} {{ if .Spec.Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxsubmit wall time %E" {{ end }} {{.Container.Commands.Prefix}} $@
            {{ else }}
            {{ if .Spec.Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxstart wall time %E" {{ end }}${asFlux} flux start -o --config /etc/flux/config ${brokerOptions} {{ if .Spec.Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxsubmit wall time %E" {{ end }} {{.Container.Commands.Prefix}} flux submit {{ if ge .Spec.Tasks .Spec.Size }} -N {{.Spec.Size}}{{ end }} -n {{.Spec.Tasks}} --quiet {{ if .Spec.Flux.OptionFlags }}{{ .Spec.Flux.OptionFlags}}{{ end }} --watch{{ if .Spec.Logging.Debug }} -vvv{{ end }} $@
            {{ end }} # end if container.launcher
            {{ end }} # end if container.batch
        fi
    else

       # Commands only run by the workers
       {{ .Container.Commands.WorkerPre}} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

        # Sleep until the broker is ready
        {{ if not .Spec.Logging.Quiet }}printf "\n🌀 {{.Container.Commands.Prefix}} flux start -o --config /etc/flux/config ${brokerOptions}\n"{{ end }}
        while true
        do
            {{ if .Spec.Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxstart wall time %E" {{ end }}${asFlux} {{.Container.Commands.Prefix}} flux start -o --config /etc/flux/config ${brokerOptions}
            retval=$?
            {{ if not .Spec.Logging.Quiet }}printf "Return value for follower worker is ${retval}\n"{{ end }}
            if [[ "${retval}" -eq 0 ]]; then
                {{ if not .Spec.Logging.Quiet }}printf "The follower worker exited cleanly. Goodbye!\n"{{ end }}
                break
            fi
            {{ if not .Spec.Logging.Quiet }}printf "\n😪 Sleeping 15s until broker is ready..."{{ end }}
            sleep 15
        done
    fi

    {{ .Container.Commands.Post}}
fi

