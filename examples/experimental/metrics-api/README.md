# Metrics API

This example will show running the [Flux Metrics API](https://github.com/converged-computing/flux-metrics-api) alongside the main index 0 broker to serve custom metrics.
The idea is that Flux is running in a container in a pod, but that it has a sidecar (with a mount of the same flux install and access to the same flux socket!) so that:

1. The main Flux container (where Flux is still added on the fly) is running some application of interest.
2. The sidecar container runs the Flux Metrics API, which (exposed via ingress) gives any cluster user (you!) programmatic access to get metrics about the queue.


We use selectors for the index and job name to ensure we are targeting the service only to the lead broker. This could target sub-instances as well (meaning other pods,
yes you could totally have some weird setup with a metrics API at several levels of the hierarchy) but this should be good for now :)
 

## Cluster

First, create a test cluster.

```bash
kind create cluster --config ./kind-config.yaml
```

Note that I struggled with creating ingress for a while because Kind needs additions for it to work (in that file).
Install the flux operator from here.

## MiniCluster

Create the minicluster.

```bash
kubectl apply -f minicluster.yaml
```

To shell into the metrics container:

```bash
kubectl exec -it flux-sample-0-zp5q9 -c metrics bash
```

### Ingress

We then want to create a service to access the registry (make sure to close the port forward):

```bash
kubectl apply -f metrics/ingress.yaml
```
```console
$ kubectl describe ingress
Name:             oras-ingress
Labels:           <none>
Namespace:        default
Address:          
Ingress Class:    <none>
Default backend:  <default>
Rules:
  Host        Path  Backends
  ----        ----  --------
  localhost   
              /   oras-service:5000 (10.244.0.12:5000)
Annotations:  <none>
Events:       <none>
```

We are then going to apply [ingress-nginx](https://kind.sigs.k8s.io/docs/user/ingress/#ingress-nginx).

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
```

### Basic Test

Then we can interact with the API! Here is to list metrics.

```bash
$ curl -k localhost/apis/custom.metrics.k8s.io/v1beta2 | jq
```
```console
{"kind":"APIResourceList","apiVersion":"v1beta2","groupVersion":"custom.metrics.k8s.io/v1beta2","resources":[{"name":"node_cores_free_count","singularName":"node_cores_free_count","namespaced":true,"kind":"MetricValueList","verbs":["get"]},{"name":"node_cores_up_count","singularName":"node_cores_up_count","namespaced":true,"kind":"MetricValueList","verbs":["get"]},{"name":"node_free_count","singularName":"node_free_count","namespaced":true,"kind":"MetricValueList","verbs":["get"]},{"name":"node_up_count","singularName":"node_up_count","namespaced":true,"kind":"MetricValueList","verbs":["get"]},{"name":"job_queue_state_new_count","singularName":"job_queue_state_new_count","namespaced":true,"kind":"MetricValueList","verbs":["get"]},{"name":"job_queue_state_depend_count","singularName":"job_queue_state_depend_count","namespaced":true,"kind":"MetricValueList","verbs":["get"]},{"name":"job_queue_state_priority_count","singularName":"job_queue_state_priority_count","namespaced":true,"kind":"MetricValueList","verbs":["get"]},{"name":"job_queue_state_sched_count","singularName":"job_queue_state_sched_count","namespaced":true,"kind":"MetricValueList","verbs":["get"]},{"name":"job_queue_state_run_count","singularName":"job_queue_state_run_count","namespaced":true,"kind":"MetricValueList","verbs":["get"]},{"name":"job_queue_state_cleanup_count","singularName":"job_queue_state_cleanup_count",{
  "kind": "APIResourceList",
  "apiVersion": "v1beta2",
  "groupVersion": "custom.metrics.k8s.io/v1beta2",
  "resources": [
    {
      "name": "node_cores_free_count",
      "singularName": "node_cores_free_count",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    },
    {
      "name": "node_cores_up_count",
      "singularName": "node_cores_up_count",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    },
    {
      "name": "node_free_count",
      "singularName": "node_free_count",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    },
    {
      "name": "node_up_count",
      "singularName": "node_up_count",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    },
    {
      "name": "job_queue_state_new_count",
      "singularName": "job_queue_state_new_count",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    },
    {
      "name": "job_queue_state_depend_count",
      "singularName": "job_queue_state_depend_count",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    },
    {
      "name": "job_queue_state_priority_count",
      "singularName": "job_queue_state_priority_count",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    },
    {
      "name": "job_queue_state_sched_count",
      "singularName": "job_queue_state_sched_count",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    },
    {
      "name": "job_queue_state_run_count",
      "singularName": "job_queue_state_run_count",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    },
    {
      "name": "job_queue_state_cleanup_count",
      "singularName": "job_queue_state_cleanup_count",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    },
    {
      "name": "job_queue_state_inactive_count",
      "singularName": "job_queue_state_inactive_count",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    }
  ]
}
```

### Submit A Job

Let's now submit a job to flux and confirm we see it. First, here is how to get the nodes being used:

```
# How many node are free?
curl -s http://localhost/apis/custom.metrics.k8s.io/v1beta2/namespaces/flux-operator/metrics/node_free_count | jq
```
```console
{
  "items": [
    {
      "metric": {
        "name": "node_free_count"
      },
      "value": 4,
      "timestamp": "2023-12-13T02:31:47+00:00",
      "windowSeconds": 0,
      "describedObject": {
        "kind": "Service",
        "namespace": "default",
        "name": "custom-metrics-apiserver",
        "apiVersion": "v1beta2"
      }
    }
  ],
  "apiVersion": "custom.metrics.k8s.io/v1beta2",
  "kind": "MetricValueList",
  "metadata": {
    "selfLink": "/apis/custom.metrics.k8s.io/v1beta2"
  }
}
```

Or more specifically, here are the free nodes before:

```bash
$ curl -s http://localhost/apis/custom.metrics.k8s.io/v1beta2/namespaces/flux-operator/metrics/node_free_count | jq -r .items[0].value
4
```
Let's use some of those nodes :) Shell into the lead broker:

```bash
 kubectl exec -it flux-sample-0-rfpgz  bash
```
```
[root@flux-sample-0 /]# . /mnt/flux/flux-view.sh 
[root@flux-sample-0 /]# flux proxy $fluxsocket bash
```

Submit a sleep job to 2 nodes.

```
flux submit -N 2 sleep 120
```

On the outside, try getting the metric again:

```bash
$ curl -s http://localhost/apis/custom.metrics.k8s.io/v1beta2/namespaces/flux-operator/metrics/node_free_count | jq -r .items[0].value
2
```

There you go! Super cool! Note that you can see our currently provided metrics [here](https://github.com/converged-computing/flux-metrics-api?tab=readme-ov-file#endpoints).
We expose metrics for flux nodes and the queue, and can add anything else really.


## Clean Up

When you are done, clean up.

```bash
kind delete cluster
```
