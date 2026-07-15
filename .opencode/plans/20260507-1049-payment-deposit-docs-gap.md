# Plan — Payment Deposit Report Docs Gap Alignment

Task ID: `20260507-1049-payment-deposit-docs-gap`
Tanggal: 2026-05-07 10:49 Asia/Jakarta
Primary source of truth: `.opencode/plans/20260507-1049-payment-deposit-docs-gap.md`

## Goal

Pastikan endpoint Payment Deposit Report dan Download sesuai `docs/Report - Payment Deposit_BE.md` point 2 dan point 3 sampai ke query pengambilan data, sekaligus hapus legacy `ReportPaymentDepositController` dan seluruh kode legacy terkait agar request dengan `emp_id` dan `deposit_type=AR,AP` tidak lagi masuk validasi lama `SalesmanID is a required field`.

## Non-goals

- Tidak mengubah kontrak endpoint lain seperti Deposit Number List (`/finance/v1/deposits`) kecuali reuse diperlukan untuk test.
- Tidak mengubah schema database atau migration kecuali ditemukan field yang benar-benar belum ada saat implementasi.
- Tidak mengubah gateway/reverse proxy production; finance service cukup expose route internal `/v1/...` sesuai keputusan user.
- Tidak menyalin bearer token dari cURL user ke test, log, commit, atau artifact.

## Scope

### In scope

- Route active Payment Deposit Report:
  - internal finance service: `/v1/reports/payment-deposit`
  - internal finance service: `/v1/reports/payment-deposit/download`
  - public URL tetap dianggap `/finance/v1/...` via gateway/service prefix.
- Request parsing dan validation untuk docs point 2:
  - `q`, `page`, `limit`, `sort`, `deposit_type`, `start_date`, `end_date`, `emp_id`, `deposit_no`.
- Compatibility alias:
  - `salesman_id` bersifat non-mandatory. Jika dikirim, boleh dipakai sebagai alias opsional untuk `emp_id`; jika tidak dikirim, request tetap valid selama parameter docs utama terpenuhi.
- Query alignment untuk:
  - AR list.
  - AP list.
  - AR+AP union list.
  - AR download invoice rows.
  - AR download expense rows.
  - AP download rows.
- Removal legacy code path:
  - controller, entity, service, repository, model, tests, and `main.go` wiring for `ReportPaymentDeposit*` legacy implementation.

### Out of scope unless validation fails

- Performance tuning/index changes.
- Asynchronous background processing redesign for download.
- Excel visual formatting beyond docs-required columns/values.

## Requirements

1. `GET /v1/reports/payment-deposit` in finance service must be handled only by `PaymentDepositReportController.List`.
2. `GET /v1/reports/payment-deposit/download` in finance service must be handled only by `PaymentDepositReportController.Download`.
3. Public docs URL `/finance/v1/reports/payment-deposit` is satisfied by gateway prefix; service route itself must not include `/finance` per user instruction.
4. Remove all legacy `ReportPaymentDepositController` code instead of leaving no-op route methods.
5. Remove `ReportPaymentDepositQueryFilter` legacy DTO with `SalesmanID validate:"required"` so no path can return the old validation error.
6. Request validation must accept cURL shape:
   - `page=1`
   - `limit=10`
   - `sort=created_date:desc`
   - `start_date=1775001600`
   - `end_date=1780271999`
   - `emp_id=421,415,381`
   - `deposit_type=AR,AP`
7. `deposit_type` must support comma-separated and repeated query values, normalize to `AR`/`AP`, reject unknown values.
8. `emp_id` must support comma-separated and repeated query values, validate integer list, and filter AR query by `acf.deposit.emp_id IN ?` when provided.
9. AP query must not require or apply `emp_id`, because AP docs SQL has no collector filter.
10. `deposit_no` must support comma-separated and repeated query values:
    - AR maps to `acf.deposit.deposit_no`.
    - AP maps to `acf.account_payable_payment.account_payable_payment_no`.
11. `start_date`/`end_date` must be parsed from epoch to date range used by SQL:
    - AR list uses `d.deposit_date BETWEEN ? AND ?`.
    - AP list uses `app.account_payable_payment_date BETWEEN ? AND ?`.
12. `sort=created_date:desc` must be accepted and normalized/mapped to `deposit_date desc` because docs default says `created_date:desc` while report rows expose `deposit_date`.
13. Response list should match docs field names:
    - `deposit_date`, `deposit_no`, `deposit_type`, `collector_id`, `collector_code`, `collector_name`, `cash_amount`, `cheque_amount`, `transfer_amount`, `return_amount`, `credit_debit_amount`, `expense_amount`, `total_payment`, `pagination`, `request_id`.
14. Decide and implement `deposit_type` display values consistently. Recommended: return docs display strings `Account Receivable`/`Account Payable`; if FE already expects `AR`/`AP`, keep `AR`/`AP` but document intentional deviation and adjust docs/FE contract.
15. `salesman_id` must be non-mandatory for list and download request validation. It may remain as an optional backward-compatible alias only.
16. Download response must match docs point 3 behavior: processing message with `data: null` when queued/processing, and ready metadata if generated synchronously.

## Acceptance Criteria

- Request matching user cURL no longer returns `SalesmanID is a required field`.
- Finance service route table has no legacy `/v1/reports/payment-deposit` handler from `ReportPaymentDepositController`.
- Codebase has no `ReportPaymentDepositController`, `ReportPaymentDepositQueryFilter`, legacy `ReportPaymentDepositService`, legacy `ReportPaymentDepositRepository`, or their model/test files unless a compile-required shared type is deliberately renamed/reused.
- `main.go` no longer constructs or routes legacy `reportPaymentDepositController`.
- Unit tests prove:
  - `emp_id=421,415,381` is accepted and normalized.
  - `deposit_type=AR,AP` is accepted.
  - `sort=created_date:desc` is accepted and mapped safely.
  - missing `salesman_id` does not fail list or download validation.
  - missing `emp_id` does not fail validation; it simply means AR branch is not filtered by collector unless product later decides otherwise.
  - unknown `deposit_type` fails.
  - invalid `emp_id` fails.
  - `end_date < start_date` fails.
- Repository SQL-shape tests prove:
  - AR query includes `FROM acf.deposit d`, `d.emp_id IN ?`, `d.cust_id = ?`, `d.deleted_at IS NULL`, `d.deposit_date BETWEEN ? AND ?`.
  - AP query includes `FROM acf.account_payable_payment app`, `app.deleted_by IS NULL`, `app.account_payable_payment_date BETWEEN ? AND ?`, and does not include `emp_id` filter.
  - AR+AP combines with `UNION ALL`.
  - `deposit_no` filters the correct columns in each branch.
- `go test ./controller ./repository ./service` passes in `finance/`.
- `go test ./...` passes in `finance/` or every failure is unrelated and documented with exact package/error.

## Existing Patterns/Reuse

- Reuse active `PaymentDepositReportController`, `PaymentDepositReportService`, `PaymentDepositReportRepository`, and `PaymentDepositReportQueryFilter` as the single implementation.
- Reuse existing helper patterns:
  - `normalizeCSVValues`
  - `normalizeAndValidateDepositTypes`
  - `validateAndNormalizeSort`
  - `normalizeDateInput`
  - repository raw SQL builders.
- Extend current tests in `finance/controller/payment_deposit_report_controller_test.go` instead of creating a parallel validation test style.
- No matching KiloCode utility was needed; project-local active implementation already exists and should be consolidated.

## Constraints

- Follow Controller → Service → Repository → DB; controller must not call repository.
- Preserve tenant filter `cust_id` in every branch.
- Avoid string matching for runtime errors where Go sentinel errors exist; not central here.
- Do not add `fmt.Println` or `log.Println` debug prints.
- Do not expose secrets from user cURL.
- Service modules are independent; run commands from `finance/`.
- Project `AGENTS.md` says use `rtk` for commands, but global OpenCode says not to prefix commands. For implementation validation, prefer commands approved by user/session policy; document exact command actually run.

## Risks

- Docs point 2 table contains AP/AR mapping contradictions; SQL cases and response attribute tables appear more reliable than the param table notes.
- Changing `deposit_type` response from `AR`/`AP` to full labels may break existing FE if it consumes codes. This needs a small decision before implementation.
- Removing legacy files can reveal compile dependencies in `main.go` or tests; compile/test will guide cleanup.
- Download docs point 3 still references `salesman_id`; active point 2 uses `emp_id`. Keeping `salesman_id` alias reduces FE/backward-compat risk.
- Live DB query parity may depend on remote DB data and permissions; use SQL-shape unit tests plus optional integration cURL.

## Decisions/Assumptions

- User decision: remove `/finance` from finance service route; use `/v1/reports/payment-deposit` internally.
- Assumption: public `https://best.scyllax.online/finance/v1/...` is routed by gateway to finance service `/v1/...`.
- Assumption: docs SQL examples are the source of truth when param table notes contradict AP/AR table mapping.
- User decision: `salesman_id` is non-mandatory in this request path. Keep it only as an optional alias for `emp_id` during transition.
- Open question before implementation: Should `deposit_type` in JSON response be full docs labels (`Account Receivable`, `Account Payable`) or keep current codes (`AR`, `AP`) for FE compatibility?
- Open question before implementation: Should `start_date`/`end_date` reject `YYYY-MM-DD` because docs says epoch, or keep existing date-string compatibility?

## TDD/Test Plan

### TDD required

Ya. Ini bug API behavior + validation + query mapping yang bisa regresi. Gunakan Red → Green → Refactor.

### Existing test patterns

- `finance/controller/payment_deposit_report_controller_test.go` already tests validation helpers directly.
- `finance/repository/report_payment_deposit_repository_test.go` exists for legacy repository; remove or replace with tests for active repository.

### First failing/regression test

1. Add controller validation test that builds `PaymentDepositReportQueryFilter` equivalent to user cURL:
   - `DepositType: []string{"AR,AP"}`
   - `EmpID: []string{"421,415,381"}`
   - `Sort: "created_date:desc"`
   - epoch dates.
2. Assert validation succeeds, `EmpID` length is 3, `DepositType` has `AP` and `AR`, and `Sort` becomes `deposit_date:desc` or repository maps `created_date` safely.
3. Add tests that missing `SalesmanID` is valid for list and download.
4. Add tests that missing `EmpID` is valid and results in no AR collector predicate being appended.

### Repository regression tests

- Test active `buildQuery` or `buildCountAndDataQueries` SQL string for AR only, AP only, AR+AP.
- Assertions should check important substrings and args count/order, not full whitespace-sensitive SQL.
- Recommended cases:
  - AR with `emp_id` and `deposit_no`.
  - AP with same filter should not include `emp_id` predicate.
  - Combined AR+AP uses `UNION ALL`.
  - `q` applies to AR deposit no/employee and AP payment no.

### Green step

- Delete legacy code and fix compile errors.
- Adjust active validation and repository mapping until regression tests pass.
- Align response `deposit_type` values after answering/openly deciding display-label question.

### Refactor step

- Remove duplicated compatibility logic where possible.
- Keep helper names clear: e.g. `EmpID` canonical, `SalesmanID` legacy alias.
- Ensure deleted legacy tests are replaced by active implementation tests.

### Edge cases

- `deposit_type=AR`, `emp_id` empty: allow and return all AR collectors for the date range; do not require `salesman_id` or `emp_id`.
- `deposit_type=AP`, `emp_id` empty: allow because AP has no collector.
- `deposit_type=AR,AP`, `emp_id` present: filter only AR branch.
- `deposit_type=AR,AP`, `emp_id` empty: allow; AR branch has no collector filter, AP branch unchanged.
- Large `limit`: current cap 9999; keep unless docs says otherwise.
- Duplicate CSV values: optional dedupe; not required unless query perf concern.

### Commands

Run from `finance/` after implementation:

```bash
go test ./controller ./repository ./service
go test ./...
```

If project/session requires RTK:

```bash
rtk go test ./controller ./repository ./service
rtk go test ./...
```

## Implementation Steps

1. **Red: add/adjust tests first**
   - Extend `finance/controller/payment_deposit_report_controller_test.go` with user-cURL validation case.
   - Add active repository SQL-shape tests for `payment_deposit_report_repository.go`.
   - Add compile guard test or route test ensuring active controller owns `/v1/reports/payment-deposit` and legacy handler is gone.

2. **Remove legacy code path**
   - Delete `finance/controller/report_payment_deposit_controller.go`.
   - Delete `finance/entity/report_payment_deposit.go`.
   - Delete `finance/service/report_payment_deposit_service.go`.
   - Delete `finance/repository/report_payment_deposit_repository.go`.
   - Delete `finance/repository/report_payment_deposit_repository_test.go` or migrate useful assertions to active repository tests.
   - Delete `finance/model/report_payment_deposit.go` if only used by legacy service/repository; otherwise migrate required struct usage to active model file.
   - Update `finance/main.go` to remove legacy repository/service/controller construction and route call.

3. **Route finalization**
   - Keep only active route group:
     ```go
     app.Group("/v1/reports/payment-deposit", middleware.JWTProtected())
     ```
   - Do not register `/finance/v1/...` inside finance service.

4. **Request contract alignment**
   - Make `EmpID` canonical for point 2.
   - Keep `SalesmanID` optional alias by copying to `EmpID` only when `EmpID` is empty.
   - Enforce `deposit_type`, `start_date`, `end_date` required.
   - Enforce integer list validation for `emp_id` and optional alias when either is present.
   - Do not require `salesman_id` or `emp_id` for list/download validation.
   - Keep `AP` allowed without `emp_id`.

5. **Sort/date alignment**
   - Accept docs default `created_date:desc` and map to `deposit_date:desc`/`t.deposit_date DESC`.
   - Keep whitelist sort fields only.
   - Parse epoch to `YYYY-MM-DD` for SQL args; optionally keep `YYYY-MM-DD` compatibility but test epoch as docs-required input.

6. **List SQL alignment**
   - AR branch:
     - `acf.deposit d`.
     - pre-aggregate `acf.deposit_payment` by `deposit_no, cust_id`.
     - pre-aggregate `acf.deposit_expense` by `deposit_no, cust_id`.
     - join `mst.m_employee` for collector code/name.
     - filter `d.cust_id`, `d.deleted_at IS NULL`, `d.deposit_date BETWEEN ? AND ?`, and append `d.emp_id IN ?` only when `emp_id` or optional `salesman_id` alias is provided.
     - compute total as cash + cheque + transfer + return + credit/debit - expense.
   - AP branch:
     - `acf.account_payable_payment app`.
     - aggregate `acf.account_payable_payment_options`.
     - filter `app.cust_id`, `app.deleted_by IS NULL`, `app.account_payable_payment_date BETWEEN ? AND ?`.
     - no `emp_id` predicate.
     - collector fields null.
   - Combined branch:
     - use `UNION ALL`, same column aliases and compatible data types.

7. **Download SQL alignment**
   - Compare `buildDownloadARQuery` to docs Query Account Receivable:
     - invoice rows from deposit payment/order/outlet/detail.
     - expense rows from deposit expense/expense/expense type.
     - filter `d.emp_id IN ?` when AR included.
   - Compare `buildDownloadAPQuery` to docs Query Account Payable:
     - payment options joined to payment detail/account payable/supplier.
     - filter date/cust/deleted.
     - no collector filter.
   - Add SQL-shape tests for download if not too brittle.

8. **Response alignment**
   - Confirm whether `deposit_type` should be `Account Receivable`/`Account Payable` or `AR`/`AP`.
   - Ensure AP collector fields are `null` for list response if following response table.
   - Ensure response wrapper places `items`, `summary` if retained, and `pagination` under `data` via `responsebuild`.

9. **Validation and cleanup**
   - Run focused tests, then full finance tests.
   - Run a local/service cURL smoke test without secret token if auth can be bypassed only in test; otherwise document manual production/gateway check.
   - Run `git diff` to verify only relevant files changed and no secrets added.

## Expected Files to Change

### Delete legacy files

- `finance/controller/report_payment_deposit_controller.go`
- `finance/entity/report_payment_deposit.go`
- `finance/service/report_payment_deposit_service.go`
- `finance/repository/report_payment_deposit_repository.go`
- `finance/repository/report_payment_deposit_repository_test.go`
- `finance/model/report_payment_deposit.go` if unused after cleanup.

### Update active implementation

- `finance/main.go`
- `finance/controller/payment_deposit_report_controller.go`
- `finance/entity/payment_deposit_report.go` if response/display values or comments need alignment.
- `finance/service/payment_deposit_report_service.go` if response mapping/download behavior needs adjustment.
- `finance/repository/payment_deposit_report_repository.go`
- `finance/controller/payment_deposit_report_controller_test.go`
- Add `finance/repository/payment_deposit_report_repository_test.go` for active SQL-shape tests.

## Agent/Tool Routing

- Implementation: `@fixer` / bounded code edits and tests.
- Architecture/query review if SQL mapping remains ambiguous: `@oracle`.
- Security/privacy review not required unless auth/tenant filtering changes beyond current `cust_id`/JWT usage.
- Release engineer not required unless deployment/gateway route config must be changed.
- No external docs/context7 needed; local docs and local code are sufficient.

## Validation Commands

From repo root before code work if following project instructions:

```bash
rtk docker compose -f docker-compose.yml ps
```

From `finance/` after implementation:

```bash
go test ./controller ./repository ./service
go test ./...
```

If RTK is required in the active session:

```bash
rtk go test ./controller ./repository ./service
rtk go test ./...
```

Optional smoke checks after deploy/gateway availability:

```bash
curl 'https://best.scyllax.online/finance/v1/reports/payment-deposit?page=1&limit=10&sort=created_date:desc&start_date=1775001600&end_date=1780271999&emp_id=421,415,381&deposit_type=AR,AP' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer <redacted>'
```

Expected smoke outcome: not `SalesmanID is a required field`; response is 200 with data/pagination or a domain/database error unrelated to request validation.

## Evidence Requirements

- Keep this plan as source of truth.
- Keep discovery evidence: `.opencode/evidence/20260507-1049-payment-deposit-docs-gap/discovery.md` because it lists docs/code gaps and contradiction notes useful for implementer.
- During implementation, add evidence notes under same evidence folder if useful:
  - test command output summary,
  - route registration proof,
  - SQL-shape test results,
  - optional cURL smoke result with token redacted.
- No GitHub, Brave, Context7, or browser evidence needed; key facts are local docs and local Go code.

## Done Criteria

- Legacy `ReportPaymentDeposit*` implementation is removed, not merely no-op.
- Active `PaymentDepositReport*` is the only Payment Deposit Report code path.
- User cURL contract with `emp_id` + `deposit_type=AR,AP` validates successfully.
- Requests without `salesman_id` validate successfully; requests without both `salesman_id` and `emp_id` also validate successfully and do not add collector filtering.
- List and download SQL builders are reviewed/tested against docs point 2 and 3.
- Tests pass or unrelated failures are documented.
- Final summary lists changed files, validation, and any remaining docs contradiction/open question.

## Final Planning Summary

- Artifacts created:
  - `.opencode/plans/20260507-1049-payment-deposit-docs-gap.md` — source of truth for implementation.
  - `.opencode/evidence/20260507-1049-payment-deposit-docs-gap/discovery.md` — kept because it records inspected files, docs contradictions, and query-level gaps.
- Draft artifacts: none created, so no draft cleanup needed.
- Key decisions:
  - Route inside finance service should be `/v1/reports/payment-deposit`, not `/finance/v1/...`.
  - Remove legacy code path completely.
  - Follow docs SQL cases for AR/AP table mapping where docs param table contradicts itself.
- Questions asked: none; user direction was clear enough to produce the plan.
- User follow-up incorporated:
  - `salesman_id` is non-mandatory.
  - `emp_id`/collector filtering is optional; absence means no collector filter.
- Assumptions still open for implementer/product confirmation:
  - Whether response `deposit_type` should be full labels or current `AR`/`AP` codes.
  - Whether to keep accepting `YYYY-MM-DD` dates in addition to docs epoch.
- Readiness: ready for implementation after resolving the two small contract questions above, or proceed with recommended assumptions if FE compatibility is known.
