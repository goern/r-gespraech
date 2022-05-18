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
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/goern/r-gespraech/api/v1alpha1"
	webhookv1alpha1 "github.com/goern/r-gespraech/api/v1alpha1"
)

// CallbackReconciler reconciles a Callback object
type CallbackReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Callback *v1alpha1.Callback
}

const (
	RequeueAfter = 10 * time.Second
)

//+kubebuilder:rbac:groups=webhook.thoth-station.ninja,resources=callbacks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webhook.thoth-station.ninja,resources=callbacks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webhook.thoth-station.ninja,resources=callbacks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Callback object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *CallbackReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	r.Callback = &v1alpha1.Callback{}

	if err := r.Get(ctx, req.NamespacedName, r.Callback); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Resource deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Unable to fetch reconciled resource")
		return ctrl.Result{Requeue: true}, err
	}

	if !r.Callback.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("Resource being delete, skipping further reconcile.")
		return ctrl.Result{}, nil
	}

	// TODO(user): your logic here
	if r.Callback.Spec.URL == "" {
		logger.Info("URL was ''")
		r.SetCondition("URL", metav1.ConditionFalse, "unparsableURL", "we cant parse the URL provided")
	}

	return r.UpdateStatusNow(ctx, nil)
}

// Set status condition helper
func (r *CallbackReconciler) SetCondition(conditionType string, status metav1.ConditionStatus, reason, message string) {
	meta.SetStatusCondition(&r.Callback.Status.Conditions, metav1.Condition{
		Type:               conditionType,
		Status:             status,
		LastTransitionTime: metav1.Time{Time: time.Now()},
		Reason:             reason,
		Message:            message,
	})
}

// SetupWithManager sets up the controller with the Manager.
func (r *CallbackReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webhookv1alpha1.Callback{}).
		Complete(r)
}

// Update object status. Returns a reconcile result
func (r *CallbackReconciler) UpdateStatusNow(ctx context.Context, originalErr error) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	if err := r.Status().Update(ctx, r.Callback); err != nil {
		logger.WithValues("reason", err.Error()).Info("Unable to update status, retrying")
		return ctrl.Result{Requeue: true}, nil
	}
	if originalErr != nil {
		return ctrl.Result{RequeueAfter: RequeueAfter}, originalErr
	} else {
		return ctrl.Result{}, nil
	}
}
