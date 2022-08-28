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

// FluxSpec defines the desired state of Flux
type FluxSpec struct {
	// Important: Run "make" and "make manifests" to regenerate code after modifying this file

	// Container image must contain flux and flux-sched install
	// This container is provided by the user via Flux, but is also passed
	// to the FluxSetup reconciler, which needs to run the same container image.
	// Likely these could be separated, but I'm not sure how that works yet.
	// +optional
	Image string `json:"image"`

	// Single user executable to provide to flux start
	Command string `json:"command"`
}

// FluxStatus defines the observed state of Flux
type FluxStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Flux is the Schema for the fluxes API
type Flux struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FluxSpec   `json:"spec,omitempty"`
	Status FluxStatus `json:"status,omitempty"`
}

// SetDefaults ensures that empty settings are defined with defaults
func (f *Flux) SetDefaults() {

	// Default container image to use
	if f.Spec.Image == "" {
		f.Spec.Image = "fluxrm/flux-sched:focal"
	}

	fmt.Println()
	fmt.Printf("🤓 Flux.Image %s\n", f.Spec.Image)
	fmt.Printf("🤓 Flux.Command %s\n", f.Spec.Command)
}

//+kubebuilder:object:root=true

// FluxList contains a list of Flux
type FluxList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Flux `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Flux{}, &FluxList{})
}
