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

					// Job is used as selector for service
					Labels: map[string]string{"name": cluster.Name, "namespace": cluster.Namespace, "job": cluster.Name},
				},
				Spec: corev1.PodSpec{
					// Following example in:
					// https://github.com/alculquicondor/enhancements/blob/master/keps/sig-apps/2214-indexed-job/README.md
					// When this is set, we see:
					// 172.17.0.7      flux-sample-0.flux-sample.flux-operator.svc.cluster.local       flux-sample-0
					// The FQDN setting doesn't seem to work
					Subdomain:  cluster.Name,
					Volumes:    getVolumes(cluster),
					Containers: containers,

					// The init containers use flux keygen to create certs
					// There are multiple now but eventually we need just one
					InitContainers: r.getMiniClusterInitContainer(cluster),
					RestartPolicy:  corev1.RestartPolicyOnFailure,
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
			Name:  cluster.Name + "-init",
			Image: (*cluster).Spec.Image,

			// TODO this should be done in a script (via CM volume) that checks the index envar
			Command:         []string{"flux", "keygen", "/mnt/curve/curve.cert"},
			Resources:       corev1.ResourceRequirements{},
			VolumeMounts:    getVolumeMounts(cluster),
			ImagePullPolicy: pullPolicy,
			Stdin:           true,
			TTY:             true,

			// Just setting this for testing
			// It's added to the JOB_COMPLETION_INDEX variable
			Env: []corev1.EnvVar{
				{
					Name:  "FLUX_OPERATOR_CONTAINER_TYPE",
					Value: "INIT",
				},
			},
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

			// This is a wrapper that is going to wait for the generation of update_hosts.sh
			// Once it's there, we update /etc/hosts, and run the command to start flux.
			Command:      []string{"/bin/bash", "/flux_operator/wait.sh", (*cluster).Spec.Command},
			WorkingDir:   (*cluster).Spec.WorkingDir,
			VolumeMounts: getVolumeMounts(cluster),
			Stdin:        true,
			TTY:          true,
			Env: []corev1.EnvVar{
				{
					Name:  "FLUX_OPERATOR_CONTAINER_TYPE",
					Value: "WORKER",
				},
			},
		},
	}
	return containers
}
