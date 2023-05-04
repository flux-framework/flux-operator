# MiniCluster

> The CRD "Custom Resource Definition" defines a MiniCluster

A CRD is a yaml file that you can apply to your cluster (with the Flux Operator
installed) to ask for a MiniCluster. Kubernetes has these [custom resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
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

The size will always be the minimum number of pods that the Flux broker is required to see online
in order to start (meaning for the time being, all of them). If you've set a `maxSize` or you want to
scale smaller, you can re-apply the CRD to scale down. Flux will see the nodes as going offline, and 
of course you will want to be careful about the state of your cluster when you do this. If you scale
down, you cannot go below 1 node, and if you scale up, you cannot exceed the maximum of the `size` or
`maxsize`.

### maxSize

The `maxSize` variable is typically used when you want elasticity. E.g., it is the largest size of cluster that
you would be able to scale to. This works by way of registering this many workers (fully qualified domain names)
with the broker.toml. If you don't set this value, the maxsize will be set to the size.

```yaml
  # Number of pods to allow the MiniCluster to scale to
  maxSize: 10
```

The `maxSize` must always be greater than the size, if set. 

### tasks

The `tasks` variable under the spec is the number of tasks that each pod in the MiniCluster should be given.

```yaml
  tasks: 4
```

This value defaults to 1.

### interactive

Interactive mode means that the Flux broker is started without a command, and this would
allow you to shell into your cluster, connect to the broker, and interact with the Flux install.

```yaml
  interactive: true
```

This would be equivalent to giving a start command of `sleep infinity` however on exit
(e.g., if there is a flux shutdown from within the Flux instance) the sleep command would
not exit with a failed code.

### launcher

If you are using an executor that launches Flux Jobs (e.g., workflow managers such as Snakemake and Nextflow do!)
then you can set launcher to true.

```yaml
  launcher: true
```

The main difference is that for a launcher, we don't wrap it in a flux submit (as we would do with a job command).

### jobLabels

To add custom labels for your job, add a set of key value pairs (strings) to a "jobLabels" section:

```yaml
  jobLabels:
    job-attribute-a: dinosaur-a
    job-attribute-b: dinosaur-b
```


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
    storageClass: hostpath
```

The "storageClass" above (which you can leave out) defaults to hostpath, and should be the storage class that your cluster provides.
The Operator createst the "hostpath" volume claim. This currently is always created as a host path volume claim in MiniKube,
and likely in the future will have different logic if it varies from that.

#### labels

You can add labels:


```yaml
volumes:
  myvolume:
    path: /full/path/to/volume
    storageClass: hostpath
    labels:
      type: "local"
```

#### delete

By default, we will cleanup the persistent volume. To not do this (e.g., for a more permanent mount) set delete
to false:

```yaml
volumes:
  myvolume:
    path: /full/path/to/volume
    storageClass: csi-gcs
    delete: false
```

#### driver

If you are using anything aside from hostpath, you'll need a reference to a storage driver (usually a plugin)
you've installed separately. This can also be referenced as a provisioner. Here is an example for Google Cloud storage:

```yaml
volumes:
  myvolume:
    path: /full/path/to/volume
    storageClass: csi-gcs
    driver: gcs.csi.ofek.dev
```

#### volumeHandle

If your volume handle differs from your storage class name, you can define it:

```yaml
volumes:
  myvolume:
    path: /full/path/to/volume
    storageClass: csi-gcs
    driver: gcs.csi.ofek.dev
    volumeName: manualbucket/path
```

#### attributes

If your volume has attributes, you can add them too:

```yaml
volumes:
  myvolume:
    attributes:
    mounter: geesefs
    capacity: 25Gi
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
    storageClass: csi-gcs
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
    storageClass: csi-gcs
    secret: "csi-gcs-secret"
```

The secret (for now) should be in the default namespace.

#### annotations

To add annotations for the volume use "annotations"

```yaml
volumes:
  myvolume:
    annotations:
      provider.svc/attribute: value
```

#### claimAnnotations

To set annotations for the claim:

```yaml
volumes:
  myvolume:
    claimAnnotations:
      gcs.csi.ofek.dev/location: us-central1
      gcs.csi.ofek.dev/project-id: my-project
      gcs.csi.ofek.dev/bucket: flux-operator-storage
```

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

#### strict

By default, we run in bash strict mode, meaning that an error in a worker entrypoint script
will cause it to exit with a non-zero exit code. However, if you want to debug (and pass over the issue)
you can set this to false (it defaults to true):

```yaml
logging:
  strict: false
```

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

#### zeromq

ZeroMQ has logging explicitly for it, and you can enable it for the operator as follows:

```yaml
logging:
  zeromq: true
```

As an example, when the cluster workers are connecting to the broker, they will retry every
so often until the headless service is ready. That looks like this:

```bash
broker.debug[1]: parent sockevent tcp://lammps-0.flux-service.flux-operator.svc.cluster.local:8050 connect delayed
broker.debug[1]: parent sockevent tcp://lammps-0.flux-service.flux-operator.svc.cluster.local:8050 closed
broker.debug[1]: parent sockevent tcp://lammps-0.flux-service.flux-operator.svc.cluster.local:8050 connect retried
broker.debug[1]: parent sockevent tcp://lammps-0.flux-service.flux-operator.svc.cluster.local:8050 connect delayed
broker.debug[1]: parent sockevent tcp://lammps-0.flux-service.flux-operator.svc.cluster.local:8050 connected
broker.debug[1]: parent sockevent tcp://lammps-0.flux-service.flux-operator.svc.cluster.local:8050 handshake succeeded
```

This parameter is currently unset in Flux (so the wait can be slow) and we will have an update to Flux
soon to make the check more frequent. If you find that your cluster creation times are slower than
you expected, this is the probably cause, and it will be resolved with this update to Flux.


### archive

If you want to save state between MiniClusters, you can set an archive path for
the MiniCluster to load and save to. Given that the path exists, in the entrypoint
script it will be loaded via `flux system reload`. At the end, a pre stop hook
will then do another `flux dump` to that same path.

```yaml
archive:
  path: /state/archive.tar.gz
```

This obviously requires that you have a persistent volume to save to that subsequent MiniClusters
can access! This also assumes we are OK updating the archive state (and don't want to save the original). This can
be adjusted if needed.


### users

If you add a listing of users, minimally you need to provide a name for each one:

```yaml
users:
  - name: peenut
  - name: squidward
  - name: avocadosaurus
```

The users will be created and added to the Flux Accounting database. If you don't provide passwords,
they will be generated randomly (and you will need to retrieve them from the operator logs).
You can also define them manually:

```yaml
users:
  - name: peenut
    password: butter
  - name: squidward
    password: underdac
  - name: avocadosaurus
    password: eathings
```

The passwords (if provided) are validated to be 8 or fewer characters.
Note that although we don't validate this in the job, multi-user mode only makes sense to
provide alongside a custom resource definition without a command, meaning you submit
directly to the Flux Restful API server.```

### pod

Variables and attributes for each pod in the Indexed job.

#### labels

To add custom labels for your pods (in the indexed job), add a set of key value pairs (strings) to a "labels" section:

```yaml
pod:
  labels:
    pod-attribute-a: dinosaur-a
    pod-attribute-b: dinosaur-b
```

Note that the "namespace" variable is controlled by the operator here, and would be over-ridden if you defined it here.

#### annotations

The same is true for annotations! Just add annotations to a pod like so:


```yaml
pod:
  annotations:
    pod-annotation-a: dinosaur-a
    pod-annotation-b: dinosaur-b
```

#### resources

Resource lists for a pod go under [Overhead](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-overhead/). Known keys include "memory" and "cpu" (should be provided in some
string format that can be parsed) and all others are considered some kind of quantity request.

```yaml
pod:
  resources:
    memory: 500M
    cpu: 4
```

#### serviceAccountName

To give a service account name to your pods, simply do:

```yaml
pod:
  serviceAccountName: my-service-account
```

#### nodeSelector

A node selector is a set of key value pairs that helps to schedule pods to the right nodes! You can
add nodeSelector attributes to your pod as follows:

```yaml
pod:
  nodeSelector:
    iam.gke.io/gke-metadata-server-enabled: "true"
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

##### init

Init happens before everything - at the very beginning of the entrypoint. If you want to customize
the `PATH`, `PYTHONPATH`, or `LD_LIBRARY_PATH` handed to asFlux you can do that here.

```yaml
containers:
  - image: my-flux-image
    ...
    commands:
      init: export LD_LIBRARY_PATH=/opt/conda/lib
```


##### pre

The "pre" command is akin to `preCommand` but only run for the Flux workers and broker. It is run after a few
early environment variables are set (e.g., asFlux). Here is an example:

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

##### post

The "post" command is run for in the entrypoint after everything finishes up.

```yaml
containers:
  - image: my-flux-image
    ...
    commands:
      post: echo "Finishing up..."
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

##### prefix

This is a "wrapper" to any of flux submit, broker, or start. It typically is needed if you need to wrap
the initial command with something else. As an example, to get a storage driver working in the context
of a command, you might need to prefix the executable (see [fusion storage](../deployment/google/fusion.md)).

```yaml
containers:
  - image: my-flux-image
    ...
    commands:
      prefix: fusion
```

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

#### batch

If you are submitting many jobs, you are better off providing them to flux at once as a batch submission.
This way, we won't stress any Kubernetes APIs to submit multiple. To do this, you can define a command as before,
but then set batch to true:

```yaml
containers:
  - image: ghcr.io/flux-framework/flux-restful-api:latest

    # Indicate this should be a batch job
    batch: true

    # This command, as a batch command, will be written to a script and given to flux batch
    command: |
      echo hello world 1
      echo hello world 2
      echo hello world 3
      echo hello world 4
      echo hello world 5
      echo hello world 6
```

By default, output will be written to "/tmp/fluxout" for each of .out and .err files, and the
jobs are numbered by the order you provide above. To change this path:

```yaml
containers:
  - image: ghcr.io/flux-framework/flux-restful-api:latest

    # Indicate this should be a batch job
    batch: true
    logs: /tmp/another-out

    # This command, as a batch command, will be written to a script and given to flux batch
    command: |
      echo hello world 1
      echo hello world 2
      echo hello world 3
      echo hello world 4
      echo hello world 5
      echo hello world 6
```

Note that the output is recommended to be a shared volume so all pods can write to it.
If you can't use the filesystem for saving output, it's recommended to have some other
service used in your jobs to send output.

#### batchRaw

By default, the commands you provide to batch will be wrapped in flux submit, and with `--output` and `--flags waitable` added.
This works for a set of commands that are intended to be launched as such, but if you want custom logic in your script (such
as using flux exec or flux filemap) you can set batchRaw to true, and then provide the full flux directives in your
minicluster.yaml. As an example, here is using [flux filemap](https://flux-framework.readthedocs.io/projects/flux-core/en/latest/man1/flux-filemap.html) 
to copy data from the broker to all nodes in a batch job, and run the job.

```yaml
command: |
  flux filemap map -C /data mpi.sif
  flux exec -x 0 -r all flux filemap get -C /data
  flux submit singularity exec /data/mpi.sif /opt/mpitest
  flux exec -x 0 -r all rm -rf /data
  flux queue idle
  flux filemap unmap
```

See our [staging tutorial](../tutorials/staging.md) for more details on how this works! You could use it for a Singularity container (that needs to be seen by 
all nodes) or for data.

#### diagnostics

Flux has a command that makes it easy to run diagnostics on a cluster, and we expose a boolean that makes it possible
to run that (instead of your job or starting the server). Since you might only want this for a specific container,
we provide this argument on the level of the container. To enable this, set this boolean to true. By default, it is false.

```yaml
  # Diagnostics runs flux commands for diagnostics, and a final sleep command
  # That makes it easy for you to shell into the pod to look around
  diagnostics: false
```

#### lifeCycle

You can define postStartExec or preStopExec hooks.

```yaml
lifeCycle:
  postStartExec: ...
  preStopExec: ...
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
If you want to change the readOnly status to true:

```yaml
volumes:
  myvolume:
    path: /data
    readOnly: true
```

### existingVolumes

Existing volumes come down (typically) to a persistent volume claim (PVC) and persistent volume (PV)
that you've already created and want to give to the operator. As an example, the IBM plugin we use
to setup takes this approach, and then we define the existing volume on the level of the container:


```yaml
# This is a list because a pod can support multiple containers
containers:

    # This image has snakemake installed, and although it has data, we will
    # provide it as a volume to the container to demonstrate that (and share it)
  - image: ghcr.io/rse-ops/atacseq:app-latest

    # This is an existing PVC (and associated PV) we created before the MiniCluster
    existingVolumes:
      data:
        path: /workflow
        claimName: data 
```

The above would add a claim named "data" to the container it is defined alongside. Note that the names
define uniqueness, so if you use a claim in two places with the same name "data" it should also
use the same path "/workflow." If this doesn't work for your use case, please [let us know](https://github.com/flux-framework/flux-operator/issues).

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

#### secretKey

We use a secretKey to encode all payloads to the server. If you don't specify one,
the Flux Operator will make one for you! If you intend to communicate with your
MiniCluster outside of this context, you can either grab this from the logs
or define your own as follows:

```yaml
  fluxRestful:
    secretKey: notsosecrethoo
```
