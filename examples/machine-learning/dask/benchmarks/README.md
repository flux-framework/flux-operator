# Dask Benchmarks

I want to test running some simple [benchmarks](https://matthewrocklin.com/blog/work/2017/07/03/scaling)
for Dask with Flux. This example was run on a cloud, but for our purposes we will apply the minicluster 4 times with the following
parameters (manually done for now):

```yaml
size: 5
command: python3 /tmp/workflow/launch.py --workers 5 --cores 2 --iter 3

size: 4
command: python3 /tmp/workflow/launch.py --workers 4 --cores 2 --iter 3

size: 3
command: python3 /tmp/workflow/launch.py --workers 3 --cores 2 --iter 3

size: 2
command: python3 /tmp/workflow/launch.py --workers 2 --cores 2 --iter 3

size: 1
command: python3 /tmp/workflow/launch.py --workers 1 --cores 2 --iter 3
```
(note for some sizes I happened to do more iterations, but then decided 3 was enough for a demo)
Since I wanted to be sure the pods were all ready with `flux resource list` I decided to do
this interactively for now (interactive: true in the minicluster.yaml)

## Usage

### Create Cluster

First, let's create a kind cluster.

```bash
$ kind create cluster --config ../../../kind-config.yaml
```

And then install the operator, create the namespace, and apply the MiniCluster YAML here.

```bash
$ kubectl apply -f ../../../dist/flux-operator.yaml
$ kubectl create namespace flux-operator
```

### Run Experiments

You'll want to tweak the minicluster.yaml size and command parameters for each run of the above.
Note that the output file saved corresponds to this size.

```bash
$ kubectl apply -f ./minicluster.yaml
```

For each run, running the above will install dependencies (dask and pandas) directly into the base image, and then mount
the current directory (in the pods as `/tmp/workflow`). We run the [launch.py](launch.py)
script from the broker, and you can inspect this script to see how we create and connect
to Flux. You can watch logs doing the following: 

```bash
$ kubectl logs -n flux-operator flux-sample-0-7tx7s -f
```

You'll want to shell into the broker node, and connect to the socket:

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-jlsp6 bash
$ sudo -u fluxuser -E $(env) -E HOME=/home/fluxuser flux proxy local:///run/flux/local bash
```

And then run each of the commands above, e.g, (depending on your size:)

```bash
python3 /tmp/workflow/launch.py --workers 5 --cores 2 --iter 3  # size 5
python3 /tmp/workflow/launch.py --workers 4 --cores 2 --iter 3  # size 4
python3 /tmp/workflow/launch.py --workers 3 --cores 2 --iter 3  # size 3
python3 /tmp/workflow/launch.py --workers 2 --cores 2 --iter 3  # size 2
python3 /tmp/workflow/launch.py --workers 1 --cores 2 --iter 3  # size 1
```

The experiments will run (across Dask worker sizes) and save output to a data frame:

```bash
dask-experiments-1-raw.csv
dask-experiments-2-raw.csv
dask-experiments-3-raw.csv
dask-experiments-4-raw.csv
dask-experiments-5-raw.csv
```

You'll probably want to run this in an actual large / scaled environment, my local
runs weren't very fruitful!

### Cleanup

When you are done, clean up:

```bash
$ kubectl delete -f minicluster.yaml
```

Make sure to clean up your shared tmpdir!

```bash
$ rm *.out
$ sudo rm -rf ./tmp/*
```
