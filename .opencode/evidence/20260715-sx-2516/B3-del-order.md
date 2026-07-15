# B3 DeleteOrderByScope

Method added:

- `DeleteOrderByScope(ctx, custId, roDates)` — `DELETE FROM sls.order WHERE cust_id=$1 AND is_sales_mapping=true AND ro_date = ANY($2::date[])`
- Returns RowsAffected
- Uses `model(ctx)` for tx-awareness
