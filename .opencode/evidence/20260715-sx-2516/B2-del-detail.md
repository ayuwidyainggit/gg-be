# B2 DeleteOrderDetailByScope

Method added:

- `DeleteOrderDetailByScope(ctx, custId, roDates)` — `DELETE FROM sls.order_detail d USING sls.order o WHERE o.ro_no=d.ro_no AND o.cust_id=$1 AND o.is_sales_mapping=true AND o.ro_date = ANY($2::date[])`
- Returns RowsAffected
- Uses `model(ctx)` for tx-awareness
