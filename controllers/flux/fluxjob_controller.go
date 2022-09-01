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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// FluxJobReconciler reconciles a FluxJob object
type FluxJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=flux-framework.org,resources=fluxes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=fluxes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=flux-framework.org,resources=fluxes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// We compare the state of the Flux object to the actual cluster state and
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *FluxJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO:
	// 1. Get the original job CRD
	// 2. Generate a unique identifier for it - the name plus mangled string
	// 3. Create a MiniCluster CRD that has status PENDING waiting for resources, with the id
	// 4. keep reconciling until we can get that MiniCluster and see it has status READY
	// 5. When ready, run the job!
	// 6. Continue reconciling until we can get STATUS-> Complete (the FLuxSetup should be watching)
	// 7. Then we're done.
	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FluxJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.FluxJob{}).
		Complete(r)
}
