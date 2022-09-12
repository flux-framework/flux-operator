# TODO

### Design 3

 - [x] figure out where to put flux hostname / config - volume needs write
 - [ ] Are `--cores` properly set?
 - [x] debug nodes finding on another (see How it works in README.md)
 - [ ] can (and should) we use generics to reduce redudancy of code? (e.g., the `get<X>` functions) (@vsoch would like to do this)
 - [ ] what should be the proper start command for the main/worker nodes
 - [ ] By what user? I am currently root but flux is an option
 - [ ] How we can print better verbose debugging output (possibly exposed by a variable)
 - [ ] And have some solid evidence the node communication is successful (or is the job running that evidence?
 - [ ] I think if a pod dies the IP address might change, so eventually we want to test that (and may need more logic for re-updating /etc/hosts)
 - [x] debug pod containers not seeing config again (e.g., mounts not creating)
 - [ ] Should there be a min/max size for the MiniCluster CRD?
 - [x] Should the secondary (non-driver) pods have a different start command? (answer is no - with the Indexed job it's all the same command)
 - [ ] MiniCluster - how should we handle deletion / update?
 - [ ] Do we want to be able to launch additional tasks? (e.g., after the original job started)
 - [ ] Currently we have no representation of quota - we need to be able to set (and check) hard limits from the scheduler (or maybe we get that out of the box)?
 - [x] Details for etc-hosts (or will this just work? - no it won't just work)
 - [ ] klog can be changed to add V(2) to handle verbository from the command line, see https://pkg.go.dev/k8s.io/klog/v2
 - [ ] At some point we want more intelligent use of labels/selectors (I haven't really read about them yet)

### Design 2 (not currently working on)

 - [ ] pkg/util/heap should implement an actual heap
 - [x] kubebuilder should be able to provide defaults in the *_types.
 - [x] Figure out logging connected to reconciler

### Design 1 (not currently working on)

- [x] consolidate configmap functions into shared functionality (less redundancy)
- [x] Debug why the configmaps aren't being populated with the hostfile (it wasn't working with kind, worked without changes with minikube)
- [x] Figure out adding namespaces to config/samples - should be flux-operator
- [x] Each of config files written (e.g., hostname, broker, cert) should have their own types and more simply generated. The strategy right now is just temporary.
- [x] Stateful set (figure out how to create properly, doesn't seem to have pods) (figured out need to create ConfigMaps for Volumes)