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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ExtractorReconciler reconciles a Extractor object
type ExtractorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ExtractorReconciler) updateExtractorStatus(ctx context.Context, req ctrl.Request, condition *metav1.Condition, extractor katorv1.Extractor) error {
	log := log.FromContext(ctx)
	log.Info("updateExtractorStatus", "Updating ExtractorStatus", condition.Status, extractor)

	extractorName := extractor.Name
	existing := &katorv1.ExtractorStatus{}
	err := r.Get(ctx, client.ObjectKey{Name: extractorName, Namespace: req.Namespace}, existing)
	if err != nil {
		log.Info("Failed to get ExtractorStatus for Extractor", extractorName)
		return err
	}

	existing.Status.Conditions[0] = metav1.Condition{}

	if err := r.Client.Create(ctx, existing); err != nil {
		log.Info(err.Error(), "updateExtractorStatus", "Failed to update ExtractorStatus", extractorName)
		return err
	}

	return nil
}

func (r *ExtractorReconciler) createExtractorDeployment(ctx context.Context, req ctrl.Request, extractor katorv1.Extractor) error {
	// extract data from the extractor and create the deployment
	log := log.FromContext(ctx)
	log.Info("createExtractorDeployment", "Creating Deployment for extractor", extractor.Name)

	var replicas int32 = 1
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      extractor.Name,
			Namespace: req.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": extractor.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": extractor.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  extractor.Name,
							Image: "default:latest",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
									Protocol:      corev1.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}

	if err := r.Client.Create(ctx, deployment); err != nil {
		log.Info("createExtractorDeployment", "Failed to create deployment for extractor", extractor.Name, err.Error())
		return err
	}

	log.Info("createExtractorDeployment", "Created Deployment for extractor", extractor.Name)
	return nil
}

// +kubebuilder:rbac:groups=kator.my.domain,resources=extractorstatuses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kator.my.domain,resources=extractorstatuses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kator.my.domain,resources=extractorstatuses/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Extractor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *ExtractorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling ExtractorStatus", "name", req.Name)

	var pipeline katorv1.ExtractorStatus
	if err := r.Get(ctx, req.NamespacedName, &pipeline); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Resource not found. Ignoring...")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get ExtractorStatus resource")
	}

	// check the status of ExtractorStatus
	// if created and in waiting state, create the deployments per extractor
	// TODO double check this
	if pipeline.Status.Conditions[0].Status == metav1.ConditionTrue && pipeline.Status.Conditions[0].Type == "Waiting" {
		var pipeline katorv1.Pipeline
		if err := r.Get(ctx, req.NamespacedName, &pipeline); err != nil {
			if errors.IsNotFound(err) {
				log.Info("Resource not found. Ignoring...")
				return ctrl.Result{}, nil
			}
			log.Error(err, "Failed to get Pipeline resource")
		}

		// collect the extractors
		// deploy the pods
		extractors := pipeline.Spec.Extractors
		for _, extractor := range extractors {
			if err := r.createExtractorDeployment(ctx, req, extractor); err != nil {
				log.Error(err, "Failed to create Deployment")
			}
			condition := &metav1.Condition{
				Type:    "Running",
				Status:  metav1.ConditionTrue,
				Reason:  "ExtractorRunning",
				Message: "Created Deployments for Extractor",
			}
			r.updateExtractorStatus(ctx, req, condition, extractor)
		}
	} else if pipeline.Status.Conditions[0].Status == metav1.ConditionTrue && pipeline.Status.Conditions[0].Type == "TransformReady" {

		// apply transform status CR with waiting
		var pipeline katorv1.Pipeline
		if err := r.Get(ctx, req.NamespacedName, &pipeline); err != nil {
			if errors.IsNotFound(err) {
				log.Info("Resource not found. Ignoring...")
				return ctrl.Result{}, nil
			}
			log.Error(err, "Failed to get Pipeline resource")
		}

		//transforms := pipeline.Spec.
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExtractorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		// For().
		Named("extractor").
		Complete(r)
}
