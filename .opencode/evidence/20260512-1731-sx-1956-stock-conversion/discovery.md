# Discovery — SX-1956 Available Stock Conversion

## Ringkasan

- Scope yang disetujui user: **`sales` only**, tidak menyentuh `pjp-sales`.
- Mapping kontrak yang dipakai sebagai sumber kebenaran: **`qty1_stok = large`, `qty2_stok = medium`, `qty3_stok = small`**.
- Strategi fix yang disetujui user: **shared helper + regression tests**.

## File yang Diinspeksi

- `docs/Sales Order Enhancement_BE.md`
- `sales/pkg/conversion/quantity.go`
- `sales/pkg/conversion/qtyunit.go`
- `sales/service/order_service.go`
- `sales/service/validate_order_service.go`
- `sales/repository/order_repository.go`
- `sales/service/order_service_test.go`
- `sales/model/order_detail.go`

## Pola Existing yang Ditemukan

### 1. Konversi total small -> breakdown existing

File: `sales/pkg/conversion/quantity.go`

- `ConvToQtyConversion()` saat ini menghitung:
  - `Qty3 = total / (conv2 * conv3)` → level terbesar / large
  - `Qty2 = remainder / conv2` → medium
  - `Qty1 = remainder akhir` → small
- Artinya output helper internal adalah **`Qty1=small`, `Qty2=medium`, `Qty3=large`**.

### 2. Konversi breakdown -> total small existing

File: `sales/pkg/conversion/qtyunit.go`

- `ToTotalQuantity()` memakai rumus:
  - `total = (conv2*conv3)*Qty3 + conv2*Qty2 + Qty1`
- Ini mengonfirmasi model internal helper qty adalah:
  - `Qty1 = small`
  - `Qty2 = medium`
  - `Qty3 = large`

### 3. Kontrak doc untuk response stock

File: `docs/Sales Order Enhancement_BE.md:133-135`

- Dokumen menyebut explicit mapping response:
  - `qty1_stok` = qty1 **(L)**
  - `qty2_stok` = qty2 **(M)**
  - `qty3_stok` = qty3 **(S)**
- Jadi kontrak response **berlawanan urutan** dengan helper internal `conversion.Qty`.

### 4. Lokasi pengisian stock breakdown di response detail

File: `sales/service/order_service.go`

- `DetailV2()` sales tab:
  - sekitar `2841-2857`
- `DetailV2()` final tab:
  - sekitar `3051-3067`
- `PurchaseDetails` dibuat dengan `copy(response.Details.Normal, ...)`, sehingga nilai stok purchase mengikuti hasil sales detail.

### 5. Lokasi pengisian stock snapshot saat persist/update

File: `sales/service/order_service.go`

- `Store` / create flow sekitar `4125-4138`
- `UpdateEnhance` refresh snapshot sekitar `5408-5417`
- Keduanya langsung assign:
  - `qty1_stok = stockConversion.Qty1`
  - `qty2_stok = stockConversion.Qty2`
  - `qty3_stok = stockConversion.Qty3`
- Dengan helper internal current, assignment ini berarti:
  - `qty1_stok` terisi **small**, bukan large
  - `qty3_stok` terisi **large**, bukan small

### 6. Lokasi persist field qty*_stok

File: `sales/repository/order_repository.go:1043-1051`

- `RefreshOrderDetailStock()` hanya persist nilai yang sudah dihitung service.
- Root cause bukan di repository, melainkan di service/helper mapping dan normalization.

## Root Cause Kandidat Terkuat

### A. Salah mapping helper internal ke kontrak response/API

- Helper internal menghasilkan urutan `small, medium, large`.
- API contract butuh urutan `large, medium, small`.
- Existing assignment langsung 1:1 menyebabkan field terbalik.

### B. Tidak ada re-normalization setelah stock current + qty order digabung

File: `sales/service/order_service.go:2850-2853`, `3060-3063`

- Existing logic untuk non-cancelled case melakukan penjumlahan per komponen:
  - `qty1_stok += detailData.Qty1`
  - `qty2_stok += detailData.Qty2`
  - `qty3_stok += detailData.Qty3`
- Ini menjaga **total small** tetap benar, tetapi **breakdown bisa invalid** karena carry antar unit tidak dinormalisasi ulang.
- Contoh pattern bug:
  - warehouse normalized + order normalized dijumlahkan per field
  - hasil total ekuivalen benar dalam small
  - tetapi tuple L/M/S bisa menjadi bentuk non-canonical seperti `8 0 1` alih-alih `2 0 3`.

## Reuse Candidates

- Reuse `sales/pkg/conversion/qtyunit.go` untuk konversi L/M/S -> total small dengan catatan mapping API harus diadaptasi lebih dulu.
- Reuse `sales/pkg/conversion/quantity.go` untuk total small -> breakdown, tetapi perlu adapter/mapping yang jelas atau helper baru agar output API menjadi `L/M/S`.
- Reuse existing `order_service_test.go` pattern untuk `DetailV2` regression coverage.

## Constraint Teknis

- Repo menggunakan beberapa flow yang sudah mengandalkan helper `conversion.Qty` dan `conversion.QtyUnit`; perubahan langsung ke semantik helper berisiko memecahkan flow lain seperti return, invoice, report, dan validate.
- Karena itu lebih aman membuat **adapter/helper khusus stock breakdown API** daripada mengubah kontrak internal package conversion secara global tanpa audit luas.
- `PurchaseDetails` adalah copy dari `Details.Normal`, jadi fix di sales detail path otomatis mempengaruhi purchase tab response.

## Risiko

- Test existing `TestDetailV2_NonCancelled_KeepsExistingDisplayedStockBehavior` saat ini mengunci behavior lama yang belum normalized; test ini kemungkinan harus diubah karena issue SX-1956 memang mengubah expected behavior.
- Bila helper conversion global diubah tanpa adapter, flow lain yang mengandalkan `Qty1=small, Qty2=medium, Qty3=large` dapat regress.
- Snapshot persist (`qty*_stok` di DB) bisa tetap salah walaupun response sudah benar bila create/update flow tidak ikut diperbaiki.

## Commands / Pencarian yang Dijalankan

- `rtk docker compose -f docker-compose.yml ps`
- search `qty1_stok|qty2_stok|qty3_stok`
- search `available[_ ]stock|warehouse_stock|oncust|qty.*stok`
- review targeted file reads untuk helper conversion, service, repository, tests, dan dokumen API

## Research Gate

- Local project discovery: **required** dan sudah dilakukan.
- Official docs/context7: **tidak diperlukan**, karena isu terikat ke logic repo internal dan kontrak repo doc lokal.
- GitHub/upstream: **tidak diperlukan**, tidak ada dependensi behavior upstream.
- Brave/web search: **tidak diperlukan**, tidak ada fakta eksternal current yang memengaruhi keputusan.
- Browser/screenshot: **tidak diperlukan**, task BE non-visual.
