# Tutorials

This quick set of tutorials shows how to accomplish specific tasks with the Flux Operator.
For these tutorials, we assume using a generic cluster. For cloud-specific tutorials, see the
[Deployment](https://flux-framework.org/flux-operator/deployment/index.html) tutorials.
If you have any questions or issues, please [let us know](https://github.com/flux-framework/flux-operator/issues)

## Isolated Tutorials

The following tutorials are provided from their respective directories (and are not documented here):

### Simulations

 - [Laghos](https://github.com/flux-framework/flux-operator/blob/main/examples/simulations/laghos-demos/minicluster.yaml)
 - [Exaworks Ball Bounce](https://github.com/flux-framework/flux-operator/tree/main/examples/simulations/exaworks-ball-bounce/minicluster.yaml)

### Launchers

 - [Parsl](https://github.com/flux-framework/flux-operator/tree/main/examples/launchers/parsl)

### Experimental

#### Bursting

 - [Bursting to GKE](https://github.com/flux-framework/flux-operator/tree/main/examples/experimental/bursting/broker-gke) from a local broker to an external Google Kubernetes Engine cluster.
 - [Bursting to EKS](https://github.com/flux-framework/flux-operator/tree/main/examples/experimental/bursting/broker-eks) from a local broker to an external Amazon Elastic Kubernetes Service
 - [Bursting to Compute Engine](https://github.com/flux-framework/flux-operator/tree/main/examples/experimental/bursting/broker-compute-engine) from a GKE broker to an external Compute Engine cluster.
 - [Bursting (nginx service)](https://github.com/flux-framework/flux-operator/tree/main/examples/experimental/bursting/nginx) design to use central router for bursting.

#### Process Namespace

 - [multiple-pods-per-node](https://github.com/flux-framework/flux-operator/tree/main/examples/experimental/multiple-pods-per-node): Allow multiple pods to be scheduled per node (controlled by cgroups)
 - [shared-process-space](https://github.com/flux-framework/flux-operator/tree/main/examples/experimental/shared-process-space): Allow flux to execute a command into another container

### Machine Learning

 - [Fireworks](https://github.com/flux-framework/flux-operator/blob/main/examples/machine-learning/fireworks)
 - [Pytorch MNIST](https://github.com/flux-framework/flux-operator/blob/main/examples/machine-learning/pytorch)
 - [Tensorflow cifar-10](https://github.com/flux-framework/flux-operator/blob/main/examples/machine-learning/tensorflow)
 - [Ray with Scikit-Learn](https://github.com/flux-framework/flux-operator/blob/main/examples/machine-learning/ray/scikit-learn)

### Message Passing Interface (MPI)

 - [mpich](https://github.com/flux-framework/flux-operator/blob/main/examples/mpi/mpich)

### Queue Interaction

These examples show how to interact with your flux queue from a sidecar container (that has access to the flux broker of the pod):

 - [flux-sidecar](https://github.com/flux-framework/flux-operator/blob/main/examples/tests/flux-sidecar) to see a sleep job in the main application queue

### Services

 - [Flux Metrics API](https://github.com/flux-framework/flux-operator/blob/main/examples/experimental/metrics-api): run a custom metrics API directly from the lead broker to help with autoscaling
 - [Nginx](https://github.com/flux-framework/flux-operator/blob/main/examples/services/sidecar/nginx): to run alongisde your MiniCluster (and possibly expose functionality)
 - [Flux Restful](https://github.com/flux-framework/flux-operator/blob/main/examples/interactive/flux-restful): to run a restful API server alongside your cluster.

### Workflows

 - [ramble](https://github.com/flux-framework/flux-operator/blob/main/examples/workflows/ramble): recommended if you require installation with spack.

We have just started this arm of our experiments and you can expect more as we go!

## Integrated Tutorials

The following tutorials are included in the rendered documentation here.

```{toctree}
:maxdepth: 2
jobs
interactive
services
scaling
elasticity
volumes
state
```
