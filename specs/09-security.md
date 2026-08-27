# Security Specification

## Security model

The project is a platform control plane plus an AI tool surface. Treat the AI layer as an untrusted decision-maker with potentially privileged tools.

## Kubernetes RBAC

Operator permissions must be generated from the resources it actually manages. Start with the narrowest required verbs.

The operator should not receive cluster-admin.

The MCP/AI service should have a separate service account from the operator.

Read-only AI mode should use read-only RBAC where possible.

## AI action boundary

Default capability tiers:

```text
Tier 0 — Public/help text
Tier 1 — Read-only platform inspection
Tier 2 — Proposal generation (no cluster mutation)
Tier 3 — Human-approved Git changes
Tier 4 — Narrow direct mutation (future, exceptional)
Tier 5 — Destructive/autonomous mutation (not planned by default)
```

The MVP stops at Tier 1/2.

## Prompt injection

Potential untrusted inputs:

- pod logs;
- Kubernetes events;
- ConfigMap content;
- annotations/labels;
- Git commit messages and diffs;
- application-generated HTTP payloads;
- user input.

The agent must treat these as untrusted data, not instructions.

## Secrets

- Never pass Kubernetes Secrets to the LLM unless a specific future feature has a reviewed requirement.
- Never expose secret data through MCP tools.
- Redact common credential patterns from logs and traces.
- Credentials should be injected through environment/secret references, not committed.

## MCP authorization

For remote/HTTP MCP deployments, implement an explicit authentication and authorization model suitable for the deployment environment. Follow the current MCP authorization specification and platform identity model.

For local stdio use, rely on the local host's process/environment security and do not assume the network authorization flow applies.

## Tool security

Each tool must:

- validate input;
- enforce namespace/resource boundaries;
- limit result size;
- timeout;
- return structured errors;
- log/audit who invoked it and what target it accessed;
- avoid executing arbitrary commands.

## Supply chain

CI should include, as practical:

- dependency scanning;
- container image scanning;
- SBOM generation;
- provenance/signing before publishing mature releases.

Do not make signing tooling a blocker for the first local MVP if it creates significant implementation friction; add it before declaring a production-grade release.

## Network boundaries

The AI service may need to reach:

- Kubernetes API;
- Prometheus;
- Argo CD API;
- GitHub API;
- model provider endpoint.

NetworkPolicy should default toward minimal egress/ingress once the deployment topology is stable.
