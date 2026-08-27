# System Architecture

## Architectural principles

1. Kubernetes remains the source of runtime truth.
2. Git is the desired-state source of truth for deployed configuration.
3. The operator is deterministic and never calls an LLM during reconciliation.
4. AI is separated from the control plane.
5. MCP is an interoperability boundary for AI clients/agents.
6. Integrate upstream systems rather than reimplementing them.
7. Local development must not require AWS.
8. Optional infrastructure integrations must not contaminate the core API.

## Logical architecture

```text
                                 +-------------------+
                                 |      GitHub       |
                                 | source + manifests|
                                 +---------+---------+
                                           |
                              GitHub Actions / CI
                                           |
                                           v
                                     +-----+-----+
                                     |    GHCR   |
                                     +-----------+
                                           |
                                  Helm / Kustomize
                                           |
                                           v
                                     +-----+-----+
                                     |   Argo CD  |
                                     +-----+-----+
                                           |
                                           v
                              +------------+-------------+
                              |     Kubernetes cluster   |
                              |                          |
                              |  +--------------------+  |
                              |  | Application CRD    |  |
                              |  +---------+----------+  |
                              |            |             |
                              |            v             |
                              |  +---------+----------+  |
                              |  | Platform Operator  |  |
                              |  |       Go            |  |
                              |  +---+-------+------+--+  |
                              |      |       |      |     |
                              |      v       v      v     |
                              | Deployment  Service  ...   |
                              |                          |
                              +------------+-------------+
                                           |
                                  metrics/traces/logs/events
                                           |
                         +-----------------+------------------+
                         |                                    |
                         v                                    v
                  +------+-------+                      +-----+------+
                  | Prometheus   |                      | OTel/Logs  |
                  +------+-------+                      +-----+------+
                         |                                    |
                         +----------------+-------------------+
                                          |
                                          v
                                  +-------+--------+
                                  |  AI Platform   |
                                  | (separate svc) |
                                  +-------+--------+
                                          |
                              +-----------+-----------+
                              |                       |
                              v                       v
                        +-----+------+          +-----+------+
                        | MCP Server |          | LLM Provider|
                        | tools      |          | adapter(s)  |
                        +-----+------+          +------------+
                              |
                              v
                  Claude / Cursor / other MCP hosts
```

## Component boundaries

### Operator

Owns:

- CRD lifecycle.
- desired/observed-state reconciliation.
- Kubernetes child resources.
- status and conditions.
- deterministic defaults and validation.

Must not own:

- LLM prompts.
- model selection.
- AI reasoning.
- arbitrary user commands.
- external billing databases.

### AI service

Owns:

- model provider abstraction;
- agent loop/orchestration;
- context gathering through tools;
- structured diagnosis/recommendation;
- AI telemetry;
- safety policy enforcement;
- MCP server endpoints.

### MCP server

Owns:

- protocol translation;
- tool/resource/prompt exposure;
- explicit schemas and authorization boundaries.

It should not contain business logic that belongs in the operator or observability backends.

### Observability stack

Prometheus stores time-series metrics. OpenTelemetry provides telemetry instrumentation/export. Grafana visualizes. Logs remain outside Kubernetes API storage.

## Deployment topology

### Local

Kind + Docker + Helm + Argo CD + Prometheus/OTel/Grafana + optional Ollama or another reachable LLM endpoint.

### Optional cloud

The same Helm chart should be deployable to EKS without changing the core CRD. Cloud-specific credentials/configuration must be externalized.
