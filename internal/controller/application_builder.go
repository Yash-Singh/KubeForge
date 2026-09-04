package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

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

func desiredHorizontalPodAutoscaler(app *platformv1alpha1.Application) *autoscalingv2.HorizontalPodAutoscaler {
	if app.Spec.HorizontalPodAutoscaler == nil {
		return nil
	}

	hpaSpec := app.Spec.HorizontalPodAutoscaler

	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			Labels:    applicationLabels(app),
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       app.Name,
			},
			MinReplicas: &hpaSpec.MinReplicas,
			MaxReplicas: hpaSpec.MaxReplicas,
		},
	}

	if len(hpaSpec.Metrics) > 0 {
		hpa.Spec.Metrics = hpaSpec.Metrics
	} else {
		metrics := []autoscalingv2.MetricSpec{}
		if hpaSpec.TargetCPUUtilizationPercentage != nil {
			target := autoscalingv2.MetricTarget{
				Type:               autoscalingv2.UtilizationMetricType,
				AverageUtilization: hpaSpec.TargetCPUUtilizationPercentage,
			}
			metrics = append(metrics, autoscalingv2.MetricSpec{
				Type:      autoscalingv2.ResourceMetricSourceType,
				Resource:  &autoscalingv2.ResourceMetricSource{Name: corev1.ResourceCPU, Target: target},
			})
		}
		if hpaSpec.TargetMemoryUtilizationPercentage != nil {
			target := autoscalingv2.MetricTarget{
				Type:               autoscalingv2.UtilizationMetricType,
				AverageUtilization: hpaSpec.TargetMemoryUtilizationPercentage,
			}
			metrics = append(metrics, autoscalingv2.MetricSpec{
				Type:      autoscalingv2.ResourceMetricSourceType,
				Resource:  &autoscalingv2.ResourceMetricSource{Name: corev1.ResourceMemory, Target: target},
			})
		}
		if len(metrics) == 0 {
			target := int32(80)
			metrics = append(metrics, autoscalingv2.MetricSpec{
				Type: autoscalingv2.ResourceMetricSourceType,
				Resource: &autoscalingv2.ResourceMetricSource{
					Name: corev1.ResourceCPU,
					Target: autoscalingv2.MetricTarget{
						Type:               autoscalingv2.UtilizationMetricType,
						AverageUtilization: &target,
					},
				},
			})
		}
		hpa.Spec.Metrics = metrics
	}

	if hpaSpec.Behavior != nil {
		hpa.Spec.Behavior = hpaSpec.Behavior
	}

	return hpa
}

func desiredKEDAScaledObject(app *platformv1alpha1.Application) *unstructured.Unstructured {
	if app.Spec.KEDAScaledObject == nil {
		return nil
	}

	kedaSpec := app.Spec.KEDAScaledObject

	triggers := make([]interface{}, 0, len(kedaSpec.Triggers))
	for _, t := range kedaSpec.Triggers {
		trigger := map[string]interface{}{
			"type": t.Type,
		}
		if len(t.Metadata) > 0 {
			metadata := map[string]interface{}{}
			for k, v := range t.Metadata {
				metadata[k] = v
			}
			trigger["metadata"] = metadata
		}
		if t.AuthenticationRef != nil {
			authRef := map[string]interface{}{
				"name": t.AuthenticationRef.Name,
			}
			if t.AuthenticationRef.Namespace != nil {
				authRef["namespace"] = *t.AuthenticationRef.Namespace
			}
			trigger["authenticationRef"] = authRef
		}
		triggers = append(triggers, trigger)
	}

	scaleTargetRef := map[string]interface{}{
		"apiversion": "apps/v1",
		"kind":       "Deployment",
		"name":       app.Name,
	}

	spec := map[string]interface{}{
		"scaleTargetRef": scaleTargetRef,
		"minReplicaCount": int64(kedaSpec.MinReplicaCount),
		"maxReplicaCount": int64(kedaSpec.MaxReplicaCount),
		"triggers":        triggers,
	}
	if kedaSpec.PollingInterval != nil {
		spec["pollingInterval"] = int64(*kedaSpec.PollingInterval)
	}
	if kedaSpec.CooldownPeriod != nil {
		spec["cooldownPeriod"] = int64(*kedaSpec.CooldownPeriod)
	}

	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "keda.sh/v1alpha1",
			"kind":       "ScaledObject",
			"metadata": map[string]interface{}{
				"name":      app.Name,
				"namespace": app.Namespace,
				"labels":    applicationLabels(app),
			},
			"spec": spec,
		},
	}
}

func desiredArgoRollout(app *platformv1alpha1.Application) *unstructured.Unstructured {
	if app.Spec.ArgoRollout == nil {
		return nil
	}

	rolloutSpec := app.Spec.ArgoRollout

	spec := map[string]interface{}{
		"replicas": int64(ptr.Deref(rolloutSpec.Replicas, 1)),
		"selector": map[string]interface{}{
			"matchLabels": map[string]interface{}{
				nameLabel: app.Name,
			},
		},
		"template": map[string]interface{}{
			"metadata": map[string]interface{}{
				"labels": applicationLabels(app),
			},
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"name":  app.Name,
						"image": app.Spec.Image,
					},
				},
			},
		},
	}

	if rolloutSpec.Strategy != nil {
		strategy := map[string]interface{}{}
		if rolloutSpec.Strategy.Canary != nil {
			canary := buildArgoCanaryStrategy(app, rolloutSpec.Strategy.Canary)
			strategy["canary"] = canary
		}
		if rolloutSpec.Strategy.BlueGreen != nil {
			bg := buildArgoBlueGreenStrategy(app, rolloutSpec.Strategy.BlueGreen)
			strategy["blueGreen"] = bg
		}
		spec["strategy"] = strategy
	}

	// Apply probes and resources from the main container if set
	containers := spec["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})
	container := containers[0].(map[string]interface{})

	if app.Spec.Resources != nil {
		res := map[string]interface{}{}
		if app.Spec.Resources.Requests != nil {
			reqs := map[string]interface{}{}
			if v, ok := app.Spec.Resources.Requests[corev1.ResourceCPU]; ok {
				reqs["cpu"] = v.String()
			}
			if v, ok := app.Spec.Resources.Requests[corev1.ResourceMemory]; ok {
				reqs["memory"] = v.String()
			}
			res["requests"] = reqs
		}
		if app.Spec.Resources.Limits != nil {
			limits := map[string]interface{}{}
			if v, ok := app.Spec.Resources.Limits[corev1.ResourceCPU]; ok {
				limits["cpu"] = v.String()
			}
			if v, ok := app.Spec.Resources.Limits[corev1.ResourceMemory]; ok {
				limits["memory"] = v.String()
			}
			res["limits"] = limits
		}
		container["resources"] = res
	}

	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "argoproj.io/v1alpha1",
			"kind":       "Rollout",
			"metadata": map[string]interface{}{
				"name":      app.Name,
				"namespace": app.Namespace,
				"labels":    applicationLabels(app),
				"ownerReferences": []interface{}{
					map[string]interface{}{
						"apiVersion":         "platform.kubeforge.io/v1alpha1",
						"kind":               "Application",
						"name":               app.Name,
						"uid":                string(app.UID),
						"controller":         true,
						"blockOwnerDeletion": true,
					},
				},
			},
			"spec": spec,
		},
	}
}

func buildArgoCanaryStrategy(app *platformv1alpha1.Application, canarySpec *platformv1alpha1.ArgoCanaryStrategy) map[string]interface{} {
	result := map[string]interface{}{}

	if len(canarySpec.Steps) > 0 {
		steps := make([]interface{}, 0, len(canarySpec.Steps))
		for _, step := range canarySpec.Steps {
			stepMap := map[string]interface{}{}
			if step.SetWeight != nil {
				stepMap["setWeight"] = int64(*step.SetWeight)
			}
			if step.Pause != nil {
				pause := map[string]interface{}{}
				if step.Pause.Duration != nil {
					pause["duration"] = *step.Pause.Duration
				}
				stepMap["pause"] = pause
			}
			if step.SetHeaderRoute != nil {
				stepMap["setHeaderRoute"] = buildArgoHeaderRoute(step.SetHeaderRoute)
			}
			if step.SetMirror != nil {
				stepMap["setMirror"] = map[string]interface{}{
					"percentage": int64(step.SetMirror.Percentage),
				}
			}
			steps = append(steps, stepMap)
		}
		result["steps"] = steps
	}

	if canarySpec.CanaryService != "" {
		result["canaryService"] = canarySpec.CanaryService
	}
	if canarySpec.StableService != "" {
		result["stableService"] = canarySpec.StableService
	}

	if canarySpec.TrafficRouting != nil {
		tr := map[string]interface{}{}
		if canarySpec.TrafficRouting.Istio != nil {
			istio := map[string]interface{}{}
			if len(canarySpec.TrafficRouting.Istio.VirtualServices) > 0 {
				vs := make([]interface{}, 0, len(canarySpec.TrafficRouting.Istio.VirtualServices))
				for _, v := range canarySpec.TrafficRouting.Istio.VirtualServices {
					vsEntry := map[string]interface{}{
						"name": v.Name,
						"routes": v.Routes,
					}
					vs = append(vs, vsEntry)
				}
				istio["virtualServices"] = vs
			}
			tr["istio"] = istio
		}
		if canarySpec.TrafficRouting.Nginx != nil {
			nginx := map[string]interface{}{
				"stableIngress": canarySpec.TrafficRouting.Nginx.StableIngress,
			}
			if len(canarySpec.TrafficRouting.Nginx.AdditionalIngressAnnotations) > 0 {
				nginx["additionalIngressAnnotations"] = canarySpec.TrafficRouting.Nginx.AdditionalIngressAnnotations
			}
			tr["nginx"] = nginx
		}
		if canarySpec.TrafficRouting.SMI != nil {
			tr["smi"] = map[string]interface{}{
				"trafficSplitName": canarySpec.TrafficRouting.SMI.TrafficSplitName,
			}
		}
		if canarySpec.TrafficRouting.ALB != nil {
			alb := map[string]interface{}{
				"rootService": canarySpec.TrafficRouting.ALB.RootService,
			}
			if len(canarySpec.TrafficRouting.ALB.Ingress) > 0 {
				alb["ingress"] = canarySpec.TrafficRouting.ALB.Ingress
			}
			tr["alb"] = alb
		}
		result["trafficRouting"] = tr
	}

	return result
}

func buildArgoHeaderRoute(route *platformv1alpha1.ArgoHeaderRoute) map[string]interface{} {
	result := map[string]interface{}{
		"name": route.Name,
	}
	if len(route.Match) > 0 {
		match := make([]interface{}, 0, len(route.Match))
		for _, m := range route.Match {
			entry := map[string]interface{}{
				"headerName": m.HeaderName,
				"headerValue": map[string]interface{}{
					"exact": m.HeaderValue.Exact,
					"regex": m.HeaderValue.Regex,
					"prefix": m.HeaderValue.Prefix,
				},
			}
			match = append(match, entry)
		}
		result["match"] = match
	}
	return result
}

func buildArgoBlueGreenStrategy(app *platformv1alpha1.Application, bgSpec *platformv1alpha1.ArgoBlueGreenStrategy) map[string]interface{} {
	result := map[string]interface{}{
		"activeService": bgSpec.ActiveService,
	}
	if bgSpec.PreviewService != "" {
		result["previewService"] = bgSpec.PreviewService
	}
	if bgSpec.AutoPromotionEnabled != nil {
		result["autoPromotionEnabled"] = *bgSpec.AutoPromotionEnabled
	}
	if bgSpec.PrePromotionAnalysis != nil {
		result["prePromotionAnalysis"] = buildArgoAnalysis(app, bgSpec.PrePromotionAnalysis)
	}
	if bgSpec.PostPromotionAnalysis != nil {
		result["postPromotionAnalysis"] = buildArgoAnalysis(app, bgSpec.PostPromotionAnalysis)
	}
	return result
}

func buildArgoAnalysis(app *platformv1alpha1.Application, analysis *platformv1alpha1.ArgoAnalysis) map[string]interface{} {
	result := map[string]interface{}{}

	if len(analysis.Templates) > 0 {
		templates := make([]interface{}, 0, len(analysis.Templates))
		for _, t := range analysis.Templates {
			templates = append(templates, map[string]interface{}{
				"templateName": t.TemplateName,
			})
		}
		result["templates"] = templates
	}

	if len(analysis.Args) > 0 {
		args := make([]interface{}, 0, len(analysis.Args))
		for _, a := range analysis.Args {
			arg := map[string]interface{}{
				"name": a.Name,
			}
			if a.Value != nil {
				arg["value"] = *a.Value
			}
			if a.ValueFrom != nil {
				arg["valueFrom"] = map[string]interface{}{
					"fieldRef": map[string]interface{}{
						"fieldPath": a.ValueFrom.FieldRef.FieldPath,
					},
				}
			}
			args = append(args, arg)
		}
		result["args"] = args
	}

	return result
}

var (
	rolloutGVR = schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "rollouts",
	}
	scaledObjectGVR = schema.GroupVersionResource{
		Group:    "keda.sh",
		Version:  "v1alpha1",
		Resource: "scaledobjects",
	}
)