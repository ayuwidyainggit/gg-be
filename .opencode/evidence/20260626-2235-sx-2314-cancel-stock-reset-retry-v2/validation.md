# Validation — SX-2314 Cancel Stock Reset Retry v2 (Post-Remediation)

Date: 2026-06-26
Service: `sales`

## Focused Tests
```
rtk go test ./repository/... ./service/... -run "TestReconcileCancel|TestBuildCancelReconcile|TestCancelAudit|TestGetCancelStockBasisQuery_IncludesReward|TestBulkUpdateStatus_Cancel" -count=1 -v
```
Result: 18 passed in 2 packages

## Full Suites
```
rtk go test ./repository/... ./service/... -count=1
```
Result: 279 passed in 2 packages (baseline 267 + 12 new)

## Full Build
```
rtk go build ./...
```
Result: Success

## Remediation-specific validations
- `CancelStockBasis` struct carries `ItemType int64` field
- `cancelStockBasisQuery` projects `od.item_type AS item_type`
- `cancelAuditValues` helper splits order/reward sums
- Structured `[CANCEL]` audit log present in `BulkUpdateStatus` cancel branch
- `buildCancelReconcileMutations` pure helper asserted for residual math, dedup, negative-clamp, reward sign, canonical dedup
- Weak panic-wait test replaced by pure helper assertions
