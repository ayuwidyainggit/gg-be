# Execution — SX-2314 cancel stock reset

Date: 2026-06-24
Task id: `20260624-0941-sx-2314-cancel-stock-reset`
Executor: `@fixer`

## Repo-local evidence used
- `.opencode/plans/20260624-0941-sx-2314-cancel-stock-reset.md`
- `.opencode/evidence/20260624-0941-sx-2314-cancel-stock-reset/discovery.md`
- `sales/repository/stock_repository.go`
- `sales/repository/stock_repository_cancel_test.go`
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`

## Notes
- Repo stack docs requested by backend lane were absent:
  - `.opencode/docs/PROJECT_STACK.md`
  - `.opencode/docs/PROJECT_COMMANDS.md`
  - `.opencode/docs/FRAMEWORK_PLAYBOOK.md`
  - `.opencode/docs/PROJECT_DETECTED_TOOLS.md`
- Used repo-local plan, AGENTS guidance, source, tests.

## TDD / validation trail
- Test-first expectation updates applied in `sales/repository/stock_repository_cancel_test.go`.
- Initial repo tests after expectation update already matched new behavior.
- Service test exposed validator gap for Need Review + missing-source + positive fallback qty.
- Patched `validateCancelStockBasis` to allow missing source only when `QtyOutSmallest > 0`.

## Commands run
1. `rtk go test ./repository/... -run TestBuildCancelStockMutations -count=1 -v`
   - pass
2. `rtk go test ./repository/... -run TestGetCancelStockBasisQuery -count=1 -v`
   - pass
3. `rtk go test ./repository/... -run "TestBuildCancelStockMutations|TestGetCancelStockBasisQuery|TestCancelStockBasisFallback" -count=1 -v`
   - pass
4. `rtk go test ./service/... -run TestBulkUpdateStatus_Cancel -count=1 -v`
   - fail first, then pass after service validator patch
5. `rtk go test ./repository/... ./service/... -count=1`
   - pass (`267 passed in 2 packages`)

## Behavior changed
- Cancel reversal now writes single `CO` stock row per basis row.
- Duplicate cancel-source `SO` row removed.
- Cancel idempotency query now counts prior `CO` reversal rows.
- Cancel basis SQL now includes smallest-unit detail fallback using `qty*_final -> qty* -> qty_po*` with `conv_unit2/conv_unit3` and `mst.m_product` fallback.
- Need Review cancel now allowed when source ledger missing but positive outstanding fallback basis exists.

## Scope check
- Source/test edits kept to diff boundary files only.
- No migrations.
- No FE/API contract change.
- No env or secret changes.
