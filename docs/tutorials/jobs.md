# Submit Jobs

Given a running Kubernetes cluster, you have several modes for interaction, or submitting jobs:

- 1. Submit a MiniCluster custom resource definition (CRD) with a command to run one job. (<a href="#single-job-minicluster-crd">#ref</a>)
  - a. Submit a minicluster.yaml via `kubectl` (<a href="#a-submit-a-minicluster-yaml-via-kubectl">#ref</a>)
  - b. Submit a MiniCluster custom resource definition via the Python SDK (<a href="#b-submit-a-minicluster-crd-via-the-python-sdk">#ref</a>)
- 2. Submit in "interactive" mode, starting a Flux broker and instance, and interacting with your Flux instance. (<a href="submit-in-interactive-mode">#ref</a>)
- 3. Submit jobs interacting with Flux via ssh-ing to the pod (advanced) (<a href="#submit-jobs-directly-to-flux-via-ssh">#ref</a>)
- 4. Submit batch jobs

The reason we have many options is because we are still testing and understanding what use cases are best matched for each!
Each of the strategies above will be discussed here.

## 1. Single Job MiniCluster CRD

### a. Submit a MiniCluster yaml via kubectl

When you define a command in a MiniCluster CRD, this will submit and run one job using that command, and clean up.
You can find many examples for doing this under our [tests](https://github.com/flux-framework/flux-operator/tree/main/examples/tests)
directory. For any file there (e.g., lammps) once you have a Kubernetes cluster with the operator installed and
the `flux-operator` namespace created, you can do:

```bash
# Create the MiniCluster
$ kubectl apply -f examples/tests/lammps/minicluster.yaml

# See the status of the pods
$ kubectl get pods -n flux-operator pods

# View output for lammps (replace with actual pod id)
$ kubectl logs -n flux-operator lammps-0-xxxxx
```

### b. Submit a MiniCluster CRD via the Python SDK

We have early in development a [python SDK](https://github.com/flux-framework/flux-operator/tree/main/sdk/python)
that interacts directly with the Kubernetes API and defines the MiniCluster CRD programmatically.
You can see an early [directory of examples](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha2/examples).
Note that using this method you need to be pedantic - we use the [openapi generator](https://openapi-generator.tech/) and not all defaults
get correctly carried through. As an example, the default deadline in seconds is a large
number but it gets passed as 0, so if you don't define it, the Job will never start.

## 2. Submit in interactive mode

We have an [interactive mode](https://flux-framework.org/flux-operator/tutorials/interactive.html?h=interactive#interactive)
that means starting a single-user MiniCluster and launching a command.
We start the broker, and essentially wait for you to interact with it. This can mean issuing commands from
the terminal with `kubectl exec` (as shown in the tutorial above) or writing a Python script using the Flux Operator
Python SDK to proggrammatically interact. This is a really nice solution if you want to automate something.

## 3. Submit jobs directly to Flux via ssh

You can also submit jobs interacting with Flux via ssh-ing to the pod! This is considered
advanced, and is good for debugging. As you did before, get your pod listing:

```bash
$ kubectl -n flux-operator pods
```
And then shell in

```bash
$ kubectl exec --stdin --tty -n flux-operator ${brokerPod} -- /bin/bash
```

Note that if you have more than one container, you'll need to include it's name with `-c <container>`.
Also note that if you need to load any custom environments (e.g., something you'd define in the preCommand for a container)
you'll likely need to do that when you shell in.

## 4. Submit batch jobs

If you are submitting many jobs, you are better off providing them to flux at once as a batch submission.
This way, we won't stress any Kubernetes APIs to submit multiple. To do this, set `batch` to the batch
script content. The content is written verbatim to a file and run with `flux batch` — you provide the
`flux submit` lines yourself (they are not wrapped or numbered for you), which gives you full control over
flags, output paths, and `#flux:` directives. Here is the "containers" section of a MiniCluster crd:

```yaml
containers:
  - image: rockylinux:9

    # The batch content is written verbatim and run with flux batch.
    batch: |
      #!/bin/bash
      flux submit --output=/tmp/fluxout/job-1.out echo hello world 1
      flux submit --output=/tmp/fluxout/job-2.out echo hello world 2
      flux submit --output=/tmp/fluxout/job-3.out echo hello world 3
      flux queue idle
```

Because you write the `flux submit` lines yourself, you control exactly where each job's output goes
(via `--output`/`--error`) and any per-job flags. Set the output to a shared volume so all pods can write
to it. If you can't use the filesystem for saving output, it's recommended to have some other service used
in your jobs to send output.