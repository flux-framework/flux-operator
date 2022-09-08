/*
Copyright 2022 Lawrence Livermore National Security, LLC
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

	api "flux-framework/flux-operator/api/v1alpha1"
)

// newMiniCluster is used to create the MiniCluster Job
func (r *MiniClusterReconciler) newMiniClusterJob(cluster *api.MiniCluster) *batchv1.Job {

	// Question from V: what are labels for?
	//	labels := labels(cluster, "flux-rank0")

	// We need to create the number of containers (and names) that the user requests
	// Before the stateful set was doing this for us, but for a batch job it's manaul
	containers := r.getMiniClusterContainers(cluster)

	// Number of retries before marking as failed
	backoffLimit := int32(100)
	completionMode := batchv1.IndexedCompletion

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{

			// In the example from Dan this was "indexed-job"
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
		},

		// https://github.com/kubernetes/api/blob/2f9e58849198f8675bc0928c69acf7e50af77551/batch/v1/types.go#L205
		Spec: batchv1.JobSpec{

			//			Selector:       &metav1.LabelSelector{MatchLabels: labels},
			BackoffLimit: &backoffLimit,
			Completions:  &cluster.Spec.Size,
			Parallelism:  &cluster.Spec.Size,

			// This would set a limit on the amount of time allowed to run
			// ActiveDeadlineSeconds: ...
			Template: corev1.PodTemplateSpec{ObjectMeta: metav1.ObjectMeta{Name: cluster.Name, Namespace: cluster.Namespace}, Spec: corev1.PodSpec{
				Volumes:       getVolumes(),
				Containers:    containers,
				RestartPolicy: corev1.RestartPolicyOnFailure,

				// Note that this spec also has variables for Host networking and DNS
				// These might be important (eventually)
				// Look at expanding this spec again to re-evaluate what we need
			}},
			// TODO I think we eventually can try CompletionMode Indexed here
			CompletionMode: &completionMode,
		},
	}
	ctrl.SetControllerReference(cluster, job, r.Scheme)
	return job
}

func (r *MiniClusterReconciler) getMiniClusterContainers(cluster *api.MiniCluster) []corev1.Container {

	// Allow the user to dictate pulling
	pullPolicy := corev1.PullIfNotPresent
	if (*cluster).Spec.PullAlways {
		pullPolicy = corev1.PullAlways
	}
	// Create the initial "driver" container to start flux
	containers := []corev1.Container{
		{
			// Call this the driver container, number 0
			Name:            cluster.Name,
			Image:           (*cluster).Spec.Image,
			ImagePullPolicy: pullPolicy,
			// Re-add config when we can reliably write it
			// Command:         []string{"flux", "start", "-o", "--config-path=/etc/flux/", (*cluster).Spec.Command},
			Command:      []string{"flux", "start", "-o", (*cluster).Spec.Command},
			WorkingDir:   (*cluster).Spec.WorkingDir,
			VolumeMounts: getVolumeMounts(),
		},
	}
	return containers
}
