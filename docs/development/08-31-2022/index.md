# Design 1

This was Dan's original design!

We don't have an excalidraw and associated page because the design was based off of shared yaml files (in a link I can no longer find).
We create a new custom controller that listens for MiniCluster resources. When a new MiniCluster is created, the controller goes through the following steps:

 1. Create ConfigMap that contains:
   - A broker.toml file that lists the pods in the worker StatefulSet plus the Flux rank 0 (in the form of flux-workers-0, flux-workers-1, ...) which defines their rank in the Flux TBON
   - /etc/hosts with a mapping between pod names and IPs (should be able to generate this before the StatefulSet  or IndexedJob per this MPI Operator line which occurs before the StatefulSet creation
 2. Create the worker StatefulSet that contains the desired replicas minus 1, as rank 0 also does work. If we use IndexedJob  we can create all desired replicas
 3. Wait for the worker pods to enter Running state.
 4. Create the launcher Job.
 5. After the launcher job finishes, set the replicas to 0 in the worker StatefulSet, or delete the IndexedJob
