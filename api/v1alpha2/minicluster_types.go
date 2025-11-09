/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package v1alpha2

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	entrypointSuffix = "-entrypoint"
)

// MiniCluster is an HPC cluster in Kubernetes you can control
// Either to submit a single job (and go away) or for a persistent single- or multi- user cluster
type MiniClusterSpec struct {
	// Important: Run "make" and "make manifests" to regenerate code after modifying this file

	// Containers is one or more containers to be created in a pod.
	// There should only be one container to run flux with runFlux
	// +listType=atomic
	Containers []MiniClusterContainer `json:"containers"`

	// ResourceClaims to be referenced in containers
	// +optional
	ResourceClaims []corev1.PodResourceClaim `json:"resourceClaims"`

	// Services are one or more service containers to bring up
	// alongside the MiniCluster.
	// +optional
	// +listType=atomic
	Services []MiniClusterContainer `json:"services"`

	// A spec for exposing or defining the cluster headless service
	//+optional
	Network Network `json:"network"`

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

	// Logging modes determine the output you see in the job log
	// +optional
	Logging LoggingSpec `json:"logging"`

	// Archive to load or save
	// +optional
	Archive MiniClusterArchive `json:"archive"`

	// Share process namespace?
	// +optional
	ShareProcessNamespace bool `json:"shareProcessNamespace"`

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

	// BackoffLimit is the number of retries for the job before failing
	// +optional
	BackoffLimit int32 `json:"backoffLimit,omitempty"`

	// MinSize (minimum number of pods that must be up for Flux)
	// Note that this option does not edit the number of tasks,
	// so a job could run with fewer (and then not start)
	// +optional
	MinSize int32 `json:"minSize,omitempty"`

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
	// +optional
	HeadlessName string `json:"headlessName,omitempty"`

	// Disable affinity rules that guarantee one network address / node
	// +optional
	DisableAffinity bool `json:"disableAffinity,omitempty"`
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
	// +kubebuilder:default=false
	// +default=false
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

	// Restart Policy
	// +optional
	RestartPolicy string `json:"restartPolicy,omitempty"`

	// Service account name for the pod
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// RuntimeClassName for the pod
	// +optional
	RuntimeClassName string `json:"runtimeClassName,omitempty"`

	// Automatically mount the service account name
	// +optional
	AutomountServiceAccountToken bool `json:"automountServiceAccountToken,omitempty"`

	// Scheduler name for the pod
	// +optional
	SchedulerName string `json:"schedulerName,omitempty"`

	// NodeSelectors for a pod
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// NodeAffinity is for a list of values assoicated with a label
	// +optional
	NodeAffinity map[string][]string `json:"nodeAffinity,omitempty"`

	// Tolerations for a pod
	// +optional
	Tolerations []Toleration `json:"tolerations,omitempty"`

	// Pod DNS policy (defaults to ClusterFirst)
	// +optional
	DNSPolicy string `json:"dnsPolicy,omitempty"`

	// Use Host IPC
	// +optional
	HostIPC bool `json:"hostIPC,omitempty"`

	// Use Host PID
	// +optional
	HostPID bool `json:"hostPID,omitempty"`

	// Resources include limits and requests
	// +optional
	Resources corev1.ResourceList `json:"resources"`

	// PodSecurity Context
	// +optional
	SecurityContext corev1.PodSecurityContext `json:"securityContext,omitempty"`
}

type Toleration struct {

	// The label key to tolerate
	// +optional
	Key string `json:"key,omitempty"`

	// The effect to have
	// +optional
	Effect string `json:"effect,omitempty"`

	// The operator to use (e.g., Equal)
	// +optional
	Operator string `json:"operator,omitempty"`

	// E.g., NoSchedule
	// +optional
	Value string `json:"value,omitempty"`
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

// Mini Cluster local volumes available to mount (these are on the host)

type ContainerVolume struct {

	// Path and claim name are always required if a secret isn't defined
	// +optional
	Path string `json:"path,omitempty"`

	// An existing hostPath to bind to path
	// +optional
	HostPath string `json:"hostPath,omitempty"`

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

	// Add an empty directory custom type
	// +optional
	EmptyDirMedium string `json:"emptyDirMedium,omitempty"`

	// Add an empty directory sizeLimit
	// +optional
	EmptyDirSizeLimit string `json:"emptyDirSizeLimit,omitempty"`

	// Add a csi driver type volume
	// +optional
	CSIDriver string `json:"csiDriver,omitempty"`

	// Add attributes for the csi driver
	// +optional
	CSIDriverAttributes map[string]string `json:"csiDriverAttributes"`

	// +kubebuilder:default=false
	// +default=false
	// +optional
	EmptyDir bool `json:"emptyDir,omitempty"`
}

type FluxSpec struct {

	// Container base for flux. Options include only:
	// ghcr.io/converged-computing/flux-view-rocky:arm-9
	// ghcr.io/converged-computing/flux-view-rocky:arn-8
	// ghcr.io/converged-computing/flux-view-rocky:tag-9
	// ghcr.io/converged-computing/flux-view-rocky:tag-8
	// ghcr.io/converged-computing/flux-view-ubuntu:tag-noble
	// ghcr.io/converged-computing/flux-view-ubuntu:tag-jammy
	// ghcr.io/converged-computing/flux-view-ubuntu:tag-focal
	// ghcr.io/converged-computing/flux-view-ubuntu:arm-jammy
	// ghcr.io/converged-computing/flux-view-ubuntu:arm-focal
	// +optional
	Container FluxContainer `json:"container,omitempty"`

	// Change the arch string - determines the binaries
	// that are downloaded to run the entrypoint
	// +optional
	Arch string `json:"arch,omitempty"`

	// Modify flux submit to be something else
	// +optional
	SubmitCommand string `json:"submitCommand,omitempty"`

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

	// Disable specifying the socket path
	// +optional
	DisableSocket bool `json:"disableSocket"`

	// Specify a custom Topology
	// +optional
	Topology string `json:"topology"`

	// Do not wait for the socket
	// +optional
	NoWaitSocket bool `json:"noWaitSocket"`

	// Complete workers when they fail
	// This is ideal if you don't want them to restart
	// +optional
	CompleteWorkers bool `json:"completeWorkers"`

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

	// Custom attributes for the fluxion scheduler
	//+optional
	Scheduler FluxScheduler `json:"scheduler"`

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

	// Environment
	// If defined, set these envars for the flux view.
	//+optional
	Environment map[string]string `json:"environment"`
}

// FluxScheduler attributes
type FluxScheduler struct {

	// Use sched-simple (no support for GPU)
	// +optional
	Simple bool `json:"simple"`

	// Scheduler queue policy, defaults to "fcfs" can also be "easy"
	// +optional
	QueuePolicy string `json:"queuePolicy"`
}

// Bursting Config
// For simplicity, we internally handle the name of the job (hostnames)
type Bursting struct {

	// The lead broker ip address to join to. E.g., if we burst
	// to cluster 2, this is the address to connect to cluster 1
	// For the first cluster, this should not be defined
	//+optional
	LeadBroker FluxBroker `json:"leadBroker"`

	// Hostlist is a custom hostlist for the broker.toml
	// that includes the local plus bursted cluster. This
	// is typically used for bursting to another resource
	// type, where we can predict the hostnames but they
	// don't follow the same convention as the Flux Operator
	//+optional
	Hostlist string `json:"hostlist"`

	// External clusters to burst to. Each external
	// cluster must share the same listing to align ranks
	//+optional
	// +listType=atomic
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

// A FluxContainer is equivalent to a MiniCluster container but has a different default image
type FluxContainer struct {

	// Disable the sidecar container, assuming that the main application container has flux
	// +kubebuilder:default=false
	// +default=false
	Disable bool `json:"disable,omitempty"`

	// Container name is only required for non flux runners
	// +kubebuilder:default="flux-view"
	// +default="flux-view"
	Name string `json:"name,omitempty"`

	// Working directory to run command from
	// +optional
	WorkingDir string `json:"workingDir"`

	// Customize python path for flux
	// +optional
	PythonPath string `json:"pythonPath"`

	// Resources include limits and requests
	// These must be defined for cpu and memory
	// for the QoS to be Guaranteed
	// +optional
	Resources corev1.ResourceRequirements `json:"resources"`

	// Allow the user to pull authenticated images
	// By default no secret is selected. Setting
	// this with the name of an already existing
	// imagePullSecret will specify that secret
	// in the pod spec.
	// +optional
	ImagePullSecret string `json:"imagePullSecret"`

	// +kubebuilder:default="ghcr.io/converged-computing/flux-view-rocky:tag-9"
	// +default="ghcr.io/converged-computing/flux-view-rocky:tag-9"
	Image string `json:"image,omitempty"`

	// Allow the user to dictate pulling
	// By default we pull if not present. Setting
	// this to true will indicate to pull always
	// +kubebuilder:default=false
	// +default=false
	// +optional
	PullAlways bool `json:"pullAlways,omitempty"`

	// Security Context
	// +optional
	SecurityContext corev1.SecurityContext `json:"securityContext,omitempty"`

	// Mount path for flux to be at (will be added to path)
	// +kubebuilder:default="/mnt/flux"
	// +default="/mnt/flux"
	MountPath string `json:"mountPath,omitempty"`
}

type MiniClusterContainer struct {

	// Allow the user to pull authenticated images
	// By default no secret is selected. Setting
	// this with the name of an already existing
	// imagePullSecret will specify that secret
	// in the pod spec.
	// +optional
	ImagePullSecret string `json:"imagePullSecret"`

	// Container image must contain flux and flux-sched install
	// +kubebuilder:default="ghcr.io/rse-ops/accounting:app-latest"
	// +default="ghcr.io/rse-ops/accounting:app-latest"
	Image string `json:"image,omitempty"`

	// Working directory to run command from
	// +optional
	WorkingDir string `json:"workingDir"`

	// Container name is only required for non flux runners
	// +optional
	Name string `json:"name"`

	// Single user executable to provide to flux start
	// +optional
	Command string `json:"command"`

	// Allow the user to dictate pulling
	// By default we pull if not present. Setting
	// this to true will indicate to pull always
	// +kubebuilder:default=false
	// +default=false
	// +optional
	PullAlways bool `json:"pullAlways,omitempty"`

	// Allow the user to dictate pulling directly
	// +optional
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`

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

	// Application container intended to run flux (broker)
	// +optional
	RunFlux bool `json:"runFlux"`

	// Do not wrap the entrypoint to wait for flux, add to path, etc?
	// +optional
	NoWrapEntrypoint bool `json:"noWrapEntrypoint"`

	// Existing volumes that can be mounted
	// +optional
	Volumes map[string]ContainerVolume `json:"volumes"`

	// Lifecycle can handle post start commands, etc.
	// +optional
	LifeCycle LifeCycle `json:"lifeCycle"`

	// Resources include limits and requests
	// +optional
	Resources corev1.ResourceRequirements `json:"resources"`

	// More specific or detailed commands for just workers/broker
	// +optional
	Commands Commands `json:"commands"`

	// Security Context
	// https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
	// +optional
	SecurityContext corev1.SecurityContext `json:"securityContext,omitempty"`
}

type LifeCycle struct {

	// +optional
	PostStartExec string `json:"postStartExec"`

	// +optional
	PreStopExec string `json:"preStopExec"`
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

	// Prefix to flux start / submit / broker
	// Typically used for a wrapper command to mount, etc.
	// +optional
	Prefix string `json:"prefix"`

	// Custom script for submit (e.g., multiple lines)
	// +optional
	Script string `json:"script"`

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

	// A command only for service start.sh tor run
	// +optional
	ServicePre string `json:"servicePre"`
}

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

// Determine if a MiniCluster container has custom commands
// if we have custom commands and a command entrypoint we can support additional custom logic
func (c *MiniClusterContainer) HasCommands() bool {
	return c.Commands.Pre != "" || c.Commands.BrokerPre != "" || c.Commands.WorkerPre != "" || c.Commands.Init != "" || c.Commands.Post != ""
}

// Determine if we should generate a start.sh entrypoint for a sidecar
// Only do so (for now) if we are customizing the command
func (c *MiniClusterContainer) GenerateEntrypoint() bool {
	return !c.RunFlux && !c.NoWrapEntrypoint
}

// Return a lookup of all container existing volumes (for the higher level Pod)
// Volumes are unique by name.
func (f *MiniCluster) ExistingContainerVolumes() map[string]ContainerVolume {
	return uniqueExistingVolumes(f.Spec.Containers)
}

// Return a lookup of all service existing volumes (for the higher level Pod)
// Volumes are unique by name.
func (f *MiniCluster) ExistingServiceVolumes() map[string]ContainerVolume {
	return uniqueExistingVolumes(f.Spec.Services)
}

// uniqueExistingVolumes is a shared function to populate existing container or service volumes
func uniqueExistingVolumes(containers []MiniClusterContainer) map[string]ContainerVolume {
	volumes := map[string]ContainerVolume{}
	for _, container := range containers {
		for name, volume := range container.Volumes {
			volumes[name] = volume
		}
	}
	return volumes
}

// Consistent functions to return config map names
func (f *MiniCluster) EntrypointConfigMapName() string {
	return f.Name + entrypointSuffix
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

	// BackoffLimit must be postive if set
	if f.Spec.BackoffLimit != 0 && f.Spec.BackoffLimit < 0 {
		fmt.Printf("😥️ BackoffLimit of cluster must be greater than 0.\n")
		return false
	}

	// If MinSize is set, it must be <= MaxSize and Size
	if f.Spec.MinSize != 0 && f.Spec.MaxSize != 0 && f.Spec.MinSize > f.Spec.MaxSize {
		fmt.Printf("😥️ MinSize of cluster must be less than MaxSize.\n")
		return false
	}
	if f.Spec.MinSize != 0 && f.Spec.MinSize > f.Spec.Size {
		fmt.Printf("😥️ MinSize of cluster must be less than size.\n")
		return false
	}

	// Set the default headless service name
	if f.Spec.Network.HeadlessName == "" {
		f.Spec.Network.HeadlessName = f.Name
	}
	if f.Spec.Flux.Container.Name == "" {
		f.Spec.Flux.Container.Name = "flux-view"
	}
	if f.Spec.Flux.Container.MountPath == "" {
		f.Spec.Flux.Container.MountPath = "/mnt/flux"
	}
	if f.Spec.Flux.Container.Image == "" {
		f.Spec.Flux.Container.Image = "ghcr.io/converged-computing/flux-view-rocky:tag-9"
	}
	if f.Spec.Flux.Scheduler.QueuePolicy == "" {
		f.Spec.Flux.Scheduler.QueuePolicy = "fcfs"
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
		if service.Commands.Pre != "" ||
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

	// If we are provided a hostlist, we don't need bursted clusters
	if f.Spec.Flux.Bursting.Hostlist != "" && len(f.Spec.Flux.Bursting.Clusters) > 0 {
		fmt.Printf("😥️ A custom hostlist cannot be provided with a bursting spec, choose one or the other!\n")
		return false
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

		// If a custom script is provided AND a command, no go
		if (container.Commands.Script != "" && container.Command != "") && container.RunFlux {
			fmt.Printf("😥️ %s has both a script and command provided, choose one\n", name)
			return false
		}
	}
	if fluxRunners != 1 {
		valid = false
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
func (f *MiniCluster) validateExistingVolumes(existing map[string]ContainerVolume) bool {
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

		// Case 3: claim name without path
		if volume.ClaimName != "" && volume.Path == "" {
			fmt.Printf("😥️ Found existing volume %s with claimName %s that is missing a path\n", key, volume.ClaimName)
			valid = false
		}

		// Case 4: hostpath without path
		if volume.HostPath != "" && volume.Path == "" {
			fmt.Printf("😥️ Found existing volume %s with hostPath %s that is missing a path\n", key, volume.HostPath)
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
