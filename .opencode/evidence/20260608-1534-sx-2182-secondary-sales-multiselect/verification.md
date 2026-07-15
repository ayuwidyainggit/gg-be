# Verification Evidence — SX-2182 Secondary Sales Multiselect BE

Task id: `20260608-1534-sx-2182-secondary-sales-multiselect`
Verified at: `2026-06-08T15:43+07:00` onward

## Implementation summary

Implemented by bounded fixer lane, then orchestrator patched one service mapping issue where `SecondarySalesReportGroupResp.Code` needed to be populated from `model.SecondarySalesReportGroup.Code`.

Changed files observed from implementation report:

- `master/controller/business_unit_controller.go`
- `master/controller/business_unit_controller_test.go`
- `sales/entity/report.go`
- `sales/model/report.go`
- `sales/controller/report_controller.go`
- `sales/controller/so_controller_test.go`
- `sales/service/report_service.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`

Artifact files created by planner:

- `.opencode/plans/20260608-1534-sx-2182-secondary-sales-multiselect.md`
- `.opencode/evidence/20260608-1534-sx-2182-secondary-sales-multiselect/discovery.md`
- `.opencode/evidence/20260608-1534-sx-2182-secondary-sales-multiselect/index.json`
- `.opencode/evidence/20260608-1534-sx-2182-secondary-sales-multiselect/verification.md`

## Key behavior verified by source inspection

- Business-unit `region_id`/`area_id` parsing uses strict parser in `master/controller/business_unit_controller.go`:
  - supports `region_id[]` and `region_id`
  - splits comma-separated values
  - trims whitespace
  - dedupes values
  - returns HTTP 400 for invalid numeric token
- Sales `cust_id` parsing uses `entity.StringListOrScalar` and `entity.NormalizeStringList` in `sales/entity/report.go` and `sales/controller/report_controller.go`:
  - accepts legacy string body
  - accepts array body
  - accepts comma/repeated query via controller helper path
  - dedupes values
  - rejects non-alphanumeric cust ids
- Sales service multi auth resolver in `sales/service/report_service.go`:
  - missing/empty `cust_id` falls back to auth cust
  - distributor user may only request auth cust
  - principal validates each non-auth child via `ExistsCustomerInParentScope`
  - unauthorized selection returns `ErrUnauthorizedCustID`
- Export query builder uses bound list filters in `sales/repository/report_repository.go`:
  - `od.cust_id IN ?`
  - `rd.cust_id IN ?`
- Dashboard queries use bound list filters:
  - `report.fact_orders.cust_id IN ?`
  - `report.fact_returns.cust_id IN ?`
  - group/trend raw queries bind `custIDs` as parameters.
- Group query now selects and groups `code`:
  - `SELECT id, code, name, COALESCE(SUM(net_sales), 0) AS net_sales`
  - `GROUP BY id, code, name`
- Service response mapping now includes `Code: r.Code`.
- Quality-gate remediation fixed export product metadata lookup for principal multi-cust export:
  - child product lookup now uses row-level `t.cust_id` instead of one bound auth/principal cust
  - parent fallback remains bound to `ParentCustID`
  - regression test asserts SQL contains `cp.pro_id = t.product_id AND cp.cust_id = t.cust_id` and does not bind a single auth cust for child lookup

## Validation commands run

Preflight from repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

Result: passed enough for preflight; `master`, `system`, and `redis` were up. Docker emitted only warnings about untrusted RTK filters and obsolete compose `version` attribute.

Targeted tests from `master/`:

```bash
rtk go test ./controller -run BusinessUnit
```

Result: passed — `Go test: 4 passed in 1 packages`.

```bash
rtk go test ./repository -run BusinessUnit
```

Result: no matching tests — `Go test: No tests found`. This is not a blocker because full module tests and controller/service coverage passed.

```bash
rtk go test ./service -run BusinessUnit
```

Result: passed — `Go test: 9 passed in 1 packages`.

Full master module:

```bash
rtk go test ./...
```

Result: passed — `Go test: 339 passed in 23 packages`.

Targeted tests from `sales/`:

```bash
rtk go test ./service -run 'SecondarySalesReport|PublishSecondarySalesReport|SubscribeSecondarySalesReport'
```

Result: passed — `Go test: 23 passed in 1 packages`.

```bash
rtk go test ./repository -run SecondarySales
```

Result: passed — `Go test: 18 passed in 1 packages`.

```bash
rtk go test ./controller -run SecondarySales
```

Result: passed — `Go test: 6 passed in 1 packages`.

Full sales module:

```bash
rtk go test ./...
```

Result: passed — `Go test: 227 passed in 22 packages`.

Quality-gate remediation rerun from `sales/`:

```bash
rtk go test ./repository -run SecondarySales
```

Result: passed — `Go test: 18 passed in 1 packages`.

```bash
rtk go test ./service -run 'SecondarySalesReport|PublishSecondarySalesReport|SubscribeSecondarySalesReport'
```

Result: passed — `Go test: 23 passed in 1 packages`.

```bash
rtk go test ./...
```

Result: passed — `Go test: 227 passed in 22 packages`.

## Diff boundary check

Within allowed plan boundary:

- `master/controller/business_unit_controller.go`
- `master/controller/business_unit_controller_test.go`
- `sales/entity/report.go`
- `sales/model/report.go`
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`
- `.opencode/plans/20260608-1534-sx-2182-secondary-sales-multiselect.md`
- `.opencode/evidence/20260608-1534-sx-2182-secondary-sales-multiselect/**`

Potential out-of-boundary but related test file:

- `sales/controller/so_controller_test.go`

Justification: this file contains the existing `mockReportServiceForController` used by report controller tests; the service interface signature changed for `SecondarySalesReportTrendSales`, so the mock needed updating. It also contains Secondary Sales controller tests despite the filename.

No `.env`, secrets, dumps, `pjp-sales`, or unrelated runtime infra changes were reported/touched.

## Remaining risks / skipped checks

- Git diff/status could not be inspected because `/Users/ujang/Projects/Geekgarden/scylla-be` is not a git repository in this environment (`fatal: not a git repository`). Diff boundary was checked by implementation report plus file inspection, not by `git diff`.
- Manual cURL/API smoke was not run because auth tokens/runtime request setup was not provided. Unit/controller/repository/service coverage passed.
- Group `code` columns are assumed to exist on report dim tables based on docs and tests; runtime schema verification was not performed through DB because no manual DB query was run.

## Plan compliance checkpoint

- Plan status was `PASS` / `ready-for-implementation` before execution.
- Non-negotiable tenant invariant preserved by service resolver and 403 mapping.
- Backward compatibility for legacy string `cust_id` covered by tests.
- Multi list binding implemented via `IN ?` / bound params.
- Empty/missing cust fallback preserved as auth cust.
- RMQ old/new payload compatibility covered by service tests.
- Group response `code` implemented and mapped.
- Quality-gate blocker on row-level export product metadata lookup remediated with repository regression coverage.
- Validation commands passed for target modules.
