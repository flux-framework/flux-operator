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
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	setupLog = ctrl.Log.WithName("setup")
)

// SetupControllers sets up all controllers.
func SetupControllers(mgr ctrl.Manager) (string, error) {

	// User facing Flux Reconciler - receives the job
	if err := (&controllers.FluxJobReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FluxJob")
		return "FluxJob", err
	}
	return "", nil

	// Admin (internal) Flux Setup Reconciler (setup first!)
	if err := (&controllers.FluxSetupReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FluxSetup")
		return "FluxSetup", err
	}
}
