# Shared logic to install hq
{{define "wait-view"}}

# This waiting script is intended to wait for the flux view, and then start running
# Ensure the flux volume addition is complete.
url=https://github.com/converged-computing/goshare/releases/download/2023-09-06/wait-fs
curl -L -O -s ${url} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }} || wget ${url} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
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
cp -R ${viewbase}/software /opt/software
{{end}}

{{define "paths"}}
foundroot=$(find $viewroot -maxdepth 2 -type d -path $viewroot/lib/python3\*)

# Ensure we use flux's python (TODO update this to use variable)
export PYTHONPATH={{ if .Spec.Flux.Container.PythonPath }}{{ .Spec.Flux.Container.PythonPath }}{{ else }}${foundroot}/site-packages{{ end }}
echo "PYTHONPATH is ${PYTHONPATH}" {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
echo "PATH is $PATH" {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

# Write an easy file we can source for the environment
cat <<EOT >> ${viewbase}/flux-view.sh
#!/bin/bash
export PATH=$PATH
export PYTHONPATH=$PYTHONPATH
export fluxsocket=local://${viewroot}/run/flux/local
EOT
{{end}}