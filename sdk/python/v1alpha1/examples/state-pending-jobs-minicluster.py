#!/usr/bin/env python3
# coding: utf-8

# This example shows saving state between two (different) miniclusters,
# where one has a lot of pending jobs that need to be continued to the next.
# We do this by way of turning on a flag to save an archive, and then load from
# that same location for the next MiniCluster.

# IMPORTANT! Remember to clear your old archive first
# minikube ssh -- rm /tmp/data/archive.tar.gz

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

# TODO this should have a different name to not confuse?
from fluxoperator.client import FluxMiniCluster
from fluxoperator.resource.pods import delete_minicluster

# Set our namespace and name
namespace = "flux-operator"
minicluster_name = "save-state"


# Here is our custom class to "wrap" an exec
class CommandExecutor(FluxMiniCluster):
    """
    A CommandExecutor wraps a FluxOperator controller. Note
    that this class assumes we have one MiniCluster to interact with -
    it will group all pods in the same namespace. If you intend to control
    multiple MiniClusters as once, use the MiniClusterManager.
    """

    fluxuser = "flux"

    def execute(self, command, print_result=True):
        """
        Wrap the kubectl_exec to add logic to issue to the broker instance.
        """
        res = self.ctrl.kubectl_exec(
            f"sudo -u {self.fluxuser} flux proxy local:///var/run/flux/local {command}",
            name=self.name,
            namespace=self.namespace,
            pod=self.broker_pod,
        )
        if print_result:
            print(res, end="")
        return res

    def kubectl_exec(self, command, print_result=True):
        res = self.ctrl.kubectl_exec(
            command, name=self.name, namespace=self.namespace, pod=self.broker_pod
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
result = crd_api.create_namespaced_custom_object(
    group="flux-framework.org",
    version="v1alpha1",
    namespace=namespace,
    plural="miniclusters",
    body=minicluster,
)

# Now let's create a flux operator client to interact
print("🥱️ Waiting for MiniCluster to be ready...")
cli = CommandExecutor()
cli.load(result)

# Let's exec commands to run a bunch of jobs!
# This is why we want interactive mod!
# By default, this selects (and waits for) the broker pod

print("✨️ Submitting a ton of jobs!")
time.sleep(5)
for iter in range(0, 30):
    res = cli.execute("flux submit sleep %s" % iter)
    assert res.startswith("ƒ")
    res = cli.execute("flux submit whoami")
    assert res.startswith("ƒ")

print("\n🥱️ Waiting for a few jobs...")
cli.execute("flux jobs -a")

print("\n🥱️ Asking flux to stop the queue...")
cli.execute("flux queue stop")
time.sleep(5)

print("\n🥱️ Waiting for running jobs...")
cli.execute("flux queue idle")

print("\n💩️ Dumping the archive...")
cli.execute("flux dump /state/archive.tar.gz")

print("\n🧐️ Inspecting jobs...")
cli.execute("flux jobs -a")

print("\n🧊️ Current state directory at /var/lib/flux...")
cli.kubectl_exec("ls -l /var/lib/flux")

print("\n🧊️ Current archive directory at /state... should be empty, not saved yet")
cli.kubectl_exec("ls -l /state")

print("Cleaning up...")
cli.delete()
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
print("Wait for MiniCluster...")
time.sleep(120)

# This also waits for the cluster to be running
print("🧊️ Current archive directory at /state... should now be populated")
cli.kubectl_exec("ls -l /state")
time.sleep(10)

print("\n🤓️ Inspecting state directory in new cluster...")
cli.kubectl_exec("ls -l /var/lib/flux")
time.sleep(10)

print("\n😎️ Looking to see if old job history exists...")
cli.execute("flux jobs -a")
cli.delete()
