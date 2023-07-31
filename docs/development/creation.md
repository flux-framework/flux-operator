# Operator Creation

These sections will walk through some of the steps that @vsoch took to create the controller using the operator-sdk, and challenges she faced.

## 1. Installation

First, [install the operator-sdk](https://sdk.operatorframework.io/docs/installation/) for your platform.
At the end of this procedure it should be on your path.

```bash
$ which operator-sdk
/usr/local/bin/operator-sdk
```

## 2. Start A Development Cluster

You can use [minikube](https://minikube.sigs.k8s.io/docs/start/):

```bash
$ minikube start

# or for the first time
$ minikube start init
```

[Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) was recommended to me, however it sent me down a rabbit hole of "something isn't working and I don't know why"
for about a day, and when I realized it was specific to kind, I went back to minikube.

## 3. Local Workspace

At this point, I made sure I was in this present working directory, and I created
a new (v2) module and then "init" the operator:

```bash
$ go mod init flux-framework/flux-operator
$ operator-sdk init
```

Note that you don't need to do this, obviously, if you are using the existing operator here!

## 4. Create Controller

Now let's create a controller, and call it Flux (again, no need to do this if you are using the one here).

```bash
$ operator-sdk create api --version=v1alpha1 --kind=MiniCluster
```
(say yes to create a resource and controller). Make sure to install all dependencies (I think this might not be necessary - I saw it happen when I ran the previous command).

```bash
$ go mod tidy
$ go mod vendor
```

And then I started working on the actual content of the files generated.
You can look at the [User Guide](http://flux-framework.org/flux-operator/getting_started/index.html) for using the operator,
or other links in the Developer Guides here. For the above steps, I found the following resources really useful:

 - [RedHat OpenShift API Spec](https://docs.openshift.com/container-platform/3.11/rest_api/objects/index.html#objectmeta-meta-v1) for digging into layers of objects
 - [Kubernetes API](https://github.com/kubernetes/api/blob/2f9e58849198f8675bc0928c69acf7e50af77551/apps/v1/types.go): top level folders apps/core/batch useful!
