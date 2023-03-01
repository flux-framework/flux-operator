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
	hostfileName      = "hostfile"
	certGeneratorName = "cert-generate"
	curveCertKey      = "curve-cert"
	curveCert         = ""
)

// This is a MiniCluster! A MiniCluster is associated with a running MiniCluster and include:
// 1. An indexed job with some number of pods
// 2. Config maps for secrets and other things.
// 3. We "launch" a job by starting the Indexed job on the connected nodes
// newMiniCluster creates a new mini cluster, a stateful set for running flux!
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

	// Generate the curve certificate config map. This creates a container to get
	// the curve certificate from
	_, result, err = r.getConfigMap(ctx, cluster, "cert", cluster.Name+curveVolumeSuffix)
	if err != nil {
		return result, err
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

	// Create the batch job that brings it all together!
	// A batchv1.Job can hold a spec for containers that use the configs we just made
	_, result, err = r.getMiniCluster(ctx, cluster)
	if err != nil {
		return result, err
	}

	// Expose pod index 0 service
	result, err = r.exposeServices(ctx, cluster)
	if err != nil {
		return result, err
	}

	// If we get here, update the status to be ready
	status := jobctrl.GetCondition(cluster)
	if status != jobctrl.ConditionJobReady {
		clusterCopy := cluster.DeepCopy()
		jobctrl.FlagConditionReady(clusterCopy)
		r.Client.Status().Update(ctx, clusterCopy)
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

	// Delete the mini cluster first
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

		if !volumeSpec.Delete {
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
	err := r.Client.Get(
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

			err = r.Client.Create(ctx, job)
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
	err := r.Client.Get(
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
				data[hostfileName] = generateFluxConfig(cluster)

			} else if configName == "cert" {

				// So we don't require read write many, we use an emphemeral pod to generate
				// the curve certificate. We don't care about this pod, we just want the cert!
				if curveCert == "" {
					curveCert, err := r.getCurveCert(ctx, cluster)
					if err != nil || curveCert == "" {
						return existing, ctrl.Result{Requeue: true}, err
					}
					data[curveCertKey] = curveCert
				}

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
			err = r.Client.Create(ctx, dep)
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

// generateFluxConfig creates the broker.toml file used to boostrap flux
func generateFluxConfig(cluster *api.MiniCluster) string {

	// Prepare suffix of fully qualified domain name
	fqdn := fmt.Sprintf("%s.%s.svc.cluster.local", restfulServiceName, cluster.Namespace)
	hosts := fmt.Sprintf("[%s]", generateRange(int(cluster.Spec.Size)))
	fluxConfig := fmt.Sprintf(brokerConfigTemplate, fqdn, cluster.Name, hosts)
	if cluster.MultiUser() {
		fluxConfig += brokerConfigJobManagerPlugin
	}
	return fluxConfig
}

// generateWaitScript generates the main script to start everything up!
// TODO try removing the extra wait logic and just sleep/retry
func generateWaitScript(cluster *api.MiniCluster, containerIndex int) (string, error) {

	// The first pod (0) should always generate the curve certificate
	container := cluster.Spec.Containers[containerIndex]
	mainHost := fmt.Sprintf("%s-0", cluster.Name)
	hosts := fmt.Sprintf("%s-[%s]", cluster.Name, generateRange(int(cluster.Spec.Size)))

	// Ensure our requested users each each have a password
	for i, user := range cluster.Spec.Users {
		cluster.Spec.Users[i].Password = getRandomToken(user.Password)

		// Passwords will be truncated to 8
		if len(cluster.Spec.Users[i].Password) > 8 {
			cluster.Spec.Users[i].Password = cluster.Spec.Users[i].Password[:8]
		}
	}

	// Only derive cores if > 1
	var cores int32
	if container.Cores > 1 {
		cores = container.Cores - 1
	}
	// The token uuid is the same across images
	wt := WaitTemplate{
		FluxUser:    getFluxUser(cluster.Spec.FluxRestful.Username),
		FluxToken:   getRandomToken(cluster.Spec.FluxRestful.Token),
		MainHost:    mainHost,
		Hosts:       hosts,
		Container:   container,
		Users:       cluster.Spec.Users,
		Size:        cluster.Spec.Size,
		Tasks:       cluster.Spec.Tasks,
		Cores:       cores,
		FluxRestful: cluster.Spec.FluxRestful,
		Logging:     cluster.Spec.Logging,
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
func generateRange(size int) string {
	var rangeString string
	if size == 1 {
		rangeString = "0"
	} else {
		rangeString = fmt.Sprintf("0-%d", size-1)
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
