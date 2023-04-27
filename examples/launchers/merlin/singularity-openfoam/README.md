### Merlin OpenFoam

This example will run openfoam via a Singularity container.

Note that we derive this example from [this repository](https://github.com/rse-ops/flux-hpc/tree/main/merlin-demos)
and use a [customized set of containers](https://github.com/rse-ops/flux-hpc/tree/main/merlin-demos-certs) that have
certificates built into the containers. This is obviously not recommended and only used here for an example.

#### Create MiniCluster

First, cd to the directory here, and create the kind cluster:

```bash
$ kind create cluster --config ../../../kind-config.yaml
```

And the Flux Operator namespace created:

```bash
$ kubectl create -n flux-operator
```

And then generate the (separate) pods to run redis and rabbitmq in the flux-operator namespace.
The containers already have shared certificates (just for this test case)!

```bash
$ kubectl create -f ../services.yaml
```

And create the MiniCluster to use them!

```bash
$ kubectl apply -f minicluster.yaml
```

The MiniCluster is created in interactive mode, and we do this because merlin isn't designed to
run and hang until it's done. If we don't run in interactive mode we will miss the execution
that is created via flux alloc. Look at pods running:

```bash
kubectl get -n flux-operator pods
```
```console
$ kubectl get -n flux-operator pods
NAME                         READY   STATUS      RESTARTS   AGE
flux-sample-0-q22pt          1/1     Running     0          20s
flux-sample-1-nmqsc          1/1     Running     0          20s
flux-sample-2-49gv6          1/1     Running     0          20s
flux-sample-3-tvvg9          1/1     Running     0          20s
rabbitmq-f8c84d986-2dff6     1/1     Running     0          32s
redis-c9469b9c5-hvkdr        1/1     Running     0          32
```

The broker is going to pull the container, and to the present working directory, so you
should see a *.sif show up.

#### Interactive Mode

We will be issuing commands in interactive mode! Shell into the broker pod:

```bash
$ kubectl exec -n flux-operator -it flux-sample-0-f9ts7 -- bash
```

You'll want to watch the main broker pod for when it's finished installing conda packages
(the socket won't be available until after). For an actual workflow, you'd want to install these
to the container apriori and not have to wait. Here is how to look at logs:

```bash
$ kubectl logs -n flux-operator flux-sample-0-f9ts7 -f
```

When the broker has started, in your other terminal, connect to it:

```bash
$ export PYTHONPATH=/home/fluxuser/.local/lib/python3.10/site-packages:/opt/conda/lib/python3.10
$ sudo -E LD_LIBRARY_PATH=$LD_LIBRARY_PATH -E PATH=$PATH -E HOME=/home/fluxuser -E PYTHONPATH=$PYTHONPATH -u fluxuser flux proxy local:///var/run/flux/local bash
```

At this point, we are the fluxuser! We can test if merlin is connected to its services via `merlin info`:

```bash
$ whoami
fluxuser
```
```bash
$ merlin info
```

**Note** that if you find `merlin info` returns an error, try deleting and re-creating the services.

```bash
$ kubectl delete -f ../services.yaml 
$ kubectl apply -f ../services.yaml 
```

<details>

<summary>Output of Merlin Info</summary>

```console
[2023-04-03 22:48:54: INFO] Reading app config from file /home/fluxuser/.merlin/app.yaml
  
                                                
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

I found I couldn't get this to work within the workflow, so I staged it beforhand:

```bash
$ singularity exec -B $PWD/:/workflow openfoam6.sif cp -rf /opt/openfoam6/tutorials/incompressible/icoFoam/cavity/cavity /workflow/cavity
```
Now let's queue our tasks and run the workers! 

```bash
$ merlin run openfoam_wf.yaml 
$ merlin run-workers openfoam_wf.yaml
```

You should be able to inspect the bound output from the host, and see the result in one simulation:

```bash
$ cat openfoam_wf_output/openfoam_wf_singularity_20230404-042704/sim_runs/00/00/sim_runs.slurm.out 
```
(Note I'm not sure why the log prefix is slurm)

<details>

<summary>Example output</summary>

```console
/workflow/openfoam_wf_output/openfoam_wf_singularity_20230404-042704/sim_runs/00/00
cavity
sim_runs.slurm.sh
/workflow/openfoam_wf_output/openfoam_wf_singularity_20230404-042704/sim_runs/00/00/cavity
0
constant
system
singularity exec --bind  /workflow/openfoam_wf_output/openfoam_wf_singularity_20230404-042704/sim_runs/00/00:/merlin_sample,/workflow/openfoam_wf_output/openfoam_wf_singularity_20230404-042704/sim_runs/00/00/cavity:/cavity /workflow/openfoam6.sif /merlin_sample/run_openfoam 16.696213872840325
***** Setting up control parameters ***** 
Running icoFoam
/*---------------------------------------------------------------------------*\
  =========                 |
  \\      /  F ield         | OpenFOAM: The Open Source CFD Toolbox
   \\    /   O peration     | Website:  https://openfoam.org
    \\  /    A nd           | Version:  6
     \\/     M anipulation  |
\*---------------------------------------------------------------------------*/
Build  : 6-e29811f5dff8
Exec   : postProcess -func enstrophy
Date   : Apr 04 2023
Time   : 04:27:41
Host   : "flux-sample-2"
PID    : 3594
I/O    : uncollated
Case   : //cavity
nProcs : 1
sigFpe : Enabling floating point exception trapping (FOAM_SIGFPE).
fileModificationChecking : Monitoring run-time modified files using timeStampMaster (fileModificationSkew 10)
allowSystemOperations : Allowing user-supplied system call operations

// * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * //
Create time

Create mesh for time = 0

Time = 0

Reading fields:
    volVectorFields: U

Executing functionObjects
    functionObjects::enstrophy enstrophy writing field: enstrophy

Time = 0.100445

Reading fields:
    volVectorFields: U

Executing functionObjects
    functionObjects::enstrophy enstrophy writing field: enstrophy

Time = 0.199551

Reading fields:
    volVectorFields: U

Executing functionObjects
    functionObjects::enstrophy enstrophy writing field: enstrophy

Time = 0.299996

Reading fields:
    volVectorFields: U

Executing functionObjects
    functionObjects::enstrophy enstrophy writing field: enstrophy

Time = 0.400441

Reading fields:
    volVectorFields: U

Executing functionObjects
    functionObjects::enstrophy enstrophy writing field: enstrophy

Time = 0.499546

Reading fields:
    volVectorFields: U

Executing functionObjects
    functionObjects::enstrophy enstrophy writing field: enstrophy

Time = 0.599991

Reading fields:
    volVectorFields: U

Executing functionObjects
    functionObjects::enstrophy enstrophy writing field: enstrophy

Time = 0.700436

Reading fields:
    volVectorFields: U

Executing functionObjects
    functionObjects::enstrophy enstrophy writing field: enstrophy

Time = 0.799542

Reading fields:
    volVectorFields: U

Executing functionObjects
    functionObjects::enstrophy enstrophy writing field: enstrophy

Time = 0.899987

Reading fields:
    volVectorFields: U

Executing functionObjects
    functionObjects::enstrophy enstrophy writing field: enstrophy

Time = 1.000432

Reading fields:
    volVectorFields: U

Executing functionObjects
    functionObjects::enstrophy enstrophy writing field: enstrophy

End

/*---------------------------------------------------------------------------*\
  =========                 |
  \\      /  F ield         | OpenFOAM: The Open Source CFD Toolbox
   \\    /   O peration     | Website:  https://openfoam.org
    \\  /    A nd           | Version:  6
     \\/     M anipulation  |
\*---------------------------------------------------------------------------*/
Build  : 6-e29811f5dff8
Exec   : foamToVTK
Date   : Apr 04 2023
Time   : 04:27:41
Host   : "flux-sample-2"
PID    : 3595
I/O    : uncollated
Case   : //cavity
nProcs : 1
sigFpe : Enabling floating point exception trapping (FOAM_SIGFPE).
fileModificationChecking : Monitoring run-time modified files using timeStampMaster (fileModificationSkew 10)
allowSystemOperations : Allowing user-supplied system call operations

// * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * //
Create time

Create mesh for time = 0

Time: 0
    volScalarFields            : p enstrophy
    volVectorFields            : U

    Internal  : "//cavity/VTK/cavity_0.vtk"
    Original cells:400 points:882   Additional cells:0  additional points:0

    Patch     : "//cavity/VTK/movingWall/movingWall_0.vtk"
    Patch     : "//cavity/VTK/fixedWalls/fixedWalls_0.vtk"
    Patch     : "//cavity/VTK/frontAndBack/frontAndBack_0.vtk"
Time: 0.100445
    volScalarFields            : p enstrophy
    volVectorFields            : U

    Internal  : "//cavity/VTK/cavity_75.vtk"
    Patch     : "//cavity/VTK/movingWall/movingWall_75.vtk"
    Patch     : "//cavity/VTK/fixedWalls/fixedWalls_75.vtk"
    Patch     : "//cavity/VTK/frontAndBack/frontAndBack_75.vtk"
    surfScalarFields  : phi
Time: 0.199551
    volScalarFields            : p enstrophy
    volVectorFields            : U

    Internal  : "//cavity/VTK/cavity_149.vtk"
    Patch     : "//cavity/VTK/movingWall/movingWall_149.vtk"
    Patch     : "//cavity/VTK/fixedWalls/fixedWalls_149.vtk"
    Patch     : "//cavity/VTK/frontAndBack/frontAndBack_149.vtk"
    surfScalarFields  : phi
Time: 0.299996
    volScalarFields            : p enstrophy
    volVectorFields            : U

    Internal  : "//cavity/VTK/cavity_224.vtk"
    Patch     : "//cavity/VTK/movingWall/movingWall_224.vtk"
    Patch     : "//cavity/VTK/fixedWalls/fixedWalls_224.vtk"
    Patch     : "//cavity/VTK/frontAndBack/frontAndBack_224.vtk"
    surfScalarFields  : phi
Time: 0.400441
    volScalarFields            : p enstrophy
    volVectorFields            : U

    Internal  : "//cavity/VTK/cavity_299.vtk"
    Patch     : "//cavity/VTK/movingWall/movingWall_299.vtk"
    Patch     : "//cavity/VTK/fixedWalls/fixedWalls_299.vtk"
    Patch     : "//cavity/VTK/frontAndBack/frontAndBack_299.vtk"
    surfScalarFields  : phi
Time: 0.499546
    volScalarFields            : p enstrophy
    volVectorFields            : U

    Internal  : "//cavity/VTK/cavity_373.vtk"
    Patch     : "//cavity/VTK/movingWall/movingWall_373.vtk"
    Patch     : "//cavity/VTK/fixedWalls/fixedWalls_373.vtk"
    Patch     : "//cavity/VTK/frontAndBack/frontAndBack_373.vtk"
    surfScalarFields  : phi
Time: 0.599991
    volScalarFields            : p enstrophy
    volVectorFields            : U

    Internal  : "//cavity/VTK/cavity_448.vtk"
    Patch     : "//cavity/VTK/movingWall/movingWall_448.vtk"
    Patch     : "//cavity/VTK/fixedWalls/fixedWalls_448.vtk"
    Patch     : "//cavity/VTK/frontAndBack/frontAndBack_448.vtk"
    surfScalarFields  : phi
Time: 0.700436
    volScalarFields            : p enstrophy
    volVectorFields            : U

    Internal  : "//cavity/VTK/cavity_523.vtk"
    Patch     : "//cavity/VTK/movingWall/movingWall_523.vtk"
    Patch     : "//cavity/VTK/fixedWalls/fixedWalls_523.vtk"
    Patch     : "//cavity/VTK/frontAndBack/frontAndBack_523.vtk"
    surfScalarFields  : phi
Time: 0.799542
    volScalarFields            : p enstrophy
    volVectorFields            : U

    Internal  : "//cavity/VTK/cavity_597.vtk"
    Patch     : "//cavity/VTK/movingWall/movingWall_597.vtk"
    Patch     : "//cavity/VTK/fixedWalls/fixedWalls_597.vtk"
    Patch     : "//cavity/VTK/frontAndBack/frontAndBack_597.vtk"
    surfScalarFields  : phi
Time: 0.899987
    volScalarFields            : p enstrophy
    volVectorFields            : U

    Internal  : "//cavity/VTK/cavity_672.vtk"
    Patch     : "//cavity/VTK/movingWall/movingWall_672.vtk"
    Patch     : "//cavity/VTK/fixedWalls/fixedWalls_672.vtk"
    Patch     : "//cavity/VTK/frontAndBack/frontAndBack_672.vtk"
    surfScalarFields  : phi
Time: 1.000432
    volScalarFields            : p enstrophy
    volVectorFields            : U

    Internal  : "//cavity/VTK/cavity_747.vtk"
    Patch     : "//cavity/VTK/movingWall/movingWall_747.vtk"
    Patch     : "//cavity/VTK/fixedWalls/fixedWalls_747.vtk"
    Patch     : "//cavity/VTK/frontAndBack/frontAndBack_747.vtk"
    surfScalarFields  : phi
End
```

Note that if you need to re-run (or try again) you need to purge the queues:

```bash
# Do this until you see no messages
$ merlin purge openfoam_wf.yaml  -f

# This can be run once after the above
$ flux job cancelall -f
$ flux job purge --force --num-limit=0
```

You'll see the tasks pick up in celery and a huge stream of output! You can also
look at the logs of redis and rabbitmq (the other containers) to see them receiving data.
We currently just run openfoam and generate data, and you can see both flux jobs running
to generate it:

```bash
$ flux jobs -a
```
```console
$ flux jobs -a
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
   ƒA86QZ2qm fluxuser merlin      R      1      1   45.89s flux-sample-2
   ƒA84PS1aX fluxuser merlin      R      1      1   45.97s flux-sample-3
   ƒ8oEEtiFh fluxuser merlin     CA      1      1   1.799m flux-sample-1
   ƒ8oGHVioH fluxuser merlin     CA      1      1   1.797m flux-sample-0
```

And some output (I installed tree 🌲️)!

```bash
$ find openfoam_wf_output/ -name MERLIN_FINISHED
```
```console
openfoam_wf_output/openfoam_wf_singularity_20230404-045711/sim_runs/00/00/MERLIN_FINISHED
openfoam_wf_output/openfoam_wf_singularity_20230404-045711/setup/MERLIN_FINISHED
```

Note there are additional steps in the workflow "combine-outputs" and "learn" 
but I found dimension errors in both. When you are done, exit and:

```bash
$ kubectl delete -f ./examples/launchers/merlin/singlarity-openfoam/minicluster.yaml
$ kubectl delete -f ./examples/launchers/merlin/services.yaml
```
