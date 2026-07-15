# Runtime Evidence SX-2234

Task id: `20260616-1740-sx-2234-sales-trend`
Waktu: `2026-06-16 Asia/Jakarta`

## Runtime status

Compose started from repo root:

```bash
rtk docker compose -f docker-compose.yml up -d
```

Sales service running, verified with:

```bash
curl -sS -i http://localhost:9004/ping
```

Result:

```text
HTTP/1.1 200 OK
...
It works
```

Container check:

```bash
docker ps --format '{{.Names}} {{.Status}} {{.Ports}}'
```

Relevant result:

```text
scylla-sales Up ... 0.0.0.0:9004->9004/tcp
```

## API validation

Token handling:

- Used synthetic local JWT generated at runtime.
- Token was kept in shell variable only.
- Token was not printed, logged, saved, or committed.
- JWT secret was loaded inside command from local container/env and was not printed.

Endpoint called:

```http
GET http://localhost:9004/v1/reports/secondary-sales/trend-sales?year=2026&cust_id=C260020001
```

Result:

```text
HTTP/1.1 200 OK
```

Response data summary:

- `data.length = 12`
- months returned: `1..12`
- month 1: `total_gross_sale=0`, `total_discount_promo=0`, `net_sales=0`
- month 4: `total_gross_sale=535000000`, `total_discount_promo=4382000`, `net_sales=584669800`
- months 7..12: zero values

Sanitized response excerpt:

```json
{
  "message": "",
  "data": [
    {"month":1,"total_gross_sale":0,"total_discount_promo":0,"net_sales":0},
    {"month":4,"total_gross_sale":535000000,"total_discount_promo":4382000,"net_sales":584669800},
    {"month":12,"total_gross_sale":0,"total_discount_promo":0,"net_sales":0}
  ],
  "request_id": "217c3020-7a92-4e8e-9b4d-6c69299b7e63"
}
```

## Direct SQL parity check

Direct SQL run against local `ggn_scyllax`, checking one empty month and one populated month:

- month 1: empty month
- month 4: populated month

Command used `PGPASSWORD=postgres psql -h localhost -U postgres -d ggn_scyllax ...` with source-table formula from implementation.

Result:

```text
1,0,0,0
4,535000000,4382000.0000,584669800
```

Parity:

| Month | API `total_gross_sale` | SQL `total_gross_sale` | API `total_discount_promo` | SQL `total_discount_promo` | API `net_sales` | SQL `net_sales` | Result |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 0 | 0 | 0 | 0 | 0 | 0 | pass |
| 4 | 535000000 | 535000000 | 4382000 | 4382000.0000 | 584669800 | 584669800 | pass |

## Risk closure

Closed risk from previous quality gate:

- live runtime/API-direct-SQL validation absent

Closure evidence:

- API returned HTTP 200.
- API returned 12 month rows in order.
- Empty month parity pass.
- Populated month parity pass.
- `total_discount_promo` source-table formula parity pass.

## Remaining risk

None known for SX-2234 scope.
