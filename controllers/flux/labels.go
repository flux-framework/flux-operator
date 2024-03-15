/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	api "github.com/flux-framework/flux-operator/api/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

type patchPayload struct {
	Operation string `json:"op"`
	Path      string `json:"path"`
	Value     string `json:"value"`
}

func (r *MiniClusterReconciler) addBrokerLabel(
	ctx context.Context,
	cluster *api.MiniCluster,
) (ctrl.Result, error) {

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(r.RESTConfig)
	if err != nil {
		return ctrl.Result{}, err
	}

	// List pods, and only get ones with undefined labels
	pods, err := clientset.CoreV1().Pods(cluster.Namespace).List(ctx, metav1.ListOptions{LabelSelector: "!job-index"})
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, pod := range pods.Items {

		// Add labels to all pods to indicate job-index
		prefix := fmt.Sprintf("%s-", cluster.Name)
		podName := strings.Replace(pod.GetName(), prefix, "", 1)
		podIndex := strings.SplitN(podName, "-", 2)[0]

		payload := []patchPayload{{
			Operation: "add",
			Path:      "/metadata/labels/job-index",
			Value:     podIndex,
		}}
		payloadBytes, _ := json.Marshal(payload)

		_, err = clientset.CoreV1().Pods(pod.GetNamespace()).Patch(ctx, pod.GetName(), types.JSONPatchType, payloadBytes, metav1.PatchOptions{})
		if err != nil {
			return ctrl.Result{}, err
		}

	}
	return ctrl.Result{}, nil
}
