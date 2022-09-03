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

	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Errors
var (
	errSetupAlreadyExists = errors.New("setup already exists")
	errSetupDoesNotExist  = errors.New("setup doesn't exist")
)

type Manager struct {
	sync.RWMutex
	cond sync.Cond

	client client.Client
	queues map[string]Queue

	// Key is cohort's name. Value is a set of associated Queue names.
	cohorts map[string]sets.String
}

func NewManager(client client.Client) *Manager {
	m := &Manager{
		client:  client,
		queues:  make(map[string]Queue),
		cohorts: make(map[string]sets.String),
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

func (m *Manager) Pending(setup *api.FluxSetup) int {
	m.RLock()
	defer m.RUnlock()
	return m.queues[setup.Name].Pending()
}

func (m *Manager) addOrUpdateJob(job *api.FluxJob) bool {

	info := jobctrl.NewInfo(job)

	// If we don't have a setup yet
	if len(m.queues) == 0 {
		return false
	}
	// Grab the first one
	for _, q := range m.queues {
		// This always returns true
		q.PushOrUpdate(info)
		break
	}
	// TODO report pending jobs here
	m.Broadcast()
	return true
}

func (m *Manager) InitQueue(ctx context.Context, setup *api.FluxSetup) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.queues[setup.Name]; ok {
		return errSetupAlreadyExists
	}

	// Get a new queue based on the strategy defined in the FluxSetup
	// custom resource definition. Currently we have Best effort FIFP
	createdQueue, err := newQueue(setup)
	if err != nil {
		return err
	}

	// Store the queue namespaced by the setup for now
	m.queues[setup.Name] = createdQueue

	// Ensuring waiting jobs are added to the heap
	queued := createdQueue.QueueWaitingJobs(ctx, m.client)

	// TODO report pending jbos to some metric server here
	if queued {
		m.Broadcast()
	}
	return nil
}

// Awake go routines waiting on the condition
func (m *Manager) Broadcast() {
	m.cond.Broadcast()
}
