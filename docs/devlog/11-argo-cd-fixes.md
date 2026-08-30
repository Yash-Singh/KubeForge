# Phase 2 — GitOps and CI/CD

## Task 11 — Argo CD project and sync fixes

### What was done

Fixed multiple Argo CD configuration issues that prevented the platform operator from syncing:

1. **Argo CD Application source repo URL** - Pointed to non-existent `kubeforge/kube-forge` instead of actual repo `Yash-Singh/KubeForge`
2. **AppProject sourceRepos** - Same issue, prevented Application from using the repo
3. **Namespace not in clusterResourceWhitelist** - `Namespace` is cluster-scoped, must be in `clusterResourceWhitelist` not `namespaceResourceWhitelist`
4. **ServiceMonitor CRD missing** - `monitoring.coreos.com/v1` API group existed but `ServiceMonitor` CRD not installed
5. **Wrong GHCR image repository** - Helm chart used `ghcr.io/Yash-Singh/KubeForge` instead of correct `ghcr.io/kubeforge/kube-forge`

### Files changed

| File | Change |
|---|---|
| `deploy/argocd/application.yaml` | `repoURL: https://github.com/Yash-Singh/KubeForge` |
| `deploy/argocd/project.yaml` | `sourceRepos: https://github.com/Yash-Singh/KubeForge`; moved `Namespace` to `clusterResourceWhitelist` |
| `charts/platform-operator/values.yaml` | `image.repository: ghcr.io/kubeforge/kube-forge` |

### Root causes & solutions

| Error | Root Cause | Fix |
|---|---|---|
| `Repository not found` | Argo CD Application `repoURL` pointed to `kubeforge/kube-forge` (doesn't exist) | Updated to `Yash-Singh/KubeForge` |
| `application repo not permitted in project` | AppProject `sourceRepos` didn't include actual repo | Added `Yash-Singh/KubeForge` to `sourceRepos` |
| `resource :Namespace is not permitted in project kubeforge` | `Namespace` is cluster-scoped but was in `namespaceResourceWhitelist` | Moved to `clusterResourceWhitelist` |
| `failed to discover server resources for group version monitoring.coreos.com/v1` | `ServiceMonitor` CRD not installed on cluster | Applied `servicemonitors.monitoring.coreos.com` CRD |
| `ImagePullBackOff: 403 Forbidden` | Helm chart used wrong GHCR repo (`Yash-Singh/KubeForge`) | Changed to correct `ghcr.io/kubeforge/kube-forge` |

### Verification

```bash
# Check app status
kubectl -n argocd get application kubeforge-platform-operator -o yaml

# Expected: status.sync.status: Synced, all resources status: Synced
# All 14 resources show status: Synced
```

### Current state

- ✅ Argo CD Application: **Synced**
- ✅ All 14 resources: **Synced**
- ⏳ Pod: `ImagePullBackOff` — waiting for image `ghcr.io/kubeforge/kube-forge:v0.1.0` to be built/pushed via release workflow

### Next steps

Trigger release workflow to build and push `ghcr.io/kubeforge/kube-forge:v0.1.0`:

```bash
git tag v0.1.0 && git push origin v0.1.0
```