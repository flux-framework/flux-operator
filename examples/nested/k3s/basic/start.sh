#!/bin/bash

name=$(hostname)
echo "Hello I am ${name}"
if [[ "${name}" == "flux-sample-0" ]]; then
    export K3S_TOKEN=secret
    export K3S_KUBECONFIG_OUTPUT=/workflow/kubeconfig.yaml
    export K3S_KUBECONFIG_MODE=666
    /bin/k3s server
else
    export K3S_URL="https://flux-sample-0.flux-service.flux-operator.svc.cluster.local:6443"
    export K3S_TOKEN=secret
    /bin/k3s agent
fi