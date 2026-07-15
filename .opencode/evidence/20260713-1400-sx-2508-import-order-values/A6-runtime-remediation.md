# A6 Runtime remediation

## Scope
Patched imported sales-order Store path only. No endpoint, migration, converter, module, secret, or transaction-boundary changes.

## Root causes
- `sales/service/order_service.go:373-376`: Store unconditionally nulled `orderModel.InvoiceDate` and `orderModel.InvoiceNo` after Automapper. Imported parser values existed at `:7146-7147`, then were erased.
- `model.OrderDetail.QtyPo` is `float64`, while the entity/parser source is `*float64`. GORM `Create` therefore persists its zero value instead of SQL `NULL`, even when imported parser input is nil.

## Patch
- Preserve parsed invoice fields when `isImportedOrder(&request)` is true.
- After `StoreDetail` creates imported Normal or Promo detail, call `UpdateDetailPartial` in the same transaction with `map[string]interface{}{"qty_po": nil}`. GORM map updates emit `qty_po = NULL` without broad model type changes.
- Non-import order behavior remains unchanged.

## Validation
| Command | Result |
|---|---|
| `cd sales && rtk gofmt -w service/order_service.go` | Pass |
| `cd sales && rtk go test ./service -run 'Test.*Import|Test.*Order'` | Pass, 107 tests |
| `cd sales && rtk go build .` | Pass |
| `cd sales && rtk go test ./...` | Existing unrelated failures: `TestDetailV2_PostRolloutWithoutSnapshot_UsesConsultV2ByTab`, `TestDetailV2_PreRolloutWithoutSnapshot_MustConsultV2ByTab`; import tests pass. See `A4-validation.md` and `A3-detail-stock.log`. |

## Runtime / SQL
- Fresh local Docker import used a copied template with unique document number `INVJUL1307-101`; source template remained unchanged.
- `POST /v1/orders/import` returned HTTP 200: `Secondary sales file imported successfully.`
- Read-only local SQL returned: `INVJUL1307-101|INVJUL1307-101|2026-07-11|1|1|0|`.
- SQL fields, in order: `ro_no|invoice_no|invoice_date|detail_count|qty_po_is_null_count|qty_po_is_not_null_count|max_qty_po`.
- This confirms `invoice_no=ro_no`, `invoice_date=ro_date`, one persisted detail, and `qty_po IS NULL` for all persisted details. No token, password, Authorization header, or full payload recorded.

## Contract impact
- Imported header: `invoice_no=ro_no`, `invoice_date=ro_date` confirmed at runtime.
- Imported detail: `qty_po=NULL` confirmed at runtime.
- Normal order invoice clearing and quantity persistence preserved.
- Mapped DetailV2, Large/Middle/Small orientation, converter, tenant and transaction boundaries unchanged.
- Rollback: revert `sales/service/order_service.go` changes; no migration/backfill required.

## Status
Runtime acceptance confirmed locally. Remaining risk: two pre-existing unrelated full-suite promotion tests remain red.
