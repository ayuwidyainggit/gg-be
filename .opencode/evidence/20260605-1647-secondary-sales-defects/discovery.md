Task: SX-2143 + SX-2165 Secondary Sales BE defect plan

Source strategy:
- Local repo evidence used: `AGENTS.md`, `.opencode/docs/ARCHITECTURE.md`, `.opencode/docs/QUALITY.md`, `sales/controller/report_controller.go`, `sales/entity/report.go`, `sales/model/report.go`, `sales/service/report_service.go`, `sales/repository/report_repository.go`, `sales/service/report_service_test.go`.
- Official docs skipped: task depends on repo SQL/contracts, not unknown library behavior.
- GitHub/web/Jira fetch skipped: Jira details supplied by user; repo-local code enough for fix plan.
- Browser evidence skipped: backend API/export defect, not UI visual parity.

Files inspected:
- `.opencode/docs/ARCHITECTURE.md`: strict Controller → Service → Repository → DB, tenant `cust_id` rules.
- `.opencode/docs/QUALITY.md`: validate in target service dir; use `rtk go test ./...` under `sales`.
- `sales/controller/report_controller.go`: routes exist for `/v1/reports/secondary-sales/sum-date`, `/group`, `/trend-sales`, export; current code already routes auth/parent cust to service and parses query request.
- `sales/entity/report.go`: dashboard payloads already have `Month`, optional `Year`, and `CustID`; response lacks `total_ppn` and `net_sales_exc_ppn`.
- `sales/service/report_service.go`: dashboard cust resolution exists with parent-scope validation; year fallback exists; return rate already avoids divide by zero when order qty is zero.
- `sales/repository/report_repository.go`: return export branch still has `0 AS special_discount`; return extract still has `0 AS special_discount`; summary query still order-only for main totals and separate return qty/net only; group queries still order-only.
- `sales/service/report_service_test.go`: tests already cover dashboard `cust_id` scope, year fallback, unauthorized cust, safe return rate.

Project patterns found:
- Service owns customer-scope authorization through `resolveSecondaryDashboardCustID` and `ExistsCustomerInParentScope`.
- Dashboard reports use `report.fact_orders`, `report.fact_returns`, and `report.dim_dates`.
- Export/list union uses SQL builder `buildSecondarySalesUnionQuery`.
- Extract uses `GetReportSecondarySalesReportOrder` and `GetReportSecondarySalesReportReturn`, then writes `model.FactOrder` / `model.FactReturn`.
- Tests exist in `sales/service/report_service_test.go` and `sales/repository/report_repository_test.go`; repository SQL tests likely use SQL mocking.

Reuse candidates:
- Reuse existing `resolveSecondaryDashboardCustID` instead of new controller-level free `cust_id` trust.
- Reuse `resolveSecondaryDashboardYear` and optional `Year *int` pattern.
- Reuse fact table summary path, but combine order/return with CTE/raw SQL.
- Reuse `buildSecondarySalesUnionQuery` for export/list promo fix.

Constraints:
- Do not hardcode `cust_id`, month, year, user, token.
- Keep tenant scope checks.
- Preserve order transaction behavior.
- Service validation commands must run from `sales/` and stay `rtk`-prefixed per repo policy.
- No source edits made by planner.

Risks:
- VAT sign convention ambiguous: Jira reference uses `os.vat + rs.vat` for exposed `vat_value` but `os.vat - rs.vat` inside `net_sales_inc_ppn`. Plan marks this as open decision or implement repo-consistent `order - return` with PM/QA validation.
- `SpecialDiscount` export acceptance says `disc_value + promo_value`, but current export has separate `SpecialDiscount` and `Discount` columns. To avoid double-count in existing report layout, plan needs explicit choice: safest repo-compatible fix is `special_discount = promo_value`, `discount = disc_value`; Jira wording may require total in one visual column.
- Dashboard fact tables depend on extract freshness; live order/return tables may have data while facts empty.
