# Product Specification

## Working name

**AI-Assisted Kubernetes Application Platform**

The public project name is intentionally not fixed yet. See `OPEN-QUESTIONS.md`.

## Problem

Kubernetes platform teams often expose a growing collection of low-level objects and tools to developers: Deployments, Services, autoscaling resources, rollout controllers, policies, observability configuration, and GitOps workflows. Engineers can automate much of the deterministic lifecycle, but diagnosing failures still requires stitching together Kubernetes state, events, metrics, logs, deployment history, and Git changes manually.

This project aims to provide a small, open, Kubernetes-native application abstraction with an operator for deterministic lifecycle management, GitOps-based delivery, standard observability, and an AI layer that can inspect operational evidence and produce grounded recommendations.

## Target users

### Primary

- Platform engineers.
- DevOps/SRE engineers.
- Kubernetes engineers.
- Infrastructure/software engineers working on internal platforms.

### Secondary

- Developers who want a simple declarative application interface.
- Open-source contributors interested in Kubernetes controllers and AI infrastructure.

## Product goals

1. Provide a meaningful Kubernetes CRD that represents application intent.
2. Implement a robust Go controller using standard Kubernetes reconciliation patterns.
3. Package and release the operator as a deployable Helm chart.
4. Demonstrate GitOps deployment through Argo CD.
5. Provide reproducible CI/CD with tests, linting, container publishing, and releases.
6. Integrate optional Kubernetes-native scaling/reliability components rather than reinventing them.
7. Provide Prometheus/OpenTelemetry/Grafana observability.
8. Provide an AI service that calls an LLM but is model-provider agnostic.
9. Expose selected platform capabilities through MCP.
10. Keep AI actions safety-bounded and human-controlled by default.

## Product non-goals

- Replacing Kubernetes.
- Replacing Argo CD, KEDA, HPA, Argo Rollouts, Prometheus, or OpenTelemetry.
- Building a general-purpose cloud cost product in the core project.
- Rebuilding Backstage in the MVP.
- Operating Kafka/Elasticsearch as core dependencies.
- Building a proprietary LLM or training a foundation model.
- Allowing autonomous destructive changes in production by default.

## User journey

### Application lifecycle

1. Developer commits an `Application` custom resource to Git.
2. GitHub Actions validates and tests the repository.
3. Argo CD synchronizes the desired state to a cluster.
4. The operator reconciles the `Application` into owned Kubernetes resources.
5. Platform telemetry is exposed to Prometheus/OTel and visualized in Grafana.
6. An AI client/agent can inspect the application through MCP tools.
7. AI returns a diagnosis/recommendation with evidence.
8. Any write action, when later enabled, goes through explicit authorization and preferably Git PR/GitOps rather than direct mutation.

## Success criteria

The first usable release is successful when a contributor can:

- start a local cluster;
- install the operator with Helm;
- apply an example `Application` CR;
- observe Deployment/Service creation;
- modify the CR and see reconciliation;
- delete a child resource and observe recovery;
- inspect status/conditions;
- run automated unit/envtest/E2E checks;
- install via Argo CD from Git;
- see operator metrics and basic traces;
- query a local or remote model through the AI service;
- connect a supported MCP client and perform safe, read-only diagnosis.

## Quality attributes

- Deterministic operator behavior.
- Model/provider agnosticism at the AI service boundary.
- Local-first development.
- Reproducible builds.
- Least privilege.
- Observable by default.
- Safe failure over clever automation.
- Backward-compatible API evolution.
- Clear separation between spec, implementation, and optional integrations.
