# Discovery Evidence — Payment Deposit Docs Gap

Task ID: `20260507-1049-payment-deposit-docs-gap`
Tanggal: 2026-05-07 10:49 Asia/Jakarta

## Files inspected

- `docs/Report - Payment Deposit_BE.md`
  - Point 2 `Report Payment Deposit`, lines 65-197.
  - Point 3 `Report Payment Download`, lines 204-283.
- `finance/controller/payment_deposit_report_controller.go`
  - Current route and validation for `PaymentDepositReportController`.
- `finance/entity/payment_deposit_report.go`
  - Current DTO/response structs for active implementation.
- `finance/repository/payment_deposit_report_repository.go`
  - Current list/download query builders for AR, AP, and combined AR+AP.
- `finance/controller/report_payment_deposit_controller.go`
  - Legacy controller still present but route now no-op.
- `finance/entity/report_payment_deposit.go`
  - Legacy DTO still has `SalesmanID validate:"required"`, matching reported production error.
- `finance/service/report_payment_deposit_service.go`
  - Legacy service still present.
- `finance/repository/report_payment_deposit_repository.go`
  - Legacy repository still present.
- `finance/model/report_payment_deposit.go`
  - Legacy model still present.
- `finance/main.go`
  - Both legacy `reportPaymentDepositController` and active `paymentDepositReportController` are constructed; both route methods are called.

## Project patterns found

- Finance service is Fiber-based and registers routes via controller `Route(app *fiber.App)` methods.
- Active implementation follows Controller → Service → Repository.
- Query building for report is raw SQL in repository, not GORM chain.
- Tenant filtering uses `cust_id`; master lookup sometimes uses `parent_cust_id` for download joins.
- Existing active controller normalizes comma-separated query params for `deposit_type`, `emp_id`, `salesman_id`, and `deposit_no`.

## Reuse candidates

- Keep and extend active files:
  - `PaymentDepositReportController`
  - `PaymentDepositReportQueryFilter`
  - `PaymentDepositReportRepository`
  - `PaymentDepositReportService`
- Remove legacy `ReportPaymentDeposit*` files and wiring rather than maintaining two implementations.
- Reuse tests in `finance/controller/payment_deposit_report_controller_test.go` and add repository SQL-shape/unit tests.

## Docs-vs-code findings

### Endpoint routing

- Docs public URL: `/finance/v1/reports/payment-deposit`.
- User decision: finance service route should remove `/finance` and expose `/v1/reports/payment-deposit`; gateway/public layer supplies `/finance` prefix.
- Active controller now uses `/v1/reports/payment-deposit`.
- Legacy controller previously owned `/v1/reports/payment-deposit`; it must be removed completely to avoid `SalesmanID is a required field` from legacy validation.

### Request params

- Docs point 2 requires: `page`, `limit`, `sort`, `deposit_type`, `start_date`, `end_date`, `emp_id`; optional `q`, `deposit_no`.
- Active DTO supports all docs params plus backward-compatible `salesman_id` alias.
- Active controller currently accepts epoch or `YYYY-MM-DD`; docs says epoch. Plan should decide whether to keep date string compatibility while tests assert epoch.
- Active controller currently accepts `created_date` sort alias and normalizes to `deposit_date`.

### Important docs inconsistency

- Docs table says:
  - `start_date`: "if deposit_type = AP acf.deposit.deposit_date if deposit_type = AR acf.account_payable_payment.account_payable_payment_date"
  - `emp_id`: "if deposit_type = AP acf.deposit.emp_id if deposit_type = AR tidak ada"
  - `deposit_no`: "if deposit_type = AP acf.deposit.deposit_no if deposit_type = AR account_payable_payment_no"
- But docs SQL cases show the opposite real table mapping:
  - AR uses `acf.deposit`, `d.deposit_date`, `d.emp_id`, `d.deposit_no`.
  - AP uses `acf.account_payable_payment`, `app.account_payable_payment_date`, `app.account_payable_payment_no`, no collector.
- Plan should follow SQL cases/response attributes as source of truth and record docs table as typo unless product confirms otherwise.

### List query gaps

- Active AR list query broadly matches docs SQL: pre-aggregated `deposit_payment`, pre-aggregated `deposit_expense`, joins employee, filters `d.cust_id`, `d.deleted_at`, date range, optional `d.emp_id IN ?`, `d.deposit_no IN ?`.
- Active AP list query broadly matches docs SQL: `account_payable_payment`, pre-aggregated payment options, filters `app.cust_id`, `app.deleted_by`, date range, optional payment number.
- Active response uses `deposit_type` values `AR`/`AP`, while docs response examples say `Account Receivable`/`Account Payable`.
- Active AP collector fields are `NULL`, matching docs response table, although one example JSON contradicts that by showing collector fields for AP.
- Active summary exists in response; docs examples do not show summary. This may be accepted existing enhancement or needs confirmation.

### Download query gaps

- Active download implementation is in the newer `PaymentDepositReport*` path and includes AR invoice rows + AR expense rows and AP invoice rows, closer to docs point 3 than legacy.
- User clarified `salesman_id` must be non-mandatory. Treat `salesman_id` only as an optional backward-compatible alias for `emp_id`; absence of both `salesman_id` and `emp_id` should not fail validation and should mean no collector filter.
- Need tests ensuring cURL style with `emp_id=421,415,381&deposit_type=AR,AP` reaches active controller validation and repository filters AR by `d.emp_id IN ?` only, not AP.

## Commands/docs checked

- `rtk docker compose -f docker-compose.yml ps` was run earlier; finance/master/system were up.
- Test commands were attempted earlier but rejected by user permission.
- No external docs were needed because this is local code vs local docs alignment.

## Constraints

- Do not expose or copy bearer tokens/secrets from user cURL.
- Repo has plaintext credentials in some files; avoid touching unrelated files.
- Each service is an independent Go module; validation should run from `finance/`.
- Project instructions conflict on RTK prefix. Active project `AGENTS.md` says use `rtk`; global OpenCode says do not prefix in OpenCode. Use direct commands only if user explicitly approves; otherwise prefer project-consistent `rtk` for verification notes.

## Risks

- Removing legacy code may require deleting constructor/service/repository/model wiring from `main.go`; if any other code imports legacy symbols, compile will catch it.
- Public docs URL `/finance/v1/...` differs from internal service route `/v1/...`; validation must include either gateway/proxy check or clear curl against local service without prefix.
- Docs contain contradictory AP/AR param mapping; implementation should follow SQL examples and response attribute table unless user/product confirms table text.
- Raw SQL tests need SQL-shape assertions or sqlmock; avoid live DB dependency for deterministic regression coverage.
