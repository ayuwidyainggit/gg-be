# Discovery Evidence — SX-1878/SX-1879 PO Final Export

## Files inspected

- `AGENTS.md` at repo root: multi-module Go repo, strict Controller → Service → Repository → DB flow, tenant `cust_id` filter, transaction requirement for writes, project conversion utilities.
- `sales/controller/order_controller.go`: route mapping for `PATCH /sales/v1/orders/enhance/:ro_no`, `POST /sales/v1/orders/conversion`, `GET /sales/v2/orders/:ro_no`, `PATCH /sales/v1/orders/final/:ro_no`.
- `sales/controller/validate_order_controller.go`: route mapping for `POST /sales/v1/validate-order/` and `/detail`.
- `sales/controller/so_controller.go`: route mapping and query parsing for `GET /sales/v1/download`, including `salesman_id[]` support.
- `sales/service/order_service.go`: detail V2 stock/qty projection, `UpdateEnhance`, `ProcessEnhanceWithoutProductEdit`, conversion helpers, stock update calls, header recompute.
- `sales/repository/order_repository.go`: `UpdateDetailPartial`, warehouse stock lookup, `SyncFinalOrderFields`.
- `sales/repository/stock_repository.go`: `GetCurrentStock` exists for current warehouse stock.
- `sales/service/so_service.go`: async export generation, amount formatter, PO/SO/Final sheet builders, `mapPoToEntity`, `mapSoToEntity`, `mapFinalToEntity`.
- `sales/repository/so_repository.go`: export queries `FindDownloadDataPo`, `FindDownloadDataSo`, `FindDownloadDataFinal`.
- `sales/service/order_service_test.go`: mocks and existing tests for `UpdateEnhance`, stock update deltas, `DetailV2` tab behavior.
- `sales/service/so_service_test.go`: existing SX-1879 formatter/mapping/sheet tests.
- `.opencode/plans/20260504-2141-sx-1879-export-issue-data.md`: prior plan for SX-1879; current code appears to already include several suggested fixes/tests.

## Project patterns found

- `sales` is an independent Go module; run tests from `sales/`.
- Endpoint flow follows controller → service → repository.
- Writes in `UpdateEnhance` already run inside `service.Transaction.WithinTransaction`.
- Existing utility `conversion.QtyUnit` converts qty1/qty2/qty3 into smallest-unit/base quantity, and `conversion.Qty` converts base quantity back into unit breakdown.
- `UpdateEnhance` currently cascades PO qty edits to `qty*`, `qty*_final`, and stock update deltas.
- `DetailV2` builds three surfaces from the same `sls.order_detail` rows: Sales Order (`Details`), Purchase Order (`PurchaseDetails`), Final Order (`DetailsFinal`).
- Export already has `formatDownloadAmount` with Indonesian thousands separator and tests in `so_service_test.go`.

## Reuse candidates

- Reuse `calculateNormalizedQty` and `conversion.QtyUnit` for multi-unit base quantity conversion; avoid creating a new conversion implementation unless existing helper cannot normalize/clamp cleanly.
- Reuse `StockRepository.GetCurrentStock` for BE defensive stock clamp in `UpdateEnhance` and new add-detail paths.
- Reuse `OrderRepository.UpdateDetailPartial` for persisting clamped PO/SO/Final qty fields.
- Reuse `SalesStockUpdates` mechanism after status is `PROCESSED`; ensure `QtyOrder` reflects clamped/processed qty.
- Reuse existing test mocks in `order_service_test.go` and `so_service_test.go`.
- Reuse existing SX-1879 tests and update only gaps/regressions.

## Commands/docs checked

- `docker compose -f docker-compose.yml ps` from repo root showed `system`, `master`, and `finance` up; `sales` was not listed as running.
- `go test ./service` from `sales/` passed: `ok sales/service 0.604s`.
- Local project discovery was sufficient; no official docs, GitHub, web, or browser research was required for this backend-only defect plan.

## Important code observations

### SX-1878

- `UpdateEnhance` PO path (`order_service.go:5224+`) calculates requested qty via `calculateNormalizedQty` and then writes the same result into:
  - `qty_po`, `qty_po1`, `qty_po2`, `qty_po3`
  - `qty`, `qty1`, `qty2`, `qty3`
  - `qty_final`, `qty1_final`, `qty2_final`, `qty3_final`
- That PO path currently does not visibly clamp requested qty against current available stock. It trusts the payload for `totalQty` and stock update `QtyOrder`.
- `DetailV2` Sales Order stock display adds current warehouse stock plus `detailData.Qty*` for non-cancelled orders.
- `DetailV2` Final Order stock display separately adds current warehouse stock plus final converted qty (`detailData.Qty*` after converting `QtyFinal`). This can diverge if final fields are stale or if the Final Order row reads fallback/original qty.
- `SyncFinalOrderFields` copies `qty*` to `qty*_final` but is not called in `UpdateEnhance` PO path; PO path directly includes final updates in the same map.
- `ProcessEnhanceWithoutProductEdit` only updates header status; it does not refresh stock snapshot, clamp, or sync detail qty fields. It may be valid for no-edit scenario, but it is a risk if final detail fields were stale before processing.

### SX-1879

- `so_service.go` already includes `formatDownloadAmount` and applies it to amount/selling-price cells.
- `mapPoToEntity`, `mapSoToEntity`, and `mapFinalToEntity` currently calculate gross from sheet-specific qty * price, set discount from `DiscValueFinal`, VAT from `VatValueFinal`, and gross equal to grossSales.
- Existing SX-1879 tests already cover Indonesian separator, PO/SO/Final financial mapping, PO number fallback/filtering, and sheet generation. The service package tests pass locally.
- Export repository filters all three sheet queries by `sls.order.ro_date` and `salesman_id IN ?`, and tenant filters are present.
- Purchase Order export data is filtered in service by `filterDownloadDataPoWithPONumber`, which excludes rows with blank `po_no` and falls back to `order_no` only for display when `po_no` is present enough to pass the filter. This may explain PO sheet null/empty if data exists but `po_no` is blank.

## Constraints

- User requested implementation, but this agent is artifact planner; only `.opencode/` artifacts are edited here.
- Do not commit or store Jira authorization tokens/credentials.
- Repo-level AGENTS requires `rtk`, but global OpenCode instructions say not to prefix commands with `rtk`; current execution used direct commands per global OpenCode instruction. Record this conflict for implementer.
- `sales` container is not up; manual endpoint reproduction requires starting it and valid env/auth.
- Staging Jira URLs may require credentials not available here.

## Risks

- BE stock clamp needs clear definition of “available stock”: current warehouse stock alone versus current stock plus previously allocated qty for the same order detail. For already processed order edits, clamp should likely allow the previous allocation to be reused, otherwise lowering/raising within the original reservation may be miscalculated.
- Multi-unit clamp must convert requested qty and available stock into the same base unit before comparing; then convert clamped base qty back to qty1/qty2/qty3.
- If `qty_change* == 0` means “no delta” from FE, BE must not use it as proof that final qty is unchanged; `purchase_order.qty_po*` should be source of truth for process payload.
- Export PO blank filtering may have been intentional to hide rows without PO number. Changing it could add rows to PO sheet; test and product/QA expectation should be explicit.
- Manual reproduction against staging may expose sensitive data. Do not paste tokens or full private export data into artifacts.
