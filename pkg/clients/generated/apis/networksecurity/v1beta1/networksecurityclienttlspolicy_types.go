// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Config Connector and manual
//     changes will be clobbered when the file is regenerated.
//
// ----------------------------------------------------------------------------

// *** DISCLAIMER ***
// Config Connector's go-client for CRDs is currently in ALPHA, which means
// that future versions of the go-client may include breaking changes.
// Please try it out and give us feedback!

package v1beta1

import (
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/clients/generated/apis/k8s/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClienttlspolicyCertificateProviderInstance struct {
	/* Required. Plugin instance name, used to locate and load CertificateProvider instance configuration. Set to "google_cloud_private_spiffe" to use Certificate Authority Service certificate provider instance. */
	PluginInstance string `json:"pluginInstance"`
}

type ClienttlspolicyClientCertificate struct {
	/* The certificate provider instance specification that will be passed to the data plane, which will be used to load necessary credential information. */
	// +optional
	CertificateProviderInstance *ClienttlspolicyCertificateProviderInstance `json:"certificateProviderInstance,omitempty"`

	/* gRPC specific configuration to access the gRPC server to obtain the cert and private key. */
	// +optional
	GrpcEndpoint *ClienttlspolicyGrpcEndpoint `json:"grpcEndpoint,omitempty"`
}

type ClienttlspolicyGrpcEndpoint struct {
	/* Required. The target URI of the gRPC endpoint. Only UDS path is supported, and should start with “unix:”. */
	TargetUri string `json:"targetUri"`
}

type ClienttlspolicyServerValidationCa struct {
	/* The certificate provider instance specification that will be passed to the data plane, which will be used to load necessary credential information. */
	// +optional
	CertificateProviderInstance *ClienttlspolicyCertificateProviderInstance `json:"certificateProviderInstance,omitempty"`

	/* gRPC specific configuration to access the gRPC server to obtain the CA certificate. */
	// +optional
	GrpcEndpoint *ClienttlspolicyGrpcEndpoint `json:"grpcEndpoint,omitempty"`
}

type NetworkSecurityClientTLSPolicySpec struct {
	/* Optional. Defines a mechanism to provision client identity (public and private keys) for peer to peer authentication. The presence of this dictates mTLS. */
	// +optional
	ClientCertificate *ClienttlspolicyClientCertificate `json:"clientCertificate,omitempty"`

	/* Optional. Free-text description of the resource. */
	// +optional
	Description *string `json:"description,omitempty"`

	/* Immutable. The location for the resource */
	Location string `json:"location"`

	/* Immutable. The Project that this resource belongs to. */
	// +optional
	ProjectRef *v1alpha1.ResourceRef `json:"projectRef,omitempty"`

	/* Immutable. Optional. The name of the resource. Used for creation and acquisition. When unset, the value of `metadata.name` is used as the default. */
	// +optional
	ResourceID *string `json:"resourceID,omitempty"`

	/* Required. Defines the mechanism to obtain the Certificate Authority certificate to validate the server certificate. */
	// +optional
	ServerValidationCa []ClienttlspolicyServerValidationCa `json:"serverValidationCa,omitempty"`

	/* Optional. Server Name Indication string to present to the server during TLS handshake. E.g: "secure.example.com". */
	// +optional
	Sni *string `json:"sni,omitempty"`
}

type NetworkSecurityClientTLSPolicyStatus struct {
	/* Conditions represent the latest available observations of the
	   NetworkSecurityClientTLSPolicy's current state. */
	Conditions []v1alpha1.Condition `json:"conditions,omitempty"`
	/* Output only. The timestamp when the resource was created. */
	CreateTime string `json:"createTime,omitempty"`
	/* ObservedGeneration is the generation of the resource that was most recently observed by the Config Connector controller. If this is equal to metadata.generation, then that means that the current reported status reflects the most recent desired state of the resource. */
	ObservedGeneration int `json:"observedGeneration,omitempty"`
	/* Output only. The timestamp when the resource was updated. */
	UpdateTime string `json:"updateTime,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkSecurityClientTLSPolicy is the Schema for the networksecurity API
// +k8s:openapi-gen=true
type NetworkSecurityClientTLSPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkSecurityClientTLSPolicySpec   `json:"spec,omitempty"`
	Status NetworkSecurityClientTLSPolicyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkSecurityClientTLSPolicyList contains a list of NetworkSecurityClientTLSPolicy
type NetworkSecurityClientTLSPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkSecurityClientTLSPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NetworkSecurityClientTLSPolicy{}, &NetworkSecurityClientTLSPolicyList{})
}
