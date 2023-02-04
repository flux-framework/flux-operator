# Developer Guide

This developer guide includes complete instructions for setting up a developer
environment.

## Setup

To work on this operator you should:

 - Have a recent version of Go installed (1.18.1)
 - Have minikube installed

**Important** For minikube, make sure to enable [DNS plugins](https://minikube.sigs.k8s.io/docs/handbook/addons/ingress-dns/).

```bash
$ minikube addons enable ingress
$ minikube addons enable ingress-dns
```

Note that for production clusters (e.g., [GKE](https://cloud.google.com/kubernetes-engine/docs/concepts/ingress)) I believe this
addon is enabled by default. The basic Flux networking (pods seeing one another) won't work if your cluster does not support DNS.
You also won't be able to expose the service with minikube service if you don't do the above (but port-forward would technically work)
You'll then also want to clone the repository.

```bash
# Clone the source code
$ git clone https://github.com/flux-framework/flux-operator
$ cd flux-operator
```

### Local Development

After cloning, you need to create your MiniKube cluster before doing anything else! 

#### 1. Quick Start

Here is a quick start for doing that, making the namespace, and installing the operator.

```console
# Start a minikube cluster
$ minikube start

# Make a flux operator namespace
$ kubectl create namespace flux-operator
namespace/flux-operator created
```

Here is how to build and install the operator:

```
# Build the operator
$ make

# How to make your manifests
$ make manifests

# And install. This places an executable "bin/kustomize"
$ make install
```

#### 2. Configs

##### Default Config

Our "default" job config, or custom resource definition "CRD" can be found in [config/samples](https://github.com/flux-framework/flux-operator/tree/main/config/samples).
We will walk through how to deploy this one, and then any of the ones in our examples gallery.
First, before launching any jobs, making sure `localDeploy` is set to true in your CRD so you don't ask for a persistent volume claim!

```yaml
spec:
# Set to true to use volume mounts instead of volume claims
  localDeploy: true
```

The default is set to true (so a local volume in `/tmp` is used) but this won't work on an actual cloud Kubernetes cluster,
and vice versa - the volume claim won't work locally. When you are sure this is good,
and you've built and installed the operator, here is how to "launch"
the Mini Cluster (or if providing a command, an ephemeral job):

```bash
$ kubectl apply -f config/samples/flux-framework.org_v1alpha1_minicluster.yaml
```

there is a courtesy function to clean, and apply the samples:

```bash
$ make clean  # remove old flux-operator namespaced items
$ make apply  # apply the setup and job config
$ make run    # make the cluster
```

or run all three for easy development!

```bash
$ make redo  # clean, apply, and run
```

To see logs for the job, you'd do:

```bash
$ kubectl logs -n flux-operator job.batch/flux-sample
```

And this is also:

```bash
$ make log
```

But the above gives you a (somewhat random?) pod. If you want to see a specific one do:

```bash
./script/log.sh flux-sample-0-b5rw6
```

E.g., rank 0 is the "main" one. List running pods (each pod is part of a batch job)

```bash
$ make list
```

And shell into one with the helper script:

```bash
./script/shell.sh flux-sample-0-b5rw6
```

Note that if you are running more than one container, you need a custom command for shell:

```bash
$ kubectl exec --stdin --tty -n flux-operator flux-sample-0-vsnvz -c flux-sample-1 -- /bin/bash
```

Where the first is the name of the pod, and `-c` references the container.

##### Example Configs

We provide an extended gallery of configs for:

 - [Flux Restful Examples](https://github.com/flux-framework/flux-operator/tree/main/examples/flux-restful) that deploy a web interface to submit jobs to.
 - [Headless Test Examples](https://github.com/flux-framework/flux-operator/tree/main/examples/tests) that run commands directly, ideal for testing.

Generally these sets of configs are the same (e.g., same container bases and other options) and only vary in providing a command entrypoint (for the headless)
or not (for the Flux Restful). To make it easy to clean up your cluster and apply a named config, you have two options:

###### Flux Restful 

We call this the main set of "examples." To clean up Minikube, apply a named config, and run the example via a Mini Cluster, first match the name of the yaml
file to a variable name. E.g., for the following files, the names would be:

```console
$ tree examples/flux-restful/
├── minicluster-conveyorlc.yaml       (- name=conveyorlc
├── minicluster-lammps.yaml           (- name=lammps
└── minicluster-osu-benchmarks.yaml   (- name=osu-benchmarks
```

This, to run the full example for conveyorlc:

```bash
$ make name=conveyorlc redo_example 
```

Be careful running this on a production cluster, as it will delete all Kubernetes objects in the namespace.
If you just want to just apply the new job and run the cluster, do:

```bash
$ make name=conveyorlc example
$ make run
```

###### Headless Tests

Our headless tests are modified examples intended to be run without the web interface.
The naming convention is the same as above, except we are concerned with files
in the examples test folder:

```console
$ tree examples/tests/
examples/tests/
├── hello-world           (- name=hello-world
│   └── minicluster.yaml
└── lammps                (- name=lammps
    └── minicluster.yaml
```

Thus, to run the full example for hello-world you can do:

```bash
$ bash script/test.sh hello-world
```
or (for a less scripted run):

```bash
$ make name=hello-world redo_test 
```

If you just want to just apply the new job without a cleanup, do:

```bash
$ make name=hello-world applytest
$ make run
```

Note that there is a default sleep of ~20 seconds for all jobs (so the worker nodes start after the broker)
so they will not run instantly. You can list pods to get an id, and then view logs:

```bash
$ make list
$ bash script/log.sh <pod>
```

Also note that these headless tests have `test: true` in the config,
meaning you will only see output with the command above from what you actually ran.
We do this because we test them in CI, and we don't want the other verbose output
to get in the way! If you want to disable this quiet mode, just comment out
this parameter in the configuration.


### Storage with Rook

We are testing using [rook](https://github.com/rook/rook) that can be used with [minikube](https://rook.io/docs/rook/v1.10/Contributing/development-environment/) or
a [production cluster](https://rook.io/docs/rook/v1.10/Storage-Configuration/Shared-Filesystem-CephFS/filesystem-storage/#shared-volume-creationls ):

```bash
# Create minikube cluster
$ minikube start --disk-size=40g --extra-disks=1 --driver kvm2

# Create aws production "eenie-meenie" cluster, with credentials in environment
$ eksctl create cluster -f examples/storage/ceph/aws/eksctl-config.yaml
```

It's strongly recommended to use the kvm driver in this way, as using docker
can bork your system in strange ways. If you can't use this driver, it's instead recommended
to use a production cluster. For consistently,
the Flux Operator ships with a particular version of the rook yaml files
to create storage, e.g., we did (and you don't need to do this, these files exist):

<details>

<summary>Reproducing creating the yaml configs</summary>

```bash
git clone --depth 1 --single-branch --branch v1.10.10 https://github.com/rook/rook.git
mkdir -p examples/storage/ceph/minikube
mkdir -p examples/storage/ceph/aws

# These are shared by any ceph backend
cp rook/deploy/examples/crds.yaml examples/storage/ceph
cp rook/deploy/examples/common.yaml examples/storage/ceph/
cp rook/deploy/examples/operator.yaml  examples/storage/ceph/

# This is for MiniKube
cp rook/deploy/examples/cluster-test.yaml examples/storage/ceph/minikube/

# This is for a production cluster
cp rook/deploy/examples/cluster.yaml examples/storage/ceph/aws/
```

</details>

For any cluster type, you first can install common CRDs for rook:

```bash
cd examples/storage/ceph
kubectl create -f crds.yaml -f common.yaml -f operator.yaml
cd -
```

Wait until the rook-ceph operator is running before proceeding:

```bash
$ kubectl -n rook-ceph get pod
```

The cluster you create will depend on using MiniKube or a production cluster:

```bash
# MiniKube
$ kubectl create -f examples/storage/ceph/minikube/cluster-test.yaml

# Production (e.g., aws or gcp)
$ kubectl create -f examples/storage/ceph/aws/cluster.yaml
```

Then wait until you can see the pod running:

```bash
$ kubectl -n rook-ceph get pod
```

#### 2. Create Filesystem / Storage class

After either creating the cluster with Minikube or a production cluster (above), wait until the pods are running.

```bash
$ kubectl -n rook-ceph get pod
NAME                                            READY   STATUS      RESTARTS   AGE
csi-cephfsplugin-8qb8k                          2/2     Running     0          9m42s
csi-cephfsplugin-provisioner-75b9f74d7b-xcj6m   5/5     Running     0          9m42s
csi-rbdplugin-provisioner-66d48ddf89-wgvbs      5/5     Running     0          9m42s
csi-rbdplugin-t5rft                             2/2     Running     0          9m42s
rook-ceph-mgr-a-76c787688-9bfdt                 1/1     Running     0          9m11s
rook-ceph-mon-a-7b46c8d7c7-cw9vh                1/1     Running     0          9m35s
rook-ceph-operator-677f8f4c47-ntpxs             1/1     Running     0          12m
rook-ceph-osd-0-576449f5d9-vtgvl                1/1     Running     0          8m42s
rook-ceph-osd-prepare-minikube-zwv77            0/1     Completed   0          8m50s
```

Then we make the filesystem.yaml (we have included at `examples/storage/ceph/filesystem.yaml`)

```bash
$ kubectl create -f examples/storage/ceph/filesystem.yaml
```
Ensure it is running:

```bash
$ kubectl -n rook-ceph get pod -l app=rook-ceph-mds
NAME                                   READY   STATUS    RESTARTS   AGE
rook-ceph-mds-myfs-a-dbc94fc7d-xrl25   1/1     Running   0          75s
rook-ceph-mds-myfs-b-d8494cddb-6r42n   1/1     Running   0          74s
```

and create the storage class (also included):

```bash
$ kubectl create -f examples/storage/ceph/storageclass.yaml
```

At this point, we've created a storage class in the rook-ceph namespace, and we need
to make it available to the flux-operator. Those steps are [outlined here](https://rook.io/docs/rook/v1.10/Storage-Configuration/Shared-Filesystem-CephFS/filesystem-storage/#consume-the-shared-filesystem-across-namespaces).
We first need to make a copy of the `rook-csi-cephfs-node` secret:

```bash
$ kubectl -n rook-ceph describe secrets rook-csi-cephfs-node
```

We will copy this to a different secret, `rook-csi-cephfs-node-user`
but use a different set of key/values. First, save to file:

```bash
$ kubectl get secret rook-csi-cephfs-node -n rook-ceph -o yaml > filesystem-secret.yaml
```

Then make the following changes:

```diff
apiVersion: v1
data:
+  userID: Y3NpLWNlcGhmcy1ub2Rl
-  adminID: Y3NpLWNlcGhmcy1ub2Rl
+  userKey: QVFBZWRkMWpYblB6Q3hBQVp2c3ZYRUkxSWtidE5pVEl1Mk5SNHc9PQ==
-  adminKey: QVFBZWRkMWpYblB6Q3hBQVp2c3ZYRUkxSWtidE5pVEl1Mk5SNHc9PQ==
kind: Secret
metadata:
  creationTimestamp: "2023-02-03T20:57:02Z"
+  name: rook-csi-cephfs-node-user
-  name: rook-csi-cephfs-node
  namespace: rook-ceph
  ownerReferences:
  - apiVersion: ceph.rook.io/v1
    blockOwnerDeletion: true
    controller: true
    kind: CephCluster
    name: rook-ceph
    uid: 8d9b2c68-51a0-48fe-afe4-4f852f83dce9
  resourceVersion: "3817"
  uid: 0c5b103b-1ef4-449e-bfe4-bb8a20cd84e7
type: kubernetes.io/rook
```

And apply

```bash
$ kubectl apply -f filesystem-secret.yaml
```

Now we can create a Persistent Volume Claim, meaning it will create a Persistent Volume for us!
Here is what that might look like:

```yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: base-pvc
  namespace: flux-operator
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 100Gi
  storageClassName: rook-cephfs
  volumeMode: Filesystem
```

Create out flux-operator namespace, and apply and wait until it shows up.

```bash
$ kubectl create namespace flux-operator
```
```bash
$ kubectl apply -f examples/storage/aws/pvc.yaml
```
```bash
$ kubectl get pvc --all-namespaces
NAMESPACE   NAME       STATUS    VOLUME   CAPACITY   ACCESS MODES   STORAGECLASS   AGE
rook-ceph   base-pvc   Pending                                      rook-cephfs    3m35s
```

When it's ready, save it's config:

```bash
$ kubectl get pv base-pvc -n flux-operator -o yaml > base-pvc.yaml
```

Note that the above does not work - I'm next going to try creating the rook storage
in the same namespace as the operator to skip the final mapping across namespaces.


### Interacting with Services

Currently, the most reliable thing to do is port forward:

#### port-forward

If we run as a ClusterIP, we can accomplish the same with a one off `kubectl port-forward`:

```console
kubectl port-forward -n flux-operator flux-sample-0-zdhkp 5000:5000
Forwarding from 127.0.0.1:5000 -> 5000
```

This means you can open [http://localhost:5000](http://localhost:5000) to see the restful API (and interact with it there).

If you want to use a minikube service, this seems to work, but is spotty - I think because
minikube is expecting the service to be available from any pod (it is only running from index 0).
If you want to try this:

```console
$ minikube service -n flux-operator flux-restful-service --url=true
```

But for now I'm developing with port forward.


## Build Images

If you want to build the "production" images - here is how to do that!
This happens in our Docker CI, however you can build (and deploy if you are an owner) them too!

```bash
$ make docker-build
$ make docker-push
```
```bash
# operator lifecycle manager
$ operator-sdk olm install
$ make bundle
$ make bundle-build
$ make bundle-push
```

And for the catalog:

```bash
$ make catalog-build
$ make catalog-push
```

## Build Operator Yaml

To generate the CRD to install to a cluster, we've added a `make build-config` command:

```bash
$ make build-config
```

That will generate a yaml to install the operator (with default container image) to a
cluster in `examples/dist`.

## Container Requirements

If you are looking to build a container to use with the Flux Operator, we have a set of
example containers [here](https://github.com/rse-ops/flux-hpc) and general guidelines are below.
Generally we recommend using the flux-sched base image 
so that install locations and users are consistent. This assumes that:

 - we are currently starting focus on supporting debian bases
 - if you created the flux user, it has uid 1000
 - sudo is available in the container (apt-get install -y sudo)
 - `/etc/flux` is used for configuration and general setup
 - `/usr/libexec/flux` has executables like flux-imp, flux-shell
 - flux-core / flux-sched with flux-security should be installed and ready to go.
 - If you haven't created a flux user, one will be created for you (with a common user id 1000)
 - Any executables that the flux user needs for your job should be on the path (if launching command directly)
 - The container (for now) should start with user root, and we run commands on behalf of flux.
 - You don't need to install the flux-restful-api (it will be installed by the operator)
 - munge should be install, and a key generated at `/etc/munge/munge.key`
  
For the last point, since all Flux running containers should have the same munge key
in that location, we simply use it. The pipeline will fail if the key is missing from any
Flux runner container.  For the curve.cert that we need to secure the cluster, we will
be running your flux runner container before the indexed job is launched, generating
the certificate, and then mapping it into the job pods via another config map.
Note that we considered generating this natively in Gom, however the underlying library to do this 
generation that is [available in Go](https://pkg.go.dev/github.com/zeromq/goczmq#section-readme)
requires system libraries, and thus would be annoying to add as a dependency.

These criteria are taken from the [flux-sched](https://github.com/flux-framework/flux-sched/blob/master/src/test/docker/focal/Dockerfile)
base image [available on Docker hub](https://hub.docker.com/r/fluxrm/flux-sched) as `fluxrm/flux-sched:focal`, 
and we strongly suggest you use this for your base container to make development easier! 
If you intend to use the [Flux RESTful API](https://github.com/flux-framework/flux-restful-api)
to interact with your cluster, ensure that flux (python bindings) are on the path, along with
either python or python3 (depending on which you used to install Flux).
If/when needed we can lift some of these constraints, but for now they are 
reasonable. If you use this image, you should have python3 and pip3 available to you,
and the active user is `fluxuser`. This means if you want to add content, either you'll
need to change the user to `root` in a build (and back to `fluxuser` at the end), use sudo, or
install to `/home/fluxuser`.


## Testing

Testing is underway! From a high level, we want three kinds of testing:

 - Unit tests, which will be more traditional `*_test.go` files alongside others in the repository (not done yet)
 - End to end "e2e" tests, also within Go, to test an entire submission of a job, also within Go. (not done yet)
 - Integration testing, likely some within Go and some external to it. (in progress)

### Integration

For the integration testing outside of Go, we currently have basic tests written that allow the following:

1. Write a custom resource definition (CRD) for a named mini cluster under `examples/tests/${name}` as `minicluster.yaml`.
2. The CRD should set `test:true` and include a command to run, and a container to do it.
3. Add your test name, container, and estimated running time to `.github/workflows/main.yaml`
4. If your tests require a working directory, it must be set in the CRD for the headless test.
5. If a test is deterministic, add a `test.out` to the output folder that we can validate results for.
6. We will validate output (if provided) and that containers exit with 0.

To run the test (and you can also do this locally) we use the `script/test.sh` and provide a name and the estimated
job time, just like in actions. The below example runs the "hello-world" test and gives it 30 seconds to finish.

```bash
./bin/bash script/test.sh hello-world 10
```
```console
...
NAME                  READY   STATUS    RESTARTS   AGE
flux-sample-0-ng8c5   1/1     Running   0          3s
flux-sample-1-nxqj9   1/1     Running   0          3s
flux-sample-2-vv7jr   1/1     Running   0          3s
flux-sample-3-kh2br   1/1     Running   0          3s
Pods: flux-sample-0-ng8c5 flux-sample-1-nxqj9 flux-sample-2-vv7jr flux-sample-3-kh2br
Pod: flux-sample-0-ng8c5
Actual:
hello world
Expected:
hello world
```
What you don't see in the above is that we also use kubectl to ensure that the exit code for all containers
(typically 4) is 0. Also note that the "sleep" time doesn't have to be exact, it's technically not necessary because we are waiting
for the output to finish coming (and the job to stop). I added it to provide a courtesy message to the user and developer.
Finally, note that for big containers it's suggested to pull them first, e.g.,:

```bash
$ minikube ssh docker pull ghcr.io/rse-ops/lammps:flux-sched-focal-v0.24.0
```
The tests above are headless, meaning they submit commands directly, and that way
we don't need to do it in the UI and can programmatically determine if they were successful.
Finally, note that if you have commands that you need to run before or after the tests,
you can add a `pre-run.sh` or `post-run.sh` in the directory. As an example, because minikube
is run inside of a VM, if you are using a host volume mount, it won't actually show up on your host!
This is because it's inside the VM. This you might want to move files there before the test, e.g.,:

```bash
#!/bin/bash

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "Copying local volume to /tmp/data-volumes in minikube"

# We don't care if this works or not - mkdir -p seems to bork
minikube ssh -- mkdir -p /tmp/data-volumes
minikube cp ${HERE}/data/pancakes.txt /tmp/data-volumes/pancakes.txt
minikube ssh ls /tmp/data-volumes
```

and then clean up after

```bash
#!/bin/bash

echo "Cleaning up /tmp/data-volumes in minikube"
minikube ssh -- sudo rm -rf /tmp/data-volumes
```

This would be the same for anytime you use minikube and want to create a local volume. It's not actually on your host,
but rather in the VM. For an example test that does this, see the `examples/tests/volumes` example.

## Documentation

The documentation is provided in the `docs` folder of the repository, and generally most content that you might want to add is under
`getting_started`. For ease of contribution, files that are likely to be updated by contributors (e.g., mostly everything but the module generated files)
are written in markdown. If you need to use [toctree](https://www.sphinx-doc.org/en/master/usage/restructuredtext/directives.html#table-of-contents) you should not use extra newlines or spaces (see index.md files for examples). The documentation is also provided in Markdown (instead of rst or restructured syntax) to make contribution easier for the community.

Finally, we recommend you use the same development environment also to build and work on
documentation. The reason is because we import the app to derive docstrings,
and this will require having Flux.

### Install Dependencies and Build

The documentation is built using sphinx, and generally you can 
create a virtual environment:

```bash
$ cd docs
$ python -m venv env 
$ source env/bin/activate
```
And then install dependencies:

```console
$ pip install -r requirements.txt

# Build the docs into _build/html
$ make html
```

### Preview Documentation

After `make html` you can enter into `_build/html` and start a local web
server to preview:

```console
$ python -m http.server 9999
```

And open your browser to `localhost:9999`
