# Scope and Roadmap

## Completed Phases

### Phase 1 — Kubernetes Operator MVP ✅

- Go module with Kubebuilder/controller-runtime
- `Application` CRD (namespaced, short name: `app`)
- Structural OpenAPI schema validation
- Deployment, Service, PDB, NetworkPolicy, HPA reconciliation
- KEDA ScaledObject and Argo Rollouts integration
- Owner references and garbage collection
- Status/conditions with phase tracking (Ready/Progressing/Degraded)
- Observed generation
- Idempotent reconciliation
- Kubernetes API error handling and safe requeue behavior
- Unit tests and envtest reconciliation tests
- Example manifests

### Phase 2 — GitOps and Delivery ✅

- Helm chart for installation
- Kustomize examples/overlays
- GitHub Actions for lint/test/build
- Container image publication to GHCR (multi-arch: amd64+arm64)
- Release tags and immutable versioning
- ArgoCD ApplicationSet for dev/prod environments
- Local Kind-based E2E pipeline
- Security scanning for source and image

### Phase 3 — Reliability and Scaling ✅

- Probes and resource defaults
- PodDisruptionBudget
- Topology spread constraints
- NetworkPolicy generation
- HPA integration
- KEDA scale-to-zero autoscaling
- Argo Rollouts for canary/blue-green deployments

### Phase 4 — Observability ✅

- Custom Prometheus metrics (reconcile rate, duration, errors, application phase)
- Grafana dashboard definitions
- Structured JSON logging
- Kubernetes events on phase transitions
- Metrics server on :8443

## Future Phases (Deferred)

### Phase 5 — AI and MCP

> Not implemented. Deferred to a separate project.

Would include:
- Separate AI service
- Configurable model provider
- MCP server exposing read-only platform tools
- Evidence-grounded incident/triage analysis

### Phase 6 — Optional Advanced Integrations

> Not implemented. Deferred for future consideration.

Would include:
- Lightweight IDP/UI
- AWS/EKS integration
- Kafka/Elasticsearch log pipeline
- Multi-cluster support
- Advanced policy/security integrations

## Scope Rule

A feature is not MVP merely because it is interesting. If it does not strengthen the core story of Kubernetes control plane + GitOps + observability, defer it.
