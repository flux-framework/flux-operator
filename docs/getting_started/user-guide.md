# User Guide

Welcome to the Flux Operator user guide! If you come here, we are assuming you have a cluster
with the Flux Operator installed and are interested to submit your own [custom resource](custom-resource-definition.md)
to create a MiniCluster, or that someone has already done it for you. If you are a developer wanting to work
on new functionality or features for the Flux Operator, see our [Developer Guides](../development/index.md) instead.

## Containers Available

All containers are provided under [ghcr.io/flux-framework/flux-operator](https://github.com/flux-framework/flux-operator/pkgs/container/flux-operator). The latest tag is the current main branch, a "bleeding edge" version, and we also provide [releases](https://github.com/flux-framework/flux-operator/releases), each of which has YAML for x86 or ARM associated with a release container. 

## v1alpha1 

For dates before June 30, 2023 (that must be used with the corresponding GitHub releases or yamls) we provide the other pinned containers in case you want a previous version:

 - [ghcr.io/flux-framework/flux-operator:feb-2023](https://github.com/flux-framework/flux-operator/pkgs/container/flux-operator): the version used for Kubecon experiments, and before storage (minikube and Google Cloud example) were added.
 - [ghcr.io/flux-framework/flux-operator:april-2023](https://github.com/flux-framework/flux-operator/pkgs/container/flux-operator): the version directly before the refactor to remove the certificate generator pod (3.3)

These were primarily experimental versions run for experiments like Kubecon!

## Install

### Quick Install

We generally recommend that you install a [release](https://raw.githubusercontent.com/flux-framework/flux-operator/main/examples/dist/flux-operator.yaml), e.g.,

```bash
VERSION=0.1.0

# For x86
kubectl apply -f https://github.com/flux-framework/flux-operator/releases/download/${VERSION}/flux-operator.yaml

# For ARM
kubectl apply -f https://github.com/flux-framework/flux-operator/releases/download/${VERSION}/flux-operator-arm.yaml
```

You can also install from the current main branch "bleeding edge" latest:

```bash
kubectl apply -f https://raw.githubusercontent.com/flux-framework/flux-operator/main/examples/dist/flux-operator.yaml
kubectl apply -f https://raw.githubusercontent.com/flux-framework/flux-operator/main/examples/dist/flux-operator-arm.yaml
```

Note that from the repository, these configs are generated with:

```bash
$ make build-config
$ make build-config-arm
```

and then saved to the main branch or release where you retrieve it from.

### Helm Install

We optionally provide an install with helm, which you can do either from the charts in the repository:

```bash
$ git clone https://github.com/flux-framework/flux-operator 
$ cd flux-operator
$ helm install ./chart
```

Or directly from GitHub packages (an OCI registry):

```bash
# helm prior to v3.8.0
$ export HELM_EXPERIMENTAL_OCI=1
$ helm pull oci://ghcr.io/flux-framework/flux-operator-helm/chart
```
```console
Pulled: ghcr.io/flux-framework/flux-operator-helm/chart:0.1.0
```

And install!

```bash
$ helm install chart-0.1.0.tgz 
```
```console
NAME: flux-operator
LAST DEPLOYED: Fri Mar 24 18:36:18 2023
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

### Development Install

For developer instructions, please see our [developer documentation](../development/index.md).

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
Then apply your CRD. You can use [an example](https://github.com/flux-framework/flux-operator/blob/main/exampless)
Here is using a default we provide:

```bash
$ kubectl apply -f https://raw.githubusercontent.com/flux-framework/flux-operator/main/config/samples/flux-framework.org_v1alpha1_minicluster.yaml
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
