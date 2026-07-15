# Discovery SX-2234 — Sales Trend Net Sales & Gross

Task id: `20260616-1740-sx-2234-sales-trend`
Waktu: `2026-06-16 17:40 Asia/Jakarta`

## Sumber dicek

- Repo docs:
  - `.opencode/docs/index.md`
  - `.opencode/docs/ARCHITECTURE.md`
  - `.opencode/docs/QUALITY.md`
  - `.opencode/docs/SERVICE_MATRIX.md`
- Export dokumen: `/Volumes/External/Downloads/BrowserOS/sx2234_docs_export.txt`
- Kode lokal:
  - `sales/controller/report_controller.go`
  - `sales/service/report_service.go`
  - `sales/repository/report_repository.go`
  - `sales/model/report.go`
  - `sales/entity/report.go`
  - `sales/repository/report_repository_test.go`
  - `sales/service/report_service_test.go`
  - `sales/controller/so_controller_test.go`

## Pola repo ditemukan

- Repo multi-module Go; target service endpoint `/sales/v1/...` ada di module `sales`.
- `sales` compose-managed Fiber service, default port `9004`.
- Layer wajib: Controller → Service → Repository → DB.
- Validasi module target: dari `sales/`, pakai `rtk go test ./...` atau targeted `rtk go test ./repository -run ...`.
- Query read-only report ada di repository; service handle auth/scope `cust_id`.
- Shell repo-local pakai prefix `rtk` sesuai `AGENTS.md` repo.

## File dan fungsi target

- Route: `sales/controller/report_controller.go:98`
  - `reportRouteV1.Get("/secondary-sales/trend-sales", controller.SecondaryReportSalesTrendSales)`
- Controller: `sales/controller/report_controller.go:570-634`
  - bind `year`, optional body/query `cust_id`, normalize `CustIDs`, call service.
- Query parser: `sales/controller/report_controller.go:770-782`
  - supports repeated `cust_id`, `cust_id[]`, and comma values via `entity.NormalizeStringList`.
- Payload: `sales/entity/report.go:234-238`
  - `Year int query:"year" validate:"required"`
  - `CustID string json:"cust_id,omitempty" validate:"omitempty,alphanum,max=20"`
  - `CustIDs []string json:"-"`
- Service: `sales/service/report_service.go:1462-1482`
  - resolves effective cust ids, calls repo, maps response.
- Auth/scope helper: `sales/service/report_service.go:1331-1360`
  - empty requested cust => auth cust fallback.
  - distributor cannot request sibling.
  - principal can request own or scoped children after `ExistsCustomerInParentScope`.
- Repository target: `sales/repository/report_repository.go:1387-1421`
  - current query reads `report.fact_orders` + `report.dim_dates`.
  - current output already 12 months via local `m.month` CTE.
- Model/response:
  - `sales/model/report.go:432-436` scan aliases `month`, `total_gross_sale`, `total_discount_promo`, `net_sales`.
  - `sales/entity/report.go:270-274` JSON shape unchanged.

## Existing tests found

- Controller tests: `sales/controller/so_controller_test.go`
  - forbidden unauthorized sibling.
  - body `cust_id` pass-through.
  - multi-query `cust_id=CHILD1,CHILD2&cust_id=CHILD3`.
  - missing year returns 400.
  - invalid `cust_id` returns 400.
- Service tests: `sales/service/report_service_test.go`
  - child cust allowed.
  - distributor sibling rejected.
  - fallback to auth cust.
- Repository tests: `sales/repository/report_repository_test.go`
  - dry-run query capture helpers already exist.
  - sum-date and group tests assert source-table SQL and date vars.

## Reuse candidates

- Reuse `SecondarySalesReportSumReportByMonth` source-table pattern for order/return summary, date range, and `custIDs IN ?` binding.
- Reuse `SalesmanActivityReportSumByMonth` gross/net formula using `sell_price_final1/2/3`, `promo_final1..5`, `vat_value_final`, return `sell_price1/2/3`, `promo_value`, `vat_value`.
- Reuse repository dry-run test helper `newReportRepoDryRunDB`, `latestRecordedQuery`.
- Reuse existing DTO/model; no response schema change needed.

## Constraint dan risiko

- Jira names only `total_gross_sale` and `net_sales`; changing `total_discount_promo` needs explicit PM/FE confirmation unless current query refactor cannot preserve it. Plan chooses minimal visible behavior: keep field present and compute from order + return only if needed as intermediate; document if value changes.
- Docs export line 21 lacks `cust_id` in sample URL, but user prompt and existing code support `cust_id`; preserve existing query/body behavior.
- Current repository query uses `report.fact_orders`; target must use `sls."order"`, `sls.order_detail`, `sls.return_det`, `sls."return"`.
- Source SQL in prompt filters return date by `r.return_date`; existing sum-date code uses `o.invoice_date` for return side. SX-2234 explicitly says return dates filter selected year, so trend query should use `r.return_date >= dateFrom AND r.return_date < dateToExclusive` unless business confirms invoice-date semantics.
- `pjp-sales` has duplicate Sales code but not root compose; target `/sales/v1` likely `sales`. Do not edit `pjp-sales` unless deployment owner says `pjp-sales` also serves same endpoint.
- No real Bearer token allowed in tests, logs, fixtures, or plan evidence.

## Source strategy

- Local repo evidence: used and sufficient for code path, tests, validation commands.
- Official docs/context7: skipped; no external framework/API behavior central to plan.
- GitHub/web search: skipped; business formula provided by Jira/local export and repo source.
- Browser evidence: skipped; backend API change only.

## Decision gate answers

User answered question gate after plan creation:

- `total_discount_promo`: use source-table order+return formula. Reason: `net_sales` needs new discount source, mixed `report.fact_orders` and `sls.*` response would be inconsistent, FE may use field for tooltip/label/chart breakdown.
- Return-side date filter: use `r.return_date`.
- Duplicate module: only `sales`; do not change `pjp-sales`.

## Open questions

- Tidak ada blocking question tersisa.
