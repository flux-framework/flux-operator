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
	"errors"
	api "flux-framework/flux-operator/api/v1alpha1"
	jobctrl "flux-framework/flux-operator/pkg/job"
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Errors
var (
	errSetupAlreadyExists = errors.New("setup already exists")
)

type Manager struct {
	sync.RWMutex
	cond sync.Cond

	client client.Client

	// Currently we want one setup == one queue
	queue Queue
}

func NewManager(client client.Client) *Manager {
	m := &Manager{
		client: client,
	}
	m.cond.L = &m.RWMutex
	return m
}

// AddOrUpdateJob adds or updates a job, triggered by the FluxJob creation
func (m *Manager) AddOrUpdateJob(job *api.FluxJob) bool {
	m.Lock()
	defer m.Unlock()
	return m.addOrUpdateJob(job)
}

// AddOrUpdateJob adds or updates a job, triggered by the FluxJob creation
func (m *Manager) IsRunningJob(job *api.FluxJob) bool {
	m.Lock()
	defer m.Unlock()
	info := jobctrl.NewInfo(job)
	return m.queue.IsRunningJob(info)
}

func (m *Manager) IsWaitingJob(job *api.FluxJob) bool {
	m.Lock()
	defer m.Unlock()
	info := jobctrl.NewInfo(job)
	return m.queue.IsWaitingJob(info)
}

func (m *Manager) JobsPending() int {
	m.RLock()
	defer m.RUnlock()
	return m.queue.Pending()
}

func (m *Manager) JobsRunning() int {
	m.RLock()
	defer m.RUnlock()
	return m.queue.Running()
}

func (m *Manager) Delete(job *api.FluxJob) bool {
	m.RLock()
	defer m.RUnlock()
	info := jobctrl.NewInfo(job)
	return m.queue.Delete(info)
}

func (m *Manager) addOrUpdateJob(job *api.FluxJob) bool {

	info := jobctrl.NewInfo(job)

	// If we don't have a setup yet
	if m.queue == nil {
		return false
	}

	// This always returns true
	m.queue.PushOrUpdate(info)

	// TODO report pending jobs here
	// TODO what does broadcast do?
	m.Broadcast()
	return true
}

func (m *Manager) HasQueue() bool {
	return m.queue != nil
}

func (m *Manager) InitQueue(ctx context.Context, setup *api.FluxSetup) error {
	m.Lock()
	defer m.Unlock()

	// We already have a queue
	if m.HasQueue() {
		return errSetupAlreadyExists
	}

	// Get a new queue based on the strategy defined in the FluxSetup
	// custom resource definition. Currently we have Best effort FIFP
	createdQueue, err := newQueue(setup)
	if err != nil {
		return err
	}

	// Store the queue namespaced by the setup for now
	m.queue = createdQueue
	queued := m.QueueWaitingJobs(ctx)

	// TODO report pending jbos to some metric server here
	if queued {
		m.Broadcast()
	}
	return nil
}

// CleanUpOnContext tracks the context. When closed, it wakes routines waiting
// on elements to be available. It should be called before doing any calls to
// Heads.
func (m *Manager) CleanUpOnContext(ctx context.Context) {
	<-ctx.Done()
	m.Broadcast()
}

// QueueWaitingJobs can be called on init or by the scheduler
func (m *Manager) QueueWaitingJobs(ctx context.Context) bool {
	// Ensuring waiting jobs are added to the heap
	queued := m.queue.QueueWaitingJobs(ctx, m.client)
	return queued
}

// Awake go routines waiting on the condition
func (m *Manager) Broadcast() {
	m.cond.Broadcast()
}
