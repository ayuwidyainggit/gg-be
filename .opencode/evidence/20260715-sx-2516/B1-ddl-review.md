# B1 — DDL Review: `20260715_add_sales_update_temp.sql`

## Scope
Create `import.sales_update_temp` staging table for secondary sales replace import.

## Confirmed source facts (A3)
- `import.sales_update_temp` does NOT exist (`to_regclass` returned NULL).
- `sls.order_detail` has no FK constraint — no CASCADE risk.
- `sls.order` columns: `cust_id VARCHAR(10)`, `ro_no VARCHAR(30)`, `ro_date DATE`, `is_sales_mapping BOOLEAN`.

## File reviewed
`sales/migration/20260715_add_sales_update_temp.sql`

## Column mapping (22 columns)

| Column | Type | Constraints | Notes |
|---|---|---|---|
| id | BIGINT | GENERATED ALWAYS AS IDENTITY PK | bigserial equivalent, identity standard |
| history_id | BIGINT | NOT NULL | FK to import.import_history |
| cust_id | VARCHAR(10) | NOT NULL | Matches sls.order.cust_id width |
| document_no | VARCHAR(30) | NOT NULL | Matches sls.order.ro_no width |
| document_date | DATE | NOT NULL | Matches sls.order.ro_date type |
| outlet_code | VARCHAR(50) | NOT NULL | |
| outlet_name | VARCHAR(255) | NOT NULL | |
| salesman_code | VARCHAR(50) | NOT NULL | |
| salesman_name | VARCHAR(255) | NOT NULL | |
| pro_code | VARCHAR(50) | NOT NULL | |
| pro_name | VARCHAR(255) | NOT NULL | |
| price | NUMERIC(20,4) | NOT NULL DEFAULT 0 | Matches sls.order_detail sell_price precision |
| unit | VARCHAR(5) | NOT NULL | Matches sls.order_detail unit_id width |
| qty | REAL | NOT NULL DEFAULT 0 | Matches sls.order_detail qty type |
| gross_sales | NUMERIC(20,4) | NOT NULL DEFAULT 0 | |
| promo | NUMERIC(20,4) | NOT NULL DEFAULT 0 | |
| discount | NUMERIC(20,4) | NOT NULL DEFAULT 0 | |
| ppn | NUMERIC(20,4) | NOT NULL DEFAULT 0 | |
| net_sales_inc_ppn | NUMERIC(20,4) | NOT NULL DEFAULT 0 | |
| status_insert | VARCHAR(20) | NOT NULL DEFAULT 'SUCCESS' | SUCCESS/FAILED |
| error_message | TEXT | nullable | NULL when status_insert=SUCCESS |
| created_at | TIMESTAMP WITHOUT TIME ZONE | NOT NULL DEFAULT now() | |

## Indexes
- `idx_sales_update_temp_history_id` on `(history_id)` — worker lookup by history.
- `idx_sales_update_temp_cust_id` on `(cust_id)` — cleanup per customer.

## Static validation
- `CREATE TABLE IF NOT EXISTS` — idempotent.
- `GENERATED ALWAYS AS IDENTITY` — PostgreSQL 10+ standard, no sequence name collision.
- All `NOT NULL` columns have explicit defaults where appropriate.
- `COMMENT ON` statements for table and key columns.
- Transaction wrapped in `BEGIN`/`COMMIT` matching local migration convention.
- No secrets, tokens, or real data.

## No runtime SQL execution
This file was created but NOT applied against any database. B5 owns apply.

## Risks
- None identified. Table is new, no existing data to preserve.
