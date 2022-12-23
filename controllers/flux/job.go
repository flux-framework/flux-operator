/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// newMiniCluster is used to create the MiniCluster Job
func (r *MiniClusterReconciler) newMiniClusterJob(cluster *api.MiniCluster) *batchv1.Job {

	// We need to create the number of containers (and names) that the user requests
	// Before the stateful set was doing this for us, but for a batch job it's manaul
	containers := r.getMiniClusterContainers(cluster)

	// Number of retries before marking as failed
	backoffLimit := int32(100)
	completionMode := batchv1.IndexedCompletion
	setAsFQDN := false

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{

			// In the example from Dan this was "indexed-job"
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
		},

		Spec: batchv1.JobSpec{

			BackoffLimit:          &backoffLimit,
			Completions:           &cluster.Spec.Size,
			Parallelism:           &cluster.Spec.Size,
			CompletionMode:        &completionMode,
			ActiveDeadlineSeconds: &cluster.Spec.DeadlineSeconds,

			// Note there is parameter to limit runtime
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cluster.Name,
					Namespace: cluster.Namespace,

					// Job is used as selector for service
					Labels: map[string]string{"name": cluster.Name, "namespace": cluster.Namespace, "job": serviceName},
				},
				Spec: corev1.PodSpec{
					// matches the service
					Subdomain:         serviceName,
					SetHostnameAsFQDN: &setAsFQDN,
					Volumes:           getVolumes(cluster),
					Containers:        containers,
					RestartPolicy:     corev1.RestartPolicyOnFailure,
					ImagePullSecrets:  getImagePullSecrets(cluster),
				}},
		},
	}
	ctrl.SetControllerReference(cluster, job, r.Scheme)
	return job
}

func (r *MiniClusterReconciler) getMiniClusterContainers(cluster *api.MiniCluster) []corev1.Container {

	// Create the containers for the pod
	containers := []corev1.Container{}

	for i, container := range cluster.Spec.Containers {

		// Allow dictating pulling on the level of the container
		pullPolicy := corev1.PullIfNotPresent
		if container.PullAlways {
			pullPolicy = corev1.PullAlways
		}

		// Fluxrunner will use the namespace name
		containerName := container.Name
		command := []string{}

		// This is a wrapper that is going to wait for the generation of update_hosts.sh
		// Once it's there, we update /etc/hosts, and run the command to start flux.
		if container.FluxRunner {

			// wait.sh path corresponds to container identifier
			waitScript := fmt.Sprintf("/flux_operator/wait-%d.sh", i)
			command = []string{"/bin/bash", waitScript, container.Command}
			containerName = fmt.Sprintf("%s-%d", cluster.Name, i)
		}

		// Do we have a postStartExec Lifecycle command?
		lifecycle := corev1.Lifecycle{}
		if container.LifeCyclePostStartExec != "" {
			r.log.Info("🌀 MiniCluster", "LifeCycle.PostStartExec", container.LifeCyclePostStartExec)
			lifecycle.PostStart = &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: []string{container.LifeCyclePostStartExec},
				},
			}
		}

		// Get volume mounts, add on container specific ones
		mounts := getVolumeMounts(cluster)
		for volumeName, volume := range container.Volumes {
			newVolume := corev1.VolumeMount{
				Name:      volumeName,
				MountPath: volume.Path,
				ReadOnly:  volume.ReadOnly,
			}
			mounts = append(mounts, newVolume)
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
			Lifecycle:       &lifecycle,
		}

		// Only add command if we actually have one
		if len(command) > 0 {
			newContainer.Command = command
		}

		ports := []corev1.ContainerPort{}
		envars := []corev1.EnvVar{}

		// If it's the FluxRunner, expose port 5000 for the service
		if container.FluxRunner {
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
		for key, value := range container.Envars {
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
	return containers
}

// getImagePullSecrets returns a list of secret object references for each container.
func getImagePullSecrets(cluster *api.MiniCluster) []corev1.LocalObjectReference {
	pullSecrets := []corev1.LocalObjectReference{}
	for _, container := range cluster.Spec.Containers {
		if container.ImagePullSecret != "" {
			newSecret := corev1.LocalObjectReference{Name: container.ImagePullSecret}
			pullSecrets = append(pullSecrets, newSecret)
		}
	}
	return pullSecrets
}
