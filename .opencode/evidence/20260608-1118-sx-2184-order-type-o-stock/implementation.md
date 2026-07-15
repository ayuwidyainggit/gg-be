# SX-2184 Implementation Evidence

Task ID: `20260608-1118-sx-2184-order-type-o-stock`
Date: 2026-06-08 Asia/Jakarta
Mode: Maintenance Stability Mode

## References used

- Plan source of truth: `.opencode/plans/20260608-1118-sx-2184-order-type-o-stock.md`
- Discovery evidence: `.opencode/evidence/20260608-1118-sx-2184-order-type-o-stock/discovery.md`
- Repo architecture rule: `.opencode/docs/ARCHITECTURE.md`
- Repo quality rule: `.opencode/docs/QUALITY.md`
- Current branch inspection: `sales/controller/order_controller.go`, `sales/service/validate_order_service.go`, `sales/service/order_service.go`, `sales/service/order_status_helper.go`, `sales/entity/order.go`, `sales/model/order.go`

## Branch and schema inspection summary

- Current `sales` branch is `dev`.
- Runtime preflight command `rtk docker compose -f docker-compose.yml ps` showed no running compose services, so live API/DB smoke was not attempted during this remediation.
- Current working tree already contained the initial SX-2184 implementation, including taking-order helpers, `order_type` support, and additive migration files.
- Remaining quality-gate gap was not schema shape anymore; it was behavior and persistence semantics:
  - `order_type = "O"` bypassed all order validations instead of only stock validation.
  - `validate_stok_message` still persisted as empty string compatibility value instead of SQL NULL.
  - Regression coverage for nil/empty/`C` create path expectations was incomplete.

## Root cause

Primary root cause for the quality-gate failure was the initial remediation strategy in create flow:

- `sales/controller/order_controller.go` skipped `ValidateOrderService.ValidateOrder` entirely for `order_type = "O"`.
- Because `sales/service/validate_order_service.go` coupled stock validation with credit-limit / overdue / outstanding validation, the controller bypass also skipped non-stock validations.
- That changed business semantics for `O` beyond the accepted scope.

A second root cause affected persistence semantics:

- `sales/service/order_type_helper.go` forced `orderModel.ValidateStokMessage = ""` for taking order.
- `sales/model/order.go` and read-path structs still represented `validate_stok_message` as non-nullable `string`, so the initial implementation chose empty-string compatibility instead of true SQL NULL.

## Changes implemented in remediation

### 1) Preserved non-stock validations for `order_type = O`

Updated `sales/service/validate_order_service.go`:

- Added `ValidateOrderWithoutStock(dataFilter entity.ValidateOrderBody)` to the service interface.
- Extracted shared validation flow into internal helper `validateOrder(dataFilter, includeStockValidation bool)`.
- When `includeStockValidation` is `false`, service now:
  - skips only `GetWarehouseStockByProducts` + stock loop
  - still runs AR, credit-limit, overdue, and outstanding validations
  - still computes `IsSuccessValidate` from the resulting non-stock validations
  - keeps stock snapshot neutral (`Validate1Success=true`, `Validate1="Sufficient Stock"`)

Updated `sales/controller/order_controller.go`:

- `order_type = "O"` now calls `ValidateOrderService.ValidateOrderWithoutStock(...)` instead of bypassing validation entirely.
- nil/empty/`SO` continue using the original `ValidateOrder(...)` path.
- Empty-string `order_type` is normalized back to `nil` before validation, preserving backward compatibility with old clients that may send `""`.

Result:

- `O` no longer fails on stock validation.
- `O` still respects credit-limit / overdue / outstanding validation semantics.
- Nil/empty/`SO` stay on the original validation path.

### 2) Persist `validate_stok_message` as nullable message

Updated nullable shape across create/read models:

- `sales/model/order.go`
  - `Order.ValidateStokMessage` changed from `string` to `*string`
  - `OrderList.ValidateStokMessage` changed from `string` to `*string`
- `sales/entity/order.go`
  - `OrderResponse.ValidateStokMessage` changed from `string` to `*string`

Updated helper/read logic:

- `sales/service/order_type_helper.go`
  - added `nullableValidationMessage(message string) *string`
  - taking order snapshot now sets `ValidateStokMessage = nil`
- `sales/service/order_status_helper.go`
  - `applyValidationResultToOrderModel` now stores `Validate1` via nullable helper
  - `validationResultFromOrderList` converts nullable DB value back to string safely with `stringFromPtr`
  - `hasStoredValidationResult` now treats non-nil `ValidateStokMessage` as stored evidence

Result:

- Taking order can persist SQL NULL for `validate_stok_message` through normal GORM create behavior.
- Existing read/status logic remains compile-safe and nil-tolerant.

### 3) Keep inventory mutation gate for `O`

Retained prior bounded behavior in `sales/service/order_service.go`:

- create stock updates are appended only when `ShouldMutateInventoryOnCreate(orderType)` is true and order status is processed
- `O` still skips `StockRepository.SalesStockUpdates`
- nil/empty/`SO` still follow existing mutation path

### 4) Expanded regression coverage

#### Controller coverage

File: `sales/controller/order_controller_test.go`

Added/updated tests:

1. `TestOrderControllerCreate_SX2184TakingOrderBypassesStockValidation`
   - verifies `ValidateOrder` is **not** called for `O`
   - verifies `ValidateOrderWithoutStock` **is** called for `O`
   - verifies preserved non-stock failure data reaches store (`Validate2Success=false`, `Validate2value=1000`)
   - verifies overall summary is not falsely forced to success when non-stock validation fails

2. `TestOrderControllerCreate_SX2184NilOrderTypeStillValidatesStock`
   - verifies nil `order_type` stays on original validation path

3. `TestOrderControllerCreate_SX2184EmptyOrderTypeStillValidatesStock`
   - verifies empty string behaves like legacy/no-type input and still validates stock

4. `TestOrderControllerCreate_SX2184OrderTypeCStillValidatesStock`
   - verifies `C` remains on the original stock validation path

5. `TestOrderControllerCreate_SX2184SalesOrderStillValidatesStock`
   - verifies `SO` still uses stock validation path and still returns bad request on insufficient stock

#### Service coverage

File: `sales/service/order_type_helper_test.go`

Added/updated tests:

1. `TestStore_SX2184TakingOrderSkipsStockMutationAndPersistsOriginalQty`
   - verifies no `SalesStockUpdates` for `O`
   - verifies `order_type` persisted
   - verifies `opr_type = O`
   - verifies `validate_stok = false`
   - verifies `validate_stok_message == nil`
   - verifies PO/original qty mapping and qty conversion remain correct

2. `TestStore_SX2184NilOrderTypeStillMutatesInventory`
   - verifies nil `order_type` remains on old stock mutation path

3. `TestStore_SX2184EmptyOrderTypeStillMutatesInventory`
   - verifies empty-string `order_type` remains on old stock mutation path

4. `TestStore_SX2184OrderTypeCStillMutatesInventory`
   - verifies `C` remains on the original stock mutation path

5. `TestStore_SX2184SalesOrderStillMutatesInventory`
   - verifies `SO` still mutates inventory when processed

#### Read/status nil handling coverage

File: `sales/service/order_status_helper_test.go`

Added:

- `TestValidationResultFromOrderList_NilValidateStokMessage`
  - verifies nil `ValidateStokMessage` round-trips safely into validation response as empty string without panics
  - verifies boolean stock result is preserved separately from nullable message

## Validation run results

### Runtime preflight

Command:

```bash
rtk docker compose -f docker-compose.yml ps
```

Result:

- Passed as environment check
- No running compose services were present

### Targeted controller tests

Command:

```bash
rtk go test ./controller -run 'TestOrderControllerCreate_SX2184'
```

Result:

- Passed after one red/green fix for empty-string normalization
- `Go test: 5 passed in 1 packages`

### Targeted service tests

Command:

```bash
rtk go test ./service -run 'TestStore_SX2184|TestValidationResultFromOrderList_NilValidateStokMessage'
```

Result:

- Passed
- `Go test: 6 passed in 1 packages`

### Full sales module test suite

Command:

```bash
rtk go test ./...
```

Result:

- Passed
- `Go test: 219 passed in 22 packages`

## API / DB smoke evidence

Completed local Docker/API/DB smoke for `order_type = "O"` after remediation.

Commands/results:

```bash
rtk docker compose -f docker-compose.yml up -d rabbitmq sales
```

- Result: `rabbitmq` healthy and `sales` started.

```bash
curl -sS -i http://localhost:9004/ping
```

- Result: HTTP 200 `It works`.

```bash
PGPASSWORD=postgres /opt/homebrew/opt/postgresql@18/bin/psql -h localhost -p 5432 -U postgres -d ggn_scyllax -v ON_ERROR_STOP=1 -f sales/migration/sls.order/add_order_type_and_original_qty_po_fields.sql
```

- Result: migration completed; columns already existed and were skipped via `IF NOT EXISTS` notices.

Fixture used from local DB:

- `cust_id = C220010001`
- `parent_cust_id = C22001`
- `salesman_id = 210`
- `outlet_id = 1390`
- `wh_id = 241`
- `pro_id = 472`
- pre-smoke `inv.warehouse_stock.qty = 0`, `qty_on_order = 108`

Local JWT:

- Generated from local `sales/.env` secret into temp file only (`/var/folders/.../T/opencode/sx2184_token`).
- Token value was not copied into code, evidence, or final notes.

API smoke:

```bash
POST http://localhost:9004/v1/orders
```

Payload used `order_type = "O"`, qty1 `10`, warehouse stock `0`.

Result:

- HTTP 201
- response `ro_no = SO2606080001`

DB verification after create:

```sql
SELECT ro_no, order_type, validate_stok, validate_stok_message, opr_type, data_status, data_source, notes
FROM sls."order"
WHERE cust_id='C220010001' AND ro_no='SO2606080001';
```

Result summary:

- `order_type = O`
- `validate_stok = false`
- `validate_stok_message = NULL`
- `opr_type = O`
- `data_status = 2`
- `data_source = 1`

```sql
SELECT original_qty_po1, original_qty_po2, original_qty_po3,
       qty_po1, qty_po2, qty_po3, qty_po,
       qty1, qty2, qty3, qty, qty_final
FROM sls.order_detail
WHERE cust_id='C220010001' AND ro_no='SO2606080001';
```

Result summary:

- `original_qty_po1/2/3 = 10/0/0`
- `qty_po1/2/3 = 10/0/0`
- `qty_po = 10`
- `qty1/2/3 = NULL/NULL/NULL`
- `qty = 0`, `qty_final = 0`

```sql
SELECT count(*) FROM inv.stock WHERE cust_id='C220010001' AND tr_no='SO2606080001';
```

Result: `0` rows.

```sql
SELECT qty, qty_on_order FROM inv.warehouse_stock
WHERE cust_id='C220010001' AND wh_id=241 AND pro_id=472;
```

Result after create: `qty = 0`, `qty_on_order = 108`; unchanged from pre-smoke.

SO API regression smoke was not executed because the local existing stock-validation formula may allow the selected zero-stock fixture through and would create an extra inventory mutation. Regression remains covered by targeted controller/service tests proving non-`O` paths still call stock validation/mutation.

After smoke, local `sales` and `rabbitmq` compose services were stopped with:

```bash
rtk docker compose -f docker-compose.yml stop sales rabbitmq
```

## Migration apply evidence

Applied locally to `ggn_scyllax` via Homebrew `psql`.

Result:

- `sls.order.order_type` already existed and was skipped by `ADD COLUMN IF NOT EXISTS`.
- `sls.order_detail.original_qty_po1/2/3` already existed and were skipped by `ADD COLUMN IF NOT EXISTS`.
- Migration completed successfully.

## Security / secrets check

- No `.env` files changed.
- No postman credentials changed.
- No tokens or auth headers copied into code, tests, or evidence.
- Scope remained inside `sales` plus the requested evidence artifact file.

## Residual risks / notes

1. Live DB smoke still pending
   - SQL NULL persistence for `validate_stok_message` is now compile-safe in code/model shape, but it was not verified against a live local DB in this pass.
   - Service-level tests cover the model state sent to persistence, not the final stored row.

2. Read API contract changed slightly
   - `OrderResponse.ValidateStokMessage` is now nullable (`*string`) instead of always string in Go shape.
   - This matches the requested SQL NULL behavior, but downstream consumers expecting always-present JSON string may need awareness if they deserialize strictly.

3. Taking-order qty compatibility behavior remains as before
   - Create `O` still zeroes sales-order qty totals and keeps PO/original qty fields populated.
   - Full runtime schema constraints for qty-null semantics remain un-smoke-tested locally.

## Changed files

- `sales/controller/order_controller.go`
- `sales/controller/order_controller_test.go`
- `sales/entity/order.go`
- `sales/model/order.go`
- `sales/service/order_service.go`
- `sales/service/order_status_helper.go`
- `sales/service/order_status_helper_test.go`
- `sales/service/order_type_helper.go`
- `sales/service/order_type_helper_test.go`
- `sales/service/validate_order_service.go`
