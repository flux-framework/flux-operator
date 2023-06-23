from kubernetes import client, watch
from kubernetes.stream import stream
from kubernetes.client.api import core_v1_api
from kubernetes.client.models.v1_pod_list import V1PodList

import kubernetes.client.exceptions
from contextlib import contextmanager
from .resource.network import port_forward
from .resource.pods import create_minicluster, delete_minicluster
from .decorator import timed

import logging
import requests
import time

# This should be better exposed
logger = logging.getLogger("fluxoperator")


class FluxMiniCluster:
    """
    A MiniCluster is a small class to hold Metadata about a minicluster!

    Each MiniCluster holds it's own controller for the Flux Operator.
    """

    def __init__(self, **kwargs):
        """
        Create a persistent client to interact with a MiniCluster

        This currently assumes the namespace exists.
        """
        self.times = {}
        self.metadata = None
        self._pods = None
        self._broker_pod = None
        self.ctrl = FluxOperator(**kwargs)

    def load(self, metadata):
        """
        Load minicluster metadata
        """
        self.metadata = metadata
        self._broker_pod = None
        return self.wait_pods()

    @timed
    def create(self, *args, **kwargs):
        """
        Create (and time the creation of) the MiniCluster
        """
        self.metadata = create_minicluster(*args, **kwargs)
        self._broker_pod = None

        # Save the pods upfront so we don't have to query for them again.
        return self.wait_pods()

    @property
    def flux_user(self):
        """
        Derive the name of the flux instance owner.
        """
        return (self.metadata or {}).get("spec", {}).get("flux_user") or "flux"

    def wait_pods(self, retry_seconds=1, quiet=False):
        """
        A wrapper to wait for pods.
        """
        self._pods = self.ctrl.wait_pods(
            namespace=self.namespace,
            name=self.name,
            size=self.size,
            retry_seconds=retry_seconds,
            quiet=quiet,
        )
        return self._pods

    def get_nodes(self):
        """
        Get nodes metadata
        """
        return self.ctrl.get_nodes()

    @property
    def broker_pod(self):
        """
        The broker pod has the name + 0
        """
        if self._broker_pod:
            return self._broker_pod

        if not self._pods:
            return
        for pod in self._pods.items:
            if pod.metadata.name.startswith(f"{self.name}-0"):
                self._broker_pod = pod.metadata.name
                return self._broker_pod

    @property
    def pods(self):
        """
        Get the pods (set of names)
        """
        return set([x.metadata.name for x in self._pods.items])

    def get_pods(self):
        """
        A wrapper to the Flux Operator get pods
        """
        # Return cached pods to not stress the kubernetes API
        if self._pods:
            return self._pods
        return self.ctrl.get_pods(name=self.name, namespace=self.namespace)

    @contextmanager
    def port_forward(self, retry_seconds=1):
        """
        Wrapper to the port forward with context
        """
        start = time.time()
        with port_forward(self.ctrl.core_v1):
            # Wait until this url is actually ready. In practice about 10-15 seconds
            url = f"http://{self.broker_pod}.pod.flux-operator.kubernetes:5000"
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

            print()
            yield url

    def stream_output(self, filename, stdout=True, timestamps=False, pod=None):
        """
        Stream output, optionally printing also to stdout.

        We allow specifying the pod if we don't want output from the broker.
        """
        pod = pod or self.broker_pod
        return self.ctrl.stream_output(
            filename=filename,
            stdout=stdout,
            namespace=self.namespace,
            timestamps=timestamps,
            name=self.name,
            pod=pod,
        )

    @timed
    def delete(self):
        """
        Deletion (and time the deletion of) the MiniCluster
        """
        return self.ctrl.delete_minicluster(name=self.name, namespace=self.namespace)

    @property
    def name(self):
        return self.metadata["metadata"]["name"]

    @property
    def namespace(self):
        return self.metadata["metadata"]["namespace"]

    @property
    def size(self):
        return self.metadata["spec"]["size"]


class FluxBrokerMiniCluster(FluxMiniCluster):
    """
    A MiniCluster with a broker can have commands exec'd to the socket.
    """

    def __init__(self, socket=None, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.socket = socket or "local:///var/run/flux/local"

    def execute(self, command, print_result=True):
        """
        Wrap the kubectl_exec to add logic to issue to the broker instance.
        """
        res = self.ctrl.kubectl_exec(
            f"sudo -u {self.flux_user} flux proxy {self.socket} {command}",
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


class FluxOperator:
    def __init__(self, namespace=None, **kwargs):
        """
        Create a persistent client to interact with a MiniCluster

        This currently assumes the namespace exists.
        """
        self._core_v1 = kwargs.get('core_v1_api')
        self.namespace = namespace

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

    def stream_output(
        self,
        filename,
        name=None,
        namespace=None,
        pod=None,
        stdout=True,
        return_output=True,
        timestamps=False,
    ):
        """
        Stream output, optionally printing also to stdout.

        Also return the output to the user.
        """
        namespace = namespace or self.namespace
        pod = pod or self.get_broker_pod(name=name, namespace=namespace).metadata.name
        watcher = watch.Watch()

        # Stream output to file and return it if desired!
        lines = []
        with open(filename, "w") as fd:
            for line in watcher.stream(
                self.core_v1.read_namespaced_pod_log,
                name=pod,
                namespace=namespace,
                timestamps=timestamps,
                follow=True,
            ):
                # Lines end with /r and we need to strip and add a newline
                fd.write(line.strip() + "\n")
                if stdout:
                    print(line)
                if return_output:
                    lines.append(line)

        # I can imagine cases where we wouldn't want to keep it
        if return_output:
            return lines

    def kubectl_exec(self, command, pod=None, quiet=False, namespace=None, name=None):
        """
        Issue a command with kubectl exec

        A pod is required. If not provided, we retrieve (and wait for)
        the broker pod to be ready. If you are doing multiple commands
        in a row, it's recommended to provide the broker pod.
        """
        namespace = namespace or self.namespace
        pod = (
            pod
            or self.get_broker_pod(
                quiet=quiet, name=name, namespace=namespace
            ).metadata.name
        )
        if not quiet:
            print(command)

        # Assemble the exec command - bash subshell always easier
        exec_command = ["/bin/sh", "-c", command]
        return stream(
            self.core_v1.connect_get_namespaced_pod_exec,
            pod,
            namespace,
            command=exec_command,
            stderr=True,
            stdin=False,
            stdout=True,
            tty=False,
        )

    def get_pods(self, namespace, name=None):
        """
        Get namespaced pods metadata, either scoped to a name or entire namespace.
        """
        namespace = namespace or self.namespace
        try:
            req = self.core_v1.list_namespaced_pod(namespace, async_req=True)
            pods = req.get()

            # If name is present, filter down to pods with that prefix
            if name is not None:
                pods = self._filter_pods(pods, name)
            return pods

        # Not found - it was deleted
        except kubernetes.client.exceptions.ApiException:
            return V1PodList(items=[])
        except:
            time.sleep(2)
            return self.get_pods(namespace, name)

    def _filter_pods(self, pods, name):
        """
        Filter a set of pods (associated with a job) to a name prefix.
        """
        filtered = []
        for pod in pods.items:
            if pod.metadata.name.startswith(name):
                filtered.append(pod)
        pods.items = filtered
        return pods

    def get_nodes(self):
        """
        Get nodes metadata
        """
        return self.core_v1.list_node()

    def wait_termination_pods(
        self, name=None, namespace=None, retry_seconds=1, quiet=False
    ):
        """
        Ensure the namespace of pods is cleaned up before contining.
        """
        namespace = namespace or self.namesapce
        ready = False
        while not ready:
            pod_list = self.get_pods(name=name, namespace=namespace)
            if len(pod_list.items) == 0:
                ready = True
                break
            time.sleep(retry_seconds)
        if not quiet:
            print("All pods are terminated.")

    def delete_minicluster(self, name=None, namespace=None, **kwargs):
        """
        Deletion (and time the deletion of) the MiniCluster
        """
        namespace = namespace or self.namespace
        res = delete_minicluster(name, namespace, **kwargs)
        self.wait_termination_pods(name, namespace)
        return res

    def wait_pods(
        self,
        size=None,
        name=None,
        namespace=None,
        states=None,
        retry_seconds=1,
        quiet=False,
    ):
        """
        Wait for all pods to be running or completed (or in a specific set of states)
        """
        namespace = namespace or self.namespace
        states = states or ["Running", "Succeeded", "Completed"]
        if not isinstance(states, list):
            states = [states]

        # We don't have a size - get from cluster
        if not size:
            time.sleep(10)
            pod_list = self.get_pods(name=name, namespace=namespace)
            size = len(pod_list.items)

        ready = set()
        while len(ready) != size:
            logger.debug(f"{len(ready)} pods are ready, out of {size}")
            pod_list = self.get_pods(name=name, namespace=namespace)

            for pod in pod_list.items:
                print(f"{pod.metadata.name} is in phase {pod.status.phase}")

                # Don't include the cert generator pod
                if "cert-generator" in pod.metadata.name:
                    continue

                # Ignore services pod
                if pod.metadata.name.endswith("-services"):
                    continue
                if pod.status.phase not in states:
                    time.sleep(retry_seconds)
                    continue

                if pod.status.phase not in ["Terminating"]:
                    ready.add(pod.metadata.name)

        if not quiet:
            states = '" or "'.join(states)
            print(f'All pods are in states "{states}"')
        return pod_list

    def get_broker_pod(self, quiet=False, name=None, namespace=None):
        """
        Given a core_v1 connection and namespace, get the broker pod.
        """
        namespace = namespace or self.namespace

        # All pods required to be ready
        self.wait_pods(quiet=quiet, name=name, namespace=namespace)

        # Go through process again and get broker, must be running!
        brokerPod = None
        while not brokerPod:
            pod_list = self.get_pods(name=name, namespace=namespace)
            for pod in pod_list.items:
                if "-0" in pod.metadata.name and pod.status.phase in ["Running"]:
                    if not quiet:
                        print(f"Found broker pod {pod.metadata.name}")
                    brokerPod = pod
        return brokerPod
