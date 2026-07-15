# Execution — SX-2258 Secondary Sales Summary

Task id: `20260618-1002-sx-2258-secondary-sales-summary`
Module: `sales`
Branch: `bugfix/SX-2258-dev`

## Summary

- Updated `sales/repository/report_repository.go` to remove `r.data_status = 6` from `buildSecondarySalesReportSummarySQL` return-summary branch.
- Kept existing summary date semantics `o.invoice_date >= ? AND o.invoice_date < ?` unchanged.
- Added regression assertions so summary SQL rejects `r.data_status = 6` while preserving subtract arithmetic and invoice-date/product filters.
- Added focused service/controller mapping tests for `qty` and `total_discount_promo` passthrough.

## Probe outcome

Initial probe in fixer session reported missing local source tables. Follow-up validation after user confirmed local DB usage found the required schema/data available.

- Local DB: `localhost:5432/ggn_scyllax`.
- Required source tables exist: `sls.order`, `sls.order_detail`, `sls.return`, `sls.return_det`.
- `C260020001` has local source data for June 2026.
- Four-variant probe result for `C260020001`, June 2026:
  - `current_status_between`: order qty `277`, return qty `12`, net qty `265`, return value `693693`.
  - `current_status_semio`: order qty `277`, return qty `12`, net qty `265`, return value `693693`.
  - `no_status_between`: order qty `277`, return qty `12`, net qty `265`, return value `693693`.
  - `no_status_semio`: order qty `277`, return qty `12`, net qty `265`, return value `693693`.
- Direct Widya reference SQL on local data also returns net qty `265`, return qty `12`, order qty `277`.
- Result: local dataset does not reproduce QA/staging expected `qty=134`. It proves the implementation matches reference SQL behavior for the local `ggn_scyllax` dataset.
- Decision applied: remove `r.data_status = 6` and keep existing month-safe `>= AND <`; `BETWEEN` produces identical result on local data.

## Changed files

- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`
- `sales/service/report_service_test.go`
- `sales/controller/so_controller_test.go`

## Validation

### Targeted repository

- Command: `cd sales && rtk go test ./repository -run 'TestSecondarySalesReportSumReportByMonth|TestSX2258' -v`
- Result: PASS (`Go test: 3 passed in 1 packages`)

### Targeted service

- Command: `cd sales && rtk go test ./service -run 'TestSecondarySalesReportSumReportByMonth' -v`
- Result: PASS (`Go test: 7 passed in 1 packages`)

### Targeted controller

- Command: `cd sales && rtk go test ./controller -run 'TestSecondaryReportSalesSumMonth' -v`
- Result: PASS (`Go test: 5 passed in 1 packages`)

### Full module

- Command: `cd sales && rtk go test ./...`
- Result: PASS (`Go test: 267 passed in 22 packages`)

### Local DB SQL probe (post-fix)

- DB: `ggn_scyllax` at `localhost:5432`.
- 4-variant query output: all four variants return `order_qty=277`, `return_qty=12`, `net_qty=265`, `net_sales_return=693693`.
- Direct execution of Widya reference query (BETWEEN + no `r.data_status` filter) returns `qty=265`, `qty_return=12`, `total_qty_sold=277`, `total_qty_return=12`.
- Local result `qty=265` does not match QA/staging expected `qty=134`. Local dataset differs from QA fixture; subtract arithmetic and reference filter alignment verified.

### Local API login

- Endpoint: `POST http://localhost:9001/v1/users/login`.
- Principal user: `princessa@gmail.com` (cust `C26002`, parent `C26002`). Login OK, access token issued. Redacted, not stored.
- Distributor user: `adminbm@gmail.com` (cust `C260020001`, parent `C26002`). Login OK, access token issued. Redacted, not stored.
- Security: token values never written to files, logs, or evidence. Replaced with `<REDACTED_BEARER_TOKEN>` in this document.

### Local API endpoint validation

- `GET /v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001`.
- Distributor account: HTTP 200 with response `{"qty":265,"qty_return":12,"total_discount_promo":1403740,"return_rate":0.01,"net_sales_return":693693,...}`.
- Principal account (parent scope, no `cust_id`): HTTP 200 with all-zero summary, indicating service correctly honors the principal fallback. Redacted.
- Both responses confirm subtract arithmetic is active and `r.data_status = 6` removal did not break the return branch.
- Numbers in local response differ from QA expected values because local dataset differs. Closure for QA `qty=134`/`total_discount_promo=1238740` requires staging fixture (data sync) and is left as a follow-up.

### Extra check

- Command: `cd sales && rtk go vet ./...`
- Result: existing non-slice issues remain in `controller/report_controller.go` duplicate json tags at lines 54 and 65. Not introduced by this task.

## Risks / residual

- Acceptance numeric proof for QA/staging values (`qty=134`, `total_discount_promo=1238740`) is not reproducible on local `ggn_scyllax`: direct Widya reference query returns `qty=265`, `qty_return=12`, `total_qty_sold=277`, `total_qty_return=12`.
- Local API endpoint validation completed and matches local SQL (`qty=265`, `qty_return=12`, `return_rate=0.01`, `net_sales_return=693693`). This closes local DB/API risk but not staging fixture parity.
- Date semantics remains `>= AND <`; local probe shows `BETWEEN` and semi-open produce identical output for the local fixture.
- Staging/demo fixture validation is still required to close exact QA expected values. If staging direct Widya reference returns `qty=134`, deployed BE should match after this filter fix. If not, next suspect is data divergence, not code branch reviewed here.

## Security / secret handling

- No bearer token, JWT, or secret written to source or evidence.
- Manual curl placeholder, if later used, must remain `Bearer $USER_TOKEN` or `<REDACTED_BEARER_TOKEN>` only.
