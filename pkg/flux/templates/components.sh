{{define "flags"}}

# If tasks are == 0, then only define nodes
node_spec="{{ if ge .Spec.Tasks .Spec.Size }}-N {{.Spec.Size}} {{ end }}-n{{.Spec.Tasks}}"
node_spec="{{ if eq .Spec.Tasks 0 }}-N {{.Spec.Size}}{{ else }}${node_spec}{{ end }}"
flags="${node_spec} {{ if .Spec.Flux.OptionFlags }}{{ .Spec.Flux.OptionFlags}}{{ end }} {{ if .Spec.Logging.Debug }} -vvv{{ end }}"
echo "Flags for flux are ${flags}" {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
{{ end }}

{{define "wait-view"}}

# This is the baseurl for a wait script.
goshareUrl=https://github.com/converged-computing/goshare/releases/download/2024-01-18

# Ensure the flux volume addition is complete. We default to linux, fall back to arm
url=$goshareUrl/wait-fs
{{ if .Spec.Flux.Arch }}
url=$goshareUrl/wait-fs-{{ .Spec.Flux.Arch }}
{{ end }}

# This waiting script is intended to wait for the flux view, and then start running
curl -L -O -s -o ./wait-fs -s ${url} {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }} || wget ${url} -q -O ./wait-fs {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }} || true
chmod +x ./wait-fs || true
${SUDO} mv ./wait-fs /usr/bin/goshare-wait-fs || true

# Ensure spack view is on the path, wherever it is mounted
viewbase="{{ .ViewBase }}"
viewroot=${viewbase}/view
software="${viewbase}/software"
viewbin="${viewroot}/bin"
fluxpath=${viewbin}/flux

# Set the flux root, don't show the viewer if view is disabled (can be confusing)
# The view is used to have configs and that is it
{{ if not .Spec.Logging.Quiet }}{{ if not .Spec.Flux.Container.Disable }}
echo
echo "Flux install root: ${viewroot}"
echo
{{ end }}{{ end }}

# Important to add AFTER in case software in container duplicated
{{ if not .Spec.Flux.Container.Disable }}export PATH=$PATH:${viewbin}{{ end }}

# Wait for marker (from spack.go) to indicate copy is done
goshare-wait-fs -p ${viewbase}/flux-operator-done.txt {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }} || true

# Copy mount software to /opt/software
# If /opt/software already exists, we need to copy into it
if [[ -e  "/opt/software" ]]; then
  ${SUDO} cp -R ${viewbase}/software/* /opt/software/ || true
else
  ${SUDO} cp -R ${viewbase}/software /opt/software || true
fi
{{end}}

{{define "custom-script"}}
{{ if .Container.Commands.Script }}
cat <<EOF > ./custom-entrypoint.sh
{{ .Container.Commands.Script }}
EOF
chmod +x ./custom-entrypoint.sh
command="/bin/bash ./custom-entrypoint.sh"
{{end}}
{{end}}

{{define "paths"}}
foundroot=$(find $viewroot -maxdepth 2 -type d -path $viewroot/lib/python3\*) {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
pythonversion=$(basename ${foundroot})
pythonversion=${viewroot}/bin/${pythonversion}
echo "Python version: $pythonversion" {{ if .Spec.Logging.Quiet }} > /dev/null 2>&1{{ end }}
echo "Python root: $foundroot" {{ if .Spec.Logging.Quiet }} > /dev/null 2>&1{{ end }}

# If we found the right python, ensure it's linked (old link does not work)
if [[ -f "${pythonversion}" ]]; then
   ${SUDO} rm -rf $viewroot/bin/python3
   ${SUDO} rm -rf $viewroot/bin/python
   ${SUDO} ln -s ${pythonversion} $viewroot/lib/python  || true
   ${SUDO} ln -s ${pythonversion} $viewroot/lib/python3 || true
fi

# Ensure we use flux's python (TODO update this to use variable)
export PYTHONPATH=$PYTHONPATH:{{ if .Spec.Flux.Container.PythonPath }}{{ .Spec.Flux.Container.PythonPath }}{{ else }}${foundroot}/site-packages{{ end }}
echo "PYTHONPATH is ${PYTHONPATH}" {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
echo "PATH is $PATH" {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}

find $viewroot . -name libpython*.so* {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
ls -l /mnt/flux/view/lib/libpython3.11.so.1.0 {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
{{ if .Spec.Flux.Scheduler.Simple }}{{ else }}export FLUX_RC_EXTRA=$viewroot/etc/flux/rc1.d{{ end }}

# Write a script to load fluxion
cat <<EOT >> ./load-fluxion.sh
flux module remove sched-simple
flux module load sched-fluxion-resource
flux module load sched-fluxion-qmanager
EOT
${SUDO} mv ./load-fluxion.sh ${viewbase}/load-fluxion.sh

# Write an easy file we can source for the environment
cat <<EOT >> ./flux-view.sh
#!/bin/bash
export PATH=$PATH
export PYTHONPATH=$PYTHONPATH
export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:$viewroot/lib
export fluxsocket=local://${viewroot}/run/flux/local
EOT
${SUDO} mv ./flux-view.sh ${viewbase}/flux-view.sh

# The same, but also connect
cat <<EOT >> ./flux-connect.sh
#!/bin/bash
. ${viewbase}/flux-view.sh
flux proxy ${fluxsocket} bash
EOT
${SUDO} mv ./flux-connect.sh ${viewbase}/flux-connect.sh


{{end}}
{{define "ensure-pip"}}
${SUDO} ${pythonversion} -m pip --version || ${SUDO} ${pythonversion} -m ensurepip || (${SUDO} wget https://bootstrap.pypa.io/get-pip.py && ${pythonversion} ./get-pip.py) {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
${SUDO} ${pythonversion} -m pip --upgrade pip {{ if .Spec.Logging.Quiet }}> /dev/null 2>&1{{ end }}
{{end}}