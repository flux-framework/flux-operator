/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"k8s.io/apimachinery/pkg/api/errors"

	api "flux-framework/flux-operator/api/v1alpha1"
	"flux-framework/flux-operator/pkg/flux"
	jobctrl "flux-framework/flux-operator/pkg/job"
	"flux-framework/flux-operator/pkg/util/uuid"

	logctrl "sigs.k8s.io/controller-runtime/pkg/log"
)

// This interface allows us to define a NotifyJobUpdate functionn
type JobUpdateWatcher interface {
	NotifyJobUpdate(*api.FluxJob)
}

// FluxJobReconciler reconciles a FluxJob object
type FluxJobReconciler struct {
	Client      client.Client
	Scheme      *runtime.Scheme
	Manager     ctrl.Manager
	log         logr.Logger
	fluxManager *flux.Manager
	watchers    []JobUpdateWatcher
}

var (
	ownerKey = ".metadata.controller"
)

func NewFluxJobReconciler(client client.Client, scheme *runtime.Scheme, q *flux.Manager, watchers ...JobUpdateWatcher) *FluxJobReconciler {
	return &FluxJobReconciler{
		log:         ctrl.Log.WithName("fluxjob-reconciler"),
		Client:      client,
		Scheme:      scheme,
		fluxManager: q,
		watchers:    watchers,
	}
}

//+kubebuilder:rbac:groups=flux-framework.org,resources=fluxjobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=fluxjobs/status,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=fluxjobs/finalizers,verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups=flux-framework.org,resources=jobs,verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups="",resources=events,verbs=create;watch;update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// We compare the state of the Flux object to the actual cluster state and
// perform operations to make the cluster state reflect the state specified by
// the user.
// LOGIC:
// 1. Get the original job CRD
// 2. Generate a unique identifier for it - the name plus mangled string
// 3. Update the status to be pending resources
// 4.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *FluxJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// Create a new FluxJob
	var fluxjob api.FluxJob

	log := logctrl.FromContext(ctx).WithValues("FluxJob", req.NamespacedName)

	// Keep developer informed what is going on.
	log.Info("🕵 Event received by FluxJob!")
	log.Info("Request: ", "req", req)

	// Does the Flux Job exist yet (based on name and namespace)
	err := r.Client.Get(ctx, req.NamespacedName, &fluxjob)
	if err != nil {

		// Create it, doesn't exist yet
		if errors.IsNotFound(err) {
			log.Info("🌀 Flux Job not found . Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Info("🌀 Failed to get Flux Job. Re-running reconcile.")
		return ctrl.Result{}, err
	}
	// If we don't have them, set minicluster conditions on the fluxjob
	// I don't think this should trigger...
	if len(fluxjob.Status.Conditions) == 0 || fluxjob.Status.JobId == "" {
		return r.newJob(ctx, &fluxjob)
	}

	// Get the current job status
	status := jobctrl.GetCondition(&fluxjob)

	// If it's running, let it keep running (for now)
	if status == jobctrl.ConditionJobRunning {
		// TODO check pods for being finished
		// TODO can we put a requeue checking time? Maybe a minute?
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	r.log.Info("🌀 Reconciling Flux Job", "Image: ", fluxjob.Spec.Image, "Command: ", fluxjob.Spec.Command, "Name:", fluxjob.Status.JobId, "Conditions:", jobctrl.GetCondition(&fluxjob))

	// If it's waiting, either it's been admitted (in the fluxmanager) or needs to continue waiting
	if status == jobctrl.ConditionJobWaiting {

		// Determine if job is running, update status to running
		if r.fluxManager.IsRunningJob(&fluxjob) {

		}
	}

	// TODO do we need to ensure that there is only one instance of the batchjobs owned by FluxJob, or FluxJob
	// for this reconciler request?

	/*jobFinishedCond, jobFinished := jobFinishedCondition(&fluxjob)
	// 2. create new workload if none exists
	if wl == nil {
		// Nothing to do if the job is finished
		if jobFinished {
			return ctrl.Result{}, nil
		}
		err := r.handleJobWithNoWorkload(ctx, &job)
		if err != nil {
			log.Error(err, "Handling job with no workload")
		}
		return ctrl.Result{}, err
	}

	// 3. handle a finished job
	if jobFinished {
		added := false
		wl.Status.Conditions, added = appendFinishedConditionIfNotExists(wl.Status.Conditions, jobFinishedCond)
		if !added {
			return ctrl.Result{}, nil
		}
		err := r.client.Status().Update(ctx, wl)
		if err != nil {
			log.Error(err, "Updating workload status")
		}
		return ctrl.Result{}, err
	}

	// 4. Handle a not finished job
	if jobSuspended(&job) {
		// 4.1 start the job if the workload has been admitted, and the job is still suspended
		if wl.Spec.Admission != nil {
			log.V(2).Info("Job admitted, unsuspending")
			err := r.startJob(ctx, wl, &job)
			if err != nil {
				log.Error(err, "Unsuspending job")
			}
			return ctrl.Result{}, err
		}

		// 4.2 update queue name if changed.
		q := queueName(&job)
		if wl.Spec.QueueName != q {
			log.V(2).Info("Job changed queues, updating workload")
			wl.Spec.QueueName = q
			err := r.client.Update(ctx, wl)
			if err != nil {
				log.Error(err, "Updating workload queue")
			}
			return ctrl.Result{}, err
		}
		log.V(3).Info("Job is suspended and workload not yet admitted by a clusterQueue, nothing to do")
		return ctrl.Result{}, nil
	}

	if wl.Spec.Admission == nil {
		// 4.3 the job must be suspended if the workload is not yet admitted.
		log.V(2).Info("Running job is not admitted by a cluster queue, suspending")
		err := r.stopJob(ctx, wl, &job, "Not admitted by cluster queue")
		if err != nil {
			log.Error(err, "Suspending job with non admitted workload")
		}
		return ctrl.Result{}, err
	}*/

	// 4.4 workload is admitted and job is running, nothing to do.
	log.Info("🌀 Job running with admitted workload, nothing to do")

	// This will reconcile and trigger the watch on the MiniCluster
	return ctrl.Result{}, nil
}

// Notify watchers (the FluxSetup) that we have a new job request
func (r *FluxJobReconciler) notifyWatchers(job *api.FluxJob) {
	for _, watcher := range r.watchers {
		watcher.NotifyJobUpdate(job)
	}
}

// joStatus (to start) can either be finished or pending
func jobStatus(job *api.FluxJob) string {
	// If the job is finished, return finished status
	if jobctrl.HasCondition(job, jobctrl.ConditionJobFinished) {
		return jobctrl.Finished
	}
	if jobctrl.HasCondition(job, jobctrl.ConditionJobRunning) {
		return jobctrl.Running
	}
	if jobctrl.HasCondition(job, jobctrl.ConditionJobRequested) {
		return jobctrl.Requested
	}
	return jobctrl.Waiting
}

// Called when a new job is created
func (r *FluxJobReconciler) Create(e event.CreateEvent) bool {
	job := e.Object.(*api.FluxJob)

	// Add conditions if don't exist yet
	if len(job.Status.Conditions) == 0 || job.Status.Conditions == nil {
		job.Status.Conditions = jobctrl.GetJobConditions()
	}

	// We will tell FluxSetup there is a new job request
	defer r.notifyWatchers(job)
	status := jobStatus(job)
	r.log.Info("🌀 Job create event", "Name:", job.Name)

	// If it's waiting or running, do nothing
	// TODO might there be some need to update something if waiting?
	// I assume after it's running you can't, but maybe yes for waiting
	if status == jobctrl.Running || status == jobctrl.Waiting {
		r.log.Info("🌀 Job is running or waiting", "Name:", job.Name)
		return false
	}

	// If it's finished we need to clean up
	if status == jobctrl.Finished {
		r.log.Info("🌀 Job is finished", "Name:", job.Name)
		return true
	}

	// If we get here it was requested. We don't need to reconcile, but we need to ensure the flux manager
	// knows about the job (and we update it on our manager queue)
	jobCopy := job.DeepCopy()

	// TODO handle any figuring out of resources?
	// https://github.com/kubernetes-sigs/kueue/blob/main/pkg/controller/core/workload_controller.go#L280

	// Add the job to the waiting queue - when it is moved to the heap
	// it is considered running.
	if !r.fluxManager.AddOrUpdateJob(jobCopy) {
		r.log.Info("🌀 Issue adding or updating job; ignored for now")
		return false
	}
	// If we get here update the job condition to be waiting
	jobctrl.FlagConditionWaiting(job)
	r.log.Info("🌀 Job was added or updated!", "Name:", job.Name, "Condition:", jobctrl.GetCondition(job))
	return true
}

func (r *FluxJobReconciler) Delete(e event.DeleteEvent) bool {

	job := e.Object.(*api.FluxJob)
	//	defer r.notifyWatchers(wl)
	//	status := workloadStatus(wl)
	log := r.log.WithValues("job", klog.KObj(job))
	log.Info("🌀 Job delete event")

	/*	if !e.DeleteStateUnknown {
			status = workloadStatus(wl)
		}
		log := r.log.WithValues("workload", klog.KObj(wl), "queue", wl.Spec.QueueName, "status", status)
		log.V(2).Info("Workload delete event")
		ctx := ctrl.LoggerInto(context.Background(), log)

		// When assigning a clusterQueue to a workload, we assume it in the cache. If
		// the state is unknown, the workload could have been assumed and we need
		// to clear it from the cache.
		if wl.Spec.Admission != nil || e.DeleteStateUnknown {
			if err := r.cache.DeleteWorkload(wl); err != nil {
				if !e.DeleteStateUnknown {
					log.Error(err, "Failed to delete workload from cache")
				}
			}

			// trigger the move of associated inadmissibleWorkloads if required.
			r.queues.QueueAssociatedInadmissibleWorkloads(ctx, wl)
		}

		// Even if the state is unknown, the last cached state tells us whether the
		// workload was in the queues and should be cleared from them.
		if wl.Spec.Admission == nil {
			r.queues.DeleteWorkload(wl)
		}*/
	return false
}

func (r *FluxJobReconciler) Update(e event.UpdateEvent) bool {

	//_ := e.ObjectOld.(*api.FluxJob)
	newJob := e.ObjectNew.(*api.FluxJob)

	//	defer r.notifyWatchers(wl)
	//	status := workloadStatus(wl)
	log := r.log.WithValues("job", klog.KObj(newJob))
	log.Info("🌀 Job update event")

	/*	oldWl := e.ObjectOld.(*kueue.Workload)
		wl := e.ObjectNew.(*kueue.Workload)
		defer r.notifyWatchers(oldWl)
		defer r.notifyWatchers(wl)

		status := workloadStatus(wl)
		log := r.log.WithValues("workload", klog.KObj(wl), "queue", wl.Spec.QueueName, "status", status)
		ctx := ctrl.LoggerInto(context.Background(), log)

		prevQueue := oldWl.Spec.QueueName
		if prevQueue != wl.Spec.QueueName {
			log = log.WithValues("prevQueue", prevQueue)
		}
		prevStatus := workloadStatus(oldWl)
		if prevStatus != status {
			log = log.WithValues("prevStatus", prevStatus)
		}
		if wl.Spec.Admission != nil {
			log = log.WithValues("clusterQueue", wl.Spec.Admission.ClusterQueue)
		}
		if oldWl.Spec.Admission != nil && (wl.Spec.Admission == nil || wl.Spec.Admission.ClusterQueue != oldWl.Spec.Admission.ClusterQueue) {
			log = log.WithValues("prevClusterQueue", oldWl.Spec.Admission.ClusterQueue)
		}
		log.V(2).Info("Workload update event")

		wlCopy := wl.DeepCopy()
		// We do not handle old workload here as it will be deleted or replaced by new one anyway.
		handlePodOverhead(r.log, wlCopy, r.client)

		switch {
		case status == finished:
			if err := r.cache.DeleteWorkload(oldWl); err != nil && prevStatus == admitted {
				log.Error(err, "Failed to delete workload from cache")
			}
			r.queues.DeleteWorkload(oldWl)

			// trigger the move of associated inadmissibleWorkloads if required.
			r.queues.QueueAssociatedInadmissibleWorkloads(ctx, wl)

		case prevStatus == pending && status == pending:
			if !r.queues.UpdateWorkload(oldWl, wlCopy) {
				log.V(2).Info("Queue for updated workload didn't exist; ignoring for now")
			}

		case prevStatus == pending && status == admitted:
			r.queues.DeleteWorkload(oldWl)
			if !r.cache.AddOrUpdateWorkload(wlCopy) {
				log.V(2).Info("ClusterQueue for workload didn't exist; ignored for now")
			}

		case prevStatus == admitted && status == pending:
			if err := r.cache.DeleteWorkload(oldWl); err != nil {
				log.Error(err, "Failed to delete workload from cache")
			}
			// trigger the move of associated inadmissibleWorkloads if required.
			r.queues.QueueAssociatedInadmissibleWorkloads(ctx, wl)

			if !r.queues.AddOrUpdateWorkload(wlCopy) {
				log.V(2).Info("Queue for workload didn't exist; ignored for now")
			}

		default:
			// Workload update in the cache is handled here; however, some fields are immutable
			// and are not supposed to actually change anything.
			if err := r.cache.UpdateWorkload(oldWl, wlCopy); err != nil {
				log.Error(err, "Updating workload in cache")
			}
		}*/

	return false
}

func (r *FluxJobReconciler) Generic(e event.GenericEvent) bool {
	r.log.V(3).Info("Ignore generic event", "obj", klog.KObj(e.Object), "kind", e.Object.GetObjectKind().GroupVersionKind())
	return false
}

// newJob inits a new job, creating both the id and original conditions
func (r *FluxJobReconciler) newJob(ctx context.Context, fluxjob *api.FluxJob) (ctrl.Result, error) {

	// We should never edit the object directly?
	fluxjobCopy := fluxjob.DeepCopy()

	// If we haven't generated a JobId yet, do that now
	if fluxjob.Status.JobId == "" {
		fluxjobCopy.Status.JobId = uuid.Generate(fluxjob.Name)
		r.Client.Status().Update(ctx, fluxjobCopy)
		return ctrl.Result{Requeue: true}, nil
	}

	// This should be done in create? Just in case...
	// Get available conditions and set on copy
	conditions := jobctrl.GetJobConditions()
	fluxjobCopy.Status.Conditions = conditions

	// Update the status of the resource on the CRD
	return ctrl.Result{Requeue: true}, r.Client.Status().Update(ctx, fluxjobCopy)
}

// SetupWithManager sets up the controller with the Manager.
func (r *FluxJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.FluxJob{}).

		// This references the Create/Delete/Update,etc functions above
		// they return a boolean to indicate if we should reconcile given the event
		Owns(&batchv1.Job{}).
		WithEventFilter(r).
		Complete(r)
}
