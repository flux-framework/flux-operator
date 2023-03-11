# Volumes

These short examples will describe advanced functionality for volumes. For examples,
see our [storage examples directory](https://github.com/flux-framework/flux-operator/tree/main/examples/storage).

## Existing Persistent Volume

It might be the case that you've already defined a persistent volume claim, and you simply want to use it.
We currently support this, and require that you manage both the PVC and the PV (in our testing,
when a PV was created beforehand and then a PVC created by the operator, it would get status "Lost").
We think it might be possible to create the PV first (and have the PVC created by the operator)
but more testing is needed. Thus, for the time being, we recommend that you create your own
PV and PVC in advance, and then provide it to the operator. Here is an example
workflow that will use a pre-defined persistent volume claim:


```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
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

Given our minicluster defined to use the claim named "data":

```yaml
apiVersion: flux-framework.org/v1alpha1
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

  # This is a list because a pod can support multiple containers
  containers:
    - image: ghcr.io/flux-framework/flux-restful-api:latest
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

And then we are ready to create the minicluster!

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