# Debugging

These sections include common errors that you might run into developing the operator.

## Errors

### Pull Errors

I had a problem where the containers I wanted to deploy seemed to time-out during pull. In the pod logs I would see:

> Context deadline exceeded

In this case, if you manually pull first with ssh (and then don't ask to re-pull) you should be able to get around this.

```bash
$ minikube ssh docker pull ghcr.io/rse-ops/lammps:flux-sched-focal-v0.24.0
```

### Service Already on Port 8080

One time I messed something up and my metrics server was still running (and I couldn't start again) and I killed it like this:

```bash
$ kill $(lsof -t -i:8080)
```

### CRD should be installed

If you see something like:

```bash
1.6605195805812113e+09	ERROR	controller-runtime.source	if kind is a CRD, it should be installed before calling Start	{"kind": "Flux.flux-framework.org", "error": "no matches for kind \"Flux\" in version \"flux-framework.org/v1alpha1\""}
```

You need to remove the previous kustomize and install the CRD again:

```bash
$ rm bin/kustomize
$ make install
```

### Configs not taking

If your resource (config/samples) don't seem to be taking, remember you need to apply them for changes to take
effect:

```bash
$ bin/kustomize build config/samples | kubectl apply -f -
```

### Fields not showing up

Since we are providing the instance as a reference (`&instance`) to get the fields from that you sometimes need to do:

```go
(*instance).Spec.Field
```
Otherwise it can show up as an empty string!

### Yaml

Since many configs are created in the operator, a strategy to debug what you are creating is to write them out verbatim in a `yaml`
folder alongside the repository. The "describe" commands of kubectl are useful for debugging specific Kubernetes objects:

```bash
# Why didn't my statefulset create?
kubectl describe -n flux-operator statefulset
```

## Troubleshooting

Here is a general debugging strategy. To view a resource (in the flux-operator namespace):

```bash
$ kubectl --namespace flux-operator get pods
```

I found this helpful for debugging the stateful set - e.g., it was stuck on ContainerCreating:

```bash
$ kubectl describe -n flux-operator pods
```

If you need to clean things up (ensuring you only have this one pod and service running first) I've found it easier to do:

```bash
$ kubectl delete pod --all

# Service shorthand
$ kubectl delete svc --all
$ kubectl delete statefulset --all

# ConfigMap shorthand
$ kubectl delete cm --all
```

And take a look at the project `Makefile` that has some commands pre-set for you, and feel free to add other
combinations that might be useful.
