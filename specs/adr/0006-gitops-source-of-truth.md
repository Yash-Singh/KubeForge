# ADR 0006 — GitOps Is the Deployment Source of Truth

## Status
Accepted

## Decision

Git stores desired installation/application configuration. Argo CD synchronizes Git to Kubernetes. The operator reconciles the `Application` CR to child resources.

## Rationale

This separates two reconciliation loops cleanly:

```text
Git -> Argo CD -> Kubernetes desired resources
Kubernetes Application -> Operator -> child resources
```

## Consequences

Manual changes are useful for debugging but are not the supported persistent workflow. Persistent changes should normally be represented in Git.
