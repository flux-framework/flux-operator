# Google Cloud

On Google Cloud, "GKE" refers to Google Kubernetes Engine. Kubernetes was originally a Google project,
so it works really well on Google Cloud, and it's super fun to use the Flux Operator there! This
guide will walk you through the basics. We will be deploying a simple [lammps](https://www.lammps.org/) workflow as
shown below.

## Setup

Whether you choose the LAMMPS or Snakemake tutorial below, the setup for Google cloud is the same!
Follow these steps and then continue on to either of the tutorials.

### Install

You should first [install gcloud](https://cloud.google.com/sdk/docs/quickstarts)
and ensure you are logged in and have kubectl installed:

```bash
$ gcloud auth login
```

Depending on your install, you can either install kubectl with gcloud:

```bash
$ gcloud components install kubectl
```

or just [on your own](https://kubernetes.io/docs/tasks/tools/). I already
had it installed so I was good to go.

### Create Cluster

Now let's use gcloud to create a cluster, and we are purposefully going to choose
a very small node type to test. Note that I choose us-central1-a because it tends
to be cheaper (and closer to me). First, here is our project for easy access:

```bash
GOOGLE_PROJECT=myproject
```
Replace the above with your project name, of course!

```bash
$ gcloud container clusters create flux-cluster --project $GOOGLE_PROJECT \
    --zone us-central1-a --cluster-version 1.23 --machine-type n1-standard-1 \
    --num-nodes=4 --enable-network-policy --tags=flux-cluster --enable-intra-node-visibility
```

Note that not all of the flags above might be necessary - I did a lot of testing to get
this working and didn't go back and try removing things after the fact!
If you want to use cloud dns instead (after [enabling it](https://console.cloud.google.com/apis/library/dns.googleapis.com))

```bash
$ gcloud beta container clusters create flux-cluster --project $GOOGLE_PROJECT \
    --zone us-central1-a --cluster-version 1.23 --machine-type n1-standard-1 \
    --num-nodes=4 --enable-network-policy --tags=flux-cluster --enable-intra-node-visibility \
    --cluster-dns=clouddns \
    --cluster-dns-scope=cluster
```

In your Google cloud interface, you should be able to see the cluster! Note
this might take a few minutes.

![img/cluster.png](img/cluster.png)

I also chose a tiny size (nodes and instances) anticipating having it up longer to figure things out.

### Get Credentials

Next we need to ensure that we can issue commands to our cluster with kubectl.
To get credentials, in the view shown above, select the cluster and click "connect."
Doing so will show you the correct statement to run to configure command-line access,
which probably looks something like this:

```bash
$ gcloud container clusters get-credentials flux-cluster --zone us-central1-a --project $GOOGLE_PROJECT
```
```console
Fetching cluster endpoint and auth data.
kubeconfig entry generated for flux-cluster.
```

Finally, use [cloud IAM](https://cloud.google.com/iam) to ensure you can create roles, etc.

```bash
$ kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user $(gcloud config get-value core/account)
```
```console
clusterrolebinding.rbac.authorization.k8s.io/cluster-admin-binding created
```

At this point you should be able to get your nodes:

```bash
$ kubectl get nodes
```
```console
NAME                                            STATUS   ROLES    AGE     VERSION
gke-flux-cluster-default-pool-f103d9d8-379m   Ready    <none>   3m41s   v1.23.14-gke.1800
gke-flux-cluster-default-pool-f103d9d8-3wf9   Ready    <none>   3m42s   v1.23.14-gke.1800
gke-flux-cluster-default-pool-f103d9d8-c174   Ready    <none>   3m42s   v1.23.14-gke.1800
gke-flux-cluster-default-pool-f103d9d8-zz1q   Ready    <none>   3m42s   v1.23.14-gke.1800
```

### Deploy Operator

To deploy the Flux Operator, [choose one of the options here](https://flux-framework.org/flux-operator/getting_started/user-guide.html#production-install) to deploy the operator. Whether you apply a yaml file, use [flux-cloud](https://converged-computing.github.io/flux-cloud) or clone the repository and `make deploy` you will see the operator install to the `operator-system` namespace.

For a quick "production deploy" from development, the Makefile has a directive that will build and push a `test` tag (you'll need to edit `DEVIMG` to be one you can push to) and then generate a
yaml file targeting that image, e.g.,

```bash
$ make test-deploy
$ kubectl apply -f examples/dist/flux-operator-dev.yaml
```

```console
...
clusterrole.rbac.authorization.k8s.io/operator-manager-role created
clusterrole.rbac.authorization.k8s.io/operator-metrics-reader created
clusterrole.rbac.authorization.k8s.io/operator-proxy-role created
rolebinding.rbac.authorization.k8s.io/operator-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/operator-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/operator-proxy-rolebinding created
configmap/operator-manager-config created
service/operator-controller-manager-metrics-service created
deployment.apps/operator-controller-manager created
```

Ensure the `operator-system` namespace was created:

```bash
$ kubectl get namespace
NAME              STATUS   AGE
default           Active   6m39s
kube-node-lease   Active   6m42s
kube-public       Active   6m42s
kube-system       Active   6m42s
operator-system   Active   39s
```
```bash
$ kubectl describe namespace operator-system
Name:         operator-system
Labels:       control-plane=controller-manager
              kubernetes.io/metadata.name=operator-system
Annotations:  <none>
Status:       Active

Resource Quotas
  Name:                              gke-resource-quotas
  Resource                           Used  Hard
  --------                           ---   ---
  count/ingresses.extensions         0     100
  count/ingresses.networking.k8s.io  0     100
  count/jobs.batch                   0     5k
  pods                               1     1500
  services                           1     500

No LimitRange resource.
```

And you can find the name of the operator pod as follows:

```bash
$ kubectl get pod --all-namespaces
```
```console
      <none>
operator-system   operator-controller-manager-56b5bcf9fd-m8wg4               2/2     Running   0          73s
```

### Create Flux Operator namespace

Make your namespace for the flux-operator custom resource definition (CRD), which is the yaml file above that generates the MiniCluster:

```bash
$ kubectl create namespace flux-operator
```

## Lammps on Google Kubernetes Engine

In this short experiment we will run the Flux Operator on Google Cloud, at
at a fairly small size intended for development.

### Custom Resource Definition

The Custom Resource Definition (CRD) defines our Mini Cluster, and is what we hand to the flux
operator to create it.  Here is the CRD for a small lammps run.

```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:
  # Number of pods to create for MiniCluster
  size: 4

  # Disable verbose output
  logging:
    quiet: true

  # This is a list because a pod can support multiple containers
  containers:
    # The container URI to pull (currently needs to be public)
    - image: ghcr.io/rse-ops/lammps:flux-sched-focal-v0.24.0

      # You can set the working directory if your container WORKDIR is not correct.
      workingDir: /home/flux/examples/reaxff/HNS
      command: lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite
```

You can save the above file as `minicluster-lammps.yaml` to get started.


### Create the Lammps Job

Now let's apply the custom resource definition to create the lammps mini cluster!
The file we generated above should be in your present working directory.
Importantly, we have set `localDeploy` to false because we need to create volume
claims and not local host mounts for shared resources.

```bash
$ kubectl apply -f minicluster-lammps.yaml
```

There are different ways to see logs for pods. First, see pods running and state.
You probably want to wait until the state changes from `ContainersCreating` to `Running`
because this is where we are pulling the chonker containers.

```bash
$ kubectl get -n flux-operator pods
```

If you need to debug (or see general output for a pod about creation) you can do:

```bash
$ kubectl -n flux-operator describe pods flux-sample-0-742bm
```

And finally, the most meaty piece of metadata is the log for the pod,
where the Flux Operator will be setting things up and starting flux.

```bash
# Add the -f to keep it hanging
$ kubectl -n flux-operator logs flux-sample-0-742bm -f
```

To shell into a pod to look around (noting where important flux stuff is)

```bash
$ kubectl exec --stdin --tty -n flux-operator flux-sample-0-742bm -- /bin/bash
```
```console
ls /mnt/curve
ls /etc/flux
ls /etc/flux/config
```

To get logs for the operator itself:

```bash
$ kubectl logs -n operator-system operator-controller-manager-56b5bcf9fd-j2g75
```

If you need to run in verbose (non-test) mode, set test to false in the [minicluster-lammps.yaml](minicluster-lammps.yaml).
And make sure to clean up first:

```bash
$ kubectl delete -f minicluster-lammps.yaml
```

and wait until the pods are gone:

```bash
$ kubectl get -n flux-operator pods
No resources found in flux-operator namespace.
```

Observations about comparing this to MiniKube (local):

 - The containers that are large actually pull!
 - The startup times of the different pods vary quite a bit.
 - A few config maps aren't found or timed out mount for up to 3-4 minutes, then it runs.
 - Sometimes it also runs quickly!

If you want to run the same workflow again, use `kubectl delete -f` with the file
and apply it again. I wound up running with test set to true, and then saving the logs:

```bash
$ kubectl -n flux-operator logs flux-sample-0-qc5z2 > lammps.out
```

For fun, here is the first successful run of Lammps using the Flux Operator on GCP
ever!

![img/lammps.png](img/lammps.png)

Then to delete your lammps MiniCluster:

```bash
$ kubectl delete -f minicluster-lammps.yaml
```

## Snakemake (requiring storage) on Google Kubernetes Engine

Akin to how we created a local volume, we can do something similar, but instead of pointing the Flux Operator
to a volume on the host (e.g., in MiniKube) we are going to point it to a storage bucket with our data.
For Google cloud, the Flux Operator currently uses the [this driver](https://github.com/ofek/csi-gcs) to
connect a cloud storage bucket to our cluster.

### Prepare Data

To start, prepare your data in a temporary directory (that we will upload into Google cloud storage):

```bash
$ git clone --depth 1 https://github.com/snakemake/snakemake-tutorial-data /tmp/workflow
```

You'll want to add the [Snakefile](https://github.com/rse-ops/flux-hpc/blob/main/snakemake/atacseq/Snakefile) for your workflow
along with a [plotting script](https://github.com/rse-ops/flux-hpc/blob/main/snakemake/atacseq/scripts/plot-quals.py):

```bash
$ wget -O /tmp/workflow/Snakefile https://raw.githubusercontent.com/rse-ops/flux-hpc/main/snakemake/atacseq/Snakefile
$ mkdir -p /tmp/workflow/scripts
$ wget -O /tmp/workflow/scripts/plot-quals.py https://raw.githubusercontent.com/rse-ops/flux-hpc/main/snakemake/atacseq/scripts/plot-quals.py
```

You should have this structure:

```bash
$ tree /tmp/workflow
```
```
/tmp/workflow/
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
├── Dockerfile
├── environment.yaml
├── README.md
├── scripts
│   └── plot-quals.py
└── Snakefile
```

We can then use Google Cloud (`gcloud`) to create a bucket and upload to it.

```bash
$ gcloud storage buckets create gs://flux-operator-storage --project=dinodev  --location=US-CENTRAL1 --uniform-bucket-level-access
```

In the above, the storage class defaults to "Standard" Once we've created the bucket, let's go to our snakemake data
and upload the data to it.

```bash
$ cd /tmp/workflow
$ gcloud storage cp --recursive . gs://flux-operator-storage/snakemake-workflow/
```

You can either view your files in the Google Storage console

![img/google-storage-console.png](img/google-storage-console.png)

or view with gcloud again:

```bash
$ gcloud storage ls gs://flux-operator-storage/snakemake-workflow/
```
```console
gs://flux-operator-storage/snakemake-workflow/.gitpod.yml
gs://flux-operator-storage/snakemake-workflow/Dockerfile
gs://flux-operator-storage/snakemake-workflow/README.md
gs://flux-operator-storage/snakemake-workflow/Snakefile
gs://flux-operator-storage/snakemake-workflow/environment.yaml
gs://flux-operator-storage/snakemake-workflow/data/
gs://flux-operator-storage/snakemake-workflow/scripts/
```

### Install the Constainer Storage Driver (CSI)

There are many [drivers](https://kubernetes-csi.github.io/docs/drivers.html) for kubernetes, and we will use
[this one](https://ofek.dev/csi-gcs/getting_started/) that requires a stateful set and daemon set to work.
Let's install those first.

```bash
$ kubectl apply -k "github.com/ofek/csi-gcs/deploy/overlays/stable?ref=v0.9.0"
$ kubectl get CSIDriver,daemonsets,pods -n kube-system | grep csi
```

And to debug:

```bash
$ kubectl logs -l app=csi-gcs -c csi-gcs -n kube-system
```

As you are working, if the mounts seem to work but you don't see files, keep
in mind you need to be aware of [implicit directories](https://ofek.dev/csi-gcs/dynamic_provisioning/#extra-flags).
The operator will do a `mkdir -p` on the working directory (and this will show content there) but if you don't
see content and expect to, you either need to interact in this way or set this flag as an annotation in
your `minicluster.yaml`.

### Permissions via Secrets

We will need to give permission for the nodes to access storage, and we can do that via [these instructions](https://ofek.dev/csi-gcs/dynamic_provisioning/#permission)
to create a service account key (a json file) from a service account. E.g., I first created a custom service account that
has these permissions:

![img/google-service-account.png](img/google-service-account.png)

And then I could find the new identifier in a listing:

```bash
$ gcloud iam service-accounts list
DISPLAY NAME                            EMAIL                                               DISABLED
Compute Engine default service account  270958151865-compute@developer.gserviceaccount.com  False
flux-operator-gke                       flux-operator-gke@project.iam.gserviceaccount.com   False
```

And create a credential file:

```bash
$ gcloud iam service-accounts keys create <FILE_NAME>.json --iam-account <EMAIL>
```

And create a secret from it! This is basically giving your cluster permission to interact with a specific bucket.

```bash
$ kubectl create secret generic csi-gcs-secret --from-literal=bucket=flux-operator-storage --from-file=key=<PATH_TO_SERVICE_ACCOUNT_KEY>
```

### Storage Class

We can then create our storage class, this file is provided in `examples/storage/google/storageclass.yaml`

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-gcs
provisioner: gcs.csi.ofek.dev
```

```bash
$ kubectl apply -f examples/storage/google/storageclass.yaml
```

### Snakemake MiniCluster

The operator can now be run, telling it to use this storage class (named `csi-gcs`) - we provide
an example Minicluster to run this snakemake tutorial to test this out, and note that if you want to debug
you can change the command first to be an ls to see what data is there! Before you apply this file,
make sure that the annotations for storage match your Google project, zone, etc.

```yaml
  # Make this kind of persistent volume and claim available to pods
  # This is a type of storage that will use Google Storage
  volumes:
    data:
      class: csi-gcs
      path: /tmp/data
      secret: csi-gcs-secret
      annotations:
        gcs.csi.ofek.dev/location: us-central1
        gcs.csi.ofek.dev/project-id: my-project
        gcs.csi.ofek.dev/bucket: flux-operator-storage
```

Also note that we are setting the `commands: -> runFluxAsRoot` to true. This isn't ideal, but it was the
only way I could get the storage to both be seen and have permission to write there. Let's create the job!

```bash
$ kubectl apply -f examples/storage/google/minicluster.yaml
```

Wait to see the certificate generator pod come up, complete, and the worker pods (that depend on it) will finish creation and
then come up:

```bash
$ kubectl get pods -n flux-operator
```

And I like to get the main pod and stream the output so I don't miss it:

```bash
# Stream to your terminal
$ kubectl logs -n flux-operator flux-sample-0-fj6td -f

# Stream to file
$ kubectl logs -n flux-operator flux-sample-0-fj6td -f > output.txt
```

<details>

<summary>Snakemake output from Log</summary>

```console
flux user identifiers:
uid=1000(flux) gid=1000(flux) groups=1000(flux)

As Flux prefix for flux commands: sudo -E PYTHONPATH=/usr/lib/flux/python3.8:/code -E PATH=/opt/micromamba/envs/snakemake/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin -E HOME=/home/flux

👋 Hello, I'm flux-sample-0
The main host is flux-sample-0
The working directory is /workflow/snakemake-workflow, contents include:
Dockerfile  README.md  Snakefile  data  environment.yaml  scripts
End of file listing, if you see nothing above there are no files.
flux R encode --hosts=flux-sample-[0-1]

📦 Resources
{"version": 1, "execution": {"R_lite": [{"rank": "0-1", "children": {"core": "0"}}], "starttime": 0.0, "expiration": 0.0, "nodelist": ["flux-sample-[0-1]"]}}

🐸 Diagnostics: false

🦊 Independent Minister of Privilege
[exec]
allowed-users = [ "flux", "root" ]
allowed-shells = [ "/usr/libexec/flux/flux-shell" ]

🐸 Broker Configuration
# Flux needs to know the path to the IMP executable
[exec]
imp = "/usr/libexec/flux/flux-imp"

[access]
allow-guest-user = true
allow-root-owner = true

# Point to resource definition generated with flux-R(1).
[resource]
path = "/etc/flux/system/R"

[bootstrap]
curve_cert = "/etc/curve/curve.cert"
default_port = 8050
default_bind = "tcp://eth0:%p"
default_connect = "tcp://%h.flux-service.flux-operator.svc.cluster.local:%p"
hosts = [
        { host="flux-sample-[0-1]"},
]
#   ****  Generated on 2023-02-17 21:05:10 by CZMQ  ****
#   ZeroMQ CURVE **Secret** Certificate
#   DO NOT PROVIDE THIS FILE TO OTHER USERS nor change its permissions.

metadata
    name = "flux-sample-cert-generator"
    time = "2023-02-17T21:05:10"
    userid = "0"
    hostname = "flux-sample-cert-generator"
curve
    public-key = ".!?zfo10Ew)m=+J:j^zehs&{Ayy#BGSV0Eets5Ne"
    secret-key = "vmk%8&dl7ICTfgx?*+0wgPb=@kFA>djvZU-Sl[T6"

🔒️ Working directory permissions:
total 3
-rw-rw-r-- 1 root 63147  233 Feb 10 22:57 Dockerfile
-rw-rw-r-- 1 root 63147  347 Feb 10 22:57 README.md
-rw-rw-r-- 1 root 63147 1144 Feb 10 22:57 Snakefile
drwxrwxr-x 1 root 63147    0 Feb 17 21:05 data
-rw-rw-r-- 1 root 63147  203 Feb 10 22:57 environment.yaml
drwxrwxr-x 1 root 63147    0 Feb 17 21:05 scripts


✨ Curve certificate generated by helper pod
#   ****  Generated on 2023-02-17 21:05:10 by CZMQ  ****
#   ZeroMQ CURVE **Secret** Certificate
#   DO NOT PROVIDE THIS FILE TO OTHER USERS nor change its permissions.

metadata
    name = "flux-sample-cert-generator"
    time = "2023-02-17T21:05:10"
    userid = "0"
    hostname = "flux-sample-cert-generator"
curve
    public-key = ".!?zfo10Ew)m=+J:j^zehs&{Ayy#BGSV0Eets5Ne"
    secret-key = "vmk%8&dl7ICTfgx?*+0wgPb=@kFA>djvZU-Sl[T6"
Extra arguments are: snakemake --cores 1 --flux

🌀 flux start -o --config /etc/flux/config -Scron.directory=/etc/flux/system/cron.d   -Stbon.fanout=256   -Srundir=/run/flux   -Sstatedir=/var/lib/flux   -Slocal-uri=local:///run/flux/local   -Slog-stderr-level=6    -Slog-stderr-mode=local flux mini submit  -n 1 --quiet  --watch snakemake --cores 1 --flux
broker.info[0]: start: none->join 10.6371ms
broker.info[0]: parent-none: join->init 0.051043ms
cron.info[0]: synchronizing cron tasks to event heartbeat.pulse
job-manager.info[0]: restart: 0 jobs
job-manager.info[0]: restart: 0 running jobs
job-manager.info[0]: restart: checkpoint.job-manager not found
broker.info[0]: rc1.0: running /etc/flux/rc1.d/01-sched-fluxion
sched-fluxion-resource.warning[0]: create_reader: allowlist unsupported
sched-fluxion-resource.info[0]: populate_resource_db: loaded resources from core's resource.acquire
broker.info[0]: rc1.0: running /etc/flux/rc1.d/02-cron
broker.info[0]: rc1.0: /etc/flux/rc1 Exited (rc=0) 0.9s
broker.info[0]: rc1-success: init->quorum 0.922753s
broker.info[0]: online: flux-sample-0 (ranks 0)
broker.info[0]: online: flux-sample-[0-1] (ranks 0-1)
broker.info[0]: quorum-full: quorum->run 0.427513s
Building DAG of jobs...
Using shell: /usr/bin/bash
Provided cores: 1 (use --cores to define parallelism)
Rules claiming more threads will be scaled down.
Job stats:
job               count    min threads    max threads
--------------  -------  -------------  -------------
all                   1              1              1
bcftools_call         1              1              1
bwa_map               2              1              1
plot_quals            1              1              1
samtools_index        2              1              1
samtools_sort         2              1              1
total                 9              1              1

Select jobs to execute...

[Fri Feb 17 21:05:22 2023]
rule bwa_map:
    input: data/genome.fa, data/samples/A.fastq
    output: mapped_reads/A.bam
    jobid: 4
    reason: Missing output files: mapped_reads/A.bam
    wildcards: sample=A
    resources: tmpdir=/tmp

[Fri Feb 17 21:05:42 2023]
Finished job 4.
1 of 9 steps (11%) done
Select jobs to execute...

[Fri Feb 17 21:05:43 2023]
rule samtools_sort:
    input: mapped_reads/A.bam
    output: sorted_reads/A.bam
    jobid: 3
    reason: Missing output files: sorted_reads/A.bam; Input files updated by another job: mapped_reads/A.bam
    wildcards: sample=A
    resources: tmpdir=/tmp

[Fri Feb 17 21:05:52 2023]
Finished job 3.
2 of 9 steps (22%) done
Select jobs to execute...

[Fri Feb 17 21:05:52 2023]
rule bwa_map:
    input: data/genome.fa, data/samples/B.fastq
    output: mapped_reads/B.bam
    jobid: 6
    reason: Missing output files: mapped_reads/B.bam
    wildcards: sample=B
    resources: tmpdir=/tmp

[Fri Feb 17 21:06:02 2023]
Finished job 6.
3 of 9 steps (33%) done
Select jobs to execute...

[Fri Feb 17 21:06:02 2023]
rule samtools_sort:
    input: mapped_reads/B.bam
    output: sorted_reads/B.bam
    jobid: 5
    reason: Missing output files: sorted_reads/B.bam; Input files updated by another job: mapped_reads/B.bam
    wildcards: sample=B
    resources: tmpdir=/tmp

[Fri Feb 17 21:06:12 2023]
Finished job 5.
4 of 9 steps (44%) done
Select jobs to execute...

[Fri Feb 17 21:06:13 2023]
rule samtools_index:
    input: sorted_reads/A.bam
    output: sorted_reads/A.bam.bai
    jobid: 7
    reason: Missing output files: sorted_reads/A.bam.bai; Input files updated by another job: sorted_reads/A.bam
    wildcards: sample=A
    resources: tmpdir=/tmp

[Fri Feb 17 21:06:22 2023]
Finished job 7.
5 of 9 steps (56%) done
Select jobs to execute...

[Fri Feb 17 21:06:22 2023]
rule samtools_index:
    input: sorted_reads/B.bam
    output: sorted_reads/B.bam.bai
    jobid: 8
    reason: Missing output files: sorted_reads/B.bam.bai; Input files updated by another job: sorted_reads/B.bam
    wildcards: sample=B
    resources: tmpdir=/tmp

[Fri Feb 17 21:06:32 2023]
Finished job 8.
6 of 9 steps (67%) done
Select jobs to execute...

[Fri Feb 17 21:06:32 2023]
rule bcftools_call:
    input: data/genome.fa, sorted_reads/A.bam, sorted_reads/B.bam, sorted_reads/A.bam.bai, sorted_reads/B.bam.bai
    output: calls/all.vcf
    jobid: 2
    reason: Missing output files: calls/all.vcf; Input files updated by another job: sorted_reads/A.bam, sorted_reads/B.bam, sorted_reads/A.bam.bai, sorted_reads/B.bam.bai
    resources: tmpdir=/tmp

[Fri Feb 17 21:06:43 2023]
Finished job 2.
7 of 9 steps (78%) done
Select jobs to execute...

[Fri Feb 17 21:06:43 2023]
rule plot_quals:
    input: calls/all.vcf
    output: plots/quals.svg
    jobid: 1
    reason: Missing output files: plots/quals.svg; Input files updated by another job: calls/all.vcf
    resources: tmpdir=/tmp

[Fri Feb 17 21:07:02 2023]
Finished job 1.
8 of 9 steps (89%) done
Select jobs to execute...

[Fri Feb 17 21:07:02 2023]
localrule all:
    input: plots/quals.svg
    jobid: 0
    reason: Input files updated by another job: plots/quals.svg
    resources: tmpdir=/tmp

[Fri Feb 17 21:07:02 2023]
Finished job 0.
9 of 9 steps (100%) done
Complete log: .snakemake/log/2023-02-17T210519.114509.snakemake.log
broker.info[0]: rc2.0: flux mini submit -n 1 --quiet --watch snakemake --cores 1 --flux Exited (rc=0) 115.1s
broker.info[0]: rc2-success: run->cleanup 1.91796m
broker.info[0]: cleanup.0: flux queue stop --quiet --all --nocheckpoint Exited (rc=0) 0.0s
broker.info[0]: cleanup.1: flux job cancelall --user=all --quiet -f --states RUN Exited (rc=0) 0.0s
broker.info[0]: cleanup.2: flux queue idle --quiet Exited (rc=0) 0.0s
broker.info[0]: cleanup-success: cleanup->shutdown 49.7479ms
broker.info[0]: children-complete: shutdown->finalize 93.2169ms
broker.info[0]: rc3.0: running /etc/flux/rc3.d/01-sched-fluxion
broker.info[0]: online: flux-sample-0 (ranks 0)
broker.info[0]: rc3.0: /etc/flux/rc3 Exited (rc=0) 0.2s
broker.info[0]: rc3-success: finalize->goodbye 0.246792s
broker.info[0]: goodbye: goodbye->exit 0.080512ms
```

</details>

After it finishes the job will cleanup (unless you've set `cleanup: false`) in your minicluster.yaml. If you check
Google storage, you'll see the output of the run (e.g., data in mapped_reads / mapped_samples /plots) and a ".snakemake"
hidden directory with logs.

![img/snakemake-output.png](img/snakemake-output.png)

Note that when this is deployed in a larger scale or a more productions sense, you'll want to account for
the details (e.g., resource limits) of the storage. This is basically done by (instead of applying the GitHub provided files)
editing them [via kustomize](https://ofek.dev/csi-gcs/getting_started/#resource-requests-limits). Note that I haven't
tried this yet.

**Note**: we would like to get this working without requiring running the workflow as root, but it hasn't been figured
out yet! If you have insight, please comment on [this issue](https://github.com/ofek/csi-gcs/issues/155).


## Clean up

Whatever tutorial you choose, don't forget to clean up at the end!
You can optionally undeploy the operator (this is again at the root of the operator repository clone)

```bash
$ make undeploy
```

Or the file you used to deploy it:

```bash
$ kubectl delete -f examples/dist/flux-operator.yaml
$ kubectl delete -f examples/dist/flux-operator-dev.yaml
```

And then to delete the cluster with gcloud:

```bash
$ gcloud container clusters delete --zone us-central1-a flux-cluster
```

I like to check in the cloud console to ensure that it was actually deleted.


## Customization and Debugging

### Firewall

When I first created my cluster, the nodes could not see one another. I added a few
flags for networking, and looked at firewalls as follows:

```bash
$ gcloud container clusters describe flux-cluster --zone us-central1-a | grep clusterIpv4Cidr
```
I didn't ultimately change anything, but I found this useful.
