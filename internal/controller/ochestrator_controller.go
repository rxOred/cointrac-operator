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

	katorv1 "github.com/kator/api/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// OchestratorReconciler reconciles a Ochestrator object
type OchestratorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *OchestratorReconciler) createExtractorStatus(ctx context.Context, req ctrl.Request, extractor katorv1.Extractor) error {
	log := log.FromContext(ctx)
	log.Info("createExtractorStatus", "Creating ExtractorStatus CR", extractor.Name)

	extractorName := extractor.Name
	existing := &katorv1.ExtractorStatus{}
	err := r.Get(ctx, client.ObjectKey{Name: extractorName, Namespace: req.Namespace}, existing)
	if err != nil {
		if errors.IsNotFound(err) {
			new := &katorv1.ExtractorStatus{
				ObjectMeta: metav1.ObjectMeta{
					Name:      extractorName,
					Namespace: req.Namespace,
				},
				Spec: katorv1.ExtractorStatusSpec{
					Name: extractorName,
				},
				Status: katorv1.ExtractorStatusStatus{
					Conditions: []metav1.Condition{
						{
							Type:    "Waiting",
							Status:  metav1.ConditionTrue,
							Reason:  "ExtractorStatusCreated",
							Message: "Created ExtractorStatus. Waiting for Deployments",
						},
					},
				},
			}

			if err := r.Client.Create(ctx, new); err != nil {
				log.Info("createExtractorStatus", "Failed to create CR")
				return err
			}
		} else if !errors.IsNotFound(err) {
			log.Info("createExtractorStatus", "ExtractorStatus already exist")
		} else {
			log.Info("createExtractorStatus", "Failed to get resources")
		}
	}

	return err
}

// +kubebuilder:rbac:groups=kator.my.domain,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kator.my.domain,resources=pipelines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kator.my.domain,resources=pipelines/finalizers,verbs=update
// +kubebuilder:rbac:groups=kator.my.domain,resources=extractorstatuses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kator.my.domain,resources=extractorstatuses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kator.my.domain,resources=extractorstatuses/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Ochestrator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *OchestratorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling Pipeline", "name", req.Name)

	var pipeline katorv1.Pipeline
	if err := r.Get(ctx, req.NamespacedName, &pipeline); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Resource not found. Ignoring...")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get Pipeline resource")
	}

	// Create Extractor CRs per Extractor in the pipeline
	extractors := pipeline.Spec.Extractors
	for _, extractor := range extractors {
		r.createExtractorStatus(ctx, req, extractor)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OchestratorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		// For().
		Named("ochestrator").
		Complete(r)
}
