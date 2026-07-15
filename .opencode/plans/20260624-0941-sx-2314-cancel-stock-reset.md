# Plan — SX-2314 Cancel Order resets Warehouse Stock and On Cust Order

Plan Quality Gate: `PASS_FOR_SLICE`
Readiness: `ready-for-slice`
Task id: `20260624-0941-sx-2314-cancel-stock-reset`
Primary source of truth: `.opencode/plans/20260624-0941-sx-2314-cancel-stock-reset.md`

## Goal
Fix `PATCH /sales/v1/orders/status` cancel path so canceling a sales order atomically updates `sls.order.data_status`, inserts the correct cancel reversal ledger row, and releases the reserved quantities in `inv.warehouse_stock` for Warehouse Stock and On Cust Order tabs. The first slice is intentionally narrow: sales service only, cancel-to-`data_status=9` only, no FE contract change, no migration, no broad rewrite. Preserve existing transaction boundary and tenant/customer filters. Implement only the smallest code diff needed to align the cancel branch with SX-2314 docs and prior SX-2291 behavior.

## Non-goals
- Do not change FE payload or endpoint contract.
- Do not redesign order status machine.
- Do not change create/final/invoice stock mutation flows except where tests require shared helper expectations.
- Do not add new dependency.
- Do not add migration unless a validation run proves schema mismatch.
- Do not touch services outside `sales`.
- Do not run destructive DB mutations against shared remote environments.

## Scope
In scope:
1. `sales/service/order_service.go` cancel branch in `BulkUpdateStatus`.
2. `sales/repository/stock_repository.go` cancel basis and cancel write functions.
3. Existing cancel-focused tests in `sales/repository/stock_repository_cancel_test.go` and `sales/service/order_service_test.go`.
4. New or changed tests that prove first-time cancel, Need Review cancel, PO qty fallback, and idempotent retry.

Out of scope next slice:
- Broader reconciliation job for already-corrupted historical orders. Promote only if QA finds existing production rows with double reversal or missing reversal after code fix.
- FE screenshot verification. Promote only after BE deploy to QA.
- Cross-service stock reporting changes.

## Requirements
1. Cancel request `{"orders":[{"ro_no":"<SO_NO>","data_status":9}]}` must keep same API contract.
2. Each order update must remain inside existing transaction.
3. When canceling an eligible order, update `sls.order.data_status` to 9.
4. Insert one cancel reversal row per order detail/product/warehouse basis with `tr_code='CO'` and `tr_no='<SO_NO>-CO'`.
5. Cancel reversal row must set `qty_in=0`, `qty_out=0`, `qty_in_order=0`, and `qty_out_order=<cancelQtySmallest>`.
6. `inv.warehouse_stock.qty` must increase by `<cancelQtySmallest>`.
7. `inv.warehouse_stock.qty_on_order` must decrease by `<cancelQtySmallest>`.
8. Qty source priority must be final order (`qty1_final/qty2_final/qty3_final`), sales order (`qty1/qty2/qty3`), purchase order (`qty_po1/qty_po2/qty_po3`).
9. Qty conversion must produce smallest-unit qty using existing conversion convention (`qty1 * conv_unit2 * conv_unit3 + qty2 * conv_unit3 + qty3`, with safe fallback for null/zero conv values using project pattern).
10. Need Review (1) → Cancelled (9) must not skip stock release when either stock ledger basis or detail qty basis proves an outstanding reservation.
11. Processed (2) → Cancelled (9) must keep existing happy path.
12. Already-cancelled orders must be idempotent: no duplicate reversal row, no duplicate warehouse_stock delta.
13. Bulk payload must keep existing transaction behavior per item; do not silently alter error semantics unless tests show existing behavior.
14. Tenant safety must use `cust_id` in all reads/writes.
15. Do not overwrite unrelated `warehouse_stock` columns.

## Acceptance Criteria
1. Canceling `SO2606230004`-style order creates `inv.stock` row where `tr_code='CO'` and `tr_no='SO2606230004-CO'`.
2. Reversal row has `qty_out_order` equal to cancel qty in smallest unit.
3. `inv.warehouse_stock.qty` increases by cancel qty.
4. `inv.warehouse_stock.qty_on_order` decreases by cancel qty.
5. `sls.order.data_status` becomes `9`.
6. Need Review cancel with outstanding reservation releases stock/order values.
7. PO-only qty detail resolves quantity from `qty_po*` fields.
8. Repeating cancel does not insert duplicate `CO` rows and does not apply duplicate `warehouse_stock` delta.
9. Existing cancel service tests still pass after expected assertion updates.
10. `rtk go test ./repository/...` and `rtk go test ./service/...` pass for cancel-focused suites.

## Existing Patterns/Reuse
Reuse:
- `orderServiceImpl.BulkUpdateStatus` transaction wrapper.
- `validateCancelTransition` status guard.
- `GetCancelStockBasis` and `CancelSalesStockUpdates` repository boundary.
- `UpsertWithExistingValueArr` delta semantics for `warehouse_stock`.
- Existing tests as regression anchors.

Do not create:
- New repository interface if existing method can be adjusted.
- New package for conversion if SQL expression suffices.
- New status enum mapping.

## Constraints
- Active repo: `/Users/ujang/Projects/Geekgarden/scylla-be`.
- Service module: `/Users/ujang/Projects/Geekgarden/scylla-be/sales`.
- Source edit currently blocked in planner lane; implementation must be handed to `@backend`/`@fixer`.
- Project-local stack docs named in global workflow were absent; repo `AGENTS.md`, service layout, source files, and tests were used as local evidence.
- Runtime DB did not contain `SO2606230004`, so DB repro not available locally.

## Risks
- Existing tests currently encode wrong `tr_code='SO'` for cancel reversal row; they must be updated with docs-backed expected behavior.
- Detail-qty fallback can mutate Need Review orders that never reserved stock. Mitigation: only fallback when current code would otherwise skip but detail qty is positive and there is no prior `CO` reversal; keep idempotency guard.
- Historical data may already contain reversal rows with `tr_code='SO'` and `tr_no='<SO>-CO'`. Mitigation: idempotency query should consider both old legacy rows and new correct `CO` rows when preventing double reversal, but new writes must use `CO`.

## Decisions/Assumptions
- Decision: New reversal rows must use `tr_code='CO'`, matching SX-2314 docs.
- Decision: Do not insert duplicate source `SO` row on cancel; cancel should only insert reversal ledger row.
- Decision: `warehouse_stock` delta remains `qty=+cancelQty`, `qty_on_order=-cancelQty`.
- Assumption: `conv_unit2` and `conv_unit3` are sufficient for Large/Middle/Small conversion in `sales` current data shape.
- Assumption: Existing item type filter `od.item_type = 1` remains correct for normal SKU stock release.
- Open question: Whether legacy `tr_code='SO'` cancel rows should be migrated. Not needed for first slice.

## Execution Source of Truth
Precedence:
1. Latest explicit user instruction and SX-2314 docs in prompt.
2. Safety/security/tenant rules in repo `AGENTS.md`.
3. Non-negotiable Implementation Invariants.
4. Acceptance Criteria and Done Criteria.
5. Implementation Steps.
6. Follow-up recommendations.

If conflict exists, executor follows higher source and records conflict in verification evidence.

## Non-negotiable Implementation Invariants
- No FE contract change.
- No source edits outside `sales` unless executor proves required.
- All cancel stock changes stay in same transaction as status update.
- Every stock query/write filters by `cust_id`.
- New reversal row uses `tr_code='CO'`, `tr_no='<SO>-CO'`.
- Idempotency must cover new `CO` rows and legacy `<SO>-CO` rows if present.
- Do not silently skip stock release when outstanding basis exists.

## Do Not / Reject If
Reject implementation if:
- It inserts duplicate original `SO` ledger row during cancel.
- It decreases `qty` instead of increasing it on cancel.
- It increases `qty_on_order` on cancel.
- It changes endpoint payload.
- It adds a new dependency.
- It removes tenant filter.
- It passes tests only by weakening assertions.

## Diff Boundary
Allowed source/test files:
- `sales/repository/stock_repository.go`
- `sales/repository/stock_repository_cancel_test.go`
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`
- Minimal adjacent test helper files only if compiler requires.

Allowed evidence files:
- `.opencode/evidence/20260624-0941-sx-2314-cancel-stock-reset/**`

Anything else must be reverted or justified in verification notes.

## TDD/Test Plan
TDD required: yes. This is production logic affecting stock and tenant data.

Red step:
1. Update `TestBuildCancelStockMutations_SingleSKU` to expect one stock row, `tr_code='CO'`, `tr_no='<SO>-CO'`, correct qty fields, correct warehouse delta.
2. Update `TestBuildCancelStockMutations_MultiSKUAndIdempotent` to expect one reversal row per SKU, not two rows.
3. Add SQL dry-run assertion that cancel aggregate detects `tr_code = 'CO'` and/or legacy fallback if implemented.
4. Add repository test proving detail qty fallback resolves PO qty to smallest unit if no source ledger exists.
5. Update service Need Review test so positive fallback basis calls `CancelSalesStockUpdates`, while truly zero basis still skips.

Green step:
- Patch repository and service with minimal conditional changes.

Refactor step:
- Remove dead `SourceTrCode` if unused only if no broad ripple; otherwise leave it.

Edge cases:
- Already Cancelled.
- Need Review with no source and no qty.
- Need Review with source or detail fallback qty.
- Processed with source stock row.
- Legacy prior reversal row using old `SO` code.
- PO-only qty.

## Implementation Steps
1. Open `sales/repository/stock_repository_cancel_test.go`.
2. Change single-SKU cancel mutation expected stock row count from 2 to 1.
3. Change expected reversal row identity to `tr_code='CO'`, `tr_no='<SO>-CO'`.
4. Remove rowA duplicate-source expectation.
5. Keep qty assertion for reversal row: `qty_in=0`, `qty_out=0`, `qty_in_order=0`, `qty_out_order=<qty>`.
6. Keep warehouse delta assertion: `Qty=<qty>`, `QtyOnOrder=-<qty>`.
7. Change multi-SKU expected stock row count from 4 to 2.
8. Change partial reverse test expected stock row count from 2 to 1.
9. Add dry-run SQL assertion for `c.tr_code = 'CO'` or legacy-safe predicate.
10. Add dry-run SQL assertion for smallest-unit fallback expression from `qty*_final`, `qty*`, `qty_po*`.
11. Run repository cancel tests and confirm red.
12. Open `sales/repository/stock_repository.go`.
13. In `buildCancelStockMutations`, remove creation of duplicate source `stockRowA`.
14. Change cancel reversal stock row `TrCode` to `CO`.
15. Keep `TrNo` as `orderNo + "-CO"`.
16. Keep reversal qty fields per docs.
17. Append only the one reversal stock row.
18. In `cancelAgg`, change prior cancel filter to find `tr_code='CO'` for `tr_no='<SO>-CO'`.
19. If legacy data matters, use `AND c.tr_code IN ('CO','SO')` only for prior-cancel detection, but new insert must be `CO`.
20. Ensure `qty_out_order_cancel` sums only reversal rows.
21. Build SQL expression for detail qty fallback in smallest unit.
22. Use final qty priority first: any final qty present > 0 means use final set.
23. Else use sales qty set when any sales qty present > 0.
24. Else use PO qty set when any PO qty present > 0.
25. Convert qty set to smallest unit with conv unit fields.
26. Use source ledger qty when present and positive.
27. Use detail fallback qty when source ledger missing but detail qty positive.
28. Keep `is_missing_source` true for observability if source ledger absent.
29. Set `qty_outstanding = cancel_basis_qty - qty_out_order_cancel`.
30. Set `qty_out_smallest = GREATEST(qty_outstanding, 0)`.
31. Preserve ambiguity guard for multiple active details per product unless tests prove too strict.
32. Run repository cancel tests and confirm green.
33. Open `sales/service/order_service_test.go`.
34. Add/adjust Need Review test with positive fallback basis so `CancelSalesStockUpdates` receives rows.
35. Keep a zero basis test where status updates without stock write if both source and detail qty are zero.
36. Run service cancel tests and confirm red/green as appropriate.
37. Open `sales/service/order_service.go` only if service skip logic still blocks valid fallback rows.
38. For Need Review, detect outstanding basis by `row.QtyOutSmallest > 0`, not by source-only signal.
39. Do not skip stock write when fallback qty produces `QtyOutSmallest > 0`.
40. Keep validation errors for ambiguous or invalid outstanding.
41. Ensure already cancelled path returns before stock write.
42. Add concise log around cancel basis and warehouse delta only if project logging pattern is present; avoid noisy logs in loops unless requested.
43. Run `rtk go test ./repository/... -run "TestBuildCancelStockMutations|TestGetCancelStockBasisQuery" -count=1 -v`.
44. Run `rtk go test ./service/... -run "TestBulkUpdateStatus_Cancel" -count=1 -v`.
45. Run `rtk go test ./repository/... ./service/... -count=1` if focused tests pass.
46. If compile errors from unused struct fields occur, remove fields minimally.
47. Run DB dry-run queries only against local DB; sample order missing locally is acceptable evidence.
48. Prepare verification note with changed behavior and tests.
49. Route to `@quality-gate` because stock/accounting data mutation is material.
50. Do not deploy without QA DB check for a real SO.

## Expected Files to Change
- `sales/repository/stock_repository.go`
- `sales/repository/stock_repository_cancel_test.go`
- `sales/service/order_service.go` only if skip logic needs adjustment.
- `sales/service/order_service_test.go`

## Agent/Tool Routing
- Implementation owner: `@backend` or `@fixer`.
- Review owner: `@quality-gate`.
- Architecture review: not required for first slice unless implementation discovers schema mismatch.
- DB evidence: executor may use local `psql`; remote QA needs user-provided token/approval.

## Execution Ownership Table
| Area | Implementation owner | Review gate |
| --- | --- | --- |
| Cancel repository basis/write | `@backend` | `@quality-gate` |
| Cancel service skip/transaction behavior | `@backend` | `@quality-gate` |
| Tests | `@backend` | `@quality-gate` |
| QA DB verification notes | `@explorer`/executor | `@quality-gate` |

## Executor Handoff Prompt
Execute plan `20260624-0941-sx-2314-cancel-stock-reset` in `/Users/ujang/Projects/Geekgarden/scylla-be`, service `sales`. Fix only SX-2314 cancel stock reset. Must preserve API payload, tenant `cust_id` filters, existing transaction boundary, and status update behavior. New cancel reversal rows must use `tr_code='CO'`, `tr_no='<SO>-CO'`, `qty_in=0`, `qty_out=0`, `qty_in_order=0`, `qty_out_order=<cancelQtySmallest>`. Warehouse delta must be `qty += cancelQtySmallest`, `qty_on_order -= cancelQtySmallest`. Do not touch FE, migrations, other services, env files, or unrelated tests. Start with failing tests, then minimal code. Return changed files, commands run, outputs, and any DB evidence. Workers execute only and report back; no replanning unless schema/test evidence contradicts this plan.

## Execution-ready Worklist / Handoff Contract
start_with: `T1`

| id | action | depends_on | owner | validation | exit criteria | status | requires_user_decision | must_preserve | do_not_touch | evidence_update | exit_verification |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T1 | Update cancel repository tests to docs-backed reversal shape | none | `@backend` | `rtk go test ./repository/... -run TestBuildCancelStockMutations -count=1` | tests fail before code | ready | no | docs reversal shape | non-sales files | note red output | failing assertion proves bug |
| T2 | Patch `buildCancelStockMutations` to emit one `CO` reversal row and correct warehouse delta | T1 | `@backend` | same as T1 | repository mutation tests pass | ready | no | no duplicate source row | service code | note diff | green output |
| T3 | Patch `cancelAgg` idempotency to count `CO` reversal rows, with legacy guard if needed | T2 | `@backend` | `rtk go test ./repository/... -run TestGetCancelStockBasisQuery -count=1` | SQL asserts pass | ready | no | idempotent cancel | FE/API | note SQL fragments | green output |
| T4 | Add detail qty fallback basis test | T3 | `@backend` | repository dry-run test | fallback expression present | ready | no | final→sales→PO priority | unrelated SQL | note test name | green output |
| T5 | Adjust service Need Review tests if current skip blocks fallback | T4 | `@backend` | `rtk go test ./service/... -run TestBulkUpdateStatus_Cancel -count=1` | service cancel tests pass | ready | no | transaction boundary | controllers | note behavior | green output |
| T6 | Run focused and package tests | T5 | `@backend` | listed validation commands | no failing focused tests | ready | no | no source scope creep | env/secrets | command outputs | final green/known failure note |
| T7 | Quality gate review | T6 | `@quality-gate` | inspect diff + tests | pass or actionable findings | ready | no | evidence over assertion | implementation edits | review result | signoff |

## Validation Commands
Run from `/Users/ujang/Projects/Geekgarden/scylla-be/sales`:
1. `rtk go test ./repository/... -run TestBuildCancelStockMutations -count=1 -v`
2. `rtk go test ./repository/... -run TestGetCancelStockBasisQuery -count=1 -v`
3. `rtk go test ./repository/... -run "TestBuildCancelStockMutations|TestGetCancelStockBasisQuery" -count=1 -v`
4. `rtk go test ./service/... -run TestBulkUpdateStatus_Cancel -count=1 -v`
5. `rtk go test ./repository/... ./service/... -count=1`
6. `rtk go test ./... -count=1` if time permits.
7. `PGPASSWORD=postgres psql "host=localhost user=postgres dbname=ggn_scyllax sslmode=disable" -c "SELECT 1"`
8. `PGPASSWORD=postgres psql "host=localhost user=postgres dbname=ggn_scyllax sslmode=disable" -c "SELECT ro_no, cust_id, wh_id, data_status FROM sls.\"order\" WHERE ro_no='SO2606230004';"`
9. `PGPASSWORD=postgres psql "host=localhost user=postgres dbname=ggn_scyllax sslmode=disable" -c "SELECT tr_code, tr_no, wh_id, pro_id, qty_in, qty_out, qty_in_order, qty_out_order FROM inv.stock WHERE tr_no IN ('SO2606230004','SO2606230004-CO') ORDER BY stock_id;"`
10. `PGPASSWORD=postgres psql "host=localhost user=postgres dbname=ggn_scyllax sslmode=disable" -c "SELECT cust_id, wh_id, pro_id, qty, qty_on_order FROM inv.warehouse_stock WHERE (wh_id, pro_id) IN (SELECT wh_id, pro_id FROM sls.order_detail WHERE ro_no='SO2606230004');"`

Expected: focused Go tests pass. Local DB sample may return zero rows; record that as environment limitation, not failure.

## Evidence Requirements
- Red test output before fix.
- Green focused test output after fix.
- `git diff` or equivalent changed-file summary.
- DB query output if real SO exists in local/QA DB.
- If sample order absent locally, state `SO2606230004 absent in local ggn_scyllax`.

## Done Criteria
- Acceptance criteria met in tests.
- Focused repository and service tests pass.
- No unrelated files changed.
- `@quality-gate` approves or findings resolved.
- QA can verify Warehouse Stock and On Cust Order reset after cancel.

## Final Planning Summary
Artifacts consulted/created:
- Consulted `AGENTS.md`, `sales/service/order_service.go`, `sales/repository/stock_repository.go`, cancel tests.
- Created `.opencode/evidence/20260624-0941-sx-2314-cancel-stock-reset/discovery.md`.
- Created this primary plan.

Key decisions:
- Plan first slice only; no broad rewrite.
- New reversal rows should be `CO` per SX-2314 docs.
- `warehouse_stock` delta direction already correct; fix row shape/idempotency and fallback path.

Assumptions:
- Detail qty fallback can be expressed in existing SQL without new dependency.
- Legacy `<SO>-CO` rows may exist; idempotency should consider them if easy.

Open questions:
- Whether historical rows need one-time remediation. Not blocking first slice.

Cleanup:
- Kept discovery evidence because it names exact files/lines and baseline tests.
- No stale draft artifacts kept.

Readiness:
- `PASS_FOR_SLICE`. Implementation blocked only by planner permission boundary; handoff ready for `@backend`/`@fixer`.
