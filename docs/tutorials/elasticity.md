# Elasticity

> This functionality requires [Kubernetes 1.27 and later](https://github.com/kubernetes/enhancements/tree/master/keps/sig-apps/3715-elastic-indexed-job#motivation).

This will combine what we've learned about [scaling](scaling.md) with elasticity, meaning that the application
is able to control scaling up and down. You can read [the scaling tutorial](scaling.md) to get a sense of how that
works. This will add the ability for the application to communicate with the API server to request scaling up
or down.

## Basic Example

> An application interacting with the custom resource API to scale up and down

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/elasticity/basic/minicluster.yaml)**

To run this example, you'll want to first cd to the example directory:


```bash
$ cd examples/elasticity/basic
```

Create a kind cluster with Kubernetes version 1.27.0

```bash
$ kind create cluster --config ./kind-config.yaml
```

Since we need custom permissions - to make requests from inside the cluster - we need to apply
extra RBAC roles:

```bash
$ kubectl apply -f ../rbac.yaml
```

Note that this is not added to the operator proper because not everyone needs this level of permission,
and we should be conservative. Likely the above can be better streamlined - I was in "get it working" mode
when I wrote it! We could also likely create a more scoped service account as opposed to adding to
the flux-operator's. Then install the operator, and create the MiniCluster:

```bash
$ kubectl apply -f ../../dist/flux-operator.yaml
$ kubectl apply -f ./minicluster.yaml
```

At this point, wait until the containers go from creating to running.
Note that it does take about a minute to pull.

```bash
$ kubectl get -n flux-operator pods
NAME                  READY   STATUS    RESTARTS   AGE
flux-sample-0-4wmmp   1/1     Running   0          6m50s
flux-sample-1-mjj7b   1/1     Running   0          6m50s
```

Then you can watch logs to see the demo!

```bash
$ kubectl logs -n flux-operator flux-sample-0-zkfp8 -f
```

<details>

<summary>Gopherlocks and the (not) three bears demo</summary>

```console
Elasticity anyone? Story anyone? Bueller? BUELLER?!
BOOTING UP THE STORY! 🥾️
Hello there! 👋️ I'm Gopherlocks! 👱️
broker.err[0]: flux-sample-1 (rank 1) reconnected after crash, dropping old connection state
broker.err[0]: accepting connection from flux-sample-1 (rank 1) status full
broker.info[0]: online: flux-sample-0 (ranks 0)
broker.info[0]: online: flux-sample-[0-1] (ranks 0-1)
broker.info[0]: online: flux-sample-[0-3] (ranks 0-3)

Oh my, am I in a container?! 👱️

Let's take a look around... who else is here?
 👀️

Hello pods... who is out there?
10.244.0.21     flux-operator   flux-sample-0-kjfpl
10.244.0.20     flux-operator   flux-sample-1-74w7r

I see it over there! I'm running in a job called flux-sample. 🌀️

<flux.core.handle.Flux object at 0x7f2d0c97d7e0>
Oh hi Flux, I guess you are here too. 👋️

Please don't lay a stinky one, I know how you job managers get! 💩️

So hmm. I think I'm running in a Flux Operator MiniCluster. 😎️
broker.err[0]: flux-sample-1 (rank 1) reconnected after crash, dropping old connection state
broker.err[0]: accepting connection from flux-sample-1 (rank 1) status full
broker.info[0]: online: flux-sample-[0,2-3] (ranks 0,2-3)

Just a guess! 🤷️

At least it is not three bears, har har har. 🐻️ 🐻️ 🐻️
broker.info[0]: online: flux-sample-[0-3] (ranks 0-3)

I wonder if I can find the spec for the cluster I AM IN RIGHT NOW...

Oh, I think I found it!
broker.err[0]: flux-sample-2 (rank 2) transitioning to LOST due to EHOSTUNREACH error on send
broker.err[0]: flux-sample-3 (rank 3) transitioning to LOST due to EHOSTUNREACH error on send
broker.info[0]: online: flux-sample-[0-1,3] (ranks 0-1,3)
broker.info[0]: online: flux-sample-[0-1] (ranks 0-1)

⭐️    cleanup: False
⭐️    command: python /data/three_bears.py
⭐️   commands: {'pre': 'pip install kubernetes', 'runFluxAsRoot': False}
⭐️   fluxUser: {'name': 'fluxuser', 'uid': 1000}
⭐️      image: ghcr.io/rse-ops/singularity:tag-mamba
⭐️   launcher: True
⭐️ pullAlways: False
⭐️    volumes: {'data': {'path': '/data', 'readOnly': False}}
⭐️ workingDir: /data
⭐️ deadlineSeconds: 31500000
⭐️ interactive: False
⭐️    maxSize: 10
⭐️       size: 2
⭐️      tasks: 1
⭐️    volumes: {'data': {'capacity': '5Gi', 'delete': True, 'path': '/tmp/workflow', 'secretNamespace': 'default', 'storageClass': 'hostpath'}}

Oh my, is it a bit, tight in here? A size 2?!

Let's see what I can do about that...

Did that work? Hello out there... do we have more friends? 🍤️

Hello pods... who is out there?
10.244.0.21     flux-operator   flux-sample-0-kjfpl
10.244.0.20     flux-operator   flux-sample-1-74w7r
10.244.0.25     flux-operator   flux-sample-2-jmhb9
10.244.0.26     flux-operator   flux-sample-3-tpg9j

Oh my, we have FOUR friends!! I'm so happy! 😹️

But actually I wanted to play some Mario Kart but I only have 4 controllers... 🕹️
broker.info[0]: online: flux-sample-[0-3] (ranks 0-3)

Sorry one of you has to leave!!... 😭️

I know I'm a terrible person. 👿️

I feel so bad. How many do we have now?

** DRAMATIC PAUSE FOR STORY ** but actually to wait for pod to terminate :)
broker.err[0]: flux-sample-3 (rank 3) transitioning to LOST due to EHOSTUNREACH error on send
broker.info[0]: online: flux-sample-[0-2] (ranks 0-2)

Hello pods... who is out there?
10.244.0.21     flux-operator   flux-sample-0-kjfpl
10.244.0.20     flux-operator   flux-sample-1-74w7r
10.244.0.25     flux-operator   flux-sample-2-jmhb9

NICE!! TIME TO DESTROY YOU IN MARIO KART! 💪️

"Player select: Peach." 🍑️

Hey now, do not judge! 😜️
```

</details>

You can watch the asciinema too:

<script async id="asciicast-585802" src="https://asciinema.org/a/585802.js"></script>

And that's it! Note that there are other ways to scale using the Horizontal Auto Scaler that I'm going
to also look into, although at brief glance it seems like this would require the application to serve
an endpoint. Make sure to cleanup when you finish:

```bash
$ kind delete cluster
```

## Horizontal Autoscaler Examples

> Using the horizontal pod autoscaler (HPA) API to scale instead

The rbac permissions required above might be a bit much, and this was suggested as an alternative approach.
For this approach we are going to use the [HPA + Scale sub-resource](https://book.kubebuilder.io/reference/generating-crd.html#scale) in our custom resource definition to scale resources automatically,
and our application running in the MiniCluster can then provide custom metrics based on which we set up the scaling. 

The [APIs for v1 and v2 of the autoscaler are subtly different](https://www.pulumi.com/registry/packages/kubernetes/api-docs/autoscaling/v1/horizontalpodautoscaler/).
In the version 1 example below, we scale based on size. In version 2 we use a different API and then add [multiple metrics](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/#autoscaling-on-multiple-metrics-and-custom-metrics). Note that we are interested in this technique because the requests can [come from an external service](https://cloud.google.com/kubernetes-engine/docs/concepts/horizontalpodautoscaler), meaning we could have a single service (paired with application logic, optionally) to handle scaling both node and MiniClusters, and coordinating the two.

### How does it work?

I think what happens is that our operator provides the autoscaler with a field (the size) and selector for pods (the `hpa-selector`) and then:

1. It pings the API endpoint for our CRD to get the current selector and size
2. It retrieves the pods based on the selector
3. It compares the actual vs. desired metrics
4. If scaling is needed, it applies a change to the field indicated, directly to the CRD
5. The change request comes into our operator!

The above could be wrong - this is my best effort guess after an afternoon of poking around.
This means that, to the Flux Operator, a request to scale coming from a user applying an updated
YAML vs. the autoscaler looks the same (and we can use the same functions). It also means that
if an autoscaler isn't installed to the cluster, we have those fields but nothing happens /
no action is taken. This means that:

### What do you need?

You'll need to ensure that:

- Your minicluster.yaml has resources defined for CPU (or the metrics you are interested in)
- You've set at maxSize in your minicluster.yaml to give Flux a heads up
- Your cluster will need a metrics server
- If it's a local metrics server, you'll want to disable ssl ([see our example](https://github.com/flux-framework/flux-operator/blob/main/examples/elasticity/horizontal-autoscaler/metrics-server.yaml))

Details for how to inspect the above (and sample files) are provided below.


### Horizontal Autoscaler (v2) Example

The version 2 API is more flexible than version 1 in allowing custom metrics. This means we can use a [prometheus-flux](https://github.com/converged-computing/prometheus-flux)
exporter running inside of an instance to interact with it. This small set of tutorials will show setting a basic autoscaling example
based on CPU, and then one based on custom metrics.

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/elasticity/horizontal-autoscaler/v2-cpu/minicluster.yaml)**

#### Setup

This is the setup for testing v2 with a CPU metric:

```bash
$ cd examples/elasticity/horizontal-autoscaler/v2-cpu
```

Create a kind cluster with Kubernetes version 1.27.0

```bash
$ kind create cluster --config ./kind-config.yaml
```

Install the operator:

```bash
$ kubectl apply -f ../../../dist/flux-operator-dev.yaml
```

And then create a very simply interactive cluster (it's small, 2 pods, but importantly has a maxsize of 10).
Note that we are going to limit the HPA to a size of 4, because we assume you are running on an average desktop computer.

```bash
$ kubectl apply -f ./minicluster.yaml
```

You'll need to wait for the container to pull (status `ContainerCreating` to `Running`).
At this point, wait until the containers go from creating to running.

```bash
$ kubectl get pods
NAME                  READY   STATUS    RESTARTS   AGE
flux-sample-0-4wmmp   1/1     Running   0          6m50s
flux-sample-1-mjj7b   1/1     Running   0          6m50s
```

Look at the scale endpoint of the MiniCluster with `kubectl` directly! Remember that we haven't installed a horizontal auto-scaler yet:

```bash
$ kubectl get --raw /apis/flux-framework.org/v1alpha2/namespaces/default/miniclusters/flux-sample/scale | jq
```
```console
{
  "kind": "Scale",
  "apiVersion": "autoscaling/v1",
  "metadata": {
    "name": "flux-sample",
    "namespace": "default",
    "uid": "3e5c568c-0367-4e76-bf7a-6e3215afdb1e",
    "resourceVersion": "603",
    "creationTimestamp": "2023-10-16T00:47:35Z"
  },
  "spec": {
    "replicas": 2
  },
  "status": {
    "replicas": 0,
    "selector": "hpa-selector=flux-sample"
  }
}
```

The above knows the selector to use to get pods (and look at current resource usage).
The output above is also telling us the `autoscaler/v1` is being used, which I used to think
means I could not use autoscaler/v2, but they seem to work OK (note that autoscaler/v2 is
installed to my current cluster with Kubernetes 1.27).

Before we deploy any autoscaler, we need a metrics server! This doesn't come out of the box with kind so
we install it:

```bash
$ kubectl apply -f metrics-server.yaml
```

I found this suggestion [here](https://gist.github.com/sanketsudake/a089e691286bf2189bfedf295222bd43). Ensure
it's running:

```bash
$ kubectl get deploy,svc -n kube-system | egrep metrics-server
```

#### Autoscaler with CPU

This first autoscaler will work based on CPU. We can create it as follows:

```bash
$ kubectl apply -f hpa-cpu.yaml
```
```console
horizontalpodautoscaler.autoscaling/flux-sample-hpa created
```

Remember that when you first created your cluster, your size was two, and we had two?

```bash
$ kubectl get pods
NAME                  READY   STATUS    RESTARTS   AGE
flux-sample-0-4wmmp   1/1     Running   0          6m50s
flux-sample-1-mjj7b   1/1     Running   0          6m50s
```

If you watch your pods (and your autoscaler and your endpoint) you'll
see first that the resource usage changes (just by way of Flux starting):
And to get it to change more, try shelling into your broker leader pod, connecting
to the broker, and issuing commands:


```bash
kubectl exec -it flux-sample-0-p85cj bash

# This is written by the operator to make it easy to add the view
. /mnt/flux/flux-view.sh
flux proxy local:///mnt/flux/view/run/flux/local bash
$ openssl speed -multi 4
```

You'll see it change (with updates between 15 seconds and 1.5 minutes!):

```
$ kubectl get hpa -w
NAME              REFERENCE                 TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
flux-sample-hpa   MiniCluster/flux-sample   0%/2%     2         4         2          3m45s
flux-sample-hpa   MiniCluster/flux-sample   0%/2%     2         4         2          3m45s
flux-sample-hpa   MiniCluster/flux-sample   34%/2%    2         4         2          4m
flux-sample-hpa   MiniCluster/flux-sample   49%/2%    2         4         4          4m15s
```

See the [autoscaler/v1](#creating-the-v1-autoscaler) example for more detail about outputs. They have
a slightly different design, but result in the same output to the terminal.
When you are done demo-ing the CPU autoscaler, you can clean it up:

```bash
$ kubectl delete -f hpa-cpu.yaml
```

#### Autoscaler with Custom Metrics

For this tutorial, it requires a little extra work to setup the api service, so we provide
the full details in the tutorial directory alongside the MiniCluster YAML file.

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/elasticity/horizontal-autoscaler/v2-custom-metric/minicluster.yaml)**

And when you are done either of the tutorials above, don't forget to clean up.

```bash
$ kind delete cluster
```

### Horizontal Autoscaler (v1) Example

This will show using the version 1 API.

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/elasticity/horizontal-autoscaler/v1/minicluster.yaml)**

You should have already completed [setup](#setup) for this example.
To run this example, you'll want to first cd to the example directory:

```bash
$ cd examples/elasticity/horizontal-autoscaler/v1
```

#### Creating the v1 autoscaler

Now create the autoscaler!

```bash
$ kubectl apply -f hpa-v1.yaml
```
```console
horizontalpodautoscaler.autoscaling/flux-sample-hpa created
```

Create the MiniCluster:

```bash
kubectl apply -f ./minicluster.yaml
```
Remember that when you first created your cluster, your size was two, and we had two?

```bash
$ kubectl get pods
NAME                  READY   STATUS    RESTARTS   AGE
flux-sample-0-4wmmp   1/1     Running   0          6m50s
flux-sample-1-mjj7b   1/1     Running   0          6m50s
```

If you watch your pods (and your autoscaler and your endpoint) you'll
see first that the resource usage changes (just by way of Flux starting):
And to get it to change more, try shelling into your broker leader pod, connecting
to the broker, and issuing commands:

```bash
kubectl exec -it flux-sample-0-p85cj bash
. /mnt/flux/flux-view.sh 
flux proxy ${fluxsocket} bash
openssl speed -multi 4
```

You'll see it change (with updates between 15 seconds and 1.5 minutes!):

```bash
$ kubectl get hpa -w
NAME              REFERENCE                 TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
flux-sample-hpa   MiniCluster/flux-sample   0%/2%     2         4         2          33m
flux-sample-hpa   MiniCluster/flux-sample   0%/2%     2         4         2          34m
flux-sample-hpa   MiniCluster/flux-sample   3%/2%     2         4         2          34m
flux-sample-hpa   MiniCluster/flux-sample   0%/2%     2         4         3          34m
```

So basically, we are seeing the cluster increase in size accordingly! Note that another way to achieve this is just to set `interactive: true` to
the minicluster.yaml. 

```bash
kubectl get pods
NAME                  READY   STATUS    RESTARTS   AGE
flux-sample-0-blgsw   1/1     Running   0          56s
flux-sample-1-k4m6s   1/1     Running   0          56s
flux-sample-2-z64n9   1/1     Running   0          17s
```

Also note that the autoscaler age is 33 minutes, but the job pods under a minute!
This is because we can install the autoscaler once, and then use it for the same
namespace and job name, even if you wind up deleting and re-creating. That's pretty
cool because it means you can create different HPAs for different job types (organized by name).
Depending on how you stress the pods, it can easily go up to the max:

```bash
kubectl get --raw /apis/flux-framework.org/v1alpha2/namespaces/default/miniclusters/flux-sample/scale | jq
```
```console
{
  "kind": "Scale",
  "apiVersion": "autoscaling/v1",
  "metadata": {
    "name": "flux-sample",
    "namespace": "default",
    "uid": "a5455ad0-ed53-4408-b8e1-67478516ba14",
    "resourceVersion": "3247",
    "creationTimestamp": "2023-10-16T01:05:04Z"
  },
  "spec": {
    "replicas": 4
  },
  "status": {
    "replicas": 4,
    "selector": "hpa-selector=flux-sample"
  }
}
```

And you see the autoscaler reflect that. Since we set the maxSize to 4, it won't ever go above that.

```bash
$ kubectl get -n flux-operator hpa -w
NAME              REFERENCE                 TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
flux-sample-hpa   MiniCluster/flux-sample   0%/2%     2         4         2          49m
flux-sample-hpa   MiniCluster/flux-sample   0%/2%     2         4         2          49m
flux-sample-hpa   MiniCluster/flux-sample   0%/2%     2         4         2          50m
flux-sample-hpa   MiniCluster/flux-sample   0%/2%     2         4         4          50m
flux-sample-hpa   MiniCluster/flux-sample   9%/2%     2         4         4          50m
```

Also note that if the load goes down (e.g., to zero again) the cluster will scale down.
I saw this when I set `interactive: true` and manually shelled in to increase load,
and then exited and let it sit. 

And that's it! This is really cool because it's an early step toward a totally automated,
self-scaling workflow. The next step is combining that with logic to scale the Kubernetes
cluster itself. When you are done:

```bash
$ kind delete cluster
```