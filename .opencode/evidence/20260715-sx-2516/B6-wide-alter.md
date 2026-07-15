# B6 — Wide ALTER Snapshot: import.sales_update_temp

## Source basis
- `sales/service/order_service.go:6783-6830` — `mapQtyAndPriceToSlot` detail builder
- `sales/service/order_service.go:7008-7178` — `CreateOrderBody` assembly from parsed lines
- `sales/entity/order.go:57+` — `CreateOrderBody` struct
- `sales/entity/order_detail.go:3+` — `CreateOrderDetBody` struct
- Plan §"Evidence note — staging shape boundary" (line 96-98)

## Design rationale
Existing `import.sales_update_temp` has only 22 raw import columns. The current `Store` path (`order_service.go:6633-6640`) calls `service.Store(req, importValidation)` which expects a fully-materialized `CreateOrderBody` with header aggregates and detail snapshot fields. The worker cannot relookup mutable master data (salesman, outlet, product). Therefore all snapshot fields the Store path consumes from the parse result must be persisted.

## Columns added (idempotent ALTER TABLE ADD COLUMN IF NOT EXISTS)

| Column | Type | Default | Source in Store path |
|--------|------|---------|---------------------|
| `parent_cust_id` | VARCHAR(10) | NULL | `CreateOrderBody.ParentCustId` (line 7136) |
| `user_id` | BIGINT | NULL | `CreateOrderBody.CreatedBy` (line 7169) |
| `ro_no` | VARCHAR(30) | NULL | `CreateOrderBody.RoNo` (line 7134) |
| `salesman_id` | BIGINT | NULL | `CreateOrderBody.SalesmanId` (line 7140) |
| `wh_id` | BIGINT | NULL | `CreateOrderBody.WhId` (line 7141) |
| `outlet_id` | BIGINT | NULL | `CreateOrderBody.OutletID` (line 7142) |
| `parent_product_id` | INTEGER | NULL | `detail.ProId` = `parentProId` (line 6786) |
| `grouping_product_key` | VARCHAR(100) | NULL | Deduplication key (line 7024-7030) |
| `subtotal` | NUMERIC(20,4) | NULL | `CreateOrderBody.SubTotal` (line 7151) |
| `subtotal_final` | NUMERIC(20,4) | NULL | `CreateOrderBody.SubTotalFinal` (line 7152) |
| `disc` | NUMERIC(20,4) | 0 | `CreateOrderBody.Disc` (line 7153) |
| `disc_value` | NUMERIC(20,4) | 0 | `CreateOrderBody.DiscValue` (line 7154) |
| `disc_value_final` | NUMERIC(20,4) | 0 | `CreateOrderBody.DiscValueFinal` (line 7155) |
| `promo_value` | NUMERIC(20,4) | 0 | `CreateOrderBody.PromoValue` (line 7156) |
| `promo_value_final` | NUMERIC(20,4) | 0 | `CreateOrderBody.PromoValueFinal` (line 7157) |
| `cash_disc_value` | NUMERIC(20,4) | 0 | `CreateOrderBody.CashDiscValue` (line 7160) |
| `tot_disc1` | NUMERIC(20,4) | 0 | `CreateOrderBody.TotDisc1` (line 7161) |
| `tot_disc2` | NUMERIC(20,4) | 0 | `CreateOrderBody.TotDisc2` (line 7162) |
| `vat` | NUMERIC(20,4) | 0 | `CreateOrderBody.Vat` (line 7163) |
| `vat_value` | NUMERIC(20,4) | 0 | `CreateOrderBody.VatValue` (line 7164) |
| `vat_value_final` | NUMERIC(20,4) | 0 | `CreateOrderBody.VatValueFinal` (line 7165) |
| `total` | NUMERIC(20,4) | NULL | `CreateOrderBody.Total` (line 7166) |
| `total_final` | NUMERIC(20,4) | NULL | `CreateOrderBody.TotalFinal` (line 7167) |
| `data_status` | INTEGER | 6 | `CreateOrderBody.DataStatus` (line 7168) |
| `data_source` | BIGINT | 3 | `CreateOrderBody.DataSource` (line 7170) |
| `pay_type` | BIGINT | 1 | `CreateOrderBody.PayType` (line 7148) |
| `is_closed` | BOOLEAN | FALSE | `CreateOrderBody.IsClosed` (line 7173) |
| `validate_stok` | BOOLEAN | FALSE | Hardcoded in plan §12 |
| `is_sales_mapping` | BOOLEAN | TRUE | `CreateOrderBody.IsSalesMapping` (line 7177) |
| `invoice_no` | VARCHAR(30) | NULL | `CreateOrderBody.InvoiceNo` (line 7175) |
| `invoice_date` | DATE | NULL | `CreateOrderBody.InvoiceDate` (line 7176) |
| `due_date` | DATE | NULL | `CreateOrderBody.DueDate` (line 7139) |
| `delivery_date` | DATE | NULL | `CreateOrderBody.DeliveryDate` (line 7144) |
| `created_by` | BIGINT | NULL | `CreateOrderBody.CreatedBy` (line 7169) |
| `staging_key` | BIGINT | NULL | Link key for detail staging table |

## Stable constants (hardcoded values NOT persisted — recorded here)
These are set in the Store path and will be reconstructed from code constants, not from staging:

| Constant | Value | Source |
|----------|-------|--------|
| `importDataStatus` | 6 (Invoicing) | `order_service.go:7220` |
| `importDataSource` | 3 (Import) | `order_service.go:7221` |
| `defaultImportVat` | 11 | `order_service.go:7222` |
| `defaultImportVatRate` | 11 | `order_service.go:7223` |
| `payType` | 1 (Cash On Delivery) | `order_service.go:7129` |
| `is_del` | false | Plan §12 |
| `is_paid_off` | false | Plan §12 |
| `is_proforma_inv` | false | Plan §12 |
| `disc` | 0 | `order_service.go:7153` |
| `vat` (percentage) | 0 | `order_service.go:7163` |
| `PromoBgValue` | nil | `order_service.go:7158` |
| `PromoBgValueFinal` | nil | `order_service.go:7159` |
| `CashDiscValue` | 0 | `order_service.go:7160` |
| `TotDisc1` | 0 | `order_service.go:7161` |
| `TotDisc2` | 0 | `order_service.go:7162` |
| `ValDate` | nil | `order_service.go:7138` |
| `OutletAddress1` | nil | `order_service.go:7143` |
| `OrderNo` | nil | `order_service.go:7145` |
| `PoNo` | nil | `order_service.go:7146` |
| `VehicleNo` | nil | `order_service.go:7147` |
| `ReffNo` | nil | `order_service.go:7149` |
| `MobileID` | nil | `order_service.go:7150` |
| `OrderType` | nil | `order_service.go:7172` |
| `Notes` | nil | `order_service.go:7174` |

## Indexes
- `idx_sales_update_temp_staging_key` on `(staging_key)` — FK target for detail table

## Rollback
`20260715_rollback_alter_sales_update_temp_wide.sql` — DROP COLUMN IF EXISTS for all added columns + DROP INDEX.

## Validation
- SQL parsed without error
- `\d import.sales_update_temp` shows all 36 new columns present
- Existing 22 raw columns preserved
- Idempotent: re-running ALTER produces no error
