# SX-2184 order_type O DetailV2 purchase details evidence

Date: 2026-06-08

## Scope
- Target service: `sales`
- Endpoint behavior: `GET /sales/v2/orders/:ro_no`
- Bug: taking order (`order_type = "O"`) rows can have `qty1/qty2/qty3` nil or zero while `qty_po1/qty_po2/qty_po3` are populated. `DetailV2` previously copied `purchase_details` from Sales Order rows, so Purchase Order could be empty.

## Root cause
- `sales/service/order_service.go` built `response.Details.Normal` using Sales Order active qty (`qty1/qty2/qty3`).
- `response.PurchaseDetails.Normal` was then copied from `response.Details.Normal`.
- For taking orders, Sales Order rows can correctly be inactive while Purchase Order rows should be active from `qty_po*`.

## Implementation
- `DetailV2` now builds `PurchaseDetails` directly from persisted `details` instead of copying sales details.
- Purchase rows are included when either:
  - purchase qty is active via `qty_po1/qty_po2/qty_po3`, or
  - sales qty is active, preserving legacy/non-taking-order display behavior.
- For rows with active purchase qty, displayed purchase `qty1/qty2/qty3` mirrors `qty_po1/qty_po2/qty_po3`.
- Purchase displayed `sell_price1/2/3` maps from `sell_price_po1/2/3` when present.
- Sales Order and Final Order active-row filters remain unchanged.

## Changed files
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`

## Validation
- `rtk go test ./service -run 'TestDetailV2_PurchaseDetailsUsesPurchaseActiveRowsForOrderTypeO|TestDetailV2_MapsPromoRemarksPerTabFieldsFromPersistedSnapshot|TestDetailV2_UsesPersistedRewardProductsWhenSnapshotExists'`
  - Result: pass, 3 tests / 1 package.
- `rtk go test ./...` from `sales`
  - Result: pass, 220 tests / 22 packages.

## Runtime and local DB validation
- `rtk docker compose -f docker-compose.yml up -d`
  - Initial command timed out while RabbitMQ health was still starting, but dependent services came up.
- `rtk docker compose -f docker-compose.yml ps`
  - Redis healthy; RabbitMQ later healthy; main services up.
- `rtk docker compose -f docker-compose.yml up -d sales && rtk docker compose -f docker-compose.yml ps sales`
  - `scylla-sales` started and exposed `0.0.0.0:9004->9004/tcp`.
- Sales service logs showed Fiber running on `http://127.0.0.1:9004` and bound to `0.0.0.0:9004`.
- `GET http://localhost:9004/ping`
  - Result: `200 It works`.
- `GET http://localhost:9004/v2/orders/SO2606080006` without JWT
  - Result: `400 {"message":"Missing or malformed JWT","request_id":""}`; endpoint route is alive but protected.
- Login through local system API with the provided test credentials, then call `GET http://localhost:9004/v2/orders/SO2606080006` using the returned bearer token. Token value was not printed or stored.
  - Result: endpoint returned `ro_no = SO2606080006`.
  - `details.normal = 0`.
  - `details_final.normal = 0`.
  - `purchase_details.normal = 1`.
  - Purchase row: `order_detail_id = 7182`, `pro_code = TP-007`, `pro_name = Topi Naga`, `qty1/qty2/qty3 = 3/0/0`, `qty_po1/qty_po2/qty_po3 = 3/0/0`, `sell_price1 = 1200000`, `sell_price_po1 = 1200000`.
- Local DB connectivity:
  - `/opt/homebrew/opt/postgresql@18/bin/psql -h localhost -U postgres -d ggn_scyllax -Atc "select current_database(), current_user"`
  - Result: `ggn_scyllax|postgres`.
- Local DB fixture validation for `SO2606080006` / `C260020001`:
  - Header: `order_type = O`, `opr_type = O`, `validate_stok = false`, `validate_stok_message = NULL/blank`, `detail_count = 1`.
  - Detail row `7182`: product `TP-007 / Topi Naga`, `qty = 0`, `qty1/qty2/qty3 = NULL`, `qty_final = 0`, `qty1_final/qty2_final/qty3_final = NULL`, `qty_po = 3`, `qty_po1 = 3`, `qty_po2 = 0`, `qty_po3 = 0`, `original_qty_po1 = 3`, `original_qty_po2 = 0`, `original_qty_po3 = 0`.
  - Tab activity query result: `sales_active_rows = 0`, `final_active_rows = 0`, `purchase_active_rows = 1`.
  - Expected Purchase display fields from DB: `purchase qty1/2/3 = 3/0/0`, `sell_price_po1 = 1200000.0000`; raw Sales/Final qty fields are NULL.
- Legacy comparison from local DB:
  - `SO2606080001` / `C260020001`: `sales_active_rows = 1`, `final_active_rows = 1`, `purchase_active_rows = 0`, validating why the new service fallback still includes sales-active rows in `purchase_details` for legacy/non-taking-order data.

## Quality gate
- Final `@quality-gate`: PASS.
- Non-blocking follow-up: run live smoke against a known `order_type = "O"` fixture when compose services are available.

## Notes
- Root repo path is not a git repo; `sales` is the service git worktree on `bugfix/SX-2184-dev`.
- No commit or push was performed in this run.
