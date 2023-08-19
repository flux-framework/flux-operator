/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"fmt"

	api "github.com/flux-framework/flux-operator/api/v1alpha1"
)

// GetCurveCert generates a pod to run a single command and get a curve certificate
func GetCurveCert(cluster *api.MiniCluster) (string, error) {
	if cluster.Spec.Flux.CurveCert != "" {
		return cluster.Spec.Flux.CurveCert, nil
	}
	curveCert, err := KeyGen("flux-cert-generator", fmt.Sprintf("%s-0", cluster.Name))
	return curveCert, err
}
