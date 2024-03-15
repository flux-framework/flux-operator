import json
import os
import sys
import time

from flux_restful_client.main import get_client

for envar in ["FLUX_USER", "FLUX_TOKEN", "FLUX_SECRET_KEY"]:
    if not os.getenv(envar):
        sys.exit(f"Please export {envar} in the environment")

# We require the hostname / port as the only argument
if len(sys.argv) < 2:
    sys.exit("Please include the <ipaddress>:<port> of your restful API service")

host = sys.argv[-1]

# Give the client the host, the rest come from envars
cli = get_client(host=host)

print("🐭️ What are we going to do tonight, Brain?")
time.sleep(3)
print("🐀️ The same thing we do every night, Pinky...")
time.sleep(3)
print("🐀️ Try to submit jobs to a remote Flux instance! 🌀️")
time.sleep(2)
print("     (diabolical laugher) 🦹️\n")
time.sleep(3)
print(" -- Cluster nodes -- ")
nodes = cli.list_nodes()
print(json.dumps(nodes, indent=4))
print()


def submit_job(command, nodes=1, sleep_seconds=3, tasks=1):
    """
    Example function to submit a job, wait some time for it
    and report output. There is a lot more functionaliy that
    can be exposed here!
    """
    print(f" -- Submit {command} to {nodes} node -- ")
    res = cli.submit(command, num_nodes=nodes, num_tasks=tasks)
    time.sleep(3)
    print(res)
    print(f"Flux job id {res['id']}\n")
    time.sleep(sleep_seconds)

    print(" -- Flux job metadata -- ")
    res = cli.jobs(res["id"])
    print(json.dumps(res, indent=4))
    time.sleep(5)

    print(" -- Output --")
    for line in cli.output(res["id"])["Output"]:
        print(line, end="")


submit_job("hostname")
time.sleep(5)
print("     (MOOOOOAR!) 🦹️\n")
time.sleep(5)
submit_job("hostname", nodes=4, tasks=4)
print("     (BWAHAHAH!) 🦹️\n")
time.sleep(3)
print("🐭️ Ok brain, but as long as we can get tacos after 🌮️")
