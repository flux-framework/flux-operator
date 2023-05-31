# Horizontal Autoscaler with Custom Metrics

> Using the horizontal pod autoscaler (HPA) (autoscaler/v2) with custom metrics from Flux

For more background on autoscaling, see [this page](https://flux-framework.org/flux-operator/tutorials/elasticity.html).
To be repetitive, we use the version 2 autoscaler here that is more flexible to custom metrics, and we are also
going to configure our MiniCluster to allow Flux to give us numbers directly using

Another strategy would be to pair the Flux Operator with Prometheus and use [prometheus-flux](https://github.com/converged-computing/prometheus-flux)
but we have not made a tutorial for that yet! Let us know if this is of interest to you. For this tutorial here we will be installing the
[flux metrics api](https://github.com/converged-computing/flux-metrics-api) that serves a Kubernetes API server.
Note that this is in Python, but I realized after I wrote this we could create [one in Go too](https://github.com/kubernetes-sigs/custom-metrics-apiserver).

Note that for this setup I've prepared it so that for a production cluster the operator will support SSL. However,
developing locally with kubectl I found I needed to disable it.

## Setup

You'll want to perform these steps from this root as the present working directory.

```bash
$ cd examples/elasticity/horizontal-autoscaler/v2-custom-metrics
```

Create a kind cluster with Kubernetes version 1.27.0

```bash
$ kind create cluster --config ./kind-config.yaml
```

Create the flux-operator namespace and install the operator:

```bash
$ kubectl create namespace flux-operator
$ kubectl apply -f ../../../dist/flux-operator-dev.yaml
```

## Secrets

We are going to be making an API server that is served directly from the Flux broker leader (the index 0 pod) and this
is going to provide an endpoint to serve metrics directly. Since this needs to have SSL, we need to generate certificates
and add them to the cluster as a secret for the MiniCluster to bind to `/etc/certs`. Let's do that first.
We will follow the pattern [here](https://github.com/kflansburg/py-custom-metrics/tree/39ad121047b8c798dde380f94abc97b1589ba4ed/scripts). 
First, cd into the scripts folder to make the certificates:

```bash
$ cd ./scripts
$ ./certs.sh
```
```console
Generating RSA private key, 2048 bit long modulus (2 primes)
...................................................................+++++
......................................................................+++++
e is 65537 (0x010001)
Generating a RSA private key
............................+++++
...+++++
writing new private key to 'server.key'
-----
Signature ok
subject=CN = custom-metrics-apiserver.custom-metrics.svc
Getting CA Private Key
```

This should generate `ca.crt` and `server.crt` that we will create a secret with.
Since this is being served from the `flux-operator` namespace we will add it there -
it's scoped and belonging to the MiniCluster (for now).

```bash
$ kubectl create secret tls -n flux-operator certs --cert=server.crt --key=server.key
```

Make sure that you can see it!

```bash
$ kubectl get secret -n flux-operator 
NAME    TYPE                DATA   AGE
certs   kubernetes.io/tls   2      4s
```

When you are done you can cd out of the "scripts" directory.

```bash
$ cd ../
```

## MiniCluster

Now let's make the MiniCluster. It's going to have a volume that maps the secret into the pods at `/etc/certs` 
(see the "existingVolume" in the minicluster.yaml). This is going to be what our API server uses. 
sNormally, we would have an API server in a separate pod (and if we are
able to figure out communication between a Flux broker and this external pod we could still do that) but for now
we are putting them together. This will create a very simply interactive cluster (it's small, 2 pods, but importantly has a maxsize of 10).
Note that we are going to limit the HPA to a size of 4, because we assume you are running on an average desktop computer.

```bash
$ kubectl apply -f ./minicluster.yaml
```

You'll need to wait for the container to pull (status `ContainerCreating` to `Running`).
At this point, wait until the containers go from creating to running.

```bash
$ kubectl get -n flux-operator pods
NAME                  READY   STATUS    RESTARTS   AGE
flux-sample-0-4wmmp   1/1     Running   0          6m50s
flux-sample-1-mjj7b   1/1     Running   0          6m50s
```

If you describe the pods in the operator namespace, you should see the secret mounted at `/etc/certs`. This is important!

```bash
$ kubectl describe -n flux-operator pods
...
Volumes:
...
  certs:
    Type:        Secret (a volume populated by a Secret)
    SecretName:  certs
    Optional:    false
```

Look at the scale endpoint of the MiniCluster with `kubectl` directly! Remember that we haven't installed a horizontal auto-scaler yet:

```bash
$ kubectl get --raw /apis/flux-framework.org/v1alpha1/namespaces/flux-operator/miniclusters/flux-sample/scale | jq
```
```console
{
  "kind": "Scale",
  "apiVersion": "autoscaling/v1",
  "metadata": {
    "name": "flux-sample",
    "namespace": "flux-operator",
    "uid": "581c708a-0eb2-48da-84b1-3da7679d349d",
    "resourceVersion": "3579",
    "creationTimestamp": "2023-05-20T05:11:28Z"
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

The above knows the selector to use to get pods (and look at current resource usage). Note that
the reference to "autoscaling/v1" does not mean we cannot use autoscaling/v1 for the one we create.
This threw me off for a bit!

## Metrics Server

TODO not sure if I need this or not

Before we deploy any autoscaler, we need a main metrics server! This doesn't come out of the box with kind so
we install it:

```bash
$ kubectl apply -f metrics-server.yaml
```

I found this suggestion [here](https://gist.github.com/sanketsudake/a089e691286bf2189bfedf295222bd43). Ensure
it's running:

```bash
$ kubectl get deploy,svc -n kube-system | egrep metrics-server
```

## Custom Metrics API Service

We next want to setup our custom metrics API service. Since this is a tutorial let's start this manually
for now, shelling into the broker pod, then connecting to the broker and running the server from there:


```bash
$ kubectl exec -it -n flux-operator flux-sample-0-p85cj -- bash
$ sudo -u fluxuser -E $(env) -E HOME=/home/fluxuser flux proxy local:///run/flux/local bash
```

This is how to start the metrics exporter (using defaults). Note that the default port should match
the one in our servica.yaml and we are pointing to the certificates to use!

```bash
$ flux-metrics-api start --port 8443 --ssl-certfile /etc/certs/tls.crt --ssl-keyfile /etc/certs/tls.key --namespace flux-operator --service-name custom-metrics-apiserver
```

Note that most of these are defaults, but it doesn't hurt to set them explicitly.
The easiest thing to do is run the above in the background OR open a separate terminal (recommended so you can continue to monitor the server output).
Once you have issued the two commands above again in a different terminal, test the endpoint (any of these should work):

```bash
$ curl -s http://0.0.0.0:8080/apis/custom.metrics.k8s.io/v1beta2/namespaces/flux-operator/metrics/node_up_count | jq
$ curl -s http://flux-sample-0:8080/apis/custom.metrics.k8s.io/v1beta2/namespaces/flux-operator/metrics/node_up_count | jq
$ curl -s http://flux-sample-0.flux-service.flux-operator.svc.cluster.local:8080/apis/custom.metrics.k8s.io/v1beta2/namespaces/flux-operator/metrics/node_up_count | jq
```
```console
{
  "items": [
    {
      "metric": {
        "name": "node_up_count"
      },
      "value": 2,
      "timestamp": "2023-05-30T22:20:07+00:00",
      "windowSeconds": 0,
      "describedObject": {
        "kind": "Service",
        "namespace": "flux-operator",
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

The output we are seeing above is what is expected (in terms of format) for a custom metrics endpoint.
It's nice to keep the second window open because you can see the requests:

```console
INFO:     127.0.0.1:38402 - "GET /apis/custom.metrics.k8s.io/v1beta2/namespaces/flux-operator/metrics/node_up_count HTTP/1.1" 200 OK
INFO:     10.244.0.6:38818 - "GET /apis/custom.metrics.k8s.io/v1beta2/namespaces/flux-operator/metrics/node_up_count HTTP/1.1" 200 OK
```

## Custom Metrics Service

OK to step back, at this point we have a custom metrics server running from inside of a pod, and we need to somehow tell
the cluster to provide a service with a particular address being served from that pod. Here are the logical steps we will take
to do that:

1. Add a label selector on the index 0 pod (the leader broker with the metrics API running)
2. Create a service that uses the selector to point the particular port service to the pod
3. Create an API service that targets that service name to create a cluster-scoped API

It's Kubernetes so the order of operations (and figuring this out to begin with) was kind of weird.
Let's do those steps one at a time. First, adding the selector label to the leader broker pod:

```bash
$ kubectl label pods -n flux-operator flux-sample-0-xxx api-server=custom-metrics
```

Now let's create the service that knows how to select that.

```bash
$ kubectl apply -f ./scripts/service.yaml
```

We want to see that there is a cluster IP address serving a secure port:

```bash
kubectl get svc -n flux-operator 
NAME                       TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
custom-metrics-apiserver   ClusterIP   10.96.20.246   <none>        443/TCP   5s
flux-service               ClusterIP   None           <none>        <none>    12m
```

Note that the name of our service is "custom-metrics-apiserver." This is exactly the name
we will use for the API service! We also would need to provide the certificate here,
but for now I'm disabling it because it was adding complexity.

```bash
# Without TLS (no changes needed to file)
$ kubectl apply -f ./scripts/api-service.yaml

# This would be WITH TLS (uncomment lines in file)
# This needs to be -b=0 for Darwin
export CA_BUNDLE=$(cat ./scripts/ca.crt | base64 --wrap=0)
cat ./scripts/api-service.yaml | envsubst | kubectl apply -f -
```

If you get an error, ensure the versions match up of the api registration:

```bash
$ kubectl api-resources | grep apiregistration
apiservices                                    apiregistration.k8s.io/v1              false        APIService
```

Here is how to debug the service.

```bash
$ kubectl describe apiservice v1beta2.custom.metrics.k8s.io
```

For example, when I first created it, I hadn't actually added SSL / certificates to my endpoints
so I saw an error that the connection was refused. When it works, you will see this endpoint get hit a LOT.

```console
INFO:     10.244.0.1:4372 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
INFO:     10.244.0.1:34610 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
INFO:     10.244.0.1:40447 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
INFO:     10.244.0.1:54072 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
INFO:     10.244.0.1:6895 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
INFO:     10.244.0.1:64937 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
INFO:     10.244.0.1:48753 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
INFO:     10.244.0.1:58257 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
INFO:     10.244.0.1:17035 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
INFO:     10.244.0.1:54047 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
```

I think it's doing a constant health check, and this is why we have to provide essentially an empty 200 response
there to get it working. As a reminder of what is happening?

- That service.yaml has a selector that finds the pods running the service
- The service.yaml name matches the name of the api-service
- The pods that are selected need to have those endpoints and the ssl!

This is what it should look like when it's working:

```bash
kubectl get apiservice v1beta2.custom.metrics.k8s.io
NAME                            SERVICE                                  AVAILABLE   AGE
v1beta2.custom.metrics.k8s.io   flux-operator/custom-metrics-apiserver   True        22m
```

And wow - we now have an API service running at this endpoint! Let's make some autoscaler stuff!

## Custom Metrics

We now should be able to retrieve a custom metrics endpoint (for our custom metrics directly from Flux)!
Let's create the autoscaler to use it:

```bash
$ kubectl apply -f hpa-flux.yaml
```
We should also be able to ping our new metrics server directly with kubectl. This is to look at the node up count:

```bash
$ kubectl get --raw /apis/custom.metrics.k8s.io/v1beta2/namespaces/flux-operator/services/custom-metrics-service/node_up_count | jq
```

Note that when I was working on this, I first hit [this error](https://github.com/kubernetes/kubernetes/blob/18e3f01deda3bc1ea62751553df0b689598de7a7/staging/src/k8s.io/metrics/pkg/client/custom_metrics/discovery.go#L101) and had to find that spot in the source code, and then realize that 
the `/apis` root endpoint was being pinged for some kind of "preferred version." I tried mocking
the endpoint (and it seemed to work?) so then I simply got the actual endpoint from within the
pod, and forwarded it along. This seemed to make the cluster happy - I started seeing the
autoscaler actually pinging my server for the metric!

```console
INFO:     10.244.0.1:10333 - "GET /openapi/v2 HTTP/1.1" 200 OK
Requested metric node_up_count in  namespace flux-operator
INFO:     10.244.0.1:31834 - "GET /apis/custom.metrics.k8s.io/v1beta2/namespaces/flux-operator/services/custom-metrics-apiserver/node_up_count HTTP/1.1" 200 OK
INFO:     10.244.0.1:12095 - "GET /apis HTTP/1.1" 200 OK
INFO:     10.244.0.1:33736 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
INFO:     10.244.0.1:43900 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
INFO:     10.244.0.1:53777 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
INFO:     10.244.0.1:28114 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
INFO:     10.244.0.1:4014 - "GET /apis/custom.metrics.k8s.io/v1beta2 HTTP/1.1" 200 OK
Requested metric node_up_count in  namespace flux-operator
INFO:     10.244.0.1:16713 - "GET /apis/custom.metrics.k8s.io/v1beta2/namespaces/flux-operator/services/custom-metrics-apiserver/node_up_count HTTP/1.1" 200 OK
```

At this point we are retrieving the metric, although we haven't really added any logic for what to do with it. I'll likely work on this next!

I noticed that a few times over a long period of time there would be an error issued that the server was "unable to handle the request" and I think this is related to
memory. Note that another strategy I haven't looked into is using [External](https://github.com/GoogleCloudPlatform/bank-of-anthos/blob/a32d4cf14a6a030705f00fc9d0dbf2d547ef1231/extras/postgres-hpa/hpa/frontend.yaml#L16) for autoscaling.


Get logs for HPA
$ kubectl get -n flux-operator hpa -w