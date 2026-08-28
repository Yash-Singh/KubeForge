# Kubernetes API / CRD Specification

## API group

**Final:** `platform.kubeforge.io`

## Kind

`Application`

## Scope

Namespaced.

## Version

`v1alpha1` for the initial public API.

## Design intent

The resource describes **application intent**, not every possible Kubernetes field. Avoid turning the CRD into a generic Deployment replacement.

## Initial API

```yaml
apiVersion: platform.kubeforge.io/v1alpha1
kind: Application
metadata:
  name: checkout
  namespace: apps
spec:
  image: ghcr.io/example/checkout:v1.0.0
  replicas: 2
  service:
    port: 8080
    targetPort: 8080
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 500m
      memory: 512Mi
  env:
    - name: LOG_LEVEL
      value: info
```

## Initial `spec`

### `image` — required

Container image reference for the primary application container.

### `replicas` — optional

Positive integer. Initial default: **1**.

Scaling integrations must define whether this field remains authoritative when autoscaling is enabled. Until that contract is finalized, autoscaling should be opt-in and must not silently fight the replicas field.

### `service` — optional

- `port` — required when `service` is present.
- `targetPort` — optional; default to `port`.

### `resources` — optional

Requests and limits for the primary container.

### `env` — optional

Non-secret environment variables only. Secret references are deferred until a secure secret reference contract is designed.

## Explicitly deferred fields

Do not add these to `v1alpha1` until a concrete use case and ownership model exists:

- arbitrary PodSpec passthrough;
- cloud provider settings;
- arbitrary annotations/labels copied everywhere;
- secret values;
- autoscaling implementation details;
- rollout step syntax;
- network policy graph syntax;
- cost/billing data;
- AI prompts or model configuration.

## Status contract

The operator should publish concise operational state.

```yaml
status:
  observedGeneration: 3
  phase: Ready
  readyReplicas: 2
  desiredReplicas: 2
  conditions:
    - type: Ready
      status: "True"
      reason: ResourcesReady
      message: Application resources are ready.
      observedGeneration: 3
      lastTransitionTime: "..."
```

Recommended condition types:

- `Ready`
- `Progressing`
- `Degraded`

Avoid excessively verbose status. Do not put logs or large event histories into `.status`.

## API conventions

- Use Kubernetes quantity/int types where appropriate.
- Apply structural schemas.
- Add kubebuilder validation markers and generated CRDs.
- Keep defaults deterministic and documented.
- Evolve via additive changes where possible.
- Use conversion/webhooks only when versioning actually requires them.

## Child resources

Initial owned resources:

- Deployment.
- Service when `spec.service` is present.

Future optional resources may include PDB, HPA/KEDA, NetworkPolicy, ServiceMonitor, or Argo Rollout, each behind explicit feature/API fields.

## CLI examples

```bash
kubectl get applications -A
kubectl get application checkout -n apps -o yaml
kubectl describe application checkout -n apps
```
