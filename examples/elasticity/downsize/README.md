# Downsize

In this example, we will create a MiniCluster using the [JobBackoffPerIndex](https://kubernetes.io/blog/2023/08/21/kubernetes-1-28-jobapi-update/#backoff-limit-per-index) feature gate (Kubernetes 1.28) that allows us to set an explicit number of failures allowed. We expose this as a boolean `restartWorkers` (defaults to true) We want to test:

- When we drain and delete a pod, it is not recreated
- What happens when we resize the cluster (does the index get recreated)?

## Cluster

Create the cluster with the feature gate.

```bash
kind create cluster --config kind-config.yaml
```

Then install the operator (however you prefer) and when it's running, create the interactive MiniCluster.

```bash
kubectl apply -f minicluster.yaml
```

Ensure you have four pods running.

```bash
kubectl get pods
```
```console
NAME                  READY   STATUS    RESTARTS   AGE
flux-sample-0-9zsjw   1/1     Running   0          40s
flux-sample-1-z94mv   1/1     Running   0          40s
flux-sample-2-rm4ks   1/1     Running   0          40s
flux-sample-3-5qbsw   1/1     Running   0          40s
```

## 1. Drain a Node

Now we are going to drain a node, preparing to delete it.
If this works, it should not recreate.

```bash
$ kubectl exec -it flux-sample-0-9zsjw bash
. /mnt/flux/flux-view.sh 
flux proxy $fluxsocket bash
```
```bash
     STATE NNODES   NCORES    NGPUS NODELIST
      free      4       40        0 flux-sample-[0-3]
 allocated      0        0        0 
      down      6       60        0 flux-sample-[4-9]
```
Remember that we told Flux it can support 10 nodes, but only want to start with 4 (hence the resources above). To drain:

```bash
flux resource drain flux-sample-2 "I am a bad person"
```

It should be down.

```console
     STATE NNODES   NCORES    NGPUS NODELIST
      free      3       30        0 flux-sample-[0-1,3]
 allocated      0        0        0 
      down      7       70        0 flux-sample-[2,4-9]
```

## 2. Delete a Pod

Now we are going to be terrible and manually delete the pod.

```bash
kubectl delete pod flux-sample-2-xxx
```

You should see the pod go away (and not come back)!

```console
NAME                  READY   STATUS    RESTARTS   AGE
flux-sample-0-pdqzl   1/1     Running   0          2m42s
flux-sample-1-rmzhf   1/1     Running   0          2m42s
flux-sample-3-cfgwr   1/1     Running   0          2m42s
```

And if you peek inside, Flux doesn't see anything different - the node (2) is still down.

```
     STATE NNODES   NCORES    NGPUS NODELIST
      free      3       30        0 flux-sample-[0-1,3]
 allocated      0        0        0 
      down      7       70        0 flux-sample-[2,4-9]
```

## 3. Resize the Cluster Down

Now we are going to resize. We would hope that it sees that we have 3 pods, and doesn't delete another index.
Change the value of 4 to 3.

```diff
spec:
-  size: 4
+  size: 3
```
```
kubectl apply -f minicluster.yaml
```

And oh no! That's not what we want, it basically just deletes the last index. It doesn't check the number of pods we have running.

```bash
NAME                  READY   STATUS        RESTARTS   AGE
flux-sample-0-pdqzl   1/1     Running       0          3m58s
flux-sample-1-rmzhf   1/1     Running       0          3m58s
flux-sample-3-cfgwr   1/1     Terminating   0          3m58s
```
So now we had asked for size 3 but we actually have size 2.

```
NAME                  READY   STATUS        RESTARTS   AGE
flux-sample-0-pdqzl   1/1     Running       0          3m58s
flux-sample-1-rmzhf   1/1     Running       0          3m58s
```

If we look inside:

```
     STATE NNODES   NCORES    NGPUS NODELIST
      free      2       20        0 flux-sample-[0-1]
 allocated      0        0        0 
      down      8       80        0 flux-sample-[2-9]
```
That is what we would expect.

## 3. Resize the Cluster Up

Now let's add 2 nodes back, back up to size 4 from 3. We would expect the flux-sample-2 to be recreated but flux to still have it flagged as drained. But if the rule is hard that these indices cannot be recreated once they are gone, then we won't see it scale at all.

```diff
spec:
-  size: 3
+  size: 4
```
Ah interesting - so the flux-sample-2 pod is never re-created, but flux-sample-3 is! This tells us that once an index is killed (and the corresponding node drained for flux) it will never be possible to recreate. So we can't really keep track of actual sizes via that count, because the indices won't be recreated. 

```
NAME                  READY   STATUS     RESTARTS   AGE
flux-sample-0-pdqzl   1/1     Running    0          7m7s
flux-sample-1-rmzhf   1/1     Running    0          7m7s
flux-sample-3-f68rs   0/1     Init:0/1   0          4s
```

We are back in this state

```
     STATE NNODES   NCORES    NGPUS NODELIST
      free      3       30        0 flux-sample-[0-1,3]
 allocated      0        0        0 
      down      7       70        0 flux-sample-[2,4-9]
```

Let's do one more resize up, and to create indices we've never had before. I predict these will be created since we've never deleted / failed them.

```diff
spec:
-  size: 4
+  size: 6
```

Checks out!

```
NAME                  READY   STATUS     RESTARTS   AGE
flux-sample-0-pdqzl   1/1     Running    0          9m33s
flux-sample-1-rmzhf   1/1     Running    0          9m33s
flux-sample-3-f68rs   1/1     Running    0          2m30s
flux-sample-4-jpgnx   0/1     Init:0/1   0          2s
flux-sample-5-rbs4k   0/1     Init:0/1   0          2s
```

And inside.

```
     STATE NNODES   NCORES    NGPUS NODELIST
      free      5       50        0 flux-sample-[0-1,3-5]
 allocated      0        0        0 
      down      5       50        0 flux-sample-[2,6-9]
```

When you are done, clean up

```bash
kind delete cluster
```

## Summary

Here is what we learned:

- When the jobBackoffPerIndex is set to 0 (meaning the pod won't be recreated) this absolutely black lists the index from ever existing again for the job.
- Scaling down/up does not change that, those indices will never be re-created
- Scaling up is flexible to allow new indices that have not been created (and thus destroyed) yet
- Thus, we cannot rely on the MiniCluster size to indicate actual cluster size
  - We should only change it, increasing, when we need more nodes
  - There is never a need to downsize, it will erroneously delete a running index and make it too small
  - The TLDR is that we need to keep track of the actual size of the cluster via our own tool for now