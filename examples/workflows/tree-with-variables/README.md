# Tree with Variables

We can use [flux tree](https://github.com/flux-framework/flux-sched/blob/master/t/t2001-tree-real.t#L43-L51)
to create instances inside of instances. For this example, we will start with a root, create
two instances under it, and two instances under each of those. We will (instead of running hostname) run
a script that demonstrates the environment available to each subinstance.
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
$ flux tree -T2x2 -J 4 -N 4 -c 4 -o /tmp/workflow/tree.out -Q easy:fcfs /bin/bash ./run-on-instance.sh
```
```console
$ cat tree.out 
TreeID                  Elapsed(sec)         Begin(Epoch)           End(Epoch)     Match(usec)           NJobs NNodes  CPN  GPN
tree                        3.646440    1682094481.024492    1682094484.670933        0.000000               4     4    4    0
tree.2                      1.847760    1682094482.167398    1682094484.015160        0.000000               2     2    4    0
tree.2.2                    0.146933    1682094483.195491    1682094483.342424        0.000000               1     1    4    0
tree.2.1                    0.098842    1682094483.068877    1682094483.167719        0.000000               1     1    4    0
tree.1                      1.789910    1682094482.071364    1682094483.861272        0.000000               2     2    4    0
tree.1.2                    0.102510    1682094483.056029    1682094483.158540        0.000000               1     1    4    0
tree.1.1                    0.119904    1682094482.937050    1682094483.056954        0.000000               1     1    4    0
```

This information is repeated from the [basic tree](../tree) example, and you can look there for details about what the above means.
For this example, we focus on the variables available in the script, and we write files that are named by the tree id! You
should be able to see them in the present working directory:

```bash
$ ls
```
```console
minicluster.yaml  README.md  run-on-instance.sh  tree.1.1-output.txt  tree.1.2-output.txt  tree.2.1-output.txt  tree.2.2-output.txt  tree.out
```

If we look in a script we can see the variables available to the instance:

```bash
$ cat tree.1.2-output.txt 
```
```console
FLUX_TREE_ID tree.1.2
FLUX_TREE_JOBSCRIPT_INDEX 1
FLUX_TREE_NNODES 1
FLUX_TREE_NCORES_PER_NODE 1
FLUX_TREE_NGPUS_PER_NODE 0
```

Note that for this example we are only running the scripts on the leaves, hence why we only see one `NNODES` above. The table above
that shows we go from `4 > 2 > 1`. You would direct custom logic in this little script to control execution of your job, likely with different instances using different resources.
It's super cool!

```bash
$ kubectl delete -f minicluster.yaml
```
