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

And then you can see the pods:

```bash
kubectl get -n flux-operator pods
NAME                                       READY   STATUS    RESTARTS   AGE
minicluster-flux-sample-broker-0-0-fsv4p   1/1     Running   0          12m
minicluster-flux-sample-worker-0-0-tjg69   1/1     Running   0          12m
minicluster-flux-sample-worker-0-1-hwmbb   1/1     Running   0          12m
minicluster-flux-sample-worker-0-2-q246d   1/1     Running   0          12m
```

And that the broker is running and ready for dootie! I mean duty! :D

```
✨ Curve certificate generated by helper pod
#   ****  Generated on 2023-04-26 22:54:42 by CZMQ  ****
#   ZeroMQ CURVE **Secret** Certificate
#   DO NOT PROVIDE THIS FILE TO OTHER USERS nor change its permissions.
    
metadata
    name = "flux-cert-generator"
    keygen.hostname = "flux-sample-0"
curve
    public-key = "WCU!!@2t4>.:Khqg%bNFN#.lf*Eh)vlbVx^@s-is"
    secret-key = "NXvRbbIU9KZ?&64Y#)09@zxSw20VT.FfH(J/sJ1-"

🌀  flux start -o --config /etc/flux/config -Scron.directory=/etc/flux/system/cron.d   -Stbon.fanout=256   -Srundir=/run/flux    -Sstatedir=/var/lib/flux   -Slocal-uri=local:///run/flux/local     -Slog-stderr-level=6    -Slog-stderr-mode=local
broker.info[2]: start: none->join 0.434417ms
broker.info[2]: parent-ready: join->init 0.442597s
broker.info[2]: configuration updated
broker.info[2]: rc1.0: running /etc/flux/rc1.d/01-sched-fluxion
broker.info[2]: rc1.0: running /etc/flux/rc1.d/02-cron
broker.info[2]: rc1.0: /etc/flux/rc1 Exited (rc=0) 0.1s
broker.info[2]: rc1-success: init->quorum 0.12074s
broker.info[2]: quorum-full: quorum->run 0.201772s
```