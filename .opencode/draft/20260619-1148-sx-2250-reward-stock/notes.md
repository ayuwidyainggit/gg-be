# Notes SX-2250 — Reward Product Stock Display

Task id: `20260619-1148-sx-2250-reward-stock`

## Current code summary (from discovery)

- `DetailV2` (sales/service/order_service.go:2856) iterates `details` 3 kali:
  1. Sales order `Details.Normal` (item_type==1 ke Normal, item_type!=1 ke Promo lalu `movePromoDetailsToNormal`).
  2. Purchase `PurchaseDetails.Normal` (mirip).
  3. Final order `DetailsFinal.Normal` (menggunakan `QtyFinal`/`Qty1Final`/dst).
- Tiap iterasi hitung `qty*_stok` lewat `computeDisplayedAvailableStockBreakdown(whStock, rowQty1, rowQty2, rowQty3, !useWarehouseCurrentOnly, conv2, conv3)` — row-scoped.
- `useWarehouseCurrentOnly = true` saat `DataStatus == CANCELLED`.
- `whStock = warehouseStockMap[int64(detail.ProId)]` — ini product-level warehouse stock current, bukan penjumlahan SO qty.

## Hipotesis defect

- Berdasarkan baris kode di branch ini, `qty*_stok` per row sudah row-scoped.
- Bug paling mungkin: kurang test yang mengunci behavior row-level untuk kasus normal+reward same `pro_id`. Tanpa test, regresi mudah masuk lagi.
- Risk kedua: kalau executor berikutnya salah baca argumen `computeDisplayedAvailableStockBreakdown`, bisa saja terjadi refactor yang malah salah.
- Concern 1 (On Customer Stock mutation): tidak ada perubahan yang diminta user sekarang untuk hal itu (komentar terbaru Jira menyatakan movement evidence sudah sesuai). Tetap dicatat di plan sebagai out-of-scope kecuali ditemukan bug saat eksekusi.

## Test minimum yang harus dikunci (eksplisit sesuai defect)

1. Same pro_id, normal row + reward row, qty sama → dua row `qty*_stok = 0 0 3` (warehouse 0 0 2 + row 0 0 1).
2. Same pro_id, qty beda → masing-masing row `qty*_stok` independen.
3. Non-promotion single row → behavior existing tidak berubah.
4. (Opsional) Reward row on-customer stock mutation — hanya jika test existing gagal.
