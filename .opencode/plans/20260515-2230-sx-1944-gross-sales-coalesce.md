# SX-1944 — Gross Sales COALESCE Plan

## Goal

Memperbaiki bagian SX-1944 yang menjadi tanggung jawab saat ini: `GrossSales` pada Secondary Sales Report tidak boleh menjadi `0`/NULL akibat salah satu field numerik NULL di query. Perbaikan hanya berupa hardening query report dengan `COALESCE` pada kalkulasi GrossSales.

## Non-goals

- Tidak memperbaiki mapping insert `sls.order_detail.amount_final` dan `sls.order_detail.vat_value_final` di `mobile/v1/orders`.
- Tidak melakukan backfill data lama.
- Tidak mengubah kontrak FE/Mobile.
- Tidak mengubah flow RMQ/export Secondary Sales.
- Tidak mengubah formula bisnis `NetSalesExcPPN`, `PPN`, atau `NetSalesIncPPN` kecuali hanya bila diperlukan untuk test isolasi GrossSales.
- Tidak mengubah schema DB, migration, atau seed data.

## Scope

- Module target: `sales/`.
- File utama:
  - `sales/repository/report_repository.go`
  - `sales/repository/report_repository_test.go`
- Query/path audit:
  - `buildSecondarySalesUnionQuery` untuk endpoint/export live `POST /sales/v1/reports/secondary-sales`.
  - `RepositoryReportImpl.SecondarySales` sebagai legacy defensive hardening.
  - `GetReportSecondarySalesReportOrder` dan `GetReportSecondarySalesReportReturn` untuk cron extract/dashboard fact table.
  - `SecondarySalesReportSumReportByMonth` untuk agregasi dashboard gross sales.

## Requirements

1. Kalkulasi `gross_sales` harus memakai `COALESCE(field, 0)` pada seluruh operand numerik yang dapat NULL.
2. Field minimal yang wajib NULL-safe pada GrossSales:
   - `qty1_final`, `qty2_final`, `qty3_final`
   - `sell_price1`, `sell_price2`, `sell_price3`
   - untuk return: `qty1`, `qty2`, `qty3`, `sell_price1`, `sell_price2`, `sell_price3`
3. Path live endpoint/export harus tetap memakai `buildSecondarySalesUnionQuery` sebagai single source of truth.
4. Legacy/dead query `SecondarySales` boleh ikut di-hardening supaya tidak jadi jebakan jika dipakai lagi nanti.
5. Dashboard extract path harus di-hardening karena menulis `report.fact_orders.gross_sale` / `report.fact_returns.gross_sale`.
6. Agregasi `total_gross_sale` dashboard harus memakai `COALESCE(SUM(...), 0)` untuk konsisten dengan trend sales.
7. Query tetap parameterized; tidak hardcode `SO`, `INV`, `cust_id`, token, atau credential QA.

## Acceptance Criteria

- `gross_sales` di semua query Secondary Sales relevan tidak lagi menjadi NULL hanya karena salah satu qty/sell price NULL.
- Existing `buildSecondarySalesUnionQuery` tetap mengandung `COALESCE` pada order dan return branch GrossSales.
- `GetReportSecondarySalesReportOrder` menghitung GrossSales dengan `COALESCE(od.qty*_final, 0) * COALESCE(od.sell_price*, 0)`.
- `GetReportSecondarySalesReportReturn` menghitung GrossSales dengan `COALESCE(rd.qty*, 0) * COALESCE(rd.sell_price*, 0)`.
- `RepositoryReportImpl.SecondarySales` legacy query ikut NULL-safe untuk GrossSales.
- `SecondarySalesReportSumReportByMonth` memakai `COALESCE(SUM(report.fact_orders.gross_sale), 0) AS total_gross_sale`.
- Regression test query-level membuktikan SQL builder/path relevan memuat `COALESCE` untuk GrossSales.
- `rtk go test ./repository -run SecondarySales` lulus dari direktori `sales/`.
- Bila waktu cukup: `rtk go test ./...` lulus dari direktori `sales/`.

## Existing Patterns/Reuse

- `buildSecondarySalesUnionQuery` sudah memakai pola benar:
  - `COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)`
  - pola sama untuk `qty2/qty3` dan return branch.
- `sales/repository/report_repository_test.go` sudah punya test query-builder berbasis substring SQL; reuse pattern ini untuk regression COALESCE.
- `SecondarySalesReportTrendSales` sudah memakai `COALESCE(SUM(fo.gross_sale), 0)`, dapat dijadikan pattern untuk `SecondarySalesReportSumReportByMonth`.
- Tidak ada kebutuhan membuat helper baru bila patch hanya mengganti ekspresi SQL inline.

## Constraints

- Ikuti layering repo: Controller → Service → Repository → DB.
- Perubahan hanya di module `sales/`.
- Jangan menyimpan token/staging credential dari Jira.
- Jangan jalankan backfill production/staging tanpa approval eksplisit.
- Jangan menggunakan `rtk trust` atau mengubah konfigurasi RTK.
- Runtime awal: `rtk docker compose -f docker-compose.yml ps` sudah dijalankan; service stack tidak berjalan, tapi tidak memblokir unit/query-level plan.

## Risks

1. **False path risk**
   - Endpoint live export sudah aman untuk GrossSales di builder utama. Jika fix hanya dilakukan di path legacy yang tidak dipakai, efek QA endpoint live tidak berubah.
   - Mitigasi: patch semua path Secondary Sales yang menghitung GrossSales, khususnya extract/dashboard path.
2. **Out-of-scope confusion**
   - `NetSalesExcPPN`, `PPN`, dan `NetSalesIncPPN` tetap bisa `0` bila `amount_final`/`vat_value_final` belum diisi oleh mobile order insert.
   - Mitigasi: sebut eksplisit bahwa itu task lain.
3. **NULL-to-zero semantics**
   - `COALESCE` mengubah hasil NULL menjadi `0`. Untuk GrossSales ini sesuai Jira.
4. **Report fact stale data**
   - Data yang sudah pernah diekstrak ke `report.fact_orders` dengan nilai salah mungkin butuh re-extract/backfill. Di luar scope kecuali diminta lead/QA.

## Decisions/Assumptions

### Decisions

- Scope final: `GrossSales` COALESCE saja.
- Tidak mengubah mobile order insert.
- Test prioritas: query-level unit tests di `sales/repository/report_repository_test.go`.

### Assumptions / Open Questions

- Asumsi: Jira comment #1 dipisah dari #2; developer lain menangani `amount_final`/`vat_value_final` mobile insert.
- Asumsi: Expected angka NetSales/PPN pada QA case tidak harus diselesaikan oleh patch ini sendiri.
- Open question non-blocking: apakah QA/lead ingin re-extract/backfill fact table setelah patch? Tidak perlu untuk code fix, tapi perlu untuk data lama yang sudah telanjur tersimpan salah.

## TDD/Test Plan

### TDD required

Ya. Ini bug report numerik/query logic.

### Reason

Perubahan SQL kecil bisa regress silent dan hasil export/dashboard menjadi salah tanpa compile error.

### Existing test patterns

- `sales/repository/report_repository_test.go`
  - `TestBuildSecondarySalesUnionQueryUsesTransCTEAndDeterministicParentProductJoin`
  - `TestBuildSecondarySalesUnionQueryPreservesFilterAliases`
- Pattern: build SQL string lalu `strings.Contains(sql, expected)`.

### First failing/regression test

Tambahkan test di `sales/repository/report_repository_test.go`, misalnya:

- `TestBuildSecondarySalesUnionQueryCoalescesGrossSalesOperands`
  - assert SQL mengandung:
    - `COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)`
    - `COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)`
    - `COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0)`
    - `COALESCE(rd.qty1, 0) * COALESCE(rd.sell_price1, 0)`
    - `COALESCE(rd.qty2, 0) * COALESCE(rd.sell_price2, 0)`
    - `COALESCE(rd.qty3, 0) * COALESCE(rd.sell_price3, 0)`

Tambahan test opsional/berguna:

- Test atau assertion helper baru untuk raw SQL extract jika dibuat helper function. Jika tidak ada helper, cukup reviewer check karena raw SQL inline sulit dites tanpa DB.

### Green step

- Patch `sales/repository/report_repository.go`:
  - `SecondarySales` gross expression: tambahkan `COALESCE` pada qty/sell_price.
  - `GetReportSecondarySalesReportOrder` gross expression: tambahkan `COALESCE` pada qty/sell_price.
  - `GetReportSecondarySalesReportReturn` gross expression: tambahkan `COALESCE` pada qty/sell_price.
  - `SecondarySalesReportSumReportByMonth`: `COALESCE(SUM(report.fact_orders.gross_sale), 0) AS total_gross_sale`.

### Refactor step

- Bila ada duplikasi ekspresi terlalu besar, boleh extract konstanta/helper SQL kecil, tetapi jangan refactor luas.
- Pastikan format SQL tetap mudah dibaca dan tidak mengubah filter/joins.

### Edge cases

- `qty1_final` NULL, `sell_price1` non-NULL → kontribusi unit 1 menjadi 0, unit lain tetap dihitung.
- `sell_price2` NULL, `qty2_final` non-NULL → kontribusi unit 2 menjadi 0.
- Semua qty/sell price NULL → GrossSales menjadi 0, bukan NULL.
- Return qty/sell price NULL → GrossSales return menjadi 0 atau jumlah unit lain yang tersedia, lalu dikalikan sign existing.
- Bulan dashboard tanpa fact order → `total_gross_sale` 0, bukan NULL.

### Commands

Dari repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

Dari `sales/`:

```bash
rtk go test ./repository -run SecondarySales
rtk go test ./...
```

## Implementation Steps

1. Buka `sales/repository/report_repository.go`.
2. Pastikan `buildSecondarySalesUnionQuery` tetap tidak diubah secara bisnis; hanya tambah/pertahankan test COALESCE untuk GrossSales.
3. Patch legacy `SecondarySales` expression:
   - dari `qty*_final*sell_price*`
   - menjadi `COALESCE(qty*_final, 0) * COALESCE(sell_price*, 0)`.
4. Patch `GetReportSecondarySalesReportOrder` expression:
   - dari `od.qty*_final * od.sell_price*`
   - menjadi `COALESCE(od.qty*_final, 0) * COALESCE(od.sell_price*, 0)`.
5. Patch `GetReportSecondarySalesReportReturn` expression:
   - dari `rd.qty* * rd.sell_price*`
   - menjadi `COALESCE(rd.qty*, 0) * COALESCE(rd.sell_price*, 0)`.
6. Patch `SecondarySalesReportSumReportByMonth`:
   - `SUM(report.fact_orders.gross_sale)` → `COALESCE(SUM(report.fact_orders.gross_sale), 0)`.
7. Tambahkan regression test query builder untuk memastikan live endpoint/export builder tetap COALESCE-safe.
8. Jalankan targeted test.
9. Jika targeted lulus, jalankan full module test jika waktu cukup.
10. Catat bahwa NetSales/PPN task terpisah dan mungkin tetap gagal sampai mapping mobile insert diperbaiki.

## Expected Files to Change

- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`

Tidak perlu ubah:

- `mobile/service/order.go`
- `mobile/service/order_canvas.go`
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- migration/DDL files

## Agent/Tool Routing

- `@artifact-planner`: plan + discovery evidence selesai.
- `@fixer`: eksekusi patch SQL dan test.
- `@quality-gate`: review akhir karena bug report numerik dan report export berdampak bisnis.
- `@explorer`: hanya bila implementer perlu discovery tambahan.
- `@architect` tidak diperlukan; scope query-level kecil.
- `@librarian` tidak diperlukan; tidak ada library behavior version-sensitive.

## Execution-ready Worklist / Handoff Contract

`start_with`: `T1`

| Task | Action | depends_on | owner/lane | validation | exit criteria | status | requires_user_decision |
| --- | --- | --- | --- | --- | --- | --- | --- |
| T1 | Tambah regression test untuk COALESCE GrossSales di `buildSecondarySalesUnionQuery` | none | `@fixer` | `rtk go test ./repository -run TestBuildSecondarySalesUnionQueryCoalescesGrossSalesOperands` dari `sales/` | Test gagal jika COALESCE hilang pada order/return branch GrossSales | ready | no |
| T2 | Patch legacy `SecondarySales` GrossSales expression dengan COALESCE | T1 | `@fixer` | Review diff + compile via targeted test | Semua operand qty/sell price di expression GrossSales legacy NULL-safe | ready | no |
| T3 | Patch `GetReportSecondarySalesReportOrder` GrossSales expression dengan COALESCE | T2 | `@fixer` | Review diff + compile via targeted test | Semua operand `od.qty*_final` dan `od.sell_price*` NULL-safe | ready | no |
| T4 | Patch `GetReportSecondarySalesReportReturn` GrossSales expression dengan COALESCE | T3 | `@fixer` | Review diff + compile via targeted test | Semua operand `rd.qty*` dan `rd.sell_price*` NULL-safe | ready | no |
| T5 | Patch `SecondarySalesReportSumReportByMonth` `total_gross_sale` aggregation dengan `COALESCE(SUM(...),0)` | T4 | `@fixer` | Review diff + compile via targeted test | Dashboard monthly gross sale menghasilkan 0, bukan NULL, saat tidak ada data | ready | no |
| T6 | Jalankan targeted repository tests | T5 | `@fixer` | `rtk go test ./repository -run SecondarySales` | Test lulus | ready | no |
| T7 | Jalankan full sales module tests bila waktu cukup | T6 | `@fixer` | `rtk go test ./...` | Full module lulus atau failure unrelated dicatat | ready | no |
| T8 | Final quality review | T7 | `@quality-gate` | Review diff + test evidence | Scope tidak melebar ke mobile/backfill, COALESCE coverage cukup, risk dicatat | ready | no |

## Validation Commands

Dari repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

Dari `sales/`:

```bash
rtk go test ./repository -run SecondarySales
rtk go test ./...
```

Manual SQL sanity bila ada akses DB aman:

```sql
SELECT
  ((COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)) +
   (COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)) +
   (COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0))) AS gross_sales
FROM sls.order_detail od
JOIN sls."order" o ON o.ro_no = od.ro_no AND o.cust_id = od.cust_id
WHERE o.invoice_no = 'INV2605080009';
```

Catatan: jangan hardcode query manual ini dalam code/test; ini hanya sanity check lokal/staging bila credential resmi tersedia.

## Evidence Requirements

Implementation dianggap siap review bila ada:

- Diff `report_repository.go` menunjukkan semua GrossSales expression di path yang disebut sudah pakai `COALESCE`.
- Diff test menunjukkan regression test COALESCE untuk builder live endpoint/export.
- Output `rtk go test ./repository -run SecondarySales`.
- Bila full test tidak dijalankan, alasan dicatat.
- Bila manual staging check tidak dilakukan, alasan dicatat (mis. tidak ada token/DB access aman).

## Done Criteria

- Semua Acceptance Criteria terpenuhi.
- Scope tetap hanya GrossSales COALESCE.
- Tidak ada credential/token/backfill/migration ditambahkan.
- Tests relevan lulus atau blocker eksplisit dicatat.
- Quality gate menyetujui atau issue minor dicatat untuk follow-up.

## Final Planning Summary

- Artifact utama dibuat: `.opencode/plans/20260515-2230-sx-1944-gross-sales-coalesce.md`.
- Evidence dibuat dan dipertahankan: `.opencode/evidence/20260515-2230-sx-1944-gross-sales-coalesce/discovery.md` karena berisi audit path query dan penting untuk implementer agar tidak salah patch path dead saja.
- Draft tidak dibuat; tidak ada draft stale untuk dibersihkan.
- Keputusan kunci: task ini hanya bagian #1 (`GrossSales` NULL → `COALESCE`), bukan mapping `amount_final`/`vat_value_final` mobile.
- Research gate: local discovery selesai; official docs/GitHub/web/browser tidak diperlukan.
- Readiness: siap dieksekusi oleh `@fixer` tanpa replanning.
