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
	MainHost          string // Main host identifier
	FluxOptionFlags   string // Option flags
	Hosts             string // List of hosts
	Diagnostics       bool   // Run diagnostics instead of job?
	PreCommand        string // Custom commands, looked up by container identifier
	FluxRestfulBranch string // branch to clone Flux Restful from, defaults to main
	FluxRestfulPort   int32  // port to run flux restful on
	ClusterSize       int32  // number of nodes in mini cluster, should be size
	TestMode          bool   // Don't print additional output
	Size              int32
	FluxLogLevel      int32
}

// CertTemplate populates cert-generate.sh
type CertTemplate struct {
	PreCommand string
}
