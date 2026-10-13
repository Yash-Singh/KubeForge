# System Architecture

## Architectural Principles

1. Kubernetes remains the source of runtime truth
2. Git is the desired-state source of truth for deployed configuration
3. The operator is deterministic and idempotent
4. Integrate upstream systems rather than reimplementing them
5. Local development must not require AWS
6. Optional infrastructure integrations must not contaminate the core API

## Logical Architecture

```text
                                 +-------------------+
                                 |      GitHub       |
                                 | source + manifests|
                                 +---------+---------+
                                           |
                              GitHub Actions / CI
                                           |
                                           v
                                     +-----+-----+
                                     |    GHCR   |
                                     +-----------+
                                           |
                                  Helm / Kustomize
                                           |
                                           v
                                     +-----+-----+
                                     |   Argo CD  |
                                     +-----+-----+
                                           |
                                           v
                              +------------+-------------+
                              |     Kubernetes cluster   |
                              |                          |
                              |  +--------------------+  |
                              |  | Application CRD    |  |
                              |  +---------+----------+  |
                              |            |             |
                              |            v             |
                              |  +---------+----------+  |
                              |  | Platform Operator  |  |
                              |  |       Go            |  |
                              |  +---+-------+------+--+  |
                              |      |       |      |     |
                              |      v       v      v     |
                              |  +---+---+ +---+ +-------+|
                              |  | Deploy| |Svc | |PDB/HPA||
                              |  +-------+ +---+ +-------+|
                              +---------------------------+
```

## Components

### Application CRD

- Namespaced custom resource
- Defines application intent (image, replicas, service, resources, etc.)
- Short name: `app`

### Platform Operator

- Go-based controller using controller-runtime
- Reconciles Application CRs to child resources
- Creates and manages: Deployment, Service, PDB, HPA, NetworkPolicy, KEDA, Argo Rollouts
- Reports status through conditions and phase tracking
- Exposes Prometheus metrics on :8443

### Child Resources

| Resource | Purpose |
|----------|---------|
| Deployment | Pod management and rolling updates |
| Service | Network exposure |
| PodDisruptionBudget | High availability during disruptions |
| HorizontalPodAutoscaler | CPU/memory-based scaling |
| NetworkPolicy | Traffic control |
| KEDA ScaledObject | Event-driven scaling |
| Argo Rollout | Canary/blue-green deployments |

### GitOps Pipeline

```text
GitHub → GitHub Actions → GHCR → Helm/Kustomize → Argo CD → Kubernetes
```

### Observability Stack

```text
Operator → Prometheus metrics (:8443)
         → Kubernetes events
         → Structured JSON logs
         → Grafana dashboards
```

## Data Flow

### Deployment Flow

1. Developer pushes Application YAML to Git
2. Argo CD syncs the manifest to Kubernetes
3. Operator detects new/changed Application
4. Operator creates/updates child resources
5. Operator updates status and records events
6. Prometheus scrapes metrics

### Reconciliation Flow

```text
Application CR
     ↓
Operator Reconcile
     ↓
+----+----+----+----+----+
|    |    |    |    |    |
v    v    v    v    v    v
Deploy Svc PDB HPA NetPol KEDA
```

## Design Decisions

See `specs/adr/` for architecture decision records.
