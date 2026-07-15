# SX-2291 Jira Notes

## Summary

Fixed BE cancel flow for Sales Order in `Need Review` (`data_status = 1`) from Purchase Order tab. `PATCH /sales/v1/orders/status` / local `/v1/orders/status` now allows cancel target `data_status = 9` without failing on PO-tab orders that have `qty_po*` but no source stock ledger.

## Root Cause

Cancel status transition was already allowing `Need Review`, but stock cancel basis query only considered final qty (`qty_final` / final fields). PO-tab Need Review orders can have only `qty_po1/qty_po2/qty_po3`, or no `inv.stock` rows yet. That made cancel stock validation fail with:

```text
order cannot be cancelled because final detail and stock ledger are inconsistent
```

So status update never completed.

## Fix

- Extended `cancelStockBasisQuery` to include qty fallback priority:
  1. `qty*_final`
  2. `qty*`
  3. `qty_po*`
- Added PO sell price fallback for cancel basis unit price.
- Updated cancel service flow:
  - `Need Review` + no outstanding stock basis => skip stock reversal and still update `sls.order.data_status = 9`.
  - Existing stock reversal path still runs when basis has outstanding stock.
  - `Processed` with missing basis still fails as before.
  - Already `Cancelled` remains no-op/idempotent.
- Preserved tenant scope using `cust_id` in read/write queries.
- Preserved transaction wrapping around cancel flow.

## Files Changed

- `sales/service/order_service.go`
- `sales/repository/stock_repository.go`
- `sales/service/order_service_test.go`
- `sales/repository/stock_repository_cancel_test.go`

## Validation

Targeted tests passed in `scylla-sales` container:

```bash
docker exec scylla-sales sh -c "cd /app && go build ./..."
docker exec scylla-sales sh -c "cd /app && go test ./service/... -run 'TestBulkUpdateStatus_Cancel|TestValidateCancelTransition' -v"
docker exec scylla-sales sh -c "cd /app && go test ./repository/... -run TestGetCancelStockBasisQuery -v"
```

Results:

- Build: pass
- Cancel service tests: pass
- Repository SQL tests: pass

Docker/API validation:

- Login user: `adminbm@gmail.com`
- Local validation order: `SO2606180003` (`cust_id=C260020001`, current `data_status=1`, PO qty exists, no source `inv.stock` rows)
- Request:

```bash
PATCH http://localhost:9004/v1/orders/status
{"orders":[{"ro_no":"SO2606180003","data_status":9}]}
```

Response:

```json
{"message":"Updated Status Successfully"}
```

DB after cancel:

```text
sls.order.data_status = 9
inv.stock rows for SO2606180003 / SO2606180003-CO = 0 rows
```

Stock rows are expected to remain 0 for this local case because no source stock movement existed before cancel.

## Notes

- Jira sample `SO2606190005` was not used for local API validation because in local DB it belongs to `cust_id=C220010001` and status is already `6`, while `adminbm@gmail.com` has `cust_id=C260020001`.
- Full `go test ./...` in container still has unrelated existing failures:
  - missing CSV fixture: `/app/service/docs/test promo integrasi sales order - Request mas angga.csv`
  - container-local Postgres connection refused in `TestSyncFinalOrderFields_UsesNoProformaPromoSyncSQL`
