# Testing Strategy

## Test Pyramid

```text
         +----------------+
         |   Kind E2E     |
         +-------+--------+
                 |
         +-------+--------+
         | envtest / integ |
         +-------+--------+
                 |
         +-------+--------+
         |   unit tests   |
         +-----------------+
```

## Unit Tests

Test pure functions independently:

- Defaulting/calculation
- Desired Deployment generation
- Desired Service generation
- Condition calculation
- Error classification
- Resource naming/labels

## envtest

Use envtest for API/reconciler semantics without requiring a full workload cluster.

### Minimum Scenarios

1. Application creates Deployment
2. Application creates Service when requested
3. Application status becomes progressing/ready appropriately
4. Spec update changes child desired state
5. Child deletion is recovered
6. Invalid spec is rejected by API validation

## Kind E2E

Run a real local cluster with:

- Operator installed via Helm
- Example Application deployed
- All child resources verified
- Status and events verified

### E2E Test Cases

1. Deploy Application with all features enabled
2. Verify Deployment, Service, PDB, HPA, NetworkPolicy created
3. Verify status phase transitions
4. Update Application spec and verify reconciliation
5. Delete Application and verify cleanup
6. Verify metrics endpoint
7. Verify events are recorded

## Running Tests

```bash
# Unit tests
go test ./...

# Integration tests (requires envtest)
make test

# E2E tests (requires Kind cluster)
make test-e2e
```

## Test Coverage

Aim for:

- 80%+ unit test coverage
- All reconciliation paths covered
- All error scenarios covered
- All CRD validation paths covered
