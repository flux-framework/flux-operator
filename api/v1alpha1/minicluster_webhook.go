/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var miniclusterlog = ctrl.Log.WithName("minicluster-webhook")

func (r *MiniCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-flux-framework-org-v1alpha1-minicluster,mutating=true,failurePolicy=fail,sideEffects=None,groups={flux-framework.org},resources=miniclusters,verbs=create;update,versions=v1alpha1,name=mminicluster.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &MiniCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *MiniCluster) Default() {
	miniclusterlog.Info("🌈 Setting defaults")
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-flux-framework-org-v1alpha1-minicluster,mutating=false,failurePolicy=fail,sideEffects=None,groups={flux-framework.org},resources=miniclusters,verbs=create;update,versions=v1alpha1,name=vminicluster.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &MiniCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MiniCluster) ValidateCreate() error {
	miniclusterlog.Info("🌈 Validating create")

	// TODO(user): fill in your validation logic upon object creation.
	return r.validateMiniCluster()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *MiniCluster) ValidateUpdate(old runtime.Object) error {
	miniclusterlog.Info("🌈 validating update")

	// TODO(user): fill in your validation logic upon object update.
	return r.validateMiniCluster()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *MiniCluster) ValidateDelete() error {
	miniclusterlog.Info("🌈 validating delete")

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *MiniCluster) validateMiniCluster() error {
	var allErrs field.ErrorList
	if err := r.validateMiniClusterSize(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		miniclusterlog.Info("🌈 validate succesful!")
		return nil
	}

	miniclusterlog.Info("⛈️ validate Failed!")
	return apierrors.NewInvalid(
		schema.GroupKind{Kind: "MiniCluster"},
		r.Name, allErrs)
}

func (r *MiniCluster) validateMiniClusterSize() *field.Error {
	if r.Spec.Size == 0 {
		return field.Invalid(field.NewPath("spec").Child("size"), r.Spec.Size, "Size must be greater than 0")
	}
	return nil
}
