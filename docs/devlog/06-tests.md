# Phase 1 — Operator MVP

## Task 6 — Tests (unit + envtest)

### What was done

Added unit tests for pure builder functions and envtest integration tests covering the reconciler semantics.

### Test files

| File | Type | What it tests |
|---|---|---|
| `internal/controller/application_builder_test.go` | Unit | `applicationLabels`, `desiredDeploymentSpec`, `convertEnvVars`, `isDeploymentDegraded` |
| `internal/controller/application_controller_test.go` | envtest | Full reconciler scenarios against a real API server |

### Unit tests

- `applicationLabels` — verifies expected labels are set
- `desiredDeploymentSpec` — verifies image, replicas, resources, env, selector, template labels
- `convertEnvVars` — verifies custom `EnvVar` → `corev1.EnvVar` conversion
- `isDeploymentDegraded` — verifies detection of `Available=False` and `ProgressDeadlineExceeded`

### envtest scenarios (minimum from `specs/10-testing-strategy.md`)

1. ✅ Application creates a Deployment
2. ✅ Application creates a Service when `spec.service` is configured
3. ✅ Application does not create a Service when `spec.service` is nil
4. ✅ Application status becomes Progressing/Ready appropriately
5. ✅ Spec update (image) changes child desired state
6. ✅ Child deletion is recovered (Deployment is recreated)

### Verification

- `go build ./...` — clean
- `go vet ./...` — clean
- `go test ./internal/...` — **16/16 pass**
- `make setup-envtest` — binaries downloaded to `bin/k8s/`

### Notes

- Both `Describe` blocks across `application_builder_test.go` and `application_controller_test.go` run under a single `TestControllers` entry point in `application_controller_test.go`. Ginkgo discovers all specs in the package automatically.
- `AfterEach` cleans up Service, Deployment, and Application in that order (Service/Deployment have owner references to Application).