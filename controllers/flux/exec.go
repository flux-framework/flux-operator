package controllers

// This file has extra (not used) functions that might be useful
// (and I didn't want to delete just yet)

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/client"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// podExec executes a command to a named pod
func (r *MiniClusterReconciler) podExec(ctx context.Context, pod corev1.Pod, cluster *api.MiniCluster, command []string) (string, error) {

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(r.RESTConfig)
	if err != nil {
		r.log.Info("🦀 PodExec", "Error with Creating Client", err)
		return "", err
	}

	// The flux runner has the same name as the namespace
	var container corev1.Container
	for _, contender := range pod.Spec.Containers {
		if strings.HasPrefix(contender.Name, cluster.Name) {
			container = contender
		}
	}

	// Prepare request TODO this will need to target flux runner
	req := clientset.CoreV1().RESTClient().Post().Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec").
		VersionedParams(
			&corev1.PodExecOptions{
				Container: container.Name,
				Command:   command,
				Stdin:     false,
				Stdout:    true,
				Stderr:    true,
				TTY:       true,
			},
			scheme.ParameterCodec,
		)

	exec, err := remotecommand.NewSPDYExecutor(r.RESTConfig, "POST", req.URL())
	if err != nil {
		r.log.Info("🦀 PodExec", "Error with Remote Command", err)
		return "", err
	}

	// Prepare buffers to stream to
	var stdout, stderr bytes.Buffer

	// Important! stdin must be none here so it isn't expecting our input
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		r.log.Info("🦀 PodExec", "Error with Stream", err)
		return "", err
	}

	outMsg := stdout.String()
	errMsg := stderr.String()

	fmt.Printf("🦀 PodExec Output\n%s", outMsg)
	if errMsg != "" {
		fmt.Printf("🦀 PodExec Error\n%s", errMsg)
	}
	return outMsg, err
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

// brokerIsReady determines if broker 0 is waiting for worker nodes
func (r *MiniClusterReconciler) brokerIsReady(ctx context.Context, cluster *api.MiniCluster) (bool, error) {

	brokerPrefix := fmt.Sprintf("%s-0", cluster.Name)

	makeReadyCommand := []string{
		"/bin/sh",
		"-c",
		"touch /etc/flux/ready",
	}

	// See if the file exists
	command := []string{
		"/bin/sh",
		"-c",
		"ls /etc/flux",
	}

	pods := r.getMiniClusterPods(ctx, cluster)
	for _, pod := range pods.Items {
		r.log.Info("🦀 Found Pod", "Pod Name", pod.Name)
		if strings.HasPrefix(pod.Name, brokerPrefix) {
			r.log.Info("🦀 Found Broker", "Pod Name", pod.Name)
			out, err := r.podExec(ctx, pod, cluster, command)

			// Right before the broker runs, it creates this file
			if !strings.Contains(out, "ready") || err != nil {
				return false, fmt.Errorf("broker is not ready")
			}
			r.log.Info("🦀 Broker Is Ready", "Pod Name", pod.Name)

			// Is the broker ready? If yes, touch files to indicate others ready
			for _, worker := range pods.Items {

				// Don't exec to the broker, pods that aren't for the job, or the certificate generator pod!
				if worker.Name == pod.Name || !strings.HasPrefix(worker.Name, cluster.Name) || strings.Contains(worker.Name, certGenSuffix) {
					continue
				}
				_, err := r.podExec(ctx, worker, cluster, makeReadyCommand)
				if err != nil {
					return false, err
				}
				r.log.Info("🦀 Worker", "Flagged Ready", worker.Name)
			}
			return true, nil
		}
	}
	return false, fmt.Errorf("broker is not ready")
}

// getPodLogs gets the pod logs (with the curve cert)
func (r *MiniClusterReconciler) getPodLogs(ctx context.Context, pod *corev1.Pod) (string, error) {

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(r.RESTConfig)
	if err != nil {
		return "", err
	}

	// Keep developer user informed what is going on.
	r.log.Info("Pod Logs", "Name", pod.Name)
	r.log.Info("Pod Logs", "Container", pod.Spec.Containers[0].Name)

	opts := corev1.PodLogOptions{
		Container: pod.Spec.Containers[0].Name,
	}

	// This will fail (and need to reconcile) while container is creating, etc.
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &opts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", err
	}
	logs := buf.String()
	return logs, err
}
