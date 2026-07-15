# Plan — SX-2314 Cancel Order Stock Reset (Reconcile Follow-up)

Plan Quality Gate: `PASS_FOR_SLICE`
Readiness: `ready-for-slice`
Task id: `20260625-1346-sx-2314-cancel-stock-reset-retry`
Parent plan: `.opencode/plans/20260624-0941-sx-2314-cancel-stock-reset.md`
Primary source of truth: `.opencode/plans/20260625-1346-sx-2314-cancel-stock-reset-retry.md`

## Goal
Resolve QA's failing re-test (`SO2606240002` still shows stale Warehouse Stock and On Cust Order) without rewriting the previous patch. Make the cancel path self-heal against partial-state environments: any previous attempt that inserted an `inv.stock` reversal row but skipped the `warehouse_stock` upsert must converge on the next PATCH. Re-cancel must never double-reverse and must never leave the warehouse stuck. The slice is still narrow: sales service only, cancel to `data_status=9` only.

## Why the previous patch is not enough
The committed patch on `bugfix/SX-2314-dev` (`20ee58f`) fixed the bug for first-time cancel in a clean environment:
- `buildCancelStockMutations` emits one `tr_code='CO' tr_no='<SO>-CO'` reversal row (no duplicate source).
- `cancelAgg` filters prior reversals by `tr_code='CO'`.
- `cancelStockBasisQuery` falls back to detail smallest-unit qty when the SO source ledger is missing.
- `validateCancelStockBasis` allows `IsMissingSource` when `QtyOutSmallest > 0` (detail fallback).

But the service still short-circuits at line 5179-5181:
```go
if *orderData.DataStatus == entity.CANCELLED {
    return nil
}
```
So when QA cancels a SO that was previously flipped to `data_status=9` (e.g. by an older code path, or by an aborted first attempt that updated status but rolled back stock), the new code does not re-evaluate the basis or re-apply any warehouse delta. Staging may also still hold legacy reversal rows with `tr_code='SO' tr_no='<SO>-CO'` from before the patch, which the new `cancelAgg` filter (`tr_code='CO'`) does not count, so the basis `QtyOutSmallest` overstates the residual. The combination of these two facts matches the QA symptom: first cancel in a clean DB works; re-cancel or first-cancel on a partially-corrupted DB leaves warehouse_stock untouched.

## Non-goals
- Do not change FE payload or endpoint contract.
- Do not redesign the order status machine.
- Do not change create/final/invoice stock mutation flows.
- Do not add new dependency, no migration.
- Do not touch services outside `sales`.
- Do not run destructive DB mutations against shared remote environments.
- Do not amend the previous commit; ship this as a follow-up commit on the same branch.

## Scope
In scope:
1. `sales/repository/stock_repository.go`:
   - Widen `cancelAgg` filter to count legacy `tr_code='SO' tr_no='<SO>-CO'` reversal rows alongside new `tr_code='CO'` rows so the residual is correct on databases that still hold legacy rows.
   - Add `ReconcileCancelStockUpdates(c, orderNo, stockDate, basis []entity.CancelStockBasis) error` that computes residual per `(wh_id, pro_id, ref_det_id)` against existing reversal rows (any tr_code with `tr_no='<SO>-CO'`), inserts only the missing residual reversal row, and always applies the residual `warehouse_stock` delta (positive when residual > 0, negative when over-reversed).
   - Add a small helper to look up already-reversed qty per key.
   - Keep `CancelSalesStockUpdates` signature intact for any future caller; new method supersedes its use in the service.
2. `sales/service/order_service.go` `BulkUpdateStatus` cancel branch:
   - Drop the early `return nil` for already-cancelled.
   - Always evaluate basis. For `NEED_REVIEW` keep the `hasOutstandingStock` skip when both source and fallback yield `QtyOutSmallest == 0`. Otherwise call the new `ReconcileCancelStockUpdates`.
   - Keep `validateCancelTransition` guard except that `CANCELLED → CANCELLED` is allowed because we only reconcile stock; no status change is needed. To keep the function semantics simple, allow `currentStatus` to be `NEED_REVIEW`, `PROCESSED`, or `CANCELLED`. If `currentStatus` is anything else, error out as before.
   - Add structured audit log: `[CANCEL] ro_no=<> current_status=<> basis_rows=<> residuals_applied=<> warehouse_deltas=<>`.
3. `sales/entity/stock.go`: no change unless we need a new field. (Not expected.)
4. Tests in `sales/repository/stock_repository_cancel_test.go` and `sales/service/order_service_test.go` for the new reconcile behavior.

Out of scope:
- Backfill / one-time cleanup of historical rows (next slice).
- Cross-service stock reporting changes.

## Requirements
1. PATCH cancel on `data_status=9` re-evaluates basis and applies residual warehouse delta.
2. PATCH cancel never inserts a duplicate reversal row when the canonical `tr_code='CO' tr_no='<SO>-CO' ref_det_id=… qty_out_order=QtyOutSmallest` row already exists.
3. PATCH cancel never decreases `warehouse_stock.qty` below zero unless the basis is negative (which it should not be for cancel).
4. `cancelAgg` counts legacy `tr_code='SO' tr_no='<SO>-CO'` rows when computing residual so prod databases with legacy rows still reconcile correctly.
5. Each order keeps per-order atomic transaction. No shared transaction.
6. Tenant `cust_id` filter applied on every stock query/write.
7. No `warehouse_stock` mutation when residual is zero for that key.

## Acceptance Criteria
1. PATCH cancel on a freshly-cancelled SO that previously had stock writes applied: no new reversal row, no warehouse delta, response 200.
2. PATCH cancel on a freshly-cancelled SO that previously had a reversal `inv.stock` row but no warehouse update: residual warehouse delta is applied (positive qty, negative qty_on_order); no second reversal row inserted.
3. PATCH cancel on a SO whose legacy `inv.stock` reversal rows used `tr_code='SO' tr_no='<SO>-CO'`: residual = basis − sum(legacy SO rows + new CO rows). No duplicate rows; warehouse delta only for true residual.
4. `cancelAgg` SQL fragment includes `c.tr_code = 'CO' OR (c.tr_code = 'SO' AND c.tr_no LIKE '%-CO')`.
5. `rtk go test ./repository/... ./service/... -count=1` passes; no regression in existing 267 tests.

## Existing Patterns/Reuse
- Keep `entity.CancelStockBasis` and `entity.CancelStockWrite` shapes.
- Reuse `repository.UpsertWithExistingValueArr` for warehouse delta.
- Reuse `repository.StoreBulk` for new reversal rows.
- Reuse `extractTx` for transaction context.
- `validateCancelTransition` stays; only allow `CANCELLED` as an extra permitted source status for reconcile.

## Constraints
- Source edits must be handed to `@backend`/`@fixer`; planner lane cannot edit source files.
- Use `rtk go test` per repo `AGENTS.md`.
- No destructive DB writes. Local DB at `host=localhost user=postgres password=postgres dbname=ggn_scyllax sslmode=disable` is acceptable for sanity queries but tests are gorm dry-run.

## Risks
- Wide `cancelAgg` filter that counts legacy `SO` reversal rows could, on a database that has never had a cancellation, miss the legacy rows (no impact because there are none). On a database with mixed rows, it correctly nets them out.
- The reconcile path also fires on `data_status=1 (NEED_REVIEW)`. If the basis is non-zero (fallback qty positive), it applies the reversal even when the SO never reserved stock in `inv.stock`. This matches the docs ("cancel of a Need Review order must still release whatever `qty_on_order` it holds") and the prior patch.
- Residual could be negative if previous attempts double-applied warehouse delta (e.g. legacy bug did warehouse update without basis). The reconcile path must clamp residual to zero or treat negative residual as "skip insert, ignore" so we never run `qty_on_order` positive again.

## Decisions/Assumptions
- Decision: cancel path always re-evaluates; idempotency is enforced by residual math, not by short-circuiting on `data_status=9`.
- Decision: legacy `tr_code='SO' tr_no='<SO>-CO'` reversal rows are counted as already-reversed for residual purposes.
- Decision: residual negative means the warehouse was already corrected (or over-corrected); do not re-apply negative warehouse delta. Insert only if residual > 0 and no exact CO row exists.
- Assumption: per-order `WithinTransaction` from `repository.Dbtransaction` is enough for atomic stock + status mutation; do not collapse loop into one big transaction.

## Execution Source of Truth
1. Latest explicit user instruction and SX-2314 follow-up prompt.
2. Safety/security/tenant rules in repo `AGENTS.md`.
3. Non-negotiable Implementation Invariants.
4. Acceptance Criteria and Done Criteria.
5. Implementation Steps.
6. Follow-up recommendations.

## Non-negotiable Implementation Invariants
- New reversal row still uses `tr_code='CO'`, `tr_no='<SO>-CO'`.
- `warehouse_stock.qty += residual`, `qty_on_order -= residual`. Sign flip on negative residual is not allowed (clamp to 0).
- `cust_id` filter on every stock query/write.
- No FE contract change.

## Do Not / Reject If
Reject if:
- Negative residual triggers a negative warehouse delta.
- New reversal row inserted when canonical CO row already exists with matching qty.
- Legacy `SO` reversal rows are still ignored.
- Loop-level atomic transaction is collapsed into a single transaction across orders.

## Diff Boundary
Allowed files:
- `sales/repository/stock_repository.go`
- `sales/repository/stock_repository_cancel_test.go`
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`
- `sales/entity/stock.go` if a tiny struct addition is required (likely not).

Allowed evidence:
- `.opencode/evidence/20260625-1346-sx-2314-cancel-stock-reset-retry/**`

## TDD/Test Plan
TDD required. Plan tests before code.

Red step:
1. Add `TestReconcileCancelStockUpdates_PartialExistingCO_OnlyAppliesResidual` (gorm dry-run or unit-level) that exercises: basis with two rows, one has an existing CO reversal of equal qty, the other has none. Expect 1 new CO row, warehouse delta = sum of residuals (just the row without an existing CO row).
2. Add `TestReconcileCancelStockUpdates_OverReversed_LeavesWarehouseAlone` where sum of existing reversal rows exceeds basis. Expect no insert, no warehouse mutation.
3. Add SQL fragment assertion: `cancelAgg` SQL contains `tr_code = 'SO'` AND `LIKE '%-CO'`.
4. Add service test `TestBulkUpdateStatus_Cancel_AlreadyCancelled_ReappliesResidualWarehouseDelta` where the mock returns `currentStatus=9` and a basis with positive qty. Assert the new reconcile method is called.

Green step:
- Add `ReconcileCancelStockUpdates` repository method.
- Update `cancelAgg` to include legacy SO reversal rows.
- Update service to call reconcile and stop short-circuiting.

Refactor:
- Optional cleanup of `buildCancelStockMutations` if no longer used elsewhere.

Edge cases:
- Already cancelled, no basis rows (no residual): no mutation.
- Already cancelled, basis qty positive, prior CO row matches: skip insert, no warehouse mutation.
- Already cancelled, basis qty positive, prior CO row exists with less qty: insert residual, apply warehouse delta.
- Already cancelled, basis qty positive, prior CO row exceeds basis (over-reversed): skip insert, no warehouse mutation.
- Need Review with fallback qty positive: behavior unchanged from prior patch.

## Implementation Steps
1. Read `.opencode/evidence/20260624-0941-sx-2314-cancel-stock-reset/{discovery,execution,validation,quality-gate}.md` and `jira-comment.md`.
2. Open `sales/repository/stock_repository_cancel_test.go` and add the new failing tests above.
3. Run `rtk go test ./repository/... -run "TestReconcileCancelStockUpdates|TestGetCancelStockBasisQuery" -count=1 -v` and confirm red.
4. Open `sales/repository/stock_repository.go`.
5. Add `reconcileCancelStockUpdates` method:
   - Query existing reversal rows for `tr_no = '<SO>-CO'` (any `tr_code`), grouped by `(wh_id, pro_id, ref_det_id)`, summing `qty_out_order`.
   - For each basis row, compute `residual = basis.QtyOutSmallest - existingQty` (clamp to >= 0).
   - If `residual > 0` and no exact `tr_code='CO' tr_no='<SO>-CO' ref_det_id=… qty_out_order=basis.QtyOutSmallest` row exists yet: append a `model.Stock` with `qty_out_order=residual`.
   - If `residual > 0`: append a `model.WarehouseStock` delta with `Qty=+residual QtyOnOrder=-residual`.
   - Skip insert and warehouse update when residual is 0 (idempotent).
   - When residual would be negative, treat as 0 (over-reversed; do not write negative warehouse delta).
6. Update `cancelAgg` filter to include legacy `tr_code='SO' tr_no='<SO>-CO'` rows. Suggested pattern:
   - Replace `Where("c.cust_id = ? AND c.tr_no = ? AND c.tr_code = 'CO'", ...)` with
     `Where("c.cust_id = ? AND c.tr_no = ? AND (c.tr_code = 'CO' OR (c.tr_code = 'SO' AND c.tr_no LIKE '%-CO'))", ...)`.
   - Note: `tr_no = '<SO>-CO'` already implies `LIKE '%-CO'`, so the legacy clause is simply `c.tr_code = 'SO'`.
7. Add `ReconcileCancelStockUpdates` to the `StockRepository` interface.
8. Update mocks in `sales/service/order_service_test.go` (`mockStockRepository`) and other test mocks to include the new method.
9. Run repository cancel tests and confirm green.
10. Open `sales/service/order_service.go`.
11. Remove the early `return nil` for `*orderData.DataStatus == entity.CANCELLED`.
12. In `validateCancelTransition` (or its caller), allow `entity.CANCELLED` to also pass (since reconcile only mutates stock, not status, when already cancelled). Add an explicit branch:
    - If `currentStatus == NEED_REVIEW` and `hasOutstandingStock == false`: skip stock write, only status update if needed.
    - Else: call `ReconcileCancelStockUpdates(txCtx, roNo, stockDate, basisRows)`.
13. Update status update path: when already cancelled, do not call `OrderRepository.Update` if status didn't change (avoid unnecessary write). Otherwise keep current update.
14. Add structured log line for cancel reconcile (per order summary) using project log style.
15. Add service test `TestBulkUpdateStatus_Cancel_AlreadyCancelled_ReappliesResidualWarehouseDelta` that mocks the basis and asserts reconcile call.
16. Run `rtk go test ./service/... -run "TestBulkUpdateStatus_Cancel" -count=1 -v` and confirm green.
17. Run `rtk go test ./repository/... ./service/... -count=1` and confirm no regressions.
18. Optionally rerun docker compose local validation per previous evidence workflow; record new evidence under `.opencode/evidence/20260625-1346-sx-2314-cancel-stock-reset-retry/validation.md`.
19. Commit on `bugfix/SX-2314-dev` (do not amend `20ee58f`).
20. Update Jira comment draft under `.opencode/evidence/20260625-1346-sx-2314-cancel-stock-reset-retry/jira-comment.md`.

## Expected Files to Change
- `sales/repository/stock_repository.go`
- `sales/repository/stock_repository_cancel_test.go`
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`

## Agent/Tool Routing
- Implementation owner: `@backend` or `@fixer`.
- Review owner: `@quality-gate`.
- Architecture review: not required unless implementation finds a schema/data mismatch.

## Execution Ownership Table
| Area | Implementation owner | Review gate |
| --- | --- | --- |
| Reconcile repository method | `@backend` | `@quality-gate` |
| CancelAgg legacy guard | `@backend` | `@quality-gate` |
| Service reconcile flow + log | `@backend` | `@quality-gate` |
| Tests | `@backend` | `@quality-gate` |
| Validation | `@explorer`/executor | `@quality-gate` |

## Executor Handoff Prompt
Execute plan `20260625-1346-sx-2314-cancel-stock-reset-retry` in `/Users/ujang/Projects/Geekgarden/scylla-be`, service `sales`. Build on commit `20ee58f` of branch `bugfix/SX-2314-dev`; do not amend. Goal: make PATCH `/v1/orders/status` cancel self-heal any partial state in `inv.stock` or `inv.warehouse_stock`. Required changes: widen `cancelAgg` filter to also count legacy `tr_code='SO' tr_no='<SO>-CO'` reversal rows; add `ReconcileCancelStockUpdates` that computes residual per `(wh, pro, ref_det)` against existing reversal rows (any `tr_code` with the same `tr_no`), inserts only missing residual `tr_code='CO' tr_no='<SO>-CO' qty_out_order=residual` rows, and applies warehouse delta `Qty=+residual QtyOnOrder=-residual` only when residual > 0; clamp negative residual to 0. In `OrderService.BulkUpdateStatus`, remove the `return nil` for `data_status=9` and route through reconcile; allow `CANCELLED` as a permitted source status (reconcile only); keep Need Review `hasOutstandingStock==false` skip. Update mocks and add tests proving partial reconciliation and over-reversal safety. Validate with `rtk go test ./repository/... ./service/... -count=1` and the focused runs. Return changed files, commit hash, validation outputs, blockers, residual risk. Workers execute only; no replanning.

## Execution-ready Worklist / Handoff Contract
start_with: `R1`

| id | action | depends_on | owner | validation | exit criteria | status | requires_user_decision | must_preserve | do_not_touch | evidence_update | exit_verification |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| R1 | Add reconcile + cancelAgg tests (red) | none | `@backend` | `rtk go test ./repository/... -run TestReconcileCancelStockUpdates -count=1` | tests fail before code | ready | no | docs reversal shape | non-sales files | note red output | failing assertion proves bug |
| R2 | Implement `ReconcileCancelStockUpdates` and widen `cancelAgg` | R1 | `@backend` | same as R1 | repository tests pass | ready | no | no negative warehouse delta | FE/API | note diff | green output |
| R3 | Update mocks in service tests for new interface method | R2 | `@backend` | `rtk go build ./...` | build succeeds | ready | no | back-compat for other callers | controller code | note changes | clean build |
| R4 | Drop short-circuit on `data_status=9`; route through reconcile | R3 | `@backend` | `rtk go test ./service/... -run TestBulkUpdateStatus_Cancel -count=1` | service cancel tests pass | ready | no | transaction boundary | unrelated services | note diff | green output |
| R5 | Add already-cancelled reconcile service test | R4 | `@backend` | same as R4 | new test asserts reconcile call | ready | no | idempotent cancel | env/secrets | note test | green output |
| R6 | Full repository + service test sweep | R5 | `@backend` | `rtk go test ./repository/... ./service/... -count=1` | no regression vs 267 prior pass | ready | no | no source scope creep | unrelated tests | command outputs | final green/known failure note |
| R7 | Local docker validation (re-cancel idempotent) | R6 | `@backend`/`@explorer` | curl re-cancel against local sales container; DB check residual | response 200; no duplicate CO row | ready | no | endpoint payload | prod env | validation.md | green logs |
| R8 | Quality gate review | R7 | `@quality-gate` | inspect diff + tests | pass or actionable findings | ready | no | evidence over assertion | implementation edits | review result | signoff |

## Validation Commands
Run from `sales` workdir:
1. `rtk go test ./repository/... -run "TestReconcileCancelStockUpdates|TestGetCancelStockBasisQuery" -count=1 -v`
2. `rtk go test ./service/... -run TestBulkUpdateStatus_Cancel -count=1 -v`
3. `rtk go test ./repository/... ./service/... -count=1`
4. `rtk go build ./...`
5. `PGPASSWORD=postgres psql "host=localhost user=postgres password=postgres dbname=ggn_scyllax sslmode=disable" -c "SELECT tr_code, tr_no, COUNT(*) FROM inv.stock WHERE tr_no LIKE '%-CO' GROUP BY tr_code, tr_no ORDER BY 1;"`
6. `PGPASSWORD=postgres psql "host=localhost user=postgres password=postgres dbname=ggn_scyllax sslmode=disable" -c "SELECT cust_id, wh_id, pro_id, qty, qty_on_order FROM inv.warehouse_stock WHERE (wh_id, pro_id) IN (SELECT wh_id, pro_id FROM sls.order_detail WHERE ro_no='<SO_NO>') ORDER BY pro_id;"`
7. Re-cancel the same order via curl with JWT; assert warehouse_stock deltas only when residual > 0.

## Evidence Requirements
- Red test output before fix.
- Green focused test output after fix.
- Validation note describing the simulated partial state (cancel once; force warehouse back to old value via SQL; re-cancel; check warehouse now matches).
- New Jira comment draft under `.opencode/evidence/20260625-1346-sx-2314-cancel-stock-reset-retry/jira-comment.md`.

## Done Criteria
- All acceptance criteria pass.
- All worklist tasks R1–R8 closed.
- Focused + full tests pass; no regressions.
- QA-gate approves.
- New commit on `bugfix/SX-2314-dev`.

## Final Planning Summary
Artifacts consulted:
- Previous plan and evidence under `.opencode/.../20260624-0941-...`.
- Repo source files: `sales/repository/stock_repository.go`, `sales/service/order_service.go`.
- QA Jira prompt (this conversation turn).

Artifacts to create:
- This primary plan.
- Follow-up evidence under `.opencode/evidence/20260625-1346-sx-2314-cancel-stock-reset-retry/` after execution.

Key decisions:
- Self-heal via residual math instead of short-circuiting on `data_status=9`.
- Count legacy `tr_code='SO' tr_no='<SO>-CO'` rows in residual so prod databases reconcile.
- Clamp negative residual to zero; never write a negative warehouse delta on cancel reconcile.

Assumptions:
- Per-order `WithinTransaction` is sufficient; do not collapse loop into shared transaction.
- Existing `entity.CancelStockBasis` and `entity.CancelStockWrite` shapes are sufficient.

Open questions:
- Whether staging has any legacy `tr_code='SO' tr_no='%SO%CO'` rows. Action item for executor: query staging DB before deploy; if none, the legacy clause is harmless. If yes, the legacy clause reconciles correctly.

Cleanup:
- No draft artifacts needed; reuse the existing discovery/validation evidence and add new follow-up evidence only.

Readiness:
- `PASS_FOR_SLICE`. Implementation blocked only because planner lane cannot edit source files. Hand off to `@backend`/`@fixer`.