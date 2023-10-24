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

#### Bursting

 - [Bursting to GKE](https://github.com/flux-framework/flux-operator/tree/main/examples/experimental/bursting/broker-gke) from a local broker to an external Google Kubernetes Engine cluster.
 - [Bursting to EKS](https://github.com/flux-framework/flux-operator/tree/main/examples/experimental/bursting/broker-eks) from a local broker to an external Amazon Elastic Kubernetes Service
 - [Bursting to Compute Engine](https://github.com/flux-framework/flux-operator/tree/main/examples/experimental/bursting/broker-compute-engine) from a GKE broker to an external Compute Engine cluster.
 - [Bursting (nginx service)](https://github.com/flux-framework/flux-operator/tree/main/examples/experimental/bursting/nginx) design to use central router for bursting.

#### Nested

 - [K3s](https://github.com/flux-framework/flux-operator/tree/main/examples/nested/k3s/basic): instiatiate k3s inside Flux, and deploy an app.

#### Process Namespace

 - [shared-process-space](https://github.com/flux-framework/flux-operator/tree/main/examples/experimental/shared-process-space): Allow flux to execute a command into another container

### Machine Learning

 - [Fireworks](https://github.com/flux-framework/flux-operator/blob/main/examples/machine-learning/fireworks)
 - [Pytorch MNIST](https://github.com/flux-framework/flux-operator/blob/main/examples/machine-learning/pytorch)
 - [Tensorflow cifar-10](https://github.com/flux-framework/flux-operator/blob/main/examples/machine-learning/tensorflow)
 - [Ray with Scikit-Learn](https://github.com/flux-framework/flux-operator/blob/main/examples/machine-learning/ray/scikit-learn)

### Message Passing Interface (MPI)

 - [openmpi](https://github.com/flux-framework/flux-operator/blob/main/examples/mpi/ompi)
 - [mpich](https://github.com/flux-framework/flux-operator/blob/main/examples/mpi/mpich)


### Queue Interaction

These examples show how to interact with your flux queue from a sidecar container (that has access to the flux broker of the pod):

 - [flux-sidecar](https://github.com/flux-framework/flux-operator/blob/main/examples/tests/flux-sidecar) to see a sleep job in the main application queue

### Services

 - [Nginx](https://github.com/flux-framework/flux-operator/blob/main/examples/services/sidecar/nginx): to run alongisde your MiniCluster (and possibly expose functionality)

### Workflows

 - [ramble](https://github.com/flux-framework/flux-operator/blob/main/examples/workflows/ramble): recommended if you require installation with spack.

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
interactive
services
scaling
elasticity
staging
volumes
state
```
