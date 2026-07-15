# Discovery Evidence — SX-2291 Cancel Need Review

## Files inspected

- `sales/controller/order_controller.go`
- `sales/main.go`
- `sales/entity/order.go`
- `sales/service/order_service.go`
- `sales/service/order_type_helper.go`
- `sales/repository/order_repository.go`
- `sales/repository/stock_repository.go`
- `sales/repository/dbtransaction.go`
- `sales/service/order_service_test.go`
- `sales/repository/stock_repository_cancel_test.go`

## Route and handler

- `sales/controller/order_controller.go:36-47`
  - `OrderController.Route` registers `roRouteV1.Patch("/status", controller.UpdateStatus)` under `/v1/orders` with `middleware.JWTProtected()`.
- `sales/main.go:110-112`
  - `orderController.Route(app)` bootstraps route.
- `sales/controller/order_controller.go:894-947`
  - `UpdateStatus` parses `entity.BulkUpdateStatusOrder`, fills `UpdatedBy` from JWT `user_id`, passes `cust_id` from JWT to `OrderService.BulkUpdateStatus`.

## Payload model

- `sales/entity/order.go:387-396`
  - `UpdateDataStatusBody` has `cust_id`, `ro_no`, `data_status`, `updated_by`.
  - `BulkUpdateStatusOrder` has `orders []UpdateDataStatusBody`.

## Status constants

- `sales/entity/order.go:3-26`
  - `NEED_REVIEW = 1`
  - `PROCESSED = 2`
  - `CANCELLED = 9`
  - `dataStatusName` maps `Need Review` and `Cancelled`.

## Transition validation

- `sales/service/order_service.go:5002-5008`
  - `validateCancelTransition` already allows `entity.NEED_REVIEW` and `entity.PROCESSED`.
- `sales/service/order_service.go:5092-5107`
  - Cancel target detected by `*request.Orders[index].DataStatus == entity.CANCELLED`.
  - Already-cancelled order returns `nil`.
  - Invalid transition returns `invalid status transition from %d to %d`.

## Transaction and update path

- `sales/service/order_service.go:5091-5150`
  - `BulkUpdateStatus` wraps cancel flow in `WithinTransaction`.
- `sales/repository/order_repository.go:435-443`
  - `Update` writes `sls.order` via `Where("ro_no=? AND cust_id = ?", RoNo, custID).Updates(data)`.
- `sales/repository/dbtransaction.go:27-55`
  - Repository model uses transaction from context.

## Tenant/customer scope

- Controller passes `cust_id` from JWT local.
- `GetOrderById`, `Update`, `GetCancelStockBasis`, and cancel stock queries scope by `cust_id`.
- No `parent_cust_id` or distributor scope found in cancel flow.

## Cancel stock path

- `sales/service/order_service.go:5110-5139`
  - Fetches basis with `StockRepository.GetCancelStockBasis`.
  - Validates consistency.
  - Builds commands via `buildCancelStockWriteCommands`.
  - Calls `CancelSalesStockUpdates`.
- `sales/service/order_service.go:5011-5029`
  - `buildCancelStockWriteCommands` skips rows with `QtyOutSmallest <= 0`.
- `sales/repository/stock_repository.go:286-364`
  - `cancelStockBasisQuery` reads `inv.stock`, cancel rows, and `sls.order_detail`.
  - Filter uses `COALESCE(od.qty_final, 0) > 0`.
  - Select uses `COALESCE(od.qty_final, 0) AS qty_final`.
  - No `qty_po1`, `qty_po2`, `qty_po3`, `qty1`, `qty2`, `qty3`, `qty1_final`, `qty2_final`, `qty3_final` fallback in this query.
- `sales/repository/stock_repository.go:231-283`
  - `buildCancelStockMutations` creates `SO` and `SO-CO` stock rows and warehouse stock deltas.
- `sales/repository/stock_repository.go:376-411`
  - `CancelSalesStockUpdates` writes stock mutations.

## Existing qty helper

- `sales/service/order_service.go:2476-2485`
  - `stockDisplayQtyByPriority` priority: `qty*_final` → `qty*` → `qty_po*`.
  - Helper not used by cancel path.
- `sales/service/order_type_helper.go:56-61`
  - `takingOrderQtySource` falls back `qty*` → `qty_po*`.
  - Helper not used by cancel path.

## Existing tests

- `sales/service/order_service_test.go:534` `TestValidateCancelTransition`
- `sales/service/order_service_test.go:556` `TestBuildCancelStockWriteCommands_MultiSKU`
- `sales/service/order_service_test.go:575` `TestBulkUpdateStatus_Cancel_ConsistentBasisShouldApplyReversal`
- `sales/service/order_service_test.go:632` `TestBulkUpdateStatus_Cancel_MissingBasisShouldFailWithoutReversal`
- `sales/service/order_service_test.go:670` `TestBulkUpdateStatus_Cancel_AmbiguousBasisShouldFailWithoutReversal`
- `sales/service/order_service_test.go:723` `TestBulkUpdateStatus_Cancel_InvalidOutstandingShouldFailWithoutReversal`
- `sales/repository/stock_repository_cancel_test.go:162` `TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula`

## Likely root cause

`Need Review` transition already allowed. Failure likely comes from cancel stock basis, because `cancelStockBasisQuery` only considers `od.qty_final > 0`. Need Review orders from Purchase Order tab can have only `qty_po*` populated or no stock ledger. That makes basis missing/inconsistent and prevents `sls.order.data_status = 9`.

## Source strategy

- Local repo evidence used.
- User-provided docs used as requirements source.
- External docs/web skipped: endpoint and business rule supplied by user, implementation is repo-local.
- Browser skipped: BE bug, no UI validation required beyond API/manual DB checks.
