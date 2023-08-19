# Helper Clients

To provide advanced functionality for the Flux Operator, we provide a set of helper clients 
that perform specific tasks. This is currently a limited set (just one!) however if you
think of a case where you would find a helper useful, please
[let us know](https://github.com/flux-framework/flux-operator/issues). 

## Building Helpers

To build helpers from source, after cloning the repository:

```bash
$ make helpers
```
They will appear as "fluxoperator-*" in the bin:

```bash
$ tree bin/
bin/
├── fluxoperator-gen
└── fluxoperator-gen
```

Or if you pull and have the container, interact as follows:

```bash
$ docker run -it --entrypoint fluxoperator-gen  ghcr.io/flux-framework/flux-operator:test --help
```
```console
Usage of fluxoperator-gen:
  -f string
        YAML filename to read it
  -i string
        Custom list of includes (csjv) for cm, svc, job, volume
  -kubeconfig string
        Paths to a kubeconfig. Only required if out-of-cluster.
```

In the above, we change the entrypoint from the `/manager` to target our fluxoperator-gen that is on the path.

## fluxoperator-gen

For some cases of using the Flux operator where dynamism is involved (either in scaling or having custom volumes set up)
you typically need the entire operator. However, for cases where you are using the Flux Operator as more of a job submission
tool (e.g., akin to what we do in [Kueue](https://kueue.sigs.k8s.io/docs/tasks/run_flux_minicluster/)) you really only
need to generate the core assets for your MiniCluster, which are an indexed job, config maps, and a service. We
realized this, that in fact a simple operator use case is just a fancy way to produce complex configs, and provide
this extra helper, "gen" that does exactly that. Our use case was exactly that, where we only need the config maps
for a [metrics operator](https://github.com/converged-computing/metrics-operator-experiments/tree/main/flux-operator) experiment,
and we now provide the tool as a more general use case. A few notes:

 - We don't currently support sidecar services. The assumption is that you can generate them yourself! If you think we should add this support, [let us know](https://github.com/flux-framework/flux-operator/issues). 

### fluxoperator-gen usage

Let's walk through some examples We will assume you have the `fluxoperator-gen` built or provided by the container.
By default, provide it a file `-f` that has a MiniCluster spec in YAML to generate (to the screen) a dump of all the objects that the Flux Operator creates. Here is an example with our LAMMPS test file:

```bash
$ fluxoperator-gen -f ./examples/tests/lammps/minicluster.yaml 
```

<details>

<summary>fluxoperator-gen output</summary>

```console
----
apiVersion: v1
data:
  hostfile: "# Flux needs to know the path to the IMP executable\n[exec]\nimp = \"/usr/libexec/flux/flux-imp\"\n\n[access]\nallow-guest-user
    = true\nallow-root-owner = true\n\n# Point to resource definition generated with
    flux-R(1).\n[resource]\npath = \"/etc/flux/system/R\"\n\n[bootstrap]\ncurve_cert
    = \"/etc/curve/curve.cert\"\ndefault_port = 8050\ndefault_bind = \"tcp://eth0:%p\"\ndefault_connect
    = \"tcp://%h..flux-operator.svc.cluster.local:%p\"\nhosts = [\n\t{ host=\"flux-sample-[0--1]\"},\n]\n[archive]\ndbpath
    = \"/var/lib/flux/job-archive.sqlite\"\nperiod = \"1m\"\nbusytimeout = \"50s\"\n\n#
    Configure the flux-sched (fluxion) scheduler policies\n# The 'lonodex' match policy
    selects node-exclusive scheduling, and can be\n# commented out if jobs may share
    nodes.\n[sched-fluxion-qmanager]\nqueue-policy = \"fcfs\""
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: flux-sample-entrypoint
  namespace: flux-operator

----
apiVersion: v1
data:
  hostfile: "# Flux needs to know the path to the IMP executable\n[exec]\nimp = \"/usr/libexec/flux/flux-imp\"\n\n[access]\nallow-guest-user
    = true\nallow-root-owner = true\n\n# Point to resource definition generated with
    flux-R(1).\n[resource]\npath = \"/etc/flux/system/R\"\n\n[bootstrap]\ncurve_cert
    = \"/etc/curve/curve.cert\"\ndefault_port = 8050\ndefault_bind = \"tcp://eth0:%p\"\ndefault_connect
    = \"tcp://%h..flux-operator.svc.cluster.local:%p\"\nhosts = [\n\t{ host=\"flux-sample-[0--1]\"},\n]\n[archive]\ndbpath
    = \"/var/lib/flux/job-archive.sqlite\"\nperiod = \"1m\"\nbusytimeout = \"50s\"\n\n#
    Configure the flux-sched (fluxion) scheduler policies\n# The 'lonodex' match policy
    selects node-exclusive scheduling, and can be\n# commented out if jobs may share
    nodes.\n[sched-fluxion-qmanager]\nqueue-policy = \"fcfs\""
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: flux-sample-flux-config
  namespace: flux-operator

----
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  namespace: flux-operator
spec:
  clusterIP: None
  selector:
    job-name: flux-sample
status:
  loadBalancer: {}

----
apiVersion: v1
kind: Job
metadata:
  creationTimestamp: null
  name: flux-sample
  namespace: flux-operator
spec:
  activeDeadlineSeconds: 0
  backoffLimit: 100
  completionMode: Indexed
  completions: 4
  parallelism: 4
  template:
    metadata:
      creationTimestamp: null
      labels:
        app.kubernetes.io/name: flux-sample
        hpa-selector: flux-sample
        namespace: flux-operator
      name: flux-sample
      namespace: flux-operator
    spec:
      containers:
      - image: ghcr.io/rse-ops/lammps:flux-sched-focal
        imagePullPolicy: IfNotPresent
        lifecycle: {}
        name: ""
        resources: {}
        securityContext:
          capabilities: {}
          privileged: false
        stdin: true
        tty: true
        volumeMounts:
        - mountPath: /mnt/curve/
          name: flux-sample-curve-mount
          readOnly: true
        - mountPath: /etc/flux/config
          name: flux-sample-flux-config
          readOnly: true
        - mountPath: /flux_operator/
          name: flux-sample-entrypoint
          readOnly: true
        workingDir: /home/flux/examples/reaxff/HNS
      restartPolicy: OnFailure
      setHostnameAsFQDN: false
      shareProcessNamespace: false
      volumes:
      - configMap:
          items:
          - key: hostfile
            path: broker.toml
          name: flux-sample-flux-config
        name: flux-sample-flux-config
      - configMap:
          name: flux-sample-entrypoint
        name: flux-sample-entrypoint
      - configMap:
          name: flux-sample-curve-mount
        name: flux-sample-curve-mount
status: {}
```
</details>

You can also ask for specific includes (by default we include everything):

- **c** configmaps
- **j** job
- **s** service
- **v** volumes

As an example, to generate only config maps (a use case I have) I do:

```bash
# This says "include c for config maps"
$ fluxoperator-gen -i c -f ./examples/tests/lammps/minicluster.yaml 
```
```console
----
apiVersion: v1
data:
  hostfile: "# Flux needs to know the path to the IMP executable\n[exec]\nimp = \"/usr/libexec/flux/flux-imp\"\n\n[access]\nallow-guest-user
    = true\nallow-root-owner = true\n\n# Point to resource definition generated with
    flux-R(1).\n[resource]\npath = \"/etc/flux/system/R\"\n\n[bootstrap]\ncurve_cert
    = \"/etc/curve/curve.cert\"\ndefault_port = 8050\ndefault_bind = \"tcp://eth0:%p\"\ndefault_connect
    = \"tcp://%h..flux-operator.svc.cluster.local:%p\"\nhosts = [\n\t{ host=\"flux-sample-[0--1]\"},\n]\n[archive]\ndbpath
    = \"/var/lib/flux/job-archive.sqlite\"\nperiod = \"1m\"\nbusytimeout = \"50s\"\n\n#
    Configure the flux-sched (fluxion) scheduler policies\n# The 'lonodex' match policy
    selects node-exclusive scheduling, and can be\n# commented out if jobs may share
    nodes.\n[sched-fluxion-qmanager]\nqueue-policy = \"fcfs\""
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: flux-sample-entrypoint
  namespace: flux-operator

----
apiVersion: v1
data:
  hostfile: "# Flux needs to know the path to the IMP executable\n[exec]\nimp = \"/usr/libexec/flux/flux-imp\"\n\n[access]\nallow-guest-user
    = true\nallow-root-owner = true\n\n# Point to resource definition generated with
    flux-R(1).\n[resource]\npath = \"/etc/flux/system/R\"\n\n[bootstrap]\ncurve_cert
    = \"/etc/curve/curve.cert\"\ndefault_port = 8050\ndefault_bind = \"tcp://eth0:%p\"\ndefault_connect
    = \"tcp://%h..flux-operator.svc.cluster.local:%p\"\nhosts = [\n\t{ host=\"flux-sample-[0--1]\"},\n]\n[archive]\ndbpath
    = \"/var/lib/flux/job-archive.sqlite\"\nperiod = \"1m\"\nbusytimeout = \"50s\"\n\n#
    Configure the flux-sched (fluxion) scheduler policies\n# The 'lonodex' match policy
    selects node-exclusive scheduling, and can be\n# commented out if jobs may share
    nodes.\n[sched-fluxion-qmanager]\nqueue-policy = \"fcfs\""
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: flux-sample-flux-config
  namespace: flux-operator
```

And that's it! This particular use case is because we wanted to generate miniclusters to be monitored by another operator, and we needed the second operator to handle the actual flux operator container. We needed to manually create the other assets, and this seemed like the logical thing to do.