# SX-2214 — Total Invoice List/PDF Harus Pakai Final Order

Task ID: `20260615-1530-sx-2214-invoice-final-order`  
Readiness: `ready-for-implementation`  
Plan Quality Gate: `PASS`  
Mode: Maintenance Stability Mode  
Primary source of truth: file ini.

## Goal

Perbaiki final invoice di `sales` BE agar Invoice List, Invoice Detail, dan source data yang dipakai invoice PDF/download memakai nilai Final Order, bukan Sales Order/stale quantity.

Target evidence `SO2606100015`:

```text
Final Order gross = 18.600.000
Final Order PPN   = 1.860.000
Final Order total = 20.460.000
```

## Non-goals

- Tidak mengubah behavior Purchase Order tab/list.
- Tidak mengubah Sales Order tab/list dan Proforma Invoice kecuali terbukti berbagi path final invoice yang sama.
- Tidak mengubah semantics `GET /sales/v2/orders/{ro_no}` Final Order kecuali menemukan bug langsung.
- Tidak membuat migration kecuali schema target ternyata belum punya final columns.
- Tidak hardcode `SO2606100015`, `C260020001`, token, password, atau staging credential.
- Tidak klaim PDF fixed jika PDF generator ternyata berada di repo/service lain dan belum diverifikasi.

## Scope

Masuk scope utama:

- `sales/service/invoice_service.go`
- `sales/repository/invoice_repository.go`
- `sales/model/invoice.go`
- `sales/model/invoice_detail.go`
- `sales/entity/invoice.go`
- `sales/entity/invoice_detail.go`
- `sales/service/invoice_service_concurrency_test.go` atau test invoice baru.

Masuk scope hanya bila perlu:

- `sales/service/order_service.go` untuk reuse/helper extraction atau final header consistency check.
- `sales/model/order.go` hanya bila executor memilih reuse struct final header, bukan required.
- `sales/controller/invoice_controller.go` hanya route verification; perubahan tidak diharapkan.

Di luar scope:

- FE/PDF renderer repo lain.
- Report secondary sales.
- Stock mutation kecuali test menunjukkan compile/interface break.

## Requirements

- Final invoice line source wajib:
  - `qty1_final`, `qty2_final`, `qty3_final`
  - `sell_price_final1`, `sell_price_final2`, `sell_price_final3`
  - `promo_final1..5`
  - `disc_value_final`
  - `vat_value_final`
- Sales Order fields tidak boleh menjadi source final invoice line:
  - `qty1..3`
  - `sell_price1..3`
  - `amount`
  - `disc_value`
  - `vat_value`
- Null numeric harus dihitung sebagai `0`.
- Formula final line:

```text
line_gross = qty1_final*sell_price_final1 + qty2_final*sell_price_final2 + qty3_final*sell_price_final3
line_promo = promo_final1 + promo_final2 + promo_final3 + promo_final4 + promo_final5
line_discount = disc_value_final
line_vat = vat_value_final
line_net = line_gross - line_promo - line_discount + line_vat
```

- Total final invoice:

```text
total_gross = SUM(line_gross)
total_promo = SUM(line_promo)
total_discount = SUM(line_discount)
total_vat = SUM(line_vat)
total_invoice = total_gross - total_promo - total_discount + total_vat
```

- Promo/discount sign convention harus dikonfirmasi dari existing data. Default: promo/discount positif lalu dikurangkan.
- `AmountFinal` tidak boleh jadi source utama kecuali DB evidence membuktikan selalu sesuai formula.
- DB final invoice state harus konsisten: final header columns di `sls.order` harus merefleksikan formula saat invoice generated.

## Acceptance Criteria

- `GET /v1/invoices?q=SO2606100015&is_invoice=true` atau filter ekuivalen mengembalikan total invoice `20.460.000` jika staging data tidak berubah.
- `GET /v1/invoices/SO2606100015` detail mengembalikan line `TP-012` dengan final qty/price/amount setara `16.500.000`, bukan stale `1.500.000`.
- Invoice list/detail footer fields konsisten dengan final formula.
- Source data untuk PDF/download tidak lagi memakai `qty1/2/3` Sales Order bila PDF mengambil dari invoice API ini.
- `GET /sales/v2/orders/SO2606100015` Final Order tidak berubah secara semantic.
- Proforma Invoice path tetap memakai Sales Order/proforma convention existing.
- Purchase Order dan Sales Order response tidak berubah.
- Null promo fields menghasilkan promo `0`, bukan null.
- Automated test mencakup mismatch Sales Order total `3.960.000` vs Final Order total `20.460.000` dan memilih Final Order.
- Tidak ada credential/token/SO sample hardcode di production code.

## Existing Patterns/Reuse

Evidence lengkap: `.opencode/evidence/20260615-1530-sx-2214-invoice-final-order/discovery.md`.

Reuse:

- `InvoiceService.List`, `Details`, `Detail` sudah fetch details via `InvoiceRepository.FindDetail`; bisa hitung final totals service-side tanpa query group kompleks.
- `InvoiceService.BulkUpdate` sudah fetch detail final untuk stock update; bisa hitung final totals di transaction yang sama.
- `model.OrderDetailRead` menunjukkan final field naming yang harus ditambahkan ke `model.InvoiceDetRead` bila belum ada.
- `model.OrderList` punya final header fields; `model.InvoiceList` belum.
- `recomputePromoStateForTab(... promoSnapshotTabFinalOrder ...)` di `order_service.go` menjadi referensi convention final header.
- `invoice_service_concurrency_test.go` memberi pola mock service/repository.
- `report_repository_test.go` memberi pola SQL-string regression bila executor memilih SQL aggregate.

Root cause repo-backed:

- `model.InvoiceDetRead` hanya expose `Qty1/2/3`, `SellPrice1/2/3`, `Amount`, `DiscValue`, `VatValue` dari Sales Order fields.
- `InvoiceService` automap raw invoice details langsung ke response.
- `InvoiceRepository.FindAllByCustId` memilih `sls.order.total`, bukan `total_final` atau calculated final total.
- `InvoiceService.BulkUpdate` generate invoice number/status tanpa recompute final invoice amount columns.

## Constraints

- Ikuti `AGENTS.md`: shell workflow pakai `rtk` prefix di repo ini.
- Validasi dari module `sales`, bukan root.
- Preserve Controller → Service → Repository → DB.
- Repository write harus tx-aware via `extractTx(ctx)` pattern existing.
- Tenant filter `cust_id` wajib tetap ada.
- Jangan simpan token/cURL Authorization/Jira credential ke source, fixture, log, commit, atau artifact publik.
- Existing models pakai `float64`; gunakan style existing untuk minimal diff. Jangan introduce decimal dependency tanpa repo evidence kuat.

## Risks

- PDF/download endpoint final invoice tidak ditemukan di `sales`; mungkin ada di FE/BFF/renderer. Fix BE API source belum otomatis membuktikan PDF final tanpa runtime trace.
- Mengoverwrite non-final header columns (`total`, `sub_total`, `disc_value`, `vat_value`) bisa mengubah Sales Order semantics. Prefer update/read final columns untuk invoice response.
- `InvoiceRepository.Update` dengan struct pointer fields dapat menulis monetary fields dari request jika tidak disanitasi. Executor harus clear/ignore stale request totals for final invoice generation.
- Existing `OrderService` fallback branch around `order_service.go:5837-5863` computes `totalFinal = subTotalFinal - discValueFinal + vatValueFinal` without promo. Jika SX-2214 data path hits branch with promo, root fix may need include `promo_final` there too.
- `promo_value_final` header may not equal sum `promo_final1..5` if stale; final invoice helper should calculate from detail fields.
- Global tests may fail from pre-existing repo issues; record targeted evidence.

## Decisions/Assumptions

- Keputusan: source canonical invoice final adalah calculated detail final formula, bukan `amount`, bukan `amount_final`, dan bukan legacy header `total`.
- Keputusan: list/detail response boleh expose final totals through existing JSON keys (`sub_total`, `promo_value`, `disc_value`, `vat_value`, `total`) because endpoint is invoice context.
- Keputusan: generated invoice should update final header columns (`sub_total_final`, `promo_value_final`, `disc_value_final`, `vat_value_final`, `total_final`) from same helper inside `BulkUpdate`.
- Keputusan: do not overwrite non-final Sales Order header money columns unless executor finds current invoice list/PDF consumers cannot read final columns and product owner accepts invoice-context overwrite.
- Asumsi slice-safe: `sls.order` has final header columns and `sls.order_detail` has final detail columns based on local models/docs.
- Asumsi slice-safe: PDF/download consumes `GET /v1/invoices/:ro_no` or `GET /v1/invoices/details`; if not, executor must locate actual consumer before claiming PDF fixed.
- Question gate: tidak ada blocking question; user supplied scope, formula, guardrails, and sample evidence. Remaining uncertainty is implementation discovery, not business decision.

## Execution Source of Truth

Precedence executor:

1. Instruksi eksplisit terbaru dari user.
2. Safety/security repo: no secrets, no tokens, no `.env`, use `rtk`, validate in `sales`.
3. Non-negotiable Implementation Invariants.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria dan Done Criteria.
6. Implementation Steps.
7. Executor recommendations within Diff Boundary.

Jika konflik, follow higher source and record conflict in verification evidence.

## Non-negotiable Implementation Invariants

- Final invoice calculation must use `*_final` detail fields.
- `COALESCE`/nil-safe per field, not around whole promo sum.
- Do not use `amount` or `amount_final` as canonical formula source.
- Do not change Proforma Invoice calculation path unless direct shared bug proven.
- Do not alter Purchase Order/Sales Order response semantics.
- Do not hardcode Jira sample identifiers in production code.
- Do not copy credentials/tokens from Jira or local env into code/test/artifacts.
- Keep DB writes in `InvoiceService.BulkUpdate` transaction.
- Do not claim PDF parity without runtime/API trace showing PDF source uses fixed BE data or direct PDF route fixed.

## Do Not / Reject If

Reject/revert implementation if:

- It only changes display formatting but still reads `qty1/2/3` for final invoice.
- It uses `COALESCE(promo_final1 + ... + promo_final5, 0)` and can null out promo sum.
- It fixes `SO2606100015` via hardcoded `ro_no`/`cust_id`.
- It writes token/password/bearer into test fixture, log, docs, or code.
- It changes Purchase Order/Proforma behavior without failing regression evidence.
- It changes stock quantity/valuation as side effect without explicit evidence.
- It claims final PDF fixed while no PDF path or API source verification exists.

## Diff Boundary

Allowed source/test changes:

- `sales/service/invoice_service.go`
- `sales/repository/invoice_repository.go`
- `sales/model/invoice.go`
- `sales/model/invoice_detail.go`
- `sales/entity/invoice.go`
- `sales/entity/invoice_detail.go`
- `sales/service/invoice_service_concurrency_test.go`
- Optional new test/helper under `sales/service/` if project style supports it.

Allowed only with direct evidence:

- `sales/service/order_service.go`
- `sales/model/order.go`
- `sales/controller/invoice_controller.go`
- `sales/repository/order_repository.go`

Evidence paths:

- `.opencode/evidence/20260615-1530-sx-2214-invoice-final-order/`

Any out-of-boundary change must be reverted or justified in final evidence before quality gate.

## TDD/Test Plan

TDD required: yes. This is production money calculation defect.

Existing patterns:

- `sales/service/invoice_service_concurrency_test.go` uses mocked invoice repo/stock repo/transaction.
- `sales/service/order_service_test.go` has many final-order amount tests, useful reference.

Red step:

- Add unit/service test before fix, e.g. `TestInvoiceFinalTotalsUsesFinalOrderFields`.
- Fixture detail has Sales Order qty/price producing `3.960.000`, and Final Order qty/price/vat producing `20.460.000` total.
- Include null `promo_final2..5` regression; expected promo `0`.
- Add service list/detail test if practical:
  - repository returns `InvoiceList.Total = 3960000` and details final total `20460000`.
  - `InvoiceService.List` response `Total` must be `20460000`.
  - detail line response `Qty1/2/3`, `SellPrice1/2/3`, `Amount`, `VatValue`, `NetValue` must reflect final fields.

Green step:

- Add final fields to `model.InvoiceDetRead`:
  - `SellPriceFinal1/2/3`
  - `PromoFinal1..5`
  - `DiscValueFinal`
  - `VatValueFinal`
  - `AmountFinal` only for diagnostic/fallback, not formula source.
- Add final header fields to `model.Invoice` and `model.InvoiceList` if needed:
  - `SubTotalFinal`
  - `DiscValueFinal`
  - `PromoValueFinal`
  - `PromoBgValueFinal`
  - `VatValueFinal`
  - `TotalFinal`
- Add helper, preferably in `sales/service/invoice_service.go` or small `invoice_amount.go`:

```go
type invoiceFinalLineAmount struct {
    Gross float64
    PromoPrimary float64
    PromoSecondary float64
    Discount float64
    VAT float64
    Net float64
}
```

- Helper input should be `model.InvoiceDetRead`, use nil/zero-safe fields where struct supports pointers or zero values.
- Apply helper in `Detail`, `List`, `Details`, and `BulkUpdate`.
- In `BulkUpdate`, compute final totals after `FindDetail`, set final header columns on `invoiceModel`, and clear/avoid stale request monetary fields if they are not final.

Refactor step:

- Keep one formula helper; avoid formula copy in list/detail/BulkUpdate.
- If SQL aggregate is chosen instead of service helper, add repository tests verifying exact final field SQL and per-field `COALESCE`.

Edge cases:

- `promo_final1 = 0`, `promo_final2..5 = NULL` → promo `0`.
- final prices null/zero with Sales Order prices present → final invoice should not silently use Sales Order unless status/business rule proves invoice generated before Final Order.
- promo/discount negative values → avoid double-negative after data convention check.
- item_type promo rows → verify whether final invoice includes/excludes via existing `ItemType` convention; do not invent.

Commands:

```bash
rtk go test ./service -run 'TestInvoice.*Final|TestInvoiceBulkUpdate'
rtk go test ./...
```

## Implementation Steps

1. Add/extend invoice test mock to capture update payload and list/detail responses.
2. Write failing tests for final formula, list total override, detail line mapping, and null promo fields.
3. Extend invoice models with final detail/header fields needed by helper and DB update.
4. Implement nil-safe final amount helper using final detail fields and promo sum per field.
5. Update `InvoiceService.Detail` to map invoice line response from final fields and override header totals from helper.
6. Update `InvoiceService.List` and `Details` to override invoice-context totals from final helper after `FindDetail`.
7. Update `InvoiceService.BulkUpdate` to compute final totals within transaction and persist final header columns before/with invoice number/status update.
8. Ensure `InvoiceRepository.FindAllByCustId` selects needed final header fields if service or response uses them.
9. Search for any remaining final invoice consumer still using Sales Order fields; fix only invoice-context path.
10. Run targeted tests, then `rtk go test ./...` from `sales`.
11. Manual verify staging/local with secure token outside artifacts.
12. Record whether PDF/download path is covered by fixed API or still external follow-up.

## Expected Files to Change

Likely:

- `sales/model/invoice.go`
- `sales/model/invoice_detail.go`
- `sales/service/invoice_service.go`
- `sales/service/invoice_service_concurrency_test.go` or new `sales/service/invoice_service_test.go`

Possible:

- `sales/repository/invoice_repository.go`
- `sales/entity/invoice.go`
- `sales/entity/invoice_detail.go`

## Agent/Tool Routing

- `@orchestrator`: route implementation and integrate evidence.
- `@fixer`: bounded code/test implementation.
- `@explorer`: optional if PDF route/source still missing after first search.
- `@quality-gate`: final signoff due money calculation + defect fix.
- `@artifact-planner`: no source edits; this plan only.

## Executor Handoff Prompt

```text
Implement SX-2214 in ScyllaX Sales BE using `.opencode/plans/20260615-1530-sx-2214-invoice-final-order.md` as source of truth.

Scope: fix final invoice list/detail/PDF source data so final invoice uses `qty*_final`, `sell_price_final*`, `promo_final1..5`, `disc_value_final`, `vat_value_final`. Preserve Purchase Order, Sales Order, and Proforma semantics.

must_preserve:
- no secrets/tokens/credentials in code/tests/artifacts
- no hardcoded `SO2606100015`/`C260020001` in production logic
- Controller → Service → Repository → DB layering
- tx-aware writes in `InvoiceService.BulkUpdate`
- per-field nil-safe promo sum
- no `amount`/`amount_final` as canonical invoice formula source

do_not_touch:
- non-invoice services
- Purchase Order/Proforma behavior unless failing evidence proves shared path
- `.env`, secrets, lockfiles unless dependency not added

Start with failing tests for final formula and invoice list/detail mismatch. Then implement shared final invoice amount helper and apply to `Detail`, `List`, `Details`, and `BulkUpdate` final header persistence. Validate in `sales` with `rtk go test ./service -run 'TestInvoice.*Final|TestInvoiceBulkUpdate'` and `rtk go test ./...`.

Return evidence: changed files, tests run, exact pass/fail output, manual API checks, PDF/download source trace, and migration/data repair note for existing generated invoices.
```

## Execution-ready Worklist / Handoff Contract

`start_with`: `T1`

### T1 — Write regression tests

- `depends_on`: none
- `owner/lane`: `@fixer`
- `action`: add failing tests for final invoice formula/list/detail/null promo.
- `validation`: `rtk go test ./service -run 'TestInvoice.*Final'`
- `exit_criteria`: tests fail before production fix for expected stale Sales Order total.
- `blocking_status`: `ready`
- `requires_user_decision`: `no`
- `must_preserve`: no secrets; no hardcoded sample in production.
- `do_not_touch`: non-invoice source files.
- `evidence_update`: record test names and red failure in evidence.
- `exit_verification`: red failure clearly shows final vs Sales Order mismatch.

### T2 — Add final invoice amount helper and model fields

- `depends_on`: `T1`
- `owner/lane`: `@fixer`
- `action`: extend invoice detail/header models and add one shared final amount helper.
- `validation`: `rtk go test ./service -run 'TestInvoice.*Final'`
- `exit_criteria`: helper unit tests pass; null promo fields pass.
- `blocking_status`: `ready`
- `requires_user_decision`: `no`
- `must_preserve`: formula uses final fields and per-field nil-safe values.
- `do_not_touch`: Proforma code path.
- `evidence_update`: record helper location and formula source.
- `exit_verification`: targeted helper tests green.

### T3 — Apply helper to invoice detail/list responses

- `depends_on`: `T2`
- `owner/lane`: `@fixer`
- `action`: update `Detail`, `List`, and `Details` response mapping/totals to final invoice values.
- `validation`: `rtk go test ./service -run 'TestInvoice.*Final|TestInvoiceList|TestInvoiceDetail'`
- `exit_criteria`: response total and line amount use Final Order fixture.
- `blocking_status`: `ready`
- `requires_user_decision`: `no`
- `must_preserve`: invoice context only; no Purchase/Sales Order response changes.
- `do_not_touch`: `OrderService.PrintProformaInvoice` unless direct shared defect found.
- `evidence_update`: record exact response fields overridden.
- `exit_verification`: targeted service tests green.

### T4 — Persist final invoice totals on generate

- `depends_on`: `T3`
- `owner/lane`: `@fixer`
- `action`: update `BulkUpdate` to compute and persist final header totals during invoice generation.
- `validation`: `rtk go test ./service -run 'TestInvoiceBulkUpdate|TestInvoice.*Final'`
- `exit_criteria`: captured repository update contains final header totals and does not persist stale request totals as canonical invoice final.
- `blocking_status`: `ready`
- `requires_user_decision`: `no`
- `must_preserve`: transaction boundary and invoice number retry behavior.
- `do_not_touch`: stock mutation behavior except compile-required interface changes.
- `evidence_update`: record DB columns updated.
- `exit_verification`: concurrency tests still pass.

### T5 — Verify PDF/download source path

- `depends_on`: `T3`
- `owner/lane`: `@explorer` then `@fixer` if source in repo.
- `action`: locate final invoice PDF/download consumer and verify it uses fixed invoice API/data; fix in-scope route if found.
- `validation`: grep/code trace plus manual endpoint/PDF check if credentials available.
- `exit_criteria`: evidence states either fixed API covers PDF or external repo follow-up needed.
- `blocking_status`: `ready`
- `requires_user_decision`: `no`
- `must_preserve`: no token stored.
- `do_not_touch`: external repo/source unless user explicitly authorizes.
- `evidence_update`: record route/source trace.
- `exit_verification`: no unsupported PDF claim.

### T6 — Run full validation and quality gate

- `depends_on`: `T1,T2,T3,T4,T5`
- `owner/lane`: `@quality-gate`
- `action`: review diff, tests, security, scope, and evidence.
- `validation`: `rtk go test ./...` from `sales`; optional API smoke with secure token.
- `exit_criteria`: pass or documented pre-existing failures with targeted green tests.
- `blocking_status`: `ready`
- `requires_user_decision`: `no`
- `must_preserve`: acceptance criteria and guardrails.
- `do_not_touch`: no artifact-only planner files except evidence update.
- `evidence_update`: final validation summary.
- `exit_verification`: quality gate PASS or explicit blockers.

## Validation Commands

From repo root first, when runtime needed:

```bash
rtk docker compose -f docker-compose.yml ps
```

From `sales` module:

```bash
rtk go test ./service -run 'TestInvoice.*Final|TestInvoiceBulkUpdate'
rtk go test ./...
```

Manual DB validation after implementation, using secure DB access only:

```sql
WITH final_calc AS (
  SELECT
    SUM(
      COALESCE(qty1_final, 0) * COALESCE(sell_price_final1, 0) +
      COALESCE(qty2_final, 0) * COALESCE(sell_price_final2, 0) +
      COALESCE(qty3_final, 0) * COALESCE(sell_price_final3, 0)
    ) AS gross,
    SUM(
      COALESCE(promo_final1, 0) + COALESCE(promo_final2, 0) +
      COALESCE(promo_final3, 0) + COALESCE(promo_final4, 0) +
      COALESCE(promo_final5, 0)
    ) AS promo,
    SUM(COALESCE(disc_value_final, 0)) AS discount,
    SUM(COALESCE(vat_value_final, 0)) AS vat
  FROM sls.order_detail
  WHERE cust_id = 'C260020001'
    AND ro_no = 'SO2606100015'
)
SELECT gross, promo, discount, vat, gross - promo - discount + vat AS expected_invoice_total
FROM final_calc;
```

Manual API validation, no token in artifact:

```bash
curl -H "Authorization: Bearer <secure-local-token>" "https://best.scyllax.online/sales/v1/invoices?q=SO2606100015&is_invoice=true"
curl -H "Authorization: Bearer <secure-local-token>" "https://best.scyllax.online/sales/v1/invoices/SO2606100015"
curl -H "Authorization: Bearer <secure-local-token>" "https://best.scyllax.online/sales/v2/orders/SO2606100015"
```

## Evidence Requirements

Implementation evidence must include:

- Root cause exact file/function/query confirmed.
- Before/after API values for Invoice List and Invoice Detail.
- Test output for targeted and full sales module tests.
- DB query comparison final formula vs stored final header columns.
- PDF/download trace: route, source API/table, or external follow-up.
- Note whether existing generated invoice rows need regeneration/backfill.
- No credential/token in evidence.

Source strategy used for plan:

- Used repo-local docs/code and user Jira evidence.
- Skipped official docs/web/GitHub because defect is repo-local business logic.
- Skipped browser capture because runtime token not available to planner; manual verification required after implementation.

## Done Criteria

- Code uses one shared final invoice formula helper.
- Invoice list/detail final totals match Final Order formula.
- Generate invoice persists final header totals.
- Tests cover mismatch and null promo regression.
- `rtk go test ./...` run from `sales`, or failures documented as pre-existing with targeted pass.
- Manual staging/local verification recorded without secrets.
- Migration/data repair note documented.
- `@quality-gate` review passes or blockers documented.

## Final Planning Summary

Artifacts created/kept:

- `.opencode/plans/20260615-1530-sx-2214-invoice-final-order.md` — primary source of truth.
- `.opencode/evidence/20260615-1530-sx-2214-invoice-final-order/discovery.md` — kept because executor needs file/function evidence.
- `.opencode/evidence/20260615-1530-sx-2214-invoice-final-order/index.json` — kept as evidence manifest.

Draft cleanup:

- No draft artifact created; nothing stale to delete.

Key decisions:

- Canonical final invoice formula comes from detail `*_final` fields.
- Invoice response may override invoice-context legacy JSON keys with final values.
- Invoice generation should persist final header columns, not rely on request totals.
- PDF claim requires source trace because local repo has no obvious final invoice PDF generator.

Assumptions:

- `sales` BE invoice API is source for Invoice List and likely PDF/download.
- Final columns exist in DB, matching local models/docs.
- Secure staging token/DB access will be provided by local environment, not stored.

Open questions:

- Actual final invoice PDF/download generator path remains to be traced during implementation.
- Existing generated invoice rows may need regeneration/backfill if stored final header columns are already stale.

Readiness: `ready-for-implementation`.
