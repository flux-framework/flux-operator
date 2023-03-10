//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	common "k8s.io/kube-openapi/pkg/common"
	spec "k8s.io/kube-openapi/pkg/validation/spec"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"./api/v1alpha1/.Commands":             schema__api_v1alpha1__Commands(ref),
		"./api/v1alpha1/.ContainerResources":   schema__api_v1alpha1__ContainerResources(ref),
		"./api/v1alpha1/.ContainerVolume":      schema__api_v1alpha1__ContainerVolume(ref),
		"./api/v1alpha1/.FluxRestful":          schema__api_v1alpha1__FluxRestful(ref),
		"./api/v1alpha1/.FluxUser":             schema__api_v1alpha1__FluxUser(ref),
		"./api/v1alpha1/.LifeCycle":            schema__api_v1alpha1__LifeCycle(ref),
		"./api/v1alpha1/.LoggingSpec":          schema__api_v1alpha1__LoggingSpec(ref),
		"./api/v1alpha1/.MiniCluster":          schema__api_v1alpha1__MiniCluster(ref),
		"./api/v1alpha1/.MiniClusterContainer": schema__api_v1alpha1__MiniClusterContainer(ref),
		"./api/v1alpha1/.MiniClusterList":      schema__api_v1alpha1__MiniClusterList(ref),
		"./api/v1alpha1/.MiniClusterSpec":      schema__api_v1alpha1__MiniClusterSpec(ref),
		"./api/v1alpha1/.MiniClusterStatus":    schema__api_v1alpha1__MiniClusterStatus(ref),
		"./api/v1alpha1/.MiniClusterUser":      schema__api_v1alpha1__MiniClusterUser(ref),
		"./api/v1alpha1/.MiniClusterVolume":    schema__api_v1alpha1__MiniClusterVolume(ref),
		"./api/v1alpha1/.PodSpec":              schema__api_v1alpha1__PodSpec(ref),
	}
}

func schema__api_v1alpha1__Commands(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"runFluxAsRoot": {
						SchemaProps: spec.SchemaProps{
							Description: "Run flux start as root - required for some storage binds",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"pre": {
						SchemaProps: spec.SchemaProps{
							Description: "pre command is run after global PreCommand, before anything else",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schema__api_v1alpha1__ContainerResources(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ContainerResources include limits and requests",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"limits": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("k8s.io/apimachinery/pkg/util/intstr.IntOrString"),
									},
								},
							},
						},
					},
					"requests": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("k8s.io/apimachinery/pkg/util/intstr.IntOrString"),
									},
								},
							},
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/util/intstr.IntOrString"},
	}
}

func schema__api_v1alpha1__ContainerVolume(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "A Container volume must reference one defined for the MiniCluster The path here is in the container",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"path": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
					"readOnly": {
						SchemaProps: spec.SchemaProps{
							Default: false,
							Type:    []string{"boolean"},
							Format:  "",
						},
					},
				},
				Required: []string{"path"},
			},
		},
	}
}

func schema__api_v1alpha1__FluxRestful(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"branch": {
						SchemaProps: spec.SchemaProps{
							Description: "Branch to clone Flux Restful API from",
							Default:     "main",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"port": {
						SchemaProps: spec.SchemaProps{
							Description: "Port to run Flux Restful Server On",
							Default:     5000,
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"secretKey": {
						SchemaProps: spec.SchemaProps{
							Description: "Secret key shared between server and client",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"username": {
						SchemaProps: spec.SchemaProps{
							Description: "These two should not actually be set by a user, but rather generated by tools and provided Username to use for RestFul API",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"token": {
						SchemaProps: spec.SchemaProps{
							Description: "Token to use for RestFul API",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schema__api_v1alpha1__FluxUser(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"name": {
						SchemaProps: spec.SchemaProps{
							Description: "Flux user name",
							Default:     "flux",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"uid": {
						SchemaProps: spec.SchemaProps{
							Description: "UID for the FluxUser",
							Default:     1000,
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
				},
			},
		},
	}
}

func schema__api_v1alpha1__LifeCycle(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"postStartExec": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
				},
			},
		},
	}
}

func schema__api_v1alpha1__LoggingSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"quiet": {
						SchemaProps: spec.SchemaProps{
							Description: "Quiet mode silences all output so the job only shows the test running",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"strict": {
						SchemaProps: spec.SchemaProps{
							Description: "Strict mode ensures any failure will not continue in the job entrypoint",
							Default:     true,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"debug": {
						SchemaProps: spec.SchemaProps{
							Description: "Debug mode adds extra verbosity to Flux",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"timed": {
						SchemaProps: spec.SchemaProps{
							Description: "Timed mode adds timing to Flux commands",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schema__api_v1alpha1__MiniCluster(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("./api/v1alpha1/.MiniClusterSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("./api/v1alpha1/.MiniClusterStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./api/v1alpha1/.MiniClusterSpec", "./api/v1alpha1/.MiniClusterStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema__api_v1alpha1__MiniClusterContainer(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"image": {
						SchemaProps: spec.SchemaProps{
							Description: "Container image must contain flux and flux-sched install",
							Default:     "ghcr.io/rse-ops/accounting:app-latest",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"name": {
						SchemaProps: spec.SchemaProps{
							Description: "Container name is only required for non flux runners",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"cores": {
						SchemaProps: spec.SchemaProps{
							Description: "Cores the container should use",
							Default:     0,
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"workingDir": {
						SchemaProps: spec.SchemaProps{
							Description: "Working directory to run command from",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"diagnostics": {
						SchemaProps: spec.SchemaProps{
							Description: "Run flux diagnostics on start instead of command",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"fluxUser": {
						SchemaProps: spec.SchemaProps{
							Description: "Flux User, if created in the container",
							Default:     map[string]interface{}{},
							Ref:         ref("./api/v1alpha1/.FluxUser"),
						},
					},
					"ports": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "atomic",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Ports to be exposed to other containers in the cluster We take a single list of integers and map to the same",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: 0,
										Type:    []string{"integer"},
										Format:  "int32",
									},
								},
							},
						},
					},
					"environment": {
						SchemaProps: spec.SchemaProps{
							Description: "Key/value pairs for the environment",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
					"imagePullSecret": {
						SchemaProps: spec.SchemaProps{
							Description: "Allow the user to pull authenticated images By default no secret is selected. Setting this with the name of an already existing imagePullSecret will specify that secret in the pod spec.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"command": {
						SchemaProps: spec.SchemaProps{
							Description: "Single user executable to provide to flux start IMPORTANT: This is left here, but not used in favor of exposing Flux via a Restful API. We Can remove this when that is finalized.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"pullAlways": {
						SchemaProps: spec.SchemaProps{
							Description: "Allow the user to dictate pulling By default we pull if not present. Setting this to true will indicate to pull always",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"runFlux": {
						SchemaProps: spec.SchemaProps{
							Description: "Main container to run flux (only should be one)",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"volumes": {
						SchemaProps: spec.SchemaProps{
							Description: "Volumes that can be mounted (must be defined in volumes)",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("./api/v1alpha1/.ContainerVolume"),
									},
								},
							},
						},
					},
					"fluxOptionFlags": {
						SchemaProps: spec.SchemaProps{
							Description: "Flux option flags, usually provided with -o optional - if needed, default option flags for the server These can also be set in the user interface to override here. This is only valid for a FluxRunner \"runFlux\" true",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"fluxLogLevel": {
						SchemaProps: spec.SchemaProps{
							Description: "Log level to use for flux logging (only in non TestMode)",
							Default:     6,
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"preCommand": {
						SchemaProps: spec.SchemaProps{
							Description: "Special command to run at beginning of script, directly after asFlux is defined as sudo -u flux -E (so you can change that if desired.) This is only valid if FluxRunner is set (that writes a wait.sh script) This is for the indexed job pods and the certificate generation container.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"lifeCycle": {
						SchemaProps: spec.SchemaProps{
							Description: "Lifecycle can handle post start commands, etc.",
							Default:     map[string]interface{}{},
							Ref:         ref("./api/v1alpha1/.LifeCycle"),
						},
					},
					"resources": {
						SchemaProps: spec.SchemaProps{
							Description: "Resources include limits and requests",
							Default:     map[string]interface{}{},
							Ref:         ref("./api/v1alpha1/.ContainerResources"),
						},
					},
					"commands": {
						SchemaProps: spec.SchemaProps{
							Description: "More specific or detailed commands for just workers/broker",
							Default:     map[string]interface{}{},
							Ref:         ref("./api/v1alpha1/.Commands"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./api/v1alpha1/.Commands", "./api/v1alpha1/.ContainerResources", "./api/v1alpha1/.ContainerVolume", "./api/v1alpha1/.FluxUser", "./api/v1alpha1/.LifeCycle"},
	}
}

func schema__api_v1alpha1__MiniClusterList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "MiniClusterList contains a list of MiniCluster",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"),
						},
					},
					"items": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("./api/v1alpha1/.MiniCluster"),
									},
								},
							},
						},
					},
				},
				Required: []string{"items"},
			},
		},
		Dependencies: []string{
			"./api/v1alpha1/.MiniCluster", "k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"},
	}
}

func schema__api_v1alpha1__MiniClusterSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"containers": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "atomic",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Containers is one or more containers to be created in a pod. There should only be one container to run flux with runFlux",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("./api/v1alpha1/.MiniClusterContainer"),
									},
								},
							},
						},
					},
					"users": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "atomic",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Users of the MiniCluster",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("./api/v1alpha1/.MiniClusterUser"),
									},
								},
							},
						},
					},
					"jobLabels": {
						SchemaProps: spec.SchemaProps{
							Description: "Labels for the job",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
					"interactive": {
						SchemaProps: spec.SchemaProps{
							Description: "Run a single-user, interactive minicluster",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"volumes": {
						SchemaProps: spec.SchemaProps{
							Description: "Volumes accessible to containers from a host",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("./api/v1alpha1/.MiniClusterVolume"),
									},
								},
							},
						},
					},
					"logging": {
						SchemaProps: spec.SchemaProps{
							Description: "Logging modes determine the output you see in the job log",
							Default:     map[string]interface{}{},
							Ref:         ref("./api/v1alpha1/.LoggingSpec"),
						},
					},
					"fluxRestful": {
						SchemaProps: spec.SchemaProps{
							Description: "Customization to Flux Restful API There should only be one container to run flux with runFlux",
							Default:     map[string]interface{}{},
							Ref:         ref("./api/v1alpha1/.FluxRestful"),
						},
					},
					"cleanup": {
						SchemaProps: spec.SchemaProps{
							Description: "Cleanup the pods and storage when the index broker pod is complete",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"size": {
						SchemaProps: spec.SchemaProps{
							Description: "Size (number of job pods to run, size of minicluster in pods)",
							Default:     1,
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"tasks": {
						SchemaProps: spec.SchemaProps{
							Description: "Total number of CPUs being run across entire cluster",
							Default:     1,
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"deadlineSeconds": {
						SchemaProps: spec.SchemaProps{
							Description: "Should the job be limited to a particular number of seconds? Approximately one year. This cannot be zero or job won't start",
							Default:     3.15e+07,
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"pod": {
						SchemaProps: spec.SchemaProps{
							Description: "Pod spec details",
							Default:     map[string]interface{}{},
							Ref:         ref("./api/v1alpha1/.PodSpec"),
						},
					},
				},
				Required: []string{"containers"},
			},
		},
		Dependencies: []string{
			"./api/v1alpha1/.FluxRestful", "./api/v1alpha1/.LoggingSpec", "./api/v1alpha1/.MiniClusterContainer", "./api/v1alpha1/.MiniClusterUser", "./api/v1alpha1/.MiniClusterVolume", "./api/v1alpha1/.PodSpec"},
	}
}

func schema__api_v1alpha1__MiniClusterStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "MiniClusterStatus defines the observed state of Flux",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"jobid": {
						SchemaProps: spec.SchemaProps{
							Description: "The Jobid is set internally to associate to a miniCluster This isn't currently in use, we only have one!",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"conditions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "atomic",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "conditions hold the latest Flux Job and MiniCluster states",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.Condition"),
									},
								},
							},
						},
					},
				},
				Required: []string{"jobid"},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1.Condition"},
	}
}

func schema__api_v1alpha1__MiniClusterUser(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"name": {
						SchemaProps: spec.SchemaProps{
							Description: "If a user is defined, the username is required",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"password": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
				},
				Required: []string{"name"},
			},
		},
	}
}

func schema__api_v1alpha1__MiniClusterVolume(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Mini Cluster local volumes available to mount (these are on the host)",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"path": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
					"labels": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
					"annotations": {
						SchemaProps: spec.SchemaProps{
							Description: "Annotations for persistent volume claim",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
					"attributes": {
						SchemaProps: spec.SchemaProps{
							Description: "Optional volume attributes",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
					"volumeHandle": {
						SchemaProps: spec.SchemaProps{
							Description: "Volume handle, falls back to storage class name if not defined",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"storageClass": {
						SchemaProps: spec.SchemaProps{
							Default: "hostpath",
							Type:    []string{"string"},
							Format:  "",
						},
					},
					"driver": {
						SchemaProps: spec.SchemaProps{
							Description: "Storage driver, e.g., gcs.csi.ofek.dev Only needed if not using hostpath",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"delete": {
						SchemaProps: spec.SchemaProps{
							Description: "Delete the persistent volume on cleanup",
							Default:     true,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"secret": {
						SchemaProps: spec.SchemaProps{
							Description: "Secret reference in Kubernetes with service account role",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"secretNamespace": {
						SchemaProps: spec.SchemaProps{
							Description: "Secret namespace",
							Default:     "default",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"capacity": {
						SchemaProps: spec.SchemaProps{
							Description: "Capacity (string) for PVC (storage request) to create PV",
							Default:     "5Gi",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"path"},
			},
		},
	}
}

func schema__api_v1alpha1__PodSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "PodSpec controlls variables for the cluster pod",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"annotations": {
						SchemaProps: spec.SchemaProps{
							Description: "Annotations for each pod",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
					"labels": {
						SchemaProps: spec.SchemaProps{
							Description: "Labels for each pod",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
					"resources": {
						SchemaProps: spec.SchemaProps{
							Description: "Resources include limits and requests",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("k8s.io/apimachinery/pkg/util/intstr.IntOrString"),
									},
								},
							},
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/util/intstr.IntOrString"},
	}
}
