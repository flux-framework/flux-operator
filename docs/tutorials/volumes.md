# Volumes

These short examples will describe advanced functionality for volumes. For examples,
see our [storage examples directory](https://github.com/flux-framework/flux-operator/tree/main/examples/storage).


## Local Volumes

For testing it's helpful to use a local volume, or basically a path on your host that is bound to either
the MiniKube Virtual Machine, or to a kind cluster container as a volume. We will show you how to do both here.

### Kind

I would first recommend kind, as I've had better luck using it. With kind, the control plane
is run as a container, and so a local volume is a simple bind to your host. You need to create
the cluster and ask for the bind. Here is a yaml file example for how to do that, binding
`/tmp/workflow` from the path on your host to the same path in the container.

```yaml
apiVersion: kind.x-k8s.io/v1alpha4
kind: Cluster
nodes:
   - role: control-plane
     extraMounts:
      - hostPath: /tmp/workflow
        containerPath: /tmp/workflow
```
You'd create this as follows:

```bash
$ kind create cluster --config kind-config.yaml
```

<details>

<summary>Output of kind create cluster</summary>

```console
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.24.0) 🖼 
 ✓ Preparing nodes 📦  
 ✓ Writing configuration 📜 
 ✓ Starting control-plane 🕹️ 
 ✓ Installing CNI 🔌 
 ✓ Installing StorageClass 💾 
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Not sure what to do next? 😅  Check out https://kind.sigs.k8s.io/docs/user/quick-start/
```

Don't worry, there are plenty of emojis, as you can see.

</details>

You can then use `/tmp/workflow` as a local volume. An example workflow is provided as 
[our volumes test](https://github.com/flux-framework/flux-operator/blob/main/examples/tests/volumes/minicluster.yaml).

### MiniKube

For minikube, the design is different, so a local volume comes down to a directory inside of the virtual
machine. For data that needs to be read and written (and not executed) I've found that defining the local volume
works fine - it will exist in the VM and be shared by the pods. In this case, you might just want to create
the volume and copy over some small data file(s):

```bash
$ minikube ssh -- mkdir -p /tmp/data
$ minikube cp ./data/pancakes.txt /tmp/data/pancakes.txt
$ minikube ssh ls /tmp/data
```

However, if you are going to execute binaries,
I've run into trouble not actually binding the path to the host. In this case, we need to do:

```console
[host]                [vm]              [pod]

/tmp/workflow   ->    /tmp/workflow  -> /data
```

And the easiest way to do this is (in another terminal) run `minikube mount`

```bash
$ minikube mount /tmp/workflow:/tmp/workflow
```

The above mounts your hostpath `/tmp/workflow` to `/tmp/workflow` in the virtual machine,
and then pods will access it via a named volume in the minicluster.yaml:

```yaml
apiVersion: flux-framework.org/v1alpha2
kind: MiniCluster
metadata:
  name: flux-sample
spec:

  # Number of pods to create for MiniCluster
  size: 2

  # Make this kind of persistent volume and claim available to pods
  # This is a path in minikube (e.g., minikube ssh or minikube mount)
  volumes:
    data:
      storageClass: hostpath
      path: /tmp/workflow

  containers:
    - image: rockylinux:9
      command: ls /data

      # The name (key) "data" corresponds to the one in "volumes" above
      # so /data in the container is bound to /tmp/workflow in the VM, and this
      # is the path you might have done minikube mount from your host
      volumes:
        data:
          path: /data
```

Beyond these simple use cases, you likely want to use an existing persistent volume / claim,
or a container storage interface.

## Existing Persistent Volume

It might be the case that you've already defined a persistent volume claim, and you simply want to use it.
We currently support this, and require that you manage both the PVC and the PV (in our testing,
when a PV was created beforehand and then a PVC created by the operator, it would get status "Lost").
We think it might be possible to create the PV first (and have the PVC created by the operator)
but more testing is needed. Thus, for the time being, we recommend that you create your own
PV and PVC in advance, and then provide it to the operator. Here is an example
workflow that will use a pre-defined persistent volume claim:


```yaml
apiVersion: flux-framework.org/v1alpha2
kind: MiniCluster
metadata:
  name: flux-sample
spec:

  # Number of pods to create for MiniCluster
  size: 2

  # show all operator output and test run output
  logging:
    quiet: false

        
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

      # This is where storage is mounted - we are only going to touch a file
      workingDir: /workflow
      command: touch test-file.txt

      # Commands just for workers / broker
      commands:

        # Running flux as root is currently required for the storage driver to work
        runFluxAsRoot: true
```

In the above, we've created a PVC called "data" and we want it to be mounted to "/workflow" in the container.
Note that we are currently running flux as root because it's the lazy way to ensure the volume mount works,
however it's not ideal from a security standpoint. More testers are needed to test different (specific)
container storage interfaces (or other volume types) to find the correct mount options to set in order
to allow ownership by the "flux" user (or for more advanced cases, breaking storage into pieces to be
owned by a set of cluster users)!


### Example

You can see an example in the [existing-volumes](https://github.com/flux-framework/flux-operator/tree/main/examples/tests/existing-volumes)
test, where we create a host path volume and copy a file there in advance, and then use it during the operator execution. Specifically,
we have the following object manifest for the persistent volume:

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: flux-sample
  namespace: flux-operator
spec:
  persistentVolumeReclaimPolicy: Delete
  capacity:
    storage: 5Gi
  accessModes:
    - ReadWriteMany
  hostPath:
    path: /tmp/data
```

And the persistent volume claim:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: data
  namespace: flux-operator
spec:
  accessModes:
    - ReadWriteMany
  storageClassName: ""
  resources:
    requests:
      storage: 1Gi
```

Given our MiniCluster defined to use the claim named "data":

```yaml
apiVersion: flux-framework.org/v1alpha2
kind: MiniCluster
metadata:
  name: flux-sample
spec:
  # suppress all output except for test run
  logging:
    quiet: true

  # Number of pods to create for MiniCluster
  size: 2

  # This is a list because a pod can support multiple containers
  containers:
    - image: rockylinux:9
      command: ls /data
      existingVolumes:
        data:
          path: /data
          claimName: data
```

We can start minikube:

```bash
$ minikube start
```

Create the operator namespace:

```bash
$ kubectl create -n flux-operator
```

Prepare the volume directory in the minikube virtual machine:

```bash
$ minikube ssh -- mkdir -p /tmp/data
$ minikube cp .pancakes.txt /tmp/data/pancakes.txt
$ minikube ssh ls /tmp/data
```

And then create each of the PV and PVC above:

```bash
$ kubectl apply -f ./pv.yaml 
$ kubectl apply -f ./pvc.yam
```

And then we are ready to create the MiniCluster!

```bash
$ kubectl apply -f minicluster.yaml
```

If you watch the logs, the command is doing an ls to data, so you should see `pancakes.txt`.

Note that for cleanup, pods typically need to be deleted before pvc and pv. Since the MiniCluster
doesn't own either of the PV nor PVC, the easiest thing to do is:

```bash
$ kubectl delete -n flux-operator pods --all --grace-period=0 --force
$ kubectl delete -n flux-operator pvc --all --grace-period=0 --force
$ kubectl delete -n flux-operator pv --all --grace-period=0 --force
$ kubectl delete -f minicluster.yaml
```

You can also delete the MiniCluster first and retain the PV and PVC. Have fun!
