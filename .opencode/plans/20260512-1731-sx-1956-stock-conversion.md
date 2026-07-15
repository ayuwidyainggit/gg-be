# Rencana Implementasi — SX-1956 Available Stock Conversion

## Goal

Memperbaiki perhitungan dan pemetaan `qty1_stok`, `qty2_stok`, `qty3_stok` pada service `sales` agar output available stock untuk `details.normal`, `purchase_details.normal`, dan `details_final.normal` selalu konsisten, canonical, dan sesuai kontrak `large / medium / small`.

## Non-goals

- Tidak mengubah module `pjp-sales` pada task ini.
- Tidak mengubah aturan bisnis total available stock source di luar kebutuhan SX-1956.
- Tidak melakukan refactor besar terhadap semua consumer `sales/pkg/conversion` di luar area stock breakdown.
- Tidak mengubah FE contract selain memperbaiki nilai field agar sesuai doc.

## Scope

- Module: `sales`
- Flow response detail order:
  - `DetailV2()` → `details.normal`
  - `DetailV2()` → `details_final.normal`
  - `DetailV2()` → `purchase_details.normal` via copy dari sales detail
- Flow persist snapshot stock yang memakai field sama:
  - create/store order
  - update enhance / refresh stock snapshot
- Unit test helper conversion baru / adapter
- Regression test detail response untuk SX-1956 dan existing SX-1878-adjacent behavior

## Requirements

1. Total available stock dalam satuan smallest/small harus tetap benar.
2. Breakdown response harus canonical mengikuti conversion product.
3. Mapping API yang dipakai:
   - `qty1_stok = large`
   - `qty2_stok = medium`
   - `qty3_stok = small`
4. Logic tidak boleh lagi menjumlahkan komponen L/M/S secara mentah tanpa renormalization.
5. `purchase_details`, `details`, dan `details_final` harus konsisten bila sumber stock sama.
6. Produk dengan conversion tidak lengkap tetap ditangani aman.
7. Regression terhadap validasi / processed scenario dari SX-1878 harus dipertahankan pada level total available stock.

## Acceptance Criteria

1. Untuk total small `13` dengan conversion yang relevan, hasil canonical menjadi `2 large, 0 medium, 3 small` dan dipetakan ke response sebagai:
   - `qty1_stok = 2`
   - `qty2_stok = 0`
   - `qty3_stok = 3`
   jika memang data master yang dipakai menghasilkan komposisi tersebut.
2. `details.normal`, `purchase_details.normal`, dan `details_final.normal` menghasilkan breakdown yang sama untuk product/source stock yang sama.
3. Reverse conversion dari `qty1_stok/qty2_stok/qty3_stok` kembali ke total small harus sama dengan total available stock source.
4. Produk tanpa medium / conversion parsial tidak panic dan fallback canonical-nya benar.
5. Existing tests yang sebelumnya mengunci breakdown non-canonical diperbarui agar memverifikasi canonical output baru.

## Existing Patterns/Reuse

- Reuse `conversion.QtyUnit.ToTotalQuantity()` untuk mengubah input L/M/S API ke total smallest, dengan adapter mapping yang benar.
- Reuse `conversion.Qty.ConvToQtyConversion()` untuk decompose total smallest ke helper internal `small/medium/large`, lalu map hasilnya ke API contract `large/medium/small`.
- Reuse `DetailV2` test harness di `sales/service/order_service_test.go`.
- Reuse `RefreshOrderDetailStock()` repository method; tidak perlu ubah repository kecuali naming helper internal butuh wrapper baru.

## Constraints

- Helper internal existing memakai urutan `Qty1=small`, `Qty2=medium`, `Qty3=large`.
- Doc lokal menjadi source-of-truth mapping response untuk task ini.
- Perubahan global semantik helper conversion berisiko memengaruhi service lain (`return`, `invoice`, `report`, `validate`).
- Karena itu fix harus memprioritaskan **adapter/helper stock-specific** daripada merombak seluruh package conversion.

## Risks

1. Ada test existing yang menganggap breakdown non-canonical sebagai behavior benar; test tersebut harus disesuaikan dengan justifikasi SX-1956.
2. Bila fix hanya dilakukan pada response detail tetapi tidak pada snapshot persist, data tersimpan `qty*_stok` tetap inconsistent.
3. Bila fix dilakukan langsung di helper umum tanpa audit, flow lain bisa regress karena semantik `Qty1/Qty2/Qty3` internal sudah dipakai luas.
4. Ambiguitas angka Jira `5,0,3` vs `2,0,3` harus dicatat sebagai data example inconsistency; implementasi mengikuti kontrak mapping + data master aktual repo.

## Decisions/Assumptions

- Keputusan user: scope **`sales only`**.
- Keputusan user: gunakan doc mapping sebagai sumber kebenaran, kecuali code evidence membuktikan kontrak lain.
- Keputusan user: targetkan **shared helper + regression tests**.
- Asumsi implementasi: issue utama berada pada dua hal sekaligus:
  1. salah mapping helper internal ke contract field, dan
  2. tidak ada renormalization setelah stock current digabung dengan qty order.
- Asumsi implementasi: `purchase_details` akan ikut benar setelah `details.normal` benar, karena source response-nya di-copy dari sana.
- Open question minor yang tersisa: contoh angka Jira `5,0,3` tampak tidak sinkron dengan narasi `2,0,3`; implementasi harus memprioritaskan conversion master dan mapping contract aktual, bukan literal angka yang saling konflik.

## TDD/Test Plan

### TDD Required

Ya, wajib. Alasannya karena ini adalah logic produksi untuk stock breakdown, menyentuh behavior API lintas tab, dan rentan regress bila hanya di-fix manual.

### Existing Test Patterns

- `sales/service/order_service_test.go`
  - `TestDetailV2_Cancelled_UsesWarehouseCurrentOnlyForDisplayedStock`
  - `TestDetailV2_NonCancelled_KeepsExistingDisplayedStockBehavior`
- Pattern mock repository `mockOrderRepositoryDetailV2`
- Existing stock snapshot/update tests di file yang sama

### First Failing / Regression Test

Tambahkan test unit/helper baru yang memverifikasi:

1. total small `13` + conversion tertentu → canonical API stock tuple `2,0,3`
2. gabungan warehouse + ordered qty harus dihitung dengan cara:
   - convert stock source ke total small
   - convert qty order ke total small
   - jumlahkan total small
   - decompose ulang ke canonical `L/M/S`

Lalu tambahkan / ubah regression test `DetailV2` untuk memastikan non-cancelled behavior tidak lagi `component-wise addition`, tetapi canonical normalized output.

### Green Step

- Implement helper stock adapter baru, misalnya fungsi service/private helper atau package helper lokal dengan shape seperti:
  - `toTotalSmallFromApiUnits(qtyLarge, qtyMedium, qtySmall, conv2, conv3)`
  - `toApiStockBreakdown(totalSmall, conv2, conv3)`
  - `computeDisplayedAvailableStockBreakdown(warehouseSmall, orderedLarge, orderedMedium, orderedSmall, includeOrder bool, conv2, conv3)`
- Pakai helper itu di semua titik perhitungan `qty*_stok` response dan persist snapshot.

### Refactor Step

- Hapus duplikasi formula stock breakdown di `DetailV2`, create/store, dan update enhance dengan helper shared yang sama.
- Tambahkan naming/comment yang eksplisit bahwa helper internal `conversion.Qty` memakai urutan `small, medium, large`, sedangkan API field memakai `large, medium, small`.

### Edge Cases

- total small = 0
- habis dibagi large
- menyisakan medium dan small
- tanpa medium / `convUnit2 <= 0` atau `convUnit3 <= 0`
- stock negative bila business rule mengizinkan
- cancelled vs non-cancelled displayed stock behavior

### Commands

- `rtk go test ./...`
- bila perlu fokus cepat:
  - `rtk go test ./service -run TestDetailV2`
  - `rtk go test ./... -run SX1956`

## Implementation Steps

1. Tambahkan helper shared khusus stock breakdown di module `sales` tanpa mengubah semantik global `conversion.Qty` / `conversion.QtyUnit` secara langsung.
2. Implement adapter yang mengubah hasil helper internal `small/medium/large` menjadi API contract `large/medium/small`.
3. Ubah perhitungan displayed stock di `DetailV2()` sales tab:
   - gunakan warehouse stock total small
   - bila non-cancelled, tambah qty order dalam total small
   - decompose ulang ke canonical API stock tuple
4. Ubah perhitungan displayed stock di `DetailV2()` final tab dengan helper yang sama.
5. Pastikan `purchase_details.normal` tetap ikut hasil corrected path karena copy dari sales detail.
6. Ubah refresh snapshot create/store flow agar `qty1_stok/qty2_stok/qty3_stok` yang dipersist juga mengikuti mapping benar.
7. Ubah refresh snapshot `UpdateEnhance()` agar memakai helper shared yang sama.
8. Tambahkan unit test helper conversion.
9. Tambahkan / perbarui regression tests `DetailV2` untuk canonical breakdown dan cancelled/non-cancelled scenarios.
10. Jalankan test module `sales`, evaluasi test yang perlu disesuaikan karena sebelumnya mengunci behavior salah.

## Expected Files to Change

- `sales/service/order_service.go`
- `sales/service/order_service_test.go`
- kemungkinan file helper baru, contoh salah satu dari:
  - `sales/service/order_stock_helper.go`, atau
  - `sales/pkg/conversion/stock_breakdown.go`

## Agent/Tool Routing

- Artifact planner: selesai pada plan ini.
- Implementasi berikutnya sebaiknya oleh `@fixer` atau workflow implementasi bounded.
- Quality/review akhir sebaiknya lewat `@quality-gate` karena perubahan menyentuh logic stock dan response contract.

## Validation Commands

- `rtk go test ./...`
- `rtk go test ./service -run TestDetailV2`
- `rtk go test ./service -run Test.*Stock.*`

## Evidence Requirements

- Buktikan dari test bahwa helper baru menghasilkan breakdown canonical.
- Buktikan dari test `DetailV2` bahwa:
  - cancelled memakai warehouse current only
  - non-cancelled memakai total small gabungan lalu normalized ulang
  - purchase / sales / final konsisten untuk source stock sama
- Catat root cause di PR:
  1. helper internal vs contract mapping mismatch
  2. component-wise addition tanpa renormalization

## Done Criteria

- Helper shared stock breakdown tersedia dan dipakai lintas titik perhitungan yang relevan.
- `qty1_stok/qty2_stok/qty3_stok` mengikuti mapping kontrak doc.
- Breakdown canonical, bukan sekadar total small yang benar.
- Response PO/SO/Final konsisten.
- Snapshot stock yang dipersist tidak lagi terbalik mapping-nya.
- Unit test dan regression test lulus.
- Catatan root cause siap dipakai untuk PR description.

## Final Planning Summary

- Sumber kebenaran implementasi: **`.opencode/plans/20260512-1731-sx-1956-stock-conversion.md`**.
- Artefak yang dibuat:
  - `.opencode/plans/20260512-1731-sx-1956-stock-conversion.md`
  - `.opencode/evidence/20260512-1731-sx-1956-stock-conversion/discovery.md`
- Keputusan utama:
  - scope `sales only`
  - mapping contract `qty1=L`, `qty2=M`, `qty3=S`
  - fix melalui shared helper + regression tests
- Asumsi utama:
  - issue berasal dari mismatch mapping dan lack of renormalization setelah penggabungan total stock
- Pertanyaan user:
  - sudah ditanyakan dan sudah dijawab; tidak ada blocker material tersisa
- Readiness:
  - siap untuk implementasi bounded
- Cleanup:
  - tidak ada draft file yang perlu dipertahankan sebagai source of truth
  - evidence discovery dipertahankan karena masih berguna sebagai jejak root-cause untuk implementer/PR
