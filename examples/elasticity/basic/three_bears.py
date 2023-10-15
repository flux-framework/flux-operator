# This one is too small...
from kubernetes import client, config
import time
import flux
import os

print()
time.sleep(2)
print("Elasticity anyone? Story anyone? Bueller? BUELLER?!")
time.sleep(3)
print("BOOTING UP THE STORY! 🥾️")
time.sleep(5)


def msg(message, sleep=3):
    """
    Print a message to the screen
    """
    print(message)
    time.sleep(sleep)
    print()


def list_pods(v1, states=None):
    """
    List pods with a newline, ensuring we only show pods in a particular state.
    """
    print("Hello pods... who is out there?")
    pods = v1.list_namespaced_pod("default", watch=False)
    states = states or ["Running"]
    for pod in pods.items:
        if pod.status.phase not in states:
            continue
        print(
            "%s\t%s\t%s"
            % (pod.status.pod_ip, pod.metadata.namespace, pod.metadata.name)
        )
    print()


def get_minicluster_spec(crd_api, jobname):
    """
    Get the MiniCluster spec from the API server
    """
    resource = crd_api.list_namespaced_custom_object(
        group="flux-framework.org",
        version="v1alpha2",
        namespace="default",
        plural="miniclusters",
    )
    for item in resource["items"]:
        if item["metadata"]["name"] == jobname:
            msg("Oh, I think I found it!")
            print_spec(item["spec"])
            break
    return item


def resize_minicluster(crd_api, jobname, size):
    """
    Resize the MiniCluster via a patch
    """
    minicluster_patch = {"spec": {"size": size}}
    return crd_api.patch_namespaced_custom_object(
        version="v1alpha2",
        namespace="default",
        plural="miniclusters",
        name=jobname,
        body=minicluster_patch,
        group="flux-framework.org",
    )


def print_spec(item, prefix=""):
    if isinstance(item, dict):
        print_spec_dict(item)
    else:
        print_spec_list(item, prefix)


def print_spec_list(item, prefix=""):
    for i in item:
        print_spec(i, prefix)


def print_spec_dict(item, prefix=""):
    for k, v in item.items():
        if k == "containers":
            print_spec(v, "  containers:")
        else:
            print(f"⭐️ {prefix}{k.rjust(10)}: {v}")
    time.sleep(4)


def main():
    config.load_incluster_config()
    msg("Hello there! 👋️ I'm Gopherlocks! 👱️")
    msg("Oh my, am I in a container?! 👱️")
    msg("Let's take a look around... who else is here?\n 👀️")

    v1 = client.CoreV1Api()
    list_pods(v1)

    # The hostname is always in the environment, and this can also
    # give us the job name by removing the index!
    hostname = os.environ.get("HOSTNAME")
    jobname = hostname.rsplit("-", 1)[0]
    msg(f"I see it over there! I'm running in a job called {jobname}. 🌀️")
    print(flux.Flux())
    msg("Oh hi Flux, I guess you are here too. 👋️")
    msg("Please don't lay a stinky one, I know how you job managers get! 💩️")
    msg("So hmm. I think I'm running in a Flux Operator MiniCluster. 😎️")
    msg("Just a guess! 🤷️", 2)
    msg("At least it is not three bears, har har har. 🐻️ 🐻️ 🐻️", 2)
    msg("I wonder if I can find the spec for the cluster I AM IN RIGHT NOW...")
    crd_api = client.CustomObjectsApi()

    # Get the spec (CRD) directly from the API server
    item = get_minicluster_spec(crd_api, jobname)

    # This is getting the current size and showing how to scale up
    size = item["spec"]["size"]
    print()
    msg(f"Oh my, is it a bit, tight in here? A size {size}?!")
    msg("Let's see what I can do about that...")

    # Make the MiniCluster larger!
    res = resize_minicluster(crd_api, jobname, 4)

    # This is showing how to scale down
    # In these runs, it takes a while for the pod to terminate (despite Flux)
    # seeing it as gone very quickly) so we want to sleep enough so that
    # the pod is completely gone.
    msg("Did that work? Hello out there... do we have more friends? 🍤️")
    list_pods(v1)
    msg("Oh my, we have FOUR friends!! I'm so happy! 😹️")

    # Do this early since pods terminate a bit slowly
    # But they disconnect from flux almost immediately
    res = resize_minicluster(crd_api, jobname, 3)

    # Normally you'd list with some backoff until it was terminated... but we just wait a long time here :)
    msg(
        "But actually I wanted to play some Mario Kart but I only have 4 controllers... 🕹️",
        7,
    )

    # And here is where we learn Gopherlocks lacks social skills :)
    msg("Sorry one of you has to leave!!... 😭️", 6)
    msg("I know I'm a terrible person. 👿️", 7)
    msg("I feel so bad. How many do we have now?")
    msg(
        "** DRAMATIC PAUSE FOR STORY ** but actually to wait for pod to terminate :)",
        10,
    )
    list_pods(v1)
    msg("NOICE!! TIME TO DESTROY YOU IN MARIO KART! 💪️")
    msg('"Player select: Peach." 🍑️')
    msg("Hey now, do not judge! 😜️")
    time.sleep(4)


if __name__ == "__main__":
    main()

