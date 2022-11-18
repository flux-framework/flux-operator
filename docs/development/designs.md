# Designs

## Active Design

![09-16-2022/design-three-team.png](09-16-2022/design-three-team.png)

 - [See the Design](09-16-2022/index.md)

 - A **MiniCluster** is an [indexed job](https://kubernetes.io/docs/tasks/job/indexed-parallel-processing-static/) so we can create N copies of the "same" base containers (each with Flux, and the connected workers in our cluster)
 - The flux config is written to a volume at `/etc/flux/config` (created via a config map) as a brokers.toml file.
 - The startup script "wait.sh" handles writing the curve certificate (`/mnt/curve/curve.cert`) and generating and sharing a munge key in the same location. 
 - Networking of the pods works by way of exposing a service that includes the Pod subdomain. We add fully qualified domain names to the pods so that the `hostname` command matches the full name, and Flux is given the full names in its broker.toml.
 - The main pod either runs `flux start` with a web service (creating a persistent "Mini Cluster" or `flux start` with a specific command (if provided in the CRD) in which case the command runs, and the jobs finish and the cluster goes away.


## Early Designs

It is often useful to think about many different designs, and iterate quickly before deciding on a direction
to take, and this is what we did in the first days of the Flux Operator.  The final design (at the top)
is what we decided to focus on, and the others are described (and illustrated) in detail at their respective
links. This is a review of various early designs we've thought about, presented in reverse order (newest to latest).

### Design 3

This design is a simple design based around a single custom resource definition

 - [See the Design](09-07-2022/index.md)


### Design 2.2

This was a set of weekend work I embarked on because (primarily) I couldn't get an internal CRD (the MiniCluster) working
and needed to keep this state elsewhere. 

 - [See the Design](09-05-2022/index.md)


### Design 2.1

This was the earliest design I came up with after talking to Eduardo!

 - [See the Design](09-01-2022/index.md)

### Design 1

This is the original design spec'd out by Dan!

 - [See the Design](08-31-2022/index.md)


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

Know that the scaffolded k8sClient includes a cache that automatically updates based on watches, and may give you out-of-date data (but this is fine because if it is out-of-date, there should be a reconcile in the queue already). Also know that there is a way to request objects bypassing a cache (look for APIReader).  This gives a read-only, but direct access to the API.  Useful for e.g. those workspace objects.
