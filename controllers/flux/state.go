/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"fmt"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// createContainerLifecycle adds lifecycle commands to help with moving cluster state
func (r *MiniClusterReconciler) createContainerLifecycle(
	cluster *api.MiniCluster,
	container api.MiniClusterContainer) *corev1.Lifecycle {

	// Empty Lifecycle by default
	lifecycle := corev1.Lifecycle{}

	// Manual lifecycles from the user before container start
	if container.LifeCycle.PostStartExec != "" {
		lifecycle.PostStart = &corev1.LifecycleHandler{
			Exec: &corev1.ExecAction{
				Command: []string{container.LifeCycle.PostStartExec},
			},
		}
	}

	// If this logic needs to be shared can be moved external to the function
	fluxuser := "flux"
	if container.FluxUser.Name != "" {
		fluxuser = container.FluxUser.Name
	}

	// Assemble the command
	asSudo := "sudo -E PYTHONPATH=$PYTHONPATH -E PATH=$PATH"
	asFlux := fmt.Sprintf("sudo -u %s -E PYTHONPATH=$PYTHONPATH -E PATH=$PATH -E HOME=/home/%s", fluxuser, fluxuser)
	if container.Commands.RunFluxAsRoot {
		asFlux = asSudo + fmt.Sprintf("-E HOME=/home/%s", fluxuser)
	}

	// If we have an archive path, we will need to save there
	// Note that copy FROM archive TO container happens in wait.sh
	if cluster.Spec.Archive.Path != "" {

		dirname := filepath.Dir(cluster.Spec.Archive.Path)
		preStop := []string{
			"/bin/bash", "-c",
			fmt.Sprintf("mkdir -p %s && %s flux proxy local:///var/run/flux/local flux dump %s",
				dirname, asFlux, cluster.Spec.Archive.Path),
		}
		r.log.Info("🌀 MiniCluster", "LifeCycle.PreStopExec", strings.Join(preStop, " "))
		lifecycle.PreStop = &corev1.LifecycleHandler{
			Exec: &corev1.ExecAction{
				Command: preStop,
			},
		}

	}
	return &lifecycle
}
