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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CallbackSpec defines the desired state of Callback
type CallbackSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// URL is the URL to be called back
	URL                string `json:"url"`
	DeliveryRetryLimit *int32 `json:"deliveryRetryLimit,omitempty"`
}

// CallbackStatus defines the observed state of Callback
type CallbackStatus struct {
	// Current condition of the Shower.
	//+operator-sdk:csv:customresourcedefinitions:type=status,displayName="Phase",xDescriptors={"urn:alm:descriptor:io.kubernetes.phase'"}
	//+optional
	Phase string `json:"phase,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Callback is the Schema for the callbacks API
type Callback struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CallbackSpec   `json:"spec,omitempty"`
	Status CallbackStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CallbackList contains a list of Callback
type CallbackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Callback `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Callback{}, &CallbackList{})
}
