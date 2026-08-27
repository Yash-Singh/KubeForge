# Observability Specification

## Goals

The platform should make it possible to answer:

1. Is the operator healthy?
2. Is an Application converging?
3. How long does reconciliation take?
4. Why is an Application degraded?
5. What changed before a failure?
6. What evidence did the AI use?

## Metrics

Use Prometheus-compatible metrics.

### Operator metrics

Recommended baseline:

- `platform_operator_reconcile_total`
- `platform_operator_reconcile_errors_total`
- `platform_operator_reconcile_duration_seconds`
- `platform_operator_applications_total`
- `platform_operator_applications_ready`
- `platform_operator_applications_degraded`

Controller-runtime metrics should be retained where useful rather than duplicating all framework metrics.

Avoid high-cardinality labels such as unrestricted Pod UID, log line, request ID, or user text.

## Logs

Emit structured JSON logs with:

- timestamp;
- level;
- controller/component;
- namespace/name where applicable;
- reconcile/request correlation ID when available;
- error classification;
- safe, bounded message.

Never log:

- credentials;
- bearer tokens;
- secrets;
- full HTTP authorization headers;
- arbitrary user prompts containing confidential data.

## Traces

Use OpenTelemetry for the operator and AI service where practical.

Important spans:

```text
Git/Argo deployment context
        |
        v
operator.reconcile
   |
   +--> kubernetes.get
   +--> kubernetes.patch
   +--> status.update
```

AI service:

```text
agent.invoke
   |
   +--> tool.kubernetes.get_pods
   +--> tool.kubernetes.get_events
   +--> tool.prometheus.query
   +--> tool.argocd.get_application
   +--> llm.invoke
```

Use OpenTelemetry GenAI semantic conventions where applicable and stable enough for the selected instrumentation. Keep sensitive prompt/completion content out of telemetry by default; capture metadata such as provider/model, operation, duration, token counts and error class according to the current conventions and privacy policy.

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

The first implementation may correlate only the identifiers that are reliably available. Do not fabricate timestamps or causal links.

## Grafana

Provide a dashboard JSON or provisioning config with:

- controller health;
- reconcile rate/errors/latency;
- Applications by state;
- queue depth if exposed;
- optional AI request/tool-call health.

## Logs backend

The core project does not require Elasticsearch or Sumo Logic. Local development should work with stdout + OTel Collector and an optional log backend.

Kafka/Elasticsearch belong to a future scale-oriented adapter/demo.
