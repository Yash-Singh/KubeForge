# Agent Build Plan

This is the execution queue for an implementation agent. Do not skip ahead unless the current phase's Definition of Done is satisfied.

## Phase 1 — Operator MVP

### Task 1 — Scaffold

Create the Kubebuilder/controller-runtime project for the provisional API group and `Application` kind.

Deliver:

- Go module;
- controller;
- API types;
- generated CRD/manifests;
- RBAC markers;
- basic manager configuration.

### Task 2 — API

Implement only the fields in `specs/04-crd-api-spec.md`.

Add schema validation and generated manifests.

### Task 3 — Deployment reconciliation

Implement a deterministic desired Deployment builder and reconcile it from the `Application` resource.

### Task 4 — Service reconciliation

Create/manage the Service only when configured.

### Task 5 — Status

Implement Ready/Progressing/Degraded conditions, replica counts, and observed generation.

### Task 6 — Tests

Add unit and envtest coverage from `specs/10-testing-strategy.md`.

### Phase 1 Definition of Done

```text
[ ] application CR accepted
[ ] deployment created
[ ] service created when requested
[ ] status is meaningful
[ ] child drift is repaired
[ ] repeated reconciliation is idempotent
[ ] tests pass
[ ] examples work on Kind
```

## Phase 2 — Delivery

### Task 7 — Helm

Package the operator as a Helm chart.

### Task 8 — CI

Implement PR checks and main-branch checks from `specs/06-gitops-cicd.md`.

### Task 9 — Image release

Publish versioned images to GHCR and create reproducible release artifacts.

### Task 10 — Argo CD

Add a reproducible Argo CD deployment example.

### Phase 2 Definition of Done

```text
[ ] PR CI is green
[ ] image published
[ ] Helm chart lint/template passes
[ ] fresh Kind cluster can be installed from chart
[ ] Argo CD can sync the platform from Git
[ ] release artifact can be reinstalled reproducibly
```

## Phase 3 — Reliability and scaling

Implement one integration at a time. Before each integration, add its API contract and an ADR if ownership/source-of-truth changes.

Suggested order:

```text
probes -> resources -> PDB -> topology -> NetworkPolicy -> HPA/KEDA -> Rollouts
```

Do not implement all of these as a single feature.

## Phase 4 — Observability

Add:

- operator metrics;
- structured logs;
- OTel traces;
- Prometheus/Grafana demo configuration;
- basic deployment correlation.

## Phase 5 — AI/MCP

### Task 1 — AI service shell

Create a separate service. It must run independently of the operator.

### Task 2 — ModelProvider abstraction

Implement the internal provider interface described in `specs/08-ai-mcp-spec.md`.

No business logic may depend on a specific vendor.

### Task 3 — Read-only platform tools

Implement tool adapters for Kubernetes first. Add Prometheus and Argo CD after the tool contract is stable.

### Task 4 — Agent workflow

Implement the triage workflow in `specs/design/agent-workflow.md`.

### Task 5 — MCP

Expose the tools using the supported MCP SDK and target the MCP 2026-07-28 specification.

### Task 6 — AI evaluations

Add fixture-based diagnosis tests and safety tests before adding any write tool.

### Phase 5 Definition of Done

```text
[ ] local model can be configured
[ ] remote provider can be configured
[ ] changing model/provider does not change tools
[ ] AI can diagnose known fixture failures
[ ] evidence is included
[ ] missing evidence is acknowledged
[ ] unsafe/unavailable actions are rejected
[ ] MCP read-only tools work from a compatible client
[ ] AI calls and tool calls are observable
```

## Phase 6 — Optional

Only implement after reviewing project adoption and contributor feedback:

- IDP/UI;
- GitHub PR creation;
- human-approved remediation;
- AWS/EKS demo;
- Kafka/Elasticsearch/OpenSearch log pipeline;
- multi-cluster;
- advanced policy/security integrations.

## Agent stop conditions

Stop implementation and ask the human when:

- a missing API contract changes public behavior;
- a secret/data boundary is unclear;
- a controller ownership conflict appears;
- a write-capable AI tool is requested without an authorization design;
- a cloud provider becomes a required dependency;
- a new CRD field has no clear semantics;
- a task contradicts an accepted ADR.

Do not resolve these by guessing.
