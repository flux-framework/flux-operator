# Amazon Web Services

This set ofl tutorials wall walk through how to run the Flux Operator on AWS! 
You can start with [setup](#setup) and then move down to [examples](#examples).

## Setup

### Build

Note that we needed to be root in the container, so we have a custom build [Dockerfile](Dockerfile)

```bash
# Yes I built the wrong name and was too lazy to change it :)
docker build -t vanessa/somosie:latest
```

And this is the container that we use here.

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
$ eksctl create cluster -f eksctl-config.yaml
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

##### S3 Storage Policy

For our testing case, we made the public and created the following Permission -> Bucket policy:

```console
{
  "Version":"2012-10-17",
  "Statement":[
    {
      "Sid":"AddPerm",
      "Effect":"Allow",
      "Principal": "*",
      "Action": "s3:*",
      "Resource":["arn:aws:s3:::my-bucket-name/*"]
    }
  ]
}
```

You should obviously create a policy that would be associated with your user or credential (IAM) account.

#### Prepare the OIDC Provider for EKS

 **[MiniCluster YAML](minicluster.yaml)**

We will be following [this guide](https://dev.to/otomato_io/mount-s3-objects-to-kubernetes-pods-12f5).
First, [create an OIDC role for the cluster](https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html):

```bash
$ aws eks describe-cluster --name flux-operator --query "cluster.identity.oidc.issuer" --output text
```

Get the identifier `EXAMPLEXXXXXXXXXXXXXXX` and check if you've already done this. If you have, the following command will have output:

```bash
$ aws iam list-open-id-connect-providers | grep EXAMPLED539D4633E53DE1B7
```

If there is new output, open up the following section to create the OIDC provider:

<details>

<summary>Create the OIDC provider</summary>

```bash
$ eksctl utils associate-iam-oidc-provider --cluster flux-operator --approve
```

```console
2023-02-19 19:56:55 [ℹ]  will create IAM Open ID Connect provider for cluster "flux-operator" in "us-east-1"
2023-02-19 19:56:56 [✔]  created IAM Open ID Connect provider for cluster "flux-operator" in "us-east-1"
```

Then create a `policy.json` with your bucket name (you only need to do this once). The s3* is important
so we have permissions to mount, read, write, etc.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "s3:*",
            "Resource": [
                "arn:aws:s3:::<your-bucket-name>"
            ]
        }
    ]
}
```

Apply the policy (here is with the example and bucket we provide):

```bash
$ aws iam create-policy --policy-name kubernetes-s3-access --policy-document file://./policy.json
```
</details>

Once you have the policy, create the otomount namespace:

```bash
$ kubectl create namespace otomount
```

And then use eksctl to create the iam service account:

```bash
$ eksctl create iamserviceaccount --name s3-mounter --namespace otomount --cluster flux-operator \
    --role-name "eks-otomounter-role" --attach-policy-arn arn:aws:iam::633731392008:policy/kubernetes-s3-access --approve
```
```console
2023-03-01 13:32:53 [ℹ]  1 iamserviceaccount (otomount/s3-mounter) was included (based on the include/exclude rules)
2023-03-01 13:32:53 [!]  serviceaccounts that exist in Kubernetes will be excluded, use --override-existing-serviceaccounts to override
2023-03-01 13:32:53 [ℹ]  1 task: { 
    2 sequential sub-tasks: { 
        create IAM role for serviceaccount "otomount/s3-mounter",
        create serviceaccount "otomount/s3-mounter",
    } }2023-03-01 13:32:53 [ℹ]  building iamserviceaccount stack "eksctl-flux-operator-addon-iamserviceaccount-otomount-s3-mounter"
2023-03-01 13:32:53 [ℹ]  deploying stack "eksctl-flux-operator-addon-iamserviceaccount-otomount-s3-mounter"
2023-03-01 13:32:54 [ℹ]  waiting for CloudFormation stack "eksctl-flux-operator-addon-iamserviceaccount-otomount-s3-mounter"
2023-03-01 13:33:24 [ℹ]  waiting for CloudFormation stack "eksctl-flux-operator-addon-iamserviceaccount-otomount-s3-mounter"
2023-03-01 13:33:24 [ℹ]  created serviceaccount "otomount/s3-mounter"
```

Check to make sure it worked:

```bash
$ aws iam get-role --role-name eks-otomounter-role --query Role.AssumeRolePolicyDocument
```

And sanity check the attached role policies:

```bash
$ aws iam list-attached-role-policies --role-name eks-otomounter-role --query AttachedPolicies[].PolicyArn --output text
```
```console
arn:aws:iam::633731392008:policy/kubernetes-s3-access
```

Save that to a variable

```bash
export policy_arn=arn:aws:iam::633731392008:policy/kubernetes-s3-access
```

Look at it again, to sanity check the annotation of the arn:

```bash
$ aws iam get-policy --policy-arn $policy_arn
$ aws iam get-policy-version --policy-arn $policy_arn --version-id v1
```

Finally, ensure the role shows up alongside the service account:

```bash
$ kubectl describe serviceaccount s3-mounter -n otomount
```
Also ensure that the region your bucket is in matches the resources you are interacting with above!
At this point you have the correct service account and policy, and next need to create
a daemonset that will create the mount using that service account!

#### Install the S3 mounter

This storage drive can be installed with [helm](https://otomato-gh.github.io/s3-mounter), but for reproducibility we will
install directly from yaml (that was generated via helm). The difference here is that we have already created
the service account. In case you need to see how we generated the original oidc.yaml:

```bash
$ helm repo add otomount https://otomato-gh.github.io/s3-mounter
$ helm template s3-mounter otomount/s3-otomount  --namespace otomount --set bucketName="somospie" \
   --set iamRoleARN=arn:aws:iam::633731392008:policy/kubernetes-s3-access --create-namespace > ./oidc.yaml
```

To install the mounter pods (that are going to run `goofys`), we create a daemon set that will do the work as follows:

```bash
$ kubectl apply -f ./oidc.yaml
```

A few notes about this file:

 - the mount permissions assume that root (uid/gid 0) is going to run the workflow (and this is the default for version 2 of the Flux Operator)
 - the volume needs to be mounted in the pod under `/tmp/*` otherwise you won't be able to cleanup
 - the goofys flags are what we used, but it is not clear if all are needed.
 - I had to edit it quite a bit to get something working. Specifically, the command needed to create /var/s3

Check the pods are running:

```bash
$ kubectl get -n otomount pods
```

You can look at their logs to debug any issues with the mount or permissions. You should
ensure they are running with no obvious error before continuing!
Then (assuming you've already installed the operator and created the flux-operator namespace):

```bash
$ kubectl create -f ./minicluster.yaml
```

The biggest factor of whether your mount will work (with permission to read and write)
is determined by the S3 storage policy and rules. For testing, we were irresponsible
and made the bucket public, but you likely don't want to do that. Next,
get pods - you'll see the containers creating and then running:

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

Depending on the command you provided, you can then shell in and look around / test running workflows, or just
put the running command directly in the terminal.

You can get the output file in the terminal (and don't worry about saving it too much, as it will save to storage).
By default cleanup is set to False so you shouldn't lose the pod to get the pods for it.

```bash
$ kubectl get flux-sample-0-f5znt
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
$ eksctl delete cluster -f examples/storage/aws/oidc/eksctl-config.yaml --wait
$ eksctl delete cluster -f eksctl-config.yaml --wait
```

Either way, it's good to check the web console too to ensure you didn't miss anything.
