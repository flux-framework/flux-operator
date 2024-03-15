# Amazon Web Services

We are going to test [the Mountpoint s3 CSI driver](https://github.com/awslabs/mountpoint-s3-csi-driver).

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
  version: "1.27"

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

If you don't need an ssh key just remove the last "ssh" block.

Given the above file `eks-cluster-config.yaml` we create the cluster as follows:

```bash
$ eksctl create cluster -f eksctl-config.yaml
```

🚧️ Warning! 🚧️ The above takes 15-20 minutes! Go have a party! Grab an avocado! 🥑️
And then come back and view your nodes.

```console
$ kubectl get nodes
NAME                             STATUS   ROLES    AGE     VERSION
ip-192-168-25-25.ec2.internal    Ready    <none>   38m   v1.27.6-eks-a5df82a
ip-192-168-69-249.ec2.internal   Ready    <none>   38m   v1.27.6-eks-a5df82a
```

Note we are following the [install instructions here](https://github.com/awslabs/mountpoint-s3-csi-driver/blob/main/docs/install.md).
And taking the easiest path (not necessarily the best). They make it really hard otherwise. :/

```bash
export ROLE_ARN=arn:aws:iam::aws:policy/AmazonS3FullAccess
export ROLE_NAME=somospie-testing
```

Do this?

```bash
aws eks update-kubeconfig --region $REGION --name $CLUSTER_NAME
eksctl utils associate-iam-oidc-provider --cluster flux-operator --approve
```

Then associate the iam service account

```bash
eksctl create iamserviceaccount \
    --name s3-csi-driver-sa \
    --namespace kube-system \
    --cluster flux-operator \
    --attach-policy-arn $ROLE_ARN \
    --approve \
    --role-name $ROLE_NAME \
    --region us-east-1 \
    --override-existing-serviceaccounts
```

Ensure it exists:

```bash
kubectl describe sa s3-csi-driver-sa --namespace kube-system
```

With your credentials exported, do:


```bash
kubectl create secret generic aws-secret \
    --namespace kube-system \
    --from-literal "key_id=${AWS_ACCESS_KEY_ID}" \
    --from-literal "access_key=${AWS_SECRET_ACCESS_KEY}"
```
```bash
$ kubectl describe secret --namespace kube-system
```

Deploy the driver and verify pods are running:

```bash
kubectl apply -k "github.com/awslabs/mountpoint-s3-csi-driver/deploy/kubernetes/overlays/stable/"
kubectl get pods -n kube-system -l app.kubernetes.io/name=aws-mountpoint-s3-csi-driver
```

### Deploy Operator

To deploy the Flux Operator, [choose one of the options here](https://flux-framework.org/flux-operator/getting_started/user-guide.html#production-install) to deploy the operator. Whether you apply a yaml file, use [flux-cloud](https://converged-computing.github.io/flux-cloud) or clone the repository and `make deploy`. You can also deploy a development image:

```bash
$ make test-deploy-recreate
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

And you can find the name of the operator system pod as follows:

```bash
$ kubectl get pod --all-namespaces
```
```console
      <none>
operator-system   operator-controller-manager-6c699b7b94-bbp5q   2/2     Running   0             80s   192.168.30.43    ip-192-168-28-166.ec2.internal   <none>           <none>
```

## Run SOMOSPIE

### Prepare S3 Storage

You should already know the S3 storage path that has the SOMOSPIE data.
You can sanity check that listing buckets shows your bucket:

```bash
$ aws s3 ls
```

Create the PV and PVC and ensure they are bound:

```bash
kubectl apply -f pv-pvc.yaml
```

Let's now create the MiniCluster!


```bash
$ kubectl create -f ./minicluster.yaml
```

Wait until the init is done:

```bash
$ kubectl  get pods --watch
NAME                  READY   STATUS     RESTARTS   AGE
flux-sample-0-dpnj7   0/1     Init:0/1   0          33s
flux-sample-1-7kdqx   0/1     Init:0/1   0          33s
```

And then the pods should create:

```bash
$ kubectl get pods
```
```console
NAME                         READY   STATUS              RESTARTS   AGE
flux-sample-0-f5znt          0/1     ContainerCreating   0          100s
flux-sample-1-th589          0/1     ContainerCreating   0          100s
```

Shell inside the pod:

```bash
kubectl exec -it flux-sample-0-xxx bash
```

And you should see the data!

```bash
$ ls /data/november-2023/
oklahoma-10m  oklahoma-27km  oklahoma-30m
```

## Clean up

Make sure you clean everything up!

```bash
$ eksctl delete cluster -f ./eksctl-config.yaml --wait
```

Either way, it's good to check the web console too to ensure you didn't miss anything.