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

	// MiniCluster State
	// Requested:
	//     The default when the job is requested
	// Waiting + in waiting queue:
	//     The job is waiting to be admitted (there are resources) - before added to heap
	// Waiting + in heap:
	//     There are resources, and the job has permission to run. We are creating the
	//     MiniCluster (in the heap) and reconciling until it's all there.
	// Ready:
	//     We finished creating the MiniCluster, and it's ready to run!
	// Running:
	//     Resources are created, and status is switched from Waiting -> Running by MiniCluster
	//     This should be when we kick off the command. The job is running
	// Finished:
	//     When resources are done (TBA how determined)
	// Finished the job is finished running!
	ConditionJobRequested string = "JobRequested"
	ConditionJobWaiting   string = "JobWaitingForResources"
	ConditionJobReady     string = "JobMiniClusterReady"
	ConditionJobRunning   string = "JobRunning"
	ConditionJobFinished  string = "JobFinished"
)

func getJobRequestedCondition(status metav1.ConditionStatus) metav1.Condition {
	now := time.Now()
	return metav1.Condition{
		Type:               ConditionJobRequested,
		Reason:             ConditionJobRequested,
		Status:             status,
		Message:            ConditionJobRequested,
		LastTransitionTime: metav1.Time{Time: now},
	}
}

func getJobReadyCondition(status metav1.ConditionStatus) metav1.Condition {
	now := time.Now()
	return metav1.Condition{
		Type:               ConditionJobReady,
		Reason:             ConditionJobReady,
		Status:             status,
		Message:            ConditionJobReady,
		LastTransitionTime: metav1.Time{Time: now},
	}
}

func getJobWaitingCondition(status metav1.ConditionStatus) metav1.Condition {
	now := time.Now()
	return metav1.Condition{
		Type:               ConditionJobWaiting,
		Reason:             ConditionJobWaiting,
		Status:             status,
		Message:            ConditionJobWaiting,
		LastTransitionTime: metav1.Time{Time: now},
	}
}

func getJobRunningCondition(status metav1.ConditionStatus) metav1.Condition {
	now := time.Now()
	return metav1.Condition{
		Type:               ConditionJobRunning,
		Reason:             ConditionJobRunning,
		Status:             status,
		Message:            ConditionJobRunning,
		LastTransitionTime: metav1.Time{Time: now},
	}
}

func getJobFinishedCondition(status metav1.ConditionStatus) metav1.Condition {
	now := time.Now()
	return metav1.Condition{
		Type:               ConditionJobFinished,
		Reason:             ConditionJobFinished,
		Status:             status,
		Message:            ConditionJobFinished,
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
		getJobReadyCondition(metav1.ConditionFalse),
		getJobFinishedCondition(metav1.ConditionFalse),
	}
}
