#!/bin/sh

# This script handles start logic for the broker
{{template "init" .}}

# Are we running diagnostics or the start command?
if [ "${diagnostics}" == "true" ]; then
    run_diagnostics
else

    # Start flux with the original entrypoint
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
    {{ .Container.Commands.Post}}
fi