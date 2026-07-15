# SX-2508 — Import Sales Order Values

- **Task ID:** `20260713-1400-sx-2508-import-order-values`
- **Readiness:** `ready-for-slice`
- **Plan Quality Gate:** `PASS_FOR_SLICE` after validator pass
- **Plan profile:** `maintenance`
- **Substance:** `non-ui`
- **Mode:** `maintenance`

## Goal

Repair SX-2508 only in `sales`: imported Sales Orders must persist required header/detail values, preserve canonical API quantity orientation (`qty1=Large`, `qty2=Middle`, `qty3=Small`), avoid re-converting mapped orders in `DetailV2`, and show available stock through existing canonical breakdown helper. First slice covers one import-to-detail happy path, existing persistence layer, no new integration, and local regression proof. Runtime staging proof remains conditional on a fresh user-supplied `SCYLLAX_TOKEN`.

## Non-goals

No migration, endpoint change, module outside `sales`, dependency, cache, global quantity converter rewrite, bulk data repair, secret handling, or new mapping architecture. No staging write beyond user-authorized existing mapping edit. No production-ready claim without local validation and redacted runtime evidence.

## Scope

1. Import parser/header/detail mapping in `sales/service/order_service.go`.
2. Mapped-order branch in `OrderService.DetailV2`.
3. Existing stock presentation helper only when focused Red test proves boundary defect.
4. Parser, DetailV2, stock helper regression tests.
5. Redacted local and conditional staging/SQL evidence.

## Requirements

1. `invoice_no` persists from `ro_no`.
2. `invoice_date` persists from `ro_date`.
3. Imported `qty_po` persists `NULL`.
4. `promo_so1` persists import promotion value.
5. `promo_final1` persists import final-promotion value.
6. `vat_value` and `vat_value_final` persist PPN value.
7. `sell_price_system2` and `sell_price_system3` derive from parent `sell_price3`.
8. Product/UOM mapping follows latest edited mapping; prove lookup behavior with test before changing repository seam.
9. Mapped `DetailV2` returns stored SO and Final triples without conversion.
10. Non-mapped `DetailV2` preserves legacy conversion.
11. Stock arithmetic converts to smallest unit before canonical Large/Middle/Small display.
12. Controller → Service → Repository → DB and tenant/transaction contracts stay intact.

## Acceptance Criteria

1. Parser regression proves invoice number/date values.
2. Parser regression proves `QtyPo == nil`.
3. Parser regression proves promotion fields.
4. Parser regression proves both VAT fields.
5. Parser regression proves both system-price fields use parent `sell_price3`.
6. Mapping regression proves current UOM data used after mapping edit.
7. Mapped `DetailV2` returns persisted Large/Middle/Small values for SO and Final quantities.
8. Non-mapped `DetailV2` control regression remains legacy-compatible.
9. Available stock regression proves canonical breakdown.
10. `cd sales && rtk go test ./...` and `cd sales && rtk go build .` pass.
11. Runtime claims omitted unless fresh-token curl and read-only SQL evidence pass.

## Existing Patterns/Reuse

- Reuse `computeDisplayedAvailableStockBreakdown()` in `sales/service/order_stock_helper.go`.
- Reuse existing import parser/store flow in `sales/service/order_service.go`.
- Reuse existing service tests: `sales/service/order_service_test.go`, `sales/service/order_import_parser_test.go`, `sales/service/order_stock_helper_test.go`.
- Do not change `sales/pkg/conversion/quantity.go`; its internal output remains Small/Middle/Large.

## Source Anatomy

| Layer | Authority | Confirmed finding |
| --- | --- | --- |
| Import parser | `sales/service/order_service.go:6629-7159` | Existing parser sets `IsSalesMapping`, invoice header values, and VAT fields partially. |
| Detail presentation | `sales/service/order_service.go:3004-3113` | `DetailV2` currently recomputes quantities unconditionally. |
| Stock | `sales/service/order_stock_helper.go` | Canonical helper remaps internal converter output to API Large/Middle/Small. |
| Converter | `sales/pkg/conversion/quantity.go` | Internal converter returns Small/Middle/Large. |
| Mapping query | `sales/repository/order_repository.go` | Freshness behavior needs focused proof. |
| DTO fields | `sales/entity/order_detail.go` | DTO contains `QtyPo`, promotion, VAT, system-price fields. |

## Reference Map

| Feature | Basis | Why sufficient |
| --- | --- | --- |
| Import field persistence | repo-backed | Exact parser and DTO inspected. |
| Mapped quantity display | repo-backed | Exact `DetailV2` conversion seam inspected. |
| Stock orientation | repo-backed | Existing canonical helper inspected. |
| Runtime protocol | user-confirmed | User supplied curl/SQL contract. |
| Go/Fiber/GORM compatibility | docs-backed | `sales/go.mod` plus librarian check; no dependency change planned. |

## Constraints

- Go `1.23.0`, toolchain `go1.24.6`; Fiber `v2.52.6`; GORM `v1.24.7-0.20230306060331-85eaf9eeda11`.
- Per-service module rules: run commands from `sales`.
- `rtk` prefix required by repo-local `AGENTS.md`.
- Never write token, Authorization header, `.env`, staging response body, or tracked secret to evidence.

## Risks

- Quantity orientation easily reversed at converter boundary.
- Product mapping query may select stale UOM. Red test must identify exact repository behavior first.
- Staging access unavailable without fresh token.
- Parent-product persistence rule may conflict with mapping data; document exception only when Red proof requires it.

## Decisions/Assumptions

- **confirmed_repo:** `ConvToQtyConversion()` output orientation is internal Small/Middle/Large.
- **confirmed_repo:** existing stock helper converts to canonical API orientation.
- **confirmed_repo:** mapped marker is set during parsing.
- **assumption:** required import columns already exist in source format; implementation must validate via parser fixture.
- **unverified:** repository lookup freshness after distributor mapping edit.
- **user_confirmed:** staging proof requires fresh `SCYLLAX_TOKEN`; token never stored.

## Execution Source of Truth

Precedence: latest explicit user instruction; security and permission rules; Non-negotiable Implementation Invariants; Handoff Contract; Acceptance/Done Criteria; Implementation Steps. Higher source wins; executor records conflict in evidence.

## Non-negotiable Implementation Invariants

- API/domain quantity order always Large/Middle/Small.
- Global converter stays unchanged.
- `IsSalesMapping` orders return persisted SO/Final triples in `DetailV2`; no re-conversion.
- Non-mapped orders retain legacy conversion.
- Stock math uses smallest unit before output breakdown.
- Controller → Service → Repository → DB; writes remain transaction-aware and tenant-safe.
- Planner artifact is source of truth. Execution lanes refresh own permissions.

## Do Not / Reject If

- Reject converter rewrite, migration, cache, dependency, endpoint, or module expansion.
- Reject swapped quantity orientation.
- Reject mapped-order conversion rerun.
- Reject staging success claim without fresh-token evidence.
- Reject evidence containing secret/token or sensitive full response.
- Reject repository refactor without failing proof.

## Diff Boundary

Allowed: `sales/service/order_service.go`; focused tests under `sales/service/`; conditional `sales/service/order_stock_helper.go` or `sales/repository/order_repository.go` only after Red proof; task evidence/state. Out-of-bound change must be reverted or justified in evidence.

## TDD/Test Plan

TDD required: parser, presentation, mapping, and stock behavior are production logic.

- **Red:** add focused parser tests for each required persisted value; mapped/non-mapped DetailV2 tests; stock orientation test; mapping freshness test.
- **Green:** smallest parser/detail/helper/repository change that passes focused tests.
- **Refactor:** remove duplication only after full service/module tests pass.
- **Edges:** middle-only, largest-only, smallest-only mapped quantities; absent `qty_po`; PPN zero/nonzero; edited UOM mapping; non-mapped control.
- **Commands:** targeted tests then `rtk go test ./...` and `rtk go build .`.

## Implementation Steps

1. Run compose status from root; do not start/modify runtime unless needed.
2. Read current import parser, Store assembly, mapping repository method, DetailV2, stock helper, existing test fixtures.
3. Record exact current mapping query order/filter in A1 evidence.
4. Add Red import parser fixture with invoice, promo, PPN, parent price, and `qty_po` assertions.
5. Add Red mapping-edit fixture proving selected UOM behavior.
6. Add Red mapped DetailV2 SO and Final triple assertions for each unit position.
7. Add Red non-mapped DetailV2 control assertion.
8. Add Red stock canonical output assertion.
9. Run targeted test and retain redacted failing output.
10. Patch parser field mapping only.
11. Set `QtyPo` explicitly nil in import request assembly.
12. Map promotion values to `PromoSo1` and `PromoFinal1`.
13. Map PPN to both VAT fields.
14. Map both system price fields from parent `SellPrice3`.
15. Keep parent-product persistence unless Red proof documents exception.
16. Fix mapping repository seam only if mapping Red test proves stale selection.
17. Add mapped branch before DetailV2 conversion; use stored triples.
18. Preserve non-mapped existing conversion branch.
19. Reuse existing stock helper; patch only if stock Red test requires it.
20. Run focused parser tests.
21. Run DetailV2 and stock tests.
22. Run full sales module tests and build.
23. Run pre-gate smoke check.
24. If fresh token exists, run user-provided staging curl cases without logging token/body secrets.
25. Run user-provided read-only SQL field assertions.
26. Write redacted A4 runtime report or exact blocked reason.
27. Run quality gate review.

## Expected Files to Change

- `sales/service/order_service.go`
- `sales/service/order_import_parser_test.go`
- `sales/service/order_service_test.go`
- `sales/service/order_stock_helper_test.go`
- Conditional: `sales/service/order_stock_helper.go`
- Conditional: `sales/repository/order_repository.go`

## Agent/Tool Routing

| Area | Owner | Review |
| --- | --- | --- |
| Parser/persistence and DetailV2 | `@backend` | `@quality-gate` |
| Stock regression | `@backend` | `@quality-gate` |
| Runtime/SQL evidence | `@backend` | `@quality-gate` |
| Coordination | `@orchestrator` | n/a |

MCP strategy: local discovery and sequential-thinking used. Context7/librarian used for Go stack compatibility. Browser/GitHub skipped: no UI/upstream dependency. Missing project stack docs recorded; plan follows local `AGENTS.md` and `sales/go.mod`.

## Executor Handoff Prompt

```text
Plan: .opencode/plans/20260713-1400-sx-2508-import-order-values.md
Scope: SX-2508 in sales import/parser, mapped DetailV2 quantity view, canonical stock presentation, regression tests, redacted runtime evidence.
Must preserve: Large/Middle/Small API qty order; unchanged global converter; mapped DetailV2 uses persisted triples; non-mapped flow stays legacy; smallest-unit stock math; service/repository boundaries; tenant/transaction safety; no secrets.
Do not touch: modules outside sales, migrations, endpoints, dependencies, caches, global converter, .env, tokens. Repository/helper edit only after focused Red proof.
Validate: Red→Green→Refactor, focused tests, `cd sales && rtk go test ./...`, `cd sales && rtk go build .`, then fresh-token curl + read-only SQL only when access exists.
Return changed files, tests/output, redacted runtime evidence or exact blocker, residual assumptions. Workers execute and report to @orchestrator. tracker updates at every status transition are mandatory, not optional bookkeeping.
```

## Execution-ready Worklist / Handoff Contract

`start_with: A1`

1. **A1** | `@backend` | Prove seams and add Red regressions. Depends on `none`. Evidence: `.opencode/evidence/20260713-1400-sx-2508-import-order-values/A1-red.log`.
2. **A2** | `@backend` | Fix parser/persistence and mapping seam. Depends on `A1`. Evidence: `.opencode/evidence/20260713-1400-sx-2508-import-order-values/A2-parser.log`.
3. **A3** | `@backend` | Fix mapped DetailV2 and stock boundary. Depends on `A2`. Evidence: `.opencode/evidence/20260713-1400-sx-2508-import-order-values/A3-detail-stock.log`.
4. **A4** | `@backend` | Run local and conditional staging checks. Depends on `A3`. Evidence: `.opencode/evidence/20260713-1400-sx-2508-import-order-values/A4-validation.md`.
5. **A5** | `@quality-gate` | Review conformance. Depends on `A4`. Evidence: `.opencode/evidence/20260713-1400-sx-2508-import-order-values/A5-quality-gate.md`.

| ID | Owner | Action | Depends | Evidence |
| --- | --- | --- | --- | --- |
| A1 | `@backend` | Prove seams; add Red regressions | none | `.opencode/evidence/20260713-1400-sx-2508-import-order-values/A1-red.log` |
| A2 | `@backend` | Fix parser/persistence and fresh mapping seam | A1 | `.opencode/evidence/20260713-1400-sx-2508-import-order-values/A2-parser.log` |
| A3 | `@backend` | Fix mapped DetailV2 and stock boundary | A2 | `.opencode/evidence/20260713-1400-sx-2508-import-order-values/A3-detail-stock.log` |
| A4 | `@backend` | Run local and conditional staging checks | A3 | `.opencode/evidence/20260713-1400-sx-2508-import-order-values/A4-validation.md` |
| A5 | `@quality-gate` | Conformance review | A4 | `.opencode/evidence/20260713-1400-sx-2508-import-order-values/A5-quality-gate.md` |

### Handoff A1
```yaml
handoff:
  task_id: 20260713-1400-sx-2508-import-order-values
  plan_id: 20260713-1400-sx-2508-import-order-values
  caller: orchestrator
  callee: backend
  scope: Add focused Red regressions and record exact persistence and mapping seams.
  claim_level: scoped
  claim_scope: Test baseline and seam findings only; no fix claim.
  source_basis: ["sales/service/order_service.go", "sales/repository/order_repository.go", "sales/service/order_service_test.go", "sales/service/order_import_parser_test.go"]
  must_preserve: ["Large Middle Small API order", "global converter", "tenant transaction boundaries"]
  do_not_touch: ["modules outside sales", "global quantity helper", ".env tokens migrations dependencies"]
  validation: ["cd sales && rtk go test ./service -run 'Test.*(Import|DetailV2|Stock)'"]
  exit_criteria: ["Focused regression cases fail for intended SX-2508 defects"]
  evidence_required: [".opencode/evidence/20260713-1400-sx-2508-import-order-values/A1-red.log"]
  depends_on: ["none"]
  context_bundle: ["sales/service/order_service.go:6629-7159", "sales/service/order_service.go:3004-3113", "sales/service/order_stock_helper.go"]
```

### Handoff A2
```yaml
handoff:
  task_id: 20260713-1400-sx-2508-import-order-values
  plan_id: 20260713-1400-sx-2508-import-order-values
  caller: orchestrator
  callee: backend
  scope: Correct import request/detail field mapping and prove current UOM lookup.
  claim_level: scoped
  claim_scope: Parser persistence only; no DetailV2 or staging claim.
  source_basis: ["sales/service/order_service.go:6629-7159", "sales/repository/order_repository.go", "sales/entity/order_detail.go"]
  must_preserve: ["parent-product persistence unless Red proof says otherwise", "no cache/dependency", "transaction-aware Store"]
  do_not_touch: ["endpoint contract", "modules outside sales", "global converter"]
  validation: ["cd sales && rtk go test ./service -run 'Test.*Import'"]
  exit_criteria: ["Import field and current-UOM tests green"]
  evidence_required: [".opencode/evidence/20260713-1400-sx-2508-import-order-values/A2-parser.log"]
  depends_on: ["A1"]
  context_bundle: ["sales/entity/order_detail.go", "sales/service/order_import_parser_test.go", ".opencode/evidence/20260713-1400-sx-2508-import-order-values/discovery.md"]
```

### Handoff A3
```yaml
handoff:
  task_id: 20260713-1400-sx-2508-import-order-values
  plan_id: 20260713-1400-sx-2508-import-order-values
  caller: orchestrator
  callee: backend
  scope: Correct mapped DetailV2 and canonical available-stock output.
  claim_level: scoped
  claim_scope: Local detail and stock tests only; no staging claim.
  source_basis: ["sales/service/order_service.go:3004-3113", "sales/service/order_stock_helper.go", "sales/pkg/conversion/quantity.go"]
  must_preserve: ["mapped persisted triples", "non-mapped legacy conversion", "smallest-unit stock arithmetic"]
  do_not_touch: ["sales/pkg/conversion/quantity.go", "parser behavior closed by A2 unless test proves need"]
  validation: ["cd sales && rtk go test ./service -run 'Test.*(DetailV2|Stock)'"]
  exit_criteria: ["Mapped/non-mapped and stock regressions green"]
  evidence_required: [".opencode/evidence/20260713-1400-sx-2508-import-order-values/A3-detail-stock.log"]
  depends_on: ["A2"]
  context_bundle: ["sales/service/order_service_test.go", "sales/service/order_stock_helper_test.go", "sales/pkg/conversion/quantity.go"]
```

### Handoff A4
```yaml
handoff:
  task_id: 20260713-1400-sx-2508-import-order-values
  plan_id: 20260713-1400-sx-2508-import-order-values
  caller: orchestrator
  callee: backend
  scope: Produce local verification and redacted staging/SQL evidence.
  claim_level: partial
  claim_scope: Staging claim only after fresh-token curl and read-only SQL checks; otherwise report exact blocker.
  source_basis: ["user supplied SX-2508 curl and SQL contract", "sales/go.mod"]
  must_preserve: ["no token or sensitive payload evidence", "read-only SQL"]
  do_not_touch: ["staging data except explicit user-approved mapping edit", ".env repository secrets"]
  validation: ["cd sales && rtk go test ./...", "cd sales && rtk go build ."]
  exit_criteria: ["Local suite/build pass and runtime evidence or blocker"]
  evidence_required: [".opencode/evidence/20260713-1400-sx-2508-import-order-values/A4-validation.md"]
  depends_on: ["A3"]
  context_bundle: [".opencode/evidence/20260713-1400-sx-2508-import-order-values/discovery.md"]
```

### Handoff A5
```yaml
handoff:
  task_id: 20260713-1400-sx-2508-import-order-values
  plan_id: 20260713-1400-sx-2508-import-order-values
  caller: orchestrator
  callee: quality-gate
  scope: Review scoped diff, evidence, quantity invariants, and security claim limits.
  claim_level: done
  claim_scope: PASS only with required evidence, validation, and no quantity regression.
  source_basis: [".opencode/plans/20260713-1400-sx-2508-import-order-values.md", ".opencode/evidence/20260713-1400-sx-2508-import-order-values"]
  must_preserve: ["quantity tenant transaction no-secret invariants"]
  do_not_touch: ["source files"]
  validation: ["review diff tests evidence staging claim scope"]
  exit_criteria: ["PASS or actionable blocker"]
  evidence_required: [".opencode/evidence/20260713-1400-sx-2508-import-order-values/A5-quality-gate.md"]
  depends_on: ["A4"]
  context_bundle: [".opencode/evidence/20260713-1400-sx-2508-import-order-values/discovery.md", ".opencode/evidence/20260713-1400-sx-2508-import-order-values/A4-validation.md"]
```

## Progress Tracking

- `tracker_path`: `.opencode/state/20260713-1400-sx-2508-import-order-values/progress.json`
- `init_command`: `python3 ~/.config/opencode/scripts/task-progress.py 20260713-1400-sx-2508-import-order-values --init --plan .opencode/plans/20260713-1400-sx-2508-import-order-values.md`
- `summary_command`: `python3 ~/.config/opencode/scripts/task-progress.py 20260713-1400-sx-2508-import-order-values --summary`
- `checklist_command`: `python3 ~/.config/opencode/scripts/task-progress.py 20260713-1400-sx-2508-import-order-values --checklist`
- `update_rules`: update before start, after completed/blocked/cancelled, whenever evidence is written, and every cross-lane handoff. Tracker updates at every status transition are mandatory, not optional bookkeeping.

| Task | Owner | Evidence | Update command |
| --- | --- | --- | --- |
| A1 | `@backend` | `.../A1-red.log` | `python3 ~/.config/opencode/scripts/task-progress.py 20260713-1400-sx-2508-import-order-values --update A1 --status in_progress --owner @backend` |
| A2 | `@backend` | `.../A2-parser.log` | `python3 ~/.config/opencode/scripts/task-progress.py 20260713-1400-sx-2508-import-order-values --update A2 --status in_progress --owner @backend` |
| A3 | `@backend` | `.../A3-detail-stock.log` | `python3 ~/.config/opencode/scripts/task-progress.py 20260713-1400-sx-2508-import-order-values --update A3 --status in_progress --owner @backend` |
| A4 | `@backend` | `.../A4-validation.md` | `python3 ~/.config/opencode/scripts/task-progress.py 20260713-1400-sx-2508-import-order-values --update A4 --status in_progress --owner @backend` |
| A5 | `@quality-gate` | `.../A5-quality-gate.md` | `python3 ~/.config/opencode/scripts/task-progress.py 20260713-1400-sx-2508-import-order-values --update A5 --status in_progress --owner @quality-gate` |

## Validation Commands

1. `rtk docker compose -f docker-compose.yml ps`
2. `cd sales && rtk go test ./service -run 'Test.*Import'`
3. `cd sales && rtk go test ./service -run 'Test.*DetailV2'`
4. `cd sales && rtk go test ./service -run 'Test.*Stock'`
5. `cd sales && rtk go test ./...`
6. `cd sales && rtk go build .`
7. `python3 ~/.config/opencode/scripts/pre-gate-smoke-check.py --project-root .`
8. `python3 ~/.config/opencode/scripts/runtime-verify.py --help`
9. Fresh-token user curl `POST /sales/v1/validate-order` case.
10. Fresh-token user curl `GET /sales/v2/orders/:ro_no` cases.
11. User supplied read-only SQL persisted-field assertions.
12. `python3 ~/.config/opencode/scripts/subagent-handoff-check.py --plan .opencode/plans/20260713-1400-sx-2508-import-order-values.md`

## Evidence Requirements

- Keep `discovery.md` and `index.json` for code/source audit.
- A1-A5 evidence paths required by worklist.
- A4 includes commands, route/assertion results, SQL fields, redaction method, and exact not-run blocker when token/access missing.
- No assets applicable.
- `SCYLLAX_TOKEN` shell-only. `BASE_URL` staging-only after user authorizes runtime check.

## Done Criteria

- All local acceptance criteria proven.
- A1-A4 evidence written; A5 quality gate PASS.
- Staging unavailable remains explicit blocker, never success claim.
- Diff remains boundary.

## Final Planning Summary

Created canonical plan and retained `.opencode/evidence/20260713-1400-sx-2508-import-order-values/discovery.md` plus `index.json` for source/claim audit. No draft retained. Key decision: preserve converter, branch mapped DetailV2 explicitly, reuse canonical stock helper. Remaining open fact: product mapping query freshness; A1/A2 prove it. Runtime proof awaits fresh token. Execution must run under next active lane permissions (`@orchestrator` then implementation lane); planner restrictions do not persist.
