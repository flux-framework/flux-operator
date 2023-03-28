/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// Get labels for any pod in the cluster
func (r *MiniClusterReconciler) getPodLabels(cluster *api.MiniCluster) map[string]string {
	podLabels := cluster.Spec.Pod.Labels
	podLabels["namespace"] = cluster.Namespace
	podLabels["app.kubernetes.io/name"] = cluster.Name
	return podLabels
}

// ensure service containers are running, currently in one pod
func (r *MiniClusterReconciler) ensureServicePod(
	ctx context.Context,
	cluster *api.MiniCluster,
) (*corev1.Pod, ctrl.Result, error) {

	// Look for an existing service container
	existing, err := r.getExistingPod(ctx, cluster)

	// Create a new job if it does not exist
	if err != nil {

		if errors.IsNotFound(err) {
			pod, err := r.newServicePod(cluster)
			if err != nil {
				return existing, ctrl.Result{}, err
			}
			r.log.Info(
				"✨ Creating a new MiniCluster Service Pod ✨",
				"Namespace:", pod.Namespace,
				"Name:", pod.Name,
			)

			err = r.Client.Create(ctx, pod)
			if err != nil {
				r.log.Error(
					err,
					"Failed to create new MiniCluster Service Pod",
					"Namespace:", pod.Namespace,
					"Name:", pod.Name,
				)
				return pod, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return pod, ctrl.Result{Requeue: true}, nil

		} else if err != nil {
			r.log.Error(err, "Failed to get MiniCluster Service Pod")
			return existing, ctrl.Result{}, err
		}

	} else {
		r.log.Info(
			"🎉 Found existing MiniCluster Service Pod 🎉",
			"Namespace:", existing.Namespace,
			"Name:", existing.Name,
		)
	}
	return existing, ctrl.Result{}, err
}

// getExistingPod gets an existing pod service
func (r *MiniClusterReconciler) getExistingPod(
	ctx context.Context,
	cluster *api.MiniCluster,
) (*corev1.Pod, error) {

	existing := &corev1.Pod{}
	err := r.Client.Get(
		ctx,
		types.NamespacedName{
			Name:      cluster.Name + "-services",
			Namespace: cluster.Namespace,
		},
		existing,
	)
	return existing, err
}

// newMiniCluster is used to create the MiniCluster Job
func (r *MiniClusterReconciler) newServicePod(
	cluster *api.MiniCluster,
) (*corev1.Pod, error) {

	setAsFQDN := false
	podLabels := r.getPodLabels(cluster)
	podServiceName := cluster.Name + "-services"

	// service selector?
	podLabels["job-name"] = cluster.Name

	// This is an indexed-job
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        podServiceName,
			Namespace:   cluster.Namespace,
			Labels:      podLabels,
			Annotations: cluster.Spec.Pod.Annotations,
		},
		Spec: corev1.PodSpec{
			// This is the headless service name
			Subdomain:          restfulServiceName,
			Hostname:           podServiceName,
			SetHostnameAsFQDN:  &setAsFQDN,
			RestartPolicy:      corev1.RestartPolicyOnFailure,
			ImagePullSecrets:   getImagePullSecrets(cluster),
			ServiceAccountName: cluster.Spec.Pod.ServiceAccountName,
			NodeSelector:       cluster.Spec.Pod.NodeSelector,
		},
	}
	mounts := []corev1.VolumeMount{}
	containers, err := r.getContainers(cluster.Spec.Services, podServiceName, mounts)
	if err != nil {
		return pod, err
	}
	pod.Spec.Containers = containers
	ctrl.SetControllerReference(cluster, pod, r.Scheme)
	return pod, err
}

/*
// exposeService will expose services - one for the port 5000 forward, and the other for job networking (headless)
func (r *MiniClusterReconciler) exposeServiceContainerService(
	ctx context.Context,
	cluster *api.MiniCluster,
	container *api.MiniClusterContainer,
) (ctrl.Result, error) {

	// This service is for the restful API
	existing := &corev1.Service{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: container.Name, Namespace: cluster.Namespace}, existing)
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = r.createServicePodService(ctx, cluster, container)
		}
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, err
}

// createMiniClusterService creates the service for the minicluster
func (r *MiniClusterReconciler) createServicePodService(
	ctx context.Context,
	cluster *api.MiniCluster,
	container *api.MiniClusterContainer,
) (*corev1.Service, error) {

	r.log.Info("Creating service with: ", container.Name, cluster.Namespace)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: container.Name, Namespace: cluster.Namespace},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector: map[string]string{
				"job-name": cluster.Name,
			},
		},
	}
	err := ctrl.SetControllerReference(cluster, service, r.Scheme)
	if err != nil {
		r.log.Error(err, "🔴 Create service", "Service", restfulServiceName)
		return service, err
	}
	err = r.Client.Create(ctx, service)
	if err != nil {
		r.log.Error(err, "🔴 Create service", "Service", restfulServiceName)
	}
	return service, err
}*/
