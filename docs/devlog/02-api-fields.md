# Phase 1 — Operator MVP

## Task 2 — API Fields

### What was done

Replaced the scaffold `Foo` field with spec-defined fields from `specs/04-crd-api-spec.md`.

### ApplicationSpec fields

| Field | Type | Required | Default | Notes |
|---|---|---|---|---|
| `image` | `string` | Yes | — | Container image reference |
| `replicas` | `*int32` | No | `1` (kubebuilder default) | Min 1 |
| `service` | `*ServiceSpec` | No | — | `port` required when present; `targetPort` defaults to `port` |
| `resources` | `*corev1.ResourceRequirements` | No | — | Standard Kubernetes requests/limits |
| `env` | `[]EnvVar` | No | — | Non-secret only — custom type (not `corev1.EnvVar`; no `ValueFrom`) |

### Custom types added

- `ServiceSpec` — `port` (int32, 1–65535), `targetPort` (optional *int32, 1–65535)
- `EnvVar` — `name`, `value` (strings only — no secret/field refs)

### ApplicationStatus fields

| Field | Type | Notes |
|---|---|---|
| `observedGeneration` | `int64` | Last reconciled generation |
| `phase` | `ApplicationPhase` | One of `""`, `Ready`, `Progressing`, `Degraded` |
| `readyReplicas` | `int32` | Observed ready pods |
| `desiredReplicas` | `int32` | Last reconciled replica count |
| `conditions` | `[]metav1.Condition` | Types: `Ready`, `Progressing`, `Degraded` |

### API group rename

The provisional `platform.example.io` was replaced with `platform.kubeforge.io` across all Go source, generated CRD/RBAC, kustomize config, and spec docs. The open question in `specs/OPEN-QUESTIONS.md` was resolved.

### Files changed

- `api/v1alpha1/application_types.go` — full spec/status types
- `api/v1alpha1/groupversion_info.go` — group name
- `api/v1alpha1/zz_generated.deepcopy.go` — regenerated
- `internal/controller/application_controller.go` — RBAC markers
- `cmd/main.go` — leader election ID
- `PROJECT` — domain + project name
- `config/samples/platform_v1alpha1_application.yaml` — sample CR
- `config/crd/bases/platform.kubeforge.io_applications.yaml` — regenerated
- `config/crd/kustomization.yaml` — CRD reference
- `config/rbac/role.yaml` — regenerated
- `config/rbac/application_*_role.yaml` — group references
- `specs/04-crd-api-spec.md` — final group name
- `specs/OPEN-QUESTIONS.md` — resolved entry

### Verification

- `go build ./...` — clean
- `go vet ./...` — clean
- `make generate` — deepcopy regenerated
- `make manifests` — CRD + RBAC regenerated