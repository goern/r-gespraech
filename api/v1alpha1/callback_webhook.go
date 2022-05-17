/*
Copyright (C) 2022 Christoph Görn

This file is part of r-gespraech.

r-gespraech is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

r-gespraech is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with r-gespraech.  If not, see <http://www.gnu.org/licenses/>.
*/

package v1alpha1

import (
	"net/url"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var callbacklog = logf.Log.WithName("callback-resource")

func (r *Callback) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-webhook-thoth-station-ninja-v1alpha1-callback,mutating=true,failurePolicy=fail,sideEffects=None,groups=webhook.thoth-station.ninja,resources=callbacks,verbs=create;update,versions=v1alpha1,name=mcallback.kb.io,admissionReviewVersions=v1
var _ webhook.Defaulter = &Callback{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Callback) Default() {
	callbacklog.Info("default", "name", r.Name)

	if r.Spec.BackoffLimit == nil {
		r.Spec.BackoffLimit = new(int32)
		*r.Spec.BackoffLimit = 6
	}
}

//+kubebuilder:webhook:verbs=create;update;delete,path=/validate-webhook-thoth-station-ninja-v1alpha1-callback,mutating=false,failurePolicy=fail,groups=webhook.thoth-station.ninja,resources=callbacks,versions=v1,name=vcallback.kb.io,sideEffects=None,admissionReviewVersions=v1
var _ webhook.Validator = &Callback{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Callback) ValidateCreate() error {
	callbacklog.Info("validate create", "name", r.Name)

	return r.validateCallback()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Callback) ValidateUpdate(old runtime.Object) error {
	callbacklog.Info("validate update", "name", r.Name)

	return r.validateCallback()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Callback) ValidateDelete() error {
	callbacklog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *Callback) validateCallback() error {
	var allErrs field.ErrorList
	if err := r.validateCallbackSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "webhooks.thoth-station.ninja", Kind: "Callback"},
		r.Name, allErrs)
}

func (r *Callback) validateCallbackSpec() *field.Error {
	return validateURLFormat(
		r.Spec.URL,
		field.NewPath("spec").Child("URL"))
}

func validateURLFormat(uRL string, fldPath *field.Path) *field.Error {
	if _, err := url.Parse(uRL); err != nil {
		return field.Invalid(fldPath, uRL, err.Error())
	}
	return nil
}
