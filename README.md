# Flux Operator

🚧️ Under Construction! 🚧️

This is currently a scaffolding that is having functionality added as @vsoch figures it out.
Come back and check later for updates! To work on this operator you should:

 - Have a recent version of Go installed (1.18.1)
 - Have minikube or kind installed
  
The sections below will describe:

 - [Organization](#organization) and design
 - [Quick Start](#quick-start) for using the operator
 - [Using the Operator](#using-the-operator): in more detail, as it is provided here
 - [Making the Operator](#making-the-operator): steps to initially create the operator

## Organization

The basic idea is that we have two controllers:

 - Flux: is a user facing controller (meaning we have a CR, a simple yaml the user can provide an image and entrypoint)
 - FluxSetup: is more an internal or admin controller, meaning no CR, but there is a CRD (custom resource definition).
 
And you can find the following:

 - [Flux Controllers](controllers/flux) are under `controllers/flux` for each of `Flux` and `FluxSetup`
 - [API Spec](api/v1alpha1/) are under `api/v1alpha1/` also for each of `Flux` and `FluxSetup`
 - [TODO.md](TODO.md) is a set of TODOs @vsoch took notes on as she developed.


## Quick Start

Know how this stuff works? Then here you go!

```bash
# Clone the source code
$ git clone https://github.com/flux-framework/flux-operator
$ cd flux-operator

# Start a minikube cluster or kind
$ minikube start
$ kind create cluster

# Make a flux operator namespace
$ kubectl create namespace flux-operator
namespace/flux-operator created

# Build the operator
$ make

# How to make your manifests
$ make manifests

# And install. This places an executable "bin/kustomize"
$ make install
```

And then edit your [config/samples](config/samples) and deploy them, and start the operator!

```bash
$ bin/kustomize build config/samples | kubectl apply -f -
```
```
flux.flux-framework.org/flux-sample configured
fluxsetup.flux-framework.org/flux-sample configured
```
```bash
$ make run
```


## Using the Operator

If you aren't starting from scratch, then you can use the code here to see how things work!

```bash
$ git clone https://github.com/flux-framework/flux-operator
$ cd flux-operator
```

### 1. Start a Cluster

First, start a cluster with minikube or kind.

```bash
$ minikube start
```
```bash
$ kind create cluster
```

I tried loading the image for flux first (into kind). I'm not sure this is required (it's definitely not required for Kubernetes).

```bash
$ docker pull fluxrm/flux-sched:focal
$ kind load docker-image fluxrm/flux-sched:focal
```

And make a flux operator namespace

```bash
$ kubectl create namespace flux-operator
namespace/flux-operator created
```

If you haven't ever installed minkube, you can see [install instructions here](https://minikube.sigs.k8s.io/docs/start/),
or [instructions here for kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installing-from-release-binaries).

### 2. Build

And then officially build.

```bash
$ make
```

To make your manifests:

```bash
$ make manifests
```

And install. Note that this places an executable [bin/kustomize](bin/kustomize) that you'll need to delete first if you make install again.

```bash
$ make install
```

### 3. Deploy

Note that you will be using the config yamls [here](config/samples/_v1alpha1_lolcow.yaml) to start, which include a greeting and port.
We will look at these later for demonstrating how the operator watches for changes. Apply your configs (kustomize is in the bin).

```bash
$ bin/kustomize build config/samples | kubectl apply -f -
flux.flux-framework.org/job created
```

And finally, run it.

```bash
$ make run
```

And you should be able to open the web-ui:

```bash
$ minikube service job
```

Note that if you get a 404 page, do `kubectl get svc` and wait until the service goes from "pending" to "ready."

### 7. Cleanup

When cleaning up, you can control+c to kill the operator from running, and then:

```bash
$ kubectl delete pod --all
$ kubectl delete svc --all
```
And one of the following:

```bash
$ minikube stop
```
```bash
$ kind delete cluster 
```

## Making the operator

This section will walk through some of the steps that @vsoch took to create the controller using the operator-sdk, and challenges she faced.

### 1. Installation

First, [install the operator-sdk](https://sdk.operatorframework.io/docs/installation/) for your platform.
At the end of this procedure it should be on your path.
At the

```bash
$ which operator-sdk
/usr/local/bin/operator-sdk
```

### 2. Start A Development Cluster

You can either use [minikube](https://minikube.sigs.k8s.io/docs/start/):

```bash
$ minikube start

# or for the first time
$ minikube start init
```

or you can use [kind](https://kind.sigs.k8s.io/docs/user/quick-start/):

```bash
# Default cluster context name is `kind`.
$ kind create cluster
```

### 3. Local Workspace

At this point, I made sure I was in this present working directory, and I created
a new (v2) module and then "init" the operator:

```bash
$ go mod init flux-framework/flux-operator
$ operator-sdk init
```

Note that you don't need to do this, obviously, if you are using the existing operator here!

### 4. Create Controller

Now let's create a controller, and call it Flux (again, no need to do this if you are using the one here).

```bash
$ operator-sdk create api --version=v1alpha1 --kind=Flux
```
(say yes to create a resource and controller). Make sure to install all dependencies (I think this might not be necessary - I saw it happen when I ran the previous command).

```bash
$ go mod tidy
$ go mod vendor
```

At this point I chat with Eduardo about operator design, and we decided to go for a simple design:

- Flux: would have a CR that defines a job and exposes an entrypoint command / container to the user
- FluxSetup would have most of the content of [here](https://lc.llnl.gov/confluence/display/HFMCCEL/Flux+Operator+Design) and be more of an internal or admin setup.

To generate FluxSetup I think I could (maybe?) have run this command again, but instead I added a new entry to the [PROJECT](PROJECT) and then generated [api/v1alpha1/fluxsetup_types.go](api/v1alpha1/fluxsetup_types.go) from the `flux_types.go` (and changing all the references from Flux to FluxSetup). 
I also needed to (manually) make [controllers/fluxsetup_controller.go](controllers/fluxsetup_controller.go) and ensure it was updated to use FluxSetup, and then adding
it's creation to [main.go](main.go). Yes, this is a bit of manual work, but I think I'll only need to do it once.
At this point, I needed to try and represent what I saw in the various config files in this types file.

**under development**

And then see the instructions above for [using the operator](#using-the-operator).

### 5. Debugging Yaml

Since many configs are created in the operator, I always write them out in yaml to the [yaml](yaml)
directory. We can remove these bits of the code after we are done. I also found the following debugging commands useful:

```bash
# Why didn't my statefulset create?
kubectl describe statefulset flux-sample
```

## Useful Resources

I found the following resources really useful:

 - [RedHat OpenShift API Spec](https://docs.openshift.com/container-platform/3.11/rest_api/objects/index.html#objectmeta-meta-v1) for digging into layers of objects
 - [Kubernetes API](https://github.com/kubernetes/api/blob/2f9e58849198f8675bc0928c69acf7e50af77551/apps/v1/types.go): top level folders apps/core/batch useful!

## Troubleshooting

To view a resource (in the flux-operator namespace):

```bash
$ kubectl --namespace flux-operator get pods
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

If you see:

```bash
1.6605195805812113e+09	ERROR	controller-runtime.source	if kind is a CRD, it should be installed before calling Start	{"kind": "Flux.flux-framework.org", "error": "no matches for kind \"Flux\" in version \"flux-framework.org/v1alpha1\""}
```

You need to remove the previous kustomize and install the CRD again:

```bash
$ rm bin/kustomize
$ make install
```

If your resource (config/samples) don't seem to be taking, remember you need to apply them for changes to take
effect:

```bash
$ bin/kustomize build config/samples | kubectl apply -f -
```

Also remember since we are providing the instance as a reference (`&instance`) to get the fields from that
you need to do:

```go
(*instance).Spec.Field
```
Otherwise it shows up as an empty string.

## Wisdom

**from the kubebuilder slack**

### Learned Knowledge

- Reconciling should only take into account the spec of your object, and the real world.  Don't use status to hold knowledge for future reconcile loops.  Use a workspace object instead.
- Status should only hold observations of the reconcile loop.  Conditions, perhaps a "Phase", IDs of stuff you've found, etc.
- Use k8s ownership model to help with cleaning up things that should automatically be reclaimed when your object is deleted.
- Use finalizers to do manual clean-up-tasks
- Send events, but be very limited in how often you send events.  We've opted now to send events, essentially only when a Condition is modified (e.g. a Condition changes state or reason).
- Try not to do too many things in a single reconcile.  One thing is fine.  e.g. see one thing out of order?  Fix that and ask to be reconciled.  The next time you'll see that it's in order and you can check the next thing.  The resulting code is very robust and can handle almost any failure you throw at it.
- Add "kubebuilder:printcolums" markers to help kubectl-users get a nice summary when they do "kubectl get yourthing".
- Accept and embrace that you will be reconciling an out-of-date object from time to time.  It shouldn't really matter.  If it does, you might want to change things around so that it doesn't matter.  Inconsistency is a fact of k8s life.
- Place extra care in taking errors and elevating them to useful conditions, and/or events.  These are the most visible part of an operator, and the go-to-place for humans when trying to figure out why your code doesn't work.  If you've taken the time to extract the error text from the underlying system into an Event, your users will be able to fix the problem much quicker.

### What is a workspace?

A workspace object is when you need to record some piece of knowledge about a thing you're doing, so that later you can use that when reconciling this object. MyObject "foo" is reconciled; so to record the thing you need to remember, create a MyObjectWorkspace — Owned by the MyObject, and with the same name + namespace.  MyObjectWorkspace doesn't need a reconciler; it's simply a tool for you to remember the thing. Next time you reconcile a MyObject, also read your MyObjectWorkspace so you can remember "what happened last time". E.g. I've made a controller to create an EC2 instance, and we needed to be completely sure that we didn't make the "launch instance" API call twice.  EC2 has a "post once only" technique whereby you specify a nonce to avoid duplicate API calls.  You would write the nonce to the workspace use the nonce to call the EC2 API write any status info of what you observed to the status. Rremove the nonce when you know that you've stored the results (e.g. instance IDs or whatever) When you reconcile, if the nonce is set, you can re-use it because it means that the EC2 call failed somehow.  EC2 uses the nonce the second time to recognise that "heh, this is the same request as before ..." Stuff like this nonce shouldn't go in your status. Put simply, the status should really never be used as input for your reconcile.

Know that the scaffolded k8sClient includes a cache that automatically updates based on watches, and may give you out-of-date data (but this is fine because if it is out-of-date, there should be a reconcile in the queue already). Also know that there is a way to request objets bypassing a cache (look for APIReader).  This gives a read-only, but direct access to the API.  Useful for e.g. those workspace objects.
