# Testing Strategy

## Test pyramid

```text
                 +----------------+
                 | AI evaluations  |
                 +-------+--------+
                         |
                 +-------+--------+
                 |   Kind E2E      |
                 +-------+--------+
                         |
                 +-------+--------+
                 | envtest / integ |
                 +-------+--------+
                         |
                 +-------+--------+
                 |    unit tests   |
                 +-----------------+
```

## Operator unit tests

Test pure functions independently:

- defaulting/calculation;
- desired Deployment generation;
- desired Service generation;
- condition calculation;
- error classification;
- resource naming/labels.

## envtest

Use envtest for API/reconciler semantics without requiring a full workload cluster.

Minimum scenarios:

1. Application creates Deployment.
2. Application creates Service when requested.
3. Application status becomes progressing/ready appropriately.
4. Spec update changes child desired state.
5. Child deletion is recovered.
6. Invalid spec is rejected by API validation.

## Kind E2E

Run a real local cluster.

Scenarios:

1. install CRDs/operator;
2. apply example Application;
3. wait for Deployment readiness;
4. verify Service;
5. update image/replicas;
6. simulate failure;
7. verify recovery;
8. uninstall/upgrade chart.

Future E2E scenarios:

- Argo CD sync;
- KEDA scaling;
- Argo Rollouts;
- OTel/Prometheus telemetry.

## Helm tests

- `helm lint`;
- `helm template`;
- schema validation where available;
- render expected RBAC/service account/metrics resources.

## AI tool contract tests

For every MCP tool:

- valid input;
- invalid input;
- authorization denied;
- timeout;
- empty result;
- oversized result;
- upstream failure;
- redaction.

## AI evaluation tests

Create a fixture dataset containing known scenarios such as:

- CrashLoopBackOff with OOMKilled;
- Pending pod due to resource shortage;
- failed image pull;
- rollout regression;
- service unavailable because endpoints are absent;
- no evidence / insufficient telemetry.

The evaluator should score:

- evidence usage;
- factual consistency with fixtures;
- correct identification of uncertainty;
- safe recommendation;
- tool selection efficiency;
- refusal to perform unauthorized actions.

Do not measure only textual similarity to a reference answer.

## Regression policy

Every production bug becomes a regression test at the lowest practical layer.

Every AI safety incident becomes a deterministic evaluation case before the affected behavior is changed again.
