# Multiple Pods Per Nodes

This is an example of how to create a cluster with several pods per node, or limiting a physical node's resources (in terms of what a flux node sees).
This was tested on a Google Cloud Kubernetes Engine cluster and has not been tested locally with kind, etc. We use this instance type:

 - [c2d-standard-112](https://cloud.google.com/compute/docs/compute-optimized-machines#c2d_machine_types)

## c2d-standard-112

### Create the Cluster

Let's test a cluster on c2d-standard-112 for size 4.

```bash
GOOGLE_PROJECT=myproject
gcloud container clusters create test-cluster \
    --threads-per-core=1 \
    --num-nodes=4 \
    --region=us-central1-a \
    --project=${GOOGLE_PROJECT} \
    --machine-type=c2d-standard-112 \
    --placement-type=COMPACT \
    --system-config-from-file=./cluster-config.yaml
```

**IMPORTANT** we will not be able to ask for granularity of resources with a special [cluster-config.yaml](cluster-config.yaml) that specifies we want a static CPU manager policy [see here](https://kubernetes.io/docs/tasks/administer-cluster/cpu-management-policies/#static-policy). And yes, I tested with and without this config file to sanity check. This was the piece that I was missing before - Antonio had added it (and it went over my head at the time) but I remembered the file and started looking into "the other stuff that can be defined there" and stumbled on that documentation.

### Install the Flux Operator

As follows:

```bash
kubectl apply -f https://raw.githubusercontent.com/flux-framework/flux-operator/main/examples/dist/flux-operator.yaml
```

### Experiments

Then create the flux operator pods. 

#### 4 pod "nodes"

For our first effort, we will ask for 4 pods, but limit resources (and should validate
that we can only see a limited subset).

```bash
kubectl apply -f minicluster-4.yaml
```

Shell in

```bash
kubectl exec -it flux-sample-0-xxx bash
```

and then

```bash
source /mnt/flux/flux-view.sh 
flux proxy $fluxsocket bash
flux resource list
```
```console
[root@flux-sample-0 /]# flux resource list
     STATE NNODES   NCORES    NGPUS NODELIST
      free      4      160        0 flux-sample-[0-3]
 allocated      0        0        0 
      down      0        0        0 
```

😱️ IT WORKS! OMG!

#### 7 "nodes" on 4 physical nodes

Do the same above, but with minicluster-7.yaml

```bash
kubectl apply -f minicluster-7.yaml
```

Shell in

```bash
kubectl exec -it flux-sample-0-xxx bash
```

and then

```bash
source /mnt/flux/flux-view.sh 
flux proxy $fluxsocket bash
flux resource list
```
```console
[root@flux-sample-0 /]# flux resource list
     STATE NNODES   NCORES    NGPUS NODELIST
      free      7      140        0 flux-sample-[0-6]
 allocated      0        0        0 
      down      0        0        0 
```

Note that there is a HUGE gotcha in here - if you give a wrong number above, it will still create (but fall back to `Burstable``) and flux will see the wrong number of resources. If you ask for too little memory, it can also get OOMKilled. I asked for 7 nodes instead of 8 because I always had one killed. So:

- Always check the pod is Guaranteed
- Ensure none of your pods are OOMKilled (init containers can have this happen too)
- Always check flux resource list

### Clean Up

When you are done:

```bash
gcloud container clusters delete test-cluster --region=us-central1-a --quiet
```
