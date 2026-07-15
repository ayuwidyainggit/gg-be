# Discovery SX-2250 — Reward Product Stock Display

Task id: `20260619-1148-sx-2250-reward-stock`
Target: `sales/` service

## Files inspected

- `.opencode/docs/index.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`
- `sales/service/order_stock_helper.go`
- `sales/service/order_stock_helper_test.go`
- `sales/pkg/conversion/quantity.go`
- `sales/pkg/conversion/qtyunit.go`
- `.opencode/plans/20260619-0831-sx-2253-available-stock-row.md`

## Project patterns found

- Repo adalah multi-module Go monorepo; validasi wajib dari target service directory `sales/`.
- `DetailV2` berada di `sales/service/order_service.go:2856`.
- Route `GET /sales/v2/orders/:ro_no` di `sales/controller/order_controller.go` memanggil `OrderService.DetailV2`.
- `DetailV2` membangun `response.Details.Normal`, `response.PurchaseDetails.Normal`, dan `response.DetailsFinal.Normal` dari `OrderRepository.FindDetail`.
- Reward/promo row `item_type != 1` awalnya masuk `Promo`, lalu `movePromoDetailsToNormal` memindahkan promo rows ke `.Normal` dan mengosongkan `.Promo`.
- Display stock helper sekarang ada di `sales/service/order_stock_helper.go`:
  - `computeDisplayedAvailableStockBreakdown`
  - `canonicalAPIStockBreakdown`
  - `applyStockBreakdownToPointers`
- Existing tests terkait:
  - `TestDetailV2_SameSKURewardDoesNotContaminateNormalRow`
  - `TestDetailV2_Cancelled_UsesWarehouseCurrentOnlyForDisplayedStock`
  - `TestDetailV2_NonCancelled_KeepsExistingDisplayedStockBehavior`
  - `order_stock_helper_test.go`

## Reuse candidates

- Reuse `computeDisplayedAvailableStockBreakdown` untuk formula `warehouse current + current row qty`.
- Reuse `mockOrderRepositoryDetailV2` untuk unit test DetailV2 tanpa staging token.
- Reuse `movePromoDetailsToNormal` behavior; jangan ubah row separation/flag semantics kecuali test membuktikan bug.
- Reuse plan SX-2253 invariant: row stock display tidak boleh aggregate qty semua row dengan `pro_id` sama.

## Key code evidence

- `DetailV2` sales-order stock display:
  - `sales/service/order_service.go:2946-2964`
  - `warehouseStockMap[int64(detail.ProId)] + current row qty`
- `DetailV2` final-order stock display:
  - `sales/service/order_service.go:3216-3234`
  - sama, memakai current row final qty.
- Promo row masuk `.Normal` setelah:
  - `sales/service/order_service.go:2467-2474` (`movePromoDetailsToNormal`)
  - `sales/service/order_service.go:2997`, `3257`
- Current helper already row-scoped; risiko utama adalah kurang test eksplisit untuk Jira SX-2250: same `pro_id`, normal row + reward row, qty sama/beda, assert `details.normal` dan `details_final.normal`.

## Commands/docs checked

- Local docs: `.opencode/docs/index.md`, `ARCHITECTURE.md`, `QUALITY.md`.
- Search evidence via content search for `DetailV2`, `qty*_stok`, promo flags, promo aggregate functions.
- No Jira/browser fetch used; prompt already includes sanitized Jira details. Credentials/tokens from Jira intentionally not accessed or copied.

## Constraints

- No source implementation done by planner.
- No staging credentials/token copied into artifacts.
- `rtk` prefix required for repo commands.
- Sales service validation must run from `sales/`.

## Risks

- Code appears already row-scoped in current branch. Implementation may be test-only if Red test already passes. Executor must run test before changing logic and avoid unnecessary refactor.
- On-customer stock movement not yet traced deeply in this discovery. Latest Jira says movement evidence already expected; only inspect/fix if failing test or code evidence shows reward rows excluded.
- `computeDisplayedAvailableStockBreakdown` argument names are confusing relative to conversion output; avoid renaming/changing semantics unless covered by existing tests.
