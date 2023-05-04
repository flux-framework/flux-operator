# Services

These tutorials will show you how to run a "sidecar" service container (one per Flux node) alongside your
flux install, along with a service for the entire cluster (a deployment next to the cluster).

## Sidecar Tutorials

### Sidecar NGINX Container

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/tests/nginx-sidecar-service/minicluster.yaml)**

This is one of the simplest examples, implemented as a test, to run a sidecar with NGINX and then curl localhost
to get a response from flux.

```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:

  logging:
    quiet: true

  # Number of pods to create for MiniCluster
  size: 2

  # This is a list because a pod can support multiple containers
  containers:
    - image: ghcr.io/flux-framework/flux-restful-api:latest
      runFlux: true
      command: curl -s localhost
      commands:
        pre: apt-get update > /dev/null && apt-get install -y curl > /dev/null
    - image: nginx
      name: nginx
      ports:
        - 80
```

Create it (after you have the flux-operator namespace):

```bash
$ kubectl create -f ./examples/tests/nginx-sidecar-service/minicluster.yaml
```

See nginx is running:

```bash
$ kubectl -n flux-operator logs flux-sample-0-zlpwx -c nginx -f
```
```console
/docker-entrypoint.sh: /docker-entrypoint.d/ is not empty, will attempt to perform configuration
/docker-entrypoint.sh: Looking for shell scripts in /docker-entrypoint.d/
/docker-entrypoint.sh: Launching /docker-entrypoint.d/10-listen-on-ipv6-by-default.sh
10-listen-on-ipv6-by-default.sh: info: Getting the checksum of /etc/nginx/conf.d/default.conf
10-listen-on-ipv6-by-default.sh: info: Enabled listen on IPv6 in /etc/nginx/conf.d/default.conf
/docker-entrypoint.sh: Launching /docker-entrypoint.d/20-envsubst-on-templates.sh
/docker-entrypoint.sh: Launching /docker-entrypoint.d/30-tune-worker-processes.sh
/docker-entrypoint.sh: Configuration complete; ready for start up
2023/03/18 05:01:31 [notice] 1#1: using the "epoll" event method
2023/03/18 05:01:31 [notice] 1#1: nginx/1.23.3
2023/03/18 05:01:31 [notice] 1#1: built by gcc 10.2.1 20210110 (Debian 10.2.1-6) 
2023/03/18 05:01:31 [notice] 1#1: OS: Linux 5.15.0-67-generic
2023/03/18 05:01:31 [notice] 1#1: getrlimit(RLIMIT_NOFILE): 1048576:1048576
2023/03/18 05:01:31 [notice] 1#1: start worker processes
2023/03/18 05:01:31 [notice] 1#1: start worker process 29
2023/03/18 05:01:31 [notice] 1#1: start worker process 30
2023/03/18 05:01:31 [notice] 1#1: start worker process 31
...
```

And then look at the main logs to see the output of curl, run by Flux:

```bash
$ kubectl -n flux-operator logs flux-sample-0-zlpwx -c flux-sample -f
```
```console
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
html { color-scheme: light dark; }
body { width: 35em; margin: 0 auto;
font-family: Tahoma, Verdana, Arial, sans-serif; }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

And that's it! In a real world use case, you'd have some service running alongside
an analysis. Clean up:

```bash
$ kubectl delete -f ./examples/tests/nginx-sidecar-service/minicluster.yaml
```

### Sidecar Registry with ORAS

> Create an interactive MiniCluster with a sidecar registry container

As an example, we will run a local container registry to push/pull artifacts
with ORAS. I don't know why, I just like ORAS :) In all seriousness, you could
imagine interesting use cases like needing an API to save and get artifacts 
for your analysis.

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/services/sidecar/minicluster-registry.yaml)**

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
$ kubectl apply -f examples/services/sidecar/minicluster-registry.yaml
```

If you are curious, the entrypoint for the service sidecar container is `registry serve /etc/docker/registry/config.yml`
to start the registry. Since it's not a flux runner, not providing an entrypoint means we use the container's default
entrypoint. We can then wait for our pods to be running

```bash
$ kubectl get -n flux-operator pods
NAME                         READY   STATUS      RESTARTS   AGE
flux-sample-0-p5xls          1/1     Running     0          7s
flux-sample-1-nmtt7          1/1     Running     0          7s
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
$ kubectl delete -f examples/services/sidecar/minicluster-registry.yaml
```

That's it. Please do something more useful than my terrible example.

## Service Containers Alongside the Cluster

### Registry Service

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/services/single/minicluster-registry.yaml)**

Unlike the example above, it's more reasonable that you would want a single registry that your pods can access. E.g., perhaps
you use it like a pull-through cache - first pulling to this service node (or pushing from another pod) and then having
your worker pods pull from there. Let's do that now. If you are using MiniKube, remember to pull large
containers first. Create the namespace:

```bash
$ kubectl create namespace flux-operator
```

And apply the MiniCluster CRD:

```bash
$ kubectl apply -f examples/services/single/minicluster-registry.yaml
```

You'll see a services pod! The current design deploys one pod to share your services,
and the services share the same cluster networking space.

```bash
$ kubectl get -n flux-operator pods
NAME                         READY   STATUS      RESTARTS   AGE
flux-sample-0-4wt26          1/1     Running     0          38s
flux-sample-1-sr5zx          1/1     Running     0          38s
flux-sample-services         1/1     Running     0          38s
```

And then shell into the broker pod, index 0, which is "flux-sample"

```bash
$ kubectl exec -it  -n flux-operator flux-sample-0-d5jbb -- bash
```

The registry hostname should be in the environment (it's defined in the CRD):

```bash
# echo $REGISTRY 
flux-sample-services.flux-service.flux-operator.svc.cluster.local
```

Both oras and Singularity are installed. Conceptually, we should be able to pull a container
with Singularity, and push it to the registry with oras.  This is a relatively small container
and should be quick:

```bash
$ singularity pull docker://vanessa/salad
```
Give it a run!

```bash
$ singularity run salad_latest.sif 
```
```console
 In Go an array is a slice. Utensil discrimination!  

          _________________  .========
         [_________________>< :======
                             '======== 
```

Now let's push to the oras registry

```bash
$ oras push $REGISTRY:5000/vanessa/salad:latest --artifact-type appliciation/vnd.sylabs.sif.layer.tar ./salad_latest.sif  --plain-http
```

Great! Now you could, theoretically, push a single SIF to your registry (as a local cache) and have the other nodes pull it!
Here is an example, shelling in to the second worker:

```bash
$ oras pull $REGISTRY:5000/vanessa/salad:latest --plain-http
```

Super cool! We will have more examples in the `examples` folder of how this can be used for workflow containers.


#### Development Notes

I did some digging into the logic, and found that the underlying submission was a flux submit -> flux exec
to start a celery worker:

```bash
$ flux alloc -N 2 --exclusive --job-name=merlin flux exec `which /bin/bash` -c "celery -A merlin worker -l INFO --concurrency 1 --prefetch-multiplier 1 -Ofair -Q \'[merlin]_flux_par\'"
```
I think this should be changed to:

```bash
$ flux alloc -N 2 --exclusive --job-name=merlin /bin/bash -c "celery -A merlin worker -l INFO --concurrency 1 --prefetch-multiplier 1 -Ofair -Q \'[merlin]_flux_par\'"
```

 - I don't think we need flux exec
 - Why would there be more than one /bin/bash?

I don't fully understand the relationship between the celery queue and Flux - I think Flux should be used to submit jobs directly to,
as opposed to just using it to start a celery working. It also seems like there is one too many layers of complexity. If we have a Flux queue
why do we also need a celery queue?
