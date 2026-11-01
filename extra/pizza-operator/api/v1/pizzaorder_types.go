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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PizzaOrderSpec defines the desired state (What the user wants)
type PizzaOrderSpec struct {
	Size     string   `json:"size,omitempty"`
	Toppings []string `json:"toppings,omitempty"`
}

// PizzaOrderStatus defines the observed state (What the Chef reports back)
type PizzaOrderStatus struct {
	Progress string `json:"progress,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PizzaOrder is the Schema for the pizzaorders API
type PizzaOrder struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of PizzaOrder
	// +required
	Spec PizzaOrderSpec `json:"spec"`

	// status defines the observed state of PizzaOrder
	// +optional
	Status PizzaOrderStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// PizzaOrderList contains a list of PizzaOrder
type PizzaOrderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []PizzaOrder `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(SchemeGroupVersion, &PizzaOrder{}, &PizzaOrderList{})
		return nil
	})
}
