/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	api "github.com/flux-framework/flux-operator/api/v1alpha1"
)

// getHostfileConfig gets an existing configmap, if it's done
func (r *MiniClusterReconciler) getConfigMap(
	ctx context.Context,
	cluster *api.MiniCluster,
	configName string,
	configFullName string,
) (*corev1.ConfigMap, ctrl.Result, error) {

	// Look for the config map by name
	existing := &corev1.ConfigMap{}
	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      configFullName,
			Namespace: cluster.Namespace,
		},
		existing,
	)

	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {

			// Data for the config map
			data := map[string]string{}

			// check if its broker.toml (the flux config)
			if configName == "flux-config" {
				brokerConfig, err := GenerateFluxConfig(cluster)
				if err != nil {
					return existing, ctrl.Result{Requeue: true}, err
				}
				data[HostfileName] = brokerConfig

			} else if configName == "cert" {

				// Use zeromq to generate the curve certificate
				curveCert, err := GetCurveCert(cluster)
				if err != nil || curveCert == "" {
					return existing, ctrl.Result{Requeue: true}, err
				}
				data[CurveCertKey] = curveCert
				r.log.Info("ConfigMap", "Curve Certificate", curveCert)

			} else if configName == "entrypoint" {

				// Get updated data with
				data, err = GenerateEntrypoints(cluster)
				if err != nil {
					return existing, ctrl.Result{}, err
				}
			}

			// Finally create the config map
			dep := r.createConfigMap(cluster, configFullName, data)
			r.log.Info(
				"✨ Creating MiniCluster ConfigMap ✨",
				"Type", configName,
				"Namespace", dep.Namespace,
				"Name", dep.Name,
			)
			err = r.New(ctx, dep)
			if err != nil {
				r.log.Error(
					err, "❌ Failed to create MiniCluster ConfigMap",
					"Type", configName,
					"Namespace", dep.Namespace,
					"Name", (*dep).Name,
				)
				return existing, ctrl.Result{}, err
			}

			// Successful - return and requeue
			return dep, ctrl.Result{Requeue: true}, nil

		} else if err != nil {
			r.log.Error(err, "Failed to get MiniCluster ConfigMap")
			return existing, ctrl.Result{}, err
		}

	} else {
		r.log.Info(
			"🎉 Found existing MiniCluster ConfigMap",
			"Type", configName,
			"Namespace", existing.Namespace,
			"Name", existing.Name,
		)
	}
	return existing, ctrl.Result{}, err
}

// GenerateEntrypoints generates the data structure (for config map) with entrypoint scripts
func GenerateEntrypoints(cluster *api.MiniCluster) (map[string]string, error) {
	data := map[string]string{}
	for i, container := range cluster.Spec.Containers {
		if container.RunFlux {
			waitScriptID := fmt.Sprintf("wait-%d", i)
			waitScript, err := generateEntrypointScript(cluster, i, "wait-sh", waitToStartTemplate)
			if err != nil {
				return data, err
			}
			data[waitScriptID] = waitScript
		}

		// Custom logic for a sidecar container alongside flux
		if container.GenerateEntrypoint() {
			startScriptID := fmt.Sprintf("start-%d", i)
			startScript, err := generateEntrypointScript(cluster, i, "start-sh", sidecarStartTemplate)
			if err != nil {
				return data, err
			}
			data[startScriptID] = startScript
		}
	}
	return data, nil
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

// generateFluxConfig creates the broker.toml file used to boostrap flux
func GenerateFluxConfig(cluster *api.MiniCluster) (string, error) {

	// If we have a config provided by user, use it.
	if cluster.Spec.Flux.BrokerConfig != "" {
		return cluster.Spec.Flux.BrokerConfig, nil
	}

	// Generate the broker.toml template, always up to the max size allowed
	fqdn := fmt.Sprintf("%s.%s.svc.cluster.local", cluster.Spec.Network.HeadlessName, cluster.Namespace)
	hosts := generateHostlist(cluster, cluster.Spec.MaxSize)

	bt := BrokerTemplate{
		Hosts:           hosts,
		FQDN:            fqdn,
		Spec:            cluster.Spec,
		ClusterName:     cluster.Name,
		FluxInstallRoot: cluster.FluxInstallRoot(),
	}

	t, err := template.New("broker-toml").Parse(brokerConfigTemplate)
	if err != nil {
		return "", err
	}

	var output bytes.Buffer
	if err := t.Execute(&output, bt); err != nil {
		return "", err
	}

	return output.String(), nil
}

// getRequiredRanks figures out the quorum that should be online for the cluster to start
func getRequiredRanks(cluster *api.MiniCluster) string {

	// Use the Flux default - all ranks must be online
	// Because our maximum size is == our starting size
	requiredRanks := ""
	if cluster.Spec.MaxSize == cluster.Spec.Size {
		return requiredRanks
	}
	// This is the quorum - the nodes required to be online - so we can start
	// This can be less than the MaxSize
	return generateRange(cluster.Spec.Size, 0)
}

// generateEntrypointScript generates an entrypoint script to start everything up!
func generateEntrypointScript(
	cluster *api.MiniCluster,
	containerIndex int,
	templateName string,
	templateScriptName string,
) (string, error) {

	container := cluster.Spec.Containers[containerIndex]
	mainHost := fmt.Sprintf("%s-0", cluster.Name)

	// The resources size must also match the max size in the cluster
	// This set of hosts explicitly gets provided to resources
	hosts := generateHostlist(cluster, cluster.Spec.MaxSize)

	// Ensure our requested users each each have a password
	for i, user := range cluster.Spec.Users {
		cluster.Spec.Users[i].Password = getRandomToken(user.Password)

		// Passwords will be truncated to 8
		if len(cluster.Spec.Users[i].Password) > 8 {
			cluster.Spec.Users[i].Password = cluster.Spec.Users[i].Password[:8]
		}
	}

	// Ensure Flux Restful has a secret key
	cluster.Spec.FluxRestful.SecretKey = getRandomToken(cluster.Spec.FluxRestful.SecretKey)

	// Only derive cores if > 1
	var cores int32
	if container.Cores > 1 {
		cores = container.Cores - 1
	}

	// Ensure if we have a batch command, it gets split up
	batchCommand := strings.Split(container.Command, "\n")

	// Required quorum - might be smaller than initial list if size != maxsize
	requiredRanks := getRequiredRanks(cluster)

	// The token uuid is the same across images
	wt := WaitTemplate{
		FluxUser:      getFluxUser(cluster.Spec.FluxRestful.Username),
		FluxToken:     getRandomToken(cluster.Spec.FluxRestful.Token),
		MainHost:      mainHost,
		Hosts:         hosts,
		Cores:         cores,
		Container:     container,
		Spec:          cluster.Spec,
		Batch:         batchCommand,
		RequiredRanks: requiredRanks,
	}
	t, err := template.New(templateName).Parse(templateScriptName)
	if err != nil {
		return "", err
	}

	var output bytes.Buffer
	if err := t.Execute(&output, wt); err != nil {
		return "", err
	}

	return output.String(), nil
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

// createConfigMap generates a config map with some kind of data
func (r *MiniClusterReconciler) createConfigMap(
	cluster *api.MiniCluster,
	configName string,
	data map[string]string,
) *corev1.ConfigMap {

	// Create the config map with respective data!
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configName,
			Namespace: cluster.Namespace,
		},
		Data: data,
	}

	// Show in the logs
	fmt.Println(cm.Data)
	ctrl.SetControllerReference(cluster, cm, r.Scheme)
	return cm
}
