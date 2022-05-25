/**
 * Copyright (C) 2022 Christoph Görn
 *
 * This file is part of rebuldah.
 *
 * rebuldah is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * rebuldah is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with rebuldah.  If not, see <http://www.gnu.org/licenses/>.
 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cfg "sigs.k8s.io/controller-runtime/pkg/config/v1alpha1"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ErinnerungConfig is the Schema for the erinnerungenconfigs API
type ErinnerungConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// ControllerManagerConfigurationSpec returns the contfigurations for controllers
	cfg.ControllerManagerConfigurationSpec `json:",inline"`

	// Namespaces is the list of Namespaces we want to operate in
	// TODO this might be an anti-pattern, if the operator is namespace-scoped, do we need to deploy it to each namespace we want to operate in?!
	Namespaces []string `json:"namespaces,omitempty"`
}

//+kubebuilder:object:root=true

func init() {
	SchemeBuilder.Register(&ErinnerungConfig{})
}
