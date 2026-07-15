# Validation Evidence — SX-2174 + SX-2182 Secondary Sales BE

Task id: `20260611-1054-sx-2174-2182-secondary-sales`
Tanggal: `2026-06-11`

## Source changes

Changed source/test files:

- `sales/model/report.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`

No config/env/migration files were changed.

## Implementation summary

SX-2174:

- `model.SumReportByMonthModel` now has `ReturnRate` mapped from `return_rate`.
- `SecondarySalesReportSumReportByMonth` service now maps repository-provided `ReturnRate`; it no longer computes `qty_return / qty * 100`.
- `SecondarySalesReportSumReportByMonth` repository query now computes target fields from transactional tables:
  - orders: `sls."order"` + `sls.order_detail`
  - returns: `sls.return_det` + `sls."return"` + linked `sls."order"`
- Quantity formula uses converted smallest unit:
  - `(qty3 * conv_unit2 * conv_unit3) + (qty2 * conv_unit2) + qty1`
- `net_sales_return` uses return include PPN (`rd.total` aggregate / `rs.net_sales_inc_ppn`).
- `return_rate` is value-based and rounded in SQL with PostgreSQL-compatible numeric cast:
  - `COALESCE(ROUND(((rs.net_sales_inc_ppn / NULLIF(os.net_sales_inc_ppn, 0)) * 100)::numeric, 2), 0)`
- Return branch date filter follows SQL acuan intent through linked order invoice date: `o.invoice_date >= dateFrom AND o.invoice_date < dateTo`.

SX-2182:

- Existing string/array `cust_id` compatibility remained intact.
- Existing multi-cust authorization and row-level metadata behavior remained intact.
- No export code rewrite was needed after regression validation passed.

## Test validation

From `sales/`:

- `rtk go test ./repository -run 'TestSecondarySalesReportSumReportByMonth'` — PASS, 1 test package.
- `rtk go test ./service -run 'TestSecondarySalesReportSumReportByMonth'` — PASS, 5 service tests.
- `rtk go test ./service -run 'TestPublishSecondarySalesReport|TestSubscribeSecondarySalesReport|TestResolveSecondaryDashboardCustIDs'` — PASS, 10 service tests.
- `rtk go test ./repository -run 'TestBuildSecondarySalesUnionQuery|TestSecondarySalesReportGroup'` — PASS, 15 repository tests.
- Combined targeted final:
  - service targeted report/export tests — PASS, 15 service tests.
  - repository targeted report/export/group tests — PASS, 16 repository tests.
- `rtk go test ./...` — PASS, `245 passed in 22 packages`.

From repo root:

- `rtk docker compose -f docker-compose.yml ps` — PASS; `scylla-sales`, `scylla-master`, `rabbitmq`, `redis`, and other compose services were up.
- `rtk docker compose -f docker-compose.yml restart sales` — PASS; sales service restarted to load updated code.

## Local DB validation

Database: local `ggn_scyllax` via Postgres `localhost:5432`, user `postgres`.

Manual SQL for `cust_id=C260020001`, period `2026-06-01` to before `2026-07-01` returned:

| field | value |
| --- | ---: |
| `qty` | `235` |
| `qty_return` | `12` |
| `return_rate` | `0.01` |
| `net_sales_return` | `693693.0000` |

## Local API validation

Service: `http://localhost:9004`.

Token handling:

- Bearer token was used only in ephemeral shell commands.
- Token was not written to repository files, evidence files, or final notes.

`GET /v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001` returned HTTP `200` with matching fields:

```json
{"qty":235,"qty_return":12,"return_rate":0.01,"net_sales_return":693693}
```

Export validation:

- `POST /v1/reports/secondary-sales` with legacy `cust_id` string returned HTTP `200`.
  - report id: `6a2a39a48775fad2c3008742`
  - initial `file_status`: `2`
  - later DB status: `1`
  - `report.list.cust_id`: `C26002`
  - `file_url`: non-empty
- `POST /v1/reports/secondary-sales` with array `cust_id` returned HTTP `200`.
  - report id: `6a2a39a48775fad2c3008743`
  - initial `file_status`: `2`
  - later DB status: `1`
  - `report.list.cust_id`: `C26002`
  - `file_url`: non-empty
- Unauthorized export request with foreign `cust_id` returned HTTP `403` and message `cust_id is outside authorized scope`.

## Quality gate

Final `@quality-gate`: PASS.

Quality gate confirmed:

- Diff boundary conformance.
- SX-2174 transactional source correctness.
- converted quantity formula.
- include-PPN return value.
- value-based rounded return rate.
- SX-2182 backward compatibility and authorization.
- SQL binding/tenant safety.

## Commit status

No local commit was created because `/Users/ujang/Projects/Geekgarden/scylla-be` is not a git repository in this environment (`git rev-parse --show-toplevel` failed).

## Remaining risks / follow-ups

- No required remediation remains.
- Optional follow-up: include this SQL-vs-API comparison in PR/Jira notes.
