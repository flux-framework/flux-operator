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
		if container.FluxRunner {
			startScript := corev1.KeyToPath{
				Key:  fmt.Sprintf("wait-%d", i),
				Path: fmt.Sprintf("wait-%d.sh", i),
				Mode: &makeExecutable,
			}
			runnerStartScripts = append(runnerStartScripts, startScript)
		}
	}

	volumes := []corev1.Volume{
		{
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
		{
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
		},
	}

	// Add claims for storage types requested
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

// createPersistentVolume creates a volume in /mnt
func (r *MiniClusterReconciler) createPersistentVolume(
	cluster *api.MiniCluster,
	volumeName string,
	volume api.MiniClusterVolume,
) *corev1.PersistentVolume {

	// We either support a hostpath (miniKube) or a Container Storage Interface (CSI)
	var pvsource corev1.PersistentVolumeSource
	if volume.StorageClassName == "hostpath" {

		pvsource = corev1.PersistentVolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: path.Join(volume.Path),
			},
		}

	} else {

		// VolumeHandle defaults to storage class name
		// unless it is explicitly different!
		volumeHandle := volume.StorageClassName
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
					Name:      volume.SecretReference,
				},
				ControllerPublishSecretRef: &corev1.SecretReference{
					Namespace: volume.SecretNamespace,
					Name:      volume.SecretReference,
				},
				NodeStageSecretRef: &corev1.SecretReference{
					Namespace: volume.SecretNamespace,
					Name:      volume.SecretReference,
				},
				VolumeAttributes: volume.Attributes,
			},
		}
	}

	newVolume := &corev1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      volumeName,
			Namespace: cluster.Namespace,
			Labels:    volume.Labels,
		},

		Spec: corev1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
			AccessModes:                   []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},

			// This is a path in the minikube vm or on the node
			PersistentVolumeSource: pvsource,
			StorageClassName:       volume.StorageClassName,
		},
	}
	// Capacity is optional for some storage like Google Cloud
	if volume.Capacity != "" {
		newVolume.Spec.Capacity = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceStorage: resource.MustParse(volume.Capacity),
		}
	}

	ctrl.SetControllerReference(cluster, newVolume, r.Scheme)
	return newVolume
}

func (r *MiniClusterReconciler) getExistingPersistentVolume(
	ctx context.Context,
	cluster *api.MiniCluster,
	volumeName string,
) (*corev1.PersistentVolume, error) {

	// First look for an existing persistent volume
	existing := &corev1.PersistentVolume{}
	err := r.Client.Get(
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

			err = r.Client.Create(ctx, volume)
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
	return existing, ctrl.Result{}, err
}

// getExistingVolumeClaim gets an existing volume claim
func (r *MiniClusterReconciler) getExistingPersistentVolumeClaim(
	ctx context.Context,
	cluster *api.MiniCluster,
	claimName string,
) (*corev1.PersistentVolumeClaim, error) {

	existing := &corev1.PersistentVolumeClaim{}
	err := r.Client.Get(
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
			err = r.Client.Create(ctx, volume)

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

		r.log.Info("🎉 Found existing MiniCluster Mounted Volume",
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

	// This can be explicitly set to an empty string
	pvcStorageClass := volume.StorageClassName
	//if volume.PVCStorageClassName != "pvc-storage-class-name-unset" {
	//	pvcStorageClass = volume.PVCStorageClassName
	//}

	// Create a new RWX persistent volume claim
	newVolume := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:        volumeName,
			Namespace:   cluster.Namespace,
			Annotations: volume.Annotations,
		},

		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &pvcStorageClass,
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
