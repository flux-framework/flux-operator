/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MiniCluster defines a Flux MiniCluster
type MiniClusterSpec struct {
	// Important: Run "make" and "make manifests" to regenerate code after modifying this file
}

// MiniClusterStatus defines the observed state of Flux
type MiniClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MiniCluster is the Schema for a mini flux instance
type MiniCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MiniClusterSpec   `json:"spec,omitempty"`
	Status MiniClusterStatus `json:"status,omitempty"`
}

// SetDefaults ensures that empty settings are defined with defaults
func (f *MiniCluster) SetDefaults() {
}

//+kubebuilder:object:root=true

// MiniClusterList lists a mini cluster
type MiniClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MiniCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MiniCluster{}, &MiniClusterList{})
}
