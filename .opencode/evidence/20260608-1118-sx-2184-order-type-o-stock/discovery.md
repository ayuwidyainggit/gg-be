# Discovery SX-2184 — order_type O skip stock validation

Task ID: `20260608-1118-sx-2184-order-type-o-stock`
Tanggal: 2026-06-08 Asia/Jakarta
Mode: Maintenance Stability Mode

## Sumber yang dicek

- Repo lokal: `/Users/ujang/Projects/Geekgarden/scylla-be`
- Harness: `AGENTS.md`, `.opencode/docs/index.md`, `.opencode/docs/ARCHITECTURE.md`, `.opencode/docs/QUALITY.md`
- Plan terkait: `.opencode/plans/20260604-1024-sx-2154-order-type.md`
- Source target `sales`:
  - `sales/controller/order_controller.go`
  - `sales/entity/order.go`
  - `sales/entity/order_detail.go`
  - `sales/entity/validate_order.go`
  - `sales/model/order.go`
  - `sales/model/order_detail.go`
  - `sales/repository/order_repository.go`
  - `sales/repository/stock_repository.go`
  - `sales/repository/validate_order_repository.go`
  - `sales/service/order_service.go`
  - `sales/service/validate_order_service.go`
  - `sales/service/order_stock_helper.go`
  - `sales/service/order_service_test.go`
  - `sales/migration/sls.order/add_order_detail_po_fields.sql`
- External docs Google Docs dicoba via `webfetch`; doc pertama hanya menampilkan halaman login/JavaScript, doc kedua terlalu besar. Isi dokumen yang dipakai adalah ringkasan user dalam prompt, bukan ekstraksi authenticated docs.
- Runtime preflight: `rtk docker compose -f docker-compose.yml ps` dijalankan dari repo root; hasil kosong service running dan ada warning `.rtk/filters.toml` untrusted serta `version` obsolete.

## Pola repo yang ditemukan

- Repo multi-module Go; target fix ada di module `sales`.
- Repo rule: Controller → Service → Repository → DB.
- Validasi/test harus dari folder `sales`, dengan command `rtk go test ./...` dan targeted tests.
- Write DB harus dalam service-layer transaction dan repository harus tx-aware lewat `extractTx`.
- `POST /sales/v1/orders` didaftarkan di `OrderController.Route`: `app.Group("/v1/orders", middleware.JWTProtected()).Post("", controller.Create)`; path efektif di service compose adalah `/sales/v1/orders` bila gateway/base prefix menambahkan `/sales`.

## Temuan implementasi saat ini

1. `sales/controller/order_controller.go:67-124`:
   - `Create` parse `entity.CreateOrderBody`.
   - Controller selalu membuat `entity.ValidateOrderBody`, map `request` dan `request.Details.Normal`, lalu memanggil `ValidateOrderService.ValidateOrder(validateOrderRequest)` sebelum `OrderService.Store`.
   - Karena validasi terjadi di controller sebelum service store, branch `order_type = O` di service tidak akan cukup untuk mencegah error stock jika controller tetap memanggil `ValidateOrder`.

2. `sales/service/validate_order_service.go:51-109`:
   - `ValidateOrder` mengambil `inv.warehouse_stock` via `GetWarehouseStockByProducts`.
   - Validasi stock membandingkan `totalQty` dengan `availableStock`.
   - Jika qty > available stock, `Validate1Success = false` dan message `Insufficient Stock`.
   - Ini root-cause candidate untuk SX-2184 karena tidak ada input/branch `order_type`.

3. `sales/service/order_service.go:281-286`:
   - `Store` menentukan `DataStatus` dari `determineSalesOrderStatus(validateOrderRequest, outletRules...)` dan menerapkan validation result ke order model.
   - Untuk taking order `O`, `validate_stok` harus false dan message nil/null, sehingga service perlu menerima validation response netral/stock-bypassed atau melakukan branch sebelum status decision.

4. `sales/service/order_service.go:377-390`, `450-463`, `467-472`:
   - Stock mutation (`SalesStockUpdates`) dikumpulkan untuk normal dan promo detail jika `isProcessedDataStatus(orderModel.DataStatus)` lalu dieksekusi setelah detail store.
   - Untuk `order_type = O`, mutation ini harus tidak dibuat/ tidak dieksekusi meskipun status decision menjadi processed; lebih aman gate dengan helper `shouldMutateInventoryOnCreate(orderType)`.

5. `sales/repository/stock_repository.go:469-475` dan `481-500`:
   - `SalesStockUpdates` melakukan upsert `inv.warehouse_stock` dan insert bulk `inv.stock`.
   - Jadi bypass inventory untuk `O` harus terjadi sebelum memanggil `SalesStockUpdates`.

6. `sales/entity/order.go:57-101` dan `sales/model/order.go:11-81`:
   - Pada working tree yang dibaca, `CreateOrderBody` belum punya `OrderType` dan `model.Order` belum punya `OrderType`.
   - Plan lama SX-2154 sudah menargetkan penambahan field ini, tetapi source lokal belum menunjukkan field tersebut.

7. `sales/entity/order_detail.go:3-91` dan `sales/model/order_detail.go:7-106`:
   - `CreateOrderDetBody` punya `QtyPo` tapi belum punya `QtyPo1/2/3` dan belum punya `OriginalQtyPo1/2/3`.
   - `model.OrderDetail` punya `QtyPo1/2/3`, `SellPricePo1/2/3`, promo PO fields, `DiscPo`, `VatValuePo`; belum punya `OriginalQtyPo1/2/3`.

8. `sales/migration/sls.order/add_order_detail_po_fields.sql`:
   - Existing PO columns `qty_po1/2/3` bertipe `FLOAT4 DEFAULT 0`.
   - Belum terlihat migration `order_type` atau `original_qty_po1/2/3` dari grep SQL.

9. `sales/service/order_service.go:310-327`:
   - Create detail menghitung `Qty`, `QtyPo`, `QtyFinal` dari request `Qty1/Qty2/Qty3`, lalu `Automapper` detail ke model.
   - Untuk `order_type = O`, docs user mengharapkan sales-order qty (`Qty`, `Qty1/2/3`, `QtyFinal`) null/belum terisi sampai process order, tetapi existing code saat ini mengisi semuanya. Ini perlu difix dengan branch O yang eksplisit, namun harus hati-hati agar promo/discount tetap tidak rusak.

10. `sales/service/order_service.go:4090-4108`:
    - Update flow existing memakai `ro.DataSource == 2` untuk menganggap mobile purchase order dan map `qty1/2/3` ke `qty_po1/2/3`.
    - Untuk create SX-2184, jangan mengandalkan `data_source` lagi; contract harus `order_type == "O"`.

## Reuse candidates

- Reuse helper konversi qty: `conversion.QtyUnit.ToTotalQuantity()`.
- Reuse tx-aware repositories: `OrderRepository.Store`, `StoreDetail`, `StockRepository.SalesStockUpdates`.
- Reuse `applyValidationResultToOrderModel`, tetapi perlu stock-bypassed validation result atau explicit override untuk `O`.
- Reuse existing test mocks di `sales/service/order_service_test.go`; tambah method mock untuk `Store`, `CountAllRoByCustId`, promo/discount funcs bila test `Store` dibuat unit-level.
- Reuse migration style idempotent: `BEGIN; ALTER TABLE IF EXISTS ... ADD COLUMN IF NOT EXISTS ... COMMENT ... COMMIT;`.

## Constraints dan risiko

- Jangan copy token/auth header dari docs.
- `order_type` kosong/null harus backward compatible dan tidak boleh bypass stock.
- `SO` tidak boleh bypass stock validation atau inventory mutation.
- `C` belum punya requirement baru; pertahankan existing/as-is.
- Google Docs/Jira tidak dapat diverifikasi tanpa auth; gunakan prompt user sebagai reference summary.
- `sales` migration command tidak terdokumentasi; SQL harus idempotent dan manual apply lokal perlu dicatat.
- Runtime compose saat dicek tidak ada service running, jadi API smoke kemungkinan perlu `rtk docker compose -f docker-compose.yml up -d` dulu oleh executor.

## Root cause hypothesis

Defect kemungkinan terjadi karena `OrderController.Create` selalu memanggil `ValidateOrderService.ValidateOrder` sebelum `OrderService.Store`, sedangkan `ValidateOrder` selalu menjalankan warehouse stock validation tanpa awareness `order_type`. Akibatnya request `order_type = "O"` dengan qty > wh stock gagal sebelum create order bisa diproses sebagai taking/purchase order. Selain itu, create store masih punya jalur `SalesStockUpdates` yang dapat menulis `inv.stock` dan update `inv.warehouse_stock` bila status processed, sehingga fix harus mencakup controller validation bypass dan inventory mutation bypass.
