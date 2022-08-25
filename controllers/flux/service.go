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
