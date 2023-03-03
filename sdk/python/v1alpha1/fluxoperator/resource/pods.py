from kubernetes.client import V1ObjectMeta
from kubernetes import client, config
import fluxoperator.models as models

import fluxoperator.defaults as defaults


def create_minicluster(name, size, image, namespace, user=None, token=None):
    """
    Create a MiniCluster of a particular size to run an image.

    This currently assumes running in single-user mode.
    """
    # The cluster should be running with the operator installed
    config.load_kube_config()
    crd_api = client.CustomObjectsApi()

    # We assume that this is a single container to run flux
    # Multi-container support can eventually be added.
    container = models.MiniClusterContainer(image=image, run_flux=True)

    # Flux Restful pre-determined user and token
    flux_restful = None
    if user is not None and token is not None:
        flux_restful = models.FluxRestful(username=user, token=token)

    # Create the MiniCluster
    minicluster = models.MiniCluster(
        kind="MiniCluster",
        api_version=f"flux-framework.org/{defaults.flux_operator_api_version}",
        metadata=V1ObjectMeta(name=name, namespace=namespace),
        spec=models.MiniClusterSpec(
            flux_restful=flux_restful,
            size=size,
            containers=[container],
        ),
    )

    # Note that you might want to do this first for minikube
    # minikube ssh docker pull ghcr.io/rse-ops/accounting:app-latest
    return crd_api.create_namespaced_custom_object(
        group="flux-framework.org",
        version=defaults.flux_operator_api_version,
        namespace=namespace,
        plural="miniclusters",
        body=minicluster,
    )


def delete_minicluster(name, namespace):
    """
    Delete a named MiniCluster.
    """
    config.load_kube_config()
    crd_api = client.CustomObjectsApi()

    return crd_api.delete_namespaced_custom_object(
        group="flux-framework.org",
        version=defaults.flux_operator_api_version,
        namespace=namespace,
        plural="miniclusters",
        name=name,
    )
