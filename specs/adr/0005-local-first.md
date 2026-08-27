# ADR 0005 — Local-First Development

## Status
Accepted

## Decision

The core development and test workflow must run without AWS or any paid cloud service.

Default local stack:

- Docker
- Kind
- Helm
- Argo CD
- Prometheus
- OpenTelemetry Collector
- Grafana
- optional local LLM runtime

AWS/EKS is an optional integration/demo environment.

## Rationale

This keeps contributor cost low and prevents cloud-provider APIs from coupling the project to AWS.
