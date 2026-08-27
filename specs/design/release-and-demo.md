# Public Demo Design

## Goal

A contributor should be able to demonstrate the project without AWS credentials.

## Demo stack

- Kind
- Helm
- Argo CD
- Platform Operator
- Example `Application`
- Prometheus
- Grafana
- OTel Collector
- optional local LLM runtime
- MCP-capable client

## Demo scenario

1. Install platform.
2. Deploy `checkout` example.
3. Verify Argo CD sync.
4. Observe Application CR and child resources.
5. Break the workload intentionally (for example, use an image tag that cannot be pulled).
6. Observe status/events/metrics.
7. Ask the AI assistant why the application is unhealthy.
8. Agent gathers evidence through MCP tools.
9. LLM produces a grounded diagnosis.
10. Restore the workload through Git and let Argo CD reconcile it.

## Why this demo

It proves the entire architecture without requiring a cloud account and demonstrates the separation between deterministic reconciliation, GitOps, observability, and AI reasoning.
