/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package flux

import (
	"bytes"
	"fmt"
	"strings"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"
)

// GenerateEntrypoints generates the data structure (for config map) with entrypoint scripts
func GenerateEntrypoints(cluster *api.MiniCluster) (map[string]string, error) {
	data := map[string]string{}

	// Each application container has a wait script that waits for flux to be ready
	for i, container := range cluster.Spec.Containers {
		if container.RunFlux {
			waitScriptID := fmt.Sprintf("wait-%d", i)

			waitScript, err := generateEntrypointScript(cluster, i)
			if err != nil {
				return data, err
			}
			data[waitScriptID] = waitScript
		}

		// Custom logic for a sidecar container alongside flux
		if container.GenerateEntrypoint() {
			startScriptID := fmt.Sprintf("start-%d", i)
			startScript, err := generateServiceEntrypoint(cluster, container)
			if err != nil {
				return data, err
			}
			data[startScriptID] = startScript
		}
	}
	// Main flux entrypoint for flux-view generation
	script, err := GenerateFluxEntrypoint(cluster)
	if err != nil {
		return data, err
	}
	data[cluster.Spec.Flux.Container.Name] = script

	// Add the curve.cert
	curveCert, err := GetCurveCert(cluster)
	data["curve.cert"] = curveCert
	return data, err
}

// generateServiceEntrypoint generates an entrypoint for a service container
func generateServiceEntrypoint(cluster *api.MiniCluster, container api.MiniClusterContainer) (string, error) {
	st := ServiceTemplate{
		ViewBase:  cluster.Spec.Flux.Container.MountPath,
		Container: container,
		Spec:      cluster.Spec,
	}

	// Wrap the named template to identify it later
	startTemplate := `{{define "start"}}` + sidecarContainerTemplate + "{{end}}"

	// We assemble different strings (including the components) into one!
	t, err := combineTemplates(startComponents, startTemplate)
	if err != nil {
		return "", err
	}
	var output bytes.Buffer
	if err := t.ExecuteTemplate(&output, "start", st); err != nil {
		return "", err
	}
	return output.String(), nil

}

// generateEntrypointScript generates an entrypoint script to start everything up!
func generateEntrypointScript(
	cluster *api.MiniCluster,
	containerIndex int,
) (string, error) {

	container := cluster.Spec.Containers[containerIndex]

	// Ensure if we have a batch command, it gets split up
	batchCommand := strings.Split(container.Command, "\n")

	// Required quorum - might be smaller than initial list if size != maxsize
	// This used to be a range, now it's a size
	requiredRanks := getRequiredRanks(cluster)

	// The token uuid is the same across images
	wt := WaitTemplate{
		RequiredRanks: requiredRanks,
		ViewBase:      cluster.Spec.Flux.Container.MountPath,
		Container:     container,
		MainHost:      cluster.MainHost(),
		Spec:          cluster.Spec,
		Batch:         batchCommand,
	}

	// Wrap the named template to identify it later
	startTemplate := `{{define "start"}}` + waitToStartTemplate + "{{end}}"

	// We assemble different strings (including the components) into one!
	t, err := combineTemplates(startComponents, startTemplate)
	if err != nil {
		return "", err
	}
	var output bytes.Buffer
	if err := t.ExecuteTemplate(&output, "start", wt); err != nil {
		return "", err
	}
	return output.String(), nil
}
