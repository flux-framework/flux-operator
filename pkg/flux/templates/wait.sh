#!/bin/sh

# If we are not in strict, don't set strict mode
{{ if .Spec.Logging.Strict }}set -eEu -o pipefail{{ end }}

# We use the actual time command and not the wrapper, otherwise we get there is no argument -f
{{ if .Spec.Logging.Timed }}which /usr/bin/time > /dev/null 2>&1 || (echo "/usr/bin/time is required to use logging.timed true" && exit 1);{{ end }}

# If any initCommand logic is defined
{{ .Container.Commands.Init}} {{ if .Spec.Logging.Quiet }}> /dev/null{{ end }}

# This waiting script is intended to wait for the flux view, and then start running
# Ensure the flux volume addition is complete.
curl -L -O -s https://github.com/converged-computing/goshare/releases/download/2023-09-06/wait-fs {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
chmod +x ./wait-fs
mv ./wait-fs /usr/bin/goshare-wait-fs

# Ensure spack view is on the path, wherever it is mounted
viewbase="{{ .ViewBase }}"
viewroot=${viewbase}/view
software="${viewbase}/software"
viewbin="${viewroot}/bin"
fluxpath=${viewbin}/flux

# Set the flux root
{{ if not .Spec.Logging.Quiet }}
echo
echo "Flux install root: ${viewroot}"
echo
{{ end }}

# Important to add AFTER in case software in container duplicated
export PATH=$PATH:${viewbin}

# Wait for marker (from spack.go) to indicate copy is done
goshare-wait-fs -p ${viewbase}/flux-operator-done.txt {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

# Copy mount software to /opt/software
cp -R ${viewbase}/software /opt/software

# And pre command logic that isn't passed to the certificate generator
{{ .Container.Commands.Pre}} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

# Set the flux user and id from the getgo
fluxuser=$(whoami)
fluxuid=$(id -u $fluxuser)

# Variables we can use again
cfg="${viewroot}/etc/flux/config"
command="{{ .Container.Command }}"

{{ if not .Spec.Logging.Quiet }}
echo 
echo "Hello user ${fluxuser}"{{ end }}

# Add fluxuser to sudoers living... dangerously!
if [[ "${fluxuser}" != "root" ]]; then
  echo "${fluxuser} ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers
fi
		
# Ensure the flux user owns the curve.cert
# We need to move the curve.cert because config map volume is read only
curvesrc=/flux_operator/curve.cert
curvepath=$viewroot/curve/curve.cert

mkdir -p $viewroot/curve
cp $curvesrc $curvepath
{{ if not .Spec.Logging.Quiet }}
echo 
echo "🌟️ Curve Certificate"
ls $viewroot/curve
cat ${curvepath}
{{ end }}

# Remove group and other read
chmod o-r ${curvepath}
chmod g-r ${curvepath}
chown -R ${fluxuid} ${curvepath}

foundroot=$(find $viewroot -maxdepth 2 -type d -path $viewroot/lib/python3\*)

# Ensure we use flux's python (TODO update this to use variable)
export PYTHONPATH={{ if .Spec.Flux.Container.PythonPath }}{{ .Spec.Flux.Container.PythonPath }}{{ else }}${foundroot}/site-packages{{ end }}
echo "PYTHONPATH is ${PYTHONPATH}" {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
echo "PATH is $PATH" {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

# Put the state directory in /var/lib on shared view
export STATE_DIR=${viewroot}/var/lib/flux
export FLUX_OUTPUT_DIR={{ if .Container.Logs }}{{.Container.Logs}}{{ else }}/tmp/fluxout{{ end }}
mkdir -p ${STATE_DIR} ${FLUX_OUTPUT_DIR}

# Main host <name>-0 and the fully qualified domain name
mainHost="{{ .MainHost }}"
workdir=$(pwd)

{{ if .Spec.Logging.Quiet }}{{ else }}
echo "👋 Hello, I'm $(hostname)"
echo "The main host is ${mainHost}"

echo "The working directory is ${workdir}, contents include:"
ls .
{{ end }}

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


function run_flux_restful() {

    # Start restful API server
    branch={{if .Spec.FluxRestful.Branch}}{{.Spec.FluxRestful.Branch}}{{else}}main{{end}}
    startServer="uvicorn app.main:app --host=0.0.0.0 --port={{or .Spec.FluxRestful.Port 5000}}"
    printf "Cloning flux-framework/flux-restful-api branch ${branch}\n"
    git clone -b ${branch} --depth 1 https://github.com/flux-framework/flux-restful-api ./flux-restful-api > /dev/null 2>&1
    cd ./flux-restful-api
            
    # Export the main flux user and token "superuser"
    export FLUX_USER=$(whoami)
    export FLUX_TOKEN={{ .FluxToken }}
    printf "🔒️ Credentials, my friend!\n    FLUX_USER: ${FLUX_USER}\n    FLUX_TOKEN: ${FLUX_TOKEN}\n\n"

    # Install python requirements, with preference for python3
    python3 -m pip install -r requirements.txt > /dev/null 2>&1 || python -m pip install -r requirements.txt > /dev/null 2>&1

    # Prepare databases!
    alembic revision --autogenerate -m "Create intital tables"
    alembic upgrade head
    python3 app/db/init_db.py init || python app/db/init_db.py init

    # Shared envars across user modes
    # For the RestFul API, we can't easily scale this up so MaxSize is largely ignored
    export FLUX_REQUIRE_AUTH=true
    export FLUX_SECRET_KEY={{ .Spec.FluxRestful.SecretKey }}
    export FLUX_NUMBER_NODES={{ .Spec.Size }}

    printf "\n 🔑 Use your Flux user and token credentials to authenticate with the MiniCluster with flux-framework/flux-restful-api\n"

    # -o is an "option" for the broker
    # -S corresponds to a shortened --setattr=ATTR=VAL 
    printf "\n🌀 {{.Container.Commands.Prefix}} flux broker --config-path ${cfg} ${brokerOptions} ${startServer}\n"
    {{.Container.Commands.Prefix}} flux broker --config-path /etc/flux/config ${brokerOptions} ${startServer}
}

# Run an interactive cluster, giving no command to flux start
function run_interactive_cluster() {
    echo "🌀 flux broker --config-path ${cfg} ${brokerOptions}"
    flux broker --config-path ${cfg} ${brokerOptions}
}

# if we are given an archive to use, load first, not required to exist
# Note that we ask the user to dump in interactive mode - I am not
# sure that doing it with a hook ensures the dump will be successful.
{{if .Spec.Archive.Path }}
if [[ -e "{{ .Spec.Archive.Path}}" ]]; then
{{ if not .Spec.Logging.Quiet }}printf "🧊️ Found existing archive at {{ .Spec.Archive.Path}} loading into state directory\nBefore:\n"{{ end }}
brokerOptions="${brokerOptions} -Scontent.restore={{ .Spec.Archive.Path}}"
fi{{ end }}

# And pre command logic that isn't passed to the certificate generator
{{ .Container.Commands.Pre}} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

# Flux option flags
{{ if not .Spec.Logging.Quiet }}echo "🚩️ Flux Option Flags defined"{{ end }}

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
    chown -R ${fluxuid} flux-job.batch{{ end }}

    # Commands only run by the broker
    {{ .Container.Commands.BrokerPre}} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

    echo "Command provided is: ${command}" {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
    if [ "${command}" == "" ]; then

       # An interactive job also doesn't require a command
       {{ if .Spec.Interactive }}run_interactive_cluster{{ else }}run_flux_restful{{ end }}

    else
    
      # If we are running a batch job, no launcher mode
      {{ if .Container.Batch }}
        {{ if not .Spec.Logging.Quiet }}printf "✨️ Prepared Batch Job:\n"
        cat flux-job.batch
        {{ end }}

        flags="{{ if ge .Spec.Tasks .Spec.Size }} -N {{.Spec.Size}}{{ end }} -n {{.Spec.Tasks}} {{ if .Spec.Flux.OptionFlags }}{{ .Spec.Flux.OptionFlags}}{{ end }} {{ if .Spec.Logging.Debug }} -vvv{{ end }}"
        {{ if not .Spec.Logging.Quiet }}          
        printf "\n🌀 Batch Mode: flux start {{ if .Spec.Flux.Wrap }}--wrap={{ .Spec.Flux.Wrap }} {{ end }}-o --config ${cfg} ${brokerOptions} {{.Container.Commands.Prefix}} sh -c 'flux batch ${flags} --flags waitable ./flux-job.batch && flux job wait --all'\n"
        {{ end }}
        {{ if .Spec.Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxstart wall time %E" {{ end }}flux start {{ if .Spec.Flux.Wrap }}--wrap={{ .Spec.Flux.Wrap }} {{ end }}-o --config ${cfg} ${brokerOptions} {{ if .Spec.Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxsubmit wall time %E" {{ end }} {{.Container.Commands.Prefix}} sh -c "flux batch ${flags} --flags waitable ./flux-job.batch && flux job wait --all" {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

      {{ else }} # else for if container.batch
        {{ if not .Spec.Logging.Quiet }} # if tasks >= size
        # Container launchers are snakemake, nextflow, that will launch their own jobs
        {{ if .Container.Launcher }}
        printf "\n🌀 Launcher Mode: flux start {{ if .Spec.Flux.Wrap }}--wrap={{ .Spec.Flux.Wrap }} {{ end }}-o --config ${cfg} ${brokerOptions} {{.Container.Commands.Prefix}} $@\n"
        {{ else }}
        printf "\n🌀 Submit Mode: flux start {{ if .Spec.Flux.Wrap }}--wrap={{ .Spec.Flux.Wrap }} {{ end }}-o --config ${cfg} ${brokerOptions} {{.Container.Commands.Prefix}} {{ if .Spec.Flux.SubmitCommand }}{{ .Spec.Flux.SubmitCommand }}{{ else }}flux submit {{ end }} {{ if ge .Spec.Tasks .Spec.Size }} -N {{.Spec.Size}}{{ end }} -n {{.Spec.Tasks}} --quiet {{ if .Spec.Flux.OptionFlags }}{{ .Spec.Flux.OptionFlags}}{{ end }} --watch{{ if .Spec.Logging.Debug }} -vvv{{ end }} ${command}\n"
        {{ end }}
      {{ end }}

      {{ if .Container.Launcher }}
      {{ if .Spec.Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxstart wall time %E" {{ end }}flux start {{ if .Spec.Flux.Wrap }}--wrap={{ .Spec.Flux.Wrap }} {{ end }}-o --config ${cfg} ${brokerOptions} {{ if .Spec.Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxsubmit wall time %E" {{ end }} {{.Container.Commands.Prefix}} ${command}
      {{ else }}
      {{ if .Spec.Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxstart wall time %E" {{ end }}flux start {{ if .Spec.Flux.Wrap }}--wrap={{ .Spec.Flux.Wrap }} {{ end }} -o --config ${cfg} ${brokerOptions} {{ if .Spec.Logging.Timed }}/usr/bin/time -f "FLUXTIME fluxsubmit wall time %E" {{ end }} {{.Container.Commands.Prefix}} {{ if .Spec.Flux.SubmitCommand }}{{ .Spec.Flux.SubmitCommand }}{{ else }}flux submit {{ end }} {{ if ge .Spec.Tasks .Spec.Size }} -N {{.Spec.Size}}{{ end }} -n {{.Spec.Tasks}} --quiet {{ if .Spec.Flux.OptionFlags }}{{ .Spec.Flux.OptionFlags}}{{ end }} --watch{{ if .Spec.Logging.Debug }} -vvv{{ end }} ${command}
      {{ end }} # end if container.launcher
      {{ end }} # end if container.batch
    fi

# Block run by workers
else

   # Commands only run by the workers
   {{ .Container.Commands.WorkerPre}} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

    # We basically sleep/wait until the lead broker is ready
    echo "🌀 flux start {{ if .Spec.Flux.Wrap }}--wrap={{ .Spec.Flux.Wrap }} {{ end }} -o --config ${viewroot}/etc/flux/config ${brokerOptions}"

    # We can keep trying forever, don't care if worker is successful or not
    while true
    do
        flux start -o --config ${viewroot}/etc/flux/config ${brokerOptions}
        retval=$?
        if [[ "${retval}" -eq 0 ]]; then
             echo "The follower worker exited cleanly. Goodbye!"
             break
        fi
        echo "Return value for follower worker is ${retval}"
        echo "😪 Sleeping 15s to try again..."
        sleep 15
    done

fi

{{ .Container.Commands.Post}}