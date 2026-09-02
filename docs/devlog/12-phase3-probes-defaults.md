# Phase 3 — Reliability and Scaling

## Task 12 — Probes and Resource Defaults

### What was done

Added liveness/readiness probe support and sensible resource defaults to the Application CRD and operator.

### Files changed

| File | Change |
|---|---|
| `api/v1alpha1/application_types.go` | Added `ProbesSpec`, `ProbeSpec`, `HTTPGetAction`, `HTTPHeader` types |
| `internal/controller/application_builder.go` | Added probe building, default resources (10m/64Mi requests, 500m/128Mi limits) |
| `config/crd/bases/platform.kubeforge.io_applications.yaml` | Regenerated with probes spec |
| `specs/04-crd-api-spec.md` | Documented probes and resource defaults |

### API changes

**New `spec.probes` field:**

```yaml
probes:
  liveness:
    httpGet:
      path: /healthz
      port: 8080
    initialDelaySeconds: 15
    periodSeconds: 20
  readiness:
    httpGet:
      path: /readyz
      port: 8080
    initialDelaySeconds: 5
    periodSeconds: 10
```

**Probe fields (all optional):**
- `httpGet.path` (required when httpGet used)
- `httpGet.port` (required when httpGet used)
- `httpGet.scheme` (default: HTTP)
- `httpGet.httpHeaders` (optional)
- `initialDelaySeconds` (default: 10)
- `periodSeconds` (default: 10)
- `timeoutSeconds` (default: 1)
- `failureThreshold` (default: 3)
- `successThreshold` (default: 1)

**Resource defaults (when `spec.resources` omitted):**
- Requests: CPU 10m, Memory 64Mi
- Limits: CPU 500m, Memory 128Mi

### Verification

```bash
# Build and test
go build ./... && go test ./...

# Regenerate CRD
make manifests

# Check CRD includes probes
kubectl apply -f config/crd/bases/platform.kubeforge.io_applications.yaml
kubectl get crd applications.platform.kubeforge.io -o yaml | grep -A 5 "probes:"
```

### Test example

```yaml
apiVersion: platform.kubeforge.io/v1alpha1
kind: Application
metadata:
  name: checkout
  namespace: apps
spec:
  image: ghcr.io/example/checkout:v1.0.0
  replicas: 2
  probes:
    liveness:
      httpGet:
        path: /healthz
        port: 8080
      initialDelaySeconds: 15
      periodSeconds: 20
    readiness:
      httpGet:
        path: /readyz
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 10
```

### Next Phase 3 items

1. ✅ Probes/defaults
2. ⏳ Resource defaults (done as part of this)
3. ⏳ PodDisruptionBudget
4. ⏳ Topology spread constraints
5. ⏳ NetworkPolicy generation (optional)
6. ⏳ HPA integration
7. ⏳ KEDA integration
8. ⏳ Argo Rollouts integration