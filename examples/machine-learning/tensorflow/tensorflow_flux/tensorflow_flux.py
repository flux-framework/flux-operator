from __future__ import print_function
from __future__ import absolute_import
from __future__ import division

import socket
import os


def tf_config_from_flux(ps_number, cluster_size=4, job_name="flux-sample", port_number=2222):
    """
    Creates configuration for a distributed tensorflow session 
    from environment variables provided by the the Flux operator.
    
    @param: ps_number number of parameter servers to run
    @param: job_name name of minicluster job
    @param: cluster_size number of nodes in cluster
    @param: port_number port number to be used for communication
    @return: a tuple containing cluster with fields cluster_spec,
             task_name and task_id 
    """
    nodelist = os.environ.get("FLUX_JOB_NODELIST") or ["%s-%s.flux-service.default.svc.cluster.local" %(job_name, i) for i in range(cluster_size)]
    nodename = os.environ.get("FLUX_NODENAME") or "%s.flux-service.default.svc.cluster.local" % socket.gethostname()
    num_nodes = int(os.getenv("FLUX_NUM_NODES") or cluster_size)
    
    if len(nodelist) != num_nodes:
        raise ValueError("Number of flux nodes {} not equal to {}".format(len(nodelist), num_nodes))
    
    if nodename not in nodelist:
        raise ValueError("Nodename({}) not in nodelist({}). This should not happen! ".format(nodename,nodelist))
    
    ps_nodes = [node for i, node in enumerate(nodelist) if i < ps_number]
    worker_nodes = [node for i, node in enumerate(nodelist) if i >= ps_number]
    
    if nodename in ps_nodes:
        my_job_name = "ps"
        my_task_index = ps_nodes.index(nodename)
    else:
        my_job_name = "worker"
        my_task_index = worker_nodes.index(nodename)
    
    worker_sockets = [":".join([node, str(port_number)]) for node in worker_nodes]
    ps_sockets = [":".join([node, str(port_number)]) for node in ps_nodes]
    cluster = {"worker": worker_sockets, "ps" : ps_sockets}
    
    return cluster, my_job_name, my_task_index
