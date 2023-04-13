# Fireworks with Flux

This will demonstrate running Fireworks on the Flux Operator to demonstrate a simple
[machine learning workflow](https://github.com/CrossFacilityWorkflows/DOE-HPC-workflow-training/tree/main/FireWorks/OLCF).

[![asciicast](https://asciinema.org/a/577459.svg)](https://asciinema.org/a/577459?speed=2)

## Usage

First, let's create a kind cluster. From the context of this directory:

```bash
$ kind create cluster --config ../../kind-config.yaml
```

And then install the operator, create the namespace, and apply the MiniCluster YAML here.

```bash
$ kubectl apply -f ../../dist/flux-operator.yaml
$ kubectl create namespace flux-operator
$ kubectl apply ./minicluster.yaml
```

You can watch the broker pod (0) to see (first) the submit of the tasks to the MongoDB, this part:

```bash
```
```console
...
2023-04-13 22:05:54,325 INFO Performing db tune-up
2023-04-13 22:05:54,434 INFO LaunchPad was RESET.
2023-04-13 22:05:54,437 INFO Added a workflow. id_map: {-3: 1, -2: 2, -1: 3}
Workflow submitted.
...
```

And then you can see the different flux workers picking up tasks, and printing output for us to see!

```console
...
🌀 Launcher Mode: flux start -o --config /etc/flux/config -Scron.directory=/etc/flux/system/cron.d   -Stbon.fanout=256   -Srundir=/run/flux    -Sstatedir=/var/lib/flux   -Slocal-uri=local:///run/flux/local  -Slog-stderr-level=6    -Slog-stderr-mode=local  python3 /tmp/workflow/run-workflow.py
broker.info[0]: start: none->join 5.66988ms
broker.info[0]: parent-none: join->init 0.031419ms
resource.err[0]: verify: rank 0 (flux-sample-0) has extra resources: core[1-3]
cron.info[0]: synchronizing cron tasks to event heartbeat.pulse
job-manager.info[0]: restart: 0 jobs
job-manager.info[0]: restart: 0 running jobs
job-manager.info[0]: restart: checkpoint.job-manager not found
broker.info[0]: rc1.0: running /etc/flux/rc1.d/01-sched-fluxion
sched-fluxion-resource.info[0]: version a9c1bd5
sched-fluxion-resource.warning[0]: create_reader: allowlist unsupported
sched-fluxion-resource.info[0]: populate_resource_db: loaded resources from core's resource.acquire
sched-fluxion-qmanager.info[0]: version a9c1bd5
broker.info[0]: rc1.0: running /etc/flux/rc1.d/02-cron
broker.info[0]: rc1.0: /etc/flux/rc1 Exited (rc=0) 3.5s
broker.info[0]: rc1-success: init->quorum 3.49085s
broker.info[0]: online: flux-sample-0 (ranks 0)
broker.info[0]: online: flux-sample-[0,6] (ranks 0,6)
broker.info[0]: online: flux-sample-[0,6-7] (ranks 0,6-7)
broker.info[0]: online: flux-sample-[0-9] (ranks 0-9)
broker.info[0]: quorum-full: quorum->run 21.2885s
2023-04-13 22:06:20,236 INFO Created new dir /tmp/workflow/launcher_2023-04-13-22-06-20-236857
2023-04-13 22:06:20,237 INFO Launching Rocket
2023-04-13 22:06:20,248 INFO RUNNING fw_id: 3 in directory: /tmp/workflow/launcher_2023-04-13-22-06-20-236857
2023-04-13 22:06:20,252 INFO Task started: ScriptTask.
                  0             1  ...             8             9
count  4.420000e+02  4.420000e+02  ...  4.420000e+02  4.420000e+02
mean  -2.511817e-19  1.230790e-17  ...  9.293722e-17  1.130318e-17
std    4.761905e-02  4.761905e-02  ...  4.761905e-02  4.761905e-02
min   -1.072256e-01 -4.464164e-02  ... -1.260971e-01 -1.377672e-01
25%   -3.729927e-02 -4.464164e-02  ... -3.324559e-02 -3.317903e-02
50%    5.383060e-03 -4.464164e-02  ... -1.947171e-03 -1.077698e-03
75%    3.807591e-02  5.068012e-02  ...  3.243232e-02  2.791705e-02
max    1.107267e-01  5.068012e-02  ...  1.335973e-01  1.356118e-01

[8 rows x 10 columns]
wrote x_diabetes file
wrote y_diabetes file
2023-04-13 22:06:21,281 INFO Task completed: ScriptTask
2023-04-13 22:06:21,298 INFO Rocket finished
2023-04-13 22:06:21,302 INFO Created new dir /tmp/workflow/launcher_2023-04-13-22-06-21-302642
2023-04-13 22:06:21,302 INFO Launching Rocket
2023-04-13 22:06:21,317 INFO RUNNING fw_id: 2 in directory: /tmp/workflow/launcher_2023-04-13-22-06-21-302642
2023-04-13 22:06:21,321 INFO Task started: ScriptTask.
Loading 'x_diabetes.npy'...
Loading 'y_diabetes.npy'...
Loading 'x_diabetes.npy'...
Loading 'y_diabetes.npy'...
Loading 'x_diabetes.npy'...
Loading 'y_diabetes.npy'...
Loading 'x_diabetes.npy'...
Loading 'y_diabetes.npy'...
Loading 'x_diabetes.npy'...
Loading 'y_diabetes.npy'...
Loading 'x_diabetes.npy'...
Loading 'y_diabetes.npy'...
Loading 'x_diabetes.npy'...
Loading 'y_diabetes.npy'...
Loading 'x_diabetes.npy'...
Loading 'y_diabetes.npy'...
Loading 'x_diabetes.npy'...
Loading 'y_diabetes.npy'...
Loading 'x_diabetes.npy'...
Loading 'y_diabetes.npy'...
[ 0.18788875  0.043062    0.58645013  0.44148176  0.21202248  0.17405359
 -0.39478925  0.43045288  0.56588259  0.38248348]
wrote all_coeffs file
2023-04-13 22:06:22,581 INFO Task completed: ScriptTask
2023-04-13 22:06:22,595 INFO Rocket finished
2023-04-13 22:06:22,599 INFO Created new dir /tmp/workflow/launcher_2023-04-13-22-06-22-599759
2023-04-13 22:06:22,600 INFO Launching Rocket
2023-04-13 22:06:22,612 INFO RUNNING fw_id: 1 in directory: /tmp/workflow/launcher_2023-04-13-22-06-22-599759
2023-04-13 22:06:22,627 INFO Task started: ScriptTask.
[ 0.188  0.043  0.586  0.441  0.212  0.174 -0.395  0.43   0.566  0.382]
Pearson correlation coefficients for each attribute
                                   0
age                         0.187889
sex                         0.043062
body_mass_index             0.586450
blood_pressure              0.441482
total_cholesterol           0.212022
ldl_cholesterol             0.174054
hdl_cholesterol            -0.394789
total/hdl_cholesterol       0.430453
log_of_serum_triglycerides  0.565883
blood_sugar_level           0.382483
2023-04-13 22:06:23,188 INFO Task completed: ScriptTask
2023-04-13 22:06:23,204 INFO Rocket finished
Done.
broker.info[0]: rc2.0: python3 /tmp/workflow/run-workflow.py Exited (rc=0) 4.2s
broker.info[0]: rc2-success: run->cleanup 4.22486s
broker.info[0]: cleanup.0: flux queue stop --quiet --all --nocheckpoint Exited (rc=0) 0.1s
broker.info[0]: cleanup.1: flux job cancelall --user=all --quiet -f --states RUN Exited (rc=0) 0.0s
broker.info[0]: cleanup.2: flux queue idle --quiet Exited (rc=0) 0.1s
broker.info[0]: cleanup-success: cleanup->shutdown 0.24966s
broker.info[0]: children-complete: shutdown->finalize 0.16094s
broker.info[0]: rc3.0: running /etc/flux/rc3.d/01-sched-fluxion
broker.info[0]: rc3.0: /etc/flux/rc3 Exited (rc=0) 0.2s
broker.info[0]: rc3-success: finalize->goodbye 0.218838s
broker.info[0]: goodbye: goodbye->exit 0.183057ms
...
```

If you want to debug something, you can set `interactive: true` to run in interactive mode, and then shell into the pod, connect to
the broker:

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-jlsp6 bash
$ sudo -u fluxuser -E $(env) -E HOME=/home/fluxuser flux proxy local:///run/flux/local bash
```

And run commands as you please! When you are done, clean up.

```bash
$ kubectl delete -f minicluster.yaml
```