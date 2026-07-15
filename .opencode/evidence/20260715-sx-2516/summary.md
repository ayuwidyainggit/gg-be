# SX-2516 Execution Summary

## Status: REMEDIATION COMPLETE; QUALITY GATE PENDING

Plan: `.opencode/plans/20260715-sx-2516.md`  
Scope: Replace data logic + validation rules; residual validation remediation only  
Claim level: scoped remediation; no final quality-gate PASS claim

## Residual fixes

- `sales/service/order_service_test.go`
  - Corrected promo fixture slab range so both consulted tabs satisfy intended promo rule.
  - Made store repository mock provide warehouse stock required by existing stock validation.
  - Removed unused test variable.
- `sales/controller/report_controller.go`
  - Changed internal scalar compatibility fields to `json:"-"`; list field remains public `json:"cust_id,omitempty"`. Duplicate JSON tags removed without changing request parsing.

## Root causes

- Promo snapshot tests used incomplete/inconsistent fixtures: one slab range excluded first-tab quantity, and store mock omitted stock data while production path validates stock.
- Vet warnings came from two fields sharing `json:"cust_id"` in each request DTO; scalar fields are internal normalization fields, not wire fields.

## Validation

Focused remediation: `2 passed`.

`rtk go vet ./...`: no issues.

Full suite: `320 passed, 9 failed`; failures are existing broader service regressions outside narrow residual slice. Exact failures recorded in `risk-remediation-full-suite.log`.

Build: success.

## Evidence

- `risk-remediation-tests.log`
- `risk-remediation-vet.log`
- `risk-remediation-full-suite.log`

## Out of scope

No migrations, dependency changes, route changes, async/staging/history changes, or modules outside `sales`.

`C1-F6-implementation.md` is superseded for this narrow remediation slice; historical unfinished C1-F6 claims remain unclaimed.
