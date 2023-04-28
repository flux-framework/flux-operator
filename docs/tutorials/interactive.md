# Interactive

The following tutorials demonstrate interactive MiniClusters, beyond using the Flux Restful API or
launching a job to run and complete.

## Persistent Cluster

> This example is for a persistent minicluster that provides a shell

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/interactive/minicluster-persistent.yaml)**

This example demonstrates bringing up a MiniCluster solely to shell in and interact with Flux. First, 
note that there is nothing special about the MiniCluster YAML except that we set interactive to true:

```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:

  # Number of pods to create for MiniCluster
  size: 2
  interactive: true

# This is a list because a pod can support multiple containers
  containers:
    - image: ghcr.io/rse-ops/pokemon:app-latest
      fluxOptionFlags: "-ompi=openmpi@5"
```

When interactive is true, we tell the Flux broker to start without a command. This means
the cluster will remain running until you shutdown Flux or `kubectl delete` the MiniCluster
itself. The container you choose should have the software you are interested in having for each node.
Given a running cluster, we can create the namespace and the MiniCluster as follows:

```bash
$ kubectl create namespace flux-operator
$ kubectl apply -f examples/interactive/minicluster-persistent.yaml
```

We can then wait for our pods to be running

```bash
$ kubectl get -n flux-operator pods
NAME                         READY   STATUS      RESTARTS   AGE
flux-sample-0-p5xls          1/1     Running     0          7s
flux-sample-1-nmtt7          1/1     Running     0          7s
flux-sample-cert-generator   0/1     Completed   0          7s
```

And then shell into the broker pod, index 0:

```bash
$ kubectl exec -it  -n flux-operator flux-sample-0-p5xls -- bash
```

At this point, remember the broker is running, and we need to connect to it. We do this via
flux proxy and targeting the socket, which is a local reference at `/run/flux/local`:

```bash
# Connect to the flux socket at /run/flux/local as the flux instance owner "flux"
$ sudo -u flux flux proxy local:///run/flux/local
```

At this point, you are the instance owner "flux":

```bash
$ whoami
flux
```
And can also see the resources known to your instance!

```bash
flux@flux-sample-0:/code$ flux resource list
     STATE NNODES   NCORES NODELIST
      free      2        2 flux-sample-[0-1]
 allocated      0        0
      down      0        0 
```

If you see a node is down, it could be the worker node has not connected to the broker
yet. When the cluster is ready, it will look like the above! At this point you 
have your own Flux install connected to your entire MiniCluster,
and you can launch and monitor jobs as you please.

```bash
$ flux submit sleep 120
ƒg4ZRVS7
```

Note that under flux jobs below, the first sleep (flux-sample-1) is the command
originally given to the broker. The second, newer sleep is the job we just launched.

```bash
flux@flux-sample-0:/code$ flux jobs
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
    ƒg4ZRVS7 flux     sleep       R      1      1   1.290s flux-sample-1
```

When you are done, there are a few ways to "get out." You could just shutdown Flux from
inside of the instance:

```bash
$ flux shutdown
```

Or you could exit from the instance, and exit from the pod, and then delete the MiniCluster.

```bash
$ kubectl delete -f examples/interactive/minicluster-persistent.yaml
```

If your MiniCluster spec does not have `clean: true` you might need to run this 
anyway with a shutdown. List pods if you aren't sure. Finally, interactive mode works
very well when using the Flux Operator Python SDK. [Here is a full example](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/examples/interactive-submit.py).

Have fun! 🦄