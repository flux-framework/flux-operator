# Ramble

Let's test [ramble](https://github.com/GoogleCloudPlatform/ramble)
in the Flux operator. This is just a hello world experiment that runs
the hostname. I did test gromacs but hit segfaults and ran away (the workspace
is in the container for anyone interested)!

> I would recommend this for groups that are already using spack for running workflows (e.g., spack environments)
> I would not recommend ramble if you are starting fresh and have a choice of what tools to use.
> Either way, I recommend that you try it out. I like the idea of having recipes for reproducible workflows, conceptually.

## Usage

When the cluster is created and the operator installed,
create volumes first:

```bash
$ kubectl create -f  volumes.yaml
```

Then the cluster:

```bash
$ kubectl apply -f ./minicluster.yaml
```

This will create the MiniCluster, and you can wait for the pods to be running (it took over 7 minutes for me) to pull because the container is big! Then copy the ramble-on.sh script:

```bash
pod=flux-sample-0-7st8l
kubectl cp ./ramble-on.sh ${pod}:/tmp/workflow/ramble-on.sh
```
Note that the volume on the host ensures that all pods can see the script and workspace!
Then shell in and connect to the instance and run the job:

```bash
kubectl exec -it $pod bash
. /mnt/flux/flux-view.sh
flux proxy $fluxsocket bash
flux run -N 4 /bin/bash /tmp/workflow/ramble-on.sh
```

Then watch ramble run!

```bash
==>     Executing phase get_inputs
==>     Executing phase make_experiments
==>     Executing phase get_inputs
==>     Executing phase make_experiments
==> Warning: The env-vars workspace section is deprecated. Environment variables
==>     Executing phase get_inputs
should be defined in the env_vars config section using the same
==>     Executing phase make_experiments
syntax. Support for env-vars will be removed in a future. See
==>     Executing phase get_inputs
the documentation for examples of the new syntax.
==>     Executing phase make_experiments
flux-sample-1
0.00user 0.00system 0:00.00elapsed 0%CPU (0avgtext+0avgdata 1380maxresident)k
0inputs+0outputs (0major+70minor)pagefaults 0swaps
flux-sample-3
0.00user 0.00system 0:00.00elapsed 0%CPU (0avgtext+0avgdata 1484maxresident)k
0inputs+0outputs (0major+71minor)pagefaults 0swaps
flux-sample-2
0.00user 0.00system 0:00.00elapsed 0%CPU (0avgtext+0avgdata 1432maxresident)k
0inputs+0outputs (0major+72minor)pagefaults 0swaps
flux-sample-0
0.00user 0.00system 0:00.00elapsed 0%CPU (0avgtext+0avgdata 1432maxresident)k
```

I would recommend ramble for groups that are already using Spack and want to put some structure and reproducibility to their workflows,
but would not recommend it for starting out fresh.  If you want to extend this example, you can work
from the Dockerfile build [here](https://github.com/rse-ops/flux-hpc/blob/main/ramble-gromacs/Dockerfile)
to prepare a custom container. Note that the build takes about 3 hours. 
If you want to debug something, you can set interactive: true to run in interactive mode, and then shell into the pod, connect to the broker:

```bash
kubectl exec -it flux-sample-0-jlsp6 bash
. /mnt/flux/flux-view.sh
flux proxy $fluxsocket bash
```

And when you are done, clean up:

```bash
kubectl delete -f minicluster.yaml
kubectl delete -f volumes.yaml
```

