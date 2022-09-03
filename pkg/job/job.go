/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/
package job

import (
	api "flux-framework/flux-operator/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Job statuses for logging purposes
	Waiting   = "waiting"
	Finished  = "finished"
	Running   = "running"
	Requested = "requested"
)

// Info holds job data in the flux manager
type Info struct {
	Obj *api.FluxJob

	// If admitted, the name of the running queue
	RunningQueue string
}

func (i *Info) JobKey() string {
	return i.Obj.Name
}

func HasCondition(job *api.FluxJob, condition string) bool {
	i := FindCondition(&job.Status, condition)
	return i != -1 && job.Status.Conditions[i].Status == metav1.ConditionTrue
}

// FindConditionIndex finds the provided condition from the given status and returns the index.
// Returns -1 if the condition is not present.
func FindCondition(status *api.FluxJobStatus, conditionType string) int {

	// Index of -1 means we have zero conditions
	if status == nil || status.Conditions == nil {
		return -1
	}
	for i := range status.Conditions {
		// We found the index!
		if status.Conditions[i].Type == conditionType {
			return i
		}
	}
	return -1
}

func UpdateCondition(job *api.FluxJob, conditionType string) {
	for i, condition := range job.Status.Conditions {
		if condition.Type == conditionType {
			job.Status.Conditions[i].Status = metav1.ConditionTrue
		} else {
			job.Status.Conditions[i].Status = metav1.ConditionFalse
		}
	}
}

func NewInfo(job *api.FluxJob) *Info {
	info := &Info{
		Obj: job,
	}
	return info
}

func FlagConditionWaiting(job *api.FluxJob) {
	UpdateCondition(job, ConditionJobWaiting)
}
