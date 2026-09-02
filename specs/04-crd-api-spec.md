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
  replicas: 3
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
  podDisruptionBudget:
    minAvailable: "50%"
  topologySpreadConstraints:
    - maxSkew: 1
      topologyKey: topology.kubernetes.io/zone
      whenUnsatisfiable: DoNotSchedule
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

When omitted, the operator applies sensible defaults:
- Requests: CPU 10m, Memory 64Mi
- Limits: CPU 500m, Memory 128Mi

These defaults are applied only when no resources are specified in the spec.

### `env` — optional

Non-secret environment variables only. Secret references are deferred until a secure secret reference contract is designed.

### `probes` — optional

Liveness and readiness probe configuration for the primary container.

- `liveness` — optional. Configures the liveness probe.
- `readiness` — optional. Configures the readiness probe.

Each probe supports:

- `httpGet` — optional. HTTP GET action to perform.
  - `path` — required. HTTP path to GET.
  - `port` — required. Port number to connect to.
  - `scheme` — optional. HTTP or HTTPS (default: HTTP).
  - `httpHeaders` — optional. Custom headers.

- `initialDelaySeconds` — optional. Seconds after container start before probe initiates (default: 10).
- `periodSeconds` — optional. How often to perform the probe (default: 10).
- `timeoutSeconds` — optional. Probe timeout in seconds (default: 1).
- `failureThreshold` — optional. Consecutive failures before probe fails (default: 3).
- `successThreshold` — optional. Consecutive successes before probe succeeds (default: 1).

When `probes` is omitted, no probes are configured. When present but fields omitted, Kubernetes defaults apply.

### `podDisruptionBudget` — optional

Configures a PodDisruptionBudget for high availability during voluntary disruptions (node drains, upgrades).

- `minAvailable` — optional. Minimum pods that must be available (absolute or percentage, e.g., "50%").
- `maxUnavailable` — optional. Maximum pods that can be unavailable (absolute or percentage, e.g., "25%").

Only one of `minAvailable` or `maxUnavailable` should be specified. When omitted, no PDB is created.

### `topologySpreadConstraints` — optional

Array of topology spread constraints to distribute pods across failure domains.

Each constraint:

- `maxSkew` — required. Maximum allowed difference in pod count between topology domains.
- `topologyKey` — required. Node label key for topology domain (e.g., `topology.kubernetes.io/zone`, `kubernetes.io/hostname`).
- `whenUnsatisfiable` — optional. Behavior when constraint can't be satisfied: `DoNotSchedule` (default) or `ScheduleAnyway`.
- `labelSelector` — optional. Pods to include in spread calculation. Defaults to Application-managed labels.

When omitted, no topology spread constraints are applied.

### `networkPolicy` — optional

Configures NetworkPolicy for application pods to control ingress/egress traffic.

- `ingress` — optional. Array of ingress rules.
  - `ports` — optional. Ports to allow (protocol, port, endPort).
  - `from` — optional. Sources allowed to access pods.
    - `podSelector` — optional. Select pods in same namespace.
    - `namespaceSelector` — optional. Select namespaces.
    - `ipBlock` — optional. CIDR range with optional exceptions.

- `egress` — optional. Array of egress rules.
  - `ports` — optional. Ports to allow.
  - `to` — optional. Destinations pods can access (same structure as `from`).

- `policyTypes` — optional. Explicit policy types: `Ingress`, `Egress`. Defaults inferred from ingress/egress presence.

When omitted, no NetworkPolicy is created.

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
- PodDisruptionBudget when `spec.podDisruptionBudget` is present.
- NetworkPolicy when `spec.networkPolicy` is present.

Future optional resources may include HPA/KEDA, ServiceMonitor, or Argo Rollout, each behind explicit feature/API fields.

## CLI examples

```bash
kubectl get applications -A
kubectl get application checkout -n apps -o yaml
kubectl describe application checkout -n apps
```
