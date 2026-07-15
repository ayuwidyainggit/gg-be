# Plan — Restore Sales Order Endpoint Code From `dev` To `bugfix/SX-2209-qa`

Task ID: `20260617-1845-restore-sales-orders-dev-to-sx-2209-qa`

Readiness: `ready-for-slice`

Primary source of truth: this file.

## Goal

Create branch `bugfix/SX-2209-qa` from `qa`, then restore only relevant code from `dev` for:

- `POST /sales/v1/orders`
- `GET /sales/v2/orders/{order_no}`

Internal route currently uses `:ro_no`, so user-facing `{order_no}` maps to code param `ro_no` unless product/API owner says otherwise.

## Non-goals

- Do not merge all `dev` into `qa`.
- Do not restore unrelated `dev` changes for invoice, report, Open API middleware/config, unrelated promotion endpoints, or non-order behavior.
- Do not change DB schema unless implementation proves endpoint behavior requires existing fields missing from `qa`.
- Do not commit, push, or create MR unless user explicitly asks.
- Do not copy secrets or `.env` values.

## Scope

In scope:

- Branch operation inside `sales/` Git repo.
- Endpoint route/controller/service/repository/entity/model/test code needed by the two endpoints.
- Regression tests for restored behavior.
- Validation inside `sales/` module.

Out of scope unless explicitly approved:

- Root monorepo changes outside `sales/`.
- Whole-branch `dev` merge.
- Remote deploy/DB migration.
- Cleanup of unrelated existing test/LSP failures unless they block endpoint validation.

## Requirements

1. Start from clean `sales/` working tree.
2. Update `qa` and `dev` refs from remote before extracting code.
3. Create `bugfix/SX-2209-qa` from `qa`.
4. Bring endpoint-relevant code from `dev` with smallest safe diff.
5. Preserve route ownership:
   - `POST /v1/orders` → `OrderController.Create`
   - `GET /v2/orders/:ro_no` → `OrderController.DetailV2`
6. Preserve auth middleware: `middleware.JWTProtected()`.
7. Preserve repo architecture: Controller → Service → Repository → DB.
8. Preserve service transaction boundary for order creation.
9. Preserve repository tx-context extraction for writes.
10. Preserve tenant rules: `cust_id`, `parent_cust_id`, distributor-specific data isolation.
11. Add or restore regression tests from `dev` for SX-2209 behavior where relevant.
12. Validate with targeted tests first, then broader `sales` module tests if feasible.

## Acceptance Criteria

- `bugfix/SX-2209-qa` exists locally from `qa`.
- Diff against `qa` contains only files needed for requested endpoints and their tests.
- `POST /sales/v1/orders` compiles and keeps create-order behavior from selected `dev` endpoint code.
- `GET /sales/v2/orders/{order_no}` compiles and returns purchase details including original purchase qty rows per SX-2209 behavior.
- Regression test for `DetailV2` original purchase qty passes.
- No unrelated Open API, invoice, report, or global middleware/config changes land unless proven required and documented.
- `rtk go test ./service -run TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero` passes after implementation.
- `rtk go test ./service` passes or every failure is classified as pre-existing/unrelated with evidence.
- Final diff reviewed against `qa` before handoff.

## Existing Patterns / Reuse

Repo-local reuse:

- `sales/controller/order_controller.go` already owns both endpoint routes.
- `sales/service/order_service.go` already owns `Store` and `DetailV2`.
- `sales/repository/order_repository.go` already owns order persistence and detail lookup.
- `sales/service/order_service_test.go` already has `mockOrderRepositoryDetailV2` and DetailV2 tests.
- `sales/controller/so_controller_test.go` shows Fiber `httptest` controller-test pattern, but service-level tests are better first slice for SX-2209.

Branch/source reuse:

- Prefer direct reuse from `dev`/`bugfix/SX-2209-dev` over reimplementation.
- Strong candidate commit: `f784cf8 fix(order): show purchase rows with original qty`.
- Changed files in candidate commit:
  - `service/order_service.go`
  - `service/order_service_test.go`

## Constraints

- `sales/` is Git repo; root `/scylla-be` is not Git repo.
- Use `rtk` prefix for shell workflows in this repo.
- Validate inside `/Users/ujang/Projects/Geekgarden/scylla-be/sales`.
- Service README is stale template; prefer repo docs, `go.mod`, source, tests.
- Missing optional framework docs:
  - `.opencode/docs/PROJECT_STACK.md`
  - `.opencode/docs/PROJECT_COMMANDS.md`
  - `.opencode/docs/FRAMEWORK_PLAYBOOK.md`
- Existing LSP diagnostics appeared in unrelated/test files during artifact write; executor must re-check after branch creation before claiming new failures.

## Risks

- `dev` includes broad unrelated changes. Blind merge/cherry-pick may drag invoice/report/Open API changes into QA bugfix branch.
- Endpoint behavior may depend on earlier `dev` commits, especially `order_type` / taking-order support. Executor must verify dependencies with tests and compile errors, not assume `f784cf8` alone is enough.
- `POST /v1/orders` has broad create-order logic; restoring too much can alter inventory, promotion, validation, or transaction behavior.
- Param naming mismatch: user says `{order_no}`, code uses `ro_no`. Changing public route param could break clients. Keep `:ro_no` unless explicitly required.
- Existing tests may already fail on `qa`; classify pre-existing failures with clean branch evidence.

## Decisions / Assumptions

- Decision: branch target is local `sales/` repo, not root folder.
- Decision: create `bugfix/SX-2209-qa` from `qa`.
- Decision: use endpoint-scoped restore, not whole `dev` merge.
- Decision: first candidate source is `f784cf8`; widen to supporting commits only if compile/tests prove dependency.
- Assumption: `{order_no}` in user request means existing internal `ro_no` route param.
- Assumption: user wants QA bugfix branch with SX-2209 endpoint behavior, not full dev parity.
- Open question: if user wants exact full endpoint parity with current `dev`, executor must ask before including broad dependencies outside endpoint scope.

## Execution Source of Truth

Precedence during implementation:

1. Latest explicit user instruction.
2. Safety/security/secret rules and branch safety.
3. Non-negotiable Implementation Invariants.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria and Done Criteria.
6. Implementation Steps.
7. Follow-up recommendations.

If conflict exists, follow higher-precedence source and record conflict in evidence.

## Non-negotiable Implementation Invariants

- Work inside `sales/` Git repo for branch/diff/source edits.
- Do not merge all `dev` into `bugfix/SX-2209-qa` unless user explicitly approves.
- Keep branch base as `qa`.
- Keep `POST /v1/orders` and `GET /v2/orders/:ro_no` protected by `middleware.JWTProtected()`.
- Keep Controller → Service → Repository → DB layering.
- Keep create-order DB writes inside service transaction.
- Keep repository writes tx-aware through `model(ctx)` / transaction context pattern.
- Keep `cust_id` and `parent_cust_id` behavior from auth locals and tenant filters.
- Keep route param `:ro_no` unless user explicitly requests public route rename.
- Any supporting `dev` code copied outside endpoint files must have evidence explaining dependency.
- Do not expand or copy tracked plaintext credentials.

## Do Not / Reject If

Reject implementation if:

- Diff includes unrelated `model/open_api_config.go`, `pkg/middleware/open_api_middleware.go`, invoice/report changes, or global config changes without documented dependency.
- Executor uses `git merge dev` as default.
- Executor changes route to `/v2/orders/:order_no` without explicit approval.
- Controller calls repository directly.
- Repository gains business logic.
- Create flow writes outside transaction.
- Tenant filters are removed or weakened.
- Test failures are claimed “unrelated” without command output and baseline comparison.
- Branch was not created from `qa`.

## Diff Boundary

Allowed file groups by default:

- `sales/service/order_service.go`
- `sales/service/order_service_test.go`
- `sales/controller/order_controller.go` only if route/controller delta from `dev` is required.
- `sales/entity/order.go` only if restored request/response fields are required.
- `sales/model/order.go` only if restored response/source fields are required.
- `sales/model/order_detail.go` only if restored detail fields like `original_qty_po*` are required and already expected by code/tests.
- `sales/repository/order_repository.go` only if `DetailV2`/`Create` required query fields from `dev` are missing.

Generated/evidence exceptions:

- `.opencode/evidence/20260617-1845-restore-sales-orders-dev-to-sx-2209-qa/**`
- `.opencode/plans/20260617-1845-restore-sales-orders-dev-to-sx-2209-qa.md`

Any file outside this boundary must be reverted or justified in verification evidence before final quality gate.

## TDD / Test Plan

TDD required: yes.

Reason: endpoint behavior and service logic change, with tenant/data correctness risk.

Existing test patterns:

- `sales/service/order_service_test.go` uses mock repositories and direct service calls.
- `sales/controller/so_controller_test.go` uses Fiber + `httptest` for controller behavior.

First failing/regression test:

- Restore/add `TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero` from `dev` commit `f784cf8` into `sales/service/order_service_test.go`.
- Run it on target branch before production code patch if practical; it should fail or not compile until behavior exists.

Green step:

- Add `hasPurchaseDisplayQty` and `shouldIncludePurchaseDetailRow` behavior from `dev`.
- Replace purchase-details inclusion condition in `DetailV2` with `shouldIncludePurchaseDetailRow(detail)`.
- Add supporting model fields only if missing on `qa`.

Refactor step:

- Keep helper names and placement consistent with `dev` unless conflict with QA code.
- Remove accidental unrelated hunks from patch.
- Re-run targeted tests after any conflict resolution.

Edge cases:

- Current `qty_po1/2/3` zero but `original_qty_po1/2/3` non-zero → included in `purchase_details.normal`.
- Current and original purchase qty all zero, sales qty zero → excluded.
- `item_type = 2` promo rows still excluded from normal purchase rows.
- Sales order details remain empty when sales qty nil/zero.
- Final details remain empty when final qty zero.
- `purchase_details.normal[*].order_status` still follows header status.

Commands:

```bash
rtk go test ./service -run TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero
rtk go test ./service -run TestDetailV2
rtk go test ./service
rtk go test ./...
```

If `rtk go test ./...` fails due pre-existing unrelated controller mock errors, capture:

```bash
rtk git status --short --branch
rtk go test ./... 2>&1 | tee ../.opencode/evidence/20260617-1845-restore-sales-orders-dev-to-sx-2209-qa/go-test-all.txt
```

Then run narrower passing service checks and document blocker.

## Implementation Steps

1. Verify clean state:

```bash
cd /Users/ujang/Projects/Geekgarden/scylla-be/sales
rtk git status --short --branch
```

2. Fetch refs:

```bash
rtk git fetch origin dev qa bugfix/SX-2209-dev
```

If `origin/bugfix/SX-2209-dev` unavailable, use local `dev` commit evidence.

3. Reset/update local `qa` from remote only if no local divergence:

```bash
rtk git checkout qa
rtk git status --short --branch
rtk git pull --ff-only origin qa
```

4. Create target branch:

```bash
rtk git checkout -b bugfix/SX-2209-qa qa
```

If branch already exists:

```bash
rtk git checkout bugfix/SX-2209-qa
rtk git merge-base --is-ancestor qa HEAD
```

If not based on `qa`, stop and ask before reset/recreate.

5. Baseline target tests before patch:

```bash
rtk go test ./service -run TestDetailV2
rtk go test ./...
```

Capture failures as baseline if any.

6. Inspect source candidate:

```bash
rtk git show --stat --oneline f784cf8
rtk git show --function-context f784cf8 -- service/order_service.go service/order_service_test.go
```

7. Apply only candidate patch first:

Preferred:

```bash
rtk git cherry-pick --no-commit f784cf8
```

If cherry-pick drags only `service/order_service.go` and `service/order_service_test.go`, continue. If conflicts or unrelated hunks appear, abort and apply manually:

```bash
rtk git cherry-pick --abort
rtk git checkout dev --patch -- service/order_service.go service/order_service_test.go
```

8. Compile/test targeted behavior:

```bash
rtk go test ./service -run TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero
rtk go test ./service -run TestDetailV2
```

9. If compile fails due missing fields/helpers, inspect source dependencies from `dev`:

```bash
rtk git diff qa..dev -- entity/order.go model/order.go model/order_detail.go repository/order_repository.go controller/order_controller.go service/order_service.go
```

Copy only missing endpoint-required hunks with `git checkout dev --patch -- <file>`.

10. If `POST /v1/orders` parity requested by test/product evidence, compare create flow:

```bash
rtk git diff qa..dev -- controller/order_controller.go service/order_service.go entity/order.go model/order.go model/order_detail.go repository/order_repository.go
```

Apply only hunks used by `OrderController.Create`, `OrderService.Store`, and direct dependencies. Reject unrelated Open API/invoice/report hunks.

11. Review diff boundary:

```bash
rtk git diff --name-status qa...HEAD
rtk git diff -- controller/order_controller.go service/order_service.go service/order_service_test.go entity/order.go model/order.go model/order_detail.go repository/order_repository.go
```

12. Run validation:

```bash
rtk go test ./service -run TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero
rtk go test ./service -run TestDetailV2
rtk go test ./service
rtk go test ./...
```

13. Optional runtime smoke if service and auth token available:

```bash
rtk docker compose -f ../docker-compose.yml ps
```

Do not call real order endpoints against remote/dev DB. Use local DB only.

14. Record evidence in `.opencode/evidence/20260617-1845-restore-sales-orders-dev-to-sx-2209-qa/verification.md`.

15. Route to `@quality-gate` for final review before push/MR.

## Expected Files To Change

Expected minimal SX-2209 slice:

- `sales/service/order_service.go`
- `sales/service/order_service_test.go`

Possible if dependencies required:

- `sales/model/order_detail.go`
- `sales/model/order.go`
- `sales/entity/order.go`
- `sales/repository/order_repository.go`
- `sales/controller/order_controller.go`

Do not expect:

- `sales/main.go`
- `sales/go.mod`
- `sales/go.sum`
- `sales/model/open_api_config.go`
- `sales/pkg/middleware/open_api_middleware.go`
- invoice/report files

## Agent / Tool Routing

- `@orchestrator`: route implementation, enforce branch/diff scope.
- `@fixer`: apply bounded code/test patch and run validation.
- `@explorer`: optional if compile errors reveal hidden dependencies.
- `@oracle`: optional if endpoint parity requires broader `dev` behavior and tradeoff review.
- `@quality-gate`: mandatory final signoff before push/MR.

No UI/designer/browser needed.

## Executor Handoff Prompt

Copyable prompt:

```text
Implement plan `.opencode/plans/20260617-1845-restore-sales-orders-dev-to-sx-2209-qa.md` in repo `/Users/ujang/Projects/Geekgarden/scylla-be/sales`.

Scope: create/use branch `bugfix/SX-2209-qa` from `qa`, restore endpoint-scoped code from `dev` for `POST /sales/v1/orders` and `GET /sales/v2/orders/{order_no}` (`:ro_no` in code). Start with commit `f784cf8 fix(order): show purchase rows with original qty`; widen only if compile/tests prove dependency.

must_preserve:
- branch base `qa`
- no full `dev` merge
- Controller → Service → Repository → DB
- service transaction boundary for create-order writes
- tx-aware repository writes
- `cust_id` / `parent_cust_id` tenant behavior
- `middleware.JWTProtected()` on both routes
- route param remains `:ro_no` unless user approves rename

do_not_touch:
- secrets, `.env`
- unrelated invoice/report/Open API files
- root repo files outside `sales/` except `.opencode/evidence/**`

validation:
- `rtk go test ./service -run TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero`
- `rtk go test ./service -run TestDetailV2`
- `rtk go test ./service`
- `rtk go test ./...` or documented baseline blocker

Return:
- branch status
- changed files
- commits/patch source used
- commands run and outputs summary
- remaining risks
- evidence path under `.opencode/evidence/20260617-1845-restore-sales-orders-dev-to-sx-2209-qa/`
```

## Execution-ready Worklist / Handoff Contract

`start_with`: `T1`

### T1 — Confirm repo and branch safety

- action: verify clean `sales/` working tree and branch refs.
- depends_on: none
- owner/lane: `@fixer`
- validation: `rtk git status --short --branch`; `rtk git branch --list dev qa bugfix/SX-2209-qa`; `rtk git branch -r --list '*/dev' '*/qa' '*/bugfix/SX-2209-*'`
- exit criteria: clean tree; `qa` and `dev` available.
- blocking status: ready
- blocker reason: none
- requires_user_decision: no
- must_preserve: no source edits before clean state.
- do_not_touch: secrets, env files.
- evidence_update: note branch status in verification evidence.
- exit_verification: command output recorded.

### T2 — Create target branch from `qa`

- action: checkout updated `qa`, create `bugfix/SX-2209-qa`.
- depends_on: T1
- owner/lane: `@fixer`
- validation: `rtk git merge-base --is-ancestor qa HEAD`; `rtk git status --short --branch`
- exit criteria: current branch is `bugfix/SX-2209-qa`, based on `qa`.
- blocking status: ready
- blocker reason: none
- requires_user_decision: no unless branch exists with divergent history.
- must_preserve: branch base `qa`.
- do_not_touch: do not reset existing branch without approval.
- evidence_update: record branch creation command and base SHA.
- exit_verification: merge-base check passes.

### T3 — Establish baseline tests

- action: run targeted and broad tests on clean target branch before patch.
- depends_on: T2
- owner/lane: `@fixer`
- validation: `rtk go test ./service -run TestDetailV2`; `rtk go test ./...`
- exit criteria: baseline pass or baseline failures captured.
- blocking status: ready
- blocker reason: none
- requires_user_decision: no
- must_preserve: classify baseline failures separately from patch failures.
- do_not_touch: no code changes yet.
- evidence_update: save failure summary if any.
- exit_verification: evidence has baseline status.

### T4 — Apply SX-2209 source patch

- action: apply `f784cf8` with `--no-commit` or manual endpoint-scoped patch.
- depends_on: T3
- owner/lane: `@fixer`
- validation: `rtk git diff --name-status qa...HEAD`
- exit criteria: only expected files changed; no unrelated dev hunks.
- blocking status: ready
- blocker reason: none
- requires_user_decision: no
- must_preserve: no full dev merge; no unrelated Open API/invoice/report changes.
- do_not_touch: `sales/main.go`, `sales/go.mod`, `sales/go.sum`, unrelated files unless proven needed.
- evidence_update: record source commit and changed files.
- exit_verification: diff boundary reviewed.

### T5 — Resolve compile/test dependencies narrowly

- action: if tests fail due missing endpoint dependencies, apply only required hunks from `dev`.
- depends_on: T4
- owner/lane: `@fixer` with optional `@explorer`
- validation: targeted test command and `rtk git diff --name-status qa...HEAD`
- exit criteria: compile passes for targeted tests; extra files justified.
- blocking status: ready
- blocker reason: none
- requires_user_decision: yes if dependency requires broad full endpoint parity or unrelated dev files.
- must_preserve: tenant and transaction invariants.
- do_not_touch: unrelated dev changes.
- evidence_update: dependency rationale per added file.
- exit_verification: targeted compile/test result captured.

### T6 — Validate endpoint behavior

- action: run service tests and optional controller/runtime checks.
- depends_on: T5
- owner/lane: `@fixer`
- validation:
  - `rtk go test ./service -run TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero`
  - `rtk go test ./service -run TestDetailV2`
  - `rtk go test ./service`
  - `rtk go test ./...`
- exit criteria: required targeted tests pass; broad failures classified if any.
- blocking status: ready
- blocker reason: none
- requires_user_decision: no unless tests require live DB/token decisions.
- must_preserve: no real remote DB use.
- do_not_touch: remote dev DB, `.env`.
- evidence_update: command output summary.
- exit_verification: validation evidence ready.

### T7 — Final diff review and quality gate

- action: review final diff, route to `@quality-gate`.
- depends_on: T6
- owner/lane: `@quality-gate`
- validation: `rtk git diff --name-status qa...HEAD`; inspect risky hunks.
- exit criteria: quality gate PASS or documented fixes.
- blocking status: ready
- blocker reason: none
- requires_user_decision: no unless quality gate flags scope expansion.
- must_preserve: plan invariants and diff boundary.
- do_not_touch: commit/push unless user asks.
- evidence_update: quality gate result path/summary.
- exit_verification: final status and remaining risks documented.

## Validation Commands

From `sales/`:

```bash
rtk git status --short --branch
rtk git diff --name-status qa...HEAD
rtk go test ./service -run TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero
rtk go test ./service -run TestDetailV2
rtk go test ./service
rtk go test ./...
```

Optional runtime readiness from repo root if needed:

```bash
rtk docker compose -f docker-compose.yml ps
```

## Evidence Requirements

Implementation evidence must include:

- Current branch and base SHA.
- Source commit(s) or source branch hunks used.
- Changed file list.
- Why each non-minimal file was needed.
- Test commands and pass/fail output summary.
- Any pre-existing failures with baseline comparison.
- Confirmation no unrelated `dev` merge occurred.

Keep evidence under:

- `.opencode/evidence/20260617-1845-restore-sales-orders-dev-to-sx-2209-qa/`

Existing planning evidence kept:

- `.opencode/evidence/20260617-1845-restore-sales-orders-dev-to-sx-2209-qa/discovery.md`

## Done Criteria

- Branch `bugfix/SX-2209-qa` from `qa` contains endpoint-scoped restore.
- SX-2209 DetailV2 regression test passes.
- `POST /v1/orders` compile path remains valid if touched.
- Diff boundary respected or exceptions justified.
- Validation evidence recorded.
- `@quality-gate` reviewed before push/MR.

## Final Planning Summary

Artifacts consulted:

- `AGENTS.md`
- `.opencode/docs/index.md`
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `sales/controller/order_controller.go`
- `sales/service/order_service.go`
- `sales/repository/order_repository.go`
- `sales/entity/order.go`
- `sales/service/order_service_test.go`
- Git branch/diff/log evidence in `sales/`

Artifacts created:

- `.opencode/evidence/20260617-1845-restore-sales-orders-dev-to-sx-2209-qa/discovery.md`
- `.opencode/plans/20260617-1845-restore-sales-orders-dev-to-sx-2209-qa.md`

Key decisions:

- Use `sales/` as Git repo.
- Create target from `qa`.
- Avoid full `dev` merge.
- Start with `f784cf8` and widen only by proven dependency.

Assumptions:

- `{order_no}` means code param `ro_no`.
- User wants endpoint-scoped restore from `dev`, not full dev parity.

Remaining open questions:

- If full exact endpoint parity with current `dev` is required, user must approve broader dependency inclusion.

Readiness:

- `ready-for-slice`: safe first implementation slice exists. Whole-dev parity remains intentionally open.

Cleanup performed:

- Draft artifacts not created.
- Evidence kept because implementation needs branch/source discovery and risk notes.
