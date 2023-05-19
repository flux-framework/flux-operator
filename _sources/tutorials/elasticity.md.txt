# Elasticity

> This functionality requires [Kubernetes 1.27 and later](https://github.com/kubernetes/enhancements/tree/master/keps/sig-apps/3715-elastic-indexed-job#motivation).

This will combine what we've learned about [scaling](scaling.md) with elasticity, meaning that the application
is able to control scaling up and down. You can read [the scaling tutorial](scaling.md) to get a sense of how that
works. This will add the ability for the application to communicate with the API server to request scaling up
or down.

## Basic Example

> Starting with a cluster at a maximum size and scaling it up and down

 **[Tutorial File](https://github.com/flux-framework/flux-operator/blob/main/examples/elasticity/basic/minicluster.yaml)**

To run this example, you'll want to first cd to the example directory:


```bash
$ cd examples/elasticity/basic
```

Create a kind cluster with Kubernetes version 1.27.0

```bash
$ kind create cluster --config ./kind-config.yaml
```

Create the flux-operator namespace:

```bash
$ kubectl create namespace flux-operator
```

Since we need custom permissions - to make requests from inside the cluster - we need to apply
extra RBAC roles:

```bash
$ kubectl apply -f ../rbac.yaml
```

Note that this is not added to the operator proper because not everyone needs this level of permission,
and we should be conservative. Likely the above can be better streamlined - I was in "get it working" mode
when I wrote it! We could also likely create a more scoped service account as opposed to adding to
the flux-operators. Then install the operator, and create the MiniCluster:

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