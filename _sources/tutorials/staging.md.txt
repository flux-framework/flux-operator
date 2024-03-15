# Staging

If you have data (or a Singularity container or similar) that you want all nodes to have access to before starting your job,
Flux comes with a utility called [filemap](https://flux-framework.readthedocs.io/projects/flux-core/en/latest/man1/flux-filemap.html) 
that can make this easy to do! Note that staging will ensure the content in _unshared_ directories
across nodes (e.g., /data on all nodes which isn't a shared volume) has the same contents. These contents are not further updated or
synced. The staging (and running of these commands) must happen in a batch script, so all of these examples will use the batch: true 
directive.

## Singularity Container Staging

> Stage a singularity container for all nodes to access

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/singularity/staging/minicluster.yaml)**

This example demonstrates pulling a Singularity container to the broker pod, and having it staged across all nodes before running
a job.

```yaml
apiVersion: flux-framework.org/v1alpha2
kind: MiniCluster
metadata:
  name: flux-sample
spec:
  # Number of pods to create for MiniCluster
  size: 4
  tasks: 4

  # Job output files will be written here
  volumes:
    data:
      storageClass: hostpath
      path: /tmp/data

  # This is a list because a pod can support multiple containers
  containers:
    - image: ghcr.io/rse-ops/singularity:tag-mamba

      # original command: mpirun -n 4 singularity exec ./mpi.sif /opt/mpitest    
      command: |
        flux filemap map -C /data mpi.sif
        flux exec -x 0 -r all flux filemap get -C /data
        flux submit -n 4 --output /tmp/fluxout/job.out --error /tmp/fluxout/job.out --flags waitable singularity exec /data/mpi.sif /opt/mpitest
        flux queue idle
        flux filemap unmap

      commands:
        post: sleep infinity
        brokerPre: |
          if [[ ! -e "mpi.sif" ]]; then
              singularity pull mpi.sif oras://ghcr.io/rse-ops/singularity-mpi:mpich
          fi

      workingDir: /data
      cores: 1

      # Output files written here
      volumes:
        data:
          path: /tmp/fluxout

      # Batch, and don't wrap in flux submit (we will do this)
      batch: true
      batchRaw: true

      fluxUser:
        name: fluxuser
       
      # Running a container in a container
      securityContext:
        privileged: true
```

Let's break down the above!


### Staging

The staging (and then running a subsequent job) happens in the batch script. Note that we have also set "batchRaw" to true,
which tells the Flux operator "Don't wrap our commands in submits." This means we need to write the entire batch script
on our own. It's a bit like an expert mode! Let's go through it in detail:

```bash
# Make mpi.sif in /data available to other nodes via mmap
flux filemap map -C /data mpi.sif

# Map the file to all nodes, but skip rank 0, since the file is already there
flux exec -x 0 -r all flux filemap get -C /data

# When it's done, submit the job (also to 4 nodes) and ensure the output/error files are written to the shared mounted volume /tmp/fluxout
# The waitable flag ensures the next command will wait for this job
flux submit -n 4 --output /tmp/fluxout/job.out --error /tmp/fluxout/job.out --flags waitable singularity exec /data/mpi.sif /opt/mpitest

# This is important to have so we wait for jobs to finish!
flux queue idle

# Clean up the filemap
flux filemap unmap
```

### Volumes

Since we want the output file to be written to in a single location by all workers, we write that to the shared mount `/tmp/fluxout`,
which is the default output location. Finally, since we are running a Singularity container, we run with privileged true.

### Commands

The broker pre command is going to pull the container _once_ from an OCI registry (GitHub packages, ghcr.io) via ORAS (OCI Registry
as Storage) to its local filesystem. This SIF binary is going to be mapped to the other workers before running it.
Since we want to keep the cluster running after it finishes, we create a "post" command with "sleep infinity." This will
allow us to shell inside and look at output, etc. More realistically if you run a batch job it's suggested to send output
to some kind of service, and not put a dependency on the cluster running.

### Running the Workflow

Okay let's run the example!

```bash
$ kubectl create namespace flux-operator
$ kubectl apply -f examples/singularity/staging/minicluster.yaml
```

We can then wait for our pods to be running

```bash
$ kubectl get -n flux-operator pods
```

And then look at the logs to see the container being pulled:

```bash
$ kubectl logs -n flux-operator flux-sample-0-p5xls -f
```

Next, let's shell into the broker to look at the output log. This will show the job running.

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-8rlps bash
```
```bash
$ cat /tmp/fluxout/job.out 
```
```console
INFO:    Converting SIF file to temporary sandbox...
INFO:    Converting SIF file to temporary sandbox...
INFO:    Converting SIF file to temporary sandbox...
INFO:    Converting SIF file to temporary sandbox...
WARNING: underlay of /etc/localtime required more than 50 (81) bind mounts
WARNING: underlay of /etc/localtime required more than 50 (81) bind mounts
WARNING: underlay of /etc/localtime required more than 50 (81) bind mounts
WARNING: underlay of /etc/localtime required more than 50 (81) bind mounts
Hello, I am rank 0/4
INFO:    Cleaning up image...
Hello, I am rank 3/4
Hello, I am rank 1/4
Hello, I am rank 2/4
INFO:    Cleaning up image...
INFO:    Cleaning up image...
INFO:    Cleaning up image...
```
The above only works because we used mmap to share the SIF with the other workers! If we don't
do the unmap command, you could also shell into any of the other workers and see it in `/data`.
Remember the mapping will only work if the directory is not shared. You could also share a directory
(with read/write) between the workers, and then pull once if you are able to.

When you are done, clean up:

```bash
$ kubectl delete -f examples/singularity/staging/minicluster.yaml
```
