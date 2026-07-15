# B7 — Detail Staging Table

## Source basis
- `sales/service/order_service.go:6783-6830` — initial `CreateOrderDetBody` snapshot
- `sales/service/order_service.go:7013-7022` — row-level financial fields
- `sales/service/order_service.go:7113-7127` — grouped detail fields consumed by `Store`
- `sales/entity/order_detail.go:3-94` — complete `CreateOrderDetBody` shape
- Plan §13 (lines 108-114) — mapping and aggregation rules

## Migration
`20260715_create_sales_update_temp_detail.sql` creates `import.sales_update_temp_detail` inside `BEGIN`/`COMMIT` with `CREATE TABLE IF NOT EXISTS`. It links each detail to its staging header through:

- `history_id` — tenant/import processing correlation
- `staging_key` — FK to `import.sales_update_temp.id`, `ON DELETE CASCADE`
- `document_no` — stable order reconstruction key
- `seq_no` — grouped detail order

## Persisted field groups
- Identity: `document_no`, `seq_no`, `pro_id`, `item_type`
- Quantities: `qty`, `qty_po`, `qty_po1..3`, `qty_final`, `qty1..5`, `qty1_final..5_final`
- Product snapshots: `purch_price1..5`, `sell_price1..5`, `sell_price_system1..5`
- Amounts and discounts: `amount`, `amount_final`, `promo_value`, `promo_value_final`, `disc_value`, `disc_value_final`
- VAT: `vat`, `vat_value`, `vat_value_final`, `vat_bg`, `vat_lg_sell`, `vat_bg_value`, `vat_lg_value`, `vat_lg_sell_value`, `vat_value_po`
- Unit snapshots: `unit_id1..5`, `conv_unit2..5`
- Promotion slots: `promo_so1..5`, `promo_final1..5`, `promo_po1..5`, promotion booleans
- Optional parsed fields: `batch_no`, `exp_date`, `notes`, `discount_id`, `promo_id`, `disc_po`

No master IDs/names are re-looked up by worker. Values are parse-time snapshots.

## Rollback
`20260715_rollback_sales_update_temp_detail.sql` drops only `import.sales_update_temp_detail` with `DROP TABLE IF EXISTS` inside `BEGIN`/`COMMIT`. Not run locally; forward migration remains applied for schema inspection.

## Validation
- Forward SQL applied on local `ggn_scyllax`
- `\d import.sales_update_temp_detail` confirms table, FK, 84 detail columns, and indexes
- Re-running is idempotent by `CREATE TABLE IF NOT EXISTS` and `CREATE INDEX IF NOT EXISTS`
