# Mini Cluster

> The CRD "Custom Resource Definition" defines a Mini Cluster

A CRD is a yaml file that you can apply to your cluster (with the Flux Operator
installed) to ask for a Mini Cluster. Kubernetes has these [custom resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
to make it easy to automate tasks, and in fact this is the goal of an operator!
A Kubernetes operator is conceptually like a human operator that takes your CRD,
looks at the cluster state, and does whatever is necessary to get your cluster state
to match your request. In the case of the Flux Operator, this means creating the resources
for a MiniCluster. This document describes the spec of our custom resource definition.
Development examples can be found under [config/samples](https://github.com/flux-framework/flux-operator/tree/main/config/samples) 
in the repository. We will have more samples soon, either in that directory or
separately in the [flux-hpc](https://github.com/rse-ops/flux-hpc) repository.

## Custom Resource Definition

### Header

The yaml spec will normally have an API version, the kind `MiniCluster` and then
a name and namespace to identify the custom resource definition followed by the spec for it.

```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:
  ...
```

### Spec

Under the spec, there are several variables to define. Descriptions are included below, and we
recommend that you look at [config/samples](https://github.com/flux-framework/flux-operator/tree/main/config/samples) 
in the repository  and the [flux-hpc](https://github.com/rse-ops/flux-hpc) repository to see
more.

### size

The `size` variable under the spec is the number of pods to create for the MiniCluster. Each pod includes
the set of containers that you describe.

```yaml
  # Number of pods to create for MiniCluster
  size: 4
```

### diagnostics

Flux has a command that makes it easy to run diagnostics on a cluster, and we expose a boolean that makes it possible
to run that (instead of your job or starting the server). To enable this, set this boolean to true. By default, it is false.

```yaml
  # Diagnostics runs flux commands for diagnostics, and a final sleep command
  # That makes it easy for you to shell into the pod to look around
  diagnostics: false
```

### deadline

This is the maximum running time for your job. If you leave unset, it is essentially infinite.
If a Job is suspended (at creation or through an update), this timer will effectively be stopped and reset when the Job is resumed again.

```yaml
  # Deadline in seconds, if not set there is no deadline
  deadlineSeconds: 100
```

### localDeploy

This is a boolean to indicate that you are doing a local deploy. What it is really determining is if we should
ask for a volume mount (e.g., binding to a temporary directory) vs. a volume claim (typically only available)
in production level clusters. If you are developing or working locally using, for example, MiniKube, you
likely want to set this to True.

```yaml
  # Set to true to use volume mounts instead of volume claims
  localDeploy: true
``` 

### containers

Early on we identified that a job could include more than one container, where there might be a primary container
running Flux, and others that provide services. 

```yaml
  containers:
    # The container URI to pull (currently needs to be public)
    - image: ghcr.io/rse-ops/lammps:flux-sched-focal-v0.24.0
    ...
```

For each container, the follow variables are available (nested under `containers` as a list, as shown above).

#### image

This is the only required attribute! You *must* provide a container base that has Flux.
The requirments of your container are defined in the README of the [flux-hpc](https://github.com/rse-ops/flux-hpc/)
repository. Generally speaking, you need to have Flux executables, Flux Python bindings,
and your own executables on the path, and should be started with root with a flux user.
If you use the [fluxrm/flux-sched](https://hub.docker.com/r/fluxrm/flux-sched) 
base containers this is usually a good start. 

#### command

Providing (or not providing) a command is going to dictate the behavior of your MiniCluster!

1. Providing a custom command means the MiniCluster is ephemeral - it will run the command and clean up.
2. Not providing a command means that we will create a persistent MiniCluster running a RESTFul API service (and GUI) to submit jobs to.

```yaml
    # Don't set a command unless you want to forgo running the restful server to submit
    # commands to! E.g., instead of starting the server, it will just run your job command.
    command: lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite
```

#### imagePullSecret

If your container has a pull secret, define it as `imagePullSecret`. If it's publicly pullable,
you don't need this. But we do hope you are able to practice open science and share your containers!

```yaml
    # Name of an already created ImagePullSecret for the image specfied above
    imagePullSecret: flux-image-secret
```

#### workingDir

The container likely has a set working directory, and if you are running the RESTful API service (meaning
you start without a command, as shown above) this will likely be the application folder. If you are lauching
a job directly with flux start and require a particular working directory, set it here!

```yaml
    # You can set the working directory if your container WORKDIR is not correct.
    workingDir: /home/flux/examples/reaxff/HNS
```    

Remember that if you don't provide a command and launch the RESTFul API, you can provide the working
directory needed on the level of each job submit, and you don't need to define it here.
In fact, if you are using the flux-restful-api server, it will be changed anyway.

#### pullAlways

For development, it can be helpful to request that an image is re-pulled. Control that using `pullAlways`:

```yaml
    # Always pull the image (if you are updating the image between runs, set to true)!
    pullAlways: false
```

#### runFlux

If you are running multiple containers in a pod, this boolean indicates the one that should
be running Flux (and the rest are providing services).
This defaults to true, so if you have one container, you largely don't need to worry about this.
However, if you set this to true for *two* container (not allowed currently) you will get an eror message.

```yaml
    # This defaults to true - this is the container we want to run flux in. This means
    # that if you have more than one container, set the non-flux runners to false.
    # For one container, you can leave this unset for the default. This will be
    # validated in case you make a mistake :)
    runFlux: true
```