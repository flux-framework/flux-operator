/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
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

var (
	restfulServiceName = "flux-service"
	servicePort        = 5000
)

// exposeService will expose services - one for the port 5000 forward, and the other for job networking (headless)
func (r *MiniClusterReconciler) exposeServices(ctx context.Context, cluster *api.MiniCluster) (ctrl.Result, error) {

	// This service is for the restful API
	existing := &corev1.Service{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: restfulServiceName, Namespace: cluster.Namespace}, existing)
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = r.createMiniClusterService(ctx, cluster)
		}
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, err
}

// createMiniClusterService creates the service for the minicluster
func (r *MiniClusterReconciler) createMiniClusterService(ctx context.Context, cluster *api.MiniCluster) (*corev1.Service, error) {

	r.log.Info("Creating service with: ", restfulServiceName, cluster.Namespace)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: restfulServiceName, Namespace: cluster.Namespace},
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
}
