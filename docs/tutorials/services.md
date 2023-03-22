# Services

These tutorials will show you how to run a "sidecar" service container (one per Flux node) alongside your
flux install, along with a service for the entire cluster (a deployment next to the cluster).

## Sidecar Tutorials

### Sidecar NGINX Container

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/tests/nginx-service/minicluster.yaml)**

This is one of the simplest examples, implemented as a test, to run a sidecar with NGINX and then curl localhost
to get a response from flux.

```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:

  logging:
    quiet: true

  # Number of pods to create for MiniCluster
  size: 2

  # This is a list because a pod can support multiple containers
  containers:
    - image: ghcr.io/flux-framework/flux-restful-api:latest
      runFlux: true
      command: curl -s localhost
      commands:
        pre: apt-get update > /dev/null && apt-get install -y curl > /dev/null
    - image: nginx
      name: nginx
      ports:
        - 80
```

Create it (after you have the flux-operator namespace):

```bash
$ kubectl create -f ./examples/tests/nginx-service/minicluster.yaml
```

See nginx is running:

```bash
$ kubectl -n flux-operator logs flux-sample-0-zlpwx -c nginx -f
```
```console
/docker-entrypoint.sh: /docker-entrypoint.d/ is not empty, will attempt to perform configuration
/docker-entrypoint.sh: Looking for shell scripts in /docker-entrypoint.d/
/docker-entrypoint.sh: Launching /docker-entrypoint.d/10-listen-on-ipv6-by-default.sh
10-listen-on-ipv6-by-default.sh: info: Getting the checksum of /etc/nginx/conf.d/default.conf
10-listen-on-ipv6-by-default.sh: info: Enabled listen on IPv6 in /etc/nginx/conf.d/default.conf
/docker-entrypoint.sh: Launching /docker-entrypoint.d/20-envsubst-on-templates.sh
/docker-entrypoint.sh: Launching /docker-entrypoint.d/30-tune-worker-processes.sh
/docker-entrypoint.sh: Configuration complete; ready for start up
2023/03/18 05:01:31 [notice] 1#1: using the "epoll" event method
2023/03/18 05:01:31 [notice] 1#1: nginx/1.23.3
2023/03/18 05:01:31 [notice] 1#1: built by gcc 10.2.1 20210110 (Debian 10.2.1-6) 
2023/03/18 05:01:31 [notice] 1#1: OS: Linux 5.15.0-67-generic
2023/03/18 05:01:31 [notice] 1#1: getrlimit(RLIMIT_NOFILE): 1048576:1048576
2023/03/18 05:01:31 [notice] 1#1: start worker processes
2023/03/18 05:01:31 [notice] 1#1: start worker process 29
2023/03/18 05:01:31 [notice] 1#1: start worker process 30
2023/03/18 05:01:31 [notice] 1#1: start worker process 31
...
```

And then look at the main logs to see the output of curl, run by Flux:

```bash
$ kubectl -n flux-operator logs flux-sample-0-zlpwx -c flux-sample -f
```
```console
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
html { color-scheme: light dark; }
body { width: 35em; margin: 0 auto;
font-family: Tahoma, Verdana, Arial, sans-serif; }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

And that's it! In a real world use case, you'd have some service running alongside
an analysis. Clean up:

```bash
$ kubectl delete -f ./examples/tests/nginx-service/minicluster.yaml
```

### Sidecar Registry with ORAS

> Create an interactive MiniCluster with a sidecar registry container

As an example, we will run a local container registry to push/pull artifacts
with ORAS. I don't know why, I just like ORAS :) In all seriousness, you could
imagine interesting use cases like needing an API to save and get artifacts 
for your analysis.

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/services/minicluster-registry.yaml)**

This example demonstrates bringing up a MiniCluster and then interacting with a service (a registry)
to push / pull artifacts. Here is our example custom resource definition:

```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:

  # Number of pods to create for MiniCluster
  size: 2

  # Interactive so we can submit commands
  interactive: true

  # This is a list because a pod can support multiple containers
  containers:
    # The container URI to pull (currently needs to be public)
    - image: ghcr.io/flux-framework/flux-restful-api:latest
      runFlux: true
      commands:

        # This is going to install oras for the demo
        pre: |
          apt-get update && apt-get install -y curl
          VERSION="1.0.0-rc.2"
          curl -LO "https://github.com/oras-project/oras/releases/download/v${VERSION}/oras_${VERSION}_linux_amd64.tar.gz"
          mkdir -p oras-install/
          tar -zxf oras_${VERSION}_*.tar.gz -C oras-install/
          sudo mv oras-install/oras /usr/local/bin/
          rm -rf oras_${VERSION}_*.tar.gz oras-install/

      # This is our registry we want to run
    - image: ghcr.io/oras-project/registry:latest
      name: registry
      ports:
        - 5000

```

It's helpful to pull containers to MiniKube first:

```bash
$ minikube ssh docker pull ghcr.io/oras-project/registry:latest
$ minikube ssh docker pull ghcr.io/flux-framework/flux-restful-api:latest
```

When interactive is true, we tell the Flux broker to start without a command. This means
the cluster will remain running until you shutdown Flux or `kubectl delete` the MiniCluster
itself. The container you choose should have the software you are interested in having for each node.
Given a running cluster, we can create the namespace and the MiniCluster as follows:

```bash
$ kubectl create namespace flux-operator
```

And apply the MiniCluster CRD:
```bash
$ kubectl apply -f examples/services/minicluster-registry.yaml
```

If you are curious, the entrypoint for the service sidecar container is `registry serve /etc/docker/registry/config.yml`
to start the registry. Since it's not a flux runner, not providing an entrypoint means we use the container's default
entrypoint. We can then wait for our pods to be running

```bash
$ kubectl get -n flux-operator pods
NAME                         READY   STATUS      RESTARTS   AGE
flux-sample-0-p5xls          1/1     Running     0          7s
flux-sample-1-nmtt7          1/1     Running     0          7s
flux-sample-cert-generator   0/1     Completed   0          7s
```

To see logs, since we have 2 containers per pod, you can either leave out the pod (and get the first or default)
or specify a container with `-c`:

```bash
$ kubectl logs -n flux-operator flux-sample-0-d5jbb -c registry
$ kubectl logs -n flux-operator flux-sample-0-d5jbb -c flux-sample
$ kubectl logs -n flux-operator flux-sample-0-d5jbb
```

And then shell into the broker pod, index 0, which is "flux-sample"

```bash
$ kubectl exec -it  -n flux-operator flux-sample-0-d5jbb -c flux-sample -- bash
```

Let's first make and push an artifact. First, just using oras natively (no flux)

```bash
cd /tmp

# Assume we would be running from inside the flux instance
sudo -u flux echo "hello dinosaur" > artifact.txt
```

And push! The registry, by way of being a container in the same pod, is on port 5000:

At this point, remember the broker is running, and we need to connect to it. We do this via
flux proxy and targeting the socket, which is a local reference at `/run/flux/local`:

```bash
# Connect to the flux socket at /run/flux/local as the flux instance owner "flux"
$ sudo -u flux flux proxy local:///run/flux/local oras push localhost:5000/dinosaur/artifact:v1 \
   --artifact-type application/vnd.acme.rocket.config \
   ./artifact.txt
```
```console
Uploading 07f469745bff artifact.txt
Uploaded  07f469745bff artifact.txt
Pushed [registry] localhost:5000/dinosaur/artifact:v1
Digest: sha256:3a6cb1d1d1b1d80d4c4de6abc66a6c9b4f7fef0b117f87be87fea9b725053ead
```
Now try pulling, deleting the original first, and again without flux:

```bash
rm -f artifact.txt
sudo -u flux flux proxy local:///run/flux/local oras pull localhost:5000/dinosaur/artifact:v1
cat artifact.txt
```
```console
hello dinosaur
```

We did this under the broker (and flux user) assuming your actual use case will be running
in the Flux instance. Feel free to play with oras outside of that context!
When you are done, exit from the instance, and exit from the pod, and then delete the MiniCluster.

```bash
$ kubectl delete -f examples/services/minicluster-registry.yaml
```

That's it. Please do something more useful than my terrible example.

## Service Containers Alongside the Cluster

### Merlin Demo Workflow

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/launchers/merlin/minicluster.yaml)**

This example will run a small "hello world" tutorial bringing up a rabbitmq and redis container from Flux.
This could (eventually) be part of a workflows tutorial, but for now can exemplify pod services alongside a MiniCluster.
Note that we derive this example from [this repository](https://github.com/rse-ops/flux-hpc/tree/main/merlin-demos)
and use a [customized set of containers](https://github.com/rse-ops/flux-hpc/tree/main/merlin-demos-certs) that have
certificates built into the containers. This is obviously not recommended and only used here for an example.
This assumes we have minikube running:

```bash
$ minikube start
```

And the Flux Operator namespace created:

```bash
$ kubectl create -n flux-operator
```

First, pull the containers to MiniKube:

```bash
$ minikube ssh docker pull ghcr.io/rse-ops/merlin-demos-certs:merlin
$ minikube ssh docker pull ghcr.io/rse-ops/merlin-demos-certs:rabbitmq
$ minikube ssh docker pull ghcr.io/rse-ops/merlin-demos-certs:redis
```

And then generate the (separate) pods to run redis and rabbitmq in the flux-operator namespace.
The containers already have shared certificates (just for this test case)!

```bash
$ kubectl apply -f examples/launchers/merlin/services.yaml
```

And create the MiniCluster to use them!

```bash
$ kubectl apply -f examples/launchers/merlin/minicluster.yaml
```
The MiniCluster is created in interactive mode, and we do this because merlin isn't designed to
run and hang until it's done. If we don't run in interactive mode we will miss the execution
that is created via flux alloc. Look at pods running:

```bash
kubectl get -n flux-operator pods
```
```console
NAME                         READY   STATUS      RESTARTS   AGE
flux-sample-0-774tg          1/1     Running     0          22s
flux-sample-1-k24tq          1/1     Running     0          22s
flux-sample-cert-generator   0/1     Completed   0          22s
rabbitmq-f8c84d986-262pg     1/1     Running     0          32m
redis-c9469b9c5-cnhgl        1/1     Running     0          32m
```

Check the logs of the broker - you should see that the merlin example and tasks were created:

<details>

<summary>Expected merlin output at top of broker log</summary>

```bash
$ kubectl logs -n flux-operator flux-sample-0-f9ts7 -f
```
```console
Flux username: fluxuser
  
                                                
       *      
   *~~~~~                                       
  *~~*~~~*      __  __           _ _       
 /   ~~~~~     |  \/  |         | (_)      
     ~~~~~     | \  / | ___ _ __| |_ _ __  
    ~~~~~*     | |\/| |/ _ \ '__| | | '_ \ 
   *~~~~~~~    | |  | |  __/ |  | | | | | |
  ~~~~~~~~~~   |_|  |_|\___|_|  |_|_|_| |_|
 *~~~~~~~~~~~                                    
   ~~~*~~~*    Machine Learning for HPC Workflows                                 
              


[2023-03-22 20:32:25: INFO] Copying example 'flux_par' to /workflow/flux
  
                                                
       *      
   *~~~~~                                       
  *~~*~~~*      __  __           _ _       
 /   ~~~~~     |  \/  |         | (_)      
     ~~~~~     | \  / | ___ _ __| |_ _ __  
    ~~~~~*     | |\/| |/ _ \ '__| | | '_ \ 
   *~~~~~~~    | |  | |  __/ |  | | | | | |
  ~~~~~~~~~~   |_|  |_|\___|_|  |_|_|_| |_|
 *~~~~~~~~~~~                                    
   ~~~*~~~*    Machine Learning for HPC Workflows                                 
              


[2023-03-22 20:32:26: INFO] Loading specification from path: /workflow/flux/flux_par.yaml
[2023-03-22 20:32:26: INFO] Made dir(s) to output path '/workflow/studies'.
[2023-03-22 20:32:26: INFO] Study workspace is '/workflow/studies/flux_par_20230322-203226'.
[2023-03-22 20:32:26: INFO] Reading app config from file /home/fluxuser/.merlin/app.yaml
[2023-03-22 20:32:26: INFO] Generating samples...
[2023-03-22 20:32:26: INFO] Generating samples complete!
[2023-03-22 20:32:26: INFO] Loading samples from 'samples.npy'...
[2023-03-22 20:32:26: INFO] 10 samples loaded.
[2023-03-22 20:32:26: INFO] Calculating task groupings from DAG.
[2023-03-22 20:32:26: INFO] Converting graph to tasks.
[2023-03-22 20:32:26: INFO] Launching tasks.
```

</details>

At this point, we can interactively shell in to look around. This would also ideally be tweaked to
be run automatically (if merlin had a "watch" or similar functionality):

```bash
$ kubectl exec -n flux-operator  -it flux-sample-0-f9ts7 -- bash
```

Connect to the flux broker:

```bash
$ export PYTHONPATH=$PYTHONPATH:/home/fluxuser/.local/lib/python3.10/site-packages
$ sudo -E LD_LIBRARY_PATH=$LD_LIBRARY_PATH -E PATH=$PATH -E HOME=/home/fluxuser -E PYTHONPATH=$PYTHONPATH -u fluxuser flux proxy local:///var/run/flux/local
```

At this point, we are the fluxuser! We can test if merlin is connected to its services via `merlin info`:

```bash
$ whoami
fluxuser
```
```bash
$ merlin info
```

<details>

<summary>Output of Merlin Info</summary>

```console
[2023-03-22 20:41:45: INFO] Reading app config from file /home/fluxuser/.merlin/app.yaml
  
                                                
       *      
   *~~~~~                                       
  *~~*~~~*      __  __           _ _       
 /   ~~~~~     |  \/  |         | (_)      
     ~~~~~     | \  / | ___ _ __| |_ _ __  
    ~~~~~*     | |\/| |/ _ \ '__| | | '_ \ 
   *~~~~~~~    | |  | |  __/ |  | | | | | |
  ~~~~~~~~~~   |_|  |_|\___|_|  |_|_|_| |_|
 *~~~~~~~~~~~                                    
   ~~~*~~~*    Machine Learning for HPC Workflows                                 
              


Merlin Configuration
-------------------------

 config_file        | /home/fluxuser/.merlin/app.yaml
 is_debug           | False
 merlin_home        | /home/fluxuser/.merlin
 merlin_home_exists | True
 broker server      | amqps://fluxuser:******@rabbitmq:5671//merlinu
 broker ssl         | {'keyfile': '/cert_rabbitmq/client_rabbitmq_key.pem', 'certfile': '/cert_rabbitmq/client_rabbitmq_certificate.pem', 'ca_certs': '/cert_rabbitmq/ca_certificate.pem', 'cert_reqs': <VerifyMode.CERT_REQUIRED: 2>}
 results server     | rediss://redis:6379/0
 results ssl        | {'ssl_keyfile': '/cert_redis/client_redis_key.pem', 'ssl_certfile': '/cert_redis/client_redis_certificate.pem', 'ssl_ca_certs': '/cert_redis/ca_certificate.pem', 'ssl_cert_reqs': <VerifyMode.CERT_REQUIRED: 2>}

Checking server connections:
----------------------------
broker server connection: OK
results server connection: OK

Python Configuration
-------------------------

 $ which python3
/opt/conda/bin/python3

 $ python3 --version
Python 3.10.9

 $ which pip3
/opt/conda/bin/pip3

 $ pip3 --version
pip 23.0 from /opt/conda/lib/python3.10/site-packages/pip (python 3.10)

"echo $PYTHONPATH"
/usr/lib/flux/python3.1:/home/fluxuser/.local/lib/python3.10/site-packages
```

</details>

Now let's queue our tasks (I guess we already did this, but I did it again)
and run the workers! Note that this is a test - so I edited 
`vim /home/fluxuser/.local/lib/python3.10/site-packages/merlin/study/batch.py`
to use `flux run` instead of `flux alloc`. E.g.,:

```diff
- launch_command = f"{flux_exe} mini alloc -o pty -N {nodes} --exclusive --job-name=merlin"
+ launch_command = f"{flux_exe} run"
```

(Yes, this isn'[t ready for production - we are figuring it out!])

```bash
$ merlin run flux/flux_par.yaml
$ merlin run-workers flux/flux_par.yaml
```

Note that if you need to re-run (or try again) you need to purge the queues:

```bash
# Do this until you see no messages
$ merlin purge flux/flux_par.yaml -f

# This can be run once after the above
$ flux job purge --force --num-limit=0
```

You'll see the tasks pick up in celery and a huge stream of output! You can also
look at the logs of redis and rabbitmq (the other containers) to see them receiving data.

<details>

<summary>Redis receiving data</summary>

```bash
$ kubectl logs -n flux-operator redis-c9469b9c5-cnhgl -f
```
```console
7:C 22 Mar 2023 19:58:09.512 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
7:C 22 Mar 2023 19:58:09.512 # Redis version=7.0.10, bits=64, commit=00000000, modified=0, pid=7, just started
7:C 22 Mar 2023 19:58:09.512 # Configuration loaded
7:M 22 Mar 2023 19:58:09.514 * monotonic clock: POSIX clock_gettime
                _._                                                  
           _.-``__ ''-._                                             
      _.-``    `.  `_.  ''-._           Redis 7.0.10 (00000000/0) 64 bit
  .-`` .-```.  ```\/    _.,_ ''-._                                  
 (    '      ,       .-`  | `,    )     Running in standalone mode
 |`-._`-...-` __...-.``-._|'` _.-'|     Port: 6379
 |    `-._   `._    /     _.-'    |     PID: 7
  `-._    `-._  `-./  _.-'    _.-'                                   
 |`-._`-._    `-.__.-'    _.-'_.-'|                                  
 |    `-._`-._        _.-'_.-'    |           https://redis.io       
  `-._    `-._`-.__.-'_.-'    _.-'                                   
 |`-._`-._    `-.__.-'    _.-'_.-'|                                  
 |    `-._`-._        _.-'_.-'    |                                  
  `-._    `-._`-.__.-'_.-'    _.-'                                   
      `-._    `-.__.-'    _.-'                                       
          `-._        _.-'                                           
              `-.__.-'                                               

7:M 22 Mar 2023 19:58:09.516 # Server initialized
7:M 22 Mar 2023 19:58:09.517 * Ready to accept connections
7:M 22 Mar 2023 20:43:59.338 * 100 changes in 300 seconds. Saving...
7:M 22 Mar 2023 20:43:59.338 * Background saving started by pid 22
22:C 22 Mar 2023 20:43:59.341 * DB saved on disk
22:C 22 Mar 2023 20:43:59.342 * Fork CoW for RDB: current 0 MB, peak 0 MB, average 0 MB
7:M 22 Mar 2023 20:43:59.438 * Background saving terminated with success
```

</details>

And for rabbitmq:

<details>

<summary>Rabbitmq receiving data</summary>

This view is truncated to just show the bottom - the log was too big to reasonably include.

```bash
$ kubectl logs -n flux-operator rabbitmq-f8c84d986-262pg -f
```
```console
...
2023-03-22 20:45:07.927630+00:00 [info] <0.2896.0> connection <0.2896.0> (172.17.0.1:41672 -> 172.17.0.5:5671): user 'fluxuser' authenticated and granted access to vhost '/merlinu'
2023-03-22 20:45:07.929943+00:00 [info] <0.2896.0> closing AMQP connection <0.2896.0> (172.17.0.1:41672 -> 172.17.0.5:5671, vhost: '/merlinu', user: 'fluxuser')
2023-03-22 20:45:07.931598+00:00 [info] <0.2882.0> closing AMQP connection <0.2882.0> (172.17.0.1:13790 -> 172.17.0.5:5671, vhost: '/merlinu', user: 'fluxuser')
2023-03-22 20:45:09.052810+00:00 [info] <0.1144.0> closing AMQP connection <0.1144.0> (172.17.0.1:49018 -> 172.17.0.5:5671, vhost: '/merlinu', user: 'fluxuser')
2023-03-22 20:45:09.057481+00:00 [info] <0.1133.0> closing AMQP connection <0.1133.0> (172.17.0.1:50219 -> 172.17.0.5:5671, vhost: '/merlinu', user: 'fluxuser')
```
</details>

We can look to see that flux job ran (and completed)

```bash
$ flux jobs -a
```
```console
$ flux jobs -a
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
   ƒDJWS5XrX fluxuser merlin     CD      1      1    1.22m flux-sample-3
```

At this point we can look at output (I installed tree 🌲️)!

```bash
$ tree studies
```
```console
studies
└── flux_par_20230322-213613
    ├── build
    │   ├── MERLIN_FINISHED
    │   ├── build.out
    │   ├── build.slurm.err
    │   ├── build.slurm.out
    │   ├── build.slurm.sh
    │   └── mpi_hello
    ├── merlin_info
    │   ├── cmd.err
    │   ├── cmd.out
    │   ├── cmd.sh
    │   ├── flux_par.expanded.yaml
    │   ├── flux_par.orig.yaml
    │   ├── flux_par.partial.yaml
    │   └── samples.npy
    └── runs
        ├── 03
        │   ├── MERLIN_FINISHED
        │   ├── flux_run.out
        │   ├── runs.slurm.err
        │   ├── runs.slurm.out
        │   └── runs.slurm.sh
        └── 07
            ├── flux_run.out
            ├── runs.slurm.err
            ├── runs.slurm.out
            └── runs.slurm.sh

6 directories, 22 files
```

And look at one of the runs to see flux output:

```bash
$ cat studies/flux_par_20230322-213613/runs/03/flux_run.out
```
```console
Hello world from processor flux-sample-0, rank 0 out of 1 processors
num args = 3
args = /workflow/studies/flux_par_20230322-213613/build/mpi_hello 0.7961944941912834 0.5904591175676233 
```

I'm not sure this is entirely correct (I was expecting more runs) but I'm fairly happy for
a first shot. To clean up:

```bash
$ kubectl delete -f ./examples/launchers/merlin/minicluster.yaml
$ kubectl delete -f ./examples/launchers/merlin/services.yaml
```

### Development Notes

I did some digging into the logic, and found that the underlying submission was a flux submit -> flux exec
to start a celery worker:

```bash
$ flux mini alloc -N 2 --exclusive --job-name=merlin flux exec `which /bin/bash` -c "celery -A merlin worker -l INFO --concurrency 1 --prefetch-multiplier 1 -Ofair -Q \'[merlin]_flux_par\'"
```
I think this should be changed to:

```bash
$ flux alloc -N 2 --exclusive --job-name=merlin /bin/bash -c "celery -A merlin worker -l INFO --concurrency 1 --prefetch-multiplier 1 -Ofair -Q \'[merlin]_flux_par\'"
```

 - I don't think we need flux exec
 - Why would there be more than one /bin/bash?

I don't fully understand the relationship between the celery queue and Flux - I think Flux should be used to submit jobs directly to,
as opposed to just using it to start a celery working. It also seems like there is one too many layers of complexity. If we have a Flux queue
why do we also need a celery queue?
