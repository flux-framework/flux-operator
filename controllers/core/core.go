/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	controllers "flux-framework/flux-operator/controllers/flux"
	"flux-framework/flux-operator/pkg/flux"

	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	setupLog = ctrl.Log.WithName("setup")
)

// SetupControllers sets up all controllers.
func SetupControllers(mgr ctrl.Manager, manager *flux.Manager) (string, error) {

	// Admin (internal) Flux Setup Reconciler (setup first!)
	setupReconciler := controllers.NewFluxSetupReconciler(mgr.GetClient(), mgr.GetScheme(), manager)
	if err := setupReconciler.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FluxSetup")
		return "FluxSetup", err
	}

	// User facing Flux Reconciler - receives the job.
	// We provide the setupReconciler as a watcher
	jobReconciler := controllers.NewFluxJobReconciler(mgr.GetClient(), mgr.GetScheme(), manager, setupReconciler)
	if err := jobReconciler.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FluxJob")
		return "FluxJob", err
	}
	return "", nil
}
