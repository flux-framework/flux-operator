# Flux Operator

![docs/09-07-2022/design-three-team.png](docs/09-07-2022/design-three-team.png)

🚧️ Under Construction! 🚧️

This is currently a scaffolding that is having functionality added as @vsoch figures it out.
Come back and check later for updates! To work on this operator you should:

 - Have a recent version of Go installed (1.18.1)
 - Have minikube installed
  
The sections below will describe:

 - [Organization](#organization) and design
 - [Quick Start](#quick-start) for using the operator
 - [Using the Operator](#using-the-operator): in more detail, as it is provided here
 - [Making the Operator](#making-the-operator): steps to initially create the operator

## Organization

The basic idea is that we present the idea of a **MiniCluster** that is a custom resource definition (CRD)
that defines a job container (that must have flux) that (when submit) will create a set of config maps,
secrets (e.g., tls), and the final Batch job that has the pod containers running with flux. Since
this is a batchv1.Job, it will have states that we can track.

And you can find the following:

 - [Flux Controllers](controllers/flux) are under `controllers/flux` for the `MiniCluster`
 - [API Spec](api/v1alpha1/) are under `api/v1alpha1/` also for `MiniCluster`
 - [Packages](pkg) include supporting packages for job conditions (state), if we eventually want that.
 - [TODO.md](TODO.md) is a set of things to be worked on, if you'd like to contribute!
 - [Documentation](docs) is currently a place to document design, and eventually can be more user-facing docs.

## Design

 - A **MiniCluster** is an [indexed job](https://kubernetes.io/docs/tasks/job/indexed-parallel-processing-static/) so we can create N copies of the "same" base containers (each with flux, and the connected workers in our cluster)
 - The flux config is written to a volume at `/etc/flux/config` (created via a config map) as a brokers.toml file.
 - We use an [initContainer](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) with an Empty volume (shared between init and worker) to generate the curve certificates (`/mnt/curve/curve.cert`). The broker sees them via the definition of that path in the broker.toml in our config directory mentioned above. Currently ever container generates its own curve.cert so this needs to be updated to have just one.
 - Networking is a bit of a hack - we have a wrapper starting script that essentially waits until a file is populated with hostnames. While it's waiting, we are waiting for the pods to be created and allocated an ip address, and then we write the addresses to this update file (that will echo into `/etc/hosts`). When the Pod is re-created with the same ip address, the second time around the file is run to update the hosts, and then we submit the job.
 - When the hosts are configured, the main rank (pod 0) does some final setup, and runs the job via the flux user. The others start flux with a sleep command.

## Quick Start

Know how this stuff works? Then here you go!

```bash
# Clone the source code
$ git clone https://github.com/flux-framework/flux-operator
$ cd flux-operator
```

### Local Development

Ensure localDeploy is set to true in your CRD so you don't ask for a persistent volume claim!

```
# Start a minikube cluster
$ minikube start

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

There is also a courtesy function to clean, and apply the samples:

```bash
$ make clean  # remove old flux-operator namespaced items
$ make apply  # apply the setup and job config
$ make run    # make the cluster (e.g.., setting up the batch job)
```

or run all three for easy development! The below command ensures to submit a job after the FluxSetup is applied so it won't be ignored (jobs submit before there is a cluster setup are ignored and require the user to re-submit).

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

### Production Deployment

## Deploy Operator 

To deploy the flux operator, you'll need to start with the same clone:

```bash
$ git clone https://github.com/flux-framework/flux-operator
$ cd flux-operator
```

A deploy will use the latest docker image [from the repository](https://github.com/orgs/flux-framework/packages?repo_name=flux-operator):

```bash
$ make deploy
```
```console
...
namespace/operator-system created
customresourcedefinition.apiextensions.k8s.io/miniclusters.flux-framework.org unchanged
serviceaccount/operator-controller-manager created
role.rbac.authorization.k8s.io/operator-leader-election-role created
clusterrole.rbac.authorization.k8s.io/operator-manager-role configured
clusterrole.rbac.authorization.k8s.io/operator-metrics-reader unchanged
clusterrole.rbac.authorization.k8s.io/operator-proxy-role unchanged
rolebinding.rbac.authorization.k8s.io/operator-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/operator-manager-rolebinding unchanged
clusterrolebinding.rbac.authorization.k8s.io/operator-proxy-rolebinding unchanged
configmap/operator-manager-config created
service/operator-controller-manager-metrics-service created
deployment.apps/operator-controller-manager created
```

Ensure the `operator-system` namespace was created:

```bash
$ kubectl get namespace
NAME              STATUS   AGE
default           Active   12m
kube-node-lease   Active   12m
kube-public       Active   12m
kube-system       Active   12m
operator-system   Active   11s
```
```bash
$ kubectl describe namespace operator-system
Name:         operator-system
Labels:       control-plane=controller-manager
              kubernetes.io/metadata.name=operator-system
Annotations:  <none>
Status:       Active

No resource quota.

No LimitRange resource.
```

And you can find the name of the operator pod as follows:

```bash
$ kubectl get pod --all-namespaces -o wide
```
```console
      <none>
operator-system   operator-controller-manager-6c699b7b94-bbp5q   2/2     Running   0             80s   192.168.30.43    ip-192-168-28-166.ec2.internal   <none>           <none>
```

Make your namespace for the flux-operator:

```bash
$ kubectl create namespace flux-operator
```

Then apply your CRD (the localDeploy can be false for an actual cluster with persistent volume claims):

```bash
$ make apply
# OR
$ kubectl apply -f config/samples/flux-framework.org_v1alpha1_minicluster.yaml 
```

And now you can get logs for the manager:

```bash
$ kubectl logs -n operator-system operator-controller-manager-6c699b7b94-bbp5q
```

And then watch your jobs as before!

```bash
$ make list
$ make apply
```

And don't forget to clean up! Leaving on resources by accident is expensive!

## Clean up

Make sure you clean everything up!

```bash
$ eksctl delete cluster -f eks-cluster-config.yaml
```
It might be better to add `--wait`, which will wait until all resources are cleaned up:

```bash
$ eksctl delete cluster -f eks-cluster-config.yaml --wait
```
Either way, it's good to check the web console too to ensure you didn't miss anything.


### Build Images

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

### Starting Fresh

If you want to blow up your minikube and start fresh (pulling the container again too):

```bash
make reset
```

## Using the Operator

If you aren't starting from scratch, then you can use the code here to see how things work!

```bash
$ git clone https://github.com/flux-framework/flux-operator
$ cd flux-operator
```

### 1. Start a Cluster

First, start a cluster with minikube:

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

### 3. Cleanup

When cleaning up, you can control+c to kill the operator from running, and then:

```bash
$ make clean
```

And then:

```bash
$ minikube stop
```

## What is Happening?

If you follow the commands above, you'll see a lot of terminal output, and it might not be clear what
is happening. Let's talk about it here. Generally, you'll first see the config maps and supporting resources 
being created. Since we are developing (for the time being) on a local machine, instead of a persistent volume
claim (which requires a Kubernetes cluster with a provisioner) you'll get a persistent volume
written to `/tmp` in the job namespace. If you try to use the latter it typically freezer.

The first time the pods are created, they won't have ips (yet) so you'll see an empty list in the logs.
As they are creating and getting ips, after that is finished you'll see the same output but with a 
lookup of hostnames to ip addresses, and after it will tell you the cluster is ready.

```
1.6629325562267003e+09  INFO    minicluster-reconciler  🌀 Mini Cluster is Ready!
```
When you are waiting and run `make log` in a separate terminal you'll see output from one of the pods 
in the job. Typically the first bit of time you'll be waiting:

```bash
$ make log
kubectl logs -n flux-operator job.batch/flux-sample
Found 6 pods, using pod/flux-sample-0-njnnd
Host updating script not available yet, waiting...
```
It's waiting for the `/flux_operator/update_hosts.sh` script. When this is available, it will be found
and the job setup will continue, first adding the found hosts to `/etc/hosts` and then (for the main node,
which typically is `<name>-0`). When this happens, you'll see the host file cat to the screen:

```console
Host updating script not available yet, waiting...
# Kubernetes-managed hosts file.
127.0.0.1       localhost
::1     localhost ip6-localhost ip6-loopback
fe00::0 ip6-localnet
fe00::0 ip6-mcastprefix
fe00::1 ip6-allnodes
fe00::2 ip6-allrouters
172.17.0.4      flux-sample-1.flux-sample.flux-operator.svc.cluster.local       flux-sample-1
172.17.0.2 flux-sample-0-flux-sample.flux-operator.svc.cluster.local flux-sample-0
172.17.0.4 flux-sample-1-flux-sample.flux-operator.svc.cluster.local flux-sample-1
172.17.0.6 flux-sample-2-flux-sample.flux-operator.svc.cluster.local flux-sample-2
172.17.0.7 flux-sample-3-flux-sample.flux-operator.svc.cluster.local flux-sample-3
172.17.0.5 flux-sample-4-flux-sample.flux-operator.svc.cluster.local flux-sample-4
172.17.0.8 flux-sample-5-flux-sample.flux-operator.svc.cluster.local flux-sample-5
flux-sample-1 is sleeping waiting for main flux node
```

And then final configs are created, the flux user is created, and the main
node creates the certificate and we start the cluster. You can look at 
[controllers/flux/templates.go](controllers/flux/templates.go)
for all the scripts and logic that are run. It takes about ~90 seconds for the
whole thing to come up and run. If `make log` doesn't show you the main node
(where we run the command) you get `make list` to get the identifier and then:

```bash
$ kubectl logs -n flux-operator flux-sample-0-zfbvm
# or
$ ./script/log.sh 
```

What is happening now is that the main rank (0) finishes and the others sort of are waiting
around:

```console
$ make list
kubectl get -n flux-operator pods
NAME                  READY   STATUS      RESTARTS   AGE
flux-sample-0-wvs8w   0/1     Completed   0          11m
flux-sample-1-9cz5c   1/1     Running     0          11m
flux-sample-2-hbcrb   1/1     Running     0          11m
flux-sample-3-lzv4s   1/1     Running     0          11m
flux-sample-4-fzxgf   1/1     Running     0          11m
flux-sample-5-9pnwf   1/1     Running     0          11m
```
I think we want to either use a different starting command, or have a cleanup in the reconciler. TBA!

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
$ operator-sdk create api --version=v1alpha1 --kind=MiniCluster
```
(say yes to create a resource and controller). Make sure to install all dependencies (I think this might not be necessary - I saw it happen when I ran the previous command).

```bash
$ go mod tidy
$ go mod vendor
```

**under development**

And then see the instructions above for [using the operator](#using-the-operator).

### 5. Debugging Yaml

Since many configs are created in the operator, I always write them out in yaml to the [yaml](yaml)
directory. We can remove these bits of the code after we are done. I also found the following debugging commands useful:

```bash
# Why didn't my statefulset create?
kubectl describe -n flux-operator statefulset
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

One time I messed something up and my metrics server was still running (and I couldn't start again) and I killed it like this:

```bash
$ kill $(lsof -t -i:8080)
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

#### License

This work is heavily inspired from [kueue](https://github.com/kubernetes-sigs/kueue) for the design. I am totally new to operator design and tried
several basic designs, and decided to mimic this setup (a simplified version) for a first shot, and for my own learning. kueue at the time of
this was also under the [Apache-2.0](https://github.com/kubernetes-sigs/kueue/blob/ec9b75eaadb5c78dab919d8ea6055d33b2eb09a2/LICENSE) license.

SPDX-License-Identifier: Apache-2.0

LLNL-CODE-764420
