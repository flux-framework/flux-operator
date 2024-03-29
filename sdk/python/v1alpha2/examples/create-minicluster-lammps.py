#!/usr/bin/env python3
# coding: utf-8

# This is an example of creating a Lammps cluster using the native API models
# directly (and not using the client.FluxOperator or client.FluxMiniCluster classes)

from kubernetes import client, config
from kubernetes.client import V1ObjectMeta

from fluxoperator.models import MiniCluster, MiniClusterContainer, MiniClusterSpec


# Here is our main container
container = MiniClusterContainer(
    cores=2,
    image="ghcr.io/rse-ops/lammps:flux-sched-focal",
    working_dir="/home/flux/examples/reaxff/HNS",
    command="lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite",
    run_flux=True,
)

# There is currently a bug where the defaults are not set/correct, so for example,
# we need to set the deadline seconds or the minicluster will not create.
minicluster = MiniCluster(
    kind="MiniCluster",
    api_version="flux-framework.org/v1alpha1",
    metadata=V1ObjectMeta(
        name="lammps",
        namespace="flux-operator",
    ),
    spec=MiniClusterSpec(
        size=4, tasks=2, deadline_seconds=31500000, containers=[container]
    ),
)

# Make sure your cluster or minikube is running
# and the operator is installed
config.load_kube_config()

crd_api = client.CustomObjectsApi()

# Note that you might want to do this first for minikube
# minikube ssh docker pull ghcr.io/rse-ops/lammps:flux-sched-focal-v0.24.0

result = crd_api.create_namespaced_custom_object(
    group="flux-framework.org",
    version="v1alpha1",
    namespace="flux-operator",
    plural="miniclusters",
    body=minicluster,
)

# At this point you can look at your pods in the result, or similar
