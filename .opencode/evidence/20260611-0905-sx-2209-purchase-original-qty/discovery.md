# Discovery SX-2209 Purchase Details Original Qty

Task ID: `20260611-0905-sx-2209-purchase-original-qty`
Tanggal: 2026-06-11 09:05 Asia/Jakarta
Mode: Maintenance Stability Mode, artifact-only planning.

## Source strategy

- Dipakai: repo-local evidence dari `AGENTS.md`, `.opencode/docs/index.md`, `.opencode/docs/ARCHITECTURE.md`, `.opencode/docs/QUALITY.md`, `sales` module code, migration, model/entity, service, controller, dan test.
- Dipakai: detail issue dari prompt user sebagai Jira evidence karena konteks, sample row, dan acceptance criteria sudah disediakan.
- Dilewati: Atlassian/web fetch karena kemungkinan butuh credential dan prompt sudah memuat data issue yang cukup.
- Dilewati: Context7/library docs karena perubahan memakai Go/GORM/repo helpers existing, bukan API library yang unfamiliar atau version-sensitive.
- Dilewati: runtime DB/API smoke saat planning karena planner tidak melakukan implementasi/source edit dan token/API env tidak tersedia.

## File dan pola yang diinspeksi

- `AGENTS.md`: aturan repo, service `sales`, perintah `rtk`, layer Controller → Service → Repository → DB, larangan secret/token.
- `.opencode/docs/ARCHITECTURE.md`: batas layer, tenant/schema rules, service ownership.
- `.opencode/docs/QUALITY.md`: validasi target service di `sales`, `rtk go test ./...`, targeted `go test`.
- `.opencode/docs/AGENT_ROUTING.md`: bounded code/test oleh `@fixer`, final signoff oleh `@quality-gate`.
- `sales/controller/order_controller.go`: `Route()` mendaftarkan `GET /v2/orders/:ro_no` ke `DetailV2`; controller memanggil `OrderService.DetailV2(params.RoNo, custId, parentCustId)` untuk path normal.
- `sales/repository/order_repository.go`: `FindDetail` memakai `Select("sls.order_detail.*, p.pro_code, p.pro_name, p.conv_unit2 as mconv_unit2,p.conv_unit3 as mconv_unit3")`, sehingga field `original_qty_po*` ikut terbaca dari `sls.order_detail.*`.
- `sales/model/order_detail.go`: `OrderDetail` dan `OrderDetailRead` sudah punya `OriginalQtyPo1/2/3`, `QtyPo1/2/3`, dan `SellPricePo1/2/3`.
- `sales/entity/order_detail.go`: `OrderDetResponse` sudah expose JSON `original_qty_po1/2/3` dan `qty_po1/2/3`.
- `sales/migration/sls.order/add_order_type_and_original_qty_po_fields.sql`: migration SX-2184 menambahkan `original_qty_po1/2/3` ke `sls.order_detail`.
- `sales/service/order_type_helper.go`: Taking Order `O` mengisi `QtyPo1/2/3` dan `OriginalQtyPo1/2/3` dari payload original.
- `sales/service/order_service.go`: `DetailV2` membangun `PurchaseDetails` dari persisted rows, bukan copy sales details.
- `sales/service/order_service_test.go`: sudah ada test `TestDetailV2_PurchaseDetailsUsesPurchaseActiveRowsForOrderTypeO`, mock `mockOrderRepositoryDetailV2`, dan pola validasi `PurchaseDetails.Normal`.

## Flow endpoint yang ditemukan

`GET /sales/v2/orders/{order_no}` pada deployment kemungkinan diprefix `/sales`; di service code route internal adalah `GET /v2/orders/:ro_no`.

Flow normal:

`OrderController.DetailV2` → `OrderService.DetailV2` → `OrderRepository.FindByNo` + `FindDetail` + `FindReward` → mapper `structs.Automapper` → build `response.PurchaseDetails.Normal`.

## Bug candidate utama

Di `sales/service/order_service.go`:

- `activeQtyForTab(detail, promoSnapshotTabPurchase)` saat ini hanya menjumlahkan `QtyPo1 + QtyPo2 + QtyPo3`.
- `isActiveDetailForTab(detail, promoSnapshotTabPurchase)` mengembalikan `activeQtyForTab(...) > 0`.
- Saat build `PurchaseDetails.Normal`, append terjadi pada kondisi:

```go
if isActiveDetailForTab(detail, promoSnapshotTabPurchase) || isActiveDetailForTab(detail, promoSnapshotTabSalesOrder) {
    response.PurchaseDetails.Normal = append(response.PurchaseDetails.Normal, detailData)
}
```

Konsekuensi: row dengan `qty_po1/2/3 = 0` dan sales qty juga kosong/zero tidak masuk, walaupun `original_qty_po3 = 3`. Ini cocok dengan sample SX-2209 `order_detail_id = 7273`.

## Reuse candidates

- Gunakan `getValueOrDefault(*float64, fallback)` existing untuk null-safe float handling.
- Tambahkan helper kecil di `sales/service/order_service.go` dekat `activeQtyForTab` agar rule bisa dites dan tidak mengubah repository query.
- Reuse `mockOrderRepositoryDetailV2` dan test pattern `TestDetailV2_PurchaseDetailsUsesPurchaseActiveRowsForOrderTypeO` untuk regression test.

## Constraints dan risiko

- Jangan ubah stock mutation/process order; target hanya display response `purchase_details`.
- Jangan mengganti `qty_po*` dengan `original_qty_po*`; response harus expose keduanya sesuai field existing.
- Hati-hati jika mengubah `activeQtyForTab` global karena helper itu dipakai untuk promo snapshot, promo consult, recompute, dan tab lain. Lebih aman tambahkan predicate display khusus untuk `purchase_details` response.
- Row original-only akan ikut `purchase_details.normal` dan kemudian dapat ikut promo consultation payload dengan current qty zero. Perlu test bahwa tidak panic dan promo amount tetap berbasis `qty_po*`, bukan `original_qty_po*`.
- `FindDetail` sudah memilih `sls.order_detail.*`; tidak perlu raw SQL scan update berdasarkan evidence saat ini.

## Rekomendasi implementasi dari discovery

- Tambah helper `hasOriginalPurchaseOrderQty(detail model.OrderDetailRead) bool` atau `hasDisplayablePurchaseOrderQty(detail model.OrderDetailRead) bool`.
- Gunakan helper hanya pada filter append `PurchaseDetails.Normal` di `DetailV2`.
- Preserve fallback existing `|| isActiveDetailForTab(detail, promoSnapshotTabSalesOrder)` kecuali test membuktikan fallback tidak diinginkan.
- Tambahkan regression test untuk case SX-2209: `qty_po1/2/3 = 0`, `original_qty_po3 = 3`, sales/final qty zero/nil, `ItemType = 1`; expected satu row di `PurchaseDetails.Normal`, original field ada, current qty tetap zero.
- Tambahkan test negative: semua current/original qty zero/null tetap tidak masuk purchase details bila sales fallback juga tidak aktif.

## Commands/docs checked

- Read docs: `AGENTS.md`, `.opencode/docs/index.md`, `.opencode/docs/ARCHITECTURE.md`, `.opencode/docs/QUALITY.md`, `.opencode/docs/AGENT_ROUTING.md`.
- Repo search: `glob **/go.mod`, `glob sales/**/*order*`, grep patterns `sales/v2/orders|purchase_details|PurchaseDetails|qty_po1|qty_po2|qty_po3|original_qty_po`, `func isActiveDetailForTab|activeQtyForTab`.
- No runtime/test command executed during planning; validation is specified in primary plan.
