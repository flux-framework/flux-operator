package controllers

// This file has extra (not used) functions that might be useful
// (and I didn't want to delete just yet)

import (
	"context"
	"path"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// getMiniClusterIPS was used when we needed to write /etc/hosts and is no longer used
func (r *MiniClusterReconciler) getMiniClusterIPS(ctx context.Context, cluster *api.MiniCluster) map[string]string {

	ips := map[string]string{}
	for _, pod := range r.getMiniClusterPods(ctx, cluster).Items {
		// Skip init pods
		if strings.Contains(pod.Name, "init") {
			continue
		}

		// The pod isn't ready!
		if pod.Status.PodIP == "" {
			continue
		}
		ips[pod.Name] = pod.Status.PodIP
	}
	return ips
}

// getPersistentVolume creates the PVC claim for the curve certificate (to be written once)
func (r *MiniClusterReconciler) getPersistentVolumeClaim(ctx context.Context, cluster *api.MiniCluster, configFullName string) (*corev1.PersistentVolumeClaim, ctrl.Result, error) {

	existing := &corev1.PersistentVolumeClaim{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: configFullName, Namespace: cluster.Namespace}, existing)
	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {
			volume := r.createPersistentVolumeClaim(cluster, configFullName)
			r.log.Info("✨ Creating MiniCluster Mounted Volume ✨", "Type", configFullName, "Namespace", volume.Namespace, "Name", volume.Name)
			err = r.Client.Create(ctx, volume)
			if err != nil {
				r.log.Error(err, "❌ Failed to create MiniCluster Mounted Volume", "Type", configFullName, "Namespace", volume.Namespace, "Name", (*volume).Name)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return volume, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			r.log.Error(err, "Failed to get MiniCluster Mounted Volume")
			return existing, ctrl.Result{}, err
		}
	} else {
		r.log.Info("🎉 Found existing MiniCluster Mounted Volume", "Type", configFullName, "Namespace", existing.Namespace, "Name", existing.Name)
	}
	return existing, ctrl.Result{}, err
}

// getPersistentVolume creates the PV for the curve certificate (to be written once)
func (r *MiniClusterReconciler) getPersistentVolume(ctx context.Context, cluster *api.MiniCluster, configFullName string) (*corev1.PersistentVolume, ctrl.Result, error) {

	existing := &corev1.PersistentVolume{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: configFullName, Namespace: cluster.Namespace}, existing)
	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {
			volume := r.createPersistentVolume(cluster, configFullName)
			r.log.Info("✨ Creating MiniCluster Mounted Volume ✨", "Type", configFullName, "Namespace", volume.Namespace, "Name", volume.Name)
			err = r.Client.Create(ctx, volume)
			if err != nil {
				r.log.Error(err, "❌ Failed to create MiniCluster Mounted Volume", "Type", configFullName, "Namespace", volume.Namespace, "Name", (*volume).Name)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return volume, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			r.log.Error(err, "Failed to get MiniCluster Mounted Volume")
			return existing, ctrl.Result{}, err
		}
	} else {
		r.log.Info("🎉 Found existing MiniCluster Mounted Volume", "Type", configFullName, "Namespace", existing.Namespace, "Name", existing.Name)
	}
	return existing, ctrl.Result{}, err
}

// createPersistentVolumeClaim generates a PVC
// This tends to choke on MiniKube, I'm not sure it has a provisioner?
func (r *MiniClusterReconciler) createPersistentVolumeClaim(cluster *api.MiniCluster, configName string) *corev1.PersistentVolumeClaim {
	volume := &corev1.PersistentVolumeClaim{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: configName, Namespace: cluster.Namespace},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},

			// No idea how much to ask for here! I made it up.
			Resources: corev1.ResourceRequirements{Requests: corev1.ResourceList{
				corev1.ResourceStorage: *resource.NewQuantity(1024, resource.BinarySI),
			}},
		},
	}
	ctrl.SetControllerReference(cluster, volume, r.Scheme)
	return volume
}

// createPersistentVolume creates a volume in /tmp, which doesn't seem to choke
func (r *MiniClusterReconciler) createPersistentVolume(cluster *api.MiniCluster, configName string) *corev1.PersistentVolume {
	volume := &corev1.PersistentVolume{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: configName, Namespace: cluster.Namespace},
		Spec: corev1.PersistentVolumeSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
			Capacity: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceStorage: *resource.NewQuantity(1024, resource.BinarySI),
			},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: path.Join("/tmp/", configName),
				},
			},
			StorageClassName: "manual",
		},
	}
	ctrl.SetControllerReference(cluster, volume, r.Scheme)
	return volume
}
