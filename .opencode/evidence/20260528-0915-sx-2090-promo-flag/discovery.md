# Discovery SX-2090 Promo Flag

## File diperiksa
- `sales/controller/order_controller.go`: route `PATCH /sales/v1/orders/enhance/:ro_no`, handler `UpdateEnhance`, `BodyParser` ke `entity.EditOrderEnhanceBody`.
- `sales/entity/edit_order_enhance.go`: `EditSalesOrderDetail.IsProductPromotionSo *bool json:"is_product_promotion_so,omitempty"`; `AddSalesOrderDetail.IsProductPromotionSo *bool json:"is_product_promotion_so,omitempty"`.
- `sales/service/order_service.go`: `UpdateEnhance`, `normalizeEnhancePromoFlags`, `createOrderDetailFromSalesOrder`, `recomputePromoStateForTab`, `applyExplicitPromoFlagOverride`.
- `sales/repository/order_repository.go`: `UpdateDetailPartial` pakai GORM `Updates(map[string]interface{})` dengan `order_detail_id` dan `cust_id`.
- `sales/model/order_detail.go`: `OrderDetail.IsProductPromotionSo *bool gorm:"column:is_product_promotion_so" json:"is_product_promotion_so"`.
- `sales/service/order_service_test.go`: mock repo, test promo snapshot, test add sales order default final promo.
- `sales/migration/sls.order_detail/add_promo_snapshot_fields.sql`: kolom `is_product_promotion_so BOOLEAN DEFAULT FALSE`.
- `.opencode/docs/AGENT_ROUTING.md`, `.opencode/docs/ARCHITECTURE.md`, `.opencode/docs/QUALITY.md`.

## Pola project ditemukan
- Repo multi-module Go; target module `sales`.
- Layer wajib: Controller → Service → Repository → DB.
- Write dalam service transaction; repository memakai context transaction lewat `model(ctx)`.
- Test service memakai mock repository langsung di `sales/service/order_service_test.go`.
- Validasi target module: `rtk go test ./service -run <TestName>` lalu `rtk go test ./...` dari `sales`.

## Kandidat reuse
- Reuse DTO pointer bool; sudah benar untuk eksplisit `false`.
- Reuse `normalizeEnhancePromoFlags` untuk alias `is_product_promotion` ke `is_product_promotion_so`.
- Reuse `explicitPromoOverrides` dan `applyExplicitPromoFlagOverride` agar recompute promo tidak menimpa flag dari FE.
- Reuse mock `mockOrderRepository`, `mockDbtransaction`, dan pola capture `updateDetailPartialFn` / `storeDetailFn`.

## Temuan utama
- Update existing `sales_order[]` sudah memasukkan `updates["is_product_promotion_so"] = *detail.IsProductPromotionSo` saat pointer tidak nil, jadi eksplisit `false` aman di update awal.
- Recompute promo setelah update existing memakai `explicitPromoOverrides` untuk mempertahankan `false` / `true` saat snapshot menulis ulang `is_product_promotion_so`.
- Insert `add_sales_order[]` memasukkan `IsProductPromotionSo: addDetail.IsProductPromotionSo`, jadi create awal menyimpan pointer `false` bila payload mengirim `false`.
- Risiko bug tersisa: setelah insert, `UpdateEnhance` menjalankan `recomputePromoStateForTab`; untuk row baru dari `add_sales_order`, kode belum mencatat `explicitPromoOverrides[stockUpdate.RefDetId].SalesOrder`, sehingga snapshot promo dapat menimpa flag FE dengan hasil consult promo/default. `add_final_order` sudah punya pola override setelah insert; `add_sales_order` belum.
- Kolom DB aktual bernama `is_product_promotion_so`, bukan `is_product_promotion` tunggal. Jira menyebut `sls.order_detail` dan expected mapping harus pakai kolom ini.

## Perintah / docs dicek
- `rtk docker compose -f docker-compose.yml ps` dari repo root; output hanya warning compose `version` obsolete yang tertangkap tool, perlu rerun bila butuh runtime DB.
- Tidak perlu Context7/GitHub/Brave; perilaku bug lokal di Go service + GORM map update.
- Tidak perlu browser; endpoint BE non-UI.

## Constraints
- Jangan pakai token/credential real; fixture sanitized saja.
- Jangan ubah source saat planning; implementasi lewat `@orchestrator` → `@fixer`.
- Pertahankan `cust_id` filter pada write/read.
- Hindari migration kecuali DB lokal membuktikan kolom belum ada.

## Risks
- Nama expected dari Jira `is_product_promotion` ambigu; code/schema memakai `is_product_promotion_so`.
- Promo recompute bisa menulis ulang flag eksplisit bila override tidak dibawa untuk row baru.
- Integration test DB mungkin berat karena butuh seed lengkap order/promo/stock; mulai dari unit/service regression test dahulu.
