#!/usr/bin/env python3

import argparse
import os
import socket
import sys

# How we provide custom parameters to a flux-burst plugin
from fluxburst_gke.plugin import BurstParameters
from fluxburst.client import FluxBurst
import fluxburst.plugins as plugins

# Save data here
here = os.path.dirname(os.path.abspath(__file__))


# This is the dataclass we create with parameters for our plugin.
# Note that this will eventually come from a config file / the environment
# We originally had most of these exposed via command line, but we are
# migrating to an approach where it comes primarily from a config file.
# Thus, the only command line stuff is the project, or ephemeral list host

def get_parser():
    parser = argparse.ArgumentParser(
        description="Experimental Bursting",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument("--project", help="Google Cloud project")
    parser.add_argument(
        "--lead-host",
        help="Lead broker service hostname or ip address",
        dest="lead_host",
    )
    parser.add_argument("--lead-size", help="Lead broker size", type=int)
    parser.add_argument(
        "--lead-port", help="Lead broker service port", dest="lead_port", default=30093
    )
    parser.add_argument("--flux-operator-yaml", dest="flux_operator_yaml")
    parser.add_argument(
        "--munge-secret-name",
        help="Name of a secret to be made in the same namespace",
        default="munge.key"
    )
    parser.add_argument(
        "--munge-key",
        help="Path to munge.key",
    )
    parser.add_argument(
        "--name", help="Name for external MiniCluster", default="burst-0"
    )
    parser.add_argument(
        "--curve-cert", dest="curve_cert", default="/mnt/curve/curve.cert"
    )
    parser.add_argument(
        "--curve-cert-secret-name", default="curve-cert"
    )
    return parser


def main():
    """
    Create an external cluster we can burst to, and optionally resize.
    """
    parser = get_parser()

    # If an error occurs while parsing the arguments, the interpreter will exit with value 2
    args, _ = parser.parse_known_args()
    if not args.project:
        sys.exit("Please define your Google Cloud Project with --project")

    # Lead host and port are required. A custom broker.toml can be provided,
    # but we are having the operator create it for us
    if not args.lead_port or not args.lead_host or not args.lead_size:
        sys.exit("All of --lead-host, --lead-size, and --lead-port must be defined.")
    print(
        "Broker lead will be expected to be accessible on {args.lead_host}:{args.lead_port}"
    )

    # These checks are done by plugin, but I wanted to do them earlier too
    if args.munge_key and not os.path.exists(args.munge_key):
        sys.exit(f"Provided munge key {args.munge_key} does not exist.")
    if args.munge_key and not args.munge_secret_name:
        args.munge_secret_name = "munge-key"

    # Create the dataclass for the plugin config
    # We use a dataclass because it does implicit validation of required params, etc.
    params = BurstParameters(
        project=args.project,
        munge_key=args.munge_key,
        munge_secret_name=args.munge_secret_name,
        curve_cert_secret_name=args.curve_cert_secret_name,
        flux_operator_yaml=args.flux_operator_yaml,
        lead_host=args.lead_host,   
        lead_port=args.lead_port,     
        lead_size=args.lead_size,
        name=args.name,
    )

    # Create the flux burst client. This can be passed a flux handle (flux.Flux())
    # and will make one otherwise.
    client = FluxBurst()

    # For debugging, here is a way to see plugins available
    # from fluxburst.plugins import burstable_plugins
    # print(plugins.burstable_plugins)
    # {'gke': <module 'fluxburst_gke' from '/home/flux/.local/lib/python3.8/site-packages/fluxburst_gke/__init__.py'>}

    # Load our plugin and provide the dataclass to it!
    client.load("gke", params)

    # Sanity check loaded
    print(f'flux-burst client is loaded with plugins for: {client.choices}')

    # We are using the default algorithms to filter the job queue and select jobs.
    # If we weren't, we would add them via:
    # client.set_ordering()
    # client.set_selector()

    # Here is how we can see the jobs that are contenders to burst!
    # client.select_jobs()

    # Now let's run the burst! The active plugins will determine if they
    # are able to schedule a job, and if so, will do the work needed to
    # burst. unmatched jobs (those we weren't able to schedule) are
    # returned, maybe to do something with?
    unmatched = client.run_burst()

    print('TODO interact with the dataclass')
    import IPython
    IPython.embed()
    sys.exit()

    ensure_flux_operator_yaml(args)
    curve_cert = ensure_curve_cert(args)


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
