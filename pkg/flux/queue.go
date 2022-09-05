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

	// Is a job running in the queue?
	IsRunningJob(info *jobctrl.Info) bool

	// Is a job waiting?
	IsWaitingJob(info *jobctrl.Info) bool

	// PushOrUpdate adds the job to the heap
	PushOrUpdate(*jobctrl.Info)

	// Delete a job from waiting or running
	Delete(*jobctrl.Info) bool

	// Pending is the count of pending jobs
	Pending() int

	// Running is the count of actively running jobs
	Running() int
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
