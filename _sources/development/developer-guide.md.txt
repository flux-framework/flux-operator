# Developer Guide

This developer guide includes complete instructions for setting up a developer
environment.

## Setup

To work on this operator you should:

 - Have a recent version of Go installed (1.18.1)
 - Have minikube installed

**Important** For minikube, make sure to enable [DNS plugins](https://minikube.sigs.k8s.io/docs/handbook/addons/ingress-dns/).

```bash
$ minikube addons enable ingress
$ minikube addons enable ingress-dns
```

The basic Flux networking (pods seeing one another) won't work if your cluster does not support DNS.
You also won't be able to expose the service with minikube service if you don't do the above (but port-forward would technically work)
You'll then also want to clone the repository.

```bash
# Clone the source code
$ git clone https://github.com/flux-framework/flux-operator
$ cd flux-operator
```

### Local Development

After cloning, you need to create your MiniKube cluster before doing anything else! 

#### 1. Quick Start

Here is a quick start for doing that, making the namespace, and installing the operator.

```console
# Start a minikube cluster
$ minikube start

# Make a flux operator namespace
$ kubectl create namespace flux-operator
namespace/flux-operator created
```

Here is how to build and install the operator:

```
# Build the operator
$ make

# How to make your manifests
$ make manifests

# And install. This places an executable "bin/kustomize"
$ make install
```

#### 2. Configs

The job configs - custom resource definitions or "CRD" can be found in [config/samples](https://github.com/flux-framework/flux-operator/tree/main/config/samples).
Before launching any jobs, making sure `localDeploy` is set to true in your CRD so you don't ask for a persistent volume claim!

```yaml
spec:
# Set to true to use volume mounts instead of volume claims
  localDeploy: true
```

The default is set to true (so a local volume in `/tmp` is used) but this won't work on an actual cloud Kubernetes cluster,
and vice versa - the volume claim won't work locally. When you are sure this is good,
here is how to "launch" the Mini Cluster (or if providing a command, an ephemeral job):

```bash
$ kubectl apply -f config/samples/flux-framework.org_v1alpha1_minicluster.yaml
```

there is a courtesy function to clean, and apply the samples:

```bash
$ make clean  # remove old flux-operator namespaced items
$ make apply  # apply the setup and job config
$ make run    # make the cluster
```

or run all three for easy development!

```bash
$ make redo  # clean, apply, and run
```

To see logs for the job, you'd do:

```bash
$ kubectl logs -n flux-operator job.batch/flux-sample
```

And this is also:

```bash
$ make log
```

But the above gives you a (somewhat random?) pod. If you want to see a specific one do:

```bash
./script/log.sh flux-sample-0-b5rw6
```

E.g., rank 0 is the "main" one. List running pods (each pod is part of a batch job)

```bash
$ make list
```

And shell into one with the helper script:

```bash
./script/shell.sh flux-sample-0-b5rw6
```

### Interacting with Services

Currently, the most reliable thing to do is port forward:

#### port-forward

If we run as a ClusterIP, we can accomplish the same with a one off `kubectl port-forward`:

```console
kubectl port-forward -n flux-operator flux-sample-0-zdhkp 5000:5000
Forwarding from 127.0.0.1:5000 -> 5000
```

This means you can open [http://localhost:5000](http://localhost:5000) to see the restful API (and interact with it there).

If you want to use a minikube service, this seems to work, but is spotty - I think because
minikube is expecting the service to be available from any pod (it is only running from index 0).
If you want to try this:

```console
$ minikube service -n flux-operator flux-restful-service --url=true
```

But for now I'm developing with port forward.


## Build Images

If you want to build the "production" images - here is how to do that!
This happens in our Docker CI, however you can build (and deploy if you are an owner) them too!

```bash
$ make docker-build
$ make docker-push
```
```bash
# operator lifecycle manager
$ operator-sdk olm install
$ make bundle
$ make bundle-build
$ make bundle-push
```

And for the catalog:

```bash
$ make catalog-build
$ make catalog-push
```


## Documentation

The documentation is provided in the `docs` folder of the repository, and generally most content that you might want to add is under
`getting_started`. For ease of contribution, files that are likely to be updated by contributors (e.g., mostly everything but the module generated files)
are written in markdown. If you need to use [toctree](https://www.sphinx-doc.org/en/master/usage/restructuredtext/directives.html#table-of-contents) you should not use extra newlines or spaces (see index.md files for examples). The documentation is also provided in Markdown (instead of rst or restructured syntax) to make contribution easier for the community.

Finally, we recommend you use the same development environment also to build and work on
documentation. The reason is because we import the app to derive docstrings,
and this will require having Flux.

### Install Dependencies and Build

The documentation is built using sphinx, and generally you can 
create a virtual environment:

```bash
$ cd docs
$ python -m venv env 
$ source env/bin/activate
```
And then install dependencies:

```console
$ pip install -r requirements.txt

# Build the docs into _build/html
$ make html
```

### Preview Documentation

After `make html` you can enter into `_build/html` and start a local web
server to preview:

```console
$ python -m http.server 9999
```

And open your browser to `localhost:9999`
