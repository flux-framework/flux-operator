#!/usr/bin/env python3

import argparse
import os
import sys
import time

# How we provide custom parameters to a flux-burst plugin
from fluxburst_gke.plugin import BurstParameters
from fluxburst.client import FluxBurst

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
        default="munge.key",
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
    parser.add_argument("--curve-cert-secret-name", default="curve-cert")
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
    # import fluxburst.plugins as plugins
    # print(plugins.burstable_plugins)
    # {'gke': <module 'fluxburst_gke' from '/home/flux/.local/lib/python3.8/site-packages/fluxburst_gke/__init__.py'>}

    # Load our plugin and provide the dataclass to it!
    client.load("gke", params)

    # Sanity check loaded
    print(f"flux-burst client is loaded with plugins for: {client.choices}")

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
    print("Sleeping for a few minutes so you can look around...")
    time.sleep(360)

    # Get a handle to the plugin so we can cleanup!
    plugin = client.plugins["gke"]
    plugin.cleanup()

if __name__ == "__main__":
    main()
