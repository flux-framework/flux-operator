# Scaling

> This functionality requires [Kubernetes 1.27 and later](https://github.com/kubernetes/enhancements/tree/master/keps/sig-apps/3715-elastic-indexed-job#motivation).

While Flux does not natively support scaling or elasticity (yet) we can do some tricks
with the Flux Operator to enable it! Specifically:

 - We tell Flux to create a cluster at the maximum size that is possible (and the broker config sees this many nodes)
 - We update the resource spec to reflect that.


## Basic Example

> Starting with a cluster at a maximum size and scaling it up and down

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/scaling/basic/minicluster.yaml)**

To run this example:

```bash
$ minikube start --kubernetes-version=1.27.
```

Install the operator, create the namespace, and create the MiniCluster:

```bash
$ kubectl apply -f ./examples/dist/flux-operator.yaml
$ kubectl create namespace flux-operator
$ kubectl apply -f examples/scaling/basic/minicluster.yaml
```

You'll need to wait for the containers to create (the image is pulling) and you can
help the pull via:

```bash
$ minikube ssh docker pull ghcr.io/flux-framework/flux-restful-api:latest
```

### Check Initial Size

Wait until the cluster finishes, and you see the pods are ready to go (Running state):

```bash
kubectl get -n flux-operator pods
NAME                  READY   STATUS    RESTARTS   AGE
flux-sample-0-kfl7p   1/1     Running   0          76s
flux-sample-1-2v57b   1/1     Running   0          76s
flux-sample-2-m7n85   1/1     Running   0          76s
flux-sample-3-zkmvq   1/1     Running   0          76s
```

We recommend (in another terminal) shelling into the broker pod and connecting to the broker's Flux instance
so that you can follow the changes.

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-xd2gc -- bash
root@flux-sample-0:/code# sudo -E $(env) -E HOME=/home/flux -u flux flux proxy local:///var/run/flux/local 
```

Here is how to look at the state of the cluster. When you first create it, we will have 4 pods, and all of them
are up.

```bash
flux@flux-sample-0:/code$ flux resource list
     STATE NNODES   NCORES NODELIST
      free      4       16 flux-sample-[0-3]
 allocated      0        0 
      down      0        0 
```

At this point we want to try changing the size.

### Ask for a Larger Size

```diff
-  size: 4
+  size: 5
```

Let's first try asking for something the operator can't give us - a larger size. After the change above, do:

```bash
$ kubectl apply -f examples/scaling/basic/minicluster.yaml
```

The reason a larger size isn't supported is because the Flux broker already has registered N nodes, known by their fully qualified
domain name, and we would need to do some tricks to update that configuration to add another
one. While this might be possible (and likely will be in the future) for now we don't support it.
Thus, if you do a request for a size that is larger than the originally created maximum size,
you'll see this in the operator logs:

```console
1.6831543179369373e+09  INFO    minicluster-reconciler  MiniCluster     {"Size": 4, "Requested Size": 5}
1.6831543179369428e+09  INFO    minicluster-reconciler  MiniCluster     {"PatchSize": 5, "Status": "Denied"}
```

### Ask for Smaller Size

Asking for a smaller size will work! Let's decrease the original CRD from 4 to 3:

```diff
-  size: 4
+  size: 3
```

Apply the CRD again:

```bash
$ kubectl apply -f examples/scaling/basic/minicluster.yaml
```

The first thing you will notice is that a pod is terminating

```
 make list
kubectl get -n flux-operator pods
NAME                  READY   STATUS        RESTARTS   AGE
flux-sample-0-xd2gc   1/1     Running       0          30m
flux-sample-1-tbj7c   1/1     Running       0          30m
flux-sample-2-wbpf9   1/1     Running       0          30m
flux-sample-3-cfs6c   1/1     Terminating   0          26m
```

When the pod is gone, if (in your second terminal) you look at the resource status, Flux
will now report this pod as down.

```bash
flux@flux-sample-0:/code$ flux resource list
     STATE NNODES   NCORES NODELIST
      free      3       12 flux-sample-[0-2]
 allocated      0        0 
      down      1        4 flux-sample-3
```

And importantly, the rest of the cluster keeps running smoothly! We haven't interrupted the Flux broker or
install by changing the size, at least superficially and not running any jobs on the pod that was terminated.

### Ask for Larger Size

Finally, let's scale back up! Restore the original CRD size to 4:

```diff
-  size: 3
+  size: 4
```

Apply again:

```bash
$ kubectl apply -f examples/scaling/basic/minicluster.yaml
```
And time time you'll see the container creating:

```bash
$ kubectl get -n flux-operator pods
NAME                  READY   STATUS              RESTARTS   AGE
flux-sample-0-xd2gc   1/1     Running             0          35m
flux-sample-1-tbj7c   1/1     Running             0          35m
flux-sample-2-wbpf9   1/1     Running             0          35m
flux-sample-3-ll76s   0/1     ContainerCreating   0          1s
```

And when it's running, Flux will notice it online again. Your full cluster is online again.
And that's it!

```bash
$ flux resource list
     STATE NNODES   NCORES NODELIST
      free      4       16 flux-sample-[0-3]
 allocated      0        0 
      down      0        0 
```

We will have a tutorial for expanding a cluster size soon. Flux doesn't allow
the hosts to be greater than nodes currently, so we haven't added this yet.

## Expand Example

> Starting with a small cluster that is able to grow to a maximum size

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/scaling/expand/minicluster.yaml)**

To run this example:

```bash
$ minikube start --kubernetes-version=1.27.
```

Install the operator, create the namespace, and create the MiniCluster:

```bash
$ kubectl apply -f ./examples/dist/flux-operator.yaml
$ kubectl create namespace flux-operator
$ kubectl apply -f examples/scaling/basic/minicluster.yaml
```

### Create the cluster

First, apply the CRD to create the MiniCluster. Note that we are asking for a size of 2, but allowing
for a maximum size of 4. 

```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:
  # Number of pods to create for MiniCluster to start
  size: 2

  # Number of pods to allow scaling to (the number that flux will see)
  maxSize: 4

  # This needs to be in interactive or launcher mode to work
  # otherwise we submit as a job (and it will be running under the smaller size number of tasks)
  interactive: true

  # This is a list because a pod can support multiple containers
  containers:
    - image: ghcr.io/flux-framework/flux-restful-api:latest
```
```bash
$ kubectl apply -f ./examples/scaling/expand/minicluster.yaml
```
Since our initial size is 2, you'll see two pods creating:

```bash
kubectl get -n flux-operator pods
NAME                  READY   STATUS              RESTARTS   AGE
flux-sample-0-r2cxt   0/1     ContainerCreating   0          1s
flux-sample-1-bxwbw   0/1     ContainerCreating   0          1s
```