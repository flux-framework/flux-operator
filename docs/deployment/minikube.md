# MiniKube

This small tutorial wall walk through how to run the Flux Operator on MiniKube,
deploying a workflow that requires storage.

## External Examples

The following (external) example tutorials are available:

 - [nsfd-materialscience](https://github.com/converged-computing/workflows/tree/main/nsdf-materialscience): to preprocess images - requires private data, but code is public.


## Workflow with Storage

The example here will use the [Snakemake workflow](https://github.com/flux-framework/flux-operator/tree/main/examples/tests/snakemake) that comes alongside the examples! Instead of using data that is
already in the [base container](https://github.com/rse-ops/flux-hpc/tree/main/snakemake/atacseq), we will mount it as a shared volume to demonstrate how that works.

### Cluster

Bring up your MiniKube cluster!

```bash
$ minikube start
```
If you are working from the repository, you can install and run the operator via developer commands:

```bash
$ make
$ make install
$ make run
```

Or external to that, you can install the latest version:

```bash
$ wget https://raw.githubusercontent.com/flux-framework/flux-operator/main/examples/dist/flux-operator.yaml
$ kubectl apply -f flux-operator.yaml
```

### Prepare Data

Let's first clone the required data to our local machine. We are going to bind this to the pods.

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
$ tree /tmp/workflow/
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

In a real production setting, you would likely have this data in cloud storage somewhere. We are emulating
that for MiniKube with a clone to a temporary directory! Note that to make it easier

### Containers

MiniKube can burp sometimes on pulling large images, so let's issue a command to it to
pull the container we need via ssh.

```bash
$ minikube ssh docker pull ghcr.io/rse-ops/atacseq:app-latest
```

### Mount Data

The Flux Operator is going to expect to find volumes on the host of a particular storage type.
Since we are early in development, we currently (as the default) define a "hostpath" storage type,
meaning the operator will expect the path to be present on the node where you are running the job.
This means that we need to mount the data on our host into MiniKube (where the cluster is running)
with `minikube mount`.

Note that in our [minicluster.yaml](https://github.com/flux-framework/flux-operator/tree/main/examples/tests/snakemake/minicluster.yaml)
we are defining the volume on the host (inside the MiniKube VM) to be at `/tmp/data` so let's tell MiniKube to mount our local path there:

```bash
minikube ssh -- mkdir -p /tmp/data

minikube mount /tmp/workflow:/tmp/data
```

Leave that process running in a window and then open another terminal to interact with the cluster.
If you want to double check the data is in the MiniKube vm:

```bash
$ minikube ssh -- ls /tmp/data
```
```console
Dockerfile  README.md  Snakefile  data	environment.yaml  scripts
```

### Run Workflow

At this point we want to run our Snakemake workflow, which means applying the custom resource definition
for the MiniCluster. Conceptually, this is going to bring up the pods on your MiniCluster nodes,
start a flux instance, run your command of interest, and then clean up. Here we can work from
the `examples/tests/snakemake` directory with the `minicluster.yaml` file:

```bash
$ kubectl apply -f minicluster.yaml
minicluster.flux-framework.org/flux-sample created
```

You can see your pods:

```bash
$ kubectl get -n flux-operator pods
NAME                         READY   STATUS      RESTARTS   AGE
flux-sample-0-sk72j          1/1     Running     0          35s
flux-sample-1-6tw4z          1/1     Running     0          35s
flux-sample-cert-generator   0/1     Completed   0          35s
```

Keep watching for when the job is finished! You can see logs (and snakemake output) as follows:

```bash
$ kubectl logs -n flux-operator flux-sample-0-pssf6 -f
```

<details>

<summary>Snakemake Output in Console</summary>

```console
broker.info[0]: quorum-full: quorum->run 0.420945s
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

[Fri Feb 10 00:59:17 2023]
rule bwa_map:
    input: data/genome.fa, data/samples/A.fastq
    output: mapped_reads/A.bam
    jobid: 4
    reason: Missing output files: mapped_reads/A.bam
    wildcards: sample=A
    resources: tmpdir=/tmp

[Fri Feb 10 00:59:27 2023]
Finished job 4.
1 of 9 steps (11%) done
Select jobs to execute...

[Fri Feb 10 00:59:27 2023]
rule samtools_sort:
    input: mapped_reads/A.bam
    output: sorted_reads/A.bam
    jobid: 3
    reason: Missing output files: sorted_reads/A.bam; Input files updated by another job: mapped_reads/A.bam
    wildcards: sample=A
    resources: tmpdir=/tmp

[Fri Feb 10 00:59:37 2023]
Finished job 3.
2 of 9 steps (22%) done
Select jobs to execute...

[Fri Feb 10 00:59:37 2023]
rule bwa_map:
    input: data/genome.fa, data/samples/B.fastq
    output: mapped_reads/B.bam
    jobid: 6
    reason: Missing output files: mapped_reads/B.bam
    wildcards: sample=B
    resources: tmpdir=/tmp

[Fri Feb 10 00:59:47 2023]
Finished job 6.
3 of 9 steps (33%) done
Select jobs to execute...

[Fri Feb 10 00:59:47 2023]
rule samtools_sort:
    input: mapped_reads/B.bam
    output: sorted_reads/B.bam
    jobid: 5
    reason: Missing output files: sorted_reads/B.bam; Input files updated by another job: mapped_reads/B.bam
    wildcards: sample=B
    resources: tmpdir=/tmp

[Fri Feb 10 00:59:57 2023]
Finished job 5.
4 of 9 steps (44%) done
Select jobs to execute...

[Fri Feb 10 00:59:57 2023]
rule samtools_index:
    input: sorted_reads/A.bam
    output: sorted_reads/A.bam.bai
    jobid: 7
    reason: Missing output files: sorted_reads/A.bam.bai; Input files updated by another job: sorted_reads/A.bam
    wildcards: sample=A
    resources: tmpdir=/tmp
...
```

</details>

Keep in mind that when `cleanup: true` is set, the pods will be cleaned up (with the storage)
when everything is finished, so you won't be able to access the output log. To access it you
can either stream it as it comes (e.g, flux-cloud supports this) or disable cleanup, and clean
up your volumes yourself after things are finished! When the job is finished, everything should register as completed, and you can clean up like:

```bash
$ kubectl delete -f minicluster.yaml
minicluster.flux-framework.org "flux-sample" deleted
```

And then stop minikube and the bound directory and you are good!

```bash
$ minikube stop
```

But we might not actually need to get the logs from the operator! Since our data
directory was bound to our host, we have all the output locally too:

```console
/tmp/workflow/
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
├── Dockerfile
├── environment.yaml
├── mapped_reads
│   ├── A.bam
│   └── B.bam
├── plots
│   └── quals.svg
├── README.md
├── scripts
│   └── plot-quals.py
├── Snakefile
└── sorted_reads
    ├── A.bam
    ├── A.bam.bai
    ├── B.bam
    └── B.bam.bai
```

And snakemake keeps logs in a hidden directory:

```bash
$ tree /tmp/workflow/.snakemake/
/tmp/workflow/.snakemake/
├── auxiliary
├── conda
├── conda-archive
├── incomplete
├── locks
├── log
│   └── 2023-02-10T005917.395176.snakemake.log
├── metadata
│   ├── bWFwcGVkX3JlYWRzL0EuYmFt
│   ├── bWFwcGVkX3JlYWRzL0IuYmFt
│   ├── c29ydGVkX3JlYWRzL0EuYmFt
│   ├── c29ydGVkX3JlYWRzL0EuYmFtLmJhaQ==
│   ├── c29ydGVkX3JlYWRzL0IuYmFt
│   ├── c29ydGVkX3JlYWRzL0IuYmFtLmJhaQ==
│   ├── cGxvdHMvcXVhbHMuc3Zn
│   └── Y2FsbHMvYWxsLnZjZg==
├── scripts
│   └── tmp5ppvd29r.plot-quals.py
├── shadow
└── singularity

10 directories, 10 files
```

Wicked! This is actually the first "realish" workflow we've run with the Flux Operator and local data
(that is public) so this is hugely cool to see!


## Development Notes

### Debugging Containers

If you load a container in your MiniKube cluster and need to debug it, it's useful to shell into the container already pulled there:

```bash
$ minikube ssh -- docker run -it --entrypoint bash ghcr.io/rse-ops/atacseq:app-latest
```

This will avoid the redundant pull to your host!

### Cleaning up

If you are testing and need to cleanup between runs, if you don't finish via the operator (`cleanup: true`) and need to cleanup
volumes:

```bash
kubectl delete -n flux-operator pods --all --grace-period=0 --force
kubectl delete -n flux-operator pvc --all --grace-period=0 --force
kubectl delete -n flux-operator pv --all --grace-period=0 --force
kubectl delete -n flux-operator jobs --all --grace-period=0 --force
kubectl delete -n flux-operator MiniCluster --all --grace-period=0 --force
```

That's technically the set of commands to do all-the-things! However, setting the `cleanup: true` boolean in
the Operator should handle this for you.
