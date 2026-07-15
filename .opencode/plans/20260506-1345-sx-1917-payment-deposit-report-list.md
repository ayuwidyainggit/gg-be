# Plan — SX-1917 Payment Deposit Report List

Task ID: `20260506-1345-sx-1917-payment-deposit-report-list`

## Goal

Implementasikan list endpoint `GET /finance/v1/reports/payment-deposit` di module `finance/` agar mendukung report Payment Deposit gabungan AR, AP, dan AR+AP dengan filter, pagination, sorting aman, response contract, dan perhitungan amount sesuai Jira SX-1917.

## Non-goals

- Tidak mengubah root service architecture atau membuat module baru.
- Tidak mengimplementasikan export/download kecuali perubahan DTO/model minimum diperlukan agar compile tetap aman.
- Tidak mengubah schema database atau migration.
- Tidak mengubah endpoint lama `/v1/reports/payment-deposit` kecuali diperlukan untuk menghindari konflik compile.
- Tidak melakukan manual test ke environment `best.scyllax.online` tanpa token yang valid.

## Scope

- Target utama file existing baru yang route-nya sudah sesuai:
  - `finance/controller/payment_deposit_report_controller.go`
  - `finance/entity/payment_deposit_report.go`
  - `finance/model/payment_deposit_report.go`
  - `finance/repository/payment_deposit_report_repository.go`
  - `finance/service/payment_deposit_report_service.go`
  - test terkait di controller/repository/service bila perlu.
- Endpoint list harus menerima:
  - `deposit_type=AR`, `deposit_type=AP`, atau `deposit_type=AR,AP`
  - `start_date`, `end_date` sebagai epoch atau `YYYY-MM-DD`
  - `emp_id` CSV/array untuk AR only filter
  - `deposit_no` CSV/array untuk AR/AP branch masing-masing
  - `page`, `limit`, `sort`, dan opsional `q` bila mengikuti pattern aman.

## Requirements

1. Ambil `cust_id` dari JWT locals `c.Locals("cust_id")`.
2. Validasi `deposit_type` hanya `AR` dan/atau `AP`; invalid harus 400.
3. Date range inclusive dan `start_date <= end_date`.
4. AR branch:
   - Source utama `acf.deposit d`.
   - Tanggal `d.deposit_date`.
   - Collector `d.emp_id` join `mst.m_employee`.
   - Payment breakdown dari pre-aggregated `acf.deposit_payment` per `deposit_no, cust_id`.
   - Expense dari pre-aggregated `acf.deposit_expense` per `deposit_no, cust_id`.
   - `total_payment = cash + cheque + transfer + return + credit_debit - expense`.
5. AP branch:
   - Source utama `acf.account_payable_payment app`.
   - Tanggal `app.account_payable_payment_date`.
   - Nomor `app.account_payable_payment_no`.
   - Payment breakdown dari `acf.account_payable_payment_options appo`.
   - `expense_amount = 0`.
   - Collector fields `NULL`.
   - Wajib `app.deleted_by IS NULL`.
6. `emp_id` hanya diterapkan ke AR branch; AP tidak difilter oleh `emp_id`.
7. `deposit_no` diterapkan ke `d.deposit_no` pada AR dan `app.account_payable_payment_no` pada AP.
8. Pagination dan count dilakukan dari hasil branch final/union final.
9. Sorting memakai whitelist final alias, bukan raw interpolation dari request.
10. Semua amount memakai `COALESCE(..., 0)`.

## Acceptance Criteria

- `deposit_type=AR` mengembalikan hanya AR dengan collector dan expense.
- `deposit_type=AP` mengembalikan hanya AP dengan collector fields `null` dan expense `0`.
- `deposit_type=AR,AP` mengembalikan gabungan via `UNION ALL`.
- Filter `emp_id` tidak mempengaruhi AP ketika `AR,AP`.
- Filter `deposit_no` bekerja untuk AR dan AP sesuai source masing-masing.
- `items.length <= limit`, `total_data` sesuai total filtered, dan `total_page` hasil ceiling.
- `total_payment` sesuai formula dan tidak bergantung pada `d.total_payment`.
- Tidak ada hardcoded `cust_id`, tanggal, atau employee dari sample SQL.
- Sort aman dari SQL injection.
- Response root tetap memakai `responsebuild.BuildResponse` dan data payload berisi `items` serta `pagination` sesuai pattern existing/spec.

## Existing Patterns/Reuse

- Gunakan `PaymentDepositReportController` karena route existing sudah `app.Group("/finance/v1/reports/payment-deposit", middleware.JWTProtected())`.
- Gunakan response builder existing `responsebuild.BuildResponse`.
- Gunakan date helper existing `normalizeDateInput` yang sudah mendukung epoch dan `YYYY-MM-DD`.
- Gunakan pola query repository existing, tetapi ubah ke raw SQL builder/branch builder agar union, count, dan binding lebih terkendali.
- Gunakan test pattern existing GORM `DryRun`/helper unit untuk validasi SQL shape dan helper validation.
- Tidak ditemukan KiloCode/project utility yang sudah menyelesaikan AR+AP union Payment Deposit Report; perlu extend implementasi existing, bukan membuat domain baru.

## Constraints

- Module finance adalah Go module mandiri; perintah validasi dijalankan dari `finance/`.
- Root instructions konflik: global melarang `rtk`, repo meminta `rtk`. Untuk repo ini ikuti `AGENTS.md` lokal saat shell command: gunakan `rtk`.
- Docker daemon tidak aktif saat discovery; validasi DB/service lokal mungkin blocked sampai Docker dijalankan.
- Jangan expose atau menyalin secrets dari compose/workflow.
- Controller tidak boleh memanggil repository langsung; tetap Controller → Service → Repository.
- Write operations harus lewat transaction; list endpoint read-only tidak memerlukan transaction.

## Risks

- Ada dua jalur implementasi Payment Deposit Report (`payment_deposit_report_*` dan `report_payment_deposit_*`). Implementer harus fokus ke route `/finance/v1/reports/payment-deposit` dan tidak menghapus jalur lama tanpa konfirmasi.
- Mengubah `PaymentDepositReportQueryFilter` dapat mempengaruhi download endpoint yang memakai service/model sama. Mitigasi: pertahankan compatibility minimum atau sesuaikan download dengan field baru tanpa memperluas scope.
- AR soft delete field belum terkonfirmasi. Implementasi existing memakai `d.deleted_at IS NULL`; spec menyebut cek jika ada. Mitigasi: inspeksi model/schema/migration bila tersedia; jika tidak, pertahankan `d.deleted_at IS NULL` dan tambahkan `d.deleted_by IS NULL` hanya bila kolom terbukti ada.
- Join AP options tanpa `cust_id` mungkin bergantung schema. Mitigasi: cek model/migration/table usage untuk `account_payable_payment_options`; jika ada `cust_id`, tambahkan pada join/filter.
- Test unit SQL tidak membuktikan correctness aggregation aktual; butuh manual/integration SQL validation saat DB tersedia.

## Decisions/Assumptions

- **Keputusan:** Implementasi utama memakai `PaymentDepositReportController` pada route `/finance/v1/reports/payment-deposit`.
- **Keputusan:** `deposit_type` menjadi required untuk list sesuai SX-1917.
- **Keputusan:** `emp_id` menggantikan kebutuhan `salesman_id` untuk list endpoint baru; `salesman_id` boleh dipertahankan sebagai backward-compatible alias hanya bila tidak mengganggu spec.
- **Keputusan:** AP collector fields direpresentasikan nullable di model/entity (`*int`, `*string`) agar JSON `null`, bukan zero value.
- **Keputusan:** Default sort `deposit_date:desc`, mapping final ke alias `t.deposit_date DESC` untuk union.
- **Asumsi:** `q` opsional cukup mencakup `deposit_no` dan AR collector name bila mudah; bila tidak ada pattern reliable, jangan paksa sampai mengganggu acceptance utama.
- **Asumsi:** Summary existing di response boleh dipertahankan jika tidak merusak consumer; spec minimum mensyaratkan `items` dan `pagination`.
- **Open question rendah risiko:** Apakah download endpoint juga harus langsung support AR/AP union pada SX-1917? Scope user menekankan list endpoint; rencana ini tidak memblokir list pada jawaban tersebut.

## TDD/Test Plan

### Apakah TDD wajib?

Wajib. Ini perubahan production logic, validation, query behavior, pagination, dan security-sensitive sort handling.

### Existing test patterns

- `finance/controller/payment_deposit_report_controller_test.go` menguji `normalizeAndValidatePaymentDepositFilter`.
- `finance/repository/payment_deposit_report_repository_test.go` menguji `buildSafeSort`, normalize CSV, dan SQL shape dengan GORM DryRun.

### Red step — first failing/regression tests

1. Controller validation:
   - missing `deposit_type` → error.
   - invalid `deposit_type=XX` → error.
   - valid `deposit_type=AR,AP`, `start_date`, `end_date` tanpa `emp_id` → valid.
   - invalid date range → error.
   - invalid sort field/direction → error.
2. Parser/normalizer helper:
   - `deposit_type=AR, AP` menjadi set `AR`, `AP`.
   - `emp_id=381,421` menjadi `[]int{381,421}`.
   - `deposit_no=DP1,PY1` menjadi `[]string{"DP1","PY1"}`.
3. Repository SQL shape:
   - `deposit_type=AR` membangun query AR tanpa `UNION ALL`, memakai `acf.deposit`, subquery `acf.deposit_payment`, subquery `acf.deposit_expense`, dan filter `d.emp_id IN` saat `emp_id` ada.
   - `deposit_type=AP` membangun query AP tanpa `acf.deposit d`, memakai `acf.account_payable_payment`, `app.deleted_by IS NULL`, collector `NULL`, expense `0`, dan tidak punya `emp_id` filter.
   - `deposit_type=AR,AP` mengandung `UNION ALL`, count wrapping, dan final `ORDER BY t.deposit_date`.
   - sort injection-like payload tidak muncul raw dalam SQL; fallback/validation terjadi.
4. Service mapping:
   - AP row dengan null collector tetap menghasilkan JSON nullable, bukan `0`/`""`.
   - `TotalPage` ceiling benar untuk `total=72`, `limit=10`.

### Green step

- Update entity/model/helper agar tests compile dan pass.
- Implement branch query builder dengan parameter binding.
- Service map rows ke response baru.

### Refactor step

- Kurangi duplikasi SQL amount expression dengan helper function/string constants internal repository.
- Pastikan helper names tetap jelas: `buildARPaymentDepositQuery`, `buildAPPaymentDepositQuery`, `buildPaymentDepositUnionQuery`, `buildSafeFinalSort`.
- Pertahankan one-way layer dependency.

### Edge cases

- `deposit_type=AR` dengan `emp_id` kosong: return semua AR collector dalam date range.
- `deposit_type=AP` dengan `emp_id`: AP tetap muncul.
- `deposit_type=AR,AP` dengan `emp_id`: AR terfilter, AP tidak terfilter.
- Empty payment options: amount `0` bukan null.
- Empty expense: expense `0`.
- `limit=0` default 20; `page=0` default 1; negative values error.
- Large limit cap mengikuti existing `9999` atau standard project.

### Commands

Jalankan dari `finance/`:

```bash
rtk go test ./controller -run TestNormalizeAndValidatePaymentDepositFilter
rtk go test ./repository -run TestBuildSafeSort
rtk go test ./repository -run TestBuild.*PaymentDeposit
rtk go test ./service -run TestPaymentDeposit
rtk go test ./...
```

Jika Docker/DB tersedia, dari root repo:

```bash
rtk docker compose -f docker-compose.yml ps
rtk docker compose -f docker-compose.yml up -d finance redis
```

Lalu manual curl dengan token valid sesuai prompt Jira.

## Implementation Steps

1. **Entity/model update**
   - Tambahkan `DepositType []string` pada filter.
   - Tambahkan `EmpID []int` pada filter.
   - Pertimbangkan `Q string` bila implementasi search aman.
   - Tambahkan `DepositType string` pada item/row.
   - Ubah collector fields menjadi nullable untuk response AP:
     - model row: `*int`, `*string` atau `sql.Null*`.
     - entity item: `*int`, `*string`.
2. **Controller validation/normalization**
   - Hapus requirement `salesman_id` untuk list.
   - Tambahkan normalizer CSV untuk `deposit_type` dan `emp_id`.
   - Validasi `deposit_type` required dan hanya `AR`/`AP`.
   - Pertahankan date normalization epoch/`YYYY-MM-DD`.
   - Perluas sort whitelist: `deposit_date`, `deposit_no`, `deposit_type`, `collector_name`, `total_payment`.
3. **Repository query builder**
   - Ganti `buildQuery` menjadi builder union-aware.
   - AR select field order harus sama dengan AP:
     - `deposit_date`, `deposit_type`, `deposit_no`, `collector_id`, `collector_code`, `collector_name`, amount fields, `expense_amount`, `total_payment`.
   - AR payment subquery aggregate per `deposit_no, cust_id`.
   - AR expense subquery aggregate per `deposit_no, cust_id`.
   - AP query aggregate dengan `GROUP BY app.account_payable_payment_date, app.account_payable_payment_no`.
   - Count: `SELECT COUNT(1) FROM (<union-or-single>) t`.
   - Data: `SELECT * FROM (<union-or-single>) t ORDER BY <safe> LIMIT ? OFFSET ?`.
   - Summary existing: `SUM` dari same subquery jika tetap diperlukan.
4. **Service mapping**
   - Map `DepositType` ke JSON.
   - Map nullable collector fields.
   - Hitung pagination dari total dan limit; guard limit default.
   - Jangan recompute `total_payment` di service bila SQL sudah formula; cukup map.
5. **Tests**
   - Tambahkan Red tests dulu sesuai TDD plan.
   - Update existing tests yang masih mewajibkan `salesman_id`.
   - Tambahkan repository SQL shape tests untuk AR/AP/union.
6. **Validation**
   - Jalankan targeted tests lalu `rtk go test ./...` dari `finance/`.
   - Jika Docker daemon tersedia, lakukan manual smoke untuk AR, AP, AR+AP dengan token valid.

## Expected Files to Change

- `finance/entity/payment_deposit_report.go`
- `finance/model/payment_deposit_report.go`
- `finance/controller/payment_deposit_report_controller.go`
- `finance/controller/payment_deposit_report_controller_test.go`
- `finance/repository/payment_deposit_report_repository.go`
- `finance/repository/payment_deposit_report_repository_test.go`
- `finance/service/payment_deposit_report_service.go`
- Opsional bila compile terdampak:
  - `finance/service/payment_deposit_report_service_test.go`
  - `finance/main.go` hanya jika route conflict ditemukan saat validation, tidak direncanakan sebagai perubahan utama.

## Agent/Tool Routing

- Implementasi berikutnya: `@fixer` / `opencode-fixer` untuk bounded implementation + tests.
- Jika query union atau nullable mapping memicu tradeoff besar: `@oracle` review sebelum implementasi final.
- `@security-privacy-reviewer` tidak wajib, tetapi dapat dipakai bila ada perubahan auth/tenant isolation di luar `cust_id` filter.
- Tidak perlu `@designer`, browser, visual asset, atau mobile specialists.

## Validation Commands

Dari root repo sebelum code work/validation:

```bash
rtk docker compose -f docker-compose.yml ps
```

Dari `finance/`:

```bash
rtk go test ./controller -run TestNormalizeAndValidatePaymentDepositFilter
rtk go test ./repository -run TestBuildSafeSort
rtk go test ./repository -run TestBuild.*PaymentDeposit
rtk go test ./service -run TestPaymentDeposit
rtk go test ./...
```

Manual smoke dengan token valid:

```bash
curl --location --request GET 'http://localhost:9005/finance/v1/reports/payment-deposit?deposit_type=AR&start_date=2026-04-24&end_date=2026-04-24&emp_id=381,421&page=1&limit=10' --header 'Accept: application/json' --header 'Authorization: Bearer <access_token>'
curl --location --request GET 'http://localhost:9005/finance/v1/reports/payment-deposit?deposit_type=AP&start_date=2026-04-24&end_date=2026-04-27&page=1&limit=10' --header 'Accept: application/json' --header 'Authorization: Bearer <access_token>'
curl --location --request GET 'http://localhost:9005/finance/v1/reports/payment-deposit?deposit_type=AR,AP&start_date=2026-04-24&end_date=2026-04-27&emp_id=381,421&page=1&limit=10&sort=deposit_date:asc' --header 'Accept: application/json' --header 'Authorization: Bearer <access_token>'
```

## Evidence Requirements

- Keep this discovery evidence during implementation:
  - `.opencode/evidence/20260506-1345-sx-1917-payment-deposit-report-list/discovery.md`
- Implementation evidence should record:
  - Test outputs for targeted tests and `rtk go test ./...`.
  - SQL shape snippets or test assertions proving AR pre-aggregation, AP query, union, count wrapping, and sort whitelist.
  - Manual curl response snippets if token/DB available, without exposing token or secrets.
- Official docs/context7: skipped because no new library behavior or version-sensitive API is central.
- GitHub/web search: skipped because no upstream dependency/reference facts are required.
- Browser evidence: not applicable; backend JSON endpoint only.

## Done Criteria

- All acceptance criteria above pass.
- TDD tests for validation/query/service mapping pass.
- `rtk go test ./...` in `finance/` passes or any unrelated failures are documented with evidence.
- No raw `sort` interpolation from request.
- No hardcoded tenant/date/employee sample values.
- AP collector JSON fields are `null`.
- AR expense is not multiplied by payment joins.
- Primary route `/finance/v1/reports/payment-deposit` works with AR, AP, and AR+AP.

## Final Planning Summary

- Artifacts created:
  - Primary source of truth: `.opencode/plans/20260506-1345-sx-1917-payment-deposit-report-list.md`
  - Kept evidence: `.opencode/evidence/20260506-1345-sx-1917-payment-deposit-report-list/discovery.md`
- Draft artifacts: none created.
- Cleanup: no stale draft/evidence deleted; discovery evidence intentionally kept because it is operationally useful for implementer.
- Key decisions: extend existing `/finance/v1/reports/payment-deposit` implementation; implement union-aware repository; make collector nullable; default sort `deposit_date:desc`; `emp_id` applies only to AR.
- Questions: no blocking question asked because prompt provides sufficient behavior; remaining low-risk ambiguity is whether download endpoint should also support AR/AP union, assumed out-of-scope for list task unless compile compatibility requires minor adjustment.
- Readiness: ready for implementation by `@fixer` with TDD. Docker/service validation is partially blocked until Docker daemon is running, but unit tests can proceed.
