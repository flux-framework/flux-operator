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

// getBrokerConfig gets the existing broker config, if it's done
func (r *FluxSetupReconciler) getBrokerConfig(ctx context.Context, instance *api.FluxSetup) (*corev1.ConfigMap, ctrl.Result, error) {

	log := logctrl.FromContext(ctx).WithValues("FluxSetup", instance.Namespace)
	existing := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: "flux-config", Namespace: instance.Namespace}, existing)
	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {
			dep := r.createBrokerConfig(instance)
			log.Info("✨ Creating a new Broker ConfigMap ✨", "Namespace", dep.Namespace, "Name", dep.Name, "Data", (*dep).Data)
			err = r.Create(ctx, dep)
			if err != nil {
				log.Error(err, "❌ Failed to create new Broker ConfigMap", "Namespace", dep.Namespace, "Name", (*dep).Name)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return existing, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Broker ConfigMap")
			return existing, ctrl.Result{}, err
		}
	} else {
		log.Info("🎉 Found existing Broker ConfigMap 🎉", "Namespace", existing.Namespace, "Name", existing.Name)
	}
	saveDebugYaml(existing, "broker.yaml")
	return existing, ctrl.Result{}, err
}

// getBrokerConfig gets the existing broker config, if it's done
func (r *FluxSetupReconciler) getEtcHostsConfig(ctx context.Context, instance *api.FluxSetup) (*corev1.ConfigMap, ctrl.Result, error) {

	log := logctrl.FromContext(ctx).WithValues("FluxSetup", instance.Namespace)
	existing := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: "etc-hosts", Namespace: instance.Namespace}, existing)
	if err != nil {
		if errors.IsNotFound(err) {
			dep := r.createEtcHostsConfig(instance)
			log.Info("✨ Creating a new etc-hosts ConfigMap ✨", "Namespace", dep.Namespace, "Name", dep.Name, "Data", (*dep).Data)
			err = r.Create(ctx, dep)
			if err != nil {
				log.Error(err, "❌ Failed to create new etc-hosts ConfigMap", "Namespace", dep.Namespace, "Name", dep.Name)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return existing, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Broker ConfigMap")
			return existing, ctrl.Result{}, err
		}
	} else {
		log.Info("🎉 Found existing etc-hosts ConfigMap 🎉", "Namespace", existing.Namespace, "Name", existing.Name)
	}
	saveDebugYaml(existing, "etc-hosts-config.yaml")
	return existing, ctrl.Result{}, err
}

// createBrokerConfig creates the stateful set
func (r *FluxSetupReconciler) createBrokerConfig(instance *api.FluxSetup) *corev1.ConfigMap {
	broker := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flux-config",
			Namespace: instance.Namespace,
		},
		Data: map[string]string{
			"hostfile": (*instance).Spec.Broker.Hostfile,
		},
	}
	fmt.Println(broker.Data)
	ctrl.SetControllerReference(instance, broker, r.Scheme)
	return broker
}

// createEtcHostsConfig creates the stateful set
func (r *FluxSetupReconciler) createEtcHostsConfig(instance *api.FluxSetup) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "etc-hosts",
			Namespace: instance.Namespace,
		},
		Data: map[string]string{
			"hostfile": (*instance).Spec.EtcHosts.Hostfile,
		},
	}
	fmt.Println(cm.Data)
	ctrl.SetControllerReference(instance, cm, r.Scheme)
	return cm
}
