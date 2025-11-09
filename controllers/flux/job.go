/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"reflect"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"
)

// newMiniCluster is used to create the MiniCluster Job
func NewMiniClusterJob(cluster *api.MiniCluster) (*batchv1.Job, error) {

	// Number of retries before marking as failed
	backoffLimit := int32(100)
	if cluster.Spec.BackoffLimit != 0 {
		backoffLimit = cluster.Spec.BackoffLimit
	}
	completionMode := batchv1.IndexedCompletion
	podLabels := getPodLabels(cluster)
	setAsFQDN := false

	// We add the selector for the horizontal auto scaler, if active
	// We can't use the job-name selector, as this would include the
	// external sidecar service!
	podLabels["hpa-selector"] = cluster.Name

	// Add tolerations
	tolerations := []corev1.Toleration{}
	for _, tspec := range cluster.Spec.Pod.Tolerations {
		toleration := corev1.Toleration{
			Effect:   corev1.TaintEffect(tspec.Effect),
			Key:      tspec.Key,
			Operator: corev1.TolerationOperator(tspec.Operator),
			Value:    tspec.Value,
		}
		tolerations = append(tolerations, toleration)
	}

	// This is an indexed-job
	// TODO don't hard code type meta
	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
			Labels:    cluster.Spec.JobLabels,
		},

		// Completions must be == to Parallelism to allow for scaling
		Spec: batchv1.JobSpec{
			BackoffLimit:          &backoffLimit,
			Completions:           &cluster.Spec.Size,
			Parallelism:           &cluster.Spec.Size,
			CompletionMode:        &completionMode,
			ActiveDeadlineSeconds: &cluster.Spec.DeadlineSeconds,

			// Note there is parameter to limit runtime
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        cluster.Name,
					Namespace:   cluster.Namespace,
					Labels:      podLabels,
					Annotations: cluster.Spec.Pod.Annotations,
				},
				Spec: corev1.PodSpec{
					// matches the service
					Subdomain:                    cluster.Spec.Network.HeadlessName,
					ShareProcessNamespace:        &cluster.Spec.ShareProcessNamespace,
					SetHostnameAsFQDN:            &setAsFQDN,
					Volumes:                      getVolumes(cluster),
					ImagePullSecrets:             getImagePullSecrets(cluster),
					ServiceAccountName:           cluster.Spec.Pod.ServiceAccountName,
					AutomountServiceAccountToken: &cluster.Spec.Pod.AutomountServiceAccountToken,
					RestartPolicy:                corev1.RestartPolicyOnFailure,
					NodeSelector:                 cluster.Spec.Pod.NodeSelector,
					Tolerations:                  tolerations,
					SchedulerName:                cluster.Spec.Pod.SchedulerName,
					HostPID:                      cluster.Spec.Pod.HostPID,
					HostIPC:                      cluster.Spec.Pod.HostIPC,
					ResourceClaims:               cluster.Spec.ResourceClaims,
				},
			},
		},
	}

	// Set the security context
	if !reflect.DeepEqual(cluster.Spec.Pod.SecurityContext, corev1.PodSecurityContext{}) {
		securityContext := corev1.PodSecurityContext{
			FSGroup: cluster.Spec.Pod.SecurityContext.FSGroup,
		}
		job.Spec.Template.Spec.SecurityContext = &securityContext
	}

	// Custom restart policy
	if cluster.Spec.Pod.RestartPolicy != "" {
		job.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicy(cluster.Spec.Pod.RestartPolicy)
	}

	// Custom DNS Policy
	if cluster.Spec.Pod.DNSPolicy != "" {
		job.Spec.Template.Spec.DNSPolicy = corev1.DNSPolicy(cluster.Spec.Pod.DNSPolicy)
	}

	// Only add runClassName if defined
	if cluster.Spec.Pod.RuntimeClassName != "" {
		job.Spec.Template.Spec.RuntimeClassName = &cluster.Spec.Pod.RuntimeClassName
	}

	// Add Affinity to map one pod / node only if the user hasn't disabled it
	if !cluster.Spec.Network.DisableAffinity {
		job.Spec.Template.Spec.Affinity = getAffinity(cluster)
	}
	job.Spec.Template.Spec.Overhead = cluster.Spec.Pod.Resources

	// Get volume mounts specific to operator, add on container specific ones
	mounts := getVolumeMounts(cluster)

	// Get the flux view container (only if requested)
	fluxViewContainer, err := getFluxContainer(cluster, mounts)
	if err != nil {
		return job, err
	}

	// Add spack view (we need to copy back here)
	spackVolume := corev1.VolumeMount{
		Name:      spackSoftware,
		MountPath: spackSoftwarePath,
		ReadOnly:  false,
	}
	mounts = append(mounts, spackVolume)

	// Prepare listing of containers for the MiniCluster
	containers, err := getContainers(
		cluster.Spec.Containers,
		cluster.Name,
		mounts,
		false,
	)

	// Add on the flux view container
	job.Spec.Template.Spec.InitContainers = []corev1.Container{fluxViewContainer}
	job.Spec.Template.Spec.Containers = containers
	return job, err
}

// getAffinity returns to pod affinity to ensure 1 address / node
func getAffinity(cluster *api.MiniCluster) *corev1.Affinity {
	affinity := &corev1.Affinity{
		// Prefer to schedule pods on the same zone
		PodAffinity: &corev1.PodAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
				{
					Weight: 100,
					PodAffinityTerm: corev1.PodAffinityTerm{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									// added in getPodLabels
									Key:      podLabelAppName,
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{cluster.Name},
								},
							},
						},
						TopologyKey: "topology.kubernetes.io/zone",
					},
				},
			},
		},
		// Prefer to schedule pods on different nodes
		PodAntiAffinity: &corev1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
				{
					Weight: 100,
					PodAffinityTerm: corev1.PodAffinityTerm{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									// added in getPodLabels
									Key:      podLabelAppName,
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{cluster.Name},
								},
							},
						},
						TopologyKey: "kubernetes.io/hostname",
					},
				},
			},
		},
	}

	// Add custom affinity if defined to select subset of nodes in list
	if len(cluster.Spec.Pod.NodeAffinity) > 0 {
		requires := []corev1.NodeSelectorRequirement{}
		for label, values := range cluster.Spec.Pod.NodeAffinity {
			requires = append(requires, corev1.NodeSelectorRequirement{
				Key:      label,
				Operator: corev1.NodeSelectorOpIn,
				Values:   values,
			})
		}
		affinity.NodeAffinity = &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{MatchExpressions: requires},
				},
			},
		}
	}
	return affinity
}

// getImagePullSecrets returns a list of secret object references for each container.
func getImagePullSecrets(cluster *api.MiniCluster) []corev1.LocalObjectReference {
	pullSecrets := []corev1.LocalObjectReference{}
	for _, container := range cluster.Spec.Containers {
		if container.ImagePullSecret != "" {
			newSecret := corev1.LocalObjectReference{
				Name: container.ImagePullSecret,
			}
			pullSecrets = append(pullSecrets, newSecret)
		}
	}
	return pullSecrets
}
