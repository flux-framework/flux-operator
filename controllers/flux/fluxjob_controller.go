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
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"k8s.io/apimachinery/pkg/api/errors"

	api "flux-framework/flux-operator/api/v1alpha1"
	"flux-framework/flux-operator/pkg/flux"
	jobctrl "flux-framework/flux-operator/pkg/job"
	"flux-framework/flux-operator/pkg/util/uuid"
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
	RESTClient  rest.Interface
	RESTConfig  *rest.Config
}

var (
	ownerKey = ".metadata.controller"
)

func NewFluxJobReconciler(client client.Client, scheme *runtime.Scheme, q *flux.Manager, restConfig rest.Config, restClient rest.Interface, watchers ...JobUpdateWatcher) *FluxJobReconciler {
	return &FluxJobReconciler{
		log:         ctrl.Log.WithName("fluxjob-reconciler"),
		Client:      client,
		Scheme:      scheme,
		fluxManager: q,
		watchers:    watchers,
		RESTClient:  restClient,
		RESTConfig:  &restConfig,
	}
}

//+kubebuilder:rbac:groups=flux-framework.org,resources=fluxjobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=fluxjobs/status,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=fluxjobs/finalizers,verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups=flux-framework.org,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources="",verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups="",resources=events,verbs=create;watch;update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// We compare the state of the Flux object to the actual cluster state and
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *FluxJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// Create a new FluxJob
	var fluxjob api.FluxJob

	// Keep developer informed what is going on.
	r.log.Info("🕵 Event received by FluxJob!")
	r.log.Info("Request: ", "req", req)

	// Does the Flux Job exist yet (based on name and namespace)
	err := r.Client.Get(ctx, req.NamespacedName, &fluxjob)
	if err != nil {

		// Create it, doesn't exist yet
		if errors.IsNotFound(err) {
			r.log.Info("🌀 Flux Job not found . Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		r.log.Info("🌀 Failed to get Flux Job. Re-running reconcile.")
		return ctrl.Result{}, err
	}
	// If we don't have them, set minicluster conditions on the fluxjob
	// I don't think this should trigger... just in case!
	if len(fluxjob.Status.Conditions) == 0 || fluxjob.Status.JobId == "" {
		return r.newJob(ctx, &fluxjob)
	}

	// Get the current job status
	status := jobctrl.GetCondition(&fluxjob)
	r.log.Info("🌀 Reconciling Flux Job", "Image: ", fluxjob.Spec.Image, "Command: ", fluxjob.Spec.Command, "Name:", fluxjob.Status.JobId, "Conditions:", status)

	// If it's just requested, ensure the waiting queue knows about it!
	if status == jobctrl.ConditionJobRequested {
		if !r.fluxManager.AddOrUpdateJob(&fluxjob) {
			r.log.Info("🌀 Issue adding job, will retry in one minute.")
			return ctrl.Result{RequeueAfter: time.Minute}, nil
		}

		// If we get here update the job condition to be waiting
		jobCopy := fluxjob.DeepCopy()
		jobctrl.FlagConditionWaiting(jobCopy)
		r.Client.Status().Update(ctx, jobCopy)
		return ctrl.Result{Requeue: true}, nil
	}

	// If it's waiting, either it's been admitted (in the fluxmanager heap)
	// or needs to continue waiting
	if status == jobctrl.ConditionJobWaiting {

		// If the FluxJob condition is waiting but the manager says running,
		// this means it was moved to the running heap
		// We need to create the MiniCluster and update the status to
		// be running
		if r.fluxManager.IsRunningJob(&fluxjob) {

			// This will either reconcile with an updated state, or cause the
			// job to re-enter this loop to create other resources for the
			// MiniCluster. When the state changes from waiting to
			// ready, then we know the mini cluster is done
			// and go beyond this loop.
			return r.newMiniCluster(ctx, &fluxjob)
		}

		// If it's waiting but not actually running, we need to check again later
		return ctrl.Result{Requeue: true}, nil
	}

	// By the time we get here we've done Waiting -> Ready
	// At this point, if the status is ready we need to submit the job
	// and call it running
	if status == jobctrl.ConditionJobReady {
		r.log.Info("🌀 Mini Cluster is Ready!")

		// Launching the job handles updating the status
		// If the job fails, we should not retry until the command or job is tweaked
		// We set the command to empty so it will stay ready (but not launch)
		return r.LaunchJob(ctx, &fluxjob)
	}

	if status == jobctrl.ConditionJobRunning {
		// TODO will need to look for running, and then finished state (and change here)
		// TODO decide if when finished if we want to return it to some kind of
		// ready state, or clean up entirely (likely desired). There should be
		// some ability to keep the cluster running, however, to accept new
		// commands if desired. Maybe this should be a FluxSetup variable?
		// finished and then possibly clean up.
		r.log.Info("🌀 Mini Cluster is Running!")
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	// If the job is finished, we need to ensure we cleanup
	if status == jobctrl.ConditionJobFinished {
		// TODO: try cleanup, if true, return reconciled below. Otherwise keep going.
	}
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

	// Only respond to job events!
	job, match := e.Object.(*api.FluxJob)
	if !match {
		return true
	}

	// Add conditions - they should never exist for a new job
	job.Status.Conditions = jobctrl.GetJobConditions()

	// We will tell FluxSetup there is a new job request
	defer r.notifyWatchers(job)
	r.log.Info("🌀 Job create event", "Name:", job.Name)

	// Continue to creation event
	r.log.Info("🌀 Job was added!", "Name:", job.Name, "Condition:", jobctrl.GetCondition(job))
	return true
}

func (r *FluxJobReconciler) Delete(e event.DeleteEvent) bool {

	job, match := e.Object.(*api.FluxJob)
	if !match {
		return true
	}

	// TODO any reason to notify watchers here?
	//	defer r.notifyWatchers(wl)
	log := r.log.WithValues("job", klog.KObj(job))
	log.Info("🌀 Job delete event")

	// If it's in the waiting or running queues, we need to delete
	if r.fluxManager.Delete(job) {
		log.Info("🌀 Job was deleted", "Name:", job.Name)
	}

	// Update status to be finished so resources are cleaned up
	jobCopy := job.DeepCopy()
	jobctrl.FlagConditionFinished(jobCopy)
	r.Client.Status().Update(context.TODO(), jobCopy)

	// Reconcile should clean up resources now
	return true
}

func (r *FluxJobReconciler) Update(e event.UpdateEvent) bool {
	oldJob, match := e.ObjectOld.(*api.FluxJob)
	if !match {
		return true
	}

	// Figure out the state of the old job
	job := e.ObjectNew.(*api.FluxJob)

	r.log.Info("🌀 Job update event")

	// If the job hasn't changed, continue reconcile
	// There aren't any explicit updates beyond conditions
	if jobctrl.JobsEqual(job, oldJob) {
		return true
	}

	// No matter what, if it's running we can't modify it
	// For now user should delete and re-submit
	if r.fluxManager.IsRunningJob(oldJob) {
		r.log.Info("🌀 Job is running and cannot be updated", "Name:", oldJob.Name)
		return false
	}

	// If it's waiting / finished / ready we can update and change status
	if r.fluxManager.IsWaitingJob(oldJob) {
		r.fluxManager.Delete(oldJob)
	}

	jobCopy := job.DeepCopy()

	if !r.fluxManager.AddOrUpdateJob(jobCopy) {
		r.log.Info("🌀 Issue updating job; ignored for now")
		return false
	}

	jobctrl.FlagConditionWaiting(job)
	r.log.Info("🌀 Job was updated!", "Name:", job.Name, "Condition:", jobctrl.GetCondition(job))
	return true
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
		WithEventFilter(r).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Pod{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
