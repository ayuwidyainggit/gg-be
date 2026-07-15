[BE] SX-2246 - Add Product di Purchase Order tab untuk status Need Review sudah dikerjakan di service `sales`.

Commit:
- `6e441fe` — `fix(sales-order): handle add purchase details`

Ringkasan pekerjaan:
- Endpoint `PATCH /sales/v1/orders/enhance/{order_no}` sekarang sudah handle payload `add_purchase_details` melalui flow enhance existing.
- Product tambahan dari tab Purchase Order akan diinsert sebagai row baru ke `sls.order_detail`.
- Field `pro_id` diambil dari payload `add_purchase_details[].pro_id`.
- Field `original_qty_po1`, `original_qty_po2`, `original_qty_po3` disimpan dari raw payload, tanpa cap stock.
- Field berikut disimpan dari hasil kalkulasi stock/UOM-aware cap:
  - `qty_po1`, `qty_po2`, `qty_po3`
  - `qty1`, `qty2`, `qty3`
  - `qty1_final`, `qty2_final`, `qty3_final`
- Stock calculation pakai existing helper `canonicalAPIStockBreakdown`, jadi bukan raw `min(qty, stock)` tanpa konversi UOM.
- Field `qty1_stok`, `qty2_stok`, `qty3_stok` disimpan saat insert product tambahan.
- Field `unit_id1`, `unit_id2`, `unit_id3` memakai value dari payload jika ada; fallback ke product master jika kosong.
- Field `is_product_promotion_po` memakai value dari payload; default `false` jika tidak dikirim.
- Ditambahkan guard untuk order yang tidak punya `wh_id` atau `ro_date` supaya tidak panic saat add product.
- Flow update existing `purchase_order` tetap dipertahankan, tidak diubah.
- Tidak ada perubahan schema, controller, repository, atau model.

Validasi yang sudah dilakukan:
- Unit test targeted:
  - `rtk go test ./service -run 'TestCreateOrderDetailFromPurchaseOrder|TestUpdateEnhance.*AddPurchase' -count=1`
  - Result: 7 passed.
- Full sales test suite:
  - `rtk go test ./... -count=1`
  - Result: 290 passed in 22 packages.
- Static check:
  - `rtk go vet ./entity/... ./service/...`
  - Result: no issues.
- DB local `ggn_scyllax` validation:
  - verified target columns in `sls.order_detail`.
  - verified sample products `8436` and `10813` from `mst.m_product`.
  - inserted temporary Need Review order, executed service insert against real DB, queried `sls.order_detail`, and verified mapped fields.
  - cleanup verified zero leftover temporary rows.
  - rollback proof with real transaction: forced failure on second add product and verified first insert was rolled back (`0` leftover `sls.order_detail` rows).

Files changed:
- `sales/entity/edit_order_enhance.go`
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`

Evidence:
- `.opencode/evidence/20260619-1754-sx-2246-add-purchase-details/verification.md`
- `.opencode/plans/20260619-1754-sx-2246-add-purchase-details.md`

Notes / follow-up:
- Current behavior tidak melakukan dedupe untuk retry/resend `add_purchase_details`, karena belum ada product rule eksplisit untuk idempotency.
- Mixed-UOM stock cap mengikuti plan decision: per-level cap setelah `canonicalAPIStockBreakdown`. Jika FE/product mengharapkan allocation by total smallest unit, perlu contoh rule tambahan untuk disepakati.
