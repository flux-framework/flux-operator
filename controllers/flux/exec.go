package controllers

import (
	"bytes"
	"context"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

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
