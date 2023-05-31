# Flux Bootrap

> All about dem' FLUB! (FLUx Bootstrap) 🤪️

Flubbing? Flubbed? Flub-ified? So many ideas! 

This tests a [new feature](https://github.com/flux-framework/flux-core/pull/5184) to add support for Flux Bootstrap,
which means we can create a cluster with a max size that can support adding another broker, and then
create a second broker and easily connect the two. In other words, we would be connecting MiniClusters.

## Usage

First, let's create a kind cluster. From the context of this directory:

```bash
$ kind create cluster --config ../../kind-config.yaml
```

And then install the operator, create the namespace, and apply the MiniCluster YAML here. Note
that you'll likely need to build the custom image for this branch and then apply the development CRD
(this feature is not merged into Flux and is considered experimental):

```bash
$ make test-deploy
$ kubectl apply -f ../../dist/flux-operator-dev.yaml
$ kubectl create namespace flux-operator
```

### Create First MiniCluster

Let's first create the base MiniCluster - number one! This cluster is going to be the one
we have up first, and we allow the second cluster to see and connect to. I'm not sure
if the ordering is important (e.g., we could create a cluster that knows about an external one
that doesn't exist yet) but for this first experiment this order seemed logical to me.
Here is how to create:

```bash
$ kubectl apply -f ./minicluster-1st.yaml
```

If you look at the [minicluster-1st.yaml](minicluster-1st.yaml) you'll see it is just
an interactive mode cluster, meaning we start the broker and do nothing else. However,
we add one attribute to explicitly define the service selector label, and we do this
so the headless service for this MiniCluster is not scoped to the one job:

```yaml
jobSelector: connected-service
```

You should wait for the pods to be running before moving on to the next step.

### Create Second MiniCluster

This second MiniCluster has extra (new!) attributes that will make it possible to 
see (and connect to) the already existing MiniCluster. It could also be the case
that this doesn't have to exist, but I haven't tested that yet. Specifically,
this part of the custom resource definition is what is important:

```yaml

  # Ensure this job can fall under the same networking namespace
  jobSelector: connected-service
  
  # Broker options
  flux:

    # Allow the minicluster-1st brokers to connect
    bootServer: flux-sample-1-0.flux-service.flux-operator.svc.cluster.local

    # The number of nodes the first server has
    bootServerSize: 2
```

You can check the logs of the broker pods (the index 0s) to make sure that everything started OK.
Since we told cluster 2 about cluster 1, let's shell into 2:

```bash
$ kubectl exec -it -n flux-operator flux-sample-2-0-xxx bash
```

Connect to the broker as the flux user:

```bash
$ sudo -u flux -E $(env) -E HOME=/home/flux flux proxy local:///run/flux/local bash
```

I currently can see the other node! 

```console
     STATE NNODES   NCORES    NGPUS NODELIST
      free      2        8        0 flux-sample-2-[0-1]
 allocated      0        0        0 
      down      2        8        0 flux-sample-1-[0-1]
```

However they are down, so likely I have a value incorrect somewhere. I suspect this is an issue
with formatting the hostnames vs. service, and the value that is needed for the flux broker option.

**Questions**:

- Do we need to change the tbon topology to be kary:2 or similar?
- I am assuming the hosts (from the second cluster) need to be defined in `[[resource.config]]` AND the flux R resource spec
- What should be the format of the broker option?

When you are done:

```bash
$ kubectl delete -f minicluster-1st.yaml
$ kubectl delete -f minicluster-2nd.yaml
```
