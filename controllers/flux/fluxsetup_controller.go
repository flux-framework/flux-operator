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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	logctrl "sigs.k8s.io/controller-runtime/pkg/log"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// FluxSetupReconciler reconciles a FluxSetup object
type FluxSetupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
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

// Reconcile moves the current state of the cluster closer to the desired state.
func (r *FluxSetupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// Create a new FluxSetup and Flux instance
	var instance api.FluxSetup
	var flux api.Flux

	// Prepare a logger to communicate to the developer user
	// Note that we could attach a named logger to the reconciler object,
	// and that might be a better approach for organization or state
	// https://github.com/kubernetes-sigs/kueue/blob/main/pkg/controller/core/queue_controller.go#L50
	log := logctrl.FromContext(ctx).WithValues("FluxSetup", req.NamespacedName)

	// Keep developed informed what is going on.
	log.Info("⚡️ Event received! ⚡️")
	log.Info("Request: ", "req", req)

	// Get the flux instance
	err := r.Get(ctx, req.NamespacedName, &flux)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Flux resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Info("Failed to get Flux resource. Re-running reconcile.")
		return ctrl.Result{}, err
	}

	// Get the fluxSetup
	err = r.Get(ctx, req.NamespacedName, &instance)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("FluxSetup resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Info("Failed to get FluxSetup resource. Re-running reconcile.")
		return ctrl.Result{}, err
	}
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
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FluxSetupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.FluxSetup{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{}).
		// Defaults to 1, putting here so we know it exists!
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}
