/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"fmt"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"

	corev1 "k8s.io/api/core/v1"
)

// Shared function to return consistent set of volume mounts
func getVolumeMounts(cluster *api.MiniCluster) []corev1.VolumeMount {
	mounts := []corev1.VolumeMount{
		// The empty volume for Flux will go here
		{
			Name:      cluster.Spec.Flux.Container.Name,
			MountPath: cluster.Spec.Flux.Container.MountPath,
			ReadOnly:  false,
		},

		// Entrypoints will go here
		{
			Name:      cluster.EntrypointConfigMapName(),
			MountPath: "/flux_operator/",
			ReadOnly:  true,
		},
	}
	return mounts
}

// getVolumes that are shared between MiniCluster and statefulset
func getVolumes(cluster *api.MiniCluster) []corev1.Volume {

	// Runner start scripts
	makeExecutable := int32(0777)
	runnerStartScripts := []corev1.KeyToPath{}

	// Prepare a custom "wait.sh" for each container based on index
	for i, container := range cluster.Spec.Containers {

		// For now, only Flux runners get the custom wait.sh script
		if container.RunFlux {
			startScript := corev1.KeyToPath{
				Key:  fmt.Sprintf("wait-%d", i),
				Path: fmt.Sprintf("wait-%d.sh", i),
				Mode: &makeExecutable,
			}
			runnerStartScripts = append(runnerStartScripts, startScript)
		}

		// A non flux container can also handle custom logic, if command is provided
		if container.GenerateEntrypoint() {
			startScript := corev1.KeyToPath{
				Key:  fmt.Sprintf("start-%d", i),
				Path: fmt.Sprintf("start-%d.sh", i),
				Mode: &makeExecutable,
			}
			runnerStartScripts = append(runnerStartScripts, startScript)
		}
	}

	// /flux_operator/curve.cert
	curveKey := corev1.KeyToPath{
		Key:  CurveCertKey,
		Path: "curve.cert",
	}

	// Add the flux init script
	fluxScript := corev1.KeyToPath{
		Key:  cluster.Spec.Flux.Container.Name,
		Path: "flux-init.sh",
		Mode: &makeExecutable,
	}
	runnerStartScripts = append(runnerStartScripts, fluxScript)
	runnerStartScripts = append(runnerStartScripts, curveKey)

	// Defaults volumes we always write - entrypoint and empty volume
	volumes := []corev1.Volume{
		{
			Name: cluster.Spec.Flux.Container.Name,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
		{
			Name: cluster.EntrypointConfigMapName(),
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{

					// Namespace based on the cluster
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cluster.EntrypointConfigMapName(),
					},
					// /flux_operator/wait-<index>.sh
					Items: runnerStartScripts,
				},
			},
		},
	}

	// Add volumes that already exist (not created by the Flux Operator)
	// These are unique names and path/claim names across containers
	// This can be a claim, secret, or config map
	existingVolumes := getExistingVolumes(cluster.ExistingContainerVolumes())
	volumes = append(volumes, existingVolumes...)
	return volumes
}

// Get Existing volumes for the MiniCluster
func getExistingVolumes(existing map[string]api.ContainerVolume) []corev1.Volume {
	volumes := []corev1.Volume{}
	for volumeName, volumeMeta := range existing {

		var newVolume corev1.Volume
		if volumeMeta.SecretName != "" {
			newVolume = corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: volumeMeta.SecretName,
					},
				},
			}

		} else if volumeMeta.EmptyDir {

			// The Flux Operator supports default and memory
			medium := corev1.StorageMediumDefault
			if volumeMeta.EmptyDirMedium == "memory" {
				medium = corev1.StorageMediumMemory
			}
			newVolume = corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{
						Medium: medium,
					},
				},
			}

		} else if volumeMeta.HostPath != "" {
			newVolume = corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					// Empty path for type means no checks are done
					HostPath: &corev1.HostPathVolumeSource{
						Path: volumeMeta.Path,
					},
				},
			}

		} else if volumeMeta.ConfigMapName != "" {

			// Prepare items as key to path
			items := []corev1.KeyToPath{}
			for key, path := range volumeMeta.Items {
				newItem := corev1.KeyToPath{
					Key:  key,
					Path: path,
				}
				items = append(items, newItem)
			}

			// This is a config map volume with items
			newVolume = corev1.Volume{
				Name: volumeMeta.ConfigMapName,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: volumeMeta.ConfigMapName,
						},
						Items: items,
					},
				},
			}

		} else {

			// Fall back to persistent volume claim
			newVolume = corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: volumeMeta.ClaimName,
					},
				},
			}
		}
		volumes = append(volumes, newVolume)
	}
	return volumes
}
