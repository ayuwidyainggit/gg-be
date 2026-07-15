# SX-2508 discovery

## Source strategy

- **Repo local, used:** `sales` module code, tests, migration, prior import plans.
- **Stack evidence, used:** `sales/go.mod`, `sales/Dockerfile`, `sales/.air.toml`; no new dependency/API planned.
- **External docs, skipped:** no version-sensitive API change; existing Go/Fiber/GORM patterns cover fix.
- **GitHub/web/browser, skipped:** ticket supplies behavior; no upstream/reference dependency. Staging curl is execution evidence, not planning evidence; token intentionally absent.
- **Harness stack docs:** `.opencode/docs/PROJECT_STACK.md`, `PROJECT_COMMANDS.md`, `FRAMEWORK_PLAYBOOK.md`, `PROJECT_DETECTED_TOOLS.md` absent. `/init-harness` command unavailable in lane; plan uses `AGENTS.md`, `ARCHITECTURE.md`, module manifest, Dockerfile, existing tests. Claim level lowered where exact runtime state remains unverified.

## Files inspected

- `AGENTS.md`
- `.opencode/docs/MCP.md`, `.opencode/docs/ARCHITECTURE.md`
- `sales/go.mod`, `sales/Dockerfile`, `sales/.air.toml`, `sales/PROJECT_MEMORY.md`
- `sales/controller/order_controller.go`, `sales/controller/validate_order_controller.go`
- `sales/service/order_service.go`, `sales/service/order_stock_helper.go`
- `sales/pkg/conversion/quantity.go`
- `sales/repository/order_repository.go`, `sales/repository/stock_repository.go`
- `sales/entity/order.go`, `sales/entity/order_detail.go`
- `sales/model/order.go`, `sales/model/order_detail.go`
- `sales/service/order_service_test.go`, `sales/service/order_import_parser_test.go`, `sales/service/order_stock_helper_test.go`, `sales/controller/order_controller_test.go`
- `.opencode/plans/20260706-sx-2435-2451-sales-order-import.md`
- `.opencode/plans/20260709-sx-2499-import-feedback.md`

## Confirmed findings

1. `POST /v1/orders/import` and URL-import both route to `OrderService.ImportOrders`; `GET /v2/orders/:ro_no` routes to `OrderService.DetailV2`. `confirmed_repo`.
2. Import parser currently marks imported request `IsSalesMapping=true`, `InvoiceNo=documentNo`, `InvoiceDate=documentDate`, and copies row PPN into `VatValue`/`VatValueFinal`; it does **not** assign `PromoSo1`/`PromoFinal1`, sets system prices 2/3 from parent `SellPrice2`/`SellPrice3`, and does not make `QtyPo` explicitly nil. `confirmed_repo`.
3. `ConvToQtyConversion()` emits `Qty1=small`, `Qty2=middle`, `Qty3=large`. `canonicalAPIStockBreakdown()` already explicitly remaps this to API `Qty1=large`, `Qty2=middle`, `Qty3=small`. `confirmed_repo`.
4. `DetailV2` always recomputes detail quantities from smallest `Qty`, then uses `stockDisplayQtyByPriority(detail)` plus `computeDisplayedAvailableStockBreakdown`. `confirmed_repo`.
5. `order_stock_helper.go` provides smallest-unit conversion and canonical API breakdown. Reuse, not new global converter. `confirmed_repo`.
6. Existing import unit lookup checks `product.UnitId1..5`; imported detail persists parent product units/conversions. Parent fetch currently uses `FindProductByID`; whether that read returns latest edited distributor mapping for every tenant is `unverified` until repository SQL/test inspection during execution.
7. Writes are service-owned and repository model extraction is transaction-aware. `confirmed_repo`.
8. No material dependency/API change planned; `sales/go.mod` pins Go 1.23 + toolchain 1.24.6, Fiber 2.52.6, GORM 1.24.7. `confirmed_repo`.

## Risks

- Do not alter global `ConvToQtyConversion()`; callers retain internal small/middle/large contract.
- `DetailV2` dereferences mapped conversions without nil guard before fix. Test fixture must populate valid values; any new guard must preserve existing fallback semantics.
- Product mapping freshness needs SQL-level proof. Do not claim cache bug fixed until test proves updated mapping read after update.
- Import loops `Store` per document. No atomicity redesign in scope.
- Staging must use fresh shell token only. Do not write token, headers, or sensitive response bodies into artifacts.

## Confirmed vs Assumed Audit

| Claim | Level | Evidence |
| --- | --- | --- |
| DetailV2 double-converts stored mapped triples | confirmed_repo | `sales/service/order_service.go:3071-3092` |
| Global conversion helper is internally small/middle/large | confirmed_repo | `sales/pkg/conversion/quantity.go:15-23` |
| Stock helper exposes canonical large/middle/small boundary | confirmed_repo | `sales/service/order_stock_helper.go:21-63` |
| Import receives current mapping on each new file | unverified | execution must inspect query and add regression test |
| Required business outcomes | user_confirmed | SX-2508 request |
| Existing staging records reproduce all listed values | unverified | no token/runtime access used during planning |
