# Agents

Runtime AI agent definitions for KubeForge. These definitions are **model-agnostic** by design:

- The **Model** is configuration, not identity. Any provider that implements the `ModelProvider` interface described in `specs/08-ai-mcp-spec.md` (local or remote) can be bound to an agent without changing its tools, instructions, or contracts.
- **Tools** are bounded, versioned, and authoritative — defined independently of prompts.
- **Instructions/guardrails** are behavioral policy, versioned in source control.

Authoritative sources:

- `specs/08-ai-mcp-spec.md` — AI architecture, provider abstraction, tool contract, guardrails.
- `specs/design/agent-workflow.md` — planner/tool loop and evidence model.
- `specs/09-security.md` — authorization boundaries.
- `AGENTS.md` — implementation rules for coding agents (distinct from runtime agents).

Adding a new agent: add a definition file here, declare its tools explicitly, and include contract, authorization, and fixture-based evaluation tests per `specs/10-testing-strategy.md`.
