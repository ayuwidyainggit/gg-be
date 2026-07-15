# DB + cURL Validation — SX-1917 Payment Deposit Report List

## Remote DB target

- Host: `103.28.219.73`
- Port: `25431`
- DB: `scylla_citus_dev`
- User: `postgres`
- Password: provided by user for this validation only; not persisted in source files.
- SSL mode: `disable`

## Service setup

- Docker compose check:
  - `rtk docker compose -f docker-compose.yml ps`
  - Result: command succeeded but no services were running/listed.
- Local finance service was started manually with env overrides:
  - `SERVER_HOST=127.0.0.1`
  - `SERVER_PORT=19005`
  - DB env pointed to remote DB.
- Health check:
  - `rtk curl -sS -i http://127.0.0.1:19005/ping`
  - Result: `HTTP/1.1 200 OK`, body `It works`.
- Service stopped after validation.

## Schema/data availability checks

Data sample counts by `cust_id` showed usable test data for `C260020001`:

- `acf.deposit`: min date `2026-04-14`, max date `2026-05-05`, total `7`.
- `acf.account_payable_payment`: min date `2026-04-24`, max date `2026-04-27`, total `3` where `deleted_by IS NULL`.

Column check confirmed relevant columns exist, including:

- `acf.deposit`: `cust_id`, `deposit_no`, `deposit_date`, `emp_id`, `deleted_at`, `deleted_by`.
- `acf.deposit_payment`: `cust_id`, `deposit_no`, `pay_type`, `payment_amount`, `deleted_at`.
- `acf.deposit_expense`: `cust_id`, `deposit_no`, `payment_amount`, `deleted_at`, `deleted_by`.
- `acf.account_payable_payment`: `cust_id`, `account_payable_payment_no`, `account_payable_payment_date`, `deleted_by`, `deleted_at`.
- `acf.account_payable_payment_options`: `cust_id`, `account_payable_payment_no`, `pay_type`, `payment_amount`.

## Direct DB query validation

Validation range:

- `cust_id = C260020001`
- `start_date = 2026-04-24`
- `end_date = 2026-04-27`

Aggregated AR/AP union query returned:

| deposit_type | rows | cash | cheque | transfer | return | credit_debit | expense | total_payment |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| AP | 3 | 20000000.0000 | 0 | 5000000.0000 | 0 | 0 | 0 | 25000000.0000 |
| AR | 2 | 22788000.0000 | 0 | 2662000.0000 | 0 | 0 | 500000.0000 | 24950000.0000 |

Sample rows returned:

| deposit_date | type | deposit_no | collector_id | collector_code | collector_name | cash | transfer | expense | total_payment |
|---|---|---|---:|---|---|---:|---:|---:|---:|
| 2026-04-24 | AR | DP2604240001 | 421 | EMP0025 | Piere Njangka | 11038000.0000 | 2662000.0000 | 500000.0000 | 13200000.0000 |
| 2026-04-24 | AR | DP2604240002 | 381 | EMP0022 | Phill Jones | 11750000.0000 | 0 | 0 | 11750000.0000 |
| 2026-04-24 | AP | PY2604240001 | NULL | NULL | NULL | 10000000.0000 | 0 | 0 | 10000000.0000 |
| 2026-04-24 | AP | PY2604240002 | NULL | NULL | NULL | 8000000.0000 | 5000000.0000 | 0 | 13000000.0000 |
| 2026-04-27 | AP | PY2604270001 | NULL | NULL | NULL | 2000000.0000 | 0 | 0 | 2000000.0000 |

Formula verified from data:

- AR row `DP2604240001`: `11038000 + 0 + 2662000 + 0 + 0 - 500000 = 13200000`.
- AP row `PY2604240002`: `8000000 + 0 + 5000000 + 0 + 0 - 0 = 13000000`.

## cURL validation

Authentication:

- Generated a temporary local HS256 smoke token using existing service `JWT_SECRET_KEY` from local `.env` and payload with:
  - `cust_id = C260020001`
  - `parent_cust_id = C26002`
  - `emp_id = 381`
  - `expires = now + 1 hour`
- Token value was not written to evidence.

Saved response artifacts under temp directory:

- `/var/folders/2r/vlp3bhfn1cl5cpdhp23rkh440000gn/T/opencode/sx1917-curl/ar.json`
- `/var/folders/2r/vlp3bhfn1cl5cpdhp23rkh440000gn/T/opencode/sx1917-curl/ap.json`
- `/var/folders/2r/vlp3bhfn1cl5cpdhp23rkh440000gn/T/opencode/sx1917-curl/union.json`
- `/var/folders/2r/vlp3bhfn1cl5cpdhp23rkh440000gn/T/opencode/sx1917-curl/deposit_no.json`

### AR only

Request:

```http
GET http://127.0.0.1:19005/finance/v1/reports/payment-deposit?deposit_type=AR&start_date=2026-04-24&end_date=2026-04-27&emp_id=381,421&page=1&limit=10&sort=deposit_date:asc
```

Result:

- HTTP `200`
- `items = 2`
- `total_data = 2`
- `total_page = 1`
- `grand_total = 24950000`
- First row collector populated: `collector_id=421`, `collector_code=EMP0025`, `collector_name=Piere Njangka`.

### AP only

Request:

```http
GET http://127.0.0.1:19005/finance/v1/reports/payment-deposit?deposit_type=AP&start_date=2026-04-24&end_date=2026-04-27&emp_id=381,421&page=1&limit=10&sort=deposit_date:asc
```

Result:

- HTTP `200`
- `items = 3`
- `total_data = 3`
- `total_page = 1`
- `grand_total = 25000000`
- First row collector fields are `null`.
- Confirms `emp_id` did not filter AP.

### AR + AP union

Request:

```http
GET http://127.0.0.1:19005/finance/v1/reports/payment-deposit?deposit_type=AR,AP&start_date=2026-04-24&end_date=2026-04-27&emp_id=381,421&page=1&limit=10&sort=deposit_date:asc
```

Result:

- HTTP `200`
- `items = 5`
- `total_data = 5`
- `total_page = 1`
- `grand_total = 49950000`
- Matches DB aggregate: AR `24950000` + AP `25000000`.

### deposit_no filter

Request:

```http
GET http://127.0.0.1:19005/finance/v1/reports/payment-deposit?deposit_type=AR,AP&start_date=2026-04-24&end_date=2026-04-27&deposit_no=DP2604240001,PY2604240001&page=1&limit=10&sort=deposit_no:asc
```

Result:

- HTTP `200`
- `items = 2`
- `total_data = 2`
- `total_page = 1`
- `grand_total = 23200000`
- Confirms branch-specific filter for AR and AP deposit/payment numbers.

### invalid deposit_type

Request:

```http
GET http://127.0.0.1:19005/finance/v1/reports/payment-deposit?deposit_type=XX&start_date=2026-04-24&end_date=2026-04-27&page=1&limit=10
```

Result:

- HTTP `400`
- Response: `{"message":"Bad Request","errors":"invalid deposit_type: XX",...}`.

## Notes / residual risks

- cURL was performed against a locally-run finance service pointed to the remote DB, not a deployed remote service.
- The token was a local smoke token signed with local `.env` secret; this validates middleware and controller/service/repository path locally.
- Response `deposit_type` currently returns code values `AR` / `AP`; prior quality gate already noted possible product/FE ambiguity if full labels are required.
