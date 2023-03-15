# State

We are experimenting with saving state for the Flux Operator, which can have several different levels
of difficulty:

 - Saving state of the jobs queue and metadata, after runs are complete (between two MiniClusters)
 - Saving state of the jobs queue and metadata, pausing a queue in the middle and resuming.
 - Saving state of the jobs queue and metadata, and including filesystem (storage) assets.

These small tutorials will walk through examples of each. The most likely use cases for doing
this will be using the Flux Operator Python SDK (since we need to create multiple clusters)
in a reasonable way) but for the purposes of explanation, minicluster.yaml files are provided
as well. One important note is that since we cannot control the timing of a pod termination,
while we can have Flux automatically load a saved archive, for the process to wait for
jobs to finish and then dump the archive anew, we rely on issuing a command to the MiniCluster
(done by a script or workflow tool). This can likely be improved upon.

## Saving Pending Jobs

> Pausing scheduling and the queue in a populated queue

This example shows (via the Python SDK) how we can pause and stop a running queue and move
the jobs to a new MiniCluster to continue.

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/sdk/python/v1alpha1/examples/state-pending-jobs-minicluster.py)**

To run this example:

```bash
$ python sdk/python/v1alpha1/examples/state-pending-jobs-minicluster.py 
```

Using this example, we are able to (with slight modification) test:

 - Starting jobs on one cluster and running on another
 - Changing the size of the cluster to be larger
 - Changing the size of the cluster to be smaller

For the different cases, you can adjust the original size (and updated size) in the script
by changing the `minicluster.size`. All cases are successful to pause and resume
on the new cluster (regardless of size). Make sure (between runs) that you delete
the previous archive so you aren't loading jobs across *all* the clusters!

```bash
$ minikube ssh -- rm /tmp/data/archive.tar.gz
```

The commands we are issuing to flux are:

```bash
# Stop the queue
flux queue stop

# This should wait for running jobs to finish
flux queue idle

# And then do the dump!
flux dump /state/archive.tar.gz
```

And this means we will stop and wait for jobs to finish, and then this state is loaded
into the next cluster. If you run the example you might want to insert an `IPython.embed()`
before the delete command at the end, and then interactively shell into the new MiniCluster
(when the node are running) and then look at jobs:

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-mbv54 -- sudo -u flux flux proxy local:///var/run/flux/local flux jobs -a
```

And always make sure to clean up your archive at the end!

```bash
$ minikube ssh -- rm /tmp/data/archive.tar.gz
```

The next (basic) example goes through the same ideas, but manually for each step so you
can learn about what the script is doing.

## Basic Saving Jobs and Metadata

> Saving state of the jobs queue and metadata, after runs are complete (between two MiniClusters)

This example will walk through creating two MiniClusters - the first running a set of jobs (and finishing)
and the second cluster then loading those states. The assets for these files are in 
[examples/state/basic-job-completion](https://github.com/flux-framework/flux-operator/tree/main/examples).
Note that in order for this to work, a shared storage location is required. Since it's easier to submit
multiple jobs interactively, we will do it that way. Here is the first minicluster.yaml to create:


```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:
  # Number of pods to create for MiniCluster
  size: 2

  # Make this interactive so we can launch a bunch of jobs!
  interactive: true

  # Define the archive load/save path here (in our volume mount that persists)
  archive:
    path: /state/archive.tar.gz

  # This volume needs to persistent between MiniClusters so we can load the archive!
  volumes:
    data:
      storageClass: hostpath
      path: /tmp/data
      labels:
        type: "local"

  # This is a list because a pod can support multiple containers
  containers:
    - image: ghcr.io/flux-framework/flux-restful-api:latest
      volumes:
        data:
          path: /state
```

Note that interactive mode is set to true - this will start a broker to keep running until we decide we are done.
Since we are defining the archive path to `/state/archive.tar.gz`, this means that before Flux is started,
we will load an archive from that path given that it exists with `flux resource reload`. This is done directly
in the entrypoint. To have better control of the reverse sequence - saving the final state to that same location,
we will run `flux dump` to that same archive as an interactive command. Note this is a simple approach
that assumes we are OK replacing a previous state with a new one - for more complex workflows (where
possibly we need to maintain an original state) we likely will need to do something differently. For
the time being, let's create this first minicluster to submit jobs to, and the plan will be
that the second minicluster can load previous job history. If you are using
Minikube, make sure to pull first:

```bash
$ minikube ssh docker pull ghcr.io/flux-framework/flux-restful-api:latest
```

Now let's create it! You can either walk through this tutorial and learn about each step (continue)
below with kubectl apply) or you can run a demo script that runs the commands on your behalf:

```bash
$ /bin/bash ./examples/state/basic-job-completion/example.sh
```

<details>

<summary>View the Interactive Example Output</summary>

<script async id="asciicast-566800" src="https://asciinema.org/a/566800.js" data-speed="2"></script>


```bash
$ bash examples/state/basic-job-completion/example.sh  
```
```console
🌀️ Creating first MiniCluster...
minicluster.flux-framework.org/flux-sample created

🥱️ Sleeping 20 seconds to wait for cluster...Broker pod is flux-sample-0-qwsqw

🤓️ Contents of /tmp/data in MiniKube

✨️ Submitting jobs
ƒQK5i1V
ƒmXbRuh
ƒw92msM
ƒ27R1o9y
ƒ2JTUStw
ƒ2UhyUuD
ƒ2eVnjqH
ƒ2prDixw
ƒ2zV94Cw

🥱️ Waiting for jobs...
Jobs finished...
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
    ƒ2zV94Cw flux     whoami      S      1      -        - 
    ƒ2prDixw flux     sleep       R      1      1   3.740s flux-sample-0
    ƒ2eVnjqH flux     whoami     CD      1      1   0.042s flux-sample-0
    ƒ2UhyUuD flux     sleep      CD      1      1   4.023s flux-sample-0
    ƒ2JTUStw flux     whoami     CD      1      1   0.045s flux-sample-0
    ƒ27R1o9y flux     sleep      CD      1      1   3.022s flux-sample-0
     ƒw92msM flux     whoami     CD      1      1   0.015s flux-sample-0
     ƒmXbRuh flux     sleep      CD      1      1   2.019s flux-sample-0

🥱️ Wait a minute to be sure we have saved...

🧊️ Current state directory at /var/lib/flux...
total 4332
-rw-r--r-- 1 flux flux  151552 Mar 12 16:36 content.sqlite
-rw-r--r-- 1 flux flux 4120032 Mar 12 16:37 content.sqlite-wal
-rw-r--r-- 1 flux flux    4096 Mar 12 16:36 job-archive.sqlite
-rw-r--r-- 1 flux flux   32768 Mar 12 16:37 job-archive.sqlite-shm
-rw-r--r-- 1 flux flux  123632 Mar 12 16:37 job-archive.sqlite-wal

🧊️ Current archive directory at /state... should be empty, not saved yet
total 0
Cleaning up...
minicluster.flux-framework.org "flux-sample" deleted
total 7
-rw-rw-r-- 1 docker docker 6165 Mar 12 16:38 archive.tar.gz

🌀️ Creating second MiniCluster
minicluster.flux-framework.org/flux-sample created

🥱️ Sleeping a minute to wait for cluster...
Broker pod is flux-sample-0-jpx76

🤓️ Contents of /tmp/data in MiniKube - should be populated with archive from first
total 7
-rw-rw-r-- 1 docker docker 6165 Mar 12 16:38 archive.tar.gz

🤓️ Inspecting state directory in new cluster...
total 1308
-rw-r--r-- 1 flux flux    4096 Mar 12 16:38 content.sqlite
-rw-r--r-- 1 flux flux 1281352 Mar 12 16:38 content.sqlite-wal
-rw-r--r-- 1 flux flux    4096 Mar 12 16:38 job-archive.sqlite
-rw-r--r-- 1 flux flux   32768 Mar 12 16:38 job-archive.sqlite-shm
-rw-r--r-- 1 flux flux   12392 Mar 12 16:38 job-archive.sqlite-wal

😎️ Looking to see if old job history exists...
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
    ƒ2zV94Cw flux     whoami     CD      1      1   0.037s flux-sample-0
    ƒ2prDixw flux     sleep      CD      1      1   5.023s flux-sample-0
    ƒ2eVnjqH flux     whoami     CD      1      1   0.042s flux-sample-0
    ƒ2UhyUuD flux     sleep      CD      1      1   4.023s flux-sample-0
    ƒ2JTUStw flux     whoami     CD      1      1   0.045s flux-sample-0
    ƒ27R1o9y flux     sleep      CD      1      1   3.022s flux-sample-0
     ƒw92msM flux     whoami     CD      1      1   0.015s flux-sample-0
     ƒmXbRuh flux     sleep      CD      1      1   2.019s flux-sample-0
Cleaning up..
minicluster.flux-framework.org "flux-sample" deleted
```

</details>

### Create the MiniCluster

This is how to create the MiniCluster:

```bash
$ kubectl apply -f examples/state/basic-job-completion/minicluster.yaml 
```

At this point you can proceed to either [Interactive Submit](#interactive-submit) or [Programmatic Submit](#programmatic-submit).

### Interactive Submit

And now we need to submit a bunch of jobs to run to completion. We can do this by shelling in (and
note this could be done by the Flux Restful API for a more proggrammatic example). First,
here is how to do this interactively:

```bash
# Shell to the pod
$ kubectl exec -it -n flux-operator flux-sample-0-gzqfl -- bash
```

Check out the state directory! This is where Flux stores job metadata:

```bash
$ ls /var/lib/flux/
```
```console
content.sqlite  content.sqlite-shm  content.sqlite-wal  job-archive.sqlite  job-archive.sqlite-shm  job-archive.sqlite-wal
```

Let's now connect to the Flux instance:

```bash
$ sudo -u flux flux proxy local:///var/run/flux/local
```

And now launch a bunch of jobs. It doesn't matter what they are, go crazy.

```bash
for i in {1..5}
do
   flux submit sleep ${i}
   flux submit whoami
done
```
These will be done very quickly! You should see them all green (to indicate success) via:

```bash
$ flux jobs -a
```

### Programmatic Submit

Or just do the entire thing from the command line! First, confirm the archive path is empty:

```bash
$ minikube ssh ls /tmp/data
# no output
```

Then submit jobs - either one or many:

```bash
kubectl exec -it -n flux-operator flux-sample-0-g6gv4 -- sudo -u flux flux proxy local:///var/run/flux/local flux submit sleep 2

for i in {1..5}
do
   kubectl exec -it -n flux-operator flux-sample-0-g6gv4 -- sudo -u flux flux proxy local:///var/run/flux/local flux submit sleep ${i}
   kubectl exec -it -n flux-operator flux-sample-0-g6gv4 -- sudo -u flux flux proxy local:///var/run/flux/local flux submit whoami
done
```

When you are done, you can see all the jobs:

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-mbv54 -- sudo -u flux flux proxy local:///var/run/flux/local flux jobs -a
```

Then you can stop the queue, wait for jobs to finish, and request the dump. Note that we do this
as an interactive command and not automatically because (for large dumps 💩️) we cannot ensure that the time
will be given for completion. 

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-mbv54 -- sudo -u flux flux proxy local:///var/run/flux/local flux queue stop
$ kubectl exec -it -n flux-operator flux-sample-0-mbv54 -- sudo -u flux flux proxy local:///var/run/flux/local flux queue idle
$ kubectl exec -it -n flux-operator flux-sample-0-mbv54 -- sudo -u flux flux proxy local:///var/run/flux/local flux queue dump /state/archive.tar.gz
```

After that, outside of the shell (if you didn't already exit) let's delete the Minicluster.

```bash
$ kubectl delete -f examples/state/basic-job-completion/minicluster.yaml 
```

At this point, it should be the case that the same flux state directory is dumped to the archive path we requested, 
which is located at `/tmp/data/archive.tar.gz` in the MiniKube vm (`/tmp/data` is bound to `/state` and the
archive inside the container is asked to be saved to `/state/archive.tar.gz`).

```bash
$ minikube ssh -- ls -l /tmp/data/
total 7
-rw-rw-r-- 1 docker docker 6231 Mar 12 07:44 archive.tar.gz
```

Next, let's bring up a second minicluster! This time, in the entry point it should find the existing archive,
load into the broker, and then we will see them. We can use the same minicluster file!

```bash
$ kubectl apply -f examples/state/basic-job-completion/minicluster.yaml 
```

We then then test if this current setup has a memory of the jobs run on the first one:

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-dpd42 -- sudo -u flux flux proxy local:///var/run/flux/local flux jobs -a
```

And of course, clean up when you are done.

```bash
$ kubectl delete -f examples/state/basic-job-completion/minicluster.yaml 
$ minikube ssh -- rm -rf /tmp/data/archive.tar.gz
```

We also have this example demonstrated [entirely in Python](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/examples/state-basic-job-completion-minicluster.py) using the Flux Operator Python SDK.

> What are next steps?

This is really cool! Intuitively what we need to do for the next example (stopping a queue of jobs that are running,
meaning waiting for running jobs to finish and pausing the rest) is to submit jobs that will take much longer to run,
and then figure out how to issue a command to the cluster to stop scheduling, stop the queue, wait for running jobs to
finish, and then to do the same archive. What isn't clear is how it will work when Flux reloads the archive - will
the jobs that weren't run yet start? What commands should be the responsibility of the Operator vs. a client like
the Python SDK? I'm not sure - we will find out! 