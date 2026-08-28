# Phase 1 — Operator MVP

## Task 1 — Scaffold

### What was done

Created the Kubebuilder/controller-runtime project for the `platform.kubeforge.io` API group and `Application` kind.

### Files changed / created

| Item | Path |
|---|---|
| Go module | `go.mod` (module `github.com/kubeforge/kube-forge`) |
| Manager entrypoint | `cmd/main.go` — metrics (secure), health probes, leader election, webhook server, controller registration |
| API types (stub) | `api/v1alpha1/application_types.go` — Application, ApplicationSpec, ApplicationStatus, ApplicationList |
| Group/version registration | `api/v1alpha1/groupversion_info.go` |
| Controller (stub) | `internal/controller/application_controller.go` — `Reconcile()` placeholder + `SetupWithManager()` watching Application |
| Envtest suite | `internal/controller/suite_test.go` — test framework for envtest |
| CRD manifest | `config/crd/bases/platform.kubeforge.io_applications.yaml` |
| RBAC (generated) | `config/rbac/role.yaml` — ClusterRole with CRUD on applications, applications/status, applications/finalizers |
| RBAC roles | `config/rbac/application_admin_role.yaml`, `application_editor_role.yaml`, `application_viewer_role.yaml` |
| Sample CR | `config/samples/platform_v1alpha1_application.yaml` |
| Manager deployment | `config/manager/manager.yaml` |
| Kustomize overlays | `config/default/` — metrics, cert-manager patches |
| Prometheus config | `config/prometheus/monitor.yaml` |
| Network policy | `config/network-policy/allow-metrics-traffic.yaml` |
| Dockerfile | `Dockerfile` |
| Makefile | `Makefile` — build, test, generate, deploy targets |
| CI workflows | `.github/workflows/lint.yml`, `test.yml`, `test-e2e.yml` |
| Linter config | `.golangci.yml` |
| Devcontainer | `.devcontainer/devcontainer.json` |
| E2E test scaffold | `test/e2e/e2e_suite_test.go`, `e2e_test.go` |

### Verification

- `go build ./...` — clean
- `go vet ./...` — clean
- `make manifests` — CRD and RBAC generated

### Relevant specs

- `specs/03-system-architecture.md` — component boundaries, deployment topology
- `specs/04-crd-api-spec.md` — kind, scope, version, API conventions
- `specs/05-operator-design.md` — controller model, ownership, resource watches