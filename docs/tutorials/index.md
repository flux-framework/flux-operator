# Tutorials

This quick set of tutorials shows how to accomplish specific tasks with the Flux Operator.
For these tutorials, we assume using a generic cluster. For cloud-specific tutorials, see the
[Deployment](https://flux-framework.org/flux-operator/deployment/index.html) tutorials.
If you have any questions or issues, please [let us know](https://github.com/flux-framework/flux-operator/issues)

## Isolated Tutorials

The following tutorials are provided from their respective directories (and are not documented here):

### Simulations

 - [Laghos](https://github.com/flux-framework/flux-operator/blob/main/examples/simulations/laghos-demos/minicluster.yaml)
 - [Lulesh](https://github.com/flux-framework/flux-operator/tree/main/examples/simulations/lulesh/minicluster.yaml)
 - [Qmcpack](https://github.com/flux-framework/flux-operator/tree/main/examples/simulations/qmcpack/minicluster.yaml)
 - [Exaworks Ball Bounce](https://github.com/flux-framework/flux-operator/tree/main/examples/simulations/exaworks-ball-bounce/minicluster.yaml)

### Launchers

 - [Parsl](https://github.com/flux-framework/flux-operator/tree/main/examples/launchers/parsl)

### Experimental

#### Nested

 - [K3s](https://github.com/flux-framework/flux-operator/tree/main/examples/nested/k3s/basic): instiatiate k3s inside Flux, and deploy an app.


### Machine Learning

 - [Fireworks](https://github.com/flux-framework/flux-operator/blob/main/examples/machine-learning/fireworks)
 - [Pytorch MNIST](https://github.com/flux-framework/flux-operator/blob/main/examples/machine-learning/pytorch)
 - [Tensorflow cifar-10](https://github.com/flux-framework/flux-operator/blob/main/examples/machine-learning/tensorflow)
 - [Dask with Scikit-Learn](https://github.com/flux-framework/flux-operator/blob/main/examples/machine-learning/dask/scikit-learn)
 - [Ray with Scikit-Learn](https://github.com/flux-framework/flux-operator/blob/main/examples/machine-learning/ray/scikit-learn)

### Services

 - [Merlin Basic](https://github.com/flux-framework/flux-operator/blob/main/examples/launchers/merlin/basic)
 - [Merlin Singularity Openfoam](https://github.com/flux-framework/flux-operator/blob/main/examples/launchers/merlin/singularity-openfoam)

### Workflows

Although some of the others above are also workflows, these examples are going to use `flux tree` (in various contexts) to
submit different job hierarchies and get around the etcd bottleneck in Kubernetes. 

 - [Basic Tree](https://github.com/flux-framework/flux-operator/blob/main/examples/workflows/tree)
 - [Instance Variables](https://github.com/flux-framework/flux-operator/blob/main/examples/workflows/tree-with-variables)

We have just started this arm of our experiments and you can expect more as we go!

## Integrated Tutorials

The following tutorials are included in the rendered documentation here.

```{toctree}
:maxdepth: 2
jobs
singularity
multi-tenancy
interactive
services
scaling
staging
volumes
state
```
