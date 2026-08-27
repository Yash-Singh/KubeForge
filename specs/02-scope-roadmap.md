# Scope and Roadmap

## Phase 1 — Kubernetes Operator MVP

### Must have

- Go module.
- Kubebuilder/controller-runtime project structure.
- `Application` CRD, namespaced.
- Structural OpenAPI schema validation.
- Deployment reconciliation.
- Service reconciliation.
- Owner references.
- Status/conditions.
- Observed generation.
- Idempotent reconciliation.
- Finalizer only if a real external cleanup requirement appears; otherwise avoid one.
- Kubernetes API error handling and safe requeue behavior.
- Unit tests.
- envtest reconciliation tests.
- Example manifests.

### Not yet

- AI.
- KEDA.
- Argo Rollouts.
- AWS.
- Kafka.
- Elasticsearch.
- UI.

## Phase 2 — GitOps and Delivery

### Must have

- Helm chart for installation.
- Kustomize examples/overlays if they add useful GitOps value.
- GitHub Actions for lint/test/build.
- Container image publication to GHCR by default.
- Release tags and immutable versioning.
- Automated release notes/changelog generation may be added after basic releases work.
- Argo CD application manifest/example.
- Local Kind-based E2E pipeline.
- Security scanning for source and image.

### Design constraint

Argo CD is the deployment reconciler. Git is the desired-state source of truth for installed configuration. The operator is the reconciler of `Application` objects inside Kubernetes.

## Phase 3 — Reliability and Scaling

Implement as independent integrations, not a single large feature.

Potential order:

1. probes/defaults;
2. resource defaults;
3. PodDisruptionBudget;
4. topology spread constraints;
5. optional NetworkPolicy generation;
6. HPA integration;
7. KEDA integration;
8. Argo Rollouts integration.

Each integration needs a clear CRD contract, lifecycle ownership model, tests, and opt-in behavior.

## Phase 4 — Observability

- Prometheus operator metrics.
- Grafana dashboard definitions.
- OpenTelemetry traces for operator and AI service.
- Deployment/change correlation.
- Structured JSON logs.
- Kubernetes event/context collection for diagnosis.

Do not add Kafka/Elasticsearch just to demonstrate them. They are optional future adapters for high-volume centralized log pipelines.

## Phase 5 — AI and MCP

### MVP

- Separate AI service.
- Configurable model provider.
- No model specified in core project.
- Tool-calling agent.
- MCP server exposing safe read-only platform tools.
- Evidence-grounded incident/triage analysis.
- Structured recommendation output.
- Audit trail for AI tool calls.
- Explicit confidence/uncertainty language.

### Later

- Git diff/change correlation.
- Resource optimization recommendations.
- GitHub issue/PR generation.
- Human-approved write tools.
- AI evaluation suite with regression datasets.

## Phase 6 — Optional advanced integrations

Only after the core project is stable:

- lightweight IDP/UI;
- AWS/EKS integration;
- local cloud emulation where technically useful;
- Kafka/Elasticsearch log path;
- multi-cluster support;
- advanced policy/security integrations;
- controlled automated remediation.

## Scope rule

A feature is not MVP merely because it is interesting. If it does not strengthen the core story of Kubernetes control plane + GitOps + observability + safe AI operations, defer it.
