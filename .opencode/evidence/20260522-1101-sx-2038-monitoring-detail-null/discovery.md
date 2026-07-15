# Evidence Discovery — SX-2038 Live Monitoring Detail Null

Task ID: `20260522-1101-sx-2038-monitoring-detail-null`
Sprint: SX Sprint 13. Env: staging (`https://best.scyllax.online`).
Endpoint: `GET /scylla-pjp/api/v1/monitoring_locations/details?emp_id=484&date=2026-05-22`.
Reproduction response: `{ "data": null, "message": "No Data", "request_id": "..." }`.

## Files Inspected

- `pjp/router/live_monitoring.go` — route registration baris 14: `router.GET("/monitoring_locations/details", controller.GetMonitoringDetail)`.
- `pjp/controller/live_monitoring/get_detail_controller.go` — handler endpoint detail.
- `pjp/data/request/live_monitoring_request.go` — `LiveMonitoringDetailRequest{ EmpID int, DistributorID *int, Date string }`. `distributor_id` opsional → kalau NULL berarti jalur **Principal**.
- `pjp/service/live_monitoring/get_detail_service.go` — `GetMonitoringDetail` (1-120) + `getPrincipalVisitInfo` (122-189) + `getDistributorVisitInfo` (191-269) + helper transform.
- `pjp/service/live_monitoring/get_detail_service_test.go` — stub repo + test hanya menutup jalur Distributor + expense filter; **belum ada test Principal happy path**.
- `pjp/repository/live_monitoring/get_detail_repository.go` — semua repo method untuk endpoint detail (visit info, counters distributor, sales/return/expense/shipment, activity time, distributor info, user fullname, child custIDs, salesman cust id).
- `pjp/repository/live_monitoring/get_principal_repository.go` — implementasi list principal sebagai pembanding (statuses param + filter region/area/distributor/empIDs).
- `pjp/repository/live_monitoring/get_principal_extra_call_repository.go` — loader extra-call principal (baca `pjp_principles.destinations_history` dengan `is_extra_call = true`).
- `pjp/repository/live_monitoring/live_monitoring_repository.go` — interface repo (memuat `GetSalesmanCustID(ctx, tx, empID)` baris 46).
- `pjp/constant/pjp_constant.go` — `MsgNoData = "No Data"`, `ApprovalStatusApproved = "Approved"`.
- `.opencode/plans/20260521-1515-sx2034-extra-call-monitoring.md` — plan bug terkait yang menargetkan list principal (paritas dengan distributor).
- `plans/live-monitoring-debug-plan.md` — referensi historis untuk endpoint serupa.
- `postman/Scylla-Live-Monitoring-Complete.postman_collection.json` — kontrak response yang diharapkan FE (`data` adalah array of object, ada `visit_information`, `sales`, `return`, `collection`, `expense`, `shipment`).

## Smoking Gun (Repo Detail Principal)

`pjp/repository/live_monitoring/get_detail_repository.go::GetVisitInformationPrincipal` (baris 13-56):

```go
tx.WithContext(ctx).Table("pjp_principles.permanent_journey_plans pjp").
  Select(`me.emp_id, me.emp_code, me.emp_name,
          COUNT(d.destination_code) AS plan,
          COUNT(CASE WHEN ovl."start"  IS NOT NULL THEN 1 END) AS on_going,
          COUNT(CASE WHEN ovl.finish    IS NOT NULL THEN 1 END) AS visited,
          COUNT(CASE WHEN ovl.skip_at   IS NOT NULL THEN 1 END) AS total_skip,
          COUNT(ovl.outlet_code) AS matched`).
  Joins("JOIN pjp_principles.route_pop_permanent rpp ON rpp.pjp_id = pjp.id").
  Joins("JOIN pjp_principles.routes r ON r.route_code = rpp.route_code").
  Joins("JOIN pjp_principles.destinations d ON d.route_code = r.route_code").
  Joins("JOIN mst.m_salesman ms2 ON pjp.salesman_id = ms2.emp_id").
  Joins("JOIN mst.m_employee me ON me.emp_id = ms2.emp_id").
  Joins("JOIN smc.m_customer mc ON ms2.cust_id = mc.cust_id").
  Joins("JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id").
  Joins("LEFT JOIN pjp_principles.outlet_visit_list ovl ON ovl.pjp_id = pjp.id AND ovl.outlet_code = d.destination_code AND DATE(ovl.date) = ?", date).
  Where("pjp.salesman_id IN (SELECT emp_id FROM mst.m_salesman ms WHERE ms.cust_id IN ?)", custIDs).
  Where("DATE(rpp.date) = ?", date).
  Where("me.emp_id = ?", empID).
  Where("pjp.approval_status = ?", constant.ApprovalStatusApproved).
  Group("me.emp_id, me.emp_code, me.emp_name")

err := query.Take(&result).Error
if errors.Is(err, gorm.ErrRecordNotFound) { return nil, nil }
```

Dipakai oleh service `getPrincipalVisitInfo`, dan kalau `nil` → `GetMonitoringDetail` mengembalikan `nil` → controller membungkus `{data: null, message: "No Data"}` (`get_detail_controller.go` baris 81-87).

## Mapping Penting (klarifikasi prompt)

- `pjp.salesman_id == mst.m_salesman.emp_id` (lihat `get_detail_repository.go:36-37` dan referensi list principal). Artinya `req.EmpID` yang dikirim FE sebagai `emp_id` dipakai langsung sebagai `pjp.salesman_id` tanpa lookup tambahan. Tidak ada bug “mapping `emp_id` → `salesman_id`” di code path saat ini — keduanya sama domainnya.
- `helper.GetCurrentCustomerId` dan `helper.GetCurrentUserId` dipakai untuk JWT scoping (parent cust_id + user). `GetChildCustIDs(parentCustID)` mengumpulkan self + children → daftar `custIDs` untuk filter tenant.
- Detail Principal hanya menerima 1 status (`Approved`), sedangkan jalur **Distributor** sudah `IN ('Approved','Need Review')` (lihat baris 115, 147 file yang sama). List principal (`get_principal_repository.go` + `get_principal_extra_call_repository.go`) menerima `statuses []string` dari controller.

## Hipotesis Root Cause (di-rangking)

1. **Filter `pjp.approval_status = 'Approved'` saja** — kalau PJP `salesman_id=484/482` pada `2026-05-22` masih di status `Need Review`, query return 0 rows → response `null`. Query referensi dev.fe pakai `IN ('Approved','Need Review')`. **High-probability**.
2. **INNER JOIN `pjp_principles.destinations d ON d.route_code = r.route_code`** — kalau salesman tidak punya `destinations` record (mis. PJP cuma berisi extra-call yang ditangani via `destinations_history`, atau data destinations belum di-seed untuk route tersebut), query 0 rows. **High-probability** karena issue ini bertetangga dengan SX-2034 (`destinations_history.destination_id` NULL & extra-call hilang dari list).
3. **Mismatch `emp_id` 484 vs 482** — curl QA pakai `emp_id=484`, komentar dev.fe pakai `salesman_id=482`. Karena `pjp.salesman_id = m_salesman.emp_id`, kalau memang data Princessa adalah `482` lalu FE/QA hit dengan `484`, hasilnya 0 rows. **Medium-probability**, butuh konfirmasi `m_salesman` & `m_user` mapping di staging.
4. **JOIN `smc.m_customer mc ON ms2.cust_id = mc.cust_id` + `JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id` (INNER)** — kalau cust_id salesman tidak punya `distributor_id` (mis. user principal murni), join drop semua row. **Medium-probability**.
5. **`childCustIDs` kosong** — kalau parent cust_id user tidak punya child mapping di `smc.m_customer`, `GetChildCustIDs` cuma return parent itu sendiri; lalu subquery `m_salesman.cust_id IN (childCustIDs)` mungkin tidak match `salesman.cust_id` aktual. **Low-medium-probability**, butuh cek data.
6. **`destination_id` NULL di `destinations_history` (SX-2034)** — Tidak relevan untuk endpoint **detail** principal saat ini, karena query detail tidak join ke `destinations_history`. Relevansi datang kalau perbaikan beralih ke `destinations_history` (rekomendasi paritas dengan list).

## Project Patterns ditemukan

- Layer Controller→Service→Repository→DB ketat. Detail tidak punya transaksi (read-only).
- `Take()` + `gorm.ErrRecordNotFound` → return `(nil, nil)` adalah konvensi “no data”.
- Statuses PJP yang valid untuk reporting: `Approved`, `Need Review` (lihat distributor branch + list principal default).
- Tenant scoping pakai `parent_cust_id` (parent + child) dari helper `GetChildCustIDs`.
- List principal sudah pisah dua loader (regular destinations + extra-call destinations_history). Pola itu siap dipakai di detail.

## Reuse Candidates

- `liveMonitoringRepository.GetPrincipalExtraCalls` — bisa direplika menjadi `GetVisitInformationPrincipalExtraCalls` (count-only) untuk menambahkan kontribusi extra-call ke `plan/extra_call/visited/skipped/on_going` di detail.
- Pola SELECT count `IN ('Approved','Need Review')` yang sudah dipakai jalur distributor di file yang sama.
- `LiveMonitoringRequest` (list endpoint) menerima `statuses []string` dari client; pertimbangkan menambah opsi serupa di detail bila perlu, atau cukup hardcode `IN ('Approved','Need Review')` agar selaras.

## Constraints

- Service Fiber/Gin (`pjp` Gin-based), GORM, schema multi-tenant, schema prefix `pjp_principles.*`, `pjp.*`, `mst.*`, `smc.*`, `mobile.*`, `tms.*`.
- Tidak boleh modifikasi schema DB. Backfill data `destinations_history` ditangani plan SX-2034, jangan diduplikasi.
- Test command: `rtk go test ./...` di `pjp/`.

## Risiko

- Mengubah filter `approval_status` ke `IN ('Approved','Need Review')` mengubah behaviour endpoint untuk semua user → cek apakah kontrak FE sebelumnya menyatakan hanya `Approved`. Sumber Postman/Docs FE mengindikasikan dua status diterima (selaras dev.fe & jalur distributor) → low risk.
- Mengganti INNER JOIN `destinations` ke loader `destinations_history` perlu hati-hati supaya tidak men-double-count plan vs extra-call. Mitigasi: loader terpisah untuk extra-call (sesuai pola SX-2034) atau switch sumber data ke `destinations_history` keseluruhan tapi kontrol `is_extra_call` di SUM.
- Kalau hipotesis utama adalah data salesman murni `Need Review`, fix paling minimal cukup mengubah filter status → segera resolve untuk Princessa tanpa menyentuh SX-2034.

## Commands & Source Checked

- Local search: `grep monitoring_locations`, `grep destinations_history`, `grep approval_status`, `grep ApprovalStatusApproved` (lihat tool log).
- Repo docs: `.opencode/docs/AGENT_ROUTING.md`, `.opencode/docs/ARCHITECTURE.md`, `.opencode/docs/QUALITY.md`, `.opencode/plans/20260521-1515-sx2034-extra-call-monitoring.md`.
- DB queries: belum dijalankan di staging — diserahkan ke implementer (`@fixer`) atau QG dengan akses DB.
- MCP/external: tidak digunakan; root cause sepenuhnya internal.
