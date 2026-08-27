# ADR 0002 — Keep AI Out of the Kubernetes Reconcile Loop

## Status
Accepted

## Decision

The operator must never call an LLM as part of its reconciliation path.

AI runs in a separate service and can inspect the platform through bounded tools.

## Rationale

Kubernetes reconciliation must be predictable, idempotent, retry-safe, and capable of converging despite repeated events. LLM inference introduces non-determinism, external dependencies, variable latency/cost, and security/data-boundary concerns.

## Consequences

AI can evolve independently from the operator. The operator remains usable when the AI service or model provider is down.
