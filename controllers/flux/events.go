/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

// Events are added to the Reconciler directly. If we don't need them:
// 1. Delete this file
// 2. Delete the AddEventFilter(r)
// 3. (Optionally) the Reconciler Client can be inherited directly

import (
	jobctrl "github.com/flux-framework/flux-operator/pkg/job"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"

	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

// Notify watchers (the FluxSetup) that we have a new job request
func (r *MiniClusterReconciler) notifyWatchers(job *api.MiniCluster) {
	for _, watcher := range r.watchers {
		watcher.NotifyMiniClusterUpdate(job)
	}
}

// Called when a new job is created
func (r *MiniClusterReconciler) Create(e event.CreateEvent) bool {

	// Only respond to job events!
	job, match := e.Object.(*api.MiniCluster)
	if !match {
		return true
	}

	// Add conditions - they should never exist for a new job
	job.Status.Conditions = jobctrl.GetJobConditions()

	// We will tell FluxSetup there is a new job request
	defer r.notifyWatchers(job)
	r.log.Info("🌀 MiniCluster create event", "Name:", job.Name)

	// Continue to creation event
	r.log.Info("🌀 MiniCluster was added!", "Name:", job.Name, "Condition:", jobctrl.GetCondition(job))
	return true
}

func (r *MiniClusterReconciler) Delete(e event.DeleteEvent) bool {

	job, match := e.Object.(*api.MiniCluster)
	if !match {
		return true
	}

	defer r.notifyWatchers(job)
	log := r.log.WithValues("job", klog.KObj(job))
	log.Info("🌀 MiniCluster delete event")

	// TODO should trigger a delete here
	// Reconcile should clean up resources now
	return true
}

func (r *MiniClusterReconciler) Update(e event.UpdateEvent) bool {
	oldMC, match := e.ObjectOld.(*api.MiniCluster)
	if !match {
		return true
	}

	// Figure out the state of the old job
	mc := e.ObjectNew.(*api.MiniCluster)

	r.log.Info("🌀 MiniCluster update event")

	// If the job hasn't changed, continue reconcile
	// There aren't any explicit updates beyond conditions
	if jobctrl.JobsEqual(mc, oldMC) {
		return true
	}

	// TODO: check if ready or running, shouldn't be able to update
	// OR if we want update, we need to completely delete and recreate
	return true
}

func (r *MiniClusterReconciler) Generic(e event.GenericEvent) bool {
	r.log.V(3).Info("Ignore generic event", "obj", klog.KObj(e.Object), "kind", e.Object.GetObjectKind().GroupVersionKind())
	return false
}
