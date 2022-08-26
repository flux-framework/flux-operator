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
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	logctrl "sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// getStatefulSet gets the existing statefulset, if it's done
func (r *FluxSetupReconciler) getStatefulSet(ctx context.Context, instance *api.FluxSetup, containerImage string) (*appsv1.StatefulSet, ctrl.Result, error) {

	log := logctrl.FromContext(ctx).WithValues("FluxSetup", instance.Namespace)
	existing := &appsv1.StatefulSet{}
	err := r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, existing)
	if err != nil {

		// Case 1: not found yet, check if deployment needs deletion
		if errors.IsNotFound(err) {
			dep := r.createStatefulSet(instance, containerImage)
			log.Info("✨ Creating a new StatefulSet ✨", "Namespace", dep.Namespace, "Name", dep.Name)
			err = r.Create(ctx, dep)
			if err != nil {
				log.Error(err, "❌ Failed to create new StatefulSet", "Namespace", dep.Namespace, "Name", dep.Name)
				return existing, ctrl.Result{}, err
			}
			// Deployment created successfully - return and requeue
			return existing, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get StatefulSet")
			return existing, ctrl.Result{}, err
		}
	} else {
		log.Info("🎉 Found existing StatefulSet 🎉", "Namespace", existing.Namespace, "Name", existing.Name, "Image", existing.Spec.Template.Spec.Containers[0].Image)
	}

	// Debugging to write yaml to yaml directory at root
	saveDebugYaml(existing, "stateful-set.yaml")
	return existing, ctrl.Result{}, err
}

// createStatefulSet creates the stateful set
func (r *FluxSetupReconciler) createStatefulSet(instance *api.FluxSetup, containerImage string) *appsv1.StatefulSet {
	labels := setupLabels(instance, "flux-workers")
	set := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{},
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
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           containerImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            instance.Name,
						VolumeMounts:    getVolumeMounts(),
					}},
					Volumes: getVolumes(),
				},
			},
		},
		Status: appsv1.StatefulSetStatus{},
	}
	ctrl.SetControllerReference(instance, set, r.Scheme)
	return set
}
