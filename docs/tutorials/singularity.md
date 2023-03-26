# Singularity

We are exploring using Singularity containers as a reasonable way to (much more easily) package 
complex workflows. The reason is because we can run containers without requiring you to build a
container with Flux Framework and your software.

## Singularity Hello World

> This example is to pull and run a "hello world" example.

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/singularity/minicluster-hello-world.yaml)**

This example demonstrates pulling a Singularity container to all nodes prior to executing the SIF binary.

```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:

  # Number of pods to create for MiniCluster
  size: 2

  containers:
    - image: ghcr.io/rse-ops/singularity:tag-mamba
      fluxUser:
        name: fluxuser
      command: singularity exec ubuntu_latest.sif echo hello world

      # This pulls the container (separately) to each worker
      commands:
        pre: singularity pull docker://ubuntu
       
      # Running a container in a container
      securityContext:
        privileged: true
```

You would run the example as follows:

```bash
$ kubectl create namespace flux-operator
$ kubectl apply -f examples/singularity/minicluster-hello-world.yaml
```

We can then wait for our pods to be running

```bash
$ kubectl get -n flux-operator pods
```

And then look at the logs to see the print of "hello world"

```bash
$ kubectl logs -n flux-operator flux-sample-0-p5xls -f
```

But we can take a better approach. Ideally we can pull the container once to be shared by
all workers. We can do this with a `brokerPre` commands block shown in
[this tutorial file](https://github.com/flux-framework/flux-operator/blob/main/examples/singularity/minicluster-prepull.yaml).

```yaml
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:
  # suppress all output except for test run
  logging:
    quiet: true

  # Number of pods to create for MiniCluster
  size: 2

  # Make this kind of persistent volume and claim available to pods
  # This is a path in minikube (e.g., minikube ssh)
  volumes:
    data:
      storageClass: hostpath
      path: /tmp/pulls

  # This is a list because a pod can support multiple containers
  containers:
    - image: ghcr.io/rse-ops/singularity:tag-mamba
      command: singularity exec ./ubuntu_latest.sif echo hello world
      workingDir: /data

      # This pulls the container (once) by the broker to workingDir /data
      commands:
        brokerPre: singularity pull docker://ubuntu

      fluxUser:
        name: fluxuser

      # Container will be pre-pulled here only by the broker
      volumes:
        data:
          path: /data
       
      # Running a container in a container
      securityContext:
        privileged: true
```

Notice that we've added a local volume, and this is the (shared) working directory for the broker
and all workers. Since the workers wait for the broker, there is no issue with however much
time the broker needs for the pull. Since the pull is done to a shared, writable space,
all workers can access the binary. Finally, note that in order for this to work,
Singularity should be installed in the container, the MiniCluster is run in privileged mode, 
and `tzdata` is also needed so there is an `/etc/localtime` to bind. We could likely improve this
to cut down permissions, if/when someone wants to work on that!