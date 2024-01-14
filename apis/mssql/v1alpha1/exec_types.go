/*
Copyright 2021 The Crossplane Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// A ExecSpec defines the desired state of a Exec.
type ExecSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       ExecParameters `json:"forProvider"`
}

// A ExecStatus represents the observed state of a Exec.
type ExecStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ExecObservation `json:"atProvider,omitempty"`
	Synced              bool            `json:"synced,omitempty"`
}

// ExecParameters define the desired state of a MSSQL Exec instance.
type ExecParameters struct {
	// +crossplane:generate:reference:type=Database
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Field 'forProvider.database' is immutable"
	Database *string `json:"database,omitempty"`
	// DatabaseRef allows you to specify custom resource name of the Database
	// to fill Database field.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Field 'forProvider.databaseRef' is immutable"
	DatabaseRef *xpv1.Reference `json:"databaseRef,omitempty"`
	// DatabaseSelector allows you to use selector constraints to select a
	// Database.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Field 'forProvider.databaseSelector' is immutable"
	DatabaseSelector *xpv1.Selector `json:"databaseSelector,omitempty"`
	// Exec is the Exec that will be executed
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Field 'forProvider.query' is immutable"
	Exec string `json:"exec"`
}

// A ExecObservation represents the observed state of a MSSQL Exec.
type ExecObservation struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true

// A Exec represents the declarative state of a MSSQL Exec.
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,sql}
type Exec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExecSpec   `json:"spec"`
	Status ExecStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ExecList contains a list of Exec
type ExecList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Exec `json:"items"`
}
