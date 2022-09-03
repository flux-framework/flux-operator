/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package job

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (

	// FluxJob State
	// Requested: The default when the job is requested
	// Waiting: The Mini Cluster was created and is waiting for resources
	// Running:
	ConditionJobRequested string = "JobRequested"
	ConditionJobWaiting   string = "JobWaitingForResources"
	ConditionJobRunning   string = "JobRunning"
	ConditionJobFinished  string = "JobFinished"
)

func getJobRequestedCondition(status metav1.ConditionStatus) metav1.Condition {
	now := time.Now()
	return metav1.Condition{
		Type:               ConditionJobRequested,
		Reason:             ConditionJobRequested,
		Status:             status,
		LastTransitionTime: metav1.Time{Time: now},
	}
}

func getJobWaitingCondition(status metav1.ConditionStatus) metav1.Condition {
	now := time.Now()
	return metav1.Condition{
		Type:               ConditionJobWaiting,
		Reason:             ConditionJobWaiting,
		Status:             status,
		LastTransitionTime: metav1.Time{Time: now},
	}
}

func getJobRunningCondition(status metav1.ConditionStatus) metav1.Condition {
	now := time.Now()
	return metav1.Condition{
		Type:               ConditionJobRunning,
		Reason:             ConditionJobRunning,
		Status:             status,
		LastTransitionTime: metav1.Time{Time: now},
	}
}

func getJobFinishedCondition(status metav1.ConditionStatus) metav1.Condition {
	now := time.Now()
	return metav1.Condition{
		Type:               ConditionJobFinished,
		Reason:             ConditionJobFinished,
		Status:             status,
		LastTransitionTime: metav1.Time{Time: now},
	}
}

// getJobConditions inits the job conditions. By default, the job
// request is true since this is the first state it hits!
func GetJobConditions() []metav1.Condition {
	return []metav1.Condition{
		getJobRequestedCondition(metav1.ConditionTrue),
		getJobWaitingCondition(metav1.ConditionFalse),
		getJobRunningCondition(metav1.ConditionFalse),
		getJobFinishedCondition(metav1.ConditionFalse),
	}
}
