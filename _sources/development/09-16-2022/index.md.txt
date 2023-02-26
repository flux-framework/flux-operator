# Design 3.2

This design is a simple design based around a single custom resource definition

## Summary

 - A **MiniCluster** is an [indexed job](https://kubernetes.io/docs/tasks/job/indexed-parallel-processing-static/) so we can create N copies of the "same" base containers (each with Flux, and the connected workers in our cluster)
 - The flux config is written to a volume at `/etc/flux/config` (created via a config map) as a brokers.toml file.
 - The startup script "wait-X.sh" handles setting up Flux and checking that the container meets all requirements.
   - If a command is provided, we give it to Flux directly (suggested), otherwise we start a Flux RestFul server to interact with
   - The container for the flux runner is expected to already have a munge.key in `/etc/munge/munge.key`. This will be the same across pods generated given generation from the same container.
 - To generate the curve certificate (`/etc/curve/curve.cert`) we use the flux runner container in a one-off pod to run `flux-keygen` and retrieve the output in the log. We then write this curve.cert as a Config Map to the indexed job pods. We do this beceause generating it natively in Go would require other libraries on the host for ZeroMQ.
 - Networking of the pods works by way of exposing a service that includes the Pod subdomain. We add fully qualified domain names to the pods so that the `hostname` command matches the full name, and Flux is given the full names in its broker.toml.
 - The main pod either runs `flux start` with a web service (creating a persistent "Mini Cluster" or `flux start` with a specific command (if provided in the CRD) in which case the command runs, and the jobs finish and the cluster goes away.

This means that:

- A MiniCluster is a CRD that includes containers, a command (optional for the ephemeral use case), and size
- Creating a MiniCluster first creates Config Maps, Volumes, and secrets and then an indexed job with pods that use them
- Index 0 is "special" in that is creates the main shared assets, and launches the main command or server (the others start with `flux start` and a sleep and essentially register to the cluster.
- Flux is required in the running container image (from the user), however the Flux RESTful API is not (it is installed when the server comes up).
- For the persistent case, jobs can be submit until the cluster is stopped or deleted. For the ephemeral case (providing the command) the job runs and the cluster goes away.

- [Link on Excalidraw](https://excalidraw.com/#json=3p1bpgBFeNWpqUJjrxDmi,wZPk1I0FHI4POAAJfIdNBg)
![design-three-team.png](design-three-team.png)
