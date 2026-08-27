# Agent Implementation Contract

This file is for coding agents working on the repository.

## Before changing code

1. Read `README.md` and the relevant document under `specs/`.
2. Inspect existing code and tests before creating new abstractions.
3. Check whether the requested behavior already exists in Kubernetes or an upstream dependency.
4. Do not invent API fields or runtime behavior not present in the specs. Use `specs/OPEN-QUESTIONS.md` for unresolved requirements.
5. Add/update an ADR when changing a major architectural decision.

## General engineering rules

- Prefer small, reviewable changes.
- Keep deterministic infrastructure logic separate from LLM logic.
- Kubernetes reconciliation must be idempotent and level-based, not event-scripted.
- Treat `.spec` as desired state and `.status` as observed state.
- Use owner references for resources directly managed by a parent custom resource when appropriate.
- Make retries safe.
- Never log secrets, tokens, full credentials, or unrestricted application logs.
- Do not store telemetry, logs, billing data, or large application datasets in Kubernetes CRs.
- Do not make an LLM call from the operator reconcile loop.
- Do not allow an LLM to directly mutate arbitrary Kubernetes resources.
- Every AI action must have an explicit tool contract, authorization boundary, audit record, timeout, and failure mode.

## AI/model abstraction rules

- The core project must not assume a particular LLM vendor or model.
- Model selection belongs to configuration, not business logic.
- Prompts/instructions must not contain vendor-specific API behavior.
- Provider adapters must expose a common interface for chat/agent inference and structured output where supported.
- Tool schemas are authoritative and versioned independently of prompts.
- The AI service must remain useful in read-only mode when no write-capable provider or credentials are configured.

## Testing rules

Every feature should include the lowest appropriate test level and, when it affects reconciliation, an integration test.

Required for operator changes:

- Go unit tests.
- envtest tests for API/reconciliation semantics where applicable.
- Kind E2E for behavior that depends on real Kubernetes components.

Required for AI changes:

- tool input/output contract tests.
- authorization tests.
- deterministic fixture-based evaluation cases.
- prompt/tool-call behavior tests where practical.

## Documentation rules

A feature is not complete until the relevant specification, examples, and operational behavior are documented.
