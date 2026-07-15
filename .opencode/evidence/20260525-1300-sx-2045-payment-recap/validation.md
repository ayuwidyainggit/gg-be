# Validation Evidence — SX-2045 Payment Recap

Task ID: `20260525-1300-sx-2045-payment-recap`
Tanggal: 2026-05-25 Asia/Jakarta

## Scope implemented

- `deposit_type=All` / `ALL` diterima backend dan dinormalisasi jadi `AR` + `AP`.
- View response sekarang punya `summary_by_deposit_type` untuk recap terpisah `Account Receivable` dan `Account Payable`.
- Download Excel recap sekarang dipisah dua blok sesuai attachment:
  - AR header/value block di `B` dan `A:B`
  - AP header/value block di `E` dan `D:E`
- AR expense detail/export/recap dibuat negative.
- AP expense dijaga `0` di data recap.
- Detail table `A:Q` tetap sama kecuali sign expense AR.

## Files changed

- `finance/controller/payment_deposit_report_controller.go`
- `finance/controller/payment_deposit_report_controller_test.go`
- `finance/entity/payment_deposit_report.go`
- `finance/model/payment_deposit_report.go`
- `finance/repository/payment_deposit_report_repository.go`
- `finance/repository/payment_deposit_report_repository_test.go`
- `finance/service/payment_deposit_report_service.go`
- `finance/service/payment_deposit_report_service_test.go`

## Root cause confirmed

- Recap lama flat/global tanpa split AR/AP.
- `deposit_type=All` ditolak validasi.
- View recap lama tidak punya dimension deposit type dan tidak memuat `discount` + `payment_balance` per type.
- Download recap lama menjumlah semua row tanpa grouping `deposit_type`.
- AR expense source di export memakai sum positive, jadi recap ikut positive.

## Commands run

From `finance/`:

```bash
rtk go test ./controller ./repository ./service
rtk go test ./...
```

Results:

```text
Go test: 76 passed in 3 packages
Go test: 79 passed in 20 packages
```

## Test coverage added/confirmed

- Controller:
  - `deposit_type=All` normalize ke `AP`, `AR`.
- Repository:
  - AR/AP query shape tetap benar.
  - AR expense SQL memakai negative expression.
- Service:
  - `ListReport` mengisi `summary_by_deposit_type` untuk AR/AP.
  - AR `total_expense` negative.
  - AP `total_expense` dipaksa `0`.
  - `generateExcel` recap layout mengikuti attachment structure:
    - recap mulai 2 row setelah detail
    - `B<start>` = `Account Receivable`
    - `E<start>` = `Account Payable`
    - AR labels/values di `A:B`
    - AP labels/values di `D:E`
    - AR `Total Expense` negative
    - AP block berhenti di `Total Payment Balance`
  - header metadata download sekarang memakai metadata eksplisit dari filter request, bukan first/last row order.
  - resolver collector label diuji untuk no filter / multi filter / single collector konsisten / single collector campuran.

## Attachment alignment

Verified against extracted workbook evidence in `.opencode/evidence/20260525-1300-sx-2045-payment-recap/attachment-layout.md`:

- Detail AP row tetap tidak berubah bentuk.
- AR expense sample negative pada detail dan recap.
- AP recap mengikuti workbook: tidak render baris `Total Expense` di blok AP, walau API grouped recap tetap membawa `total_expense = 0`.

## Quality gate

Verdict: `PASS`

Accepted:
- `deposit_type=All` behavior OK.
- Grouped recap AR/AP OK untuk view + download.
- AR negative expense / AP zero OK.
- Tenant filters dan parameter binding OK.
- Header metadata risk fixed; date range sekarang dari filter metadata, collector label deterministic dan tested.
- DB live sample totals match workbook.

Remaining low risk:
- Async generation masih fire-and-forget goroutine; jika fetch/generate gagal, status report belum punya fail-state update/retry path. Ini operational risk lama, bukan blocker SX-2045.

## Database validation

Read-only DB validation executed with provided connection.

Connection used:

```text
host=103.28.219.73 port=25431 user=postgres dbname=scylla_citus_dev sslmode=disable
```

Sample workbook row set matches `cust_id = C260020001`.
Parent customer for lookup joins: `C26002`.

Validated recap totals for date range `2026-05-05` to `2026-05-08`:

```text
C260020001|Account Payable|20000000.0000|0|0|0|0|0.0000|0.0000|0
C260020001|Account Receivable|35080000.0000|0|3195000.0000|0|0|1009000.0000|1000.0000|-16000.0000
```

Field order above:

```text
cust_id|deposit_type|cash|cheque_giro|transfer|return_amount|credit_debit|discount|payment_balance|expense
```

These live totals match workbook attachment evidence:

- AR cash `35.080.000`
- AR transfer `3.195.000`
- AR discount `1.009.000`
- AR payment balance `1.000`
- AR expense `-16.000`
- AP cash `20.000.000`
- AP transfer `0`
- AP discount `0`
- AP payment balance `0`

Validated live detail rows also match workbook sample, including:

- `DP2605050001` `Jaka` `INV2604290004` cash `1.400.000`
- `DP2605050002` `Piere Njangka` transfer `145.000`
- `DP2605080002` expense doc `E20260508002` expense `-16.000` `Makan CGR`
- `PY2605080001` supplier `Mainan Toys` doc `testreport` cash `20.000.000`

## Manual runtime validation

Service-level live curl/download validation against running finance app not executed in this session.

Remaining blocked runtime checks:
- generate live endpoint response from running app with scenario issue
- generate live workbook from endpoint and compare output file end-to-end with issue sample
