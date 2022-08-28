/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	api "flux-framework/flux-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

/*
apiVersion: v1
kind: Service
metadata:
  name: flux-workers
spec:
  clusterIP: None
  selector:
    app: flux-workers*/
func (r *FluxSetupReconciler) createService(instance *api.FluxSetup) *corev1.Service {

	labels := setupLabels(instance, "flux-workers")

	// We shouldn't need this, as the port comes from the manifest
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
		},
	}
	ctrl.SetControllerReference(instance, service, r.Scheme)
	return service
}
