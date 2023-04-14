# Pytorch with Flux

This will test running [pytorch](https://github.com/google/learn-oss-with-google/tree/main/kubernetes/job-examples/ml_training_pytorch) training with Flux, 
which we will do by way of pulling a Singularity container. It's a very simple (and a bit hacky setup) but it was fairly straight forward to do.

## Usage

First, let's create a kind cluster. From the context of this directory:

```bash
$ kind create cluster --config ../../kind-config.yaml
```

And then install the operator, create the namespace, and apply the MiniCluster YAML here.

```bash
$ kubectl apply -f ../../dist/flux-operator.yaml
$ kubectl create namespace flux-operator
$ kubectl apply -f ./minicluster.yaml
```

This will pull the singularity container to your bound directory (in the pods as `/tmp/workflow`) and then run
Pytorch. Pulling the container is the slowest step (took 10+ minutes), and if you run this many times you'd only need to do it once given
a persistent volume. You can see all logs from the broker with this command:

```bash
$ kubectl logs -n flux-operator flux-sample-0-7tx7s -f
```

And first you'll see the broker start and wait for the workers, and then each converting the SIF to a sandbox. 
Arguably we could (and should) do this first so each of the worker nodes doesn't have to redundantly do it.

```bash
broker.info[0]: online: flux-sample-[0-3] (ranks 0-3)
broker.info[0]: quorum-full: quorum->run 27.6138s
INFO:    Converting SIF file to temporary sandbox...
INFO:    Converting SIF file to temporary sandbox...
```

I tried doing a build into a sandbox, but ran into [this issue](https://github.com/apptainer/singularity/issues/2352).
This would important to fix if running this beyond this demo. Once the sandboxes are converted, you
should see pytorch training run, and it will take a while to get through all the batch and epochs, but then finish!

```console
...
2023-04-14 02:54:00,885 [INFO] root:121: Rank 0: epoch 8/8, average epoch loss=0.0881
2023-04-14 02:54:00,885 [INFO] root:121: Rank 1: epoch 8/8, average epoch loss=0.0881
2023-04-14 02:54:00,885 [INFO] root:122: Rank 0: training completed.
2023-04-14 02:54:00,885 [INFO] root:122: Rank 1: training completed.
2023-04-14 02:54:00,890 [INFO] root:121: Rank 2: epoch 8/8, average epoch loss=0.0881
2023-04-14 02:54:00,891 [INFO] root:121: Rank 3: epoch 8/8, average epoch loss=0.0881
2023-04-14 02:54:00,891 [INFO] root:122: Rank 2: training completed.
2023-04-14 02:54:00,896 [INFO] root:122: Rank 3: training completed.
INFO:    Cleaning up image...
INFO:    Cleaning up image...
broker.info[0]: rc2.0: flux submit -N 2 -n 2 --quiet --watch singularity exec ./pytorch.sif /bin/bash ./launch.sh flux-sample 8080 2 Exited (rc=0) 1196.1s
broker.info[0]: rc2-success: run->cleanup 19.9356m
broker.info[0]: cleanup.0: flux queue stop --quiet --all --nocheckpoint Exited (rc=0) 0.1s
broker.info[0]: cleanup.1: flux job cancelall --user=all --quiet -f --states RUN Exited (rc=0) 0.0s
broker.info[0]: cleanup.2: flux queue idle --quiet Exited (rc=0) 0.1s
broker.info[0]: cleanup-success: cleanup->shutdown 0.247919s
broker.info[0]: children-complete: shutdown->finalize 68.6433ms
broker.info[0]: rc3.0: running /etc/flux/rc3.d/01-sched-fluxion
broker.info[0]: online: flux-sample-0 (ranks 0)
broker.info[0]: rc3.0: /etc/flux/rc3 Exited (rc=0) 0.2s
broker.info[0]: rc3-success: finalize->goodbye 0.197565s
broker.info[0]: goodbye: goodbye->exit 0.063946ms
```

And that's it! You'll see data populated into your local volume:

```bash
data_0
└── MNIST
    └── raw
        ├── t10k-images-idx3-ubyte
        ├── t10k-images-idx3-ubyte.gz
        ├── t10k-labels-idx1-ubyte
        ├── t10k-labels-idx1-ubyte.gz
        ├── train-images-idx3-ubyte
        ├── train-images-idx3-ubyte.gz
        ├── train-labels-idx1-ubyte
        └── train-labels-idx1-ubyte.gz
data_1
└── MNIST
    └── raw
        ├── t10k-images-idx3-ubyte
        ├── t10k-images-idx3-ubyte.gz
        ├── t10k-labels-idx1-ubyte
        ├── t10k-labels-idx1-ubyte.gz
        ├── train-images-idx3-ubyte
        ├── train-images-idx3-ubyte.gz
        ├── train-labels-idx1-ubyte
        └── train-labels-idx1-ubyte.gz
data_2
└── MNIST
    └── raw
        ├── t10k-images-idx3-ubyte
        ├── t10k-images-idx3-ubyte.gz
        ├── t10k-labels-idx1-ubyte
        ├── t10k-labels-idx1-ubyte.gz
        ├── train-images-idx3-ubyte
        ├── train-images-idx3-ubyte.gz
        ├── train-labels-idx1-ubyte
        └── train-labels-idx1-ubyte.gz
data_3
└── MNIST
    └── raw
        ├── t10k-images-idx3-ubyte
        ├── t10k-images-idx3-ubyte.gz
        ├── t10k-labels-idx1-ubyte
        ├── t10k-labels-idx1-ubyte.gz
        ├── train-images-idx3-ubyte
        ├── train-images-idx3-ubyte.gz
        ├── train-labels-idx1-ubyte
        └── train-labels-idx1-ubyte.gz
```

If you want to debug something, you can set interactive: true to run in interactive mode, and then shell into the pod, connect to the broker:

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-jlsp6 bash
$ sudo -u fluxuser -E $(env) -E HOME=/home/fluxuser flux proxy local:///run/flux/local bash
```

And when you are done, clean up:

```bash
$ kubectl delete -f minicluster.yaml
```