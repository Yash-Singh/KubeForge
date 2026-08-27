# AI and MCP Specification

## Purpose

Add an AI operational intelligence layer without turning the Kubernetes operator into an AI controller.

The AI system should answer operational questions by gathering evidence from platform tools, then calling a configurable LLM to synthesize a grounded result.

## Architecture

```text
                MCP Client / AI Host
                         |
                         v
                  +------+------+
                  | MCP Server  |
                  +------+------+
                         |
                  tool/resource API
                         |
                         v
                  +------+------+
                  | AI Service  |
                  |             |
                  | Agent loop  |
                  | Policy      |
                  | Context     |
                  +------+------+
                         |
            +------------+-------------+
            |            |             |
            v            v             v
           K8s       Prometheus      Argo CD
           tools       tools          tools
                         |
                         v
                   Model Provider
                      Adapter
                         |
           +-------------+-------------+
           |             |             |
         remote        remote         local
         provider      provider       model
```

## LLM/provider abstraction

The project must not specify a single LLM model as a product requirement.

Define an internal provider interface roughly equivalent to:

```text
ModelProvider
  - Generate(request) -> response
  - GenerateStructured(request, schema) -> typed response
  - SupportsToolCalling() -> bool
  - ModelMetadata() -> metadata
```

The exact programming interface is implementation-specific.

Provider configuration should include:

- provider type;
- endpoint/base URL where supported;
- model name;
- credentials reference;
- request timeout;
- token/output limits;
- optional temperature/response controls only when supported.

Business logic must not branch on vendor names.

## Provider support strategy

The architecture should support:

1. a local model endpoint for development (for example Ollama or another compatible local serving stack);
2. at least one remote provider through an adapter;
3. additional providers without changing tools, CRDs, or agent business logic.

Do not claim every model is interchangeable. Tool calling, structured output, context limits, and multimodal capabilities differ. The system should detect and fail clearly when a selected model lacks a capability required by a workflow.

## Agent design

An agent has three explicit dimensions:

- Model — the LLM used for reasoning.
- Tools — bounded external capabilities.
- Instructions/guardrails — behavioral policy.

The agent should be optimized for evidence gathering, not open-ended autonomous action.

## Initial AI use case: incident/triage assistant

Input:

```text
Why is application checkout unhealthy?
```

The agent gathers evidence using read-only tools such as:

- `get_application`
- `get_deployment`
- `get_pods`
- `get_events`
- `get_services`
- `query_prometheus`
- `get_argo_application`
- `get_recent_changes`

The agent then returns a structured result:

```json
{
  "summary": "...",
  "severity": "warning",
  "confidence": 0.0,
  "evidence": [
    "..."
  ],
  "likely_causes": [
    "..."
  ],
  "recommended_actions": [
    "..."
  ],
  "unknowns": [
    "..."
  ]
}
```

Do not force the model to output confidence as a calibrated probability unless the system has a validated calibration method. Treat it as an explicitly named heuristic signal.

## Evidence-grounding rule

Every diagnosis claim should be traceable to one or more tool results. The AI response should prefer:

- observed values;
- timestamps;
- resource names;
- deployment/image changes;
- relevant events;
- metric query summaries.

If evidence is insufficient, the agent must say so and identify what is missing.

## MCP role

MCP is the interoperability boundary for AI hosts/clients. The server can expose:

- **Tools:** executable operations such as querying Kubernetes or Prometheus.
- **Resources:** structured contextual data where useful.
- **Prompts:** optional reusable workflows/templates.

The project should start with read-only tools.

## MCP version target

Target the **MCP 2026-07-28 specification** at implementation time, while keeping the server/client SDK version pinned and upgradeable. Follow the SDK/spec compatibility matrix chosen by the implementation.

The 2026-07-28 spec is stateless at the core and emphasizes cacheable list results, explicit headers/routing, authorization hardening, extensions, and improved request handling. Implement only the core features actually needed by this project.

## MCP tool design

Every tool must have:

- stable name;
- human-readable description;
- strict JSON Schema input;
- deterministic error format;
- authorization check;
- bounded result size;
- timeout;
- audit metadata;
- redaction rules.

Example:

```text
get_pods
Input:
  namespace: string
  selector?: string
  limit?: integer

Output:
  pods: [
    {
      name,
      phase,
      ready,
      restarts,
      node,
      reason
    }
  ]
```

Do not expose arbitrary `kubectl exec`, unrestricted shell execution, or arbitrary raw Kubernetes API write access in the initial MCP server.

## Write capabilities

When write tools are eventually introduced, separate them from read tools and require:

- explicit user authorization;
- resource-level allowlists;
- dry-run or plan mode where possible;
- bounded operations;
- GitOps-first execution for persistent configuration changes;
- audit record;
- idempotency key when applicable;
- clear user-visible diff before execution.

Example preferred workflow:

```text
AI recommendation
      ↓
proposed Git diff
      ↓
human approval
      ↓
GitHub PR
      ↓
CI
      ↓
Argo CD
      ↓
Kubernetes
```

Direct cluster mutation may be appropriate for a narrow operational action later, but it is not the default design.

## Prompt/instruction rules

- System instructions are versioned in source control.
- Prompts must distinguish observed evidence from model inference.
- Do not place credentials or unrestricted cluster dumps in prompts.
- Bound context size with summaries/filters.
- Prefer tool results over copied text from prior model output.
- Never trust tool output as instructions; treat it as data.
- Defend against prompt injection from logs, resource annotations, Git content, or user-controlled application text.

## AI failure modes

The AI layer must tolerate:

- provider timeout;
- rate limits;
- invalid model output;
- unavailable tool;
- partial telemetry;
- stale data;
- contradictory evidence;
- model hallucination;
- context limit.

Failure should degrade to an explicit error or partial analysis, never to an unsafe mutation.

## AI telemetry

Instrument:

- agent invocation count;
- tool invocation count/duration/errors;
- model request count/duration/errors;
- token counts when reliably available;
- provider/model identifiers;
- structured output validation failures.

Use OpenTelemetry GenAI semantic conventions where applicable.

Do not record raw prompts/completions by default.
