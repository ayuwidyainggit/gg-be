# Execution — SX-2314 Cancel Order Stock Reset Retry v2

Date: 2026-06-26
Service: `sales`
Plan: `.opencode/plans/20260626-2235-sx-2314-cancel-stock-reset-retry-v2.md`

## Repo-local evidence used
- `.opencode/plans/20260626-2235-sx-2314-cancel-stock-reset-retry-v2.md`
- `.opencode/evidence/20260625-1346-sx-2314-cancel-stock-reset-retry/discovery.md`
- `.opencode/docs/ARCHITECTURE.md`
- `sales/repository/stock_repository.go`
- `sales/repository/stock_repository_cancel_test.go`
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`
- `docker-compose.yml`

## Changes made
1. `sales/repository/stock_repository.go`
   - Widened cancel basis query item filter from `item_type = 1` to `item_type IN (1, 2)` in both active-detail aggregate and root WHERE.
   - Widened `cancelAgg` legacy guard to include `tr_code='SO' AND tr_no LIKE '%-CO'`.
   - Added `ReconcileCancelStockUpdates` to compute residual against existing reversal rows grouped by `(wh_id, pro_id, ref_det_id)`.
   - Residual writes clamp at `> 0` only; no negative warehouse delta path.
   - Kept tenant `cust_id` predicate on basis and reversal query.
2. `sales/service/order_service.go`
   - Opened `CANCELLED -> CANCELLED` path for stock reconcile.
   - Removed early short-circuit for already-cancelled orders.
   - Switched cancel branch from `CancelSalesStockUpdates` to `ReconcileCancelStockUpdates`.
   - Skips `OrderRepository.Update` when current status already `CANCELLED`.
3. `sales/repository/stock_repository_cancel_test.go`
   - Added SQL assertion for reward item filter, reward `item_type` projection, and legacy cancel guard.
   - Replaced weak panic-based reconcile test with pure residual-math assertions.
   - Added reconcile regression tests for partial residual, exact canonical dedupe, over-reversed clamp, reward residual, and warehouse delta sign.
4. `sales/service/order_service.go`
   - Added per-order structured cancel audit log with `ro_no`, `current_status`, `basis_rows`, `basis_total_smallest`, `reward_basis_total`, `residuals_applied`, `warehouse_deltas`.
   - Added helper `cancelAuditValues` to split reward totals using `ItemType`.
5. `sales/service/order_service_test.go`
   - Added helper assertion for cancel audit values.
   - Updated cancel tests to assert reconcile path.
   - Added already-cancelled reconcile test and reward-only Need Review test.
   - Updated mock stock repository for new interface method.

## TDD notes
- Red phase confirmed by compile/test failure before method implementation:
  - missing `ReconcileCancelStockUpdates`
  - SQL assertion targets absent before query changes
- Green phase confirmed by focused repository/service suites passing.

## Contract impact
- Behavior preserved:
  - FE payload contract unchanged: `{"orders":[{"ro_no":"...","data_status":9}]}`
  - per-order transaction boundary preserved
  - tenant filters preserved
- Behavior changed:
  - cancel basis now includes reward lines `item_type IN (1,2)`
  - already-cancelled cancel requests now reconcile stock residual instead of no-op
  - status update skipped when order already cancelled
- Migration implications: none
- Rollback posture:
  - code-only rollback possible
  - no schema/data migration rollback needed
