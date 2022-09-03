/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package flux

import (
	api "flux-framework/flux-operator/api/v1alpha1"
	"fmt"
)

type QueueBestEffortFIFO struct {
	*QueueTemplate
}

var _ Queue = &QueueBestEffortFIFO{}

const BestEffortFIFO = api.BestEffortFIFO

// Sorting algorithms
// sort jobs based on priority, and fallback to creation time
// We aren't currently using this, our heap is faux!
func byCreationTime(a, b interface{}) bool {
	/*	A := a.(*workload.Info)
		B := b.(*workload.Info)
		p1 := utilpriority.Priority(A.Obj)
		p2 := utilpriority.Priority(B.Obj)

		if p1 != p2 {
			return p1 > p2
		}
		return A.Obj.CreationTimestamp.Before(&B.Obj.CreationTimestamp)*/
	return true
}

func keyFunc(obj interface{}) string {
	job := obj.(*api.FluxJob)
	return fmt.Sprintf("%s/%s", job.Namespace, job.Name)
}

// Given a FluxSetup, return the Best Effort FIFO queue
func newQueueBestEffortFIFO(setup *api.FluxSetup) (Queue, error) {
	template := newQueueTemplate(keyFunc, byCreationTime)
	fifo := &QueueBestEffortFIFO{
		QueueTemplate: template,
	}

	err := fifo.Update(setup)
	return fifo, err
}

/*
func (q *QueueBestEffortFIFO) RequeueIfNotPresent(wInfo *workload.Info, reason RequeueReason) bool {
	return cq.ClusterQueueImpl.requeueIfNotPresent(wInfo, reason == RequeueReasonFailedAfterNomination)
}*/
