# Multiple Applications in a Pod

This example is related to [multiple-pods-per-node](../multiple-pods-per-node) except we are testing a flipped variant - having several application containers that are submitting to the same
flux queue. For this example, we will use two simple applications from our examples - LAMMPS and the OSU benchmarks. We will try:

1. Creating an interactive MiniCluster that has a shared Flux install
2. Two containers that each are running Flux.

The key for the above is that while the two containers both have flux (meaning the view is mounted and available) only _one_ will start the flux broker and see the entire resources of the node.

## Experiment

### Create the Cluster

We should be able to use kind for this.

```bash
kind create cluster --config ../../kind-config.yaml
```

### Install the Flux Operator

As follows:

```bash
kubectl apply -f https://raw.githubusercontent.com/flux-framework/flux-operator/main/examples/dist/flux-operator.yaml
```

Note that I ran into issues with kind and the containers pulling - it

### Fluxion Application Scheduler

#### 1. Create the MiniCluster

Then create the flux operator pods!

```bash
kubectl apply -f minicluster.yaml
```

Wait for everything to be running:

```console
NAME                   READY   STATUS    RESTARTS   AGE
flux-sample-0-zgsgp    4/4     Running   0          15m
flux-sample-1-gqdf9    4/4     Running   0          15m
flux-sample-2-l7774    4/4     Running   0          15m
flux-sample-3-8whls    4/4     Running   0          15m
flux-sample-services   1/1     Running   0          15m
```

Here is the neat thing - each container running inside each pod is an independent broker that sees all resources! The lead broker (for each) is at index 0. You can confirm this by selecting to see logs for any specific container:

```bash
# This is running the queue orchestrator
kubectl logs flux-sample-0-zgsgp -c queue

# These are application containers
kubectl logs flux-sample-0-zgsgp -c lammps
kubectl logs flux-sample-0-zgsgp -c chatterbug
kubectl logs flux-sample-0-zgsgp -c ior
```

And this is the fluxion graph server, which is running as the scheduler for the entire cluster!

```bash
$ kubectl logs flux-sample-services 
🦩️ This is the fluxion graph server
[GRPCServer] gRPC Listening on [::]:4242
```

#### 2. Load the bypass plugin

When the "queue" broker comes up, it loads a plugin on each of the application brokers that
ensures we can give scheduling decisions directly to those brokers from the fluxion service:

```bash
for socket in $(ls /mnt/flux/view/run/flux/)
  do
  flux proxy local:///mnt/flux/view/run/flux/$socket flux jobtap load alloc-bypass.so 
done
```

This will allow us to bypass the scheduler, and pass forward exactly the decision from fluxion. We do this so that
we can schedule down to the CPU and not have resource oversubscription. When all the containers are running and the queue starts, you should see:

```bash
job-manager.err[0]: jobtap: job.new: callback returned error
⭐️ Found application      queue: index 0
⭐️ Found application chatterbug: index 3
⭐️ Found application        ior: index 2
⭐️ Found application     lammps: index 1
✅️ Init of Fluxion resource graph success!
 * Serving Flask app 'fluxion_controller'
 * Debug mode: off
WARNING: This is a development server. Do not use it in a production deployment. Use a production WSGI server instead.
 * Running on all addresses (0.0.0.0)
 * Running on http://127.0.0.1:5000
 * Running on http://10.244.0.50:5000
Press CTRL+C to quit
```

We are ready to submit jobs!

#### 3. Connect to the Queue

The "queue" container of the set is special because it doesn't have a specific application - it's mostly a thin layer provided to interact with other containers (and we will run our application to handle orchestration from there).
So let's shell into this controller pod container - which is the one that doesn't have an application, but has access to all the resources available! Let's shell in:

```bash
kubectl exec -it flux-sample-0-xxx bash
```
```console
kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl exec [POD] -- [COMMAND] instead.
Defaulted container "queue" out of: queue, lammps, ior, chatterbug, flux-view (init)
```

Notice in the message above we see all the containers running - and we are shelling into the first (queue). Also note that since we are installing some Python stuff, 
you need to wait for that to finish before you see the flux socket for the queue show up. When it's done, it will be the index 0 here:

```bash
ls /mnt/flux/view/run/flux/
local-0  local-1  local-2  local-3
```
The indices correspond with the other containers. You can see the mapping here in the "meta" directory:

```bash
ls /mnt/flux/view/etc/flux/meta/
0-queue  1-lammps  2-ior  3-chatterbug
```

That's a pretty simple (maybe dumb) approach, but it will be how we get names for the containers when we run the Fluxion controller. Let's do that next!

#### 4. Connect to Fluxion

Since we can see all the instances running, this allows us to easily (meaning programatically) write a script that orchestrates interactions between the different brokers, where there is one broker per container set, where each container set is
running across all pods (physical nodes). Since we are doing this interactively, let's connect to the queue broker. It doesn't actually matter, and theoretically this could run from any container that is running Flux.
If this script that (TBA written) is run on startup, we won't need to do this.

```bash
source /mnt/flux/flux-view.sh 
flux proxy $fluxsocket bash
flux resource list
```
```console
[root@flux-sample-0 /]# flux resource list
     STATE NNODES   NCORES    NGPUS NODELIST
      free      4       40        0 flux-sample-[0-3]
 allocated      0        0        0 
      down      0        0        0 
```

What we are seeing in the above is the set of resources that need to be shared across the containers (brokers). We don't want to oversubscribe, or for example, tell any specific broker that it can use all the resources while we tell the same to the others. We have to be careful that we use the Python install that is alongside the Flux install. Note that *you should not run this* but I want to show you how the queue was started. You can issue `--help` to see all the options to customize:

```bash
/mnt/flux/view/bin/python3.11 /mnt/flux/view/fluxion_controller.py start --help

# This is how it was started using the defaults (do not run this again)
/mnt/flux/view/bin/python3.11 /mnt/flux/view/fluxion_controller.py start
```

To submit a job, (and you can do this from any of the flux container brokers) - it will be hitting a web service that the Python script is exposing from the queue!

```bash
/mnt/flux/view/bin/python3.11 /mnt/flux/view/fluxion_controller.py submit --help

# The command is the last bit here (ior)                                           # command
/mnt/flux/view/bin/python3.11 /mnt/flux/view/fluxion_controller.py submit --cpu 4 --container ior ior
```

And then we see from where we submit:

```console
⭐️ Found application      queue: index 0
⭐️ Found application chatterbug: index 3
⭐️ Found application        ior: index 2
⭐️ Found application     lammps: index 1
{'annotations': {}, 'bank': '', 'container': 'ior', 'cwd': '', 'dependencies': [], 'duration': 3600.0, 'exception': {'note': '', 'occurred': False, 'severity': '', 'type': ''}, 'expiration': 0.0, 'fluxion': 1, 'id': 19001371525120, 'name': 'ior', 'ncores': 4, 'nnodes': 1, 'nodelist': 'flux-sample-3', 'ntasks': 1, 'priority': 16, 'project': '', 'queue': '', 'ranks': '3', 'result': 'COMPLETED', 'returncode': 0, 'runtime': 0.5983412265777588, 'state': 'INACTIVE', 'status': 'COMPLETED', 'success': True, 't_cleanup': 1719904964.2517486, 't_depend': 1719904963.6396549, 't_inactive': 1719904964.254762, 't_remaining': 0.0, 't_run': 1719904963.6534073, 't_submit': 1719904963.6277533, 'urgency': 16, 'userid': 0, 'username': 'root', 'waitstatus': 0}
```

And from the Fluxion service script:

```console
{'command': ['ior'], 'cpu': '4', 'container': 'ior'}
🙏️ Requesting to submit: ior
✅️ Match of jobspec to Fluxion graph success!
10.244.0.18 - - [01/Jul/2024 23:55:42] "POST /submit HTTP/1.1" 200 -
👉️ Job on ior 1 is complete.
✅️ Cancel of jobid 1 success!
```

I am calling this "pancake elasticity" since we can theoretically deploy many application containers and then use them when needed, essentially expanding the one running out (resource wise) while the others remain flat (not using resources). This isn't entirely ready yet (still testing) but a lot of the automation is in place.

It's so super cool!! :D This is going to likely inspire the next round of work for thinking about scheduling and fluxion.
