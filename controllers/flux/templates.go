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
	"fmt"
	"text/template"

	_ "embed"
)

//go:embed templates/broker.toml
var brokerConfigTemplate string

//go:embed templates/job-manager.toml
var brokerConfigJobManagerPlugin string

//go:embed templates/components.sh
var startComponents string

//go:embed templates/archive.toml
var brokerArchiveSection string

//go:embed templates/broker.sh
var brokerStartTemplate string

//go:embed templates/worker.sh
var workerStartTemplate string

// StartTemplate populates broker.sh or worker.sh
type StartTemplate struct {
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

// combineTemplates into one "start"
func combineTemplates(listing ...string) (t *template.Template, err error) {
	t = template.New("start")

	for i, templ := range listing {
		_, err = t.New(fmt.Sprint("_", i)).Parse(templ)
		if err != nil {
			return t, err
		}
	}
	return t, nil
}
