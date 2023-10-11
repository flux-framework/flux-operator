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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"
	"github.com/flux-framework/flux-operator/pkg/flux"
)

// getConfigMap gets the entrypoint config map
func (r *MiniClusterReconciler) getConfigMap(
	ctx context.Context,
	cluster *api.MiniCluster,
	configName string,
) (*corev1.ConfigMap, ctrl.Result, error) {

	// Look for the config map by name
	r.log.Info("👀️ Looking for ConfigMap 👀️", "Type", configName)
	existing := &corev1.ConfigMap{}
	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      configName,
			Namespace: cluster.Namespace,
		},
		existing,
	)

	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {

			// Generate all entrypoints for the config map
			data, err := flux.GenerateEntrypoints(cluster)
			if err != nil {
				return existing, ctrl.Result{}, err
			}

			// Finally create the config map
			dep := r.createConfigMap(cluster, configName, data)
			r.log.Info(
				"✨ Creating MiniCluster ConfigMap ✨",
				"Type", configName,
				"Namespace", dep.Namespace,
				"Name", dep.Name,
			)
			err = r.New(ctx, dep)
			if err != nil {
				r.log.Error(
					err, "❌ Failed to create MiniCluster ConfigMap",
					"Type", configName,
					"Namespace", dep.Namespace,
					"Name", (*dep).Name,
				)
				return existing, ctrl.Result{}, err
			}

			// Successful - return and requeue
			return dep, ctrl.Result{Requeue: true}, nil

		} else if err != nil {
			r.log.Error(
				err, "Failed to get MiniCluster ConfigMap",
				"Type", configName,
			)
			return existing, ctrl.Result{}, err
		}

	} else {
		r.log.Info(
			"🎉 Found existing MiniCluster ConfigMap",
			"Type", configName,
			"Namespace", existing.Namespace,
			"Name", existing.Name,
		)
	}
	return existing, ctrl.Result{}, err
}

// createConfigMap generates a config map with some kind of data
func (r *MiniClusterReconciler) createConfigMap(
	cluster *api.MiniCluster,
	configName string,
	data map[string]string,
) *corev1.ConfigMap {

	// Create the config map with respective data!
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configName,
			Namespace: cluster.Namespace,
		},
		Data: data,
	}

	// Show in the logs
	fmt.Println(cm.Data)
	ctrl.SetControllerReference(cluster, cm, r.Scheme)
	return cm
}
