# Flux Operator

![docs/development/the-operator.jpg](docs/development/the-operator.jpg)
[![DOI](https://zenodo.org/badge/528650707.svg)](https://zenodo.org/badge/latestdoi/528650707)


The Flux Operator is a Kubernetes Cluster [Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
that you can install to your cluster to create and control [Flux Framework](https://flux-framework.org/) "Mini Clusters"
to launch jobs to.

Read more, including user and developer guides, and project background in our 💛 [Documentation](https://flux-framework.org/flux-operator) 💛

🚧️ Under Construction! 🚧️

**Important!** We recently removed a one-off container that ran before the MiniCluster creation to generate a certificate.
We have found [through testing](https://github.com/kubernetes-sigs/jobset/issues/104) that this somehow served as a warmup
for networking, and this means if you use the latest operator here, you may see slow times in creating the initial
broker setup. More details are available in [this post](https://github.com/converged-computing/operator-experiments/tree/main/google/service-timing).
We have fixed the zeromq timeout bug, and will hopefully be able to reproduce the issue outside of the operator
soon to report upstream.

## Presentations

 - [Kubecon 2023](https://t.co/vjRydPx1rb)
 - [HPC Knowledge Meeting 2023](https://hpckp.org/talks/cloud-and-hpc-convergence-flux-for-job-management-on-kubernetes/)

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

And the following external resources might be useful:

 - [Flux Cloud](https://github.com/converged-computing/flux-cloud): automation of experiments using the Flux Operator
 - [Flux RESTful API](https://github.com/flux-framework/flux-restful-api): provides the interface for submitting jobs, if no command provided to the operator.
 - [Python SDK](sdk/python): for deploying MiniClusters and port forwarding.
 - [Flux HPC Examples](https://github.com/rse-ops/flux-hpc) containers and CRD for the operator to run Flux with HPC workloads (under development)

**Note** this project is actively under development, and you can expect change and improvements!
We apologize for bugs you run into, and hope you tell us soon so we can work on resolving them.

## TODO

 - add [jar](wget https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/6.6.0/openapi-generator-cli-6.6.0.jar -O openapi-generator-cli.jar) to makefile
 - Ensure that curve.cert is a variable (path) in case we create a custom one.
 - Convert ML (e.g., mnist and pytorch) to use base containers instead of Singularity
 - When JobSet is available we need it to say job is successful when main application container is done.
 - Investigate why initContainers approach is not working with emptyDir (asked but no response in wg-batch slack)

## License

HPCIC DevTools is distributed under the terms of the MIT license.
All new contributions must be made under this license.

See [LICENSE](https://github.com/converged-computing/cloud-select/blob/main/LICENSE),
[COPYRIGHT](https://github.com/converged-computing/cloud-select/blob/main/COPYRIGHT), and
[NOTICE](https://github.com/converged-computing/cloud-select/blob/main/NOTICE) for details.

SPDX-License-Identifier: (MIT)

LLNL-CODE- 842614
