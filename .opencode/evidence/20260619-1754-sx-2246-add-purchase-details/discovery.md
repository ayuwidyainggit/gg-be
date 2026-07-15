# Discovery — SX-2246 Add Purchase Details

Task ID: `20260619-1754-sx-2246-add-purchase-details`

## Files inspected

- `AGENTS.md`
- `.opencode/docs/ARCHITECTURE.md`
- `sales/controller/order_controller.go`
- `sales/entity/edit_order_enhance.go`
- `sales/service/order_service.go`
- `sales/service/order_type_helper.go`
- `sales/service/order_stock_helper.go`
- `sales/model/order_detail.go`
- `sales/service/order_service_test.go`

## Endpoint and flow found

- Route: `sales/controller/order_controller.go:48`
  - `roRouteV1.Patch("/enhance/:ro_no", controller.UpdateEnhance)`
- Handler: `sales/controller/order_controller.go:997-1044`
  - Parses `entity.UpdateOrderParams` and `entity.EditOrderEnhanceBody`.
  - Injects `cust_id`, `parent_cust_id`, and `updated_by` from request locals.
  - Calls `OrderService.UpdateEnhance(ctx, params.RoNo, request)`.
- Service: `sales/service/order_service.go:5371-5923`
  - Runs under `service.Transaction.WithinTransaction`.
  - Updates existing `purchase_order` rows by `order_detail_id`.
  - Inserts add rows through `createOrderDetailFromPurchaseOrder`.
  - Recomputes promo/header state after changes.
  - Applies stock updates only after status decision returns `PROCESSED`.

## Current `add_purchase_details` behavior

- `entity.EditOrderEnhanceBody` already includes `AddPurchaseDetails []AddPurchaseOrderDetail` with JSON key `add_purchase_details`.
- `normalizeEnhancePromoFlags` maps `AddPurchaseDetails` into `AddPurchaseOrder` at `sales/service/order_service.go:5960-5962`.
- `UpdateEnhance` only derives `hasPurchaseOrder` from `PurchaseOrder` or `AddPurchaseOrder` after normalization, so `add_purchase_details` reaches insert path.
- Insert helper is `createOrderDetailFromPurchaseOrder` at `sales/service/order_service.go:6024-6108`.

## Gaps found

- `AddPurchaseOrderDetail` lacks fields from FE payload:
  - `unit_id1`, `unit_id2`, `unit_id3`
  - `is_product_promotion_po`
- `createOrderDetailFromPurchaseOrder` does not set:
  - `original_qty_po1`, `original_qty_po2`, `original_qty_po3`
  - `is_product_promotion_po`
  - `qty1_stok`, `qty2_stok`, `qty3_stok`
- Added purchase rows currently store requested qty directly into:
  - `qty_po1/2/3`
  - `qty1/2/3`
  - `qty1_final/2_final/3_final`
- No UOM-aware stock cap exists in purchase-add helper.

## Existing patterns to reuse

- Transaction: keep `UpdateEnhance` transaction wrapper.
- Product master lookup: `OrderRepository.FindProductByID`.
- Detail insert: `OrderRepository.StoreDetail`.
- Stock source: `StockRepository.GetCurrentStock`.
- Stock display conversion: `canonicalAPIStockBreakdown` from `order_stock_helper.go`.
- Total qty conversion: `calculateNormalizedQty`.
- Taking-order raw original qty pattern: `applyTakingOrderDetailFields` stores `OriginalQtyPo*` from raw PO payload.
- Tests: extend `sales/service/order_service_test.go` near existing `TestCreateOrderDetailFromPurchaseOrder_InheritsUOMFromProductMaster` and `TestUpdateEnhance_NormalizePurchaseDetailsAlias`.

## Constraints and risks

- Must preserve Controller → Service → Repository → DB boundaries.
- Writes must stay transaction-aware through `txCtx`.
- Must keep `cust_id` filters.
- Must not alter existing `purchase_order` update behavior except shared helper improvements for add rows.
- Stock rule must not be raw `min(qty, stock)` without UOM conversion. Use existing conversion helpers.
- `ro.WhId` and `ro.RoDate` are dereferenced in current add path. Implementation should guard or return clear error if missing before add insert.
- Duplicate prevention requirement unclear. Avoid destructive delete/replace without explicit business rule.

## Source strategy

- Local project discovery used.
- Jira and Google Docs URLs were provided but not fetched because authenticated/private access likely required and user pasted required payload/rules.
- Official docs not needed; Go/GORM/Fiber behavior not central or unfamiliar.
- GitHub/web search not needed; issue is repo-local backend logic.
- Browser/screenshot not needed; backend task.
