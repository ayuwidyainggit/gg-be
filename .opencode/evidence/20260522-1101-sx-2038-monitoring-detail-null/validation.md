# Validation — SX-2038 Monitoring Detail Null Data

Task ID: `20260522-1101-sx-2038-monitoring-detail-null`
Date: 2026-05-22
Mode: local docker services + direct DB verification

## Scope

Validasi dilakukan pada service lokal docker `system` dan `pjp`, dengan DB dev yang user berikan. Tidak menyimpan password DB atau JWT di artifact ini.

## Local services

Container aktif:
- `scylla-system` → `0.0.0.0:9001->9001/tcp`
- `scylla-pjp` → `0.0.0.0:9010->9010/tcp`

Action:
- restart `system` dan `pjp` via compose agar Air compile kode terbaru
- login ke `http://127.0.0.1:9001/v1/users/login`
- hit endpoint `http://127.0.0.1:9010/api/v1/monitoring_locations/details`

## Login verification

Request user:
- email: `princessa@gmail.com`
- password: provided by user in chat

Result:
- login lokal berhasil
- JWT berhasil diterbitkan
- user mapping di DB:
  - `sys.m_user.user_id = 140`
  - `sys.m_user.cust_id = C26002`
  - `sys.m_user.emp_id = 380`

## DB verification

### Salesman mapping

Query ringkas:

```sql
SELECT emp_id, sales_name, cust_id
FROM mst.m_salesman
WHERE emp_id IN (484, 482);
```

Result:
- `482 | Jihan Fahira | C26002`
- `484 | Syaiful | C26002`

### PJP principal valid

Query ringkas:

```sql
SELECT id, pjp_code, salesman_id, cust_id, approval_status
FROM pjp_principles.permanent_journey_plans
WHERE salesman_id IN (482, 484)
  AND approval_status IN ('Approved', 'Need Review');
```

Result:
- `id=62, pjp_code=1265, salesman_id=482, cust_id=C26002, approval_status=Approved`
- `id=64, pjp_code=1224, salesman_id=484, cust_id=C26002, approval_status=Approved`

### destinations_history pada 2026-05-22

Query ringkas:

```sql
SELECT dh.pjp_id, dh.destination_id, dh.destination_code, dh.destination_type, dh.is_extra_call, dh.cust_id, dh.date::date
FROM pjp_principles.destinations_history dh
JOIN pjp_principles.permanent_journey_plans pjp ON pjp.id = dh.pjp_id
WHERE pjp.salesman_id IN (482, 484)
  AND dh.date::date = '2026-05-22';
```

Result penting:
- ada 16 row untuk salesman `482` dan `484`
- terdapat row `is_extra_call = true` untuk kedua salesman
- row principal reguler + extra-call memang hidup di `pjp_principles.destinations_history`

## Old query vs new query

### Old principal detail query (route_pop_permanent + destinations)

Untuk `emp_id=482`, query lama return:
- `0 rows`

Ini cocok dengan symptom bug `data: null`.

### New principal detail query (destinations_history)

Untuk `emp_id=482`, query baru return:
- `plan = 10`
- `on_going = 6`
- `extra_call = 1`
- `visited = 0`
- `total_skip = 5`

Untuk `emp_id=484`, query baru return:
- `plan = 3`
- `on_going = 0`
- `extra_call = 2`
- `visited = 2`
- `total_skip = 0`

## Local endpoint verification

### Request 1

```http
GET http://127.0.0.1:9010/api/v1/monitoring_locations/details?emp_id=482&date=2026-05-22
Authorization: Bearer <local login token>
```

Response summary:
- `message = Success`
- `data[0].visit_information.emp_id = 482`
- `planned = 10`
- `on_going = 6`
- `extra_call = 1`
- `visited = 0`
- `skipped = 5`

### Request 2

```http
GET http://127.0.0.1:9010/api/v1/monitoring_locations/details?emp_id=484&date=2026-05-22
Authorization: Bearer <local login token>
```

Response summary:
- `message = Success`
- `data[0].visit_information.emp_id = 484`
- `planned = 3`
- `on_going = 0`
- `extra_call = 2`
- `visited = 2`
- `skipped = 0`

## Match check

Hasil endpoint lokal cocok dengan hasil query DB baru untuk kedua salesman:
- `482` → cocok
- `484` → cocok

## Build and tests

Dari module `pjp/`:
- `rtk go build ./...` → success
- `rtk go test ./service/live_monitoring/... -v -run TestGetMonitoringDetail` → 5 passed
- `rtk go test ./...` → 49 passed, 0 failed, 48 packages

## Quality gate

`@quality-gate` verdict: `PASS_WITH_RISKS`

Residual risk utama:
- JOIN `ovl.outlet_id = dh.destination_id OR ovl.outlet_code = dh.destination_code` perlu dipantau untuk potensi double-count pada data anomali
- semantik `visited`/`on_going` principal bisa ditinjau lagi bila schema `pjp_principles.outlet_visit_list` memakai kombinasi `finish` / `skip_at` berbeda dari distributor

## Conclusion

Fix SX-2038 tervalidasi di local docker service terhadap DB dev:
- bug lama ter-reproduce via query lama (`0 rows`)
- query baru menghasilkan summary valid
- endpoint lokal `monitoring_locations/details` sekarang mengembalikan `data` non-null untuk `emp_id=482` dan `emp_id=484` pada `2026-05-22`
