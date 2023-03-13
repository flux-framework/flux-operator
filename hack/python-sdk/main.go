package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/validation/spec"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// Generate OpenAPI spec definitions for MPIJob Resource
func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Supply the Flux Operator API version")
	}
	version := os.Args[1]

	filter := func(name string) spec.Ref {
		return spec.MustCreateRef(
			"#/definitions/" + common.EscapeJsonPointer(swaggify(name)))
	}

	oAPIDefs := api.GetOpenAPIDefinitions(filter)
	defs := spec.Definitions{}
	for defName, val := range oAPIDefs {
		//		fmt.Println(val)
		defs[swaggify(defName)] = val.Schema
	}
	swagger := spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger:     "2.0",
			Info:        &spec.Info{InfoProps: spec.InfoProps{Title: "fluxoperator", Description: "Python SDK for Flux-Operator", Version: version}},
			Paths:       &spec.Paths{Paths: map[string]spec.PathItem{}},
			Definitions: defs,
			ExternalDocs: &spec.ExternalDocumentation{
				Description: "The Flux Operator",
				URL:         "https://flux-framework.org/flux-operator",
			},
		},
	}

	jsonBytes, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(string(jsonBytes))
}

// Our strategy here is to replace specific needs with classes we will define
func swaggify(name string) string {

	// These are specific to the Flux Operator
	name = strings.Replace(name, "github.com/flux-framework/flux-operator/api/v1alpha/", "", -1)
	name = strings.Replace(name, "../api/v1alpha1/.", "", -1)
	name = strings.Replace(name, "./api/v1alpha1/.", "", -1)

	// k8s.io/apimachinery/pkg/apis/meta/v1.Condition -> v1Condition
	name = strings.Replace(name, "k8s.io/apimachinery/pkg/apis/meta/v1.Condition", "v1Condition", -1)

	// k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta
	name = strings.Replace(name, "k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta", "v1ListMeta", -1)

	// k8s.io/apimachinery/pkg/util/intstr.IntOrString -> IntOrString
	name = strings.Replace(name, "k8s.io/apimachinery/pkg/util/intstr.", "", -1)

	// k8s.io/api/core/v1.SecurityContext -> v1SecurityContext
	name = strings.Replace(name, "k8s.io/api/core/v1.SecurityContext", "v1SecurityContext", -1)

	// k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta
	name = strings.Replace(name, "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta", "v1ObjectMeta", -1)
	return name
}
