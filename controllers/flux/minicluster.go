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

	jobctrl "flux-framework/flux-operator/pkg/job"

	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	api "flux-framework/flux-operator/api/v1alpha1"
)

var (
	hostfileName   = "hostfile"
	curveCertKey   = "curve.cert"
	mungeMountName = "munge-key"
)

// This is a MiniCluster! A MiniCluster is associated with a running MiniCluster and include:
// 1. An indexed job with some number of pods
// 2. Config maps for secrets and other things.
// 3. We "launch" a job by starting the Indexed job on the connected nodes
// ensureMiniCluster creates a new MiniCluster, a stateful set for running flux!
func (r *MiniClusterReconciler) ensureMiniCluster(
	ctx context.Context,
	cluster *api.MiniCluster,
) (ctrl.Result, error) {

	// Ensure the configs are created (for volume sources)
	_, result, err := r.getConfigMap(ctx, cluster, "flux-config", cluster.Name+fluxConfigSuffix)
	if err != nil {
		return result, err
	}

	// Add initial config map with entrypoint scripts (wait.sh, start.sh, empty update_hosts.sh)
	_, result, err = r.getConfigMap(ctx, cluster, "entrypoint", cluster.Name+entrypointSuffix)
	if err != nil {
		return result, err
	}

	// Generate the curve certificate config map, unless already exists
	if cluster.Spec.Flux.CurveCertSecret != "" {
		_, result, err = r.getConfigMap(ctx, cluster, "cert", cluster.Name+curveVolumeSuffix)
		if err != nil {
			return result, err
		}
	}

	// Prepare volumes, if requested, to be available to containers
	for volumeName, volume := range cluster.Spec.Volumes {
		_, result, err = r.getPersistentVolume(ctx, cluster, volumeName, volume)
		if err != nil {
			return result, err
		}
		_, result, err = r.getPersistentVolumeClaim(ctx, cluster, volumeName, volume)
		if err != nil {
			return result, err
		}
	}

	// Any extra service containers (running alongside the cluster)
	// For now run these in the same pod, one service pod
	if len(cluster.Spec.Services) > 0 {
		_, result, err = r.ensureServicePod(ctx, cluster)
		if err != nil {
			return result, err
		}
	}

	// Create headless service for the MiniCluster OR single service for the broker
	selector := map[string]string{"job-name": cluster.Name}

	// If we are adding a minimal service to the index 0 pod only
	if cluster.Spec.Flux.MinimalService {
		selector = map[string]string{"job-index": "0"}
	}

	result, err = r.exposeServices(ctx, cluster, cluster.Spec.Network.HeadlessName, selector)
	if err != nil {
		return result, err
	}

	// Create the batch job that brings it all together!
	// A batchv1.Job can hold a spec for containers that use the configs we just made
	mc, result, err := r.getMiniCluster(ctx, cluster)
	if err != nil {
		return result, err
	}

	// If the sizes are different, we patch to update.
	// This would be an explicit update from a user or application via the CRD to scale up/down
	// The Flux Operator can't tell the difference between these two, but honors maxSize
	// (the size the Flux broker leader knows about) and updates both Spec.Size and Status.Size
	if *mc.Spec.Parallelism != cluster.Spec.Size {
		r.log.Info("MiniCluster", "Size", mc.Spec.Parallelism, "Requested Size", cluster.Spec.Size)
		result, err := r.resizeCluster(ctx, mc, cluster)
		if err != nil {
			return result, err
		}
	}

	// Add selector to allow horizontal pod autoscaler
	// This would be done via a request to a running metrics server
	// If there is no autoscaler, has no impact. The .Status.Size
	// should already be updated via the function above.
	result, err = r.addScaleSelector(ctx, cluster)
	if err != nil {
		return result, err
	}

	// Expose other sidecar container services
	for _, container := range cluster.Spec.Containers {

		// Assume now services only available TO flux runner
		if container.RunFlux || len(container.Ports) == 0 {
			continue
		}

		// Service name corresponds to container, but selector is pod-specific
		selector := map[string]string{"app.kubernetes.io/name": cluster.Name}
		result, err = r.exposeService(ctx, cluster, container.Name, selector, container.Ports)
		if err != nil {
			return result, err
		}
	}

	// Add the single label for the broker pod
	result, err = r.addBrokerLabel(ctx, cluster)
	if err != nil {
		return result, err
	}

	// If we get here, update the status to be ready
	status := jobctrl.GetCondition(cluster)
	if status != jobctrl.ConditionJobReady {
		clusterCopy := cluster.DeepCopy()
		jobctrl.FlagConditionReady(clusterCopy)
		r.Status().Update(ctx, clusterCopy)
	}

	// Determine if the job is completed, and flag the MiniCluster
	if !cluster.Status.Completed {
		if jobctrl.IsFinished(mc) {
			cluster.Status.Completed = true
			clusterCopy := cluster.DeepCopy()
			jobctrl.FlagConditionFinished(clusterCopy)
			r.Status().Update(ctx, clusterCopy)
		}
	}

	// And we re-queue so the Ready condition triggers next steps!
	return ctrl.Result{Requeue: true}, nil
}

// cleanupPodsStorage looks for the existing job, and cleans up if completed
func (r *MiniClusterReconciler) cleanupPodsStorage(
	ctx context.Context,
	cluster *api.MiniCluster,
) (ctrl.Result, error) {

	// Find the broker pod and determine if finished
	completed := false
	for _, pod := range r.getMiniClusterPods(ctx, cluster).Items {
		if !strings.HasPrefix(pod.Name, fmt.Sprintf("%s-0", cluster.Name)) {
			continue
		}

		// If it's succeeded or failed, we call that finished
		// https://pkg.go.dev/k8s.io/api@v0.25.0/core/v1#PodPhase
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			completed = true
			break
		}
	}

	// Cut out early if not completed
	if !completed {
		r.log.Info("MiniCluster", "Job Status", "Not Completed")
		return ctrl.Result{Requeue: true}, nil
	}
	r.log.Info("MiniCluster", "Job Status", "Completed")

	// Delete the MiniCluster first
	// If we don't, it will keep re-creating the assets and loop forever :)
	r.Client.Delete(ctx, cluster)

	// The job deletion should handle pods, next delete pvc and pv per each volume
	// Must be deleted in that order, per internet advice :)
	for volumeName := range cluster.Spec.Volumes {
		volumeSpec := cluster.Spec.Volumes[volumeName]

		claimName := fmt.Sprintf("%s-claim", volumeName)

		// Only delete if we retrieve without error and user has requested
		claim, err := r.getExistingPersistentVolumeClaim(ctx, cluster, claimName)
		if err != nil {
			r.log.Info("Volume Claim", "Deletion", claim.Name)
			r.Client.Delete(ctx, claim)
		}

		// Different request to delete
		if volumeSpec.Delete {
			pv, err := r.getExistingPersistentVolume(ctx, cluster, volumeName)
			if err != nil {
				r.log.Info("Volume", "Deletion", pv.Name)
				r.Client.Delete(ctx, pv)
			}
		}
	}
	return ctrl.Result{Requeue: false}, nil
}

// getExistingJob gets an existing job that matches the MiniCluster CRD
func (r *MiniClusterReconciler) getExistingJob(
	ctx context.Context,
	cluster *api.MiniCluster,
) (*batchv1.Job, error) {

	existing := &batchv1.Job{}
	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
		},
		existing,
	)
	return existing, err
}

// getMiniCluster does an actual check if we have a batch job in the namespace
func (r *MiniClusterReconciler) getMiniCluster(
	ctx context.Context,
	cluster *api.MiniCluster,
) (*batchv1.Job, ctrl.Result, error) {

	// Look for an existing job
	existing, err := r.getExistingJob(ctx, cluster)

	// Create a new job if it does not exist
	if err != nil {

		if errors.IsNotFound(err) {
			job, err := r.newMiniClusterJob(cluster)
			if err != nil {
				r.log.Error(
					err, "Failed to create new MiniCluster Batch Job",
					"Namespace:", job.Namespace,
					"Name:", job.Name,
				)
				return job, ctrl.Result{}, err
			}

			r.log.Info(
				"✨ Creating a new MiniCluster Batch Job ✨",
				"Namespace:", job.Namespace,
				"Name:", job.Name,
			)

			err = r.New(ctx, job)
			if err != nil {
				r.log.Error(
					err,
					"Failed to create new MiniCluster Batch Job",
					"Namespace:", job.Namespace,
					"Name:", job.Name,
				)
				return job, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return job, ctrl.Result{Requeue: true}, nil

		} else if err != nil {
			r.log.Error(err, "Failed to get MiniCluster Batch Job")
			return existing, ctrl.Result{}, err
		}

	} else {
		r.log.Info(
			"🎉 Found existing MiniCluster Batch Job 🎉",
			"Namespace:", existing.Namespace,
			"Name:", existing.Name,
		)
	}
	return existing, ctrl.Result{}, err
}

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
				brokerConfig, err := generateFluxConfig(cluster)
				if err != nil {
					return existing, ctrl.Result{Requeue: true}, err
				}
				data[hostfileName] = brokerConfig

			} else if configName == "cert" {

				// Use zeromq to generate the curve certificate
				curveCert, err := r.getCurveCert(ctx, cluster)
				if err != nil || curveCert == "" {
					return existing, ctrl.Result{Requeue: true}, err
				}
				data[curveCertKey] = curveCert

			} else if configName == "entrypoint" {

				// The main logic for generating the Curve certificate, start commands, is here
				// We create a custom script for each container that warrants one,
				// meaning a Flux Runner.
				for i, container := range cluster.Spec.Containers {
					if container.RunFlux {
						waitScriptID := fmt.Sprintf("wait-%d", i)
						waitScript, err := generateWaitScript(cluster, i)
						if err != nil {
							return existing, ctrl.Result{}, err
						}
						data[waitScriptID] = waitScript
					}
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

// generateHostlist for a specific size given the cluster namespace and a size
func generateHostlist(cluster *api.MiniCluster, size int32) string {

	// If we don't have a leadbroker address, we are at the root
	var hosts string
	if cluster.Spec.Flux.Bursting.LeadBroker.Address == "" {
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

	// Now regardless of where we are, we add the bursted jobs in the same order.
	// Any cluster with bursting must share all the bursted hosts across clusters
	// This ensures that the ranks line up
	for _, bursted := range cluster.Spec.Flux.Bursting.Clusters {
		burstedHosts := fmt.Sprintf("%s-[%s]", bursted.Name, generateRange(bursted.Size, 0))
		hosts = fmt.Sprintf("%s,%s", hosts, burstedHosts)
	}
	return hosts
}

// generateFluxConfig creates the broker.toml file used to boostrap flux
func generateFluxConfig(cluster *api.MiniCluster) (string, error) {

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

// generateWaitScript generates the main script to start everything up!
func generateWaitScript(cluster *api.MiniCluster, containerIndex int) (string, error) {

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
	t, err := template.New("wait-sh").Parse(waitToStartTemplate)
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
