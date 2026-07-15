# Sales Order Enhancement Final Order Bug Plan

## Ringkasan Masalah

Bug terjadi pada flow edit enhancement Sales Order ketika existing item di [`sales_order`](../sales/entity/edit_order_enhance.go:12) diubah menjadi qty 0 dan product baru ditambahkan lewat [`add_sales_order`](../sales/entity/edit_order_enhance.go:17).

Ekspektasi final yang sudah dikonfirmasi user:
- row qty 0 diperlakukan sebagai deleted-effective
- row qty 0 boleh tetap tersimpan di database
- row qty 0 tidak boleh ikut final calculation
- row qty 0 tidak boleh ikut stock visibility state, promo snapshot, VAT/PPN final, dan tampilan Final Order
- item baru dari [`add_sales_order`](../sales/entity/edit_order_enhance.go:17) harus ikut penuh pada kalkulasi final terbaru

## Scope Implementasi

Perubahan dibatasi pada service layer:
- [`UpdateEnhance`](../sales/service/order_service.go:4764)
- [`recomputePromoStateForTab`](../sales/service/order_service.go:2019)
- [`syncRewardProductState`](../sales/service/order_service.go:1833)
- [`DetailV2`](../sales/service/order_service.go:2475)
- [`createOrderDetailFromSalesOrder`](../sales/service/order_service.go:5379)
- tests di [`order_service_test.go`](../sales/service/order_service_test.go)

Update plan ini **tidak** memasukkan hard delete.

## Diagnosis Final (Disepakati)

Akar masalah bukan pada insert item baru, tetapi pada filtering state aktif:
1. Row qty 0 masih ikut dataset recompute promo/final pada [`recomputePromoStateForTab`](../sales/service/order_service.go:2019)
2. Builder response final pada [`DetailV2`](../sales/service/order_service.go:2475) belum mengecualikan row qty_final 0 dari [`DetailsFinal.Normal`](../sales/service/order_service.go:2765)

## Keputusan Desain Fix

Tidak ada utilitas KiloCode spesifik untuk deleted-effective filtering pada domain ini.

Maka implementasi menggunakan helper internal service layer (minimal-risk, reusable):
- `activeQtyForTab`
- `isActiveDetailForTab`
- `filterActiveNormalDetailsForTab`

Tujuan helper:
- mengecualikan row normal qty efektif 0 dari perhitungan aktif
- mempertahankan row di DB (soft behavior)
- menjaga item baru dari add flow tetap ikut penuh

## Rencana TDD

### Fase RED (Test dulu)
Tambahkan/ubah test di [`order_service_test.go`](../sales/service/order_service_test.go):
1. Skenario inti bug:
   - existing row A diubah qty jadi 0 via [`UpdateEnhance`](../sales/service/order_service.go:4764)
   - row baru B ditambahkan via `add_sales_order`
   - final header (`SubTotalFinal`, `VatValueFinal`, `TotalFinal`) hanya menghitung row B
   - stock release row A tetap benar (delta old->0 tetap terkirim)
2. Skenario response final:
   - [`DetailV2`](../sales/service/order_service.go:2475) hanya menampilkan row aktif di [`DetailsFinal.Normal`](../sales/service/order_service.go:2765)
   - row qty_final 0 tidak muncul
   - VAT/PPN final yang tampil hanya berasal dari row aktif
3. Helper-level test (jika ditambahkan):
   - validasi rule active filtering per tab

### Fase GREEN (Implementasi minimal)
1. Tambah helper deleted-effective filtering di [`order_service.go`](../sales/service/order_service.go)
2. Pakai helper pada dataset recompute di [`recomputePromoStateForTab`](../sales/service/order_service.go:2019)
3. Pakai helper pada response build final di [`DetailV2`](../sales/service/order_service.go:2475)
4. Tinjau [`syncRewardProductState`](../sales/service/order_service.go:1833) agar reward/promo tidak memakai row normal deleted-effective
5. Pastikan stock release row qty 0 tetap berjalan dari delta existing update (tidak diregresikan)

### Fase REFACTOR
- Rapikan logic yang redundant/usang akibat filtering baru bila aman
- Pastikan tidak ada penambahan logic hard delete

## Kriteria Done

- Semua test baru lulus
- Tidak ada hard delete baru
- Final calculation tidak memasukkan row qty 0
- Final response tidak menampilkan row qty 0
- VAT/PPN final berasal dari row aktif saja
- Item baru dari `add_sales_order` tetap masuk penuh ke final calculation

## Risiko & Mitigasi

- Risiko: consumer lain membaca semua row normal tanpa filter qty aktif
  - Mitigasi: pusatkan rule active filtering di helper service dan gunakan pada jalur final-critical terlebih dulu
- Risiko: perubahan filtering mengganggu stock behavior
  - Mitigasi: test khusus memastikan stock release row qty 0 tetap benar

## Catatan Implementasi

- Arsitektur tetap Controller → Service → Repository
- Business rule tetap di layer service
- Repository tidak ditambah business logic baru
- Tidak menambah hard delete DB
