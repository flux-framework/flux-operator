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

### tasks

The `tasks` variable under the spec is the number of tasks that each pod in the MiniCluster should be given. 

```yaml
  tasks: 4
```

This value defaults to 1.

### jobLabels

To add custom labels for your job, add a set of key value pairs (strings) to a "jobLabels" section:

```yaml
  jobLabels:
    job-attribute-a: dinosaur-a
    job-attribute-b: dinosaur-b
```

### podLabels

To add custom labels for your pods (in the indexed job), add a set of key value pairs (strings) to a "podLabels" section:

```yaml
  pobLabels:
    pod-attribute-a: dinosaur-a
    pod-attribute-b: dinosaur-b
```

Note that the "namespace" variable is controlled by the operator here, and would be over-ridden if you defined it here.


### deadline

This is the maximum running time for your job. If you leave unset, it is essentially infinite.
If a Job is suspended (at creation or through an update), this timer will effectively be stopped and reset when the Job is resumed again.

```yaml
  # Deadline in seconds, if not set there is no deadline
  deadlineSeconds: 100
```

### volumes

Volumes can be defined on the level of the MiniCluster that are then used by containers.

 - For MiniKube, these volumes are expected to be inside of the VM, e.g., accessed via `minikube ssh`
 - For an actual cluster, they should be on the node running the pod.

#### volume ids

For each volume under "volumes" we enforce a unique name by way of using key values - e.g., "myvolume"
in the example below can then be referenced for a container:

```yaml
volumes:
  myvolume:
    path: /full/path/to/volume
    class: hostpath
```

The "class" above (which you can leave out) defaults to hostpath, and should be the storage class that your cluster provides.
The Operator createst the "hostpath" volume claim. This currently is always created as a host path volume claim in MiniKube,
and likely in the future will have different logic if it varies from that. 

#### labels

You can add labels:


```yaml
volumes:
  myvolume:
    path: /full/path/to/volume
    class: hostpath
      labels:
        type: "local"
```

#### request storage size

By default, a capacity request is "5Gi", and we only do this because the field is required. However, keep in mind
for many some cloud storage interfaces there is no concept of a max. 
This is defined as a string to be parsed. To tweak that, meaning
that this container will request this amount of storage for the container (and here we show a different storageclass
for Google Cloud)

```yaml
volumes:
  myvolume:
    path: /full/path/to/volume
    capacity: 5Gi
    class: csi-gcs
```

Since storage classes are created separately (not by the operator) you should check with your storage
class to ensure resource limits work with your selection above.

#### secret

For a CSI (container storage interface) you usually need to provide a secret reference. For example,
for GKE we create a service account with the appropriate permissions, and then apply them as a secret named `csi-gks-secret`:

```yaml
volumes:
  myvolume:
    path: /full/path/to/volume
    capacity: 1Gi
    class: csi-gcs
    secret: "csi-gcs-secret"
```

The secret (for now) should be in the default namespace.

### cleanup

If you add any kind of persistent volume to your MiniCluster, it will likely need a cleanup after the fact
(after you bring the MiniCluster down). By default, the operator will perform this cleanup, checking if the
Job status is "Completed" and then removing the pods, Persistent Volume Claims, and Persistent Volumes.
If you want to disable this cleanup:

```yaml
  cleanup: false
```

If you are streaming the logs with `kubectl logs` the steam would stop when the broker pod is completed,
so typically you will get the logs as long as you are streaming when the job starts running.

### logging

We provide simple types of "logging" within the main script that is run for the job.
If you don't set any variables, you'll get the most verbosity with timing of the main
command. You can set any subset of these for a custom output. Note that these logging levels
are not specific to operator logs, but the indexed job running in the pod. 

#### quiet

Quiet mode turns off all verbose output (yes, the emojis too) so only the output of 
your job will be printed to the console. This way, you can retrieve the job lob
and then determine if the test was successful based on this output.

```yaml
logging:
  quiet: true
```

By default quiet is false. In early stages of the operator this was called `test`.

#### timed

Timed mode adds timing for the main Flux command and a few other interactions in the script.

```yaml
logging:
  timed: true
```

By default timed is set to `false` above, and this is because if you turn it on your Flux runner
container is required to have `time` installed. We target `/usr/bin/time` and not the `time`
wrapper because we want to set a format with `-f` (which won't be supported by the wrapper).
By default we ask for `-f E` which means:

> Elapsed real (wall clock) time used by the process, in [hours:]minutes:seconds.

Also note that `timed` and `quiet` can influence one another - e.g., if quiet is `true` and
there are some timed sections under a section that is no longer included when the job
is quiet, you will not see those times. Here is an example of timing a hello-world run:

```bash
hello world
FLUXTIME fluxsubmit wall time 0:04.73
```
All timed command lines are prefixed with `FLUXTIME` and the main submit will be `fluxsubmit`
and the worker pods flux start will have `fluxstart`.

#### debug

Debug mode adds verbosity to flux to see additional information about the job submission.

```yaml
logging:
  debug: true
```


### pod

Variables and attributes for each pod in the Indexed job.

#### resources

Resource lists for a pod go under [Overhead](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-overhead/). Known keys include "memory" and "cpu" (should be provided in some
string format that can be parsed) and all others are considered some kind of quantity request. 

```yaml
pod:
  resources:
    memory: 500M
    cpu: 4
```


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

#### cores

The number of cores to provide to the container as a variable. This does not actually allocate or control cores
for the container, but exposes the variable for your container template (e.g., for the Flux wait.sh script). 

```yaml
  cores: 4
```

This value when unset defaults to 1.


#### command

Providing (or not providing) a command is going to dictate the behavior of your MiniCluster!

1. Providing a custom command means the MiniCluster is ephemeral - it will run the command and clean up.
2. Not providing a command means that we will create a persistent MiniCluster running a RESTFul API service (and GUI) to submit jobs to.

```yaml
    # Don't set a command unless you want to forgo running the restful server to submit
    # commands to! E.g., instead of starting the server, it will just run your job command.
    command: lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite
```

#### resources

Resources can include limits and requests. Known keys include "memory" and "cpu" (should be provided in some
string format that can be parsed) and all others are considered some kind of quantity request. 

```yaml
resources:
  limits:
    memory: 500M
    cpu: 4
```
If you wanted to, for example, request a GPU, that might look like:

```yaml
resources:
  limits:
    gpu-vendor.example/example-gpu: 1
```

Or for a particulat type of networking fabric:

```yaml
resources:
  limits:
    vpc.amazonaws.com/efa: 1
```

Both limits and resources are flexible to accept a string or an integer value, and you'll get an error if you
provide something else. If you need something else, [let us know](https://github.com/flux-framework/flux-operator/issues).
If you are requesting GPU, [this documentation](https://kubernetes.io/docs/tasks/manage-gpus/scheduling-gpus/) is helpful.

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
for a string (multiple lines possible) of custom logic to do that. This
"global" preCommand will be run for both flux workers (including the broker)
and the certificate generation script.

```yaml
  # The pipe preserves line breaks
  preCommand: |
    ### Heading

    * Bullet
    * Points
```

As a good example use case, we use an `asFlux` prefix to run any particular flux command
as the flux user. This defaults to the following giving you have the default `runAsFluxUser` to true:

```bash
asFlux="sudo -u flux -E PYTHONPATH=$PYTHONPATH -E PATH=$PATH"
```

However, let's say you have a use case that warrants passing on a custom set of environment
variables. For example, when we want to use Flux with MPI + libfabric (EFA networking in AWS)
we want these extra variables:

```bash
asFlux="sudo -u flux -E PYTHONPATH=$PYTHONPATH -E PATH=$PATH -E FI_EFA_USE_DEVICE_RDMA=1 -E RDMAV_FORK_SAFE=1"
```

Thus, we would define this line in our `preCommand` section. Since this runs directly after the default asFlux is defined,
it will be over-ridden to use our variant. As a final example, for a snakemake workflow we are expected to write
assets to a home directory, so we need to customize the entrypoint for that.

```bash
# Ensure the cache targets our flux user home
asFlux="sudo -u flux -E PYTHONPATH=$PYTHONPATH -E PATH=$PATH -E HOME=/home/flux"
```

#### commands

A special "commands" section is available for commands that you want to run in the broker and workers containers,
but not during certificate generation. As an example, if you print extra output to the certificate generator,
it will mangle the certificate output. Instead, you could write debug statements in this section.

##### pre

The "pre" command is akin to `preCommand` but only run for the Flux workers and broker. Here is an example:

```yaml
containers:
  - image: my-flux-image
    ...
    commands:
      pre: |
        # Commands that might print to stdout/stderr to debug, etc.
        echo "I am running the pre-command"
        ls /workdir
```

##### runFluxAsRoot

For different storage interfaces (e.g., CSI means "Container Storage Interface") you might need to 
run flux as root (and not change permission of the mounted working directory) to be owned by the flux user. You
can set this flag to enable that:

```yaml
containers:
  - image: my-flux-image
    ...
    commands:
      runFluxAsRoot: true
```

This defaults to false, meaning we run everything as the Flux user, and you are encouraged to try to figure out
setting up your storage to be owned by that user.


#### fluxUser

If you need to change the uid or name for the flux user, you can define that here.

```yaml
containers:
  - image: my-flux-image
    ...
    fluxuser:
      # Defaults to 1000
      uid: 1002
      # Defaults to flux
      name: flux
```

Note that if the "flux" user already exists in your container, the uid will be discovered and you don't need 
to set this. These parameters are only if you want the flux user to be created with a different unique id.


#### diagnostics

Flux has a command that makes it easy to run diagnostics on a cluster, and we expose a boolean that makes it possible
to run that (instead of your job or starting the server). Since you might only want this for a specific container,
we provide this argument on the level of the container. To enable this, set this boolean to true. By default, it is false.

```yaml
  # Diagnostics runs flux commands for diagnostics, and a final sleep command
  # That makes it easy for you to shell into the pod to look around
  diagnostics: false
```

### volumes

Volumes that are defined on the level of the container must be defined at the top level of the MiniCluster.
As an example, here is how we tell the container to use the already defined volume "myvolume" to be mounted
in the container as "/data":

```yaml
volumes:
  myvolume:
    path: /data
```

The `myvolume` key must be defined in the MiniCluster set of volumes, and this is checked.
If you want to change the readonly status to true:

```yaml
volumes:
  myvolume:
    path: /data
    readonly: true
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