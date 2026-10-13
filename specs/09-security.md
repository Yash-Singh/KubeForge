# Security Specification

## Security Model

The project is a Kubernetes platform control plane. Security focuses on RBAC, supply chain, and runtime safety.

## Kubernetes RBAC

Operator permissions are generated from the resources it actually manages. The operator does not receive cluster-admin.

### Operator Permissions

- Applications (CRD): full CRUD
- Deployments, Services, PDBs, NetworkPolicies, HPAs: full CRUD
- Events: create/patch (for recording phase transitions)
- Secrets: not accessed by the operator

### ArgoCD Permissions

ArgoCD manages the operator deployment through standard Helm/Kustomize patterns.

## Supply Chain Security

### Image Security

- Multi-arch images (amd64+arm64) built in CI
- Images published to GHCR with immutable tags
- Dockerfile uses distroless base images
- No secrets or credentials in images

### Code Security

- GolangCI-Lint with security-focused linters
- Dependency scanning via Dependabot/Renovate
- No hardcoded credentials or secrets

## Network Security

### Network Policies

- Applications can opt-in to NetworkPolicy configuration through the CRD
- Default deny-all with explicit allow rules
- Ingress and egress can be configured separately

### Service Exposure

- Operator metrics endpoint can be secured with TLS
- ServiceMonitor for Prometheus scraping (when enabled)

## Runtime Security

### Container Security

- Security context enforced on operator pods
- Non-root execution where possible
- Read-only root filesystem where supported

### Resource Limits

- CPU and memory limits on operator pods
- Resource requests for guaranteed scheduling

## RBAC Best Practices

1. **Least privilege**: Operator only gets permissions for resources it manages
2. **No cluster-admin**: Uses namespace-scoped roles where possible
3. **Separate service accounts**: Operator and ArgoCD use different SAs
4. **Audit logging**: Kubernetes events provide audit trail

## Future Considerations

If AI/MCP is implemented in the future:

- AI service should have separate service account
- Read-only AI mode should use read-only RBAC
- Never pass Kubernetes Secrets to LLM
- Treat pod logs, events, and annotations as untrusted input
