/*
Copyright 2022 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"context"
	"fmt"
	"path"
	"sort"
	"strings"

	jobctrl "flux-framework/flux-operator/pkg/job"

	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// This is a MiniCluster! A MiniCluster is associated with a running MiniCluster and include:
// 1. An indexed job with some number of pods
// 2. Config maps for secrets and other things.
// 3. We "launch" a job by starting the Indexed job on the connected nodes
// newMiniCluster creates a new mini cluster, a stateful set for running flux!
func (r *MiniClusterReconciler) ensureMiniCluster(ctx context.Context, cluster *api.MiniCluster) (ctrl.Result, error) {

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

	// A local host (developer machine) does not support provisioning, so for the meantime we use a
	// persistent volume instead (running on same host)
	if cluster.Spec.LocalDeploy {
		r.log.Info("MiniCluster", "localDeploy", "true (persistent volume in /tmp)")
		_, result, err = r.getPersistentVolume(ctx, cluster, cluster.Name+curveVolumeSuffix)
		if err != nil {
			return result, err
		}

		// Otherwise we can ask for a persistent volume claim
		// (not running on the same host)
	} else {
		r.log.Info("MiniCluster", "localDeploy", "false (persistent volume claim)")
		_, result, err = r.getPersistentVolumeClaim(ctx, cluster, cluster.Name+curveVolumeSuffix)
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

	// Reconcile until pods ips are ready
	// In the pods, it's waiting to see the update_hosts.sh file to be written.
	// We can do this because ips are written on the first creation and don't change
	ips := r.getMiniClusterIPS(ctx, cluster)
	r.log.Info("MiniCluster", "ips", ips)

	// Continue reconciling until we have pod ips
	if len(ips) != int(cluster.Spec.Size) {
		return ctrl.Result{Requeue: true}, fmt.Errorf("IPs are not ready yet! Only found %d", len(ips))
	}

	// At this point we've created job pods that have a waiting entrypoint for the update_hosts.sh
	// to exist. This is where we update the ConfigMap so it exists
	// Yes, this is a hack. Better ideas appreciated!
	_, result, err = r.addDiscoveryHostsFile(ctx, cluster)
	if err != nil {
		return result, err
	}

	// Expose pod index 0 service
	result, err = r.exposeService(ctx, cluster)
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

// getMiniCluster does an actual check if we have a batch job in the namespace
func (r *MiniClusterReconciler) getMiniCluster(ctx context.Context, cluster *api.MiniCluster) (*batchv1.Job, ctrl.Result, error) {
	existing := &batchv1.Job{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: cluster.Name, Namespace: cluster.Namespace}, existing)
	if err != nil {
		if errors.IsNotFound(err) {
			job := r.newMiniClusterJob(cluster)
			r.log.Info("✨ Creating a new MiniCluster Batch Job ✨", "Namespace:", job.Namespace, "Name:", job.Name)
			err = r.Client.Create(ctx, job)
			if err != nil {
				r.log.Error(err, " Failed to create new MiniCluster Batch Job", "Namespace:", job.Namespace, "Name:", job.Name)
				return job, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return job, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			r.log.Error(err, "Failed to get MiniCluster Batch Job")
			return existing, ctrl.Result{}, err
		}
	} else {
		r.log.Info("🎉 Found existing MiniCluster Batch Job 🎉", "Namespace:", existing.Namespace, "Name:", existing.Name)
	}
	if cluster.Spec.LocalDeploy {
		saveDebugYaml(existing, "batch-job.yaml")
	}
	return existing, ctrl.Result{}, err
}

func (r *MiniClusterReconciler) getMiniClusterIPS(ctx context.Context, cluster *api.MiniCluster) map[string]string {

	ips := map[string]string{}
	for _, pod := range r.getMiniClusterPods(ctx, cluster).Items {
		// Skip init pods
		if strings.Contains(pod.Name, "init") {
			continue
		}

		// The pod isn't ready!
		if pod.Status.PodIP == "" {
			continue
		}
		ips[pod.Name] = pod.Status.PodIP
	}
	return ips
}

// getMiniClusterPods returns a sorted (by name) podlist in the MiniCluster
func (r *MiniClusterReconciler) getMiniClusterPods(ctx context.Context, cluster *api.MiniCluster) *corev1.PodList {

	podList := &corev1.PodList{}
	opts := []client.ListOption{
		client.InNamespace(cluster.Namespace),
	}
	err := r.Client.List(ctx, podList, opts...)
	if err != nil {
		r.log.Error(err, "🌀 Error listing MiniCluster pods", "Name:", podList)
		return podList
	}

	// Ensure they are consistently sorted by name
	sort.Slice(podList.Items, func(i, j int) bool {
		return podList.Items[i].Name < podList.Items[j].Name
	})
	return podList
}

// discoverHosts generates a file that the pod can use to discover hosts.
// We assume the pods are sorted by name for a consistent output!
func (r *MiniClusterReconciler) generateDiscoverHostsFile(cluster *api.MiniCluster, pods *corev1.PodList, ips map[string]string) string {
	content := "#!/bin/sh"

	// NOTE: host will is duplicated, if that makes things wonky.
	for _, pod := range pods.Items {

		// flux-sample-N-xxxx -> flux-sample-N
		hostname := strings.Join(strings.SplitN(pod.Name, "-", 4)[0:3], "-")
		ip_address := ips[pod.Name]
		fqdn := fmt.Sprintf("%s-%s.%s.svc.cluster.local", hostname, cluster.Name, cluster.Namespace)
		if ip_address == "" {
			continue
		}
		content = fmt.Sprintf("%s\necho %s 	%s	%s >> /etc/hosts", content, ip_address, fqdn, hostname)
	}

	// Add one more newline for better readability
	content += "\n"

	// This is wrapping the entrypoint, so the last command needs to take args and start flux
	// The last set of arguments from the call should be the container entrypoint
	r.log.Info("🌀 MiniCluster Discover Hosts", "/flux_operator/update_hosts.sh", content)
	return content
}

// getPersistentVolume creates the PVC claim for the curve certificate (to be written once)
func (r *MiniClusterReconciler) getPersistentVolumeClaim(ctx context.Context, cluster *api.MiniCluster, configFullName string) (*corev1.PersistentVolumeClaim, ctrl.Result, error) {

	existing := &corev1.PersistentVolumeClaim{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: configFullName, Namespace: cluster.Namespace}, existing)
	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {
			volume := r.createPersistentVolumeClaim(cluster, configFullName)
			r.log.Info("✨ Creating MiniCluster Mounted Volume ✨", "Type", configFullName, "Namespace", volume.Namespace, "Name", volume.Name)
			err = r.Client.Create(ctx, volume)
			if err != nil {
				r.log.Error(err, "❌ Failed to create MiniCluster Mounted Volume", "Type", configFullName, "Namespace", volume.Namespace, "Name", (*volume).Name)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return volume, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			r.log.Error(err, "Failed to get MiniCluster Mounted Volume")
			return existing, ctrl.Result{}, err
		}
	} else {
		r.log.Info("🎉 Found existing MiniCluster Mounted Volume", "Type", configFullName, "Namespace", existing.Namespace, "Name", existing.Name)
	}
	if cluster.Spec.LocalDeploy {
		saveDebugYaml(existing, configFullName+".yaml")
	}
	return existing, ctrl.Result{}, err
}

// getPersistentVolume creates the PV for the curve certificate (to be written once)
func (r *MiniClusterReconciler) getPersistentVolume(ctx context.Context, cluster *api.MiniCluster, configFullName string) (*corev1.PersistentVolume, ctrl.Result, error) {

	existing := &corev1.PersistentVolume{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: configFullName, Namespace: cluster.Namespace}, existing)
	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {
			volume := r.createPersistentVolume(cluster, configFullName)
			r.log.Info("✨ Creating MiniCluster Mounted Volume ✨", "Type", configFullName, "Namespace", volume.Namespace, "Name", volume.Name)
			err = r.Client.Create(ctx, volume)
			if err != nil {
				r.log.Error(err, "❌ Failed to create MiniCluster Mounted Volume", "Type", configFullName, "Namespace", volume.Namespace, "Name", (*volume).Name)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return volume, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			r.log.Error(err, "Failed to get MiniCluster Mounted Volume")
			return existing, ctrl.Result{}, err
		}
	} else {
		r.log.Info("🎉 Found existing MiniCluster Mounted Volume", "Type", configFullName, "Namespace", existing.Namespace, "Name", existing.Name)
	}
	if cluster.Spec.LocalDeploy {
		saveDebugYaml(existing, configFullName+".yaml")
	}
	return existing, ctrl.Result{}, err
}

// getHostfileConfig gets an existing configmap, if it's done
func (r *MiniClusterReconciler) getConfigMap(ctx context.Context, cluster *api.MiniCluster, configName string, configFullName string) (*corev1.ConfigMap, ctrl.Result, error) {

	existing := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: configFullName, Namespace: cluster.Namespace}, existing)
	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {

			// Data for the config map
			data := map[string]string{}

			// check if its broker.toml (the flux config)
			if configName == "flux-config" {
				data["hostfile"] = generateFluxConfig(cluster)
			}

			// Initial "empty" set of start/wait scripts until we have host ips
			if configName == "entrypoint" {

				// The main logic for generating the Curve certificate, start commands, is here
				data["wait"] = generateWaitScript(cluster)

				// This will be updated after initial creation and we have host ips!
				data["update-hosts"] = ""
			}
			dep := r.createConfigMap(cluster, configFullName, data)
			r.log.Info("✨ Creating MiniCluster ConfigMap ✨", "Type", configName, "Namespace", dep.Namespace, "Name", dep.Name)
			err = r.Client.Create(ctx, dep)
			if err != nil {
				r.log.Error(err, "❌ Failed to create MiniCluster ConfigMap", "Type", configName, "Namespace", dep.Namespace, "Name", (*dep).Name)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return dep, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			r.log.Error(err, "Failed to get MiniCluster ConfigMap")
			return existing, ctrl.Result{}, err
		}
	} else {
		r.log.Info("🎉 Found existing MiniCluster ConfigMap", "Type", configName, "Namespace", existing.Namespace, "Name", existing.Name)
	}
	if cluster.Spec.LocalDeploy {
		saveDebugYaml(existing, configName+".yaml")
	}
	return existing, ctrl.Result{}, err
}

// generateFluxConfig creates the broker.toml file used to boostrap flux
func generateFluxConfig(cluster *api.MiniCluster) string {
	hosts := fmt.Sprintf("[%s]", generateRange(int(cluster.Spec.Size)))
	fluxConfig := fmt.Sprintf(brokerConfigTemplate, cluster.Name, hosts)
	return fluxConfig
}

// generateWaitScript generates the main script to start everything up!
func generateWaitScript(cluster *api.MiniCluster) string {

	// Generate a token uuid
	fluxToken := uuid.New()

	// The first pod (0) should always generate the curve certificate
	mainHost := fmt.Sprintf("%s-0", cluster.Name)
	hosts := fmt.Sprintf("%s-[%s]", cluster.Name, generateRange(int(cluster.Spec.Size)))
	waitScript := fmt.Sprintf(waitToStartTemplate, fluxToken.String(), mainHost, hosts, cluster.Spec.Diagnostics)
	return waitScript
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

// getHostfileConfig gets an existing configmap, if it's done
func (r *MiniClusterReconciler) addDiscoveryHostsFile(ctx context.Context, cluster *api.MiniCluster) (*corev1.ConfigMap, ctrl.Result, error) {

	configName := cluster.Name + entrypointSuffix
	cm := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: configName, Namespace: cluster.Namespace}, cm)

	// This is a bit redundant, but probably ok
	pods := r.getMiniClusterPods(ctx, cluster)
	ips := r.getMiniClusterIPS(ctx, cluster)

	// To update it we need to have found it
	if err == nil {

		cmCopy := cm.DeepCopy()
		cmCopy.Data["update-hosts"] = r.generateDiscoverHostsFile(cluster, pods, ips)
		err = r.Client.Update(ctx, cmCopy)
		if err != nil {
			r.log.Error(err, "❌ Error Adding Discovery Hosts File", "Namespace", cmCopy.Namespace, "Name", (*cmCopy).Name)
			return cmCopy, ctrl.Result{}, err
		}
		if cluster.Spec.LocalDeploy {
			saveDebugYaml(cmCopy, configName+".yaml")
		}
		return cmCopy, ctrl.Result{Requeue: true}, nil
	}
	return cm, ctrl.Result{}, err
}

// createConfigMap generates a config map with some kind of data
func (r *MiniClusterReconciler) createConfigMap(cluster *api.MiniCluster, configName string, data map[string]string) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configName,
			Namespace: cluster.Namespace,
		},
		Data: data,
	}
	fmt.Println(cm.Data)
	ctrl.SetControllerReference(cluster, cm, r.Scheme)
	return cm
}

// createPersistentVolumeClaim generates a PVC
// This tends to choke on MiniKube, I'm not sure it has a provisioner?
func (r *MiniClusterReconciler) createPersistentVolumeClaim(cluster *api.MiniCluster, configName string) *corev1.PersistentVolumeClaim {
	volume := &corev1.PersistentVolumeClaim{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: configName, Namespace: cluster.Namespace},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},

			// No idea how much to ask for here! I made it up.
			Resources: corev1.ResourceRequirements{Requests: corev1.ResourceList{
				corev1.ResourceStorage: *resource.NewQuantity(1024, resource.BinarySI),
			}},
		},
	}
	ctrl.SetControllerReference(cluster, volume, r.Scheme)
	return volume
}

// createPersistentVolume creates a volume in /tmp, which doesn't seem to choke
func (r *MiniClusterReconciler) createPersistentVolume(cluster *api.MiniCluster, configName string) *corev1.PersistentVolume {
	volume := &corev1.PersistentVolume{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: configName, Namespace: cluster.Namespace},
		Spec: corev1.PersistentVolumeSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
			Capacity: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceStorage: *resource.NewQuantity(1024, resource.BinarySI),
			},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: path.Join("/tmp/", configName),
				},
			},
			StorageClassName: "manual",
		},
	}
	ctrl.SetControllerReference(cluster, volume, r.Scheme)
	return volume
}
