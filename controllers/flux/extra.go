package controllers

// This file has extra (not used) functions that might be useful
// (and I didn't want to delete just yet)!
// - pod exec commands
// - mini cluster ingress
// - pod listing, and ips
// - persistent volume claims

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	api "flux-framework/flux-operator/api/v1alpha1"
)

// Use a global variable for now to quickly return when broker ready
var (
	brokerIsReady = false
)

// podExec executes a command to a named pod
func (r *MiniClusterReconciler) podExec(
	ctx context.Context,
	pod corev1.Pod,
	cluster *api.MiniCluster,
	command []string,
) (string, error) {

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
func (r *MiniClusterReconciler) getMiniClusterPods(
	ctx context.Context,
	cluster *api.MiniCluster,
) *corev1.PodList {

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
func (r *MiniClusterReconciler) brokerIsReady(
	ctx context.Context,
	cluster *api.MiniCluster,
) (bool, error) {

	// Cut out quickly if we've already done this
	if brokerIsReady {
		return true, nil
	}
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
				if worker.Name == pod.Name || !strings.HasPrefix(worker.Name, cluster.Name) {
					continue
				}
				_, err := r.podExec(ctx, worker, cluster, makeReadyCommand)
				if err != nil {
					return false, err
				}
				r.log.Info("🦀 Worker", "Flagged Ready", worker.Name)
			}
			brokerIsReady = true
			return true, nil
		}
	}
	return false, fmt.Errorf("broker is not ready")
}

// getMiniClusterIPS was used when we needed to write /etc/hosts and is no longer used
func (r *MiniClusterReconciler) getMiniClusterIPS(
	ctx context.Context,
	cluster *api.MiniCluster,
) map[string]string {

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

// createMiniClusterIngress exposes the service for the minicluster
func (r *MiniClusterReconciler) createMiniClusterIngress(
	ctx context.Context,
	cluster *api.MiniCluster,
	service *corev1.Service,
) error {

	pathType := networkv1.PathTypePrefix
	ingressBackend := networkv1.IngressBackend{
		Service: &networkv1.IngressServiceBackend{
			Name: service.Name,
			Port: networkv1.ServiceBackendPort{
				Number: service.Spec.Ports[0].NodePort,
			},
		},
	}
	ingress := &networkv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
		},
		Spec: networkv1.IngressSpec{
			DefaultBackend: &ingressBackend,
			Rules: []networkv1.IngressRule{
				{
					IngressRuleValue: networkv1.IngressRuleValue{
						HTTP: &networkv1.HTTPIngressRuleValue{
							Paths: []networkv1.HTTPIngressPath{
								{
									PathType: &pathType,
									Backend:  ingressBackend,
									Path:     "/",
								},
							},
						},
					},
				},
			},
		},
	}
	err := ctrl.SetControllerReference(cluster, ingress, r.Scheme)
	if err != nil {
		r.log.Error(err, "🔴 Create ingress", "Service", cluster.Spec.ServiceName)
		return err
	}
	err = r.Client.Create(ctx, ingress)
	if err != nil {
		r.log.Error(err, "🔴 Create ingress", "Service", cluster.Spec.ServiceName)
		return err
	}
	return nil
}
