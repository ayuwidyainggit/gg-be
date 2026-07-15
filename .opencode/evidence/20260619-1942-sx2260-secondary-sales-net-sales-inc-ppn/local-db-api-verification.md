# Local DB/API Verification SX-2260

## Runtime
- DB: local PostgreSQL `ggn_scyllax` via `127.0.0.1:5432`, user `postgres`.
- Docker services checked:
  - `scylla-system` up on `9001`
  - `scylla-sales` up on `9004`
- Login user: `adminbm@gmail.com`; token not stored.

## Login
- `POST http://localhost:9001/v1/users/login`
- Result: HTTP 200, `cust_id=C260020001`, `parent_cust_id=C26002`.

## DB formula verification
Filter:
- `cust_id='C260020001'`
- `invoice_date >= '2026-06-01'`
- `invoice_date < '2026-07-01'`
- `data_status IN (6,7)`

DB totals:
- order rows: 26
- return rows: 1
- order old exclude PPN total: `5459620890`
- order new include PPN total: `6005582979`
- order PPN total: `545962089`
- return old exclude PPN total: `-630630`
- return new include PPN total: `-693693`
- return PPN total: `63063`
- final expected net_sales include PPN: `6004889286`

Formula check:
- `6005582979 + (-693693) = 6004889286`
- selisih dari old formula = `545962089 - 63063 = 545899026`.

## Local sales API verification
Correct local route prefix is `/v1/reports`, not `/sales/v1/reports` when hitting service container directly.

Endpoint calls:
- `GET http://localhost:9004/v1/reports/secondary-sales/group?month=6&year=2026&cust_id=C260020001&group_by=outlet` -> HTTP 200, rows 9, sum net_sales `6004889286`.
- `GET http://localhost:9004/v1/reports/secondary-sales/group?month=6&year=2026&cust_id=C260020001&group_by=salesman` -> HTTP 200, rows 7, sum net_sales `6004889286`.
- `GET http://localhost:9004/v1/reports/secondary-sales/group?month=6&year=2026&cust_id=C260020001&group_by=product_category` -> HTTP 200, rows 2, sum net_sales `6004889286`.
- `GET http://localhost:9004/v1/reports/secondary-sales/group?month=6&year=2026&cust_id=C260020001&group_by=product` -> HTTP 200, rows 14, sum net_sales `6004889286`.

Top API rows:
- outlet: `{id:1841, code:"BMI260029", name:"TK Mawar Melati", net_sales:5521175000}`
- salesman: `{id:479, code:"BM300", name:"Herman Lee", net_sales:5765755886}`
- product_category: `{id:74, code:"01", name:"Mainan", net_sales:5878871900}`
- product: `{id:10745, code:"AF-003", name:"Action Figure Messi", net_sales:5524568500}`

## Result
PASS. Local DB expected total equals local API group totals for all four `group_by` branches.
