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

	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	api "flux-framework/flux-operator/api/v1alpha1"
)

var (
	serviceName = "flux-restful-service"
	servicePort = 5000
)

// exposeService will expose a service
func (r *MiniClusterReconciler) exposeService(ctx context.Context, cluster *api.MiniCluster) (ctrl.Result, error) {

	existing := &corev1.Service{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: serviceName, Namespace: cluster.Namespace}, existing)
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = r.createMiniClusterService(ctx, cluster)
		}
		return ctrl.Result{}, err
	}

	ingress := &networkv1.Ingress{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: serviceName, Namespace: cluster.Namespace}, ingress)
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.createMiniClusterIngress(ctx, cluster, existing)
		}
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, err
}

// createMiniClusterService creates the service for the minicluster
func (r *MiniClusterReconciler) createMiniClusterService(ctx context.Context, cluster *api.MiniCluster) (*corev1.Service, error) {

	r.log.Info("Creating service with: ", serviceName, cluster.Namespace)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: serviceName, Namespace: cluster.Namespace},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:       serviceName,
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(servicePort),
					TargetPort: intstr.FromInt(servicePort),
				},
			},
			ExternalIPs: []string{"192.168.0.194"},
		},
	}
	err := ctrl.SetControllerReference(cluster, service, r.Scheme)
	if err != nil {
		r.log.Error(err, "🔴 Create service", "Service", serviceName)
		return service, err
	}
	err = r.Client.Create(ctx, service)
	if err != nil {
		r.log.Error(err, "🔴 Create service", "Service", serviceName)
	}
	return service, err
}

// createMiniClusterIngress exposes the service for the minicluster
func (r *MiniClusterReconciler) createMiniClusterIngress(ctx context.Context, cluster *api.MiniCluster, service *corev1.Service) error {

	pathType := networkv1.PathTypePrefix
	ingressBackend := networkv1.IngressBackend{
		Service: &networkv1.IngressServiceBackend{
			Name: service.Name,
			Port: networkv1.ServiceBackendPort{
				Number: service.Spec.Ports[0].NodePort,
			},
		},
	}
	ingress := &networkv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
		},
		Spec: networkv1.IngressSpec{
			DefaultBackend: &ingressBackend,
			Rules: []networkv1.IngressRule{
				{
					IngressRuleValue: networkv1.IngressRuleValue{
						HTTP: &networkv1.HTTPIngressRuleValue{
							Paths: []networkv1.HTTPIngressPath{
								{
									PathType: &pathType,
									Backend:  ingressBackend,
									Path:     "/",
								},
							},
						},
					},
				},
			},
		},
	}
	err := ctrl.SetControllerReference(cluster, ingress, r.Scheme)
	if err != nil {
		r.log.Error(err, "🔴 Create ingress", "Service", serviceName)
		return err
	}
	err = r.Client.Create(ctx, ingress)
	if err != nil {
		r.log.Error(err, "🔴 Create ingress", "Service", serviceName)
		return err
	}
	return nil
}
