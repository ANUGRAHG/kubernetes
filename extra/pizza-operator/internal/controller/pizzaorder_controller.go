/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	restaurantv1 "pizza-operator/api/v1"
)

// PizzaOrderReconciler reconciles a PizzaOrder object
type PizzaOrderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=restaurant.example.com,resources=pizzaorders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=restaurant.example.com,resources=pizzaorders/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=restaurant.example.com,resources=pizzaorders/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PizzaOrder object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.0/pkg/reconcile
func (r *PizzaOrderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	// 1. Fetch the PizzaOrder the user submitted
	var pizza restaurantv1.PizzaOrder
	if err := r.Get(ctx, req.NamespacedName, &pizza); err != nil {
		// If it's not found, it means it was deleted. We just ignore it.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// 2. The Action: Do the work! (In real life, this might call a cloud API)
	log.Info("👨‍🍳 Chef received order!", "Pizza Name", pizza.Name)
	log.Info("🔥 Baking pizza...", "Size", pizza.Spec.Size, "Toppings", pizza.Spec.Toppings)
	// 3. The Status: Update Kubernetes on our progress
	pizza.Status.Progress = "Baking in the oven!"
	if err := r.Status().Update(ctx, &pizza); err != nil {
		log.Error(err, "Failed to update PizzaOrder status")
		return ctrl.Result{}, err
	}
	// Tell Kubernetes we are done and don't need to retry
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PizzaOrderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&restaurantv1.PizzaOrder{}).
		Named("pizzaorder").
		Complete(r)
}
