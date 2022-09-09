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
	"os"

	jobctrl "flux-framework/flux-operator/pkg/job"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	api "flux-framework/flux-operator/api/v1alpha1"
)

var (
	brokerConfigTemplate = `
[bootstrap]
curve_cert = "/mnt/curve/curve.cert"
default_port = 8050
default_bind = "tcp://eth0:%%p"
default_connect = "tcp://%%h:%%p"
hosts = [
	{ host="%s-%s"},
]
`
	// Dummy hostfile that we eventually need to generate dymically
	// 172.17.0.3      flux-sample-0
	dummyHostfile = `
flux-workers-0.flux-workers 10.0.0.1
flux-workers-1.flux-workers 10.0.0.2`
)

// This is a MiniCluster! A MiniCluster is associated with a running MiniCluster and include:
// 1. A stateful set with some number of pods
// 2. A service to expose the mini-cluster (still needed?)
// 3. Config maps for secrets and other things.
// 4. We "launch" a job by starting the Indexed job on the connected nodes
// newMiniCluster creates a new mini cluster, a stateful set for running flux!
func (r *MiniClusterReconciler) ensureMiniCluster(ctx context.Context, cluster *api.MiniCluster) (ctrl.Result, error) {

	// Ensure the configs are created (for volume sources)
	// The hostfile here is empty because we generate it entirely
	_, result, err := r.getHostfileConfig(ctx, cluster, "flux-config", "", cluster.Name+fluxConfigSuffix)
	if err != nil {
		return result, err
	}

	// This is using a dummy hostfile - this will need to be generated by way of
	// creating the pods, then going back and getting the hostnames and recreatingi t.
	_, result, err = r.getHostfileConfig(ctx, cluster, "etc-hosts", dummyHostfile, cluster.Name+etcHostsSuffix)
	if err != nil {
		return result, err
	}

	// Create the batch job that brings it all together!
	// A batchv1.Job can hold a spec for containers that use the configs
	// and secrets that we just created.
	_, result, err = r.getMiniCluster(ctx, cluster)
	if err != nil {
		return result, err
	}

	// Create the actual hostfile from the pods
	//hostfile := r.createEtcHosts(set, cluster)
	//r.log.Info("✨ Hostfile created", "Hostfile", hostfile)

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
			r.log.Error(err, "Failed to get Broker ConfigMap")
			return existing, ctrl.Result{}, err
		}
	} else {
		r.log.Info("🎉 Found existing Secret 🎉", "Namespace:", existing.Namespace, "Name:", existing.Name)
	}
	saveDebugYaml(existing, "batch-job.yaml")
	return existing, ctrl.Result{}, err
}

// LaunchJob is a wrapper to launch job, returning the reconcile result/error
// We aren't currently using this, but I thought it could be eventually useful if we want to
// exec a command to a specific container!
func (r *MiniClusterReconciler) LaunchJob(ctx context.Context, cluster *api.MiniCluster) (ctrl.Result, error) {

	// If the job command is empty, don't continue
	// This will keep the cluster running (for debugging?) but not exec a command
	if cluster.Spec.Command == "" {
		return ctrl.Result{Requeue: true}, nil
	}

	// Launch the job! Flag as running if successful
	err := r.fauxLaunchJob(ctx, cluster)
	jobCopy := cluster.DeepCopy()

	// No error, flag it running!
	if err == nil {
		jobctrl.FlagConditionRunning(jobCopy)
		r.log.Error(err, "🌀 Mini Cluster launched job!")

	} else {
		// This is a quasi flag to say "don't try running this again"
		jobCopy.Spec.Command = ""
		r.log.Error(err, "🌀 Mini Cluster Error launching job, try updating command.")
	}

	r.Client.Status().Update(ctx, jobCopy)
	return ctrl.Result{Requeue: true}, nil
}

// fauxLaunchJob is here just to fake launching a job, for the time being.
func (r *MiniClusterReconciler) fauxLaunchJob(ctx context.Context, cluster *api.MiniCluster) error {
	return nil
}

// launchJob actually tries executing the job command to a pod node
// TODO this isn't working yet - we don't have permission to execute on the pods
// could we figure out how to add the permission, or create a service to interact
// with instead?
func (r *MiniClusterReconciler) launchJob(ctx context.Context, cluster *api.MiniCluster) error {

	// Retrieve the stateful set, we will get the first pod name
	set := &appsv1.StatefulSet{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: cluster.Name, Namespace: cluster.Namespace}, set)
	if err != nil {
		r.log.Error(err, "🌀 Mini Cluster Error finding StatefulSet")
		return err
	}

	command := []string{
		"sh",
		"-c",
		"sleep 1",
	}

	// Programatically get pods
	pods := r.getMiniClusterPods(ctx, cluster)
	if len(pods.Items) == 0 {
		err = fmt.Errorf("No pods found in listing.")
		r.log.Error(err, "🌀 No pods are running in this MiniCluster, cannot launch a job.")
		return err
	}
	// Target the first pod to exec a command to
	pod := pods.Items[0]
	container := pod.Spec.Containers[0]
	r.log.Error(err, "🌀 Executing command to pod in statefulset", "Name:", pod.Name, "Container:", container.Name)

	// Prepare a request to execute to the pod in the statefulset
	execReq := r.RESTClient.Post().Namespace(cluster.Namespace).Resource("pods").
		Name(pod.GetName()).
		Namespace(cluster.Namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container.Name,
			Command:   command,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
		}, runtime.NewParameterCodec(r.Scheme))

	exec, err := remotecommand.NewSPDYExecutor(r.RESTConfig, "POST", execReq.URL())
	if err != nil {
		r.log.Error(err, "🌀 Error executing command to pod in statefulset", "Name:", pod.Name)
		return err
	}

	// This is just for debugging for now :)
	return exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})
}

// Get a list of string (expected) hostnames
func (r *MiniClusterReconciler) getHostnames(ctx context.Context, cluster *api.MiniCluster) []string {

	// We always should have one host
	hosts := []string{fmt.Sprintf("%s-0", cluster.Name)}

	// Generate list of hosts
	for i := 1; i < int(cluster.Spec.Size); i++ {
		hosts = append(hosts, fmt.Sprintf("%s-%d", cluster.Name, i))
	}
	return hosts
}

func (r *MiniClusterReconciler) getMiniClusterPods(ctx context.Context, cluster *api.MiniCluster) *corev1.PodList {

	podList := &corev1.PodList{}
	opts := []client.ListOption{
		client.InNamespace(cluster.Namespace),
		//		client.MatchingLabels{"instance": cluster.Name},
		//		client.MatchingFields{"status.phase": "Running"},
	}
	err := r.Client.List(ctx, podList, opts...)
	if err != nil {
		r.log.Error(err, "🌀 Error listing MiniCluster pods", "Name:", podList)
		return podList
	}

	// This is just for debugging
	for _, pod := range podList.Items {
		r.log.Error(err, "🌀 Found Pod", "Name:", pod.Name, "Container:", pod.Spec.Containers)
	}
	return podList
}

// This is what an actual pod in a stateful set entry looks like
// TODO we need to figure out how to generate an etc hosts so the pods
// can see one another! Here are some ideas;
// 1. Create a service DNS https://stackoverflow.com/questions/63415324/is-it-possible-for-a-pod-running-in-a-satrefulset-to-get-the-hostname-of-the-all
// 2. Create the pods, get the ips, and then somehow update them
//    This might be an issue if they are re-created (and the changes lost)
// 172.17.0.3      flux-sample-0
func (r *MiniClusterReconciler) createEtcHosts(set *appsv1.StatefulSet, cluster *api.MiniCluster) map[string]string {
	lookup := map[string]string{}
	return lookup
}

/*
apiVersion: v1
kind: Service
metadata:
  name: flux-workers
spec:
  clusterIP: None
  selector:
    app: flux-workers*/
func (r *MiniClusterReconciler) createService(cluster *api.MiniCluster) *corev1.Service {

	labels := setupLabels(cluster, "flux-workers")

	// We shouldn't need this, as the port comes from the manifest
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
		},
	}
	ctrl.SetControllerReference(cluster, service, r.Scheme)
	return service
}

// getHostfileConfig gets an existing configmap, if it's done
func (r *MiniClusterReconciler) getHostfileConfig(ctx context.Context, cluster *api.MiniCluster, configName string, hostfile string, configFullName string) (*corev1.ConfigMap, ctrl.Result, error) {

	existing := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: configFullName, Namespace: cluster.Namespace}, existing)
	if err != nil {

		// Case 1: not found yet, and hostfile is ready (recreate)
		if errors.IsNotFound(err) {
			// check if its broker.toml TODO : Convert all configMaps to use template strings
			if configName == "flux-config" {
				hostfile = generateFluxConfig(cluster.Name, cluster.Spec.Size)
			}
			dep := r.createHostfileConfig(cluster, configFullName, hostfile)
			r.log.Info("✨ Creating MiniCluster ConfigMap ✨", "Type", configName, "Namespace", dep.Namespace, "Name", dep.Name, "Data", (*dep).Data)
			err = r.Client.Create(ctx, dep)
			if err != nil {
				r.log.Error(err, "❌ Failed to create MiniCluster ConfigMap", "Type", configName, "Namespace", dep.Namespace, "Name", (*dep).Name)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return existing, ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			r.log.Error(err, "Failed to get MiniCluster ConfigMap")
			return existing, ctrl.Result{}, err
		}
	} else {
		r.log.Info("🎉 Found existing MiniCluster ConfigMap", "Type", configName, "Namespace", existing.Namespace, "Name", existing.Name, "Data", (*existing).Data)
	}
	saveDebugYaml(existing, configName+".yaml")
	return existing, ctrl.Result{}, err
}

// generateFluxConfig creates the broker.toml file used to boostrap flux
func generateFluxConfig(name string, size int32) string {
	var hosts string
	if size == 1 {
		hosts = "0"
	} else {
		hosts = fmt.Sprintf("[0-%d]", size-1)
	}
	fluxConfig := fmt.Sprintf(brokerConfigTemplate, name, hosts)

	return fluxConfig
}

// createBrokerConfig creates the stateful set
func (r *MiniClusterReconciler) createHostfileConfig(cluster *api.MiniCluster, configName string, hostfile string) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configName,
			Namespace: cluster.Namespace,
		},
		Data: map[string]string{
			"hostfile": hostfile,
		},
	}
	fmt.Println(cm.Data)
	ctrl.SetControllerReference(cluster, cm, r.Scheme)
	return cm
}
