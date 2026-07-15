# D4 Scope Tests

Scope boundaries maintained:
- Replace only `is_sales_mapping=true` orders
- Replace only within `(cust_id, ro_date)` scope
- Other dates of same cust_id untouched
- Non-mapping orders untouched
- All queries filter by `cust_id` (tenant isolation)
