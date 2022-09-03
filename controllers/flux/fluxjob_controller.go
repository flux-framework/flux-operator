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

	"github.com/go-logr/logr"
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

//+kubebuilder:rbac:groups=flux-framework.org,resources=miniclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=miniclusters/status,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=miniclusters/finalizers,verbs=get;list;watch;create;update;patch;delete

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

	// Keep developed informed what is going on.
	log.Info("⚡️ Event received! ⚡️")
	log.Info("Request: ", "req", req)

	// Does the Flux Job exist yet (based on name and namespace)
	err := r.Client.Get(ctx, req.NamespacedName, &fluxjob)
	if err != nil {

		// Create it, doesn't exist yet
		if errors.IsNotFound(err) {
			log.Info("Flux Job not found . Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Info("Failed to get Flux Job. Re-running reconcile.")
		return ctrl.Result{}, err
	}
	// If we don't have them, set minicluster conditions on the fluxjob
	if len(fluxjob.Status.Conditions) == 0 || fluxjob.Status.JobId == "" {
		return r.newJob(ctx, &fluxjob)
	}

	log.Info("🏃 Found Flux Job 🏃 ", "Image: ", fluxjob.Spec.Image, "Command: ", fluxjob.Spec.Command, "Name", fluxjob.Status.JobId)

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
	log := r.log.WithValues("job", klog.KObj(job))
	log.Info("🌀 Job create event")

	// If it's waiting or running, do nothing
	// TODO might there be some need to update something if waiting?
	// I assume after it's running you can't, but maybe yes for waiting
	if status == jobctrl.Running || status == jobctrl.Waiting {
		log.Info("🌀 Job is running or waiting", "Name:", job.Name)
		return false
	}

	// If it's finished we need to clean up
	if status == jobctrl.Finished {
		log.Info("🌀 Job is finished", "Name:", job.Name)
		return true
	}

	// If we get here it was requested. We don't need to reconcile, but we need to ensure the flux manager
	// knows about the job (and we update it on our manager queue)
	jobCopy := job.DeepCopy()

	// TODO handle any figuring out of resources?
	// https://github.com/kubernetes-sigs/kueue/blob/main/pkg/controller/core/workload_controller.go#L280

	// As an alternative we could store a variable on the spec to indicate admitted,
	// but I think using conditions are more best practice
	if !r.fluxManager.AddOrUpdateJob(jobCopy) {
		log.Info("🌀 Issue adding or updating job; ignored for now")
		return false
	}
	// If we get here update the job condition to be waiting
	log.Info("🌀 Job was added or updated!", "Name:", job.Name)
	jobctrl.FlagConditionWaiting(job)
	return true
}

func (r *FluxJobReconciler) Delete(e event.DeleteEvent) bool {

	job := e.Object.(*api.FluxJob)
	//	defer r.notifyWatchers(wl)
	//	status := workloadStatus(wl)
	log := r.log.WithValues("job", klog.KObj(job))
	log.Info("Job delete event")

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
	log.Info("Workload update event")

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
		WithEventFilter(r).
		Complete(r)
}
