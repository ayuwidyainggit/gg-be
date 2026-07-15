# Quality Gate — SX-2314

Verdict: `PASS_WITH_RISKS`
Date: 2026-06-24
Reviewer: `@quality-gate`

## Status
All non-negotiable invariants met. One non-blocking follow-up closed via local DB evidence.

## Findings closed
- Legacy `tr_code='SO' tr_no='<SO>-CO'` rows: none present in local `ggn_scyllax` (`SELECT tr_code, tr_no, COUNT(*) FROM inv.stock WHERE tr_no LIKE '%-CO' GROUP BY tr_code, tr_no` returned only `tr_code='CO'`). Legacy idempotency guard therefore not required for first slice. Re-check on prod/QA DB before any production deploy; if any `tr_code='SO' tr_no='%SO%CO'` row exists, apply the dual `IN ('CO','SO')` guard to `cancelAgg` and rerun package tests.

## Invariants confirmed
- No FE contract change.
- Cancel writes remain in same transaction as status update (`txCtx` inside `WithinTransaction`).
- All stock queries/writes filter by `cust_id` (sourceAgg, cancelAgg, activeDetailAgg, root WHERE, whDeltas, stock rows).
- New reversal row `tr_code='CO'`, `tr_no='<SO>-CO'`, `qty_in=0`, `qty_out=0`, `qty_in_order=0`, `qty_out_order=cancelQtySmallest`.
- `warehouse_stock.qty += cancelQtySmallest`, `qty_on_order -= cancelQtySmallest`.
- Qty priority final→sales→PO enforced via COALESCE chain.
- No silent skip when outstanding basis exists; Need Review skip only when no row has `QtyOutSmallest>0`.
- No duplicate original SO row on cancel.
- No qty direction inversion.
- No endpoint payload change.
- No new dependency.
- No tenant filter removed.
- No assertion weakening.

## Required before production deploy
- Re-run this check against prod/QA DB for legacy rows; apply guard if found.
- QA cancel test on real SO per plan step 50 (Warehouse Stock and On Cust Order reset).
