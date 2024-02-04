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

// A ExecuteSpec defines the desired state of a Execute.
type ExecuteSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       ExecuteParameters `json:"forProvider"`
}

// A ExecuteStatus represents the observed state of a Execute.
type ExecuteStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ExecuteObservation `json:"atProvider,omitempty"`
	Synced              bool             `json:"synced,omitempty"`
}

// ExecuteParameters define the desired state of a MSSQL Execute instance.
type ExecuteParameters struct {
	// +crossplane:generate:reference:type=Database
	Database *string `json:"database,omitempty"`
	// DatabaseRef allows you to specify custom resource name of the Database
	// to fill Database field.
	DatabaseRef *xpv1.Reference `json:"databaseRef,omitempty"`
	// DatabaseSelector allows you to use selector constraints to select a
	// Database.
	DatabaseSelector *xpv1.Selector `json:"databaseSelector,omitempty"`
	// Execute is the Execute that will be Queried
	// TODO (REL): check if its the syntax
	Execute string `json:"execute"`
}

// A ExecuteObservation represents the observed state of a MSSQL Execute.
type ExecuteObservation struct {
	Results []map[string]string `json:"results,omitempty"`
	Error   string              `json:"error,omitempty"`
	Message string              `json:"message,omitempty"`
}

// +kubebuilder:object:root=true

// A Execute represents the declarative state of a MSSQL Execute.
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,sql}
type Execute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExecuteSpec   `json:"spec"`
	Status ExecuteStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ExecuteList contains a list of Execute
type ExecuteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Execute `json:"items"`
}
