#!/usr/bin/env python3
# coding: utf-8

# This is an example of creating a Lammps cluster using the native API models
# directly (and not using the client.FluxOperator or client.FluxMiniCluster classes)

import statistics
from kubernetes import client, config
from fluxoperator.client import FluxMiniCluster
import copy
import time
import json

# Here is our main container
container = {
    "image": "ghcr.io/rse-ops/lammps:flux-sched-focal",
    "working_dir": "/home/flux/examples/reaxff/HNS",
    "command": "lmp -v x 1 -v y 1 -v z 1 -in in.reaxc.hns -nocite",
}

# Make sure your cluster or minikube is running
# and the operator is installed
config.load_kube_config()

crd_api = client.CustomObjectsApi()

# Note that you might want to do this first for minikube
# minikube ssh docker pull ghcr.io/rse-ops/lammps:flux-sched-focal-v0.24.0

# Interact with the Flux Operator Python SDK
mc = {
    # Note that when you enable a service pod, the indexed job can come up 3-4 seconds faster
    # it seems like a bug, unfortunatly
    "size": 2,
    "namespace": "flux-operator",
    "name": "lammps",
    "logging": {"zeromq": True},
    "flux": {"connect_timeout": "5s", "log_level": 7},
}


def run_experiments(results_file=None, with_services=False, stdout=False):
    """
    Shared script to run experiments, so we can try across many different cases
    """
    minicluster = copy.deepcopy(mc)
    if with_services:
        minicluster["services"] = [
            {
                "image": "nginx",
                "name": "nginx",
                "ports": [80],
            }
        ]

    times = {}
    means = {}
    stds = {}

    # This will test timeouts for the connection between 0 and 10
    # Along with the case of not setting one
    # A time of zero will be unset (default to 30 seconds, long!)
    for timeout in [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, None]:
        if timeout is None:
            slug = "no-timeout-set"
            minicluster["flux"]["connect_timeout"] = ""
        else:
            slug = f"{timeout}s"
            minicluster["flux"]["connect_timeout"] = slug

        # Be pedantic and lazy and save everything as we go
        times[slug] = []
        means[slug] = []
        stds[slug] = []

        # 20 times each
        for iter in range(0, 20):
            print(f"TESTING {slug} iteration {iter}")
            operator = FluxMiniCluster()
            start = time.time()
            operator.create(**minicluster, container=container)

            # Ensure we keep it hanging until the job finishes
            operator.stream_output("lammps.out", stdout=stdout, timestamps=True)
            end = time.time()
            runtime = end - start
            times[slug].append(runtime)
            print(f"Runtime for teeny LAMMPS iteration {iter} is {runtime} seconds")
            operator.delete()

        print(json.dumps(times[slug]))
        means[slug] = statistics.mean(times[slug])
        stds[slug] = statistics.stdev(times[slug])
        print(f"Mean: {means[slug]}")
        print(f"Std: {stds[slug]}")

    print("\nFinal Results:")
    print(json.dumps(times))
    print(f"Means: {means}")
    print(f"Stds: {stds}")
    results = {"times": times, "means": means, "stds": stds}
    with open(results_file, "w") as fd:
        fd.write(json.dumps(results, indent=4))


with_stdout = False
run_experiments("lammps-no-services.json", False, with_stdout)
run_experiments("lammps-with-services.json", True, with_stdout)
