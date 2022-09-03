/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package flux

import (
	"context"
	"fmt"

	api "flux-framework/flux-operator/api/v1alpha1"
	jobctrl "flux-framework/flux-operator/pkg/job"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

/*
type RequeueReason string

const (
	RequeueReasonFailedAfterNomination RequeueReason = "FailedAfterNomination"
	RequeueReasonNamespaceMismatch     RequeueReason = "NamespaceMismatch"
	RequeueReasonGeneric               RequeueReason = ""
)*/

// A Queue a Jobs waiting to be scheduled
type Queue interface {

	// Ensure waiting jobs are added to the heap
	QueueWaitingJobs(ctx context.Context, client client.Client) bool

	// PushOrUpdate adds the job to the heap
	PushOrUpdate(*jobctrl.Info)

	// Pending is the count of pending jobs
	Pending() int

	// PendingActive is the count of active jobs
	PendingActive() int

	// PendingWaiting are count of waiting jobs.
	PendingWaiting() int

	/*
		// Update updates the properties of this ClusterQueue.
		Update(*kueue.ClusterQueue) error
		// Cohort returns the Cohort of this ClusterQueue.
		Cohort() string

		// AddFromQueue pushes all workloads belonging to this queue to
		// the ClusterQueue. If at least one workload is added, returns true.
		// Otherwise returns false.
		AddFromLocalQueue(*LocalQueue) bool
		// DeleteFromQueue removes all workloads belonging to this queue from
		// the ClusterQueue.
		DeleteFromLocalQueue(*LocalQueue)

		// Delete removes the workload from ClusterQueue.
		Delete(*kueue.Workload)
		// Pop removes the head of the queue and returns it. It returns nil if the
		// queue is empty.
		Pop() *workload.Info

		// RequeueIfNotPresent inserts a workload that was not
		// admitted back into the ClusterQueue. If the boolean is true,
		// the workloads should be put back in the queue immediately,
		// because we couldn't determine if the workload was admissible
		// in the last cycle. If the boolean is false, the implementation might
		// choose to keep it in temporary placeholder stage where it doesn't
		// compete with other workloads, until cluster events free up quota.
		// The workload should not be reinserted if it's already in the ClusterQueue.
		// Returns true if the workload was inserted.
		RequeueIfNotPresent(*workload.Info, RequeueReason) bool

		// Dump produces a dump of the current workloads in the heap of
		// this ClusterQueue. It returns false if the queue is empty.
		// Otherwise returns true.
		Dump() (sets.String, bool)
		DumpInadmissible() (sets.String, bool)
		// Info returns workload.Info for the workload key.
		// Users of this method should not modify the returned object.
		Info(string) *workload.Info*/

}

// From a FluxSetup, based on the QueueStrategy return a Queue
var registry = map[api.QueueStrategy]func(setup *api.FluxSetup) (Queue, error){
	BestEffortFIFO: newQueueBestEffortFIFO,
}

func newQueue(setup *api.FluxSetup) (Queue, error) {
	strategy := setup.Spec.QueueStrategy
	function, exist := registry[strategy]
	if !exist {
		return nil, fmt.Errorf("invalid QueueStrategy %q", setup.Spec.QueueStrategy)
	}
	return function(setup)
}
