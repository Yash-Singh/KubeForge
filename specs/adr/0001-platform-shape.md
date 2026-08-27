# ADR 0001 — Platform Shape

## Status
Accepted

## Context

The project could be built as a standalone MCP server, an AI agent, a Kubernetes operator, an IDP, or a collection of unrelated infrastructure tools.

The strongest project story for platform/core-infrastructure engineering is a Kubernetes-native control plane with GitOps, observability, and AI capabilities built around it.

## Decision

Build one coherent platform project with these layers:

1. Go Kubernetes operator as deterministic core.
2. Helm/GitHub Actions/Argo CD for delivery.
3. Prometheus/OpenTelemetry/Grafana for observability.
4. Separate AI service for analysis/recommendations.
5. MCP as the interoperability interface for AI clients/agents.
6. Optional integrations added incrementally.

## Rationale

This demonstrates both established infrastructure engineering and modern AI-assisted operations without making the LLM responsible for Kubernetes convergence.

## Consequences

Positive:

- strong Kubernetes/control-plane engineering story;
- realistic GitOps lifecycle;
- AI is real and calls an LLM;
- MCP adds modern interoperability;
- optional components do not block MVP.

Negative:

- more components than a toy project;
- requires strict boundaries to avoid scope explosion.
