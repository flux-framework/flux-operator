#!/usr/bin/env python3
# coding: utf-8

# This example shows saving state between two (different) miniclusters,
# where one has a lot of pending jobs that need to be continued to the next.
# We do this by way of turning on a flag to save an archive, and then load from
# that same location for the next MiniCluster.

import time
from kubernetes import client, config
from kubernetes.client import V1ObjectMeta
from fluxoperator.models import (
    MiniCluster,
    MiniClusterContainer,
    MiniClusterSpec,
    MiniClusterArchive,
    MiniClusterVolume,
    ContainerVolume,
)
from fluxoperator.client import FluxOperator
from fluxoperator.resource.pods import delete_minicluster

# Set our namespace and name
namespace = "flux-operator"
minicluster_name = "save-state"

# Here is our custom class to "wrap" an exec
class CommandExecutor(FluxOperator):
    fluxuser = "flux"

    def execute(self, command, print_result=True):
        """
        Wrap the kubectl_exec to add logic to issue to the broker instance.
        """
        res = self.kubectl_exec(
            f"sudo -u {self.fluxuser} flux proxy local:///var/run/flux/local {command}"
        )
        if print_result:
            print(res, end="")
        return res


# Here is our main container, we will use this for both clusters
container = MiniClusterContainer(
    image="ghcr.io/flux-framework/flux-restful-api:latest",
    volumes={"data": ContainerVolume(path="/state")},
    cores=2,
    run_flux=True,
)

# In order to save state we need a persistent volume between the MiniClusters
# it will be bound to /state, and the archive saved as "archive.tar.gz"
volumes = {"data": MiniClusterVolume(storage_class="hostpath", path="/tmp/data")}
archive = MiniClusterArchive(path="/state/archive.tar.gz")

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
        archive=archive,
        volumes=volumes,
    ),
)

# Make sure your cluster or minikube is running
# and the operator is installed
config.load_kube_config()

crd_api = client.CustomObjectsApi()

# Note that you might want to do this first for minikube
# minikube ssh docker pull ghcr.io/flux-framework/flux-restful-api:latest
# And create the cluster
crd_api.create_namespaced_custom_object(
    group="flux-framework.org",
    version="v1alpha1",
    namespace=namespace,
    plural="miniclusters",
    body=minicluster,
)

# Now let's create a flux operator client to interact
cli = CommandExecutor(namespace)

# Just call this so we know to wait
print("🥱️ Waiting for MiniCluster to be ready...")
cli.get_broker_pod()

# Let's exec commands to run a bunch of jobs!
# This is why we want interactive mod!
# By default, this selects (and waits for) the broker pod


print("✨️ Submitting a ton of jobs!")
time.sleep(5)
for iter in range(0, 50):
    res = cli.execute("flux submit sleep %s" % iter)
    assert res.startswith("ƒ")
    res = cli.execute("flux submit whoami")
    assert res.startswith("ƒ")

print("\n🥱️ Waiting for a few jobs...")
cli.execute("flux jobs -a")

print("\n🥱️ Asking flux to stop scheduling...")
cli.execute('flux queue disable "Pausing minicluster to save state"')

print("\n🥱️ Asking flux to stop the queue - this usually waits for running jobs...")
cli.execute("flux queue stop")
time.sleep(5)
print("\n🧐️ Inspecting jobs...")
cli.execute("flux jobs -a")

print("\n🧊️ Current state directory at /var/lib/flux...")
print(cli.kubectl_exec("ls -l /var/lib/flux"), end="")

print("\n🧊️ Current archive directory at /state... should be empty, not saved yet")
print(cli.kubectl_exec("ls -l /state"), end="")

print("Cleaning up...")
delete_minicluster(minicluster_name, namespace)
time.sleep(10)

# Increase size by 1
minicluster.spec.size = 3

print("\n🌀️ Creating second MiniCluster")
crd_api.create_namespaced_custom_object(
    group="flux-framework.org",
    version="v1alpha1",
    namespace=namespace,
    plural="miniclusters",
    body=minicluster,
)


# This also waits for the cluster to be running
print("🧊️ Current archive directory at /state... should now be populated")
print(cli.kubectl_exec("ls -l /state"), end="")
time.sleep(10)

print("\n🤓️ Inspecting state directory in new cluster...")
print(cli.kubectl_exec("ls -l /var/lib/flux"), end="")

print("\n😎️ Looking to see if old job history exists...")
cli.execute("flux jobs -a")

delete_minicluster(minicluster_name, namespace)
