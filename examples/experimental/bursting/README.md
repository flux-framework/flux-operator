# Bursting Experiment

> Experimental setup to burst to Google Cloud

I was reading about [external-dns](https://github.com/kubernetes-sigs/external-dns), 
and although this could be a more "production" way to set up bursting, we need something that
can be used on the fly (when we don't have a DNS name to use in some cloud). For this
experiment we are going to try defining a cluster service, nginx, that is on the same
headless network, and that (hopefully) we can forward as a port. Note that to start this
work, I started with work on the [nginx service](../../services/sidecar/nginx/) example.

## Design

### What we might need

I expect for this kind of bursting to work I will need to be able to:

 - Define the lead broker as a hostname external to the cluster AND internal (can we have two aliases?)
 - Have a custom external broker config that we generate
 - Be able to read in send the curve certificate as a configuration variable
 - Create a service to run (in place of the nginx container) to ensure the request goes to the right broker

### What we we eventually want to improve upon

We will eventually need to do the following:

 - Have an elegant way to decide:
   - When to burst
   - What jobs are marked for bursting, and how assigned to an external cluster
 - Unique external cluster names (within a cloud namespace) that are linked to their spec (content hash of broker.toml?)
 - When to configure external cluster to allow for scaling of flux (right now we set min == max so constant size)
 - Create a more scoped permission (Google service account) for running inside a cluster

### What we do now

- Allow the user to provide a custom broker.toml and curve.cert (the idea being the external clusters are created with a lead broker pointing to that service, and the same curve.cert that is on the local cluster they were created from) - these are changes to the operator CRD
- Submit jobs that are flagged for burstable (an attribute) and ask for more nodes than the cluster has (but below the max number so it doesn't immediately fail - we will want to fix this in flux because with bursting we should be able to ask for more than the potential size)
- Have a Python script that connects to the flux handle, finds the burstable jobs, and creates a minicluster spec (with the same nodes, tasks, and command - right now the machine is an argument)
- The Python script generates a broker.toml on the fly that defines the lead broker to be the service of the minicluster where it's running, and the curve.cert read directly from the filesystem (being used by the current cluster)
- The external cluster is created from the first, directly from the flux lead broker!
- Next - we need to figure out networking the two.

### Questions

 - What would happen if we gave the lead broker two hosts with the same name?
 - Is there ever interaction from any node aside from the lead broker of the second cluster?


## Credentials

Since we are interacting with Google from within the MiniCluster, you need to have your default application credentials
shared there. This can probably be scoped to a service account, but for now we are being lazy. You must ABSOLUTELY
be sure you don't add these to git.

```bash
cp $HOME/.config/gcloud/application_default_credentials.json .
```

**DO NOT DO THIS FOR ANYTHING OTHER THAN DEVELOPMENT OR TESTING.**

## Google Cloud Setup

Since we want to have two clusters communicating (and they will need public addresses) we will deploy both to GKE.
Let's create the first cluster:

```bash
CLUSTER_NAME=flux-cluster
GOOGLE_PROJECT=myproject
```
```bash
$ gcloud container clusters create ${CLUSTER_NAME} --project $GOOGLE_PROJECT \
    --zone us-central1-a --machine-type n2-standard-8 \
    --num-nodes=4 --enable-network-policy --tags=flux-cluster --enable-intra-node-visibility
```

And be sure to activate your credentials!

```bash
gcloud container clusters get-credentials ${CLUSTER_NAME}
```
Install the operator and create the minicluster

```bash
kubectl apply -f ../../dist/flux-operator-dev.yaml
kubectl create namespace flux-operator
# This is currently not used (but could be in the future)
kubectl apply -f service/nginx.yaml
kubectl apply -f minicluster.yaml
# Expose broker pod port 8050 to 30093
kubectl apply -f service/broker-service.yaml
```

We need to open up the firewall to that port - this creates the rule (you only need to do this once)

```bash
gcloud compute firewall-rules create flux-cluster-test-node-port --allow tcp:30093
```

Then figure out the node that the service is running from:

```bash
$ kubectl get pods -o wide -n flux-operator 
NAME                   READY   STATUS    RESTARTS   AGE     IP           NODE                                          NOMINATED NODE   READINESS GATES
flux-sample-0-kktl7    1/1     Running   0          7m22s   10.116.2.4   gke-flux-cluster-default-pool-4dea9d5c-0b0d   <none>           <none>
flux-sample-1-s7r69    1/1     Running   0          7m22s   10.116.1.4   gke-flux-cluster-default-pool-4dea9d5c-1h6q   <none>           <none>
flux-sample-services   1/1     Running   0          7m22s   10.116.0.4   gke-flux-cluster-default-pool-4dea9d5c-lc1h   <none>           <none>
```

Get the external ip for that node (for nginx, flux-services, and for the lead broker, a flux-sample-0-xx)

```bash
$ kubectl get nodes -o wide | grep gke-flux-cluster-default-pool-4dea9d5c-0b0d 
gke-flux-cluster-default-pool-4dea9d5c-0b0d   Ready    <none>   69m   v1.25.8-gke.500   10.128.0.83   34.171.113.254   Container-Optimized OS from Google   5.15.89+         containerd://1.6.18
```

If you applied [service/nginx.yaml](service/nginx.yaml) (a testing setup I used as a hello world for a service) you can try opening `34.135.221.11:30093`. For the broker (34.171.113.254) above, that might be harder to test. If you do the first, you should see the "Welcome to nginx!" page. Finally, when the broker index 0 pod is running, copy your scripts and configs over to it:

```bash
# This should be the index 0
POD=$(kubectl get pods -n flux-operator -o json | jq -r .items[0].metadata.name)

kubectl cp -n flux-operator ./run-burst.py ${POD}:/tmp/workflow/run-burst.py -c flux-sample

# We will find a better way than this
kubectl cp -n flux-operator ./application_default_credentials.json ${POD}:/tmp/workflow/application_default_credentials.json -c flux-sample

# Make directory
kubectl exec -it -n flux-operator ${POD} -- mkdir -p /tmp/workflow/external-config

# Copy configs
kubectl cp -n flux-operator ./external-config/flux-operator-dev.yaml ${POD}:/tmp/workflow/external-config/flux-operator-dev.yaml -c flux-sample
```

At this point, jump down to [burstable job](#burstable-job). If you are having trouble seeing communication for the service, you should come back to this step, and redo with nginx (service/service.yaml) and the flux-sample-service pod. 

## Development with Kind

I first tested on kind, and was able to get up to the point of needing to connect (and could not without the host).

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

## Burstable Job

Now let's create a job that cannot be run because we don't have the resources. In the future we would want some other logic to determine
this, but for now we are going to ensure the job doesn't run locally, and give it a label that our external application can sniff out and grab. Shell into the broker pod:

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-kvg5t bash
```

Connect to the broker socket

```bash
$ sudo -u flux -E $(env) -E HOME=/home/flux flux proxy local:///run/flux/local bash
```

Install the libraries we need:

```bash
# This is from /tmp/workflow
git clone -b add/gke-kubectl-client https://github.com/converged-computing/kubescaler
cd kubescaler
python3 -m pip install -e .[all]
cd -
git clone --depth 1 -b bursting https://github.com/flux-framework/flux-operator ./op
cd ./op/sdk/python/v1alpha1
python3 -m pip install -e .
cd -
python3 -m pip install IPython
```

We can eventually package these in a container base.
Resources we have available?

```bash
$ flux resource list
     STATE NNODES   NCORES NODELIST
      free      2        8 flux-sample-[0-1]
 allocated      0        0 
      down      6       24 flux-sample-[2-3],burst-0-[0-3]
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

Get a variant of the munge key we can see:

```bash
sudo cp /etc/munge/munge.key ./munge.key
sudo chown $USER munge.key
```

Now we can run our script to find the jobs based on this attribute!

```bash
GOOGLE_PROJECT=myproject
LEAD_HOST="34.28.41.20"
LEAD_PORT=30093
python3 run-burst.py --project ${GOOGLE_PROJECT} --cluster-name flux-external-cluster --flux-operator-yaml ./external-config/flux-operator-dev.yaml \
        --lead-host ${LEAD_HOST} --lead-port ${LEAD_PORT} --lead-size 4 \
        --munge-key ./munge.key --name burst-0
```

Important notes for the above:

- The name is same that would be automatically generated name by Flux given a bursted cluster (that isn't explicitly given a name) but we are being pedantic. It's also in the [minicluster.yaml](minicluster.yaml)
- The lead name is derived from the hostname where it is running (e.g., flux-sample) so we don't need to provide it
- We set the lead size to the max size, because the ranks indices need to line up. We are using a size that won't fail the job (which needs 4)

When you do the above (and the second MiniCluster launches) you should be able to see on your local cluster the external
MiniCluster resources, and the result of hostname will include the external hosts!

```bash
flux@flux-sample-0:/tmp/workflow$ flux resource list
     STATE NNODES   NCORES NODELIST
      free      6       24 flux-sample-[0-1],burst-0-[0-3]
 allocated      0        0 
      down      2        8 flux-sample-[2-3]
```
Our job has run:

```bash
flux@flux-sample-0:/tmp/workflow$ flux jobs -a 
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
  ƒ2S9iZpoJw flux     hostname   CD      4      4   0.027s flux-sample-1,burst-0-[0,2-3]
   ƒ3i2dgyDq flux     hostname   CD      4      4   0.035s flux-sample-[0-1,3],burst-0-0
```

And we can see output! Note that the error is because the working directory where it was launched doesn't exist on the remote.

```bash
flux@flux-sample-0:/tmp/workflow$ flux job attach ƒ2S9iZpoJw
214.588s: flux-shell[2]: ERROR: host burst-0-2: Could not change dir to /tmp/workflow: No such file or directory. Going to /tmp instead
214.589s: flux-shell[3]: ERROR: host burst-0-3: Could not change dir to /tmp/workflow: No such file or directory. Going to /tmp instead
214.588s: flux-shell[1]: ERROR: host burst-0-0: Could not change dir to /tmp/workflow: No such file or directory. Going to /tmp instead
flux-sample-1
burst-0-0
burst-0-3
burst-0-2
```

You can also launch a new job that asks to hit all the nodes (6):

```bash
flux@flux-sample-0:/tmp/workflow$ flux run -N 6 hostname
0.039s: flux-shell[5]: ERROR: host burst-0-3: Could not change dir to /tmp/workflow: No such file or directory. Going to /tmp instead
0.038s: flux-shell[4]: ERROR: host burst-0-2: Could not change dir to /tmp/workflow: No such file or directory. Going to /tmp instead
0.040s: flux-shell[3]: ERROR: host burst-0-1: Could not change dir to /tmp/workflow: No such file or directory. Going to /tmp instead
0.038s: flux-shell[2]: ERROR: host burst-0-0: Could not change dir to /tmp/workflow: No such file or directory. Going to /tmp instead
flux-sample-0
burst-0-2
burst-0-1
burst-0-3
burst-0-0
flux-sample-1
```

And that's bursting! At this point we should think about how to better start / stop a burst, since the cluster
will typically come up (and stay up).

## Debugging

### Mismatch of zeromq version

Note that if you see a bug like this - the issue is a mismatch of zeromq versions between containers.

```console
broker.debug[0]: child sockevent tcp://10.116.2.8:8050 unknown socket event
broker.debug[0]: child sockevent tcp://10.116.2.8:8050 disconnected
broker.debug[0]: child sockevent tcp://10.116.2.8:8050 accepted
broker.debug[0]: child sockevent tcp://10.116.2.8:8050 unknown socket event
```

and the external cluster:

```console
broker.err[1]: parent sockevent tcp://34.171.113.254:30093 handshake failed: Broken pipe
broker.debug[1]: parent sockevent tcp://34.171.113.254:30093 disconnected
broker.debug[1]: parent sockevent tcp://34.171.113.254:30093 connect retried
```
### Kubectl for External

Note that it's possible to get the kubectl config for you external cluster, and I recommended running
this from the Google Cloud console (shell) so you don't muck around with your default kubectl:

```bash
$  gcloud container clusters get-credentials flux-cluster --zone us-central1-a --project llnl-flux \
>  && kubectl get service kubernetes -o yaml
```

If you click on the job name, there should be a "kubectl" link at the top.

### Cleanup

Note that you'll only see the exposure with the kind docker container with `docker ps`.
When you are done, clean up

```bash
kubectl delete -f minicluster.yaml
kubectl delete -f nginx.yaml
kubectl delete -f service.yaml
gcloud container clusters delete flux-cluster
```


