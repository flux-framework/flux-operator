/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package flux

import (
	"fmt"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"
	"github.com/google/uuid"
)

// getFluxUser returns a requested user name, or the default
func getFluxUser(requested string) string {
	if requested != "" {
		return requested
	}
	return "flux"
}

// getRandomToken returns a requested token, or a generated one
func getRandomToken(requested string) string {
	if requested != "" {
		return requested
	}
	return uuid.New().String()
}

// generateHostlist for a specific size given the cluster namespace and a size
func generateHostlist(cluster *api.MiniCluster, size int32) string {

	var hosts string
	if cluster.Spec.Flux.Bursting.Hostlist != "" {

		// In case 1, we are given a custom hostlist
		// This is usually the case when we are bursting to a different resource
		// Where the hostlists are not predictable.
		hosts = cluster.Spec.Flux.Bursting.Hostlist

	} else if cluster.Spec.Flux.Bursting.LeadBroker.Address == "" {

		// If we don't have a leadbroker address, we are at the root
		hosts = fmt.Sprintf("%s-[%s]", cluster.Name, generateRange(size, 0))

	} else {

		// Otherwise, we need to put the lead broker first, replacing the previous
		// index 0, and adding the rest of the range of jobs.
		// The hosts array must be consistent in ordering of ranks across workers
		adjustedSize := cluster.Spec.Flux.Bursting.LeadBroker.Size - 1
		hosts = fmt.Sprintf(
			"%s,%s-[%s]",
			cluster.Spec.Flux.Bursting.LeadBroker.Address,
			cluster.Spec.Flux.Bursting.LeadBroker.Name,

			// Index starts at 1
			generateRange(adjustedSize, 1),
		)
	}

	// For cases where the Flux Operator determines the hostlist, we need to
	// add the bursted jobs in the same order.
	// Any cluster with bursting must share all the bursted hosts across clusters
	// This ensures that the ranks line up
	if cluster.Spec.Flux.Bursting.Hostlist == "" {
		for _, bursted := range cluster.Spec.Flux.Bursting.Clusters {
			burstedHosts := fmt.Sprintf("%s-[%s]", bursted.Name, generateRange(bursted.Size, 0))
			hosts = fmt.Sprintf("%s,%s", hosts, burstedHosts)
		}
	}
	return hosts
}

// getRequiredRanks figures out the quorum that should be online for the cluster to start
func getRequiredRanks(cluster *api.MiniCluster) string {

	// Use the Flux default - all ranks must be online
	// Because our maximum size is == our starting size
	requiredRanks := ""

	// If the user has requested a custom, different number of ranks
	if cluster.Spec.MinSize != 0 {
		return fmt.Sprintf("%d", cluster.Spec.MinSize)
	}

	if cluster.Spec.MaxSize == cluster.Spec.Size {
		return requiredRanks
	}

	// This is the quorum - the nodes required to be online - so we can start
	// This can be less than the MaxSize
	return fmt.Sprintf("%d", cluster.Spec.Size)
}

// generateRange is a shared function to generate a range string
func generateRange(size int32, start int32) string {
	var rangeString string
	if size == 1 {
		rangeString = fmt.Sprintf("%d", start)
	} else {
		rangeString = fmt.Sprintf("%d-%d", start, (start+size)-1)
	}
	return rangeString
}
