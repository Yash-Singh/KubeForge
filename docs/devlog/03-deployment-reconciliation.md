# Phase 1 — Operator MVP

## Tasks 3–5 — Deployment, Service, and Status Reconciliation

### What was done

Implemented the full reconcile loop: deterministic Deployment builder, Service builder, and status reporting including Ready/Progressing/Degraded phases.

### Architecture

Following `specs/05-operator-design.md`:

```
Reconcile(Application)
  ├── fetch Application
  ├── if not found: return
  ├── reconcile Deployment (create-or-patch)
  ├── reconcile Service if spec.service present
  ├── observe owned resources
  ├── calculate conditions / replica counts
  ├── update status
  └── return success or safe retry/error
```

### Builder (`internal/controller/application_builder.go`)

| Function | Purpose |
|---|---|
| `applicationLabels()` | Standard labels: `app.kubernetes.io/name`, `app.kubernetes.io/managed-by=kubeforge-operator` |
| `desiredDeploymentSpec()` | Calculates the full `DeploymentSpec` from Application spec — image, replicas, resources, env vars |
| `convertEnvVars()` | Converts `[]EnvVar` (custom non-secret type) to `[]corev1.EnvVar` |

### Reconcile loop (`internal/controller/application_controller.go`)

- **Idempotent create-or-patch** via `controllerutil.CreateOrPatch` — creates if absent, patches only mutable fields if changed
- **Owner references** via `controllerutil.SetControllerReference` — enables Kubernetes garbage collection when Application is deleted
- **Watches** owned Deployment and Service via `.Owns()` in `SetupWithManager`
- **Status update** observes the owned Deployment and sets `Application.Status`

### Service reconciliation (`internal/controller/application_controller.go`)

- Creates a `corev1.Service` named after the Application, owned by the Application
- `spec.ports[0].port` = `spec.service.port`, `targetPort` defaults to `port` when not set
- Selects pods via `app.kubernetes.io/name` label
- Skipped entirely when `spec.service` is nil

### Status reporting (`internal/controller/application_controller.go`)

| Field | Behavior |
|---|---|
| `ObservedGeneration` | Set to `app.Generation` at the start of every reconcile |
| `DesiredReplicas` | Read from `dep.Spec.Replicas` |
| `ReadyReplicas` | Read from `dep.Status.ReadyReplicas` |
| `Phase` | `Progressing` → `Ready` (when readyReplicas ≥ desired) |
| `Conditions` | `Ready`, `Progressing` (always present); `Degraded` when available |

### Degraded detection (`internal/controller/application_controller.go`)

`isDeploymentDegraded(dep)` returns `true` when the owned Deployment reports:
- `Available=False`, or
- `Progressing` condition with reason `ProgressDeadlineExceeded`

The `Degraded` condition is set with reason `DeploymentFailed` and an actionable message.

### Status phase matrix

| Phase | Condition | Trigger |
|---|---|---|
| `Progressing` | `Reconciling` | Deployment not yet created |
| `Progressing` | `Creating` | Deployment exists, 0 ready replicas |
| `Progressing` | `Scaling` | Some replicas ready but not all |
| `Ready` | `ResourcesReady` | All desired replicas ready |
| `Degraded` | `DeploymentFailed` | Deployment available=False or ProgressDeadlineExceeded |

### Files changed

| File | Change |
|---|---|
| `internal/controller/application_builder.go` | New — DeploymentSpec builder, label helpers, env converter |
| `internal/controller/application_controller.go` | Full reconcile loop, create-or-patch, service reconcile, status with Degraded |
| `config/rbac/role.yaml` | Regenerated — includes `apps/deployments` and `services` RBAC |

### Verification

- `go build ./...` — clean
- `go vet ./...` — clean
- `make manifests` — RBAC and CRD regenerated