# A5 SX-2508 quality gate

Reviewer: `@quality-gate` (read-only).
Date: 2026-07-13.
Source basis: plan + A1-A4 evidence + current code at `sales/`.

## Verdict
`PASS_WITH_RISKS`

## Scope checked
- Plan: `.opencode/plans/20260713-1400-sx-2508-import-order-values.md`
- Evidence: `discovery.md`, `A1-red.log`, `A2-parser.log`, `A3-detail-stock.log`, `A4-validation.md`
- Code: `sales/service/order_service.go`, `sales/model/order.go`, `sales/repository/order_repository.go`, `sales/entity/order_detail.go`, `sales/service/order_service_test.go`, `sales/pkg/conversion/quantity.go`
- Progress tracker: `.opencode/state/20260713-1400-sx-2508-import-order-values/progress.json`

## Decision
- Scoped SX-2508 fix conforms: parser + DetailV2 quantity contract + stock display + no secret leak.
- Full `go test ./...` still red on 2 pre-existing promo-consult tests. Handoff allowed `PASS_WITH_RISKS` for this case only.
- Runtime/staging proof absent by design. Not blocker per handoff.

## Findings
- HIGH: full module test suite not green — 2 failing DetailV2 promotion-consult tests remain.
  - basis: A4-validation.md:11,19-33,44; A3-detail-stock.log:4-21
  - why not blocker: failure assertions are promo remarks / promo injection, not qty orientation; pre-existing.
- MEDIUM: cannot prove `sales/repository/order_repository.go` was untouched by this task (no VCS).
  - basis: file mtime Jul 13 19:48; current content already matches plan claim.
  - note: evidence limit, not code defect.
- LOW: tracker/evidence bookkeeping incomplete.
  - basis: progress.json evidence arrays empty for A2/A3/A4; index.json still `status: planning`.
  - impact: weakens artifact hygiene only; does not change code correctness.

## Independent verification
- `sales/pkg/conversion/quantity.go:15-23` — internal order `Qty1=small, Qty2=middle, Qty3=large`. Strongest available evidence: unchanged.
- `DetailV2` quantity contract at `sales/service/order_service.go:3071-3095`:
  - mapped (`ro.IsSalesMapping == true`): skip `ConvToQtyConversion`, preserve stored `Qty1/2/3` and `Qty1Final/2Final/3Final`.
  - non-mapped: remap converter output to API orientation: `Qty1=qtyConversion.Qty3` (large), `Qty2=qtyConversion.Qty2` (middle), `Qty3=qtyConversion.Qty1` (small).
  - logic correct against converter contract.
- `IsSalesMapping` model field added minimally to `model.Order` and `model.OrderList`. Tag: ``gorm:"column:is_sales_mapping" json:"is_sales_mapping"``.
- Parser persistence at `sales/service/order_service.go:6758-6785, 6982-6989, 7086-7093`:
  - `QtyPo: nil` explicit.
  - `SellPriceSystem2/3` from parent `SellPrice3`.
  - `PromoSo1/PromoFinal1` from imported promo.
  - `VatValue/VatValueFinal` from imported PPN.
- Secret/token scan: no token, Authorization header, bearer, password, or secret in A1-A4 evidence.
- Diff boundary: in-scope changes only in `sales/service/order_service.go`, `sales/service/order_service_test.go`, `sales/model/order.go`. `sales/pkg/conversion/quantity.go` and `sales/service/order_stock_helper.go` untouched. `order_repository.go` non-modification not provable without VCS.

## Required before absolute PASS
- Get full `cd sales && rtk go test ./...` green, or formally quarantine the 2 pre-existing tests with repo-approved policy.
- Recommended: update progress/evidence bookkeeping (done by orchestrator in this run).

## Remediation worklist
- HIGH — owner `@backend`: investigate/fix or formally classify the 2 promotion-consult tests; validate with `cd sales && rtk go test ./...`. requires_user_decision: no if classified as out-of-scope ticket.
- MEDIUM — owner `@orchestrator`: attach external diff/snapshot for non-git worktree, or downgrade boundary claim. requires_user_decision: no.
- LOW — owner `@orchestrator`: refresh tracker + index (handled in this run). requires_user_decision: no.

## Escalation
None.
