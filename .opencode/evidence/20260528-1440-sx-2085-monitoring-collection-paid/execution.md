# Execution Evidence SX-2085

Task ID: `20260528-1440-sx-2085-monitoring-collection-paid`

## Local runtime

- `docker ps` menunjukkan `scylla-system`, `scylla-master`, `scylla-sales`, `scylla-redis`, `scylla-rabbitmq` aktif.
- `scylla-pjp` dinaikkan via `rtk docker compose -f docker-compose.yml up -d pjp`.
- Health local PJP: `GET http://127.0.0.1:9010/api/v1/health` → `{"code":200,"message":"OK"}`.
- Health system login source implicit OK karena login berhasil.

## Local DB access

- `PGPASSWORD=postgres psql -h 127.0.0.1 -p 5432 -U postgres -d ggn_scyllax -c "select current_database(), current_user;"`
- Result: database `ggn_scyllax`, user `postgres`.

## QA sample from prompt

- SQL sample `date=2026-05-28`, `emp_id=421` on local DB returned `0 rows`.
- Expected aggregation SQL for same sample also returned `0 rows`.
- Conclusion: prompt sample data not present in local DB snapshot; local functional validation used real local tenant/sample data below.

## Local user login

- Login request: `POST http://127.0.0.1:9001/v1/users/login`
- Payload used:

```json
{"email":"dist@sda.idetama.id","password":"admin"}
```

- Login success for local user `dist@sda.idetama.id` / `cust_id=C220010001` / `distributor_id=67`.

## Local validation sample actually present

Chosen sample from local DB:
- `cust_id=C220010001`
- `emp_id=228`
- `date=2026-05-22`
- `distributor_id=67`

Distributor mapping SQL result:
- `emp_id=228`
- `sales_name=Charles Leclerc`
- `distributor_id=67`
- `distributor_code=3434`
- `distributor_name=Distributor iDetama`

## Duplicate-risk SQL

SQL:

```sql
SELECT d.deposit_no, d.cust_id, COUNT(DISTINCT dd.invoice_no) AS invoice_count,
       COUNT(dp.*) AS payment_join_count, COALESCE(SUM(dp.payment_amount),0) AS joined_payment_sum
FROM acf.deposit d
JOIN acf.deposit_detail dd ON dd.deposit_no = d.deposit_no AND dd.cust_id = d.cust_id
LEFT JOIN acf.deposit_payment dp ON dp.deposit_no = d.deposit_no AND dp.cust_id = d.cust_id
WHERE d.deposit_date = DATE '2026-05-22'
  AND d.collection_no IS NOT NULL
  AND d.emp_id = 228
  AND d.cust_id = 'C220010001'
GROUP BY d.deposit_no, d.cust_id
ORDER BY d.deposit_no;
```

Result:

```text
DP2605220003 | C220010001 | invoice_count=1 | payment_join_count=1 | joined_payment_sum=200000.0000
```

## Expected aggregation SQL

SQL final sesuai rule `payment per invoice`:

```sql
WITH payment_per_invoice AS (
    SELECT deposit_no, cust_id, invoice_no, SUM(COALESCE(payment_amount, 0)) AS payment_amount
    FROM acf.deposit_payment
    GROUP BY deposit_no, cust_id, invoice_no
)
SELECT mo.outlet_id, mo.outlet_code, mo.outlet_name,
       SUM(COALESCE(ppi.payment_amount, 0)) AS collection_total
FROM acf.deposit d
JOIN acf.deposit_detail dd ON dd.deposit_no = d.deposit_no AND dd.cust_id = d.cust_id
JOIN sls."order" o ON o.invoice_no = dd.invoice_no AND o.cust_id = dd.cust_id
LEFT JOIN payment_per_invoice ppi
    ON ppi.deposit_no = dd.deposit_no
   AND ppi.cust_id = dd.cust_id
   AND ppi.invoice_no = dd.invoice_no
JOIN mst.m_outlet mo ON mo.outlet_id = o.outlet_id AND mo.cust_id = o.cust_id
WHERE d.deposit_date = DATE '2026-05-22'
  AND d.collection_no IS NOT NULL
  AND d.emp_id = 228
  AND d.cust_id = 'C220010001'
GROUP BY mo.outlet_id, mo.outlet_code, mo.outlet_name
ORDER BY mo.outlet_code;
```

Result:

```text
outlet_id=913 | outlet_code=B000009 | outlet_name=BOEDY 9 | collection_total=200000.0000
```

Conclusion:
- allocation sekarang mengikuti `deposit_payment.invoice_no`
- payment dijumlah per invoice dulu, lalu outlet mengikuti invoice di `sls.order`
- ini sesuai kebutuhan FE: payment per outlet

## Re-test setelah revisi query ke invoice grain

Commands:

```bash
rtk go test ./service/live_monitoring -run 'TestGetMonitoringDetail_.*Collection'
rtk go test ./...
```

Observed outcomes:
- `Go test: 3 passed in 1 packages`
- `Go test: 52 passed in 48 packages`

Re-validated SQL + API local sample:
- SQL invoice-grain result tetap `B000009 / BOEDY 9 / 200000.0000`
- API local tetap return `collection_total: 200000`
- Match.

Residual risk updated:
- risk utama lama `deposit-level payment duplicated across outlets` sudah dihapus oleh query invoice-grain
- risk sisa sekarang pindah ke kualitas data source, misal duplicate abnormal di `deposit_detail` untuk kombinasi `deposit_no,cust_id,invoice_no`, atau invoice detail tanpa pasangan payment row
- sesuai klarifikasi business, kasus itu dianggap issue data deposit, bukan rule alokasi payment

## API validation

## API validation

Request:

```text
GET http://127.0.0.1:9010/api/v1/monitoring_locations/details?emp_id=228&distributor_id=67&date=2026-05-22
Authorization: Bearer <local token from dist@sda.idetama.id>
```

Response snippet:

```json
{
  "message": "Success",
  "data": [
    {
      "visit_information": {
        "activity_date": "2026-05-22",
        "company_name": "Distributor iDetama",
        "company_code": "3434",
        "level": "Distributor",
        "emp_id": 228,
        "collection_summary": {
          "count": 1,
          "status": "completed"
        }
      },
      "collection": [
        {
          "outlet_id": 913,
          "outlet_code": "B000009",
          "outlet_name": "BOEDY 9",
          "collection_total": 200000
        }
      ]
    }
  ]
}
```

Comparison:
- API `collection[0].collection_total = 200000`
- SQL expected `collection_total = 200000.0000`
- Match.

## Tests run

From `pjp/`:

```bash
rtk go test ./service/live_monitoring -run 'TestGetMonitoringDetail_.*Collection'
rtk go test ./service/live_monitoring && rtk go test ./...
```

Observed outcomes:
- First targeted run failed before stub updates because other test stubs lacked `GetCollections`.
- After fix:
  - `Go test: 3 passed in 1 packages`
  - `Go test: 24 passed in 1 packages`
  - `Go test: 52 passed in 48 packages`

## Notes

- Local DB sample from prompt (`2026-05-28`, `421`) unavailable; local validation used real existing tenant data.
- Query still has business ambiguity if one deposit maps to multiple outlets; current code follows safe per-deposit aggregation then groups per outlet, but if source data truly spans many outlets per deposit, business allocation rule may still be needed.
