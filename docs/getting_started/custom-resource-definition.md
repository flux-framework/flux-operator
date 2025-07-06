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
apiVersion: flux-framework.org/v1alpha2
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:
  ...
```

### Spec

Under the spec, there are several variables to define. Descriptions are included below, and we
recommend that you look at [config/samples](https://github.com/flux-framework/flux-operator/tree/main/examples)
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

Note that if you don't define tasks, you'll by default get a submit/batch or start command that defaults to 1 task, assuming testing, e.g., `-n 1`. If you define more tasks than nodes (size), you'll get `-N <size> and -n <tasks>`. If you only want to get `-N <size>` then sets tasks explicitly to 0.

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

This value defaults to 1. Note that if you don't define tasks, you'll by default get a submit/batch or start command that defaults to 1, assuming testing, e.g., `-n 1`. If you define more tasks than nodes (size), you'll get `-N <size> and -n <tasks>`. If you only want to get `-N <size>` then sets tasks explicitly to 0.

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

Note that by default, each pod will be labeled with a label `job-index` that corresponds to the
particular pod index. E.g., the lead broker would have `job-index=0` and this could be used as a service
selector.

### deadline

This is the maximum running time for your job. If you leave unset, it is essentially infinite.
If a Job is suspended (at creation or through an update), this timer will effectively be stopped and reset when the Job is resumed again.

```yaml
  # Deadline in seconds, if not set there is no deadline
  deadlineSeconds: 100
```

### network

The network section exposes networking options for the Flux MiniCluster.

#### disableAffinity

By default, the Flux Operator uses [Affinity](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity) and AntiAffinity (specifically `topology.kubernetes.io/zone` and `kubernetes.io/hostname` to ensure
that one pod is mapped per node. However, advanced use cases might warrant removing these. You can disable them as follows:

```yaml
network:
  disableAffinity: true
```

We put this under network due to the second rule that is about hostname, and (abstractly) you can imagine we are talking in the scope of at what level to assign a single worker associated with a hostname.
If you need finer-tuned control than disabling entirely, please open an issue to let us know.


#### headlessName

Change the default headless service name (defaults to `flux-service`).

```yaml
network:
  headlessName: my-network
```

### flux

The operator works to add Flux to your application container dynamically by way of using a provisioner container -
one that is run as a sidecar alongside your container, and then the view is copied over and flux run as your
active user. Settings under the Flux directive typically refer to flux options, e.g., for the broker or similar.

#### arch

If you are using an arm based container, ensure to add the architecture flag to designate that.

```yaml
flux:
   arch: "arm"
```

Note that this doesn't edit the container, but rather the binaries installed for it (e.g., to wait for files).


#### container

You can customize the flux container, and most attributes that are available for a standard container are available here.
As an example:

```yaml
flux:
   container:
     image: ghcr.io/converged-computing/flux-view-rocky:tag-9
     pythonPath: /mnt/flux/view/lib/python3.11
```

When enabled, meaning that we use flux from a view within the container, these containers are expected to have Flux built into a spack view, and in a particular way, so if you want to tweak or contribute a new means it's recommended to look at the [build repository](https://github.com/converged-computing/flux-views). This means that (if desired) you can customize this container base. We provide the following bases of interest:

 - [ghcr.io/converged-computing/flux-view-rocky:tag-9](https://github.com/converged-computing/flux-views/pkgs/container/flux-view-rocky)
 - [ghcr.io/converged-computing/flux-view-rocky:tag-8](https://github.com/converged-computing/flux-views/pkgs/container/flux-view-rocky)
 - [ghcr.io/converged-computing/flux-view-ubuntu:tag-noble](https://github.com/converged-computing/flux-views/pkgs/container/flux-view-ubuntu)
 - [ghcr.io/converged-computing/flux-view-ubuntu:tag-jammy](https://github.com/converged-computing/flux-views/pkgs/container/flux-view-ubuntu)
 - [ghcr.io/converged-computing/flux-view-ubuntu:tag-focal](https://github.com/converged-computing/flux-views/pkgs/container/flux-view-ubuntu)


Note that we have [arm builds](https://github.com/converged-computing/flux-views/tree/main/arm) available for each of rocky and ubuntu as well.
If you don't want to use Flux from a view (and want to use the v1apha1 design of the Flux Operator that had the application alongside Flux) you can do that by way of disabling the flux view:

```yaml
flux:
   container:
     image: ubuntu:focal
     disable: true
```

In the above, the Flux Operator won't expect Flux to be installed in the container you specify. We require a container to be specified,
however, because we still use the init strategy to set up configuration files. Everything except for the flux resources `R` file
is generated in that step (the broker configs and paths for archives and the broker socket). Keep in mind that if you intend
to deploy a sidecar container alongside your application container, you will still have access to this shared location to connect
to the socket, however you will need to provide your own install of Flux (ideally to match the one your application uses) to connect
to it. For a simple example of running lammps with this setup (without the sidecar) see the [disable-view](https://github.com/flux-framework/flux-operator/blob/main/examples/tests/disable-view)
example.

Please let us know if you'd like a base not provided.  Finally, the flux container can also take a specification of resources, just like a regular flux MiniCluster container:

```yaml
flux:
   container:
      resources:
        requests:
          cpu: "40"
          memory: "200M"
        limits:
          cpu: "40"
          memory: "200M"         
```

Note that if you expect to be able to schedule more than one pod per node, you will need to create pods with [Guaranteed](https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/#create-a-pod-that-gets-assigned-a-qos-class-of-guaranteed) QoS, 
in addition to creating your cluster with a config that specifies the cpu manager policy to be [static](https://kubernetes.io/docs/tasks/administer-cluster/cpu-management-policies/#static-policy):


```yaml
kubeletConfig:
  cpuManagerPolicy: static
linuxConfig:
 sysctl:
   net.core.somaxconn: '2048'
   net.ipv4.tcp_rmem: '4096 87380 16777216'
   net.ipv4.tcp_wmem: '4096 16384 16777216'
```

And then your pod containers also both need to have memory and cpu defined.  In summary:

1. Ensure cpuManagerPolicy is static
2. Create all pod containers (including the init container) in the MiniCluster to have a cpu and memory definition.

### completeWorkers

By default, when a follower broker is killed it is attempted to restart. While we could use [JobBackoffPerIndex](https://kubernetes.io/blog/2023/08/21/kubernetes-1-28-jobapi-update/#backoff-limit-per-index) to prevent it from restarting under
any conditions, this currently requires a feature gate (Kubernetes 1.28) so we are opting for a more simple approach. You can set `completeWorkers` to true, in which case when a lead broker is killed, it will Complete and not recreate.

```yaml
spec:
  flux:
    completeWorkers: true
```

This can be useful for cases of autoscaling in the down direction when you need to drain a node, and then delete the pod.

#### topology

By default, Flux will have a flat topology with one lead broker (rank 0) and some number of children. You can customize this with the `topology` field:

```yaml
flux:
  topology: kary:2
```

For example, you might chooes `kary:1` (or another value) or `binomial`. You can then use `flux overlay status` after connecting to your cluster to see it.

#### submitCommand

If you need to use a container with a different version of flux (or an older one)
you might want to edit the submit command. You can do that as follows:

```yaml
flux:
   submitCommand: "flux mini submit"
```

#### optionFlags

Often when you run flux, you need to provide an option flag. E.g.,:

```bash
$ flux submit -ompi=openmpi@5
```

While these can be provided in the user interface of the Flux RESTFul API,
depending on your container image you might want to set some flags as default.
You can do this by setting this particular config parameter, and you should
set the flags just as you would to the command, starting with `-o`:

```yaml
flux:
    # optional - if needed, default option flags for the server (e.g., -ompi=openmpi@5)
    optionFlags: "-ompi=openmpi@5"
```

Note that if you run with the user interface, setting a flag in the interface
that is defined for the server will override it here. These options are
currently defined for your entire cluster and cannot be provided to specific containers.
Also remember that your base container can equally provide these flags (and you
could equally override them, but if they are set and you don't define them here
they should not be touched).

#### logLevel

The log level to provide to flux, given that test mode is not on.

```yaml
flux:
  logLevel: 7
```

#### wrap

If you want to use the flux `--wrap` directive, which will add an additional `--wrap`
with some listing of arguments, you can specify it as follows:

```yaml
flux:
  wrap: "strace,-e,network,-tt"
```

In the above, we would add `--wrap=strace,-e,network,-tt` to flux start commands.

#### scheduler

Under flux->scheduler you can define attributes for the scheduler. We currently allow
setting a custom queue policy. The default (if unset) looks like this for first come first serve:

```yaml
flux:
  scheduler:
    queuePolicy: fcfs
```

To change to a policy that supports backfilling, you can change this to "easy":

```yaml
flux:
  scheduler:
    queuePolicy: easy
```

And the broker.toml config will update appropriately:

```toml
[sched-fluxion-qmanager]
queue-policy = "easy"
```

By default, we use the simple sched algorithm. To load fluxion, you can connect to the lead broker and load the correct modules.

```bash
kubectl exec -it flux-sample-0-7mqg5 bash
. /mnt/flux/flux-view.sh 
flux proxy $fluxsocket bash
bash /mnt/flux/load-fluxion.sh 
```

Verify it is loaded - you should see two modules for fluxion, qmanager and resource.

```bash
$ flux module list
Module                   Idle  S Service
job-info                   21  R 
sched-fluxion-qmanager      5  R sched
sched-fluxion-resource      5  R 
job-exec                   21  R 
connector-local             0  R 
kvs-watch                  21  R 
heartbeat                   1  R 
content                    20  R 
kvs                        20  R 
content-sqlite             20  R content-backing
job-list                   21  R 
job-ingest                 20  R 
job-archive                21  R 
resource                    5  R 
barrier                    21  R 
job-manager                 5  R 
cron                       21  R 
```

You can learn more about queues [here](https://flux-framework.readthedocs.io/en/latest/guides/admin-guide.html?h=system#adding-queues). Please [open an issue](https://github.com/flux-framework/flux-operator/issues) if you want support for something that you don't see. Also note that you can set an entire [broker config](#broker-config) for more detailed customization.

#### minimalService

By default, the Flux MiniCluster will be created with a headless service across the cluster,
meaning that all pods can ping one another via a fully qualified hostname. As an example,
the 0 index (the lead broker) of an indexed job will be available at:

```console
flux-sample-0.flux-service.flux-operator.svc.cluster.local
```

Where "flux-sample" is the name of the job. Index 1 would be at:

```console
flux-sample-1.flux-service.flux-operator.svc.cluster.local
```

However, it's the case that only the lead broker (index 0) needs to be reachable
by the others. If you set `minimalService` to true, this will be honored, so
the networking setup will be more minimal.

```yaml
flux:
  minimalService: true
```

The drawback is that you cannot ping the other nodes by hostname.

#### brokerConfig

If you want to manually specify the broker config in entirety, you can define
it as a full string under flux->brokerConfig. This is useful for experiments
or development work.

```yaml
flux:
  brokerConfig: |
     [exec]
     ...
```


#### curveCert

The same goes for the curve certificate! If you are bursting and want your
new cluster to talk to the cluster it's bursted from, you can share this certificate.

```yaml
flux:
  curveCert: |
     [exec]
     ...
```

#### connectTimeout

For Flux versions 0.50.0 and later, you can customize the zeromq timeout. This is done
via an extra broker option provided on the command line:

```bash
-Stbon.connect_timeout=5s
```

But note that (if you are interested) you could provide it in the broker.toml as well:

```toml
[tbon]
connect_timeout = "10s"
```

Note the above is a string, so the default for this value (if you leave it unset is):

```yaml
flux:
  connectTimeout: 5s
```

To disable adding this parameter, period, set it to an empty string:

```yaml
flux:
  connextTimeout: ""
```

A timeout of 0s will not be honored, and will default to the empty string (this is done internal in Flux).

> Why is this important?

When the Flux network is coming online (the headless service for the pods) zeromq will immediately try
to connect, and then (by default if we don't set it) start an exponential backoff of retry. Since we've
changed the design to not have a certificate generation pod, we've found that the network is not reliably
up by the time the broker starts, which, combined with this backoff, can lead to slow startup times.
Part of our strategy to help with that is to set the connect timeout (as shown above) to a set value.
In our experiments using [this script](https://github.com/flux-framework/flux-operator/blob/main/sdk/python/v1alpha1/examples/time-minicluster-lammps.py)
with a LAMMPS image, we found an optimal connect timeout to be lower between 1s and 5s ([see this work](https://github.com/converged-computing/operator-experiments/tree/main/google/service-timing)). 
Note that this could vary depending on your setup, and we recommend you run the script with your MiniCluster setup in the Python script to test.
As another option, you can add an arbitrary service container, which seems to "warmup" the network
and skip most of this delay. This service combined with the connect timeout seems to be the optimal way
to minimize the startup time - for LAMMPS, this means going back to the original 11-12 seconds we saw
with our original Flux Operator experiments for Kubecon 2023.
More details are available in [this post](https://github.com/converged-computing/operator-experiments/tree/main/google/service-timing).
Although we have fixed the zeromq timeout bug, there still seems to be some underlying issue in Kubernetes.

Note that we think this higher level networking bug is still an issue, and are going to be creating a dummy case
to reproduce it and share with the Kubernetes networking team.

### bursting

The bursting spec defines when you expect to burst to an external cluster. Since the initial Flux cluster needs
to know all future hosts in advance, we ask that you provide sizes. More specifically:

- A bursted cluster should have a `leadBroker` defined - the ip address or hostname, port, and job name of the initial MiniCluster
- The main "root" cluster should not have a `leadBroker` defined, but needs to know about the intended bursted clusters and sizes.

An example for the root cluster might look like the following:

```yaml
  flux:
    # Declare that this cluster will allow for a bursted cluster
    # It would automatically be named burst-0, but we explicitly set
    # for clarity. The leadBroker is left out because this IS it.
    bursting:
      clusters:
        - size: 4
          name: burst-0
```

And then the cluster launched from this cluster "bursted to" would define:


```yaml
  flux:
    leadBroker:
      # This is the name of the first minicluster.yaml spec
      name: flux-sample
      # In a cloud environment this would be a NodePort
      address: 24.123.50.123
      port: 30093
    
    bursting:
      clusters:
        - size: 4
          name: burst-0
```

Using the above, both the main and bursted to cluster will have almost the same spec for their flux resources
and broker.toml (config). The main difference will be that the bursted cluster knows about the first one via
it's ip address or hostname, and not, for example `flux-sample-0`. Also note that when bursting, you don't
explicitly give a command to the bursted cluster - the jobs are launched on the main cluster and sent
to these external resources when they come up and are available (and needed). 

Finally, for advanced bursting cases where the pattern of hostnames does not match the convention
deployed by the Flux Operator, we allow the CRD to define a custom list. As an example, here is how
we might burst to compute engine:

```yaml
  flux:
    leadBroker:
      # This is the name of the first minicluster.yaml spec
      name: flux-sample
      # In a cloud environment this would be a NodePort
      address: 24.123.50.123
      port: 30093
      hostlist: "flux-sample-[0-3],gffw-compute-a-[001-003]" 
```

In the above case, the clusters are not used. The bursting plugin you use will determine
how the hostnames and address are provided to the remote (second) cluster.

For full examples, see [the bursting](https://github.com/flux-framework/flux-operator/tree/main/examples/experimental/bursting) 
examples directory.


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

#### schedulerName

Request that the pod be scheduled by a custom scheduler.

```yaml
pod:
  schedulerName: fluence
```

For fluence you can see the custom scheduler plugin at [flux-framework/flux-k8s](https://github.com/flux-framework/flux-k8s).

#### serviceAccountName

To give a service account name to your pods, simply do:

```yaml
pod:
  serviceAccountName: my-service-account
```

#### restartPolicy

To customize the restartPolicy for the pod:

```yaml
pod:
  restartPolicy: Never
```

#### runtimeClassName

To add a runtime class name:

```yaml
pod:
  runtimeClassName: nvidia
```

#### automountServiceAccountToken

If you want to automatically mount a service account token:

```yaml
pod:
  automountServiceAccountToken: true
```


#### nodeSelector

A node selector is a set of key value pairs that helps to schedule pods to the right nodes! You can
add nodeSelector attributes to your pod as follows:

```yaml
pod:
  nodeSelector:
    iam.gke.io/gke-metadata-server-enabled: "true"
```

#### hostIPC

A boolean to use the hostIPC namespace. We have not tested the used cases for this yet.

```yaml
pod:
  hostIPC: true
```


#### hostPID

The same, but to add `hostPID` to the pod.

```yaml
pod:
  hostPID: true
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

You do not need to provide a container base that has Flux, but you must make sure your view (with a particular operator system) that will add Flux matches your container. We don't require you to start as root, but if you
have a container with a non-root user, the user needs to have sudo available (to act as root).
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
2. Not providing a command means that we will create a persistent MiniCluster running in interactive mode.

```yaml
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

The container likely has a set working directory, and if you are running an interactive cluster (meaning
you start without a command, as shown above) this will likely be the application folder. If you are launching
a job directly with flux start and require a particular working directory, set it here!

```yaml
    # You can set the working directory if your container WORKDIR is not correct.
    workingDir: /home/flux/examples/reaxff/HNS
```

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

A special "commands" section is available for commands that you want to run in the broker and workers containers. 

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

##### script

If you want to write a custom entrypoint script (that will be chmod +x and provided to flux submit) you
can do so with script. For example:

```yaml
containers:
  - image: my-flux-image
    ...
    commands:
      script: |
        #!/bin/bash
        echo "This is my custom script"
        echo "This is rank \${FLUX_TASK_RANK}
```

Note that the environment variable `$` is escaped.
For a container that is running flux you are not allowed to provide a script and command, and validation will check this.

##### workerPre and brokerPre

This is akin to pre, but only for the broker OR the workers.

```
containers:
  - image: my-flux-image
    ...
    commands:
      preWorker: echo hello I am a worker
      preBroker: echo hello I am the lead broker 
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
  - image: rockylinux:9

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
  - image: rockylinux:9

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

As of version 0.2.0, the Flux Operator controlling its own volumes has been removed, and the concept of an "existingVolume" is now just a volume.
We did this because volumes are complex and it was challenging to support every use case. Thus, our strategy is to allow the user to create
the volumes and persistent volume claims that the MiniCluster needs, and simply tell it about them. A volume (that must exist) can be:

 - a hostpath (good for local development)
 - an empty directory (and of a particular custom type, such as Memory)
 - a persistent volume claim (PVC) and persistent volume (PV) that you've created
 - a secret that you've created
 - a config map that you've created

and for all of the above, you want to provide it to the operator, which will work either for a worker
pod (in the indexed job) or a service. 

#### hostpath example

You might start by creating a pv and a pvc:

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: data
spec:
  storageClassName: manual
  capacity:
    storage: 5Gi
  accessModes:
    - ReadWriteMany
  hostPath:
    path: /data
---
apiVersion: v1 
kind: PersistentVolumeClaim
metadata: 
  name: data
spec: 
  accessModes: 
    - ReadWriteMany
  resources: 
    requests: 
      storage: 1Gi
```
And then to add this PVC to your MiniCluster:

```yaml
apiVersion: flux-framework.org/v1alpha2
kind: MiniCluster
metadata:
  name: flux-sample
spec:
  size: 2
  containers:
    - image: rockylinux:9
      command: ls /data
      volumes:
        data:
          path: /data
          hostPath: /data
```

An example is provided in the [volumes test](https://github.com/flux-framework/flux-operator/tree/main/examples/tests/volumes).

#### emptyDir example

A standard empty directory might look like this:

```yaml
apiVersion: flux-framework.org/v1alpha2
kind: MiniCluster
metadata:
  name: flux-sample
spec:
  size: 2
  containers:
    - image: rockylinux:9
      command: df -h /dev/shm
      volumes:
        # must be lowercase!
        my-empty-dir:
          emptyDir: true
```

And one for shared memory (to inherit the host) like this:


```yaml
apiVersion: flux-framework.org/v1alpha2
kind: MiniCluster
metadata:
  name: flux-sample
spec:
  size: 2
  containers:
    - image: rockylinux:9
      command: ls /data
      volumes:
        # must be lowercase!
        my-empty-dir:
          emptyDir: true
          emptyDirMedium: "memory"
```

The default binds to the path `/dev/shm` and is not customizable. This can be changed if needed. When you have the "memory" medium added,
you should see all the shared memory from the host, which is [calculated here](https://github.com/kubernetes/kubernetes/blob/e6616033cb844516b1e91b3ec7cd30f8c5d1ea50/pkg/volume/emptydir/empty_dir.go#L148-L157). In addition, you can set a sizeLimit:

```yaml
        # must be lowercase!
        my-empty-dir:
          emptyDir: true
          emptyDirMedium: "memory"
          emptyDirSizeLimit: "64Gi"
```

As an example, here is output from a local run with kind when shared memory is added:

```console
$ kubectl logs flux-sample-0-smflk 
Defaulted container "flux-sample" out of: flux-sample, flux-view (init)
Filesystem      Size  Used Avail Use% Mounted on
tmpfs            32G     0   32G   0% /dev/shm
```

And here is the same MiniCluster with the volume removed (64M is the default):

```console
$ kubectl logs flux-sample-0-4bwjf -f
Defaulted container "flux-sample" out of: flux-sample, flux-view (init)
Filesystem      Size  Used Avail Use% Mounted on
shm              64M     0   64M   0% /dev/shm
```

#### persistent volume claim example

As an example, the IBM plugin we use to setup takes this approach, and then we define the existing volume on the level of the container:

```yaml
# This is a list because a pod can support multiple containers
containers:

    # This image has snakemake installed, and although it has data, we will
    # provide it as a volume to the container to demonstrate that (and share it)
  - image: ghcr.io/rse-ops/atacseq:app-latest

    # This is an existing PVC (and associated PV) we created before the MiniCluster
    volumes:
      data:
        path: /workflow
        claimName: data 
```

The above would add a claim named "data" to the container it is defined alongside. Note that the names
define uniqueness, so if you use a claim in two places with the same name "data" it should also
use the same path "/workflow." If this doesn't work for your use case, please [let us know](https://github.com/flux-framework/flux-operator/issues).

#### config map example

Here is an example of providing a config map to a service container (that runs as a sidecar container alongside the MiniCluster)
that we want to add to the service container. In layman's terms, we are deploying vanilla nginx, but adding a configuration file
to `/etc/nginx/conf.d`

```yaml
# Add an nginx service with an existing config map
services:
  - image: nginx
    name: nginx
    volumes:
      nginx-conf:
        configMapName: nginx-conf 
        path: /etc/nginx/conf.d
        items:
          flux.conf: flux.conf
```

Your config map would be created separately (likely first, before the MiniCluster). Here is an example:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-conf
  namespace: flux-operator
data:
  flux.conf: |
    server {
        listen       80;
        server_name  localhost;
        location / {
          root   /usr/share/nginx/html;
          index  index.html index.htm;
        }        
    }
```

Note that the above block `volumes` is valid to be under the
MiniCluster->containers section too. Either MiniCluster containers OR service
containers can have existingVolumes of any type.

#### secret example

Here is an example of providing an existing secret (in the flux-operator namespace)
to the indexed job container:

```yaml
containers:
  - image: ghcr.io/rse-ops/singularity:tag-mamba
    workingDir: /data
    volumes:
      certs:
        path: /etc/certs
        secretName: certs
```

The above shows an existing secret named "certs" that we will mount into `/etc/certs`.
