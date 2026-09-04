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
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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

	// probes configures liveness and readiness probes for the primary container.
	// +optional
	Probes *ProbesSpec `json:"probes,omitempty"`

	// podDisruptionBudget configures a PodDisruptionBudget for high availability.
	// +optional
	PodDisruptionBudget *PodDisruptionBudgetSpec `json:"podDisruptionBudget,omitempty"`

	// topologySpreadConstraints configures topology spread constraints for the pods.
	// +optional
	TopologySpreadConstraints []TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`

	// networkPolicy configures NetworkPolicy for the application pods.
	// +optional
	NetworkPolicy *NetworkPolicySpec `json:"networkPolicy,omitempty"`

	// horizontalPodAutoscaler configures HorizontalPodAutoscaler for the application.
	// Requires resource requests to be set on the container.
	// +optional
	HorizontalPodAutoscaler *HorizontalPodAutoscalerSpec `json:"horizontalPodAutoscaler,omitempty"`

	// kedaScaledObject configures a KEDA ScaledObject for event-driven autoscaling.
	// Requires KEDA to be installed in the cluster.
	// +optional
	KEDAScaledObject *KEDAScaledObjectSpec `json:"kedaScaledObject,omitempty"`

	// argoRollout configures an Argo Rollout instead of a standard Deployment.
	// Requires Argo Rollouts to be installed in the cluster.
	// +optional
	ArgoRollout *ArgoRolloutSpec `json:"argoRollout,omitempty"`
}

// PodDisruptionBudgetSpec defines the PodDisruptionBudget configuration.
type PodDisruptionBudgetSpec struct {
	// minAvailable specifies the minimum number of pods that must be available.
	// Can be an absolute number or percentage (e.g., "50%").
	// +optional
	MinAvailable *intstr.IntOrString `json:"minAvailable,omitempty"`

	// maxUnavailable specifies the maximum number of pods that can be unavailable.
	// Can be an absolute number or percentage (e.g., "25%").
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`
}

// TopologySpreadConstraint defines a topology spread constraint.
type TopologySpreadConstraint struct {
	// maxSkew describes the degree to which pods may be unevenly distributed.
	// +required
	// +kubebuilder:validation:Minimum=1
	MaxSkew int32 `json:"maxSkew"`

	// topologyKey is the key of node labels.
	// +required
	TopologyKey string `json:"topologyKey"`

	// whenUnsatisfiable indicates how to deal with a pod if it doesn't satisfy the spread constraint.
	// Options: DoNotSchedule (default), ScheduleAnyway.
	// +optional
	// +kubebuilder:default=DoNotSchedule
	WhenUnsatisfiable corev1.UnsatisfiableConstraintAction `json:"whenUnsatisfiable,omitempty"`

	// labelSelector is used to find matching pods.
	// If not specified, the Application's managed labels are used.
	// +optional
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty"`
}

// NetworkPolicySpec defines NetworkPolicy configuration for the application.
type NetworkPolicySpec struct {
	// ingress defines allowed ingress traffic to the application pods.
	// +optional
	Ingress []NetworkPolicyIngressRule `json:"ingress,omitempty"`

	// egress defines allowed egress traffic from the application pods.
	// +optional
	Egress []NetworkPolicyEgressRule `json:"egress,omitempty"`

	// policyTypes specifies the types of policies to apply.
	// Valid values: "Ingress", "Egress". Defaults to ["Ingress"] if only ingress defined, ["Egress"] if only egress defined, or both if both defined.
	// +optional
	PolicyTypes []string `json:"policyTypes,omitempty"`
}

// NetworkPolicyIngressRule defines an ingress rule for NetworkPolicy.
type NetworkPolicyIngressRule struct {
	// ports specifies the ports to allow.
	// +optional
	Ports []NetworkPolicyPort `json:"ports,omitempty"`

	// from specifies the sources allowed to access the pods.
	// +optional
	From []NetworkPolicyPeer `json:"from,omitempty"`
}

// NetworkPolicyEgressRule defines an egress rule for NetworkPolicy.
type NetworkPolicyEgressRule struct {
	// ports specifies the ports to allow.
	// +optional
	Ports []NetworkPolicyPort `json:"ports,omitempty"`

	// to specifies the destinations the pods can access.
	// +optional
	To []NetworkPolicyPeer `json:"to,omitempty"`
}

// NetworkPolicyPort defines a port for NetworkPolicy.
type NetworkPolicyPort struct {
	// protocol is the protocol (TCP, UDP, SCTP).
	// +optional
	Protocol *corev1.Protocol `json:"protocol,omitempty"`

	// port is the port number or name.
	// +optional
	Port *intstr.IntOrString `json:"port,omitempty"`

	// endPort is the end of a port range.
	// +optional
	EndPort *int32 `json:"endPort,omitempty"`
}

// NetworkPolicyPeer defines a peer for NetworkPolicy ingress/egress.
type NetworkPolicyPeer struct {
	// podSelector selects pods in the same namespace.
	// +optional
	PodSelector *metav1.LabelSelector `json:"podSelector,omitempty"`

	// namespaceSelector selects namespaces.
	// +optional
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`

	// ipBlock selects IP blocks (CIDR).
	// +optional
	IPBlock *IPBlock `json:"ipBlock,omitempty"`
}

// IPBlock defines an IP block for NetworkPolicy.
type IPBlock struct {
	// cidr is the CIDR range.
	// +required
	CIDR string `json:"cidr"`

	// except specifies CIDR ranges to exclude.
	// +optional
	Except []string `json:"except,omitempty"`
}

// ProbesSpec defines liveness and readiness probe configuration.
type ProbesSpec struct {
	// liveness configures the liveness probe.
	// +optional
	Liveness *ProbeSpec `json:"liveness,omitempty"`

	// readiness configures the readiness probe.
	// +optional
	Readiness *ProbeSpec `json:"readiness,omitempty"`
}

// ProbeSpec defines a container probe.
// Defaults match Kubernetes defaults when not specified.
type ProbeSpec struct {
	// httpGet specifies the HTTP GET action to perform.
	// +optional
	HTTPGet *HTTPGetAction `json:"httpGet,omitempty"`

	// initialDelaySeconds is the number of seconds after the container starts before the probe is initiated.
	// Defaults to 10.
	// +optional
	// +kubebuilder:default=10
	InitialDelaySeconds *int32 `json:"initialDelaySeconds,omitempty"`

	// periodSeconds is how often to perform the probe.
	// Defaults to 10.
	// +optional
	// +kubebuilder:default=10
	PeriodSeconds *int32 `json:"periodSeconds,omitempty"`

	// timeoutSeconds is the number of seconds after which the probe times out.
	// Defaults to 1.
	// +optional
	// +kubebuilder:default=1
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty"`

	// failureThreshold is the number of consecutive failures before the probe is considered failed.
	// Defaults to 3.
	// +optional
	// +kubebuilder:default=3
	FailureThreshold *int32 `json:"failureThreshold,omitempty"`

	// successThreshold is the number of consecutive successes before the probe is considered successful.
	// Defaults to 1.
	// +optional
	// +kubebuilder:default=1
	SuccessThreshold *int32 `json:"successThreshold,omitempty"`
}

// HTTPGetAction specifies an HTTP GET action for a probe.
type HTTPGetAction struct {
	// path is the HTTP path to GET.
	// +required
	Path string `json:"path"`

	// port is the port number to connect to.
	// +required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`

	// scheme is the HTTP scheme to use (HTTP or HTTPS).
	// Defaults to HTTP.
	// +optional
	Scheme corev1.URIScheme `json:"scheme,omitempty"`

	// httpHeaders are custom headers to set in the request.
	// +optional
	HTTPHeaders []HTTPHeader `json:"httpHeaders,omitempty"`
}

// HTTPHeader represents an HTTP header.
type HTTPHeader struct {
	// name of the header.
	// +required
	Name string `json:"name"`

	// value of the header.
	// +required
	Value string `json:"value"`
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

// HorizontalPodAutoscalerSpec defines the HPA configuration.
type HorizontalPodAutoscalerSpec struct {
	// minReplicas is the minimum number of pod replicas.
	// +required
	// +kubebuilder:validation:Minimum=1
	MinReplicas int32 `json:"minReplicas"`

	// maxReplicas is the maximum number of pod replicas.
	// +required
	// +kubebuilder:validation:Minimum=1
	MaxReplicas int32 `json:"maxReplicas"`

	// targetCPUUtilizationPercentage is the target average CPU utilization percentage.
	// Defaults to 80 if not specified and no custom metrics are provided.
	// +optional
	TargetCPUUtilizationPercentage *int32 `json:"targetCPUUtilizationPercentage,omitempty"`

	// targetMemoryUtilizationPercentage is the target average memory utilization percentage.
	// +optional
	TargetMemoryUtilizationPercentage *int32 `json:"targetMemoryUtilizationPercentage,omitempty"`

	// metrics contains custom metrics for scaling.
	// If provided, targetCPU/targetMemory are ignored.
	// +optional
	Metrics []autoscalingv2.MetricSpec `json:"metrics,omitempty"`

	// behavior configures the scaling behavior of the target.
	// +optional
	Behavior *autoscalingv2.HorizontalPodAutoscalerBehavior `json:"behavior,omitempty"`
}

// KEDAScaledObjectSpec defines the KEDA ScaledObject configuration.
type KEDAScaledObjectSpec struct {
	// minReplicaCount is the minimum number of replicas KEDA will scale to.
	// +required
	// +kubebuilder:validation:Minimum=0
	MinReplicaCount int32 `json:"minReplicaCount"`

	// maxReplicaCount is the maximum number of replicas KEDA will scale to.
	// +required
	// +kubebuilder:validation:Minimum=1
	MaxReplicaCount int32 `json:"maxReplicaCount"`

	// pollingInterval is the interval to check each trigger source on.
	// Defaults to 30 seconds.
	// +optional
	// +kubebuilder:default=30
	PollingInterval *int32 `json:"pollingInterval,omitempty"`

	// cooldownPeriod is the period to wait after the last trigger reported active before scaling the resource back to 0.
	// Defaults to 300 seconds.
	// +optional
	// +kubebuilder:default=300
	CooldownPeriod *int32 `json:"cooldownPeriod,omitempty"`

	// triggers defines the scaling triggers.
	// +required
	Triggers []KEDATrigger `json:"triggers"`
}

// KEDATrigger defines a KEDA scaling trigger.
type KEDATrigger struct {
	// type is the type of trigger (e.g., "prometheus", "kafka", "cpu", "memory").
	// +required
	Type string `json:"type"`

	// metadata contains the trigger metadata as key-value pairs.
	// +required
	Metadata map[string]string `json:"metadata"`

	// authenticationRef specifies a TriggerAuthentication reference.
	// +optional
	AuthenticationRef *KEDAAuthenticationRef `json:"authenticationRef,omitempty"`
}

// KEDAAuthenticationRef references a KEDA TriggerAuthentication.
type KEDAAuthenticationRef struct {
	// name is the name of the TriggerAuthentication resource.
	// +required
	Name string `json:"name"`

	// namespace is the namespace of the TriggerAuthentication. If omitted, uses the ScaledObject's namespace.
	// +optional
	Namespace *string `json:"namespace,omitempty"`
}

// ArgoRolloutSpec defines an Argo Rollout instead of a standard Deployment.
type ArgoRolloutSpec struct {
	// replicas is the desired number of pod replicas.
	// +optional
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	Replicas *int32 `json:"replicas,omitempty"`

	// strategy defines the rollout strategy.
	// +optional
	Strategy *ArgoRolloutStrategy `json:"strategy,omitempty"`
}

// ArgoRolloutStrategy defines the Argo Rollout strategy.
type ArgoRolloutStrategy struct {
	// canary defines a canary rollout strategy.
	// +optional
	Canary *ArgoCanaryStrategy `json:"canary,omitempty"`

	// blueGreen defines a blue-green rollout strategy.
	// +optional
	BlueGreen *ArgoBlueGreenStrategy `json:"blueGreen,omitempty"`
}

// ArgoCanaryStrategy defines a canary rollout strategy.
type ArgoCanaryStrategy struct {
	// steps defines the canary rollout steps.
	// +optional
	Steps []ArgoCanaryStep `json:"steps,omitempty"`

	// canaryService is the name of the Service used to route traffic to canary pods.
	// +optional
	CanaryService string `json:"canaryService,omitempty"`

	// stableService is the name of the Service used to route traffic to stable pods.
	// +optional
	StableService string `json:"stableService,omitempty"`

	// trafficRouting defines the traffic routing configuration.
	// +optional
	TrafficRouting *ArgoTrafficRouting `json:"trafficRouting,omitempty"`
}

// ArgoCanaryStep defines a single canary rollout step.
type ArgoCanaryStep struct {
	// setWeight sets the traffic weight percentage.
	// +optional
	SetWeight *int32 `json:"setWeight,omitempty"`

	// pause defines a pause step.
	// +optional
	Pause *ArgoPauseStep `json:"pause,omitempty"`

	// setHeaderRoute sets a header-based routing.
	// +optional
	SetHeaderRoute *ArgoHeaderRoute `json:"setHeaderRoute,omitempty"`

	// setMirrorStep enables traffic mirroring.
	// +optional
	SetMirror *ArgoMirrorStep `json:"setMirror,omitempty"`
}

// ArgoPauseStep defines a pause step.
type ArgoPauseStep struct {
	// duration is the duration to pause (e.g., "5m", "1h").
	// If empty, waits indefinitely until manually resumed.
	// +optional
	Duration *string `json:"duration,omitempty"`
}

// ArgoHeaderRoute defines header-based routing.
type ArgoHeaderRoute struct {
	// name is the name of the routing rule.
	// +required
	Name string `json:"name"`

	// match defines the header match criteria.
	// +required
	Match []ArgoHeaderMatch `json:"match"`
}

// ArgoHeaderMatch defines a header match rule.
type ArgoHeaderMatch struct {
	// headerName is the header name to match.
	// +required
	HeaderName string `json:"headerName"`

	// headerValue is the header value to match.
	// +required
	HeaderValue ArgoHeaderValue `json:"headerValue"`
}

// ArgoHeaderValue defines the header value match.
type ArgoHeaderValue struct {
	// exact is the exact header value.
	// +optional
	Exact string `json:"exact,omitempty"`

	// regex is the regex pattern to match.
	// +optional
	Regex string `json:"regex,omitempty"`

	// prefix is the prefix to match.
	// +optional
	Prefix string `json:"prefix,omitempty"`
}

// ArgoMirrorStep defines traffic mirroring.
type ArgoMirrorStep struct {
	// percentage is the percentage of traffic to mirror.
	// +required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	Percentage int32 `json:"percentage"`
}

// ArgoBlueGreenStrategy defines a blue-green rollout strategy.
type ArgoBlueGreenStrategy struct {
	// activeService is the Service used to route traffic to the active (blue) pods.
	// +required
	ActiveService string `json:"activeService"`

	// previewService is the Service used to route traffic to the preview (green) pods.
	// +optional
	PreviewService string `json:"previewService,omitempty"`

	// autoPromotionEnabled automatically promotes the new version when ready.
	// Defaults to true.
	// +optional
	// +kubebuilder:default=true
	AutoPromotionEnabled *bool `json:"autoPromotionEnabled,omitempty"`

	// prePromotionAnalysis defines analysis to run before promotion.
	// +optional
	PrePromotionAnalysis *ArgoAnalysis `json:"prePromotionAnalysis,omitempty"`

	// postPromotionAnalysis defines analysis to run after promotion.
	// +optional
	PostPromotionAnalysis *ArgoAnalysis `json:"postPromotionAnalysis,omitempty"`
}

// ArgoAnalysis defines an analysis template for Argo Rollouts.
type ArgoAnalysis struct {
	// templates references AnalysisTemplate resources.
	// +required
	Templates []ArgoAnalysisTemplateRef `json:"templates"`

	// args are the arguments to pass to the AnalysisTemplate.
	// +optional
	Args []ArgoAnalysisArgument `json:"args,omitempty"`
}

// ArgoAnalysisTemplateRef references an AnalysisTemplate by name.
type ArgoAnalysisTemplateRef struct {
	// templateName is the name of the AnalysisTemplate.
	// +required
	TemplateName string `json:"templateName"`
}

// ArgoAnalysisArgument defines an argument for AnalysisTemplate.
type ArgoAnalysisArgument struct {
	// name is the argument name.
	// +required
	Name string `json:"name"`

	// value is the argument value.
	// +optional
	Value *string `json:"value,omitempty"`

	// valueFrom references a Pod value.
	// +optional
	ValueFrom *ArgoArgumentValueFrom `json:"valueFrom,omitempty"`
}

// ArgoArgumentValueFrom references a value from the pod.
type ArgoArgumentValueFrom struct {
	// fieldRef references a field in the pod spec.
	// +required
	FieldRef Argov1FieldRef `json:"fieldRef"`
}

// Argov1FieldRef defines a field reference.
type Argov1FieldRef struct {
	// fieldPath is the path to the field (e.g., "metadata.labels['version']").
	// +required
	FieldPath string `json:"fieldPath"`
}

// ArgoTrafficRouting defines traffic routing configuration.
type ArgoTrafficRouting struct {
	// istio configures Istio traffic routing.
	// +optional
	Istio *ArgoIstioTrafficRouting `json:"istio,omitempty"`

	// nginx configures NGINX traffic routing.
	// +optional
	Nginx *ArgoNginxTrafficRouting `json:"nginx,omitempty"`

	// smi configures SMI traffic routing.
	// +optional
	SMI *ArgoSMITrafficRouting `json:"smi,omitempty"`

	// alb configures AWS ALB traffic routing.
	// +optional
	ALB *ArgoALBTrafficRouting `json:"alb,omitempty"`
}

// ArgoIstioTrafficRouting defines Istio traffic routing.
type ArgoIstioTrafficRouting struct {
	// virtualServices references Istio VirtualServices.
	// +optional
	VirtualServices []ArgoIstioVirtualService `json:"virtualServices,omitempty"`
}

// ArgoIstioVirtualService references an Istio VirtualService.
type ArgoIstioVirtualService struct {
	// name is the VirtualService name.
	// +required
	Name string `json:"name"`

	// routes are the routes within the VirtualService.
	// +required
	Routes []string `json:"routes"`
}

// ArgoNginxTrafficRouting defines NGINX traffic routing.
type ArgoNginxTrafficRouting struct {
	// stableIngress is the name of the stable Ingress.
	// +required
	StableIngress string `json:"stableIngress"`

	// additionalIngressAnnotations are additional annotations for the canary Ingress.
	// +optional
	AdditionalIngressAnnotations map[string]string `json:"additionalIngressAnnotations,omitempty"`
}

// ArgoSMITrafficRouting defines SMI traffic routing.
type ArgoSMITrafficRouting struct {
	// trafficSplitName is the name of the TrafficSplit resource.
	// +required
	TrafficSplitName string `json:"trafficSplitName"`
}

// ArgoALBTrafficRouting defines AWS ALB traffic routing.
type ArgoALBTrafficRouting struct {
	// rootService is the root service for the ALB.
	// +required
	RootService string `json:"rootService"`

	// ingress is the Ingress annotation settings.
	// +optional
	Ingress map[string]string `json:"ingress,omitempty"`
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
