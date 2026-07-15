# Discovery Evidence — SX-2045 Payment Recap

Task ID: `20260525-1300-sx-2045-payment-recap`
Tanggal: 2026-05-25 13:00 Asia/Jakarta

## Files inspected

- `AGENTS.md`
- `.opencode/docs/index.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `docs/Report - Payment Deposit_BE.md`
- `.opencode/plans/20260507-1049-payment-deposit-docs-gap.md`
- `.opencode/evidence/20260507-1049-payment-deposit-docs-gap/validation.md`
- `finance/controller/payment_deposit_report_controller.go`
- `finance/entity/payment_deposit_report.go`
- `finance/model/payment_deposit_report.go`
- `finance/repository/payment_deposit_report_repository.go`
- `finance/service/payment_deposit_report_service.go`
- `finance/controller/payment_deposit_report_controller_test.go`
- `finance/repository/payment_deposit_report_repository_test.go`
- `finance/service/payment_deposit_report_service_test.go`

## Commands checked

```bash
rtk docker compose -f docker-compose.yml ps
rtk go test ./controller ./repository ./service
```

Results:

```text
compose: no running services listed
Go test: 66 passed in 3 packages
```

## Attachments checked

Workbook now available and inspected:

- `DownloadDepositPayment-250526-003.xlsx`

Screenshot supplied inline by user and used as visual confirmation.

Extraction artifact:

- `.opencode/evidence/20260525-1300-sx-2045-payment-recap/attachment-layout.md`

Key layout result:

- Detail rows keep existing `A:Q` shape.
- Recap starts two blank rows after detail table.
- AR block uses `B` header and `A:B` labels/values.
- AP block uses `E` header and `D:E` labels/values.
- AR expense detail and recap are negative.

## Current active flow

- Route active: `GET /v1/reports/payment-deposit` → `PaymentDepositReportController.List` → `PaymentDepositReportService.ListReport` → `PaymentDepositReportRepository.FindAllPaymentDeposit` + `FindPaymentDepositSummary`.
- Route active: `GET /v1/reports/payment-deposit/download` → `PaymentDepositReportController.Download` → async `FindAllPaymentDepositDownload` → `generateExcel` → `report.list.file_base64`.
- `buildQuery` branches by `DepositType`:
  - `AR` uses `acf.deposit`, `acf.deposit_payment`, `acf.deposit_expense`.
  - `AP` uses `acf.account_payable_payment`, `acf.account_payable_payment_options`.
- `buildDownloadQuery` branches by `DepositType`:
  - AR invoice rows include `discount` and `payment_balance`.
  - AR expense rows append extra rows with `expense`.
  - AP rows include `discount` and `payment_balance`.

## Root cause found

1. View recap uses `FindPaymentDepositSummary`, which returns one flat `PaymentDepositReportSummaryRow`.
2. SQL in `FindPaymentDepositSummary` sums `cash_amount`, `cheque_amount`, `transfer_amount`, `return_amount`, `credit_debit_amount`, and `expense_amount` across whole `buildQuery` result without `GROUP BY deposit_type`.
3. Entity response `PaymentDepositReportSummary` has only one bucket; no `Account Receivable` / `Account Payable` recap dimension.
4. View summary has no `discount` and `payment_balance` fields, even SX-2045 requires `Discount` and `Total Payment Balance` in recap.
5. Download recap in `generateExcel` sums all `rows` into one set of labels and values. It ignores `row.DepositType` even though download rows already contain `Account Receivable` / `Account Payable`.
6. AR expense row in download currently selects `COALESCE(SUM(de.payment_amount), 0) AS expense`, so recap displays positive expense when summed. SX-2045 requires expense displayed negative and used as deduction.
7. Controller validation accepts only `AR` and `AP`; `deposit_type=All` from issue flow would be rejected unless FE omits deposit type on download. SX-2045 requires `All` to mean AR + AP with separated recap buckets.

## Reuse candidates

- Reuse active `PaymentDepositReport*` flow only; previous docs-gap work removed legacy flow.
- Reuse `buildDownloadQuery`/download row shape for recap metrics because it already contains `discount`, `payment_balance`, `expense`, and full deposit type labels.
- Add small reusable recap builder in service for download rows, plus repository grouped recap for view if performance needs SQL aggregation.
- Reuse existing test files and direct helper tests.

## Constraints

- Keep Controller → Service → Repository → DB.
- Keep tenant filters: `cust_id`, and `parent_cust_id` for parent master joins.
- No hardcoded sample nominal, date range, credential, or attachment values in source/test.
- Keep AP expense at zero.
- Keep AR expense negative in output recap.
- Avoid double count from detail joins; aggregate detail/payment rows at stable invoice/payment keys before summary.

## Risks

- Workbook layout is now confirmed, but source conflict remains: issue asks AP `Total Expense` separated while workbook leaves AP expense row blank.
- Changing `summary.total_expense` sign from positive to negative can affect FE if it already subtracts this value client-side.
- Adding `summary.by_deposit_type` is safer than replacing existing summary fields, but FE must know new field.
- Using download-detail query for view recap may be heavier; if performance matters, implement dedicated grouped SQL using same aggregation rules.
- `deposit_type=All` behavior needs controller validation update; otherwise issue reproduction cannot run.
