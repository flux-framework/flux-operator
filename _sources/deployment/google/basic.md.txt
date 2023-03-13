# Lammps on Google Cloud

This basic tutorial will walk through creating a MiniCluster to run LAMMPS! You should have
already [setup your workspace](setup.md)


```{include} includes/basic-setup.md
```


Now let's run a short experiment with LAMMPS!

## Custom Resource Definition

The Custom Resource Definition (CRD) defines our Mini Cluster, and is what we hand to the flux
operator to create it.  Here is the CRD for a small lammps run.

```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:
  # Number of pods to create for MiniCluster
  size: 4

  # Disable verbose output
  logging:
    quiet: true

  # This is a list because a pod can support multiple containers
  containers:
    # The container URI to pull (currently needs to be public)
    - image: ghcr.io/rse-ops/lammps:flux-sched-focal-v0.24.0

      # You can set the working directory if your container WORKDIR is not correct.
      workingDir: /home/flux/examples/reaxff/HNS
      command: lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite
```

You can save the above file as `minicluster-lammps.yaml` to get started.

## Create the Lammps Job

Now let's apply the custom resource definition to create the lammps mini cluster!
The file we generated above should be in your present working directory.
Importantly, we have set `localDeploy` to false because we need to create volume
claims and not local host mounts for shared resources.

```bash
$ kubectl apply -f minicluster-lammps.yaml
```

There are different ways to see logs for pods. First, see pods running and state.
You probably want to wait until the state changes from `ContainersCreating` to `Running`
because this is where we are pulling the chonker containers.

```bash
$ kubectl get -n flux-operator pods
```

If you need to debug (or see general output for a pod about creation) you can do:

```bash
$ kubectl -n flux-operator describe pods flux-sample-0-742bm
```

And finally, the most meaty piece of metadata is the log for the pod,
where the Flux Operator will be setting things up and starting flux.

```bash
# Add the -f to keep it hanging
$ kubectl -n flux-operator logs flux-sample-0-742bm -f
```

To shell into a pod to look around (noting where important flux stuff is)

```bash
$ kubectl exec --stdin --tty -n flux-operator flux-sample-0-742bm -- /bin/bash
```
```console
ls /mnt/curve
ls /etc/flux
ls /etc/flux/config
```

To get logs for the operator itself:

```bash
$ kubectl logs -n operator-system operator-controller-manager-56b5bcf9fd-j2g75
```

If you need to run in verbose (non-test) mode, set test to false in the [minicluster-lammps.yaml](minicluster-lammps.yaml).
And make sure to clean up first:

```bash
$ kubectl delete -f minicluster-lammps.yaml
```

and wait until the pods are gone:

```bash
$ kubectl get -n flux-operator pods
No resources found in flux-operator namespace.
```

Observations about comparing this to MiniKube (local):

 - The containers that are large actually pull!
 - The startup times of the different pods vary quite a bit.

If you want to run the same workflow again, use `kubectl delete -f` with the file
and apply it again. I wound up running with test set to true, and then saving the logs:

```bash
$ kubectl -n flux-operator logs flux-sample-0-qc5z2 > lammps.out
```

For fun, here is the first successful run of Lammps using the Flux Operator on GCP
ever!

![img/lammps.png](img/lammps.png)

Then to delete your lammps MiniCluster:

```bash
$ kubectl delete -f minicluster-lammps.yaml
```

```{include} includes/cleanup.md
```

## Customization and Debugging

### Firewall

When I first created my cluster, the nodes could not see one another. I added a few
flags for networking, and looked at firewalls as follows:

```bash
$ gcloud container clusters describe flux-cluster --zone us-central1-a | grep clusterIpv4Cidr
```
I didn't ultimately change anything, but I found this useful.
