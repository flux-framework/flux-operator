# Bursting Experiment to Compute Engine

> Experimental setup to burst to Google Cloud Compte Engine

This setup will expose a lead broker (index 0 of the MiniCluster job) as a service,
and then deploy a second cluster that can connect back to the first. However unlike
the other examples that burst to another Kubernetes cluster, this burst goes to Compute
Engine, and is driven by [flux-burst-compute-engine](https://github.com/converged-computing/flux-burst-compute-engine). 
For the overall design, see the top level [README](../README.md). In summary, for this setup:

1. The main cluster will be run on GKE
2. The bursted cluster will be on Compute Engine

This is a more complex setup because it requires not just the terraform configs provided by
flux-burst-compute-engine, but also the image built from [flux-terraform-gcp](https://github.com/converged-computing/flux-terraform-gcp/tree/main/build-images/bursted).
(the repository that hosts the terraform modules). 

### What should be the same?

During this setup, we learned that the following must be the same (or be available) for the bust to fully work.

 - The flux user id must match between two instances (e.g., here we use 1004, built into VMs and set here)
 - The flux lib directory (e.g., `/usr/lib` and `/usr/lib64`) should match (you'll probably be OK installing on same OS with same method)
 - The flux install location should be the same (e.g., `/usr` and `/usr/local` will have an error)


## Build Machine Image with Flux

You can prepare that image as follows:

```bash
git clone https://github.com/converged-computing/flux-terraform-gcp
cd build-images/basic
make bursted
```

Unlike the "basic" setup in that same respository, this is one simple image that includes Flux,
and expects customization to happen via the startup script. This allows for fewer images to maintain,
and less chance of needing to update the base image build.

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
    --zone us-central1-a --machine-type n2-standard-4 \
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
$ kubectl get pods -o wide
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

Take note of this ip address for later - we will need it for running the bursting script.
Finally, when the broker index 0 pod is running, copy your scripts and configs over to it:

```bash
# This should be the index 0
POD=$(kubectl get pods -o json | jq -r .items[0].metadata.name)

# This will copy configs / create directories for it
kubectl cp ./run-burst.py ${POD}:/tmp/workflow/run-burst.py -c flux-sample
kubectl cp ./application_default_credentials.json ${POD}:/tmp/workflow/application_default_credentials.json -c flux-sample
kubectl exec -it ${POD} -- mkdir -p /tmp/workflow/external-config
kubectl cp ../../../dist/flux-operator-dev.yaml ${POD}:/tmp/workflow/external-config/flux-operator-dev.yaml -c flux-sample
```

## Burstable Job

Now let's create a job that cannot be run because we don't have the resources. The `flux-burst` Python module, using it's simple
default, will just look for jobs with `burstable=True` and then look for a place to assign them to burst. Since this is a plugin
framework, in the future we can implement more intelligent algorithms for either filtering the queue (e.g., "Which jobs need bursting?"
and then determining if a burst can be scheduled for some given burst plugin (e.g., GKE)). For this simple setup and example,
we ensure the job doesn't run locally because we've asked for more nodes than we have. Shell into your broker pod:

```bash
$ kubectl exec -it ${POD} bash
```

Connect to the broker socket. If this issues an error, it's likely the install scripts are still running (you can check
the logs and wait a minute!)

```bash
source /mnt/flux/flux-view.sh
flux proxy local:///run/flux/local bash
```

The libraries we need should be installed in the minicluster.yaml.
You might want to add others for development (e.g., IPython).
Resources we have available?

```bash
$ flux resource list
     STATE NNODES   NCORES NODELIST
      free      2        8 flux-sample-[0-1]
 allocated      0        0 
      down      6       24 flux-sample-[2-3],gffw-compute-a-[001-003]
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
Also note that once it's assigned to a plugin to be bursted, it will lose that attribute
(and note be able to be scheduled again). You can see it is scheduled and waiting for resources:

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
cp /etc/munge/munge.key ./munge.key
```

Now we can run our script to find the jobs based on this attribute!

```bash
# Our Google Project name
GOOGLE_PROJECT=myproject

# This is the address of the lead host we discovered above
LEAD_HOST="35.202.211.23"

# Note that the lead host will be added here as a prefix
hostnames="flux-sample-[1-3],gffw-compute-a-[001-003]"

# This is the node port we've exposed on the cluster
LEAD_PORT=30093
python3 run-burst.py --project ${GOOGLE_PROJECT} \
        --lead-host ${LEAD_HOST} --lead-port ${LEAD_PORT} --lead-hostnames ${hostnames} \
        --munge-key ./munge.key --curve-cert /mnt/curve/curve.cert
```

When you do the above you'll see the terraform configs apply, and the second Flux cluster will be launched when they finish. 
You'll then be prompted to press ENTER when you want to destroy the burst. This is when you can open another terminal
to see the outcome. Here is how to shell into the cluster from another terminal:

```bash
POD=$(kubectl get pods -o json | jq -r .items[0].metadata.name)
kubectl exec -it ${POD} bash
source /mnt/flux/flux-view.sh
flux proxy $fluxsocket bash
```

Resources are now allocated:

```bash
$ flux resource list
     STATE NNODES   NCORES NODELIST
      free      5       10 flux-sample-[0-1],gffw-compute-a-[001-003]
 allocated      0        0 
      down      2        4 flux-sample-[2-3]
```
Our job has run:

```bash
$ flux jobs -a
```
```console
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
    ƒCJTuUPR flux     hostname   CD      4      4   0.623s flux-sample-1,gffw-compute-a-[001-003]
```
And we can see the result with hostnames from the local and bursted cluster.
Note that we get an error about resources (I think) because we haven't done any work to ensure they are correct.
This is probably OK for now and we will need to tweak further to allow the template to include them.

```bash
$ flux job attach ƒXmSxZFm
```
```console
1352.660s: flux-shell[3]:  WARN: rlimit: nofile exceeds current max, raising value to hard limit
1352.660s: flux-shell[3]:  WARN: rlimit: nproc exceeds current max, raising value to hard limit
1352.634s: flux-shell[1]:  WARN: rlimit: nofile exceeds current max, raising value to hard limit
1352.635s: flux-shell[1]:  WARN: rlimit: nproc exceeds current max, raising value to hard limit
1352.627s: flux-shell[2]:  WARN: rlimit: nofile exceeds current max, raising value to hard limit
1352.627s: flux-shell[2]:  WARN: rlimit: nproc exceeds current max, raising value to hard limit
flux-sample-1
gffw-compute-a-003
gffw-compute-a-002
gffw-compute-a-001
```

You can also launch a new job to see it interactively:

```bash
$ flux run -N 6 --cwd /tmp hostname
...
flux-sample-0
gffw-compute-a-002
gffw-compute-a-003
gffw-compute-a-001
```

Note that without `--cwd` you would see an error that the bursted cluster can't CD to `/tmp/workflow` (it doesn't exist there).
And that's bursting! And WOW did we learn a lot by using different operating systems!

## Debugging

### Checking Startup Script

I had to debug these quite a bit, and I found the following helpful. First, shell into a compute node
(you can usually copy paste this from the Google Cloud console)

```bash
gcloud compute ssh --zone "us-central1-a" "gffw-compute-a-003" --tunnel-through-iap --project "${GOOGLE_PROJECT}"
```

And then to see if the startup script ran:

```bash
# Get logs from journalctl
sudo journalctl -u google-startup-scripts.service

# Try running it again
sudo google_metadata_script_runner startup
```

I would put an IPython.embed right after the unmatched line in the [run-burst.py](run-burst.py)
script, and then you can get the plugin and get the content of the script to debug:

```bash
plugin = client.plugins["compute_engine"]
print(plugin.params.compute_boot_script)
```

### Checking Service

On the instance, you can check the status of the service (and see where that script is too.)

```bash
$ sudo systemctl status flux-start.service
```
```console
● flux-start.service - Flux message broker
   Loaded: loaded (/etc/systemd/system/flux-start.service; enabled; vendor preset: disabled)
   Active: active (running) since Mon 2023-07-10 00:13:30 UTC; 34s ago
 Main PID: 5050 (flux-broker-6)
    Tasks: 10 (limit: 100606)
   Memory: 58.2M
   CGroup: /system.slice/flux-start.service
           └─5050 broker --config-path /usr/etc/flux/system/conf.d -Scron.directory=/usr/etc/flux/system/conf.d -Stbon.fanout=256 -Srundir=/run/flux -Sbroker.rc2_no>

Jul 10 00:13:31 gffw-compute-a-003 flux[5050]: broker.debug[6]: insmod resource
Jul 10 00:13:31 gffw-compute-a-003 flux[5050]: broker.debug[6]: insmod job-info
Jul 10 00:13:32 gffw-compute-a-003 flux[5050]: broker.debug[6]: insmod job-ingest
Jul 10 00:13:32 gffw-compute-a-003 flux[5050]: job-ingest.debug[6]: configuring validator with plugins=(null), args=(null) (enabled)
Jul 10 00:13:32 gffw-compute-a-003 flux[5050]: job-ingest.debug[6]: fluid ts=1421183ms
Jul 10 00:13:32 gffw-compute-a-003 flux[5050]: broker.info[6]: rc1.0: running /etc/flux/rc1.d/01-sched-fluxion
Jul 10 00:13:32 gffw-compute-a-003 flux[5050]: broker.info[6]: rc1.0: running /etc/flux/rc1.d/02-cron
Jul 10 00:13:32 gffw-compute-a-003 flux[5050]: broker.info[6]: rc1.0: /etc/flux/rc1 Exited (rc=0) 0.9s
Jul 10 00:13:32 gffw-compute-a-003 flux[5050]: broker.info[6]: rc1-success: init->quorum 0.936237s
Jul 10 00:13:32 gffw-compute-a-003 flux[5050]: broker.info[6]: quorum-full: quorum->run 0.063886ms
```

That should be running as the flux user.

### Accidentally Exit

If you exit the run burst script and haven't cleaned up, you do so manually.

 - Delete the instances generated, including compute and the ns node.
 - Then search for cloud routers, and there should be one associated with the foundation network
 - Then under networks, find `foundation-net`
  - Delete the subnet first
  - Then delete the entire VPC

### Cleanup

Note that you'll only see the exposure with the kind docker container with `docker ps`.
When you are done, clean up

```bash
kubectl delete -f minicluster.yaml
kubectl delete -f nginx.yaml
kubectl delete -f service.yaml
gcloud container clusters delete flux-cluster
```
