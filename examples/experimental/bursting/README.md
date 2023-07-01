# Design

### What we might need

I expect for this kind of bursting to work I will need to be able to:

 - Define the lead broker as a hostname external to the cluster AND internal (we will have two aliases for the same broker)
 - Have a custom external broker config that we generate
 - Be able to read in send the curve certificate as a configuration variable
 - Be able to create the munge.key as a secret
 - Create the service to run directly from the lead broker pod

### What we we eventually want to improve upon

We will eventually need to do the following:

 - Have an elegant way to decide:
   - When to burst
   - What jobs are marked for bursting, and how assigned to an external cluster
 - Unique external cluster names (within a cloud namespace) that are linked to their spec (content hash of broker.toml?)
 - When to configure external cluster to allow for scaling of flux (right now we set min == max so constant size)
 - Create a more scoped permission (Google service account) for running inside a cluster

### What we do now

- Allow the user to provide a custom broker.toml and curve.cert (the idea being the external clusters are created with a lead broker pointing to that service, and the same curve.cert that is on the local cluster they were created from) - these are changes to the operator CRD
- Submit jobs that are flagged for burstable (an attribute) and ask for more nodes than the cluster has (but below the max number so it doesn't immediately fail - we will want to fix this in flux because with bursting we should be able to ask for more than the potential size)
- Have a Python script that connects to the flux handle, finds the burstable jobs, and creates a minicluster spec (with the same nodes, tasks, and command - right now the machine is an argument)
- The Python script generates a broker.toml on the fly that defines the lead broker to be the service of the minicluster where it's running, and the curve.cert read directly from the filesystem (being used by the current cluster)
- The external cluster is created from the first, directly from the flux lead broker!
- The order of hosts in the broker.toml and resource spec must be consistent
- The two are networked via the exposed service on the lead broker pod

### History

This example was originally implemented as a standalone script, and that setup is preserved in [v1](v1).
The version here has been updated to use the [flux-burst](https://github.com/converged-computing/flux-burst)
module and the GKE plugin for it, [flux-burst-gke](https://github.com/converged-computing/flux-burst-gke).
