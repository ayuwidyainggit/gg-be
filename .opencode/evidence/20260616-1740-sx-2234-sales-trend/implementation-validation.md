# Implementation Validation SX-2234

Task id: `20260616-1740-sx-2234-sales-trend`
Waktu: `2026-06-16 Asia/Jakarta`

## Changed files

- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`

## What changed

- `SecondarySalesReportTrendSales` no longer uses `report.fact_orders`.
- Added `buildSecondarySalesReportTrendSalesSQL()`.
- Trend query now uses:
  - `sls."order" o`
  - `sls.order_detail od`
  - `sls.return_det rd`
  - `sls."return" r`
- Year range now half-open UTC:
  - `dateFrom := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)`
  - `dateTo := dateFrom.AddDate(1, 0, 0)`
- Order filter:
  - `o.invoice_date >= ? AND o.invoice_date < ?`
- Return filter and grouping:
  - `r.return_date >= ? AND r.return_date < ?`
  - `EXTRACT(MONTH FROM r.return_date)::INTEGER`
- 12 months preserved with `WITH months AS (...)`.
- Aliases preserved:
  - `month`
  - `total_gross_sale`
  - `total_discount_promo`
  - `net_sales`

## Formula proof

Implemented formulas:

- `total_gross_sale = COALESCE(os.gross_sales, 0) - COALESCE(rs.gross_sales, 0)`
- `total_discount_promo = COALESCE(os.discount_promo, 0) + COALESCE(rs.discount_promo, 0)`
- `net_sales = ((COALESCE(os.gross_sales, 0) - COALESCE(os.discount_promo, 0)) - (COALESCE(rs.gross_sales, 0) - COALESCE(rs.discount_promo, 0))) + (COALESCE(os.ppn, 0) - COALESCE(rs.ppn, 0))`

Order source fields:

- gross: `qty1_final * sell_price_final1`, `qty2_final * sell_price_final2`, `qty3_final * sell_price_final3`
- discount/promo: `disc_value_final`, `promo_final1..5`
- ppn: `vat_value_final`

Return source fields:

- gross: `qty1 * sell_price1`, `qty2 * sell_price2`, `qty3 * sell_price3`
- discount/promo: `disc_value`, `promo_value`
- ppn: `vat_value`

## Tests run

From `/Users/ujang/Projects/Geekgarden/scylla-be/sales`:

```bash
rtk go test ./repository -run 'TestSecondarySalesReportTrendSalesSQLUsesSourceTablesAndNetSalesFormula|TestSecondarySalesReportSumReportByMonthSQLUsesSourceTablesAndDateRange'
```

Result:

```text
Go test: 2 passed in 1 packages
```

```bash
rtk go test ./service -run 'TestSecondarySalesReportTrendSales' && rtk go test ./controller -run 'TestSecondaryReportSalesTrendSales'
```

Result:

```text
Go test: 3 passed in 1 packages
Go test: 5 passed in 1 packages
```

```bash
rtk go test ./repository -run 'TestSecondarySalesReportTrendSalesSQLUsesSourceTablesAndNetSalesFormula|TestSecondarySalesReportSumReportByMonthSQLUsesSourceTablesAndDateRange' && rtk go test ./service -run 'TestSecondarySalesReportTrendSales' && rtk go test ./controller -run 'TestSecondaryReportSalesTrendSales' && rtk go test ./...
```

Result:

```text
Go test: 2 passed in 1 packages
Go test: 3 passed in 1 packages
Go test: 5 passed in 1 packages
Go test: 262 passed in 22 packages
```

## Runtime/API validation

Compose status command from repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

Result:

```text
[rtk] WARNING: untrusted project filters (.rtk/filters.toml)
[rtk] Filters NOT applied. Run `rtk trust` to review and enable.
time="2026-06-16T18:06:52+07:00" level=warning msg="/Users/ujang/Projects/Geekgarden/scylla-be/docker-compose.yml: the attribute `version` is obsolete, it will be ignored, please remove it to avoid potential confusion"
NAME      IMAGE     COMMAND   SERVICE   CREATED   STATUS    PORTS
```

Initial runtime API/direct SQL comparison was not run because compose had no active services and no sanitized auth/runtime context was available. This gap was later closed in `.opencode/evidence/20260616-1740-sx-2234-sales-trend/runtime-validation.md` after starting compose and using a synthetic local JWT kept only in shell memory.

## Diff boundary check

Allowed files changed:

- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`

No source edits in `pjp-sales`, other services, migrations, `.env`, or secrets.

## Secret check

No real Bearer token added. Test uses synthetic `CUST-1` only.

## Notes

- `total_discount_promo` visible value is now source-table order+return, per user decision gate. MR summary must mention this behavior change.
- Existing `report.fact_orders` references still exist in unrelated salesman activity functions below the changed trend code; trend query itself no longer uses it.
