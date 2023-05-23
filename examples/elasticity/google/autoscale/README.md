## Google GKE Autoscaling

> Using Google Cloud APIs to autoscale + Flux Operator autoscaling

These classes and assets are intended to help with autoscaling experiments.

Tests of elasticity, for GKE and the Flux Operator, will help us move us a tiny step toward the goal of 
trying to control both the GKE cluster scaling and MiniCluster scaling from the same
script. This way, theoretically, some other algorithm (that is smarter than my manual calls) can do it too.
See the [elasticity page](https://flux-framework.org/flux-operator/tutorials/elasticity.html) for early work.

To start, however, we just want to look at some basic elasticity or scaling for Kubernetes by itself.
All experiments use the shared class in [fluxcluster.py](fluxcluster.py), and an example experiment
is [test-scale.py](test-scale.py).


### Dependencies

We will need the google cloud APIs

```bash
$ python -m venv env
$ source env/bin/activate
# https://github.com/googleapis/python-container
$ pip install google-cloud-container
$ pip install kubernetes
```

### Experiments

The scripts are provided here to do experiments, and you can see results [here](https://github.com/converged-computing/operator-experiments/tree/main/google/autoscale/run1).
