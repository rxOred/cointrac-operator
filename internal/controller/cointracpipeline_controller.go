/*
Copyright 2025.

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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cointracv1 "github.com/cointrac/api/v1"
)

// CointracPipelineReconciler reconciles a CointracPipeline object
type CointracPipelineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cointrac.cointrac.operator,resources=cointracpipelines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cointrac.cointrac.operator,resources=cointracpipelines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cointrac.cointrac.operator,resources=cointracpipelines/finalizers,verbs=update
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CointracPipeline object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *CointracPipelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling CointracPipeline", "name", req.Name)

	var pipeline cointracv1.CointracPipeline
	if err := r.Get(ctx, req.NamespacedName, &pipeline); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Cointrac Resource bot found. Ignoring...")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get CointracPipeline Resource")
	}

	endpoint := pipeline.Spec.API.Endpoint
	schedule := pipeline.Spec.Schedule

	log.Info("extracted data", "endpoint", endpoint, "schedule", schedule)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CointracPipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cointracv1.CointracPipeline{}).
		Named("cointracpipeline").
		Complete(r)
}
