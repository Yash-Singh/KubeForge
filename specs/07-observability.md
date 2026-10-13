# Observability Specification

## Goals

The platform should make it possible to answer:

1. Is the operator healthy?
2. Is an Application converging?
3. How long does reconciliation take?
4. Why is an Application degraded?
5. What changed before a failure?

## Metrics

Use Prometheus-compatible metrics.

### Operator Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `kubeforge_reconcile_total` | Counter | namespace, name, result | Total reconciliations |
| `kubeforge_reconcile_duration_seconds` | Histogram | namespace, name | Reconciliation duration |
| `kubeforge_reconcile_errors_total` | Counter | namespace, name, error_type | Reconcile errors |
| `kubeforge_application_phase` | Gauge | namespace, name | Current phase (1=Ready, 2=Progressing, 3=Degraded) |
| `kubeforge_application_replicas` | Gauge | namespace, name, state | Replica counts |

Controller-runtime metrics are retained where useful.

Avoid high-cardinality labels such as unrestricted Pod UID, log line, request ID, or user text.

## Logs

Emit structured JSON logs with:

- timestamp
- level
- controller/component
- namespace/name where applicable
- reconcile/request correlation ID when available
- error classification
- safe, bounded message

Never log:

- credentials
- bearer tokens
- secrets
- full HTTP authorization headers

## Kubernetes Events

Record events on Application CRs during phase transitions:

- `PhaseReady` — Application is ready
- `PhaseProgressing` — Application is progressing
- `PhaseDegraded` — Application is degraded
- `ReconcileError` — Reconciliation failed

## Traces

OpenTelemetry traces are deferred to a future phase.

Important spans for future implementation:

```text
operator.reconcile
   |
   +--> kubernetes.get
   +--> kubernetes.patch
   +--> status.update
```

## Correlation

The long-term goal is to correlate:

```text
Git commit
   ↓
Argo sync
   ↓
Application generation
   ↓
Operator reconcile
   ↓
Deployment rollout
   ↓
Pod readiness
   ↓
Runtime metrics/logs
```

The first implementation correlates only the identifiers that are reliably available.

## Grafana

Pre-built dashboards are available in `charts/platform-operator/dashboards/`:

- Controller health
- Reconcile rate/errors/latency
- Applications by state
- Replica counts

## Logs Backend

The core project does not require Elasticsearch or Sumo Logic. Local development works with stdout + OTel Collector and an optional log backend.

Kafka/Elasticsearch belong to a future scale-oriented adapter/demo.
