# ADR 0008 — Keep Cloud APIs Outside the Core Contract

## Status
Accepted

## Decision

The core `Application` API and operator must remain cloud-neutral. AWS-specific behavior lives in optional adapters/integrations.

## Rationale

The project is intended to run on Kind and generic Kubernetes, while still being deployable to EKS and useful to platform engineers.
