/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package flux

import (
	_ "embed"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"
)

//go:embed templates/wait.sh
var waitToStartTemplate string

//go:embed templates/start.sh
var sidecarContainerTemplate string

// ServiceTemplate is for a separate service container
type ServiceTemplate struct {
	Container api.MiniClusterContainer
	Spec      api.MiniClusterSpec
}

// WaitTemplate populates wait.sh for an application container entrypoint
type WaitTemplate struct {
	ViewBase  string // Where the mounted view with flux is expected to be
	MainHost  string // Main host identifier
	CurveCert string // curve certificate string
	FluxToken string // Token to log into the UI, should be consistent across containers
	Container api.MiniClusterContainer
	Spec      api.MiniClusterSpec

	// Broker initial quorum that must be online to start
	// This is used if the cluster MaxSize > Size
	RequiredRanks string

	// Batch commands split up
	Batch []string
}
