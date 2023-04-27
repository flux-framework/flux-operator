# Merlin Demo Workflow


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
$ kubectl apply -f examples/launchers/merlin/basic/minicluster.yaml
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
and run the workers! 

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
$ kubectl delete -f ./examples/launchers/merlin/basic/minicluster.yaml
$ kubectl delete -f ./examples/launchers/merlin/services.yaml
```