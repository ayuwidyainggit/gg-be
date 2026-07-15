# SX-2258 local DB and curl validation

Tanggal: 2026-06-17 Asia/Jakarta

## Scope

- Local DB: `ggn_scyllax` via `localhost:5432`
- Docker local services:
  - `system` on `localhost:9001`
  - `sales` on `localhost:9004`
- Login users validated transiently via `POST /v1/users/login`.
- Tokens and passwords were not saved in evidence.

## Docker readiness

- `docker ps` showed:
  - `scylla-system` up on `9001`
  - `scylla-sales` up on `9004`
  - `scylla-redis` healthy
  - `scylla-rabbitmq` healthy
- `GET /ping`:
  - `system`: `200`
  - `sales`: `200`

## Login validation

Sanitized results:

```text
principal login status=200 token_present=true cust_id=C26002 parent_cust_id=C26002
distributor login status=200 token_present=true cust_id=C260020001 parent_cust_id=C26002
```

## DB validation: monthly summary, `C260020001`, June 2026

Direct SQL against `ggn_scyllax`, formula matching implemented query:

```text
total_gross_sale=5460394000
total_discount_promo=1403740
total_ppn=545899026
net_sales_exc_ppn=5458990260
net_sales=6004889286
total_salesman=7
total_outlet=9
total_product=14
qty=265
qty_return=12
return_rate=0.01
net_sales_return=693693
```

Formula proof:

```text
qty = qty_order - qty_return = 277 - 12 = 265
total_discount_promo = discount_order - discount_return = 1423110 - 19370 = 1403740
```

## cURL/API validation: monthly summary

Principal token, child distributor scope:

```text
GET /v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001
status=200
qty=265
qty_return=12
total_discount_promo=1403740
net_sales=6004889286
```

Distributor token, own scope:

```text
GET /v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001
status=200
qty=265
qty_return=12
total_discount_promo=1403740
net_sales=6004889286
```

Distributor token, parent scope denied:

```text
GET /v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C26002
status=403
message="cust_id is outside authorized scope"
```

## cURL/API validation: optional filters

Filter sample from local DB:

```text
cust_id=C260020001
month=6
year=2026
outlet_ids=1840
salesman_ids=479
pro_ids=10751
```

Direct DB expected:

```text
qty=60
total_discount_promo=0
qty_return=0
```

API result:

```text
GET /v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001&outlet_ids=1840&salesman_ids=479&pro_ids=10751
status=200
qty=60
total_discount_promo=0
qty_return=0
```

## cURL/API validation: trend

Principal token, child distributor scope:

```text
GET /v1/reports/secondary-sales/trend-sales?year=2026&cust_id=C260020001
status=200
rows=12
june_total_discount_promo=1403740
june_total_gross_sale=5460394000
june_net_sales=6004889286
```

## QA expected-number search in local DB

A broad local DB search for monthly `qty=134` or `total_discount_promo=1238740` returned no rows.

The original Jira expected numbers are therefore not present in this local `ggn_scyllax` dataset with available scope. Local validation proves the fixed formula and API response match local DB results, but cannot reproduce staging QA exact numbers on this DB snapshot.

## Result

Local DB/API risk remediated for available dataset:

- DB formula and API output match for monthly summary.
- DB formula and API output match for optional filters.
- Trend endpoint returns net `total_discount_promo` matching summary month.
- Auth scope works for principal and distributor users.
- No credentials/tokens written to source or evidence.
