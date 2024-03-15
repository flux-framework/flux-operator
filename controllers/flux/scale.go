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

	batchv1 "k8s.io/api/batch/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"
)

// addScaleSelector populates the fields the horizontal auto scaler needs.
// Meaning: job-name is used to select pods to check. The size variable
// is updated later.
func (r *MiniClusterReconciler) addScaleSelector(
	ctx context.Context,
	cluster *api.MiniCluster,
) (ctrl.Result, error) {

	// Update the pod selector to allow horizontal autoscaling
	selector := "hpa-selector=" + cluster.Name
	if cluster.Status.Selector == selector {
		r.log.Info("MiniCluster", "ScaleSelector", selector, "Status", "Ready")
		return ctrl.Result{}, nil
	}
	cluster.Status.Selector = selector
	err := r.Status().Update(ctx, cluster)
	r.log.Info("MiniCluster", "ScaleSelector", selector, "Status", "Updating")
	return ctrl.Result{Requeue: true}, err
}

// disallowScale is called when the size is > the maximum size allowed, and we only scale up to that
func (r *MiniClusterReconciler) disallowScale(
	ctx context.Context,
	job *batchv1.Job,
	cluster *api.MiniCluster,
) (ctrl.Result, error) {

	r.log.Info("MiniCluster", "PatchSize", cluster.Spec.Size, "Status", "Denied")
	patch := client.MergeFrom(cluster.DeepCopy())
	cluster.Spec.Size = cluster.Status.MaximumSize
	cluster.Status.Size = cluster.Status.MaximumSize

	// Apply the patch to restore to the original size
	err := r.Patch(ctx, cluster, patch)

	// First update fixes the status
	r.Status().Update(ctx, cluster)
	return ctrl.Result{Requeue: true}, err
}

// allowscale is called when we determine the size is >=1 and <= maxSize
func (r *MiniClusterReconciler) allowScale(
	ctx context.Context,
	job *batchv1.Job,
	cluster *api.MiniCluster,
) (ctrl.Result, error) {

	r.log.Info("MiniCluster", "PatchSize", cluster.Spec.Size, "Status", "Accepted")
	patch := client.MergeFrom(job.DeepCopy())
	job.Spec.Parallelism = &cluster.Spec.Size
	job.Spec.Completions = &cluster.Spec.Size
	cluster.Status.Size = cluster.Spec.Size

	err := r.Patch(ctx, job, patch)
	// I don't check for error because I want both changes to go in at once
	r.Status().Update(ctx, cluster)
	return ctrl.Result{Requeue: true}, err
}

// restoreOriginalSize is called when the request for a size is < 1, and we don't allow it
func (r *MiniClusterReconciler) restoreOriginalSize(
	ctx context.Context,
	job *batchv1.Job,
	cluster *api.MiniCluster,
) (ctrl.Result, error) {

	r.log.Info("MiniCluster", "PatchSize", cluster.Spec.Size, "Status", "Denied")
	patch := client.MergeFrom(cluster.DeepCopy())
	cluster.Spec.Size = *job.Spec.Parallelism
	cluster.Status.Size = cluster.Spec.Size

	// Apply the patch to restore to the original size
	r.Status().Update(ctx, cluster)
	err := r.Patch(ctx, cluster, patch)
	return ctrl.Result{Requeue: true}, err
}

// resizeCluster will patch the cluster to make a larger (or smaller) size
func (r *MiniClusterReconciler) resizeCluster(
	ctx context.Context,
	job *batchv1.Job,
	cluster *api.MiniCluster,
) (ctrl.Result, error) {

	// We absolutely don't allow a size less than 1
	// If this happens, restore to current / original size
	if cluster.Spec.Size < 1 {
		return r.restoreOriginalSize(ctx, job, cluster)
	}

	// ensure we don't go above the max original size, which should be saved on init
	// If we do, we need to patch it back down to the maximum - this isn't allowed
	if cluster.Spec.Size > cluster.Status.MaximumSize {
		return r.disallowScale(ctx, job, cluster)
	}

	// If we get here, the size is smaller and we allow it!
	return r.allowScale(ctx, job, cluster)
}
