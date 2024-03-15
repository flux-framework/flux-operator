# Filestore (NFS)

 **[MiniCluster YAML](https://github.com/flux-framework/flux-operator/blob/main/examples/storage/google/filestore/minicluster.yaml)**

This tutorial will walk through creating a more persistent MiniCluster on Google Cloud
using Filestore. We will be following the guidance [here](https://cloud.google.com/filestore/docs/csi-driver).
First, make sure you have the Filestore and GKE (Google Kubernetes Engine) APIs enabled,
and the other introductory steps at that link. If you've used Google Cloud for Kubernetes and
Filestore before, you should largely be good to go.

## Create the Cluster

First, create your cluster with the Filestore CSI Driver enabled:

```bash
GOOGLE_PROJECT=myproject
```
```bash
$ gcloud container clusters create flux-cluster --project $GOOGLE_PROJECT \
    --zone us-central1-a --machine-type n1-standard-2 \
    --addons=GcpFilestoreCsiDriver \
    --num-nodes=4 --enable-network-policy --tags=flux-cluster --enable-intra-node-visibility
```


## Install the Flux Operator

Let's next install the operator. You can [choose one of the options here](https://flux-framework.org/flux-operator/getting_started/user-guide.html#production-install). 
E.g., to deploy from the cloned repository:

```bash
$ kubectl apply -f examples/dist/flux-operator.yaml
```

Next, we want to create our persistent volume claim. This will use the
storage drivers installed already to our cluster (via the creation command).

```yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: data
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Ti
  storageClassName: standard-rwx
```

Note that we've selected the storage class "standard-rwx" from [this list](https://cloud.google.com/filestore/docs/csi-driver#storage-class).
You can also see the `storageclass` available for your cluster via this command:

```bash
$ kubectl get storageclass
```

<details>

<summary>storageclass Available on our Filestore Cluster</summary>

```console
NAME                        PROVISIONER                    RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
enterprise-multishare-rwx   filestore.csi.storage.gke.io   Delete          WaitForFirstConsumer   true                   5m58s
enterprise-rwx              filestore.csi.storage.gke.io   Delete          WaitForFirstConsumer   true                   5m57s
premium-rwo                 pd.csi.storage.gke.io          Delete          WaitForFirstConsumer   true                   5m25s
premium-rwx                 filestore.csi.storage.gke.io   Delete          WaitForFirstConsumer   true                   5m58s
standard                    kubernetes.io/gce-pd           Delete          Immediate              true                   5m24s
standard-rwo (default)      pd.csi.storage.gke.io          Delete          WaitForFirstConsumer   true                   5m25s
standard-rwx                filestore.csi.storage.gke.io   Delete          WaitForFirstConsumer   true                   5m58s
```

</details>

A Filestore storage class will *not* be the default (see above output) so this step is important to take! We 
are going to create a persistent volume claim that says "Please use the `standard-rwx` storageclass from
Filestore to be available as a persistent volume claim - and I want all of it - the entire 1TB!

```bash
$ kubectl apply -f examples/storage/google/filestore/pvc.yaml
```

And check on the status:

```bash
$ kubectl get pvc
NAME   STATUS    VOLUME   CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data   Pending                                      standard-rwx   6s
```

It will be pending under we make a request to use it! Let's do that by creating the MiniCluster:

```bash
$ kubectl apply -f examples/storage/google/filestore/minicluster.yaml
```

We can now look at the pods:

```bash
kubectl get pods
```
```console
NAME                         READY   STATUS      RESTARTS   AGE
flux-sample-0-vb59t          1/1     Running     0          3m22s
flux-sample-1-xnf56          1/1     Running     0          3m21s
flux-sample-2-7ws6d          1/1     Running     0          3m21s
flux-sample-3-tw8m2          1/1     Running     0          3m21s
```
Note that I saw this error at first:

```
Warning  FailedScheduling  32s   default-scheduler  running PreBind plugin "VolumeBinding": Operation cannot be fulfilled on persistentvolumeclaims "data": the object has been modified; please apply your changes to the latest version and try again
```

I tried applying it again (and actually just waited) and then it worked - I don't think I actually did anything. I would
give it a few minutes to (hopefully) resolve as it did for me!

## Test your Storage

Let's shell into the broker (index 0) and play with our Filesystem!

```bash
$ kubectl exec -it flux-sample-0-vb59t bash
```

Is it there? How big is it?

```bash
$ df -a | grep /workflow
10.247.103.226:/vol1 1055763456        0 1002059776   0% /workflow
```

That is a LOT of space (that is overkill for this tutorial).
Let's add a dinosaur fart there to test if we can see it from another pod.

```bash
$ touch /workflow/dinosaur-fart
```

Did we make it?

```bash
$ ls /workflow/
dinosaur-fart  lost+found
```

And now exit and shell into another pod...

```bash
$ kubectl exec -it flux-sample-1-xnf56 bash
```
Is the dinosaur fart there?

```bash
root@flux-sample-1:/code# ls /workflow/
dinosaur-fart  lost+found
```

We have a dinosaur fart! I repeat - we have a dinosaur fart!! 🦖🌬️

## Run Snakemake

Let's now run Snakemake! Since we haven't properly set up permissions, we need to set up the workspace as root,
and then give ownership to the flux user. Let's prepare snakemake data in `/workflow` and run it.

```bash
git clone --depth 1 https://github.com/snakemake/snakemake-tutorial-data /workflow/snakemake-workflow
```

You'll want to add the [Snakefile](https://github.com/rse-ops/flux-hpc/blob/main/snakemake/atacseq/Snakefile) for your workflow
along with a [plotting script](https://github.com/rse-ops/flux-hpc/blob/main/snakemake/atacseq/scripts/plot-quals.py):

```bash
wget -O /workflow/snakemake-workflow/Snakefile https://raw.githubusercontent.com/rse-ops/flux-hpc/main/snakemake/atacseq/Snakefile
mkdir -p /workflow/snakemake-workflow/scripts
wget -O /workflow/snakemake-workflow/scripts/plot-quals.py https://raw.githubusercontent.com/rse-ops/flux-hpc/main/snakemake/atacseq/scripts/plot-quals.py
```

Now let's run it! Remember that we are in interactive mode, so the broker is already running and we need to connect to it. Let's do that first.

```bash
source /mnt/flux/flux-view.sh
flux proxy $fluxsocket bash
```

Before snakemake, let's run a test job. You should see some subset of nodes:

```bash
$ flux run -N 4 hostname
flux-sample-0
flux-sample-2
flux-sample-3
flux-sample-1
```

And that the resources are available:

```bash
$ flux resource list
     STATE NNODES   NCORES NODELIST
      free      4        4 flux-sample-[0-3]
 allocated      0        0 
      down      0        0 
```

Let's go to the worklow directory and run Snakemake using Flux.


```bash
$ cd /workflow/snakemake-workflow
$ snakemake --cores 1 --exeutor flux --jobs 4
```

<details>

<summary>Snakemake Workflow Output</summary>

```console
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

[Thu Apr  6 23:21:57 2023]
rule bwa_map:
    input: data/genome.fa, data/samples/B.fastq
    output: mapped_reads/B.bam
    jobid: 6
    reason: Missing output files: mapped_reads/B.bam
    wildcards: sample=B
    resources: tmpdir=/tmp

Job 6 has been submitted with flux jobid ƒ8oNbgfi7 (log: .snakemake/flux_logs/bwa_map/sample_B.log).

[Thu Apr  6 23:21:57 2023]
rule bwa_map:
    input: data/genome.fa, data/samples/A.fastq
    output: mapped_reads/A.bam
    jobid: 4
    reason: Missing output files: mapped_reads/A.bam
    wildcards: sample=A
    resources: tmpdir=/tmp

Job 4 has been submitted with flux jobid ƒ8oRe7CF9 (log: .snakemake/flux_logs/bwa_map/sample_A.log).
[Thu Apr  6 23:22:07 2023]
Finished job 6.
1 of 9 steps (11%) done
[Thu Apr  6 23:22:07 2023]
Finished job 4.
2 of 9 steps (22%) done
Select jobs to execute...

[Thu Apr  6 23:22:07 2023]
rule samtools_sort:
    input: mapped_reads/B.bam
    output: sorted_reads/B.bam
    jobid: 5
    reason: Missing output files: sorted_reads/B.bam; Input files updated by another job: mapped_reads/B.bam
    wildcards: sample=B
    resources: tmpdir=/tmp

Job 5 has been submitted with flux jobid ƒ8skKRkvs (log: .snakemake/flux_logs/samtools_sort/sample_B.log).

[Thu Apr  6 23:22:07 2023]
rule samtools_sort:
    input: mapped_reads/A.bam
    output: sorted_reads/A.bam
    jobid: 3
    reason: Missing output files: sorted_reads/A.bam; Input files updated by another job: mapped_reads/A.bam
    wildcards: sample=A
    resources: tmpdir=/tmp

Job 3 has been submitted with flux jobid ƒ8snWvhAX (log: .snakemake/flux_logs/samtools_sort/sample_A.log).
[Thu Apr  6 23:22:17 2023]
Finished job 5.
3 of 9 steps (33%) done
[Thu Apr  6 23:22:17 2023]
Finished job 3.
4 of 9 steps (44%) done
Select jobs to execute...

[Thu Apr  6 23:22:17 2023]
rule samtools_index:
    input: sorted_reads/A.bam
    output: sorted_reads/A.bam.bai
    jobid: 7
    reason: Missing output files: sorted_reads/A.bam.bai; Input files updated by another job: sorted_reads/A.bam
    wildcards: sample=A
    resources: tmpdir=/tmp

Job 7 has been submitted with flux jobid ƒ8x9yquZq (log: .snakemake/flux_logs/samtools_index/sample_A.log).

[Thu Apr  6 23:22:17 2023]
rule samtools_index:
    input: sorted_reads/B.bam
    output: sorted_reads/B.bam.bai
    jobid: 8
    reason: Missing output files: sorted_reads/B.bam.bai; Input files updated by another job: sorted_reads/B.bam
    wildcards: sample=B
    resources: tmpdir=/tmp

Job 8 has been submitted with flux jobid ƒ8xC8NsEo (log: .snakemake/flux_logs/samtools_index/sample_B.log).
[Thu Apr  6 23:22:27 2023]
Finished job 7.
5 of 9 steps (56%) done
[Thu Apr  6 23:22:27 2023]
Finished job 8.
6 of 9 steps (67%) done
Select jobs to execute...

[Thu Apr  6 23:22:27 2023]
rule bcftools_call:
    input: data/genome.fa, sorted_reads/A.bam, sorted_reads/B.bam, sorted_reads/A.bam.bai, sorted_reads/B.bam.bai
    output: calls/all.vcf
    jobid: 2
    reason: Missing output files: calls/all.vcf; Input files updated by another job: sorted_reads/A.bam, sorted_reads/A.bam.bai, sorted_reads/B.bam, sorted_reads/B.bam.bai
    resources: tmpdir=/tmp

Job 2 has been submitted with flux jobid ƒ92ZZp6Mm (log: .snakemake/flux_logs/bcftools_call.log).
[Thu Apr  6 23:22:37 2023]
Finished job 2.
7 of 9 steps (78%) done
Select jobs to execute...

[Thu Apr  6 23:22:37 2023]
rule plot_quals:
    input: calls/all.vcf
    output: plots/quals.svg
    jobid: 1
    reason: Missing output files: plots/quals.svg; Input files updated by another job: calls/all.vcf
    resources: tmpdir=/tmp

Job 1 has been submitted with flux jobid ƒ96y5LKJf (log: .snakemake/flux_logs/plot_quals.log).
[Thu Apr  6 23:22:47 2023]
Finished job 1.
8 of 9 steps (89%) done
Select jobs to execute...

[Thu Apr  6 23:22:47 2023]
localrule all:
    input: plots/quals.svg
    jobid: 0
    reason: Input files updated by another job: plots/quals.svg
    resources: tmpdir=/tmp

[Thu Apr  6 23:22:47 2023]
Finished job 0.
9 of 9 steps (100%) done
Complete log: .snakemake/log/2023-04-06T232156.920090.snakemake.log
```

</details>

Wow it runs really fast with 4 jobs! When you finish, you should see all of the output files you'd expect:

```bash
flux@flux-sample-1:/workflow/snakemake-workflow$ tree .
.
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

7 directories, 23 files
```

And that's it! This is a really exciting development, because Filestore can easily provide an NFS filesystem, 
meaning that we could very easily create an entire Flux Framework cluster with users and a user-respecting filesystem!
We would provide this NFS Filesystem at `/home`, and then create the user directories under it, and the only check
you'd need to do is that your container base isn't keeping anything important there.

I'm pumped, yo!

## Clean Up

Don't forget to clean up! Delete the MiniCluster and PVC first:

```bash
$ kubectl delete -f examples/storage/google/filestore/minicluster.yaml
$ kubectl delete -f examples/storage/google/filestore/pvc.yaml
```

And then delete your Kubernetes cluster. Technically  you could probably just do this, but we might as well be proper!

```bash
$ gcloud container clusters delete --zone us-central1-a flux-cluster
```

The other really nice detail (seriously, amazing job Google!) is that you largely don't have to think about setting
up the storage. The Filestore instance is created with your PVC, and it's cleaned up too.

> I seriously love this - I think (costs aside, assuming someone else is paying for this beast) this is my favorite Kubernetes storage solution I've encountered thus far! - [Vanessasaurus](https://github.com/vsoch)