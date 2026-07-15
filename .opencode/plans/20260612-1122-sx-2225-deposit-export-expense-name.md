# Plan — SX-2225 Payment Deposit Report export expense name

Task ID: `20260612-1122-sx-2225-deposit-export-expense-name`
Readiness: `ready-for-implementation`
Quality Gate: `PASS_FOR_SLICE`
Primary source of truth: this file.

## Goal
Fix backend export Payment Deposit Report agar expense dari mobile, termasuk `E20260611001`, menampilkan Expense Name dengan format `code - name`, contoh `000 - Uang Parkir`, di file Excel.

## Non-goals
- Tidak mengubah flow create expense mobile pada slice ini.
- Tidak membuat backfill/migration, karena `acf.expense` hanya menyimpan `expense_type_id` dan fix read-side akan berlaku untuk data historis.
- Tidak mengubah route, auth, response envelope, download history flow, atau format kolom selain value `Expense Name`.
- Tidak menyentuh module selain `finance` kecuali validasi discovery tambahan membuktikan wajib.

## Scope
Target utama:
- `finance/repository/payment_deposit_report_repository.go`
- `finance/service/payment_deposit_report_service_test.go`
- Tambahan test repository bila pola project memungkinkan: `finance/repository/payment_deposit_report_repository_test.go`

Scope operasional:
- Query export AR expense section di `buildDownloadARQuery`.
- Excel mapping tetap memakai kolom existing `ExpenseName` di `PaymentDepositReportDownloadRow`.

## Requirements
- Expense dari mobile dan web harus menampilkan Expense Name di export Excel.
- Format final kolom `Expense Name` adalah `{expense_type_code} - {expense_type_name}`.
- Join expense type pada export harus memakai `expense_type_id` sebagai PK global, tanpa filter `etr.cust_id = parentCustId`.
- Baris expense tidak boleh hilang ketika expense type tidak ditemukan; tetap `LEFT JOIN` dan fallback string kosong.
- Tidak boleh menambah N+1 query; resolusi tetap di query export.
- Recap/amount/payment columns tidak berubah.

## Acceptance Criteria
- Untuk fixture expense mobile dengan `expense_type_id` yang valid tetapi `cust_id` expense type tidak sama dengan `parentCustId`, query export tetap dapat menghasilkan `000 - Uang Parkir`.
- Baris web-created expense tetap tampil, dengan format baru `code - name`.
- Kolom Excel `Expense Name` tetap kolom Q sesuai layout existing.
- `E20260611001` pada filter `11/06/2026 - 11/06/2026`, deposit type All, collector All, menampilkan `000 - Uang Parkir` setelah export ulang.
- Tidak ada perubahan pada total `Expense`, `Cash`, `Cheque / Giro`, `Transfer`, `Return`, `Credit / Debit`, `Discount`, dan `Payment Balance`.

## Existing Patterns/Reuse
- Reuse `buildDownloadARQuery` sebagai pusat fix query; tidak perlu layer service baru.
- Reuse `PaymentDepositReportDownloadRow.ExpenseName` dan `generateExcel`; service sudah menulis value `ExpenseName` ke kolom Q.
- Reuse test service existing `finance/service/payment_deposit_report_service_test.go` untuk memvalidasi Excel output.
- Reuse dry-run SQL test style dari repository test existing bila menambah test query builder.
- Reuse repo tenant rule dari `.opencode/docs/ARCHITECTURE.md`: transaksi memakai `custId`, master parent-company biasanya memakai `parentCustId`; untuk kasus ini `expense_type_id` adalah PK global sehingga join PK-only dipilih berdasarkan keputusan user.

## Constraints
- Jalankan command dari direktori `finance`, bukan root.
- Shell workflow repo ini harus memakai prefix `rtk`.
- Jangan copy atau memperluas credentials/env tracked.
- Pertahankan Controller → Service → Repository → DB.
- `acf.deposit`, `acf.deposit_expense`, dan `acf.expense` filter `cust_id`/`deleted_at` existing tidak boleh dilonggarkan.
- `LEFT JOIN acf.expense_type` tetap dipakai agar row expense tidak hilang.

## Risks
- Menghapus kondisi `etr.cust_id = parentCustId` berarti nama expense type diselesaikan murni dari PK. Ini aman jika `expense_type_id` benar-benar unik global seperti model `primaryKey;autoIncrement`, tapi tetap catat bukti SQL/DB saat implementasi.
- Perubahan format ke `code - name` akan mengubah output web yang sebelumnya kemungkinan hanya `expense_type_name`; user sudah memilih format ini.
- Mobile lookup tetap tidak di-scope `cust_id`; itu root cause data path, tapi tidak diperbaiki di slice ini. Follow-up disarankan bila ingin mencegah pilihan expense type tenant lain di masa depan.
- Jika ada data historis dengan `expense_type_id` NULL/0/tidak valid, output tetap kosong. Itu di luar bukti saat ini dan perlu data repair terpisah.

## Decisions/Assumptions
- Keputusan user: format kolom harus `code - name`.
- Keputusan user: fix export join pada PK saja; drop kondisi `cust_id` di join `expense_type`.
- Keputusan user: tidak perlu backfill.
- Asumsi berdasarkan code: `acf.expense_type.expense_type_id` adalah PK auto-increment global dan cukup unik untuk join.
- Open question non-blocking: apakah mobile expense-type lookup perlu discope parent customer di tiket terpisah.

## Source Strategy
- Local project discovery: digunakan; file finance/mobile terkait sudah diinspeksi.
- Official docs/context7: diskip karena tidak ada perilaku library eksternal yang material.
- GitHub/web search: diskip karena bug repo-local.
- Browser/screenshot: diskip untuk planning; manual UI repro/export direncanakan setelah implementasi.
- Evidence retained: `.opencode/evidence/20260612-1122-sx-2225-deposit-export-expense-name/discovery.md` dan `index.json`.

## Execution Source of Truth
Prioritas executor:
1. Instruksi eksplisit terbaru dari user.
2. Safety/security/tenant rules repo.
3. Non-negotiable Implementation Invariants di plan ini.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria dan Done Criteria.
6. Implementation Steps.
7. Follow-up non-blocking.

Jika konflik, ikuti prioritas lebih tinggi dan catat konflik di evidence verifikasi.

## Non-negotiable Implementation Invariants
- Planner artifact-only; plan ini belum mengubah source.
- Perubahan produksi utama harus berada di repository/query export finance.
- `deposit_expense`, `expense`, dan `deposit` joins/filter existing tidak boleh dilonggarkan.
- `expense_type` harus tetap `LEFT JOIN`, bukan `JOIN`.
- Join `expense_type` harus resolve berdasarkan `expense_type_id` tanpa `etr.cust_id = parentCustId`.
- Output `expense_name` harus `code - name` ketika code dan name tersedia.
- Tidak boleh ada N+1 query atau lookup per-row di service.
- Tidak boleh membuat migration/backfill untuk slice ini.

## Do Not / Reject If
- Reject jika fix hanya mengubah Excel mapper tanpa memperbaiki query join.
- Reject jika `Expense Name` masih hanya `expense_type_name` saat code tersedia.
- Reject jika join `expense_type` masih mensyaratkan `etr.cust_id = ?` sehingga mobile record bisa kosong lagi.
- Reject jika `LEFT JOIN acf.expense ex` atau filter `ex.cust_id = d.cust_id` dihapus.
- Reject jika row expense hilang karena `LEFT JOIN` diubah menjadi `JOIN`.
- Reject jika perubahan menyentuh env, credentials, package files, lockfiles, atau module unrelated.

## Diff Boundary
Allowed:
- `finance/repository/payment_deposit_report_repository.go`
- `finance/service/payment_deposit_report_service_test.go`
- `finance/repository/*payment_deposit_report*_test.go` bila dibutuhkan
- `.opencode/evidence/20260612-1122-sx-2225-deposit-export-expense-name/**` untuk bukti implementasi

Generated report exception:
- Temporary local `.xlsx`/base64 artifacts boleh dibuat hanya di temp/evidence bila perlu, jangan commit binary export kecuali user minta.

Out-of-boundary changes harus direvert atau dijustifikasi di evidence sebelum final quality gate.

## TDD/Test Plan
TDD required: yes. Ini bug backend export/query.

Existing test patterns:
- Service Excel tests ada di `finance/service/payment_deposit_report_service_test.go`.
- Repository dry-run SQL pattern ada di test repository finance.

First failing/regression test:
1. Tambah/ubah test repository untuk `buildDownloadARQuery` yang memastikan `LEFT JOIN acf.expense_type etr ON etr.expense_type_id = ex.expense_type_id` tidak lagi mengandung `AND etr.cust_id = ?`.
2. Tambah test service Excel yang memastikan `PaymentDepositReportDownloadRow{ExpenseName: ptr("000 - Uang Parkir")}` muncul di cell Q data row.

Green step:
- Ubah SELECT expense_name di AR expense subquery menjadi format `code - name`.
- Ubah join expense_type menjadi PK-only.

Suggested SQL expression:
```sql
COALESCE(
  NULLIF(
    CONCAT_WS(' - ', NULLIF(etr.expense_type_code, ''), NULLIF(etr.expense_type_name, '')),
    ''
  ),
  ''
) AS expense_name
```
Jika project DB/Postgres behavior perlu lebih eksplisit, pakai `CASE` agar ketika code/name salah satu kosong, fallback tetap masuk akal.

Refactor step:
- Jika ekspresi format terlalu panjang/duplikatif, pertahankan inline saja bila hanya satu lokasi. Jangan over-engineer helper SQL global.

Edge cases:
- code dan name tersedia → `000 - Uang Parkir`.
- code kosong, name tersedia → `Uang Parkir` atau behavior ekspresi yang disepakati; jangan menghasilkan ` - Uang Parkir`.
- expense_type tidak ditemukan → string kosong dan row tetap ada.
- web-created expense tetap terisi format baru.

Commands:
- `rtk go test ./repository -run TestPaymentDepositReport`
- `rtk go test ./service -run TestPaymentDepositReportService`
- `rtk go test ./...`

## Implementation Steps
1. Dari repo root, cek runtime posture bila manual repro DB/UI akan dilakukan: `rtk docker compose -f docker-compose.yml ps`.
2. Dari `finance`, jalankan targeted tests awal untuk melihat baseline: `rtk go test ./service -run TestPaymentDepositReportService` dan repository target bila ada.
3. Tambahkan regression test query builder untuk join expense_type PK-only dan format `expense_type_code` + `expense_type_name` pada `expense_name`.
4. Tambahkan/ubah test Excel service untuk memastikan value `000 - Uang Parkir` ditulis ke kolom `Expense Name`.
5. Update `finance/repository/payment_deposit_report_repository.go` pada AR expense branch:
   - SELECT `expense_name` memakai `expense_type_code` + `expense_type_name`.
   - Join `acf.expense_type etr` hanya pada `etr.expense_type_id = ex.expense_type_id`.
   - Update `GROUP BY` agar mencakup `etr.expense_type_code` dan `etr.expense_type_name` sesuai ekspresi SELECT.
   - Hapus argumen `parentCustId` yang sebelumnya hanya dipakai join expense_type di branch expense; pastikan urutan args tetap benar.
6. Jalankan targeted tests sampai green.
7. Jalankan `rtk go test ./...` dari `finance`.
8. Jika runtime/DB tersedia, manual repro filter SX-2225 dan unduh Excel; bukti harus menyebut cell/row `E20260611001` menampilkan `000 - Uang Parkir`.
9. Catat root-cause note untuk Jira.
10. Route final review ke `@quality-gate` karena ini bug report/export dengan tenant-data implication.

## Expected Files to Change
- `finance/repository/payment_deposit_report_repository.go`
- `finance/service/payment_deposit_report_service_test.go`
- `finance/repository/payment_deposit_report_repository_test.go` atau file test repository existing yang sesuai
- Evidence implementation under `.opencode/evidence/20260612-1122-sx-2225-deposit-export-expense-name/`

## Agent/Tool Routing
- `@fixer`: implementasi dan tests bounded di `finance`.
- `@explorer`: opsional bila implementor perlu cek DB/query tambahan.
- `@quality-gate`: final review setelah tests/manual evidence.
- `@architect`: tidak diperlukan untuk slice ini karena keputusan sudah sempit dan repo-local.

## Executor Handoff Prompt
Implement SX-2225 using `.opencode/plans/20260612-1122-sx-2225-deposit-export-expense-name.md` as source of truth. Scope is finance export only. Must preserve deposit/expense cust_id filters, keep `expense_type` as `LEFT JOIN`, remove `etr.cust_id = parentCustId` from that join, and output Expense Name as `code - name` when both fields exist. Do not touch mobile write flow, migrations, env, package files, or unrelated reports. Add regression tests before/with the fix. Validate from `finance` with targeted `rtk go test` commands and `rtk go test ./...`. Return changed files, test output, and if runtime is available, manual export evidence for `E20260611001` showing `000 - Uang Parkir`. If any requirement conflicts with newer user instruction or safety rules, stop and report the conflict.

## Execution-ready Worklist / Handoff Contract
start_with: `T1`

### T1 — Add failing repository query regression
- action: Add a test proving AR expense export query resolves `expense_name` from `expense_type_code` + `expense_type_name` and does not scope `expense_type` by `cust_id`.
- depends_on: none
- owner/lane: `@fixer`
- validation: `rtk go test ./repository -run TestPaymentDepositReport`
- exit criteria: Test fails before query fix or explicitly captures current mismatch.
- blocking status: ready
- requires_user_decision: no
- must_preserve: no production source change yet except test.
- do_not_touch: env, migrations, mobile.
- evidence_update: record test name and failure/green output.
- exit_verification: failing/targeted test evidence captured.

### T2 — Add Excel value regression
- action: Add/extend service test proving `ExpenseName = "000 - Uang Parkir"` is written into the export row under `Expense Name`.
- depends_on: none
- owner/lane: `@fixer`
- validation: `rtk go test ./service -run TestPaymentDepositReportService`
- exit criteria: Test covers final display string.
- blocking status: ready
- requires_user_decision: no
- must_preserve: Excel column order remains unchanged.
- do_not_touch: production mapper unless needed after T3.
- evidence_update: record test name/output.
- exit_verification: targeted service test evidence captured.

### T3 — Fix export query
- action: Update `buildDownloadARQuery` AR expense branch to join `expense_type` by PK only and select `expense_name` as `code - name`.
- depends_on: T1, T2
- owner/lane: `@fixer`
- validation: targeted repository + service tests.
- exit criteria: Query args order is correct, `GROUP BY` includes needed fields, tests pass.
- blocking status: ready
- requires_user_decision: no
- must_preserve: `LEFT JOIN`, existing `cust_id` filters on deposit/deposit_expense/expense, no N+1.
- do_not_touch: mobile module, migrations, env.
- evidence_update: record diff summary and tests.
- exit_verification: `rtk go test ./repository -run TestPaymentDepositReport` and `rtk go test ./service -run TestPaymentDepositReportService` pass.

### T4 — Full finance validation
- action: Run full service test suite for finance.
- depends_on: T3
- owner/lane: `@fixer`
- validation: `rtk go test ./...`
- exit criteria: Full suite passes, or failures are unrelated and documented with evidence.
- blocking status: ready
- requires_user_decision: no
- must_preserve: no package/lockfile changes unless explicitly justified.
- do_not_touch: unrelated modules.
- evidence_update: record full output summary.
- exit_verification: full command result in final summary/evidence.

### T5 — Manual repro/export check if runtime is available
- action: Re-run SX-2225 filter and inspect Excel row for `E20260611001`.
- depends_on: T4
- owner/lane: `@fixer` with optional browser/manual support
- validation: exported `.xlsx` row shows `000 - Uang Parkir`.
- exit criteria: Manual evidence captured, or runtime/login/download blocker documented.
- blocking status: ready
- requires_user_decision: no
- must_preserve: do not commit exported binary unless requested.
- do_not_touch: production DB data except read-only export action.
- evidence_update: record row/cell observation and screenshot/path if available.
- exit_verification: manual confirmation or blocker recorded.

### T6 — Quality gate
- action: Run final review against plan invariants and evidence.
- depends_on: T4, T5 if available
- owner/lane: `@quality-gate`
- validation: compare diff, tests, evidence, tenant implications.
- exit criteria: PASS or actionable findings.
- blocking status: ready
- requires_user_decision: no
- must_preserve: no unreviewed scope creep.
- do_not_touch: source files from quality gate.
- evidence_update: quality-gate result.
- exit_verification: final signoff status.

## Validation Commands
From `finance`:
```bash
rtk go test ./repository -run TestPaymentDepositReport
rtk go test ./service -run TestPaymentDepositReportService
rtk go test ./...
```

Optional runtime/manual from repo root first:
```bash
rtk docker compose -f docker-compose.yml ps
```

## Evidence Requirements
- Test evidence for repository query and Excel output.
- Diff summary proving only allowed files changed.
- If manual repro is possible, row evidence for `E20260611001` showing `000 - Uang Parkir`.
- If manual repro is not possible, exact blocker: service down, credentials unavailable, download unavailable, or DB inaccessible.
- Jira root-cause note.

## Done Criteria
- Tests added and passing.
- Export query no longer loses expense type name because of `parentCustId` join mismatch.
- `Expense Name` format is `code - name`.
- No backfill/migration added.
- No mobile/source unrelated changes unless separately approved.
- Quality gate has enough evidence to pass or flag remaining risk.

## Final Planning Summary
Artifacts created/consulted:
- Created primary plan: `.opencode/plans/20260612-1122-sx-2225-deposit-export-expense-name.md`.
- Kept evidence: `.opencode/evidence/20260612-1122-sx-2225-deposit-export-expense-name/discovery.md` and `index.json` because implementation needs the code-root-cause details.
- Consulted repo-local docs: `.opencode/docs/ARCHITECTURE.md`, `.opencode/docs/QUALITY.md`.
- Consulted source files listed in discovery.

Key decisions:
- Fix read-side in finance export.
- Join `expense_type` by PK only.
- Format output as `code - name`.
- No backfill.

Assumptions/open questions:
- Assumption: `expense_type_id` is globally unique as PK auto-increment.
- Open non-blocking follow-up: mobile expense-type lookup should likely be scoped by customer/parent customer in a separate hardening ticket.

Readiness:
- `ready-for-implementation` for a bounded finance slice.
- `PASS_FOR_SLICE` because the mobile lookup hardening remains a separate non-blocking risk.

Cleanup performed:
- No draft artifacts created.
- Evidence intentionally kept for replayability and implementation handoff.
