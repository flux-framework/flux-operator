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
	etcHostsSuffix   = "-etc-hosts"
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
			ReadOnly:  true,
		},

		// Disabled for now too - /etc/flux also needs write
		{
			Name:      cluster.Name + fluxConfigSuffix,
			MountPath: "/etc/flux/config",
			ReadOnly:  true,
		},

		// Disabled for now - not sure we want to do this because the container
		// is mounting stuff there too, and wouldn't this be controlled by the operator?
		//		{
		//			Name:      cluster.Name + etcHostsSuffix,
		//			MountPath: "/etc/",
		//			ReadOnly:  false,
		//		},
	}
}

// getVolumes that are shared between MiniCluster and statefulset
func getVolumes(cluster *api.MiniCluster) []corev1.Volume {
	permMode := int32(0600)
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
		Name: cluster.Name + etcHostsSuffix,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{

				// Namespace based on the cluster
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cluster.Name + etcHostsSuffix,
				},
				// /etc/hosts
				Items: []corev1.KeyToPath{{
					Key:  "hostfile",
					Path: "hosts",
				}},
			},
		},
	}, {
		Name: cluster.Name + curveAuthSuffix,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: cluster.Name + curveAuthSuffix,

				// /mnt/curve/curve.cert
				// /mnt/curve/curve.key
				Items: []corev1.KeyToPath{
					{
						Key:  "tls.crt",
						Path: "curve.cert",
						Mode: &permMode,
					},
					{
						Key:  "tls.key",
						Path: "curve.key",
						Mode: &permMode,
					},
				},
			},
		},
	}}
}
