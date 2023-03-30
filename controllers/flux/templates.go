/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	api "flux-framework/flux-operator/api/v1alpha1"

	_ "embed"
)

//go:embed templates/broker.toml
var brokerConfigTemplate string

//go:embed templates/job-manager.toml
var brokerConfigJobManagerPlugin string

//go:embed templates/archive.toml
var brokerArchiveSection string

//go:embed templates/wait.sh
var waitToStartTemplate string

//go:embed templates/cert-generate.sh
var generateCertTemplate string

// WaitTemplate populates wait.sh
type WaitTemplate struct {
	FluxToken string // Token to log into the UI, should be consistent across containers
	FluxUser  string // Username for Flux Restful API
	MainHost  string // Main host identifier
	Hosts     string // List of hosts
	Cores     int32
	Container api.MiniClusterContainer
	Spec      api.MiniClusterSpec

	// Batch commands split up
	Batch []string
}

// CertTemplate populates cert-generate.sh
type CertTemplate struct {
	PreCommand string
}
