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

// A QuerySpec defines the desired state of a Query.
type QuerySpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       QueryParameters `json:"forProvider"`
}

// A QueryStatus represents the observed state of a Query.
type QueryStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          QueryObservation `json:"atProvider,omitempty"`
	Synced              bool             `json:"synced,omitempty"`
}

// QueryParameters define the desired state of a MSSQL Query instance.
type QueryParameters struct {
	// +crossplane:generate:reference:type=Database
	Database *string `json:"database,omitempty"`
	// DatabaseRef allows you to specify custom resource name of the Database
	// to fill Database field.
	DatabaseRef *xpv1.Reference `json:"databaseRef,omitempty"`
	// DatabaseSelector allows you to use selector constraints to select a
	// Database.
	DatabaseSelector *xpv1.Selector `json:"databaseSelector,omitempty"`
	// Query is the Query that will be Queried
	// TODO (REL): check if its the syntax
	Query string `json:"query"`
}

// A QueryObservation represents the observed state of a MSSQL Query.
type QueryObservation struct {
	Results []map[string]string `json:"results,omitempty"`
	Error   string              `json:"error,omitempty"`
	Message string              `json:"message,omitempty"`
}

// +kubebuilder:object:root=true

// A Query represents the declarative state of a MSSQL Query.
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,sql}
type Query struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QuerySpec   `json:"spec"`
	Status QueryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// QueryList contains a list of Query
type QueryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Query `json:"items"`
}
