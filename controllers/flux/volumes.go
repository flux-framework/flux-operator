/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"context"
	api "flux-framework/flux-operator/api/v1alpha1"
	"fmt"
	"path"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	entrypointSuffix  = "-entrypoint"
	fluxConfigSuffix  = "-flux-config"
	curveVolumeSuffix = "-curve-mount"
)

// Shared function to return consistent set of volume mounts
func getVolumeMounts(cluster *api.MiniCluster) []corev1.VolumeMount {
	mounts := []corev1.VolumeMount{
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
	// Are we expecting a munge secret?
	if cluster.Spec.Flux.MungeSecret != "" {
		mungeMount := corev1.VolumeMount{
			Name:      cluster.Spec.Flux.MungeSecret,
			MountPath: "/etc/munge",
			ReadOnly:  true,
		}
		mounts = append(mounts, mungeMount)
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
	}

	// If we have Multi User mode, we need to set permission 0644
	brokerFile := corev1.KeyToPath{
		Key:  "hostfile",
		Path: "broker.toml",
	}

	// If we need the munge.key
	mungeKey := corev1.KeyToPath{
		Key:  "munge.key",
		Path: "munge.key",
	}

	// /mnt/curve/curve.cert
	curveKey := corev1.KeyToPath{
		Key:  curveCertKey,
		Path: "curve.cert",
	}

	if cluster.MultiUser() {
		mode := int32(0644)
		brokerFile = corev1.KeyToPath{
			Key:  "hostfile",
			Path: "broker.toml",
			Mode: &mode,
		}
	}

	// Defaults volumes we always write - entrypoint and configs
	volumes := []corev1.Volume{
		{
			Name: cluster.Name + fluxConfigSuffix,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cluster.Name + fluxConfigSuffix,
					},
					// /etc/flux/config
					Items: []corev1.KeyToPath{brokerFile},
				},
			},
		},
		{
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
		},
	}

	// We either generate a curve.cert config map, or get it from secret
	curveVolumeName := cluster.Name + fluxConfigSuffix
	if cluster.Spec.Flux.CurveCertSecret != "" {
		curveVolume := corev1.Volume{
			Name: curveVolumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: cluster.Spec.Flux.CurveCertSecret,
					Items:      []corev1.KeyToPath{curveKey},
				},
			},
		}
		volumes = append(volumes, curveVolume)

	} else {
		curveVolume := corev1.Volume{
			Name: curveVolumeName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{

					// Namespace based on the cluster
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cluster.Name + curveVolumeSuffix,
					},
				},
			},
		}
		volumes = append(volumes, curveVolume)
	}

	// Are we expecting a munge config map?
	if cluster.Spec.Flux.MungeSecret != "" {
		mungeVolume := corev1.Volume{
			Name: cluster.Name + fluxConfigSuffix,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: cluster.Spec.Flux.MungeSecret,
					Items:      []corev1.KeyToPath{mungeKey},
				},
			},
		}
		volumes = append(volumes, mungeVolume)
	}
	// Add volumes that already exist (not created by the Flux Operator)
	// These are unique names and path/claim names across containers
	// This can be a claim, secret, or config map
	existingVolumes := getExistingVolumes(cluster.ExistingContainerVolumes())
	volumes = append(volumes, existingVolumes...)

	// Add claims for storage types requested
	// These are created by the Flux Operator and should be Container Storage Interfaces (CSI)
	for volumeName := range cluster.Spec.Volumes {
		newVolume := corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: fmt.Sprintf("%s-claim", volumeName),
				},
			},
		}
		volumes = append(volumes, newVolume)
	}
	return volumes
}

// Get Existing volumes for the MiniCluster
func getExistingVolumes(existing map[string]api.MiniClusterExistingVolume) []corev1.Volume {
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

// createPersistentVolume creates a volume in /mnt
func (r *MiniClusterReconciler) createPersistentVolume(
	cluster *api.MiniCluster,
	volumeName string,
	volume api.MiniClusterVolume,
) *corev1.PersistentVolume {

	// We either support a hostpath (miniKube) or a Container Storage Interface (CSI)
	var pvsource corev1.PersistentVolumeSource
	if volume.StorageClass == "hostpath" {

		pvsource = corev1.PersistentVolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: path.Join(volume.Path),
			},
		}

	} else {

		// VolumeHandle defaults to storage class name
		// unless it is explicitly different!
		volumeHandle := volume.StorageClass
		if volume.VolumeHandle != "" {
			volumeHandle = volume.VolumeHandle
		}
		pvsource = corev1.PersistentVolumeSource{
			CSI: &corev1.CSIPersistentVolumeSource{

				// Choose for the user for now.
				Driver: volume.Driver,

				// Name in storageclass metadata, also what we use for name
				VolumeHandle: volumeHandle,
				NodePublishSecretRef: &corev1.SecretReference{
					Namespace: volume.SecretNamespace,
					Name:      volume.Secret,
				},
				ControllerPublishSecretRef: &corev1.SecretReference{
					Namespace: volume.SecretNamespace,
					Name:      volume.Secret,
				},
				NodeStageSecretRef: &corev1.SecretReference{
					Namespace: volume.SecretNamespace,
					Name:      volume.Secret,
				},
				VolumeAttributes: volume.Attributes,
			},
		}
	}

	newVolume := &corev1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:        volumeName,
			Namespace:   cluster.Namespace,
			Annotations: volume.Annotations,
			Labels:      volume.Labels,
		},

		Spec: corev1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
			AccessModes:                   []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},

			// This is a path in the minikube vm or on the node
			PersistentVolumeSource: pvsource,
			StorageClassName:       volume.StorageClass,
		},
	}
	// Capacity is optional for some storage like Google Cloud
	if volume.Capacity != "" {
		newVolume.Spec.Capacity = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceStorage: resource.MustParse(volume.Capacity),
		}
	}
	return newVolume
}

func (r *MiniClusterReconciler) getExistingPersistentVolume(
	ctx context.Context,
	cluster *api.MiniCluster,
	volumeName string,
) (*corev1.PersistentVolume, error) {

	// First look for an existing persistent volume
	existing := &corev1.PersistentVolume{}
	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      volumeName,
			Namespace: cluster.Namespace},
		existing,
	)
	return existing, err
}

// getPersistentVolume creates the PV for the curve certificate (to be written once)
func (r *MiniClusterReconciler) getPersistentVolume(
	ctx context.Context,
	cluster *api.MiniCluster,
	volumeName string,
	volume api.MiniClusterVolume) (*corev1.PersistentVolume, ctrl.Result, error,
) {

	existing, err := r.getExistingPersistentVolume(ctx, cluster, volumeName)

	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {

			// The volume "<name>-volume" will be under /mnt/<name>
			volume := r.createPersistentVolume(cluster, volumeName, volume)
			r.log.Info(
				"✨ Creating MiniCluster Mounted Volume ✨",
				"Type", volumeName,
				"Namespace", cluster.Namespace,
				"Name", volumeName,
			)

			err = r.New(ctx, volume)
			if err != nil {
				r.log.Error(
					err, "Failed to create MiniCluster Mounted Volume",
					"Type", volumeName,
					"Namespace", cluster.Namespace,
					"Name", volumeName,
				)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return volume, ctrl.Result{Requeue: true}, nil

		} else if err != nil {
			r.log.Error(err, "Failed to get MiniCluster Mounted Volume")
			return existing, ctrl.Result{}, err
		}

	} else {

		r.log.Info(
			"🎉 Found existing MiniCluster Mounted Volume",
			"Type", volumeName,
			"Namespace", existing.Namespace,
			"Name", existing.Name,
		)
	}

	// Always set owner to controller, whether created or found
	ctrl.SetControllerReference(cluster, existing, r.Scheme)
	return existing, ctrl.Result{}, err
}

// getExistingVolumeClaim gets an existing volume claim
func (r *MiniClusterReconciler) getExistingPersistentVolumeClaim(
	ctx context.Context,
	cluster *api.MiniCluster,
	claimName string,
) (*corev1.PersistentVolumeClaim, error) {

	existing := &corev1.PersistentVolumeClaim{}
	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      claimName,
			Namespace: cluster.Namespace},
		existing,
	)
	return existing, err
}

// getPersistentVolume creates the PVC claim for the curve certificate (to be written once)
func (r *MiniClusterReconciler) getPersistentVolumeClaim(
	ctx context.Context,
	cluster *api.MiniCluster,
	volumeName string,
	volume api.MiniClusterVolume,
) (*corev1.PersistentVolumeClaim, ctrl.Result, error) {

	claimName := fmt.Sprintf("%s-claim", volumeName)
	existing, err := r.getExistingPersistentVolumeClaim(ctx, cluster, claimName)

	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {
			volume := r.createPersistentVolumeClaim(cluster, claimName, volume)
			r.log.Info(
				"✨ Creating MiniCluster Mounted Volume ✨",
				"Type", claimName,
				"Namespace", cluster.Namespace,
				"Name", claimName,
			)
			err = r.New(ctx, volume)

			// This is a creation error we need to report back
			if err != nil {
				r.log.Error(
					err, "❌ Failed to create MiniCluster Mounted Volume",
					"Type", claimName,
					"Namespace", volume.Namespace,
					"Name", claimName,
				)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return volume, ctrl.Result{Requeue: true}, nil

		} else if err != nil {
			r.log.Error(err, "Failed to get MiniCluster Mounted Volume")
			return existing, ctrl.Result{}, err
		}

	} else {

		r.log.Info("🎉 Found existing MiniCluster Mounted Volume Claim",
			"Type", claimName,
			"Namespace", existing.Namespace,
			"Name", existing.Name,
		)
	}
	return existing, ctrl.Result{}, err
}

// createPersistentVolumeClaim generates a PVC
func (r *MiniClusterReconciler) createPersistentVolumeClaim(
	cluster *api.MiniCluster,
	volumeName string,
	volume api.MiniClusterVolume,
) *corev1.PersistentVolumeClaim {

	volumeMode := corev1.PersistentVolumeFilesystem

	// Create a new RWX persistent volume claim
	newVolume := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:        volumeName,
			Namespace:   cluster.Namespace,
			Annotations: volume.ClaimAnnotations,
		},

		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &volume.StorageClass,
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			VolumeMode: &volumeMode,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(volume.Capacity),
				},
			},
		},
	}
	ctrl.SetControllerReference(cluster, newVolume, r.Scheme)
	return newVolume
}
