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

	corev1 "k8s.io/api/core/v1"
)

const (
	entrypointSuffix = "-entrypoint"
	fluxConfigSuffix = "-flux-config"
	curveAuthSuffix  = "-curve-auth"
)

// Shared function to return consistent set of volume mounts
// for the MiniCluster and Flux Statefulset
func getVolumeMounts(cluster *api.MiniCluster) []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      cluster.Name + curveAuthSuffix,
			MountPath: "/mnt/curve/",
		},
		{
			Name:      cluster.Name + fluxConfigSuffix,
			MountPath: "/etc/flux/config",
			ReadOnly:  true,
		},

		// Entrypoint that helps to discover hosts, added after creation
		{
			Name:      cluster.Name + entrypointSuffix,
			MountPath: "/flux_operator/",
			ReadOnly:  false,
		},
	}
}

// getVolumes that are shared between MiniCluster and statefulset
func getVolumes(cluster *api.MiniCluster) []corev1.Volume {
	makeExecutable := int32(0777)
	return []corev1.Volume{{
		Name: cluster.Name + fluxConfigSuffix,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cluster.Name + fluxConfigSuffix,
				},
				// /etc/flux/config
				Items: []corev1.KeyToPath{{
					Key:  "hostfile",
					Path: "brokers.toml",
				}},
			},
		},
	}, {

		// We use an empty volume (that can be shared by the container and init container)
		// to run flux keygen and generate the /mnt/curve/curve.crt
		Name: cluster.Name + curveAuthSuffix,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}, {
		Name: cluster.Name + entrypointSuffix,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{

				// Namespace based on the cluster
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cluster.Name + entrypointSuffix,
				},
				// /flux_operator/start.sh
				// /flux_operator/wait.sh
				Items: []corev1.KeyToPath{{
					Key:  "start-flux",
					Path: "start.sh",
					Mode: &makeExecutable,
				}, {
					Key:  "wait",
					Path: "wait.sh",
					Mode: &makeExecutable,
				}, {
					Key:  "update-hosts",
					Path: "update_hosts.sh",
					Mode: &makeExecutable,
				}},
			},
		},
	}}
}
