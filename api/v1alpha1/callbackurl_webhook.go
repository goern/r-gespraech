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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var callbackurllog = logf.Log.WithName("callbackurl-resource")

func (r *CallbackUrl) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-erinnerung-thoth-station-ninja-v1alpha1-callbackurl,mutating=false,failurePolicy=fail,sideEffects=None,groups=erinnerung.thoth-station.ninja,resources=callbackurls,verbs=create;update,versions=v1alpha1,name=vcallbackurl.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &CallbackUrl{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *CallbackUrl) ValidateCreate() error {
	callbackurllog.Info("validate create", "name", r.Name)

	return r.validateCallbackUrl()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *CallbackUrl) ValidateUpdate(old runtime.Object) error {
	callbackurllog.Info("validate update", "name", r.Name)

	return r.validateCallbackUrl()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *CallbackUrl) ValidateDelete() error {
	return nil
}

func (r *CallbackUrl) validateCallbackUrl() error {
	var allErrs field.ErrorList

	if err := r.validateCallbackUrlSpec(); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: ErinnerungGroupName, Kind: "CallbackUrl"},
		r.Name, allErrs)
}

func (r *CallbackUrl) validateCallbackUrlSpec() *field.Error {
	if _, err := url.Parse(r.Spec.URL); err != nil {
		return field.Invalid(field.NewPath("spec").Child("url"), r.Spec.URL, err.Error())
	}

	return nil
}
