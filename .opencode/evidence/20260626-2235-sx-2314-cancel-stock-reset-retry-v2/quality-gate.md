# Quality Gate — SX-2314 Cancel Retry v2 (Post-Remediation)

Status: `PASS_WITH_RISKS`
Date: 2026-06-26
Reviewer: `@fixer` (remediation pass)

## Invariants confirmed
- Reward lines included in cancel basis (`item_type IN (1, 2)`)
- Reward item type projected in basis query (`od.item_type AS item_type`)
- Structured audit log per order with required fields (ro_no, current_status, basis_rows, basis_total_smallest, reward_basis_total, residuals_applied, warehouse_deltas)
- `CancelStockBasis` entity carries `ItemType` for reward total split
- Residual math pure helper tested (partial existing, exact match, over-reversed clamp, reward row, canonical dedup)
- Weak panic-wait test replaced by strong pure helper assertions
- Already-cancelled reconcile skips `OrderRepository.Update`
- Legacy cancelAgg guard includes `tr_code='SO' AND tr_no LIKE '%-CO'`
- Warehouse delta sign correct (Qty positive add-back, QtyOnOrder negative reduce)
- Tenant filter preserved on all stock queries/writes
- Per-order transaction boundary preserved

## Test results
- Focused cancel suite: 18 passed
- Full sales suite: 279 passed (267 baseline + 12 new)
- Full build: clean

## Remaining blockers
- **Staging DB access** unavailable in current workspace environment. Cannot verify:
  - Legacy `%-CO` reversal row inventory in staging
  - Reward detail qty rows for the three QA SOs (`SO2606230004`, `SO2606240002`, `SO2606260002`)
  - End-to-end cancel behavior on staging data with real reward rows
- **Action required**: staging access must be provided before deploy to confirm real-world behavior

## Risks
- Reward rows in staging may have empty `qty*_final` columns; COALESCE chain `qty*_final → qty* → qty_po*` covers this but must be validated
- Legacy `tr_code='SO' tr_no LIKE '%-CO'` rows may exist in staging; residual math is conservative but staging evidence needed
- Audit log format uses `log.Infof`; ensure downstream log aggregation parses it correctly

## Follow-ups (next slice)
- Historical backfill script for legacy rows if staging inventory shows non-zero
- Staging end-to-end cancel validation with the three QA SOs
- Verify reward qty column population in staging `sls.order_detail`
