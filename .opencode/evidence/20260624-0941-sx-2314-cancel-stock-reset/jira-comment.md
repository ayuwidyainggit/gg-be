# SX-2314 — BE Fix Summary (comment draft for Jira)

## Status
- Code shipped on branch `bugfix/SX-2314-dev` (commit `20ee58f`).
- Repository: `sales` service.
- Local validation passed against `ggn_scyllax` via the running `scylla-sales` container.

## Root cause
Two issues in `PATCH /v1/orders/status` cancel path (`OrderController.UpdateStatus` → `OrderService.BulkUpdateStatus`):

1. `buildCancelStockMutations` was inserting two `inv.stock` rows for every cancel basis row: a duplicate source row (`tr_code='SO' tr_no=<SO> qty_in=…`) plus a reversal row whose `tr_code` was also `'SO'` and `tr_no='<SO>-CO'`. The reversal row violated the cancel docs (`tr_code` must be `'CO'`), and the duplicate source row skewed idempotency and warehouse_stock totals.
2. The cancel-aggregate subquery (`cancelAgg` inside `cancelStockBasisQuery`) used `tr_code='SO'` to detect prior cancel rows, so legitimate `CO` reversals were never counted and a second cancel of the same order would double-reverse.
3. For orders sitting at `data_status=1` (Need Review) where no `inv.stock` SO source row was ever written, the service short-circuited to "update status only" even though `sls.order_detail` still held a positive qty. Detail qty was never used as a fallback basis.

## Fix
File-by-file:

- `sales/repository/stock_repository.go`
  - `buildCancelStockMutations` now emits a single reversal stock row: `tr_code='CO'`, `tr_no='<SO>-CO'`, `qty_in=0`, `qty_out=0`, `qty_in_order=0`, `qty_out_order=row.QtyOutSO` (smallest unit). Removed the duplicate `SO` source row.
  - `cancelStockBasisQuery` `cancelAgg` filter switched to `c.tr_code='CO'` so prior `CO` reversals reduce the outstanding amount. No more double reversal on re-cancel.
  - `cancelStockBasisQuery` now `LEFT JOIN mst.m_product mp` for product conv units and computes `qty_outstanding` / `qty_out_smallest` with a CASE:
    - Use `inv.stock` SO source (`SUM(qty_out - qty_in)`) when positive.
    - Otherwise fall back to the detail qty in smallest unit with priority `qty1_final → qty1 → qty_po1` (and same for `qty2*`, `qty3*`) using `conv_unit2`, `conv_unit3` from `sls.order_detail` (with `mst.m_product` fallback to 1).
  - `is_missing_source` is preserved as true when the source ledger is absent, but the fallback qty is what drives the reversal amount.

- `sales/service/order_service.go`
  - `validateCancelStockBasis` allows `IsMissingSource` only when `QtyOutSmallest > 0` from the fallback (so Need Review cancel no longer fails for orders with detail qty but no source ledger).
  - `validateFinalOrderStockBasis` tightened to `(IsMissingSource && QtyOutSmallest == 0)` so existing final/invoice flows still reject zero-basis without source.

## Tests
- Updated: `TestBuildCancelStockMutations_SingleSKU`, `TestBuildCancelStockMutations_MultiSKUAndIdempotent`, `TestBuildCancelStockMutations_PartialReverseExistingRefStillBuildsOutstanding` to expect the new `CO` reversal shape.
- Updated: `TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula` to assert `c.tr_code = 'CO'` in the dry-run SQL.
- Added: `TestCancelStockBasisFallback_DetailOnlyPOCase` covering PO-only detail with full conv-unit fallback expression.
- Updated: `TestBulkUpdateStatus_Cancel_NeedReview_MissingSourceBasis…` renamed to `…WithOutstandingFallbackShouldApplyReversal`. It now exercises the detail-fallback path and asserts `CancelSalesStockUpdates` is called once with the expected row.
- Result: `rtk go test ./repository/... ./service/... -count=1` → 267 passed in 2 packages.

## Local validation
- `rtk docker compose -f docker-compose.yml up -d sales` (container `scylla-sales` on `0.0.0.0:9004`).
- JWT signed with `sales/.env` `JWT_SECRET_KEY=secret` (HS256, adminbm `cust_id=C220010001`).
- Test 1: cancel `SO2606190007` (Need Review, wh 63) → response 200, 1 new `inv.stock` row `tr_code='CO' tr_no='SO2606190007-CO' qty_out_order=1`. Retry: 0 new rows, still 200.
- Test 2: cancel `SO2606180001` (Need Review, wh 63, pro 478) → `inv.warehouse_stock (cust=C220010001 wh=63 pro=478)`: `qty` `−199 → −197` (+2), `qty_on_order` `−2521 → −2523` (−2). One `inv.stock` row `qty_out_order=2`. Retry idempotent.
- Tenant safety: detail rows with the same `ro_no` but different `cust_id` are excluded from the basis (verified with `SO2606190007` pro 10743 row belonging to `C260020001`).
- Local `ggn_scyllax` does not contain `SO2606230004`, so the original ticket's sample was not directly reproduced; substituted with Need Review orders that have a valid basis to exercise the same code path.

## Plan and quality gate
- Plan: `.opencode/plans/20260624-0941-sx-2314-cancel-stock-reset.md` (`PASS_FOR_SLICE`).
- Evidence: `.opencode/evidence/20260624-0941-sx-2314-cancel-stock-reset/{discovery,execution,quality-gate,validation}.md`.
- Quality gate verdict: `PASS_WITH_RISKS`. Only follow-up: confirm prod/QA DB has no legacy `tr_code='SO' tr_no='%SO%CO'` rows before deploy; local DB shows none, so no dual `IN ('CO','SO')` guard needed for first slice.

## Files changed
- `sales/repository/stock_repository.go` (+32 / −21)
- `sales/repository/stock_repository_cancel_test.go` (+50 / −18)
- `sales/service/order_service.go` (+2 / −2)
- `sales/service/order_service_test.go` (+15 / −6)

Total 4 files, 99 insertions, 47 deletions.

## Risk / follow-up
- Re-run the legacy `tr_code='SO' tr_no='%SO%CO'` query on the prod/QA DB before merge. If any rows exist, extend `cancelAgg` to `c.tr_code IN ('CO','SO')` so legacy reversals still count for idempotency.
- QA should reproduce the original `SO2606230004` scenario post-deploy to confirm Warehouse Stock and On Cust Order reset on first cancel and stay steady on retry.

## Out of scope (next slice)
- One-time remediation of historical `inv.stock` reversal rows created with the wrong `tr_code`.
- Any cross-service stock reporting changes.
- FE-side verification (no FE change required).

Commit: `20ee58f` on `bugfix/SX-2314-dev`.
