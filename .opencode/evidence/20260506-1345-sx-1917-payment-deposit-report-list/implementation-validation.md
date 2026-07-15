# Implementation Validation — SX-1917 Payment Deposit Report List

## Dokumen tambahan

- `docs/Report - Payment Deposit_BE.md` dibaca.
- `docs/Enhance Payment Deposit report.xlsx` diekstrak via `@document-specialist`.
  - Sheet: `List`, `Export`, `hasil cek Export`.
  - Konfirmasi list fields: `deposit_date`, `deposit_type`, `deposit_no`, collector, payment breakdown, `expense`, `total`.
  - Konfirmasi formula: `cash + cheque/bg + transfer + return + credit/debit - expense`.
  - Konfirmasi AR memakai collector/employee; AP collector null/blank.

## Validasi environment

- Root command: `rtk docker compose -f docker-compose.yml ps`
- Hasil akhir: command berhasil, tetapi tidak ada service yang sedang berjalan/listed.
- Manual DB/API smoke tidak dijalankan karena service/token tidak tersedia dalam sesi.

## Quality gate pertama

Status: `FAIL`.

Temuan blocker/risiko yang diperbaiki:
1. SQL alias `collector_*` tidak match model tag `salesman_*`.
2. Download route berisiko regresi karena ikut validasi list dan membutuhkan `deposit_type`.
3. AP options aggregate memakai `JOIN`, berisiko drop AP payment tanpa options.
4. Test service mapping belum menangkap collector nullable/non-null behavior.

## Perbaikan setelah quality gate

- Model/service mapping diselaraskan ke `collector_*`.
- Download dibuat backward-compatible dengan default `deposit_type=AR` saat absent.
- AP payment options aggregate diubah menjadi `LEFT JOIN`.
- Tests ditambah/diperbarui untuk AR collector mapping, AP null collector mapping, download default behavior, repository SQL shape, dan service mapping.

## Validasi test akhir

Dijalankan dari `/Users/ujang/Projects/Geekgarden/scylla-be/finance`:

```bash
rtk go test ./controller -run TestNormalizeAndValidatePaymentDepositFilter
```

Hasil: `Go test: 9 passed in 1 packages`

```bash
rtk go test ./repository -run 'TestBuildSafeSort|TestBuild.*PaymentDeposit|TestNormalizeDepositNoFilter'
```

Hasil: `Go test: 18 passed in 1 packages`

```bash
rtk go test ./service -run TestPaymentDeposit
```

Hasil: `Go test: 2 passed in 1 packages`

```bash
rtk go test ./...
```

Hasil: `Go test: 81 passed in 20 packages`

## Git/commit note

- `rtk git status --short` dari `/Users/ujang/Projects/Geekgarden/scylla-be` gagal: `fatal: not a git repository`.
- `rtk git status --short` dari `/Users/ujang/Projects/Geekgarden` juga gagal dengan pesan sama.
- Auto-commit tidak dapat dilakukan karena workspace ini tidak terdeteksi sebagai git repository.
