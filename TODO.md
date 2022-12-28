# TODO

 - [ ] interface still needs debugging for 2+ processes
 - [ ] make deployments own section of docs
 - [ ] try webhooks again

### Design 3

 - [ ] can (and should) we use generics to reduce redudancy of code? (e.g., the `get<X>` functions) (@vsoch would like to do this!)
 - [ ] Events: deletion should clean up, and update should not be allowed (given rank 0 started)
 - [ ] Currently we have no representation of quota - we need to be able to set (and check) hard limits from the scheduler (or maybe we get that out of the box)?
 - [ ] klog can be changed to add V(2) to handle verbository from the command line, see https://pkg.go.dev/k8s.io/klog/v2
 - [ ] At some point we want more intelligent use of labels/selectors (I haven't really read about them yet)
 - [ ] We might eventually want a variable to control quorum expectation (e.g., rank 0 waiting or not)
 - [ ] Look into slurm feature (salloc option) to just start (locate resource and keep it going?) (do we still need this?)

#### Completed

 - [x] diagnostics should be on the level of the container
 - [x] the spec needs to support a local volume
 - [x] Eventually; nice pretty, branded user docs that describe creating CRD, and cases of sleep infinity vs command
 - [x] test better method from Aldo for networking
 - [x] allow to specify that the app restful server is installed, and don't install again
 - [x] we need to test that N=1 case works as expected (not waiting for any workers) and 0 spits an error (for now it doesn't make sense)
 - [x] docs need spell checking!
 - [x] Convert markdown docs into pretty, organized, rendered web-docs
 - [x] Remove automated builds from here in favor of https://github.com/rse-ops/flux-hpc
 - [x] Maximum time for job (seconds) set by CRD
 - [x] Do we want to be able to launch additional tasks? (e.g., after the original job started) (for now, no, but this can be re-addressed if a case comes up)
 - [x] Should there be a min/max size for the MiniCluster CRD (probably 2 if we want to have main/worker)? (right not just cannot be zero)
 - [x] We will want an event/watcher to shut down all jobs when the main command is done. Otherwise the others sometimes keep running (currently we require all ranks to be ready and then they clean up)
 - [x] what should be the proper start command for the main/worker nodes (this is important because it will determine when a job is complete) (rank 0 runs user command, workers just start)
 - [x] figure out where to put flux hostname / config - volume needs write
 - [x] Are `--cores` properly set (yes, not setting uses the default set by hwloc and that's resonable)
 - [x] debug nodes finding on another (see How it works in README.md)
 - [x] By what user? I am currently root but flux is an option (decided to use root to setup, but then run as flux user)
 - [x] How we can print better verbose debugging output (possibly exposed by a variable) (done, debugging boolean is added)
 - [x] And have some solid evidence the node communication is successful (or is the job running that evidence? (this now appears to be printing)
 - [x] debug pod containers not seeing config again (e.g., mounts not creating)
 - [x] Should the secondary (non-driver) pods have a different start command? (answer is no - with the Indexed job it's all the same command)
 - [x] Details for etc-hosts (or will this just work? - no it won't just work)

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
