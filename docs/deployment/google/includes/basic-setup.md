## Create Cluster

Let's use gcloud to create a cluster, and we are purposefully going to choose
a very small node type to test. Note that I choose us-central1-a because it tends
to be cheaper (and closer to me). First, here is our project for easy access:

```bash
GOOGLE_PROJECT=myproject
```

Then create your cluster as follows:

```bash
$ gcloud container clusters create flux-cluster --project $GOOGLE_PROJECT \
    --zone us-central1-a --machine-type n1-standard-1 \
    --num-nodes=4 --enable-network-policy --tags=flux-cluster --enable-intra-node-visibility
```

If you need a particular Kubernetes version:

```bash
$ gcloud container clusters create flux-cluster --project $GOOGLE_PROJECT \
    --zone us-central1-a --cluster-version 1.23 --machine-type n1-standard-1 \
    --num-nodes=4 --enable-network-policy --tags=flux-cluster --enable-intra-node-visibility
```

Note that not all of the flags above might be necessary - I did a lot of testing to get
this working and didn't go back and try removing things after the fact!
If you want to use cloud dns instead (after [enabling it](https://console.cloud.google.com/apis/library/dns.googleapis.com))

```bash
$ gcloud beta container clusters create flux-cluster --project $GOOGLE_PROJECT \
    --zone us-central1-a --cluster-version 1.23 --machine-type n1-standard-1 \
    --num-nodes=4 --enable-network-policy --tags=flux-cluster --enable-intra-node-visibility \
    --cluster-dns=clouddns \
    --cluster-dns-scope=cluster
```

In your Google cloud interface, you should be able to see the cluster! Note
this might take a few minutes.

![img/cluster.png](img/cluster.png)

I also chose a tiny size (nodes and instances) anticipating having it up longer to figure things out.

## Get Credentials

Next we need to ensure that we can issue commands to our cluster with kubectl.
To get credentials, in the view shown above, select the cluster and click "connect."
Doing so will show you the correct statement to run to configure command-line access,
which probably looks something like this:

```bash
$ gcloud container clusters get-credentials flux-cluster --zone us-central1-a --project $GOOGLE_PROJECT
```
```console
Fetching cluster endpoint and auth data.
kubeconfig entry generated for flux-cluster.
```

Finally, use [cloud IAM](https://cloud.google.com/iam) to ensure you can create roles, etc.

```bash
$ kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user $(gcloud config get-value core/account)
```
```console
clusterrolebinding.rbac.authorization.k8s.io/cluster-admin-binding created
```

At this point you should be able to get your nodes:

```bash
$ kubectl get nodes
```
```console
NAME                                            STATUS   ROLES    AGE     VERSION
gke-flux-cluster-default-pool-f103d9d8-379m   Ready    <none>   3m41s   v1.23.14-gke.1800
gke-flux-cluster-default-pool-f103d9d8-3wf9   Ready    <none>   3m42s   v1.23.14-gke.1800
gke-flux-cluster-default-pool-f103d9d8-c174   Ready    <none>   3m42s   v1.23.14-gke.1800
gke-flux-cluster-default-pool-f103d9d8-zz1q   Ready    <none>   3m42s   v1.23.14-gke.1800
```

### Deploy Operator

To deploy the Flux Operator, [choose one of the options here](https://flux-framework.org/flux-operator/getting_started/user-guide.html#production-install) to deploy the operator. Whether you apply a yaml file, use [flux-cloud](https://converged-computing.github.io/flux-cloud) or clone the repository and `make deploy` you will see the operator install to the `operator-system` namespace.

For a quick "production deploy" from development, the Makefile has a directive that will build and push a `test` tag (you'll need to edit `DEVIMG` to be one you can push to) and then generate a
yaml file targeting that image, e.g.,

```bash
$ make test-deploy
$ kubectl apply -f examples/dist/flux-operator-dev.yaml
```

or the production version:

```bash
$ kubectl apply -f examples/dist/flux-operator.yaml
```

```console
...
clusterrole.rbac.authorization.k8s.io/operator-manager-role created
clusterrole.rbac.authorization.k8s.io/operator-metrics-reader created
clusterrole.rbac.authorization.k8s.io/operator-proxy-role created
rolebinding.rbac.authorization.k8s.io/operator-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/operator-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/operator-proxy-rolebinding created
configmap/operator-manager-config created
service/operator-controller-manager-metrics-service created
deployment.apps/operator-controller-manager created
```

Ensure the `operator-system` namespace was created:

```bash
$ kubectl get namespace
NAME              STATUS   AGE
default           Active   6m39s
kube-node-lease   Active   6m42s
kube-public       Active   6m42s
kube-system       Active   6m42s
operator-system   Active   39s
```
```bash
$ kubectl describe namespace operator-system
Name:         operator-system
Labels:       control-plane=controller-manager
              kubernetes.io/metadata.name=operator-system
Annotations:  <none>
Status:       Active

Resource Quotas
  Name:                              gke-resource-quotas
  Resource                           Used  Hard
  --------                           ---   ---
  count/ingresses.extensions         0     100
  count/ingresses.networking.k8s.io  0     100
  count/jobs.batch                   0     5k
  pods                               1     1500
  services                           1     500

No LimitRange resource.
```

And you can find the name of the operator pod as follows:

```bash
$ kubectl get pod --all-namespaces
```
```console
      <none>
operator-system   operator-controller-manager-56b5bcf9fd-m8wg4               2/2     Running   0          73s
```

### Create Flux Operator namespace

Make your namespace for the flux-operator custom resource definition (CRD), which is the yaml file above that generates the MiniCluster:

```bash
$ kubectl create namespace flux-operator
```