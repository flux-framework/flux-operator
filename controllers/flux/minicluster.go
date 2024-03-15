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
	"strings"

	jobctrl "github.com/flux-framework/flux-operator/pkg/job"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"
)

var (
	CurveCertKey = "curve.cert"
)

// This is a MiniCluster! A MiniCluster is associated with a running MiniCluster and include:
// 1. An indexed job with some number of pods
// 2. Config maps for secrets and other things.
// 3. We "launch" a job by starting the Indexed job on the connected nodes
// ensureMiniCluster creates a new MiniCluster, a stateful set for running flux!
func (r *MiniClusterReconciler) ensureMiniCluster(
	ctx context.Context,
	cluster *api.MiniCluster,
) (ctrl.Result, error) {

	// Initial config with entrypoints for flux and main containers
	_, result, err := r.getConfigMap(ctx, cluster, cluster.EntrypointConfigMapName())
	if err != nil {
		return result, err
	}

	// Any extra service containers (running alongside the cluster)
	// For now run these in the same pod, one service pod
	if len(cluster.Spec.Services) > 0 {
		_, result, err = r.ensureServicePod(ctx, cluster)
		if err != nil {
			return result, err
		}
	}

	// Create headless service for the MiniCluster OR single service for the broker
	selector := map[string]string{"job-name": cluster.Name}

	// If we are adding a minimal service expose the index 0 pod only
	// LabelSelectors are ANDed
	if cluster.Spec.Flux.MinimalService {
		selector["job-index"] = "0"
	}

	result, err = r.exposeServices(ctx, cluster, cluster.Spec.Network.HeadlessName, selector)
	if err != nil {
		return result, err
	}

	// Create the batch job that brings it all together!
	// A batchv1.Job can hold a spec for containers that use the configs we just made
	mc, result, err := r.getMiniCluster(ctx, cluster)
	if err != nil {
		return result, err
	}

	// If the sizes are different, we patch to update.
	// This would be an explicit update from a user or application via the CRD to scale up/down
	// The Flux Operator can't tell the difference between these two, but honors maxSize
	// (the size the Flux broker leader knows about) and updates both Spec.Size and Status.Size
	if *mc.Spec.Parallelism != cluster.Spec.Size {
		r.log.Info("MiniCluster", "Size", mc.Spec.Parallelism, "Requested Size", cluster.Spec.Size)
		result, err := r.resizeCluster(ctx, mc, cluster)
		if err != nil {
			return result, err
		}
	}

	// Add selector to allow horizontal pod autoscaler
	// This would be done via a request to a running metrics server
	// If there is no autoscaler, has no impact. The .Status.Size
	// should already be updated via the function above.
	result, err = r.addScaleSelector(ctx, cluster)
	if err != nil {
		return result, err
	}

	// Expose other sidecar container services
	for _, container := range cluster.Spec.Containers {

		// Assume now services only available TO flux runner
		if container.RunFlux || len(container.Ports) == 0 {
			continue
		}

		// Service name corresponds to container, but selector is pod-specific
		selector := map[string]string{podLabelAppName: cluster.Name}
		result, err = r.exposeService(ctx, cluster, container.Name, selector, container.Ports)
		if err != nil {
			return result, err
		}
	}

	// Add the single label for the broker pod
	result, err = r.addBrokerLabel(ctx, cluster)
	if err != nil {
		return result, err
	}

	// If we get here, update the status to be ready
	status := jobctrl.GetCondition(cluster)
	if status != jobctrl.ConditionJobReady {
		clusterCopy := cluster.DeepCopy()
		jobctrl.FlagConditionReady(clusterCopy)
		r.Status().Update(ctx, clusterCopy)
	}

	// And we re-queue so the Ready condition triggers next steps!
	return ctrl.Result{Requeue: true}, nil
}

// cleanupPodsStorage looks for the existing job, and cleans up if completed
func (r *MiniClusterReconciler) cleanupPodsStorage(
	ctx context.Context,
	cluster *api.MiniCluster,
) (ctrl.Result, error) {

	// Find the broker pod and determine if finished
	completed := false
	for _, pod := range r.getMiniClusterPods(ctx, cluster).Items {
		if !strings.HasPrefix(pod.Name, fmt.Sprintf("%s-0", cluster.Name)) {
			continue
		}

		// If it's succeeded or failed, we call that finished
		// https://pkg.go.dev/k8s.io/api@v0.25.0/core/v1#PodPhase
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			completed = true
			break
		}
	}

	// Cut out early if not completed
	if !completed {
		r.log.Info("MiniCluster", "Job Status", "Not Completed")
		return ctrl.Result{Requeue: true}, nil
	}
	r.log.Info("MiniCluster", "Job Status", "Completed")

	// Delete the MiniCluster first
	// If we don't, it will keep re-creating the assets and loop forever :)
	r.Client.Delete(ctx, cluster)
	return ctrl.Result{Requeue: false}, nil
}

// getExistingJob gets an existing job that matches the MiniCluster CRD
func (r *MiniClusterReconciler) getExistingJob(
	ctx context.Context,
	cluster *api.MiniCluster,
) (*batchv1.Job, error) {

	existing := &batchv1.Job{}
	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
		},
		existing,
	)
	return existing, err
}

// getMiniCluster does an actual check if we have a batch job in the namespace
func (r *MiniClusterReconciler) getMiniCluster(
	ctx context.Context,
	cluster *api.MiniCluster,
) (*batchv1.Job, ctrl.Result, error) {

	// Look for an existing job
	existing, err := r.getExistingJob(ctx, cluster)

	// Create a new job if it does not exist
	if err != nil {

		if errors.IsNotFound(err) {
			job, err := NewMiniClusterJob(cluster)
			ctrl.SetControllerReference(cluster, job, r.Scheme)

			if err != nil {
				r.log.Error(
					err, "Failed to create new MiniCluster Batch Job",
					"Namespace:", job.Namespace,
					"Name:", job.Name,
				)
				return job, ctrl.Result{}, err
			}

			r.log.Info(
				"✨ Creating a new MiniCluster Batch Job ✨",
				"Namespace:", job.Namespace,
				"Name:", job.Name,
			)

			err = r.New(ctx, job)
			if err != nil {
				r.log.Error(
					err,
					"Failed to create new MiniCluster Batch Job",
					"Namespace:", job.Namespace,
					"Name:", job.Name,
				)
				return job, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return job, ctrl.Result{Requeue: true}, nil

		} else if err != nil {
			r.log.Error(err, "Failed to get MiniCluster Batch Job")
			return existing, ctrl.Result{}, err
		}

	} else {
		r.log.Info(
			"🎉 Found existing MiniCluster Batch Job 🎉",
			"Namespace:", existing.Namespace,
			"Name:", existing.Name,
		)
	}
	return existing, ctrl.Result{}, err
}
