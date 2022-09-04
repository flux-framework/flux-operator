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
	//	"fmt"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	//	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	//	logctrl "sigs.k8s.io/controller-runtime/pkg/log"

	api "flux-framework/flux-operator/api/v1alpha1"
	"flux-framework/flux-operator/pkg/defaults"
	"flux-framework/flux-operator/pkg/flux"
	jobctrl "flux-framework/flux-operator/pkg/job"
)

// Buffer for job update channel
const updateChBuffer = 10

// FluxSetupReconciler reconciles a FluxSetup object
type FluxSetupReconciler struct {
	// Client is separate here since we implement our own Create/etc functions
	Client           client.Client
	Scheme           *runtime.Scheme
	log              logr.Logger
	fluxManager      *flux.Manager
	jobUpdateChannel chan event.GenericEvent
}

func NewFluxSetupReconciler(client client.Client, scheme *runtime.Scheme, q *flux.Manager) *FluxSetupReconciler {
	return &FluxSetupReconciler{
		log:              ctrl.Log.WithName("setup-reconciler"),
		Client:           client,
		Scheme:           scheme,
		fluxManager:      q,
		jobUpdateChannel: make(chan event.GenericEvent, updateChBuffer),
	}
}

//+kubebuilder:rbac:groups=flux-framework.org,resources=fluxsetups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=fluxsetups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=flux-framework.org,resources=fluxsetups/finalizers,verbs=update

//+kubebuilder:rbac:groups=flux-framework.org,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=jobs,verbs=get;list;watch;create;update;patch;delete

// Reconcile moves the current state of the cluster closer to the desired state.
func (r *FluxSetupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// Create a new FluxSetup and FluxJob instance
	var setup api.FluxSetup

	// Keep developed informed what is going on.
	r.log.Info("🌞 Reconciling FluxSetup")
	r.log.Info("🕵 Event received by FluxSetup!")
	r.log.Info("Request: ", "req", req)

	// Get the fluxSetup
	err := r.Client.Get(ctx, req.NamespacedName, &setup)
	if err != nil {
		if errors.IsNotFound(err) {
			r.log.Info("🌞 FluxSetup resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		r.log.Info("🌞 Failed to get FluxSetup resource. Re-running reconcile.")
		return ctrl.Result{}, err
	}

	// This currently just shows defaults
	setup.SetDefaults()
	r.log.Info("🌞 Found FluxSetup", "Name: ", setup.Name, "Namespace:", setup.Namespace, "MaxSize:", setup.Spec.MaxSize)

	// TODO this Status function needs some way to get total jobs (kueue uses a cache)
	// This should be linked with the ability to control / limit resources
	status, err := r.Status(&setup)
	if !equality.Semantic.DeepEqual(status, setup.Status) {
		setup.Status = status
		err := r.Client.Status().Update(ctx, &setup)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	return ctrl.Result{}, nil
}

/*
	flux.SetDefaults()
	instance.SetDefaults()

	log.Info("🥑️ Found instance 🥑️", "Flux Image: ", flux.Spec.Image, "Size: ", fmt.Sprint(instance.Spec.Size))
	fmt.Printf("\n🪵 EtcHosts Hostfile \n%s\n", instance.Spec.EtcHosts.Hostfile)

	// Ensure the configs are created (for volume sources)
	// The hostfile here is empty because we generate it entirely
	_, result, err := r.getHostfileConfig(ctx, &instance, "flux-config", "")
	if err != nil {
		return result, err
	}
	_, result, err = r.getHostfileConfig(ctx, &instance, "etc-hosts", instance.Spec.EtcHosts.Hostfile)
	if err != nil {
		return result, err
	}

	// And generate the secret curve cert
	_, result, err = r.getCurveCert(ctx, &instance)
	if err != nil {
		return result, err
	}

	// Get existing deployment (statefulset, a result, and error)
	_, result, err = r.getStatefulSet(ctx, &instance, flux.Spec.Image)
	if err != nil {
		return result, err
	}*/

// STATUS
func (r *FluxSetupReconciler) Status(setup *api.FluxSetup) (api.FluxSetupStatus, error) {
	return api.FluxSetupStatus{
		UsedResources: 1, // TODO
		WaitingJobs:   int32(r.fluxManager.JobsPending()),
		RunningJobs:   int32(r.fluxManager.JobsRunning()),
	}, nil
}

// WATCHERS
// NotifyJobUpdate is called from FluxJob when there is a new Job
// It adds an event to the update channel
func (r *FluxSetupReconciler) NotifyJobUpdate(job *api.FluxJob) {
	r.log.Info("🌞 FluxSetup is being notified of a FluxJob update")
	r.jobUpdateChannel <- event.GenericEvent{Object: job}
}

// EVENTS
// The functions below are added via WithEventFilter, and determine
// if we call reconcile or not (by returning true/false)

// Create is responsible for registering a new Queue for the flux manager
// This could eventually be extended to create more than one, but I'm
// starting with one for now.
func (r *FluxSetupReconciler) Create(e event.CreateEvent) bool {

	// This should only be responding to FluxSetup events
	setup, match := e.Object.(*api.FluxSetup)
	if !match {
		// No need to interact with the cache for other objects.
		return true
	}
	log := r.log.WithValues("FluxSetup", klog.KObj(setup))
	log.Info("🌞 FluxSetup create event")
	ctx := ctrl.LoggerInto(context.Background(), log)

	// Add the new setup to the manager
	if err := r.fluxManager.InitQueue(ctx, setup); err != nil {
		log.Error(err, "🌞 Failed to init Flux Manager queue")
	}
	log.Info("🌞 Flux Manager queue created in FluxSetup Create, asking for Reconcile")
	return true
}

func (r *FluxSetupReconciler) Delete(e event.DeleteEvent) bool {
	setup, match := e.Object.(*api.FluxSetup)
	if !match {
		// No need to interact with the cache for other objects.
		return true
	}
	log := r.log.WithValues("FluxSetup", klog.KObj(setup))
	log.Info("🌞 FluxSetup delete event")

	/*defer r.notifyWatchers(cq, nil)

	r.log.V(2).Info("ClusterQueue delete event", "clusterQueue", klog.KObj(cq))
	r.cache.DeleteClusterQueue(cq)
	r.qManager.DeleteClusterQueue(cq)*/
	return true
}

func (r *FluxSetupReconciler) Update(e event.UpdateEvent) bool {
	_, match := e.ObjectOld.(*api.FluxSetup)
	newSetup, newMatch := e.ObjectNew.(*api.FluxSetup)
	r.log.Info("🌞 FluxSetup Update Event", "setup", klog.KObj(newSetup))

	if !match || !newMatch {
		// No need to interact with the cache for other objects.
		return true
	}

	log := r.log.WithValues("FluxSetup", klog.KObj(newSetup))
	log.Info("🌞 FluxSetup update event")

	if newSetup.DeletionTimestamp != nil {
		return true
	}
	/*defer r.notifyWatchers(oldCq, newCq)

	if err := r.cache.UpdateClusterQueue(newCq); err != nil {
		log.Error(err, "Failed to update clusterQueue in cache")
	}
	if err := r.qManager.UpdateClusterQueue(context.Background(), newCq); err != nil {
		log.Error(err, "Failed to update clusterQueue in queue manager")
	}*/
	return true
}

func (r *FluxSetupReconciler) Generic(e event.GenericEvent) bool {
	r.log.Info("🌞 FluxSetup Generic event", "setup", klog.KObj(e.Object))
	return true
}

// jobHandler signals the controller to reconcile the FluxSetup
// It is watched via the reconciler, and triggered by updated to the jobupdate channel
// The events come from a channel Source, so only the Generic handler will get events.
type jobHandler struct {
	fluxManager *flux.Manager
	log         logr.Logger
}

func (h *jobHandler) Create(event.CreateEvent, workqueue.RateLimitingInterface) {}
func (h *jobHandler) Update(event.UpdateEvent, workqueue.RateLimitingInterface) {}
func (h *jobHandler) Delete(event.DeleteEvent, workqueue.RateLimitingInterface) {}

// Generic adds a request for a new job
func (h *jobHandler) Generic(e event.GenericEvent, q workqueue.RateLimitingInterface) {
	job := e.Object.(*api.FluxJob)

	condition := jobctrl.GetCondition(job)
	h.log.Info("🎧 FluxSetup Job Update Channel received job", "Name:", job.Name, "Size", job.Spec.Size, "Condition:", condition)

	// If the condition is waiting for resources, we should check if we can admit
	if condition == jobctrl.ConditionJobWaiting {

		// This is where the job should come in "waiting for resources" and we should
		// provide them and TODO update the status of the cluster if they are available
		// TODO get current queue status and compare to size of job
		// Right now we assume we can accept infinite!
		if h.fluxManager.AddOrUpdateJob(job) {
			h.log.Info("🎧 FluxSetup Job Update Channel Job Accepted", "Name:", job.Name, "Size", job.Spec.Size, "Condition:", condition)
			req := h.requestForJob(job)
			if req != nil {
				q.AddAfter(*req, defaults.UpdatesBatchPeriod)
			}
		}
	}
}

// requestForJob ensures that we reconcile when there is a new job request created
func (h *jobHandler) requestForJob(job *api.FluxJob) *reconcile.Request {
	// TODO likely we want to set defaults here
	// TODO should we set a uuid job name here?
	return &reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name: job.Name,
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *FluxSetupReconciler) SetupWithManager(mgr ctrl.Manager) error {

	// We pass the flux manager to the update handler
	jobUpdateHandler := jobHandler{
		fluxManager: r.fluxManager,
		log:         r.log,
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.FluxSetup{}).
		//		Watches(&source.Kind{Type: &api.FluxJob{}}, &handler.EnqueueRequestForObject{}).
		Owns(&batchv1.Job{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{}).
		// Defaults to 1, putting here so we know it exists!
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).

		// This channel is updated by FluxJob, so watching it ensures we watch FluxJob
		Watches(&source.Channel{Source: r.jobUpdateChannel}, &jobUpdateHandler).
		WithEventFilter(r).
		Complete(r)
}
