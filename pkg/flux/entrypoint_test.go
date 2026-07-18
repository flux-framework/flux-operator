/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package flux

import (
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"
)

// fluxCluster returns a minimal single flux-runner MiniCluster whose one
// container is customized by the caller.
func fluxCluster(c api.MiniClusterContainer) *api.MiniCluster {
	c.Image = "rockylinux:9"
	c.RunFlux = true
	return &api.MiniCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "flux-operator"},
		Spec: api.MiniClusterSpec{
			Size:       1,
			Containers: []api.MiniClusterContainer{c},
		},
	}
}

// waitScript renders the entrypoints and returns the wait-0 script.
func waitScript(t *testing.T, mc *api.MiniCluster) string {
	t.Helper()
	data, err := GenerateEntrypoints(mc)
	if err != nil {
		t.Fatalf("GenerateEntrypoints error: %v", err)
	}
	script, ok := data["wait-0"]
	if !ok {
		t.Fatalf("no wait-0 entrypoint generated (keys: %v)", keys(data))
	}
	return script
}

func keys(m map[string]string) []string {
	out := []string{}
	for k := range m {
		out = append(out, k)
	}
	return out
}

// A batch workload writes the content VERBATIM and runs it with flux batch.
// It must NOT split lines or wrap them in the old per-line flux submit form.
func TestBatchRendersVerbatim(t *testing.T) {
	batch := "#!/bin/bash\nflux submit -n2 app-one\nflux submit -n2 app-two\nflux queue idle"
	script := waitScript(t, fluxCluster(api.MiniClusterContainer{Batch: batch}))

	if !strings.Contains(script, "flux batch") {
		t.Errorf("batch script should invoke 'flux batch'; got:\n%s", script)
	}
	if !strings.Contains(script, "flux-job.batch") {
		t.Errorf("batch content should be written to flux-job.batch")
	}
	// verbatim content present
	for _, line := range []string{"flux submit -n2 app-one", "flux submit -n2 app-two", "flux queue idle"} {
		if !strings.Contains(script, line) {
			t.Errorf("batch content line missing from script: %q", line)
		}
	}
	// the OLD auto-wrapping must be gone
	if strings.Contains(script, "--flags waitable --error=${FLUX_OUTPUT_DIR}/job-") {
		t.Errorf("batch should NOT auto-wrap lines in flux submit --flags waitable (removed behavior)")
	}
}

// A script workload is written to a file and run as the command (not batch).
func TestScriptWritesFileAndRunsAsCommand(t *testing.T) {
	body := "echo hello\necho goodbye"
	script := waitScript(t, fluxCluster(api.MiniClusterContainer{Script: body}))

	if !strings.Contains(script, "custom-entrypoint.sh") {
		t.Errorf("script should be written to a custom-entrypoint file; got:\n%s", script)
	}
	if !strings.Contains(script, "echo hello") || !strings.Contains(script, "echo goodbye") {
		t.Errorf("script content missing from generated entrypoint")
	}
	if !strings.Contains(script, `command="/bin/bash /tmp/custom-entrypoint.sh"`) {
		t.Errorf("script should set command to run the written file")
	}
	// script is not a batch job
	if strings.Contains(script, "flux-job.batch") {
		t.Errorf("script must not render as a batch job")
	}
}

// The top-level script takes precedence over the (legacy) commands.script.
func TestTopLevelScriptBeatsCommandsScript(t *testing.T) {
	c := api.MiniClusterContainer{Script: "echo TOPLEVEL"}
	c.Commands.Script = "echo LEGACY"
	script := waitScript(t, fluxCluster(c))
	if !strings.Contains(script, "echo TOPLEVEL") {
		t.Errorf("top-level script should be used")
	}
	if strings.Contains(script, "echo LEGACY") {
		t.Errorf("commands.script should be overridden by top-level script")
	}
}

// A plain command still renders through the submit path and is not written to
// a batch or script file (guards that command remains first-class).
func TestCommandStillSubmits(t *testing.T) {
	script := waitScript(t, fluxCluster(api.MiniClusterContainer{Command: "hostname"}))
	if strings.Contains(script, "flux-job.batch") {
		t.Errorf("command must not render as a batch job")
	}
	if strings.Contains(script, "custom-entrypoint.sh") {
		t.Errorf("command must not be written to a script file")
	}
}
