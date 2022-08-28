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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	logctrl "sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// getCurveCert gets (or creates) the curve secret
// TODO this looks like a way to check for needing to update?
// https://github.com/redhat-cop/cert-utils-operator/blob/dbf1df07a63460852a159943bb16650e139af6eb/controllers/route/route_controller.go#L291
func (r *FluxSetupReconciler) getCurveCert(ctx context.Context, instance *api.FluxSetup) (*corev1.Secret, ctrl.Result, error) {

	log := logctrl.FromContext(ctx).WithValues("FluxSetup", instance.Namespace)
	existing := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: "secret-tls", Namespace: instance.Namespace}, existing)
	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {
			dep := r.createCurveSecret(instance)
			log.Info("✨ Creating a new Curve Secret ✨", "Namespace", dep.Namespace, "Name", dep.Name, "Data", (*dep).Data)
			err = r.Create(ctx, dep)
			if err != nil {
				log.Error(err, "❌ Failed to create new Curve Secret", "Namespace", dep.Namespace, "Name", (*dep).Name)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return existing, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Broker ConfigMap")
			return existing, ctrl.Result{}, err
		}
	} else {
		log.Info("🎉 Found existing Broker ConfigMap 🎉", "Namespace", existing.Namespace, "Name", existing.Name, "Data", (*existing).Data)
	}
	saveDebugYaml(existing, "secret.yaml")
	return existing, ctrl.Result{}, err
}

// createCurveSecret creates the secret
// I think we need to do https://kubernetes.io/docs/tasks/configmap-secret/managing-secret-using-kustomize/
// This is useful https://github.com/jetstack/kustomize-cert-manager-demo
// And https://www.jetstack.io/blog/kustomize-cert-manager/
/*
   apiVersion: v1
   kind: Secret
   metadata:
     name: secret-tls
   type: kubernetes.io/tls
   data:
     # the data is abbreviated in this example
     tls.crt: |
         MIIC2DCCAcCgAwIBAgIBATANBgkqh ...
     tls.key: |
         MIIEpgIBAAKCAQEA7yn3bRHQ5FHMQ ...
*/
func (r *FluxSetupReconciler) createCurveSecret(instance *api.FluxSetup) *corev1.Secret {
	cert := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret-tls",
			Namespace: instance.Namespace,
		},
	}
	ctrl.SetControllerReference(instance, cert, r.Scheme)
	return cert
}
