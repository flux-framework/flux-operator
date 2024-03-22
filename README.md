# Flux Operator

![docs/development/the-operator.jpg](docs/development/the-operator.jpg)
[![DOI](https://zenodo.org/badge/528650707.svg)](https://zenodo.org/badge/latestdoi/528650707)

The Flux Operator is a Kubernetes Cluster [Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
that you can install to your cluster to create and control a [Flux Framework](https://flux-framework.org/) "MiniCluster"
to launch jobs to.

Read more, including user and developer guides, and project background in our 💛 [Documentation](https://flux-framework.org/flux-operator) 💛

## Presentations

 - [Kubecon 2023](https://t.co/vjRydPx1rb)
 - [HPC Knowledge Meeting 2023](https://hpckp.org/talks/cloud-and-hpc-convergence-flux-for-job-management-on-kubernetes/)

## Organization

The basic idea is that we present the idea of a **MiniCluster** that is a custom resource definition (CRD)
that defines a job container (that does not need to have Flux) that (when submit) will create a set of config maps,
secrets, and the final Indexed Job that has the pod containers running with Flux. Since
this is a batchv1.Job, it will have states that we can track.

And you can find the following here:

 - [Flux Controllers](controllers/flux) are under `controllers/flux` for the `MiniCluster`
 - [API Spec](api/v1alpha1/) are under `api/v1alpha2/` also for `MiniCluster`
 - [Packages](pkg) include supporting packages for job conditions (state), if we eventually want that.
 - [Config](config) includes mostly automatically generated yaml configuration files needed by Kubernetes

And the following external resources might be useful:

 - [Flux Framework](https://flux-framework.org)
 - [Flux RESTful API](https://github.com/flux-framework/flux-restful-api): provides the interface for submitting jobs, if no command provided to the operator.
 - [Python SDK](sdk/python): for deploying MiniClusters and port forwarding.
 - [Flux HPC Examples](https://github.com/rse-ops/flux-hpc) containers and CRD for the operator to run Flux with HPC workloads (under development)
 - [Flux Cloud](https://github.com/converged-computing/flux-cloud): automation of experiments using the Flux Operator

**Note** we welcome contributions to code or to suggest features or identify bugs!

## Citation

You can follow the CITATION.cff (right sidebar in GitHub) to cite, or [view the paper directly here](https://doi.org/10.12688/f1000research.147989.1)
A direct (copy paste) citation is the following:

> Sochat V, Culquicondor A, Ojea A and Milroy D. The Flux Operator (version 1). F1000Research 2024, 13:203 (https://doi.org/10.12688/f1000research.147989.1)


## License

HPCIC DevTools is distributed under the terms of the MIT license.
All new contributions must be made under this license.

See [LICENSE](https://github.com/converged-computing/cloud-select/blob/main/LICENSE),
[COPYRIGHT](https://github.com/converged-computing/cloud-select/blob/main/COPYRIGHT), and
[NOTICE](https://github.com/converged-computing/cloud-select/blob/main/NOTICE) for details.

SPDX-License-Identifier: (MIT)

LLNL-CODE- 842614
