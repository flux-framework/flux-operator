# TODO

### Design 2

 - [ ] klog can be changed to add V(2) to handle verbository from the command line, see https://pkg.go.dev/k8s.io/klog/v2
 - [ ] pkg/util/heap should implement an actual heap
 - [ ] kubebuilder should be able to provide defaults in the *_types.
 - [ ] FluxSetup create should do some handling to ensure we only have one setup. If there is an existing setup it should update with the new one and then the old be deleted (without messing up running jobs).
 - [ ] Currently we have no representation of quota - we need to be able to set (and check) hard limits defined in the setup. Kueue has a scheduler entity. 
 - [x] Figure out logging connected to reconciler

### Design 1 (not currently working on)

- [x] consolidate configmap functions into shared functionality (less redundancy)
- [x] Debug why the configmaps aren't being populated with the hostfile (it wasn't working with kind, worked without changes with minikube)
- [x] Figure out adding namespaces to config/samples - should be flux-operator
- [x] Each of config files written (e.g., hostname, broker, cert) should have their own types and more simply generated. The strategy right now is just temporary.
- [ ] Cert needs to be separated / generated
- [x] Stateful set (figure out how to create properly, doesn't seem to have pods) (figured out need to create ConfigMaps for Volumes)
- [ ] Stateful set is created but it cannot see configmaps "MountVolume.SetUp failed for volume "etc-hosts" : configmap references non-existent config key: etc-hosts"
- [ ] A means to generate / update certs - I don't think manually doing it is the right approach, but there is a comment that cert-manager isn't supported?
  - https://kubernetes.io/docs/tasks/configmap-secret/managing-secret-using-kustomize/
  - This is useful https://github.com/jetstack/kustomize-cert-manager-demo
  - And https://www.jetstack.io/blog/kustomize-cert-manager/
  - https://github.com/kubernetes-sigs/kustomize/blob/master/examples/secretGeneratorPlugin.md
  
## Design 1

We create a new custom controller that listens for FluxJob resources. When a new FluxJob is created, the controller goes through the following steps:

 1. Create ConfigMap that contains:
   - A broker.toml file that lists the pods in the worker StatefulSet plus the Flux rank 0 (in the form of flux-workers-0, flux-workers-1, ...) which defines their rank in the Flux TBON
   - /etc/hosts with a mapping between pod names and IPs (should be able to generate this before the StatefulSet  or IndexedJob per this MPI Operator line which occurs before the StatefulSet creation
 2. Create the worker StatefulSet that contains the desired replicas minus 1, as rank 0 also does work. If we use IndexedJob  we can create all desired replicas
 3. Wait for the worker pods to enter Running state.
 4. Create the launcher Job.
 5. After the launcher job finishes, set the replicas to 0 in the worker StatefulSet, or delete the IndexedJob 

## Wisdom

**from the kubebuilder slack**

### Learned Knowledge

- Reconciling should only take into account the spec of your object, and the real world.  Don't use status to hold knowledge for future reconcile loops.  Use a workspace object instead.
- Status should only hold observations of the reconcile loop.  Conditions, perhaps a "Phase", IDs of stuff you've found, etc.
- Use k8s ownership model to help with cleaning up things that should automatically be reclaimed when your object is deleted.
- Use finalizers to do manual clean-up-tasks
- Send events, but be very limited in how often you send events.  We've opted now to send events, essentially only when a Condition is modified (e.g. a Condition changes state or reason).
- Try not to do too many things in a single reconcile.  One thing is fine.  e.g. see one thing out of order?  Fix that and ask to be reconciled.  The next time you'll see that it's in order and you can check the next thing.  The resulting code is very robust and can handle almost any failure you throw at it.
- Add "kubebuilder:printcolums" markers to help kubectl-users get a nice summary when they do "kubectl get yourthing".
- Accept and embrace that you will be reconciling an out-of-date object from time to time.  It shouldn't really matter.  If it does, you might want to change things around so that it doesn't matter.  Inconsistency is a fact of k8s life.
- Place extra care in taking errors and elevating them to useful conditions, and/or events.  These are the most visible part of an operator, and the go-to-place for humans when trying to figure out why your code doesn't work.  If you've taken the time to extract the error text from the underlying system into an Event, your users will be able to fix the problem much quicker.

### What is a workspace?

A workspace object is when you need to record some piece of knowledge about a thing you're doing, so that later you can use that when reconciling this object. MyObject "foo" is reconciled; so to record the thing you need to remember, create a MyObjectWorkspace — Owned by the MyObject, and with the same name + namespace.  MyObjectWorkspace doesn't need a reconciler; it's simply a tool for you to remember the thing. Next time you reconcile a MyObject, also read your MyObjectWorkspace so you can remember "what happened last time". E.g. I've made a controller to create an EC2 instance, and we needed to be completely sure that we didn't make the "launch instance" API call twice.  EC2 has a "post once only" technique whereby you specify a nonce to avoid duplicate API calls.  You would write the nonce to the workspace use the nonce to call the EC2 API write any status info of what you observed to the status. Rremove the nonce when you know that you've stored the results (e.g. instance IDs or whatever) When you reconcile, if the nonce is set, you can re-use it because it means that the EC2 call failed somehow.  EC2 uses the nonce the second time to recognise that "heh, this is the same request as before ..." Stuff like this nonce shouldn't go in your status. Put simply, the status should really never be used as input for your reconcile.

Know that the scaffolded k8sClient includes a cache that automatically updates based on watches, and may give you out-of-date data (but this is fine because if it is out-of-date, there should be a reconcile in the queue already). Also know that there is a way to request objets bypassing a cache (look for APIReader).  This gives a read-only, but direct access to the API.  Useful for e.g. those workspace objects.