# Phase 3 — Reliability and Scaling

## Task 14 — NetworkPolicy Generation

### What was done

Added NetworkPolicy generation support to the Application CRD and operator for network security and microsegmentation.

### Files changed

| File | Change |
|---|---|
| `api/v1alpha1/application_types.go` | Added `NetworkPolicySpec`, `NetworkPolicyIngressRule`, `NetworkPolicyEgressRule`, `NetworkPolicyPort`, `NetworkPolicyPeer`, `IPBlock` types |
| `internal/controller/application_builder.go` | Added `desiredNetworkPolicy()`, `buildNetworkPolicyIngress()`, `buildNetworkPolicyEgress()`, `buildNetworkPolicyPorts()`, `buildNetworkPolicyPeers()` |
| `internal/controller/application_controller.go` | Added `reconcileNetworkPolicy()`, NetworkPolicy ownership/watch, networkingv1 import |
| `config/crd/bases/platform.kubeforge.io_applications.yaml` | Regenerated with NetworkPolicy spec |
| `specs/04-crd-api-spec.md` | Documented NetworkPolicy |

### API changes

**New `spec.networkPolicy` field:**

```yaml
networkPolicy:
  ingress:
    - ports:
        - protocol: TCP
          port: 8080
      from:
        - podSelector:
            matchLabels:
              app: frontend
  egress:
    - ports:
        - protocol: TCP
          port: 53
      to:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: kube-system
```

### NetworkPolicy structure

| Field | Description |
|---|---|
| `ingress[]` | Ingress rules (allowed incoming traffic) |
| `egress[]` | Egress rules (allowed outgoing traffic) |
| `policyTypes[]` | Explicit types: `Ingress`, `Egress` (defaults inferred) |

**Ingress/Egress Rule:**
- `ports[]` - Array of port specifications
  - `protocol` - TCP, UDP, SCTP
  - `port` - Port number or name (IntOrString)
  - `endPort` - End of port range
- `from[]` / `to[]` - Peer specifications
  - `podSelector` - Select pods in same namespace
  - `namespaceSelector` - Select namespaces
  - `ipBlock` - CIDR range with optional `except`

### Verification

```bash
# Build and test
go build ./... && go test ./...

# Regenerate CRD
make manifests

# Apply CRD
kubectl apply -f config/crd/bases/platform.kubeforge.io_applications.yaml

# Check CRD
kubectl get crd applications.platform.kubeforge.io -o yaml | grep -A 20 "networkPolicy:"

# Test with example
kubectl apply -f - <<EOF
apiVersion: platform.kubeforge.io/v1alpha1
kind: Application
metadata:
  name: checkout
  namespace: apps
spec:
  image: ghcr.io/example/checkout:v1.0.0
  replicas: 3
  networkPolicy:
    ingress:
      - ports:
          - protocol: TCP
            port: 8080
        from:
          - podSelector:
              matchLabels:
                app: frontend
    egress:
      - ports:
          - protocol: TCP
            port: 53
        to:
          - namespaceSelector:
              matchLabels:
                kubernetes.io/metadata.name: kube-system
EOF

# Verify NetworkPolicy created
kubectl get networkpolicy -n apps checkout -o yaml
```

### Phase 3 Progress

| Item | Status |
|---|---|
| 1. Probes/defaults | ✅ Done |
| 2. Resource defaults | ✅ Done |
| 3. PodDisruptionBudget | ✅ Done |
| 4. Topology spread constraints | ✅ Done |
| 5. NetworkPolicy generation | ✅ Done |
| 6. HPA integration | ⏳ Next |
| 7. KEDA integration | ⏳ |
| 8. Argo Rollouts integration | ⏳ |