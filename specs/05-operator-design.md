# Operator Design

## Controller model

Use Kubebuilder/controller-runtime.

The reconciliation loop is a **level-based control loop**. A reconcile must converge the actual cluster state toward the desired state represented by the `Application` resource. It must be safe to execute repeatedly and after unrelated/redundant events.

## Reconciliation outline

```text
Reconcile(Application)
        |
        +--> fetch Application
        |
        +--> if not found: return
        |
        +--> validate/default desired state
        |
        +--> reconcile Deployment
        |
        +--> reconcile Service if enabled
        |
        +--> observe owned resources
        |
        +--> calculate conditions / replica counts
        |
        +--> update status if changed
        |
        +--> return success or safe retry/error
```

## Idempotency

Repeated reconcile calls must not create duplicate objects or continuously update unchanged fields.

Prefer declarative builders/helpers that calculate desired child state and apply only necessary changes.

Avoid imperative event scripts such as "when PodDeleted, create Pod". The Deployment controller already owns pod lifecycle. The application operator should own the Deployment desired state.

## Ownership

Use owner references for directly owned namespaced children when appropriate, enabling Kubernetes garbage collection and watch relationships.

The owner relationship should be:

```text
Application
  ├── Deployment
  └── Service
```

When a future integration delegates lifecycle to another controller (for example KEDA or Argo Rollouts), explicitly document which controller owns which object and avoid controller fights.

## Update strategy

Changes to `.spec` should lead to child resource convergence.

The controller should observe `metadata.generation` and publish `status.observedGeneration` after it has processed the requested generation.

## Status handling

Status must distinguish:

- desired state accepted;
- resources progressing;
- resources ready;
- resources degraded/error.

Status writes should avoid hot loops. Only update status when values actually change.

## Error handling

Classify errors:

### Permanent/user configuration errors

Example: invalid configuration that passed unexpectedly through API validation.

Action: set `Degraded=True`, give an actionable message, avoid tight retries.

### Transient API/network errors

Action: return/requeue using controller-runtime error/retry behavior.

### Dependent resource not ready

Action: update `Progressing=True` and requeue/watch the relevant resource.

## Resource watches

The controller should primarily watch:

- `Application` objects;
- owned Deployment;
- owned Service if relevant.

Add broader watches only when required for correct level-based reconciliation.

## Concurrency and scaling

The controller must be safe under concurrent reconciles as configured by controller-runtime. Do not rely on process-local mutable state as source of truth.

Metrics should allow evaluation of:

- reconcile rate;
- reconcile errors;
- reconcile latency;
- workqueue depth/latency;
- number of Applications by phase.

## No AI in reconcile

Never invoke an LLM from `Reconcile()`.

Reasons:

- unpredictable latency;
- non-deterministic output;
- external provider failure;
- token/cost variability;
- possible prompt/data leakage;
- difficulty guaranteeing convergence.

AI belongs in a separate service that can inspect platform state and produce advisory output.

## Optional integrations rule

Each future integration must answer:

1. What user intent does this field express?
2. Which controller owns the generated object?
3. What is the source of truth?
4. What happens on disable/removal?
5. How is drift detected?
6. What happens during dependency failure?
7. How is the behavior observed/tested?

Do not merge an integration that cannot answer all seven.
