#!/usr/bin/env python3
# coding: utf-8

# In this example we will create a persistent Minicluster
# the flux operator namespace, and port forward and issue commands.
# We will show getting a page (html) manually, and an API endpoint
# also manually, then with the flux-restful-client
# pip install flux-restful-client
# E.g.,:
# minikube start

from kubernetes import client, config

from fluxoperator.models import MiniCluster
from fluxoperator.models import MiniClusterSpec
from fluxoperator.models import MiniClusterContainer
from fluxoperator.models import MiniClusterUser
from fluxoperator.models import FluxRestful
from kubernetes.client import V1ObjectMeta
from fluxoperator.client import FluxOperator

import requests
import time
import json
import uuid

try:
    from flux_restful_client.main import get_client
except ImportError:
    print("flux-restful-client not installed, will skip API examples.")
    get_client = None

# Make sure your cluster or minikube is running
# and the operator is installed
config.load_kube_config()
crd_api = client.CustomObjectsApi()


def create_minicluster():
    """
    Here is how to create the minicluster
    """

    # Here is our main container with Flux accounting
    # run_flux True is required here
    container = MiniClusterContainer(
        image="ghcr.io/rse-ops/accounting:app-latest", run_flux=True
    )

    # Two users (set their passwords so we know)
    users = [
        MiniClusterUser(name="peenut", password="peenut"),
        MiniClusterUser(name="treenut", password="treenut"),
    ]

    flux_restful = FluxRestful(secret_key=str(uuid.uuid4()))
    # Create the MiniCluster
    minicluster = MiniCluster(
        kind="MiniCluster",
        api_version="flux-framework.org/v1alpha1",
        metadata=V1ObjectMeta(
            name="multi-tenant",
            namespace="flux-operator",
        ),
        spec=MiniClusterSpec(
            size=4,
            containers=[container],
            users=users,
            flux_restful=flux_restful
        ),
    )

    # Note that you might want to do this first for minikube
    # minikube ssh docker pull ghcr.io/rse-ops/accounting:app-latest
    return crd_api.create_namespaced_custom_object(
        group="flux-framework.org",
        version="v1alpha1",
        namespace="flux-operator",
        plural="miniclusters",
        body=minicluster,
    )


def delete_minicluster(result):
    return crd_api.delete_namespaced_custom_object(
        group="flux-framework.org",
        version="v1alpha1",
        namespace="flux-operator",
        plural="miniclusters",
        name=result["metadata"]["name"],
    )


# Create the MiniCluster
print("Creating the MiniCluster...")
result = create_minicluster()

# Create a client to interact with FluxOperator MiniCluster
cli = FluxOperator("flux-operator")

# First find the broker pod. This also calls cli.wait_pods()
broker = cli.get_broker_pod()

# And then port portfward to it - this waits until the server is ready
with cli.port_forward(broker) as forward_url:
    print(f"Port forward opened to {forward_url}")

    # This would be the welcome page
    response = requests.get(forward_url)
    print(response.status_code)

    # Get the secret key created for the server
    secret_key = result["spec"]["fluxRestful"]["secretKey"]

    # This would be the RESTFUl API
    # See endpoints here https://flux-framework.org/flux-restful-api/getting_started/api.html
    # You will need to authenticate for the multi-user example here
    if get_client is not None:

        restcli = get_client(host=forward_url, user="peenut", token="peenut", secret_key=secret_key)
        print('Submitting "whoami" job as user peenut.')
        res = restcli.submit("whoami")
        print(f"Jobid {res['id']} submit!")
        time.sleep(3)

        # Here is how to get all jobs for peenut
        jobs = restcli.jobs()
        print(json.dumps(jobs, indent=4))

        # And get output for the job
        output = restcli.output(res["id"])
        print(f"Job Output: {output.get('Output', '')}")

# How to cleanup
print("Cleaning up MiniCluster!")
delete_minicluster(result)
