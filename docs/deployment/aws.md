# Amazon Web Services

This small tutorial wall walk through how to run the Flux Operator (from a development
standpoint) on AWS! You can start with setup (regardless of the workflow you choose)
and then move on to a cloud-specific workflow.

## Setup

### Install

You should first [install eksctrl](https://github.com/weaveworks/eksctl) and make sure you have access to an AWS cloud (e.g.,
with credentials or similar in your environment). E.g.,:

```bash
export AWS_ACCESS_KEY_ID=xxxxxxxxxxxxxxxxxxx
export AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export AWS_SESSION_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```
The last session token may not be required depending on your setup.
We assume you also have [kubectl](https://kubernetes.io/docs/tasks/tools/). 


### Setup SSH

This step is optional if you want to access your nodes.
You'll need an ssh key for EKS. Here is how to generate it:

```bash
ssh-keygen
# Ensure you enter the path to ~/.ssh/id_eks
```

This is used so you can ssh (connect) to your workers!

### Create Cluster

Next, let's create our cluster using eksctl "eks control." **IMPORTANT** you absolutely
need to choose a size that has [IsTrunkingCompatible](https://github.com/aws/amazon-vpc-resource-controller-k8s/blob/master/pkg/aws/vpc/limits.go)
true. Here is an example configuration. Note that we are choosing zones that works
for our account (this might vary for you) and an instance size that is appropriate
for our workloads. Also note that we are including the path to the ssh key we just
generated.

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: flux-operator
  region: us-east-1
  version: "1.22"

availabilityZones: ["us-east-1a", "us-east-1b", "us-east-1d"]
managedNodeGroups:
  - name: workers
    instanceType: c5.xlarge
    minSize: 4
    maxSize: 4
    labels: { "fluxoperator": "true" }
    ssh:
      allow: true
      publicKeyPath: ~/.ssh/id_eks.pub
```

If you don't need an ssh key:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: flux-operator
  region: us-east-1
  version: "1.23"

availabilityZones: ["us-east-1a", "us-east-1b", "us-east-1d"]
managedNodeGroups:
  - name: workers
    instanceType: c5.xlarge
    minSize: 4
    maxSize: 4
    labels: { "fluxoperator": "true" }
```

Given the above file `eks-cluster-config.yaml` we create the cluster as follows:

```bash
$ eksctl create cluster -f eks-cluster-config.yaml
```

🚧️ Warning! 🚧️ The above takes 15-20 minutes! Go have a party! Grab an avocado! 🥑️
And then come back and view your nodes.

```console
$ kubectl get nodes
NAME                             STATUS   ROLES    AGE     VERSION
ip-192-168-28-166.ec2.internal   Ready    <none>   4m58s   v1.22.12-eks-be74326
ip-192-168-4-145.ec2.internal    Ready    <none>   4m27s   v1.22.12-eks-be74326
ip-192-168-49-92.ec2.internal    Ready    <none>   5m3s    v1.22.12-eks-be74326
ip-192-168-79-92.ec2.internal    Ready    <none>   4m57s   v1.22.12-eks-be74326
```

### Deploy Operator 

To deploy the Flux Operator, [choose one of the options here](https://flux-framework.org/flux-operator/getting_started/user-guide.html#production-install) to deploy the operator. Whether you apply a yaml file, use [flux-cloud](https://converged-computing.github.io/flux-cloud) or clone the repository and `make deploy`. You can also deploy a development image:

```bash
$ make test-deploy
$ kubectl apply -f examples/dist/flux-operator-dev.yaml 
```

 you will see the operator install to the `operator-system` namespace.

```console
...
namespace/operator-system created
customresourcedefinition.apiextensions.k8s.io/miniclusters.flux-framework.org unchanged
serviceaccount/operator-controller-manager created
role.rbac.authorization.k8s.io/operator-leader-election-role created
clusterrole.rbac.authorization.k8s.io/operator-manager-role configured
clusterrole.rbac.authorization.k8s.io/operator-metrics-reader unchanged
clusterrole.rbac.authorization.k8s.io/operator-proxy-role unchanged
rolebinding.rbac.authorization.k8s.io/operator-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/operator-manager-rolebinding unchanged
clusterrolebinding.rbac.authorization.k8s.io/operator-proxy-rolebinding unchanged
configmap/operator-manager-config created
service/operator-controller-manager-metrics-service created
deployment.apps/operator-controller-manager created
```

Ensure the `operator-system` namespace was created:

```bash
$ kubectl get namespace
NAME              STATUS   AGE
default           Active   12m
kube-node-lease   Active   12m
kube-public       Active   12m
kube-system       Active   12m
operator-system   Active   11s
```
```bash
$ kubectl describe namespace operator-system
Name:         operator-system
Labels:       control-plane=controller-manager
              kubernetes.io/metadata.name=operator-system
Annotations:  <none>
Status:       Active

No resource quota.

No LimitRange resource.
```

And you can find the name of the operator pod as follows:

```bash
$ kubectl get pod --all-namespaces
```
```console
      <none>
operator-system   operator-controller-manager-6c699b7b94-bbp5q   2/2     Running   0             80s   192.168.30.43    ip-192-168-28-166.ec2.internal   <none>           <none>
```

Make your namespace for the flux-operator custom resource definition (CRD):

```bash
$ kubectl create namespace flux-operator
```

Then apply your CRD to generate the MiniCluster (default should be size 4, the max nodes of our cluster):

```bash
$ make apply
# OR
$ kubectl apply -f config/samples/flux-framework.org_v1alpha1_minicluster.yaml 
```

And now you can get logs for the manager:

```bash
$ kubectl logs -n operator-system operator-controller-manager-6c699b7b94-bbp5q
```

You'll see "errors" that the ip addresses aren't ready yet, and the operator
will reconcile until they are. You can add `-f` so the logs hang to watch:

```bash
$ kubectl logs -n operator-system operator-controller-manager-6c699b7b94-bbp5q -f
```

Once the logs indicate they are ready, you can look at the listing of nodes and the log for
the indexed job (which choosing one pod randomly to show):

```bash
$ make list
kubectl get -n flux-operator pods
NAME                  READY   STATUS    RESTARTS   AGE
flux-sample-0-zfmzc   1/1     Running   0          2m11s
flux-sample-1-p2hh5   1/1     Running   0          2m11s
flux-sample-2-zs4h6   1/1     Running   0          2m11s
flux-sample-3-prtn9   1/1     Running   0          2m11s
```

And when the containers are running, in the logs you'll see
see lots of cute emojis to indicate progress, and then the
start of your server! You'll need an exposed host to see the user
interface, or you can interact to submit jobs via the RESTful API.
A Python client is [available here](https://flux-framework.org/flux-restful-api/getting_started/user-guide.html#python).


## Clean up

Make sure you clean everything up!

```bash
$ make undeploy
```

And then:

```bash
$ eksctl delete cluster -f eks-cluster-config.yaml
```
It might be better to add `--wait`, which will wait until all resources are cleaned up:

```bash
$ eksctl delete cluster -f eks-cluster-config.yaml --wait
```
Either way, it's good to check the web console too to ensure you didn't miss anything.

