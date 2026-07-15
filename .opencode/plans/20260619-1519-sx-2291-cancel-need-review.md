# Plan — SX-2291 Cancel Sales Order Need Review from PO Tab

## Goal

Allow `PATCH /sales/v1/orders/status` to cancel Sales Orders that are currently in `Need Review` (`data_status = 1`), including orders shown on the **Purchase Order** tab where stock basis may only contain `qty_po1/qty_po2/qty_po3`. Status must reach `Cancelled` (`data_status = 9`) and stock reversal must follow docs priority: `qty*_final` → `qty*` → `qty_po*`.

## Non-goals

- FE contract changes. Payload stays `{"orders":[{"ro_no":"...","data_status":9}]}`.
- No remap of `Need Review` to `Processed`.
- No new endpoint, no new auth, no new fields.
- No changes to `pjp*` services.
- No changes to `inventory` or `master` service code.

## Scope

In scope, limited to `sales` service:
- Cancel transition for `Need Review` orders (already allowed in code; verify, do not regress).
- Cancel stock basis query and helper to resolve qty for orders lacking `qty_final`.
- Service-level guard so a `Need Review` order with no stock movement cancels safely without false-negative "inconsistent basis" error.
- Idempotency guard to avoid duplicate `SO-CO` rows when order already cancelled.
- Tests: add regression tests for Need Review + PO qty only, Need Review + no stock movement, Need Review + already cancelled.

## Requirements

- `PATCH /sales/v1/orders/status` accepts `data_status = 9` for current `data_status = 1`.
- `sls.order.data_status` updates to `9`.
- For PO-tab `Need Review` orders:
  - Cancel must work when only `qty_po1/qty_po2/qty_po3` are populated.
  - Cancel must work when no `inv.stock` row exists.
- Stock reversal priority:
  1. `qty1_final`, `qty2_final`, `qty3_final`
  2. `qty1`, `qty2`, `qty3`
  3. `qty_po1`, `qty_po2`, `qty_po3`
- No duplicate cancel stock rows when re-cancelling an already-cancelled order.
- Tenant (`cust_id`) from JWT must still scope every read/write. Do not introduce broader updates.

## Acceptance Criteria

- Cancel `Need Review` order succeeds with response shape identical to existing success response.
- After success, `SELECT data_status FROM sls.order WHERE ro_no = 'SO2606190005'` returns `9`.
- For PO-tab order with only `qty_po*` populated, cancel writes a single new `inv.stock` row with `tr_code = 'SO'`, `tr_no = <SO_NO>-CO` (per existing `buildCancelStockMutations`) and updates `inv.warehouse_stock` accordingly.
- For `Need Review` order with no `inv.stock` row and no `qty_final`, cancel updates `sls.order.data_status = 9` and does not raise the "inconsistent basis" error.
- Re-cancelling an order already at `data_status = 9` is a no-op and does not write a second `SO-CO` row.
- All other cancelable statuses (`Processed`, etc.) still cancel with identical stock behavior.

## Existing Patterns/Reuse

- Transition validator: `sales/service/order_service.go:5002-5008` `validateCancelTransition` already includes `entity.NEED_REVIEW`. Keep, do not duplicate.
- Status constants: `sales/entity/order.go:3-26`. Reuse `entity.NEED_REVIEW`, `entity.CANCELLED`.
- Cancel basis query: `sales/repository/stock_repository.go:286-364` `cancelStockBasisQuery`. Reuse and extend filter logic.
- Cancel stock writer: `sales/repository/stock_repository.go:231-283` `buildCancelStockMutations` and `:376-411` `CancelSalesStockUpdates`. Reuse.
- Service tx wrapper: `sales/repository/dbtransaction.go:27-55`.
- Order update: `sales/repository/order_repository.go:435-443` `Update` is tenant-scoped.
- Existing qty priority helper for taking/display: `sales/service/order_type_helper.go:56-61` `takingOrderQtySource` (`qty*` → `qty_po*`). Extend or wrap, do not duplicate ad-hoc.
- Existing unit tests live in `sales/service/order_service_test.go` and `sales/repository/stock_repository_cancel_test.go`. Follow `testify` patterns.

## Constraints

- Backend bug; no UI work.
- Multi-module Go monorepo. Per-service `go.mod`. Validate only inside `sales/`.
- Postgres on `host.docker.internal:5432` for local validation per repo AGENTS.
- Do not commit secrets; do not edit `docker-compose.yml` creds.
- No new external dependency; reuse `gorm` and existing `testify`.
- Keep error strings stable enough to not break FE error mapping.

## Risks

- Risk: changing `cancelStockBasisQuery` filter may include rows that should not generate stock reversal. Mitigation: keep `item_type = 1` filter, reuse `qty_out_so` formula, add new unit test for PO-only basis row.
- Risk: silent skip of stock reversal for `Need Review` could mask data issues. Mitigation: emit existing log line and add a unit test that asserts no `inv.stock` row is written and status still updates.
- Risk: idempotency check relying on `data_status == 9` early return still leaves a partial transaction if a previous attempt wrote `SO-CO` but failed before status update. Mitigation: detect existing `SO-CO` row in basis query and treat as already cancelled.
- Risk: `qty_po*` is nullable. Mitigation: coalesce to 0 and skip rows with all zero qty.
- Risk: scope creep into other status transitions. Mitigation: only edit cancel branches.

## Decisions/Assumptions

- Decision: keep `validateCancelTransition` as-is. It already allows `Need Review`. No change to transition list.
- Decision: introduce a single qty-resolver helper in `sales/repository/stock_repository.go` (or a private helper in `sales/service/order_service.go`) that maps one `sls.order_detail` row to `{qty1, qty2, qty3}` using priority `qty*_final` → `qty*` → `qty_po*`.
- Decision: extend `cancelStockBasisQuery` filter from `COALESCE(od.qty_final, 0) > 0` to a coalesced priority expression that is true when any of `qty_final`, `qty`, `qty_po` is non-zero. Keep `item_type = 1` and tenant scope.
- Decision: in `BulkUpdateStatus`, when target is `CANCELLED` and basis returns zero rows for a `Need Review` order, log a warning and proceed to update `sls.order.data_status = 9` without writing stock rows. This mirrors the "no stock movement" hypothesis and avoids false-negative "inconsistent basis" errors.
- Decision: idempotency. `CancelSalesStockUpdates` already checks existing cancel rows via `GetCancelStockBasis` (uses `tr_no = '<SO_NO>-CO'`). Re-cancel still calls this and writes nothing if already present. Keep the early-return in service for `data_status == 9`.
- Assumption: parent `SX-162` does not impose stricter rules; treating this as a regression bug not a product re-think.
- Assumption: bulk payload stays array; one failure rolls back the whole request per current `WithinTransaction` semantics. Keep that behavior.
- Assumption: existing `cust_id` JWT local is sufficient scope; no `parent_cust_id` required for this endpoint.

## Execution Source of Truth

Precedence for implementation:
1. Latest explicit user instruction in SX-2291 prompt.
2. Safety, tenancy, and security rules in `AGENTS.md` and `.opencode/docs/ARCHITECTURE.md`.
3. Non-negotiable Implementation Invariants below.
4. Execution-ready Worklist / Handoff Contract below.
5. Acceptance Criteria and Done Criteria.
6. Implementation Steps.
7. Follow-ups.

If a conflict is found during implementation, executor must follow higher source and record the conflict in evidence.

## Non-negotiable Implementation Invariants

- Transition must keep allowing `Need Review` (`1`) when target is `CANCELLED` (`9`).
- Do not modify the FE contract.
- Do not relax tenant (`cust_id`) scope on any read/write.
- Status update and stock mutation must run in a single transaction. No partial success.
- Do not introduce new `data_status` numeric values; reuse existing `entity.*` constants.
- Do not write to `inv.stock` when basis is empty for `Need Review` cancel; only update `sls.order.data_status`.
- Idempotency: never insert a second `SO-CO` row for the same `ro_no` after the first successful cancel.

## Do Not / Reject If

- Do not allow `Need Review` cancel when `data_status = 1` but `qty_final IS NULL` and no `qty_po*` exists AND there is an outstanding `inv.stock` `tr_code = 'SO'` row that still has `qty_out > qty_in + qty_in_order` discrepancy. Reject with existing "invalid outstanding" error.
- Do not allow `Need Review` cancel when basis is ambiguous (mixed sales and PO rows with different unit prices for same `pro_id`). Reject with existing "ambiguous basis" error.
- Do not silently swallow errors during stock writes. If a stock write fails, roll back the whole transaction.
- Do not introduce a new endpoint. Do not change route path or method.
- Do not duplicate `entity.NEED_REVIEW` / `entity.CANCELLED` magic numbers in this endpoint path.
- Do not edit `pjp*`, `inventory`, or `master` services.
- Do not change the request payload shape.

## Diff Boundary

- Allowed: `sales/service/order_service.go`, `sales/repository/stock_repository.go`, `sales/entity/order.go` (only if a new named constant is added; not required), `sales/service/order_service_test.go`, `sales/repository/stock_repository_cancel_test.go`.
- Generated test reports or coverage outputs under `sales/test_results/` allowed only if produced by existing test command.
- Evidence: `.opencode/evidence/20260619-1519-sx-2291-cancel-need-review/` and `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/`.
- Out of boundary: any file outside `sales/` (except `.opencode/`), any FE, any infra, any docker compose, any `.env`.

## TDD/Test Plan

TDD required: yes.

Reason: existing tests already cover cancel transition and stock basis logic. A regression test for `Need Review` + PO qty + no stock movement is the cleanest way to lock the new contract.

Existing test patterns to reuse:
- `sales/service/order_service_test.go` uses `testify` suite pattern with `OrderServiceTestSuite` and `SetupTest` mocks.
- `sales/repository/stock_repository_cancel_test.go` uses `testify` and `sqlmock` for SQL assertions.

First failing test (Red):
- New test in `sales/service/order_service_test.go`:
  - `TestBulkUpdateStatus_Cancel_NeedReview_POQtyOnly_AllowsStatusUpdate`
  - Setup: order with `data_status = 1`, only `qty_po1` populated, no `inv.stock` rows.
  - Assert: `BulkUpdateStatus` returns nil, mocked `Update` called with `data_status = 9`, no `CancelSalesStockUpdates` call (or call with empty commands).
- New test in `sales/repository/stock_repository_cancel_test.go`:
  - `TestCancelStockBasisQuery_POQtyFallback`
  - Assert: SQL filter matches when `COALESCE(qty_final, 0) > 0 OR COALESCE(qty1, 0) > 0 OR COALESCE(qty_po1, 0) > 0`.

Green:
- Extend `cancelStockBasisQuery` filter and SELECT to include `qty*` and `qty_po*` with priority `qty_final` → `qty` → `qty_po*`.
- In service, if basis empty and current status is `NEED_REVIEW`, skip stock writes and proceed with status update + log line.

Refactor:
- Extract qty-resolver to a single helper to avoid duplication.
- Keep function bodies short and named.

Edge cases:
- Already cancelled: existing test `TestBulkUpdateStatus_Cancel_*` covers. Re-verify idempotency: no new `SO-CO` row.
- Mixed basis (sales + PO with same `pro_id`): reject with existing ambiguous error. Keep existing test.
- Tenant mismatch: ensure no cross-tenant read. Reuse existing tenant-scoped tests.

Commands:
- `cd sales && rtk go mod download && rtk go mod tidy`
- `cd sales && rtk go test ./service/... -run TestBulkUpdateStatus_Cancel_NeedReview -v`
- `cd sales && rtk go test ./repository/... -run TestCancelStockBasisQuery -v`
- `cd sales && rtk go test ./...`

## Implementation Steps

1. Read full `cancelStockBasisQuery` and `buildCancelStockMutations` to confirm assumption that `qty_final` is the only qty source.
2. Add qty-resolver helper in `sales/repository/stock_repository.go` (private func) that returns a `cancelOrderStockBasis`-shaped struct using priority `qty_final` → `qty` → `qty_po*`. Use it in `GetCancelStockBasis` result shaping.
3. Update `cancelStockBasisQuery` SQL:
   - Replace `COALESCE(od.qty_final, 0) > 0` with `COALESCE(od.qty_final, 0) > 0 OR COALESCE(od.qty1, 0) > 0 OR COALESCE(od.qty_po1, 0) > 0` (and `qty2/qty3` variants where applicable for that detail row).
   - Add fallback SELECTs: `COALESCE(od.qty_final, COALESCE(od.qty1, COALESCE(od.qty_po1, 0))) AS qty1_smallest`, and analogous for `qty2_smallest`, `qty3_smallest`, `sell_price1`, `sell_price2`, `sell_price3` (priority `sell_price_final` → `sell_price` → `sell_price_po*`).
4. In `sales/service/order_service.go` `BulkUpdateStatus`, when target is `CANCELLED` and basis is empty AND `currentStatus == NEED_REVIEW`, log a warning (existing `log` package) and skip `CancelSalesStockUpdates`. Continue to status update.
5. Keep early return for `*orderData.DataStatus == entity.CANCELLED` (idempotency).
6. Update unit tests:
   - Add Red test for `NEED_REVIEW` + PO qty only.
   - Add Red test for `NEED_REVIEW` + no basis.
   - Add SQL assertion test for extended filter.
7. Run `rtk go test ./...` inside `sales/`.
8. Manual DB validation per docs section `Cancel Order` queries against `SO2606190005` if accessible; otherwise rely on unit tests.

## Expected Files to Change

- `sales/service/order_service.go` — relax cancel guard for `NEED_REVIEW` when basis is empty.
- `sales/repository/stock_repository.go` — qty-resolver helper + extended `cancelStockBasisQuery` filter/SELECT.
- `sales/service/order_service_test.go` — new regression tests.
- `sales/repository/stock_repository_cancel_test.go` — new SQL assertion test.

No other files.

## Agent/Tool Routing

- `@fixer` for implementation and tests.
- `@explorer` if SQL filter or helper details need deeper read.
- `@oracle` for review of the cancel stock basis change before commit.
- `@quality-gate` final signoff for security, tenant scope, and regression risk.
- `@orchestrator` for sequencing.

## Executor Handoff Prompt

Scope: enable `Need Review` cancel for SO from PO tab on `PATCH /sales/v1/orders/status` inside `sales/` only.

must_preserve:
- Transition allowlist including `entity.NEED_REVIEW`.
- Tenant scope `cust_id` on every read/write.
- Transaction wrapping `BulkUpdateStatus` via `WithinTransaction`.
- Existing response shape and error messages used by FE.
- Idempotency: re-cancel does not write a second `SO-CO` row.

do_not_touch:
- `pjp`, `pjp-principle`, `pjp-sales`, `inventory`, `master`, `finance`, `tms`, `mobile`, `system`, `cronjob` services.
- FE, docker, env, infra, migrations.
- `inv.stock` and `inv.warehouse_stock` table schema.
- Request payload shape.

validation:
- `cd sales && rtk go mod download && rtk go mod tidy`
- `cd sales && rtk go test ./...`
- `cd sales && rtk go test ./service/... -run TestBulkUpdateStatus_Cancel -v`
- `cd sales && rtk go test ./repository/... -run TestCancelStockBasisQuery -v`
- Manual: against dev DB, run the cancel curl from this prompt for `SO2606190005`; assert `sls.order.data_status = 9` and a single new `inv.stock` row with `tr_code = 'SO'`, `tr_no = 'SO2606190005-CO'`. Skip manual step if DB not reachable from this environment.

return/evidence expectations:
- `git diff` summary for `sales/` only.
- `go test ./...` output saved to `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/test-output.txt`.
- One-paragraph note on any deviation from this plan, recorded in `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/notes.md`.
- If a real DB was reachable, save SQL outputs to `.opencode/evidence/20260619-1519-sx-2291-cancel-need-review/db-checks.txt`.

claim limits:
- Do not claim "fixed in production" unless a deploy to `best.scyllax.online` was performed and verified.
- Do not claim "regression test passes" unless `go test ./...` shows green output for the new tests.

## Execution-ready Worklist / Handoff Contract

Tasks are atomic and ordered. `start_with` points to the first non-blocked task.

- T1 — Add Red test `TestBulkUpdateStatus_Cancel_NeedReview_POQtyOnly_AllowsStatusUpdate` in `sales/service/order_service_test.go`.
  - depends_on: none
  - owner: @fixer
  - validation: `rtk go test ./service/... -run TestBulkUpdateStatus_Cancel_NeedReview_POQtyOnly_AllowsStatusUpdate -v` shows FAIL
  - exit: test file added; current build still passes for all other tests
  - status: ready
  - requires_user_decision: no
  - must_preserve: existing test pattern in `order_service_test.go`
  - do_not_touch: anything outside `sales/service/order_service_test.go`
  - evidence_update: paste test output to `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/red-1.txt`
  - exit_verification: file compiles under `rtk go build ./...`
  - start_with: T1

- T2 — Add Red test `TestCancelStockBasisQuery_POQtyFallback` in `sales/repository/stock_repository_cancel_test.go`.
  - depends_on: none
  - owner: @fixer
  - validation: `rtk go test ./repository/... -run TestCancelStockBasisQuery_POQtyFallback -v` shows FAIL
  - exit: test file added; existing tests still pass
  - status: ready
  - requires_user_decision: no
  - must_preserve: existing `sqlmock` pattern
  - do_not_touch: anything outside `sales/repository/stock_repository_cancel_test.go`
  - evidence_update: paste test output to `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/red-2.txt`
  - exit_verification: file compiles under `rtk go build ./...`
  - start_with: T1

- T3 — Add qty-resolver helper in `sales/repository/stock_repository.go` and extend `cancelStockBasisQuery` filter/SELECT per Implementation Steps 2-3.
  - depends_on: T1, T2
  - owner: @fixer
  - validation: `rtk go test ./repository/... -run TestCancelStockBasisQuery -v` shows both tests pass
  - exit: helper exported as private; SQL filter includes `qty1`, `qty_po1` fallback; SQL select exposes priority-resolved `qty1_smallest`/`qty2_smallest`/`qty3_smallest` and `sell_price*`
  - status: ready
  - requires_user_decision: no
  - must_preserve: existing `item_type = 1`, tenant scope, `tr_code = 'SO'` join keys
  - do_not_touch: inv.stock and inv.warehouse_stock schema; cancel mutation writer for `SO-CO` row shape
  - evidence_update: append diff hunk to `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/diff-stock-repo.txt`
  - exit_verification: `rtk go test ./repository/...` all green
  - start_with: T1

- T4 — Update `BulkUpdateStatus` in `sales/service/order_service.go` so `NEED_REVIEW` cancel with empty basis updates status and skips stock writes, with a log line.
  - depends_on: T3
  - owner: @fixer
  - validation: `rtk go test ./service/... -run TestBulkUpdateStatus_Cancel -v` shows new test green and existing tests still green
  - exit: idempotency for `data_status == 9` retained; new branch only fires for `NEED_REVIEW` + empty basis + `CANCELLED` target
  - status: ready
  - requires_user_decision: no
  - must_preserve: transaction wrap, response shape, error message text for `invalid status transition` and `inconsistent basis`
  - do_not_touch: routes, controller, entity constants unless strictly required
  - evidence_update: append diff hunk to `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/diff-order-service.txt`
  - exit_verification: `rtk go test ./...` all green in `sales/`
  - start_with: T1

- T5 — Run full `sales` test suite and capture output.
  - depends_on: T4
  - owner: @fixer
  - validation: `rtk go test ./...` in `sales/` shows all green
  - exit: zero failures, no skipped new tests
  - status: ready
  - requires_user_decision: no
  - must_preserve: n/a
  - do_not_touch: n/a
  - evidence_update: save full output to `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/test-output.txt`
  - exit_verification: file exists and contains `ok` for every tested package
  - start_with: T1

- T6 — Manual DB check (only if dev DB reachable).
  - depends_on: T5
  - owner: @fixer
  - validation: cancel curl against dev returns success; `SELECT data_status FROM sls.order WHERE ro_no = 'SO2606190005'` returns 9; one new `inv.stock` row with `tr_code = 'SO'`, `tr_no = 'SO2606190005-CO'` exists
  - exit: SQL outputs saved
  - status: blocked if DB not reachable; record blocker in notes
  - requires_user_decision: no
  - must_preserve: original order data; do not modify unrelated rows
  - do_not_touch: production DB; only run against the local dev DB per `AGENTS.md`
  - evidence_update: save outputs to `.opencode/evidence/20260619-1519-sx-2291-cancel-need-review/db-checks.txt`
  - exit_verification: SQL outputs present
  - start_with: T5

- T7 — Quality gate review.
  - depends_on: T5, T6 (or T5 only if T6 blocked)
  - owner: @quality-gate
  - validation: no new `data_status` magic numbers; tenant scope preserved; idempotency preserved; no schema change
  - exit: PASS or BLOCKED with reason
  - status: ready
  - requires_user_decision: no
  - must_preserve: invariants
  - do_not_touch: source code beyond approved diff
  - evidence_update: gate verdict in `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/qg-verdict.md`
  - exit_verification: file exists
  - start_with: T1

## Validation Commands

- `cd sales && rtk go mod download && rtk go mod tidy`
- `cd sales && rtk go build ./...`
- `cd sales && rtk go test ./...`
- `cd sales && rtk go test ./service/... -run TestBulkUpdateStatus_Cancel -v`
- `cd sales && rtk go test ./repository/... -run TestCancelStockBasisQuery -v`
- `cd sales && rtk go test ./service/... -run TestValidateCancelTransition -v`

## Evidence Requirements

- `.opencode/evidence/20260619-1519-sx-2291-cancel-need-review/discovery.md` (kept).
- `.opencode/evidence/20260619-1519-sx-2291-cancel-need-review/index.json` (kept).
- `.opencode/evidence/20260619-1519-sx-2291-cancel-need-review/db-checks.txt` (if T6 runs).
- `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/red-1.txt`, `red-2.txt` (T1, T2 outputs).
- `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/diff-stock-repo.txt`, `diff-order-service.txt` (T3, T4 diffs).
- `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/test-output.txt` (T5).
- `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/qg-verdict.md` (T7).
- Source strategy note: local repo evidence; user-supplied docs; no external web or browser evidence. No current API docs lookup needed; behavior fully determined by repo + user prompt.

## Done Criteria

- All tasks T1..T7 marked complete or T6 explicitly blocked with reason.
- `go test ./...` in `sales/` green.
- New regression tests present and passing.
- Tenant scope, idempotency, and tx safety preserved.
- Diff stays inside `sales/` and the listed test files.

## Final Planning Summary

- Artifacts created: `.opencode/plans/20260619-1519-sx-2291-cancel-need-review.md`, `.opencode/evidence/20260619-1519-sx-2291-cancel-need-review/discovery.md`, `.opencode/evidence/20260619-1519-sx-2291-cancel-need-review/index.json`.
- Draft folder `.opencode/draft/20260619-1519-sx-2291-cancel-need-review/` reserved for executor evidence; not yet populated (will be filled during implementation).
- Key decisions: keep `validateCancelTransition` as-is; extend `cancelStockBasisQuery` filter/SELECT with `qty*` and `qty_po*` fallback; for `NEED_REVIEW` cancel with empty basis, skip stock writes and update status; preserve idempotency for `data_status = 9`.
- Assumptions: tenant scope sufficient; no parent_cust_id needed for this endpoint; bulk array behavior unchanged.
- Open questions: none blocking. If executor finds `parent_cust_id` enforcement required for compliance, route back to planner.
- Readiness: `ready-for-implementation`.
- Cleanup: drafts under `.opencode/draft/<task-id>/` will be deleted after final synthesis if not operationally useful.
