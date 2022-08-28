/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	corev1 "k8s.io/api/core/v1"
)

// Shared function to return consistent set of volume mounts
// for the FluxJob and Flux Statefulset
func getVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{{
		Name:      "curve-auth",
		MountPath: "/mnt/curve/",
		ReadOnly:  true,
	},
		{
			Name:      "flux-config",
			MountPath: "/etc/flux/config/",
			ReadOnly:  true,
		},
		{
			Name:      "etc-hosts",
			MountPath: "/etc/hosts",
			ReadOnly:  true,
		},
	}
}

// getVolumes that are shared between FluxJob and statefulset
func getVolumes() []corev1.Volume {

	return []corev1.Volume{{
		Name: "flux-config",
		VolumeSource: corev1.VolumeSource{

			// There wasn't a Name here so I reproduced the paths and key
			// See link and comment in TODO.md
			ConfigMap: &corev1.ConfigMapVolumeSource{
				Items: []corev1.KeyToPath{{
					Key:  "flux-config",
					Path: "/etc/flux/config",
				}},
			},
		},
	}, {
		Name: "etc-hosts",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				Items: []corev1.KeyToPath{{
					Key:  "etc-hosts",
					Path: "/etc/hosts",
				}},
			},
		},
	}, {
		Name: "curve-auth",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: "secret-tls",
				Items: []corev1.KeyToPath{{
					Key:  "curve-cert",
					Path: "curve.cert",
				},
					{
						Key:  "curve-key",
						Path: "curve.key",
					}},
			},
		},
	}}
}
