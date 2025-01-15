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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ExtractorStatusSpec defines the desired state of ExtractorStatus.
type ExtractorStatusSpec struct {
	Name string `json:"name"`
}

// ExtractorStatusStatus defines the observed state of ExtractorStatus.
type ExtractorStatusStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ExtractorStatus is the Schema for the extractorstatuses API.
type ExtractorStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExtractorStatusSpec   `json:"spec,omitempty"`
	Status ExtractorStatusStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ExtractorStatusList contains a list of ExtractorStatus.
type ExtractorStatusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExtractorStatus `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ExtractorStatus{}, &ExtractorStatusList{})
}
