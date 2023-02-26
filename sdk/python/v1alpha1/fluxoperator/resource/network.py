from kubernetes.stream import portforward
from contextlib import contextmanager

import urllib3


@contextmanager
def port_forward(core_v1, *args, **kwds):
    """
    Create context (with statement) to temporarily monkey patch socket.create_connection.
    """
    # Keep a handle on the old connection
    socket_create_connection = urllib3.util.connection.create_connection

    def kubernetes_create_connection(address, *args, **kwargs):
        """
        Function provided to urllib3.util.connection.create_connection.
        """
        # Ensure the request is for kubernetes
        # otherwise handle connection with original urllib3 create_connection
        dns_name = address[0]
        if isinstance(dns_name, bytes):
            dns_name = dns_name.decode()
        dns_name = dns_name.split(".")
        if dns_name[-1] != "kubernetes":
            return socket_create_connection(address, *args, **kwargs)

        # requred to be <pod-name>.pod.<namespace>.kubernetes
        if len(dns_name) not in (3, 4):
            raise RuntimeError("Unexpected kubernetes DNS name.")

        # Sploot out pieces we want to interact with
        namespace = dns_name[-2]
        name = dns_name[0]
        port = address[1]

        # This is currently only supporting pods, we have service across
        if len(dns_name) == 4:
            if dns_name[1] != "pod":
                raise RuntimeError("port-forward currently just supports pods.")
        pf = portforward(
            core_v1.connect_get_namespaced_pod_portforward,
            name,
            namespace,
            ports=str(port),
        )
        return pf.socket(port)

    try:
        # Do the monkey patch
        urllib3.util.connection.create_connection = kubernetes_create_connection
        yield
    finally:
        # Undo the monkey patch
        urllib3.util.connection.create_connection = socket_create_connection
