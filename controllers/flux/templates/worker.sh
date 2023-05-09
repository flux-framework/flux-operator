#!/bin/sh

# The start script for a Flux cluster worker
{{template "init" .}}

# Are we running diagnostics or the start command?
if [ "${diagnostics}" == "true" ]; then
    run_diagnostics
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
    {{ .Container.Commands.Post}}
fi

