#!/usr/bin/env python3

import argparse
import os
import sys

# How we provide custom parameters to a flux-burst plugin
from fluxburst_compute_engine.plugin import BurstParameters
from fluxburst.client import FluxBurst

# Save data here
here = os.path.dirname(os.path.abspath(__file__))


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
    parser.add_argument("--lead-hostnames", help="Custom hostnames for cluster")
    parser.add_argument(
        "--lead-port", help="Lead broker service port", dest="lead_port", default=30093
    )
    parser.add_argument(
        "--munge-key",
        help="Path to munge.key",
    )
    parser.add_argument(
        "--curve-cert", dest="curve_cert", default="/mnt/curve/curve.cert"
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

    # Lead port, lead host, and lead size and existing munge key should be
    # checked / validated by the plugin

    # Create the dataclass for the plugin config
    # We use a dataclass because it does implicit validation of required params, etc.
    params = BurstParameters(
        project=args.project,
        munge_key=args.munge_key,
        curve_cert=args.curve_cert,
        lead_host=args.lead_host,
        lead_port=args.lead_port,
        lead_hostnames=args.lead_hostnames,
        # This is a single VM that has flux installed
        # and the build from terraform-gcp/basic/bursted
        compute_family="flux-fw-bursted-x86-64",
        terraform_plan_name="burst",
        compute_machine_type="n2-standard-4",
    )

    # Create the flux burst client. This can be passed a flux handle (flux.Flux())
    # and will make one otherwise. Note that by default mock=False
    client = FluxBurst()

    # For debugging, here is a way to see plugins available
    # import fluxburst.plugins as plugins
    # print(plugins.burstable_plugins)
    # {'gke': <module 'fluxburst_gke' from '/home/flux/.local/lib/python3.8/site-packages/fluxburst_gke/__init__.py'>}

    # Load our plugin and provide the dataclass to it!
    client.load("compute_engine", params)

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

    assert not unmatched

    # Get a handle to the plugin so we can cleanup!
    plugin = client.plugins["compute_engine"]
    input("Press Enter to when you are ready to destroy...")
    plugin.cleanup()


if __name__ == "__main__":
    main()
