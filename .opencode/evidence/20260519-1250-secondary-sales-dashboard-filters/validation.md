# Validation — Task 118/119 Secondary Sales Dashboard Filters

Task id: `20260519-1250-secondary-sales-dashboard-filters`
Tanggal: `2026-05-19`

## Commands run

Dari `/Users/ujang/Projects/Geekgarden/scylla-be/sales`:

```bash
rtk go test ./service -run 'TestSecondarySalesReport(SumReportByMonth|GroupSales)'
rtk go test ./repository -run 'TestSecondarySalesReport(SumReportByMonth|ReturnSumReportByMonth|Group)'
rtk go test ./controller -run 'TestSecondaryReportSales(SumMonth|Group)ReturnsForbiddenForUnauthorizedCustID|TestParseDownloadSalesmanIDs'
rtk go test ./repository -run 'TestExistsCustomerInParentScopeSQLRequiresActiveChild|TestSecondarySalesReport(SumReportByMonth|ReturnSumReportByMonth|Group)'
rtk go test ./...
```

## Test results

- `rtk go test ./service -run 'TestSecondarySalesReport(SumReportByMonth|GroupSales)'`
  - `Go test: 9 passed in 1 packages`
- `rtk go test ./repository -run 'TestSecondarySalesReport(SumReportByMonth|ReturnSumReportByMonth|Group)'`
  - `Go test: 7 passed in 1 packages`
- `rtk go test ./controller -run 'TestSecondaryReportSales(SumMonth|Group)ReturnsForbiddenForUnauthorizedCustID|TestParseDownloadSalesmanIDs'`
  - `Go test: 8 passed in 1 packages`
- `rtk go test ./repository -run 'TestExistsCustomerInParentScopeSQLRequiresActiveChild|TestSecondarySalesReport(SumReportByMonth|ReturnSumReportByMonth|Group)'`
  - `Go test: 8 passed in 1 packages`
- `rtk go test ./...`
  - `Go test: 169 passed in 22 packages`

## DB schema validation

Connection used:

```text
host=103.28.219.73 port=25431 user=postgres dbname=scylla_citus_dev sslmode=disable
```

Observed schema:

- `report.dim_dates`
  - `id bigint`
  - `day smallint`
  - `month smallint`
  - `year smallint`
- `report.fact_orders`
  - contains `cust_id`, `date_id`, `gross_sale`, `special_discount`, `discount`, `net_sales_exclude_ppn`, `salesman_id`, `outlet_id`, `pro_id`, `qty`, `extracted_at`
- `report.fact_returns`
  - contains `cust_id`, `date_id`, `net_sales_exclude_ppn`, `qty`, `extracted_at`
- `smc.m_customer`
  - contains `cust_id`, `parent_cust_id`, `is_active`, `is_del`

## DB scope validation

Child scope sample under parent `C26002`:

```text
C260020001|C26002|f|t
```

Interpretation:

- `cust_id = C260020001`
- `parent_cust_id = C26002`
- `is_del = false`
- `is_active = true`

Rule after follow-up fix:

- child `cust_id` accepted only when `cust_id = ? AND parent_cust_id = ? AND is_del = false AND is_active = true`
- inactive child rows must be rejected even if parent relation matches and `is_del = false`
- active scope requirement now explicit in repository SQL and repository test coverage

## Live data validation

### Sample sum-date target

Chosen live sample:

- `cust_id = C260020001`
- `month = 4`
- `year = 2026`

Orders aggregate:

```text
469500000.0000|2132000.0000|477368000.0000|1|8|13|235|2026-04-21 00:01:01.800669+00
```

Mapped:

- `total_gross_sale = 469500000`
- `total_discount_promo = 2132000`
- `net_sales = 477368000`
- `total_salesman = 1`
- `total_outlet = 8`
- `total_product = 13`
- `qty = 235`
- `last_update = 2026-04-21 00:01:01.800669+00`

Returns aggregate:

```text
0|0|
```

Mapped:

- `qty_return = 0`
- `net_sales_return = 0`
- `last_update = null`

Expected merged service behavior:

- `return_rate = 0`
- `last_update` stays from orders because return last_update null

### Sample group target

Chosen live sample:

- `cust_id = C260020001`
- `month = 4`
- `year = 2026`
- `group_by = outlet`

Top 5 outlet rows:

```text
1730|Toko tosca|180500000.0000
1724|Toko abu abu|99868000.0000
1729|Toko Bersih|72050000.0000
1727|Toko hijau gelap|54000000.0000
1728|Toko orange|33000000.0000
```

Interpretation:

- Query with `cust_id + month + year` returns ordered outlet summary as expected.
- Ordering `net_sales DESC` confirmed from live DB result ordering.

## Group fallback validation

- `group_by` kosong tetap fallback ke branch `product`.
- `group_by` unknown non-empty juga tetap fallback ke branch `product`.
- Coverage ada di service test `TestSecondarySalesReportGroupSalesUsesFallbackYearForAllBranches` termasuk case `unknown group falls back to product`.

## Before vs after query intent

Before:

- `report.fact_orders.cust_id = ? AND dt.month = ?`
- `report.fact_returns.cust_id = ? AND dt.month = ?`
- group branches only `cust_id + month`

After:

- `report.fact_orders.cust_id = ? AND dt.month = ? AND dt."year" = ?`
- `report.fact_returns.cust_id = ? AND dt.month = ? AND dt."year" = ?`
- all group branches use `cust_id + month + dt."year"`

## Files validated

- `sales/entity/report.go`
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`
- `sales/controller/so_controller_test.go`

## Remaining note

- Scope check now requires both `is_del = false` and `is_active = true` for child `cust_id` override.
- If business later wants a different BU rule than active child distributor under `parent_cust_id`, scope resolver must be revised explicitly.
