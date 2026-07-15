# Discovery — SX-2258 Secondary Sales Summary

Task id: `20260618-1002-sx-2258-secondary-sales-summary`
Mode: Maintenance Stability Mode

## Repo files inspected

- `AGENTS.md`
- `.opencode/docs/index.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`
- `sales/entity/report.go`
- `sales/model/report.go`

## Endpoint route

- Route: `GET /sales/v1/reports/secondary-sales/sum-date`
- Controller handler: `sales/controller/report_controller.go:521`, `SecondaryReportSalesSumMonth`
- Route registration: `sales/controller/report_controller.go:96`, `reportRouteV1.Get("/secondary-sales/sum-date", controller.SecondaryReportSalesSumMonth)`
- Service method: `sales/service/report_service.go:1374`, `SecondarySalesReportSumReportByMonth`
- Repository method: `sales/repository/report_repository.go:1277`, `SecondarySalesReportSumReportByMonth`
- Query builder: `sales/repository/report_repository.go:1174`, `buildSecondarySalesReportSummarySQL`

## Current implementation state

Current local code already contains SX-2258 fix in `sales/repository/report_repository.go`:

- `total_discount_promo` uses subtraction:
  - `COALESCE(os.discount_promo, 0) - COALESCE(rs.discount_promo, 0) AS total_discount_promo`
- `qty` uses subtraction:
  - `COALESCE(os.qty, 0) - COALESCE(rs.qty_return, 0) AS qty`
- Return branch joins order by invoice:
  - `JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id`
- Date filter uses `o.invoice_date` for order and return branches.
- Product filter uses `od.pro_id` for order and `rd.product_id` for return.
- `o.data_status IN (6,7)` exists in both branches.
- `r.data_status = 6` exists in return branch.

Git evidence:

- Recent commit in `sales`: `4ebacfe fix(secondary-sales): subtract returns from sold metrics`
- Commit summary says it changed `repository/report_repository.go`, `repository/report_repository_test.go`, `service/report_service.go`, `service/report_service_test.go`, and `entity/report.go`.
- Diff evidence shows old bug patterns:
  - `(os.discount_promo + rs.discount_promo) AS total_discount_promo`
  - `os.qty AS qty`
- Diff evidence shows replaced with subtraction and optional filters.

## Existing tests found

- `sales/repository/report_repository_test.go:386` checks source summary SQL, date filters, subtract arithmetic fragments, and rejects old plus/order-only patterns.
- `sales/repository/report_repository_test.go:459` checks optional filters for order/return branches and rejects `r.return_date` for summary filters.
- `sales/repository/report_repository_test.go:514` validates sample arithmetic in Go, but it is not an executed SQL/integration test.
- `sales/service/report_service_test.go` validates service mapping and filter pass-through.

## Reuse candidates

- Reuse `newReportRepoDryRunDB`, `latestRecordedQuery`, and `assertSecondarySalesSummaryDateVars` in `sales/repository/report_repository_test.go` for SQL-shape regression.
- Reuse `mockReportRepositoryForService` in `sales/service/report_service_test.go` for service response mapping tests.
- For deeper DB validation, add a guarded integration test or manual SQL verification using local/staging database and the exact request payload.

## Constraints

- Planner must not edit source/test files. Implementation goes to `@fixer` after plan.
- Repo requires `rtk` prefix for shell workflows in this project.
- Validate inside `sales/`, not repo root.
- Do not use Jira credentials or tokens.
- Do not copy secrets from tracked infra or comments.
- Preserve Controller → Service → Repository → DB layering.
- Preserve tenant/customer scope rules.

## Latest feedback (Widya / Yogie)

- Feedback terbaru fokus ke `Number of Product Sold` / response `qty` untuk `cust_id=C260020001`, `month=6`, `year=2026`.
- User pasted a staging `curl` with bearer token in chat. Token value was not copied into this artifact. Security action required: rotate/revoke that token before any validation.
- Reference SQL from Widya differs from current repo SQL in two important filter semantics:
  - Reference return branch does NOT include `r.data_status = 6`; current repo SQL includes it.
  - Reference date predicate uses `o.invoice_date BETWEEN :date_from AND :date_to`; current repo SQL uses `o.invoice_date >= ? AND o.invoice_date < ?`.
- Current repo already subtracts `qty`, so remaining defect likely comes from return branch filter mismatch, not final arithmetic.

## Risks

- Current code appears partially fixed; subtract arithmetic is correct, but filter mismatch can still produce wrong `qty`.
- Existing arithmetic test is not strong enough as DB regression proof.
- Staging expected values depend on data state; local DB may not contain `C260020001` June 2026 fixtures unless seeded/synced.
- Manual endpoint validation requires valid auth and local/staging environment access.
- Bearer token leaked in chat. Do not use, store, or log it. Rotate before validation.
