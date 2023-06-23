/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC

	(c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/
package job

import (
	api "flux-framework/flux-operator/api/v1alpha1"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/mitchellh/hashstructure/v2"
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
	Obj *api.MiniCluster

	// If admitted, the name of the running queue
	RunningQueue string
}

func (i *Info) JobKey() string {
	return i.Obj.Name
}

func HasCondition(job *api.MiniCluster, condition string) bool {
	i := FindCondition(&job.Status, condition)
	return i != -1 && job.Status.Conditions[i].Status == metav1.ConditionTrue
}

// FindConditionIndex finds the provided condition from the given status and returns the index.
// Returns -1 if the condition is not present.
func FindCondition(status *api.MiniClusterStatus, conditionType string) int {

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

// UpdateCondition sets all conditions to false except for the selected
func UpdateCondition(job *api.MiniCluster, conditionType string) {
	for i, condition := range job.Status.Conditions {
		if condition.Type == conditionType {
			job.Status.Conditions[i].Status = metav1.ConditionTrue
		} else {
			job.Status.Conditions[i].Status = metav1.ConditionFalse
		}
	}
}

// GetCondition gets the active condition
// If we eventually allow more than one condition this can return multiple
func GetCondition(job *api.MiniCluster) string {
	for _, condition := range job.Status.Conditions {
		if condition.Status == metav1.ConditionTrue {
			return condition.Reason
		}
	}
	return "AllFalse"
}

func NewInfo(job *api.MiniCluster) *Info {
	info := &Info{
		Obj: job,
	}
	return info
}

// JobsEqual takes a hash of the specs and assesses equality
func JobsEqual(jobA *api.MiniCluster, jobB *api.MiniCluster) bool {

	// For any failures, assume not equal
	hashA, err := hashstructure.Hash(jobA.Spec, hashstructure.FormatV2, nil)
	if err != nil {
		return false
	}
	hashB, err := hashstructure.Hash(jobB.Spec, hashstructure.FormatV2, nil)
	if err != nil {
		return false
	}
	return hashA == hashB
}

func FlagConditionWaiting(job *api.MiniCluster) {
	UpdateCondition(job, ConditionJobWaiting)
}

func FlagConditionReady(job *api.MiniCluster) {
	UpdateCondition(job, ConditionJobReady)
}
func FlagConditionRunning(job *api.MiniCluster) {
	UpdateCondition(job, ConditionJobRunning)
}

func FlagConditionFinished(job *api.MiniCluster) {
	UpdateCondition(job, ConditionJobFinished)
}

// Determine if the job is finished
func IsFinished(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == "Complete" {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}
