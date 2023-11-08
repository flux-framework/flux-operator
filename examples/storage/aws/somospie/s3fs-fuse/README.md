# Amazon Web Services

We are going to test [this CSI driver](https://github.com/s3fs-fuse/s3fs-fuse).

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

And you can find the name of the operator system pod as follows:

```bash
$ kubectl get pod --all-namespaces
```
```console
      <none>
operator-system   operator-controller-manager-6c699b7b94-bbp5q   2/2     Running   0             80s   192.168.30.43    ip-192-168-28-166.ec2.internal   <none>           <none>
```

## Run SOMOSPIE

#### S3 Secret

Let's first create a secret with our credentials. Normally you'd put this in a namespace you own,
but since this is our tiny cluster we will use default. You can first export them in the environment:

```bash
export AWS_ACCESS_KEY_ID=xxxxxxxx
export AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxx
export AWS_SESSION_TOKEN=xxx
```

Create the base54 encoded variants.

```bash
export AWS_ID=$(echo -n ${AWS_ACCESS_KEY_ID} | base64)
export AWS_SECRET=$(echo -n ${AWS_SECRET_ACCESS_KEY} | base64)
export AWS_TOKEN=$(echo -n ${AWS_SESSION_TOKEN} | base64 -w 0)
```

And then create the secret from them:

```bash
# How to test
$ cat secret.yaml | envsubst

# Save to temporary file (you might need to do this if credentials interfere with kubectl usage)
$ cat secret.yaml | envsubst > _secret.yaml
$ unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY
$ kubectl apply -f _secret.yaml
```

### Prepare S3 Storage

You should already know the S3 storage path that has the SOMOSPIE data.
You can sanity check that listing buckets shows your bucket:

```bash
$ aws s3 ls
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

Depending on the command you provided, you can then shell in and look around / test running workflows. Check the logs to see if the mount
had any issue, and see the minicluster.yaml for how to reproduce setting the envars and running the mount command.
Ensure you unmount before exiting:

```bash
umount /tmp/data
```

## Clean up

Make sure you clean everything up!
Detach roles:

```bash
$ eksctl delete iamserviceaccount --name s3-mounter --namespace otomount --cluster flux-operator
```

Delete the flux operator in some way:

```bash
$ make undeploy
$ kubectl delete -f ../../dist/flux-operator-dev.yaml
```

If you created roles, you probably want to clean these up too:

```bash
$ aws iam delete-role --role-name eks-otomounter-role
$ aws iam delete-policy --policy-arn arn:aws:iam::633731392008:policy/kubernetes-s3-access
$ aws iam delete-open-id-connect-provider --open-id-connect-provider-arn "arn:aws:iam::633731392008:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/xxxxxxxxxxxxxx"
```

And then delete your cluster (e.g., one of the following)

```bash
$ eksctl delete cluster -f ./eksctl-config.yaml --wait
```

Either way, it's good to check the web console too to ensure you didn't miss anything.
