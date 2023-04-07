# Fusion Storage

This basic tutorial will walk through creating a MiniCluster to run a Snakemake workflow! You should have
already [setup your workspace](setup.md), including preparing the Snakemake data in
Google Storage. The cluster creation commands are slightly modified here, discussed below.

## Create Cluster

Now let's use gcloud to create a cluster, and we are purposefully going to choose
a very small node type to test. Note that I choose us-central1-a because it tends
to be cheaper (and closer to me). First, here is our project for easy access:

```bash
GOOGLE_PROJECT=myproject
```

Replace the above with your project name, of course! If you are doing the fusion demo, it
has a [few requirements](https://github.com/seqeralabs/wave-showcase/tree/master/example-gke):
Here is how we can create our cluster:

```bash
$ gcloud container clusters create flux-cluster --project $GOOGLE_PROJECT \
    --zone us-central1-a --machine-type n1-standard-2 --cluster-version 1.25 \
    --num-nodes=2 --enable-network-policy --tags=flux-cluster --enable-intra-node-visibility \
    --ephemeral-storage-local-ssd count=1 --workload-pool=${GOOGLE_PROJECT}.svc.id.goog \
    --workload-metadata=GKE_METADATA
```

 - The "clusters create" command is for a standard cluster
 - The n1-standard-2 has two vCPU ([see this page](https://cloud.google.com/compute/docs/general-purpose-machines))
 - We create ssd storage with `--ephemeral-storage-local-ssd` and the count is the [number per node](https://cloud.google.com/kubernetes-engine/docs/how-to/persistent-volumes/local-ssd)
 - We enable [workfload identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity) via ` --workload-pool`
 - We enable the metadata server via `--workload-metadata`

Note for the ssd, we need to use Kubernetes 1.25 or greater. Then get credentials:

```bash
$ gcloud container clusters get-credentials flux-cluster --zone us-central1-a --project $GOOGLE_PROJECT
```

Create the "flux-operator" namespace:

```bash
$ kubectl create namespace flux-operator
```

Create a service account:

```bash
$ kubectl create serviceaccount flux-operator-sa --namespace flux-operator
```

At this point you can either use an existing Google Cloud service account (different from the above, and shown below)
or [follow instructions to create a new one](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#gcloud_3).
For step 5 - I usually use the Google Cloud [interface for roles](https://console.cloud.google.com/iam-admin/iam) 
to ensure my service account has all that are needed. Next, we need to connect the kubernetes service account
that we created about to our Google Cloud service account:

```bash
GOOGLE_SERVICE_ACCOUNT=GSA_NAME@GSA_PROJECT.iam.gserviceaccount.com
NAMESPACE="flux-operator"
KSA_NAME="flux-operator-sa"
gcloud iam service-accounts add-iam-policy-binding ${GOOGLE_SERVICE_ACCOUNT} \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:${GOOGLE_PROJECT}.svc.id.goog[${NAMESPACE}/${KSA_NAME}]"
```

We then "annotate" the Kubernetes service account with the email address of the IAM service account.

```bash
kubectl annotate serviceaccount ${KSA_NAME} \
    --namespace ${NAMESPACE} \
    iam.gke.io/gcp-service-account=${GOOGLE_SERVICE_ACCOUNT}
```

Basically, this means our Kubernetes service account inherits the Google Service account.
Next, you'll want to choose a mens to [install the operator](https://flux-framework.org/flux-operator/getting_started/user-guide.html#install) 
and then continue.

## Running Snakemake

Fusion is cool because it binds from the pod. We will first need to apply a Daemonset that enables a device
manager so fuse to be used:

```bash
$ kubectl apply -f examples/storage/google/fusion/daemonset.yaml
```

A few things I learned for this daemonset! If you create it in a namespace other than kube-system (and have the
priority class defined) you'll get an error about quotas. The fix (for now) is to deploy to the kube-system
namespace. I also adjusted the config.yaml at the bottom to equal the number of nodes in my cluster (2 for this demo)
and I'm not sure how much that matters (the previous setting was at 20). 

And then label the nodes so they can use the manager:

```bash
for n in $(kubectl get nodes | tail -n +2 | cut -d' ' -f1 ); do
    kubectl label node $n smarter-device-manager=enabled
done 
```

After this, you should have two daemonset pods running:

```bash
$ kubectl get -n kube-system pods | grep device
```
```console
smarter-device-manager-qv7xf                             1/1     Running   0          39m
smarter-device-manager-trtbl                             1/1     Running   0          39m
```

You can do some jq fu to verify they are there:

```bash
$ kubectl get nodes -o json | jq ".items[].metadata.labels" | grep smarter
  "smarter-device-manager": "enabled",
  "smarter-device-manager": "enabled",
```

Now we will create the MiniCluster. Note that:

 - It's annotated both to use the Google Service account and be able to use the fuse device!
 - We install fusion and mount at default /fusion in a commands->pre block
 - We have to run flux as root / the container as privileged for fuse.
 - pod resource requests ask for the fuse mount
 - The organization of a mount is based on the service (e.g., `/fusion/s3` or `/fusion/gs` and both will appear with `ls`)

It's probably easiest if you [look at the file](https://github.com/flux-framework/flux-operator/blob/main/examples/storage/google/fusion/minicluster.yaml), 
which is well-commented. For the permissions, ideally someone can test this out and design a configuration that will allow
running as the Flux user and not being root to do the mount. Let's do it - create the MiniCluster!

```bash
$ kubectl apply -f examples/storage/google/fusion/minicluster.yaml
```

Wait to see the certificate generator pod come up, complete, and the worker pods (that depend on it) will finish creation and
then come up:

```bash
$ kubectl get pods -n flux-operator
```

The use of "prefix" in the container commands ensures we wrap all the flux commands with `fusion` so
the filesystem will be mounted in all nodes. Since we are running in interactive mode, 
shell in to run the snakemake workflow.


```bash
# Shell into the broker pod
$ kubectl exec -it -n flux-operator flux-sample-0-87ll7 -- bash 

# cd there
$ cd /fusion/gs/flux-operator-storage/snakemake-workflow

# Run the workflow!
$ flux start snakemake --cores 2  --jobs 2 --flux
```

And likely this needs refining in terms of how we automate (and optimizing the bind and permissions), but
it's a pretty good feat to figure out given that other storage options have taken up to a week!
When you exit the container:

```bash
$ kubectl delete -f examples/storage/google/fusion/minicluster.yaml
```


```{include} includes/cleanup.md
```



## How does it work?

You can read [more about fusion here](https://seqera.io/fusion/). Specifically, I was confused why it didn't require me to define
a bucket anywhere. The reason is because the request for a bucket path happens when you list the filesystem. E.g.,
although the original root at `/fusion` looked empty, this is what happened when I listed - I saw my workflow 
in storage!

```bash
$ tree /fusion/gs/flux-operator-storage
```
```console
/fusion/gs/flux-operator-storage
└── snakemake-workflow
    ├── Dockerfile
    ├── README.md
    ├── Snakefile
    ├── calls
    │   └── all.vcf
    ├── data
    │   ├── genome.fa
    │   ├── genome.fa.amb
    │   ├── genome.fa.ann
    │   ├── genome.fa.bwt
    │   ├── genome.fa.fai
    │   ├── genome.fa.pac
    │   ├── genome.fa.sa
    │   └── samples
    │       ├── A.fastq
    │       ├── B.fastq
    │       └── C.fastq
    ├── environment.yaml
    ├── mapped_reads
    │   ├── A.bam
    │   └── B.bam
    ├── plots
    │   └── quals.svg
    ├── scripts
    │   └── plot-quals.py
    └── sorted_reads
        ├── A.bam
        ├── A.bam.bai
        ├── B.bam
        └── B.bam.bai
```

It's so cool how easy that was to get working! Nice! I know it's a lot of steps, but trust me
this took an afternoon to figure out after some initial reading, and other storage drivers
have taken up to a week to figure out.
