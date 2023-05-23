#!/usr/bin/env python3

import argparse
import time
import json
import sys
import os

# Save data here
here = os.path.dirname(os.path.abspath(__file__))

# Insert path to import fluxcluster as local module
sys.path.insert(0, here)
from fluxcluster import FluxOperatorCluster, read_json

# Create data output directory
data = os.path.join(here, "data")


def get_parser():
    parser = argparse.ArgumentParser(
        description="K8s Scaling Experiment Runner",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument(
        "cluster_name", nargs="?", help="Cluster name suffix", default="flux-cluster"
    )
    parser.add_argument(
        "--outdir",
        help="output directory for results",
        default=data,
    )
    parser.add_argument(
        "--experiment", help="Experiment name (defaults to script name)", default=None
    )
    parser.add_argument(
        "--start-iter", help="start at this iteration", type=int, default=0
    )
    parser.add_argument(
        "--end-iter", help="end at this iteration", type=int, default=10, dest="iters"
    )
    parser.add_argument(
        "--max-node-count", help="maximum node count", type=int, default=32
    )
    parser.add_argument(
        "--start-node-count",
        help="start at this many nodes and go up",
        type=int,
        default=1,
    )
    parser.add_argument(
        "--machine-type-memory-gb",
        help="memory of desired machine, in GB",
        type=int,
        default=32,
    )
    parser.add_argument(
        "--machine-type-vcpu",
        help="Virtual cpu (vcpu) on Google Cloud",
        type=int,
        default=8,
    )
    parser.add_argument(
        "--machine-type", help="Google Cloud machine type", default="c2-standard-8"
    )
    parser.add_argument(
        "--increment", help="Increment by this value", type=int, default=1
    )
    parser.add_argument(
        "--down", action="store_true", help="Test scaling down", default=False
    )
    return parser


def main():
    """
    This experiment will test scaling a cluster, three times, each
    time going from 2 nodes to 32. We want to understand if scaling is
    impacted by cluster size.
    """
    parser = get_parser()

    # If an error occurs while parsing the arguments, the interpreter will exit with value 2
    args, _ = parser.parse_known_args()

    # Pull cluster name out of argument
    cluster_name = args.cluster_name

    # Derive the experiment name, either named or from script
    experiment_name = args.experiment
    if not experiment_name:
        experiment_name = sys.argv[0].replace(".py", "")
    time.sleep(2)

    # Shared tags for logging and output
    if args.down:
        direction = "decrease"
        tag = "down"
    else:
        direction = "increase"
        tag = "up"

    # Update cluster name to include tag and increment
    experiment_name = f"{experiment_name}-{tag}-{args.increment}"
    print(f"📛️ Experiment name is {experiment_name}")

    # Prepare an output directory, named by cluster
    outdir = os.path.join(args.outdir, experiment_name, cluster_name)
    if not os.path.exists(outdir):
        print(f"📁️ Creating output directory {outdir}")
        os.makedirs(outdir)

    # Define stopping conditions for two directions
    def less_than_max(node_count):
        return node_count <= args.max_node_count

    def greater_than_zero(node_count):
        return node_count > 0

    # Update cluster name to include experiment name
    cluster_name = f"{experiment_name}-{cluster_name}"
    print(f"📛️ Cluster name is {cluster_name}")

    # Name cannot be greater than 40
    if len(cluster_name) > 40:
        print("⚠️ Warning: cluster name is too long, must be <= 40. Will cut!")
        cluster_name = cluster_name[:39]
        print(f"📛️ Cluster name is {cluster_name}")

    # Create 10 clusters, each going up to 32 nodes
    for iter in range(args.start_iter, args.iters):
        results_file = os.path.join(outdir, f"scaling-{iter}.json")

        # Start at the max if we are going down, otherwise the starting count
        node_count = args.max_node_count if args.down else args.start_node_count
        print(
            f"⭐️ Creating the initial cluster, iteration {iter} with size {node_count}..."
        )
        cli = FluxOperatorCluster(
            "llnl-flux",
            name=cluster_name,
            machine_type=args.machine_type,
            node_count=node_count,
            max_nodes=args.max_node_count,
            machine_type_memory_gb=args.machine_type_memory_gb,
            machine_type_vcpu=args.machine_type_vcpu,
        )

        # Load a result if we have it
        if os.path.exists(results_file):
            result = read_json(results_file)
            cli.times = result["times"]

        # Create the cluster (this times it)
        res = cli.create_cluster()
        print(f"📦️ The cluster has {res.initial_node_count} nodes!")

        # Flip between functions to decide to keep going based on:
        # > 0 (we are decreasing from the max node count)
        # <= max nodes (we are going up from a min node count)
        keep_going = less_than_max
        if args.down:
            keep_going = greater_than_zero

        # Continue scaling until we reach stopping condition
        while keep_going(node_count):
            old_size = node_count

            # Are we doing down or up?
            if args.down:
                node_count -= args.increment
            else:
                node_count += args.increment

            print(
                f"⚖️ Iteration {iter}: scaling to {direction} by {args.increment}, from {old_size} to {node_count}"
            )
            start = time.time()

            # Slightly different logic for scaling
            if args.down:
                res = cli.scale_down(node_count)
            else:
                res = cli.scale_up(node_count)

            end = time.time()
            seconds = round(end - start, 3)
            cli.times[f"scale_{tag}_{old_size}_to_{node_count}"] = seconds
            print(
                f"📦️ Scaling from {old_size} to {node_count} took {seconds} seconds, and the cluster now has {res.initial_node_count} nodes!"
            )

            # Save the times as we go
            print(json.dumps(cli.data, indent=4))
            cli.save(results_file)

        # Delete the cluster and clean up
        cli.delete_cluster()
        print(json.dumps(cli.data, indent=4))
        cli.save(results_file)


if __name__ == "__main__":
    main()
