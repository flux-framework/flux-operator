# IBM Cloud

For this tutorial you will need to create or have access to an account on [IBM Cloud](https://cloud.ibm.com/).
We are going to use the [IBM Kubernetes Service](https://cloud.ibm.com/docs/containers?topic=containers-getting-started).

## Setup

### Credentials

IBM works by way of [access keys](https://cloud.ibm.com/docs/account?topic=account-userapikey&interface=ui#userapikey) for the API.
If you click on your profile in the top right of the console, you can click "Log in to CLI and API" and a box will pop up.
You will want to copy paste the first command that logs in as follows:

```bash
$ ibmcloud login -a https://cloud.ibm.com -u passcode -p xxxxxxx
```

And then choose a zone (e.g., us-east). 

### Install

You'll first need to [install the ibmcloud](https://cloud.ibm.com/docs/cli?topic=cli-getting-started) command line
client.

```bash
$ curl -fsSL https://clis.cloud.ibm.com/install/linux > install.sh
$ chmod +x install.sh
$ ./install.sh
```

This should create the `ibmcloud` executable on your path.

```bash
$ which ibmcloud
```
```console
/usr/local/bin/ibmcloud
```

Ensure the [storage plugin](https://cloud.ibm.com/docs/cloud-object-storage?topic=cloud-object-storage-cli-plugin-ic-cos-cli) is installed
along with the [kubernetes plugin](https://cloud.ibm.com/docs/containers?topic=containers-cs_cli_install):

```bash
$ ibmcloud plugin install cloud-object-storage
$ ibmcloud plugin install container-service

# And view your current config
$ ibmcloud cos config list
```

If you are doing this for the first time, you'll notice the CRN is blank. We will create a CRN when we push objects to the bucket.
In addition to this client, you'll need [aws](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) and [kubectl](https://kubernetes.io/docs/tasks/tools/).


## Snakemake (requiring storage) on IBM Cloud

Akin to how we created a local volume, we can do something similar, but instead of pointing the Flux Operator
to a volume on the host (e.g., in MiniKube) we are going to point it to a storage bucket with our data.
IBM Cloud has an s3-like storage option we can use.

### Prepare Data

To start, prepare your data in a temporary directory (that we will upload into IBM cloud storage):

<details>

<summary>Instructions for preparing Snakemake data</summary>

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

Delete GitHub/git assets:

```bash
$ rm -rf .git .github
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

</details>

Let's first [create a bucket](https://cloud.ibm.com/docs/cloud-object-storage?topic=cloud-object-storage-cli-plugin-ic-cos-cli#ic-create-bucket),
and this will use [this plugin](https://github.com/IBM/ibmcloud-cos-cli).
Following [these instructions](https://ibm.github.io/kubernetes-storage/Lab5/cos-with-s3fs/COS/).
The first step is to grant a service authorization in the cloud console, which looks like this:

![img/ibm-storage.png](img/ibm-storage.png)

<details>

<summary>Creating the IBM Cloud Storage</summary>

We first need to create a service account:

```bash
# This is usually "Default" - try doing "ibmcloud resource groups"
RESOURCEGROUP=$(ibmcloud resource groups --output json | jq -r '.[0].name') 
COS_NAME_RANDOM=$(date | md5sum | head -c10)
COS_NAME=$COS_NAME_RANDOM-cos-1
COS_CREDENTIALS=$COS_NAME-credentials
COS_PLAN=Lite
COS_BUCKET_NAME=flux-operator-storage
REGION=us-east
COS_PRIVATE_ENDPOINT=s3.private.$REGION.cloud-object-storage.appdomain.cloud
```

And now we can make the service instance:

```bash
$ ibmcloud resource service-instance-create $COS_NAME cloud-object-storage $COS_PLAN global -g $RESOURCEGROUP
```

And list to ensure it was created (repeats the output above for the most part, but can be run separately or
after the fact):

```bash
$ ibmcloud resource service-instance $COS_NAME
```

Next, set the GUID of the object storage instance,

```bash
COS_GUID=$(ibmcloud resource service-instance $COS_NAME --output json | jq -r '.[0].guid')
echo $COS_GUID
```

And add credentials so you can authenticate with IAM. 

```bash
$ ibmcloud resource service-key-create $COS_CREDENTIALS Writer --instance-name $COS_NAME --parameters '{"HMAC":true}'
```

And generate the API key in json:

```bash
COS_APIKEY=$(ibmcloud resource service-key $COS_CREDENTIALS --output json | jq -r '.[0].credentials.apikey')
echo $COS_APIKEY
```

Now we can (finally) create a bucket, first getting the CRN:

```bash
COS_CRN=$(ibmcloud resource service-key $COS_CREDENTIALS --output json | jq -r '.[0].credentials.resource_instance_id')
echo $COS_CRN
```
Add the CRN:

```bash
$ ibmcloud cos config crn --crn $COS_CRN
```

Verify that it shows up:

```bash
$ ibmcloud cos config list
```

And then create the bucket:

```bash
$ ibmcloud cos bucket-create --bucket $COS_BUCKET_NAME --region $REGION
```

And verify it was created:

```bash
$ ibmcloud cos list-buckets --ibm-service-instance-id $COS_CRN
```

</details>

When you have your storage, you then upload the workflow using the `aws` client. We will first need to make a service
credential. I did this by clicking the top left hamburger menu, and then "[Resource List](https://cloud.ibm.com/resources)" and clicking the arrow
to expand storage, clicking on the instance ID, and then I saw my bucket! Note that for small data, you can click "Upload" on the right
side and select a folder. But likely you want to do it from the command line - you can find your AWS acces token and secret in this json
payload under credentials->cos_hmac_keys:

```bash
$ ibmcloud resource service-key $COS_CREDENTIALS --output json
```

Export them to the environment.

```bash
export AWS_ACCESS_KEY_ID=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

And then you should be able to list your (currently empty) bucket:

```bash
COS_BUCKET_NAME=flux-operator-storage
REGION=us-east
```
```bash
$ aws s3 ls $COS_BUCKET_NAME --endpoint-url https://s3.${REGION}.cloud-object-storage.appdomain.cloud
```

And then

```bash
$ aws s3 cp --recursive /tmp/workflow/ s3://${COS_BUCKET_NAME}/snakemake-workflow --endpoint-url https://s3.${REGION}.cloud-object-storage.appdomain.cloud
```

Do a listing again to ensure everything is there!

```bash
$ aws s3 ls s3://${COS_BUCKET_NAME}/snakemake-workflow/ --endpoint-url https://s3.${REGION}.cloud-object-storage.appdomain.cloud
```

You can also see the created files in the web interface:

![img/ibm-storage-create.png](img/ibm-storage-create.png)

### Create Cluster

Next let's create a cluster called "flux-operator," following [these instructions](https://www.kubeflow.org/docs/distributions/ibm/create-cluster/).

```bash
# Do ibmcloud ks versions to see versions available
# I left unset to use the default
# export KUBERNETES_VERSION=1.26

# ibmcloud ks locations
export CLUSTER_ZONE=dal12
export WORKER_NODE_PROVIDER=classic
export CLUSTER_NAME=flux-operator
```

Choose a worker node instance type:

```bash
$ ibmcloud ks flavors --zone dal12 --provider classic
egxport WORKER_NODE_FLAVOR="u3c.2x4"
```

And then create the cluster!

```bash
$ ibmcloud ks cluster create ${WORKER_NODE_PROVIDER} \
  --name=$CLUSTER_NAME \
  --zone=$CLUSTER_ZONE \
  --flavor ${WORKER_NODE_FLAVOR} \
  --workers=4 
```

If you get an error about providing a VLAN, then do:


```bash
$ ibmcloud ks vlans --zone ${CLUSTER_ZONE}
```

and set the public and private VLAN ids to environment variables `PUBLIC_VLAN_ID` and `PRIVATE_VLAN_ID`

```bash
$ ibmcloud ks cluster create ${WORKER_NODE_PROVIDER} \
  --name=$CLUSTER_NAME \
  --zone=$CLUSTER_ZONE \
  --flavor ${WORKER_NODE_FLAVOR} \
  --workers=4 \
  --private-vlan ${PRIVATE_VLAN_ID} \
  --public-vlan ${PUBLIC_VLAN_ID} 
```

The instructions in the linked tutorial mention to use a command line tool to check status, but this didn't work for me.
Instead I watched until the cluster was ready on the [IBM Cloud Kubernetes](https://cloud.ibm.com/kubernetes/clusters) page.
Once it's ready, you should be able to list:

```bash
$ ibmcloud ks cluster ls
```
```console
OK
Name            ID                     State    Created          Workers   Location   Version       Resource Group Name   Provider
flux-operator   cg54590d0b7conrbfmtg   normal   25 minutes ago   4         Dallas     1.25.6_1534   Default               classic
```

And then switch the kubernetes context to use it:

```bash
$ ibmcloud ks cluster config --cluster ${CLUSTER_NAME}
```
```console
The configuration for flux-operator was downloaded successfully.

Added context for flux-operator to the current kubeconfig file.
You can now execute 'kubectl' commands against your cluster. For example, run 'kubectl get nodes'.
```

Make sure all worker nodes are up with the command below

```bash
$ kubectl get nodes
```

and make sure all the nodes are in Ready state. Next we can move on to installing the operator!


### Deploy Operator

To deploy the Flux Operator, [choose one of the options here](https://flux-framework.org/flux-operator/getting_started/user-guide.html#production-install) to deploy the operator. Whether you apply a yaml file, use [flux-cloud](https://converged-computing.github.io/flux-cloud) or clone the repository and `make deploy` you will see the operator install to the `operator-system` namespace.

For a quick "production deploy" from development, the Makefile has a directive that will build and push a `test` tag (you'll need to edit `DEVIMG` to be one you can push to) and then generate a
yaml file targeting that image, e.g.,

```bash
$ make test-deploy
$ kubectl apply -f examples/dist/flux-operator-dev.yaml
```

Ensure the `operator-system` namespace was created:

```bash
$ kubectl get namespace
```
```console
NAME              STATUS   AGE
default           Active   28m
ibm-cert-store    Active   7m29s
ibm-operators     Active   27m
ibm-system        Active   27m
kube-node-lease   Active   28m
kube-public       Active   28m
kube-system       Active   28m
operator-system   Active   11s
```

And you can find the name of the operator pod as follows:

```bash
$ kubectl get pod --all-namespaces
```
```console
operator-system   operator-controller-manager-56b5bcf9fd-m8wg4               2/2     Running   0          73s
```

And wait until that is Running.

### Create Flux Operator namespace

Make your namespace for the flux-operator custom resource definition (CRD), which is the yaml file above that generates the MiniCluster:

```bash
$ kubectl create namespace flux-operator
```

### Install the Container Storage Plugin

We did a combined approach of using helm and our own YAML.

<details>

<summary>Instructions for Generation of YAML</summary>

If you do want to use helm to reproduce this, you will need to [install helm](https://helm.sh/docs/intro/install/). First,
add the [IBM repository](https://github.com/IBM/charts/tree/5870731bfebc867a47fec79507f2f7f616688e25/stable/ibm-object-storage-plugin) to helm:

```bash
$ helm repo add ibm-helm https://raw.githubusercontent.com/IBM/charts/master/repo/ibm-helm
```
Then you'll want to update the helm IBM repository:

```bash
$ helm repo update
$ helm plugin remove ibmc || echo "No ibmc plugin installed"
```

And pull and install the latest:

```bash
$ helm fetch --untar ibm-helm/ibm-object-storage-plugin
$ helm plugin install ./ibm-object-storage-plugin/helm-ibmc
```

Ensure you set the right permissions:
```bash
$ chmod 755 /home/${USER}/.local/share/helm/plugins/helm-ibmc/ibmc.sh
```
And then berify installation:

```bash
$ helm ibmc --help
```

The output version of the Helm version needs to be >3.0. Then we can inspect
variables:

```bash
$ helm show values ibm-helm/ibm-object-storage-plugin
```

```console
replicaCount: 1
maxUnavailableNodeCount: 1

# Change license to true to indicate have read and agreed to license agreement
# https://www.apache.org/licenses/LICENSE-2.0
license: false

image:
  providerImage:
    # This image is required only for IBM Cloud clusters
    ibmContainerRegistry: icr.io/ibm/ibmcloud-provider-storage:1.30.6
  pluginImage:
    ibmContainerRegistry: icr.io/ibm/ibmcloud-object-storage-plugin@sha256:ce509caa7a47c3329cb4d854e0b3763081ac725901e40d5e57fe93b6cd125243
    publicRegistry: icr.io/cpopen/ibmcloud-object-storage-plugin@sha256:ce509caa7a47c3329cb4d854e0b3763081ac725901e40d5e57fe93b6cd125243
  driverImage:
    ibmContainerRegistry: icr.io/ibm/ibmcloud-object-storage-driver@sha256:8c91974660bf98efc772f369b828b3dbcea5d6829cbd85e6321884f4c4eabe09
    publicRegistry: icr.io/cpopen/ibmcloud-object-storage-driver@sha256:8c91974660bf98efc772f369b828b3dbcea5d6829cbd85e6321884f4c4eabe09
  pullPolicy: Always

# IAM endpoint url
iamEndpoint: https://iam.cloud.ibm.com
iamEndpointVPC: https://private.iam.cloud.ibm.com

# IBMC || IBMC-VPC || RHOCP || SATELLITE
provider: RHOCP

# Container platform [ K8S vs OpenShift ]
platform: OpenShift

# Datacenter name where cluster is deployed (required only for IKS)
dcname: ""
region: ""
# Worker node's OS [ redhat || debian ]
workerOS: redhat

# COS endpoints and COS storageClass configuration
# For satellite clusters, to use aws, wasabi or ibm object storage, please provide the s3Provider(aws/wasabi/ibm) and respective storageClass(region) values.
# If user provides all 3 values endpoint, storageClass and s3Provider, precedence is given to storageClass and s3Provider.
# To input endpoint explicitly, input only endpoint and storageClass and skip s3Provider.
# For non-satellite rhocp clusters, please provide cos endpoint and cos storageclass.

cos:
  # The s3 endpoint url for the targeted object storage service; format - https://<Endpoint URL>
  endpoint: "NA"
  # The region in which the bucket has to be created as per the object storage service provider, ex - us-south
  storageClass: ""
  # Supported object storage service providers are aws, ibm, wasabi
  s3Provider: ""

secondaryValueFile: ibm/values.yaml
secondaryValueFileSat: satellite/values.yaml

arch: amd64

resource:
  memory: 500Mi
  cpu: 500m
  ephemeralStorageReq: 5Mi
  ephemeralStorageLimit: 105Mi

# /etc/kubernetes for RHOCP else /usr/libexec/kubernetes
kubeDriver: /usr/libexec/kubernetes

bucketAccessPolicy: false

quotaLimit: false

allowCrossNsSecret: true

# current and previous
s3fsVersion: current
```

Of the values we can set, we likely want to set the following:

```diff
# Change license to true to indicate have read and agreed to license agreement
# https://www.apache.org/licenses/LICENSE-2.0
- license: false
+ license: true

# IBMC || IBMC-VPC || RHOCP || SATELLITE
- provider: RHOCP
+ provider: IBMC

# Container platform [ K8S vs OpenShift ]
- platform: OpenShift
+ platform: K8S

# Worker node's OS [ redhat || debian ]
- workerOS: redhat
+ workerOS: debian

# COS endpoints and COS storageClass configuration
# For satellite clusters, to use aws, wasabi or ibm object storage, please provide the s3Provider(aws/wasabi/ibm) and respective storageClass(region) values.
# If user provides all 3 values endpoint, storageClass and s3Provider, precedence is given to storageClass and s3Provider.
# To input endpoint explicitly, input only endpoint and storageClass and skip s3Provider.
# For non-satellite rhocp clusters, please provide cos endpoint and cos storageclass.

cos:
  # The s3 endpoint url for the targeted object storage service; format - https://<Endpoint URL>
-  endpoint: "NA"
+ endpoint: "https://s3.us-east.cloud-object-storage.appdomain.cloud/"
  # Supported object storage service providers are aws, ibm, wasabi
-  s3Provider: ""
+  s3Provider: "ibm"
```
Using the above, we would install the plugin directly with our customizations:

```bash
$ helm ibmc install ibm-object-storage-plugin ibm-helm/ibm-object-storage-plugin --set license=true --set provider=IBMC --set platform=K8s --set workerOS=debian --set cos.s3Provider=ibm
```

Give a minute or so, and then check that the storage classes were created correctly:

```bash
$ kubectl get storageclass | grep 'ibmc-s3fs'
```

Then make sure the plugin pods are running (in "RUNNING" state):

```bash
$ kubectl get pods -n ibm-object-s3fs -o wide | grep object
```

Ensure the secret exists in the default namespace:

```bash
$ kubectl get secrets -n default | grep icr-io
```

We will want to copy the secret to ibm-object-s3fs namespace

```bash
$ kubectl get secret -n default all-icr-io -o yaml | sed 's/default/ibm-object-s3fs/g' | kubectl -n ibm-object-s3fs create -f -
```

And then make sure that the image pull secret is available in the ibm-object-s3fs  namespace.

```bash
$ kubectl get secrets -n ibm-object-s3fs | grep icr-io
```

Verify that the state of the plugin pods changes to "Running".

```bash
$ kubectl get pods -n ibm-object-s3fs | grep object
```

Finally, create a secret with our service ID to give access to storage. To get your API key:

```bash
$ ibmcloud resource service-key $COS_CREDENTIALS --output json
```

If you are doing the full tutorial, this is the variable under `$COS_APIKEY` and the 
the service instance id should be `$COS_NAME`.

```bash
$ kubectl create secret generic s3-secret --namespace=flux-operator --type=ibm/ibmc-s3fs  --from-literal=api-key=${COS_APIKEY} --from-literal=service-instance-id=${COS_NAME}
```

Finally, we need to create a storage class that points to the correct credentials
and namespace. Note this file names it "ibm-s3-storage":

```bash
$ kubectl apply -f examples/storage/ibm/storageclass.yaml 
```

At this point we have the storage driver running, along with the storage class and secret, and we should
attempt to use it with the Flux Operator.

### Snakemake MiniCluster

Note that I found all these annotation options [here](https://github.com/IBM/ibmcloud-object-storage-plugin/blob/455032aea2f820b8b3ad927e9af1eef6942dc2d5/provisioner/ibm-s3fs-provisioner_test.go#L73).
Also note that we are setting the `commands: -> runFluxAsRoot` to true. This isn't ideal, but it was the
only way I could get the storage to both be seen and have permission to write there. Let's create the job!
Since the storage plugin uses a FlexDriver, and this is being deprecated, we need to create the persistent
volume manually. It will be discovered and used by the Flux Operator. I tried adding the
owner reference "uid" to this file first:

```bash
$  kubectl get -n operator-system pod operator-controller-manager-858c9ccfb4-2k79n -o yaml 
#    uid: 87b73461-b5fc-4746-a959-e84518096ed4
```

Although it didn't seem to make a difference - a second PV (or PVC?) was always made.
Next make it:

```yaml
$ kubectl create -f examples/storage/ibm/pv.yaml
```

And then right after, create the MiniCluster

```bash
$ kubectl apply -f examples/storage/ibm/minicluster.yaml
```

The pods will take a bit to pull the containers, in the meantime you can check out the pv and pvc:

STOPPED HERE - the pvc seems to create a second PV but then the pod is in pending because the "data" one
we created (which wasn't used) is the one available.:

```console
$ kubectl get -n flux-operator pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM                      STORAGECLASS     REASON   AGE
data                                       25Gi       RWX            Delete           Available   flux-operator/data-claim   ibm-s3-storage            3m24s
pvc-0400f3cc-2662-4cb2-a83b-3bda9e0ea0be   25Gi       RWX            Delete           Pending     flux-operator/data-claim   ibm-s3-storage            113s
```
Why was the second created?

```console
$ kubectl get -n flux-operator pvc
NAME         STATUS   VOLUME   CAPACITY   ACCESS MODES   STORAGECLASS     AGE
data-claim   Lost                                        ibm-s3-storage   5m55s
```
And then the data claim is lost...

```console
$ kubectl describe -n flux-operator pvc
Name:          data-claim
Namespace:     flux-operator
StorageClass:  ibm-s3-storage
Status:        Lost
Volume:        
Labels:        <none>
Annotations:   ibm.io/auto-create-bucket: false
               ibm.io/auto-delete-bucket: false
               ibm.io/bucket: flux-operator-storage
               ibm.io/chunk-size-mb: 40
               ibm.io/curl-debug: false
               ibm.io/debug-level: warn
               ibm.io/iam-endpoint: https://iam.cloud.ibm.com
               ibm.io/kernel-cache: true
               ibm.io/multireq-max: 20
               ibm.io/object-store-endpoint: https://s3.direct.us-east.cloud-object-storage.appdomain.cloud
               ibm.io/object-store-storage-class: us-east-standard
               ibm.io/parallel-count: 20
               ibm.io/s3fs-fuse-retry-count: 5
               ibm.io/secret-name: s3-secret
               ibm.io/secret-namespace: flux-operator
               ibm.io/stat-cache-size: 100000
               pv.kubernetes.io/bind-completed: yes
               pv.kubernetes.io/bound-by-controller: yes
               volume.beta.kubernetes.io/storage-provisioner: ibm.io/ibmc-s3fs
               volume.kubernetes.io/storage-provisioner: ibm.io/ibmc-s3fs
Finalizers:    [kubernetes.io/pvc-protection]
Capacity:      
Access Modes:  
VolumeMode:    Filesystem
Used By:       flux-sample-0-wbqft
               flux-sample-1-rcf97
Events:
  Type     Reason                 Age   From                                                                                                  Message
  ----     ------                 ----  ----                                                                                                  -------
  Normal   Provisioning           6m4s  ibm.io/ibmc-s3fs_ibmcloud-object-storage-plugin-bd89679b7-lgmx9_55bf4496-14e8-47dd-999e-7e9398a9bfd8  External provisioner is provisioning volume for claim "flux-operator/data-claim"
  Warning  ClaimLost              6m3s  persistentvolume-controller                                                                           Bound claim has lost reference to PersistentVolume. Data on the volume is lost!
  Normal   ProvisioningSucceeded  6m3s  ibm.io/ibmc-s3fs_ibmcloud-object-storage-plugin-bd89679b7-lgmx9_55bf4496-14e8-47dd-999e-7e9398a9bfd8  Successfully provisioned volume pvc-0400f3cc-2662-4cb2-a83b-3bda9e0ea0be
(env) (base) vanessa@vanessa-ThinkPad-T490s:~/Desktop/Code/flux/operator$ kubectl get -n flux-operator pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM                      STORAGECLASS     REASON   AGE
data                                       25Gi       RWX            Delete           Available   flux-operator/data-claim   ibm-s3-storage            8m39s
pvc-0400f3cc-2662-4cb2-a83b-3bda9e0ea0be   25Gi       RWX            Delete           Pending     flux-operator/data-claim   ibm-s3-storage            7m8s
```
But reports using the second (pending) one?

## Clean up

To clean up:

```bash
# Delete the cluster
$ ibmcloud ks cluster rm --force-delete-storage -c ${CLUSTER_NAME}

# Delete storage
$ ibmcloud cos bucket-delete --bucket $name_bucket
```
