/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// BestEffortFIFO: Best effort first in, first out
	// order jobs by creation time, but don't block new jobs
	BestEffortFIFO QueueStrategy = "BestEffortFIFO"
)

type MiniClusterSpec struct {

	// Default size of the statefulset replias
	// +kubebuilder:default=1
	// +optional
	Size int32 `json:"size"`
}

// FluxSetupSpec defines the desired state of Flux
type FluxSetupSpec struct {
	// Run "make manifests" and "make" to regenerate code after modifying here

	MiniCluster MiniClusterSpec `json:"mini-cluster"`

	// +kubebuilder:default=20
	// +optional
	MaxSize int32 `json:"max-size"`

	// namespaces that are allowed to submit jobs to the queue
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`

	// QueueStrategy indicates the queueing strategy of the jobs
	//
	// +kubebuilder:default=BestEffortFIFO
	// +kubebuilder:validation:BestEffortFIFO
	QueueStrategy QueueStrategy `json:"queueingStrategy,omitempty"`
}

// FluxSetupStatus defines the observed state of a FluxSetup
type FluxSetupStatus struct {

	// usedResources should be a more advanced data structure, but now is just a count
	// +kubebuilder:default=1
	// +optional
	UsedResources int32 `json:"usedResources"`

	// WaitingJobs is the number of FluxJob currently waiting (not running)
	// +optional
	WaitingJobs int32 `json:"pendingWorkloads"`

	// RunningJobs is the number of FluxJob currently running (not finished)
	// +optional
	RunningJobs int32 `json:"admittedWorkloads"`
	// Conditions show observations of current state.
	// +optional
	// Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// The Flux Host config is a ConfigMap with Hostanme data
type FluxHostConfig struct {
	Hostfile string `json:"hostfile"`
}

// SetDefaults ensures that empty settings are defined with defaults
func (s *FluxSetup) SetDefaults() {

	// If FluxSetup doesn't have a size, default to 1
	if s.Spec.MaxSize == 0 {
		s.Spec.MaxSize = 20
	}
	fmt.Printf("🤓 FluxSetup.MaxSize %d\n", (*s).Spec.MaxSize)
	fmt.Printf("🤓 FluxSetup.MiniCluster.Size %d\n", (*s).Spec.MiniCluster.Size)
}

// ConfigMap describes configuration options
type ConfigMap struct {
	// Data holds the configuration file data
	ConfigData string `json:"config"`
}

// Data returns a valid ConfigMap name
func (c *ConfigMap) Data() string {
	return c.ConfigData
}

type QueueStrategy string

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FluxSetup is the Schema for the fluxes setups API
type FluxSetup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FluxSetupSpec   `json:"spec,omitempty"`
	Status FluxSetupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FluxSetupList contains a list of Flux
type FluxSetupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FluxSetup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FluxSetup{}, &FluxSetupList{})
}
