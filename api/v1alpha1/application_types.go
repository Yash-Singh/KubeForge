/*
Copyright 2026 KubeForge.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApplicationSpec defines the desired state of Application.
// See specs/04-crd-api-spec.md for field semantics.
type ApplicationSpec struct {
	// image is the container image reference for the primary application container.
	// +required
	Image string `json:"image"`

	// replicas is the desired number of pod replicas.
	// Defaults to 1.
	// +optional
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	Replicas *int32 `json:"replicas,omitempty"`

	// service configures a Service for the application.
	// +optional
	Service *ServiceSpec `json:"service,omitempty"`

	// resources specifies resource requests and limits for the primary container.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// env defines non-secret environment variables for the primary container.
	// Secret references are deferred — only name/value pairs are supported.
	// +optional
	Env []EnvVar `json:"env,omitempty"`
}

// ServiceSpec describes the Service to create for the Application.
type ServiceSpec struct {
	// port is the service port exposed by the Service.
	// +required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`

	// targetPort is the container port to forward traffic to.
	// Defaults to port when not set.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	TargetPort *int32 `json:"targetPort,omitempty"`
}

// EnvVar represents a non-secret environment variable.
// This is intentionally not corev1.EnvVar — secret/field references are deferred.
type EnvVar struct {
	// name of the environment variable.
	// +required
	Name string `json:"name"`

	// value of the environment variable.
	// +optional
	Value string `json:"value,omitempty"`
}

// ApplicationPhase describes the current operational phase.
type ApplicationPhase string

const (
	// PhaseReady means all owned resources are healthy and reconciled.
	PhaseReady ApplicationPhase = "Ready"
	// PhaseProgressing means resources are being created or updated.
	PhaseProgressing ApplicationPhase = "Progressing"
	// PhaseDegraded means the application failed to reach or maintain desired state.
	PhaseDegraded ApplicationPhase = "Degraded"
)

// Condition type constants for Application conditions.
const (
	ConditionReady       = "Ready"
	ConditionProgressing = "Progressing"
	ConditionDegraded    = "Degraded"
)

// ApplicationStatus defines the observed state of Application.
type ApplicationStatus struct {
	// observedGeneration is the generation of the Application most recently reconciled.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// phase summarizes the overall operational state.
	// One of: "", Ready, Progressing, Degraded.
	// +optional
	Phase ApplicationPhase `json:"phase,omitempty"`

	// readyReplicas is the number of ready pod replicas observed.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// desiredReplicas is the number of replicas the Application was last reconciled to.
	// +optional
	DesiredReplicas int32 `json:"desiredReplicas,omitempty"`

	// conditions represent the observed state of the Application.
	// Condition types: Ready, Progressing, Degraded.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=app

// Application is the Schema for the applications API.
type Application struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of Application
	// +required
	Spec ApplicationSpec `json:"spec"`

	// status defines the observed state of Application
	// +optional
	Status ApplicationStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
