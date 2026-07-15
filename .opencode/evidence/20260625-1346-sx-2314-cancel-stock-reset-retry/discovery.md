# Discovery — SX-2314 Retry

Task id: `20260625-1346-sx-2314-cancel-stock-reset-retry`
Parent: `.opencode/plans/20260624-0941-sx-2314-cancel-stock-reset.md`
Commit baseline: `20ee58f` on branch `bugfix/SX-2314-dev`

## Why QA still FAILED
QA re-tested `SO2606240002` after `20ee58f` and Warehouse Stock + On Cust Order still did not reset. Re-reading the code on the current branch surfaced two gaps that match the symptom exactly:

1. **Service short-circuits on already-cancelled** (`sales/service/order_service.go:5179-5181`):
   ```go
   if *orderData.DataStatus == entity.CANCELLED {
       return nil
   }
   ```
   Once the order is flipped to `data_status=9`, any subsequent PATCH that re-evaluates the basis and applies the residual warehouse delta is skipped. If a prior attempt (older code or rolled-back transaction) inserted the reversal row but failed to apply `warehouse_stock`, the warehouse stays stuck forever.

2. **`cancelAgg` ignores legacy reversal rows** (`sales/repository/stock_repository.go:296`):
   ```go
   Where("c.cust_id = ? AND c.tr_no = ? AND c.tr_code = 'CO'", custID, cancelTrNo)
   ```
   Pre-patch environments may still have reversal rows with `tr_code='SO' tr_no='<SO>-CO'`. Those rows are not counted by `cancelAgg`, so the basis `QtyOutSmallest` overstates the true residual. The warehouse delta is then applied with an inflated positive qty, which can move the warehouse in the wrong direction relative to the desired target.

The two gaps combine to: a partially corrupted staging environment (from older code paths) does not self-heal, because the new code never re-runs the path on a cancelled order.

## Source anatomy for the fix
- `sales/repository/stock_repository.go`
  - `cancelStockBasisQuery` (lines 270-398): SQL that combines `inv.stock` SO source, `inv.stock` prior-cancel aggregate, and `sls.order_detail` fallback.
  - `cancelAgg` subquery: filter on prior reversal rows. Currently `tr_code='CO'`. Needs to include legacy `tr_code='SO'` rows whose `tr_no LIKE '%-CO'`.
  - `CancelSalesStockUpdates` (lines 409-445): current bulk insert + upsert path. Does not consult existing reversal rows.
  - `UpsertWithExistingValueArr` and `StoreBulk`: low-level helpers, reusable.
- `sales/service/order_service.go` `BulkUpdateStatus` (lines 5131-5252): per-order transaction loop. The cancel branch short-circuits on `data_status=9`.
- `sales/entity/stock.go` `CancelStockBasis` / `CancelStockWrite`: enough fields for residual math.

## Existing tests as anchors
- `sales/repository/stock_repository_cancel_test.go`:
  - `TestBuildCancelStockMutations_SingleSKU` already asserts 1 reversal row, `tr_code='CO'`, `tr_no='<SO>-CO'`.
  - `TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula` asserts `c.tr_code = 'CO'`. Will need updating to also assert legacy inclusion.
  - `TestCancelStockBasisFallback_DetailOnlyPOCase`: detail qty fallback to smallest unit.
- `sales/service/order_service_test.go`:
  - `TestBulkUpdateStatus_Cancel_NeedReview_MissingSourceBasisWithOutstandingFallbackShouldApplyReversal`: covers fallback path. Will need a sibling test for already-cancelled reconcile.

Baseline: `rtk go test ./repository/... ./service/... -count=1` → 267 passed in 2 packages (already verified after `20ee58f`).

## Constraints
- Per repo `AGENTS.md`: `rtk`-prefixed shell, `cust_id` filter, `WithinTransaction` per order, tenant safety, no destructive writes against shared remote environments, no FE contract changes.
- Source edits remain blocked in planner lane; implementation must be carried out by `@backend`/`@fixer`.

## Decision for the fix
1. Remove the early `return nil` for already-cancelled. Always evaluate basis.
2. `ReconcileCancelStockUpdates(c, orderNo, stockDate, basis []CancelStockBasis)` does:
   - Look up existing reversal rows for `tr_no='<SO>-CO'` (any tr_code), grouped by `(wh, pro, ref_det)`, sum `qty_out_order`.
   - For each basis row: residual = basis.QtyOutSmallest − existingSum (clamp ≥ 0).
   - Insert a single new `tr_code='CO' tr_no='<SO>-CO' qty_out_order=residual` row only when no exact existing canonical row is present and residual > 0.
   - Apply `warehouse_stock` delta only when residual > 0.
   - Negative residual: do nothing (treat as already-corrected; never write negative warehouse delta).
3. Widen `cancelAgg` filter to count legacy `tr_code='SO' tr_no LIKE '%-CO'` rows.
4. Update service to allow `CANCELLED` as a permitted source status for reconcile.

## Risk
- Negative residual handling is critical: if a previous attempt over-applied warehouse_stock delta (e.g. legacy bug wrote positive qty twice), residual = basis − existing sum is negative. We must not write a negative warehouse delta in that case, or the warehouse would drift in the opposite direction.
- The wider cancelAgg filter is harmless on databases without legacy rows (none match the additional predicate).