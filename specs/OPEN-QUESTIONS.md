# Open Questions — Do Not Guess

This file intentionally records choices that should not be silently invented by a coding agent.

## Project identity

- [ ] Final project name.
- [ ] Final public API group/domain for the CRD.
- [ ] License.

Until decided, use `platform.example.io` only in examples/templates and mark generated artifacts clearly as provisional.

## AI implementation

- [ ] Final AI service language/runtime.
- [ ] First provider adapter to implement.
- [ ] Local model runtime for development (Ollama or another runtime).
- [ ] Whether a unified provider library is used internally or providers are implemented directly.
- [ ] Supported MCP SDK/language version.

The architecture must remain provider-agnostic regardless of these choices.

## Deployment

- [ ] Exact Kubernetes version matrix to test.
- [ ] Minimum supported cluster profile.
- [ ] Whether operator watches one namespace by default or cluster-wide.

## Integrations

- [ ] First KEDA trigger type.
- [ ] Whether Argo Rollouts is included in the first reliability release.
- [ ] Whether ServiceMonitor is part of the core chart or an example.
- [ ] Whether Argo CD auto-prune/self-heal is enabled in the demo environment.

## Security

- [ ] Identity/authentication mechanism for remote MCP deployment.
- [ ] GitHub app/token model for future PR creation.
- [ ] Whether any direct write tool will ever be supported.

## How to resolve

When an implementation needs one of these answers:

1. Do not silently choose a production contract.
2. Use a temporary internal/default value only if it is explicitly marked provisional.
3. Add an ADR when the decision has architectural consequences.
4. Update this file and the affected specification once decided.
