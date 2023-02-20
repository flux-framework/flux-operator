/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// MiniCluster is an HPC cluster in Kubernetes you can control
// Either to submit a single job (and go away) or for a persistent single- or multi- user cluster

type MiniClusterSpec struct {
	// Important: Run "make" and "make manifests" to regenerate code after modifying this file

	// Containers is one or more containers to be created in a pod.
	// There should only be one container to run flux with runFlux
	// +listType=atomic
	Containers []MiniClusterContainer `json:"containers"`

	// Users of the MiniCluster
	// +optional
	// +listType=atomic
	Users []MiniClusterUser `json:"users"`

	// Labels for the job
	// +optional
	JobLabels map[string]string `json:"jobLabels"`

	// Volumes accessible to containers from a host
	// +optional
	Volumes map[string]MiniClusterVolume `json:"volumes"`

	// Logging modes determine the output you see in the job log
	// +optional
	Logging LoggingSpec `json:"logging"`

	// Customization to Flux Restful API
	// There should only be one container to run flux with runFlux
	// +optional
	FluxRestful FluxRestful `json:"fluxRestful"`

	// Cleanup the pods and storage when the index broker pod is complete
	// +kubebuilder:default=false
	// +optional
	Cleanup bool `json:"cleanup"`

	// Size (number of job pods to run, size of minicluster in pods)
	// +kubebuilder:default=1
	// +optional
	Size int32 `json:"size"`

	// Total number of CPUs being run across entire cluster
	// +kubebuilder:default=1
	// +optional
	Tasks int32 `json:"tasks"`

	// Should the job be limited to a particular number of seconds?
	// Approximately one year. This cannot be zero or job won't start
	// +kubebuilder:default=31500000
	// +optional
	DeadlineSeconds int64 `json:"deadlineSeconds"`

	// Pod spec details
	// +optional
	Pod PodSpec `json:"pod"`
}

type MiniClusterUser struct {

	// If a user is defined, the username is required
	Name string `json:"name"`

	// +optional
	Password string `json:"password"`
}

type LoggingSpec struct {

	// Quiet mode silences all output so the job only shows the test running
	// +kubebuilder:default=false
	// +optional
	Quiet bool `json:"quiet"`

	// Strict mode ensures any failure will not continue in the job entrypoint
	// +kubebuilder:default=true
	// +optional
	Strict bool `json:"strict"`

	// Debug mode adds extra verbosity to Flux
	// +kubebuilder:default=false
	// +optional
	Debug bool `json:"debug"`

	// Timed mode adds timing to Flux commands
	// +kubebuilder:default=false
	// +optional
	Timed bool `json:"timed"`
}

// PodSpec controlls variables for the cluster pod
type PodSpec struct {

	// Annotations for each pod
	// +optional
	Annotations map[string]string `json:"annotations"`

	// Labels for each pod
	// +optional
	Labels map[string]string `json:"labels"`

	// Resources include limits and requests
	// +optional
	Resources ContainerResource `json:"resources"`
}

// MiniClusterStatus defines the observed state of Flux
type MiniClusterStatus struct {

	// The Jobid is set internally to associate to a miniCluster
	// This isn't currently in use, we only have one!
	Jobid string `json:"jobid"`

	// conditions hold the latest Flux Job and MiniCluster states
	// +listType=atomic
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type FluxRestful struct {

	// Branch to clone Flux Restful API from
	// +kubebuilder:default="main"
	// +optional
	Branch string `json:"branch"`

	// Port to run Flux Restful Server On
	// +kubebuilder:default=5000
	// +optional
	Port int32 `json:"port"`

	// These two should not actually be set by a user,
	// but rather generated by tools and provided
	// Username to use for RestFul API
	// +optional
	Username string `json:"username"`

	// Token to use for RestFul API
	// +optional
	Token string `json:"token"`
}

// Mini Cluster local volumes available to mount (these are on the host)
type MiniClusterVolume struct {
	Path string `json:"path"`

	// +optional
	Labels map[string]string `json:"labels"`

	// Annotations for persistent volume claim
	// +optional
	Annotations map[string]string `json:"annotations"`

	// Optional volume attributes
	// +optional
	Attributes map[string]string `json:"attributes"`

	// Volume handle, falls back to storage class name
	// if not defined
	// +optional
	VolumeHandle string `json:"volumeHandle"`

	// +kubebuilder:default="hostpath"
	// +optional
	Class string `json:"class"`

	// +kubebuilder:default="pvc-storage-class-name-unset"
	// +optional
	PVCStorageClassName string `json:"PVCclass"`

	// Storage driver, e.g., gcs.csi.ofek.dev
	// Only needed if not using hostpath
	// +optional
	Driver string `json:"driver"`

	// Secret reference in Kubernetes with service account role
	// +optional
	Secret string `json:"secret"`

	// Secret namespace
	// +kubebuilder:default="default"
	// +optional
	SecretNamespace string `json:"secretNamespace"`

	// Capacity (string) for PVC (storage request) to create PV
	// +kubebuilder:default="5Gi"
	// +optional
	Capacity string `json:"capacity"`
}

// A Container volume must reference one defined for the MiniCluster
// The path here is in the container
type ContainerVolume struct {
	Path string `json:"path"`

	// +kubebuilder:default=false
	// +optional
	ReadOnly bool `json:"readOnly"`
}

type MiniClusterContainer struct {

	// Container image must contain flux and flux-sched install
	// +kubebuilder:default="fluxrm/flux-sched:focal"
	Image string `json:"image"`

	// Container name is only required for non flux runners
	// +optional
	Name string `json:"name"`

	// Cores the container should use
	// +optional
	Cores int32 `json:"cores"`

	// Working directory to run command from
	// +optional
	WorkingDir string `json:"workingDir"`

	// Run flux diagnostics on start instead of command
	// +optional
	Diagnostics bool `json:"diagnostics"`

	// Flux User, if created in the container
	// +optional
	FluxUser FluxUser `json:"fluxUser"`

	// Ports to be exposed to other containers in the cluster
	// We take a single list of integers and map to the same
	// +optional
	// +listType=atomic
	Ports []int32 `json:"ports"`

	// Key/value pairs for the environment
	// +optional
	Environment map[string]string `json:"environment"`

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
	// +optional
	RunFlux bool `json:"runFlux"`

	// Volumes that can be mounted (must be defined in volumes)
	// +optional
	Volumes map[string]ContainerVolume `json:"volumes"`

	// Flux option flags, usually provided with -o
	// optional - if needed, default option flags for the server
	// These can also be set in the user interface to override here.
	// This is only valid for a FluxRunner "runFlux" true
	// +optional
	FluxOptionFlags string `json:"fluxOptionFlags"`

	// Log level to use for flux logging (only in non TestMode)
	// +kubebuilder:default=6
	// +optional
	FluxLogLevel int32 `json:"fluxLogLevel"`

	// Special command to run at beginning of script, directly after asFlux
	// is defined as sudo -u flux -E (so you can change that if desired.)
	// This is only valid if FluxRunner is set (that writes a wait.sh script)
	// This is for the indexed job pods and the certificate generation container.
	// +optional
	PreCommand string `json:"preCommand"`

	// Lifecycle can handle post start commands, etc.
	// +optional
	LifeCycle LifeCycle `json:"lifeCycle"`

	// Resources include limits and requests
	// +optional
	Resources ContainerResources `json:"resources"`

	// More specific or detailed commands for just workers/broker
	// +optional
	Commands Commands `json:"commands"`
}

type LifeCycle struct {

	// +optional
	PostStartExec string `json:"postStartExec"`
}

type FluxUser struct {

	// Flux user name
	// +kubebuilder:default="flux"
	// +optional
	Name string `json:"name"`

	// UID for the FluxUser
	// +optional
	// +kubebuilder:default=1000
	Uid int `json:"uid"`
}

type Commands struct {

	// Run flux start as root - required for some storage binds
	// +kubebuilder:default=false
	// +optional
	RunFluxAsRoot bool `json:"runFluxAsRoot"`

	// pre command is run after global PreCommand, before anything else
	// +optional
	Pre string `json:"pre"`
}

// ContainerResources include limits and requests
type ContainerResources struct {

	// +optional
	Limits ContainerResource `json:"limits"`

	// +optional
	Requests ContainerResource `json:"requests"`
}

type ContainerResource map[string]intstr.IntOrString

// MiniCluster is the Schema for a Flux job launcher on K8s

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
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
		f.Spec.Containers[0].RunFlux = true
	}

	// If pod and job labels aren't defined, create label set
	if f.Spec.Pod.Labels == nil {
		f.Spec.Pod.Labels = map[string]string{}
	}
	if f.Spec.Pod.Annotations == nil {
		f.Spec.Pod.Annotations = map[string]string{}
	}
	if f.Spec.JobLabels == nil {
		f.Spec.JobLabels = map[string]string{}
	}

	// Validate user passwords. If provided, need to be 8 or fewer characters
	for _, user := range f.Spec.Users {
		if user.Password != "" && len(user.Password) > 8 {
			fmt.Printf("😥️ %s has a password that is too long, can be no longer than 8 characters\n", user.Name)
			return false
		}
	}
	for i, container := range f.Spec.Containers {
		name := fmt.Sprintf("MiniCluster.Container.%d", i)
		fmt.Printf("🤓 %s.Image %s\n", name, container.Image)
		fmt.Printf("🤓 %s.Command %s\n", name, container.Command)
		fmt.Printf("🤓 %s.FluxRunner %t\n", name, container.RunFlux)

		// Count the FluxRunners
		if container.RunFlux {
			fluxRunners += 1

			// Non flux-runners are required to have a name
		} else {
			if container.Name == "" {
				fmt.Printf("😥️ %s is missing a name\n", name)
				return false
			}
		}

		// If we have volumes defined, they must be provided in the global
		// volumes for the MiniCluster CRD.
		for key, _ := range container.Volumes {
			_, found := f.Spec.Volumes[key]
			if !found {
				fmt.Printf("😥️ %s defines a named volume %s but it is not defined for the MiniCluster\n", name, key)
				return false
			}
		}
	}
	if fluxRunners != 1 {
		valid = false
	}

	// For each volume, if it's not a hostvolume, we require a secret reference
	for key, volume := range f.Spec.Volumes {
		if volume.Class != "hostpath" && volume.Secret == "" {
			fmt.Printf("😥️ Found non-hostpath volume %s that is missing a secret\n", key)
			valid = false
		}
	}

	return valid
}

//+kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MiniClusterList contains a list of MiniCluster
type MiniClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MiniCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MiniCluster{}, &MiniClusterList{})
}
