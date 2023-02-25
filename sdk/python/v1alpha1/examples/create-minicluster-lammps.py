#!/usr/bin/env python3
# coding: utf-8

from kubernetes import client, config
from kubernetes.client.models.v1_object_meta import V1ObjectMeta

# These names are long and ugly, but thy are versioned with the FluxOperator!
from fluxoperator.models import ApiV1alpha1MiniCluster
from fluxoperator.models import ApiV1alpha1MiniClusterSpec
from fluxoperator.models import ApiV1alpha1MiniClusterContainer


# Here is our main container
container = ApiV1alpha1MiniClusterContainer(
    cores = 2,
    image = "ghcr.io/rse-ops/lammps:flux-sched-focal-v0.24.0",
    working_dir = "/home/flux/examples/reaxff/HNS",
    command = "lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite"
)

minicluster = ApiV1alpha1MiniCluster(
    kind="MiniCluster",
    api_version="flux-framework.org/v1alpha1",
    metadata=V1ObjectMeta(
        name="lammps",
    ),
    spec=ApiV1alpha1MiniClusterSpec(
        size=4,
        tasks=2,
        containers = [container]
    )
)

# Make sure your cluster or minikube is running 
# and the operator is installed
config.load_kube_config()

crd_api = client.CustomObjectsApi()


crd_api.create_namespaced_custom_object(
    group="flux-framework.org",
    version="v1alpha1",
    namespace="flux-operator",
    plural="miniclusters",
    body=minicluster
)
