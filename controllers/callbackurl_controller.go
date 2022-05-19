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

package controllers

import (
	"context"
	"fmt"
	"net/url"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/goern/r-gespraech/api/v1alpha1"
	erinnerungv1alpha1 "github.com/goern/r-gespraech/api/v1alpha1"
)

const (
	RequeueAfter = 10 * time.Second
)

// CallbackUrlReconciler reconciles a CallbackUrl object
type CallbackUrlReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	CallbackUrl *v1alpha1.CallbackUrl
}

//+kubebuilder:rbac:groups=erinnerung.thoth-station.ninka,resources=callbackpayloads,verbs=get;list;watch
//+kubebuilder:rbac:groups=erinnerung.thoth-station.ninka,resources=callbackpayloads/status,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=erinnerung.thoth-station.ninja,resources=callbackurls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=erinnerung.thoth-station.ninja,resources=callbackurls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=erinnerung.thoth-station.ninja,resources=callbackurls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CallbackUrl object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *CallbackUrlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	r.CallbackUrl = &v1alpha1.CallbackUrl{}
	if err := r.Get(ctx, req.NamespacedName, r.CallbackUrl); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Unable to fetch reconciled resource")
		return ctrl.Result{Requeue: true}, err
	}

	if !r.CallbackUrl.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("Resource being delete, skipping further reconcile.")
		return ctrl.Result{}, nil
	}

	// FIXME this is not persisted to the object?!
	r.CallbackUrl.Status.Phase = r.CallbackUrl.AggregatePhase()

	// TODO: this needs to be refactored into a validating webhook
	if r.CallbackUrl.Spec.URL == "" {
		r.SetCondition("URL", metav1.ConditionFalse, "EmptyUrl", "the provided URL is empty")
		return r.UpdateStatusNow(ctx, nil)
	}
	if _, err := url.Parse(r.CallbackUrl.Spec.URL); err != nil {
		logger.Error(err, "URL not parsable")
		return r.UpdateStatusNow(ctx, err)
	} else {
		r.SetCondition("URL", metav1.ConditionTrue, "GoodUrl", "the provided URL good")
	}

	var associatedPayloads erinnerungv1alpha1.CallbackPayloadList
	payloadSelector, err := metav1.LabelSelectorAsSelector(&r.CallbackUrl.Spec.Selector)
	if err != nil {
		return r.UpdateStatusNow(ctx, err)
	}

	options := client.ListOptions{
		LabelSelector: payloadSelector,
		Namespace:     req.Namespace,
		Raw:           &metav1.ListOptions{},
	}

	if err := r.List(ctx, &associatedPayloads, &options); err != nil {
		logger.Error(err, "unable to list associated CallbackPayloads")
		return r.UpdateStatusNow(ctx, err)
	}

	if len(associatedPayloads.Items) == 0 {
		meta.RemoveStatusCondition(&r.CallbackUrl.Status.Conditions, "AssociatedPayloads") // TODO err handler
		r.SetCondition("NoAssociatedPayloads", metav1.ConditionTrue, "NoAssociatedPayloads", "there is not associated CallbackPayload for this CallbackURL")
	} else {
		meta.RemoveStatusCondition(&r.CallbackUrl.Status.Conditions, "NoAssociatedPayloads") // TODO err handler
		r.SetCondition("AssociatedPayloads", metav1.ConditionTrue, "AssociatedPayloads", fmt.Sprintf("there is %v associated CallbackPayload for this CallbackURL", len(associatedPayloads.Items)))
	}

	// now we know we have some payloads associated with this url, let's see if we need to send a payload
	var unsendPayloads []*v1alpha1.CallbackPayload
	var sentPayloads []*v1alpha1.CallbackPayload

	isPayloadSent := func(p *v1alpha1.CallbackPayload) (bool, v1alpha1.CallbackPayloadConditionType) {
		for _, c := range p.Status.Conditions {
			if (c.Type == v1alpha1.CallbackPayloadComplete || c.Type == v1alpha1.CallbackPayloadFailed) && c.Status == corev1.ConditionTrue {
				return true, c.Type
			}
		}

		return false, ""
	}

	for i, p := range associatedPayloads.Items {
		_, finishedType := isPayloadSent(&p)

		switch finishedType {
		case "": // ongoing
			unsendPayloads = append(unsendPayloads, &associatedPayloads.Items[i])
		case v1alpha1.CallbackPayloadComplete:
		case v1alpha1.CallbackPayloadFailed:
			sentPayloads = append(sentPayloads, &associatedPayloads.Items[i])
		}

	}

	// let's send out the unsent payloads
	for _, p := range unsendPayloads {
		logger.WithValues("unsentPayload", p.ObjectMeta).Info("unsent")
		// TODO
		// p.SetStatusCondition(v1alpha1.CallbackPayloadSending, metav1.ConditionTrue)
	}

	return r.UpdateStatusNow(ctx, nil)
}

// SetupWithManager sets up the controller with the Manager.
func (r *CallbackUrlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&erinnerungv1alpha1.CallbackUrl{}).
		Watches(
			&source.Kind{Type: &erinnerungv1alpha1.CallbackPayload{}},
			handler.EnqueueRequestsFromMapFunc(r.findObjectsCallbackPayload),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Complete(r)
}

// findObjectsCallbackPayload is getting a []reconcile.Reqeust based on the LabelSelector of the Payload
func (r *CallbackUrlReconciler) findObjectsCallbackPayload(payload client.Object) []reconcile.Request {
	var urls erinnerungv1alpha1.CallbackUrlList

	// TODO is this enough of an sanity check?!
	if r.CallbackUrl == nil {
		return []reconcile.Request{}
	}

	// 1. get all CallbackURLs based on the LabelSelector
	selector, _ := metav1.LabelSelectorAsSelector(&r.CallbackUrl.Spec.Selector) //TODO err handler
	options := client.ListOptions{
		LabelSelector: selector,
		Namespace:     payload.GetNamespace(), // TODO double check if this is the way to go
		Raw:           &metav1.ListOptions{},
	}

	if err := r.List(context.TODO(), &urls, &options); err != nil {
		// quietly return nothing and ignore the error
		return []reconcile.Request{}
	}

	// 2. let's create the ist of reconcile.Requests{}
	requests := make([]reconcile.Request, len(urls.Items))
	for i, item := range urls.Items {
		requests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      item.GetName(),
				Namespace: item.GetNamespace(),
			},
		}
	}

	// 3. return all the reconcile.Requests{}
	return requests
}

// Force object status update. Returns a reconcile result
func (r *CallbackUrlReconciler) UpdateStatusNow(ctx context.Context, originalErr error) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	if err := r.Status().Update(ctx, r.CallbackUrl); err != nil {
		logger.WithValues("reason", err.Error()).Info("Unable to update status, retrying")
		return ctrl.Result{Requeue: true}, nil
	}
	if originalErr != nil {
		return ctrl.Result{RequeueAfter: RequeueAfter}, originalErr
	} else {
		return ctrl.Result{}, nil
	}
}

// Set status condition helper
func (r *CallbackUrlReconciler) SetCondition(conditionType string, status metav1.ConditionStatus, reason, message string) {
	meta.SetStatusCondition(&r.CallbackUrl.Status.Conditions, metav1.Condition{
		Type:    conditionType,
		Status:  status,
		Reason:  reason,
		Message: message,
	})
}
