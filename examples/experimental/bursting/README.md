# Bursting Experiment

> Experimental setup to burst to Google Cloud

I was reading about [external-dns](https://github.com/kubernetes-sigs/external-dns), 
and although this could be a more "production" way to set up bursting, we need something that
can be used on the fly (when we don't have a DNS name to use in some cloud). For this
experiment we are going to try defining a cluster service, nginx, that is on the same
headless network, and that (hopefully) we can forward as a port. Note that to start this
work, I started with work on the [nginx service](../../services/sidecar/nginx/) example.
I expect for this kind of bursting to work I will need to be able to:

 - Define the lead broker as a hostname external to the cluster AND internal (can we have two aliases?)
 - Have a custom external broker config that we generate
 - Be able to read in send the curve certificate as a configuration variable
 - Create a service to run (in place of the nginx container) to ensure the request goes to the right broker

We will eventually need to do the following:

 - Have an elegant way to decide:
   - When to burst
   - What jobs are marked for bursting, and how assigned to an external cluster
 - Unique external cluster names (within a cloud namespace) that are linked to their spec (content hash of broker.toml?)
 - When to configure external cluster to allow for scaling of flux (right now we set min == max so constant size)
 - Create a more scoped permission (Google service account) for running inside a cluster

We will use the following tricks to start (and each can be worked on to improve)

 - WRITE ME

Questions:

 - What would happen if we gave the lead broker two hosts with the same name?

## Usage

### Credentials

Since we are interacting with Google from within the MiniCluster, you need to have your default application credentials
shared there. This can probably be scoped to a service account, but for now we are being lazy. You must ABSOLUTELY
be sure you don't add these to git.

```bash
cp $HOME/.config/gcloud/application_default_credentials.json .
```

**DO NOT DO THIS FOR ANYTHING OTHER THAN A LOCAL TEST.**

### Setup Cluster

Create the cluster

```bash
$ kind create cluster --config ./kind-config.yaml
```

Note this config ensures that the cluster has an external IP on localhost.
Then install the operator, create the namespace, and the minicluster.

```bash
kubectl apply -f ../../dist/flux-operator-dev.yaml
kubectl create namespace flux-operator
```

Create an existing config map for nginx (that it will expect to be there, and we 
have defined in our minicluster.yaml under the nginx service existingVolumes).

```bash
$ kubectl apply -f service/nginx.yaml
```

TODO need to make portal to serve brokers from here...
And then create the MiniCluster

```bash
$ kubectl apply -f minicluster.yaml
```

And when the containers are ready, create the node port service for "flux-services":

```bash
$ kubectl apply -f service/service.yaml
```

Here is a nice block to easily copy paste all three :)

```bash
kubectl apply -f service/nginx.yaml
kubectl apply -f minicluster.yaml
kubectl apply -f service/service.yaml
```

### Burstable Job

Now let's create a job that cannot be run because we don't have the resources. In the future we would want some other logic to determine
this, but for now we are going to ensure the job doesn't run locally, and give it a label that our external application can sniff out and grab. Shell into the broker pod:

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-4r488 bash
```
Connect to the broker socket

```bash
$ sudo -u fluxuser -E $(env) -E HOME=/home/fluxuser flux proxy local:///run/flux/local bash
```

Resources we have available?

```bash
$ flux resource list
     STATE NNODES   NCORES NODELIST
      free      2        8 flux-sample-[0-1]
 allocated      0        0 
      down      8       32 flux-sample-[2-9]
```

And now let's create a burstable job, and ask for more nodes than we have :)

```bash
# Set burstable=1
# this will be for 4 nodes, 8 cores each
$ flux submit -N 4 --setattr=burstable hostname
```

You should see it's scheduled (but not running). Note that if we asked for a resource totally unknown
to the cluster (e.g. 4 nodes and 32 tasks) it would just fail so:

> TODO we need in our "mark as burstable" method a way to tell Flux not to fail in this case.

You can see it is scheduled and waiting for resources:

```bash
$ flux jobs -a
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
   ƒQURAmBXV fluxuser hostname    S      8      8        - 
```
```bash
$ flux job attach $(flux job last)
flux-job: ƒQURAmBXV waiting for resources  
```

Now we can run our script (which is bound locally in `/data` from the present working directory) to find the jobs based on this attribute!

```bash
GOOGLE_PROJECT=llnl-flux
python run-burst.py ${GOOGLE_PROJECT} --flux-operator-yaml ./external-config/flux-operator-dev.yaml
```

### Cleanup

Note that you'll only see the exposure with the kind docker container with `docker ps`.
When you are done, clean up

```bash
kubectl delete -f minicluster.yaml
kubectl delete -f nginx.yaml
kubectl delete -f service.yaml
```

kubectl apply -f nginx.yaml
kubectl apply -f service.yaml
kubectl apply -f minicluster.yaml


