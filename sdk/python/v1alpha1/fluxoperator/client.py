from kubernetes import client, watch
from kubernetes.stream import stream
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
        self._size = None

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

    def stream_output(self, filename, pod=None, stdout=True):
        """
        Stream output, optionally printing also to stdout.
        """
        if pod is None:
            pod = self.get_broker_pod()

        watcher = watch.Watch()

        # Stream output to file (should we return output too?)
        with open(filename, "w") as fd:
            for line in watcher.stream(
                self.core_v1.read_namespaced_pod_log,
                name=pod.metadata.name,
                namespace=self.namespace,
                follow=True,
            ):
                # Lines end with /r and we need to strip and add a newline
                fd.write(line.strip() + "\n")
                if stdout:
                    print(line)

    def kubectl_exec(self, command, pod=None, quiet=False):
        """
        Issue a command with kubectl exec

        A pod is required. If not provided, we retrieve (and wait for)
        the broker pod to be ready.
        """
        if not pod:
            pod = self.get_broker_pod(quiet=quiet)
        if not quiet:
            print(command)

        # Assemble the exec command - bash subshell always easier
        exec_command = ["/bin/sh", "-c", command]
        return stream(
            self.core_v1.connect_get_namespaced_pod_exec,
            pod.metadata.name,
            self.namespace,
            command=exec_command,
            stderr=True,
            stdin=False,
            stdout=True,
            tty=False,
        )

    def get_pods(self):
        """
        Get namespaced pods metadata
        """
        return self.core_v1.list_namespaced_pod(self.namespace)

    def get_nodes(self):
        """
        Get nodes metadata
        """
        return self.core_v1.list_node()

    def wait_termination_pods(self, retry_seconds=1, quiet=False):
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
        if not quiet:
            print("All pods are terminated.")

    @timed
    def create_minicluster(self, *args, **kwargs):
        """
        Create (and time the creation of) the MiniCluster
        """
        self._size = kwargs.get("size")
        res = create_minicluster(*args, **kwargs)
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

    def wait_pods(self, states=None, retry_seconds=1, quiet=False):
        """
        Wait for all pods to be running or completed (or in a specific set of states)
        """
        states = states or ["Running", "Succeeded"]
        if not isinstance(states, list):
            states = [states]

        # We don't have a size - get from cluster
        if not self._size:
            time.sleep(10)
            pod_list = self.core_v1.list_namespaced_pod(self.namespace)
            self._size = len(pod_list.items)

        ready = set()
        while len(ready) != self._size:
            pod_list = self.core_v1.list_namespaced_pod(self.namespace)

            # Upate size with more pods
            if len(pod_list.items) > self._size:
                self._size = len(pod_list.items)
            for pod in pod_list.items:
                if pod.status.phase not in states:
                    time.sleep(retry_seconds)
                else:
                    ready.add(pod.metadata.name)

        if not quiet:
            states = '" or "'.join(states)
            print(f'All pods are in states "{states}"')

    def get_broker_pod(self, quiet=False):
        """
        Given a core_v1 connection and namespace, get the broker pod.
        """
        # All pods required to be ready
        self.wait_pods(quiet=quiet)

        # Go through process again and get broker, must be running!
        brokerPod = None
        while not brokerPod:
            pod_list = self.core_v1.list_namespaced_pod(self.namespace)
            for pod in pod_list.items:
                if "-0" in pod.metadata.name and pod.status.phase in ["Running"]:
                    if not quiet:
                        print(f"Found broker pod {pod.metadata.name}")
                    brokerPod = pod
        return brokerPod
