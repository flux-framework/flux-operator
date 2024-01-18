# Downsize

In this example, we will create a MiniCluster and request that workers do not restart that fail. We expose this as a boolean `completetWorkers` (defaults to false) We want to test:

- When we drain and delete a pod, it is not recreated
- What happens when we resize the cluster (does the index get recreated)?

## Cluster

Create the cluster with the feature gate.

```bash
kind create cluster
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
kubectl exec -it flux-sample-0-9zsjw bash
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
flux resource drain flux-sample-2 "Goodbye rank 2"
```

It should be down.

```console
     STATE NNODES   NCORES    NGPUS NODELIST
      free      3       30        0 flux-sample-[0-1,3]
 allocated      0        0        0 
      down      7       70        0 flux-sample-[2,4-9]
```

## 2. Kill the Worker

We previously had killed a pod (and then it wasn't restarted):

```bash
# DO NOT RUN THIS
kubectl delete pod flux-sample-2-xxx
```

But there is a better way! We can just kill the broker and have the pod not complete.
Note that this should also be done from the lead broker. Note that we are still inside the lead broker's
flux instance. Sanity check that rank 2 `-r 2` is `flux-sample-2`

```bash
flux exec -r 2 hostname
flux-sample-2
```

Get the flux broker pid for rank 2:

```bash
flux exec -r 2 flux getattr broker.pid
47
```

Sanity check

```
poorUnfortunateSoul=$(flux exec -r 2 flux getattr broker.pid)
[root@flux-sample-0 /]# echo $poorUnfortunateSoul 
47
```

And then kill the flux broker. Note that this is for Ubuntu and we are going to use [SIGTERM](https://docs.rockylinux.org/books/admin_guide/08-process/#process-management-controls)
to nicely kill it.

```bash
# This is for Ubuntu that has pkill
flux exec -r 2 pkill -15 flux-broker

# This is for rockylinux
flux exec -r 2 kill -15 ${poorUnfortunateSoul}

# NOT RECOMMENDED: if you can't kill or pkill, but gives the follower no chance to cleanup
flux overlay disconnect 2
```

On the inside, the worker will still be down. From the outside, you should see the pod Complete.

```console
 kubectl get pods
NAME                  READY   STATUS      RESTARTS   AGE
flux-sample-0-2fw2l   1/1     Running     0          10m
flux-sample-1-zh5k7   1/1     Running     0          10m
flux-sample-2-njp8d   0/1     Completed   0          10m
flux-sample-3-rcngx   1/1     Running     0          10m
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
flux-sample-0-2fw2l   1/1     Running       0          11m
flux-sample-1-zh5k7   1/1     Running       0          11m
flux-sample-2-njp8d   0/1     Completed     0          11m
flux-sample-3-rcngx   1/1     Terminating   0          11m
```

So now we had asked for size 3 but we actually have size 2.

```
NAME                  READY   STATUS      RESTARTS   AGE
flux-sample-0-2fw2l   1/1     Running     0          12m
flux-sample-1-zh5k7   1/1     Running     0          12m
flux-sample-2-njp8d   0/1     Completed   0          12m
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

```console
NAME                  READY   STATUS      RESTARTS   AGE
flux-sample-0-2fw2l   1/1     Running     0          12m
flux-sample-1-zh5k7   1/1     Running     0          12m
flux-sample-2-njp8d   0/1     Completed   0          12m
flux-sample-3-hlvrn   0/1     Init:0/1    0          3s
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
NAME                  READY   STATUS      RESTARTS   AGE
flux-sample-0-2fw2l   1/1     Running     0          14m
flux-sample-1-zh5k7   1/1     Running     0          14m
flux-sample-2-njp8d   0/1     Completed   0          14m
flux-sample-3-hlvrn   1/1     Running     0          94s
flux-sample-4-bt2cl   0/1     Init:0/1    0          3s
flux-sample-5-ccjsm   0/1     Init:0/1    0          3s
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

One idea I just had - if we had support in indexed jobs to explicitly ask for a pod to restart, this might be an option to allow Completed -> Running. This is not supported currently.

## Summary

Here is what we learned:

- When we allow a worker to complete (with exit code 0) the pod will not restart. We don't need the `backoffLimitPerIndex` set to 0 as we tried before.
- Scaling down/up does not change that, those indices will never be re-created, they will stay completed
- Scaling up is flexible to allow new indices that have not been created (and thus destroyed) yet
- Thus, we cannot rely on the MiniCluster size to indicate actual cluster size
  - We should only change it, increasing, when we need more nodes
  - There is never a need to downsize, it will erroneously delete a running index and make it too small
  - The TLDR is that we need to keep track of the actual size of the cluster via our own tool for now
