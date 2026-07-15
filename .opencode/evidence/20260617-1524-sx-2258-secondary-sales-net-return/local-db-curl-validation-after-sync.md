# SX-2258 local validation after staging data sync

Tanggal: 2026-06-17 Asia/Jakarta

## Scope

- Local DB: `ggn_scyllax` via `localhost:5432`
- Local Docker services:
  - `system` on `localhost:9001`
  - `sales` on `localhost:9004`
- Credentials were used only transiently for local login; no tokens/passwords stored.

## Service readiness

```text
scylla-system: up, port 9001
scylla-sales: up, port 9004
scylla-redis: healthy
scylla-rabbitmq: healthy
system /ping: 200
sales /ping: 200
```

## Dataset found after sync

Broad local DB daily search found expected SX-2258 values at:

```text
cust_id=C260020001
invoice_date=2026-06-03
qty=134
total_discount_promo=1238740
qty_order=146
qty_return=12
```

## Direct DB validation

Query range:

```text
cust_id=C260020001
invoice_date >= 2026-06-03
invoice_date < 2026-06-04
```

Direct SQL result:

```text
total_gross_sale=5201300000
total_discount_promo=1238740
total_ppn=520006126
net_sales_exc_ppn=5200061260
net_sales=5720067386
total_salesman=2
total_outlet=3
total_product=6
qty=134
qty_return=12
return_rate=0.01
net_sales_return=693693
```

Formula proof:

```text
qty = qty_order - qty_return = 146 - 12 = 134
total_discount_promo = discount_order - discount_return = 1238740
```

## cURL/API validation

Request params:

```text
GET /v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001&from=1780444800&to=1780531200
```

Principal token result:

```text
status=200
qty=134
total_discount_promo=1238740
qty_return=12
net_sales=5720067386
```

Distributor token result:

```text
status=200
qty=134
total_discount_promo=1238740
qty_return=12
net_sales=5720067386
```

Distributor parent-scope guard:

```text
GET /v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C26002&from=1780444800&to=1780531200
status=403
message="cust_id is outside authorized scope"
```

## Result

SX-2258 expected values are now reproduced on local synced data:

```text
Number of Product Sold / qty = 134
Discount and Promo / total_discount_promo = 1238740
```

DB direct SQL and API/cURL results match.
