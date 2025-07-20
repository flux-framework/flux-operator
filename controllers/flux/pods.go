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

	api "github.com/flux-framework/flux-operator/api/v1alpha2"
)

const podLabelAppName = "app.kubernetes.io/name"

// Get labels for any pod in the cluster
func getPodLabels(cluster *api.MiniCluster) map[string]string {
	if cluster.Spec.Pod.Labels == nil {
		cluster.Spec.Pod.Labels = map[string]string{}
	}
	podLabels := cluster.Spec.Pod.Labels
	podLabels["namespace"] = cluster.Namespace

	// Identify the pod as being created by the Flux Operator
	podLabels["app"] = "flux-operator"
	podLabels[podLabelAppName] = cluster.Name
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

			err = r.New(ctx, pod)
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
	err := r.Get(
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
	podLabels := getPodLabels(cluster)
	podServiceName := cluster.Name + "-services"

	// service selector?
	podLabels["job-name"] = cluster.Name

	// Services can have existing volumes
	existingVolumes := getExistingVolumes(cluster.ExistingServiceVolumes())

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
			Subdomain:                    cluster.Spec.Network.HeadlessName,
			Hostname:                     podServiceName,
			SetHostnameAsFQDN:            &setAsFQDN,
			Volumes:                      existingVolumes,
			ImagePullSecrets:             getImagePullSecrets(cluster),
			RestartPolicy:                corev1.RestartPolicyOnFailure,
			ServiceAccountName:           cluster.Spec.Pod.ServiceAccountName,
			AutomountServiceAccountToken: &cluster.Spec.Pod.AutomountServiceAccountToken,
			NodeSelector:                 cluster.Spec.Pod.NodeSelector,
		},
	}

	// Custom DNS Policy
	if cluster.Spec.Pod.DNSPolicy != "" {
		pod.Spec.DNSPolicy = corev1.DNSPolicy(cluster.Spec.Pod.DNSPolicy)
	}

	// Custom restart policy
	if cluster.Spec.Pod.RestartPolicy != "" {
		pod.Spec.RestartPolicy = corev1.RestartPolicy(cluster.Spec.Pod.RestartPolicy)
	}

	// Only add runClassName if defined
	if cluster.Spec.Pod.RuntimeClassName != "" {
		pod.Spec.RuntimeClassName = &cluster.Spec.Pod.RuntimeClassName
	}

	// Assemble existing volume mounts - they are added with getContainers
	mounts := []corev1.VolumeMount{}

	// Get containers for the service pods
	// true boolean indicates this is a service container - do not generate
	// a custom entrypoint, etc.
	containers, err := getContainers(
		cluster.Spec.Services,
		podServiceName,
		mounts,
		true,
	)

	if err != nil {
		return pod, err
	}
	pod.Spec.Containers = containers
	ctrl.SetControllerReference(cluster, pod, r.Scheme)
	return pod, err
}
