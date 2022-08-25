/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
