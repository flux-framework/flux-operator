/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
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
	err := r.Client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, existing)
	if err != nil {

		// Case 1: not found yet, check if deployment needs deletion
		if errors.IsNotFound(err) {
			dep := r.createStatefulSet(instance, containerImage)
			log.Info("✨ Creating a new StatefulSet ✨", "Namespace", dep.Namespace, "Name", dep.Name)
			err = r.Client.Create(ctx, dep)
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
						Command:         []string{"/bin/bash"},
						Args:            []string{"-c", "sleep infinity"},
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
