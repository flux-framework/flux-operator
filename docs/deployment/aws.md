# Amazon Web Services

This set ofl tutorials wall walk through how to run the Flux Operator on AWS! 
You can start with [setup](#setup) and then move down to [examples](#examples).

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
$ eksctl create cluster -f eksctl-config.yaml

# use the provided
$ eksctl create cluster -f ./examples/storage/aws/oidc/eksctl-config.yaml
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

### Create the flux-operator namespace

Make your namespace for the flux-operator custom resource definition (CRD):

```bash
$ kubectl create namespace flux-operator
```

## Examples

After setup, choose one of the following examples to run.

### Run LAMMPS

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

And when the containers are running, in the logs you'll see see lots of cute emojis to indicate progress, and then the
start of your server! You'll need an exposed host to see the user interface, or you can interact to submit jobs via the RESTful API.
A Python client is [available here](https://flux-framework.org/flux-restful-api/getting_started/user-guide.html#python).
To use Flux Cloud to programmatically submit jobs, [see the guides here](https://converged-computing.github.io/flux-cloud/getting_started/aws.html).

### Run Snakemake with a Shared Filesystem

This small tutorial will run a Snakemake workflow on AWS that requires a shared
filesystem.

#### Prepare S3 Storage

Let's first get our Snakemake pre-requisite analysis files into S3.
To start, prepare your data in a temporary directory.

```bash
$ git clone --depth 1 https://github.com/snakemake/snakemake-tutorial-data /tmp/workflow
```

You'll want to add the [Snakefile](https://github.com/rse-ops/flux-hpc/blob/main/snakemake/atacseq/Snakefile) for your workflow
along with a [plotting script](https://github.com/rse-ops/flux-hpc/blob/main/snakemake/atacseq/scripts/plot-quals.py):

```bash
$ wget -O /tmp/workflow/Snakefile https://raw.githubusercontent.com/rse-ops/flux-hpc/main/snakemake/atacseq/Snakefile
$ mkdir -p /tmp/workflow/scripts
$ wget -O /tmp/workflow/scripts/plot-quals.py https://raw.githubusercontent.com/rse-ops/flux-hpc/main/snakemake/atacseq/scripts/plot-quals.py
```

You should have this structure:

```bash
$ tree /tmp/workflow
```
```
/tmp/workflow/
├── data
│   ├── genome.fa
│   ├── genome.fa.amb
│   ├── genome.fa.ann
│   ├── genome.fa.bwt
│   ├── genome.fa.fai
│   ├── genome.fa.pac
│   ├── genome.fa.sa
│   └── samples
│       ├── A.fastq
│       ├── B.fastq
│       └── C.fastq
├── Dockerfile
├── environment.yaml
├── README.md
├── scripts
│   └── plot-quals.py
└── Snakefile
```

We can then use the `aws` command line client to make a bucket "mb" and upload to it.

```bash
$ aws s3 mb s3://flux-operator-bucket --region us-east-1
```

Sanity check that listing buckets shows your bucket:

```bash
$ aws s3 ls
```
```console
2023-02-18 18:14:32 flux-operator-workflows
...
```
Now copy the entire workflow to a faux "subdirectory" there:

```bash
# Present working directory 
$ aws s3 cp --recursive . s3://flux-operator-bucket/snakemake-workflow --exclude ".git*"
```

Sanity check again by listing that path in the bucket

```bash
$ aws s3 ls s3://flux-operator-bucket/snakemake-workflow/
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
$ aws iam create-policy --policy-name kubernetes-s3-access --policy-document file://./examples/storage/aws/oidc/policy.json
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
$ helm template s3-mounter otomount/s3-otomount  --namespace otomount --set bucketName="flux-operator-bucket" \
   --set iamRoleARN=arn:aws:iam::633731392008:policy/kubernetes-s3-access --create-namespace > ./examples/aws/oidc/oidc.yaml
```

To install the mounter pods (that are going to run `goofys`), we create a daemon set that will do the work as follows:

```bash
$ kubectl apply -f ./examples/storage/aws/oidc/oidc.yaml
```

A few notes about this file:

 - the mount permissions assume that root (uid/gid 0) is going to run the workflow, `runFluxAsRoot` is set to true
 - the volume needs to be mounted in the pod under `/tmp/*` otherwise you won't be able to cleanup
 - the goofys flags are what we used, but it is not clear if all are needed.

Check the pods are running:

```bash
$ kubectl get -n otomount pods
```

You can look at their logs to debug any issues with the mount or permissions. You should
ensure they are running with no obvious error before continuing!
Then (assuming you've already installed the operator and created the flux-operator namespace):

```bash
$ kubectl create -f ./examples/storage/aws/oidc/minicluster.yaml
```

The biggest factor of whether your mount will work (with permission to read and write)
is determined by the S3 storage policy and rules. For testing, we were irresponsible
and made the bucket public, but you likely don't want to do that. Next,
get pods - you'll see the containers creating and then running - first the cert-generator
and then the MiniCluster pods:

```bash
$ kubectl get -n flux-operator pods
```
```console
NAME                         READY   STATUS              RESTARTS   AGE
flux-sample-0-f5znt          0/1     ContainerCreating   0          100s
flux-sample-1-th589          0/1     ContainerCreating   0          100s
flux-sample-cert-generator   0/1     Completed           0          100s
```

You can get the output file in the terminal (and don't worry about saving it too much, as it will save to storage).
By default cleanup is set to False so you shouldn't lose the pod to get the pods for it.

```bash
$ kubectl get -n flux-operator flux-sample-0-f5znt
```

Finally, note that with snakemake, once the output file in plots, called_reads and mapped_reads exist,
if you run it a second time, snakemake will determine there isn't anything to do.

## Clean up

Make sure you clean everything up!
Detach roles:

```bash
$ eksctl delete iamserviceaccount --name s3-mounter --namespace otomount --cluster flux-operator
```

Delete the flux operator in some way:

```bash
$ make undeploy
$ kubectl delete -f examples/dist/flux-operator-dev.yaml
$ kubectl delete -f examples/dist/flux-operator.yaml
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