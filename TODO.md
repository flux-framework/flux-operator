# TODO

- [x] Figure out adding namespaces to config/samples - should be flux-operator
- [ ] ConfigMap -> Name doesn't match any [spec I can find](https://github.com/kubernetes/api/blob/e9a69791a998e7ead3a95fec1e420d52d62aa0f8/core/v1/types.go#L1605).
- [ ] Each of config files written (e.g., hostname, broker, cert) should have their own types and more simply generated. The strategy right now is just temporary.
- [ ] Stateful set (figure out how to create properly, doesn't seem to have pods)

## Design

We create a new custom controller that listens for FluxJob resources. When a new FluxJob is created, the controller goes through the following steps:

 1. Create ConfigMap that contains:
   - A broker.toml file that lists the pods in the worker StatefulSet plus the Flux rank 0 (in the form of flux-workers-0, flux-workers-1, ...) which defines their rank in the Flux TBON
   - /etc/hosts with a mapping between pod names and IPs (should be able to generate this before the StatefulSet  or IndexedJob per this MPI Operator line which occurs before the StatefulSet creation
 2. Create the worker StatefulSet that contains the desired replicas minus 1, as rank 0 also does work. If we use IndexedJob  we can create all desired replicas
 3. Wait for the worker pods to enter Running state.
 4. Create the launcher Job.
 5. After the launcher job finishes, set the replicas to 0 in the worker StatefulSet, or delete the IndexedJob 
