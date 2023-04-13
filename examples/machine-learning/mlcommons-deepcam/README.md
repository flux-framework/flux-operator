# Deepcam

> Deep Learning Climate Segmentation Benchmark

This shows a  PyTorch implementation for the climate segmentation benchmark, based on the
Exascale Deep Learning for Climate Analytics paper: https://arxiv.org/abs/1810.01993.
The workflow is provided from [mlcommons/deepcam](https://github.com/mlcommons/hpc/tree/main/deepcam).

## Create MiniCluster

First, cd to the directory here, and create the kind cluster:

```bash
$ kind create cluster --config ../../kind-config.yaml
```

And the Flux Operator namespace created:

```bash
$ kubectl create namespace flux-operator
```

And install the flux operator (from the repository here):

```bash
$ kubectl apply -f ../../dist/flux-operator.yaml
```

We don't want to create the minicluster quite yet! We want to prepare the data first.

## Dataset

You can read [more about the dataset here](https://github.com/mlcommons/hpc/tree/main/deepcam#dataset).
You'll need to download the dataset from [this globus endpoint](https://app.globus.org/file-manager?origin_id=0b226e2c-4de0-11ea-971a-021304b0cca7&origin_path=%2F) and into the current directory.
Note that I did this by setting up [Globus Connect Personal](https://www.globus.org/globus-connect-personal) and
then downloading to a scoped location on my computer, and then moving to the directory here.
First, extract the data (make sure you have ~50GB of space):

```bash
$ tar -xzvf deepcam-data-n512.tgz
$ chmod +x install_mini_dataset.sh
```

This will extract the data to a directory, `deepcam-data-n512` and then we can run the script to prepare it:

```bash
$ mkdir -p ./data
$ ./install_mini_dataset.sh ./deepcam-data-n512 ./data
```

This will basically copy the data over, and create the needed structure for training, etc.
It should look like this, with most of the files under "training":

```bash
$ ls ./data
stats.h5  train  validation
```

Note that the root directory here is bound to /tmp/workflow in our cluster, so it should
show up as `/tmp/workflow/data`.

## Training

Now that we have our data ready, we can create the minicluster (which will pull the container to run the job)
Note that we will use default parameters, but you can learn more about the defaults and parameters
[in the repository](https://github.com/mlcommons/hpc/tree/main/deepcam).

Then create the MiniCluster to use them! Let's hope your computer doesn't run out of space, or something like that.

```bash
$ kubectl apply -f minicluster.yaml
```

**WIP** need larger computer...
