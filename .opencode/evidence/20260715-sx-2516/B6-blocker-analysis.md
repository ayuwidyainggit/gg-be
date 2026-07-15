# B6-F6 — Reschema Staging Blocker Analysis

## Trigger
C1-F6 fail-closed (`.opencode/evidence/20260715-sx-2516/C1-F6-remediation.md`): 22-column staging `import.sales_update_temp` cannot rebuild `model.Order` + `model.OrderDetail` mapped payload. Worker must reject and mark history FAILED until staging persists enough mapping fields to reconstruct orders.

## Schema gap (verified from A3 + model files)
`import.sales_update_temp` currently has only: `id`, `history_id`, `cust_id`, `document_no`, `document_date`, `outlet_code`, `outlet_name`, `salesman_code`, `salesman_name`, `pro_code`, `pro_name`, `price`, `unit`, `qty`, `gross_sales`, `promo`, `discount`, `ppn`, `net_sales_inc_ppn`, `status_insert`, `error_message`, `created_at`.

Missing required fields to rebuild a real `model.Order`:
- header IDs: `salesman_id`, `wh_id`, `outlet_id`, `delivery_date`, `val_date`, `due_date`, `order_type`, `opr_type`
- header money: `sub_total`, `sub_total_final`, `disc`, `disc_value`, `disc_value_final`, `promo_value`, `promo_value_final`, `promo_bg_value`, `promo_bg_value_final`, `cash_disc_value`, `tot_disc1`, `tot_disc2`, `vat`, `vat_value`, `vat_value_final`, `total`, `total_final`, `pay_type`, `invoice_no`, `invoice_date`, `notes`
- flags: `is_printed`, `is_paid_off`, `is_proforma_inv`, `is_del`, `is_closed`, `validate_stok`, `data_status`, `data_source`, `created_by`, `created_at`

Missing required fields to rebuild `model.OrderDetail`:
- detail IDs and type: `seq_no`, `pro_id`, `item_type`
- quantities: `qty_final`, `qty_po`, `qty1..5`, `qty1_final..qty5_final`
- prices: `purch_price1..5`, `sell_price1..5`, `sell_price_po1..3`, `sell_price_system1..5`, `sell_price_final1..3`
- per-line money: `amount`, `amount_final`, `disc_value`, `disc_value_final`, `promo_value`, `promo_value_final`, `vat`, `vat_value`, `vat_value_final`, `vat_bg_value`, `vat_lg_sell_value`
- units: `unit_id1..5`, `conv_unit2..5`
- links: `discount_id`, `promo_id`, `batch_no`, `exp_date`, `notes`

That is ~80 columns beyond current 22. Net new staging column count is on the order of 80, not 22.

## Why a different design is recommended
The plan's 22-column staging was based on the export template header only. The implementation surfaces the gap because parsed rows in `parseImportOrders` are mapped through outlet/salesman/product/unit resolution into `CreateOrderBody`; that body has the full 80-column shape. A flat row-per-product staging is the wrong shape for the worker.

Two viable designs:

### Design A — JSONB payload on history
- Add one column to `import.import_history`: `payload JSONB`.
- Enqueue serializes the full parsed `[]CreateOrderBody` to JSONB.
- Worker deserializes JSONB and uses it as the mapped input.
- Staging table `import.sales_update_temp` can stay as the audit-trail row store for distinct (cust_id, document_date) scope and validation results, but the mapped payload lives in JSONB.
- Replace SQL: lock + delete + insert new `sls.order`/`sls.order_detail` from JSONB.
- Pros: minimal schema growth, full fidelity, easy to extend.
- Cons: large text column per history row, 4,200 rows × few hundred bytes each = 2-5 MB history rows; but history rows are infrequent (one per import request), so size is OK.

### Design B — Wide flat staging (this PR's intent)
- ALTER `import.sales_update_temp` with ~80 extra columns.
- Each row is one detail line.
- Staging grows 1 row per (document_no, pro_code) line.
- Worker reads staging and rebuilds `model.Order` + `model.OrderDetail` rows.
- Pros: SQL-driven, easy to query.
- Cons: large DDL change, fragile to parser evolution, ALTER cost on existing data, hard to maintain.

## Recommendation
Switch to Design A. Add one `payload JSONB` column on `import.import_history`. Staging can shrink to only validation-status rows (or be removed entirely if not needed for audit). This keeps the existing B1 staging table and B3 history ALTER in place, plus adds B6: `ALTER TABLE import.import_history ADD COLUMN IF NOT EXISTS payload JSONB;` plus a rollback B7.

## Decision needed
- Confirm Design A as the path.
- Confirm staging table `import.sales_update_temp` is still useful for validation rows (yes) or can be dropped (no — keep for now to limit diff).
- After sign-off, redo C1-F6 implementation: parser maps to CreateOrderBody as before, enqueue serializes body to JSONB and inserts into history, worker reads JSONB, runs WithinTransaction with lock+delete+insert from JSONB.

## Scope impact if confirmed
- B6: forward migration `20260715_add_history_payload.sql`.
- B7: rollback for B6.
- B8: apply both to local DB.
- F2-F6: rewrite worker to deserialize JSONB and run the replace core; add focused tests.
- E1/H1/I1: no change to 7-day rule, controller branch, or smoke.

## Status
- 9 of 44 tasks complete.
- 1 blocked: C1 (and its dependents C2-C5, D1-D8, E1-E2, F1-F6).
- 34 pending.

Hard stop. Awaiting user sign-off before B6 onward.
