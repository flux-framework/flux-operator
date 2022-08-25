/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FluxSetupSpec defines the desired state of Flux
type FluxSetupSpec struct {
	// Run "make manifests" and "make" to regenerate code after modifying here

	// Size of the statefulset replias
	// +optional
	Size int32 `json:"size"`

	// TODO these three things should eventually be properly mapped into go structs
	// CurveCert is just a placeholder for what eventually will be done by the operator
	Hostfile ConfigMap `json:"hostfile"`

	// CurveCert is just a placeholder for what eventually will be done by the operator
	Cert ConfigMap `json:"cert"`

	// Broker describes the broker config (and any other attributes we can eventually add)
	Broker ConfigMap `json:"broker"`
}

// FluxSetupStatus defines the observed state of a FluxSetup
type FluxSetupStatus struct {
}

// SetDefaults ensures that empty settings are defined with defaults
func (s *FluxSetup) SetDefaults() {

	// If FluxSetup doesn't have a size, default to 1
	if s.Spec.Size == 0 {
		s.Spec.Size = 1
	}
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
