# GitOps and CI/CD Specification

## Repository model

A single public repository may initially contain:

```text
/
├── api/
├── cmd/
├── internal/
├── controllers/
├── config/
├── charts/
│   └── platform-operator/
├── examples/
├── deploy/
│   ├── argocd/
│   └── kustomize/
├── test/
│   └── e2e/
├── specs/
└── .github/workflows/
```

Split repositories only if operational needs justify it.

## CI on pull requests

Required checks:

- formatting;
- static analysis/lint;
- unit tests;
- generated code/manifest verification;
- CRD schema validation;
- envtest;
- Helm lint/template checks;
- container build;
- security/dependency scanning;
- optional Kind E2E for main branches or gated PRs depending on runtime cost.

## Image strategy

Default public registry: **GHCR**.

Images must be tagged with immutable release versions. `latest` may exist as a convenience tag but must never be the deployment source for reproducible releases.

## Release strategy

- Semantic versioning once the API is publicly consumable.
- Signed or provenance-attested releases when the tooling is stable enough to maintain.
- Release notes should identify API changes and compatibility implications.
- CRD changes require an explicit compatibility review.

## Helm

Helm is the supported installation package.

The chart should provide configuration for:

- image repository/tag;
- replica count for the operator itself;
- resources;
- service account/RBAC;
- metrics ServiceMonitor when the Prometheus Operator is present;
- optional PodMonitor/telemetry configuration if supported;
- leader election configuration;
- namespace-scoped vs cluster-scoped watch configuration if the implementation supports both.

Do not encode application `Application` resources into the operator chart by default.

## Kustomize

Kustomize may be used for GitOps overlays/examples. Do not duplicate Helm and Kustomize logic unnecessarily.

## Argo CD

Argo CD manages installation and application manifests from Git. The repository should contain an example `Application`/AppProject configuration that lets a user reproduce deployment.

The operator itself remains a normal Kubernetes Deployment installed by Argo CD.

## Git source of truth

For GitOps-managed resources:

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

Manual `kubectl` edits are allowed for debugging, but the expected documented workflow is Git change → CI → merge → Argo CD sync.

## Promotion model

Initially keep promotion simple. A single environment is enough for the project. Later environments may be represented with Kustomize/Helm values, but do not create dev/qa/staging/prod complexity before there is a real need.
