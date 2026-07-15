# Plan — SX-2258 Secondary Sales Summary (Order − Return)

- Task id: `20260618-1002-sx-2258-secondary-sales-summary`
- Issue: `SX-2258` `[Defect][FE] "Number of Product Sold" and "Discount and Promo" should calculated from sum of order type then subtracted from the return type`
- Endpoint: `GET /sales/v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001`
- Module: `sales` (Secondary Sales Report dashboard summary)
- Mode: Maintenance Stability Mode
- Readiness target: `ready-for-implementation`
- Latest feedback: PIC Widya (dengan Yogie) konfirmasi response `qty` masih salah meskipun commit `4ebacfe` sudah ada. Beda repo vs referensi baru: (a) `r.data_status = 6` masih ada di return branch, referensi terbaru hilangkan; (b) date filter repo pakai `>= AND <` (semi-open), referensi pakai `BETWEEN` (inclusive both ends). Perubahan filter ini yang bikin angka `134` tidak tercapai.

## Goal

Pertahankan dan kunci invariansi kalkulasi summary `secondary-sales/sum-date`:

- `Number of Product Sold` (response field `qty`) = `sum(order qty) - sum(return qty)`
- `Discount and Promo` (response field `total_discount_promo`) = `sum(order discount_promo) - sum(return discount_promo)`
- Field summary lain yang dependen (`gross_sale`, `ppn`, `net_sales_exc_ppn`, `net_sales`, `net_sales_return`, `return_rate`, `qty_return`) tetap konsisten dengan formula `order - return` sesuai acuan QA.

## Non-goals

- Tidak mengubah business logic report lain (`trend-sales`, `group`, `extract`, `pjp-sales`).
- Tidak mengubah format response JSON.
- Tidak migrasi schema atau tambah tabel.
- Tidak menulis/mengirim kredensial Jira atau token auth produksi ke commit.
- Tidak menyimpan JWT bearer token yang muncul di chat room ke file, log, evidence, atau commit. Token staging di rotasi sebelum plan lanjut.

## Scope

- File service target: `sales/repository/report_repository.go`, `sales/service/report_service.go`, `sales/controller/report_controller.go`.
- Test file target: `sales/repository/report_repository_test.go`, `sales/service/report_service_test.go`, `sales/controller/so_controller_test.go`.
- Validasi manual endpoint: environment local atau staging sesuai hak akses user; gunakan auth/env user sendiri, bukan kredensial Jira.

## Requirements

1. Filter query summary harus sama dengan acuan Widya terbaru:
   - `o.cust_id IN ?` (order)
   - `rd.cust_id IN ?` (return)
   - `o.data_status IN (6, 7)` di order branch dan return branch via join order.
   - `r.data_status = 6` harus DIHAPUS atau minimal dibuktikan tidak memengaruhi angka; referensi Widya terbaru tidak punya filter ini.
   - Date filter: `o.invoice_date` di kedua branch (bukan `r.return_date`).
   - Date semantics harus divalidasi: referensi Widya pakai `BETWEEN :date_from AND :date_to`; repo saat discovery pakai `>= ? AND < ?`. Pilih yang menghasilkan expected `qty=134` untuk `C260020001` Juni 2026 dan dokumentasikan alasan.
   - Optional: `outlet_ids`, `salesman_ids`, `pro_ids` (order: `o.outlet_id`, `o.salesman_id`, `od.pro_id`; return: `r.outlet_id`, `r.salesman_id`, `rd.product_id`).
2. Subtract arithmetic wajib:
   - `qty = os.qty - rs.qty_return`
   - `total_discount_promo = os.discount_promo - rs.discount_promo`
3. NULL safety: `COALESCE(..., 0)` di semua operand dan hasil akhir.
4. Return rate formula tetap: `ROUND(((rs.net_sales_inc_ppn / NULLIF(os.net_sales_inc_ppn, 0)) * 100)::numeric, 2)`.
5. Service mapping `model.SumReportByMonthModel` → `entity.SumReportByMonthModelResp` harus meneruskan field hasil subtract tanpa transformasi tambahan.

## Acceptance Criteria

- `GET /sales/v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001` mengembalikan:
  - `qty = 134`
  - `qty_return` sesuai hasil reference SQL Widya.
  - `net_sales_return` dan `return_rate` sesuai reference SQL Widya.
  - `total_discount_promo = 1238740` tetap dijaga dari scope awal (format UI: `1.238.740`) kecuali QA terbaru memisahkan ke tiket lain.
- Test SQL regression:
  - `os.qty - rs.qty_return` ada, `os.qty AS qty` TIDAK ada.
  - `os.discount_promo - rs.discount_promo` ada, pola `+` untuk `total_discount_promo` TIDAK ada.
- Test integrasi (sqlmock atau DB test) tervalidasi terhadap fixture yang menghasilkan nilai expected di atas.
- Test service memverifikasi `data.Qty` dan `data.TotalDiscountPromo` sama dengan output repository (tanpa overwrite).
- `rtk go test ./...` lulus di module `sales`.

## Existing Patterns / Reuse

- Reuse `newReportRepoDryRunDB`, `latestRecordedQuery`, `assertSecondarySalesSummaryDateVars` di `sales/repository/report_repository_test.go` untuk SQL-shape test.
- Reuse `mockReportRepositoryForService` di `sales/service/report_service_test.go` untuk service mapping test.
- Reuse `valueOrZero` helper di test yang sama untuk null-safety assertion.
- Reuse `SecondarySalesReportDashboardSumPayload` entity sebagai sumber parameter filter.

## Constraints

- Tetap di module `sales/`; jangan sentuh service lain.
- `rtk` prefix untuk semua shell command di sesi ini.
- Jangan commit `.env`, secret, atau dump DB.
- Planner TIDAK mengedit source/test files; hanya plan + evidence.

## Risks

- Branch aktif: `gate` (per `git status`). Implementasi akan commit ke branch kerja `@fixer`.
- Jika DB lokal tidak punya fixture `C260020001` untuk `2026-06`, integrasi test gagal karena data kosong. Solusi: gunakan sqlmock atau seed terisolasi.
- Staging environment butuh auth user; validasi manual harus pakai akun user, bukan token Jira.
- Test aritmatika existing (`TestSX2258SecondarySalesReportSummaryArithmeticRegression`) tidak mengeksekusi SQL; perlu test executed-SQL untuk bukti kuat.

## Decisions / Assumptions

- Asumsi awal (revisi): subtract arithmetic di final select SUDAH benar (commit `4ebacfe`), tapi filter di return branch (`r.data_status = 6` dan date semantics) belum match referensi Widya. Plan ini merevisi asumsi: deliverable BUKAN hanya `regression-proof`, tapi juga code change untuk alignment filter.
- Asumsi: nama field BE `qty` = `Number of Product Sold` di UI; `total_discount_promo` = `Discount and Promo` di UI (mapping via service `data.Qty` dan `data.TotalDiscountPromo`).
- Asumsi: `r.data_status = 6` di return branch akan dihapus atau dibuktikan tidak relevan; bukti di evidence sebelum eksekusi.
- Asumsi: date semantics akan disesuaikan ke `BETWEEN` atau ke `>= AND <` setelah validasi angka `134` tercapai; keputusan ditulis di evidence.
- Asumsi: FE tidak butuh perubahan karena field response sudah benar; yang berubah hanya angka internal BE.
- Decision gate answered: token bearer yang bocor tidak akan di-rotate; risiko diterima user. Eksekutor tetap DILARANG menulis token itu ke file/log/evidence/commit.
- Decision gate answered: `r.data_status = 6` boleh dihapus mengikuti reference Widya.
- Decision gate answered: date semantics diprobe dulu di local DB; pilih `BETWEEN` atau `>= AND <` berdasarkan hasil yang menghasilkan `qty=134`.
- Decision gate answered: scope validasi mencakup semua field return dari reference Widya (`qty`, `qty_return`, `net_sales_return`, `return_rate`) plus `total_discount_promo` dari scope awal.
- Decision gate answered: 4-query probe dijalankan di local DB.
- Open question (non-blocking, slice-safe): apakah FE akan consume `total_discount_promo` langsung atau ada normalizer lokal? (Di luar scope BE; tidak ada perubahan yang direncanakan di BE.)

## Execution Source of Truth

Urutan prioritas saat eksekusi:

1. Petunjuk eksplisit user terbaru.
2. Aturan keamanan/rahasia repo.
3. Non-negotiable Implementation Invariants.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria + Done Criteria.
6. Implementation Steps.
7. Follow-up/rekomendasi.

Jika ada konflik antara dua sumber di atas, eksekutor wajib mengikuti sumber dengan prioritas lebih tinggi dan mencatat konflik di evidence verifikasi.

## Non-negotiable Implementation Invariants

- `qty` di summary = `os.qty - rs.qty_return` (subtract). TIDAK boleh order-only atau plus.
- `total_discount_promo` di summary = `os.discount_promo - rs.discount_promo` (subtract). TIDAK boleh plus.
- Filter summary sesuai referensi Widya (subject T0 konfirmasi):
  - `o.cust_id IN ?` order, `rd.cust_id IN ?` return.
  - `o.data_status IN (6, 7)` di kedua branch (return ikut via join order).
  - `o.invoice_date` (bukan `r.return_date`).
  - `r.data_status = 6` HANYA dipakai jika data membuktikan harus ada; defaultnya dihapus mengikuti referensi.
  - Date semantics: `BETWEEN` atau `>= AND <` dipilih yang menghasilkan `qty=134`; keputusan didokumentasikan.
- `outlet_ids`, `salesman_ids`, `pro_ids` harus konsisten untuk order dan return; produk return pakai `rd.product_id`, order pakai `od.pro_id`.
- Service `SecondarySalesReportSumReportByMonth` TIDAK menambahkan/mengubah nilai `Qty` atau `TotalDiscountPromo` dari hasil repository (mapping 1:1).
- Response field JSON `qty` dan `total_discount_promo` tidak diubah namanya.
- TIDAK ADA JWT, token, kredensial, atau bearer apapun yang ditulis ke file, log, evidence, atau commit. Token staging yang bocor di chat harus di-rotasi sebelum eksekusi.

## Do Not / Reject If

- Jangan ganti `qty` atau `total_discount_promo` jadi penjumlahan (`+`) antara order dan return.
- Jangan pakai `r.return_date` sebagai date filter summary (QA acuan pakai `o.invoice_date`).
- Jangan pakai `od.product_id` (kolom tidak ada; kolom order produk = `od.pro_id`).
- Jangan pertahankan `r.data_status = 6` kalau reference-query validation menunjukkan nilai `134` hanya tercapai tanpa filter itu.
- Jangan tulis ulang query dari nol; reuse `buildSecondarySalesReportSummarySQL` dan helper filter.
- Jangan commit `.env`/secret/Jira token/bearer token dari chat.
- Jangan paste token bocor ke evidence. Redact sebagai `<REDACTED_BEARER_TOKEN>`.
- Jangan skip TDD step; test Red harus ditulis sebelum perubahan apapun (lihat TDD Plan).
- Jangan rebase/replace branch `gate`/`sx2258-gate`/`quality-gate`/`master-lookup-export`; pakai working tree normal.
- Jangan claim "fix selesai" hanya karena code compile; harus lewat `rtk go test ./...` di module `sales` dan validasi manual endpoint.

## Diff Boundary

Allowed file groups:

- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`
- `sales/service/report_service.go`
- `sales/service/report_service_test.go`
- `sales/controller/report_controller.go`
- `sales/controller/so_controller_test.go`
- `.opencode/evidence/20260618-1002-sx-2258-secondary-sales-summary/**`
- `.opencode/plans/20260618-1002-sx-2258-secondary-sales-summary.md` (sudah ditulis oleh planner)

Generated report exception: tidak ada generated file di luar source/test.

Evidence paths: `.opencode/evidence/20260618-1002-sx-2258-secondary-sales-summary/discovery.md` (sudah ada) + execution log nanti di folder yang sama.

Out-of-boundary: perubahan module lain (`inventory`, `master`, `finance`, `tms`, `pjp`, dll), perubahan `go.mod`, perubahan file compose/docker, perubahan migration. Jika ternyata dibutuhkan, revert atau justifikasi di evidence sebelum `@quality-gate`.

## TDD / Test Plan

- TDD: WAJIB (logika summary, regression SX-2258).
- Reason: bug sebelumnya halus (perubahan `+` ke `-` saja); test yang ada tidak cukup mengeksekusi SQL. Butuh executed-SQL test untuk mengunci.
- Existing patterns: dry-run GORM + callback untuk capture SQL; mock repository untuk service.
- First failing test (Red):
  - Tulis `TestSecondarySalesReportSumReportByMonthExecutesSubtractArithmetic` di `sales/repository/report_repository_test.go`.
  - Pakai `sqlmock` (preferred) atau stub `*gorm.DB` agar `Take(&data)` mengembalikan baris hasil subtract.
  - Assertion: `data.Qty == 134`, `data.TotalDiscountPromo == 1238740`, `data.QtyReturn == 16`, dan SQL string mengandung `os.qty - rs.qty_return` + `os.discount_promo - rs.discount_promo`.
- Green step: tidak perlu code change karena query sudah subtract, tapi test akan gagal sampai fixture/return row benar; jika gagal karena mock salah, perbaiki mock dulu.
- Refactor step: kalau ada duplikasi builder atau helper filter, refactor tanpa ubah SQL output. Jalankan semua test existing untuk pastikan tidak ada regresi.
- Edge cases (ditambah ke test):
  - Order summary kosong: `qty = 0 - qty_return`, `total_discount_promo = 0 - rs.discount_promo`.
  - Return summary kosong: hasil = order summary (tidak minus).
  - `NULL` di `promo_final2` dst: kontribusi `0` (lihat `valueOrZero`).
  - `os.net_sales_inc_ppn = 0`: return rate = 0 (bukan division by zero).
- Commands: `cd sales && rtk go test ./... -run SecondarySalesReportSumReportByMonth` dan `cd sales && rtk go test ./...`.

## Implementation Steps

0. Security preflight: minta rotasi token bearer yang sudah bocor di chat. Jangan gunakan token itu untuk test, log, atau evidence.
1. Konfirmasi current code sudah subtract (sudah terverifikasi via diff `4ebacfe`), lalu validasi kenapa response masih salah: bandingkan current SQL dengan reference Widya.
2. Jalankan 4 query pembanding di DB aman/local/staging (tanpa menyimpan token):
   - current repo filter: `r.data_status = 6`, date `>= AND <`.
   - hapus `r.data_status = 6`, date `>= AND <`.
   - current repo filter: `r.data_status = 6`, date `BETWEEN`.
   - hapus `r.data_status = 6`, date `BETWEEN`.
   Catat varian mana yang menghasilkan `qty=134` dan `qty_return` sesuai Widya.
3. Ubah `buildSecondarySalesReportSummarySQL` supaya return branch match reference Widya: default perubahan yang diantisipasi adalah hapus `r.data_status = 6`; date semantics ikut hasil step 2.
4. Tambah regression test SQL-shape: reject `r.data_status = 6` jika keputusan step 2 menghapusnya; assert date clause baru; assert `qty=order-return` tetap ada.
5. Tambah executed-SQL regression test dengan `sqlmock` untuk `SecondarySalesReportSumReportByMonth` agar angka `134` dan `1238740` tervalidasi.
6. Tambah/ekstensi test service untuk mapping `Qty` dan `TotalDiscountPromo` dari repository (mock return nilai subtract).
7. Tambah test controller yang memastikan `SecondaryReportSalesSumMonth` mengembalikan response dengan `qty=134`, `total_discount_promo=1238740` saat service mock memberi nilai itu.
8. Jalankan `cd sales && rtk go test ./...` dan pastikan semua test lulus.
9. Validasi manual endpoint di environment local/staging dengan auth user baru yang sudah aman:
   `curl -H "Authorization: Bearer $USER_TOKEN" "$BASE_URL/sales/v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001"`
10. Verifikasi response: `qty=134`, `total_discount_promo=1238740` jika masih in-scope. Catat output redacted di evidence.

## Expected Files to Change

- `sales/repository/report_repository.go` (kemungkinan besar: hapus `r.data_status = 6` di return branch, dan/atau ubah date clause jadi `BETWEEN`).
- `sales/repository/report_repository_test.go` (tambah executed-SQL regression test; reject pola filter lama jika dihapus).
- `sales/service/report_service_test.go` (tambah assertion mapping subtract).
- `sales/controller/so_controller_test.go` (tambah assertion response subtract).
- `.opencode/evidence/20260618-1002-sx-2258-secondary-sales-summary/execution.md` (log eksekusi @fixer).

## Agent / Tool Routing

- `@fixer`: implementasikan TDD Red → Green → Refactor. Step awal: jalankan 4 query pembanding (lihat Implementation Step 2) untuk konfirmasi filter mana yang harus dihapus/ditambah; update `buildSecondarySalesReportSummarySQL` sesuai hasil; tambah/update test.
- `@quality-gate`: final signoff conformance + security/secret check. WAJIB cek bahwa evidence tidak mengandung token bocor.
- `@orchestrator`: delegasi eksekusi, kumpulkan evidence, jalankan validasi.

## Executor Handoff Prompt (copyable)

```
Task: SX-2258 Secondary Sales Summary (Order - Return)
Mode: Maintenance Stability
Read first: .opencode/plans/20260618-1002-sx-2258-secondary-sales-summary.md
            .opencode/evidence/20260618-1002-sx-2258-secondary-sales-summary/discovery.md

Scope (must_preserve):
- qty = os.qty - rs.qty_return (subtract, not plus)
- total_discount_promo = os.discount_promo - rs.discount_promo (subtract, not plus)
- date filter: o.invoice_date (not r.return_date)
- product filter: od.pro_id (order) / rd.product_id (return)
- data_status: o.data_status IN (6,7); r.data_status = 6 harus dibuktikan atau dihapus karena referensi Widya tidak memakainya
- compare date semantics: repo >= AND < vs reference BETWEEN; pilih yang menghasilkan qty=134
- service mapping Qty/TotalDiscountPromo 1:1 dari repository
- response field name qty dan total_discount_promo tidak berubah
- never write leaked bearer token to any file/log/evidence/commit

Do not touch:
- module lain di luar sales/
- sales/go.mod, sales/go.sum
- file compose/docker/migrations
- .env, secret, Jira token

Validation:
- cd sales && rtk go test ./... (semua test harus PASS)
- rtk go test ./repository -run TestSecondarySalesReportSumReportByMonth
- rtk go test ./service -run TestSecondarySalesReportSumReportByMonth
- rtk go test ./controller -run TestSecondaryReportSalesSumMonth
- manual: curl ke /sales/v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001
  harus mengembalikan qty=134 dan total_discount_promo=1238740

Return/evidence:
- ringkasan perubahan (file + jumlah baris)
- output rtk go test
- output curl + response
- simpan di .opencode/evidence/20260618-1002-sx-2258-secondary-sales-summary/execution.md

Claim limit: jangan claim DONE sebelum rtk go test ./... PASS dan
response curl menunjukkan qty=134 + total_discount_promo=1238740.
```

## Execution-ready Worklist / Handoff Contract

Task ID: T0
- Action: Security preflight. User risk-accept token bearer bocor (tidak di-rotate). Eksekutor WAJIB menghindari copy/paste token ke file/log/evidence/commit. Saat dokumentasi curl, gunakan placeholder `$USER_TOKEN` dan redact apapun yang mengandung `eyJ...`.
- Depends on: none
- Owner/lane: `@orchestrator`
- Validation: grep terhadap evidence dan source tree: `grep -RIn "eyJ" sales/ .opencode/evidence/20260618-1002-sx-2258-secondary-sales-summary/` harusnya tidak menemukan token utuh. Placeholder `$USER_TOKEN` atau `<REDACTED_BEARER_TOKEN>` saja yang muncul.
- Exit criteria: tidak ada token utuh di repo.
- Blocking: advisory (tidak memblokir eksekusi; risiko sudah di-accept user).
- requires_user_decision: no (keputusan sudah dijawab di decision gate).
- must_preserve: tidak ada token tertulis di file/log.
- do_not_touch: token bearer.
- evidence_update: catat status risk-accept di `execution.md` tanpa nilai token.
- exit_verification: grep secret tidak menemukan token utuh.
- start_with: yes (mulai dari sini)

Task ID: T1
- Action: Audit diff `sales/repository/report_repository.go` vs reference Widya. Tentukan tiga hal: (a) apakah `r.data_status = 6` harus dihapus; (b) date semantics `>= AND <` atau `BETWEEN`; (c) subtotal/return fields lain yang ikut referensi (Number of Product Return, Return Value, Return Rate %).
- Depends on: T0
- Owner/lane: `@fixer` (read-only)
- Validation: grep + diff visual.
- Exit criteria: tiga keputusan tercatat.
- Blocking: ready (setelah T0 done)
- requires_user_decision: no (decisions bisa default ke reference Widya; escalate hanya jika ada konflik dengan return-status domain rule)
- must_preserve: invariants subtract.
- do_not_touch: source code.
- evidence_update: tambahkan catatan audit di `discovery.md` atau `execution.md`.
- exit_verification: keputusan T1(a/b/c) tercatat.

Task ID: T2
- Action: Jalankan 4 query pembanding di DB environment aman (local/staging-readonly). Query variants: (i) current filter; (ii) hapus `r.data_status = 6`; (iii) date `BETWEEN`; (iv) hapus `r.data_status = 6` + date `BETWEEN`. Catat `qty`, `qty_return`, `net_sales_return`, `return_rate` per variant. Pilih variant yang menghasilkan `qty=134` sesuai QA.
- Depends on: T1
- Owner/lane: `@fixer`
- Validation: output 4 variant; salah satu menghasilkan `qty=134`.
- Exit criteria: variant target terdokumentasi.
- Blocking: blocked until T0 done
- requires_user_decision: no
- must_preserve: read-only di DB production.
- do_not_touch: production data.
- evidence_update: simpan output di `execution.md` (redacted).
- exit_verification: variant target + nilai tercatat.

Task ID: T3
- Action: Tulis test Red `TestSecondarySalesReportSumReportByMonthReturnsSubtractValues` di `sales/repository/report_repository_test.go` menggunakan `sqlmock`. Stub row hasil variant target dari T2. Assertion: `Qty=134`, `QtyReturn`, `NetSalesReturn`, `ReturnRate` sesuai reference Widya.
- Depends on: T2
- Owner/lane: `@fixer`
- Validation: `cd sales && rtk go test ./repository -run TestSecondarySalesReportSumReportByMonthReturnsSubtractValues`.
- Exit criteria: test FAIL (Red) karena SQL current tidak match variant target. Capture failure log.
- Blocking: ready (setelah T2 done)
- requires_user_decision: no
- must_preserve: tidak mengubah source repository.
- do_not_touch: source repository code.
- evidence_update: simpan log Red di `execution.md`.
- exit_verification: Red test failure captured.

Task ID: T4
- Action: Update `buildSecondarySalesReportSummarySQL` di `sales/repository/report_repository.go` sesuai variant target dari T2 (kemungkinan besar: hapus `r.data_status = 6`; ubah date clause ke `BETWEEN`; tetap subtract arithmetic). Lalu update existing test SQL-shape agar tidak reject pola baru.
- Depends on: T3
- Owner/lane: `@fixer`
- Validation: `cd sales && rtk go test ./repository -run TestSecondarySalesReportSumReportByMonth -v`. Test T3 harus PASS (Green).
- Exit criteria: T3 PASS, existing test tidak regress.
- Blocking: ready
- requires_user_decision: no
- must_preserve: subtract arithmetic + response field name.
- do_not_touch: module lain.
- evidence_update: simpan log Green di `execution.md`.
- exit_verification: PASS log.

Task ID: T5
- Action: Tambah/ekstensi test service `TestSecondarySalesReportSumReportByMonthPropagatesSubtractFromRepository` di `sales/service/report_service_test.go` agar mapping `data.Qty` dan `data.TotalDiscountPromo` dari repository terjaga.
- Depends on: T4
- Owner/lane: `@fixer`
- Validation: `cd sales && rtk go test ./service -run TestSecondarySalesReportSumReportByMonthPropagatesSubtractFromRepository`.
- Exit criteria: test PASS.
- Blocking: ready
- requires_user_decision: no
- must_preserve: signature service dan entity response.
- do_not_touch: source service.
- evidence_update: log PASS.
- exit_verification: PASS log.

Task ID: T6
- Action: Tambah/ekstensi test controller `TestSecondaryReportSalesSumMonthReturnsSubtractValuesFromService` di `sales/controller/so_controller_test.go` untuk assertion JSON `qty=134` dan field summary lain yang relevan.
- Depends on: T5
- Owner/lane: `@fixer`
- Validation: `cd sales && rtk go test ./controller -run TestSecondaryReportSalesSumMonthReturnsSubtractValuesFromService`.
- Exit criteria: test PASS.
- Blocking: ready
- requires_user_decision: no
- must_preserve: route, response builder.
- do_not_touch: source controller.
- evidence_update: log PASS.
- exit_verification: PASS log.

Task ID: T7
- Action: Jalankan full test module `sales` dan pastikan semua PASS.
- Depends on: T4, T5, T6
- Owner/lane: `@fixer`
- Validation: `cd sales && rtk go test ./...`.
- Exit criteria: zero failures, zero unexpected skips.
- Blocking: ready
- requires_user_decision: no
- must_preserve: tidak ada test existing yang dihapus/di-skip.
- do_not_touch: file di luar module sales.
- evidence_update: log ringkasan di `execution.md`.
- exit_verification: PASS log + summary counts.

Task ID: T8
- Action: Validasi manual endpoint di environment local/staging dengan token baru yang sudah aman (bukan token bocor). User Phill Jones atau Yogie menjalankan curl manual; evidence disimpan redacted.
- Depends on: T7, T0
- Owner/lane: `@fixer` (operator)
- Validation: `curl -H "Authorization: Bearer $USER_TOKEN" "$BASE_URL/sales/v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001"`.
- Exit criteria: response `qty=134` dan field summary lain sesuai variant target.
- Blocking: blocked (butuh env + token baru)
- requires_user_decision: yes
- must_preserve: endpoint contract.
- do_not_touch: response shape.
- evidence_update: simpan response redacted (token + PII) di `execution.md`.
- exit_verification: response JSON valid + nilai match.

Task ID: T9
- Action: Submit ke `@quality-gate` untuk final signoff. WAJIB cek evidence tidak mengandung token bocor/PII.
- Depends on: T8
- Owner/lane: `@orchestrator` → `@quality-gate`
- Validation: quality gate review evidence + secret scan.
- Exit criteria: gate verdict PASS atau PASS_FOR_SLICE.
- Blocking: blocked until T8 complete
- requires_user_decision: no
- must_preserve: invariants.
- do_not_touch: n/a (review only).
- evidence_update: link ke `execution.md` dan `discovery.md`.
- exit_verification: gate verdict.

## Validation Commands

- `cd sales && rtk go mod download && rtk go mod tidy`
- `cd sales && rtk go test ./...`
- `cd sales && rtk go test ./repository -run TestSecondarySalesReportSumReportByMonth -v`
- `cd sales && rtk go test ./service -run TestSecondarySalesReportSumReportByMonth -v`
- `cd sales && rtk go test ./controller -run TestSecondaryReportSalesSumMonth -v`
- Manual: `curl -i -H "Authorization: Bearer $USER_TOKEN" "$BASE_URL/sales/v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001"`

## Evidence Requirements

- Source: `.opencode/evidence/20260618-1002-sx-2258-secondary-sales-summary/discovery.md` (sudah ditulis).
- Execution log: `.opencode/evidence/20260618-1002-sx-2258-secondary-sales-summary/execution.md` (akan diisi @fixer).
- Test output: PASS/FAIL log per task.
- Manual endpoint response: simpan JSON response (redacted) di `execution.md`.

Source strategy (ringkas): repo evidence (sudah dipakai), git history commit `4ebacfe`, file test existing. Tidak ada internet research; library behavior (gorm + sqlmock) diasumsikan stable. Jika ternyata butuh library docs, route ke `@librarian` (di luar scope plan ini).

## Done Criteria

- Semua task T1–T6 PASS/selesai.
- T7 quality gate verdict PASS atau PASS_FOR_SLICE.
- Acceptance criteria numerik (`qty=134`, `total_discount_promo=1238740`) tervalidasi oleh test executed-SQL + validasi manual.
- Tidak ada file out-of-boundary berubah tanpa justifikasi di evidence.

## Final Planning Summary

- Artifacts created/kept:
  - `.opencode/plans/20260618-1002-sx-2258-secondary-sales-summary.md` (primary plan, source of truth, direvisi setelah feedback Widya).
  - `.opencode/evidence/20260618-1002-sx-2258-secondary-sales-summary/discovery.md` (discovery evidence; perlu append audit T1 di iterasi berikutnya).
  - `.opencode/draft/20260618-1002-sx-2258-secondary-sales-summary/` (folder tetap untuk konsistensi task id; tidak ada draft yang disimpan).
- Key decisions (revisi setelah decision gate):
  - Decision 1: token bearer bocor tidak di-rotate, risiko diterima user. Tetap wajib: tidak ada nilai token di file/log/evidence/commit.
  - Decision 2: `r.data_status = 6` dihapus mengikuti reference Widya.
  - Decision 3: date semantics `>= AND <` vs `BETWEEN` diputuskan via 4-query probe di local DB; pilih varian yang menghasilkan `qty=134`.
  - Decision 4: scope validasi mencakup semua field return (`qty`, `qty_return`, `net_sales_return`, `return_rate`) plus `total_discount_promo` dari scope awal.
  - Decision 5: probe dijalankan di local DB; fixture `C260020001` Juni 2026 harus tersedia atau di-seed; jika tidak tersedia, probe dialihkan ke staging readonly (tetap risk-accepted untuk token).
  - Code change ke `sales/repository/report_repository.go` diantisipasi (hapus `r.data_status = 6`, dan/atau ubah date clause ke `BETWEEN`).
- Assumptions:
  - FE mapping `qty` → "Number of Product Sold" dan `total_discount_promo` → "Discount and Promo" tidak berubah.
  - `total_discount_promo = 1238740` masih in-scope sampai QA mengkonfirmasi pemisahan tiket.
  - Reference Widya adalah sumber kebenaran final untuk filter summary, menggantikan versi sebelumnya.
  - Token bearer yang ter-paste di chat tidak akan dipakai oleh eksekutor; jika dipakai manual, eksekutor WAJIB replace dengan placeholder.
- Open questions (sudah terjawab di decision gate; tetap dicatat untuk audit):
  - `r.data_status = 6` => HAPUS.
  - Date semantics => PUTUSKAN dari probe local.
  - Scope field => SEMUA field return + `total_discount_promo`.
  - Lokasi probe => LOCAL DB (fallback staging-readonly).
- Readiness: `ready-for-implementation` (semua blocker decision sudah terjawab). T0 (security preflight) downgrade dari blocking menjadi advisory karena user risk-accept.
- Cleanup: tidak ada draft/evidence yang perlu dihapus saat finalisasi (keduanya masih relevan untuk execution dan audit).
