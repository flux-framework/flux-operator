# Testing JobSet

In this tutorial we will test an implementation using JobSet. The installation of the operator is the equivalent,
but you will need to [install JobSet first](https://github.com/kubernetes-sigs/jobset).

## Setup

Create a kind cluster.

```bash
$ kind create cluster
```

Install the JobSet (it won't work if you don't!):

```bash
VERSION=v0.1.3
kubectl apply --server-side -f https://github.com/kubernetes-sigs/jobset/releases/download/$VERSION/manifests.yaml
```

Now let's try creating a hello world example MiniCluster with the JobSet.

```bash
$ kubectl create namespace flux-operator
$ kubectl apply -f examples/dist/flux-operator-dev.yaml
$ kubectl apply -f examples/tests/jobset/minicluster.yaml
```

This is a WIP! The JobSet (when installed as above) isn't seen by the operator:

```bash
2023-05-09T00:09:51Z    ERROR   Reconciler error        {"controller": "minicluster", "controllerGroup": "flux-framework.org", "controllerKind": "MiniCluster", "MiniCluster": {"name":"flux-sample","namespace":"flux-operator"}, "namespace": "flux-operator", "name": "flux-sample", "reconcileID": "1db5223d-1f8a-4e49-a219-8d3ba55de826", "error": "no kind is registered for the type v1alpha1.JobSet in scheme \"pkg/runtime/scheme.go:100\""}
```