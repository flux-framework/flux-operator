package controllers

import (
	api "flux-framework/flux-operator/api/v1alpha1"
)

// setup labels fetches and sets labels for setup
func setupLabels(v *api.FluxSetup, tier string) map[string]string {
	return map[string]string{
		"app":             "flux-workers",
		"visitorssite_cr": v.Name,
		"tier":            tier,
	}
}

// flux labels fetches and sets labels for Flux
func labels(v *api.Flux, tier string) map[string]string {
	return map[string]string{
		"app":             "flux-rank0",
		"visitorssite_cr": v.Name,
		"tier":            tier,
	}
}
