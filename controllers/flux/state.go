/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	corev1 "k8s.io/api/core/v1"

	api "github.com/flux-framework/flux-operator/api/v1alpha1"
)

// createContainerLifecycle adds lifecycle commands to help with moving cluster state
func (r *MiniClusterReconciler) createContainerLifecycle(
	container api.MiniClusterContainer) *corev1.Lifecycle {

	// Empty Lifecycle by default
	lifecycle := corev1.Lifecycle{}

	// Manual lifecycles from the user before container start
	if container.LifeCycle.PostStartExec != "" {
		r.log.Info("🌀 MiniCluster", "LifeCycle.PostStartExec", container.LifeCycle.PostStartExec)
		lifecycle.PostStart = &corev1.LifecycleHandler{
			Exec: &corev1.ExecAction{
				Command: []string{container.LifeCycle.PostStartExec},
			},
		}
	}

	if container.LifeCycle.PreStopExec != "" {
		r.log.Info("🌀 MiniCluster", "LifeCycle.PreStopExec", container.LifeCycle.PreStopExec)
		lifecycle.PreStop = &corev1.LifecycleHandler{
			Exec: &corev1.ExecAction{
				Command: []string{container.LifeCycle.PreStopExec},
			},
		}
	}
	return &lifecycle
}
