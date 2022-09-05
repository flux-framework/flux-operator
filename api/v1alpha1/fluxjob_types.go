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

// FluxJobSpec defines the desired state of Flux
type FluxJobSpec struct {
	// Important: Run "make" and "make manifests" to regenerate code after modifying this file

	// Container image must contain flux and flux-sched install
	// +kubebuilder:default="fluxrm/flux-sched:focal"
	Image string `json:"image"`

	// Size (number of jobs to run)
	// +kubebuilder:default=1
	// +optional
	Size int32 `json:"size"`

	// Single user executable to provide to flux start
	// +optional
	Command string `json:"command"`
}

// FluxJobStatus defines the observed state of Flux
type FluxJobStatus struct {

	// The JobUid is set internally to associate to a miniCluster
	JobId string `json:"jobid"`

	// conditions hold the latest Flux Job and MiniCluster states
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FluxJob is the Schema for the fluxes API
type FluxJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FluxJobSpec   `json:"spec,omitempty"`
	Status FluxJobStatus `json:"status,omitempty"`
}

// SetDefaults ensures that empty settings are defined with defaults
func (f *FluxJob) SetDefaults() {
	fmt.Println()
	fmt.Printf("🤓 FluxJob.Image %s\n", f.Spec.Image)
	fmt.Printf("🤓 FluxJob.Command %s\n", f.Spec.Command)
	fmt.Printf("🤓 FluxJob.Size %s\n", fmt.Sprint(f.Spec.Size))
}

//+kubebuilder:object:root=true

// FluxJobList contains a list of Flux
type FluxJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FluxJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FluxJob{}, &FluxJobList{})
}
