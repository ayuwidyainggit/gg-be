# Discovery — SX-2016 Monitoring Principal kosong / salah

## DB validation (read-only, BEGIN READ ONLY/COMMIT)

- `princessa@gmail.com` → `sys.m_user.user_id=140`, `cust_id=C26002`, `emp_id=380`.
- `smc.m_customer.C26002.parent_cust_id=C26002` (self-parent).
- Salesman `emp_id=482 / emp_code=MS9990 / Jihan Fahira` → `cust_id=C26002` (sama dengan login user).
- `pjp_principles.permanent_journey_plans` untuk salesman 482: `id=62, pjp_code=1265, approval_status=Approved, cust_id=C26002`.
- `pjp_principles.outlet_visit_list` di `date='2026-05-20'` untuk pjp_id=62: 5 baris (id 135..139). Yang punya `arrive_at` + lat/long aktual:
  - id=136 outlet `BMI260003` `arrive_at=1779236763189`, `leave_at=1779237085474`, lat 37.421998, long -122.084.
  - id=137 outlet `BMI260004` `arrive_at=1779268286992`, lat -6.252474, long 106.818921.
- `pjp_principles.destinations` untuk `route_code=7010`: 5 baris (id 1159..1163). 4 di antaranya punya `longitude=1, latitude=2` (placeholder), hanya `162612 PT Besi Makmur` punya koordinat asli.

## Smoke endpoint

Login (system service):
```
POST https://best.scyllax.online/scylla-system/api/v1/users/login
```
JWT principal: `cust_id=C26002`, `emp_id=380`, `parent_cust_id=C26002`.

`GET https://best.scyllax.online/scylla-pjp/api/v1/live-monitoring-principal?date=1779278400&status[]=Approved&status[]=Need+Review&emp_id=482`:
- Response 200 dengan satu employee 482, satu route 7010, **10 destinasi semuanya `BMI260005`** (destination yang sama berulang).
- Tidak ada `BMI260003` / `BMI260004` yang punya arrive aktual; tidak ada `arrive_longitude/arrive_latitude` (selalu 0).

## Root cause (revised)

1. JOIN salah di `pjp/repository/live_monitoring/get_principal_repository.go`:
   - `LEFT JOIN pjp_principles.outlet_visit_list ovl ON pjp.pjp_code = ovl.pjp_code`.
   - Tidak ada filter `ovl.date`, tidak ada relasi destinasi → OVL. Jadi tiap row destinasi di-cross-join dengan semua OVL di pjp_code itu (5 destinasi × 11 OVL = 55 row).
   - Hasil: cross product, bukan satu OVL per destinasi.
2. Pagination service-level pakai `Limit(limit).Offset(offset)` di SQL untuk row mentah, bukan di level employee. Default `limit=10` motong 55 row jadi 10 row pertama, semuanya `destination_id=1159 (BMI260005)` karena ORDER BY `d.id`.
3. Field `arrive_longitude/arrive_latitude` untuk principal tidak pernah di-SELECT/diset di repository/service principal. Distributor sudah punya, principal belum.
4. Bug child cust_id (rencana awal) **tidak relevan** untuk case ini, karena `princessa.cust_id == salesman.cust_id == C26002`. Hardening child resolve tetap defensible (parity dengan distributor) tapi bukan akar masalah SX-2016.
5. Bug data master: 4 dari 5 destinasi punya placeholder `(1,2)`. Itu masalah master `pjp_principles.destinations`, bukan bug query — tapi otomatis bikin titik dashboard “salah lokasi”. Tidak in-scope SX-2016 fix kode.

## Fix candidate (DB-verified)

Ubah JOIN OVL principal jadi:

```sql
LEFT JOIN pjp_principles.outlet_visit_list ovl
  ON ovl.pjp_id = pjp.id
 AND ovl.date = DATE(rpp.date)
 AND ovl.outlet_code = d.destination_code
```

Hasil di DB (2026-05-20, salesman 482): 5 row, satu per destinasi, dengan `arrive_at/leave_at` benar untuk `BMI260003` & `BMI260004` (mengandung lat/long mobile).

Tambahkan ke SELECT (principal repo + model + response transform):
- `ovl.leave_at`,
- `COALESCE(CAST(NULLIF(ovl.longitude,'') AS DOUBLE PRECISION),0) AS arrive_longitude`,
- `COALESCE(CAST(NULLIF(ovl.latitude,'') AS DOUBLE PRECISION),0) AS arrive_latitude`.

Pagination harus di-page di level employee, bukan di level row mentah, supaya destinasi satu employee tidak terpotong (pola distributor service sudah benar; principal harus diselaraskan).

## Constraints + risks

- Layer service→repo→DB harus dipertahankan.
- Tidak boleh lewat batas `pjp` module saat fix.
- Jangan log/store secrets/JWT/DB password ke artifact.
- Risiko regresi: dashboard frontend mungkin bergantung pada bentuk response saat ini (10 row destinasi duplikat). Need check FE expectation singkat saat QA.
- Risiko data master: titik peta untuk destinasi placeholder `(1,2)` tetap salah; perlu ticket terpisah ke data team.

## Files inspected

- `pjp/router/live_monitoring.go`
- `pjp/middleware/jwt.go`, `pjp/helper/current.go`
- `pjp/controller/live_monitoring/get_principal_controller.go`
- `pjp/service/live_monitoring/get_principal_service.go`
- `pjp/service/live_monitoring/get_distributor_service.go` (pattern reuse)
- `pjp/repository/live_monitoring/get_principal_repository.go`
- `pjp/repository/live_monitoring/get_distributor_repository.go` (pattern reuse)
- `pjp/repository/live_monitoring/live_monitoring_repository.go` (interface)
- `pjp/model/live_monitoring.go`
- `pjp/data/response/live_monitoring_response.go`
- `pjp/service/visit_service.go` (arrive write path)
- `pjp/repository/outlet_visit_principle/update_visit_list_column_at.go`
- `system/service/user_service.go` (login endpoint route)
