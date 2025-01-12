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
type API struct {
	Endpoint string            `json:"endpoint"`
	Params   map[string]string `json:"params"`
}

// CointracPipelineSpec defines the desired state of CointracPipeline.
type CointracPipelineSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	API      API    `json:"api"`
	Schedule string `json:"schedule"`
}

// CointracPipelineStatus defines the observed state of CointracPipeline.
type CointracPipelineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CointracPipeline is the Schema for the cointracpipelines API.
type CointracPipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CointracPipelineSpec   `json:"spec,omitempty"`
	Status CointracPipelineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CointracPipelineList contains a list of CointracPipeline.
type CointracPipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CointracPipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CointracPipeline{}, &CointracPipelineList{})
}
