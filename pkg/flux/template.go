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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	api "flux-framework/flux-operator/api/v1alpha1"
	jobctrl "flux-framework/flux-operator/pkg/job"
	"flux-framework/flux-operator/pkg/util/heap"
	//	"sigs.k8s.io/kueue/pkg/util/heap"
	//	"sigs.k8s.io/kueue/pkg/workload"
)

// The QueueTemplate interface can be overwritten by specific Queue classes
type QueueTemplate struct {
	QueueStrategy api.QueueStrategy

	// This is probably a bad design and we should have a continuously
	// running controller, but for now I have a waiting jobs lookup,
	// and a heap that waiting jobs are moved to
	heap              heap.Heap
	cohort            string
	namespaceSelector labels.Selector

	// Waiting jobs have tried to be requested -> waiting -> run but we didn't have resources
	waitingJobs map[string]*jobctrl.Info
}

func newQueueTemplate(keyFunc func(obj interface{}) string, lessFunc func(a, b interface{}) bool) *QueueTemplate {
	return &QueueTemplate{

		// Note this isn's a true heap yet!
		heap:        heap.New(keyFunc, lessFunc),
		waitingJobs: make(map[string]*jobctrl.Info),
	}
}

var _ Queue = &QueueTemplate{}

func (c *QueueTemplate) Update(setup *api.FluxSetup) error {
	c.QueueStrategy = setup.Spec.QueueStrategy
	//	c.cohort = setup.Spec.Cohort
	nsSelector, err := metav1.LabelSelectorAsSelector(setup.Spec.NamespaceSelector)
	if err != nil {
		return err
	}
	c.namespaceSelector = nsSelector
	return nil
}

// QueueWaitingJobs moves waiting jobs to the heap, returning true/false if a workflow is moved
func (c *QueueTemplate) QueueWaitingJobs(ctx context.Context, client client.Client) bool {

	// Cut out early if we don't have waiting jobs
	if len(c.waitingJobs) == 0 {
		return false
	}

	log := ctrl.LoggerFrom(context.TODO())
	log.Info("Queue", "Waiting Jobs:", c.waitingJobs)
	waitingJobs := make(map[string]*jobctrl.Info)
	wasMoved := false
	for key, jobInfo := range c.waitingJobs {
		ns := corev1.Namespace{}
		err := client.Get(ctx, types.NamespacedName{Name: jobInfo.Obj.Namespace}, &ns)
		if err != nil || !c.namespaceSelector.Matches(labels.Set(ns.Labels)) {
			log.Info("Queue", "Job Still Waiting", jobInfo.Obj.Name)
			waitingJobs[key] = jobInfo
		} else {
			// Note that this is a stupid function that just adds the job info
			// if it doesn't exist yet, the actual heap functionality
			// needs to be implemented
			wasMoved = c.heap.PushIfNotPresent(jobInfo) || wasMoved
			log.Info("Queue", "Job Added to Heap", jobInfo.Obj.Name)
		}
	}
	c.waitingJobs = waitingJobs
	return wasMoved
}

func (c *QueueTemplate) PushOrUpdate(info *jobctrl.Info) {
	key := info.JobKey()
	oldInfo := c.waitingJobs[key]

	// We already have seen the job, it is waiting!
	if oldInfo != nil {
		// update in place if the job didn't change
		if equality.Semantic.DeepEqual(oldInfo.Obj.Spec, info.Obj.Spec) {
			c.waitingJobs[key] = info
			return
		}
		// If they aren't equal, update in place.
		delete(c.waitingJobs, key)
	}
	c.heap.PushOrUpdate(info)
}

func (c *QueueTemplate) IsRunningJob(info *jobctrl.Info) bool {
	key := info.JobKey()

	// Is it waiting?
	if c.waitingJobs[key] != nil {
		return false
	}

	// Is it currently in the heap (running)
	return c.heap.Exists(info)
}

func (c *QueueTemplate) Running() int {
	return c.heap.Len()
}

func (c *QueueTemplate) Pending() int {
	return len(c.waitingJobs)
}
