# SX-2214 — Total Invoice List/PDF Harus Sama dengan Final Order

Task ID: `20260617-1236-sx-2214-final-invoice-total`  
Readiness: `ready-for-implementation`  
Plan Quality Gate: `PASS_FOR_SLICE`  
Mode: Maintenance Stability Mode  
Primary source of truth: file ini.

## Goal

Pastikan final invoice di ScyllaX Sales BE memakai angka Final Order untuk Invoice List, Invoice Detail, dan source data yang dipakai invoice PDF/download.

Target evidence `SO2606100015` bila staging data belum berubah:

```text
Final Order gross = 18.600.000
Final Order PPN   = 1.860.000
Final Order total = 20.460.000
```

## Non-goals

- Tidak ubah Purchase Order behavior.
- Tidak ubah Sales Order behavior.
- Tidak ubah Proforma Invoice behavior kecuali terbukti shared broken path.
- Tidak ubah `GET /sales/v2/orders/{ro_no}` Final Order semantics.
- Tidak tambah schema/migration kecuali target branch belum punya final columns.
- Tidak hardcode `SO2606100015`, `C260020001`, token, password, bearer.
- Tidak klaim PDF fixed tanpa trace PDF source/consumer.

## Scope

Masuk scope:

- `sales/service/invoice_service.go`
- `sales/service/invoice_amount.go`
- `sales/service/invoice_service_test.go`
- `sales/repository/invoice_repository.go`
- `sales/model/invoice.go`
- `sales/model/invoice_detail.go`
- `sales/entity/invoice.go`

Masuk scope hanya bila target branch belum punya current fix:

- `sales/entity/invoice_detail.go`
- `sales/controller/invoice_controller.go` untuk route verification saja.
- `sales/service/order_service.go` hanya bila final header generation dependency langsung terbukti.

Di luar scope:

- FE changes.
- PDF renderer repo lain.
- Secondary sales report.
- Stock valuation behavior kecuali compile/test menunjukkan direct break.

## Requirements

Final invoice line source wajib:

```text
qty1_final, qty2_final, qty3_final
sell_price_final1, sell_price_final2, sell_price_final3
promo_final1, promo_final2, promo_final3, promo_final4, promo_final5
disc_value_final
vat_value_final
```

Final invoice line source dilarang:

```text
qty1, qty2, qty3
sell_price1, sell_price2, sell_price3
amount
disc_value
vat_value
```

Canonical formula:

```text
line_gross =
  COALESCE(qty1_final, 0) * COALESCE(sell_price_final1, 0) +
  COALESCE(qty2_final, 0) * COALESCE(sell_price_final2, 0) +
  COALESCE(qty3_final, 0) * COALESCE(sell_price_final3, 0)

line_promo_primary = COALESCE(promo_final1, 0)
line_promo_secondary =
  COALESCE(promo_final2, 0) +
  COALESCE(promo_final3, 0) +
  COALESCE(promo_final4, 0) +
  COALESCE(promo_final5, 0)
line_discount = COALESCE(disc_value_final, 0)
line_vat = COALESCE(vat_value_final, 0)
line_net = line_gross - line_promo_primary - line_promo_secondary - line_discount + line_vat
```

Total invoice:

```text
total_gross = SUM(line_gross)
total_promo = SUM(line_promo_primary + line_promo_secondary)
total_discount = SUM(line_discount)
total_vat = SUM(line_vat)
total_invoice = total_gross - total_promo - total_discount + total_vat
```

## Acceptance Criteria

- Invoice List untuk `SO2606100015` total `20.460.000` jika staging data belum berubah.
- Invoice Detail untuk `SO2606100015` line `TP-012` memakai final amount `16.500.000`, bukan stale `1.500.000`.
- Invoice List, Detail, and generated invoice final header total konsisten dari helper sama.
- Null promo fields menghasilkan promo `0`, bukan null.
- `GET /sales/v2/orders/SO2606100015` Final Order tetap sama.
- Proforma Invoice tests tidak regresi.
- Purchase Order and Sales Order response tidak regresi.
- No credential/token/SO hardcode in production code.
- PDF/download hanya boleh diklaim fixed setelah source PDF memakai fixed final invoice API/data verified.

## Existing Patterns/Reuse

Discovery evidence: `.opencode/evidence/20260617-1236-sx-2214-final-invoice-total/discovery.md`.

Current local code sudah punya fix candidate:

- `sales/service/invoice_amount.go`:
  - `calculateInvoiceFinalLineAmount`
  - `calculateInvoiceFinalTotals`
- `sales/service/invoice_service.go`:
  - `Detail` recompute response totals from final details.
  - `List` recompute row totals from final details.
  - `Details` recompute row totals from final details.
  - `BulkUpdate` writes `SubTotalFinal`, `PromoValueFinal`, `DiscValueFinal`, `VatValueFinal`, `TotalFinal`.
  - `mapInvoiceFinalDetailResponse` maps qty/price/amount/net from final fields.
- `sales/model/invoice_detail.go` `InvoiceDetRead` already has final fields.
- `sales/service/invoice_service_test.go` already has invoice final regression tests.

Reuse order: Reuse existing helper > extend tests/verification > only create new code if target branch lacks fix.

## Constraints

- Repo rule: shell workflows in this repo use `rtk` prefix.
- Validate inside `sales` module.
- Preserve Controller → Service → Repository → DB.
- Repository writes stay tx-aware via `extractTx(ctx)`.
- Tenant filter `cust_id` stays required.
- No secrets in code, fixture, logs, commits, artifacts.
- Repo local is not git worktree; executor must verify real target branch separately.
- Existing project uses `float64`; do not add decimal dependency in this defect unless broader money decision exists.

## Risks

- PDF/download generator not found in `sales`; likely FE/BFF/another service consumes invoice API. Claim limit required.
- If deployed target branch lacks current local fix, source edits needed.
- If `promo_value_final` header differs from line `promo_final1..5`, helper must follow line fields per Jira formula.
- If final fields null for older invoices, final invoice generation may need invalid-state handling; do not silently fall back to Sales Order for final invoice without product confirmation.
- Full test suite may expose unrelated failures; document targeted pass and unrelated failures.

## Decisions/Assumptions

- Decision: final invoice source of truth is calculated detail final fields, not `amount`, not `amount_final`, not stale header `total`.
- Decision: invoice API may return final totals in existing JSON keys `sub_total`, `promo_value`, `disc_value`, `vat_value`, `total`, because endpoint context is invoice final.
- Decision: generated invoice writes final header columns only; non-final Sales Order header money should not be overwritten unless target code requires and approval exists.
- Assumption slice-safe: PDF/download consumes invoice API or `sls.order` final columns. If not true, executor must locate PDF path before final claim.
- Question gate: no blocking user question. User provided formula, sample, scope, guardrails.

## Execution Source of Truth

Precedence for executor:

1. Latest explicit user instruction.
2. Security/repo rules: no secrets, use `rtk`, validate in `sales`, preserve tenant rules.
3. Non-negotiable Implementation Invariants.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria and Done Criteria.
6. Implementation Steps.
7. Follow-up recommendations.

If conflict exists, follow higher source and record conflict in verification evidence.

## Non-negotiable Implementation Invariants

- Final invoice calculation must use `*_final` detail fields.
- Nil/COALESCE must be per promo field, not around whole promo expression.
- `amount` and `amount_final` are not canonical formula source.
- Proforma Invoice path remains unchanged unless direct shared broken code path proven.
- Purchase Order and Sales Order response semantics remain unchanged.
- `SO2606100015` and `C260020001` may appear only in manual verification notes, never production logic.
- No Jira token, bearer, password, `.env` value in source/test/artifacts.
- `InvoiceService.BulkUpdate` writes stay inside service transaction.
- PDF fixed claim requires proof of PDF source/consumer.

## Do Not / Reject If

Reject/rework if:

- Code still maps final invoice response from `qty1/2/3` Sales Order fields.
- Code uses `COALESCE(promo_final1 + ... + promo_final5, 0)` and can null out whole sum.
- Fix hardcodes sample SO/customer.
- Fix copies credentials into tests/logs/artifacts.
- Fix changes Proforma/Purchase/Sales Order without regression evidence.
- Fix changes stock valuation as side effect without explicit evidence.
- Final response claims PDF fixed without PDF route/consumer verification.

## Diff Boundary

Allowed source/test file groups:

- `sales/service/invoice_service.go`
- `sales/service/invoice_amount.go`
- `sales/service/invoice_service_test.go`
- `sales/model/invoice.go`
- `sales/model/invoice_detail.go`
- `sales/entity/invoice.go`
- `sales/repository/invoice_repository.go`

Allowed evidence paths:

- `.opencode/evidence/20260617-1236-sx-2214-final-invoice-total/`

Out-of-boundary changes must be reverted or justified in final verification evidence.

## TDD/Test Plan

TDD required: yes. Money calculation production defect.

Existing tests already present in current local code:

- `TestInvoiceFinalLineAmountUsesFinalOrderFieldsAndNullPromo`
- `TestInvoiceFinalLineAmountSubtractsAllFinalPromoFields`
- `TestInvoiceDetailUsesFinalOrderFieldsAndHeaderTotals`
- `TestInvoiceListUsesFinalOrderTotals`
- `TestInvoiceDetailsUsesFinalOrderTotals`
- `TestInvoiceBulkUpdatePersistsFinalHeaderTotals`

Red step if target branch lacks tests:

- Add fixture where Sales Order total is `3.960.000`, Final Order total is `20.460.000`.
- Assert invoice list/detail/generate uses Final Order total.
- Assert null promo fields sum to `0`.

Green step:

- Add/reuse shared final invoice calculation helper.
- Map detail response qty/price/amount/net from final fields.
- Recompute list/detail totals from helper.
- Persist generated invoice final header columns from helper.

Refactor step:

- Keep helper single-purpose and invoice-scoped.
- Avoid SQL aggregate duplication unless performance issue measured.

Commands:

```bash
cd sales
rtk go test ./service -run 'TestInvoice'
rtk go test ./...
```

Current targeted validation already run:

```text
rtk go test ./service -run 'TestInvoice'
Go test: 8 passed in 1 packages
```

## Implementation Steps

1. Confirm target branch has current local fix candidate.
2. If missing, add `InvoiceDetRead` final fields and final header fields in invoice models.
3. Add/reuse `sales/service/invoice_amount.go` helper.
4. Update `InvoiceService.Detail`, `List`, `Details` to call helper and map final detail response.
5. Update `InvoiceService.BulkUpdate` to persist final header totals from helper inside transaction.
6. Add/reuse regression tests listed in TDD plan.
7. Run targeted and full module tests.
8. Manually verify staging/local API with secure token from env, not artifact.
9. Trace PDF/download source; fix consumer path if it does not use invoice API/final fields.
10. Route final signoff to `@quality-gate`.

## Expected Files to Change

If target branch already equals current local state:

- No source changes needed.
- Evidence update only.

If target branch lacks fix:

- `sales/service/invoice_amount.go`
- `sales/service/invoice_service.go`
- `sales/service/invoice_service_test.go`
- `sales/model/invoice.go`
- `sales/model/invoice_detail.go`

Likely no controller/repository logic change, because `FindDetail` already selects `sls.order_detail.*` and tenant filters exist.

## Agent/Tool Routing

- `@orchestrator`: integrate execution, confirm target branch, coordinate verification.
- `@fixer`: bounded source/test changes if target branch lacks current fix.
- `@explorer`: locate PDF/download route if not obvious.
- `@quality-gate`: final security/regression signoff.
- `@librarian`: not needed unless external PDF/client docs needed.

## Executor Handoff Prompt

```text
Implement/verify SX-2214 in ScyllaX Sales BE. Use `.opencode/plans/20260617-1236-sx-2214-final-invoice-total.md` as source of truth. Preserve final invoice formula from Final Order fields only: `qty*_final`, `sell_price_final*`, `promo_final1..5`, `disc_value_final`, `vat_value_final`. Do not use Sales Order fields for final invoice totals/lines. Reuse current local helper `sales/service/invoice_amount.go` if present. Keep Purchase Order, Sales Order, and Proforma behavior unchanged. Do not hardcode `SO2606100015`, `C260020001`, tokens, passwords, or bearer values. Validate from `sales` with `rtk go test ./service -run 'TestInvoice'` and `rtk go test ./...`. Manually verify `SO2606100015` via secure env token only; do not record token. Do not claim PDF fixed until PDF source/consumer is traced and shown to use fixed final invoice data. Return changed files, test output, manual expected vs actual totals, and any remaining PDF/data repair risk.
```

## Execution-ready Worklist / Handoff Contract

`start_with`: `T1`

| ID | Action | depends_on | owner/lane | validation | exit criteria | status | requires_user_decision | must_preserve | do_not_touch | evidence_update | exit_verification |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T1 | Confirm target branch state against current local fix candidate | none | `@explorer`/`@orchestrator` | inspect listed files | Know whether source edits needed | ready | no | no source claim without target branch check | credentials, `.env` | note file presence/diff status | evidence says current/fix missing |
| T2 | Add/reuse invoice final calculation helper | T1 | `@fixer` | unit test helper | helper computes `20.460.000` fixture and null promo `0` | ready | no | final fields only | Proforma code | helper/test notes | targeted unit test pass |
| T3 | Map invoice detail/list/details responses from final helper | T2 | `@fixer` | service tests | list/detail line totals match Final Order | ready | no | no Sales Order fields for final invoice | Purchase/Sales/Proforma paths | response mapping notes | service tests pass |
| T4 | Persist generated invoice final header totals in `BulkUpdate` | T3 | `@fixer` | bulk update service test | `SubTotalFinal`, `PromoValueFinal`, `DiscValueFinal`, `VatValueFinal`, `TotalFinal` set from helper | ready | no | transaction boundary | non-final header overwrite unless approved | generated invoice notes | bulk update test pass |
| T5 | Run module validation | T4 | `@orchestrator` | `rtk go test ./service -run 'TestInvoice'`; `rtk go test ./...` | targeted pass; full pass or documented unrelated failures | ready | no | use `rtk` | secrets/logging tokens | command outputs | command evidence saved/summarized |
| T6 | Manual API verify sample SO | T5 | `@orchestrator` | secure local/staging request | list/detail totals equal `20.460.000` when data unchanged | ready | no | no token in artifact | credential copy | expected vs actual totals | response fields summarized only |
| T7 | Trace PDF/download source | T5 | `@explorer`/`@orchestrator` | route/source search or runtime trace | PDF source uses fixed final invoice data, or separate fix task opened | ready | no | no PDF fixed claim without proof | unrelated renderer changes | source/route notes | PDF claim allowed or blocked |
| T8 | Final quality gate | T6,T7 | `@quality-gate` | evidence review | PASS or explicit remaining risk | ready | no | security + scope | none | quality-gate notes | signoff recorded |

## Validation Commands

From repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

From `sales` module:

```bash
rtk go test ./service -run 'TestInvoice'
rtk go test ./...
```

Optional manual SQL/API checks must not record secrets:

```sql
WITH line AS (
  SELECT
    COALESCE(qty1_final, 0) * COALESCE(sell_price_final1, 0) +
    COALESCE(qty2_final, 0) * COALESCE(sell_price_final2, 0) +
    COALESCE(qty3_final, 0) * COALESCE(sell_price_final3, 0) AS gross,
    COALESCE(promo_final1, 0) + COALESCE(promo_final2, 0) + COALESCE(promo_final3, 0) +
    COALESCE(promo_final4, 0) + COALESCE(promo_final5, 0) AS promo,
    COALESCE(disc_value_final, 0) AS discount,
    COALESCE(vat_value_final, 0) AS vat
  FROM sls.order_detail
  WHERE cust_id = 'C260020001'
    AND ro_no = 'SO2606100015'
)
SELECT SUM(gross) AS total_gross,
       SUM(promo) AS total_promo,
       SUM(discount) AS total_discount,
       SUM(vat) AS total_vat,
       SUM(gross - promo - discount + vat) AS total_invoice
FROM line;
```

## Evidence Requirements

Executor must record:

- Files changed or confirmation no source changes needed.
- Targeted test output.
- Full module test output or unrelated failure detail.
- Manual API expected vs actual for `SO2606100015`, with no token/password.
- PDF/download source trace result.
- Backfill/regeneration need for existing invoices, if any.

Source strategy:

- Used: repo-local docs/code, existing SX-2214 artifacts, user Jira summary/formula, local BrowserOS extracted docs.
- Skipped: Context7, GitHub, web search; issue is repo-local business logic.
- Skipped browser screenshot; manual API/PDF trace planned by executor.

## Done Criteria

- Invoice list/detail tests pass with Final Order fixture.
- `BulkUpdate` final header test passes.
- `rtk go test ./service -run 'TestInvoice'` passes.
- `rtk go test ./...` passes or unrelated failures documented.
- Manual sample shows list/detail total `20.460.000` if data unchanged.
- PDF/download source verified or remaining risk explicitly documented.
- No secrets or sample identifiers in production logic.
- `@quality-gate` signoff completed for material BE money calculation change.

## Final Planning Summary

Artifacts created:

- `.opencode/plans/20260617-1236-sx-2214-final-invoice-total.md`
- `.opencode/evidence/20260617-1236-sx-2214-final-invoice-total/discovery.md`
- `.opencode/evidence/20260617-1236-sx-2214-final-invoice-total/index.json`

Key decisions:

- Final invoice must use calculated final detail fields, not stale Sales Order fields or `amount_final`.
- Current local code already appears to contain a SX-2214 fix candidate and tests.
- Keep `PASS_FOR_SLICE` because PDF/download route not found in `sales`; BE API slice is ready, PDF claim needs trace.

Assumptions:

- Target implementation branch should match current local files or receive same changes.
- Staging `SO2606100015` data remains unchanged for `20.460.000` expected total.

Open questions:

- Where exactly is invoice PDF/download generated if not in `sales`? This is not blocking BE API fix, but blocks PDF parity claim.

Validation already performed:

```text
rtk go test ./service -run 'TestInvoice'
Go test: 8 passed in 1 packages
```

Cleanup:

- Empty draft directory `.opencode/draft/20260617-1236-sx-2214-final-invoice-total/` removed after synthesis.
- Evidence kept because it contains operational validation and PDF trace gap.
