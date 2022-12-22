# Flux Operator

![docs/development/the-operator.jpg](docs/development/the-operator.jpg)

The Flux Operator is a Kubernetes Cluster [Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) 
that you can install to your cluster to create and control [Flux Framework](https://flux-framework.org/) "Mini Clusters"
to launch jobs to.

Read more, including user and developer guides, and project background in our 💛 [Documentation](https://flux-framework.org/flux-operator) 💛

🚧️ Under Construction! 🚧️

## Organization

The basic idea is that we present the idea of a **MiniCluster** that is a custom resource definition (CRD)
that defines a job container (that must have Flux) that (when submit) will create a set of config maps,
secrets (e.g., tls), and the final Batch job that has the pod containers running with flux. Since
this is a batchv1.Job, it will have states that we can track.

And you can find the following here:

 - [Flux Controllers](controllers/flux) are under `controllers/flux` for the `MiniCluster`
 - [API Spec](api/v1alpha1/) are under `api/v1alpha1/` also for `MiniCluster`
 - [Packages](pkg) include supporting packages for job conditions (state), if we eventually want that.
 - [Config](config) includes mostly automatically generated yaml configuration files needed by Kubernetes
 - [TODO.md](TODO.md) is a set of things to be worked on, if you'd like to contribute!

And the following external resources might be useful:

 - [Flux HPC Examples](https://github.com/rse-ops/flux-hpc) containers and CRD for the operator to run Flux with HPC workloads (under development)

**Note** this project is actively under development, and you can expect change and improvements!
We apologize for bugs you run into, and hope you tell us soon so we can work on resolving them.

#### License

This work is licensed under the [Apache-2.0](https://github.com/kubernetes-sigs/kueue/blob/ec9b75eaadb5c78dab919d8ea6055d33b2eb09a2/LICENSE) license.

SPDX-License-Identifier: Apache-2.0

LLNL-CODE-764420