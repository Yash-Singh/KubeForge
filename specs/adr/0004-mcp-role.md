# ADR 0004 — MCP Is the AI Interoperability Layer

## Status
Accepted

## Decision

Implement MCP as a separate interface for exposing selected platform capabilities to AI hosts/agents. MCP does not replace the Kubernetes operator and is not itself the LLM.

Start with read-only tools. Add write tools only after explicit authorization and GitOps-first workflows are established.

## Rationale

MCP standardizes how an AI host discovers and invokes tools/resources. It gives the project an open integration boundary while keeping platform business logic in normal services/controllers.

## Consequences

The project can be used from multiple MCP-capable AI clients without changing the operator. The MCP surface must be carefully secured because tools can expose sensitive cluster data or actions.
