#!/usr/bin/env python3
# coding: utf-8

# This example shows a interactive submit to a single-user MiniCluster. This
# does not require the Flux Restful API - instead we submit via a kubectl exec.
# In addition, we write a custom class that wraps the command to issue a
# request to the broker.

import time
from kubernetes import client, config
from kubernetes.client import V1ObjectMeta
from fluxoperator.models import (
    MiniCluster,
    MiniClusterContainer,
    MiniClusterSpec,
)
from fluxoperator.client import FluxBrokerMiniCluster
from fluxoperator.resource.pods import delete_minicluster

# Set our namespace and name
namespace = "flux-operator"
minicluster_name = "interactive-submit"


# Here is our main container, we will use this for both clusters
container = MiniClusterContainer(
    image="ghcr.io/flux-framework/flux-restful-api:latest",
    cores=2,
    run_flux=True,
)

# This is creating the full minicluster
# Interactive is set to true so the broker starts Flux
# and then we can interact / submit as we please!
minicluster = MiniCluster(
    kind="MiniCluster",
    api_version="flux-framework.org/v1alpha1",
    metadata=V1ObjectMeta(
        name=minicluster_name,
        namespace=namespace,
    ),
    spec=MiniClusterSpec(
        size=2,
        containers=[container],
        interactive=True,
    ),
)

# Make sure your cluster or minikube is running
# and the operator is installed
config.load_kube_config()

crd_api = client.CustomObjectsApi()

# Note that you might want to do this first for minikube
# minikube ssh docker pull ghcr.io/flux-framework/flux-restful-api:latest
# And create the cluster. This can also be done with cli.create(**minicluster)
result = crd_api.create_namespaced_custom_object(
    group="flux-framework.org",
    version="v1alpha1",
    namespace=namespace,
    plural="miniclusters",
    body=minicluster,
)

# Now let's create a flux operator client to interact
# This will wait for pods to be ready
print("🥱️ Waiting for MiniCluster to be ready...")
cli = FluxBrokerMiniCluster()
cli.load(result)

# Just call this so we know to wait
# Let's exec commands to run a bunch of jobs!
# This is why we want interactive mod!
# By default, this selects (and waits for) the broker pod
print("✨️ Submitting jobs!")
time.sleep(5)
for iter in range(0, 5):
    res = cli.execute("flux submit sleep %s" % iter)
    assert res.startswith("ƒ")
    res = cli.execute("flux submit whoami")
    assert res.startswith("ƒ")

print("\n🥱️ Waiting for jobs...")
print("Jobs finished...")
cli.execute("flux jobs -a")

print("\n🥱️ Wait to be sure we have finished...")
time.sleep(50)
cli.execute("flux jobs -a")

print("Cleaning up...")
cli.delete()
