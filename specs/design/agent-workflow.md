# AI Agent Workflow Design

## Primary workflow: diagnose application health

### Input

```text
Diagnose namespace/apps/application/checkout.
```

### Planner/tool loop

1. Resolve the application.
2. Inspect its status/conditions.
3. Inspect owned Deployment.
4. Inspect pods and container states.
5. Inspect recent Kubernetes events.
6. Query a bounded set of relevant Prometheus metrics.
7. Inspect Argo CD state if configured.
8. Inspect recent Git/deployment metadata if available.
9. Synthesize evidence.
10. Return diagnosis and recommended next steps.

## Tool selection

Tools should be narrow and composable. Prefer:

```text
get_application
get_deployment
get_pods
get_events
query_prometheus
get_argocd_status
get_recent_change
```

over a single unrestricted `run_kubectl` tool.

## Evidence model

Each evidence item should carry enough metadata to be independently inspected:

```json
{
  "source": "kubernetes.event",
  "resource": "apps/checkout",
  "observed_at": "...",
  "fact": "FailedScheduling: insufficient memory",
  "reference": "..."
}
```

The agent should not fabricate `reference` values. If the backend cannot provide a stable reference, omit it.

## Decision output

The model output should distinguish:

- **Facts** — directly observed tool results.
- **Inference** — model interpretation.
- **Recommendation** — proposed action.
- **Unknowns** — missing evidence.

## Safety rule

A recommendation is not an action.

The default agent is successful when it produces an accurate, evidence-grounded diagnosis, even if no mutation is available.
