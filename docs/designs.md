# Designs

This is a review of various designs we've thought about, presented in reverse order (newest to latest)

## Design 3

 - [See the Design](09-07-2022)

This design is a simple design based around a single custom resource definition:

 - A **MiniCluster**: is a specification provided by a user with a container image, command, and size to create (over several reconciles) ConfigMaps, Secrets, and a batchv1.Job with pods (containers) and volumes. Once this is created, we conceptually have a MiniCluster within K8s running the user's job! This is a CRD, so the user submitting it owns that MiniCluster.

We discussed this at our meeting on 9-07-22 and wanted to remove the complexity of design 2, and have a CRD that directly manages resources. In more detail, this means that:

- A MiniCluster is a CRD that includes an image, command, and size
- Creating a MiniCluster first creates Config Maps and secrets (and possibly etc hosts) and then a batch job with pods that use them
- Flux is required in the running container image (from the user)
- The job runs to completion and has some final state


## Design 2.2

 - [See the Design](09-05-2022)

This was a set of weekend work I embarked on because (primarily) I couldn't get an internal CRD (the MiniCluster) working
and needed to keep this state elsewhere. This resulted in the following entities:

 - A **FluxManager** was created on start and passed to both controllers. It would hold a queue of waiting/running jobs.
 - A **MiniCluster**: is a concept that is mapped to a MiniCluster. A MiniCluster, when it's given permission to run, creates a MiniCluster, which (currently) is a StatefulSet, ConfigMaps and Secrets. The idea is that the Job submitter owns that MiniCluster to submit jobs to.
 - A **scheduler**: currently is created alongside the MiniCluster and FluxSetup, and can access the same manager and queue. The scheduler does nothing aside from moving jobs from waiting to running, and currently is a simple loop. We eventually want this to be more.
 - A **FluxSetup** and **FluxJob** would be the two CRDs that both could see the manager.

The scheduler above should not be confused with anything intelligent to schedule jobs - it's imagined as a simple check to see if we have resources on our cluster for a new MiniClusters, and grant them permission for create if we do. For Design 2 I was trying to mirror a design that is used by kueue, where the two controllers are aware of the state of the cluster via a third "manager" that can hold a queue of jobs. The assumption here is that the user is submitting jobs (MiniCluster that will map to a "MiniCluster") and the cluster (FluxSetup) has some max amount of room available, and we need those two things to be aware of one another. Currently the flow I'm working on is the following:

1. We create the FluxSetup
 - The Create event creates a queue (and asks for reconcile). 
 -  The reconcile should always get the current number of running batch jobs and save to status and check against quota (not done yet). We also need Create/Update to manage cleaning up old setups to replace with new (not done yet). Right now we assume unlimited quota for the namespace, which of course isn't true!
2. A MiniCluster is submit (e.g., by a user) requesting a MiniCluster
3. The FluxSetup is watching for the Create event, and when it sees a new job, it adds it to the manager queue (under waiting jobs). The job is flagged as waiting. We eventually want the setup to do more checks here for available resources, but currently we allow everything to be put into waiting.
4. Moving from waiting -> the manager->queue->heap indicates running, and this is done by the scheduler.
5. During this time, the MiniCluster reconcile loop is continually running, and controls the lifecycle of the MiniCluster, per request of the user and status of the cluster.
 - If the job status is requested, we add it to the waiting queue and flag as waiting. This will happen right after the job is created, and then run reconcile again.
 - If a job is waiting, this could mean two things.
   - It's in the queue heap (allowed to run), in which case we create the MiniCluster and update status to Ready
   - It's still waiting, we ask to reconcile until we see it's allowed to run
 - If a job is ready, it just had its MiniCluster created! At this point we need to "launch" our job (not done yet).
 - If it's running, we need to check to see if it's finished, and flag as finished if yes (not done yet). If it's not finished, we keep the status as running and re-reconcile.
 - If the job is finished, we need to clean up
 - When we reach the bottom of the loop, we don't need to reconcile again.

Some current additional notes:

1. MiniCluster submit before the FluxSetup is ready are ignored. We don't have a Queue init yet to support them, and the assumption is that the user can re-submit.
2. This design wasn't very good, but I learned a ton (and was on cloud 9 working on it!)


## Design 2.1

 - [See the Design](09-01-2022)

At this point I chat with Eduardo about operator design, and we decided to go for a simple design:

- Flux: would have a CR that defines a job and exposes an entrypoint command / container to the user
- FluxSetup would have most of the content of [here](https://lc.llnl.gov/confluence/display/HFMCCEL/Flux+Operator+Design) and be more of an internal or admin setup.

To generate FluxSetup I think I could (maybe?) have run this command again, but instead I added a new entry to the [PROJECT](PROJECT) and then generated [api/v1alpha1/fluxsetup_types.go](api/v1alpha1/fluxsetup_types.go) from the `flux_types.go` (and changing all the references from Flux to FluxSetup). 
I also needed to (manually) make [controllers/fluxsetup_controller.go](controllers/fluxsetup_controller.go) and ensure it was updated to use FluxSetup, and then adding
it's creation to [main.go](main.go). Yes, this is a bit of manual work, but I think I'll only need to do it once.
At this point, I needed to try and represent what I saw in the various config files in this types file.


## Design 1

**This is the original design spec'd out by Dan!**

We create a new custom controller that listens for MiniCluster resources. When a new MiniCluster is created, the controller goes through the following steps:

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