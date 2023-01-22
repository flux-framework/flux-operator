/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	_ "embed"
)

//go:embed templates/broker.toml
var brokerConfigTemplate string

//go:embed templates/wait.sh
var waitToStartTemplate string

//go:embed templates/cert-generate.sh
var generateCertTemplate string

// WaitTemplate populates wait.sh
type WaitTemplate struct {
	FluxToken         string // Token to log into the UI, should be consistent across containers
	FluxUser          string // Username for Flux Restful API
	MainHost          string // Main host identifier
	FluxOptionFlags   string // Option flags
	Hosts             string // List of hosts
	Diagnostics       bool   // Run diagnostics instead of job?
	PreCommand        string // Custom commands, looked up by container identifier
	FluxRestfulBranch string // branch to clone Flux Restful from, defaults to main
	FluxRestfulPort   int32  // port to run flux restful on
	Cores             int32
	Tasks             int32
	Size              int32 // size of the Minicluster (nodes / pods in indexed jobs)

	// Logging Modes
	QuietMode    bool // Don't print additional output
	TimedMode    bool // Add times when appropriate
	FluxLogLevel int32
}

// CertTemplate populates cert-generate.sh
type CertTemplate struct {
	PreCommand string
}
