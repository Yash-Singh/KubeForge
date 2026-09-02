# Phase 3 — Reliability and Scaling

## Task 13 — PodDisruptionBudget and Topology Spread Constraints

### What was done

Added PDB and topology spread constraint support to the Application CRD and operator for high availability.

### Files changed

| File | Change |
|---|---|
| `api/v1alpha1/application_types.go` | Added `PodDisruptionBudgetSpec`, `TopologySpreadConstraint` types |
| `internal/controller/application_builder.go` | Added `desiredPodDisruptionBudget`, `buildTopologySpreadConstraints` |
| `internal/controller/application_controller.go` | Added `reconcilePodDisruptionBudget`, PDB ownership/watch |
| `config/crd/bases/platform.kubeforge.io_applications.yaml` | Regenerated with PDB and topology spec |
| `specs/04-crd-api-spec.md` | Documented PDB and topology spread constraints |

### API changes

**New `spec.podDisruptionBudget` field:**

```yaml
podDisruptionBudget:
  minAvailable: "50%"    # or absolute: 2
  # or
  maxUnavailable: "25%"  # or absolute: 1
```

**New `spec.topologySpreadConstraints` field:**

```yaml
topologySpreadConstraints:
  - maxSkew: 1
    topologyKey: topology.kubernetes.io/zone
    whenUnsatisfiable: DoNotSchedule  # or ScheduleAnyway
    labelSelector:  # optional, defaults to Application labels
      matchLabels:
        app.kubernetes.io/name: checkout
```

### PDB semantics

- Only one of `minAvailable` OR `maxUnavailable` should be set
- Values can be absolute (`int32`) or percentage (`string` with `%`)
- When omitted, no PDB is created
- PDB is owned by Application (garbage collected on delete)

### Topology Spread semantics

- `maxSkew`: Maximum difference in pod count between domains (required, >= 1)
- `topologyKey`: Node label key for domain (e.g., `topology.kubernetes.io/zone`, `kubernetes.io/hostname`)
- `whenUnsatisfiable`: `DoNotSchedule` (default, pod stays Pending) or `ScheduleAnyway` (schedule but log warning)
- `labelSelector`: Optional, defaults to Application's managed labels

### Verification

```bash
# Build and test
go build ./... && go test ./...

# Regenerate CRD
make manifests

# Apply CRD
kubectl apply -f config/crd/bases/platform.kubeforge.io_applications.yaml

# Check CRD
kubectl get crd applications.platform.kubeforge.io -o yaml | grep -A 20 "podDisruptionBudget:"
kubectl get crd applications.platform.kubeforge.io -o yaml | grep -A 20 "topologySpreadConstraints:"

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
  podDisruptionBudget:
    minAvailable: "50%"
  topologySpreadConstraints:
    - maxSkew: 1
      topologyKey: topology.kubernetes.io/zone
      whenUnsatisfiable: DoNotSchedule
EOF

# Verify PDB created
kubectl get pdb -n apps checkout -o yaml

# Verify deployment has topologySpreadConstraints
kubectl get deployment -n apps checkout -o yaml | grep -A 10 "topologySpreadConstraints"
```

### Phase 3 Progress

| Item | Status |
|---|---|
| 1. Probes/defaults | ✅ Done |
| 2. Resource defaults | ✅ Done |
| 3. PodDisruptionBudget | ✅ Done |
| 4. Topology spread constraints | ✅ Done |
| 5. NetworkPolicy generation | ⏳ Next |
| 6. HPA integration | ⏳ |
| 7. KEDA integration | ⏳ |
| 8. Argo Rollouts integration | ⏳ |