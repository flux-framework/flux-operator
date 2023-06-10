# Nginx

> Deploy an nginx service alongside your cluster

The idea of this example is that you could have some application (that connects) to the rest
of your cluster) served through a sidecar container. For this dummy example, we will:

1. Create a kind cluster that can be exposed on localhost
2. Deploy an nginx service pod alongside it
3. Define an nginx.conf custom file that you could further customize

And for the purposes of the example, we will also show how you can interact
with the Flux Restful API (running from the broker) from the nginx service
pod. Our use case (for the time being) is that we want our local MiniCluster
to be able to interact with external clusters, although that is another thing
to figure out.

## Usage

### Setup Cluster

Create the cluster

```bash
$ kind create cluster --config ./kind-config.yaml
```

Note this config ensures that the cluster has an external IP on localhost.
Then install the operator, create the namespace, and the minicluster.

```bash
$ kubectl apply -f ../../../dist/flux-operator-dev.yaml
$ kubectl create namespace flux-operator
```

Create an existing config map for nginx (that it will expect to be there, and we 
have defined in our minicluster.yaml under the nginx service existingVolumes).

```bash
$ kubectl apply -f nginx.yaml
```

And then create the MiniCluster

```bash
$ kubectl apply -f minicluster.yaml
```

And when the containers are ready, create the node port service for "flux-services":

```bash
$ kubectl apply -f service.yaml
```

Here is a nice block to easily copy paste all three :)

```bash
kubectl apply -f nginx.yaml
kubectl apply -f service.yaml
kubectl apply -f minicluster.yaml
```

It will take a few minutes for the containers to pull and pods to create.
Move on to the next step when this is finished. Depending on the order that nginx
and the broker index 0 pod come up, nginx may need to restart.

### External Service

This one is easy - you can go to [http://0.0.0.0:30093/](http://0.0.0.0:30093/)
to see the welcome to Nginx! page. Any actual application that you could be serving would
be forwarded there.


### Flux Restful from the Services Pod

Note that since we have disabled interactive and launcher modes and not provided a command in the
minicluster.yaml, the restful API is going to start. We have also defined our username
and token and secret key in the MiniCluster.yaml so we have them handy.
We can test that the flux-service container is able to interact with the restful API
from our service container. First, shell in:

```bash
$ kubectl exec -it -n flux-operator flux-sample-services bash
```

You may need to install ping, and jq, and we will make our lives easier by installing python with pip:

```bash
apt-get update && apt-get install -y iputils-ping jq python3-pip
```

We are now going to interact with the flux broker on the headless network, just to show you that we can!
Install the Flux Restful API client:

```bash
$ pip install flux-restful-client
```

Export your credentials (these are in the minicluster.yaml)

```bash
export FLUX_USER=flux
export FLUX_TOKEN="theogreisanonion"
export FLUX_SECRET_KEY=theonionisanogre
export FLUX_RESTFUL_HOST=http://flux-sample-0.flux-service.flux-operator.svc.cluster.local:5000
```

And try listing nodes for the host we know the restful API is running from:

```bash
$ flux-restful-cli list-nodes
```
```console
{
    "nodes": [
        "flux-sample-0",
        "flux-sample-1"
    ]
}
```
Huzzah! Amazing! This means the whole authentication flow is working, and the pods
can internally see one another. 


### Cleanup

Note that you'll only see the exposure with the kind docker container with `docker ps`.
When you are done, clean up

```bash
kubectl delete -f minicluster.yaml
kubectl delete -f nginx.yaml
kubectl delete -f service.yaml
```

