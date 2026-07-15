# Plan — SX-2314 Cancel Order Stock Reset (v2: self-heal + reward coverage)

Plan Quality Gate: `PASS_FOR_SLICE`
Readiness: `ready-for-slice`
Task id: `20260626-2235-sx-2314-cancel-stock-reset-retry-v2`
Parent plan: `.opencode/plans/20260625-1346-sx-2314-cancel-stock-reset-retry.md`
Prior plan: `.opencode/plans/20260624-0941-sx-2314-cancel-stock-reset.md`
Primary source of truth: `.opencode/plans/20260626-2235-sx-2314-cancel-stock-reset-retry-v2.md`
Branch baseline: `bugfix/SX-2314-dev` (post `20ee58f`)

## Goal
Make PATCH `/v1/orders/status` (`data_status=9`) on `sales` self-heal any partial state, include product-reward lines in the reversal, and stay idempotent under re-cancel. Slice stays narrow: `sales` service only, cancel to `data_status=9` only, sales-rep cancel flow. The retry must close both residual bugs from QA (`SO2606230004`, `SO2606240002`, `SO2606260002`) and the new reward-coverage gap raised in this prompt.

## Why the previous plan (`20260625-1346`) is not enough
The previous retry plan focused on self-heal (drop short-circuit, widen `cancelAgg`, add `ReconcileCancelStockUpdates`). It did not address the new requirement surfaced in this prompt:

1. **Reward lines missing from cancel basis.** `cancelStockBasisQuery` filters `od.item_type = 1` in both the `activeDetailAgg` and the root `WHERE` (`sales/repository/stock_repository.go:307` and `:345`). Reward product lines live on `item_type = 2`. They are therefore skipped in `QtyOutSmallest` math and never get a reversal `inv.stock` row. For `SO2606260002` the `(0 0 6)` reward component of the reversal is dropped, leaving `qty_on_order` overstated by the reward qty even after the order item is reversed.
2. **Plan did not bind to the three QA-reported SO numbers**, so the executor had no concrete regression targets beyond local `ggn_scyllax`.
3. **Plan did not require a staging evidence check** for the legacy `tr_code='SO' tr_no LIKE '%-CO'` rows; the previous quality gate only checked local DB. Production/staging has a different data profile.

The reconcile/skip-`data_status=9` mechanics from the previous plan still apply and are retained below. The reward fix and the staging-legacy check are net-new.

## Non-goals
- No FE/API contract change. Payload stays `{"orders":[{"ro_no":"...","data_status":9}]}`.
- No redesign of the order status machine.
- No change to create/final/invoice stock mutation flows beyond the cancel-specific reward coverage.
- No new dependency, no migration, no cross-service edit.
- No destructive DB writes against shared remote environments.
- No amendment of prior commits; new commit on top of `bugfix/SX-2314-dev` after `20ee58f`.

## Scope

In scope:
1. `sales/repository/stock_repository.go`
   - Widen `cancelStockBasisQuery` to include `item_type IN (1, 2)` in both the `activeDetailAgg` filter and the root `WHERE` so reward lines enter the basis. Keep `QtyOutSmallest` priority chain `qty*_final → qty* → qty_po*` (reward rows have these columns too — confirm during execution via staging `sls.order_detail` dump for the three QA SOs).
   - Widen `cancelAgg` to count legacy reversal rows in addition to canonical `tr_code='CO'`. Final shape:
     `c.cust_id = ? AND c.tr_no = ? AND (c.tr_code = 'CO' OR (c.tr_code = 'SO' AND c.tr_no LIKE '%-CO'))`
     This protects DBs that still hold pre-`20ee58f` legacy rows.
   - Add `ReconcileCancelStockUpdates(c, orderNo, stockDate, basis []entity.CancelStockBasis) error`. For each basis row it looks up existing reversal rows (any `tr_code` whose `tr_no = '<SO>-CO'`) grouped by `(wh_id, pro_id, ref_det_id)`, sums `qty_out_order`, computes `residual = basis.QtyOutSmallest - existingSum`, clamps residual to `>= 0`, inserts a single canonical `tr_code='CO' tr_no='<SO>-CO' qty_out_order=residual` row only when `residual > 0` and no exact existing canonical row matches, and applies `warehouse_stock` delta `Qty=+residual QtyOnOrder=-residual` only when `residual > 0`. Never write a negative residual delta.
   - Keep `CancelSalesStockUpdates` signature intact for any non-cancel caller (currently none, but cheap to retain).
2. `sales/service/order_service.go` `BulkUpdateStatus` cancel branch
   - Drop the early `return nil` for `*orderData.DataStatus == entity.CANCELLED`.
   - Allow `entity.CANCELLED` as a permitted source status for reconcile (no status change needed; stock reconcile only). Keep `validateCancelTransition` rejecting all other unexpected source statuses.
   - Keep the existing `NEED_REVIEW + hasOutstandingStock==false` skip, but re-evaluate outstanding against the new wider basis (which now includes rewards).
   - Add structured audit log per order: `[CANCEL] ro_no=<> current_status=<> basis_rows=<> basis_total_smallest=<> reward_basis_total=<> residuals_applied=<> warehouse_deltas=<>`. `basis_total_smallest` and `reward_basis_total` are summed by `IsAmbiguous` / item type from the basis rows so QA can cross-check the reward component.
3. `sales/entity/stock.go`
   - No change required. `CancelStockBasis` and `CancelStockWrite` are sufficient. Add `Reward bool` or similar only if a future helper needs it; not required for v2.
4. Tests
   - `sales/repository/stock_repository_cancel_test.go`: extend for reward coverage, reconcile behavior, and SQL fragment assertion.
   - `sales/service/order_service_test.go`: extend for reward cancel and already-cancelled reconcile; update mocks to include the new interface method.

Out of scope (next slice):
- Backfill script for historical legacy rows.
- Generic `mst.m_product` lookup refactor.
- Cross-service stock reporting changes.

## Requirements
1. PATCH cancel on a fresh SO writes exactly one `tr_code='CO' tr_no='<SO>-CO' qty_out_order=<outstanding> in smallest unit>` reversal row per `(wh, pro, ref_det)` basis row.
2. Reward lines (`item_type=2`) are included in the basis and the reversal rows. `qty_out_order` for a reward line uses the same priority chain as order items.
3. PATCH cancel on `data_status=9` re-evaluates basis and applies residual warehouse delta via reconcile, never short-circuits.
4. PATCH cancel on an order whose `inv.stock` reversal rows already match the basis (full or partial) inserts only the missing residual reversal row, never duplicates the canonical one, and applies warehouse delta only when residual > 0.
5. PATCH cancel never decreases `warehouse_stock.qty` below zero; residual < 0 is treated as already-corrected (no insert, no warehouse delta).
6. `cancelAgg` SQL fragment includes `c.tr_code = 'CO' OR (c.tr_code = 'SO' AND c.tr_no LIKE '%-CO')` so DBs with legacy reversal rows still reconcile.
7. `cancelStockBasisQuery` accepts `item_type IN (1, 2)` and the activeDetailAgg does the same. `IsAmbiguous` is still set when multiple active detail rows share `(cust_id, ro_no, pro_id)`; the existence of both an order and a reward line on the same product must NOT mark the row ambiguous just because item_type differs (group already keys by `pro_id`, so existing logic is correct — verify with test).
8. Each order keeps per-order `WithinTransaction`. No cross-order transaction.
9. Tenant `cust_id` filter applied on every stock query and write.
10. Final/invoice validation (`validateFinalOrderStockBasis`) stays strict; only cancel path may accept missing source when fallback qty is positive.

## Acceptance Criteria
1. PATCH cancel on `SO2606230004`, `SO2606240002`, `SO2606260002` in staging returns HTTP 200, status flips to `data_status=9`, Warehouse Stock and On Cust Order match QA's expected values (per `evidence/20260624-0941-sx-2314-cancel-stock-reset/discovery.md` and the new prompt). For `SO2606260002` specifically: after cancel, reward `(0 0 6)` contribution is reversed, so `On Cust Order` returns to `0 0 0`.
2. Each `(wh, pro, ref_det)` for these three SOs has at most one `tr_code='CO' tr_no='<SO>-CO'` reversal row in `inv.stock`. No `tr_code='SO' tr_no='<SO>-CO'` row is inserted by the cancel path.
3. Re-cancel on the same SO is a no-op for `inv.stock` and `inv.warehouse_stock` (no second reversal row, no second warehouse delta).
4. Cancel of a Need Review order with positive detail qty and no `inv.stock` SO source row still creates a reversal row from the detail fallback (priority final→sales→PO, smallest unit).
5. `cancelStockBasisQuery` SQL fragment contains both `item_type` values in the activeDetailAgg and root WHERE.
6. `cancelAgg` SQL fragment contains the legacy `tr_code='SO' AND c.tr_no LIKE '%-CO'` clause.
7. `rtk go test ./repository/... ./service/... -count=1` passes; no regression versus the 267-test baseline from `evidence/20260624-0941/quality-gate.md`.
8. Staging evidence file documents the `SELECT tr_code, tr_no, COUNT(*) FROM inv.stock WHERE tr_no LIKE '%-CO' GROUP BY tr_code, tr_no` result.

## Existing Patterns/Reuse
- Keep `entity.CancelStockBasis` and `entity.CancelStockWrite` shapes.
- Reuse `repository.UpsertWithExistingValueArr` for warehouse delta.
- Reuse `repository.StoreBulk` for new reversal rows.
- Reuse `extractTx` for transaction context.
- Reuse `validateCancelTransition` semantics; only the `CANCELLED → CANCELLED` branch is opened for reconcile.
- The existing COALESCE fallback chain in `cancelStockBasisQuery` (`qty1_final → qty1 → qty_po1` etc.) is preserved and applies to reward lines identically.

## Constraints
- `rtk`-prefixed shell per repo `AGENTS.md`.
- `cust_id` filter on every stock query and write; do not loosen tenant safety.
- Source edits remain blocked in planner lane. `@backend` / `@fixer` owns implementation.
- No destructive DB writes against shared remote environments. Use local `ggn_scyllax` for unit-level gorm tests, and read-only queries against staging for evidence.
- No FE/API contract change. Do not touch controllers.
- Do not amend prior commit `20ee58f`; new commit on top.

## Risks
- **Reward `qty*_final` may be empty for some product reward lines.** The COALESCE chain `qty*_final → qty* → qty_po*` will then fall through to `qty_po*`, which may also be empty. Unit tests must cover reward with all-empty qty, reward with only `qty*` filled, and reward with `qty*_final` filled.
- **Legacy `tr_code='SO' tr_no LIKE '%-CO'` rows from pre-`20ee58f` cancels may live in staging.** Widening `cancelAgg` is required to avoid overstating residual, but the wider filter also risks netting out a row that was actually a stale source `SO` row, not a reversal. Mitigation: keep the filter coupled to `tr_no LIKE '%-CO'`, and treat any sum including those rows as "already reversed" — over-cancellation in either direction is safer than under-cancellation here. Staging evidence must include the legacy-row count before deploy.
- **Reward rows in `sls.order_detail` may not have `wh_id`** if the order was created without reward warehouse context. The `qty_outstanding` formula uses `o.wh_id` and `od.pro_id`; if a reward row's wh resolution is wrong the warehouse delta could land in the wrong `inv.warehouse_stock` row. Mitigation: when widen, also confirm reward rows resolve to the same warehouse as their parent order; if not, log and skip warehouse delta for that key.
- **Re-running cancel on a SO whose reward line was partially reversed by an older code path** could leave residual negative. Reconcile path must clamp residual to zero, never apply a negative warehouse delta.

## Decisions/Assumptions
- Decision: include reward lines in the cancel basis via `item_type IN (1, 2)`. Same priority chain, same smallest-unit conversion.
- Decision: keep the `ReconcileCancelStockUpdates` design from the previous plan; it is correct for the residual problem and now also benefits reward coverage.
- Decision: net-legacy reversal rows in `cancelAgg` so staging/prod with pre-`20ee58f` data reconciles without a one-time backfill.
- Decision: Need Review skip is preserved but re-evaluated against the wider basis. If the only positive `QtyOutSmallest` is from reward lines, the skip still fires and no reversal is written — this matches "no source reservation" semantics for Need Review without any reward-impacting source either. Verify with a dedicated test.
- Assumption: reward lines store qty in the same `qty*_final / qty* / qty_po*` columns as order lines. Verify with staging evidence `SELECT order_detail_id, item_type, pro_id, qty, qty_final, qty_po, qty1, qty2, qty3, qty1_final, qty2_final, qty3_final, qty_po1, qty_po2, qty_po3 FROM sls.order_detail WHERE ro_no = 'SO2606260002'`.
- Assumption: per-order `WithinTransaction` from `repository.Dbtransaction` is sufficient; do not collapse the bulk loop into one big transaction.

## Execution Source of Truth
1. Latest explicit user instruction (this prompt).
2. Safety / security / tenant rules in repo `AGENTS.md`.
3. Non-negotiable Implementation Invariants.
4. Acceptance Criteria and Done Criteria.
5. Implementation Steps and Worklist.
6. Follow-up recommendations.

## Non-negotiable Implementation Invariants
- New reversal row still uses `tr_code='CO'`, `tr_no='<SO>-CO'`. Reward and order items both go through the same row shape.
- `warehouse_stock.qty += residual`, `qty_on_order -= residual`. Sign flip on negative residual is not allowed (clamp to 0).
- `cust_id` filter on every stock query and write.
- No FE contract change.
- `CancelSalesStockUpdates` may stay for backward compatibility but cancel path uses `ReconcileCancelStockUpdates` exclusively.
- Reward lines must not be silently excluded again; a test asserts the `item_type` filter accepts 1 and 2.

## Do Not / Reject If
Reject if:
- Negative residual triggers a negative warehouse delta.
- New reversal row inserted when canonical CO row already exists with matching qty.
- Reward `item_type=2` lines are still excluded from the basis or reversal.
- Legacy `SO` reversal rows are still ignored (over-counts residual).
- Loop-level atomic transaction is collapsed into a single transaction across orders.
- Status update is written for an already-cancelled order (no-op OK; unnecessary UPDATE not OK because of trigger side-effects in some DBs — verify with a unit test that asserts `Update` is NOT called when reconcile runs on `data_status=9`).

## Diff Boundary
Allowed files:
- `sales/repository/stock_repository.go`
- `sales/repository/stock_repository_cancel_test.go`
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`
- `sales/entity/stock.go` only if a tiny struct addition is required (not expected).

Allowed evidence:
- `.opencode/evidence/20260626-2235-sx-2314-cancel-stock-reset-retry-v2/`
  - `discovery.md` (already partially in 20260625-1346; reuse/extend)
  - `staging-legacy-rows.md` (new — count of legacy reversal rows in staging)
  - `staging-reward-detail.md` (new — qty columns populated for reward lines in the three QA SOs)
  - `execution.md` (new — after implementation)
  - `validation.md` (new — after validation)
  - `quality-gate.md` (new — after review)
  - `jira-comment.md` (new — draft reply for QA)

## Source Anatomy
- `sales/controller/order_controller.go:40-50` — defines `PATCH /v1/orders/status` route (group `roRouteV1` with `JWTProtected`). `UpdateStatus` (line 894) parses `entity.BulkUpdateStatusOrder` and calls `OrderService.BulkUpdateStatus`.
- `sales/service/order_service.go:5131-5241` — `BulkUpdateStatus` per-order transaction loop with cancel branch. Cancel branch is the target. Short-circuit at lines 5168-5170 must be removed. Helper functions `validateCancelTransition` (5068), `validateCancelStockBasis` (5097), `buildCancelStockWriteCommands` (5077), `cancelOrderStockBasis` (5057) live in the same file.
- `sales/service/order_status_helper.go:200-240` — `determineStatusForExistingOrder` for `data_status=2` path; not on the cancel critical path but adjacent.
- `sales/repository/stock_repository.go:220-445` — `cancelStockBaseRow`, `buildCancelStockMutations`, `cancelStockBasisQuery`, `GetCancelStockBasis`, `CancelSalesStockUpdates`. Reward widen touches the two `item_type = 1` predicates (lines 307, 345); legacy widen touches line 296.
- `sales/repository/stock_repository.go:514-542` — `UpsertWithExistingValueArr` and `UpdateOnCustomerOrder` reused for warehouse delta.
- `sales/model/stock.go` — `model.Stock` row shape (TrCode, TrNo, QtyIn, QtyOut, QtyInOrder, QtyOutOrder, UnitPrice, Cogs, RefDetId).
- `sales/model/warehouse_stock.go` — `model.WarehouseStock` (CustID, WhID, ProID, Qty, QtyOnOrder).
- `sales/model/order_detail.go:7-109` — `OrderDetail` fields; reward lines share the same columns.
- `sales/entity/stock.go:46-69` — `CancelStockBasis` and `CancelStockWrite`. No change required.
- Existing tests: `sales/repository/stock_repository_cancel_test.go` (TestBuildCancelStockMutations_SingleSKU, TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula, TestCancelStockBasisFallback_DetailOnlyPOCase) and `sales/service/order_service_test.go` (TestBulkUpdateStatus_Cancel_*).

## Reference Map
- Repo evidence: `.opencode/evidence/20260624-0941-sx-2314-cancel-stock-reset/{discovery,execution,validation,quality-gate,jira-comment}.md` — `confirmed_repo` for prior patch behavior and `confirmed_repo` for the `20ee58f` test baseline (267 passed in 2 packages).
- Repo evidence: `.opencode/evidence/20260625-1346-sx-2314-cancel-stock-reset-retry/discovery.md` — `confirmed_repo` for the residual bug analysis.
- Jira prompt: `https://scyllax-pratesis.atlassian.net/browse/SX-2314?focusedCommentId=16919` — `user_confirmed` for QA scenario.
- Reference docs (Google Doc): `https://docs.google.com/document/d/1MvflBz2qJWWJIZkua41SUYYUK6BVjIT8OKkZcj17zjc/edit?tab=t.0#heading=h.xrtssomlg36r` — `docs-backed` for the Sales Order Enhancement BE flow. Not opened in this planner pass; treat as `assumption` until executor cross-checks `cancelStockBasisQuery` SQL against the doc.
- First principles: idempotency via residual math, never write a negative delta on cancel, reward lines share the same `sls.order_detail` schema as order items.

## TDD/Test Plan
TDD required. Plan tests before code.

Red step (write these failing tests first):
1. `TestGetCancelStockBasisQuery_IncludesRewardItemType` — basis query dry-run with a synthetic `item_type=2` row, assert the row is included in the result. Fails today because `item_type = 1` excludes it.
2. `TestReconcileCancelStockUpdates_PartialExistingCO_OnlyAppliesResidual` — two basis rows, one has an existing canonical CO reversal of equal qty, the other has none. Expect 1 new CO row + warehouse delta = sum of residuals.
3. `TestReconcileCancelStockUpdates_OverReversed_LeavesWarehouseAlone` — sum of existing reversal rows exceeds basis. Expect no insert, no warehouse mutation.
4. `TestReconcileCancelStockUpdates_RewardResidual` — basis row for a reward product (pro_id X) with no existing CO row. Expect 1 CO row with `tr_code='CO' tr_no='<SO>-CO' qty_out_order=<residual>` and warehouse delta for that product.
5. SQL fragment assertion: `cancelAgg` SQL contains `c.tr_code = 'SO' AND c.tr_no LIKE '%-CO'`. (Extend existing `TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula`.)
6. SQL fragment assertion: `activeDetailAgg` SQL and root WHERE both contain `item_type IN (1, 2)`.
7. `TestBulkUpdateStatus_Cancel_AlreadyCancelled_ReappliesResidualWarehouseDelta` — mock returns `currentStatus=9` and a basis with positive qty (one order, one reward). Assert the new reconcile method is called and `OrderRepository.Update` is NOT called for the status field.
8. `TestBulkUpdateStatus_Cancel_NeedReview_RewardOnlyBasis_AppliesReversal` — currentStatus=1, basis has only reward rows positive, order rows zero. Expect reversal applied for the reward.

Green step:
- Widen `cancelStockBasisQuery` item_type filter to `IN (1, 2)`.
- Widen `cancelAgg` filter to include legacy SO reversal rows.
- Add `ReconcileCancelStockUpdates` to the `StockRepository` interface and implementation.
- Update mocks in `sales/service/order_service_test.go` for the new interface method.
- Drop the short-circuit in `BulkUpdateStatus` cancel branch.
- Allow `CANCELLED → CANCELLED` reconcile path; skip `OrderRepository.Update` when status is already 9.

Refactor (optional):
- Drop `CancelSalesStockUpdates` if no other caller remains. (Verify with grep; if any non-cancel path uses it, keep it.)

Edge cases covered by tests:
- Need Review with no order qty but positive reward qty: apply reversal for reward.
- Already cancelled, no basis rows: no mutation.
- Already cancelled, basis positive, prior CO row matches: no insert, no warehouse mutation.
- Already cancelled, basis positive, prior CO row exists with less qty: insert residual, apply warehouse delta.
- Already cancelled, basis positive, prior CO row exceeds basis (over-reversed): no insert, no warehouse mutation.
- Reward-only basis with over-reversal: same clamping.
- Bulk payload where one SO succeeds and another fails: succeeded one must commit (status + stock) within its own per-order transaction.

## Implementation Steps
1. Read `.opencode/evidence/20260625-1346-sx-2314-cancel-stock-reset-retry/discovery.md` and the `20260624-0941-...` evidence.
2. Read `sales/repository/stock_repository.go`, `sales/service/order_service.go`, `sales/repository/stock_repository_cancel_test.go`, `sales/service/order_service_test.go` in their current state.
3. Capture staging read-only evidence: `SELECT tr_code, tr_no, COUNT(*) FROM inv.stock WHERE tr_no LIKE '%-CO' GROUP BY tr_code, tr_no` and `SELECT order_detail_id, item_type, pro_id, qty, qty_final, qty_po, qty1, qty2, qty3, qty1_final, qty2_final, qty3_final, qty_po1, qty_po2, qty_po3, conv_unit2, conv_unit3 FROM sls.order_detail WHERE ro_no IN ('SO2606230004','SO2606240002','SO2606260002') ORDER BY ro_no, order_detail_id`. Save under `evidence/20260626-2235-.../staging-legacy-rows.md` and `staging-reward-detail.md`.
4. Add the red tests in `sales/repository/stock_repository_cancel_test.go` and `sales/service/order_service_test.go` per the TDD plan. Run the focused suites and confirm red.
5. Open `sales/repository/stock_repository.go`.
6. Widen `cancelStockBasisQuery`:
   - Replace `od.item_type = 1` with `od.item_type IN (1, 2)` in both `activeDetailAgg` (line 307) and the root `WHERE` (line 345).
   - Update `cancelAgg` filter at line 296 to: `c.cust_id = ? AND c.tr_no = ? AND (c.tr_code = 'CO' OR (c.tr_code = 'SO' AND c.tr_no LIKE '%-CO'))`.
7. Add `ReconcileCancelStockUpdates(c, orderNo, stockDate, basis)` per the design. Implementation notes:
   - Query existing reversal rows for `tr_no = orderNo + "-CO"` (any tr_code) grouped by `(wh_id, pro_id, ref_det_id)`, sum `qty_out_order`.
   - For each basis row: `residual = basis.QtyOutSmallest - existingSum`. Clamp `>= 0`.
   - If `residual > 0` and no exact canonical `tr_code='CO' tr_no='<SO>-CO' ref_det_id=basis.RefDetID qty_out_order=basis.QtyOutSmallest` row exists: append a `model.Stock` with `QtyOutOrder=residual`, `QtyIn=0`, `QtyOut=0`, `QtyInOrder=0`, `TrCode="CO"`, `TrNo=orderNo+"-CO"`, `WhID=basis.WhID`, `ProID=basis.ProID`, `ItemCdn=1`, `RefDetId=basis.RefDetID`, `UnitPrice=basis.UnitPrice`, `StockDate=stockDate`, `Cogs=0`.
   - If `residual > 0`: append a `model.WarehouseStock{Qty: residual, QtyOnOrder: -residual, CustID, WhID, ProID}`.
   - Apply `UpsertWithExistingValueArr` and `StoreBulk` at the end.
8. Add `ReconcileCancelStockUpdates` to the `StockRepository` interface (line 26-34).
9. Update mocks in `sales/service/order_service_test.go` (`mockStockRepository`) to include the new method. Update other test mocks as needed (`mockStockRepositoryFinal`, `mockStockRepositoryConcurrency`).
10. Open `sales/service/order_service.go` `BulkUpdateStatus` cancel branch (line 5158 onward).
11. Remove the `if *orderData.DataStatus == entity.CANCELLED { return nil }` short-circuit (lines 5168-5170).
12. Update `validateCancelTransition` (or its caller) to allow `entity.CANCELLED` as a permitted source status for the reconcile path. Add a new helper `validateCancelTransitionAllowReconcile` or pass an option, to keep the original strict variant available for non-cancel callers. Cleanest: change `validateCancelTransition` to accept `entity.CANCELLED` and add an explicit comment.
13. In the cancel branch, when `currentStatus == CANCELLED`:
    - Run `GetCancelStockBasis` and `ReconcileCancelStockUpdates` (with the existing skip rules if `currentStatus == NEED_REVIEW`).
    - Do NOT call `OrderRepository.Update` (status is already 9).
14. For non-cancelled source statuses, keep current logic: basis, validate, reconcile, then `OrderRepository.Update` with the new status.
15. Add structured per-order log: `[CANCEL] ro_no=<> current_status=<> basis_rows=<> basis_total_smallest=<> reward_basis_total=<> residuals_applied=<> warehouse_deltas=<>`.
16. Run `rtk go test ./repository/... ./service/... -run "TestReconcileCancelStockUpdates|TestGetCancelStockBasisQuery|TestBulkUpdateStatus_Cancel" -count=1 -v` and confirm green.
17. Run `rtk go test ./repository/... ./service/... -count=1` and confirm no regression vs the 267-test baseline.
18. Run `rtk go build ./...` to confirm clean build.
19. Local docker validation: bring `sales` up via compose, simulate a partial state on a fresh SO (cancel once, force warehouse back via SQL, re-cancel, assert warehouse matches). Record under `evidence/20260626-2235-.../validation.md`.
20. Commit on `bugfix/SX-2314-dev` (do not amend `20ee58f`).
21. Update Jira comment draft under `evidence/20260626-2235-.../jira-comment.md`.

## Expected Files to Change
- `sales/repository/stock_repository.go`
- `sales/repository/stock_repository_cancel_test.go`
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`

## Agent/Tool Routing
- Implementation owner: `@backend` or `@fixer`.
- Review owner: `@quality-gate`.
- Architecture review: not required unless implementation finds a schema/data mismatch (e.g. reward rows do not have the expected qty columns in staging).
- Staging evidence query: `@backend` or `@explorer` with read-only DB access.

## Execution Ownership Table
| Area | Implementation owner | Review gate |
| --- | --- | --- |
| Reward item_type widen in `cancelStockBasisQuery` | `@backend` | `@quality-gate` |
| Legacy `SO` reversal guard in `cancelAgg` | `@backend` | `@quality-gate` |
| `ReconcileCancelStockUpdates` repository method | `@backend` | `@quality-gate` |
| Service cancel branch reconcile + log | `@backend` | `@quality-gate` |
| Tests (red then green) | `@backend` | `@quality-gate` |
| Staging evidence (legacy rows + reward detail) | `@explorer` | `@quality-gate` |
| Local docker validation | `@backend`/`@explorer` | `@quality-gate` |
| Final review | `@quality-gate` | n/a |

## Executor Handoff Prompt
Execute plan `20260626-2235-sx-2314-cancel-stock-reset-retry-v2` in `/Users/ujang/Projects/Geekgarden/scylla-be`, service `sales`. Build on top of `20ee58f` of branch `bugfix/SX-2314-dev`; do not amend. Goal: make PATCH `/v1/orders/status` cancel self-heal any partial state, include reward lines in the reversal, and stay idempotent on re-cancel. Required changes: widen `cancelStockBasisQuery` to accept `item_type IN (1, 2)` in both `activeDetailAgg` and root WHERE so reward lines enter the basis; widen `cancelAgg` filter to `c.tr_code = 'CO' OR (c.tr_code = 'SO' AND c.tr_no LIKE '%-CO')` so staging/prod databases with legacy reversal rows still reconcile; add `ReconcileCancelStockUpdates(c, orderNo, stockDate, basis []entity.CancelStockBasis) error` that computes residual per `(wh, pro, ref_det)` against existing reversal rows (any tr_code with the same tr_no), inserts only missing residual `tr_code='CO' tr_no='<SO>-CO' qty_out_order=residual` rows, applies warehouse delta `Qty=+residual QtyOnOrder=-residual` only when residual > 0, and clamps negative residual to 0; in `OrderService.BulkUpdateStatus` drop the `if *orderData.DataStatus == entity.CANCELLED { return nil }` short-circuit and route through reconcile; allow `CANCELLED` as a permitted source status for reconcile; keep Need Review `hasOutstandingStock==false` skip; skip `OrderRepository.Update` when status is already 9; add structured per-order audit log. Update mocks and add tests proving partial reconciliation, over-reversal safety, reward coverage, and that the SQL fragments include the new item_type and legacy clauses. Capture staging read-only evidence: legacy reversal-row count and reward-detail qty columns for `SO2606230004`, `SO2606240002`, `SO2606260002`. Validate with `rtk go test ./repository/... ./service/... -count=1` and the focused runs. Return changed files, commit hash, validation outputs, blockers, residual risk. Workers execute only; no replanning.

## Execution-ready Worklist / Handoff Contract
start_with: `R1`

| id | action | depends_on | owner | validation | exit criteria | status | requires_user_decision | must_preserve | do_not_touch | evidence_update | exit_verification |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| R1 | Capture staging read-only evidence: legacy reversal rows + reward detail qty for the three QA SOs | none | `@explorer` | `psql` reads against staging; saved to evidence | two evidence files populated | ready | no | read-only queries | prod env | staging-legacy-rows.md, staging-reward-detail.md | files non-empty with row counts |
| R2 | Add red tests: reward item_type coverage, reconcile (partial, over-reversed, reward residual), SQL fragment assertions, already-cancelled reconcile, Need Review reward-only basis | none | `@backend` | `rtk go test ./repository/... ./service/... -run "TestGetCancelStockBasisQuery_IncludesRewardItemType|TestReconcileCancelStockUpdates|TestBulkUpdateStatus_Cancel" -count=1 -v` | tests fail before code | ready | no | existing test count | non-sales files | note red output | failing assertion proves bug |
| R3 | Widen `cancelStockBasisQuery` item_type filter to `IN (1, 2)` in `activeDetailAgg` and root WHERE | R2 | `@backend` | same as R2 for reward coverage tests | reward coverage tests pass | ready | no | COALESCE fallback chain | unrelated queries | note diff | green output for reward tests |
| R4 | Widen `cancelAgg` filter to include legacy `tr_code='SO' tr_no LIKE '%-CO'` rows | R2 | `@backend` | SQL fragment test | SQL fragment test passes | ready | no | other aggregator scopes | entity structs | note diff | green SQL assertion |
| R5 | Implement `ReconcileCancelStockUpdates` repository method (residual math, clamp to 0, dedupe canonical CO row) | R3, R4 | `@backend` | reconcile tests pass | reconcile tests pass; new CO row only when residual > 0 | ready | no | no negative warehouse delta | controller code | note diff | green reconcile tests |
| R6 | Add `ReconcileCancelStockUpdates` to `StockRepository` interface and update all mocks | R5 | `@backend` | `rtk go build ./...` | clean build | ready | no | back-compat for `CancelSalesStockUpdates` if reused | other service tests | note changes | build success |
| R7 | Drop short-circuit on `data_status=9`; allow CANCELLED as permitted source status; skip `OrderRepository.Update` when status already 9; add per-order audit log | R6 | `@backend` | `rtk go test ./service/... -run TestBulkUpdateStatus_Cancel -count=1 -v` | already-cancelled and Need Review reward-only tests pass | ready | no | per-order transaction boundary | unrelated services | note diff | green service tests |
| R8 | Full test sweep | R7 | `@backend` | `rtk go test ./repository/... ./service/... -count=1` | no regression vs 267-test baseline | ready | no | no source scope creep | unrelated tests | command outputs | final green |
| R9 | Local docker validation: simulate partial state, re-cancel, assert idempotency + reward coverage | R8 | `@backend`/`@explorer` | curl re-cancel + SQL checks | no duplicate CO row; warehouse matches expected | ready | no | endpoint payload | prod env | validation.md | green logs |
| R10 | Quality gate review | R9 | `@quality-gate` | inspect diff + tests + evidence | pass or actionable findings | ready | no | evidence over assertion | implementation edits | quality-gate.md | signoff |

## Validation Commands
Run from `sales` workdir:
1. `rtk go test ./repository/... -run "TestReconcileCancelStockUpdates|TestGetCancelStockBasisQuery" -count=1 -v`
2. `rtk go test ./service/... -run TestBulkUpdateStatus_Cancel -count=1 -v`
3. `rtk go test ./repository/... ./service/... -count=1`
4. `rtk go build ./...`
5. Read-only staging queries (executor captures output in `evidence/.../staging-*.md`):
   - `psql "host=<staging> ..." -c "SELECT tr_code, tr_no, COUNT(*) FROM inv.stock WHERE tr_no LIKE '%-CO' GROUP BY tr_code, tr_no ORDER BY 1;"`
   - `psql "host=<staging> ..." -c "SELECT order_detail_id, item_type, pro_id, qty, qty_final, qty_po, qty1, qty2, qty3, qty1_final, qty2_final, qty3_final, qty_po1, qty_po2, qty_po3, conv_unit2, conv_unit3 FROM sls.order_detail WHERE ro_no IN ('SO2606230004','SO2606240002','SO2606260002') ORDER BY ro_no, order_detail_id;"`
   - `psql "host=<staging> ..." -c "SELECT cust_id, wh_id, pro_id, qty, qty_on_order FROM inv.warehouse_stock WHERE (wh_id, pro_id) IN (SELECT od.wh_id, od.pro_id FROM sls.order_detail od WHERE od.ro_no IN ('SO2606230004','SO2606240002','SO2606260002')) ORDER BY pro_id;"`
6. Local docker re-cancel test (executor captures in `validation.md`): create or use a fresh SO, cancel once, force warehouse back via SQL, re-cancel, assert warehouse matches expected and only one `tr_code='CO' tr_no='<SO>-CO'` row per `(wh, pro, ref_det)`.

## Evidence Requirements
- Red test output before fix (failing assertions prove the bug).
- Green focused test output after fix.
- Staging read-only evidence: `staging-legacy-rows.md` and `staging-reward-detail.md`.
- Local validation note `validation.md` showing simulated partial-state recovery.
- Quality gate signoff `quality-gate.md`.
- Jira comment draft `jira-comment.md` (no tokens).

## Done Criteria
- All acceptance criteria pass.
- All worklist tasks R1–R10 closed.
- Focused + full tests pass; no regressions versus the 267-test baseline.
- Staging evidence captured before implementation; local docker validation captured after.
- Quality gate approves.
- New commit on `bugfix/SX-2314-dev` (on top of `20ee58f`).
- Jira comment draft ready for QA handoff.

## Final Planning Summary

Artifacts consulted:
- `.opencode/plans/20260624-0941-sx-2314-cancel-stock-reset.md`
- `.opencode/plans/20260625-1346-sx-2314-cancel-stock-reset-retry.md`
- `.opencode/evidence/20260624-0941-sx-2314-cancel-stock-reset/{discovery,execution,validation,quality-gate,jira-comment}.md`
- `.opencode/evidence/20260625-1346-sx-2314-cancel-stock-reset-retry/discovery.md`
- `sales/repository/stock_repository.go`
- `sales/service/order_service.go`
- `sales/repository/stock_repository_cancel_test.go`
- `sales/service/order_service_test.go`
- `sales/model/order_detail.go`
- `sales/entity/stock.go`
- QA Jira prompt (this conversation turn).

Artifacts to create:
- This primary plan.
- `.opencode/evidence/20260626-2235-sx-2314-cancel-stock-reset-retry-v2/{staging-legacy-rows.md,staging-reward-detail.md,execution.md,validation.md,quality-gate.md,jira-comment.md}` after execution.

Key decisions:
- Include reward lines in cancel basis via `item_type IN (1, 2)`.
- Self-heal via residual math; never write negative warehouse delta.
- Widen `cancelAgg` to net legacy `tr_code='SO' tr_no LIKE '%-CO'` reversal rows.
- Skip `OrderRepository.Update` when status is already 9 to avoid unnecessary write and side-effects.
- Add per-order structured audit log so QA can cross-check reward vs order basis.

Assumptions:
- Reward lines share `sls.order_detail` qty columns with order lines.
- Per-order `WithinTransaction` is sufficient; do not collapse loop.
- Staging database is reachable for read-only queries from the executor.

Open questions:
- Staging has any legacy `tr_code='SO' tr_no LIKE '%-CO'` rows? Resolved by `R1` (staging evidence capture).
- Reward rows in `sls.order_detail` populate which qty columns for the three QA SOs? Resolved by `R1`.
- Does `OrderRepository.Update` have a DB trigger side-effect on `data_status=9 → 9`? Resolved by `R7` test that asserts Update is NOT called for an already-cancelled order.

Cleanup:
- No new draft artifacts needed; the previous retry's `discovery.md` is reused and the new evidence folder will accumulate execution/validation/quality-gate outputs.
- The previous plan `20260625-1346` becomes superseded by this v2; keep it for history but mark in `Final Planning Summary` of the new evidence files that v2 supersedes it.

Readiness:
- `PASS_FOR_SLICE`. Implementation blocked only because the planner lane cannot edit source files. Hand off to `@backend` / `@fixer` starting with `R1` (staging evidence capture, read-only).
