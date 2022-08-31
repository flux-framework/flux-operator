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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	logctrl "sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	api "flux-framework/flux-operator/api/v1alpha1"
)

var BrokerConfigTemplate = `
[bootstrap]
curve_cert = "/mnt/curve/curve.cert"
default_port = 8050
default_bind = "tcp://eth0:%%p"
default_connect = "tcp://%%h:%%p"
hosts = [
	{ host="%s-%s"},
]
`

// getHostfileConfig gets an existing configmap, if it's done
func (r *FluxSetupReconciler) getHostfileConfig(ctx context.Context, instance *api.FluxSetup, configName string, hostfile string) (*corev1.ConfigMap, ctrl.Result, error) {

	log := logctrl.FromContext(ctx).WithValues("FluxSetup", instance.Namespace)
	existing := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: configName, Namespace: instance.Namespace}, existing)
	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {
			// check if its broker.toml TODO : Convert all 3 configMaps to use template stirngs
			if configName == "flux-config" {
				hostfile = generateFluxConfig(instance.Name, instance.Spec.Size)
			}
			dep := r.createHostfileConfig(instance, configName, hostfile)
			log.Info("✨ Creating a new ConfigMap ✨", "Type", configName, "Namespace", dep.Namespace, "Name", dep.Name, "Data", (*dep).Data)
			err = r.Create(ctx, dep)
			if err != nil {
				log.Error(err, "❌ Failed to create new ConfigMap", "Type", configName, "Namespace", dep.Namespace, "Name", (*dep).Name)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return existing, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get ConfigMap")
			return existing, ctrl.Result{}, err
		}
	} else {
		log.Info("🎉 Found existing ConfigMap 🎉", "Type", configName, "Namespace", existing.Namespace, "Name", existing.Name, "Data", (*existing).Data)
	}
	saveDebugYaml(existing, configName+".yaml")
	return existing, ctrl.Result{}, err
}

func generateFluxConfig(name string, size int32) string {
	var hosts string
	if size == 1 {
		hosts = "0"
	} else {
		hosts = fmt.Sprintf("[0-%d]", size-1)
	}
	fluxConfig := fmt.Sprintf(BrokerConfigTemplate, name, hosts)

	return fluxConfig
}

// createBrokerConfig creates the stateful set
func (r *FluxSetupReconciler) createHostfileConfig(instance *api.FluxSetup, configName string, hostfile string) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configName,
			Namespace: instance.Namespace,
		},
		Data: map[string]string{
			"hostfile": hostfile,
		},
	}
	fmt.Println(cm.Data)
	ctrl.SetControllerReference(instance, cm, r.Scheme)
	return cm
}
