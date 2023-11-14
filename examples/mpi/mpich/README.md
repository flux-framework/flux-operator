# Mpich Example

You should be able to create a MiniKube cluster, install the operator with creating the namespace:

```bash
$ minikube start
$ kubectl apply -f ../../dist/flux-operator.yaml
```

You might want to pre-pull the container:

```bash
$ minikube ssh docker pull ghcr.io/rse-ops/mpich:tag-mamba
```

And then create the MiniCluster:

```bash
$ kubectl create -f minicluster.yaml
```

And watch the example run!

```bash
$ kubectl logs flux-sample-0-5gjqt -f
```

A successful run will show four MPI ranks...

```console
broker.info[0]: rc1.0: running /etc/flux/rc1.d/02-cron
broker.info[0]: rc1.0: /etc/flux/rc1 Exited (rc=0) 0.5s
broker.info[0]: rc1-success: init->quorum 0.543544s
broker.info[0]: online: flux-sample-0 (ranks 0)
broker.info[0]: online: flux-sample-[0-3] (ranks 0-3)
broker.info[0]: quorum-full: quorum->run 0.369278s
Hello, world!  I am 1 of 4(Open MPI v4.0.3, package: Debian OpenMPI, ident: 4.0.3, repo rev: v4.0.3, Mar 03, 2020, 87)
Hello, world!  I am 0 of 4(Open MPI v4.0.3, package: Debian OpenMPI, ident: 4.0.3, repo rev: v4.0.3, Mar 03, 2020, 87)
Hello, world!  I am 2 of 4(Open MPI v4.0.3, package: Debian OpenMPI, ident: 4.0.3, repo rev: v4.0.3, Mar 03, 2020, 87)
Hello, world!  I am 3 of 4(Open MPI v4.0.3, package: Debian OpenMPI, ident: 4.0.3, repo rev: v4.0.3, Mar 03, 2020, 87)
broker.info[0]: rc2.0: flux submit -N 4 -n 4 --quiet --watch ./hello_cxx Exited (rc=0) 0.8s
broker.info[0]: rc2-success: run->cleanup 0.843814s
broker.info[0]: cleanup.0: flux queue stop --quiet --all --nocheckpoint Exited (rc=0) 0.1s
broker.info[0]: cleanup.1: flux cancel --user=all --quiet --states RUN Exited (rc=0) 0.1s
broker.info[0]: cleanup.2: flux queue idle --quiet Exited (rc=0) 0.1s
broker.info[0]: cleanup-success: cleanup->shutdown 0.320065s
broker.info[0]: children-complete: shutdown->finalize 61.2525ms
broker.info[0]: rc3.0: running /etc/flux/rc3.d/01-sched-fluxion
broker.info[0]: rc3.0: /etc/flux/rc3 Exited (rc=0) 0.3s
broker.info[0]: rc3-success: finalize->goodbye 0.310701s
broker.info[0]: goodbye: goodbye->exit 0.037999ms
```

And the job will be completed.

```bash
kubectl get -n flux-operator pods
```
```console
NAME                  READY   STATUS      RESTARTS   AGE
flux-sample-0-5gjqt   0/1     Completed   0          2m40s
flux-sample-1-j4zlc   0/1     Completed   0          2m40s
flux-sample-2-wdzz7   0/1     Completed   0          2m40s
flux-sample-3-vp8rx   0/1     Completed   0          2m40s
```