from kubernetes.client import V1ObjectMeta
from kubernetes import client, config
import fluxoperator.models as models

import fluxoperator.defaults as defaults
import sys

# These are known objects we will parse
_objects = ["logging", "volumes", "resources", "flux_restful", "container", "resources", "flux"]


def _get_logging_spec(logging):
    """
    Return models.Logging
    """
    logging_defaults = {
        "debug": False,
        "quiet": False,
        "strict": True,
        "timed": False,
        "zeromq": False,
    }

    for k, v in (logging or {}).items():
        if k in logging_defaults and v in [True, False]:
            logging_defaults[k] = v

    return models.LoggingSpec(**logging_defaults)


def _get_container_volumes(volumes):
    """
    Prepare container volumes.
    """
    volumeset = {}
    for name, volume in (volumes or {}).items():
        volume_spec = {}
        for attr in models.ContainerVolume.attribute_map:
            if attr in volume:
                volume_spec[attr] = volume[attr]
        volumeset[name] = models.ContainerVolume(**volume_spec)
    return volumeset


def _get_flux_spec(flux):
    """
    Get flux spec
    """
    fluxspec = {}
    for k in models.FluxSpec.attribute_map:
        if k in flux:
            fluxspec[k] = flux[k]
    return models.FluxSpec(**fluxspec)

def _get_container_spec(container):
    """
    Get the container spec.
    """
    # For now only one container support, it must run Flux
    container_kwargs = {"run_flux": True}
    for k in models.MiniClusterContainer.attribute_map:
        if k in container and k not in _objects:
            container_kwargs[k] = container[k]
        elif k in container and k == "volumes":
            container_kwargs[k] = _get_container_volumes(container[k])
        elif k in container and k == "resources":
            container_kwargs["resources"] = _get_container_resources_spec(container[k])
    return models.MiniClusterContainer(**container_kwargs)


def _get_container_resources_spec(resources):
    """
    Get container resources spec.
    """
    resources = resources or {}
    resource_spec = {}
    for k in models.ContainerResources.attribute_map:
        if k in resources:
            resource_spec[k] = resources[k]
    return models.ContainerResources(**resource_spec)


def _get_minicluster_spec(kwargs):
    """
    Get the main spec for the minicluster
    """
    minicluster_kwargs = {}
    for k in models.MiniClusterSpec.attribute_map:
        if k in kwargs and k not in _objects:
            minicluster_kwargs[k] = kwargs[k]
    return minicluster_kwargs


def _get_volumes_spec(volumes):
    """
    Prepare container volumes.
    """
    volumeset = {}
    for name, volume in (volumes or {}).items():
        volume_spec = {}
        for attr in models.MiniClusterVolume.attribute_map:
            if attr in volume:
                volume_spec[attr] = volume[attr]
        volumeset[name] = models.MiniClusterVolume(**volume_spec)
    return volumeset


def _get_flux_restful_spec(restful):
    """
    Get FluxRestful Spec
    """
    # Flux Restful pre-determined user and token
    flux_restful = restful or {}
    if "username" in flux_restful and "token" in flux_restful:
        flux_restful = models.FluxRestful(**flux_restful)
    return flux_restful


def create_minicluster(*args, **kwargs):
    """
    Create a MiniCluster of a particular size to run an image.

    The command is optional - if not provided will start the Flux RestFul API.
    This currently assumes running in single-user mode. The args/kwargs are
    left generic to be able to somewhat allow passing arbitrary dicts.
    """
    # Required to be in kwargs
    requireds = ["namespace", "name", "container"]
    for required in requireds:
        if required not in kwargs:
            sys.exit(f'A "{required}" field is required as a keyword argument.')

    container = kwargs["container"]
    namespace = kwargs["namespace"]
    name = kwargs["name"]
    del kwargs["container"]

    # This allows the client to provide a custom crd that already has credentials
    crd_api = kwargs.get('crd_api')
    if not crd_api:
        config.load_kube_config()
        crd_api = client.CustomObjectsApi()

    # We assume that this is a single container to run flux
    # Multi-container support can eventually be added.
    # TODO when requested, add pod resources
    container = _get_container_spec(container)

    flux_spec = _get_flux_spec(kwargs.get("flux"))
    logging_spec = _get_logging_spec(kwargs.get("logging"))
    flux_restful = _get_flux_restful_spec(kwargs.get("flux_restful"))
    volumes = _get_volumes_spec(kwargs.get("volumes"))
    minicluster_kwargs = _get_minicluster_spec(kwargs)

    # Create the MiniCluster
    minicluster = models.MiniCluster(
        kind="MiniCluster",
        api_version=f"flux-framework.org/{defaults.flux_operator_api_version}",
        metadata=V1ObjectMeta(name=name, namespace=namespace),
        spec=models.MiniClusterSpec(
            **minicluster_kwargs,
            flux=flux_spec,
            logging=logging_spec,
            containers=[container],
            flux_restful=flux_restful,
            volumes=volumes,
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


def delete_minicluster(name, namespace, **kwargs):
    """
    Delete a named MiniCluster.
    """
    # This allows the client to provide a custom crd that already has credentials
    crd_api = kwargs.get('crd_api')
    if not crd_api:
        config.load_kube_config()
        crd_api = client.CustomObjectsApi()

    return crd_api.delete_namespaced_custom_object(
        group="flux-framework.org",
        version=defaults.flux_operator_api_version,
        namespace=namespace,
        plural="miniclusters",
        name=name,
    )
