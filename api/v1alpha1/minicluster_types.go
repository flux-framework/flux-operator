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

	// Containers is one or more containers to be created in a pod.
	// There should only be one container to run flux with runFlux
	Containers []MiniClusterContainer `json:"containers"`

	// Test mode silences all output so the job only shows the test running
	// +kubebuilder:default=false
	// +optional
	TestMode bool `json:"test"`

	// Customization to Flux Restful API
	// There should only be one container to run flux with runFlux
	// +optional
	FluxRestful FluxRestful `json:"fluxRestful"`

	// Size (number of jobs to run)
	// +kubebuilder:default=1
	// +optional
	Size int32 `json:"size"`

	// Run flux diagnostics on start instead of command
	// +optional
	Diagnostics bool `json:"diagnostics"`

	// Should the job be limited to a particular number of seconds?
	// Approximately one year. This cannot be zero or job won't start
	// +kubebuilder:default=31500000
	// +optional
	DeadlineSeconds int64 `json:"deadlineSeconds"`

	// localDeploy should be true for development, or deploying in the
	// case that there isn't an actual kubernetes cluster (e.g., you
	// are not using make deploy. It uses a persistent volume instead of
	// a claim
	// +kubebuilder:default=false
	// +optional
	LocalDeploy bool `json:"localDeploy"`
}

// MiniClusterStatus defines the observed state of Flux
type MiniClusterStatus struct {

	// The JobUid is set internally to associate to a miniCluster
	JobId string `json:"jobid"`

	// conditions hold the latest Flux Job and MiniCluster states
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type FluxRestful struct {

	// Branch to clone Flux Restful API from
	// +kubebuilder:default="main"
	// +optional
	Branch string `json:"branch"`
}

type MiniClusterContainer struct {

	// Container image must contain flux and flux-sched install
	// +kubebuilder:default="fluxrm/flux-sched:focal"
	Image string `json:"image"`

	// Working directory to run command from
	// +optional
	WorkingDir string `json:"workingDir"`

	// Allow the user to pull authenticated images
	// By default no secret is selected. Setting
	// this with the name of an already existing
	// imagePullSecret will specify that secret
	// in the pod spec.
	// +optional
	ImagePullSecret string `json:"imagePullSecret"`

	// Single user executable to provide to flux start
	// IMPORTANT: This is left here, but not used in
	// favor of exposing Flux via a Restful API. We
	// Can remove this when that is finalized.
	// +optional
	Command string `json:"command"`

	// Allow the user to dictate pulling
	// By default we pull if not present. Setting
	// this to true will indicate to pull always
	// +kubebuilder:default=false
	// +optional
	PullAlways bool `json:"pullAlways"`

	// Main container to run flux (only should be one)
	// +kubebuilder:default=true
	// +optional
	FluxRunner bool `json:"runFlux"`

	// Flux option flags, usually provided with -o
	// optional - if needed, default option flags for the server
	// These can also be set in the user interface to override here.
	// This is only valid for a FluxRunner
	// +optional
	FluxOptionFlags string `json:"fluxOptionFlags"`

	// Special command to run at beginning of script, directly after asFlux
	// is defined as sudo -u flux -E (so you can change that if desired.)
	// This is only valid if FluxRunner is set (that writes a wait.sh script)
	// +optional
	PreCommand string `json:"preCommand"`

	// Lifecycle can handle post start commands, etc.
	// +optional
	LifeCyclePostStartExec string `json:"postStartExec"`
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

// Validate ensures we have data that is needed, and sets defaults if needed
func (f *MiniCluster) Validate() bool {
	fmt.Println()

	// Global (entire cluster) settings
	fmt.Printf("🤓 MiniCluster.DeadlineSeconds %d\n", f.Spec.DeadlineSeconds)
	fmt.Printf("🤓 MiniCluster.Size %s\n", fmt.Sprint(f.Spec.Size))

	// We should have only one flux runner
	valid := true
	fluxRunners := 0

	// If we only have one container, assume we want to run flux with it
	// This makes it easier for the user to not require the flag
	if len(f.Spec.Containers) == 1 {
		f.Spec.Containers[0].FluxRunner = true
	}

	for i, container := range f.Spec.Containers {
		name := fmt.Sprintf("MiniCluster.Container.%d", i)
		fmt.Printf("🤓 %s.Image %s\n", name, container.Image)
		fmt.Printf("🤓 %s.Command %s\n", name, container.Command)
		fmt.Printf("🤓 %s.FluxRunner %t\n", name, container.FluxRunner)
		if container.FluxRunner {
			fluxRunners += 1
		}
	}
	if fluxRunners != 1 {
		valid = false
	}
	return valid
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
