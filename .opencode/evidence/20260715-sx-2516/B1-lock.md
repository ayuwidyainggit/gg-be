# B1 LockOrderByScope

Method added to `OrderRepository` interface and `RepositoryOrderImpl`:

- `LockOrderByScope(ctx, custId, roDates)` — loops dates, calls `pg_advisory_xact_lock(hashtextextended($1, 0))` per scope row, then `SELECT 1 FROM sls.order WHERE cust_id=$1 AND is_sales_mapping=true AND ro_date = ANY($2::date[]) FOR UPDATE`
- Uses `model(ctx)` for tx-awareness
