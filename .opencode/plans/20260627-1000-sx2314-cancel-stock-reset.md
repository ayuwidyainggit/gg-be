# SX-2314 Plan v4 — Local DB Simulation Confirms Reward-Line Bug & Legacy Risk

## Goal
Fix two confirmed defects in `sales/repository/stock_repository.go` `cancelStockBasisQuery` discovered by simulating the documented cancel payload against the local `ggn_scyllax` database (synced 1:1 from staging):
1. **Reward-line bug**: the cancel only processes `sls.order_detail.item_type = 1` (order line) and skips `item_type = 2` (reward line). For `SO2606260002` (order `0 0 4` PCS + reward `0 0 6` PCS of product `10838`), the cancel would emit a reversal of only `qty_out_order=4` PCS, leaving 6 PCS reserved on `On Cust Order` projection after the cancel. QA expects both to reset to `0 0 0`.
2. **Legacy `tr_code='SO' tr_no LIKE '%-CO%'` rows**: 118 such rows exist in the local DB across 4 tenants and 92 distinct orders (verified via legacy audit). The current `cancelAgg` filter `c.tr_code = 'CO' AND c.tr_no = '<SO>-CO'` misses these rows, which can cause `qty_outstanding` to compute incorrectly on retry / re-cancel of affected orders.

The two fixes are independent and both required for QA's "no double-reversal, no partial reset" acceptance criteria. They land in the same `cancelStockBasisQuery` function and are shipped together as a single small bounded change.

## Non-goals
- No change to controllers, migrations, manifests, `go.mod`/`go.sum`, or any other file in `scylla-be`.
- No change to the public HTTP contract, payload shape, or auth/tenant ownership check.
- No change to `CancelSalesStockUpdates`, `buildCancelStockMutations`, or `BulkUpdateStatus` (verified correct via simulation).
- No change to FE.
- No data backfill of legacy rows; we just make the cancel defensive against them.

## Scope
- `sales/repository/stock_repository.go` `cancelStockBasisQuery`:
  - Drop or relax the `od.item_type = 1` filter in the main `Where` and in `activeDetailAgg` `Where` so reward lines are also processed. The earlier `orderDetailAgg` (in `BulkUpdateStatus` flow, line 4696) and other unrelated `item_type = 1` filters in other repos stay untouched.
  - Tighten `cancelAgg` `Where` so the legacy `tr_code='SO' tr_no LIKE '%-CO%'` rows are also matched (otherwise `qty_outstanding` will be inflated for orders that already have a partial legacy reversal). Concretely: match `c.tr_no = '<SO>-CO'` AND `(c.tr_code = 'CO' OR (c.tr_code = 'SO' AND c.tr_no LIKE '%-CO'))`. The legacy rows have `qty_out_order > 0` (set by the legacy create-time path that wrote `<SO>-CO` with `tr_code='SO'`), so including them in the cancel sum correctly accounts for prior reversals.
- `sales/repository/stock_repository_cancel_test.go`:
  - Update `TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula` to assert the new `cancelAgg` clause and the absence of `item_type = 1` from the main query.
  - Add `TestGetCancelStockBasisQuery_IncludesRewardLine` proving `item_type=2` rows are surfaced.
  - Add `TestGetCancelStockBasisQuery_LegacySORowsExcludedFromCancelAgg` proving the legacy `tr_code='SO' tr_no='<SO>-CO'` rows are matched by `cancelAgg`.
- `sales/repository/stock_repository_cancel_test.go` and downstream mocks (no change needed; the change is internal to the query builder).

## Requirements
1. After cancel of `SO2606260002` (order `0 0 4` + reward `0 0 6`, product `10838`, cust `C260020001`, wh `350`):
   - `inv.stock` gains two new rows, both `tr_code='CO' tr_no='SO2606260002-CO'`, one with `qty_out_order=4 ref_det_id=7540`, one with `qty_out_order=6 ref_det_id=7541`.
   - `inv.warehouse_stock.qty` for `(C260020001, 350, 10838)` rises by 10 PCS (from 14 to 24 — yes 24, not 18; previously the code missed the reward, so the increase was only 4).
   - `inv.warehouse_stock.qty_on_order` for the same key falls by 10 PCS (from 10 to 0).
   - `inv.stock` `on_cust_projection = SUM(qty_in_order) - SUM(qty_out_order) = (4+6) - (4+6) = 0` for the product key.
   - `inv.stock` `wh_stock_projection = SUM(qty_in) - SUM(qty_out) = 12 - 10 = 2` (unchanged; the original SO rows' `qty_out=10` is still in the ledger, so Wh Stock stays at 2 PCS).
2. Re-cancel of the same SO is idempotent: no second reversal row.
3. For an order that has a legacy `tr_code='SO' tr_no='<SO>-CO'` row (e.g. `SO2606240001` with `stock_id=26818, qty_out_order=4` from local audit), the cancel's `qty_outstanding` correctly subtracts that legacy row's `qty_out_order` so no double reversal occurs.
4. Local `cd sales && rtk go test ./repository/... ./service/... -count=1` passes.
5. No new error in any service log.

## Acceptance Criteria
1. `rtk go test ./repository/... ./service/... -count=1` all green, including:
   - `TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula` (updated).
   - `TestGetCancelStockBasisQuery_IncludesRewardLine` (new).
   - `TestGetCancelStockBasisQuery_LegacySORowsExcludedFromCancelAgg` (new).
   - Existing `TestBuildCancelStockMutations_*` and `TestBulkUpdateStatus_Cancel_*` still pass.
2. Manual SQL simulation against the local `ggn_scyllax` DB (re-run the documented probes in V-track) shows:
   - `qty_outstanding` for `SO2606260002` returns `QtyOutSmallest=4` for ref_det_id 7540 AND `QtyOutSmallest=6` for ref_det_id 7541 (sum 10).
   - The two reversal rows are produced.
   - `inv.warehouse_stock` delta = `+10 / -10`.
   - `inv.stock` on_cust projection = 0.
   - `inv.stock` wh_stock projection = 2 (unchanged).
3. Idempotency: running the cancel twice on `SO2606260002` produces exactly one `<SO>-CO` reversal pair (no duplicates).
4. For a legacy-touch order (e.g. `SO2606240001`), re-cancel does not double-reverse (validated via simulation with the legacy row present in `cancelAgg`).
5. The production code path in `sales/service/order_service.go` `BulkUpdateStatus` cancel branch is unchanged. Only the query layer changes.

## Existing patterns / reuse
- Reuse the existing `cancelStockBasisQuery` builder pattern; just edit the `Where` clauses.
- Reuse the existing `entity.CancelStockBasis` shape; no DTO change.
- Reuse `cancelStockBaseRow`, `buildCancelStockMutations`, `CancelSalesStockUpdates` — none change.
- Reuse `validateCancelStockBasis` — none change.
- Reuse the existing test scaffolding in `stock_repository_cancel_test.go`.

## Source anatomy (confirmed from this repo + DB simulation)
- `sales/repository/stock_repository.go`
  - `cancelStockBasisQuery` (line 270+): three sub-aggregates and one main SELECT rooted on `sls.order_detail`. **Bug 1**: main `Where` at line 345 has `od.item_type = 1` filter; `activeDetailAgg` `Where` at line 307 also has `od.item_type = 1` filter. Both must be removed/replaced.
  - `cancelAgg` builder (line 287-297): `Where("c.cust_id = ? AND c.tr_no = ? AND c.tr_code = 'CO'")`. **Bug 2**: must include legacy `tr_code='SO' tr_no LIKE '%-CO%'` rows.
  - `buildCancelStockMutations` (line 231-268): correct, no change.
  - `CancelSalesStockUpdates` (line 409+): correct, no change.
  - `UpsertWithExistingValueArr` (line 514+): correct, no change.
- `sales/service/order_service.go`
  - `validateCancelTransition` (line 5068+): correct.
  - `validateCancelStockBasis` (line 5097+): correct.
  - `BulkUpdateStatus` cancel branch (line 5131+): correct, no change.
- `sales/repository/stock_repository_cancel_test.go`
  - `TestBuildCancelStockMutations_*`: pass, no change.
  - `TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula`: must be updated to assert the new clauses.
  - `TestGetCancelStockBasisQuery_POQtyFallback`: pass, no change.
  - `TestCancelStockBasisFallback_DetailOnlyPOCase`: pass, no change.

## Reference map
| Feature | Source basis |
| --- | --- |
| Reward-line filter bug | local-DB-backed: simulation of cancel for `SO2606260002` shows `qty_outstanding=4` (only the order line) and `On Cust` projection after cancel = 6 PCS, not 0 |
| Legacy-row risk | local-DB-backed: audit query returned 118 rows matching `tr_code='SO' tr_no LIKE '%-CO%'` across 4 tenants |
| Qty resolution priority | repo-backed: `cancelStockBasisQuery` GREATEST/COALESCE chain (unchanged) |
| Cancel reversal row format | repo-backed: `buildCancelStockMutations` (unchanged) |
| Warehouse stock delta | repo-backed: `UpsertWithExistingValueArr` (unchanged) |
| On Cust Order projection | inventory repo-backed: `inv.stock SUM(qty_in_order) - SUM(qty_out_order)` (unchanged) |
| Wh Stock projection | inventory repo-backed: `inv.stock SUM(qty_in) - SUM(qty_out)` (unchanged) |

## Confirmed vs assumed audit

| Claim | Status | Evidence |
| --- | --- | --- |
| `SO2606260002` exists with `data_status=2` | confirmed_db | `SELECT * FROM sls.order WHERE ro_no='SO2606260002'` |
| Order has 2 detail rows (order + reward) of product 10838 | confirmed_db | `SELECT * FROM sls.order_detail WHERE ro_no='SO2606260002'` |
| 4 `inv.stock` rows exist for `SO2606260002` (2 SO + 2 CO companion) | confirmed_db | stock query |
| `warehouse_stock.qty=14, qty_on_order=10` for the product key | confirmed_db | warehouse query |
| `inv.stock` on_cust_projection = 10 (matches `qty_on_order`) | confirmed_db | SUM aggregate |
| `inv.stock` wh_stock_projection = 2 (matches Wh Stock `2 0 2`) | confirmed_db | SUM aggregate |
| Legacy `tr_code='SO' tr_no LIKE '%-CO%'` rows = 118 across 4 tenants | confirmed_db | legacy audit query |
| `cancelStockBasisQuery` returns `qty_outstanding=4` (only order line) | confirmed_db | full query simulation |
| Simulated cancel writes only 1 reversal row of `qty_out_order=4` | confirmed_db | in-tx simulation, rolled back |
| Simulated cancel leaves `On Cust` projection at 6 (not 0) | confirmed_db | in-tx simulation, rolled back |
| Simulated cancel keeps Wh Stock at 2 | confirmed_db | in-tx simulation, rolled back |
| The fix in commit `20ee58f1` is correct as described in QA's notes | confirmed_repo | QA verification note (file) + Jenkins build #526 |
| Local clone has the same cancelStockBasisQuery code as described | confirmed_repo | direct file read line 270-398 |
| Staging image runs commit `20ee58f1` | unverified | staging container unreachable from this session (docker daemon I/O errors); needs ops confirmation |
| The pre-fix local clone has the bug | confirmed_repo | local code still has `od.item_type = 1` in `cancelStockBasisQuery` at line 307 and 345 |

## Constraints
- Per-item transaction atomicity must be preserved (no change to `BulkUpdateStatus` tx wrapper).
- Tenant filter `cust_id` from auth token stays the only filter on every read/write.
- `CancelSalesStockUpdates` signature, behavior, and SQL are unchanged.
- No new dependency.
- The fix must not regress cancel of single-line orders (no reward), e.g. `SO2606230004` referenced in the original SX-2314 issue.

## Risks
- **R1**: Removing `item_type = 1` from the main query will surface `item_type = 2` (promo/reward) rows that the existing `is_ambiguous` check and `activeDetailAgg` might not handle. Mitigation: also relax `activeDetailAgg` filter; the existing `(qty1_final>0 OR qty1>0 OR qty_po1>0 ...)` predicate is the right gate. The `is_ambiguous` flag is set per `(cust_id, ro_no, pro_id)` and counts all item types — for `SO2606260002` we have one order line + one reward line for the same product, so `is_ambiguous = (count > 1) = true`, and `validateCancelStockBasis` would reject the order with the "inconsistent" error. **This is a new failure mode for SO+reward-of-same-product orders**. We need to either (a) make `activeDetailAgg` count only `item_type = 1` while the main query counts all types, or (b) make the ambiguity check exclude reward items. Pick (a) as the smaller change. Document the decision.
- **R2**: Tightening `cancelAgg` to include legacy rows changes the idempotency math. If the cancel previously ran (in legacy times) and wrote `<SO>-CO` with `tr_code='CO'`, the new query now also picks up the legacy `tr_code='SO' tr_no='<SO>-CO'` row → `qty_outstanding` could go negative. Mitigation: clamp with `GREATEST(..., 0)` (already in the SQL). Verify in the unit test.
- **R3**: For orders that have multiple `<SO>-CO` rows in the new `cancelAgg` (one CO, one legacy SO), the basis query will produce the right outstanding amount but two writes are NOT inserted — only the new reversal row. The legacy row stays. Acceptable; the legacy row represents a previous reversal that was applied.
- **R4**: Local DB is `ggn_scyllax` per AGENTS.md posture. The `host.docker.internal:5432` path is unreachable from this session; we use `127.0.0.1:5432` directly. Confirm with dev/be that this is the same DB the running sales container points to (it is, per `docker inspect scylla-sales` env).

## Decisions / assumptions
- We DO change the code. The simulation proves the local pre-fix code has the bug.
- The fix is the `item_type` filter removal in `cancelStockBasisQuery` (main + `activeDetailAgg`) + the `cancelAgg` legacy filter widening.
- We do NOT backfill or delete legacy `tr_code='SO' tr_no LIKE '%-CO%'` rows. They stay in the ledger; the new query treats them as prior reversals.
- Idempotency on re-cancel is preserved by the `qty_outstanding = qty_out_so - qty_out_order_cancel` arithmetic; the existing `GREATEST(..., 0)` clamp prevents negative.
- `activeDetailAgg` keeps `item_type = 1` for the ambiguity count (decision per R1 mitigation (a)). The main query drops the filter. This is the smallest correct change.

## Execution source of truth
1. Latest explicit user instruction (this prompt).
2. Non-negotiable Implementation Invariants (below).
3. Execution-ready Worklist / Handoff Contract.
4. Acceptance Criteria + Done Criteria.
5. Implementation Steps.
6. Follow-up notes / Out of scope items.

## Non-negotiable implementation invariants
1. `BulkUpdateStatus` cancel branch must remain inside `service.Transaction.WithinTransaction`.
2. `CancelSalesStockUpdates`, `buildCancelStockMutations`, `UpsertWithExistingValueArr` must not change signature or SQL behavior.
3. `validateCancelStockBasis` and `validateFinalOrderStockBasis` must not change.
4. The HTTP contract, payload shape, and tenant filter stay unchanged.
5. `cancelStockBasisQuery` main `Where` must drop `od.item_type = 1` so reward lines (item_type=2) are surfaced. The `activeDetailAgg` subquery keeps `item_type = 1` (per R1 decision).
6. `cancelAgg` `Where` must include both `c.tr_code = 'CO' AND c.tr_no = '<SO>-CO'` AND `c.tr_code = 'SO' AND c.tr_no LIKE '%-CO'` (legacy rows). The combined predicate is `c.tr_no = '<SO>-CO' AND (c.tr_code = 'CO' OR (c.tr_code = 'SO' AND c.tr_no LIKE '%-CO'))`.
7. The execution lane must refresh its own active permissions/context before touching source — planner read-only restrictions do not persist.

## Do not / reject if
- Reverting the simulated cancel evidence (the in-tx simulation must remain rolled back, with results preserved in the evidence file).
- Adding a new dependency or migration.
- Changing the public HTTP contract or the FE.
- Modifying `validateCancelStockBasis` to silently pass `is_ambiguous=true` for order+reward-of-same-product cases. That is the right error to surface to the caller.
- Deleting legacy `tr_code='SO' tr_no LIKE '%-CO%'` rows. They are historical data; we make the query defensive.
- Marking the plan complete without re-running the SQL simulation in the new state and showing `On Cust` projection = 0 after cancel.

## Diff boundary
- Allowed: `sales/repository/stock_repository.go` (two `Where` clause edits), `sales/repository/stock_repository_cancel_test.go` (one test update + two new tests).
- Forbidden: every other file in `scylla-be`, including controllers, migrations, package files, `go.mod`/`go.sum`, docs outside `.opencode/`, and any other service (`pjp`, `inventory`, `master`, etc.).
- Any out-of-boundary change must be reverted or recorded in the slice evidence with explicit justification.

## TDD / Test plan
- **Failing test first**: add `TestGetCancelStockBasisQuery_IncludesRewardLine` and `TestGetCancelStockBasisQuery_LegacySORowsExcludedFromCancelAgg` before the implementation; confirm both fail against current code (the main `Where` still has `item_type=1`, the `cancelAgg` still has only `tr_code='CO'`).
- **Update test**: `TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula` must be updated to assert the new main `Where` (no `item_type=1`) and the new `cancelAgg` `Where` (includes `tr_code='SO' tr_no LIKE '%-CO'` legacy clause). Confirm it fails against current code.
- **Green step**: apply the two `Where` edits in `cancelStockBasisQuery`; rerun the three tests; all green.
- **Regression**: rerun `rtk go test ./repository/... ./service/... -count=1`; all green including `TestBuildCancelStockMutations_*`, `TestBulkUpdateStatus_Cancel_*`, `TestGetCancelStockBasisQuery_POQtyFallback`, `TestCancelStockBasisFallback_DetailOnlyPOCase`.
- **DB-level confirmation**: rerun the in-tx simulation (`/tmp/simulate_cancel.sql`) but with TWO reversal rows (one for `ref_det_id=7540 qty_out_order=4`, one for `ref_det_id=7541 qty_out_order=6`). Confirm `On Cust` projection = 0, Wh Stock projection = 2, `warehouse_stock` delta = `+10 / -10`.
- Commands:
  - `cd sales && rtk go test ./repository/... -run CancelStock -v`
  - `cd sales && rtk go test ./service/... -run BulkUpdateStatus_Cancel -v`
  - `cd sales && rtk go test ./repository/... ./service/... -count=1`
  - `cd sales && rtk go vet ./...`
  - `PGSSL=ignore psql -h 127.0.0.1 -p 5432 -U postgres -d ggn_scyllax -f /tmp/simulate_cancel_v2.sql` (new simulation with two reversal rows)

## Implementation steps
1. **R0** (pre-check): read the local `cancelStockBasisQuery` and confirm the two `item_type = 1` filters at line 307 and 345, and the `cancelAgg` `Where` at line 287-296. Save the raw line numbers to the evidence folder.
2. **T1**: add the two new failing tests and the updated existing test in `stock_repository_cancel_test.go`. Run them; confirm all three fail.
3. **T2**: apply the two `Where` edits in `cancelStockBasisQuery`:
   - Main query (around line 345): remove `AND od.item_type = 1` from the `Where` clause; keep the `qty*_final/qty*/qty_po* > 0` predicate.
   - `activeDetailAgg` (around line 307): keep `od.item_type = 1` (per R1 decision). No change here.
   - `cancelAgg` (around line 287-296): change the `Where` to `c.cust_id = ? AND c.tr_no = ? AND (c.tr_code = 'CO' OR (c.tr_code = 'SO' AND c.tr_no LIKE '%-CO'))`.
4. **T3**: run the three new/updated tests; confirm all green.
5. **T4**: run the full `sales` test suite; all green.
6. **T5**: rerun the in-tx DB simulation (`/tmp/simulate_cancel_v2.sql`) with two reversal rows. Save output to evidence. Confirm:
   - `On Cust` projection = 0.
   - Wh Stock projection = 2.
   - `warehouse_stock.qty` after = 24, `qty_on_order` after = 0.
7. **T6**: rerun the legacy audit query. Save the count; confirm it's unchanged (the legacy rows are not mutated).
8. **T7**: run `rtk go vet ./...`; no warnings.
9. **R1** (`@oracle`): review the two `Where` edits and the new tests; confirm the change is minimal and correct. Save review to `.opencode/draft/20260627-1000-sx2314-cancel-stock-reset/review-oracle.md`.
10. **R2** (`@quality-gate`): final conformance review. Save to `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/review-quality-gate.md`.

## Expected files to change
- `sales/repository/stock_repository.go` (two `Where` clause edits, ~6 LOC).
- `sales/repository/stock_repository_cancel_test.go` (one test update + two new tests, ~80 LOC).
- `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/before-cancel-order.sql.txt` (already saved)
- `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/before-cancel-detail.sql.txt` (already saved)
- `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/before-cancel-stock.sql.txt` (already saved)
- `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/before-cancel-warehouse.sql.txt` (already saved)
- `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/legacy-audit.sql.txt` (already saved, 118 rows)
- `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/simulate-cancel.txt` (already saved, in-tx simulation rolled back)
- `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/simulate-cancel-v2.txt` (new, post-fix simulation with two reversal rows)
- `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/tests-fail.txt` (new, T1 output)
- `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/tests-pass.txt` (new, T3 output)
- `.opencode/draft/20260627-1000-sx2314-cancel-stock-reset/review-oracle.md` (new, R1)
- `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/review-quality-gate.md` (new, R2)

## Agent / Tool routing
- Discovery + DB read (R0, evidence collection): `@fixer` (already done in this plan).
- Implementation (T1-T8): `@fixer`.
- Review (R1): `@oracle`.
- Final signoff (R2): `@quality-gate`.
- No `@designer` work needed.

## Executor handoff prompt
```
You are the implementation lane for SX-2314 v4. Source of truth:
.opencode/plans/20260627-1000-sx2314-cancel-stock-reset.md.

Two defects confirmed by local DB simulation:
1. cancelStockBasisQuery main Where has od.item_type = 1, missing reward
   lines (item_type = 2). Drop the filter so SO+reward-of-same-product
   cancels fully.
2. cancelAgg Where only matches c.tr_code = 'CO'. The local DB has 118
   legacy rows with tr_code='SO' tr_no LIKE '%-CO%'. Widen the filter to
   (c.tr_code='CO' OR (c.tr_code='SO' AND c.tr_no LIKE '%-CO')).

Scope:
- sales/repository/stock_repository.go (two Where edits)
- sales/repository/stock_repository_cancel_test.go (one test update, two new)

Must preserve:
- HTTP contract, payload shape, tenant filter (cust_id)
- per-item transaction atomicity in BulkUpdateStatus
- buildCancelStockMutations, CancelSalesStockUpdates, UpsertWithExistingValueArr
  signatures and behavior
- activeDetailAgg subquery (keep item_type = 1 per R1 decision to avoid
  flagging order+reward-of-same-product as ambiguous)

Do not touch: controllers, migrations, manifests, go.mod/go.sum, FE, other
services, legacy data rows in inv.stock.

TDD order:
1. Add failing tests (T1).
2. Apply edits (T2).
3. Rerun tests; confirm green (T3, T4).
4. Rerun in-tx DB simulation with two reversal rows (T5).
5. Re-audit legacy rows (T6).
6. go vet (T7).
7. Route to @oracle and @quality-gate.

Evidence: save all outputs under
.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/.

Report back to @orchestrator with files changed, test results, DB simulation
output, and the go/no-go for commit.
```

## Execution-ready worklist / handoff contract
1. **R0** | `@fixer` | Read local `cancelStockBasisQuery` and capture the two `item_type=1` line numbers and the `cancelAgg` `Where` to evidence.
   - depends_on: none
   - exit_verification: `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/code-baseline.txt` records the exact line numbers and current SQL.
2. **T1** | `@fixer` | Add `TestGetCancelStockBasisQuery_IncludesRewardLine`, `TestGetCancelStockBasisQuery_LegacySORowsExcludedFromCancelAgg`, and update `TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula` to assert new clauses. Run them; confirm all three FAIL.
   - depends_on: R0
   - exit_verification: `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/tests-fail.txt` shows the three tests failing with assertion mismatch.
3. **T2** | `@fixer` | Edit `cancelStockBasisQuery` main `Where` (drop `od.item_type=1`) and `cancelAgg` `Where` (add legacy `tr_code='SO' tr_no LIKE '%-CO'` clause).
   - depends_on: T1
   - exit_verification: `rtk go build ./...` from `sales/` succeeds.
4. **T3** | `@fixer` | Re-run the three new/updated tests; all green.
   - depends_on: T2
   - exit_verification: `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/tests-pass.txt` shows the three tests passing.
5. **T4** | `@fixer` | Run full `sales` test suite; all green.
   - depends_on: T3
   - exit_verification: `rtk go test ./...` from `sales/` passes.
6. **T5** | `@fixer` | Rerun in-tx DB simulation with two reversal rows (one per `ref_det_id`). Save output.
   - depends_on: T2 (logic change); T4 (regression); can run after T2
   - exit_verification: `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/simulate-cancel-v2.txt` shows `On Cust` projection = 0, Wh Stock = 2, `warehouse_stock` delta = `+10/-10`.
7. **T6** | `@fixer` | Re-run legacy audit query. Row count unchanged.
   - depends_on: T2
   - exit_verification: `legacy-audit.sql.txt` updated timestamp, count still 118.
8. **T7** | `@fixer` | `rtk go vet ./...`; no warnings.
   - depends_on: T3
   - exit_verification: `go-vet.txt` clean.
9. **R1** | `@oracle` | Review the two `Where` edits and the new tests.
   - depends_on: T5, T6
   - exit_verification: `.opencode/draft/20260627-1000-sx2314-cancel-stock-reset/review-oracle.md` with PASS or NEEDS_FIX.
10. **R2** | `@quality-gate` | Final conformance review.
    - depends_on: R1
    - exit_verification: `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/review-quality-gate.md` with PASS or BLOCKED.

- `start_with`: R0 in parallel with nothing.
- All non-blocked tasks are ready.
- `requires_user_decision`: no (decisions are documented in the plan).
- `must_preserve`: per Non-negotiable Implementation Invariants.
- `do_not_touch`: per Diff Boundary.
- `evidence_update`: each task writes its own evidence file under `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/` and updates the task tracker via `python3 ~/.config/opencode/scripts/task-progress.py <task-id> --update <id> --status ...`.

## Validation commands
- `cd sales && rtk go mod download && rtk go mod tidy`
- `cd sales && rtk go test ./repository/... -run CancelStock -v`
- `cd sales && rtk go test ./service/... -run BulkUpdateStatus_Cancel -v`
- `cd sales && rtk go test ./repository/... ./service/... -count=1`
- `cd sales && rtk go test ./...`
- `cd sales && rtk go vet ./...`
- DB simulation (rolls back, safe to re-run):
  - `PGSSL=ignore psql -h 127.0.0.1 -p 5432 -U postgres -d ggn_scyllax -f /tmp/simulate_cancel_v2.sql` (post-fix simulation with two reversal rows).
- Legacy audit (read-only):
  - `PGSSL=ignore psql -h 127.0.0.1 -p 5432 -U postgres -d ggn_scyllax -c "SELECT count(*) FROM inv.stock WHERE tr_code='SO' AND tr_no LIKE '%-CO%';"`

## Evidence requirements
- `discovery.md` (carried over from v1).
- `code-baseline.txt` (R0).
- `before-cancel-*.sql.txt` (carried over).
- `legacy-audit.sql.txt` (carried over, 118 rows).
- `simulate-cancel.txt` (carried over, pre-fix in-tx simulation).
- `simulate-cancel-v2.txt` (T5, post-fix in-tx simulation).
- `tests-fail.txt` (T1).
- `tests-pass.txt` (T3).
- `go-vet.txt` (T7).
- `review-oracle.md` (R1).
- `review-quality-gate.md` (R2).

## Done criteria
- All worklist tasks R0..R2 completed with exit verifications.
- `rtk go test ./...` passes in `sales/`.
- DB post-fix simulation shows `On Cust` projection = 0, Wh Stock projection = 2, `warehouse_stock` delta = `+10/-10` for `SO2606260002`.
- Legacy audit row count unchanged.
- All evidence files present and readable.
- Task tracker fully updated.

## Final planning summary
- **Artifacts consulted**: `prompts/sx_be_doc.txt`, `prompts/sx-2131-docs/be_doc.txt`, `prompts/sx-2291_BE_debug_implementation_prompt.md`, `prompts/SX-1241`, `plans/sx-1241-cancel-order-patch-plan.md`, `sales/repository/stock_repository.go`, `sales/service/order_service.go`, `sales/repository/stock_repository_cancel_test.go`, `inventory/repository/stock_repository.go`, `SX-2314_BE_verification_notes.md` (commit `20ee58f1`), `SX-2314_jenkins_build_verification.md` (build #526, SHA `5c27ec91`).
- **Artifacts created**: `.opencode/evidence/20260627-1000-sx2314-cancel-stock-reset/discovery.md` (v1), `before-cancel-*.sql.txt`, `legacy-audit.sql.txt` (118 rows), `simulate-cancel.txt` (pre-fix in-tx simulation), this plan file v4.
- **Key decisions**: cancel of `SO2606260002` against the pre-fix local code emits only 1 reversal row (order line, `qty_out_order=4`), missing the reward line (`qty_out_order=6`). Net `On Cust` projection after cancel = 6 PCS, not 0 as QA expects. The pre-fix code in the local clone is the pre-`20ee58f1` version. The fix is in `cancelStockBasisQuery`: drop the `od.item_type = 1` filter in the main `Where` (keep it in `activeDetailAgg` to avoid flagging order+reward-of-same-product as ambiguous) and widen the `cancelAgg` `Where` to include legacy `tr_code='SO' tr_no LIKE '%-CO%'` rows.
- **Assumptions**: the dev commit `20ee58f1` already merged into `dev` and shipped in build #526 contains both fixes (per QA verification note and the post-fix code described therein). The local clone here is the pre-fix version; the executor will verify this with the first failing test (T1). The `activeDetailAgg` `item_type = 1` retention is a deliberate R1 mitigation, not an oversight.
- **Open questions**: none blocking. The staging redeploy / live curl to staging is deferred (this clone cannot reach the staging container from this session due to docker daemon I/O errors); the local DB simulation is the primary evidence.
- **Readiness**: `PASS_FOR_SLICE` — bounded two-edit fix with full TDD + DB evidence, no schema/contract change, no FE.
- **Cleanup**: nothing to clean. The legacy `tr_code='SO' tr_no LIKE '%-CO%'` rows stay in the DB; the new query treats them defensively.
- **Active-lane reset**: this plan was authored under `@artifact-planner`'s read-only posture. The execution lane (`@orchestrator` -> `@fixer` / `@oracle` / `@quality-gate`) must refresh its own active permissions/context before touching source files; planner restrictions do not persist.