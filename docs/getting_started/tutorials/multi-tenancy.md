# Multi-tenancy

"Multi-tenancy" means having multiple tenants, or users, in a cluster! We can accomplish this
fairly easily via [Flux Accounting](https://flux-framework.readthedocs.io/en/latest/guides/accounting-guide.html)
that will allow us to create users for our MiniCluster. This small tutorial will walk through 
creating an example with MiniKube.


## 1. Create Cluster

Bring up your MiniKube cluster.

```bash
$ minikube start
```

Then choose your step of choice to [install the operator](https://flux-framework.org/flux-operator/getting_started/user-guide.html#install)
You'll want a custom resource definition that turns on multi-user mode, and defines the users you
want to create, in addition to root (sets up the pods) and flux (owns the main flux instance). We
have a minicluster.yaml provided as an example:

```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MiniCluster
metadata:
  name: flux-sample
  namespace: flux-operator
spec:
  # suppress all output except for test run
  logging:
    quiet: false

  # Number of pods to create for MiniCluster
  size: 4

  # Define one or more users for your cluster
  users:
    - name: peenut
    - name: squidward
    - name: avocadosaurus

  containers:
    # This image has flux-accounting installed
    - image: ghcr.io/rse-ops/accounting:app-latest
    # Multi-user mode should not have a command
```

Note that we don't provide a command to run the Flux Restful API, and we use a pcontainer base that has
flux-accounting installed. Then, create the MiniCluster

```bash
$ kubectl create -f examples/tests/multi-tenant/minicluster.yaml
```

And then in the logs, we can see the different blocks of creating users (and passwords):

```console
Adding peenut with password 5784768d
Adding squidward with password cb21922b
Adding avocadosaurus with password 094cf1b6
...
🧾️ Creating flux accounting database
flux account create-db
flux account add-bank root 1
flux account add-bank --parent-bank=root user_bank 1
flux account add-user --username=peenut --bank=user_bank
flux account add-user --username=squidward --bank=user_bank
flux account add-user --username=avocadosaurus --bank=user_bank
```

And the Flux Restful Server is ready to go!

```console
🍓 Require auth: True
🍓    Flux user: ****
🍓   Flux token: ************************************
INFO:     Started server process [189]
INFO:     Waiting for application startup.
INFO:     Application startup complete.
INFO:     Uvicorn running on http://0.0.0.0:5000 (Press CTRL+C to quit)
```

Next we will add authentication mechanisms in the server to use the Linux users.
Stay tuned!