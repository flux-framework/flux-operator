# Design 2.2

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


![design-2-v.png](design-2-v.png)

The drawing session [can be found here](https://excalidraw.com/#json=QU4SQU-NMBWZS6dFiqa_1,3pI-im0G_WGhF7UdgJdUOg)