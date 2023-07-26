/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	api "github.com/flux-framework/flux-operator/api/v1alpha1"
)

// newMiniCluster is used to create the MiniCluster Job
func (r *MiniClusterReconciler) newMiniClusterJob(
	cluster *api.MiniCluster,
) (*batchv1.Job, error) {

	// Number of retries before marking as failed
	backoffLimit := int32(100)
	completionMode := batchv1.IndexedCompletion
	podLabels := r.getPodLabels(cluster)
	setAsFQDN := false

	// We add the selector for the horizontal auto scaler, if active
	// We can't use the job-name selector, as this would include the
	// external sidecar service!
	podLabels["hpa-selector"] = cluster.Name

	// This is an indexed-job
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
			Labels:    cluster.Spec.JobLabels,
		},

		Spec: batchv1.JobSpec{

			BackoffLimit:          &backoffLimit,
			Completions:           &cluster.Spec.Size,
			Parallelism:           &cluster.Spec.Size,
			CompletionMode:        &completionMode,
			ActiveDeadlineSeconds: &cluster.Spec.DeadlineSeconds,

			// Note there is parameter to limit runtime
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        cluster.Name,
					Namespace:   cluster.Namespace,
					Labels:      podLabels,
					Annotations: cluster.Spec.Pod.Annotations,
				},
				Spec: corev1.PodSpec{
					// matches the service
					Subdomain:             cluster.Spec.Network.HeadlessName,
					ShareProcessNamespace: &cluster.Spec.ShareProcessNamespace,
					SetHostnameAsFQDN:     &setAsFQDN,
					Volumes:               getVolumes(cluster),
					RestartPolicy:         corev1.RestartPolicyOnFailure,
					ImagePullSecrets:      getImagePullSecrets(cluster),
					ServiceAccountName:    cluster.Spec.Pod.ServiceAccountName,
					NodeSelector:          cluster.Spec.Pod.NodeSelector,
				}},
		},
	}

	// Get resources for the pod
	resources, err := r.getPodResources(cluster)
	r.log.Info("🌀 MiniCluster", "Pod.Resources", resources)
	if err != nil {
		r.log.Info("🌀 MiniCluster", "Pod.Resources", resources)
		return job, err
	}
	job.Spec.Template.Spec.Overhead = resources

	// Get volume mounts specific to operator, add on container specific ones
	mounts := getVolumeMounts(cluster)

	// Prepare listing of containers for the MiniCluster
	containers, err := r.getContainers(
		cluster.Spec.Containers,
		cluster.Name,
		cluster.Spec.FluxRestful.Port,
		mounts,
	)
	job.Spec.Template.Spec.Containers = containers
	ctrl.SetControllerReference(cluster, job, r.Scheme)
	return job, err
}

// getImagePullSecrets returns a list of secret object references for each container.
func getImagePullSecrets(cluster *api.MiniCluster) []corev1.LocalObjectReference {
	pullSecrets := []corev1.LocalObjectReference{}
	for _, container := range cluster.Spec.Containers {
		if container.ImagePullSecret != "" {
			newSecret := corev1.LocalObjectReference{
				Name: container.ImagePullSecret,
			}
			pullSecrets = append(pullSecrets, newSecret)
		}
	}
	return pullSecrets
}
