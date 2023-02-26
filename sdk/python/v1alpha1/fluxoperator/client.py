from kubernetes import client
from kubernetes.client.api import core_v1_api
from contextlib import contextmanager
from .resource.network import port_forward

import requests
import time


class FluxOperator:
    def __init__(self, namespace):
        """
        Create a persistent client to interact with a MiniCluster

        This currently assumes the namespace exists.
        """
        self.namespace = namespace
        self.c = client.Configuration.get_default_copy()
        self.c.assert_hostname = False
        client.Configuration.set_default(self.c)
        self.core_v1 = core_v1_api.CoreV1Api()

    @contextmanager
    def port_forward(self, pod):
        """
        Wrapper to the port forward with context
        """
        with port_forward(self.core_v1):
            # Wait until this url is actually ready. In practice about 10-15 seconds
            url = f"http://{pod.metadata.name}.pod.flux-operator.kubernetes:5000"
            print()
            print(f"Waiting for {url} to be ready")
            sleep = 2
            ready = False
            while not ready:
                time.sleep(sleep)
                try:
                    response = requests.get(url)
                    if response.status_code == 200:
                        print("🪅️ RestFUL API server is ready!")
                        ready = True

                # There will be a few connection errors before everything is ready
                except Exception:
                    pass
                print(".", end="\r")
                sleep = sleep * 2

            print()
            yield url

    def wait_pods(self):
        """
        Wait for all pods to be running or completed
        """
        ready = False
        while not ready:
            pod_list = self.core_v1.list_namespaced_pod(self.namespace)
            ready = False
            for pod in pod_list.items:
                if pod.status.phase not in ["Running", "Completed"]:
                    time.sleep(2)
                    continue
            # If we get down here, pods are ready
            ready = True

        print('All pods are "Running" or "Completed"')

    def get_broker_pod(self):
        """
        Given a core_v1 connection and namespace, get the broker pod.
        """
        # All pods required to be ready
        self.wait_pods()

        # Go through process again and get broker
        brokerPod = None
        while not brokerPod:
            pod_list = self.core_v1.list_namespaced_pod(self.namespace)
            for pod in pod_list.items:
                if "-0" in pod.metadata.name:
                    print(f"Found broker pod {pod.metadata.name}")
                    brokerPod = pod
        return brokerPod
