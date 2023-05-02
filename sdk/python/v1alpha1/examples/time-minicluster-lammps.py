#!/usr/bin/env python3
# coding: utf-8

# This is an example of creating a Lammps cluster using the native API models
# directly (and not using the client.FluxOperator or client.FluxMiniCluster classes)

import statistics
from kubernetes import client, config
from fluxoperator.client import FluxMiniCluster
import time
import json

# Here is our main container
container = {
    "image": "vanessa/lammps:test-zeromq",
#    "image": "ghcr.io/rse-ops/lammps:flux-sched-focal",
    "working_dir": "/home/flux/examples/reaxff/HNS",
    "command": "lmp -v x 1 -v y 1 -v z 1 -in in.reaxc.hns -nocite",
    "flux_log_level": 7,
}

# Make sure your cluster or minikube is running
# and the operator is installed
config.load_kube_config()

crd_api = client.CustomObjectsApi()

# Note that you might want to do this first for minikube
# minikube ssh docker pull ghcr.io/rse-ops/lammps:flux-sched-focal-v0.24.0

# Interact with the Flux Operator Python SDK
minicluster = {
    "size": 2,
    "namespace": "flux-operator",
    "name": "lammps",
    "logging": {"zeromq": True},
}

times = []

for iter in range(0, 20):
    operator = FluxMiniCluster()
    start = time.time()
    operator.create(**minicluster, container=container)
    # Ensure we keep it hanging until the job finishes
    operator.stream_output("lammps.out", stdout=True, timestamps=True)
    end = time.time()
    runtime = end - start
    times.append(runtime)
    print(f"Runtime for teeny LAMMPS iteration {iter} is {runtime} seconds")
    operator.delete()

print(json.dumps(times))
print(f"Mean: {statistics.mean(times)}")
print(f"Std: {statistics.stdev(times)}")
