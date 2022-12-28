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

### test

Test mode turns off all verbose output (yes, the emojis too) so only the output of 
your job will be printed to the console. This way, you can retrieve the job lob
and then determine if the test was successful based on this output.

```yaml
  test: true
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

### volumes

Volumes can be defined on the level of the MiniCluster that are then used by containers.
These volumes are local host volumes, and should be named (the key for the section)
with a path:

```yaml
volumes:
  myvolume:
    path: /full/path/to/volume
    readOnly: false
```

By default they will be read only unless you set `readOnly` to false.
Since we haven't implemented this for a cloud resource yet, this currently just works
with localDeploy is set to true, and we can adjust this when we test in a cloud.

### containers

Early on we identified that a job could include more than one container, where there might be a primary container
running Flux, and others that provide services. Note that currently we only allow one container to be a FluxRunner,
however we anticipate this could change (and allow for specifying custom logic for a flux runner entrypoint, a script
called "wait.sh") on the level of the container.

```yaml
  containers:
    # The container URI to pull (currently needs to be public)
    - image: ghcr.io/rse-ops/lammps:flux-sched-focal-v0.24.0
    ...
```

For each container, the follow variables are available (nested under `containers` as a list, as shown above).


#### name

For all containers that aren't flux runners, a name is required. Validation will check that it is defined.

```yaml
name: rabbit
```

#### image

This is the only required attribute! You *must* provide a container base that has Flux.
The requirements of your container are defined in the README of the [flux-hpc](https://github.com/rse-ops/flux-hpc/)
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

### volumes

Volumes that are defined on the level of the MiniCluster (named) can be mounted into containers.
As an example, here is how we specify the volume `myvolume` to be mounted to the container at `/data`.

```yaml
volumes:
  myvolume:
    path: /data
```

The `myvolume` key must be defined in the MiniCluster set of volumes, and this is checked.


#### imagePullSecret

If your container has a pull secret, define it as `imagePullSecret`. If it's publicly pullable,
you don't need this. But we do hope you are able to practice open science and share your containers!

```yaml
    # Name of an already created ImagePullSecret for the image specified above
    imagePullSecret: flux-image-secret
```

#### workingDir

The container likely has a set working directory, and if you are running the RESTful API service (meaning
you start without a command, as shown above) this will likely be the application folder. If you are launching
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
However, if you set this to true for *two* container (not allowed currently) you will get an error message.

```yaml
    # This defaults to true - this is the container we want to run flux in. This means
    # that if you have more than one container, set the non-flux runners to false.
    # For one container, you can leave this unset for the default. This will be
    # validated in case you make a mistake :)
    runFlux: true
```

#### environment

If you have environment variables to add, you can use an environment section with key value pairs:

```yaml
environment:
   RABBITMQ_DEFAULT_USER: aha
   RABBITMQ_DEFAULT_PASS: aharabbit
```

#### ports

The same goes for ports! Since we are implementing fairly simple use cases, for now ports
are provided as a single list of numbers, and the ideas is that the containerPort will
be assigned this number (and you can forward to your host however you like):

```yaml
ports:
  - 15672
  - 5671
  - 5672
```

#### fluxOptionFlags

Often when you run flux, you need to provide an option flag. E.g.,:

```bash
$ flux mini submit -ompi=openmpi@5
```

While these can be provided in the user interface of the Flux RESTFul API,
depending on your container image you might want to set some flags as default.
You can do this by setting this particular config parameter, and you should
set the flags just as you would to the command, starting with `-o`:

```yaml
	# optional - if needed, default option flags for the server (e.g., -ompi=openmpi@5)
	fluxOptionFlags: "-ompi=openmpi@5" 
```

Note that if you run with the user interface, setting a flag in the interface
that is defined for the server will override it here. These options are
currently defined for your entire cluster and cannot be provided to specific containers.
Also remember that your base container can equally provide these flags (and you
could equally override them, but if they are set and you don't define them here
they should not be touched).

#### fluxLogLevel

The log level to provide to flux, given that test mode is not on.

```yaml
	fluxLogLevel: 7 
```


#### preCommand

It might be that you want some custom logic at the beginning of your script.
E.g., perhaps you need to source an environment of interest! To support this we allow
for a string (multiple lines possible) of custom logic to do that. Remember
that since this is written into a flux runner wait.sh, this will only be
used for a Flux runner script. If you need custom logic in a service container
that is not a flux runner, you should write it into your own entrypoint.

```yaml
  # The pipe preserves line breaks
  preCommand: |
    ### Heading

    * Bullet
    * Points
```

#### diagnostics

Flux has a command that makes it easy to run diagnostics on a cluster, and we expose a boolean that makes it possible
to run that (instead of your job or starting the server). Since you might only want this for a specific container,
we provide this argument on the level of the container. To enable this, set this boolean to true. By default, it is false.

```yaml
  # Diagnostics runs flux commands for diagnostics, and a final sleep command
  # That makes it easy for you to shell into the pod to look around
  diagnostics: false
```



### fluxRestful

The "fluxRestful" section has a few parameters to dictate the installation of the
[Flux Restful API](https://github.com/flux-framework/flux-restful-api), which provides
a user interface to submit jobs.

#### branch

The branch parameter controls if you want to clone a custom branch (e.g., for testing).
It defaults to main.

```yaml
  fluxRestful:
    branch: feature-branch
```

#### port

The port parameter controls the port you want to run the FluxRestful server on,
within the cluster. Remember that you can always forward this to something else!
It defaults to 5000.

```yaml
  fluxRestful:
    port: 5000
```