# Goal

Memperbaiki endpoint `POST /sales/v1/reports/secondary-sales` pada modul `sales/` agar report/export Secondary Sales menampilkan transaksi Canvas yang sudah ter-invoice, tetap mempertahankan flow normal TO hingga invoicing, menjaga total nilai sesuai ekspektasi QA, dan menghindari duplicate row/double counting.

# Non-goals

- Tidak mengubah modul `pjp-sales/`.
- Tidak mengubah report lain di luar Secondary Sales.
- Tidak mengubah mekanisme RMQ/export async yang sudah ada.
- Tidak mengubah rule bisnis return `PPN` di luar asumsi existing tanpa bukti bisnis baru.
- Tidak mengubah timezone helper global kecuali benar-benar dibutuhkan oleh bug ini.

# Scope

- Endpoint scope: `POST /v1/reports/secondary-sales` di `sales/`.
- Layer scope:
  - `sales/controller/report_controller.go` hanya untuk verifikasi alur input/auth context bila perlu
  - `sales/service/report_service.go` untuk memastikan export tetap memakai source yang sama
  - `sales/repository/report_repository.go` sebagai titik utama perbaikan query
  - `sales/service/report_service_test.go` dan/atau test repository baru untuk regression coverage
- Evidence scope:
  - validasi mapping timestamp QA
  - validasi reuse query list/export
  - validasi risiko dedup parent/child product fallback

# Requirements

1. Report harus mengambil order Canvas yang sudah sampai invoice.
2. Report tetap harus mengambil flow existing TO hingga invoicing.
3. Filter existing harus tetap bekerja:
   - `distributor_ids`
   - `outlet_ids`
   - `salesman_ids`
   - `pro_ids`
   - `from` / `to`
4. Export Excel dan query source report harus tetap konsisten.
5. Query harus parameterized, tanpa hardcode `cust_id`, `parent_cust_id`, invoice, atau nomor dokumen contoh.
6. Harus ada regression coverage minimal level unit/query-level.

# Acceptance Criteria

1. Request dengan filter QA menampilkan invoice Canvas `INV2605080009` bila data staging/fixture tersedia.
2. Total untuk case QA sesuai:
   - `GrossSales = 122.100`
   - `NetSalesExcPPN = 110.000`
   - `PPN = 12.100`
   - `NetSalesIncPPN = 122.100`
3. Transaksi non-Canvas tetap muncul.
4. Tidak ada duplicate row akibat fallback join produk parent/child.
5. Filter outlet, salesman, distributor, dan product tetap aktif.
6. Export workbook tetap membaca data dari source query yang sama dengan report path.
7. Ada test/regression check yang menangkap risiko Canvas + dedup + filter.

# Existing Patterns/Reuse

- **Reuse utama:** `buildSecondarySalesUnionQuery(dataFilter, withPagination)` di `sales/repository/report_repository.go` sudah menjadi single source of truth untuk query export dan paginated/report variant.
- `SecondarySalesUnionPagination(...)` dan `SecondarySalesUnion(...)` sama-sama memanggil builder tersebut; ini harus dipertahankan.
- Auth context sudah diinject oleh controller lewat:
  - `request.CustID = c.Locals("cust_id")`
  - `request.ParentCustID = c.Locals("parent_cust_id")`
- Export workbook di `SubscribeSecondarySalesReport(...)` sudah membaca `ReportRepository.SecondarySalesUnion(...)`; jangan buat query export terpisah.
- Timestamp helper existing yang dipakai lintas repo adalah `str.UnixTimestampToUtcTime(...)`; QA timestamp yang diberikan sudah cocok dengan hari lokal yang dimaksud, jadi helper ini bisa direuse.
- Test export existing di `sales/service/report_service_test.go` sudah membuktikan workbook memakai field hasil query; cukup diperluas, bukan diganti total.

# Constraints

- Scope implementasi: `sales/` saja.
- Wajib mengikuti arsitektur repo: Controller -> Service -> Repository -> DB.
- Jangan menyimpan atau menyalin kredensial Jira/staging ke repo.
- Query harus tetap aman terhadap `NULL` via `COALESCE` atau ekuivalen saat perlu.
- Jangan memecah source list vs export menjadi query berbeda.
- Dalam planning ini, tidak ada edit implementasi di luar `.opencode/`.

# Risks

1. **Duplicate join risk**
   - Join `mst.m_product pp` dengan kondisi `pro_id OR pro_code` dapat menggandakan baris.
2. **False fix risk**
   - Canvas muncul, tetapi total menjadi salah karena fallback parent product menghasilkan multiplikasi row.
3. **Date-boundary regression**
   - Perubahan konversi tanggal dapat menggeser hasil lintas hari.
4. **Return-sign ambiguity**
   - Query referensi FE mempertahankan `PPN` return positif; jika diubah tanpa validasi bisnis, hasil historis bisa berubah.
5. **Filter regression**
   - Perubahan struktur SQL/CTE bisa memutus filter `distributor_ids` atau `pro_ids` bila alias tidak konsisten.

# Decisions/Assumptions

## Decisions
- Modul target: `sales/` saja.
- Level regression minimal: **unit/query-level**.
- Export dan report harus tetap menggunakan helper query yang sama.
- Sumber utama fix berada di repository query builder, bukan di controller atau format Excel.

## Assumptions / Open Questions
- **Assumption:** Unix timestamp QA (`1778086800` s/d `1778173199`) memang dimaksudkan untuk tanggal lokal `2026-05-07` Asia/Jakarta, dan helper UTC existing sudah memetakan boundary itu dengan benar.
- **Assumption:** Behavior existing untuk return `PPN` tetap dipertahankan sampai ada bukti bisnis yang mengharuskan sign negatif.
- **Open question for implementation validation:** apakah dataset staging yang disebut Jira masih tersedia dan dapat diakses saat verifikasi manual/integration (`SO2605080006`, `PO260508EMP00210006`, `INV2605080009`)?

# TDD/Test Plan

## TDD requirement
TDD/regression **wajib** karena ini bug report production logic yang memengaruhi hasil numerik export dan response report.

## Reason
Bug berada pada perilaku query dan potensi double counting. Tanpa test regression, perubahan kecil pada join/fallback/filter bisa mengembalikan bug yang sama.

## Existing test patterns
- `sales/service/report_service_test.go`
  - sudah ada pattern mock repository untuk export workbook
  - sudah ada assertion terhadap cell workbook hasil query
- Belum ada coverage kuat untuk `buildSecondarySalesUnionQuery(...)`; ini perlu ditambahkan.

## First failing/regression test
Tambahkan test query-level yang memverifikasi builder SQL:

1. Query order branch harus mengandung:
   - `FROM sls.order_detail od`
   - `JOIN/LEFT JOIN sls."order" o ON o.ro_no = od.ro_no`
   - `o.data_status IN (6,7)`
   - `o.invoice_date BETWEEN ? AND ?`
   - `o.invoice_no IS NOT NULL`
2. Query return branch harus mengandung:
   - `FROM sls.return_det rd`
   - `JOIN/LEFT JOIN sls."return" r`
   - `JOIN/LEFT JOIN sls."order" o ON o.invoice_no = r.invoice_no`
3. Query fallback product **harus diubah menjadi deterministik** dan test harus memastikan tidak lagi memakai join raw yang rawan duplicate tanpa pembatas tunggal.
4. Test params harus memastikan urutan filter tetap benar untuk `cust_id`, `parent_cust_id`, date range, dan filter ids.

## Green step
- Refactor `buildSecondarySalesUnionQuery(...)` untuk memakai CTE/derived table yang memilih maksimal satu parent-product fallback row per transaksi detail.
- Pastikan semua filter tetap berada pada alias yang benar.
- Pastikan `SecondarySalesUnion(...)` dan `SecondarySalesUnionPagination(...)` masih bergantung pada helper yang sama.

## Refactor step
- Extract helper SQL kecil bila perlu, misalnya helper fallback product join yang deterministik.
- Rapikan alias dan `COALESCE` agar order/return branch simetris.
- Tambahkan test service/export yang memastikan workbook tetap memakai hasil query source yang sama setelah refactor.

## Edge cases
- Parent product tidak ada, child product ada.
- Child product ada tetapi metadata tidak lengkap; fallback parent harus mengisi `pro_code`, `pro_name`, unit, conversion, supplier.
- Parent product memiliki lebih dari satu kandidat by `pro_code`.
- Filter `pro_ids` kosong vs berisi.
- Transaksi promo (`item_type != 1`) harus tetap `net_sales_* = 0` seperti existing/query reference.
- Return dalam range invoice yang sama tidak boleh double count.

## Commands
- Dari repo root, verifikasi service:
  - `rtk docker compose -f docker-compose.yml ps`
- Dari direktori modul `sales/`:
  - `rtk go test ./...`
  - atau minimal targeted:
  - `rtk go test ./service -run SecondarySales`
  - `rtk go test ./repository -run SecondarySales`

# Implementation Steps

1. **Trace current runtime path**
   - Pastikan endpoint `POST /v1/reports/secondary-sales` tetap berakhir di `ReportRepository.SecondarySalesUnion(...)` untuk export.

2. **Harden query builder**
   - Fokus di `buildSecondarySalesUnionQuery(...)`.
   - Pertahankan source order dari `sls.order_detail` + `sls."order"` dengan filter:
     - `cust_id`
     - `o.data_status in (6,7)`
     - `o.invoice_date between from/to`
     - `o.invoice_no is not null`
   - Pertahankan source return dari `sls.return_det` + `sls."return"` + `sls."order"` invoice.

3. **Implement deterministic product fallback**
   - Ganti join parent product yang sekarang rawan duplicate.
   - Prioritas yang disarankan:
     1. match `pp.pro_id = t.product_id`
     2. fallback `pp.pro_code = cp.pro_code` bila child code tersedia
   - Gunakan subquery/CTE/lateral/`DISTINCT ON` yang menjamin hanya satu row parent product terpilih per transaksi detail.

4. **Preserve numeric/business semantics**
   - Tetap gunakan perhitungan order final fields:
     - `qty*_final`
     - `sell_price*`
     - `amount_final`
     - `vat_value_final`
     - `promo_value_final`
     - `disc_value_final`
   - Return tetap negatif untuk qty/gross/net sesuai existing/reference.
   - Biarkan `PPN` return mengikuti behavior existing kecuali ditemukan bukti bisnis kuat saat implementasi.

5. **Re-verify filters**
   - `distributor_ids` harus filter distributor hasil join customer/distributor.
   - `outlet_ids` harus filter outlet transaksi.
   - `salesman_ids` harus filter salesman transaksi.
   - `pro_ids` harus filter ID produk transaksi, bukan hanya hasil fallback metadata.

6. **Regression tests**
   - Tambah test query-level untuk builder.
   - Tambah/adjust service export test hanya bila perlu memastikan source workbook tidak berubah.

7. **Validation run**
   - Jalankan targeted Go tests.
   - Bila environment staging/internal DB tersedia, lakukan verifikasi manual terhadap invoice contoh dan total QA.

# Expected Files to Change

- `sales/repository/report_repository.go`
- `sales/service/report_service_test.go`
- Kemungkinan file test repository baru, misalnya:
  - `sales/repository/report_repository_test.go`

# Agent/Tool Routing

- **Planner (current):** menyusun artifact plan dan evidence.
- **Implementation agent / engineer berikutnya:** eksekusi perubahan code di repository dan test.
- **Explorer evidence used:** discovery struktur route/service/repository dan risiko duplicate join.
- **Document-specialist evidence used:** pembacaan file `SecondarySales-080526-003 (2).xlsx` untuk memastikan struktur kolom output dan total row.
- **No official docs / GitHub / Brave / browser used:** tidak material untuk bug SQL internal ini.

# Validation Commands

Jalankan dari repo root atau modul sesuai konteks:

- `rtk docker compose -f docker-compose.yml ps`
- `rtk go test ./...`  *(jalankan dari `sales/` bila ingin terbatas modul)*
- `rtk go test ./service -run SecondarySales`
- `rtk go test ./repository -run SecondarySales`

Jika ada akses DB/manual verification yang aman:

- Jalankan query verifikasi untuk invoice `INV2605080009`
- Bandingkan total:
  - `GrossSales`
  - `NetSalesExcPPN`
  - `PPN`
  - `NetSalesIncPPN`

# Evidence Requirements

Implementation nanti dianggap siap direview jika ada evidence berikut:

1. Cuplikan/hasil test query-level yang menunjukkan query builder sudah mengandung source Canvas invoice path yang benar.
2. Bukti tidak ada duplicate row dari fallback product join.
3. Hasil test Go yang lulus untuk modul terkait.
4. Jika akses data tersedia, bukti manual untuk invoice `INV2605080009` dan total expected QA.
5. Catatan eksplisit bila return `PPN` tetap dipertahankan positif sesuai existing behavior.

# Done Criteria

- Query Secondary Sales di `sales/` memakai source transaksi invoice yang mencakup Canvas dan flow normal.
- Fallback product parent/child tidak lagi rawan duplicate row.
- Export workbook tetap menggunakan source query yang sama dengan report path.
- Regression test level unit/query tersedia dan lulus.
- Tidak ada perubahan behavior unrelated report.
- Tidak ada credential/token yang masuk ke repo.

# Final Planning Summary

- **Primary plan path:** `.opencode/plans/20260511-1530-sx-1944-secondary-sales-report.md`
- **Source of truth untuk implementasi:** file plan ini.
- **Artifacts created:**
  - `.opencode/plans/20260511-1530-sx-1944-secondary-sales-report.md`
  - `.opencode/evidence/20260511-1530-sx-1944-secondary-sales-report/discovery.md`
- **Key decisions:**
  - scope `sales/` saja
  - regression minimal unit/query-level
  - pertahankan single query source untuk export dan report
  - fokus fix di repository query builder dengan dedup fallback product yang deterministik
- **Assumptions:**
  - timestamp QA valid untuk boundary tanggal lokal yang dimaksud
  - return `PPN` mengikuti existing behavior sampai ada bukti bisnis lain
- **Open questions remaining:**
  - apakah data staging untuk dokumen contoh masih tersedia saat implementasi/validasi manual
- **Readiness for implementation:** siap diimplementasikan oleh engineer/agent berikutnya
- **Cleanup performed:**
  - tidak membuat draft artifact tambahan karena pertanyaan material sudah dijawab
  - evidence discovery dipertahankan karena masih operasional untuk handoff implementasi
