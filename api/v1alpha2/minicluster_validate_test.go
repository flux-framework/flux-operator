/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package v1alpha2

import "testing"

// one returns a single-container MiniCluster (RunFlux is auto-set in Validate).
func one(c MiniClusterContainer) *MiniCluster {
	c.Image = "rockylinux:9"
	return &MiniCluster{
		Spec: MiniClusterSpec{
			Size:       1,
			Containers: []MiniClusterContainer{c},
		},
	}
}

// command, script, and batch are three ways to define the same workload; at
// most one may be set. Each alone is valid; any pair is not.
func TestValidateWorkloadExclusivity(t *testing.T) {
	cases := []struct {
		name string
		c    MiniClusterContainer
		want bool
	}{
		{"command only", MiniClusterContainer{Command: "hostname"}, true},
		{"script only", MiniClusterContainer{Script: "echo hi"}, true},
		{"batch only", MiniClusterContainer{Batch: "flux submit hostname"}, true},
		{"none (interactive)", MiniClusterContainer{}, true},
		{"command+script", MiniClusterContainer{Command: "hostname", Script: "echo hi"}, false},
		{"command+batch", MiniClusterContainer{Command: "hostname", Batch: "flux submit hostname"}, false},
		{"script+batch", MiniClusterContainer{Script: "echo hi", Batch: "flux submit hostname"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := one(tc.c).Validate(); got != tc.want {
				t.Errorf("Validate() = %v, want %v for %s", got, tc.want, tc.name)
			}
		})
	}
}

// The legacy commands.script counts as "script" for exclusivity.
func TestValidateLegacyCommandsScriptConflictsWithCommand(t *testing.T) {
	c := MiniClusterContainer{Command: "hostname"}
	c.Commands.Script = "echo hi"
	if one(c).Validate() {
		t.Errorf("command + commands.script should be invalid")
	}
}

// batch and launcher are mutually exclusive.
func TestValidateBatchLauncherExclusive(t *testing.T) {
	if one(MiniClusterContainer{Batch: "flux submit hostname", Launcher: true}).Validate() {
		t.Errorf("batch + launcher should be invalid")
	}
	// launcher with a command remains valid
	if !one(MiniClusterContainer{Command: "snakemake", Launcher: true}).Validate() {
		t.Errorf("command + launcher should be valid")
	}
}
