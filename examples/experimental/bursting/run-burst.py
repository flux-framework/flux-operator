#!/usr/bin/env python3

import argparse
import base64
import os
import socket
import sys

import flux
import flux.job
import requests

import kubescaler.utils as utils

# This will allow us to create and interact with our cluster
from kubescaler.scaler import GKECluster

# Save data here
here = os.path.dirname(os.path.abspath(__file__))

from fluxoperator.client import FluxMiniCluster
from kubernetes import client as kubernetes_client
from kubernetes import utils as k8sutils
from kubernetes.client.rest import ApiException

here = os.path.abspath(os.path.dirname(__file__))

handle = flux.Flux()

# Default flux operator yaml URL
default_flux_operator_yaml = "https://raw.githubusercontent.com/flux-framework/flux-operator/main/examples/dist/flux-operator.yaml"


# Here is our main container
def get_minicluster(
    command,
    curve_cert,
    size=4,
    tasks=16,  # nodes * cpu per node, where cpu per node is vCPU / 2
    cpu_limit=7,
    memory_limit="20G",
    flags=None,
    name=None,
    namespace=None,
    image=None,
    wrap=None,
    log_level=7,
    flux_user=None,
    lead_host=None,
    lead_port=None,
    broker_toml=None,
    munge_config_map=None,
    lead_size=None,
    lead_jobname=None,
):
    """
    Get a MiniCluster CRD as a dictionary

    Limits should be slightly below actual pod resources. The curve cert and broker config
    are required, since we need this external cluster to connect to ours!
    """
    flags = flags or "-ompi=openmpi@5 -c 1 -o cpu-affinity=per-task"
    image = image or "ghcr.io/flux-framework/flux-restful-api"
    container = {
        "image": image,
        "command": command,
        "resources": {
            "limits": {"cpu": cpu_limit, "memory": memory_limit},
            "requests": {"cpu": cpu_limit, "memory": memory_limit},
        },
    }

    # Do we have a custom flux user for the container?
    if flux_user:
        container["flux_user"] = {"name": flux_user}

    # The MiniCluster has the added name and namespace
    mc = {
        "size": size,
        "tasks": tasks,
        "namespace": namespace,
        "name": name,
        "interactive": False,
        "logging": {"zeromq": True, "quiet": False, "strict": False},
        "flux": {
            "option_flags": flags,
            "connect_timeout": "5s",
            "curve_cert": curve_cert,
            "log_level": log_level,
            # Providing the lead broker and port points back to the parent
            "bursting": {
                "lead_broker": {
                    "address": lead_host,
                    "port": int(lead_port),
                    "name": lead_jobname,
                    "size": int(lead_size),
                },
                "clusters": [{"size": size, "name": name}],
            },
        },
    }

    if munge_config_map:
        mc["flux"]["mungeConfigMap"] = munge_config_map
    if broker_toml:
        mc["flux"]["broker_config"] = broker_toml

    # eg., this would require strace "strace,-e,network,-tt"
    if wrap is not None:
        mc["flux"]["wrap"] = wrap
    return mc, container


def create_munge_configmap(path, name, namespace):
    """
    Create a binary data config map
    """
    # Configureate ConfigMap metadata
    metadata = kubernetes_client.V1ObjectMeta(
        name=name,
        namespace=namespace,
    )
    # Get File Content
    with open(path, "rb") as f:
        content = f.read()

    # base64 encoded string
    content = base64.b64encode(content).decode("utf-8")

    # Instantiate the configmap object
    return kubernetes_client.V1ConfigMap(
        api_version="v1",
        kind="ConfigMap",
        binary_data={"munge.key": content},
        metadata=metadata,
    )


def get_parser():
    parser = argparse.ArgumentParser(
        description="Experimental Bursting",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument("--project", help="Google Cloud project")
    parser.add_argument("--cluster-name", help="Cluster name", default="flux-cluster")
    parser.add_argument(
        "--machine-type", help="Google machine type", default="c2-standard-8"
    )
    parser.add_argument(
        "--cpu-limit", dest="cpu_limit", help="CPU limit", default=7, type=int
    )
    parser.add_argument(
        "--memory-limit", dest="memory_limit", help="Memory limit", default="20G"
    )
    parser.add_argument("--image", help="Container image for MiniCluster")
    parser.add_argument(
        "--lead-host",
        help="Lead broker service hostname or ip address",
        dest="lead_host",
    )
    parser.add_argument("--lead-size", help="Lead broker size")
    parser.add_argument(
        "--munge-key",
        help="Name of a config map to be made in the same namespace",
    )
    parser.add_argument(
        "--munge-config-map",
        help="Path to munge.key",
    )
    parser.add_argument(
        "--lead-port", help="Lead broker service port", dest="lead_port", default=30093
    )
    parser.add_argument(
        "--log-level",
        help="Logging level for flux",
        default=7,
        type=int,
    )
    parser.add_argument(
        "--name", help="Name for external MiniCluster", default="burst-0"
    )
    parser.add_argument(
        "--namespace", help="Namespace for external cluster", default="flux-operator"
    )
    parser.add_argument(
        "--broker-toml",
        help="Broker toml template",
    )
    parser.add_argument("--flux-operator-yaml", dest="flux_operator_yaml")
    parser.add_argument(
        "--curve-cert", dest="curve_cert", default="/mnt/curve/curve.cert"
    )
    parser.add_argument("--flux-user", help="custom flux user (defaults to flux)")
    parser.add_argument(
        "--wrap", help='arguments to flux wrap, e.g., "strace,-e,network,-tt'
    )
    return parser


def get_job_info(jobid):
    """
    Get job info based on an id

    Also retrieve the full job info and jobspec.
    This is not yet currently perfectly json serializable, need to
    handle EmptyObject if that is desired.
    """
    fluxjob = flux.job.JobID(jobid)
    payload = {"id": fluxjob, "attrs": ["all"]}
    rpc = flux.job.list.JobListIdRPC(handle, "job-list.list-id", payload)
    job = rpc.get_job()

    # Job info, timing, priority, etc.
    job["info"] = rpc.get_jobinfo().__dict__
    job["info"]["_exception"] = job["info"]["_exception"].__dict__
    job["info"]["_annotations"] = job["info"]["_annotations"].__dict__

    # the KVS will have annotations!
    kvs = flux.job.job_kvs(handle, jobid)
    job["spec"] = kvs.get("jobspec")
    return job


def is_burstable(info):
    """
    Determine if a job is explicitly labeled to be burstable
    """
    return "burstable" in info["spec"]["attributes"]["system"]


def ensure_flux_operator_yaml(args):
    """
    Ensure we are provided with the installation yaml and it exists!
    """
    # flux operator yaml default is current from main
    if not args.flux_operator_yaml:
        args.flux_operator_yaml = utils.get_tmpfile(prefix="flux-operator") + ".yaml"
        r = requests.get(default_flux_operator_yaml, allow_redirects=True)
        utils.write_file(r.content, args.flux_operator_yaml)

    # Ensure it really really exists
    args.flux_operator_yaml = os.path.abspath(args.flux_operator_yaml)
    if not os.path.exists(args.flux_operator_yaml):
        sys.exit(f"{args.flux_operator_yaml} does not exist.")


def ensure_curve_cert(args):
    """
    Ensure we are provided with an existing curve certificate we can load.
    """
    if not args.curve_cert or not os.path.exists(args.curve_cert):
        sys.exit(
            f"Curve cert (provided as {args.curve_cert}) needs to be defined and exist."
        )
    return utils.read_file(args.curve_cert)


def write_minicluster_yaml(mc):
    """
    Write the MiniCluster spec to yaml to apply
    """
    # this could be saved for reproducibility, if needed.
    minicluster_yaml = utils.get_tmpfile(prefix="minicluster") + ".yaml"
    utils.write_yaml(mc, minicluster_yaml)
    return minicluster_yaml


def main():
    """
    Create an external cluster we can burst to, and optionally resize.
    """
    parser = get_parser()

    # If an error occurs while parsing the arguments, the interpreter will exit with value 2
    args, _ = parser.parse_known_args()
    if not args.project:
        sys.exit("Please define your Google Cloud Project with --project")

    # Pull cluster name out of argument
    # TODO: likely we will start Flux with an ability to say "allow this external flux cluster"
    # and then it will have a name that can be derived from that.
    cluster_name = args.cluster_name
    print(f"📛️ New cluster name will be {cluster_name}")

    ensure_flux_operator_yaml(args)
    curve_cert = ensure_curve_cert(args)

    # Lead host and port are required
    if not args.lead_port or not args.lead_host:
        sys.exit("--lead-host and --lead-port must be defined.")
    print(
        "Broker lead will be expected to be accessible on {args.lead_host}:{args.lead_port}"
    )

    if args.munge_key and not os.path.exists(args.munge_key):
        sys.exit(f"Provided munge key {args.munge_key} does not exist.")
    if args.munge_key and not args.munge_config_map:
        args.munge_config_map = "munge-key"

    # Create a spec for what we need to burst.
    # This will be just for one moment in time, obviously there would be different
    # ways to do this (to decide when to burst, based on what metrics, etc.)
    # For now we will just assume one cluster + burst per job!
    burstable = []
    listing = flux.job.job_list(handle).get()
    for job in listing.get("jobs", []):
        info = get_job_info(job["id"])
        if not is_burstable(info):
            continue
        print(f"🧋️  Job {job['id']} is marked for bursting.")
        burstable.append(info)

    if not burstable:
        sys.exit("No jobs were found marked for burstable.")

    # Assume we just have one configuration to create for now
    # We ideally want something more elegant
    info = burstable[0]

    # Determine if the cluster exists, and if not, create it
    # For now, ensure lead broker in both is same hostname
    podname = socket.gethostname()
    hostname = podname.rsplit("-", 1)[0]

    # Try creating the cluster (this is just the GKE cluster)
    # n2-standard-8 has 4 actual cores, so 4x4 == 16 tasks
    cli = GKECluster(
        project=args.project,
        name=cluster_name,
        node_count=info["nnodes"],
        machine_type=args.machine_type,
        min_nodes=info["nnodes"],
        max_nodes=info["nnodes"],
    )
    # Create the cluster (this times it)
    try:
        cli.create_cluster()
    # What other cases might be here?
    except:
        print("🥵️ Issue creating cluster, assuming already exists.")

    # Create a client from it
    print(f"📦️ The cluster has {cli.node_count} nodes!")
    kubectl = cli.get_k8s_client()

    # Install the operator!
    try:
        k8sutils.create_from_yaml(kubectl.api_client, args.flux_operator_yaml)
        print("Installed the operator.")
    except Exception as exc:
        print(f"Issue installing the operator: {exc}, assuming already exists")

    # NOTE we previously populated a broker.toml template here, and we don't
    # need to do that anymore - the operator will generate the config

    # Assemble the command from the requested job
    command = " ".join(info["spec"]["tasks"][0]["command"])
    print(f"Command is {command}")

    # TODO: we are using defaults for now, but will update this to be likely
    # configured based on the algorithm that chooses the best spec
    minicluster, container = get_minicluster(
        command,
        name=args.name,
        memory_limit=args.memory_limit,
        cpu_limit=args.cpu_limit,
        namespace=args.namespace,
        curve_cert=curve_cert,
        broker_toml=args.broker_toml,
        tasks=info["ntasks"],
        size=info["nnodes"],
        image=args.image,
        wrap=args.wrap,
        log_level=args.log_level,
        flux_user=args.flux_user,
        lead_host=args.lead_host,
        lead_port=args.lead_port,
        munge_config_map=args.munge_config_map,
        lead_jobname=hostname,
        lead_size=args.lead_size,
    )

    # Create the namespace
    try:
        kubectl.create_namespace(
            kubernetes_client.V1Namespace(
                metadata=kubernetes_client.V1ObjectMeta(name=args.namespace)
            )
        )
    except:
        print(f"🥵️ Issue creating namespace {args.namespace}, assuming already exists.")

    # Let's assume there could be bugs applying this differently
    crd_api = kubernetes_client.CustomObjectsApi(kubectl.api_client)

    # kubectl create configmap --namespace flux-operator munge-key --from-file=/etc/munge/munge.key
    # WORKING HERE
    # TODO create from file in the same namespace?
    import IPython

    IPython.embed()

    if args.munge_key:
        cm = create_munge_configmap(
            args.munge_key, args.munge_config_map, args.namespace
        )
        try:
            kubectl.create_namespaced_config_map(
                namespace=args.namespace,
                body=cm,
            )
        except ApiException as e:
            print(
                "Exception when calling CoreV1Api->create_namespaced_config_map: %s\n"
                % e
            )

    # Create the MiniCluster! This also waits for it to be ready
    # TODO we need a check here for completed - it will hang
    # Need to fix this so it doesn't hang. We need to decide when to
    # bring down the minicluster.
    print(f"⭐️ Creating the minicluster {args.name} in {args.namespace}...")
    operator = FluxMiniCluster()
    operator.create(**minicluster, container=container, crd_api=crd_api)

    # Eventually to clean up...
    cli.delete_cluster()


if __name__ == "__main__":
    main()
