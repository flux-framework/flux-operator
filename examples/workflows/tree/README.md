# Basic Tree

We can use [flux tree](https://github.com/flux-framework/flux-sched/blob/master/t/t2001-tree-real.t#L43-L51)
to create instances inside of instances. For this basic example, we will start with a root, create
two instances under it, and two instances under each of those. We can describe the shape of this topology as 2x2.
You can read more about [the utility here](https://github.com/flux-framework/flux-sched/blob/master/resource/utilities/README.md).

## Usage

First, let's create a kind cluster. From the context of this directory:

```bash
$ kind create cluster --config ../../kind-config.yaml
```

And then install the operator, create the namespace, and apply the MiniCluster YAML here.

```bash
$ kubectl apply -f ../../dist/flux-operator.yaml
$ kubectl create namespace flux-operator
$ kubectl apply -f ./minicluster.yaml
```

The cluster creation has the present working directory (where you are reading this file)
bound to `/tmp/workflow`, and we are running the `flux tree` command there. You can check the logs
for the run via:

```bash
$ kubectl logs -n flux-operator flux-sample-0-7tx7s -f
```

And when it's done, the tree.out (written to `/tmp/workflow` in the cluster) will be written to `tree.out`.
In here you will see:

```bash
$ flux tree -T2x2 -J 4 -N 1 -c 4 -o /tmp/workflow/tree.out -Q easy:fcfs hostname 
```
```console
TreeID                  Elapsed(sec)         Begin(Epoch)           End(Epoch)     Match(usec)           NJobs NNodes  CPN  GPN
tree                        3.580370    1682057760.514143    1682057764.094512        0.000000               4     1    4    0
tree.2                      2.022420    1682057761.434294    1682057763.456710        0.000000               2     1    2    0
tree.2.2                    0.154145    1682057762.485498    1682057762.639643        0.000000               1     1    1    0
tree.2.1                    0.168346    1682057762.445322    1682057762.613668        0.000000               1     1    1    0
tree.1                      2.081930    1682057761.270239    1682057763.352172        0.000000               2     1    2    0
tree.1.2                    0.136858    1682057762.364354    1682057762.501212        0.000000               1     1    1    0
tree.1.1                    0.126923    1682057762.225927    1682057762.352850        0.000000               1     1    1    0
```

What is happening is that the `flux tree` command is creating a hierarchy of instances. Based on their names you can tell that:

 - `2x2` in the command is the topology
 - It says to create two flux instances, and make them each spawn two more.
 - `tree` is the root
 - `tree.1` is the first instance
 - `tree.2` is the second instance
 - `tree.1.1` and `tree.1.2` refer to the nested instances under `tree.1`
 - `tree.2.1` and `tree.2.2` refer to the nested instances under `tree.2`
 
And we provided the command `hostname` to this script, but a more complex example would generate more interested hierarchies,
and with different functionality for each. When you are done, please clean up:

```bash
$ kubectl delete -f minicluster.yaml
```
