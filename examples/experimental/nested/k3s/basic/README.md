# K3s

This is a small experiment to try and run kubernetes within the Flux operator,
which of course is within Kubernetes! This is a "quasi" Kubernetes
because we are going to try using k3s. 

## Background

For some background, I was during a survey of [tooling available](https://github.com/converged-computing/operator-experiments/tree/main/google/rootless-kubernetes)
and stumbled on a few attributes of k3s that I liked:

- I could run it in rootless mode
- I could get it working in docker-compose
- It would have separate commands to start a main server and register an agent

These qualities (I think) made it my first contender to try within Flux.
Right now I don't have a good sense of the limitations of an HPC environment,
but I can comment on the general challenges that I saw for each approach:

 - cgroups2 is required
 - if containerization is required, might be challenging due to systemctl usage

However, I do think there might be an avenue to pursue making this work on a more
traditional HPC system, and likely with a container. I don't understand
or know these environments well, so for now I decided to try in the Flux operator.
My strategy was first to try running k3s through a singularity container,
and that didn't work because of permissions. Then I decided to get k3s installed in a container with Flux
that could be used as a base image for the operator, and that's the approach
I'm taking here. We are currently running flux as root because I'm trying
to reproduce the docker-compose exampel (that works) and hopefully we can
step away from this!

## Usage

First, let's create a kind cluster. From the context of this directory:

```bash
$ kind create cluster --config ../../../kind-config.yaml
```

And then install the operator, create the namespace, and apply the MiniCluster YAML here.

```bash
$ kubectl apply -f ../../../dist/flux-operator.yaml
$ kubectl create namespace flux-operator
$ kubectl apply -f ./minicluster.yaml
```

If you watch the broker logs, a "successful state" is when you see a ton of output but nothing
is exiting. Once I saw it look consistent, I shelled in to the broker to try and apply
a yaml file to the k3s cluster within the Flux cluster (weird, right?)

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-jlsp6 bash
$ flux proxy local:///run/flux/local bash
```

These commands are from the bottom of the [start.sh](start.sh)
script, which technically will never be reached if the agents and server are running
without issue. Note that we can first try to look at our cluster:

```bash
$ kubectl --kubeconfig=./kubeconfig.yaml get nodes
NAME            STATUS   ROLES                  AGE     VERSION
flux-sample-2   Ready    <none>                 3m38s   v1.27.1+k3s-b32bf495
flux-sample-3   Ready    <none>                 3m38s   v1.27.1+k3s-b32bf495
flux-sample-0   Ready    control-plane,master   3m44s   v1.27.1+k3s-b32bf495
flux-sample-1   Ready    <none>                 3m38s   v1.27.1+k3s-b32bf495
```

And you can sanity check it's not hitting the local machine!

```bash
# kubectl get -n flux-operator pods
No resources found in flux-operator namespace.
```

Note that it would show this message even for namespaces that don't exist.
This still feels really weird because I'm looking at the same pods that I'm running on my machine as they
are seen in the container... but they are Kubernetes pods running a Flux node, and the Flux node running
another Kubernetes agent (kubelet). So they are [same, same, but different, but still same](https://youtu.be/7tTfL-DtpXk)! 😂️ 
An easy indicator is looking at the timestamps andn node names - they were created after, and they don't have the indexed job 
random letters after the hostname, they are just the hostnames seen by Flux. Here are the pods from the outside
(on my local machine) for comparison:

```bash
kubectl get -n flux-operator pods
NAME                  READY   STATUS    RESTARTS   AGE
flux-sample-0-k6xvs   1/1     Running   0          109s
flux-sample-1-xk5x9   1/1     Running   0          109s
flux-sample-2-bfwkr   1/1     Running   0          109s
flux-sample-3-kktlk   1/1     Running   0          109s
```

I also ran into a weird bug I couldn't explain - it had the default namespace as "flux-operator"
and this is probably inherited somewhere (and I didn't look into it) but instead I decided
to just set it back to be default:

```bash
kubectl --kubeconfig=./kubeconfig.yaml config set-context --current --namespace=default
```

I'm not super worried because We likely won't be dealing with this on an actual Flux cluster (that doesn't have a second layer of Kubernetes) 
and this is where I'd like to test it next. Despite that weirdness,
everything actually works:

```bash
$ kubectl --kubeconfig=./kubeconfig.yaml apply -f my-echo.yaml
until kubectl --kubeconfig=./kubeconfig.yaml rollout status deployment my-echo; do sleep 1; done
deployment "my-echo" successfully rolled out
```

And I actually found I don't need to target the kubeconfig.yaml, it seems to hit the right cluster
either way:

```bash
root@flux-sample-0:/workflow# kubectl get deploy
NAME      READY   UP-TO-DATE   AVAILABLE   AGE
my-echo   1/1     1            1           29s
```
```bash
$ kubectl get pods
NAME                       READY   STATUS    RESTARTS   AGE
my-echo-74dc6c4f7b-lh89p   1/1     Running   0          27s
```

But not on my local machine! So we are getting somewhere:

```bash
$ kubectl get deploy
No resources found in default namespace.
```

I don't know how an echoserver is supposed to work, but we would need to
expose the port to see the service. For now we can inspect that the server is running:

```bash
$ curl http://127.0.0.1:8080
```
```console
Hostname: my-echo-74dc6c4f7b-lh89p

Pod Information:
        -no pod information available-

Server values:
        server_version=nginx: 1.13.3 - lua: 10008

Request Information:
        client_address=10.42.0.4
        method=GET
        real path=/
        query=
        request_version=1.1
        request_uri=http://127.0.0.1:8080/

Request Headers:
        accept=*/*
        host=127.0.0.1:8080
        user-agent=curl/7.68.0

Request Body:
        -no body in request-
```

To step back - we have Flux running inside Kubernetes, and now a dummy Kubernetes
running inside Flux. The four nodes in the cluster are registered, with one control
plane and three agents. We will want to test this next on a Flux cluster that doesn't
have an external Kubernetes wrapper, and then likely run something that is more like an HPC job.
Once that is working, we will want to slowly step back and figure out the steps necessary
to run this entirely rootless. I've gotten [k3s working with rootless](https://github.com/converged-computing/operator-experiments/tree/main/google/rootless-kubernetes/k3s) 
but it has a few more steps! You can then run kubernetes commands as you please! 

I had a death wish so I installed the operator again...

```bash
cd /tmp
wget https://raw.githubusercontent.com/flux-framework/flux-operator/main/examples/dist/flux-operator.yaml
kubectl --kubeconfig=/workflow/kubeconfig.yaml apply -f flux-operator.yaml 
```

Is the operator pod running?

```bash
root@flux-sample-0:/tmp# kubectl --kubeconfig=/workflow/kubeconfig.yaml get -n operator-system pods
NAME                                           READY   STATUS    RESTARTS   AGE
operator-controller-manager-658b4c6787-7stwv   2/2     Running   0          46s
```

!!! Let's try applying a Minicluster with hello world...

```bash
$ wget https://raw.githubusercontent.com/flux-framework/flux-operator/main/examples/tests/hello-world/minicluster.yaml
```
```bash
$ kubectl --kubeconfig=/workflow/kubeconfig.yaml apply -f minicluster.yaml 
minicluster.flux-framework.org/flux-sample created
```

Is it running?

```bash
root@flux-sample-0:/tmp# kubectl --kubeconfig=/workflow/kubeconfig.yaml get pods -n flux-operator
NAME                       READY   STATUS              RESTARTS   AGE
my-echo-74dc6c4f7b-snkt9   1/1     Running             0          6m1s
flux-sample-0-x45q6        0/1     ContainerCreating   0          16s
flux-sample-2-7wltv        0/1     ContainerCreating   0          16s
flux-sample-3-l668b        0/1     ContainerCreating   0          16s
flux-sample-1-8t6fd        0/1     ContainerCreating   0          16s
```

WHAT IS HAPPENING! 🤣️ I'm going to stop here because I'm afraid of it actually pulling this
second layer of (rather large) container with Flux, already in a container, and we will
embark on this second layer of the onion once we have addressed the issues above and
tested in different environments. Let's clean up before we do something that we will regret!

```bash
$ kubectl delete -f minicluster.yaml
```
