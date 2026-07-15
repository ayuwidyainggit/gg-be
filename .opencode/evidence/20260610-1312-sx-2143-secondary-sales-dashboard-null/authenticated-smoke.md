# Authenticated HTTP Smoke Evidence — SX-2143

Task ID: `20260610-1312-sx-2143-secondary-sales-dashboard-null`
Tanggal: 2026-06-10 Asia/Jakarta

## Scope

Authenticated local HTTP smoke against Docker services after source-table remediation for Secondary Sales dashboard.

Tokens were acquired transiently from local `system` service and were not printed or persisted.

## Local services

```text
rtk docker compose -f docker-compose.yml ps
system: 127.0.0.1:9001, status Up
sales: 127.0.0.1:9004, status Up
```

## Login endpoint

Repo-local route/payload evidence:

- `system/controller/user_controller.go`: `POST v1/users/login`
- `system/entity/user.go`: payload fields `email`, `password`; response token under `data.token.access`

Smoke result:

```text
principal login: status=200, token_present=true, cust_id=C26002, parent_cust_id=C26002
distributor login: status=200, token_present=true, cust_id=C260020001, parent_cust_id=C26002
```

## Endpoint smoke

### Principal token — child distributor scope

Request shape:

```text
GET /v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001
GET /v1/reports/secondary-sales/sum-date?month=6&cust_id=C260020001
GET /v1/reports/secondary-sales/group?month=6&year=2026&cust_id=C260020001&group_by=outlet
GET /v1/reports/secondary-sales/group?month=6&cust_id=C260020001&group_by=outlet
```

Result:

```text
sum-date year=2026: status=200, numeric null_fields=[]
sum-date no year: status=200, numeric null_fields=[]
group year=2026: status=200, rows=11, sum_net_sales=5403946260
group no year: status=200, rows=11, sum_net_sales=5403946260
```

Representative `sum-date` values:

```json
{
  "total_gross_sale": 5405350000,
  "total_discount_promo": 1442480,
  "total_ppn": 540394626,
  "net_sales_exc_ppn": 5403946260,
  "net_sales": 5944340886,
  "total_salesman": 5,
  "total_outlet": 9,
  "total_product": 12,
  "qty": 48,
  "qty_return": 1,
  "return_rate": 2.083333333333333,
  "net_sales_return": 630630
}
```

Representative first outlet group row:

```json
{
  "id": 479,
  "code": "BMI260029",
  "name": "Herman Lee > TK Mawar Melati",
  "net_sales": 5000000000
}
```

### Distributor token — own scope and denied parent scope

Request shape:

```text
GET /v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001
GET /v1/reports/secondary-sales/group?month=6&year=2026&cust_id=C260020001&group_by=outlet
GET /v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C26002
GET /v1/reports/secondary-sales/group?month=6&year=2026&cust_id=C26002&group_by=outlet
```

Result:

```text
distributor own sum-date: status=200, numeric null_fields=[]
distributor own group: status=200, rows=11, sum_net_sales=5403946260
distributor parent sum-date: status=403, message="cust_id is outside authorized scope"
distributor parent group: status=403, message="cust_id is outside authorized scope"
```

## Acceptance evidence

- Dashboard summary now returns non-empty numeric data from local source tables for `C260020001`, Juni 2026.
- Group endpoint now returns 11 outlet rows locally and total grouped `net_sales` matches summary `net_sales_exc_ppn` (`5403946260`).
- FE-compatible request without explicit `year` still returns data for current year 2026.
- Distributor token can access own `cust_id` and is denied parent `cust_id`, preserving scope enforcement.
- No tokens, credentials, or secrets were written to source/evidence.

## Remaining note

`SecondarySalesReportReturnSumReportByMonth` remains fact-based but is not used as the authoritative source for dashboard `sum-date`/`group` data after this remediation. It can be simplified separately if desired.
