# Plan SX-2143 + SX-2165 Secondary Sales BE

Readiness: `ready-for-implementation`

Plan Quality Gate: `PASS_FOR_SLICE`

Source of truth: `.opencode/plans/20260605-1647-secondary-sales-defects.md`

## Goal

Fix backend defect untuk Secondary Sales dashboard dan export:

- SX-2143: dashboard summary/group tidak kosong untuk principal/distributor filter `cust_id`, `month`, `year`, dan response punya field summary FE.
- SX-2165: row `trx_type = RETURN` di export dan extract fact return membawa promo return dari `sls.return_det.promo_value`.

## Non-goals

- Tidak ubah auth middleware/JWT contract.
- Tidak hardcode data Jira (`C260020001`, month/year, token, user).
- Tidak ubah behavior transaksi order kecuali null-safe `COALESCE` untuk field diskon/promo.
- Tidak buat migrasi schema kecuali implementer menemukan field model/DB belum ada.
- Tidak copy credential/evidence Jira ke repo.

## Scope

Target module: `sales`.

Target endpoint/flow:

- `GET /v1/reports/secondary-sales/sum-date`
- `GET /v1/reports/secondary-sales/group`
- `POST /v1/reports/secondary-sales` export publish + RMQ subscriber Excel generation
- Extract path `POST /v1/extract/secondary-sales` untuk `report.fact_returns`

## Requirements

- Dashboard request menerima `month`, optional `year`, optional query `cust_id`.
- Requested `cust_id` harus tetap divalidasi scope parent via existing `resolveSecondaryDashboardCustID`.
- Summary response tambah:
  - `total_ppn`
  - `net_sales_exc_ppn`
  - `net_sales` sebagai net sales include PPN.
- Summary menghitung order dan return dari fact table:
  - gross: order gross minus return gross.
  - discount/promo: order discount+special plus return discount+special, sesuai Jira reference.
  - net exc PPN: order net exc PPN minus return net exc PPN.
  - net sales: order net inc PPN minus return net inc PPN.
  - `qty_return`, `net_sales_return`, `return_rate` tetap tersedia.
- Group endpoint untuk `outlet`, `salesman`, `product_category`, `product` mengikutkan return sebagai negative net sales.
- Export return row mengisi promo return.
- Extract return fact mengisi promo return supaya dashboard fact tidak hilang promo.

## Acceptance Criteria

- `GET /sales/v1/reports/secondary-sales/sum-date?month=5&cust_id=C260020001` memakai effective `cust_id` dari query jika user punya akses parent-scope.
- Request Month=6 Year=2026 BU `C260020001` membaca `dt.month = 6` dan `dt.year = 2026`.
- Response sum-date punya JSON keys:
  - `total_gross_sale`
  - `total_discount_promo`
  - `total_ppn`
  - `net_sales_exc_ppn`
  - `net_sales`
  - metrics existing salesman/outlet/product/qty/return.
- No return data tidak menyebabkan panic/error dan summary tetap 1 row dengan `0` return values.
- Group endpoint semua `group_by` tetap menghasilkan data dan return mengurangi net sales group terkait.
- Return export row dengan `rd.promo_value = 6500` tidak lagi menghasilkan promo/special discount `0`.
- Order export row tetap formula existing.
- Null `disc_value`/`promo_value` diperlakukan `0`.

## Existing Patterns/Reuse

- Reuse `sales/service/report_service.go:resolveSecondaryDashboardCustID` untuk principal → BU access guard.
- Reuse `sales/service/report_service.go:resolveSecondaryDashboardYear` untuk `Year *int` fallback.
- Reuse repository report interface signatures already taking `(custID, month, year)`.
- Reuse SQL union builder `buildSecondarySalesUnionQuery` for export/list.
- Reuse tests in `sales/service/report_service_test.go` and add repository SQL/model tests near `sales/repository/report_repository_test.go`.

No matching utility found for combined order/return dashboard aggregation; create repo helper/raw SQL per existing repository style.

## Constraints

- Repo policy: commands under `sales/` and `rtk`-prefixed.
- Tenant rule: keep `cust_id` filters and parent-customer validation.
- Layering: controller parses, service decides/access, repository queries.
- Do not store Jira token/password/QA credential.

## Risks

- VAT sign: Jira snippet inconsistent. Default plan uses `order - return` for exposed `total_ppn` and `net_sales`; add QA/PM check if numbers differ.
- `SpecialDiscount` ambiguity: Jira says `disc_value + promo_value`, but repo export has separate `SpecialDiscount` and `Discount` columns. Default repo-compatible fix: `special_discount = COALESCE(rd.promo_value,0)`, `discount = COALESCE(rd.disc_value,0)`. If QA expects one Excel column to contain total discount+promo, ask PM before changing column semantics.
- Fact tables can be stale if extract/cron failed; implementation must validate extract too.

## Decisions/Assumptions

- Assumption slice-safe: current code already fixed dashboard query `cust_id` parsing/service auth enough; implementation should verify, not duplicate in controller.
- Decision: add summary response fields in entity/model/service; no breaking removal of old fields.
- Decision: group charts use `net_sales_exclude_ppn` and return as negative values.
- Open question: final VAT sign and SpecialDiscount visual semantics require QA/PM numeric validation if local data mismatch.

## TDD/Test Plan

TDD required: yes. Backend query/report logic and export mapping are production behavior.

Existing test patterns:

- `sales/service/report_service_test.go` has mock repository tests for cust scope, year fallback, return rate.
- `sales/repository/report_repository_test.go` has SQL-scope test for `ExistsCustomerInParentScope`.

First failing/regression tests:

1. Repository SQL test for `buildSecondarySalesUnionQuery`: return branch contains `COALESCE(rd.promo_value, 0) AS special_discount` and `COALESCE(rd.disc_value, 0) AS discount`, no `0 AS special_discount`.
2. Repository SQL test for `GetReportSecondarySalesReportReturn` query or extracted helper: no `0 AS special_discount`; includes return promo.
3. Service mapping test: `SumReportByMonthModel` with PPN/net exc maps to `SumReportByMonthModelResp.TotalPPN` and `NetSalesExcPPN`.
4. Repository/dashboard test if feasible with sqlmock: summary SQL uses both `report.fact_orders` and `report.fact_returns` and returns combined aliases.
5. Group SQL tests for all group functions include `report.fact_returns` or union return negative branch.

Green step:

- Update SQL/model/entity/service until tests pass.

Refactor step:

- Extract repeated group union SQL builder if duplication becomes high.
- Keep queries readable; do not over-abstract before tests pass.

Edge cases:

- order only.
- order + return + promo return.
- no return.
- return promo/discount null.
- order qty zero with return qty non-zero.
- unauthorized requested `cust_id`.

Commands:

```bash
rtk go test ./service -run 'TestSecondarySalesReport'
rtk go test ./repository -run 'Test.*SecondarySales'
rtk go test ./...
```

## Implementation Steps

1. Update response/model fields:
   - `sales/entity/report.go`: add `TotalPPN float64 json:"total_ppn"`, `NetSalesExcPPN float64 json:"net_sales_exc_ppn"` to `SumReportByMonthModelResp`.
   - `sales/model/report.go`: add `TotalPPN`, `NetSalesExcPPN` to `SumReportByMonthModel` with matching `gorm` columns.

2. Fix SX-2165 export/list union:
   - In `sales/repository/report_repository.go` `buildSecondarySalesUnionQuery` return branch replace:
     - `0 AS special_discount` → `COALESCE(rd.promo_value, 0) AS special_discount`
     - `rd.disc_value AS discount` → `COALESCE(rd.disc_value, 0) AS discount`
   - Confirm both `SecondarySalesUnionPagination` and `SecondarySalesUnion` use this builder.

3. Fix return extract fact:
   - In `GetReportSecondarySalesReportReturn`, replace `0 AS special_discount` with `COALESCE(rd.promo_value, 0) AS special_discount`.
   - Replace `rd.disc_value AS discount` with `COALESCE(rd.disc_value, 0) AS discount`.
   - Null-safe `rd.vat_value`/`rd.total` calculations if tests show null risk.

4. Fix dashboard summary query:
   - Replace order-only `SecondarySalesReportSumReportByMonth` with combined order/return aggregate SQL.
   - Return aliases: `total_gross_sale`, `total_discount_promo`, `total_ppn`, `net_sales_exc_ppn`, `net_sales`, `total_salesman`, `total_outlet`, `total_product`, `qty`, `qty_return`, `net_sales_return`, `last_update` if model extended or keep return method for qty/return.
   - Prefer one repository query for both order and return to avoid inconsistent last-update/return values.
   - If keeping separate return method, still compute main totals using both order and return.

5. Map new summary fields:
   - `sales/service/report_service.go`: map `sumReportModel.TotalPPN` → `data.TotalPPN`; `sumReportModel.NetSalesExcPPN` → `data.NetSalesExcPPN`; `sumReportModel.NetSales` stays include PPN.
   - Keep return rate guard `if sumReportModel.Qty > 0`.

6. Fix group queries:
   - For `SecondarySalesReportGroupOutlet`, `SecondarySalesReportGroupSalesman`, `SecondarySalesReportProductCategory`, `SecondarySalesReportProduct`, use union:
     - order: `fo.net_sales_exclude_ppn`
     - return: `fr.net_sales_exclude_ppn * -1`
   - Join matching dimensions in each branch.
   - Preserve year filter.

7. Add/update tests from TDD plan.

8. Run validation commands and capture output.

9. Manual validation if runtime/DB available:
   - Run compose check from repo root.
   - Hit sum-date/group endpoints with non-secret auth.
   - Generate export and inspect RETURN row.
   - Query DB promo return without storing credentials.

## Expected Files to Change

- `sales/entity/report.go`
- `sales/model/report.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`

Optional if repository tests need helper exposure:

- `sales/repository/report_repository.go` helper extraction only.

## Agent/Tool Routing

- `@fixer`: implementation and tests.
- `@explorer`: extra repo discovery if tests reveal hidden helpers.
- `@oracle`: review VAT/sign/column semantics if PM/QA numbers mismatch.
- `@quality-gate`: final security/tenant/query/test review.

## Execution-ready Worklist / Handoff Contract

`start_with`: `T1`

| Task | Action | depends_on | owner/lane | validation | exit criteria | status | requires_user_decision |
| --- | --- | --- | --- | --- | --- | --- | --- |
| T1 | Add failing tests for return promo export/extract SQL | none | `@fixer` | `rtk go test ./repository -run 'Test.*SecondarySales'` | Tests fail before SQL fix | ready | no |
| T2 | Add failing tests for summary response mapping fields | none | `@fixer` | `rtk go test ./service -run 'TestSecondarySalesReportSumReport'` | Tests fail before model/entity/service fix | ready | no |
| T3 | Update return promo SQL in export union and extract return | T1 | `@fixer` | `rtk go test ./repository -run 'Test.*SecondarySales'` | RETURN branch uses `COALESCE(rd.promo_value, 0)` and no `0 AS special_discount` remains for return export/extract | ready | no |
| T4 | Add model/entity fields and service mapping for `total_ppn`, `net_sales_exc_ppn` | T2 | `@fixer` | `rtk go test ./service -run 'TestSecondarySalesReportSumReport'` | Response struct emits new fields; existing fields preserved | ready | no |
| T5 | Replace summary aggregate with order+return calculation | T4 | `@fixer` | `rtk go test ./repository -run 'Test.*SecondarySalesReportSum'` | Summary totals use fact orders and fact returns, null-safe no-return case | ready | no |
| T6 | Replace group queries with order positive + return negative union | T5 | `@fixer` | `rtk go test ./repository -run 'Test.*SecondarySalesReportGroup'` | All four group functions include return negative branch and year filter | ready | no |
| T7 | Run full sales validation | T3,T4,T5,T6 | `@fixer` | `rtk go test ./...` | Full module tests pass or failures documented unrelated | ready | no |
| T8 | Manual API/export validation with safe credentials/runtime if available | T7 | `@fixer` | curl/export/DB query evidence | sum-date has fields, group returns data, RETURN export promo non-zero | ready | no |
| T9 | Final review | T8 | `@quality-gate` | review changed files + test evidence | Tenant/security/query risks accepted or flagged | ready | no |

Blocked branch:

- If QA insists Excel `SpecialDiscount` must equal `disc_value + promo_value` while `Discount` column also exists, pause before T3 finalization and ask PM/QA to choose column semantics. Current default keeps repo semantics: `SpecialDiscount=promo_value`, `Discount=disc_value`.

## Validation Commands

From repo root first:

```bash
rtk docker compose -f docker-compose.yml ps
```

From `sales/`:

```bash
rtk go test ./service -run 'TestSecondarySalesReport'
rtk go test ./repository -run 'Test.*SecondarySales'
rtk go test ./...
```

Manual API, no token stored:

```bash
curl '{{base_url}}/sales/v1/reports/secondary-sales/sum-date?month=5&cust_id=C260020001' \
  -H 'Accept: application/json'
```

DB check, no password/token stored:

```sql
SELECT
    rd.return_no,
    rd.product_id,
    rd.disc_value,
    rd.promo_value,
    rd.vat_value
FROM sls.return_det rd
WHERE rd.cust_id = 'C260020001'
  AND COALESCE(rd.promo_value, 0) <> 0;
```

## Evidence Requirements

Implementation must produce:

- Root cause summary with exact changed files.
- Before/after API response or log for sum-date showing new fields.
- Group endpoint sample for all `group_by` values or test evidence.
- Export sample row before/after for `trx_type = RETURN`, showing promo not `0`.
- Test command output from `sales/`.
- Note on VAT sign and SpecialDiscount/Discount column semantics.

Source strategy used by this plan:

- Repo-local evidence only. Enough because defects map to local SQL/service contracts.
- External docs skipped: no version-sensitive API behavior.
- Browser/UI evidence skipped: backend data/export issue.

## Done Criteria

- Tests pass for repository/service changes.
- `0 AS special_discount` no longer appears in return export/extract paths.
- Sum-date response includes required fields and no divide-by-zero return rate.
- Group queries include return branch.
- Requested `cust_id` remains access-scoped.
- Manual validation evidence captured or blocker documented.
- `@quality-gate` signoff complete for tenant/query/security risk.

## Final Planning Summary

Artifacts created:

- `.opencode/plans/20260605-1647-secondary-sales-defects.md` — source of truth.
- `.opencode/evidence/20260605-1647-secondary-sales-defects/discovery.md` — kept because it records repo evidence and unresolved sign/column risks useful for implementation.

Artifacts deleted:

- None. No stale draft created.

Key decisions:

- Use existing service-level `cust_id` authorization guard.
- Treat return promo as `special_discount = rd.promo_value` and return discount as `discount = rd.disc_value` unless PM/QA overrides column semantics.
- Use order positive + return negative for group net sales.

Assumptions:

- VAT exposed as order minus return for repo-consistent net sales. Validate against QA numbers.
- Fact table dashboard remains intended data source; extract fix required so future facts include promo return.

Open questions:

- Does QA want Excel column `SpecialDiscount` to show only promo (`promo_value`) or total discount+promo (`disc_value + promo_value`) while `Discount` column still exists?
- Should `total_ppn` expose order minus return PPN, or Jira snippet's inconsistent sum `os.vat + rs.vat`?

Readiness:

- `ready-for-implementation` with two validation-sensitive open questions. Safe first slice can implement repo-compatible defaults and confirm with QA/PM if numeric mismatch appears.
