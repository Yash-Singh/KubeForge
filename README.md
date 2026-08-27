# KubeForge — AI-Assisted Kubernetes Application Platform

> Spec-driven development source of truth for an open-source Kubernetes platform project.

## Purpose

This repository specification describes a Kubernetes-native application platform built around a Go-based operator, GitOps delivery, Kubernetes/Prometheus/OpenTelemetry observability, and an optional AI intelligence layer exposed through MCP. The project is intended to be useful as a real platform-engineering system, a learning vehicle, and a credible open-source project.

## Core principle

**The Kubernetes operator is the deterministic control plane. The AI layer is advisory and tool-driven. MCP is the interoperability layer between AI clients/agents and platform capabilities.**

Do not replace deterministic Kubernetes reconciliation with LLM decisions.

## Specification authority

Priority order when implementation choices conflict:

1. Kubernetes API conventions and upstream behavior.
2. Stable upstream project contracts (controller-runtime/Kubebuilder, Argo CD, Prometheus, OpenTelemetry, MCP).
3. This specification.
4. Implementation convenience.

When a requirement is ambiguous, do not invent behavior. Record the ambiguity in `specs/OPEN-QUESTIONS.md`, preserve backward compatibility, and prefer the smallest safe implementation.

## Documents

- `specs/01-product-spec.md` — problem, users, goals, non-goals, success criteria.
- `specs/02-scope-roadmap.md` — phases and strict scope boundaries.
- `specs/03-system-architecture.md` — component and runtime architecture.
- `specs/04-crd-api-spec.md` — initial Kubernetes API contract.
- `specs/05-operator-design.md` — reconciliation, ownership, failure, status, scaling.
- `specs/06-gitops-cicd.md` — GitHub, CI, image/release, Helm/Kustomize, Argo CD.
- `specs/07-observability.md` — metrics, traces, logs, events and correlations.
- `specs/08-ai-mcp-spec.md` — model-agnostic AI architecture, MCP server, tools, guardrails.
- `specs/09-security.md` — RBAC, AI action safety, secrets, supply chain, boundaries.
- `specs/10-testing-strategy.md` — unit, envtest, integration, E2E, AI evals.
- `specs/OPEN-QUESTIONS.md` — deliberately unresolved decisions; do not silently guess.
- `specs/adr/` — architecture decision records.
- `specs/design/` — agent workflow and release/demo design.
- `agents/` — model-agnostic runtime agent definitions.
- `BUILD.md` — phased execution plan for implementation agents.
- `AGENTS.md` — implementation rules for coding agents.

## Recommended implementation order

1. Phase 1: operator + CRD + tests.
2. Phase 2: Helm/Kustomize + GitHub Actions + GHCR + Argo CD + E2E.
3. Phase 3: reliability/scaling integrations (KEDA/HPA/PDB/probes/topology/rollouts as individually justified features).
4. Phase 4: Prometheus + Grafana + OTel + correlated deployment telemetry.
5. Phase 5: AI service + MCP read-only tools + incident/triage analysis.
6. Phase 6: optional human-approved write actions, IDP, AWS/EKS, Kafka/Elasticsearch, multi-cluster.

## Current explicit non-goals

- No requirement for AWS in local development.
- No requirement to run Kafka or Elasticsearch in the MVP.
- No requirement to build an IDP/UI in the MVP.
- No autonomous production-changing AI actions.
- No cloud-specific API in the operator core.
- No model vendor hard dependency in core AI interfaces.

## External references

The specification is aligned with current upstream guidance, including Kubernetes custom resources/controllers, Kubebuilder reconciliation practices, Argo CD GitOps, OpenTelemetry semantic conventions, and the MCP 2026-07-28 specification.

- Kubernetes Custom Resources: https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/
- Kubernetes Operator pattern: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
- Kubebuilder good practices: https://book.kubebuilder.io/reference/good-practices
- Argo CD: https://argo-cd.readthedocs.io/en/stable/
- OpenTelemetry semantic conventions: https://opentelemetry.io/docs/specs/semconv/
- MCP 2026-07-28 specification: https://modelcontextprotocol.io/specification/2026-07-28
