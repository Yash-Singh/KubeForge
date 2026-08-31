# Phase 2 — GitOps and CI/CD

## Task 11 — Argo CD project and sync fixes

### What was done

Fixed multiple Argo CD configuration issues that prevented the platform operator from syncing:

1. **Argo CD Application source repo URL** - Pointed to non-existent `kubeforge/kube-forge` instead of actual repo `Yash-Singh/KubeForge`
2. **AppProject sourceRepos** - Same issue, prevented Application from using the repo
3. **Namespace not in clusterResourceWhitelist** - `Namespace` is cluster-scoped, must be in `clusterResourceWhitelist` not `namespaceResourceWhitelist`
4. **ServiceMonitor CRD missing** - `monitoring.coreos.com/v1` API group existed but `ServiceMonitor` CRD not installed
5. **Wrong GHCR image repository** - Helm chart used `ghcr.io/Yash-Singh/KubeForge` instead of correct `ghcr.io/yash-singh/kube-forge` (GHCR normalizes to lowercase)
6. **GHCR package private** - Released images were private (401/403), added `imagePullSecrets` and workflow step to make package public
7. **Image tag mismatch** - Workflow pushes tag matching Git tag (v0.1.2), values.yaml updated to match

### Files changed

| File | Change |
|---|---|
| `deploy/argocd/application.yaml` | `repoURL: https://github.com/Yash-Singh/KubeForge` |
| `deploy/argocd/project.yaml` | `sourceRepos: https://github.com/Yash-Singh/KubeForge`; moved `Namespace` to `clusterResourceWhitelist` |
| `charts/platform-operator/values.yaml` | `image.repository: ghcr.io/yash-singh/kube-forge`, `tag: v0.1.2`, added `imagePullSecrets` |
| `.github/workflows/release.yml` | Added step to make GHCR package public after push |

### Root causes & solutions

| Error | Root Cause | Fix |
|---|---|---|
| `Repository not found` | Argo CD Application `repoURL` pointed to `kubeforge/kube-forge` (doesn't exist) | Updated to `Yash-Singh/KubeForge` |
| `application repo not permitted in project` | AppProject `sourceRepos` didn't include actual repo | Added `Yash-Singh/KubeForge` to `sourceRepos` |
| `resource :Namespace is not permitted in project kubeforge` | `Namespace` is cluster-scoped but was in `namespaceResourceWhitelist` | Moved to `clusterResourceWhitelist` |
| `failed to discover server resources for group version monitoring.coreos.com/v1` | `ServiceMonitor` CRD not installed on cluster | Applied `servicemonitors.monitoring.coreos.com` CRD |
| `ImagePullBackOff: 403 Forbidden` | GHCR package private, wrong repo case | Added `imagePullSecrets`, fixed repo to lowercase `ghcr.io/yash-singh/kube-forge` |
| `ImagePullBackOff: not found` | Workflow pushed v0.1.2 but values.yaml had v0.1.0 | Updated tag to v0.1.2 matching Git tag |

### Verification

```bash
# Check app status
kubectl -n argocd get application kubeforge-platform-operator -o yaml

# Expected: status.sync.status: Synced, all resources status: Synced
# All 14 resources show status: Synced

# Check operator pod
kubectl -n kubeforge-system get pods
# Expected: 1/1 Running
```

### Current state

- ✅ Argo CD Application: **Synced**
- ✅ All 14 resources: **Synced**
- ✅ Platform Operator pod: **Running** (1/1)
- ✅ Leader election: **Acquired lease**

### How to release

```bash
# Tag and push - workflow builds, pushes to GHCR, makes package public
git tag v0.1.3 && git push origin v0.1.3

# Verify image pulled
kubectl -n kubeforge-system get pods
```

### For local Kind development

```bash
# Build and load image locally (bypasses GHCR auth)
docker build -t ghcr.io/yash-singh/kube-forge:v0.1.2 .
kind load docker-image ghcr.io/yash-singh/kube-forge:v0.1.2

# Restart deployment to pick up local image
kubectl -n kubeforge-system rollout restart deployment/kubeforge-platform-operator
```