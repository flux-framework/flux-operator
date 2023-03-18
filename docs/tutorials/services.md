# Services

This small tutorial will show you how to run a "sidecar" service container alongside your
flux install. 

## Sidecar Registry with ORAS

> Create an interactive MiniCluster with a sidecar registry container

As an example, we will run a local container registry to push/pull artifacts
with ORAS. I don't know why, I just like ORAS :) In all seriousness, you could
imagine interesting use cases like needing an API to save and get artifacts 
for your analysis.

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/services/minicluster-registry.yaml)**

This example demonstrates bringing up a MiniCluster and then interacting with a service (a registry)
to push / pull artifacts. Here is our example custom resource definition:

```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:

  # Number of pods to create for MiniCluster
  size: 2

  # Interactive so we can submit commands
  interactive: true

  # This is a list because a pod can support multiple containers
  containers:
    # The container URI to pull (currently needs to be public)
    - image: ghcr.io/flux-framework/flux-restful-api:latest
      runFlux: true
      commands:

        # This is going to install oras for the demo
        pre: |
          apt-get update && apt-get install -y curl
          VERSION="1.0.0-rc.2"
          curl -LO "https://github.com/oras-project/oras/releases/download/v${VERSION}/oras_${VERSION}_linux_amd64.tar.gz"
          mkdir -p oras-install/
          tar -zxf oras_${VERSION}_*.tar.gz -C oras-install/
          sudo mv oras-install/oras /usr/local/bin/
          rm -rf oras_${VERSION}_*.tar.gz oras-install/

      # This is our registry we want to run
    - image: ghcr.io/oras-project/registry:latest
      name: registry
      ports:
        - 5000

```

It's helpful to pull containers to MiniKube first:

```bash
$ minikube ssh docker pull ghcr.io/oras-project/registry:latest
$ minikube ssh docker pull ghcr.io/flux-framework/flux-restful-api:latest
```

When interactive is true, we tell the Flux broker to start without a command. This means
the cluster will remain running until you shutdown Flux or `kubectl delete` the MiniCluster
itself. The container you choose should have the software you are interested in having for each node.
Given a running cluster, we can create the namespace and the MiniCluster as follows:

```bash
$ kubectl create namespace flux-operator
```

And apply the MiniCluster CRD:
```bash
$ kubectl apply -f examples/services/minicluster-registry.yaml
```

If you are curious, the entrypoint for the service sidecar container is `registry serve /etc/docker/registry/config.yml`
to start the registry. Since it's not a flux runner, not providing an entrypoint means we use the container's default
entrypoint. We can then wait for our pods to be running

```bash
$ kubectl get -n flux-operator pods
NAME                         READY   STATUS      RESTARTS   AGE
flux-sample-0-p5xls          1/1     Running     0          7s
flux-sample-1-nmtt7          1/1     Running     0          7s
flux-sample-cert-generator   0/1     Completed   0          7s
```

To see logs, since we have 2 containers per pod, you can either leave out the pod (and get the first or default)
or specify a container with `-c`:

```bash
$ kubectl logs -n flux-operator flux-sample-0-d5jbb -c registry
$ kubectl logs -n flux-operator flux-sample-0-d5jbb -c flux-sample
$ kubectl logs -n flux-operator flux-sample-0-d5jbb
```

And then shell into the broker pod, index 0, which is "flux-sample"

```bash
$ kubectl exec -it  -n flux-operator flux-sample-0-d5jbb -c flux-sample -- bash
```

Let's first make and push an artifact. First, just using oras natively (no flux)

```bash
cd /tmp

# Assume we would be running from inside the flux instance
sudo -u flux echo "hello dinosaur" > artifact.txt
```

And push! The registry, by way of being a container in the same pod, is on port 5000:

At this point, remember the broker is running, and we need to connect to it. We do this via
flux proxy and targeting the socket, which is a local reference at `/run/flux/local`:

```bash
# Connect to the flux socket at /run/flux/local as the flux instance owner "flux"
$ sudo -u flux flux proxy local:///run/flux/local oras push localhost:5000/dinosaur/artifact:v1 \
   --artifact-type application/vnd.acme.rocket.config \
   ./artifact.txt
```
```console
Uploading 07f469745bff artifact.txt
Uploaded  07f469745bff artifact.txt
Pushed [registry] localhost:5000/dinosaur/artifact:v1
Digest: sha256:3a6cb1d1d1b1d80d4c4de6abc66a6c9b4f7fef0b117f87be87fea9b725053ead
```
Now try pulling, deleting the original first, and again without flux:

```bash
rm -f artifact.txt
sudo -u flux flux proxy local:///run/flux/local oras pull localhost:5000/dinosaur/artifact:v1
cat artifact.txt
```
```console
hello dinosaur
```

We did this under the broker (and flux user) assuming your actual use case will be running
in the Flux instance. Feel free to play with oras outside of that context!
When you are done, exit from the instance, and exit from the pod, and then delete the MiniCluster.

```bash
$ kubectl delete -f examples/services/minicluster-registry.yaml
```

That's it. Please do something more useful than my terrible example.