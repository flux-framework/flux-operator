from kubernetes import client
from kubernetes.client.api import core_v1_api
from contextlib import contextmanager
from .resource.network import port_forward
from .resource.pods import create_minicluster, delete_minicluster
from .decorator import timed

import requests
import time


class FluxOperator:
    def __init__(self, namespace):
        """
        Create a persistent client to interact with a MiniCluster

        This currently assumes the namespace exists.
        """
        self.namespace = namespace
        self._core_v1 = None
        self.times = {}

    @property
    def core_v1(self):
        """
        Instantiate a core_v1 api (if not done yet)

        We have this here because we typically need to create the MiniCluster
        first.
        """
        if self._core_v1 is not None:
            return self._core_v1

        self.c = client.Configuration.get_default_copy()
        self.c.assert_hostname = False
        client.Configuration.set_default(self.c)
        self._core_v1 = core_v1_api.CoreV1Api()
        return self._core_v1

    @contextmanager
    def port_forward(self, pod, retry_seconds=1):
        """
        Wrapper to the port forward with context
        """
        start = time.time()
        with port_forward(self.core_v1):
            # Wait until this url is actually ready. In practice about 10-15 seconds
            url = f"http://{pod.metadata.name}.pod.flux-operator.kubernetes:5000"
            print()
            print(f"Waiting for {url} to be ready")
            ready = False
            while not ready:
                time.sleep(retry_seconds)
                try:
                    response = requests.get(url)
                    if response.status_code == 200:
                        end = time.time()
                        print("🪅️  RestFUL API server is ready!")
                        ready = True
                        self.times["restful_api_server_ready_seconds"] = end - start

                # There will be a few connection errors before everything is ready
                except Exception:
                    pass
                print(".", end="\r")

            print()
            yield url

    def wait_termination_pods(self, retry_seconds=1):
        """
        Ensure the namespace of pods is cleaned up before contining.
        """
        ready = False
        while not ready:
            pod_list = self.core_v1.list_namespaced_pod(self.namespace)
            if len(pod_list.items) == 0:
                ready = True
                break
            time.sleep(retry_seconds)
        print("All pods are terminated.")

    @timed
    def create_minicluster(self, name, size, image, namespace, user=None, token=None):
        """
        Create (and time the creation of) the MiniCluster
        """
        res = create_minicluster(name, size, image, namespace, user=user, token=token)
        self.wait_pods()
        return res

    @timed
    def delete_minicluster(self, name, namespace):
        """
        Deletion (and time the deletion of) the MiniCluster
        """
        res = delete_minicluster(name, namespace)
        self.wait_termination_pods()
        return res

    def wait_pods(self, states=None, retry_seconds=1):
        """
        Wait for all pods to be running or completed (or in a specific set of states)
        """
        states = states or ["Running", "Completed"]
        if not isinstance(states, list):
            states = [states]

        ready = False
        while not ready:
            pod_list = self.core_v1.list_namespaced_pod(self.namespace)
            ready = False
            for pod in pod_list.items:
                if pod.status.phase not in states:
                    time.sleep(retry_seconds)
                    continue
            # If we get down here, pods are ready
            ready = True

        states = '" or "'.join(states)
        print(f'All pods are in states "{states}"')

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
