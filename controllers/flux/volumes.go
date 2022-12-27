/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	api "flux-framework/flux-operator/api/v1alpha1"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

const (
	entrypointSuffix  = "-entrypoint"
	fluxConfigSuffix  = "-flux-config"
	curveVolumeSuffix = "-curve-mount"
)

// Shared function to return consistent set of volume mounts
// for the MiniCluster and Flux Statefulset
func getVolumeMounts(cluster *api.MiniCluster) []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      cluster.Name + curveVolumeSuffix,
			MountPath: "/mnt/curve/",
			ReadOnly:  true,
		},
		{
			Name:      cluster.Name + fluxConfigSuffix,
			MountPath: "/etc/flux/config",
			ReadOnly:  true,
		},
		{
			Name:      cluster.Name + entrypointSuffix,
			MountPath: "/flux_operator/",
			ReadOnly:  true,
		},
	}
}

// getVolumes that are shared between MiniCluster and statefulset
func getVolumes(cluster *api.MiniCluster) []corev1.Volume {

	// Runner start scripts
	makeExecutable := int32(0777)
	runnerStartScripts := []corev1.KeyToPath{}

	// Prepare a custom "wait.sh" for each container based on index
	for i, container := range cluster.Spec.Containers {

		// For now, only Flux runners get the custom wait.sh script
		if container.FluxRunner {
			startScript := corev1.KeyToPath{Key: fmt.Sprintf("wait-%d", i), Path: fmt.Sprintf("wait-%d.sh", i), Mode: &makeExecutable}
			runnerStartScripts = append(runnerStartScripts, startScript)
		}
	}

	volumes := []corev1.Volume{{
		Name: cluster.Name + fluxConfigSuffix,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cluster.Name + fluxConfigSuffix,
				},
				// /etc/flux/config
				Items: []corev1.KeyToPath{{
					Key:  "hostfile",
					Path: "broker.toml",
				}},
			},
		},
	}, {
		Name: cluster.Name + entrypointSuffix,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{

				// Namespace based on the cluster
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cluster.Name + entrypointSuffix,
				},
				// /flux_operator/wait-<index>.sh
				Items: runnerStartScripts,
			},
		},
	}, {
		Name: cluster.Name + curveVolumeSuffix,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{

				// Namespace based on the cluster
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cluster.Name + curveVolumeSuffix,
				},
				// /mnt/curve/curve.cert
				Items: []corev1.KeyToPath{{
					Key:  curveCertKey,
					Path: "curve.cert",
				}},
			},
		},
	}}

	// If we are using a localDeploy (volume on the host) vs. a cluster deploy
	// (where we need a persistent volume claim)
	if cluster.Spec.LocalDeploy {
		directoryType := corev1.HostPathDirectoryOrCreate

		// Add local volumes available to containers
		for volumeName, volume := range cluster.Spec.Volumes {
			localVolume := corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: volume.Path,
						Type: &directoryType,
					},
				},
			}
			volumes = append(volumes, localVolume)
		}
	}
	return volumes
}
