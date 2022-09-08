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
	jobctrl "flux-framework/flux-operator/pkg/job"
	"flux-framework/flux-operator/pkg/util/uuid"
)

// This interface allows us to define a NotifyMiniClusterUpdate function
type MiniClusterUpdateWatcher interface {
	NotifyMiniClusterUpdate(*api.MiniCluster)
}

// MiniClusterReconciler reconciles a MiniCluster object
type MiniClusterReconciler struct {
	Client     client.Client
	Scheme     *runtime.Scheme
	Manager    ctrl.Manager
	log        logr.Logger
	watchers   []MiniClusterUpdateWatcher
	RESTClient rest.Interface
	RESTConfig *rest.Config
}

func NewMiniClusterReconciler(client client.Client, scheme *runtime.Scheme, restConfig rest.Config, restClient rest.Interface, watchers ...MiniClusterUpdateWatcher) *MiniClusterReconciler {
	return &MiniClusterReconciler{
		log:        ctrl.Log.WithName("minicluster-reconciler"),
		Client:     client,
		Scheme:     scheme,
		watchers:   watchers,
		RESTClient: restClient,
		RESTConfig: &restConfig,
	}
}

//+kubebuilder:rbac:groups=flux-framework.org,resources=miniclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=miniclusters/status,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=miniclusters/finalizers,verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups=flux-framework.org,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
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
func (r *MiniClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// Create a new MiniCluster
	var cluster api.MiniCluster

	// Keep developer informed what is going on.
	r.log.Info("🦕 Event received by MiniCluster!")
	r.log.Info("Request: ", "req", req)

	// Does the Flux Job exist yet (based on name and namespace)
	err := r.Client.Get(ctx, req.NamespacedName, &cluster)
	if err != nil {

		// Create it, doesn't exist yet
		if errors.IsNotFound(err) {
			r.log.Info("🌀 MiniCluster not found . Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		r.log.Info("🌀 Failed to get MiniCluster. Re-running reconcile.")
		return ctrl.Result{Requeue: true}, err
	}

	// Get the current job status
	status := jobctrl.GetCondition(&cluster)

	// TODO how can we use Status (Conditions) here?
	r.log.Info("🌀 Reconciling Mini Cluster", "Image: ", cluster.Spec.Image, "Command: ", cluster.Spec.Command, "Name:", cluster.Status.JobId, "Status:", status)

	// Ensure we have the minicluster (get or create!)
	result, err := r.ensureMiniCluster(ctx, &cluster)
	if err != nil {
		return result, err
	}

	// By the time we get here we have a Job + pods + config maps!
	// What else do we want to do?
	r.log.Info("🌀 Mini Cluster is Ready!")
	return ctrl.Result{}, nil
}

// Notify watchers (the FluxSetup) that we have a new job request
func (r *MiniClusterReconciler) notifyWatchers(job *api.MiniCluster) {
	for _, watcher := range r.watchers {
		watcher.NotifyMiniClusterUpdate(job)
	}
}

// joStatus (to start) can either be finished or pending
func jobStatus(job *api.MiniCluster) string {
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
func (r *MiniClusterReconciler) Create(e event.CreateEvent) bool {

	// Only respond to job events!
	job, match := e.Object.(*api.MiniCluster)
	if !match {
		return true
	}

	// Add conditions - they should never exist for a new job
	job.Status.Conditions = jobctrl.GetJobConditions()

	// We will tell FluxSetup there is a new job request
	defer r.notifyWatchers(job)
	r.log.Info("🌀 MiniCluster create event", "Name:", job.Name)

	// Continue to creation event
	r.log.Info("🌀 MiniCluster was added!", "Name:", job.Name, "Condition:", jobctrl.GetCondition(job))
	return true
}

func (r *MiniClusterReconciler) Delete(e event.DeleteEvent) bool {

	job, match := e.Object.(*api.MiniCluster)
	if !match {
		return true
	}

	defer r.notifyWatchers(job)
	log := r.log.WithValues("job", klog.KObj(job))
	log.Info("🌀 MiniCluster delete event")

	// TODO should trigger a delete here
	// Reconcile should clean up resources now
	return true
}

func (r *MiniClusterReconciler) Update(e event.UpdateEvent) bool {
	oldMC, match := e.ObjectOld.(*api.MiniCluster)
	if !match {
		return true
	}

	// Figure out the state of the old job
	mc := e.ObjectNew.(*api.MiniCluster)

	r.log.Info("🌀 MiniCluster update event")

	// If the job hasn't changed, continue reconcile
	// There aren't any explicit updates beyond conditions
	if jobctrl.JobsEqual(mc, oldMC) {
		return true
	}

	// TODO: check if ready or running, shouldn't be able to update
	// OR if we want update, we need to completely delete and recreate
	return true
}

func (r *MiniClusterReconciler) Generic(e event.GenericEvent) bool {
	r.log.V(3).Info("Ignore generic event", "obj", klog.KObj(e.Object), "kind", e.Object.GetObjectKind().GroupVersionKind())
	return false
}

// newJob inits a new job, creating both the id and original conditions
func (r *MiniClusterReconciler) newJob(ctx context.Context, cluster *api.MiniCluster) (ctrl.Result, error) {

	// We should never edit the object directly?
	clusterCopy := cluster.DeepCopy()

	// If we haven't generated a JobId yet, do that now
	// This might be eventually useful for labels / selector of some kind
	if cluster.Status.JobId == "" {
		clusterCopy.Status.JobId = uuid.Generate(cluster.Name)
		r.Client.Status().Update(ctx, clusterCopy)
		return ctrl.Result{Requeue: true}, nil
	}

	// This should be done in create? Just in case...
	// Get available conditions and set on copy
	conditions := jobctrl.GetJobConditions()
	clusterCopy.Status.Conditions = conditions

	// Update the status of the resource on the CRD
	return ctrl.Result{Requeue: true}, r.Client.Status().Update(ctx, clusterCopy)
}

// SetupWithManager sets up the controller with the Manager.
func (r *MiniClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.MiniCluster{}).

		// This references the Create/Delete/Update,etc functions above
		// they return a boolean to indicate if we should reconcile given the event
		WithEventFilter(r).
		Owns(&batchv1.Job{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
