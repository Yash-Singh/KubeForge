# ADR 0003 — Model-Provider Agnostic AI Layer

## Status
Accepted

## Context

The project should demonstrate real LLM usage without becoming tied to a single vendor or model.

## Decision

Define an internal provider abstraction and put model/provider selection in configuration. Business logic, MCP tools, API schemas, and operational workflows must not depend on vendor-specific SDK behavior.

The initial implementation should support a local development model endpoint and at least one remote provider through adapters. Additional providers must be addable without changing the platform/operator APIs.

## Consequences

Positive:

- local development without mandatory paid API access;
- easier experimentation with different models;
- cleaner open-source contribution boundary.

Negative:

- provider capability differences must be handled explicitly;
- some features such as tool calling/structured output are not uniformly available.
