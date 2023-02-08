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

	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	setupLog = ctrl.Log.WithName("setup")
)

// SetupControllers sets up all controllers.
func SetupControllers(mgr ctrl.Manager, restClient rest.Interface) (string, error) {

	jobReconciler := controllers.NewMiniClusterReconciler(
		mgr.GetClient(),
		mgr.GetScheme(),
		*(mgr.GetConfig()),
		restClient,
		// other watching reconcilers could be added here!
	)

	if err := jobReconciler.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MiniCluster")
		return "MiniCluster", err
	}
	return "", nil
}
