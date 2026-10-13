# KubeForge

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.28+-326CE5?style=flat&logo=kubernetes&logoColor=white)](https://kubernetes.io/)
[![Helm](https://img.shields.io/badge/Helm-3.14+-0F1689?style=flat&logo=helm&logoColor=white)](https://helm.sh/)
[![ArgoCD](https://img.shields.io/badge/ArgoCD-2.x-EF7B4D?style=flat&logo=argo&logoColor=white)](https://argo-cd.readthedocs.io/)
[![CI](https://github.com/Yash-Singh/KubeForge/actions/workflows/test.yml/badge.svg)](https://github.com/Yash-Singh/KubeForge/actions/workflows/test.yml)

**A Kubernetes-native application platform that simplifies deployments through a single Custom Resource.**

KubeForge lets developers deploy applications to Kubernetes by declaring intent in one YAML file. The operator handles the rest — creating Deployment, Service, HPA, PDB, NetworkPolicy, KEDA ScaledObjects, and Argo Rollouts automatically.

## Features

- **Single CRD Deployment** — Declare your application in one `Application` CR instead of managing 5+ separate Kubernetes resources
- **Automatic Resource Creation** — Operator creates and manages Deployment, Service, PDB, HPA, NetworkPolicy, KEDA, and Argo Rollouts
- **GitOps Ready** — Helm charts, ArgoCD ApplicationSets for dev/prod environments
- **Production Observability** — Built-in Prometheus metrics, structured JSON logging, Kubernetes events, and Grafana dashboards
- **Reliability Features** — PodDisruptionBudgets, topology spread constraints, network policies, and KEDA scale-to-zero
- **Canary/Blue-Green Deployments** — Native Argo Rollouts integration for safe rollouts
- **Status Tracking** — Real-time phase reporting (Ready/Progressing/Degraded) with conditions

## Quick Start

### Prerequisites

- Kubernetes 1.28+
- Helm 3.14+
- kubectl

### Installation

```bash
# Add the Helm repository (or use local chart)
helm install kubeforge ./charts/platform-operator \
  --namespace kubeforge-system \
  --create-namespace
```

### Deploy Your First Application

```yaml
apiVersion: platform.kubeforge.io/v1alpha1
kind: Application
metadata:
  name: my-app
  namespace: default
spec:
  image: nginx:alpine
  replicas: 3
  service:
    port: 80
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 200m
      memory: 256Mi
```

```bash
kubectl apply -f my-app.yaml
kubectl get applications.platform.kubeforge.io
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Kubernetes Cluster                      │
│                                                             │
│  ┌──────────────┐     ┌──────────────┐    ┌──────────────┐ │
│  │  Application  │     │   Operator   │    │  Prometheus  │ │
│  │      CR       │────▶│   Reconciler │───▶│   Metrics    │ │
│  │  (Desired)    │     │              │    │   (:8443)    │ │
│  └──────────────┘     └──────┬───────┘    └──────────────┘ │
│                              │                              │
│              ┌───────────────┼───────────────┐              │
│              ▼               ▼               ▼              │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐       │
│  │  Deployment  │ │   Service    │ │     PDB      │       │
│  │  (Owned)     │ │  (Owned)     │ │   (Owned)    │       │
│  └──────────────┘ └──────────────┘ └──────────────┘       │
│                                                             │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐       │
│  │ NetworkPolicy│ │     HPA      │ │ KEDA ScaledO │       │
│  │  (Owned)     │ │  (Owned)     │ │   (Owned)    │       │
│  └──────────────┘ └──────────────┘ └──────────────┘       │
└─────────────────────────────────────────────────────────────┘
```

## CRD Specification

The `Application` CRD supports:

| Field | Description |
|-------|-------------|
| `spec.image` | Container image reference (required) |
| `spec.replicas` | Number of replicas (default: 1) |
| `spec.service` | Service configuration |
| `spec.resources` | CPU/memory requests and limits |
| `spec.env` | Environment variables |
| `spec.probes` | Liveness/readiness probes |
| `spec.podDisruptionBudget` | PDB configuration |
| `spec.topologySpreadConstraints` | Pod placement constraints |
| `spec.networkPolicy` | Network policy rules |
| `spec.horizontalPodAutoscaler` | HPA configuration |
| `spec.kedaScaledObject` | KEDA scale-to-zero config |
| `spec.argoRollout` | Argo Rollout configuration |

## Observability

### Metrics

The operator exposes custom Prometheus metrics on `:8443/metrics`:

- `kubeforge_reconcile_total` — Total reconciliations
- `kubeforge_reconcile_duration_seconds` — Reconciliation duration
- `kubeforge_application_phase` — Current application phase
- `kubeforge_application_replicas` — Replica counts

### Grafana Dashboard

Pre-built dashboards are available in `charts/platform-operator/dashboards/`.

### Kubernetes Events

Phase transitions are recorded as Kubernetes events:

```bash
kubectl get events --field-selector involvedObject.name=my-app
```

## Development

### Prerequisites

- Go 1.22+
- Docker
- Kind (for local development)
- kubectl

### Local Development

```bash
# Create Kind cluster
kind create cluster --name kubeforge

# Install CRD
kubectl apply -f config/crd/bases/platform.kubeforge.io_applications.yaml

# Run operator locally
make run

# Run tests
make test
```

### Building

```bash
# Build binary
make build

# Build Docker image
make docker-build IMG=kubeforge:latest

# Load into Kind
kind load docker-image kubeforge:latest --name kubeforge
```

## Testing

```bash
# Unit tests
go test ./...

# Integration tests (requires envtest)
make test

# E2E tests (requires Kind cluster)
make test-e2e
```

## Deployment

### Helm

```bash
# Install
helm install kubeforge ./charts/platform-operator \
  --namespace kubeforge-system \
  --create-namespace

# Upgrade
helm upgrade kubeforge ./charts/platform-operator \
  --namespace kubeforge-system
```

### ArgoCD

```bash
# Apply ApplicationSet
kubectl apply -f deploy/argocd/applicationset.yaml
```

## Project Structure

```
├── api/v1alpha1/          # CRD type definitions
├── internal/controller/   # Reconciler logic
├── cmd/                   # Entry point
├── charts/                # Helm charts
├── config/                # Kubernetes manifests
├── deploy/                # ArgoCD and Kustomize configs
├── test/                  # E2E tests
└── docs/                  # Development logs
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.

## Acknowledgments

- Built with [Kubebuilder](https://book.kubebuilder.io/)
- Uses [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)
- GitOps powered by [ArgoCD](https://argo-cd.readthedocs.io/)
- Autoscaling with [KEDA](https://keda.sh/)
- Progressive delivery with [Argo Rollouts](https://argoproj.github.io/rollouts/)
