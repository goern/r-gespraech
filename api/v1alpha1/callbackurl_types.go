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

const (
	PhaseFailed           = "Failed"
	PhaseAwaitingPayloads = "AwaitingPayloads"
	PhasePending          = "Pending"
	PhaseOk               = "Ready"
)

// CallbackUrlSpec defines the desired state of CallbackUrl
type CallbackUrlSpec struct {
	// Url is the Url to call back.
	URL      string               `json:"url"`
	Selector metav1.LabelSelector `json:"selector"`
}

// CallbackUrlStatus defines the observed state of CallbackUrl
type CallbackUrlStatus struct {
	// Status is and aggregated view of the Conditions
	//+operator-sdk:csv:customresourcedefinitions:type=status,displayName="Phase",xDescriptors={"urn:alm:descriptor:io.kubernetes.phase'"}
	//+optional
	Phase string `json:"phase,omitempty"`

	// Conditions is the list of error conditions for this resource
	//+operator-sdk:csv:customresourcedefinitions:type=status,displayName="Conditions",xDescriptors={"urn:alm:descriptor:io.kubernetes.conditions"}
	//+optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
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

// Aggregate phase from conditions
func (u *CallbackUrl) AggregatePhase() string {
	if len(u.Status.Conditions) == 0 {
		return PhasePending
	}

	for _, c := range u.Status.Conditions {
		switch c.Type {
		case "NoAssociatedPayloads":
			if c.Status == metav1.ConditionTrue {
				return PhaseAwaitingPayloads
			}
		}
	}
	return PhaseOk
}

func init() {
	SchemeBuilder.Register(&CallbackUrl{}, &CallbackUrlList{})
}
