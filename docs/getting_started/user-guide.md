# User Guide

Welcome to the Flux Operator user guide! If you come here, we are assuming you have a cluster
with the Flux Operator installed and are interested to submit your own [custom resource](custom-resource-definition.md)
to create a MiniCluster, or that someone has already done it for you. If you are a developer wanting to work
on new functionality or features for the Flux Operator, see our [Developer Guides](../development/index.md) instead.

## Containers Available

All containers are provided under [ghcr.io/flux-framework/flux-operator](https://github.com/flux-framework/flux-operator/pkgs/container/flux-operator). The latest tag is the current main branch, a "bleeding edge" version,
and we provide the other pinned containers in case you want a previous version:

 - [ghcr.io/flux-framework/flux-operator:feb-2023](https://github.com/flux-framework/flux-operator/pkgs/container/flux-operator): the version used for Kubecon experiments, and before storage (minikube and Google Cloud example) were added.

## Install

### Quick Install

This works best for production Kubernetes clusters, and comes down to downloading the latest yaml config, and applying it.

```bash
wget https://raw.githubusercontent.com/flux-framework/flux-operator/main/examples/dist/flux-operator.yaml
kubectl apply -f flux-operator.yaml
```

Note that from the repository, this config is generated with:

```bash
$ make build-config
```

and then saved to the main branch where you retrieve it from.

### Development Install

If you are trying this out on your own, here is a quick start to getting the operator installed on MiniKube (or similar),
and installing directly from the repository (with the development use case in mind).
This assumes some experience with Kubernetes (or applying yaml configs) and using MiniKube or similar.
More advanced users can try out the [Production](#production) install detailed below.

This setup is intended if you want to clone the codebase and use the same tools that we use
to develop! You'll first want to clone the codebase.

```bash
$ git clone https://github.com/flux-framework/flux-operator
$ cd flux-operator
```

And then start a cluster with minikube:

```bash
$ minikube start
```

And make a flux operator namespace

```bash
$ kubectl create namespace flux-operator
namespace/flux-operator created
```
You can set this namespace to be the default (if you don't want to enter `-n flux-operator` for future commands:

```bash
$ kubectl config set-context --current --namespace=flux-operator
```

If you haven't ever installed minkube, you can see [install instructions here](https://minikube.sigs.k8s.io/docs/start/).
And then officially build the operator,

```bash
$ make
```

(optionally) to make your manifests:

```bash
$ make manifests
```

And install. Note that this places an executable `bin/kustomize` that you'll need to delete first if you make install again.

```bash
$ make install
```

At this point, you can kubectl apply your custom resource definition to define your MiniCluster to your cluster to
either run a job or start a Flux Mini Cluster.

```bash
$ kubectl apply -f config/samples/flux-framework.org_v1alpha1_minicluster.yaml
```

Note that we have other examples (using the web interface in [examples/flux-restful](https://github.com/flux-framework/flux-operator/tree/main/examples/flux-restful)
and headless examples for testing in [examples/tests](https://github.com/flux-framework/flux-operator/tree/main/examples/tests)).
When you are all done, cleanup with `kubectl delete` commands and/or!

```bash
$ make clean
```

And to stop MiniKube.

```bash
$ minikube stop
```

### Tool Install

We maintain a tool [Flux Cloud](https://github.com/converged-computing/flux-cloud) that is able to bring up clusters, install the operator,
and optionally run experiments and bring them down. We currently support a handful of clouds
(AWS and Google) and if you find yourself wanting a way to easily generate and save results
for experiments, this might be the way to go. If you have a cloud or environment you
want to deploy to that isn't supported, please [let us know](https://github.com/converged-computing/flux-cloud/issues).


## Next Steps

For next steps, you might do the following:

### 1. Verify Install

Regardless of what you chose above, from this point on (after the operator is installed)
there are some checks that you can do to see that everything worked.
First, ensure the `operator-system` namespace was created:

```bash
$ kubectl get namespace
```
```console
NAME              STATUS   AGE
default           Active   12m
kube-node-lease   Active   12m
kube-public       Active   12m
kube-system       Active   12m
operator-system   Active   11s
```
```bash
$ kubectl describe namespace operator-system
```
```console
Name:         operator-system
Labels:       control-plane=controller-manager
              kubernetes.io/metadata.name=operator-system
Annotations:  <none>
Status:       Active

No resource quota.

No LimitRange resource.
```

You can find the name of the operator pod as follows:

```bash
$ kubectl get pod --all-namespaces -o wide
```
```console
      <none>
operator-system   operator-controller-manager-6c699b7b94-bbp5q   2/2     Running   0             80s   192.168.30.43    ip-192-168-28-166.ec2.internal   <none>           <none>
```

### 2. Create Namespace

You'll likely Make your namespace for the flux-operator:

```bash
$ kubectl create namespace flux-operator
```

### 3. Validate your container (optional)

Your main container (with flux installed) has a basic [set of requirements](https://flux-framework.org/flux-operator/development/developer-guide.html?h=container#container-requirements) and we provide a simple tool to sanity check the most simple of these requirements, the [Flux Operator Validator](https://github.com/converged-computing/flux-operator-validator).
You are encouragd to run this script, although it's not required- you can just as easily go
through the list and verify the points on your own.

### 4. Apply your custom resource definition

Ensure that your custom resource definition matches the namespace you just created.
Then apply your CRD. You can use the default [testing one from the repository](https://github.com/flux-framework/flux-operator/blob/main/config/samples/flux-framework.org_v1alpha1_minicluster.yaml)
or any in our [examples](https://github.com/flux-framework/flux-operator/tree/main/examples) folder. Here is using the default we provide:

```bash
$ wget https://raw.githubusercontent.com/flux-framework/flux-operator/main/config/samples/flux-framework.org_v1alpha1_minicluster.yaml
$ kubectl apply -f flux-framework.org_v1alpha1_minicluster.yaml
```

Please [let us know](https://github.com/flux-framework/flux-operator) if you would like an example type added - we have plans for many more
but are prioritizing them as we see them needed. And now you can get logs for the manager:

```bash
$ kubectl logs -n operator-system operator-controller-manager-6c699b7b94-bbp5q
```

And then watch your jobs!

```bash
$ kubectl get -n flux-operator pods
```

And don't forget to clean up! Leaving on resources by accident is expensive! This command
will vary depending on the cloud you are using. Either way, it's good to check the web console too to ensure you didn't miss anything.
Next, you might be interested in [ways to submit jobs](jobs.md) or how to build images in our [Developer Guides](../development/developer-guide.md).
