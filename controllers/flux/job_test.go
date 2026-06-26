/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"
)

// miniCluster returns a minimal valid MiniCluster for builder tests.
func miniCluster(size int32) *api.MiniCluster {
	return &api.MiniCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "flux-operator"},
		Spec: api.MiniClusterSpec{
			Size: size,
			Containers: []api.MiniClusterContainer{
				{Image: "rockylinux:9"},
			},
		},
	}
}

func TestNewMiniClusterJob_Basics(t *testing.T) {
	mc := miniCluster(4)
	job, err := NewMiniClusterJob(mc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.Name != "test" || job.Namespace != "flux-operator" {
		t.Fatalf("job identity wrong: %s/%s", job.Namespace, job.Name)
	}
	if job.Spec.Completions == nil || *job.Spec.Completions != 4 {
		t.Fatalf("completions should equal size 4, got %v", job.Spec.Completions)
	}
	if job.Spec.Parallelism == nil || *job.Spec.Parallelism != 4 {
		t.Fatalf("parallelism should equal size 4, got %v", job.Spec.Parallelism)
	}
}

func TestNewMiniClusterJob_SchedulerNamePropagates(t *testing.T) {
	mc := miniCluster(2)
	mc.Spec.Pod.SchedulerName = "my-scheduler"
	job, err := NewMiniClusterJob(mc)
	if err != nil {
		t.Fatal(err)
	}
	if got := job.Spec.Template.Spec.SchedulerName; got != "my-scheduler" {
		t.Fatalf("schedulerName not propagated: got %q", got)
	}
}

// A user can attach arbitrary labels and a custom scheduler name to the pods via
// spec.pod; the operator must propagate both verbatim. This is what lets an
// external scheduler select these pods (e.g. by a group label it reads) without
// the operator knowing anything about that scheduler. The label keys below are
// arbitrary examples — the operator is agnostic to their meaning.
func TestNewMiniClusterJob_UserPodLabelsPassThrough(t *testing.T) {
	mc := miniCluster(4)
	mc.Spec.Pod.SchedulerName = "my-scheduler"
	mc.Spec.Pod.Labels = map[string]string{
		"example.com/group": "group-a",
		"my.org/custom":     "value",
	}
	job, err := NewMiniClusterJob(mc)
	if err != nil {
		t.Fatal(err)
	}
	labels := job.Spec.Template.ObjectMeta.Labels
	if labels["example.com/group"] != "group-a" {
		t.Fatalf("user group label not propagated to pods: %v", labels)
	}
	if labels["my.org/custom"] != "value" {
		t.Fatalf("arbitrary user label not propagated: %v", labels)
	}
}

func TestNewMiniClusterJob_Tolerations(t *testing.T) {
	mc := miniCluster(2)
	mc.Spec.Pod.Tolerations = []api.Toleration{
		{Key: "launcher", Operator: "Exists", Effect: "NoSchedule"},
	}
	job, err := NewMiniClusterJob(mc)
	if err != nil {
		t.Fatal(err)
	}
	tols := job.Spec.Template.Spec.Tolerations
	if len(tols) != 1 || string(tols[0].Key) != "launcher" || string(tols[0].Effect) != "NoSchedule" {
		t.Fatalf("toleration not propagated correctly: %+v", tols)
	}
}
