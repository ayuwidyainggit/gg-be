# Discovery SX-2154 order_type

## File diperiksa

- `.opencode/docs/index.md`: repo doc index.
- `.opencode/docs/ARCHITECTURE.md`: service boundary, transaksi, tenant rules.
- `.opencode/docs/QUALITY.md`: validasi per service, `sales` pakai `rtk go test ./...` dari folder `sales`.
- `sales/controller/order_controller.go`: `Create` handler parse `entity.CreateOrderBody`, validasi struct, call `ValidateOrderService.ValidateOrder`, lalu `OrderService.Store`.
- `sales/entity/order.go`: `CreateOrderBody` belum punya `order_type`.
- `sales/entity/order_detail.go`: `CreateOrderDetBody` punya `qty1/2/3`, belum punya `qty_po1/2/3`; response struct sudah punya `qty_po1/2/3`.
- `sales/model/order.go`: `model.Order` belum punya `OrderType`; GORM `Create` menyimpan model langsung ke `sls.order`.
- `sales/model/order_detail.go`: `model.OrderDetail` punya `QtyPo1/2/3`, belum punya `OriginalQtyPo1/2/3`; table name `sls.order_detail`.
- `sales/service/order_service.go`: `Store` map request ke `model.Order` dengan `structs.Automapper`; detail normal/promo juga automap ke `model.OrderDetail`, lalu `StoreDetail`.
- `sales/repository/order_repository.go`: `Store` dan `StoreDetail` pakai tx-aware `repository.model(c).Create(data)`.
- `sales/migration/sls.order/add_order_detail_po_fields.sql`: migration style `BEGIN`, `ALTER TABLE IF EXISTS`, `ADD COLUMN IF NOT EXISTS`, comments.
- `sales/migration/sls.order/add_opr_type_column.sql`: `sls.order` nullable text-ish column pattern, no enum type.
- `sales/migration/sls.order_detail/add_promo_snapshot_fields.sql`: order_detail additive fields use SQL comments, default/backfill when needed.

## Pola ditemukan

- Target service `sales`; jangan ubah `pjp-sales` atau `mobile` kecuali compile/shared copy terbukti perlu.
- Flow endpoint: `routes -> OrderController.Create -> ValidateOrderService.ValidateOrder -> OrderService.Store -> OrderRepository.Store/StoreDetail`.
- Repository write sudah tx-aware via `extractTx(ctx)`; implementation harus tetap service-layer transaction.
- GORM `Create` berarti field baru cukup ditambahkan ke model bila kolom DB ada.
- Validasi enum project memakai `validate:"omitempty,oneof=..."` untuk optional enum.
- Migration `sales` tidak punya documented migrate command; artifact perlu minta validasi SQL secara manual atau lewat runtime lokal bila tersedia.

## Reuse kandidat

- Reuse `structs.Automapper` untuk mapping `CreateOrderBody.OrderType -> model.Order.OrderType`.
- Reuse existing `CreateOrderDetBody` qty source: jika request belum punya `qty_po1/2/3`, mapping original qty untuk taking order harus ambil dari `detail.QtyPo1/2/3` setelah DTO ditambah; fallback aman bisa dari `Qty1/2/3` hanya jika schema actual mobile mengirim qty utama, namun AC minta `qty_po1/2/3`.
- Reuse test file `sales/service/order_service_test.go` untuk unit mapper/service tests.
- Reuse migration folder `sales/migration/sls.order/` atau `sales/migration/sls.order_detail/`; satu SQL additive file cukup, rollback opsional mengikuti existing rollback pattern.

## Constraints

- `order_type` nullable dan backward compatible.
- `SO` dan tanpa `order_type` tidak boleh mengubah kalkulasi, stock validation, status, promo, tax, discount, response shape.
- `O` hanya menambah persistence `order_type` dan `original_qty_po1/2/3`; jangan skip stock/validate flow kecuali ada requirement baru.
- `C` accept dan persist only; no extra canvas flow unless existing logic ditemukan.

## Risiko

- `CreateOrderDetBody` saat ini tidak punya `qty_po1/2/3`; tanpa tambah DTO, payload AC tidak akan map ke detail model.
- Test integration DB mungkin berat karena create order perlu product/outlet/discount fixtures.
- Ada duplicate-like service `pjp-sales`; issue endpoint target `sales/v1/orders`, jadi perubahan ke `pjp-sales` bisa memperbesar risiko.
- PostgreSQL `sls.order` tidak quoted di migration existing walau `order` keyword-ish; ikuti existing style `sls.order`.

## Research gate

- Local discovery: dilakukan, cukup untuk plan.
- Official docs/context7: tidak diperlukan; issue sudah memberi DB/BE docs summary dan perubahan memakai Go/GORM/Postgres pattern existing.
- GitHub: tidak diperlukan; tidak bergantung upstream.
- Web search: tidak diperlukan; tidak bergantung fakta eksternal.
- Browser/screenshot: tidak diperlukan; backend-only.
