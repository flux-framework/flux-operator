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
	"fmt"
	"text/template"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"
)

//go:embed templates/wait.sh
var waitToStartTemplate string

//go:embed templates/start.sh
var sidecarContainerTemplate string

//go:embed templates/components.sh
var startComponents string

// ServiceTemplate is for a separate service container
type ServiceTemplate struct {
	ViewBase       string // Where the mounted view with flux is expected to be
	Container      api.MiniClusterContainer
	ContainerIndex int
	Spec           api.MiniClusterSpec
}

// WaitTemplate populates wait.sh for an application container entrypoint
type WaitTemplate struct {
	ViewBase  string // Where the mounted view with flux is expected to be
	MainHost  string // Main host identifier
	FluxToken string // Token to log into the UI, should be consistent across containers
	Container api.MiniClusterContainer

	// Index for container, for generation of unique socket path
	ContainerIndex int
	Spec           api.MiniClusterSpec

	// Broker initial quorum that must be online to start
	// This is used if the cluster MaxSize > Size
	RequiredRanks string

	// Batch commands split up
	Batch []string
}

// combineTemplates into one common start
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
