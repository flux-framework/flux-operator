# Bursting Experiment to EKS

> Experimental setup to burst to Amazon Web Services (Cloud)

This setup will expose a lead broker (index 0 of the MiniCluster job) as a service,
and then deploy a second cluster that can connect back to the first. 
For the overall design, see the top level [README](../README.md)

## AWS Elastic Kubernetes Service Setup

You should first [install eksctrl](https://github.com/weaveworks/eksctl) and make sure you have access to an AWS cloud. This means EITHER

### Environment Credentials

...with credentials or similar in your environment). E.g.,:

```bash
export AWS_ACCESS_KEY_ID=xxxxxxxxxxxxxxxxxxx
export AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export AWS_SESSION_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```
The last session token may not be required depending on your setup. OR

### AWS config

...with the aws config at `$HOME/.aws/config`.

**Important** whether you choose credentials or a config file (to map into the machine)
you *MUST* be consistent and ideally choose one! E.g., if you want to interact with the cluster in one context
but you created it with a different credential, you'll lead to bugs. This was tested
with [minicluster.yaml](minicluster.yaml) and not [minicluster-env.yaml](minicluster-env.yaml).
We also assume you also have [kubectl](https://kubernetes.io/docs/tasks/tools/). 

### Create Cluster

Next, let's create our cluster using eksctl "eks control." **IMPORTANT** you absolutely
need to choose a size that has [IsTrunkingCompatible](https://github.com/aws/amazon-vpc-resource-controller-k8s/blob/master/pkg/aws/vpc/limits.go)
true. Create the cluster as follows:

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

## Environment Credentials

**If you are using a config file, skip this step**

For this strategy, we will want to create a secret from our aws environment variable credentials to add to the MiniCluster.
You can [read about that here](https://kubernetes.io/docs/tasks/inject-data-application/distribute-credentials-secure/#define-container-environment-variables-using-secret-data).
First, base64 encode your variables doing the following:

```bash
$ echo -n $AWS_ACCESS_KEY_ID | base64
RkFLRUFXU0FDQ0VTU0tFWUlE

$ echo -n $AWS_SECRET_ACCESS_KEY | base64
RkFLRUFXU1NFQ1JFVEFDQ0VTU0tFWQ==
```

Create a secret YAML file that references these values:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: flux-operator-burst
type: Opaque
data:
  AWS_ACCESS_KEY_ID: RkFLRUFXU0FDQ0VTU0tFWUlE
  AWS_SECRET_ACCESS_KEY: RkFLRUFXU1NFQ1JFVEFDQ0VTU0tFWQ==
```

Apply to EKS cluster:

```bash
$ kubectl apply -f secret.yaml
```

## MiniCluster

Create the namespace, install the operator (assuming you are using a development version) and create the minicluster:

```bash
kubectl apply -f ../../../dist/flux-operator-dev.yaml
kubectl apply -f minicluster.yaml
# Expose broker pod port 8050 to 30093
kubectl apply -f service/broker-service.yaml
```

We need to open up the firewall to that port - this creates the rule for your security group:

```bash
# Get the group ID, and create the rule for it
SECURITY_GROUP_ID=$(aws eks describe-cluster --name flux-operator --query cluster.resourcesVpcConfig.clusterSecurityGroupId | jq -r)
aws ec2 authorize-security-group-ingress --region us-east-1 --group-id ${SECURITY_GROUP_ID} --protocol tcp --port 30093 --cidr "0.0.0.0/0"
```
```console
{
    "Return": true,
    "SecurityGroupRules": [
        {
            "SecurityGroupRuleId": "sgr-0afb218e929630d5a",
            "GroupId": "sg-07bacc27a95755641",
            "GroupOwnerId": "633731392008",
            "IsEgress": false,
            "IpProtocol": "tcp",
            "FromPort": 30093,
            "ToPort": 30093,
            "CidrIpv4": "0.0.0.0/0"
        }
    ]
}
```

Then figure out the node that the service is running from (we are interested in lead broker flux-sample-0-*)

```bash
$ kubectl get pods -o wide
NAME                  READY   STATUS    RESTARTS   AGE   IP              NODE                             NOMINATED NODE   READINESS GATES
flux-sample-0-4dk46   1/1     Running   0          53s   192.168.59.57   ip-192-168-48-220.ec2.internal   <none>           <none>
flux-sample-1-564n4   1/1     Running   0          53s   192.168.4.44    ip-192-168-27-240.ec2.internal   <none>           <none>
```

Then (using that node name) get the external ip for that node flux-sample-0-xx)

```bash
$ kubectl get nodes -o wide | grep ip-192-168-48-220.ec2.internal
ip-192-168-48-220.ec2.internal   Ready    <none>   34m   v1.23.17-eks-0a21954   192.168.48.220   44.211.52.254    Amazon Linux 2   5.4.242-156.349.amzn2.x86_64   docker://20.10.23
```

I set it to the environment to be useful later:

```bash
export LEAD_BROKER_HOST=44.211.52.254
```

Finally, when the broker index 0 pod is running, copy your scripts and configs over to it:

```bash
# This should be the index 0
POD=$(kubectl get pods -o json | jq -r .items[0].metadata.name)

# This will copy configs / create directories for it
kubectl cp ./run-burst.py ${POD}:/tmp/workflow/run-burst.py -c flux-sample
kubectl exec -it ${POD} -- mkdir -p /tmp/workflow/external-config
kubectl cp ../../../dist/flux-operator-dev.yaml ${POD}:/tmp/workflow/external-config/flux-operator-dev.yaml -c flux-sample
```

## Config File Credentials

**If you are using environment credentials, skip this step**

As an alternative to environment variables, you can copy over your config. You should only do this if you are sure
nobody else can access the machine, and it's a development context.

```bash
kubectl exec -it ${POD} -- mkdir -p /home/flux/.aws
kubectl cp $HOME/.aws/config ${POD}:/home/flux/.aws/config -c flux-sample
```

**Do not do this for a production cluster!**

## Burstable Job

Now let's create a job that cannot be run because we don't have the resources. The `flux-burst` Python module, using it's simple
default, will just look for jobs with `burstable=True` and then look for a place to assign them to burst. Since this is a plugin
framework, in the future we can implement more intelligent algorithms for either filtering the queue (e.g., "Which jobs need bursting?"
and then determining if a burst can be scheduled for some given burst plugin (e.g., GKE)). For this simple setup and example,
we ensure the job doesn't run locally because we've asked for more nodes than we have. Shell into your broker pod:

```bash
$ kubectl exec -it ${POD} bash
```

If you used environment credentials, double check your AWS credentials are there!

```bash
$ env | grep AWS
```

Connect to the broker socket (and note if you can't yet, the nodes are probably still installing stuff):

```bash
source /mnt/flux/flux-view.sh
flux proxy $fluxsocket bash
```

The libraries we need should be installed in the minicluster.yaml.
You might want to add others for development (e.g., IPython).
Resources we have available?

```bash
$ flux resource list
     STATE NNODES   NCORES NODELIST
      free      2        8 flux-sample-[0-1]
 allocated      0        0 
      down      6       24 flux-sample-[2-3],burst-0-[0-3]
```

The above shows us that the broker running here can accept burstable resources (`burst-0-[0-3]`), and even
can accept the local cluster expanding (`flux-sample[2-3]`) for a total of 24 cores. The reason
that the remote burst prefix has an extra "0" is that we could potentially have different sets of
burstable remotes, namespaced by this prefix. And now let's create a burstable job, and ask for more nodes than we have :)

```bash
# Set burstable=1
# this will be for 4 nodes, 8 cores each
$ flux submit -N 4 -o cpu-affinity=off --cwd /tmp --setattr=burstable hostname
```

You should see it's scheduled (but not running). Note that if we asked for a resource totally unknown
to the cluster (e.g. 4 nodes and 32 tasks) it would just fail. Note that because of this,
we need in our "mark as burstable" method a way to tell Flux not to fail in this case.
You can see it is scheduled and waiting for resources:

```bash
$ flux jobs -a
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
   ƒQURAmBXV fluxuser hostname    S      8      8        - 
```
```bash
$ flux job attach $(flux job last)
flux-job: ƒQURAmBXV waiting for resources  
```

Get a variant of the munge key we can see (it's owned by root so this ensures we can see/own it as the flux user)

```bash
cp /etc/munge/munge.key ./munge.key
```

Now we can run our script to find the jobs based on this attribute!

```bash
# This is the address of the lead host we discovered above
LEAD_HOST="54.90.176.102"

# This is the node port we've exposed on the cluster
LEAD_PORT=30093
python3 run-burst.py --flux-operator-yaml ./external-config/flux-operator-dev.yaml \
        --lead-host ${LEAD_HOST} --lead-port ${LEAD_PORT} --lead-size 4 \
        --munge-key ./munge.key --name burst-0
```

Important notes for the above:

- The curve path and secret name have defaults set.
- The name is same that would be automatically generated name by Flux given a bursted cluster (that isn't explicitly given a name) but we are being pedantic. It's also in the [minicluster.yaml](minicluster.yaml)
- The lead name is derived from the hostname where it is running (e.g., flux-sample) so we don't need to provide it
- We set the lead size to the max size, because the ranks indices need to line up. We are using a size that won't fail the job (which needs 4)

When you do the above (and the second MiniCluster launches) you should be able to see on your local cluster the external
MiniCluster resources, and the result of hostname will include the external hosts! Here is how to shell into the cluster
from another terminal:

```bash
POD=$(kubectl get pods -o json | jq -r .items[0].metadata.name)
kubectl exec -it ${POD} bash
source /mnt/flux/flux-view.sh
flux proxy $fluxsocket bash
```

```bash
$ flux resource list
     STATE NNODES   NCORES NODELIST
      free      6       24 flux-sample-[0-1],burst-0-[0-3]
 allocated      0        0 
      down      2        8 flux-sample-[2-3]
```

Note that flux sample 2-3 have been left if we wanted to expand the local cluster, just as an example.
Also notice that (along with the burst resources being online), our job has run:

```bash
$ flux jobs -a 
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
   ƒ2XwYQ37M flux     hostname   CD      4      4   0.049s flux-sample-[0-1],burst-0-[2-3]
```

And we can see output! Note that the error is because the working directory where it was launched doesn't exist on the remote.

```bash
$ flux job attach ƒ2XwYQ37M
flux-sample-1
burst-0-3
burst-0-1
burst-0-2
```

You can also launch a new job that asks to hit all the nodes (6):

```bash
$ flux run -o cpu-affinity=off --cwd /tmp  -N 6 hostname
flux-sample-0
flux-sample-1
burst-0-0
burst-0-1
burst-0-2
burst-0-3
```

Note that if you don't disable it, there will be error messages about affinity. 

```bash
0.046s: flux-shell[5]: ERROR: cpu-affinity: affinity: core1 not in topology
```

The same is true for `cwd` (current working directory) - if you don't specify it, it defaults to the present working directory,
`/tmp/workflow` that doesn't exist on the bursted cluster.
And that's bursting! At this point we should think about how to better start / stop a burst, since the cluster
will typically come up (and stay up). For debugging, see [broker-gke](../broker-gke).

### Cleanup

Note that you'll only see the exposure with the kind docker container with `docker ps`.
When you are done, clean up

```bash
kubectl delete -f minicluster.yaml
kubectl delete -f service.yaml
eksctl delete cluster -f eksctl-config.yaml
```
