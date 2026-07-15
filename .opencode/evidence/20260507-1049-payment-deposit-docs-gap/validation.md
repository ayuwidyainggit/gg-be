# Validation Evidence — Payment Deposit Docs Gap

Task ID: `20260507-1049-payment-deposit-docs-gap`
Tanggal: 2026-05-07 Asia/Jakarta

## Scope implemented

- Removed legacy `ReportPaymentDeposit*` code path from finance service.
- Kept only active `PaymentDepositReport*` code path for payment deposit report.
- Internal finance route is `/v1/reports/payment-deposit` and `/v1/reports/payment-deposit/download`.
- `salesman_id` is non-mandatory for list and download validation.
- `salesman_id` remains only an optional alias into `emp_id` when `emp_id` is empty.
- Missing both `salesman_id` and `emp_id` is valid and does not append collector filter.

## Route/legacy proof

- Search for legacy Go symbols returned no matches in `finance/**/*.go`:
  - `ReportPaymentDeposit`
  - `reportPaymentDeposit`
  - `NewReportPaymentDeposit`
  - `ReportPaymentDepositRepository`
  - `ReportPaymentDepositService`
  - `ReportPaymentDepositRow`
- Deleted legacy files:
  - `finance/controller/report_payment_deposit_controller.go`
  - `finance/entity/report_payment_deposit.go`
  - `finance/service/report_payment_deposit_service.go`
  - `finance/repository/report_payment_deposit_repository.go`
  - `finance/repository/report_payment_deposit_repository_test.go`
  - `finance/model/report_payment_deposit.go`
- `finance/main.go` no longer constructs/routes legacy repository/service/controller.

## Test evidence

Commands run from `finance/`:

```bash
rtk go test ./controller ./repository ./service
rtk go test ./...
```

Results:

```text
Go test: 66 passed in 3 packages
Go test: 69 passed in 20 packages
```

Coverage added/updated:

- cURL-shaped list filter with `emp_id=421,415,381`, `deposit_type=AR,AP`, and `sort=created_date:desc` validates successfully.
- `salesman_id` missing validates successfully for list and download.
- `emp_id` missing validates successfully and does not append AR collector filter.
- Active repository SQL-shape tests assert:
  - AR uses `acf.deposit d`.
  - AR uses optional `d.emp_id IN ?` only when collector filter exists.
  - AP uses `acf.account_payable_payment app`.
  - AP does not apply `emp_id` filter.
  - AR+AP uses `UNION ALL`.
  - `created_date` sort maps to `t.deposit_date`.

## SELECT-only DB validation

Only `SELECT` was used for DB validation.

Query shape validated against remote development DB with credentials provided by user, redacted here:

```sql
SELECT 'AR' AS deposit_type, COUNT(*) AS rows
FROM acf.deposit d
WHERE d.cust_id = 'C260020001'
  AND d.deleted_at IS NULL
  AND d.deposit_date BETWEEN to_timestamp(1775001600)::date AND to_timestamp(1780271999)::date
  AND d.emp_id IN (421,415,381)
UNION ALL
SELECT 'AP' AS deposit_type, COUNT(*) AS rows
FROM acf.account_payable_payment app
WHERE app.cust_id = 'C260020001'
  AND app.deleted_by IS NULL
  AND app.account_payable_payment_date BETWEEN to_timestamp(1775001600)::date AND to_timestamp(1780271999)::date;
```

Result:

```text
 deposit_type | rows
--------------+------
 AR           |    8
 AP           |    3
```

This proves the docs-aligned AR/AP source tables contain rows for the cURL date/customer/collector filter.

## cURL smoke evidence

Token was provided by user and is intentionally redacted here.

### Local finance service

Request:

```bash
curl 'http://localhost:9005/v1/reports/payment-deposit?page=1&limit=10&sort=created_date:desc&start_date=1775001600&end_date=1780271999&emp_id=421,415,381&deposit_type=AR,AP' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer <redacted>'
```

Result:

```text
HTTP_STATUS:200
message: Data berhasil ditampilkan
items returned: 10
pagination.total_data: 11
pagination.total_page: 2
old validation error: not present
```

### Production public URL

Request:

```bash
curl 'https://best.scyllax.online/finance/v1/reports/payment-deposit?page=1&limit=10&sort=created_date:desc&start_date=1775001600&end_date=1780271999&emp_id=421,415,381&deposit_type=AR,AP' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer <redacted>'
```

Result:

```text
HTTP_STATUS:400
message: Bad Request
error: SalesmanID is a required field
```

Interpretation:

- Local code now fixes the issue and returns 200.
- Production URL still runs old deployed code/gateway target and has not picked up the local changes yet.
- This is expected until the finance service is redeployed/reloaded with the updated implementation.

## Git/repo note

- The active working directory is not a Git repository (`git status` reports: `fatal: not a git repository`).
- No local commit could be created from this workspace.
- Because there is no `.git`, file diff/staging checks are not available in this workspace.

## Remaining risk

- Public production cURL will continue to fail with the old `SalesmanID` validation until deployment uses this updated local code.
- Response `deposit_type` currently remains `AR`/`AP` in list output. Docs examples show full labels; this was left unchanged to avoid FE compatibility risk unless product requires full labels.

## Quality gate result

Final read-only quality gate result: `FAIL` due production/public URL still returning the old validation error.

Local code/test evidence passed, but production acceptance is blocked until the updated finance service is deployed or the gateway target points to the updated runtime.

Required production follow-up:

1. Deploy/reload finance service with this updated code.
2. Re-run production cURL against `https://best.scyllax.online/finance/v1/reports/payment-deposit?...`.
3. Expected result: no `SalesmanID is a required field`; ideally HTTP 200 with data/pagination.
