# Phase 2 — GitOps and CI/CD

## Task 10 — Argo CD example

### What was done

Added Argo CD `AppProject` and `Application` manifests under `deploy/argocd/` to enable GitOps workflow where Argo CD installs the platform operator and manages application manifests from Git.

### Files created

| File | Description |
|---|---|
| `deploy/argocd/project.yaml` | `AppProject` named `kubeforge` with source repos, destination clusters, and resource whitelists |
| `deploy/argocd/application.yaml` | `Application` named `kubeforge-platform-operator` syncing `charts/platform-operator` from GitHub to `kubeforge-system` namespace |

### AppProject configuration

- **Name**: `kubeforge`
- **Source repos**: `https://github.com/Yash-Singh/KubeForge`
- **Destinations**: 
  - `kubeforge-system` namespace (primary, for Operator deployment)
  - `argocd` namespace (for the Application itself)
- **ClusterResourceWhitelist**: CRDs, ClusterRole, ClusterRoleBinding
- **NamespaceResourceWhitelist**: Deployment, Service, ServiceAccount, Role, RoleBinding, ServiceMonitor

### Application configuration

- **Project**: `kubeforge`
- **Source**: GitHub repo `https://github.com/Yash-Singh/KubeForge`, path `charts/platform-operator`, Helm unstructured
- **Destination**: `https://kubernetes.default.svc`, namespace `kubeforge-system`
- **Sync policy**: automated with `prune: true`, `selfHeal: true`, `CreateNamespace=true`, `ServerSideApply=true`
- **revisionHistoryLimit**: 3

### GitOps workflow (per `specs/06-gitops-cicd.md`)

```text
Git desired state
    ↓
Argo CD
    ↓
Kubernetes Application CR
    ↓
Operator
    ↓
Child resources
```

1. User makes a Git change to the repository (e.g., updates `charts/platform-operator/values.yaml`)
2. CI runs `helm template` and `helm lint` on PR
3. On merge, Argo CD detects the change and syncs
4. Argo CD applies the Helm chart → Operator CRDs are installed → Operator reconciles the `Application` CR
5. Operator creates Deployments, Services, ConfigMaps, etc. as desired state

### Usage

```bash
kubectl apply -n argocd -f deploy/argocd/project.yaml
kubectl apply -n argocd -f deploy/argocd/application.yaml
# Or via `kubectl apply -k deploy/argocd/` if using kustomize