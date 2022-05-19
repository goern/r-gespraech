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
	"net/url"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CallbackUrlSpec defines the desired state of CallbackUrl
type CallbackUrlSpec struct {
	// Url is the Url to call back.
	Url      url.URL              `json:"url"`
	Selector metav1.LabelSelector `json:"selector"`
}

// CallbackUrlStatus defines the observed state of CallbackUrl
type CallbackUrlStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CallbackUrl is a web service's URL to receive a Callback. The Callback Payload to be send to the
// web service is determined via the metav1.LabelSelector `selector`.
type CallbackUrl struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CallbackUrlSpec   `json:"spec,omitempty"`
	Status CallbackUrlStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CallbackUrlList contains a list of CallbackUrl
type CallbackUrlList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CallbackUrl `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CallbackUrl{}, &CallbackUrlList{})
}
