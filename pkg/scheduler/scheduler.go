/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

// Currently I'm making the scheduler to ONLY exist to queue
// waiting jobs, and eventually it can do more checks about
// when is appropriate to do so.

package scheduler

import (
	"context"
	"flux-framework/flux-operator/pkg/flux"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	errCouldNotAdmitWL = "Could not admit workload and assigning flavors in apiserver"
)

type Scheduler struct {
	fluxManager *flux.Manager
	client      client.Client
	recorder    record.EventRecorder
	log         logr.Logger
}

func New(manager *flux.Manager, cl client.Client, recorder record.EventRecorder) *Scheduler {
	s := &Scheduler{
		fluxManager: manager,
		client:      cl,
		log:         ctrl.Log.WithName("scheduler"),
		recorder:    recorder,
	}
	//s.applyAdmission = s.applyAdmissionWithSSA
	return s
}

func (s *Scheduler) Start(ctx context.Context) {
	log := ctrl.LoggerFrom(ctx).WithName("dummy-scheduler")
	ctx = ctrl.LoggerInto(ctx, log)
	//	wait.UntilWithContext(ctx, s.schedule, 0)
	// Just run this forever for now until I know how to do otherwise!
	s.log.Info("📅 Scheduler", "Status:", "Running...")
	for true {
		s.schedule(ctx)
		time.Sleep(8 * time.Second)
	}
}

// Continue running this until we don't have waiting jobs?
func (s *Scheduler) schedule(ctx context.Context) {

	// We haven't created the queue yet
	if !s.fluxManager.HasQueue() {
		s.log.Info("📅 Scheduler", "Queue:", "No queue exists yet.")
		return
	}

	// Ensure we queue waiting jobs, no checks for now
	if s.fluxManager.QueueWaitingJobs(ctx) {
		s.log.Info("📅 Scheduler", "Status:", "New jobs were queued")
	}
	s.log.Info("📅 Scheduler", "Status:", fmt.Sprintf("%d jobs are pending", s.fluxManager.JobsPending()))
	s.log.Info("📅 Scheduler", "Status:", fmt.Sprintf("%d jobs are running", s.fluxManager.JobsRunning()))
}

// I think this will be important to patch something...
//func (s *Scheduler) applyAdmissionWithSSA(ctx context.Context, job *api.FluxJob) error {
//	return s.client.Patch(ctx, job, client.Apply, client.FieldOwner(constants.AdmissionName))
//}
