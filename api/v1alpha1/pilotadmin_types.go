package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PilotAdmin is the Schema for the pilotadmins API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=pilotadmins,scope=Namespaced
type PilotAdmin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PilotAdminSpec   `json:"spec,omitempty"`
	Status PilotAdminStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PilotAdminList contains a list of PilotAdmin
type PilotAdminList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PilotAdmin `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PilotAdmin{}, &PilotAdminList{})
}
