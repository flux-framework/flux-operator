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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	logctrl "sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// getMiniCluster gets the MiniCluster CRD
func (r *FluxSetupReconciler) getMiniCluster(ctx context.Context, instance *api.FluxSetup) (*api.MiniCluster, ctrl.Result, error) {

	log := logctrl.FromContext(ctx).WithValues("FluxSetup", instance.Namespace)
	cluster := &api.MiniCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}
	err := r.Get(ctx, types.NamespacedName{Name: cluster.Name, Namespace: cluster.Namespace}, cluster)
	if err != nil {

		// Case 1: not found yet, check if deployment needs deletion
		if errors.IsNotFound(err) {
			ctrl.SetControllerReference(instance, cluster, r.Scheme)
			log.Info("✨ Creating a new MiniCluster ✨", "Namespace", cluster.Namespace, "Name", cluster.Name)
			err = r.Create(ctx, cluster)
			if err != nil {
				log.Error(err, "❌ Failed to create new MiniCluster", "Namespace", cluster.Namespace, "Name", cluster.Name)
				return cluster, ctrl.Result{}, err
			}
			return cluster, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get MiniCluster")
			return cluster, ctrl.Result{}, err
		}
	} else {
		log.Info("🎉 Found existing MiniCluster 🎉", "Namespace", cluster.Namespace, "Name", cluster.Name)
	}

	// Debugging to write yaml to yaml directory at root
	saveDebugYaml(cluster, "mini-cluster.yaml")
	return cluster, ctrl.Result{}, err
}
