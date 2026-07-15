# Evidence Discovery — SX-2034 Extra Call tidak muncul di Monitoring Activity

Task ID: `20260521-1515-sx2034-extra-call-monitoring`
Tanggal: 2026-05-21
Sumber: Jira SX-2034, kode lokal `scylla-be`.

## Files inspected

- `mobile/controller/m_outlet.go` (route `POST /from-list`, handler `CreateFromList`)
- `mobile/entity/extra_call_outlet.go` (payload extra call)
- `mobile/service/m_outlet.go` (`StoreFromList`, branch `IsDistributor` vs principle)
- `mobile/repository/m_outlet.go` (`StoreFromList`, `StoreFromListPrinciple`, `StoreFromListOutletVisitList`, `StoreFromListOutletVisitListPrinciple`)
- `mobile/model/m_outlet.go` (`MOutletCreadFromList`)
- `mobile/repository/pjp_principal.go` (pola JOIN outlet/distributor lewat `destination_id`)
- `mobile/service/pjp_principal.go` (pola insert `destinations_history` saat submit PJP — sudah benar mengisi `DestinationId`)
- `pjp-principle/model/destination_history.go` (`DestinationID int`)
- `pjp-principle/repository/destination_history/*` (CreateBulk dll, sudah benar kalau struct diisi)
- `pjp-principle/service/pjp_enhance/create_service.go` (insert pakai struct `DestinationHistory{DestinationID: destination.ID, ...}` — benar)
- `pjp/router/live_monitoring.go` (`GET /live-monitoring-principal`)
- `pjp/controller/live_monitoring/get_principal_controller.go`
- `pjp/service/live_monitoring/get_principal_service.go`
- `pjp/repository/live_monitoring/get_principal_repository.go`
- `pjp/model/live_monitoring.go`
- `pjp/data/response/live_monitoring_response.go`
- `pjp/repository/live_monitoring/get_distributor_repository.go` (pola distributor: sudah join `destinations_history` + cabang extra call via `roh.is_extra_call`)
- `pjp/service/live_monitoring/get_distributor_service.go` (split `RouteData` vs `ExtraCallData` lewat `row.IsExtraCall`)

## Smoking gun #1 — Insert `destination_id` salah kolom (Task 1)

`mobile/repository/m_outlet.go::StoreFromListPrinciple` (lihat baris 274-297):

INSERT urutan kolom:
```
route_code, route_name, verified_date, "date", week, "year", index_day, start_week, is_in_current_year, is_additional,
destination_id, destination_code, destination_status, destination_name, destination_address, destination_type, longitude, latitude,
pjp_id, pjp_code, old_pjp_id, old_pjp_code, old_route_code, old_route_name, photo, avg_sales_week, cust_id, is_extra_call
```

Bind values:
```
$1  outlet.RouteCode
$2  outlet.RouteName
$3  outlet.VerifiedDate
$4  outlet.Date
$5  outlet.Week
$6  outlet.Year
$7  outlet.IndexDay
$8  outlet.StartWeek
$9  outlet.IsInCurrentYear
$10 outlet.IsAdditional
$11 outlet.OldPjpId        <-- BUG: kolom destination_id diisi OldPjpId
$12 outlet.OutletCode
$13 outlet.OutletStatus
$14 outlet.OutletName
$15 outlet.OutletAddress
$16 "outlet"               (literal)
$17 outlet.Longitude
$18 outlet.Latitude
$19 outlet.PjpId
$20 outlet.PjpCode
$21 outlet.OldPjpId        (kolom old_pjp_id, benar)
$22 outlet.OldPjpCode
$23 outlet.OldRouteCode
$24 outlet.OldRouteName
$25 outlet.Photo
$26 outlet.AvgSalesWeek
$27 outlet.CustId
$28 outlet.IsExtraCall
```

`outlet.OldPjpId` selalu nil pada extra call (lihat `mobile/service/m_outlet.go::StoreFromList` baris 281-304: field `OldPjpId` tidak di-set di `MOutletCreadFromList`). Hasil: `pjp_principles.destinations_history.destination_id = NULL` setiap kali extra call principal dibuat. JOIN ke `mst.m_outlet`/`mst.m_distributor` lewat `destination_id` jadi gagal.

Catatan: `StoreFromList` (versi distributor schema `pjp.route_outlet_history`) bind `outlet.OutletId` ke kolom `outlet_id` — itu benar dan tidak terdampak. Bug hanya di varian `StoreFromListPrinciple`.

`StoreFromListOutletVisitListPrinciple` baris 318-333 sudah meng-insert `outlet_id = outlet.OutletId` dengan benar; jadi `outlet_visit_list` tidak ikut bermasalah.

## Smoking gun #2 — Live monitoring principal tidak punya jalur extra call (Task 2)

`pjp/repository/live_monitoring/get_principal_repository.go::GetPrincipalMonitoring` baris 81-138:
- JOIN dari `pjp_principles.permanent_journey_plans` ke `pjp_principles.routes` ke `pjp_principles.destinations` (BUKAN `destinations_history`) ke `outlet_visit_list`.
- `pjp_principles.destinations` adalah master destinasi PJP yang tidak menyimpan extra call (extra call ditulis ke `destinations_history` saja, sesuai perilaku `mobile/service/m_outlet.go::StoreFromList`).
- Tidak ada cabang `is_extra_call`, dan tidak ada UNION yang mengambil row extra call.
- `LiveMonitoringPrincipalRow` belum punya field `IsExtraCall`, dan response `LiveMonitoringPjpData` punya `ExtraCallData` namun selalu kosong dari jalur principal.
- Bandingkan dengan `pjp/repository/live_monitoring/get_distributor_repository.go` baris 128-158: distributor pakai `roh.is_extra_call` + LEFT JOIN ke `outlet_visit_list ovl` dengan `ovl.is_extra_call = roh.is_extra_call`, lalu `pjp/service/live_monitoring/get_distributor_service.go` membagi ke `RouteData`/`ExtraCallData`. Pola ini yang harus direuse untuk principal.

## Project patterns relevan (Reuse-first)

- Pattern split route vs extra call: `pjp/service/live_monitoring/get_distributor_service.go` (`row.IsExtraCall` → `extraRouteMap`/`pjp.ExtraCallData`).
- Pattern resolve outlet vs distributor di destination: `mobile/repository/pjp_principal.go` baris 148-176 dan baris 196-219 (`LEFT JOIN mst.m_outlet ON mo.outlet_id = dh.destination_id AND dh.destination_type = 'outlet'`, `LEFT JOIN mst.m_distributor ON md.distributor_id = dh.destination_id AND dh.destination_type = 'distributor'`).
- Pattern correct insert `destinations_history`: `mobile/service/pjp_principal.go::SubmitPjpPrincipal` (set `DestinationId: outlet.OutletId` untuk outlet, `DestinationId: distributor.DistributorId` untuk distributor).

## Constraints

- Multi-module repo. Layer wajib Controller→Service→Repository→DB (`AGENTS.md`).
- Insert tulis-balik harus berada di service-layer transaction; tx-context harus dihormati. `StoreFromListPrinciple` sudah dipanggil dari transaction `mOutletTransaction`, jadi fix cukup di repo + service tanpa mengubah arsitektur.
- Fiber based service untuk `mobile`, `pjp`. Validasi via `entity.ExtraCallOutlet`.
- Schema `pjp_principles.destinations_history` adalah sumber data untuk monitoring principal; query monitoring saat ini mengabaikan tabel itu — perlu re-platform query, bukan hanya tambal kolom.
- `mobile/entity/extra_call_outlet.go::OutletIDs` adalah `[]int`. Belum ada path distributor di endpoint `from-list` (current `IsDistributor` membedakan schema, bukan destination type). Distributor extra call belum punya jalur create yang jelas → masuk Open Question.

## Risks

- Memperbaiki insert tanpa memperbaiki query monitoring tetap menghasilkan `extra_call_data` kosong.
- Memperbaiki query tanpa memperbaiki insert: data baru tetap NULL → kosong.
- Data lama `destination_id IS NULL` perlu backfill aman (idempoten, hanya untuk `is_extra_call = true`).
- Mengganti JOIN principal dari `destinations` → `destinations_history` berisiko regresi pada PJP normal: harus dipertahankan paritasnya (route, destination, koordinat, ovl) — pakai pendekatan UNION atau tetap di `destinations_history` dengan filter date di `dh.date`.
- `salesman_id` di Jira (`482`) → `pjp.salesman_id` adalah `emp_id`. Pastikan join `m_employee.emp_id` tetap konsisten.

## Tools/docs checked

- Lokal repo only. Tidak panggil context7/brave/GitHub MCP karena masalah jelas dari kode lokal dan deskripsi Jira; tidak ada library eksternal yang materially mempengaruhi keputusan.
