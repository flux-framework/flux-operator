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
		defs[swaggify(defName)] = val.Schema
	}
	swagger := spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger:     "2.0",
			Definitions: defs,
			Paths:       &spec.Paths{Paths: map[string]spec.PathItem{}},
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       "fluxoperator",
					Description: "Python SDK for Flux-Operator",
					Version:     version,
				},
			},
		},
	}

	jsonBytes, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(string(jsonBytes))
}

func swaggify(name string) string {
	name = strings.Replace(name, "github.com/flux-framework/flux-operator/api/v1alpha/", "", -1)
	name = strings.Replace(name, "../api/v1alpha1/.", "", -1)
	name = strings.Replace(name, "./api/v1alpha1/.", "", -1)
	name = strings.Replace(name, "k8s.io/apimachinery/pkg/runtime/v1/", "", -1)
	name = strings.Replace(name, "k8s.io/apimachinery/pkg/util/", "", -1)
	name = strings.Replace(name, "k8s.io/apimachinery/pkg/apis/meta/", "", -1)
	name = strings.Replace(name, "k8s.io/kubernetes/pkg/controller/", "", -1)
	name = strings.Replace(name, "k8s.io/client-go/listers/core/", "", -1)
	name = strings.Replace(name, "/", ".", -1)
	return name
}
