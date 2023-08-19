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
└── fluxoperator-gen
```

Or if you pull and have the container, interact as follows:

```bash
$ docker run -it --entrypoint fluxoperator-gen  ghcr.io/flux-framework/flux-operator:latest --help
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
Here is how to run with a local file:

```bash
$ docker run -it --entrypoint fluxoperator-gen -v $PWD:/code ghcr.io/flux-framework/flux-operator:latest -f /code/examples/tests/lammps/minicluster.yaml
```

For the above, we bind the present working directory to `/code` and then provide a path with `-f` relative to it.
You could provide additional arguments after the flux operator container.

> **Important** if you run the gen.go file or build without copying over the keygen template, you will miss generating
the curve certificate, which is done because zeromq is compiled into the container. If you see it missing, double check your generation logic! There should be three config maps total - one for entrypoints, one for flux configs, and one for the curve certificate.

## fluxoperator-gen

For some cases of using the Flux operator where dynamism is involved (either in scaling or having custom volumes set up)
you typically need the entire operator. However, for cases where you are using the Flux Operator as more of a job submission
tool (e.g., akin to what we do in [Kueue](https://kueue.sigs.k8s.io/docs/tasks/run_flux_minicluster/)) you really only
need to generate the core assets for your MiniCluster, which are an indexed job, config maps, (optionally volumes) and a service. If you think about it, there are actually two cases of operator types:

 - *helicopter parent* meaning that your objects warrant constant monitoring for updating. For this case, the operator needs to create, delete, and perform other update operations that would be challenging (or annoying) to do manually.
 - *80s/90s parent* they might drop you off at the birthday party, but you are on your own after that, and maybe even need to walk yourself home! For this case, the operator only exists to create and delete.

After you realize this distinction, you also realize that the second case - the more "I will make you and let you be" case is well-suited to be served by static YAML files. Indeed we want the operator to generate them because the logic is really hairy, but we don't need it to do anything beyond that. The operator is just a fancy, programmatic way to produce complex configs. Thus, given this case for our Flux MiniClusters that don't require scaling or otherwise changing, we think it is useful to provide
an extra helper, "fluxoperator-gen" that does exactly that. If you are curious, for our use case we only needed the config maps
for a [metrics operator](https://github.com/converged-computing/metrics-operator-experiments/tree/main/flux-operator) experiment,
but now we now provide the tool for you too! A quick note:

 - We don't currently support sidecar services. The assumption is that you can generate them yourself! If you think we should add this support, [let us know](https://github.com/flux-framework/flux-operator/issues). 
 - The generation outputs null creationTimestamp and an empty status block that is largely not needed, and can be deleted.

### fluxoperator-gen usage

Let's walk through some examples. We will assume you have the `fluxoperator-gen` built or provided by the container.
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
  curve.cert: "#   ****  Generated on 2023-04-26 22:54:42 by CZMQ  ****\n#   ZeroMQ
    CURVE **Secret** Certificate\n#   DO NOT PROVIDE THIS FILE TO OTHER USERS nor
    change its permissions.\n    \nmetadata\n    name = \"flux-cert-generator\"\n
    \   keygen.hostname = \"flux-sample-0\"\ncurve\n    public-key = \"[@WjAzG&(B:Yf84Ge/#MrQ89N[]AtCL/v*(R7P2y\"\n
    \   secret-key = \"h^Qx2ID84@uQQpy5u+@lJ-yv!O9vRrDX{up^<CxI\"\n"
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
  curve.cert: "#   ****  Generated on 2023-04-26 22:54:42 by CZMQ  ****\n#   ZeroMQ
    CURVE **Secret** Certificate\n#   DO NOT PROVIDE THIS FILE TO OTHER USERS nor
    change its permissions.\n    \nmetadata\n    name = \"flux-cert-generator\"\n
    \   keygen.hostname = \"flux-sample-0\"\ncurve\n    public-key = \"[@WjAzG&(B:Yf84Ge/#MrQ89N[]AtCL/v*(R7P2y\"\n
    \   secret-key = \"h^Qx2ID84@uQQpy5u+@lJ-yv!O9vRrDX{up^<CxI\"\n"
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
data:
  curve.cert: "#   ****  Generated on 2023-04-26 22:54:42 by CZMQ  ****\n#   ZeroMQ
    CURVE **Secret** Certificate\n#   DO NOT PROVIDE THIS FILE TO OTHER USERS nor
    change its permissions.\n    \nmetadata\n    name = \"flux-cert-generator\"\n
    \   keygen.hostname = \"flux-sample-0\"\ncurve\n    public-key = \"[@WjAzG&(B:Yf84Ge/#MrQ89N[]AtCL/v*(R7P2y\"\n
    \   secret-key = \"h^Qx2ID84@uQQpy5u+@lJ-yv!O9vRrDX{up^<CxI\"\n"
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
  name: flux-sample-curve-mount
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
  curve.cert: "#   ****  Generated on 2023-04-26 22:54:42 by CZMQ  ****\n#   ZeroMQ
    CURVE **Secret** Certificate\n#   DO NOT PROVIDE THIS FILE TO OTHER USERS nor
    change its permissions.\n    \nmetadata\n    name = \"flux-cert-generator\"\n
    \   keygen.hostname = \"flux-sample-0\"\ncurve\n    public-key = \"&g$oLyZJSr3/(MUD+w?:p9!YnW-ydG8Iccs.zM/[\"\n
    \   secret-key = \"Pc$.HO&)5P^:^C7UgKV[.+AG]w/(Jv8ZQePr/{(n\"\n"
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
  curve.cert: "#   ****  Generated on 2023-04-26 22:54:42 by CZMQ  ****\n#   ZeroMQ
    CURVE **Secret** Certificate\n#   DO NOT PROVIDE THIS FILE TO OTHER USERS nor
    change its permissions.\n    \nmetadata\n    name = \"flux-cert-generator\"\n
    \   keygen.hostname = \"flux-sample-0\"\ncurve\n    public-key = \"&g$oLyZJSr3/(MUD+w?:p9!YnW-ydG8Iccs.zM/[\"\n
    \   secret-key = \"Pc$.HO&)5P^:^C7UgKV[.+AG]w/(Jv8ZQePr/{(n\"\n"
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
data:
  curve.cert: "#   ****  Generated on 2023-04-26 22:54:42 by CZMQ  ****\n#   ZeroMQ
    CURVE **Secret** Certificate\n#   DO NOT PROVIDE THIS FILE TO OTHER USERS nor
    change its permissions.\n    \nmetadata\n    name = \"flux-cert-generator\"\n
    \   keygen.hostname = \"flux-sample-0\"\ncurve\n    public-key = \"&g$oLyZJSr3/(MUD+w?:p9!YnW-ydG8Iccs.zM/[\"\n
    \   secret-key = \"Pc$.HO&)5P^:^C7UgKV[.+AG]w/(Jv8ZQePr/{(n\"\n"
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
  name: flux-sample-curve-mount
  namespace: flux-operator
```

And that's it! This particular use case is because we wanted to generate miniclusters to be monitored by another operator, and we needed the second operator to handle the actual flux operator container. We needed to manually create the other assets, and this seemed like the logical thing to do.