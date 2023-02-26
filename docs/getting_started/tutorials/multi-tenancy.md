# Multi-tenancy


<div class="result docutils container">
<div class="warning admonition">
<p class="admonition-title">Warning</p>
    <p>Multi-tenancy is early in development and should not be used in production.
    We create a multi-user Flux install, a set of users, and then authenticate
    to submit jobs via PAM, meaning that no credential is stored anywhere.
    A future approach might also include a secret for each user so the payloads
    are encrypted in transport.</p>
</div>
</div>

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


## 2. Custom Resource Definition

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
  # These users will be created, and user/pass required for auth
  # passwords are optional, if not provided will be generated
  users:
    - name: peenut
    - name: squidward
    - name: avocadosaurus

  containers:
    # This image has flux-accounting installed
    - image: ghcr.io/rse-ops/accounting:app-latest
    # Multi-user mode should not have a command
```

Note that we don't provide a command to run the Flux Restful API, and we use a container base that has
flux-accounting installed. Then, create the flux-operator namespace and MiniCluster:

```bash
$ kubectl create namespace flux-operator
$ kubectl create -f examples/flux-restful/minicluster-multi-tenant.yaml
```

## 3. View Logs

And then in the logs, we can see the different blocks of creating users (and passwords):

```bash
# get pods
$ kubectl get -n flux-operator pods

# get logs for broker pod
$ kubectl logs -n flux-operator flux-sample-0-fz64b
```

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

And the Flux Restful Server is ready to go! Note that pam is enabled,
and we don't have a Flux user or flux token provided. When you authenticate with PAM
you are not allowed to have a global Flux user, as it wouldn't coincide with an actual
user account.

```console
broker.err[0]: accepting connection from flux-sample-3 (rank 3) status full
🍓 Require auth: True
🍓     PAM auth: True
🍓    Flux user: unset
🍓   Flux token: unset
INFO:     Started server process [170]
INFO:     Waiting for application startup.
INFO:     Application startup complete.
INFO:     Uvicorn running on http://0.0.0.0:5000 (Press CTRL+C to quit)
```

Make sure that you grab a user name and password from the terminal.

## 4. Expose Web Interface

Next we will want to expose the service.

```bash
$ kubectl port-forward -n flux-operator flux-sample-0-zdhkp 5000:5000
```

And open your browser to [http://localhost:5000](http://localhost:5000). You
should be able to login with a user and password. Note that nothing is stored
in any kind of database for Flux Restful API - all authentication is done
against the server using PAM. If you want more programmattic access, look at 
our Python SDK examples for "port forwarding" that will do the same above, 
but not require the manual extra work (all is done in Python).

## 5. Admin Actions

### Adding a User

It is possible to dynamically add a user! You will need to create them, add them
to the accounting database, and then updating. You can exec a set of commands
to the broker pod as follows:

```bash
# Get the broker pod kubectl get -n flux-operator pods
pod="flux-sample-0-fz64b"

# Create the user account
kubectl exec --stdin --tty -n flux-operator ${pod} -- sudo useradd -m -p $(openssl passwd 'greatpw') greatuser

# Add them to flux accounting - the user bank is just called "user_bank"
kubectl exec --stdin --tty -n flux-operator ${pod} -- sudo -u flux flux account add-user --username=greatuser --bank=user_bank

# Update. If you don't run this, the user jobs will be stuck in priority-wait!
kubectl exec --stdin --tty -n flux-operator ${pod} -- sudo -u flux flux account-priority-update
```

To sanity check your user is added, you can do:

```bash
kubectl exec --stdin --tty -n flux-operator ${pod} -- sudo -u flux flux account view-user greatuser
```
```console
creation_time   mod_time        active          username        userid          bank            default_bank    shares          job_usage       fairshare       max_running_jobsmax_active_jobs max_nodes       queues          projects        default_project
1677276554      1677276554      1               greatuser       1006            user_bank       user_bank       1               0.0             0.5             5               7               2147483647                      *               *
```

Most of the fields should be populated, as shown above. If you see empty fields the user was not added properly.