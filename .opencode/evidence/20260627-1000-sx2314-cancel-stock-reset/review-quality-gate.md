# SX-2314 Quality Gate Review

**Status:** PASS  
**Reviewer:** @quality-gate  
**Date:** 2026-06-27  
**Plan:** `.opencode/plans/20260627-1000-sx2314-cancel-stock-reset.md` (v4)  
**Evidence:** `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/`

---

## Scope Checked

- `sales/repository/stock_repository.go` cancelStockBasisQuery (2 Where edits)
- `sales/repository/stock_repository_cancel_test.go` (1 test update + 2 new tests)
- Test results, DB simulation evidence, code diffs

---

## Decision

**PASS** — All requested checks pass. Diff is exactly bounded. Tests cover the new clauses. DB simulation is internally consistent. No scope creep.

---

## Findings

### Check 1: stock_repository.go diff is exactly 2 Where edits ✓

**File:** `sales/repository/stock_repository.go`

**Edit 1 — Line 296 (cancelAgg Where):**
```go
// BEFORE:
Where("c.cust_id = ? AND c.tr_no = ? AND c.tr_code = 'CO'", custID, cancelTrNo)

// AFTER:
Where("c.cust_id = ? AND c.tr_no = ? AND (c.tr_code = 'CO' OR (c.tr_code = 'SO' AND c.tr_no LIKE '%-CO%'))", custID, cancelTrNo)
```

**Edit 2 — Lines 345-356 (main Where):**
```go
// BEFORE:
Where(`od.cust_id = ? AND od.ro_no = ? AND od.item_type = 1
    AND (
        COALESCE(od.qty1_final, 0) > 0
        ...

// AFTER:
Where(`od.cust_id = ? AND od.ro_no = ?
    AND (
        COALESCE(od.qty1_final, 0) > 0
        ...
```

**Note:** Line 307 `activeDetailAgg` correctly retains `item_type = 1` per R1 decision to avoid flagging order+reward-of-same-product as ambiguous.

**Confirms:** Exactly 2 Where clause edits. No other changes to the query builder, no signature changes, no other files touched.

---

### Check 2: Test file diff is exactly 1 update + 2 new tests ✓

**File:** `sales/repository/stock_repository_cancel_test.go`

**Updated test:** `TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula` (lines 197-206)
- Added 3 new assertions for the legacy clause:
  - Asserts `c.tr_code = 'CO'` still present
  - Asserts `c.tr_code = 'SO'` present (legacy rows)
  - Asserts `c.tr_no LIKE '%-CO%'` present

**New test 1:** `TestGetCancelStockBasisQuery_IncludesRewardLine` (lines 209-232)
- Asserts `qty_outstanding` and `active_detail_count` still in SELECT
- Asserts `od.item_type = 1` appears at most once (only in `activeDetailAgg`, not in main query)

**New test 2:** `TestGetCancelStockBasisQuery_LegacySORowsIncludedInCancelAgg` (lines 234-255)
- Asserts the full combined predicate: `c.tr_no = ? AND (c.tr_code = 'CO' OR (c.tr_code = 'SO' AND c.tr_no LIKE '%-CO%'))`
- Renamed from `...ExcludedFromCancelAgg` to `...IncludedInCancelAgg` per @oracle nit

**Confirms:** One test updated, two new tests added. No other test changes.

---

### Check 3: DB simulation evidence is internally consistent ✓

**File:** `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/simulate-cancel-v2.txt`

**BEFORE:**
- `warehouse_stock.qty=14, qty_on_order=10`
- `inv.stock` on_cust_projection=10, wh_stock_projection=2

**AFTER (post-fix):**
- `warehouse_stock.qty=24, qty_on_order=0`
- `inv.stock` on_cust_projection=0, wh_stock_projection=2

**Two new reversal rows:**
- `ref_det_id=7540, qty_out_order=4` (order line)
- `ref_det_id=7541, qty_out_order=6` (reward line)

**Math verification:**
- `warehouse_stock.qty`: 14 + (4 + 6) = 24 ✓
- `warehouse_stock.qty_on_order`: 10 - (4 + 6) = 0 ✓
- `inv.stock` on_cust_projection: SUM(qty_in_order) - SUM(qty_out_order) = 10 - 10 = 0 ✓
- `inv.stock` wh_stock_projection: unchanged at 2 (CO rows have qty_in=0, qty_out=0) ✓
- Transaction rolled back: ROLLBACK confirmed ✓

**Confirms:** DB simulation is internally consistent. The fix correctly processes both order and reward lines.

---

### Check 4: No scope creep ✓

**Files touched:**
- `sales/repository/stock_repository.go` (2 Where edits)
- `sales/repository/stock_repository_cancel_test.go` (1 update + 2 new)

**Files NOT touched (per plan):**
- No controllers
- No migrations
- No manifests / docker-compose
- No `go.mod` / `go.sum`
- No other services (inventory, master, pjp, etc.)
- No FE changes

**Confirms:** Diff is exactly bounded to the plan. No unintended changes.

---

### Check 5: Single-line orders (no reward) still handled correctly ✓

**Scenario:** Order with only `item_type=1` rows (no reward lines).

**Analysis:**
- Main Where now matches all rows where `qty*_final > 0 OR qty* > 0 OR qty_po* > 0` (no `item_type` filter)
- For single-line orders, only `item_type=1` rows exist, so they still pass the qty predicate
- `cancelAgg` OR is additive: `c.tr_code = 'CO'` clause still present, so existing CO reversal rows are still matched
- `TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula` still asserts `c.tr_code = 'CO'` (line 173-175), confirming the original clause is preserved

**Confirms:** Single-line orders unaffected. The OR clause is additive, not replacement.

---

### Check 6: Legacy rows now correctly subtract from qty_outstanding ✓

**Mechanism:**
- `cancelAgg` now matches both:
  - New CO rows: `c.tr_code = 'CO' AND c.tr_no = '<SO>-CO'`
  - Legacy SO rows: `c.tr_code = 'SO' AND c.tr_no LIKE '%-CO%'`
- `cancelAgg` SELECT: `SUM(c.qty_out_order) AS qty_out_order_cancel`
- Main SELECT: `qty_outstanding = qty_out_so - qty_out_order_cancel`
- `GREATEST(..., 0)` clamp prevents negative outstanding

**Test coverage:** `TestGetCancelStockBasisQuery_LegacySORowsIncludedInCancelAgg` asserts the full combined predicate (line 249).

**Confirms:** Legacy rows are now included in `cancelAgg`, their `qty_out_order` is summed, and the sum is subtracted from outstanding. Re-cancel is idempotent.

---

### Check 7: Staging redeploy / live curl properly deferred ✓

**Plan line 334:** "The staging redeploy / live curl to staging is deferred (this clone cannot reach the staging container from this session due to docker daemon I/O errors); the local DB simulation is the primary evidence."

**Evidence:** `simulate-cancel-v2.txt` is the primary proof. No staging live curl attempted.

**Confirms:** Staging verification is explicitly deferred as a follow-up ops item. Local DB simulation is the primary evidence, which is consistent with the plan.

---

## Source Basis Checked

- Plan v4 (`.opencode/plans/20260627-1000-sx2314-cancel-stock-reset.md`)
- Code diffs (`sales/repository/stock_repository.go`, `sales/repository/stock_repository_cancel_test.go`)
- Test results (`.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/tests-pass.txt`, `tests-fail.txt`, `tests-full.txt`)
- DB simulation (`.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/simulate-cancel-v2.txt`)
- go vet output (`.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/go-vet.txt`)
- Oracle review verdict (PASS, one minor naming nit fixed)

---

## Required Before PASS

None. All checks pass.

---

## Remediation Worklist

N/A — Status is PASS.

---

## Recommended Follow-ups

1. **Staging redeploy and live curl (deferred ops item):**
   - Deploy the patch to staging
   - Run the cancel flow for `SO2606260002` (or similar SO+reward order) in staging
   - Verify `warehouse_stock.qty` rises by 10, `qty_on_order` falls by 10, `on_cust_projection` = 0
   - This is a follow-up ops item, not a blocker for local PASS

2. **Legacy row audit (informational):**
   - 118 legacy rows exist across 4 tenants (per `legacy-audit-post.txt`)
   - No backfill or deletion needed; the new query treats them defensively
   - Monitor for any unexpected behavior in production cancel flows for orders that already have legacy rows

---

## Residual Risks

1. **Staging verification gap:** Local DB simulation proves the fix works, but staging has not been live-tested. Risk: low, since the local DB is a 1:1 copy of staging and the simulation is in-tx rolled back. Mitigation: deploy to staging and verify before production rollout.

2. **Legacy row edge cases:** 118 legacy rows exist. The new query includes them defensively. Risk: very low, since `GREATEST(..., 0)` prevents negative outstanding, and the legacy rows represent prior reversals that should be accounted for. Mitigation: monitor production cancel flows for orders with legacy rows.

3. **Test naming inconsistency (cosmetic):** `tests-pass.txt` shows old test name `LegacySORowsExcludedFromCancelAgg` (pre-oracle-rename), but current code has `LegacySORowsIncludedInCancelAgg`. This is a cosmetic evidence sequencing issue — `tests-full.txt` was captured post-rename and passes. Not a blocker.

---

## Escalation

None. All checks pass. Ready for commit and staging deploy.

---

## Summary

The patch is minimal, correct, and well-tested. It fixes two confirmed defects:
1. Reward lines (item_type=2) are now surfaced in cancel basis
2. Legacy tr_code='SO' tr_no LIKE '%-CO%' rows are now matched by cancelAgg

The diff is exactly bounded to the plan. No scope creep. No signature changes. No migrations. Tests cover the new clauses. DB simulation is internally consistent and proves the fix works as expected.

**Verdict: PASS**
