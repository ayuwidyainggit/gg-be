# SX-2314 Discovery — Cancel Order not reset Warehouse Stock & On Cust Order

## Source anatomy (confirmed by reading code)

- Controller: `sales/controller/order_controller.go` line 895+ — `PATCH /sales/v1/orders/status` parses `entity.BulkUpdateStatusOrder` and calls `OrderService.BulkUpdateStatus(custId, request)` (line 933).
- Service: `sales/service/order_service.go` `BulkUpdateStatus` (line 5131).
  - Per-item transaction in `service.Transaction.WithinTransaction`.
  - Cancel branch (line 5158) loads order via `OrderRepository.FindByNo`, validates transition via `validateCancelTransition` (only `NEED_REVIEW=1` and `PROCESSED=2` can transition to `CANCELLED=9`).
  - Then loads `StockRepository.GetCancelStockBasis(txCtx, custId, roNo)`.
  - **Suspect branch** (line 5182): if `orderData.DataStatus == NEED_REVIEW` and no basis row has `QtyOutSmallest > 0`, sets `skipCancelStockWrite = true` and only logs a warn, then still updates the order to 9.
  - If not skipped, builds `cancelOrderStockBasis` and calls `StockRepository.CancelSalesStockUpdates(txCtx, roNo, stockDate, commands)`.
  - Finally `OrderRepository.Update(...)` flips `data_status` to 9.
- Repository qty source priority (`sales/repository/stock_repository.go` `cancelStockBasisQuery`, line 270+):
  1. `inv.stock` `tr_code='SO'` aggregate `SUM(qty_out - qty_in) AS qty_out_so` (line 282). If `> 0`, used as authoritative outstanding qty.
  2. Else fallback from `sls.order_detail` per unit: `qty*_final → qty* → qty_po*` converted to smallest unit.
  3. Subtracts any prior `tr_code='CO' qty_out_order` for idempotency.
- Repository reversal write (`CancelSalesStockUpdates`, line 409): calls `buildCancelStockMutations` (line 231) which emits:
  - `inv.stock` row: `tr_code='CO'`, `tr_no='<SO>-CO'`, `qty_in=0`, `qty_out=0`, `qty_in_order=0`, `qty_out_order=row.QtyOutSO`, `unit_price=row.UnitPrice`, `ref_det_id=row.RefDetID`.
  - `inv.warehouse_stock` upsert delta: `qty += QtyOutSO`, `qty_on_order += -QtyOutSO` (via `UpsertWithExistingValueArr` at line 514, which uses `qty + EXCLUDED.qty` and `qty_on_order + EXCLUDED.qty_on_order`).
- Then `StoreBulk` inserts the reversal stock rows.

## Confirmed vs assumed audit

| Claim | Status | Evidence |
| --- | --- | --- |
| Endpoint `PATCH /sales/v1/orders/status` is wired to `BulkUpdateStatus` | confirmed_repo | `sales/controller/order_controller.go:895,933` |
| `BulkUpdateStatus` runs cancel inside a single `WithinTransaction` | confirmed_repo | `sales/service/order_service.go:5157` |
| Cancel only allows transitions from `NEED_REVIEW (1)` and `PROCESSED (2)` | confirmed_repo | `validateCancelTransition` at `order_service.go:5068` |
| Need Review cancel with empty basis silently skips stock write but still updates status | confirmed_repo | `order_service.go:5182-5192` and existing test `TestBulkUpdateStatus_Cancel_NeedReview_EmptyBasisShouldSkipStockWriteAndUpdateStatus` at `order_service_test.go:656` |
| Reversal writes `tr_code='CO'` and `tr_no='<SO>-CO'` with `qty_out=0, qty_out_order=QtyOutSO` | confirmed_repo | `buildCancelStockMutations` at `stock_repository.go:240-257` |
| `warehouse_stock` upsert adds `QtyOutSO` to `qty` and subtracts `QtyOutSO` from `qty_on_order` | confirmed_repo | `UpsertWithExistingValueArr` at `stock_repository.go:514-528` |
| Qty source priority is `qty_out_so` → `qty*_final → qty* → qty_po*` | confirmed_repo | `cancelStockBasisQuery` GREATEST clause at `stock_repository.go:381-394` |
| `qty_out_so` is `SUM(qty_out - qty_in)` so any prior `qty_in` row for the SO can zero or negate it | confirmed_repo | `stock_repository.go:282` |
| QA's `SO2606230004` currently has `data_status=9` and no reversal row written | unverified | needs DB query against `sls.order` and `inv.stock` |
| QA's SO has `inv.stock` `tr_code='SO'` rows at all | unverified | needs DB query |
| QA's SO has `qty*_final/qty*/qty_po*` populated in `sls.order_detail` | unverified | needs DB query |

## Most likely root cause (ranked)

1. **Need Review with no outstanding basis → silent skip.** If `SO2606230004` was at `data_status=1` and `GetCancelStockBasis` returned 0 rows with `QtyOutSmallest > 0`, the code takes the `skipCancelStockWrite = true` branch. Status flips to 9, no reversal row, no `warehouse_stock` change. Matches the QA symptom exactly.
2. **Net `qty_out - qty_in` for the SO's ledger is 0.** Even if the SO was previously processed, a prior partial reversal/return row with `qty_in > 0` makes `qty_out_so = 0`, and the fallback `qty*_final/qty*/qty_po*` may also be empty (e.g. PO tab only had values, but cancelled before any `inv.stock` `tr_code='SO'` row recorded `qty_out_order` distinctly). Code treats it as "nothing to reverse" and skips. Same visible symptom as #1.
3. **Idempotency guard interacted with branch.** Less likely from code reading, but if a prior cancel ran and a `CO` row already exists, `qty_outstanding = qty_out_so - qty_out_order_cancel` is 0; combined with Need Review, the skip path triggers.

## Fix plan (single bounded patch)

### 1. Repro and capture evidence
Run the SQL probes in the issue against `best.scyllax.online` to confirm which path triggered for `SO2606230004`:
- `sls.order.data_status` and `ro_date`.
- `sls.order_detail` for `qty*_final/qty*/qty_po*` and `conv_unit*`.
- `inv.stock` rows with `tr_no IN ('SO2606230004', 'SO2606230004-CO')`.
- `inv.warehouse_stock` for the affected `(wh_id, pro_id)`.

Save outputs under `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/repro-*.sql.txt`.

### 2. Code change in `BulkUpdateStatus` cancel branch
File: `sales/service/order_service.go` around line 5182.

Replace the silent skip with a fallback that still writes a reversal using the SO's detail-level basis (the `qty*_final/qty*/qty_po*` fallback that `GetCancelStockBasis` already computes but currently discards when `QtyOutSmallest == 0`). Concretely:

- If Need Review and `QtyOutSmallest` of every basis row is 0, do **not** skip; instead, if `GetCancelStockBasis` returned at least one row, attempt reversal with `QtyOutSmallest = 0` filtered out (so commands list is empty and `buildCancelStockWriteCommands` would skip them → still no write). This means we need a second resolver: re-query the basis but force the fallback path even when `qty_out_so` is 0. The cleanest implementation:
  - Add a new repository method `GetCancelStockBasisForNeedReviewFallback` (or extend `GetCancelStockBasis` with a flag) that sets `qty_out_so = 0` in the `sourceAgg` and reads only the detail fallback as the basis.
  - When Need Review and no outstanding source qty, call this fallback resolver and, if it returns rows with `QtyOutSmallest > 0`, build commands and call `CancelSalesStockUpdates`.
  - Only when both source and fallback return nothing (truly zero qty), keep the current skip-and-log behavior.

Simpler alternative if a new method is too much: change the `is_missing_source` logic in `cancelStockBasisQuery` so that when `qty_out_so <= 0` we still surface the detail-fallback qty (no `GREATEST(..., 0)` clamp at the outer level) and only treat `QtyOutstanding = 0 AND QtyOutSmallest = 0` as "no basis". Then the existing `validateCancelStockBasis` and Need Review branch can use the same code path. The downside: small surface area change in a heavily-tested query (`stock_repository_cancel_test.go`). Prefer this only after confirming with the test file's existing expectations.

I will go with the **new method** approach to keep the existing test invariants intact.

### 3. Validation
- Add/update unit tests:
  - `TestBulkUpdateStatus_Cancel_NeedReview_DetailOnlyFallbackShouldApplyReversal` — Need Review SO with no `inv.stock` SO ledger but detail rows with `qty*_final > 0`; expect `CancelSalesStockUpdates` called once with the fallback qty.
  - Update `TestBulkUpdateStatus_Cancel_NeedReview_EmptyBasisShouldSkipStockWriteAndUpdateStatus` to assert the truly-empty case (no source and no detail fallback) still skips.
- Manual DB validation against the SQL probes in the issue for `SO2606230004` after the fix is deployed.

### 4. Out of scope
- No change to the controller, the URL, the payload, or the auth/tenant contract.
- No new migrations.
- No change to the "summary" / projection tables consumed by the FE — the prior `SX-1241` plan covers that and remains parked.
