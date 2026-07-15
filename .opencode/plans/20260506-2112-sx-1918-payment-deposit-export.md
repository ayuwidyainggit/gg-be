# Plan — SX-1918 Payment Deposit Report Export

## Goal

Implementasi dan debug endpoint backend `GET /finance/v1/reports/payment-deposit/download` agar FE dapat membuat file Excel Payment Deposit Report, menyimpannya sebagai base64 di `report.list.file_base64`, dan menampilkan status di Download History.

## Non-goals

- Tidak mengubah kontrak route legacy `/v1/reports/payment-deposit/download` kecuali ada task konsolidasi terpisah.
- Tidak membuat UI/FE Download History.
- Tidak menyalin asset/template spreadsheet eksternal yang tidak tersedia di repo.
- Tidak mengubah status convention global `report.list` tanpa konfirmasi FE/DB bila status failed belum jelas.

## Scope

- Modul target: `finance/`.
- Endpoint target: `GET /finance/v1/reports/payment-deposit/download`.
- Jalur source of truth: `payment_deposit_report_controller.go`, `payment_deposit_report_service.go`, `payment_deposit_report_repository.go`, `payment_deposit_report.go`, `report_list.go`.
- Tambah migration finance untuk `report.list.file_base64` bila kolom belum ada di environment.
- Perbaiki query dataset export agar detail-level sesuai AR invoice, AR expense, dan AP spec.
- Perbaiki lifecycle async report list: create processing, generate, update ready/base64, failure handling sesuai decision.

## Requirements

1. Endpoint `GET /finance/v1/reports/payment-deposit/download` tersedia dan memakai `middleware.JWTProtected()`.
2. Query params didukung:
   - `page`, `limit`, `sort` dengan default aman.
   - `start_date`, `end_date` wajib, mendukung epoch dan `YYYY-MM-DD` sesuai helper existing.
   - `salesman_id` comma-separated atau array; wajib untuk AR sesuai requirement, tidak diterapkan ke AP.
   - `deposit_no` optional comma-separated atau array.
   - `deposit_type` optional; support `AR`, `AP`, atau `AR,AP`.
3. `cust_id` wajib dari JWT local, bukan query.
4. Saat request valid, insert record `report.list` dengan `report_name` format `DownloadDepositPayment-DDMMYY-NNN`, `file_status=0`, `file_url=NULL`, `file_base64=NULL/empty`.
5. Generate Excel async agar request mengembalikan response processing.
6. Excel memakai kolom detail sesuai requirement: `deposit_date`, `deposit_type`, `deposit_no`, `collector`, `document_date`, `code`, `business_name`, `document_no`, `cash`, `cheque_bg`, `transfer`, `return`, `credit_debit`, `discount`, `payment_balance`, `expense`, `expense_name`.
7. Setelah sukses, update `report.list.file_status=1` dan `file_base64=<base64 XLSX>`.
8. Response processing:
   - HTTP `200`
   - `message = "Processing time may vary by file size. Please check Download History to access the file"`
   - `data = null`
9. Response error mengikuti pola controller finance existing: parser error `422`, validation/service error `400`.

## Acceptance Criteria

- [ ] Request tanpa token ditolak oleh middleware auth.
- [ ] Request valid pada `/finance/v1/reports/payment-deposit/download` membuat row `report.list`.
- [ ] `report_name` sesuai `DownloadDepositPayment-DDMMYY-001` dan increment per tanggal untuk prefix `DownloadDepositPayment-DDMMYY`.
- [ ] `file_status=0` saat job dibuat.
- [ ] Setelah generate selesai, `file_status=1` dan `file_base64` terisi string base64 valid untuk XLSX.
- [ ] Response processing sesuai spec dan `data=null`.
- [ ] `start_date`/`end_date` tersimpan dari filter report, bukan tanggal sekarang/+30 hari.
- [ ] `deposit_type` `AR`, `AP`, dan `AR,AP` memilih dataset yang benar; default bila tidak dikirim mengikuti decision di bawah.
- [ ] AR invoice query menghasilkan detail per invoice/document, bukan hanya aggregate per deposit.
- [ ] AR expense query menghasilkan row expense dengan amount non-zero.
- [ ] AP query menghasilkan detail per invoice/payment detail dan collector kosong.
- [ ] Semua query memakai binding parameter dan tenant filter `cust_id`.
- [ ] `salesman_id` multiple diparse aman dan hanya memfilter AR `acf.deposit.emp_id`.
- [ ] `deposit_no` optional tidak membuat `IN ()`.
- [ ] Data kosong tetap menghasilkan Excel dengan header.
- [ ] Test relevan pass di modul `finance`.

## Existing Patterns/Reuse

- Reuse route target `controller/payment_deposit_report_controller.go` karena sudah persis `/finance/v1/reports/payment-deposit/download` dan authenticated.
- Reuse `payment_deposit_report_service.go` async model: insert processing, goroutine generate, update base64 ready.
- Reuse repository tx extraction `model(ctx)` dan `service.Transaction.WithinTransaction` untuk writes.
- Reuse `excelize/v2` untuk XLSX generation.
- Reuse `model.ReportList` yang sudah memiliki `FileBase64`.
- Reuse sales migration pattern `ALTER TABLE report.list ADD COLUMN IF NOT EXISTS file_base64 TEXT`.
- Reuse sales download pattern untuk report naming dan failure handling concept, tetapi jangan langsung mengadopsi sales status value tanpa decision.
- Existing code/pattern cukup; tidak ditemukan KiloCode/project utility lain yang menggantikan query/export logic ini.

## Constraints

- Ikuti invariant Controller → Service → Repository → DB.
- Repository tidak boleh memuat business workflow; hanya query/persistence.
- Tulis ke `report.list` harus via transaction di service layer.
- Jangan hardcode `cust_id`, tanggal, salesman, atau deposit number dari dokumen.
- Gunakan parameter binding untuk semua filter user.
- Finance module memakai Go 1.18; hindari API Go lebih baru.
- Jangan copy/expose plaintext credential dari repo.
- Ada conflict instruksi shell: global OpenCode melarang `rtk`, repo `AGENTS.md` meminta `rtk`; untuk implementasi oleh OpenCode, ikuti instruksi global kecuali user eksplisit meminta RTK.

## Risks

- Dua implementasi paralel dapat membingungkan reviewer. Rencana ini memilih `payment_deposit_report_*` sebagai source of truth SX-1918.
- Current target implementation masih aggregate-level, sehingga harus ubah model/query/Excel ke detail-level.
- Default `deposit_type` belum final; existing target default `AR`, requirement menyarankan aman `AR+AP` jika tidak dikirim.
- `file_status` failed belum jelas di finance; tanpa failed status, async job bisa stuck processing.
- Running number count-based rentan race pada concurrent request.
- Goroutine async tanpa worker durable bisa hilang jika service restart setelah row processing dibuat.
- Base64 XLSX besar disimpan di DB text; perlu pertimbangkan ukuran record dan timeout query Download History.

## Decisions/Assumptions

- **Decision: source of truth** — SX-1918 diimplementasikan pada `payment_deposit_report_*` karena itulah route `/finance/v1/...`. `report_payment_deposit_*` dianggap legacy/alternate dan tidak diperluas.
- **Decision: async** — download tetap async dan response utama adalah processing, karena report bisa besar dan message spec mengarahkan user ke Download History.
- **Decision: persisted date range** — `report.list.start_date` dan `end_date` menyimpan filter `start_date`/`end_date`, bukan tanggal sekarang/+30, agar Download History merefleksikan range report.
- **Assumption: default deposit_type** — ubah default download menjadi `AR,AP` bila tidak dikirim, karena requirement menyebut query AR, AP, AR+AP dan opsi aman default gabungan. Jika FE ternyata mengirim `deposit_type` selalu, perubahan ini tetap kompatibel.
- **Assumption: salesman_id** — wajib untuk AR sesuai requirement; untuk AP tidak dipakai. Bila default AR+AP dan `salesman_id` kosong, controller harus menolak karena AR masuk scope. Bila `deposit_type=AP` saja, `salesman_id` boleh kosong.
- **Open question: failed status** — Finance current Payment Deposit memakai `0=Processing`, `1=Ready`; sales memakai `3=Failed`. Sebelum implement failed status, konfirmasi apakah finance/Download History mengenal failed status. Jika tidak terkonfirmasi, minimal log error dan buat follow-up; lebih baik tambah status failed sesuai standar report global jika FE mendukung.
- **Open question: sync success response** — Existing endpoint async mengembalikan processing. Response success spec kemungkinan untuk Download History/list setelah ready, bukan immediate download response. Implementasi endpoint download tidak perlu menunggu ready.

## TDD/Test Plan

### TDD Required

Ya. Ini menyentuh API behavior, query report, validation, async lifecycle, dan persistence `report.list`.

### Reason

Risiko utama adalah regresi kontrak endpoint, salah source route, salah query AR/AP, dan job stuck processing. Test harus membuktikan perilaku sebelum production code diubah.

### Existing Test Patterns

- `finance/controller/payment_deposit_report_controller_test.go`: table test helper normalisasi/validasi.
- `finance/service/payment_deposit_report_service_test.go`: mock repository/service mapping.
- `finance/repository/payment_deposit_report_repository_test.go`: DryRun GORM dan assertion SQL string/args.

### First Failing / Regression Tests

1. Controller validation:
   - `normalizeAndValidatePaymentDepositDownloadFilter` default deposit type menjadi `[]string{"AP", "AR"}` atau equivalent sorted.
   - `salesman_id` required saat AR termasuk scope; tidak required saat `deposit_type=AP`.
2. Repository SQL:
   - AR query contains invoice/detail joins from spec (`acf.deposit_payment`, `sls.order`, `mst.m_outlet`, `acf.deposit_detail`) and projects detail columns.
   - AR expense query includes `acf.deposit_expense`, `acf.expense`, `acf.expense_type`, `expense_name`, and non-zero expense filter.
   - AP query includes `acf.account_payable_payment_detail`, `acf.account_payable`, `mst.m_supplier`, `discount`, `payment_balance`.
   - `AR+AP` uses `UNION ALL` with consistent selected columns.
3. Service lifecycle:
   - `DownloadReport` inserts `report.list` with filter date range and `file_status=0`.
   - Async success updates `file_status=1` and non-empty base64. Use a channel-enabled mock to avoid flaky sleeps.
   - Report name sequence is prefix/date scoped.

### Green Step

- Extend model row/entity fields for detail dataset.
- Refactor repository builder into composable AR invoice, AR expense, AP branch queries with shared column schema.
- Update Excel writer columns/order and null numeric defaults.
- Add migration file under `finance/migration/report.list/add_file_base64.sql`.
- Update service naming/running-number and async failure path.

### Refactor Step

- Remove duplicated date/CSV normalization where safe, but avoid broad consolidation with legacy route.
- Extract constants for report name prefix, processing message, status values, and Excel headers.
- Keep SQL builders testable and avoid string-concatenating untrusted values.

### Edge Cases

- Empty dataset creates workbook with title/header and summary zero.
- `deposit_no` empty tokens ignored.
- Multiple `salesman_id` parsed safely.
- Date epoch seconds and `YYYY-MM-DD` normalized to inclusive day range; if SQL column is timestamp, consider end date end-of-day.
- Null numeric fields become `0`.
- AP collector remains empty string/nil in Excel.
- Async fetch/generate/update failure does not silently leave stuck processing if failed status is decided.
- Concurrent two requests for same customer/date do not generate duplicate `NNN`; if DB cannot guarantee, document residual risk and recommend unique index/advisory lock follow-up.

### Commands

Run from `finance/`:

```bash
go test ./controller -run PaymentDeposit
go test ./repository -run PaymentDeposit
go test ./service -run PaymentDeposit
go test ./...
```

Optional runtime smoke after services are up:

```bash
curl --location --request GET 'http://localhost:9005/finance/v1/reports/payment-deposit/download?start_date=2024-11-01&end_date=2024-11-30&salesman_id=1&deposit_no=DP251212001,DP251212002&page=1&limit=10' --header 'Authorization: Bearer <access_token>'
```

## Implementation Steps

1. **Lock route target**
   - Verify `finance/main.go` still registers `paymentDepositReportController.Route(app)`.
   - Do not route SX-1918 through `reportPaymentDepositController`.
2. **Add constants**
   - In entity/service package, define `DownloadDepositPayment` prefix, processing message, and finance file status constants: `0 Processing`, `1 Ready`; add failed only after decision.
3. **Migration**
   - Add `finance/migration/report.list/add_file_base64.sql` with `ALTER TABLE report.list ADD COLUMN IF NOT EXISTS file_base64 TEXT;` and comment.
4. **DTO/model expansion**
   - Extend `model.PaymentDepositReportRow` or add dedicated `PaymentDepositReportDownloadRow` for detail columns in required Excel.
   - Keep list endpoint aggregate fields intact if used by FE; prefer adding download-specific row rather than breaking list response.
5. **Controller validation**
   - Change download default `deposit_type` to AR+AP unless FE/spec decision says otherwise.
   - Require `salesman_id`/`emp_id` only when AR is included.
   - Keep parsing `deposit_no`, `salesman_id`, `deposit_type` as CSV/array.
6. **Repository detail queries**
   - Implement download-specific builder returning consistent columns for:
     - AR invoice rows from `acf.deposit`, `acf.deposit_payment`, `sls.order`, `mst.m_outlet`, `acf.deposit_detail`, `mst.m_employee`.
     - AR expense rows from `acf.deposit`, `acf.deposit_expense`, `acf.expense`, `acf.expense_type`, `mst.m_employee`.
     - AP rows from `acf.account_payable_payment`, `acf.account_payable_payment_options`, `acf.account_payable_payment_detail`, `acf.account_payable`, `mst.m_supplier`.
   - Use `UNION ALL` when multiple branches requested.
   - Apply `ORDER BY deposit_date, deposit_no` on outer query.
   - Keep all filters parameterized.
7. **Service report lifecycle**
   - Generate `report_name` by counting only existing `DownloadDepositPayment-DDMMYY-%` for same customer, not all reports created today.
   - Insert report list with filter date range and `file_status=0` in transaction.
   - Return processing response immediately.
   - Generate Excel in background using all rows, not page/limit.
   - Update ready/base64 in transaction.
   - Add failed-status update if status convention confirmed; otherwise log via existing logger and document limitation.
8. **Excel generation**
   - Update headers/order to match required dataset.
   - Format dates and numeric columns consistently.
   - Include zero/empty values for missing fields.
   - Keep summary optional; if existing template includes summary, adapt to detail columns without breaking header.
9. **Tests**
   - Add failing tests first per TDD plan.
   - Implement until all PaymentDeposit tests pass.
10. **Validation**
   - Run targeted package tests, then `go test ./...` in `finance`.
   - If DB accessible, manually call endpoint and inspect `report.list` row status/base64.

## Expected Files to Change

- `finance/controller/payment_deposit_report_controller.go`
- `finance/controller/payment_deposit_report_controller_test.go`
- `finance/service/payment_deposit_report_service.go`
- `finance/service/payment_deposit_report_service_test.go`
- `finance/repository/payment_deposit_report_repository.go`
- `finance/repository/payment_deposit_report_repository_test.go`
- `finance/model/payment_deposit_report.go`
- `finance/entity/payment_deposit_report.go`
- `finance/migration/report.list/add_file_base64.sql`
- Optional: `finance/model/report_list.go` only if tag/type needs adjustment.

## Agent/Tool Routing

- Implementation: `@fixer` / `opencode-fixer` for bounded code edits and tests.
- Architecture review after implementation: `@oracle` if query/lifecycle decisions changed materially.
- Security/privacy review: `@security-privacy-reviewer` if auth/tenant filtering or Download History access changes beyond this endpoint.
- Release/ops: `@release-engineer` if migration deployment/backfill/status convention needs coordinated rollout.
- Quality gate: `@quality-gate` after implementation and tests before commit/PR.

## Validation Commands

From repo root, first check services if runtime smoke needed:

```bash
docker compose -f docker-compose.yml ps
```

From `finance/`:

```bash
go test ./repository -run PaymentDeposit
go test ./service -run PaymentDeposit
go test ./...
```

Optional DB/manual validation:

```sql
SELECT report_id, report_name, start_date, end_date, file_status, length(file_base64) AS file_base64_len
FROM report.list
WHERE report_name LIKE 'DownloadDepositPayment-%'
ORDER BY created_at DESC
LIMIT 5;
```

## Evidence Requirements

- Keep this plan as source of truth.
- Discovery evidence was gathered in `.opencode/evidence/20260506-2112-sx-1918-payment-deposit-export/discovery.md` and consolidated here.
- No official docs/context7 required because implementation uses established local `excelize`, Fiber, and GORM patterns.
- No GitHub/web required because no upstream behavior or current external facts affect implementation.
- No browser/screenshot required because task is backend API/export.
- Implementation evidence should include:
  - test output for commands above,
  - sample response body for processing,
  - DB row showing status transition/base64 length,
  - note whether failed status was implemented or deferred.

## Done Criteria

- Endpoint `/finance/v1/reports/payment-deposit/download` returns processing response for valid authenticated request.
- Report row inserted with correct tenant, report name, filter dates, status processing.
- Background generation writes valid base64 XLSX and status ready.
- Excel columns and data match AR invoice, AR expense, AP required schema.
- Migration for `file_base64` exists in finance.
- Tests for controller validation, repository SQL builder, and service lifecycle pass.
- Remaining decisions, especially failed status, are either implemented or explicitly documented as follow-up with risk.

## Final Planning Summary

- Artifacts created:
  - `.opencode/plans/20260506-2112-sx-1918-payment-deposit-export.md` — primary source of truth.
  - `.opencode/evidence/20260506-2112-sx-1918-payment-deposit-export/discovery.md` — discovery details were created temporarily during planning and consolidated into this plan.
- Key decisions:
  - Implement SX-1918 on `payment_deposit_report_*`, not legacy `report_payment_deposit_*`.
  - Keep async download lifecycle.
  - Store filter date range in `report.list`.
  - Plan default `deposit_type` as AR+AP unless FE contradicts.
- Questions asked: none, because provided requirement plus existing route patterns are enough for an implementation-ready plan. Material unknowns are documented as open questions/assumptions.
- Remaining open questions:
  - Confirm finance failed `file_status` value for async generation failures.
  - Confirm FE expectation if `deposit_type` omitted; plan assumes AR+AP.
- Readiness: ready for implementation with one conditional decision on failed status. If failed status is not known, implementation can still deliver processing→ready flow and document failure-status follow-up.
- Cleanup performed: discovery findings were consolidated into this primary plan and `.opencode/evidence/20260506-2112-sx-1918-payment-deposit-export/` was deleted as stale. No draft/evidence artifacts are kept.
