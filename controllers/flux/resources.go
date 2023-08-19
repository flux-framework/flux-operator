/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"fmt"

	api "github.com/flux-framework/flux-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// getResourceGroup can return a ResourceList for either requests or limits
func getResourceGroup(items api.ContainerResource) (corev1.ResourceList, error) {

	list := corev1.ResourceList{}
	for key, unknownValue := range items {
		if unknownValue.Type == intstr.Int {

			value := unknownValue.IntVal
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
func getContainerResources(container *api.MiniClusterContainer) (corev1.ResourceRequirements, error) {

	// memory int, setCPURequest, setCPULimit, setGPULimit int64
	resources := corev1.ResourceRequirements{}

	// Limits
	limits, err := getResourceGroup(container.Resources.Limits)
	if err != nil {
		return resources, err
	}
	resources.Limits = limits

	// Requests
	requests, err := getResourceGroup(container.Resources.Requests)
	if err != nil {
		return resources, err
	}
	resources.Requests = requests
	return resources, nil
}

// getPodResources determines if any resources are requested via the spec
func getPodResources(cluster *api.MiniCluster) (corev1.ResourceList, error) {

	// memory int, setCPURequest, setCPULimit, setGPULimit int64
	resources, err := getResourceGroup(cluster.Spec.Pod.Resources)
	if err != nil {
		return resources, err
	}
	return resources, nil
}
