#!/usr/bin/env python3
# coding: utf-8

# Test the multi-tenant use case. We assume the fluxoperator sdk is installed
# and a minikube cluster (or other Kuberentes cluster) is running

from kubernetes import client, config

from fluxoperator.models import MiniCluster
from fluxoperator.models import MiniClusterSpec
from fluxoperator.models import MiniClusterContainer
from fluxoperator.models import MiniClusterUser
from kubernetes.client import V1ObjectMeta
from fluxoperator.client import FluxOperator

import pytest
import time
import json
import sys

try:
    from flux_restful_client.main import get_client
except ImportError:
    sys.exit("please pip install flux-restful-client")

# Make sure your cluster or minikube is running
# and the operator is installed
config.load_kube_config()
crd_api = client.CustomObjectsApi()

test_users = ["peenut", "treenut"]


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
    users = []
    for user in test_users:
        users.append(MiniClusterUser(name=user, password=user))

    # Create the minicluster
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


def test_multi_tenant():

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

        # This connection without auth should fail
        restcli = get_client(host=forward_url)
        res = restcli.submit("whoami")
        assert "detail" in res
        assert "Not authenticated" in res["detail"]

        # Correct user and wrong token
        try:
            restcli = get_client(host=forward_url, user="peenut", token="nope")
            raise ValueError("Request with wrong token should fail")
        except:
            pass

        for user in test_users:
            restcli = get_client(host=forward_url, user=user, token=user)
            print(f'Submitting "whoami" job as user {user}.')
            res = restcli.submit("whoami")
            assert "id" in res
            print(f"Jobid {res['id']} submit!")
            time.sleep(3)

            # We should get able to get the job by id
            job = restcli.jobs(res["id"])
            assert job["name"] == "whoami"

            # We should only have one job - we only get back jobs for the user that requested
            jobs = restcli.jobs()
            assert len(jobs) == 1

            # And get output for the job
            output = restcli.output(res["id"]).get("Output", "")
            print(f"Job Output: {output}")
            assert output and user in output

    # How to cleanup
    print("Cleaning up MiniCluster!")
    delete_minicluster(result)