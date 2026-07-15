# Jira Comment Draft — SX-2314 cancel retry v2

Update ready in `sales` for plan `.opencode/plans/20260626-2235-sx-2314-cancel-stock-reset-retry-v2.md`.

What changed:
- cancel basis now includes reward rows (`item_type IN (1,2)`)
- cancel basis now projects `item_type` so audit can split reward totals
- cancel reconcile still counts legacy `tr_code='SO' tr_no LIKE '%-CO'` reversal rows
- already-cancelled PATCH cancel now re-runs reconcile instead of no-op
- per-order structured audit log added in cancel path:
  - `ro_no`
  - `current_status`
  - `basis_rows`
  - `basis_total_smallest`
  - `reward_basis_total`
  - `residuals_applied`
  - `warehouse_deltas`

Test status:
- `rtk go test ./repository/... ./service/... -count=1` ✅
- `rtk go build ./...` ✅

Local regression coverage added for:
- reward basis projection
- residual math
- exact canonical dedupe
- over-reversed clamp-to-zero behavior
- warehouse delta sign
- already-cancelled reconcile without status update
- reward-only Need Review cancel

Blocker still open:
- staging DB access unavailable in current workspace, so no staging query results attached yet for:
  - legacy `%-CO` reversal inventory
  - reward detail qty rows for `SO2606230004`, `SO2606240002`, `SO2606260002`

When staging access available, please run:
```sql
SELECT tr_code, tr_no, COUNT(*)
FROM inv.stock
WHERE tr_no LIKE '%-CO'
GROUP BY tr_code, tr_no
ORDER BY 1;
```

```sql
SELECT order_detail_id, item_type, pro_id, qty, qty_final, qty_po,
       qty1, qty2, qty3, qty1_final, qty2_final, qty3_final,
       qty_po1, qty_po2, qty_po3, conv_unit2, conv_unit3
FROM sls.order_detail
WHERE ro_no IN ('SO2606230004','SO2606240002','SO2606260002')
ORDER BY ro_no, order_detail_id;
```
