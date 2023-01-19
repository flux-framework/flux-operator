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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// getResourceGroup can return a ResourceList for either requests or limits
func (r *MiniClusterReconciler) getResourceGroup(items api.ContainerResource) (corev1.ResourceList, error) {

	r.log.Info("🍅️ Resource", "items", items)
	list := corev1.ResourceList{}
	for key, unknownValue := range items {
		if unknownValue.Type == intstr.Int {

			value := unknownValue.IntVal
			r.log.Info("🍅️ ResourceKey", "Key", key, "Value", value)
			limit, err := resource.ParseQuantity(fmt.Sprintf("%d", value))
			if err != nil {
				return list, err
			}

			if key == "memory" {
				list[corev1.ResourceMemory] = limit
			} else if key == "cpu" {
				list[corev1.ResourceCPU] = limit
			} else {
				list[corev1.ResourceName(key)] = limit
			}
		} else if unknownValue.Type == intstr.String {

			value := unknownValue.StrVal
			r.log.Info("🍅️ ResourceKey", "Key", key, "Value", value)
			if key == "memory" {
				list[corev1.ResourceMemory] = resource.MustParse(value)
			} else if key == "cpu" {
				list[corev1.ResourceCPU] = resource.MustParse(value)
			} else {
				list[corev1.ResourceName(key)] = resource.MustParse(value)
			}
		}
	}
	return list, nil
}

// getContainerResources determines if any resources are requested via the spec
func (r *MiniClusterReconciler) getContainerResources(cluster *api.MiniCluster, container *api.MiniClusterContainer) (corev1.ResourceRequirements, error) {

	// memory int, setCPURequest, setCPULimit, setGPULimit int64
	resources := corev1.ResourceRequirements{}

	// Limits
	limits, err := r.getResourceGroup(container.Resources.Limits)
	if err != nil {
		r.log.Error(err, "🍅️ Resources for Container.Limits")
		return resources, err
	}
	resources.Limits = limits

	// Requests
	requests, err := r.getResourceGroup(container.Resources.Requests)
	if err != nil {
		r.log.Error(err, "🍅️ Resources for Container.Requests")
		return resources, err
	}
	resources.Requests = requests
	return resources, nil
}

// getPodResources determines if any resources are requested via the spec
func (r *MiniClusterReconciler) getPodResources(cluster *api.MiniCluster) (corev1.ResourceList, error) {

	// memory int, setCPURequest, setCPULimit, setGPULimit int64
	resources, err := r.getResourceGroup(cluster.Spec.Pod.Resources)
	if err != nil {
		r.log.Error(err, "🍅️ Resources for Pod.Resources")
		return resources, err
	}
	return resources, nil
}
