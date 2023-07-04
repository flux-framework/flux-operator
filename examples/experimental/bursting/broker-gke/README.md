# Bursting Experiment to GKE

> Experimental setup to burst to Google Cloud

This setup will expose a lead broker (index 0 of the MiniCluster job) as a service,
and then deploy a second cluster that can connect back to the first. For a different
design that could be used to do similar but via a central router (less developed)
start with [nginx](../nginx). For the overall design, see the top level [README](../README.md)

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
$ gcloud container clusters get-credentials ${CLUSTER_NAME}
```

Create the namespace, install the operator (assuming you are using a development version) and create the minicluster.
Note that if you aren't using a development version, you can apply `flux-operator.yaml` instead.

```bash
kubectl apply -f ../../../dist/flux-operator-dev.yaml
kubectl create namespace flux-operator
kubectl apply -f minicluster.yaml
# Expose broker pod port 8050 to 30093
kubectl apply -f service/broker-service.yaml
```

We need to open up the firewall to that port - this creates the rule (you only need to do this once)

```bash
gcloud compute firewall-rules create flux-cluster-test-node-port --allow tcp:30093
```

Then figure out the node that the service is running from (we are interested in lead broker flux-sample-0-*)

```bash
$ kubectl get pods -o wide -n flux-operator 
NAME                   READY   STATUS    RESTARTS   AGE     IP           NODE                                          NOMINATED NODE   READINESS GATES
flux-sample-0-kktl7    1/1     Running   0          7m22s   10.116.2.4   gke-flux-cluster-default-pool-4dea9d5c-0b0d   <none>           <none>
flux-sample-1-s7r69    1/1     Running   0          7m22s   10.116.1.4   gke-flux-cluster-default-pool-4dea9d5c-1h6q   <none>           <none>
flux-sample-services   1/1     Running   0          7m22s   10.116.0.4   gke-flux-cluster-default-pool-4dea9d5c-lc1h   <none>           <none>
```

Then (using that node name) get the external ip for that node (for nginx, flux-services, and for the lead broker, a flux-sample-0-xx)

```bash
$ kubectl get nodes -o wide | grep gke-flux-cluster-default-pool-4dea9d5c-0b0d 
gke-flux-cluster-default-pool-4dea9d5c-0b0d   Ready    <none>   69m   v1.25.8-gke.500   10.128.0.83   34.171.113.254   Container-Optimized OS from Google   5.15.89+         containerd://1.6.18
```

I set it to the environment to be useful later:

```bash
export LEAD_BROKER_HOST=34.171.113.254
```

Finally, when the broker index 0 pod is running, copy your scripts and configs over to it:

```bash
# This should be the index 0
POD=$(kubectl get pods -n flux-operator -o json | jq -r .items[0].metadata.name)

# This will copy configs / create directories for it
kubectl cp -n flux-operator ./run-burst.py ${POD}:/tmp/workflow/run-burst.py -c flux-sample
kubectl cp -n flux-operator ./application_default_credentials.json ${POD}:/tmp/workflow/application_default_credentials.json -c flux-sample
kubectl exec -it -n flux-operator ${POD} -- mkdir -p /tmp/workflow/external-config
kubectl cp -n flux-operator ../../../dist/flux-operator-dev.yaml ${POD}:/tmp/workflow/external-config/flux-operator-dev.yaml -c flux-sample
```

## Burstable Job

Now let's create a job that cannot be run because we don't have the resources. The `flux-burst` Python module, using it's simple
default, will just look for jobs with `burstable=True` and then look for a place to assign them to burst. Since this is a plugin
framework, in the future we can implement more intelligent algorithms for either filtering the queue (e.g., "Which jobs need bursting?"
and then determining if a burst can be scheduled for some given burst plugin (e.g., GKE)). For this simple setup and example,
we ensure the job doesn't run locally because we've asked for more nodes than we have. Shell into your broker pod:

```bash
$ kubectl exec -it -n flux-operator ${POD} bash
```

Connect to the broker socket. If this issues an error, it's likely the install scripts are still running (you can check
the logs and wait a minute!)

```bash
$ sudo -u flux -E $(env) -E HOME=/home/flux flux proxy local:///run/flux/local bash
```

The libraries we need should be installed in the minicluster.yaml.
You might want to add others for development (e.g., IPython).
Resources we have available?

```bash
$ flux resource list
     STATE NNODES   NCORES NODELIST
      free      2        8 flux-sample-[0-1]
 allocated      0        0 
      down      6       24 flux-sample-[2-3],burst-0-[0-3]
```

The above shows us that the broker running here can accept burstable resources (`burst-0-[0-3]`), and even
can accept the local cluster expanding (`flux-sample[2-3]`) for a total of 24 cores. The reason
that the remote burst prefix has an extra "0" is that we could potentially have different sets of
burstable remotes, namespaced by this prefix. And now let's create a burstable job, and ask for more nodes than we have :)

```bash
# Set burstable=1
# this will be for 4 nodes, 8 cores each
$ flux submit -N 4 --cwd /tmp --setattr=burstable hostname
```

You should see it's scheduled (but not running). Note that if we asked for a resource totally unknown
to the cluster (e.g. 4 nodes and 32 tasks) it would just fail. Note that because of this,
we need in our "mark as burstable" method a way to tell Flux not to fail in this case.
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

Get a variant of the munge key we can see (it's owned by root so this ensures we can see/own it as the flux user)

```bash
sudo cp /etc/munge/munge.key ./munge.key
sudo chown $USER munge.key
```

Now we can run our script to find the jobs based on this attribute!

```bash
# Our Google Project name
GOOGLE_PROJECT=myproject

# This is the address of the lead host we discovered above
LEAD_HOST="34.72.223.15"

# This is the node port we've exposed on the cluster
LEAD_PORT=30093
python3 run-burst.py --project ${GOOGLE_PROJECT} --flux-operator-yaml ./external-config/flux-operator-dev.yaml \
        --lead-host ${LEAD_HOST} --lead-port ${LEAD_PORT} --lead-size 4 \
        --munge-key ./munge.key --name burst-0
```

Important notes for the above:

- The curve path and secret name have defaults set.
- The name is same that would be automatically generated name by Flux given a bursted cluster (that isn't explicitly given a name) but we are being pedantic. It's also in the [minicluster.yaml](minicluster.yaml)
- The lead name is derived from the hostname where it is running (e.g., flux-sample) so we don't need to provide it
- We set the lead size to the max size, because the ranks indices need to line up. We are using a size that won't fail the job (which needs 4)
- mock mode is set to false for the `FluxBurst` client, meaning the cluster will attempt to connect to our first one

When you do the above (and the second MiniCluster launches) you should be able to see on your local cluster the external
MiniCluster resources, and the result of hostname will include the external hosts! Here is how to shell into the cluster
from another terminal:

```bash
$ POD=$(kubectl get pods -n flux-operator -o json | jq -r .items[0].metadata.name)
$ kubectl exec -it -n flux-operator ${POD} bash
$ sudo -u flux -E $(env) -E HOME=/home/flux flux proxy local:///run/flux/local bash
```

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
flux-sample-0
flux-sample-1
burst-0-2
burst-0-0
```

You can also launch a new job that asks to hit all the nodes (6):

```bash
flux@flux-sample-0:/tmp/workflow$ flux run -N 6 --cwd /tmp hostname
flux-sample-0
burst-0-3
burst-0-2
burst-0-1
burst-0-0
flux-sample-1
```

Note that without `--cwd` you would see an error that the bursted cluster can't CD to `/tmp/workflow` (it doesn't exist there).
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