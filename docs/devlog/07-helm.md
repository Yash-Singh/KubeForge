# Phase 2 — GitOps and CI/CD

## Task 7 — Helm chart

### What was done

Created the `charts/platform-operator/` Helm chart with the following files:

| File | Description |
|---|---|
| `Chart.yaml` | Chart metadata (name, version, icon, maintainers) |
| `values.yaml` | Default configuration for replica count, image, RBAC, metrics, service account, resources, probes, leader election |
| `templates/_helpers.tpl` | Helper template definitions for chart name, fullname, labels, service account name |
| `templates/deployment.yaml` | Operator Deployment with leader election, liveness/readiness probes, resource requests/limits |
| `templates/serviceaccount.yaml` | ServiceAccount for the operator in `system` namespace |
| `templates/rbac.yaml` | ClusterRole + ClusterRoleBinding for manager-role; Role + RoleBinding for leader election; metrics-auth-role and metrics-reader ClusterRoles |
| `templates/service.yaml` | Metrics Service exposing port 8443 |
| `templates/servicemonitor.yaml` | ServiceMonitor (enabled when Prometheus Operator is present) |
| `crds/platform.kubeforge.io_applications.yaml` | Copied from `config/crd/bases/` for chart distribution |

### Configuration points

Per `specs/06-gitops-cicd.md`, the chart provides configuration for:

- **image** repository/tag, replica count, resources
- **service account** and **RBAC** (ClusterRole, RoleBinding, leader election)
- **metrics** Service + ServiceMonitor (optional, gated by `metrics.serviceMonitor.enabled`)
- **leader election** configuration (enabled via `--leader-elect` flag)
- **pod annotations** and **security context**

The chart does **not** encode Argo CD `Application` resources by default; those are provided separately in `deploy/argocd/`.

### Usage

```bash
helm install kubeforge-platform-operator charts/platform-operator
# or with custom values:
helm install kubeforge-platform-operator charts/platform-operator -f values.yaml
```

### Verification

- `helm lint charts/platform-operator` — passes
- `helm template platform-operator charts/platform-operator` — renders without errors
- All Go tests pass (16/16 envtest, unit tests)
- CRD is included in the chart for `kubectl apply -f charts/platform-operator/crds/`