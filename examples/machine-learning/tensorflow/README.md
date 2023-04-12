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

You can then inspect logs, and see the training happening! Note that we adjusted 1000 epochs (takes a long time)
to 2 (is very quick)!

```bash
$ kubectl logs -n flux-operator flux-sample-0-7tx7s -f
```
```console
208 2 1.3872336 0.515625
209 0 1.3110054 0.5625
209 1 1.2829641 0.5234375
209 2 1.496664 0.4296875
210 1 1.3997178 0.4375
210 0 1.4337707 0.484375
210 2 1.2814229 0.5703125
211 0 1.2377738 0.5078125
211 1 1.497496 0.515625
211 2 1.2085061 0.5625
212 1 1.4642732 0.484375
212 2 1.2024314 0.5859375
212 0 1.6061184 0.453125
...
```

In the above, you are looking at a row of `step task_number loss accuracy` where the task number corresponds
to different flux cluster nodes. 2 epochs likely isn't enough to get a good result, but this is just an example.
It's intended to get you started. If you have a tensorflow example to share,
we hope you [let us know](https://github.com/flux-framework/flux-operator/issues).

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