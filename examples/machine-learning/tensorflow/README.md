# Tensorflow with Flux

This is an experiment to distribute tensorflow jobs on a Flux Operator minicluster. We are
basically going to try to match a distributed Tensorflow setup to the minicluster size, etc.

> Tensorflow, Tensornooooo!

I know, I know. But it's a fairly popular library. I promise it's worth a try!

## Usage

First, let's create a kind cluster. From the context of this directory:

```bash
$ kind create cluster --config ../../kind-config.yaml
```

And then install the operator, create the namespace, and apply the MiniCluster YAML here.

```bash
$ kubectl apply -f ../../dist/flux-operator.yaml
$ kubectl create namespace flux-operator
$ kubectl apply ./minicluster.yaml
```

## Example

A complete usage example using the CIFAR-10 dataset is included in the examples directory,
and this is what the minicluster is going to run! We basically launch main.py with flux,
and we create a distributed Server via Flux, so a flux node running the script will take
on a worker or "main" node task depending on the hostname. Note that the example original epochs
were 1000, and I reduced to 10 to make it feasible to run. If you want to create your own
distributed training using this:

```python
import tensorflow as tf
from tensorflow_flux import tf_config_from_flux

# Adjust your minicluster parameters here!
cluster, my_job_name, my_task_index = tf_config_from_flux(
    ps_number=1, cluster_size=4, job_name="flux-sample", port_number=2222
)

cluster_spec = tf.train.ClusterSpec(cluster)
server = tf.distribute.Server(
    server_or_cluster_def=cluster_spec,
    job_name=my_job_name,
    task_index=my_task_index
)

# Your code here
```