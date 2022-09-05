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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	api "flux-framework/flux-operator/api/v1alpha1"
	"flux-framework/flux-operator/pkg/certs"
)

// generateCurveCert makes a new Secret if it doesn't exist
func (r *FluxJobReconciler) getCurveCert(ctx context.Context, fluxjob *api.FluxJob) (*corev1.Secret, ctrl.Result, error) {

	existing := &corev1.Secret{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: "secret-tls", Namespace: fluxjob.Namespace}, existing)
	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {
			dep := r.createCurveSecret(fluxjob)
			r.log.Info("✨ Creating a new Secret ✨", "Namespace", dep.Namespace, "Name", dep.Name, "Data", (*dep).Data)
			err = r.Client.Create(ctx, dep)
			if err != nil {
				r.log.Error(err, "❌ Failed to create new Curve Secret", "Namespace", dep.Namespace, "Name", (*dep).Name)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return existing, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			r.log.Error(err, "Failed to get Broker ConfigMap")
			return existing, ctrl.Result{}, err
		}
	} else {
		r.log.Info("🎉 Found existing Secret 🎉", "Namespace", existing.Namespace, "Name", existing.Name, "Data", (*existing).Data)
	}
	saveDebugYaml(existing, "secret.yaml")
	return existing, ctrl.Result{}, err
}

// createCurveSecret creates the secret
func (r *FluxJobReconciler) createCurveSecret(fluxjob *api.FluxJob) *corev1.Secret {

	// TODO do we need hosts here?
	c := certs.NewCertificate([]string{}, false)
	c.Generate()

	cert := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret-tls",
			Namespace: fluxjob.Namespace,
		},
		Data: map[string][]byte{
			"tls.key": []byte(c.Public),
			"tls.crt": []byte(c.Private),
		},
	}
	ctrl.SetControllerReference(fluxjob, cert, r.Scheme)
	return cert
}
