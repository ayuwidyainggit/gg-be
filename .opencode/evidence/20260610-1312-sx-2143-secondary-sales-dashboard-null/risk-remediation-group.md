# Risk Remediation Evidence — SX-2143 Group Source Alignment

Task ID: `20260610-1312-sx-2143-secondary-sales-dashboard-null`
Tanggal: 2026-06-10 Asia/Jakarta

## Risk addressed

Quality gate sebelumnya menyisakan risk: `sum-date` sudah memakai source `sls.*`, tetapi `group` masih memakai `report.fact_orders` / `report.fact_returns`. Local DB menunjukkan facts Juni 2026 kosong, sehingga whole dashboard risk masih ada bila `group` tetap fact-based.

## Remediation decision

`group` endpoint sekarang ikut memakai source-table `sls.*` dengan date range yang sama seperti `sum-date`.

## Changed files/functions

- `sales/repository/report_repository.go`
  - `buildSecondarySalesReportGroupQuery`
    - switched group union from `report.fact_orders` / `report.fact_returns` to source tables:
      - `sls."order" o`
      - `sls.order_detail od`
      - `sls.return_det rd`
      - `sls."return" r`
    - added source-table date filters:
      - `o.invoice_date >= ?`
      - `o.invoice_date < ?`
    - group net sales now uses formula net sales exclude PPN and return branch multiplies by `-1`.
    - preserved final aliases `id`, `code`, `name`, `net_sales`.
    - preserved branch mapping for `outlet`, `salesman`, `product_category`, and `product`.
  - `SecondarySalesReportGroupOutlet`
  - `SecondarySalesReportGroupSalesman`
  - `SecondarySalesReportProductCategory`
  - `SecondarySalesReportProduct`
    - now build `dateFrom/dateTo` from `month/year` and bind date range params.
- `sales/repository/report_repository_test.go`
  - Updated group dry-run SQL tests from fact/year-filter expectations to source-table/date-range expectations.
  - Tests assert no `FROM report.fact_orders fo` / `FROM report.fact_returns fr` references in group SQL.

## DB smoke sample

Read-only local DB check for `cust_id=C260020001`, Juni 2026, `group_by=outlet` using the new source-table grouping pattern:

```text
outlet_group_rows = 11
grouped_net_sales = 5403946260
max_row_net_sales = 5000000000
```

This matches the `sum-date` `net_sales_exc_ppn` sample (`5403946260`) from implementation evidence, confirming source-based `group` can return rows locally while facts are empty.

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
semgrep scan changed source file: 0 findings, 0 errors
```

## Diff boundary

Changed after this risk remediation:

- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`

No env, token, migration, compose, package, or unrelated module files were changed.

## Remaining risk

- No authenticated HTTP smoke was run because no valid token was available in the session.
- `SecondarySalesReportReturnSumReportByMonth` remains fact-based but is no longer authoritative for summary/group data. It can be simplified in a separate cleanup.
