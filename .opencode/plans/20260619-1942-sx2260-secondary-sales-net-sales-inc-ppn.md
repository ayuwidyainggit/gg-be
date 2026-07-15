# Plan: SX-2260 — Secondary Sales Group net_sales include PPN

## Goal
Endpoint `GET /sales/v1/reports/secondary-sales/group` mengembalikan `net_sales` include PPN untuk semua `group_by` (`outlet`, `salesman`, `product_category`, `product`) sesuai rumus `Nett Sales Inc VAT`.

## Non-goals
- Tidak ubah endpoint report lain (detail, summary, trend, dst).
- Tidak ubah response contract: field `id`, `code`, `name`, `net_sales` tetap.
- Tidak migrasi data atau ubah schema.
- Tidak tambah dependency baru.

## Scope
Hanya fungsi group di `sales/repository/report_repository.go`:
- `buildSecondarySalesReportGroupQuery(groupBy string)`
- `SecondarySalesReportGroupOutlet`
- `SecondarySalesReportGroupSalesman`
- `SecondarySalesReportProductCategory`
- `SecondarySalesReportProduct`

Test terkait:
- `sales/repository/report_repository_test.go` (fragmen SQL + mock DB).
- `sales/service/report_service_test.go` (mock interface).

## Requirements
- Formula `net_sales` per row, sumber `sls.order_detail` dan `sls.return_det`:
  - order: `gross = (qty1_final * sell_price1) + (qty2_final * sell_price2) + (qty3_final * sell_price3)`, `net_inc_ppn = gross - promo_value_final - disc_value_final + vat_value_final`.
  - return: `gross = (qty1 * sell_price1) + (qty2 * sell_price2) + (qty3 * sell_price3)`, `net_inc_ppn = gross - promo_value - disc_value + vat_value` lalu `* -1`.
- Filter tetap: `cust_id IN ?`, `data_status IN (6,7)`, `invoice_date >= dateFrom AND invoice_date < dateTo`.
- Outer select tetap: `COALESCE(SUM(net_sales), 0) AS net_sales`, `GROUP BY id, code, name`, `ORDER BY net_sales DESC`.
- Tipe numeric cocok dengan `model.SecondarySalesReportGroup.NetSales` (`float64`).

## Acceptance Criteria
- AC1: Group query `outlet/salesman/product_category/product` memakai formula include PPN.
- AC2: `net_sales` akhir = `SUM(order_net_inc_ppn) + SUM(return_net_inc_ppn)`.
- AC3: `id`, `code`, `name` source per `group_by` tidak berubah.
- AC4: `ORDER BY net_sales DESC` dan empty result aman.
- AC5: Test repo/service passing, termasuk assertion fragmen SQL untuk `vat_value_final` (order) dan `vat_value` (return).
- AC6: Validasi lokal: `cd sales && rtk go test ./...` hijau.

## Existing Patterns / Reuse
- Query builder `buildSecondarySalesReportGroupQuery` sudah ada; tinggal update dua variabel `orderNetSales` dan `returnNetSales`.
- Mapping dimensi (`outlet`/`salesman`/`product_category`/`product`) di `orderSelect`, `returnSelect`, `orderJoin`, `returnJoin` reuse penuh.
- Pola `* -1` di return branch reuse (multiplier konsep).
- Pola test dry-run `newReportRepoDryRunDB` dan `latestRecordedQuery` reuse untuk assert fragmen SQL.
- Pola service test mock interface `mockReportRepositoryForService` reuse.

## Constraints
- Tidak tambah tabel/kolom baru; gunakan `od.vat_value_final` dan `rd.vat_value`.
- Tidak ubah route/controller/service signature.
- Tidak sentuh endpoint Secondary Sales lain (summary, detail, trend).
- Test tidak boleh hardcode token/credential; pakai mock DB.

## Risks
- Variasi nama kolom harga: existing code pakai `sell_price1/2/3` (bukan `sell_price_final1/2/3`). Pertahankan pola existing; tambahkan klausa `vat_value_final` saja. Jika QA evidence tidak match, pindah ke `sell_price_final1/2/3` adalah langkah lanjut, bukan default.
- Query sls bisa lebih berat dari `report.fact_orders`; gunakan index `cust_id + invoice_date` + `data_status`. Risiko perf tidak diangkat ke user kecuali muncul di QA.
- `ppn_return` harus ikut tanda minus; rumus `* -1` diterapkan ke seluruh ekspresi (gross - promo - disc + ppn) bukan hanya `* -1` di luar.

## Decisions / Assumptions
- Keputusan: pakai source table `sls` (sudah dipakai) + tambah `vat_value_final` / `vat_value`. Tidak migrasi ke `report.fact_orders`.
- Asumsi 1: `od.vat_value_final` dan `rd.vat_value` valid di schema DB.
- Asumsi 2: Pola `sell_price1/2/3` (bukan `sell_price_final*`) benar sesuai existing repo. Reverse jika evidence QA tidak match.
- Asumsi 3: Tidak ada return dengan `item_type` yang perlu filter tambahan untuk PPN; `* -1` cukup.
- Asumsi 4: Tanda PPN order = `+`, PPN return = `-` lewat multiplier `-1` (sesuai doc section 2).
- Open question: apakah prompt doc `sell_price_final1/2/3` punya evidence kuat; di-defer sampai evidence compare QA.

## Execution Source of Truth
Urutan precedence saat implementasi:
1. Instruction user terbaru (prompt SX-2260 + komentar FE).
2. Aturan safety/security/tenant di `.opencode/docs/ARCHITECTURE.md`.
3. Non-negotiable Implementation Invariants (di bawah).
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria & Done Criteria.
6. Implementation Steps & follow-ups.

## Non-negotiable Implementation Invariants
- II-1: Ubah hanya `buildSecondarySalesReportGroupQuery` dan test terkait group. Jangan sentuh fungsi summary/detail/trend/principal.
- II-2: Response contract `id/code/name/net_sales` tidak boleh hilang/tambah field.
- II-3: Return branch `* -1` diterapkan ke seluruh `(gross - promo - disc + ppn)`, bukan hanya di luar ekspresi tanpa PPN.
- II-4: Pertahankan filter `cust_id IN ?`, `data_status IN (6,7)`, dan date range month/year.
- II-5: Assertion test baru harus mengikat fragmen `COALESCE(od.vat_value_final, 0)` (order) dan `COALESCE(rd.vat_value, 0)` (return) di SQL group.
- II-6: Validasi lokal `cd sales && rtk go test ./...` harus hijau.

## Do Not / Reject If
- DN-1: Jangan pakai `SUM(report.fact_orders.net_sales_exclude_ppn)` di endpoint group.
- DN-2: Jangan ganti nama response field, tipe data, atau urutan.
- DN-3: Jangan gabung `* -1` setelah `+ ppn` saja; kalikan seluruh ekspresi.
- DN-4: Jangan hardcode cust_id/month/year dari evidence QA di kode produksi atau test.
- DN-5: Jangan commit tanpa test repo/service yang mengikat fragmen PPN di SQL group.

## Diff Boundary
- Allowed:
  - `sales/repository/report_repository.go` (fungsi group builder & 4 method group).
  - `sales/repository/report_repository_test.go` (assertion fragmen SQL include PPN untuk group).
- Generated/evidence:
  - `.opencode/plans/20260619-1942-sx2260-secondary-sales-net-sales-inc-ppn.md` (this file).
  - `.opencode/evidence/<task-id>/discovery.md` (sudah ada).
  - `.opencode/evidence/<task-id>/index.json` (sudah ada).
- Out-of-bound (harus revert/dijustifikasi): service/controller/model lain, endpoint report non-group, file migrasi, env, compose, README, package/lockfile.

## TDD / Test Plan
- Wajib TDD untuk production logic.
- First failing test (Red):
  - Tambah assertion di `TestSecondarySalesReportGroupQueriesUseSourceTablesAndDateRange` atau test baru yang mengikat:
    - `COALESCE(od.vat_value_final, 0)` muncul di SELECT order branch.
    - `COALESCE(rd.vat_value, 0)` muncul di SELECT return branch dengan `) * -1` di akhir ekspresi.
- Green:
  - Update `orderNetSales` dan `returnNetSales` di `buildSecondarySalesReportGroupQuery` agar menjumlahkan `+ COALESCE(...vat..., 0)` sebelum `* -1` untuk return.
- Refactor:
  - Pastikan `String()` SQL tetap valid, tidak duplikat ekspresi.
- Edge cases:
  - Empty result tetap `[]`/`[]model.SecondarySalesReportGroup`.
  - PPN null/0 -> `COALESCE` jadi 0.
  - Multi-cust: `cust_id IN ?` masih valid.
- Commands:
  - `cd sales && rtk go test ./repository -run TestSecondarySalesReportGroupQueriesUseSourceTablesAndDateRange -v`
  - `cd sales && rtk go test ./...`
  - `cd sales && rtk go build ./...`

## Implementation Steps
1. Tambah assertion Red di test repo untuk fragmen `vat_value_final`/`vat_value` di group SQL. Jalankan test, harus gagal.
2. Edit `buildSecondarySalesReportGroupQuery`:
   - `orderNetSales` tambah `+ COALESCE(od.vat_value_final, 0)` sebelum `) AS net_sales`.
   - `returnNetSales` tambah `+ COALESCE(rd.vat_value, 0)` sebelum `) * -1 AS net_sales`.
3. Jalankan test Red + Green, harus hijau.
4. Jalankan `cd sales && rtk go test ./...` dan `cd sales && rtk go build ./...`.
5. Reuse existing test `TestSecondarySalesReportGroupOutletUsesOutletDisplayMapping`, `TestSecondarySalesReportGroupProductCategoryUsesMasterCategoryFallback`, `TestSecondarySalesReportGroupProductUsesMasterProductFallback` harus tetap hijau.
6. Catat hasil di evidence (test output ringkas, file yang berubah).

## Expected Files to Change
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`

## Agent / Tool Routing
- Owner eksekusi: `@fixer` (bounded implementation + test).
- Read-only input:
  - `@explorer` jika ada tambahan file/pattern yang muncul saat eksekusi.
- Final signoff: `@quality-gate` setelah test hijau.
- Planner tidak eksekusi source edit; handoff via prompt di bawah.

## Executor Handoff Prompt
```
Scope: sales/repository group query saja (SecondarySalesReportGroupOutlet, SecondarySalesReportGroupSalesman, SecondarySalesReportProductCategory, SecondarySalesReportProduct, plus builder buildSecondarySalesReportGroupQuery).
must_preserve:
  - response field id/code/name/net_sales
  - filter cust_id IN ?, data_status IN (6,7), invoice_date >= dateFrom AND < dateTo
  - mapping dimensi per group_by (outlet/salesman/product_category/product)
do_not_touch:
  - service/controller/model lain
  - endpoint Secondary Sales non-group
  - migration/compose/env/README/lockfile
validation:
  - cd sales && rtk go test ./...  (green)
  - cd sales && rtk go build ./... (green)
  - assertion baru mengikat COALESCE(od.vat_value_final, 0) di order branch dan COALESCE(rd.vat_value, 0) di return branch dengan * -1
return: ringkasan diff + output test terakhir.
claim_limits: perbaikan hanya pada endpoint secondary-sales/group sesuai acceptance; tidak klaim menyentuh endpoint lain.
```

## Execution-ready Worklist / Handoff Contract
- T1: add failing assertions
  - action: tambah fragment check `COALESCE(od.vat_value_final, 0)` di order dan `COALESCE(rd.vat_value, 0)` di return pada test group.
  - depends_on: none
  - owner: `@fixer`
  - validation: `cd sales && rtk go test ./repository -run TestSecondarySalesReportGroup -v` (sengaja gagal)
  - exit_criteria: test baru terdaftar, run menunjukkan fail pada fragmen PPN
  - blocking: ready
  - requires_user_decision: no
  - must_preserve: SQL fragmen existing untuk source-table + date range
  - do_not_touch: service test, controller, non-group query
  - evidence_update: tulis ulang `discovery.md` jika ada file test tambahan
  - exit_verification: output test menunjukkan fail ekspektasi
  - start_with: T1
- T2: implement formula
  - action: update `orderNetSales` (+`COALESCE(od.vat_value_final, 0)`) dan `returnNetSales` (+`COALESCE(rd.vat_value, 0)` sebelum `* -1`).
  - depends_on: T1
  - owner: `@fixer`
  - validation: `cd sales && rtk go test ./repository -run TestSecondarySalesReportGroup -v` hijau
  - exit_criteria: semua test group hijau, termasuk assertion PPN
  - blocking: ready
  - requires_user_decision: no
  - must_preserve: response field, filter, dimensi
  - do_not_touch: scope di atas
  - evidence_update: catat log test
  - exit_verification: output test hijau
- T3: full module test + build
  - action: jalankan `rtk go test ./...` dan `rtk go build ./...` dari `sales/`.
  - depends_on: T2
  - owner: `@fixer`
  - validation: exit code 0 di kedua perintah
  - exit_criteria: tidak ada regresi
  - blocking: ready
  - requires_user_decision: no
  - must_preserve: tidak ubah file di luar scope
  - do_not_touch: scope di atas
  - evidence_update: simpan ringkas output di evidence
  - exit_verification: log test/build bersih
- T4: handoff to quality-gate
  - action: rangkum diff + test log untuk reviewer.
  - depends_on: T3
  - owner: `@fixer` → `@quality-gate`
  - validation: ringkasan PR
  - exit_criteria: file diff tersedia, test hijau
  - blocking: ready
  - requires_user_decision: no
  - must_preserve: klaim scope sesuai acceptance
  - do_not_touch: -
  - evidence_update: finalize `discovery.md` jika ada catatan tambahan
  - exit_verification: rangkuman siap

## Validation Commands
- `cd sales && rtk go test ./repository -run TestSecondarySalesReportGroup -v`
- `cd sales && rtk go test ./...`
- `cd sales && rtk go build ./...`

## Evidence Requirements
- Lokasi: `.opencode/evidence/20260619-1942-sx2260-secondary-sales-net-sales-inc-ppn/`.
- File sudah ada: `discovery.md`, `index.json`.
- Tambahan saat eksekusi: log test/build diringkas di `discovery.md` atau file baru `verification.md`.
- Catatan risiko `sell_price_final*` vs `sell_price*` harus muncul di `discovery.md` saat eksekusi bila muncul ketidakcocokan.

## Done Criteria
- Test repo dan service hijau.
- Build hijau.
- Diff hanya di `sales/repository/report_repository.go` + `sales/repository/report_repository_test.go`.
- `buildSecondarySalesReportGroupQuery` mengembalikan SQL yang menyertakan PPN untuk order & return.
- Handoff prompt siap di-copy ke `@orchestrator`/`@fixer`.

## Final Planning Summary
- Artifacts consulted/created:
  - `AGENTS.md`, `.opencode/docs/ARCHITECTURE.md`, `.opencode/docs/QUALITY.md`.
  - `sales/repository/report_repository.go`, `sales/repository/report_repository_test.go`, `sales/service/report_service.go`.
  - Created: `.opencode/plans/20260619-1942-sx2260-secondary-sales-net-sales-inc-ppn.md`.
  - Created: `.opencode/evidence/20260619-1942-sx2260-secondary-sales-net-sales-inc-ppn/discovery.md`, `index.json`.
- Key decisions:
  - Edit hanya `buildSecondarySalesReportGroupQuery` + test.
  - Tambah `vat_value_final` (order) & `vat_value` (return) ke ekspresi.
  - `* -1` tetap di akhir ekspresi return.
- Assumptions:
  - Kolom `od.vat_value_final` & `rd.vat_value` valid.
  - `sell_price1/2/3` (bukan `sell_price_final*`) sesuai existing repo.
- Open questions:
  - Apakah QA evidence memakai `sell_price_final*` atau `sell_price*`. Defer ke eksekusi.
- Readiness: ready-for-implementation (Test Plan + TDD langkah jelas, scope sempit, acceptance terukur).
- Cleanup performed: none (artifacts ini baru, tidak ada draft/evidence lama yang dihapus).
