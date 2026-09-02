package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	platformv1alpha1 "github.com/Yash-Singh/KubeForge/api/v1alpha1"
)

const (
	managedByLabel = "app.kubernetes.io/managed-by"
	managedByValue = "kubeforge-operator"
	nameLabel      = "app.kubernetes.io/name"

	defaultCPURequest    = "10m"
	defaultMemoryRequest = "64Mi"
	defaultCPULimit      = "500m"
	defaultMemoryLimit   = "128Mi"
)

func applicationLabels(app *platformv1alpha1.Application) map[string]string {
	return map[string]string{
		nameLabel:      app.Name,
		managedByLabel: managedByValue,
	}
}

func desiredDeploymentSpec(app *platformv1alpha1.Application, replicas int32) appsv1.DeploymentSpec {
	labels := applicationLabels(app)

	spec := appsv1.DeploymentSpec{
		Replicas: &replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				nameLabel: app.Name,
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: labels,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  app.Name,
						Image: app.Spec.Image,
					},
				},
			},
		},
	}

	container := &spec.Template.Spec.Containers[0]

	if app.Spec.Resources != nil {
		container.Resources = *app.Spec.Resources
	} else {
		container.Resources = corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(defaultCPURequest),
				corev1.ResourceMemory: resource.MustParse(defaultMemoryRequest),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(defaultCPULimit),
				corev1.ResourceMemory: resource.MustParse(defaultMemoryLimit),
			},
		}
	}

	if len(app.Spec.Env) > 0 {
		container.Env = convertEnvVars(app.Spec.Env)
	}

	if app.Spec.Probes != nil {
		if app.Spec.Probes.Liveness != nil {
			container.LivenessProbe = buildProbe(app.Spec.Probes.Liveness)
		}
		if app.Spec.Probes.Readiness != nil {
			container.ReadinessProbe = buildProbe(app.Spec.Probes.Readiness)
		}
	}

	if len(app.Spec.TopologySpreadConstraints) > 0 {
		spec.Template.Spec.TopologySpreadConstraints = buildTopologySpreadConstraints(app)
	}

	return spec
}

func buildProbe(p *platformv1alpha1.ProbeSpec) *corev1.Probe {
	probe := &corev1.Probe{
		InitialDelaySeconds: getOrDefault(p.InitialDelaySeconds, 10),
		PeriodSeconds:       getOrDefault(p.PeriodSeconds, 10),
		TimeoutSeconds:      getOrDefault(p.TimeoutSeconds, 1),
		FailureThreshold:    getOrDefault(p.FailureThreshold, 3),
		SuccessThreshold:    getOrDefault(p.SuccessThreshold, 1),
	}

	if p.HTTPGet != nil {
		probe.ProbeHandler = corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: p.HTTPGet.Path,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: p.HTTPGet.Port,
				},
				Scheme: p.HTTPGet.Scheme,
			},
		}
		if len(p.HTTPGet.HTTPHeaders) > 0 {
			probe.HTTPGet.HTTPHeaders = make([]corev1.HTTPHeader, len(p.HTTPGet.HTTPHeaders))
			for i, h := range p.HTTPGet.HTTPHeaders {
				probe.HTTPGet.HTTPHeaders[i] = corev1.HTTPHeader{Name: h.Name, Value: h.Value}
			}
		}
	}

	return probe
}

func desiredPodDisruptionBudget(app *platformv1alpha1.Application) *policyv1.PodDisruptionBudget {
	if app.Spec.PodDisruptionBudget == nil {
		return nil
	}

	pdbSpec := app.Spec.PodDisruptionBudget

	pdb := &policyv1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			Labels:    applicationLabels(app),
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					nameLabel: app.Name,
				},
			},
		},
	}

	if pdbSpec.MinAvailable != nil {
		pdb.Spec.MinAvailable = pdbSpec.MinAvailable
	}
	if pdbSpec.MaxUnavailable != nil {
		pdb.Spec.MaxUnavailable = pdbSpec.MaxUnavailable
	}

	return pdb
}

func buildTopologySpreadConstraints(app *platformv1alpha1.Application) []corev1.TopologySpreadConstraint {
	constraints := make([]corev1.TopologySpreadConstraint, len(app.Spec.TopologySpreadConstraints))
	for i, tsc := range app.Spec.TopologySpreadConstraints {
		constraints[i] = corev1.TopologySpreadConstraint{
			MaxSkew:           tsc.MaxSkew,
			TopologyKey:       tsc.TopologyKey,
			WhenUnsatisfiable: getOrDefaultUnsatisfiable(tsc.WhenUnsatisfiable),
			LabelSelector:     tsc.LabelSelector,
		}
	}
	return constraints
}

func getOrDefaultUnsatisfiable(action corev1.UnsatisfiableConstraintAction) corev1.UnsatisfiableConstraintAction {
	if action == "" {
		return corev1.DoNotSchedule
	}
	return action
}

func desiredNetworkPolicy(app *platformv1alpha1.Application) *networkingv1.NetworkPolicy {
	if app.Spec.NetworkPolicy == nil {
		return nil
	}

	npSpec := app.Spec.NetworkPolicy

	np := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			Labels:    applicationLabels(app),
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					nameLabel: app.Name,
				},
			},
		},
	}

	// Determine policy types
	policyTypes := []networkingv1.PolicyType{}
	hasIngress := len(npSpec.Ingress) > 0
	hasEgress := len(npSpec.Egress) > 0

	if len(npSpec.PolicyTypes) > 0 {
		for _, pt := range npSpec.PolicyTypes {
			policyTypes = append(policyTypes, networkingv1.PolicyType(pt))
		}
	} else {
		if hasIngress {
			policyTypes = append(policyTypes, networkingv1.PolicyTypeIngress)
		}
		if hasEgress {
			policyTypes = append(policyTypes, networkingv1.PolicyTypeEgress)
		}
	}
	np.Spec.PolicyTypes = policyTypes

	// Build ingress rules
	if hasIngress {
		np.Spec.Ingress = buildNetworkPolicyIngress(npSpec.Ingress)
	}

	// Build egress rules
	if hasEgress {
		np.Spec.Egress = buildNetworkPolicyEgress(npSpec.Egress)
	}

	return np
}

func buildNetworkPolicyIngress(rules []platformv1alpha1.NetworkPolicyIngressRule) []networkingv1.NetworkPolicyIngressRule {
	result := make([]networkingv1.NetworkPolicyIngressRule, len(rules))
	for i, rule := range rules {
		result[i] = networkingv1.NetworkPolicyIngressRule{
			Ports: buildNetworkPolicyPorts(rule.Ports),
			From:  buildNetworkPolicyPeers(rule.From),
		}
	}
	return result
}

func buildNetworkPolicyEgress(rules []platformv1alpha1.NetworkPolicyEgressRule) []networkingv1.NetworkPolicyEgressRule {
	result := make([]networkingv1.NetworkPolicyEgressRule, len(rules))
	for i, rule := range rules {
		result[i] = networkingv1.NetworkPolicyEgressRule{
			Ports: buildNetworkPolicyPorts(rule.Ports),
			To:    buildNetworkPolicyPeers(rule.To),
		}
	}
	return result
}

func buildNetworkPolicyPorts(ports []platformv1alpha1.NetworkPolicyPort) []networkingv1.NetworkPolicyPort {
	if len(ports) == 0 {
		return nil
	}
	result := make([]networkingv1.NetworkPolicyPort, len(ports))
	for i, p := range ports {
		np := networkingv1.NetworkPolicyPort{}
		if p.Protocol != nil {
			np.Protocol = p.Protocol
		}
		if p.Port != nil {
			np.Port = p.Port
		}
		if p.EndPort != nil {
			np.EndPort = p.EndPort
		}
		result[i] = np
	}
	return result
}

func buildNetworkPolicyPeers(peers []platformv1alpha1.NetworkPolicyPeer) []networkingv1.NetworkPolicyPeer {
	if len(peers) == 0 {
		return nil
	}
	result := make([]networkingv1.NetworkPolicyPeer, len(peers))
	for i, peer := range peers {
		result[i] = networkingv1.NetworkPolicyPeer{}
		if peer.PodSelector != nil {
			result[i].PodSelector = peer.PodSelector
		}
		if peer.NamespaceSelector != nil {
			result[i].NamespaceSelector = peer.NamespaceSelector
		}
		if peer.IPBlock != nil {
			result[i].IPBlock = &networkingv1.IPBlock{
				CIDR:   peer.IPBlock.CIDR,
				Except: peer.IPBlock.Except,
			}
		}
	}
	return result
}

func getOrDefault(ptr *int32, def int32) int32 {
	if ptr != nil {
		return *ptr
	}
	return def
}

func convertEnvVars(env []platformv1alpha1.EnvVar) []corev1.EnvVar {
	if len(env) == 0 {
		return nil
	}
	result := make([]corev1.EnvVar, len(env))
	for i, e := range env {
		result[i] = corev1.EnvVar{Name: e.Name, Value: e.Value}
	}
	return result
}