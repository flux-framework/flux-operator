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

// MiniCluster defines the desired state of a Flux MiniCluster
// "I am a Flux user and I want to launch a MiniCluster for my job!"
// A MiniCluster corresponds to a Batch Job -> StatefulSet + ConfigMaps
// A "task" within that cluster is flux running something.
type MiniClusterSpec struct {
	// Important: Run "make" and "make manifests" to regenerate code after modifying this file

	// Container image must contain flux and flux-sched install
	// +kubebuilder:default="fluxrm/flux-sched:focal"
	Image string `json:"image"`

	// Size (number of jobs to run)
	// +kubebuilder:default=1
	// +optional
	Size int32 `json:"size"`

	// Working directory to run command from
	// +optional
	WorkingDir string `json:"workingDir"`

	// Run flux diagnostics on start instead of command
	// +optional
	Diagnostics bool `json:"diagnostics"`

	// Should the job be limited to a particular number of seconds?
	// Approximately one year. This cannot be zero or job won't start
	// +kubebuilder:default=31500000
	// +optional
	DeadlineSeconds int64 `json:"deadlineSeconds"`

	// Allow the user to dictate pulling
	// By default we pull if not present. Setting
	// this to true will indicate to pull always
	// +kubebuilder:default=false
	// +optional
	PullAlways bool `json:"pullAlways"`

	// Allow the user to pull authenticated images
	// By default no secret is selected. Setting
	// this with the name of an already existing
	// imagePullSecret will specify that secret
	// in the pod spec.
	// +optional
	ImagePullSecrets []string `json:"imagePullSecrets"`

	// Single user executable to provide to flux start
	// +optional
	Command string `json:"command"`
}

// MiniClusterStatus defines the observed state of Flux
type MiniClusterStatus struct {

	// The JobUid is set internally to associate to a miniCluster
	JobId string `json:"jobid"`

	// conditions hold the latest Flux Job and MiniCluster states
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MiniCluster is the Schema for a Flux job launcher on K8s
type MiniCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MiniClusterSpec   `json:"spec,omitempty"`
	Status MiniClusterStatus `json:"status,omitempty"`
}

// SetDefaults ensures that empty settings are defined with defaults
func (f *MiniCluster) PrintDefaults() {
	fmt.Println()
	fmt.Printf("🤓 MiniCluster.Image %s\n", f.Spec.Image)
	fmt.Printf("🤓 MiniCluster.DeadlineSeconds %d\n", f.Spec.DeadlineSeconds)
	fmt.Printf("🤓 MiniCluster.Command %s\n", f.Spec.Command)
	fmt.Printf("🤓 MiniCluster.Size %s\n", fmt.Sprint(f.Spec.Size))
}

//+kubebuilder:object:root=true

// MiniClusterList contains a list of MiniCluster
type MiniClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MiniCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MiniCluster{}, &MiniClusterList{})
}
