# CSI for Cloud Storage

This basic tutorial will walk through creating a MiniCluster to run a Snakemake workflow! 
We will use a "Container Storage Interface" (CSI) to connect to Snakemake assets in Google Cloud Storage.
You should have already [setup your workspace](setup.md), including preparing the Snakemake data in
Google Storage.


```{include} includes/basic-setup.md
```


Akin to how we created a local volume, we can do something similar, but instead of pointing the Flux Operator
to a volume on the host (e.g., in MiniKube) we are going to point it to a storage bucket with our data.
For this tutorial, we will use the [csi-gcs driver](https://github.com/ofek/csi-gcs) to connect a cloud storage bucket to our cluster.

## Install the CSI

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

## Storage Class

We can then create our storage class, this file is provided in `examples/storage/google/storageclass.yaml`

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-gcs
provisioner: gcs.csi.ofek.dev
```

```bash
$ kubectl apply -f examples/storage/google/gcs-csi/storageclass.yaml
```

## Snakemake MiniCluster

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
      claimAnnotations:
        gcs.csi.ofek.dev/location: us-central1
        gcs.csi.ofek.dev/project-id: my-project
        gcs.csi.ofek.dev/bucket: flux-operator-storage
```

Also note that we are setting the `commands: -> runFluxAsRoot` to true. This isn't ideal, but it was the
only way I could get the storage to both be seen and have permission to write there. Let's create the job!

```bash
$ kubectl apply -f examples/storage/google/gcs-csi/minicluster.yaml
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

```{include} includes/cleanup.md
```