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

	// Services are one or more service containers to bring up
	// alongside the MiniCluster.
	// +optional
	// +listType=atomic
	Services []MiniClusterContainer `json:"services"`

	// A spec for exposing or defining the cluster headless service
	//+optional
	Network Network `json:"network"`

	// Users of the MiniCluster
	// +optional
	// +listType=atomic
	Users []MiniClusterUser `json:"users"`

	// Labels for the job
	// +optional
	JobLabels map[string]string `json:"jobLabels"`

	// Run a single-user, interactive minicluster
	// +kubebuilder:default=false
	// +optional
	Interactive bool `json:"interactive"`

	// Flux options for the broker, shared across cluster
	// +optional
	Flux FluxSpec `json:"flux"`

	// Volumes accessible to containers from a host
	// Not all containers are required to use them
	// +optional
	Volumes map[string]MiniClusterVolume `json:"volumes"`

	// Logging modes determine the output you see in the job log
	// +optional
	Logging LoggingSpec `json:"logging"`

	// Archive to load or save
	// +optional
	Archive MiniClusterArchive `json:"archive"`

	// Customization to Flux Restful API
	// There should only be one container to run flux with runFlux
	// +optional
	FluxRestful FluxRestful `json:"fluxRestful"`

	// Cleanup the pods and storage when the index broker pod is complete
	// +kubebuilder:default=false
	// +default=false
	// +optional
	Cleanup bool `json:"cleanup,omitempty"`

	// Size (number of job pods to run, size of minicluster in pods)
	// This is also the minimum number required to start Flux
	// +kubebuilder:default=1
	// +default=1
	// +optional
	Size int32 `json:"size,omitempty"`

	// MaxSize (maximum number of pods to allow scaling to)
	// +optional
	MaxSize int32 `json:"maxSize,omitempty"`

	// Total number of CPUs being run across entire cluster
	// +kubebuilder:default=1
	// +default=1
	// +optional
	Tasks int32 `json:"tasks,omitempty"`

	// Should the job be limited to a particular number of seconds?
	// Approximately one year. This cannot be zero or job won't start
	// +kubebuilder:default=31500000
	// +default=31500000
	// +optional
	DeadlineSeconds int64 `json:"deadlineSeconds,omitempty"`

	// Pod spec details
	// +optional
	Pod PodSpec `json:"pod"`
}

type Network struct {

	// Name for cluster headless service
	// +kubebuilder:default="flux-service"
	// +default="flux-service"
	// +optional
	HeadlessName string `json:"headlessName,omitempty"`
}

type MiniClusterUser struct {

	// If a user is defined, the username is required
	Name string `json:"name"`

	// +optional
	Password string `json:"password"`
}

type MiniClusterArchive struct {

	// Save or load from this directory path
	// +optional
	Path string `json:"path,omitempty"`
}

type LoggingSpec struct {

	// Quiet mode silences all output so the job only shows the test running
	// +kubebuilder:default=false
	// +default=false
	// +optional
	Quiet bool `json:"quiet,omitempty"`

	// Strict mode ensures any failure will not continue in the job entrypoint
	// +kubebuilder:default=true
	// +default=true
	// +optional
	Strict bool `json:"strict,omitempty"`

	// Debug mode adds extra verbosity to Flux
	// +kubebuilder:default=false
	// +default=false
	// +optional
	Debug bool `json:"debug,omitempty"`

	// Enable Zeromq logging
	// +kubebuilder:default=false
	// +default=false
	// +optional
	Zeromq bool `json:"zeromq,omitempty"`

	// Timed mode adds timing to Flux commands
	// +kubebuilder:default=false
	// +default=false
	// +optional
	Timed bool `json:"timed,omitempty"`
}

// PodSpec controlls variables for the cluster pod
type PodSpec struct {

	// Annotations for each pod
	// +optional
	Annotations map[string]string `json:"annotations"`

	// Labels for each pod
	// +optional
	Labels map[string]string `json:"labels"`

	// Service account name for the pod
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// NodeSelectors for a pod
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Resources include limits and requests
	// +optional
	Resources ContainerResource `json:"resources"`
}

// MiniClusterStatus defines the observed state of Flux
type MiniClusterStatus struct {

	// These are for the sub-resource scale functionality
	Size     int32  `json:"size"`
	Selector string `json:"selector"`

	// The Jobid is set internally to associate to a miniCluster
	// This isn't currently in use, we only have one!
	Jobid string `json:"jobid"`

	// We keep the original size of the MiniCluster request as
	// this is the absolute maximum
	MaximumSize int32 `json:"maximumSize"`

	// conditions hold the latest Flux Job and MiniCluster states
	// +listType=atomic
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type FluxRestful struct {

	// Branch to clone Flux Restful API from
	// +kubebuilder:default="main"
	// +default="main"
	// +optional
	Branch string `json:"branch,omitempty"`

	// Port to run Flux Restful Server On
	// +kubebuilder:default=5000
	// +default=5000
	// +optional
	Port int32 `json:"port,omitempty"`

	// Secret key shared between server and client
	// +optional
	SecretKey string `json:"secretKey"`

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

	// Annotations for the volume
	// +optional
	Annotations map[string]string `json:"annotations"`

	// Annotations for the persistent volume claim
	// +optional
	ClaimAnnotations map[string]string `json:"claimAnnotations"`

	// Optional volume attributes
	// +optional
	Attributes map[string]string `json:"attributes"`

	// Volume handle, falls back to storage class name
	// if not defined
	// +optional
	VolumeHandle string `json:"volumeHandle"`

	// +kubebuilder:default="hostpath"
	// +default="hostpath"
	// +optional
	StorageClass string `json:"storageClass,omitempty"`

	// Storage driver, e.g., gcs.csi.ofek.dev
	// Only needed if not using hostpath
	// +optional
	Driver string `json:"driver"`

	// Delete the persistent volume on cleanup
	// +kubebuilder:default=true
	// +default=true
	// +optional
	Delete bool `json:"delete,omitempty"`

	// Secret reference in Kubernetes with service account role
	// +optional
	Secret string `json:"secret"`

	// Secret namespace
	// +kubebuilder:default="default"
	// +default="default"
	// +optional
	SecretNamespace string `json:"secretNamespace,omitempty"`

	// Capacity (string) for PVC (storage request) to create PV
	// +kubebuilder:default="5Gi"
	// +default="5Gi"
	// +optional
	Capacity string `json:"capacity,omitempty"`
}

// Mini Cluster local volumes available to mount (these are on the host)
type MiniClusterExistingVolume struct {

	// Path and claim name are always required if a secret isn't defined
	// +optional
	Path string `json:"path,omitempty"`

	// Config map name if the existing volume is a config map
	// You should also define items if you are using this
	// +optional
	ConfigMapName string `json:"configMapName,omitempty"`

	// Items (key and paths) for the config map
	// +optional
	Items map[string]string `json:"items"`

	// Claim name if the existing volume is a PVC
	// +optional
	ClaimName string `json:"claimName,omitempty"`

	// An existing secret
	// +optional
	SecretName string `json:"secretName,omitempty"`

	// +kubebuilder:default=false
	// +default=false
	// +optional
	ReadOnly bool `json:"readOnly,omitempty"`
}

// A Container volume must reference one defined for the MiniCluster
// The path here is in the container
type ContainerVolume struct {
	Path string `json:"path"`

	// +kubebuilder:default=false
	// +default=false
	// +optional
	ReadOnly bool `json:"readOnly,omitempty"`
}

type FluxSpec struct {

	// Install root location
	// +kubebuilder:default="/usr"
	// +default="/usr"
	InstallRoot string `json:"installRoot,omitempty"`

	// Commands for flux start --wrap
	// +optional
	Wrap string `json:"wrap,omitempty"`

	// Single user executable to provide to flux start
	// +kubebuilder:default="5s"
	// +default="5s"
	ConnectTimeout string `json:"connectTimeout,omitempty"`

	// Flux option flags, usually provided with -o
	// optional - if needed, default option flags for the server
	// These can also be set in the user interface to override here.
	// This is only valid for a FluxRunner "runFlux" true
	// +optional
	OptionFlags string `json:"optionFlags"`

	// Only expose the broker service (to reduce load on DNS)
	// +optional
	MinimalService bool `json:"minimalService"`

	// Log level to use for flux logging (only in non TestMode)
	// +kubebuilder:default=6
	// +default=6
	// +optional
	LogLevel int32 `json:"logLevel,omitempty"`

	// Optionally provide an already existing curve certificate
	// This is not recommended in favor of providing the secret
	// name as curveCertSecret, below
	//+optional
	CurveCert string `json:"curveCert"`

	// Expect a secret for a curve cert here.
	// This is ideal over the curveCert (as a string) above.
	//+optional
	CurveCertSecret string `json:"curveCertSecret"`

	// Expect a secret (named according to this string)
	// for a munge key. This is intended for bursting.
	// Assumed to be at /etc/munge/munge.key
	// This is binary data.
	//+optional
	MungeSecret string `json:"mungeSecret"`

	// Bursting - one or more external clusters to burst to
	// We assume a single, central MiniCluster with an ipaddress
	// that all connect to.
	//+optional
	Bursting Bursting `json:"bursting"`

	// Optionally provide a manually created broker config
	// this is intended for bursting to remote clusters
	//+optional
	BrokerConfig string `json:"brokerConfig"`
}

// Bursting Config
// For simplicity, we internally handle the name of the job (hostnames)
type Bursting struct {

	// The lead broker ip address to join to. E.g., if we burst
	// to cluster 2, this is the address to connect to cluster 1
	// For the first cluster, this should not be defined
	//+optional
	LeadBroker FluxBroker `json:"leadBroker"`

	// External clusters to burst to. Each external
	// cluster must share the same listing to align ranks
	Clusters []BurstedCluster `json:"clusters"`
}

type BurstedCluster struct {

	// The hostnames for the bursted clusters
	// If set, the user is responsible for ensuring
	// uniqueness. The operator will set to burst-N
	//+optional
	Name string `json:"name"`

	// Size of bursted cluster.
	// Defaults to same size as local minicluster if not set
	// +optional
	Size int32 `json:"size,omitempty"`
}

// A FluxBroker defines a broker for flux
type FluxBroker struct {

	// Lead broker address (ip or hostname)
	Address string `json:"address"`

	// We need the name of the lead job to assemble the hostnames
	Name string `json:"name"`

	// Lead broker size
	Size int32 `json:"size"`

	// Lead broker port - should only be used for external cluster
	// +kubebuilder:default=8050
	// +default=8050
	// +optional
	Port int32 `json:"port,omitempty"`
}

type MiniClusterContainer struct {

	// Container image must contain flux and flux-sched install
	// +kubebuilder:default="ghcr.io/rse-ops/accounting:app-latest"
	// +default="ghcr.io/rse-ops/accounting:app-latest"
	Image string `json:"image,omitempty"`

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

	// Secrets that will be added to the environment
	// The user is expected to create their own secrets for the operator to find
	// +optional
	Secrets map[string]Secret `json:"secrets"`

	// Allow the user to pull authenticated images
	// By default no secret is selected. Setting
	// this with the name of an already existing
	// imagePullSecret will specify that secret
	// in the pod spec.
	// +optional
	ImagePullSecret string `json:"imagePullSecret"`

	// Single user executable to provide to flux start
	// +optional
	Command string `json:"command"`

	// Indicate that the command is a launcher that will
	// ask for its own jobs (and provided directly to flux start)
	// +optional
	Launcher bool `json:"launcher"`

	// Indicate that the command is a batch job that will be written to a file to submit
	// +optional
	Batch bool `json:"batch"`

	// Don't wrap batch commands in flux submit (provide custom logic myself)
	// +optional
	BatchRaw bool `json:"batchRaw"`

	// Log output directory
	// +optional
	Logs string `json:"logs"`

	// Allow the user to dictate pulling
	// By default we pull if not present. Setting
	// this to true will indicate to pull always
	// +kubebuilder:default=false
	// +default=false
	// +optional
	PullAlways bool `json:"pullAlways,omitempty"`

	// Main container to run flux (only should be one)
	// +optional
	RunFlux bool `json:"runFlux"`

	// Volumes that can be mounted (must be defined in volumes)
	// +optional
	Volumes map[string]ContainerVolume `json:"volumes"`

	// Existing Volumes to add to the containers
	// +optional
	ExistingVolumes map[string]MiniClusterExistingVolume `json:"existingVolumes"`

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

	// Security Context
	// https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`
}

type SecurityContext struct {

	// Privileged container
	// +optional
	Privileged bool `json:"privileged,omitempty"`
}

type LifeCycle struct {

	// +optional
	PostStartExec string `json:"postStartExec"`

	// +optional
	PreStopExec string `json:"preStopExec"`
}

type FluxUser struct {

	// Flux user name
	// +kubebuilder:default="flux"
	// +default="flux"
	// +optional
	Name string `json:"name,omitempty"`

	// UID for the FluxUser
	// +kubebuilder:default=1000
	// +default=1000
	// +optional
	Uid int `json:"uid,omitempty"`
}

// Secret describes a secret from the environment.
// The envar name should be the key of the top level map.
type Secret struct {

	// Name under secretKeyRef->Name
	Name string `json:"name"`

	// Key under secretKeyRef->Key
	Key string `json:"key"`
}

type Commands struct {

	// Run flux start as root - required for some storage binds
	// +kubebuilder:default=false
	// +default=false
	// +optional
	RunFluxAsRoot bool `json:"runFluxAsRoot,omitempty"`

	// Prefix to flux start / submit / broker
	// Typically used for a wrapper command to mount, etc.
	// +optional
	Prefix string `json:"prefix"`

	// init command is run before anything
	// +optional
	Init string `json:"init"`

	// pre command is run after global PreCommand, after asFlux is set (can override)
	// +optional
	Pre string `json:"pre"`

	// post command is run in the entrypoint when the broker exits / finishes
	// +optional
	Post string `json:"post"`

	// A command only for workers to run
	// +optional
	WorkerPre string `json:"workerPre"`

	// A single command for only the broker to run
	// +optional
	BrokerPre string `json:"brokerPre"`
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
// +kubebuilder:subresource:scale:specpath=.spec.size,statuspath=.status.size,selectorpath=.status.selector
type MiniCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MiniClusterSpec   `json:"spec,omitempty"`
	Status MiniClusterStatus `json:"status,omitempty"`
}

// MultuUser returns boolean to indicate if we are in multi-user mode
func (f *MiniCluster) MultiUser() bool {
	return len(f.Spec.Users) > 0
}

// Return a lookup of all container existing volumes (for the higher level Pod)
// Volumes are unique by name.
func (f *MiniCluster) ExistingContainerVolumes() map[string]MiniClusterExistingVolume {
	return uniqueExistingVolumes(f.Spec.Containers)
}

// Return a lookup of all service existing volumes (for the higher level Pod)
// Volumes are unique by name.
func (f *MiniCluster) ExistingServiceVolumes() map[string]MiniClusterExistingVolume {
	return uniqueExistingVolumes(f.Spec.Services)
}

// uniqueExistingVolumes is a shared function to populate existing container or service volumes
func uniqueExistingVolumes(containers []MiniClusterContainer) map[string]MiniClusterExistingVolume {
	volumes := map[string]MiniClusterExistingVolume{}
	for _, container := range containers {
		for name, volume := range container.ExistingVolumes {
			volumes[name] = volume
		}
	}
	return volumes
}

// fluxInstallRoot returns the flux install root
func (f *MiniCluster) FluxInstallRoot() string {
	root := f.Spec.Flux.InstallRoot
	if root == "" {
		root = "/usr"
	}
	return root
}

// Validate ensures we have data that is needed, and sets defaults if needed
func (f *MiniCluster) Validate() bool {
	fmt.Println()

	// Global (entire cluster) settings
	fmt.Printf("🤓 MiniCluster.DeadlineSeconds %d\n", f.Spec.DeadlineSeconds)
	fmt.Printf("🤓 MiniCluster.Size %s\n", fmt.Sprint(f.Spec.Size))

	// If MaxSize is set, it must be greater than size
	if f.Spec.MaxSize != 0 && f.Spec.MaxSize < f.Spec.Size {
		fmt.Printf("😥️ MaxSize of cluster must be greater than size.\n")
		return false
	}

	// Set the Flux install root
	if f.Spec.Flux.InstallRoot == "" {
		f.Spec.Flux.InstallRoot = "/usr"
	}

	// Set the default headless service name
	if f.Spec.Network.HeadlessName == "" {
		f.Spec.Network.HeadlessName = "flux-service"
	}

	// Flux Restful default port
	if f.Spec.FluxRestful.Port == 0 {
		f.Spec.FluxRestful.Port = 5000
	}

	// If the MaxSize isn't set, ensure it's equal to the size
	if f.Spec.MaxSize == 0 {
		f.Spec.MaxSize = f.Spec.Size
	}

	// If we haven't seen a MaxSize (in the status) yet, set it
	// This needs to be the absolute max that is allowed
	if f.Status.MaximumSize == 0 {
		f.Status.MaximumSize = f.Spec.Size
		if f.Spec.MaxSize > f.Spec.Size {
			f.Status.MaximumSize = f.Spec.MaxSize
		}
	}
	fmt.Printf("🤓 MiniCluster.MaximumSize %s\n", fmt.Sprint(f.Status.MaximumSize))

	// We should have only one flux runner
	valid := true
	fluxRunners := 0

	// Commands and PreCommand not supported for services
	for _, service := range f.Spec.Services {
		if service.Name == "" {
			fmt.Printf("😥️ Service containers always require a name.\n")
			return false
		}
		if service.PreCommand != "" || service.Commands.Pre != "" ||
			service.Commands.BrokerPre != "" || service.Commands.WorkerPre != "" {
			fmt.Printf("😥️ Services do not support Commands.\n")
			return false
		}
	}

	// If we have a LeadBroker address, this is a child cluster, and
	// we also need a port
	if f.Spec.Flux.Bursting.LeadBroker.Port == 0 {
		f.Spec.Flux.Bursting.LeadBroker.Port = 8050
	}

	// Set default port if unset
	for b, bursted := range f.Spec.Flux.Bursting.Clusters {

		// If bursted size not set, default to the current MiniCluster size
		if bursted.Size == 0 {
			bursted.Size = f.Spec.Size
		}

		// Set default name if not set to burst-N
		if bursted.Name == "" {
			bursted.Name = fmt.Sprintf("burst-%d", b)
		}
	}

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

		// Launcher mode does not work with batch
		if container.Launcher && container.Batch {
			fmt.Printf("😥️ %s is indicated for batch and launcher, choose one.\n", name)
			return false
		}

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
		if volume.StorageClass != "hostpath" && volume.Secret == "" {
			fmt.Printf("😥️ Found non-hostpath volume %s that is missing a secret\n", key)
			valid = false
		}
	}

	// For existing volumes, if it's a claim, a path is required.
	if !f.validateExistingVolumes(f.ExistingContainerVolumes()) {
		fmt.Printf("😥️ Existing container volumes are not valid\n")
		return false
	}
	if !f.validateExistingVolumes(f.ExistingServiceVolumes()) {
		fmt.Printf("😥️ Existing service volumes are not valid\n")
		return false
	}

	return valid
}

// validateExistingVolumes ensures secret names vs. volume paths are valid
func (f *MiniCluster) validateExistingVolumes(existing map[string]MiniClusterExistingVolume) bool {
	valid := true
	for key, volume := range existing {

		// Case 1: it's a secret and we only need that
		if volume.SecretName != "" {
			continue
		}

		// Case 2: it's a config map (and will have items too, but we don't hard require them)
		if volume.ConfigMapName != "" {
			continue
		}

		// Case 3: claim desired without path
		if volume.ClaimName == "" && volume.Path != "" {
			fmt.Printf("😥️ Found existing volume %s with path %s that is missing a claim name\n", key, volume.Path)
			valid = false
		}
		// Case 4: reverse of the above
		if volume.ClaimName != "" && volume.Path == "" {
			fmt.Printf("😥️ Found existing volume %s with claimName %s that is missing a path\n", key, volume.ClaimName)
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
