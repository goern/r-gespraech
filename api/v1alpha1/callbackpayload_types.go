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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CallbackPayloadSpec defines the desired state of CallbackPayload
type CallbackPayloadSpec struct {
	Data     string               `json:"data"`
	Selector metav1.LabelSelector `json:"selector"`
}

type CallbackPayloadConditionType string

// These are built-in conditions of a job.
const (
	// CallbackPayloadSending means that the payload is in the process of being send.
	CallbackPayloadSending CallbackPayloadConditionType = "Sending"
	// CallbackPayloadComplete means the payload has been successfully sent.
	CallbackPayloadComplete CallbackPayloadConditionType = "Complete"
	// CallbackPayloadFailed means the payload has failed sending.
	CallbackPayloadFailed CallbackPayloadConditionType = "Failed"
)

// CallbackPayloadCondition describes current state of a payload.
type CallbackPayloadCondition struct {
	// Type of condition, Complete or Failed.
	Type CallbackPayloadConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// Last time the condition was checked.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// Last time the condition transit from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// (brief) reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// CallbackPayloadStatus defines the observed state of CallbackPayload
type CallbackPayloadStatus struct {
	// Conditions is the list of error conditions for this resource
	//+operator-sdk:csv:customresourcedefinitions:type=status,displayName="Conditions",xDescriptors={"urn:alm:descriptor:io.kubernetes.conditions"}
	//+optional
	Conditions []CallbackPayloadCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CallbackPayload is storing the actual Payload Data we want to send back to any Callback URL. The
// web services receiving the Payload are determined via metav1.LabelSelector `selector`.
type CallbackPayload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CallbackPayloadSpec   `json:"spec,omitempty"`
	Status CallbackPayloadStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CallbackPayloadList contains a list of CallbackPayload
type CallbackPayloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CallbackPayload `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CallbackPayload{}, &CallbackPayloadList{})
}
