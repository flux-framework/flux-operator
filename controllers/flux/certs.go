/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

// This entire setup only exists because we want to generate the curve.cert with
// the flux in the container, which already has zeromq libraries installed.
// If we did this in Go, we would require the runner host to have them,
// which is a big ask. With this strategy we:
// 1. Define a global varible for the curve cert as an empty string
// 2. If it's empty, we generate a pod from the same Flux Runner container (with flux)
// 3. The pod config map entrypoint runs the command to generate the cert, and print to the terminal
// 4. We do not restart the pod on failure, but we can retrieve the curve.cert to populate the global variable
// 5. We use the curve.cert to write a config map for the batch job nodes to share.

// This likley isn't an elegant design - other ideas appreciated!

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	api "flux-framework/flux-operator/api/v1alpha1"
)

const (
	certGenSuffix = "-cert-generator"
)

// getCurveCert generates a pod to run a single command and get a curve certificate
func (r *MiniClusterReconciler) getCurveCert(ctx context.Context, cluster *api.MiniCluster) (string, error) {

	// Find the first Flux runner, we just need a container with flux to generate it
	// We have already validated at creation that we have at least one!
	var container api.MiniClusterContainer
	for _, contender := range cluster.Spec.Containers {
		if contender.RunFlux {
			container = contender
			break
		}
	}

	// This is a one time entrypoint to generate the flux curve certificate in a single pod
	_, _, err := r.getCurveGenerateConfigMap(ctx, cluster, container)
	if err != nil {
		return "", err
	}

	existing := &corev1.Pod{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: cluster.Name, Namespace: cluster.Namespace}, existing)
	if err != nil {
		command := []string{"/bin/bash", "/flux_operator/cert-generate.sh"}
		pod := r.newPodCommandRunner(cluster, container, command)
		r.log.Info("✨ Creating a new Pod Command Runner ✨", "Namespace:", pod.Namespace, "Name:", pod.Name)

		// We are being bad and not checking if there are errors - we just want to get the certificate
		r.Client.Create(ctx, pod)
		existing = pod
	}
	// If we get here, try to get the log output with the curve.cert
	curveCert, err := r.getPodLogs(ctx, existing)
	if curveCert != "" {
		fmt.Printf("🌵 Generated Curve Certificate\n%s\n", curveCert)
	}
	return curveCert, err
}

// generateCertGeneratorEntrypoint creates the entrypoint to create the curve.cert
func generateCertGeneratorEntrypoint(cluster *api.MiniCluster, container api.MiniClusterContainer) (string, error) {

	// We really only need any pre-command that adds flux to the path, etc.
	temp := CertTemplate{
		PreCommand: container.PreCommand,
	}
	t, err := template.New("cert-generate").Parse(generateCertTemplate)
	if err != nil {
		return "", err
	}

	var output bytes.Buffer
	if err := t.Execute(&output, temp); err != nil {
		return "", err
	}
	return output.String(), nil
}

// getCurveGenerateConfigMap generate the config map entrypoint for the pod to generate the curve cert
func (r *MiniClusterReconciler) getCurveGenerateConfigMap(
	ctx context.Context,
	cluster *api.MiniCluster,
	container api.MiniClusterContainer,
) (*corev1.ConfigMap, ctrl.Result, error) {

	existing := &corev1.ConfigMap{}
	configFullName := cluster.Name + certGenSuffix
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

			data := map[string]string{}
			genScript, err := generateCertGeneratorEntrypoint(cluster, container)
			if err != nil {
				return existing, ctrl.Result{}, err
			}
			data[certGeneratorName] = genScript
			dep := r.createConfigMap(cluster, configFullName, data)

			r.log.Info(
				"✨ Creating Curve Certificate Pod Generator Entrypoint ✨",
				"Namespace", dep.Namespace,
				"Name", dep.Name,
			)
			err = r.Client.Create(ctx, dep)
			if err != nil {
				r.log.Error(
					err, "❌ Failed to create Curve Certificate Pod Generator Entrypoint",
					"Namespace", dep.Namespace,
					"Name", (*dep).Name,
				)
				return existing, ctrl.Result{}, err
			}
			// Successful - return and requeue
			return dep, ctrl.Result{Requeue: true}, nil

		} else if err != nil {
			r.log.Error(err, "Failed to get Curve Certificate Pod Generator Entrypoint")
			return existing, ctrl.Result{}, err
		}
	} else {
		r.log.Info(
			"🎉 Found existing Curve Certificate Pod Generator Entrypoint",
			"Namespace", existing.Namespace,
			"Name", existing.Name,
		)
	}
	return existing, ctrl.Result{}, err
}

// createPersistentVolume creates a volume in /tmp, which doesn't seem to choke
func (r *MiniClusterReconciler) newPodCommandRunner(
	cluster *api.MiniCluster,
	container api.MiniClusterContainer,
	command []string,
) *corev1.Pod {

	makeExecutable := int32(0777)
	pullPolicy := corev1.PullIfNotPresent
	if container.PullAlways {
		pullPolicy = corev1.PullAlways
	}

	// Since the hostname needs to match the broker, we find the flux runner
	var containerName string
	for i, container := range cluster.Spec.Containers {
		if container.RunFlux {
			containerName = fmt.Sprintf("%s-%d", cluster.Name, i)
		}
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name + certGenSuffix,
			Namespace: cluster.Namespace,
		},
		Spec: corev1.PodSpec{
			RestartPolicy:    corev1.RestartPolicyOnFailure,
			ImagePullSecrets: getImagePullSecrets(cluster),
			Volumes: []corev1.Volume{{
				Name: cluster.Name + certGenSuffix,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: cluster.Name + certGenSuffix,
						},
						Items: []corev1.KeyToPath{{
							Key:  "cert-generate",
							Path: "cert-generate.sh",
							Mode: &makeExecutable,
						}},
					},
				},
			}},
			Containers: []corev1.Container{{
				Name:            containerName,
				Image:           container.Image,
				ImagePullPolicy: pullPolicy,
				WorkingDir:      container.WorkingDir,
				Stdin:           true,
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      cluster.Name + certGenSuffix,
						MountPath: "/flux_operator/",
						ReadOnly:  true,
					}},
				TTY:     true,
				Command: command,
			}},
		},
	}
	ctrl.SetControllerReference(cluster, pod, r.Scheme)
	return pod
}
