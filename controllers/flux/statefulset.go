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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// createDeployment creates the stateful set
func (r *FluxSetupReconciler) createDeployment(instance *api.FluxSetup, containerImage string) *appsv1.StatefulSet {
	labels := setupLabels(instance, "flux-workers")
	set := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &instance.Spec.Size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			PodManagementPolicy: appsv1.ParallelPodManagement,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      instance.Name,
					Namespace: instance.Namespace,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{

						// This comes from the Flux custom resource (from the user)
						Image:           containerImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            instance.Name,
						VolumeMounts:    getVolumeMounts(),
					}},
					Volumes: getVolumes(),
				},
			},
		},
	}
	ctrl.SetControllerReference(instance, set, r.Scheme)
	return set
}
