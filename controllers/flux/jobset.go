/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "flux-framework/flux-operator/api/v1alpha1"

	ctrl "sigs.k8s.io/controller-runtime"
	jobset "sigs.k8s.io/jobset/api/v1alpha1"
)

func (r *MiniClusterReconciler) newJobSet(
	cluster *api.MiniCluster,
) (*jobset.JobSet, error) {

	// Name for the broker job for the failure policy
	brokerJobName := cluster.Name + "-broker"

	// When suspend is true we have a hard time debugging jobs, so keep false
	suspend := false
	jobs := jobset.JobSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "minicluster",
			Namespace: cluster.Namespace,
			Labels:    cluster.Spec.JobLabels,
		},
		Spec: jobset.JobSetSpec{

			// The job is successful when the broker job finishes with completed (0)
			SuccessPolicy: &jobset.SuccessPolicy{
				Operator:             jobset.OperatorAny,
				TargetReplicatedJobs: []string{brokerJobName},
			},

			// This might be the control for child jobs (worker)
			// But I don't think we need this anymore.
			Suspend: &suspend,
			// TODO decide on FailurePolicy here
			// default is to fail if all jobs in jobset fail
		},
	}

	// Get leader broker job, the parent in the JobSet (worker or follower pods)
	// We have to do this as indexed for the predictable hostname
	//                    cluster, size, entrypoint, indexed
	leaderJob, err := r.getJob(cluster, 1, "broker", true)
	if err != nil {
		return &jobs, err
	}
	workerJob, err := r.getJob(cluster, cluster.Spec.Size-1, "worker", true)
	if err != nil {
		return &jobs, err
	}
	jobs.Spec.ReplicatedJobs = []jobset.ReplicatedJob{leaderJob, workerJob}
	ctrl.SetControllerReference(cluster, &jobs, r.Scheme)
	return &jobs, nil
}

// getJob creates a job for a main leader (broker) or worker (followers)
func (r *MiniClusterReconciler) getJob(
	cluster *api.MiniCluster,
	size int32,
	entrypoint string,
	indexed bool,
) (jobset.ReplicatedJob, error) {

	backoffLimit := int32(100)
	podLabels := r.getPodLabels(cluster)
	enableDNSHostnames := false
	completionMode := batchv1.NonIndexedCompletion

	if indexed {
		completionMode = batchv1.IndexedCompletion
	}

	job := jobset.ReplicatedJob{
		Name: cluster.Name + "-" + entrypoint,

		// This would allow pods to be reached by their hostnames!
		// It doesn't work for the Flux broker config at the moment,
		// but could if we are allowed to specify the service name.
		// <jobSet.name>-<spec.replicatedJob.name>-<job-index>-<pod-index>.<jobSet.name>-<spec.replicatedJob.name>
		Network: &jobset.Network{
			EnableDNSHostnames: &enableDNSHostnames,
		},

		Template: batchv1.JobTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cluster.Name,
				Namespace: cluster.Namespace,
				Labels:    cluster.Spec.JobLabels,
			},
		},
		// This is the default, but let's be explicit
		Replicas: 1,
	}

	// Create the JobSpec for the job -> Template -> Spec
	jobspec := batchv1.JobSpec{
		BackoffLimit:          &backoffLimit,
		Completions:           &size,
		Parallelism:           &size,
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
				Subdomain:          restfulServiceName,
				Volumes:            getVolumes(cluster, entrypoint),
				RestartPolicy:      corev1.RestartPolicyOnFailure,
				ImagePullSecrets:   getImagePullSecrets(cluster),
				ServiceAccountName: cluster.Spec.Pod.ServiceAccountName,
				NodeSelector:       cluster.Spec.Pod.NodeSelector,
			},
		},
	}
	// Get resources for the pod
	resources, err := r.getPodResources(cluster)
	r.log.Info("🌀 MiniCluster", "Pod.Resources", resources)
	if err != nil {
		r.log.Info("🌀 MiniCluster", "Pod.Resources", resources)
		return job, err
	}
	jobspec.Template.Spec.Overhead = resources

	// Get volume mounts, add on container specific ones
	mounts := getVolumeMounts(cluster)
	containers, err := r.getContainers(
		cluster.Spec.Containers,
		cluster.Name,
		mounts,
		entrypoint,
	)
	jobspec.Template.Spec.Containers = containers
	job.Template.Spec = jobspec
	return job, err
}
