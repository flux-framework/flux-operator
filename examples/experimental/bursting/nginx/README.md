# Bursting Nginx Experiment

> Experimental setup with nginx as a proxy

This original experiment used a separate nginx service to serve as a proxy to a broker.
In this experiment I was able to demonstrate connecting to a lead broker (via restful API)
from an nginx service, and if we want to continue this approach, we would need to develop
some application to run there that receives requests and forwards them to the correct location.
This work started with the [nginx service](../../services/sidecar/nginx/) example.

## Google Cloud Setup

Let's create the first cluster:

```bash
CLUSTER_NAME=flux-cluster
GOOGLE_PROJECT=myproject
```
```bash
$ gcloud container clusters create ${CLUSTER_NAME} --project $GOOGLE_PROJECT \
    --zone us-central1-a --machine-type n2-standard-8 \
    --num-nodes=4 --enable-network-policy --tags=flux-cluster --enable-intra-node-visibility
```

And be sure to activate your credentials!

```bash
gcloud container clusters get-credentials ${CLUSTER_NAME}
```
Install the operator and create the minicluster

```bash
kubectl apply -f ../../dist/flux-operator-dev.yaml
kubectl apply -f service/nginx.yaml
kubectl apply -f minicluster.yaml
kubectl apply -f service/service.yaml
```

We need to open up the firewall to that port - this creates the rule (you only need to do this once)

```bash
gcloud compute firewall-rules create flux-cluster-test-node-port --allow tcp:30093
```

Then figure out the node that the service is running from:

```bash
$ kubectl get pods -o wide -n flux-operator 
NAME                   READY   STATUS    RESTARTS   AGE     IP           NODE                                          NOMINATED NODE   READINESS GATES
flux-sample-0-kktl7    1/1     Running   0          7m22s   10.116.2.4   gke-flux-cluster-default-pool-4dea9d5c-0b0d   <none>           <none>
flux-sample-1-s7r69    1/1     Running   0          7m22s   10.116.1.4   gke-flux-cluster-default-pool-4dea9d5c-1h6q   <none>           <none>
flux-sample-services   1/1     Running   0          7m22s   10.116.0.4   gke-flux-cluster-default-pool-4dea9d5c-lc1h   <none>           <none>
```

Get the external ip for that node (for nginx, flux-services, and for the lead broker, a flux-sample-0-xx)

```bash
$ kubectl get nodes -o wide | grep gke-flux-cluster-default-pool-4dea9d5c-0b0d 
gke-flux-cluster-default-pool-4dea9d5c-0b0d   Ready    <none>   69m   v1.25.8-gke.500   10.128.0.83   34.171.113.254   Container-Optimized OS from Google   5.15.89+         containerd://1.6.18
```

You can try opening `34.135.221.11:30093`. For the broker (34.171.113.254) above and you should see the "Welcome to nginx!" page. 
At this point, you could setup some kind of bursting application to be run from that service instead,
and then manage communication to different brokers.  I didn't proceed with this design because the [broker](../broker) design
(with a direct connection to a broker) was more straight forward for this early testing.

### Cleanup

Note that you'll only see the exposure with the kind docker container with `docker ps`.
When you are done, clean up

```bash
kubectl delete -f minicluster.yaml
kubectl delete -f nginx.yaml
kubectl delete -f service.yaml
gcloud container clusters delete flux-cluster
```


