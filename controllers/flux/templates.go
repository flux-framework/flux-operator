/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	api "github.com/flux-framework/flux-operator/api/v1alpha1"

	_ "embed"
)

//go:embed templates/broker.toml
var brokerConfigTemplate string

//go:embed templates/job-manager.toml
var brokerConfigJobManagerPlugin string

//go:embed templates/wait.sh
var waitToStartTemplate string

// WaitTemplate populates wait.sh
type WaitTemplate struct {
	FluxToken string // Token to log into the UI, should be consistent across containers
	FluxUser  string // Username for Flux Restful API
	MainHost  string // Main host identifier
	Hosts     string // List of hosts
	Cores     int32
	Container api.MiniClusterContainer
	Spec      api.MiniClusterSpec

	// Broker initial quorum that must be online to start
	// This is used if the cluster MaxSize > Size
	RequiredRanks string

	// Batch commands split up
	Batch []string
}

// BrokerTemplate defines the broker templates (broker.toml)
type BrokerTemplate struct {
	FluxInstallRoot string
	FQDN            string
	Spec            api.MiniClusterSpec
	ClusterName     string
	Hosts           string
}
