# Implementation Evidence — SX-2143

Task ID: `20260610-1312-sx-2143-secondary-sales-dashboard-null`
Tanggal: 2026-06-10 Asia/Jakarta

## Root cause verified

Local `ggn_scyllax` menunjukkan `report.fact_orders` dan `report.fact_returns` kosong untuk `cust_id=C260020001` pada Juni 2026, sementara source `sls."order"`, `sls.order_detail`, `sls."return"`, dan `sls.return_det` memiliki data valid. Karena summary sebelumnya membaca fact tables, endpoint `sum-date` dapat kosong walaupun source order/return ada.

## Keputusan implementasi

Berdasarkan persetujuan user, `sum-date` diubah agar summary memakai source-table `sls.*` dengan date range dari `month/year`:

- `date_from = first day of month 00:00:00 UTC`
- `date_to = first day of next month 00:00:00 UTC`
- filter SQL: `invoice_date >= date_from AND invoice_date < date_to`

Group endpoint tidak diubah karena sudah memiliki `year` filter, return subtraction, dan `code` enhancement.

## Changed files/functions

- `sales/repository/report_repository.go`
  - `SecondarySalesReportSumReportByMonth`
    - switched from `report.fact_orders` / `report.fact_returns` to source `sls.*` aggregation.
    - order branch: `sls."order" o` + `sls.order_detail od`.
    - return branch: `sls.return_det rd` + `sls."return" r` + `sls."order" o`.
    - filters: `cust_id IN ?`, `data_status IN (6,7)`, date range.
    - numeric aggregates use `COALESCE`.
    - final summary subtracts returns for gross, ppn, net sales; keeps discount as order + return discount/promo.
    - `last_update` selection handles nulls before `GREATEST`.
- `sales/repository/report_repository_test.go`
  - Added `assertSecondarySalesSummaryDateVars` helper.
  - Replaced fact-based summary SQL test with `TestSecondarySalesReportSumReportByMonthSQLUsesSourceTablesAndDateRange`.
- `sales/service/report_service.go`
  - `SecondarySalesReportSumReportByMonth` now preserves summary query `LastUpdate`; legacy return-summary helper only fills it if summary last update is nil.
- `sales/service/report_service_test.go`
  - Extended mapping test to ensure summary `LastUpdate` is not overwritten by legacy return helper.
- `sales/controller/so_controller_test.go`
  - Added success and invalid-query controller tests for `sum-date` and `group` query parsing/validation.

## SQL behavior covered by tests

Repository dry-run test asserts summary SQL contains:

- `FROM sls."order" o`
- `JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = o.cust_id`
- `FROM sls.return_det rd`
- `JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id`
- `JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id`
- `o.invoice_date >= ?` / `o.invoice_date < ?` date range placeholders
- COALESCE gross/discount/PPN fragments for order and return
- final return subtraction math
- no `FROM report.fact_orders fo`
- no `FROM report.fact_returns fr`

Vars assert June 2026 date range:

- `2026-06-01 00:00:00 UTC`
- `2026-07-01 00:00:00 UTC`

## Local DB smoke sample

Read-only local calculation using the same source-table formula for `C260020001`, Juni 2026 returned:

```text
total_gross_sale        = 5405350000
total_discount_promo    = 1442480.0000
total_ppn               = 540394626.0000
net_sales_exc_ppn       = 5403946260
net_sales                = 5944340886
total_salesman          = 5
total_outlet            = 9
total_product           = 12
qty                     = 48
qty_return              = 1
return_rate             = 2.08333333333333333300
net_sales_return        = 630630
last_update             = 2026-06-10 10:23:37.054127+07
```

Manual HTTP API smoke was not run because no valid Authorization token was provided in the session. DB smoke verifies the source-table query has non-empty local data.

## Validation commands

From `sales`:

```bash
rtk go test ./repository -run 'TestSecondarySalesReport(SumReportByMonth|ReturnSumReportByMonth|Group)'
# Go test: 10 passed in 1 packages

rtk go test ./service -run 'TestSecondarySalesReport(SumReportByMonth|GroupSales)'
# Go test: 11 passed in 1 packages

rtk go test ./controller -run 'TestSecondaryReportSales'
# Go test: 13 passed in 1 packages

rtk go test ./service ./repository ./controller
# Go test: 241 passed in 3 packages

rtk go test ./...
# Go test: 244 passed in 22 packages

git diff --check
# no output
```

Security/static scan:

```text
semgrep scan changed source files: 0 findings, 0 errors
```

## Diff boundary check

Changed source/test files are within allowed boundary:

- `sales/controller/so_controller_test.go`
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`
- `sales/service/report_service.go`
- `sales/service/report_service_test.go`

No env, token, migration, compose, package, or unrelated module files were changed.

## Quality gate

Final `@quality-gate` result: `PASS_WITH_RISKS`.

- No blocker.
- Controller-test remediation completed.
- Remaining risk is scoped: `sum-date` is fixed and validated; `group` remains fact-based by approved scope decision, so do not claim full dashboard parity for Juni 2026 unless facts/backfill or group source-query work is handled.

## Remaining risk / follow-up

- `SecondarySalesReportReturnSumReportByMonth` remains legacy fact-based. The service still calls it for fallback `last_update`, but source summary now carries return qty/net/last_update itself. This is not blocking; it can be simplified later.
- Group endpoint remains fact-based. It has correct `year`, return subtraction, and `code`, but local group response for Juni 2026 will still depend on facts/backfill. User-approved implementation decision was to fix `sum-date` source-table first.
- FE should send `&year=2026`; missing year fallback remains `time.Now().Year()` for compatibility.
