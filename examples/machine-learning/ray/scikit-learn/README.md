# Ray with scikit-learn and Flux!

This will test running [Ray and scikit-learn](https://docs.ray.io/en/latest/ray-more-libs/joblib.html) with Flux.
Our goal is to do this simply, and eventually extend this to a more complex hierarchy of jobs.
Note that I found most of the needed logic [here](https://docs.ray.io/en/latest/cluster/vms/user-guides/launching-clusters/on-premises.html#on-prem).

## Usage

### Create Cluster

First, let's create a kind cluster.

```bash
$ kind create cluster --config ../../../kind-config.yaml
```

And then install the operator, and apply the MiniCluster YAML here.

```bash
$ kubectl apply -f ../../../dist/flux-operator.yaml
$ kubectl apply -f ./minicluster.yaml
```

This will use a base container with Ray and Scikit-learn, and mount
the current directory (in the pods as `/tmp/workflow`). We run the [start.sh](start.sh)
to start the cluster. You can watch doing the following:

```bash
$ kubectl logs flux-sample-0-7tx7s -f
```

<details>

<summary>Ray Expected Output</summary>

```console
2023-05-10 22:52:44,254 INFO scripts.py:892 -- Local node IP: 10.244.0.57
2023-05-10 22:52:47,567 SUCC scripts.py:904 -- --------------------
2023-05-10 22:52:47,567 SUCC scripts.py:905 -- Ray runtime started.
2023-05-10 22:52:47,567 SUCC scripts.py:906 -- --------------------
2023-05-10 22:52:47,567 INFO scripts.py:908 -- To terminate the Ray runtime, run
2023-05-10 22:52:47,567 INFO scripts.py:909 --   ray stop
2023-05-10 22:52:44,300 INFO usage_lib.py:372 -- Usage stats collection is disabled.
2023-05-10 22:52:44,300 INFO scripts.py:710 -- Local node IP: flux-sample-0.flux-service.default.svc.cluster.local
2023-05-10 22:52:47,605 SUCC scripts.py:747 -- --------------------
2023-05-10 22:52:47,605 SUCC scripts.py:748 -- Ray runtime started.
2023-05-10 22:52:47,605 SUCC scripts.py:749 -- --------------------
2023-05-10 22:52:47,605 INFO scripts.py:751 -- Next steps
2023-05-10 22:52:47,605 INFO scripts.py:754 -- To add another node to this Ray cluster, run
2023-05-10 22:52:47,605 INFO scripts.py:757 --   ray start --address='flux-sample-0.flux-service.default.svc.cluster.local:6379'
2023-05-10 22:52:47,605 INFO scripts.py:766 -- To connect to this Ray cluster:
2023-05-10 22:52:47,605 INFO scripts.py:768 -- import ray
2023-05-10 22:52:47,606 INFO scripts.py:769 -- ray.init(_node_ip_address='flux-sample-0.flux-service.default.svc.cluster.local')
2023-05-10 22:52:47,606 INFO scripts.py:781 -- To submit a Ray job using the Ray Jobs CLI:
2023-05-10 22:52:47,606 INFO scripts.py:782 --   RAY_ADDRESS='http://127.0.0.1:8265' ray job submit --working-dir . -- python my_script.py
2023-05-10 22:52:47,622 INFO scripts.py:791 -- See https://docs.ray.io/en/latest/cluster/running-applications/job-submission/index.html 
2023-05-10 22:52:47,622 INFO scripts.py:795 -- for more information on submitting Ray jobs to the Ray cluster.
2023-05-10 22:52:47,622 INFO scripts.py:800 -- To terminate the Ray runtime, run
2023-05-10 22:52:47,622 INFO scripts.py:801 --   ray stop
2023-05-10 22:52:47,622 INFO scripts.py:804 -- To view the status of the cluster, use
2023-05-10 22:52:47,623 INFO scripts.py:805 --   ray status
2023-05-10 22:52:47,623 INFO scripts.py:809 -- To monitor and debug Ray, view the dashboard at 
2023-05-10 22:52:47,623 INFO scripts.py:810 --   127.0.0.1:8265
2023-05-10 22:52:47,623 INFO scripts.py:817 -- If connection to the dashboard fails, check your firewall settings and network configuration.
2023-05-10 22:52:44,270 INFO scripts.py:892 -- Local node IP: 10.244.0.59
2023-05-10 22:52:47,654 SUCC scripts.py:904 -- --------------------
2023-05-10 22:52:47,654 SUCC scripts.py:905 -- Ray runtime started.
2023-05-10 22:52:47,655 SUCC scripts.py:906 -- --------------------
2023-05-10 22:52:47,655 INFO scripts.py:908 -- To terminate the Ray runtime, run
2023-05-10 22:52:47,655 INFO scripts.py:909 --   ray stop
[2023-05-10 22:52:47,666 I 108 108] global_state_accessor.cc:356: This node has an IP address of 10.244.0.56, but we cannot find a local Raylet with the same address. This can happen when you connect to the Ray cluster with a different IP address or when connecting to a container.
2023-05-10 22:52:44,305 INFO scripts.py:892 -- Local node IP: 10.244.0.56
2023-05-10 22:52:47,718 SUCC scripts.py:904 -- --------------------
2023-05-10 22:52:47,719 SUCC scripts.py:905 -- Ray runtime started.
2023-05-10 22:52:47,719 SUCC scripts.py:906 -- --------------------
2023-05-10 22:52:47,719 INFO scripts.py:908 -- To terminate the Ray runtime, run
2023-05-10 22:52:47,719 INFO scripts.py:909 --   ray stop
```

</details>

### Run Workflow

Let's shell into the broker pod and connect to the broker flux instance:

```bash
kubectl exec -it flux-sample-0-jlsp6 bash
source /mnt/flux/flux-view.sh
flux proxy $fluxsocket bash
```

At this point, since this main pod is running the show, we can run our Python example that
uses ray:

```bash
$ python3 ray_tune.py
```

Note that you could also give this to the broker as the command directly (and no need for bash above):

```bash
kubectl exec -it flux-sample-0-jlsp6 bash
source /mnt/flux/flux-view.sh
flux proxy $fluxsocket python3 ray_tune.py
```

For either approach, you should see the training logs across the cluster:

```console
(PoolActor pid=564) [CV 1/5; 125/300] END C=0.00011721022975334806, class_weight=balanced, gamma=23.95026619987481, tol=0.0017433288221999873;, score=0.100 total time=   2.2s [repeated 5x across cluster]
(PoolActor pid=564) [CV 5/5; 129/300] START C=1000000.0, class_weight=balanced, gamma=0.14873521072935117, tol=0.03039195382313198 [repeated 9x across cluster]
(PoolActor pid=565) [CV 5/5; 126/300] START C=489.3900918477499, class_weight=None, gamma=1.8873918221350996, tol=0.014873521072935119 [repeated 10x across cluster]
(PoolActor pid=559) [CV 3/5; 130/300] END C=72.78953843983146, class_weight=None, gamma=1082.6367338740563, tol=0.0008531678524172806;, score=0.103 total time=   2.1s [repeated 9x across cluster]
(PoolActor pid=564) [CV 5/5; 129/300] END C=1000000.0, class_weight=balanced, gamma=0.14873521072935117, tol=0.03039195382313198;, score=0.103 total time=   2.2s [repeated 10x across cluster]
(PoolActor pid=563) [CV 1/5; 136/300] START C=489.3900918477499, class_weight=balanced, gamma=28072162.039411698, tol=0.003562247890262444 [repeated 12x across cluster]
(PoolActor pid=564) [CV 1/5; 133/300] START C=0.2395026619987486, class_weight=None, gamma=100000000.0, tol=0.003562247890262444 [repeated 9x across cluster]
(PoolActor pid=565) [CV 1/5; 130/300] END C=72.78953843983146, class_weight=None, gamma=1082.6367338740563, tol=0.0008531678524172806;, score=0.100 total time=   2.8s [repeated 10x across cluster]
(PoolActor pid=559) [CV 2/5; 135/300] END C=2.592943797404667e-06, class_weight=balanced, gamma=28072162.039411698, tol=0.07880462815669913;, score=0.100 total time=   2.3s [repeated 8x across cluster]
(PoolActor pid=563) [CV 2/5; 139/300] START C=1.7433288221999873e-05, class_weight=balanced, gamma=0.5298316906283702, tol=0.06210169418915616 [repeated 10x across cluster]
...
```

And that's it! If you create a more substantial example using Ray please [let us know](https://github.com/flux-framework/flux-operator/issues).

### Cleanup

When you are done, clean up:

```bash
$ kubectl delete -f minicluster.yaml
```

Make sure to clean up your shared tmpdir!

```bash
$ rm *.out
$ sudo rm -rf ./tmp/*
```