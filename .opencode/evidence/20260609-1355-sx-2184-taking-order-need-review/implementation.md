# SX-2184 comment 16637 taking order Need Review evidence

Date: 2026-06-09

## Scope
- Target service: `sales`
- Endpoint: `POST /sales/v1/orders`
- Jira feedback: SX-2184 focused comment 16637
- Requirement: create order with `order_type = "O"` must persist `sls.order.data_status = 1` / Need Review, while preserving prior SX-2184 stock bypass behavior and non-`O` behavior.

## Root cause
- In create flow, `request.data_status` is automapped into `orderModel`, but then always overwritten by `determineSalesOrderStatus(...)` in `sales/service/order_service.go`.
- For taking orders, stock validation is skipped, so the normal status resolver can return `PROCESSED` (`2`) when non-stock validations pass.
- There was no explicit taking-order create status rule before this patch.

## Implementation
- Added `resolveCreateOrderDataStatus(orderType, statusDecision)` in `sales/service/order_type_helper.go`.
- For `order_type = "O"`, it returns `entity.NEED_REVIEW` (`1`).
- For all other order types, it returns the existing status decision unchanged.
- `sales/service/order_service.go` now uses this helper when assigning `orderModel.DataStatus` during create.
- Existing taking-order logic remains intact:
  - `opr_type = O`
  - `validate_stok = false`
  - `validate_stok_message = nil`
  - stock mutation gated by `ShouldMutateInventoryOnCreate(...)`
  - `qty_po*` / `original_qty_po*` persisted by `applyTakingOrderDetailFields(...)`

## Changed files
- `sales/service/order_service.go`
- `sales/service/order_type_helper.go`
- `sales/service/order_type_helper_test.go`

## Tests
- Targeted:
  - `rtk go test ./service -run 'TestResolveCreateOrderDataStatusForTakingOrder|TestStore_SX2184TakingOrderSkipsStockMutationAndPersistsOriginalQty|TestStore_SX2184NilOrderTypeStillMutatesInventory|TestStore_DeterminesProcessedWhenValidationPassesRegardlessOfPayloadStatus|TestStore_DeterminesNeedReviewForRestrictedCreditLimit'`
  - Result: pass, 8 tests/subtests.
- Full service suite:
  - `rtk go test ./...`
  - Result: pass, 225 tests / 22 packages.

## Local runtime/API/DB validation
- Runtime:
  - `rtk docker compose -f docker-compose.yml ps sales system rabbitmq redis`
  - `sales`, `system`, `rabbitmq`, and `redis` up; RabbitMQ/Redis healthy.
  - `rtk docker compose -f docker-compose.yml restart sales`
  - `GET http://localhost:9004/ping` returned `200 It works`.
- API create:
  - Login through local system API with the provided test credentials, then POST to `http://localhost:9004/v1/orders` using the returned bearer token. Token value was not printed or stored.
  - Payload used sanitized SX-2184 feedback 16637 values with `order_type = "O"`, `data_status = 1`, `wh_id = 350`, `pro_id = 10812`, `qty1 = 1`, stock qty fields `0`.
  - Result: `Created Successfully`, `ro_no = SO2606090001`.
- DB validation for `SO2606090001` / `C260020001`:
  - Header: `order_type = O`, `opr_type = O`, `data_status = 1`, `data_source = 1`, `validate_stok = false`, `validate_stok_message = NULL/blank`.
  - Detail: product `TP-012 / Topi Badut`, `qty = 0`, `qty1/qty2/qty3 = NULL`, `qty_final = 0`, `qty_po = 1`, `qty_po1 = 1`, `qty_po2 = 0`, `qty_po3 = 0`, `original_qty_po1 = 1`, `original_qty_po2 = 0`, `original_qty_po3 = 0`.
  - `inv.stock` rows for the created order: `0`.
  - `inv.warehouse_stock` row for `wh_id = 350`, `pro_id = 10812`, `cust_id = C260020001`: none before/after, so no warehouse stock row was created or updated.
- API detail validation:
  - `GET http://localhost:9004/v2/orders/SO2606090001`
  - Result: `data_status = 1`, `data_status_name = Need Review`, `details.normal = 0`, `details_final.normal = 0`, `purchase_details.normal = 1`.
  - Purchase row: `TP-012 / Topi Badut`, `qty1/qty2/qty3 = 1/0/0`, `qty_po1/qty_po2/qty_po3 = 1/0/0`.

## Quality gate
- Final `@quality-gate`: PASS.
- Required remediations: none.

## Security
- The Jira token/header was not copied into code, tests, evidence, logs, or commit.
- Runtime validation used a locally acquired token from the provided test credentials; token value was not printed or stored.

## Notes
- No commit was made in this run yet.
