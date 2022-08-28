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

// createJob is used by Flux to create the Job spec
func (r *FluxSetupReconciler) createJob(instance *api.Flux, containerImage string) *batchv1.Job {
	labels := labels(instance, "flux-rank0")
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},

		// https://github.com/kubernetes/api/blob/2f9e58849198f8675bc0928c69acf7e50af77551/batch/v1/types.go#L205
		Spec: batchv1.JobSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      instance.Name,
					Namespace: instance.Namespace,
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
					Containers: []corev1.Container{{
						Name:            "driver",
						Image:           (*instance).Spec.Image,
						ImagePullPolicy: corev1.PullAlways,
						Command:         []string{"flux", "start", "-o", "--config-path=/etc/flux/", (*instance).Spec.Command},
						VolumeMounts:    getVolumeMounts(),
					}},
					Volumes: getVolumes(),
				},
			},
		},
	}
	ctrl.SetControllerReference(instance, job, r.Scheme)
	return job
}
