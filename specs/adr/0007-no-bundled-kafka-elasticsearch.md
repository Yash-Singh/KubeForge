# ADR 0007 — Do Not Require Kafka/Elasticsearch for the Core Project

## Status
Accepted

## Decision

Kafka and Elasticsearch are optional advanced integrations. They are not required dependencies of the operator, AI service, or MVP observability stack.

## Rationale

They add substantial operational complexity without improving the core Kubernetes/AI control-plane problem enough to justify making them mandatory.

## Future use

A later reference architecture may demonstrate high-volume log/event pipelines using OTel/Fluent Bit -> Kafka -> Elasticsearch/OpenSearch or another suitable backend.
