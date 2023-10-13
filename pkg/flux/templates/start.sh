#!/bin/sh

# A custom startscript can be supported for a non flux runner given that
# the container also provides the entrypoint command to run. To be consitent,
# we provide the same blocks of commands as we do to wait.sh.

# If any initCommand logic is defined
{{ .Container.Commands.Init}} {{ if .Spec.Logging.Quiet }}> /dev/null{{ end }}

# If we are not in strict, don't set strict mode
{{ if .Spec.Logging.Strict }}set -eEu -o pipefail{{ end }}

# Shared logic to wait for view
{{template "wait-view" .}}

{{template "paths" .}}

{{ .Container.Commands.BrokerPre}} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
{{ .Container.Commands.WorkerPre}} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
{{ .Container.Commands.Pre}} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

# Ensure socket path is envar for user
fluxsocket=${viewroot}/run/flux/local

# Wait for it to exist (application is running)
goshare-wait-fs -p ${fluxsocket} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

# Ensure fluxsocket has local
fluxsocket="local://$fluxsocket"

{{ .Container.Command }}

{{ .Container.Commands.Post}}
