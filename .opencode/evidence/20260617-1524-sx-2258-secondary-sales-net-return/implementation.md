# SX-2258 implementation evidence

## Scope
- Target module: `sales`
- `pjp-sales` not changed because target confirmed `sales`
- Source of truth: `.opencode/plans/20260617-1524-sx-2258-secondary-sales-net-return.md`
- Repo docs used: `.opencode/docs/ARCHITECTURE.md`
- Stack/playbook docs status: `.opencode/docs/PROJECT_STACK.md`, `.opencode/docs/PROJECT_COMMANDS.md`, `.opencode/docs/FRAMEWORK_PLAYBOOK.md`, `.opencode/docs/PROJECT_DETECTED_TOOLS.md` not present; plan already documented absence

## Changed files
- `sales/entity/report.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`
- `sales/service/report_service_test.go`

## Before / after formulas

### Summary endpoint `GET /v1/reports/secondary-sales/sum-date`
- `qty`
  - before: `os.qty AS qty`
  - after: `COALESCE(os.qty, 0) - COALESCE(rs.qty_return, 0) AS qty`
- `total_discount_promo`
  - before: `(os.discount_promo + rs.discount_promo) AS total_discount_promo`
  - after: `COALESCE(os.discount_promo, 0) - COALESCE(rs.discount_promo, 0) AS total_discount_promo`
- order discount/promo source
  - before: `COALESCE(od.promo_value_final, 0) + COALESCE(od.disc_value_final, 0)`
  - after: `COALESCE(od.disc_value_final, 0) + COALESCE(od.promo_final1, 0) + COALESCE(od.promo_final2, 0) + COALESCE(od.promo_final3, 0) + COALESCE(od.promo_final4, 0) + COALESCE(od.promo_final5, 0)`
- return discount/promo source
  - before: unchanged semantics but old summary still added result to order side
  - after: `COALESCE(rd.disc_value, 0) + COALESCE(rd.promo_value, 0)` and subtracted from order side

### Trend endpoint `GET /v1/reports/secondary-sales/trend-sales`
- `total_discount_promo`
  - before: `COALESCE(os.discount_promo, 0) + COALESCE(rs.discount_promo, 0) AS total_discount_promo`
  - after: `COALESCE(os.discount_promo, 0) - COALESCE(rs.discount_promo, 0) AS total_discount_promo`

## Filter handling
- Added backward-compatible optional summary payload fields:
  - `from`
  - `to`
  - `outlet_ids`
  - `salesman_ids`
  - `pro_ids`
- Date precedence:
  - if `from` and `to` present: use `str.UnixTimestampToUtcTime`
  - else: derive first day of `month/year` and next month boundary
- Order filters use:
  - `o.invoice_date`
  - `o.outlet_id`
  - `o.salesman_id`
  - `od.pro_id`
- Return filters use:
  - `o.invoice_date`
  - `r.outlet_id`
  - `r.salesman_id`
  - `rd.product_id`
- Query remains parameterized with `Raw(query, params...)`

## PPN / return semantics
- `total_ppn` preserved
- `net_sales_return` preserved
- `return_rate` preserved
- `qty_return` preserved as return qty, not net qty

## TDD / validation

### Red / regression coverage added or updated
- repository SQL regression rejects old summary formula `os.qty AS qty`
- repository SQL regression rejects old summary plus formula for `total_discount_promo`
- repository SQL regression requires return invoice join
- repository SQL regression requires optional summary filter aliases on order and return sides
- arithmetic regression test covers:
  - `150 - 16 = 134`
  - `1_500_000 - 261_260 = 1_238_740`
  - null promo part treated as zero
- trend SQL regression rejects plus formula and requires subtract formula
- service regression checks optional summary filters propagate to repository

### Commands and results
1. Repo root
   - command: `rtk docker compose -f docker-compose.yml ps`
   - result: command ran; compose showed warning about obsolete `version` key and empty service table in current environment
2. `sales/`
   - command: `rtk go test ./repository -run 'SecondarySalesReport(SumReportByMonth|TrendSales)|SX2258'`
   - result: pass (`4 passed`)
3. `sales/`
   - command: `rtk go test ./service -run SecondarySalesReportSumReportByMonth`
   - result: pass (`6 passed`)
4. `sales/`
   - command: `rtk go test ./...`
   - result: pass (`265 passed in 22 packages`)

## Staging / local dataset verification
- exact dataset check status: blocked by safe DB/API access not available in this session
- staging/local exact dataset response not run; no credentials copied, stored, or inferred

## Quality gate snapshot
- status: `PASS_WITH_RISKS`
- pass basis:
  - plan scope followed in `sales` only
  - summary and trend formulas match plan invariants
  - optional filters added backward-compatible and tested
  - requested validations passed
- residual risk:
  - exact QA dataset verification still pending safe internal DB/API access
