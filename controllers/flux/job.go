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
	createJobDNS := true

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{

			// In the example from Dan this was "indexed-job"
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
		},

		// https://github.com/kubernetes/api/blob/2f9e58849198f8675bc0928c69acf7e50af77551/batch/v1/types.go#L205
		Spec: batchv1.JobSpec{

			//			Selector:       &metav1.LabelSelector{MatchLabels: labels},
			BackoffLimit:   &backoffLimit,
			Completions:    &cluster.Spec.Size,
			Parallelism:    &cluster.Spec.Size,
			CompletionMode: &completionMode,

			// This would set a limit on the amount of time allowed to run
			// ActiveDeadlineSeconds: ...
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cluster.Name,
					Namespace: cluster.Namespace,
				},
				Spec: corev1.PodSpec{
					Volumes:    getVolumes(cluster),
					Containers: containers,

					// The init containers use flux-keygen to create certs
					InitContainers: r.getMiniClusterInitContainer(cluster),
					RestartPolicy:  corev1.RestartPolicyOnFailure,

					// Create a Service-style DNS entry like:
					// pod-instance-1.default-subdomain.my-namespace.svc.cluster-domain.example
					SetHostnameAsFQDN: &createJobDNS,
				}},
		},
	}
	ctrl.SetControllerReference(cluster, job, r.Scheme)
	return job
}

// The init containers create the curve.cert if it does not exist using flux
// We do this because libsodium / zmq libraries are already in the container
func (r *MiniClusterReconciler) getMiniClusterInitContainer(cluster *api.MiniCluster) []corev1.Container {

	// Allow the user to dictate pulling
	pullPolicy := corev1.PullIfNotPresent
	if (*cluster).Spec.PullAlways {
		pullPolicy = corev1.PullAlways
	}

	containers := []corev1.Container{
		{
			// Call this the driver container, number 0
			Name:            cluster.Name + "-init",
			Image:           (*cluster).Spec.Image,
			ImagePullPolicy: pullPolicy,

			// Don't provide the name here - it will get from the host
			Command:      []string{"flux", "keygen", "/mnt/curve/curve.cert"},
			VolumeMounts: getVolumeMounts(cluster),
		},
	}
	return containers
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

			// config is a directory with any number of toml files to be used, we use a brokers.toml
			Command:      []string{"flux", "start", "-o", "--config-path=/etc/flux/config", (*cluster).Spec.Command},
			WorkingDir:   (*cluster).Spec.WorkingDir,
			VolumeMounts: getVolumeMounts(cluster),
		},
	}
	return containers
}
