import json
import sys
import os
import time
import argparse
import requests
import threading
from flask import Flask, jsonify, request

try:
    import flux
    import flux.job
except:
    sys.exit("Cannot import flux - are the Flux Python bindings on the path?")

try:
    # This is flux-sched fluxion
    from fluxion.resourcegraph.V1 import FluxionResourceGraphV1

    # This is from the fluxion service grpc
    from fluxion.protos import fluxion_pb2
    from fluxion.client import FluxionClient
    from flux.hostlist import Hostlist
except:
    sys.exit("Cannot import fluxion - is the fluxion module enabled?")

# Global fluxion client to receive server submit
ctrl = None
app = Flask(__name__)


@app.route("/submit", methods=["POST"])
def submit_job():
    """
    Submit a job to the running Fluxion server.
    """
    global ctrl

    data = request.get_json()
    print(data)
    for required in ["command", "cpu", "container"]:
        if required not in data or not data[required]:
            return jsonify({"error": f"{required} is required."})

    response = ctrl.submit_job(data["container"], data["cpu"], data["command"])
    return jsonify(response)


# We assume a jobspec on a flux operator node is only customizable via nodes and cores
# We will hard code things for now (e.g., duration) and these could be customized
jobspec_template = """
{
    "version": 1,
    "resources": [
        {
            "type": "slot",
            "count": 1,
            "label": "default",
            "with": [
                {
                    "type": "core",
                    "count": %s
                }
            ]
        }
    ],
    "attributes": {
        "system": {
            "duration": 3600
        }
    },
    "tasks": [
        {
            "command": [],
            "slot": "default",
            "count": {
                "per_slot": 1
            }
        }
    ]
}
"""


class FluxionController:
    """
    This is a controller that can:

    1. Detect the individual application brokers running in the MiniCluster
    2. Create a single graph to schedule to in the higher level fluxion service
    3. Discover the broker sockets and create handles for each.
    4. Receive a request for work and schedule properly!

    The high level idea is that we are sharing resources between the different
    application containers (each with a different broker that isn't aware of the
    others) and we don't want to oversubscribe. The Fluxion controller will handle
    this orchestration, scheduling work with the fluxion service and receiving
    callbacks from the individual clusters to determine when work is done.
    """

    def __init__(
        self,
        resource_dir=None,
        socket_dir=None,
        fluxion_host=None,
        meta_dir=None,
        heartbeat_seconds=1,
    ):
        # These paths assume a flux operator install
        self.socket_dir = socket_dir or "/mnt/flux/view/run/flux"
        self.resource_dir = resource_dir or "/mnt/flux/view/etc/flux/system"
        self.meta_dir = meta_dir or "/mnt/flux/view/etc/flux/meta"

        # This is the default for the headles service
        self.fluxion_host = (
            fluxion_host
            or "flux-sample-services.flux-service.default.svc.cluster.local:4242"
        )

        # This is imperfect - but we will keep a set of job ids for each container
        self.jobs = {}

        # How often to check for completed jobs
        self.heartbeat_seconds = heartbeat_seconds
        self.handles = {}
        self.resources = {}
        self.containers = {}

        self.discover_containers()
        self.discover_sockets()
        self.discover_resources()

    def populate_jobs(self):
        """
        Given running queues, populate with current jobs
        """
        pass
        # TODO how do we do this? We essentially need to restore state
        import IPython

        IPython.embed()

    def discover_containers(self):
        """
        Determine the container names and indices
        """
        for container in os.listdir(self.meta_dir):
            idx, container = container.split("-", 1)
            print(f"⭐️ Found application {container.rjust(10)}: index {idx}")
            self.containers[int(idx)] = container

    def discover_resources(self):
        """
        Discover physical node resources.

        Each container application is going to provide a unique R file, and this
        is done in the case that we can/want to vary this in the future. However,
        for the time being these are essentially the same so we can just read in
        the first.
        """
        for resource_file in os.listdir(self.resource_dir):
            if not resource_file.startswith("R-"):
                continue
            _, idx = resource_file.split("-", 1)

            # Index based on the container
            container = self.containers[int(idx)]
            self.resources[container] = read_json(
                os.path.join(self.resource_dir, resource_file)
            )

        # If we dont' have any resources, bail out early - something is wrong
        if not self.resources:
            sys.exit(
                "No resource files found in {self.resource_dir} - this should not happen."
            )

    def discover_sockets(self):
        """
        Discover sockets to create a flux handle to each

        We read in the associated container name via the meta directory in the
        flux install, which is created by the flux operator.
        """
        for socket_path in os.listdir(self.socket_dir):
            # In practice there should not be anything else in here
            if "local" not in socket_path:
                continue

            # The socket has the index for the container in it
            _, idx = socket_path.split("-", 1)

            # Use it to identify the container...
            container = self.containers[int(idx)]
            socket_fullpath = os.path.join(self.socket_dir, socket_path)

            # And generate the handle!
            uri = f"local://{socket_fullpath}"
            handle = flux.Flux(uri)
            self.handles[container] = handle

    def init_fluxion(self):
        """
        Connect to the fluxion service and create the graph.
        """
        # Grab the first R to generate the resource graph from
        # They are all the same
        key = list(self.resources.keys())[0]
        rv1 = self.resources[key]
        graph = FluxionResourceGraphV1(rv1)

        # Dump of json graph format for fluxion
        jgf = graph.to_JSON()

        # Init the fluxion graph - it only sees one of the entire cluster
        self.cli = FluxionClient(host=self.fluxion_host)

        # Fluxion spits out an error that properties must be an object or null
        for node in jgf["graph"]["nodes"]:
            if "properties" in node["metadata"] and not node["metadata"]["properties"]:
                node["metadata"]["properties"] = {}

        response = self.cli.init(json.dumps(jgf))
        if response.status == fluxion_pb2.InitResponse.ResultType.INIT_SUCCESS:
            print("✅️ Init of Fluxion resource graph success!")
        else:
            sys.exit(f"Issue with init, return code {response.status}")

        # Now run indefinitely, at least until we are done with the cluster
        t1 = threading.Thread(target=self.run)
        t1.start()
        app.run(host="0.0.0.0")

    def run(self):
        """
        Run fluxion, meaning we basically:

        1. Check over known submit jobs for each handle.
        2. When they are done on a cluster, cancel in the overhead graph.
        This is obviously imperfect in terms of state. What we can do to
        prevent race conditions is to ensure that a job is running when
        we submit it, that way we don't have two different brokers fighting
        for the same resources.
        """
        while True:
            for container, handle in self.handles.items():
                jobs = []
                for jobset in self.jobs.get(container, []):
                    # Get the status of the job from the handle
                    info = flux.job.get_job(handle, jobset["container"])
                    if info["result"] == "COMPLETED":
                        print(f"👉️ Job on {container} {jobset['fluxion']} is complete.")
                        self.cancel(jobset["fluxion"])
                        continue
                    # Otherwise add back to jobs set
                    jobs.append(jobset)
                self.jobs[container] = jobs

            # Do a sleep between the timeout
            time.sleep(self.heartbeat_seconds)

    def cancel(self, jobid):
        """
        Cancel a fluxion jobid
        """
        # An inactive RPC cannot cancel
        try:
            response = self.cli.cancel(jobid=jobid)
            if response.status == fluxion_pb2.CancelResponse.ResultType.CANCEL_SUCCESS:
                print(f"✅️ Cancel of jobid {jobid} success!")
            else:
                print(f"Issue with cancel, return code {response.status}")
        except:
            print(f"✅️ jobid {jobid} is already inactive.")

    def submit_job(self, container, cpu_count, command):
        """
        Demo of submitting a job. We will want a more robust way to do this.
        TODO: add working directory, duration, environment, etc.

        This currently just asks for the command and total cores across nodes.
        We let fluxion decide how to distribute that across physical nodes.
        """
        if not cpu_count:
            sys.exit("A cpu count is required.")
        if not container:
            sys.exit("An application target container is required (--container)")
        if not command:
            sys.exit("Please provide a command to submit")

        # They are asking for a broker container handle that doesn't exist
        if container not in self.handles:
            choices = ",".join(list(self.handles.keys()))
            sys.exit(
                f"Application container handle for {container} does not exist - choices are {choices}."
            )

        # Our broker hook to the container
        handle = self.handles[container]

        # Generate the jobspec, and see if we can match
        jobspec = json.loads(jobspec_template % str(cpu_count))
        print(f"🙏️ Requesting to submit: {' '.join(command)}")
        jobspec["tasks"][0]["command"] = command

        # This asks fluxion if we can schedule it
        self.cli = FluxionClient(host=self.fluxion_host)
        response = self.cli.match(json.dumps(jobspec))
        if response.status == fluxion_pb2.MatchResponse.ResultType.MATCH_SUCCESS:
            print("✅️ Match of jobspec to Fluxion graph success!")
        else:
            msg = (
                f"Issue with match, return code {response.status}, cannot schedule now"
            )
            print(msg)
            return {"error": msg}

        # We need the exact allocation to pass forward to the container broker
        alloc = json.loads(response.allocation)

        # https://flux-framework.readthedocs.io/projects/flux-rfc/en/latest/spec_31.html
        # We are going to use ranks instead of hosts, since that is matched here
        nodes = [
            x["metadata"]["name"]
            for x in alloc["graph"]["nodes"]
            if x["metadata"]["type"] == "node"
        ]
        ranks = [x.rsplit("-", 1)[-1] for x in nodes]

        # With the bypass plugin we can give a resource specification exactly to run
        # https://flux-framework.readthedocs.io/en/latest/faqs.html#how-can-i-oversubscribe-tasks-to-resources-in-flux
        # https://flux-framework.readthedocs.io/projects/flux-rfc/en/latest/spec_20.html
        # We cannot use constraint because we cannot limit cores

        # Create a constraint with AND for each host and the exact ranks assigned
        # Note that this currently isn't supported so we just give the hostlist
        # We need to be able to provide the exact hosts and cores on them.
        resource_spec = {
            "version": 1,
            "execution": {
                "R_lite": [],
                "starttime": 0.0,
                "expiration": 0.0,
                "nodelist": ["flux-sample-[0-1]"],
            },
        }

        # flux jobtap load system.alloc-bypass.R
        # Example R_lite list: {'rank': '0', 'children': {'core': '0-4'}}, {'rank': '1', 'children': {'core': '6-8'}}
        # nodelist: ['flux-sample-[0-1]']

        r_lite = []
        for node in nodes:
            ranks = [
                str(x["metadata"]["id"])
                for x in alloc["graph"]["nodes"]
                if x["metadata"]["type"] == "node"
            ]
            cores = [
                str(x["metadata"]["id"])
                for x in alloc["graph"]["nodes"]
                if x["metadata"]["type"] == "core"
                and node in x["metadata"]["paths"]["containment"]
            ]
            r_lite.append(
                {"rank": ",".join(ranks), "children": {"core": ",".join(cores)}}
            )

        hl = Hostlist(handle.attr_get("hostlist"))
        hostlist = [hl[int(x)] for x in ranks]
        resource_spec["execution"]["nodelist"] = hostlist
        resource_spec["execution"]["R_lite"] = r_lite

        # Set the resource_spec on the plugin
        jobspec["attributes"]["system"]["alloc-bypass"] = {"R": resource_spec}

        # Now we need to submit to the actual cluster, and store the mapping of our
        # fluxion jobid to the cluster jobid.
        fluxjob = flux.job.submit_async(handle, json.dumps(jobspec))

        # Wait until it's running (and thus don't submit other jobs)
        # This assumes running one client to submit, and prevents race
        jobid = fluxjob.get_id()
        while True:
            info = flux.job.get_job(handle, jobid)

            # These should be all states that come before running or finished
            if info["state"] in ["DEPEND", "PRIORITY", "SCHED"]:
                time.sleep(self.heartbeat_seconds)
                continue
            break

        # Keep a record of the fluxion job id
        if container not in self.jobs:
            self.jobs[container] = []
        self.jobs[container].append({"fluxion": response.jobid, "container": jobid})

        # Update the info and return back
        info["container"] = container
        info["fluxion"] = response.jobid
        return info


def read_json(filename):
    """
    Read content from a json file
    """
    with open(filename, "r") as fd:
        content = json.loads(fd.read())
    return content


def write_json(obj, filename):
    """
    Write content to a json file
    """
    with open(filename, "w") as fd:
        fd.write(json.dumps(obj, indent=4))


def get_parser():
    parser = argparse.ArgumentParser(
        description="Fluxion Application Scheduler Controller",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    subparsers = parser.add_subparsers(
        help="actions",
        title="actions",
        description="Fluxion application scheduler controller subparsers",
        dest="command",
    )
    start = subparsers.add_parser(
        "start", description="initialize and start fluxion (only do this once)!"
    )
    submit = subparsers.add_parser(
        "submit",
        description="submit a JobSpec for a specific application broker",
        formatter_class=argparse.RawTextHelpFormatter,
    )

    submit.add_argument("--container", help="Application container to submit to")
    submit.add_argument("--cpu", help="Total CPU across N nodes to request under slot")
    submit.add_argument(
        "--host",
        help="MiniCluster hostname running the service",
        default="flux-sample-0.flux-service.default.svc.cluster.local:5000",
    )

    for command in [start, submit]:
        command.add_argument("--fluxion-host", help="Fluxion service host")
        command.add_argument(
            "--resource-dir", help="MiniCluster resource (R) directory"
        )
        command.add_argument("--socket-dir", help="MiniCluster socket directory")
        command.add_argument("--meta-dir", help="MiniCluster Flux meta directory")
    return parser


def main():
    """
    Create a fluxion graph handler for the application broker cluster.
    """
    global ctrl

    parser = get_parser()
    args, command = parser.parse_known_args()
    ctrl = FluxionController(
        socket_dir=args.socket_dir,
        resource_dir=args.resource_dir,
        meta_dir=args.meta_dir,
    )

    if not args.command:
        parser.print_help()

    # Init creates the resource graph and must be called once first
    if args.command == "start":
        ctrl.init_fluxion()

    # The submit issues a post to the running server
    elif args.command == "submit":
        response = requests.post(
            f"http://{args.host}/submit",
            json={"command": command, "cpu": args.cpu, "container": args.container},
        )
        print(response.json())


if __name__ == "__main__":
    main()
