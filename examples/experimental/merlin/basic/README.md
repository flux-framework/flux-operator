# Merlin Demo Workflow

This example will run a small "hello world" tutorial bringing up a rabbitmq and redis container from Flux.
This could (eventually) be part of a workflows tutorial, but for now can exemplify pod services alongside a MiniCluster.
Note that we derive this example from [this repository](https://github.com/rse-ops/flux-hpc/tree/main/merlin-demos)
and use a [customized set of containers](https://github.com/rse-ops/flux-hpc/tree/main/merlin-demos-certs) that have
certificates built into the containers. This is obviously not recommended and only used here for an example.
Create your cluster and install the operator.

```bash
kind create cluster
kubectl apply -f ../../../dist/flux-operator.yaml
```

And then generate the (separate) pods to run redis and rabbitmq in the flux-operator namespace.
The containers already have shared certificates (just for this test case)!

```bash
kubectl apply -f ../services.yaml
```

And create the MiniCluster to use them!

```bash
kubectl apply -f minicluster.yaml
```
The MiniCluster is created in interactive mode, and we do this because merlin isn't designed to
run and hang until it's done. If we don't run in interactive mode we will miss the execution
that is created via flux alloc. Look at pods running:

```bash
kubectl get pods
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
$ kubectl logs flux-sample-0-f9ts7 -f
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
$ kubectl exec -it flux-sample-0-f9ts7 -- bash
```

Connect to the flux broker:

```bash
source /mnt/flux/flux-view.sh
flux proxy $fluxsocket bash
```

At this point, we are the fluxuser! We can test if merlin is connected to its services via `merlin info`:

```bash
$ merlin info
```

<details>

<summary>Output of Merlin Info</summary>

```console
[2023-10-24 05:52:02: INFO] Reading app config from file /root/.merlin/app.yaml
  
                                                
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

 config_file        | /root/.merlin/app.yaml
 is_debug           | False
 merlin_home        | /root/.merlin
 merlin_home_exists | True
 broker server      | amqps://root:******@rabbitmq:5671//merlinu
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
/usr/bin/python3

 $ python3 --version
Python 3.10.12

 $ which pip3
/usr/bin/pip3

 $ pip3 --version
pip 22.0.2 from /usr/lib/python3/dist-packages/pip (python 3.10)

"echo $PYTHONPATH"
/opt/software/linux-rocky9-x86_64/gcc-11.3.1/flux-core-0.54.0-y6v7ctnpc4i3rd4tiv6d7qiqnqtqdzoy/lib/flux/python3.11:/mnt/flux/view/lib/python3.11/site-packages
```

</details>

Now let's queue our tasks (I guess we already did this, but I did it again)
and run the workers! 

```bash
export C_FORCE_ROOT=true
merlin run flux/flux_par.yaml
merlin run-workers flux/flux_par.yaml
```

Note that if you need to re-run (or try again) you need to purge the queues:

```bash
# Do this until you see no messages
merlin purge flux/flux_par.yaml -f

# This can be run once after the above
flux job purge --force --num-limit=0
```

Note that this is not currently working. I'm not sure we are pursuing Merlin (and I've spent 4 hours here already) so I am not proceeding.
