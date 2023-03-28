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

	corev1 "k8s.io/api/core/v1"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// getContainers gets containers for a MiniCluster job or external service
func (r *MiniClusterReconciler) getContainers(
	specs []api.MiniClusterContainer,
	defaultName string,
	mounts []corev1.VolumeMount,
) ([]corev1.Container, error) {

	// Create the containers for the pod
	containers := []corev1.Container{}

	for i, container := range specs {

		// Allow dictating pulling on the level of the container
		pullPolicy := corev1.PullIfNotPresent
		if container.PullAlways {
			pullPolicy = corev1.PullAlways
		}

		// Fluxrunner will use the namespace name
		containerName := container.Name
		command := []string{}

		// A Flux runner gets a custom wait.sh script for the container
		// And also needs to have a consistent name to the cert generator
		if container.RunFlux {

			// wait.sh path corresponds to container identifier
			waitScript := fmt.Sprintf("/flux_operator/wait-%d.sh", i)
			command = []string{"/bin/bash", waitScript, container.Command}
			containerName = defaultName
		}

		// Prepare lifescycle commands for the container
		lifecycle := r.createContainerLifecycle(container)

		for volumeName, volume := range container.Volumes {
			mount := corev1.VolumeMount{
				Name:      volumeName,
				MountPath: volume.Path,
				ReadOnly:  volume.ReadOnly,
			}
			mounts = append(mounts, mount)
		}

		// Add on existing volumes/claims
		for volumeName, volume := range container.ExistingVolumes {
			mount := corev1.VolumeMount{
				Name:      volumeName,
				MountPath: volume.Path,
				ReadOnly:  volume.ReadOnly,
			}
			mounts = append(mounts, mount)
		}

		r.log.Info("🌀 MiniCluster", "Container.Mounts", mounts)

		// Prepare container resources
		resources, err := r.getContainerResources(&container)
		r.log.Info("🌀 MiniCluster", "Container.Resources", resources)
		if err != nil {
			return containers, err
		}
		securityContext := corev1.SecurityContext{
			Privileged: &container.SecurityContext.Privileged,
		}
		newContainer := corev1.Container{

			// Call this the driver container, number 0
			Name:            containerName,
			Image:           container.Image,
			ImagePullPolicy: pullPolicy,
			WorkingDir:      container.WorkingDir,
			VolumeMounts:    mounts,
			Stdin:           true,
			TTY:             true,
			Lifecycle:       lifecycle,
			Resources:       resources,
			SecurityContext: &securityContext,
		}

		// Only add command if we actually have one
		if len(command) > 0 {
			newContainer.Command = command
		}

		ports := []corev1.ContainerPort{}
		envars := []corev1.EnvVar{}

		// If it's the FluxRunner, expose port 5000 for the service
		if container.RunFlux {
			newPort := corev1.ContainerPort{
				ContainerPort: int32(servicePort),
				Protocol:      "TCP",
			}
			ports = append(ports, newPort)
		}

		// For now we will take ports and have container port == exposed port
		for _, port := range container.Ports {
			newPort := corev1.ContainerPort{
				ContainerPort: int32(port),
				Protocol:      "TCP",
			}
			ports = append(ports, newPort)
		}
		// Add environment variables
		for key, value := range container.Environment {
			newEnvar := corev1.EnvVar{
				Name:  key,
				Value: value,
			}
			envars = append(envars, newEnvar)
		}

		newContainer.Ports = ports
		newContainer.Env = envars

		r.log.Info("🌀 Container", "Ports", container.Ports)
		containers = append(containers, newContainer)
	}
	return containers, nil
}
