package controllers

/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

import (
	"context"

	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/api/errors"

	api "flux-framework/flux-operator/api/v1alpha1"
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

func NewMiniClusterReconciler(
	client client.Client,
	scheme *runtime.Scheme,
	restConfig rest.Config,
	restClient rest.Interface,
	watchers ...MiniClusterUpdateWatcher,
) *MiniClusterReconciler {

	return &MiniClusterReconciler{
		log:        ctrl.Log.WithName("minicluster-reconciler"),
		Client:     client,
		Scheme:     scheme,
		watchers:   watchers,
		RESTClient: restClient,
		RESTConfig: &restConfig,
	}
}

// RBAC rules to access cluster-api resources
//+kubebuilder:rbac:groups=flux-framework.org,resources=clusters;clusters/status,verbs=get;list;watch
//+kubebuilder:rbac:groups=flux-framework.org,resources=machines;machines/status;machinedeployments;machinedeployments/status;machinesets;machinesets/status;machineclasses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=clusters;clusters/status,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=nodes;events,verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups=flux-framework.org,resources=miniclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=miniclusters/status,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=miniclusters/finalizers,verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/log,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/exec,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=persistentvolumes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources="",verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=batch,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=core,resources=networks,verbs=create;patch
//+kubebuilder:rbac:groups=core,resources="services",verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources="ingresses",verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups="",resources=events,verbs=create;watch;update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete;exec
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;list;watch;create;update;patch;delete;exec

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// We compare the state of the Flux object to the actual cluster state and
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *MiniClusterReconciler) Reconcile(
	ctx context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {

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

	// Don't continue if they provided 0 size, that makes no sense!
	if cluster.Spec.Size == 0 {
		r.log.Info("🌀 A MiniCluster without nodes? Is this a cluster for invisible ants? Canceling!")
		return ctrl.Result{}, nil
	}

	// Show parameters provided and validate one flux runner
	if !cluster.Validate() {
		r.log.Info("🌀 Your MiniCluster config did not validate! see the sad faces above for details. Canceling!")
		return ctrl.Result{}, nil
	}
	r.log.Info("🌀 Reconciling Mini Cluster", "Containers: ", len(cluster.Spec.Containers))

	// Ensure we have the minicluster (get or create!)
	result, err := r.ensureMiniCluster(ctx, &cluster)
	if err != nil {
		return result, err
	}

	// By the time we get here we have a Job + pods + config maps!
	// What else do we want to do?
	r.log.Info("🌀 Mini Cluster is Ready!")

	// Check until the job finishes to clean up volumes if needed
	if cluster.Spec.Cleanup {
		result, err := r.cleanupPodsStorage(ctx, &cluster)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MiniClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.MiniCluster{}).

		// This references the Create/Delete/Update,etc functions above
		// they return a boolean to indicate if we should reconcile given the event
		// If we don't need these extra filters we can delete this line and events.go
		WithEventFilter(r).
		Owns(&networkv1.Ingress{}).
		Owns(&batchv1.Job{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
