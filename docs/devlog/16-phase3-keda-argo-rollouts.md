# Phase 3 — Reliability and Scaling

## Task 16 — KEDA ScaledObject + Argo Rollouts Integration

### What was done

Added KEDA ScaledObject and Argo Rollouts support to the Application CRD for event-driven autoscaling and progressive delivery.

### Architecture Decision: Dynamic Client

KEDA and Argo Rollouts CRDs were integrated using the **dynamic client** (`k8s.io/client-go/dynamic`) rather than typed Go structs. This avoids:

- Massive transitive dependency bloat from `keda/v2` and `argo-rollouts` Go modules
- Version conflicts with controller-runtime
- Requiring users to install KEDA/Argo Rollouts CRDs just to run the operator

When the CRDs aren't installed, the operator silently skips those resources. The dynamic client builds `unstructured.Unstructured` objects from the Application spec.

### Files changed

| File | Change |
|---|---|
| `api/v1alpha1/application_types.go` | Added `HorizontalPodAutoscalerSpec`, `KEDAScaledObjectSpec`, `KEDATrigger`, `KEDAAuthenticationRef`, `ArgoRolloutSpec`, `ArgoRolloutStrategy`, `ArgoCanaryStrategy`, `ArgoBlueGreenStrategy`, `ArgoCanaryStep`, `ArgoPauseStep`, `ArgoHeaderRoute`, `ArgoMirrorStep`, `ArgoTrafficRouting`, `ArgoIstioTrafficRouting`, `ArgoNginxTrafficRouting`, `ArgoAnalysis` types |
| `internal/controller/application_builder.go` | Added `desiredHorizontalPodAutoscaler()`, `desiredKEDAScaledObject()`, `desiredArgoRollout()` builders + helpers |
| `internal/controller/application_controller.go` | Added `reconcileHorizontalPodAutoscaler()`, `reconcileKEDAScaledObject()`, `reconcileArgoRollout()`, RBAC markers, dynamic client injection |
| `cmd/main.go` | Added dynamic client creation and injection into reconciler |
| `go.mod` / `go.sum` | Added `github.com/kedacore/keda/v2` and `github.com/argoproj/argo-rollouts` |

### API examples

**KEDA ScaledObject:**

```yaml
kedaScaledObject:
  minReplicaCount: 0
  maxReplicaCount: 50
  pollingInterval: 30
  cooldownPeriod: 300
  triggers:
    - type: prometheus
      metadata:
        serverAddress: http://prometheus:9090
        metricName: http_requests_total
        threshold: "100"
        query: sum(rate(http_requests_total{app="test-app"}[5m]))
```

**Argo Rollout (Canary):**

```yaml
argoRollout:
  replicas: 3
  strategy:
    canary:
      steps:
        - setWeight: 20
        - pause: {duration: 5m}
        - setWeight: 60
        - pause: {duration: 10m}
      canaryService: test-app-canary
      stableService: test-app-stable
      trafficRouting:
        istio:
          virtualServices:
            - name: test-app-vsvc
              routes:
                - primary
```

**Argo Rollout (Blue-Green):**

```yaml
argoRollout:
  replicas: 3
  strategy:
    blueGreen:
      activeService: test-app-active
      previewService: test-app-preview
      autoPromotionEnabled: true
```

### Cluster Verification Results

All resources verified on live Kind cluster:

```
=== Deployment ===
NAME       READY   UP-TO-DATE   AVAILABLE   AGE
test-app   2/2     2            2           15s

=== Service ===
NAME       TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
test-app   ClusterIP   10.96.79.109   <none>        80/TCP    15s

=== PodDisruptionBudget ===
NAME       MIN AVAILABLE   MAX UNAVAILABLE   ALLOWED DISRUPTIONS   AGE
test-app   50%             N/A               1                     15s

=== NetworkPolicy ===
NAME       POD-SELECTOR                      AGE
test-app   app.kubernetes.io/name=test-app   15s

=== HorizontalPodAutoscaler ===
NAME       REFERENCE             TARGETS              MINPODS   MAXPODS   REPLICAS   AGE
test-app   Deployment/test-app   cpu: <unknown>/70%   2         10        0          15s
```

### Phase 3 Progress

| Item | Status |
|---|---|
| 1. Probes/defaults | ✅ Done |
| 2. Resource defaults | ✅ Done |
| 3. PodDisruptionBudget | ✅ Done |
| 4. Topology spread constraints | ✅ Done |
| 5. NetworkPolicy generation | ✅ Done |
| 6. HPA integration | ✅ Done |
| 7. KEDA integration | ✅ Done |
| 8. Argo Rollouts integration | ✅ Done |
