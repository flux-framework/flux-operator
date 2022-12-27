package controllers

// This file has extra (not used) functions that might be useful
// (and I didn't want to delete just yet)

import (
	"context"
	"os"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/client"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// podExec executes a command to a named pod
// This is not currenty in use. This seems to run but I don't see expected output
func (r *MiniClusterReconciler) podExec(pod corev1.Pod, ctx context.Context, cluster *api.MiniCluster) error {

	command := []string{
		"/bin/sh",
		"-c",
		"echo",
		"hello",
		"world",
	}

	// Prepare a request to execute to the pod in the statefulset
	execReq := r.RESTClient.Post().Namespace(cluster.Namespace).Resource("pods").
		Name(pod.Name).
		Namespace(cluster.Namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command:   command,
			Container: pod.Spec.Containers[0].Name,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, runtime.NewParameterCodec(r.Scheme))

	exec, err := remotecommand.NewSPDYExecutor(r.RESTConfig, "POST", execReq.URL())
	if err != nil {
		r.log.Error(err, "🌀 Error preparing command to execute to pod", "Name:", pod.Name)
		return err
	}

	// This is just for debugging for now :)
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: nil,
		Tty:    true,
	})
	r.log.Info("🌀 PodExec", "Container", pod.Spec.Containers[0].Name)
	return err
}

// getMiniClusterIPS was used when we needed to write /etc/hosts and is no longer used
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
