# Discovery — SX-2421 & SX-2422 (Arrival Photo in Outlet Details)

Task: tambah `file_url` di response `GET /scylla-pjp/v1/live-monitoring-distributor` dan
`GET /scylla-pjp/v1/live-monitoring-principal`, dengan source of truth `mobile.visits.file_url`.

Mode: maintenance-stability (shortest diff, tidak refactor besar).

## 1. Confirmed vs Assumed Audit

| # | Claim | Source | Level |
|---|---|---|---|
| 1 | Service `pjp` adalah satu-satunya handler untuk endpoint live monitoring | `pjp/router/live_monitoring.go` + `pjp/router/router.go` (lines 92-95) — `pjp-principle` tidak punya handler live monitoring (`grep` no result) | confirmed_repo |
| 2 | Endpoint distributor = `GET /api/v1/live-monitoring-distributor` | `pjp/router/live_monitoring.go:13` | confirmed_repo |
| 3 | Endpoint principal = `GET /api/v1/live-monitoring-principal` | `pjp/router/live_monitoring.go:12` | confirmed_repo |
| 4 | Controller distributor hanya pass-through ke service | `pjp/controller/live_monitoring/get_distributor_controller.go:63` (`c.service.GetDistributorMonitoring`) | confirmed_repo |
| 5 | Controller principal hanya pass-through ke service | `pjp/controller/live_monitoring/get_principal_controller.go:63` | confirmed_repo |
| 6 | Service distributor memanggil 6 query repo (employee scope, monitoring rows, employee/route/outlet meta, **latest visit coordinates**, attendance, current coordinate) | `pjp/service/live_monitoring/get_distributor_service.go:48-141` | confirmed_repo |
| 7 | Service principal memanggil 5 query repo (employee scope, monitoring rows, extra call rows, attendance, current coordinate) — tidak ada explicit `mobile.visits` join | `pjp/service/live_monitoring/get_principal_service.go:38-114` | confirmed_repo |
| 8 | Distributor `arrive_longitude`/`arrive_latitude` saat ini bersumber dari `mobile.visits mv` (latest per day via `ROW_NUMBER() OVER`) dan di-enrich di `enrichDistributorRowsWithLatestVisits` | `pjp/repository/live_monitoring/get_distributor_repository.go:181-223` + `pjp/service/live_monitoring/get_distributor_service.go:120-159` | confirmed_repo |
| 9 | Principal `arrive_longitude`/`arrive_latitude` regular route bersumber dari `pjp_principles.outlet_visit_list ovl` (bukan `mobile.visits`) | `pjp/repository/live_monitoring/get_principal_repository.go:103-104` (`COALESCE(CAST(NULLIF(ovl.longitude, '') AS DOUBLE PRECISION), 0) AS arrive_longitude`) | confirmed_repo |
| 10 | Principal extra-call juga ambil `arrive_longitude`/`arrive_latitude` dari `pjp_principles.outlet_visit_list ovl` | `pjp/repository/live_monitoring/get_principal_extra_call_repository.go:47-48` | confirmed_repo |
| 11 | Response shape saat ini: `data[].pjp_data[].{route_data, extra_call_data}[].destination_data[]` dengan field `arrive_longitude`, `arrive_latitude`, dan tanpa `file_url` | `pjp/data/response/live_monitoring_response.go:50-67` | confirmed_repo |
| 12 | Model `LatestVisitCoordinateRow` saat ini hanya berisi `cust_id, emp_code, outlet_code, arrive_longitude, arrive_latitude` (no file_url) | `pjp/model/live_monitoring.go:73-79` | confirmed_repo |
| 13 | Model `LiveMonitoringDistributorRow` tidak punya `file_url`; `LiveMonitoringPrincipalRow` juga tidak | `pjp/model/live_monitoring.go:4-71` | confirmed_repo |
| 14 | Field `file_url` di schema `mobile.visits` ada (dipakai oleh endpoint visit) | `pjp/data/request/visit_request.go:33,58` + `pjp/service/visit_service.go:99` + `pjp/repository/outlet_visit_list_repository.go:592,943,1010` (select dengan `file_url = ?`) | confirmed_repo |
| 15 | Doc `monitoring_activity_be_doc.txt` (di `docs/Monitoring Activity - BE.md`) menyatakan endpoint principal join `mobile.visits` untuk arrive_longitude/arrival photo | `docs/Monitoring Activity - BE.md:437` (query di section 5 principal) | confirmed_repo |
| 16 | Implementasi principal saat ini **tidak** melakukan join `mobile.visits` | `pjp/repository/live_monitoring/get_principal_repository.go` (line 81-148, tidak ada `mobile.visits`) | confirmed_repo |
| 17 | Mismatch (15) vs (16) perlu diputuskan sebelum implementasi principal | doc vs code reality | user-decision-needed |

## 2. Source-of-truth Mismatch (PENTING)

`docs/Monitoring Activity - BE.md` section 5 (Location Monitoring Principal) masih menulis
query dengan:

```sql
left join mobile.visits v
    on v.emp_code = pjp.salesman_code
   and v.outlet_code = mo.outlet_code
   and v.created_at::date = '2026-05-21'
...
v.longitude as arrive_longitude,
v.latitude as arrive_latitude
```

Kenyataan di `pjp/repository/live_monitoring/get_principal_repository.go` adalah
`arrive_longitude`/`arrive_latitude` regular route + extra-call sekarang diambil dari
`pjp_principles.outlet_visit_list ovl.longitude/latitude`. Doc tidak match dengan implementasi
yang sudah refactor.

Untuk `file_url` ada dua opsi semantik, dan keduanya valid; default difavoritkan adalah opsi A
agar sesuai dengan keinginan Jira/source-of-truth yang eksplisit
(`response field file_url` dari `mobile.visits.file_url`):

- **A. Ambil `file_url` dari `mobile.visits` (sesuai doc + Jira) — implementation default.**
  Tidak mengubah `arrive_longitude`/`arrive_latitude` (tetap dari `ovl`); tambah join `mobile.visits v`
  di principal repo dengan partisi `cust_id, emp_code, destination_code/outlet_code` per date, dan pilih
  `v.file_url` dari row terbaru. `mobile.visits` punya `outlet_code` (lihat query distributor lines 181-200),
  sedangkan `pjp_principles` punya `d.destination_code` (lihat principal repo line 113). Join dapat
  menggunakan `v.outlet_code = d.destination_code`.
- **B. Ambil `file_url` dari `pjp_principles.outlet_visit_list.file_url`** (kolom sudah ada —
  `pjp/model/outlet_visit_list_principle.go:36`). Opsi ini lebih minimal diff tetapi menyalahi instruksi
  Jira/source-of-truth dan doc.

**Rekomendasi:** opsi A dengan kondisi:
1. Tidak mengubah asal `arrive_longitude`/`arrive_latitude` (tetap dari `ovl`).
2. Tambahkan jalur enrichment terpisah untuk `file_url` mengikuti pola distributor
   (`GetPrincipalLatestVisitFileURL` map), tidak mencampur dengan `GetPrincipalMonitoring`
   (untuk menjaga plan query kecil & tidak mengganggu planning/Pagination, sesuai
   catatan critical invariants: pagination employee scope + semua destination rows tanpa LIMIT).

Jika user memilih opsi B, plan ini hanya perlu satu kolom select di repo + struct + response —
sangat kecil.

## 3. Path Map (yang akan disentuh)

Distributor (sudah memakai `mobile.visits`):
1. `pjp/model/live_monitoring.go` — tambah `FileURL *string` di `LatestVisitCoordinateRow` dan
   `LiveMonitoringDistributorRow`.
2. `pjp/repository/live_monitoring/get_distributor_repository.go`:
   - `GetDistributorLatestVisitCoordinates` base query: tambah `mv.file_url AS file_url` di select.
   - Wrapper select: tambah `file_url`.
3. `pjp/service/live_monitoring/get_distributor_service.go`:
   - `enrichDistributorRowsWithLatestVisits`: assign `rows[index].FileURL = visitCoordinate.FileURL`.
   - `transformDistributorRows`: assign `FileURL: row.FileURL` di `LiveMonitoringDestinationData`.
4. `pjp/data/response/live_monitoring_response.go` — tambah `FileURL *string \`json:"file_url"\``
   di `LiveMonitoringDestinationData`.
5. `pjp/service/live_monitoring/get_distributor_service_test.go` — tambah assertion `FileURL` di test
   `TestEnrichDistributorRowsWithMetadata` (atau test baru khusus enrichment visit file_url).

Principal (Opsi A — pakai `mobile.visits.file_url`):
1. `pjp/model/live_monitoring.go` — tambah `FileURL *string` di `LiveMonitoringPrincipalRow`.
2. `pjp/repository/live_monitoring/get_principal_repository.go`:
   - Tambah `LEFT JOIN mobile.visits v ON v.emp_code = me.emp_code AND v.cust_id = me.cust_id
     AND v.outlet_code = d.destination_code AND v.created_at >= ? AND v.created_at < ?`
     (pakai `buildLiveMonitoringDayRange` pattern yang sudah ada di distributor repo).
   - Tambah `COALESCE(MAX(v.file_url) FILTER (WHERE ROW_NUMBER()...), '')` — ini lebih aman
     dipakai CTE/window function; lihat detail di plan utama (subbagian TDD/Test Plan).
3. `pjp/repository/live_monitoring/get_principal_extra_call_repository.go`:
   - Sama: tambah join `mobile.visits` untuk `dh.destination_code`, dan select `v.file_url`.
4. `pjp/service/live_monitoring/get_principal_service.go`:
   - `transformPrincipalRows`: assign `FileURL: row.FileURL` di `LiveMonitoringDestinationData`.
5. `pjp/data/response/live_monitoring_response.go` — field `FileURL` sudah di distributor; reuse.
6. `pjp/service/live_monitoring/get_principal_service_test.go` — tambah assertion pada test
   `TestTransformPrincipalRows_DoesNotDuplicateCrossJoinedDestinationsAfterRepoFix` atau test baru.

## 4. Files Inspected

- `pjp/router/live_monitoring.go`
- `pjp/router/router.go` (prefix `/api/v1`)
- `pjp/controller/live_monitoring/get_distributor_controller.go`
- `pjp/controller/live_monitoring/get_principal_controller.go`
- `pjp/service/live_monitoring/get_distributor_service.go`
- `pjp/service/live_monitoring/get_principal_service.go`
- `pjp/service/live_monitoring/get_distributor_service_test.go`
- `pjp/service/live_monitoring/get_principal_service_test.go`
- `pjp/repository/live_monitoring/get_distributor_repository.go`
- `pjp/repository/live_monitoring/get_principal_repository.go`
- `pjp/repository/live_monitoring/get_principal_extra_call_repository.go`
- `pjp/model/live_monitoring.go`
- `pjp/data/response/live_monitoring_response.go`
- `pjp/data/request/visit_request.go` (cross-check field `file_url`)
- `pjp/model/outlet_visit_list_principle.go` (cross-check `file_url` di `pjp_principles.outlet_visit_list`)
- `docs/Monitoring Activity - BE.md` (lines 403-573 — endpoint principal + distributor spec)

## 5. Constraints From Project Harness

- `AGENTS.md` (root) menyebutkan service `pjp` Gin-based — confirmed. Service dir =
  `/Users/ujang/Projects/Geekgarden/scylla-be/pjp/`.
- Tidak ada `PROJECT_STACK.md`/`PROJECT_COMMANDS.md`/`FRAMEWORK_PLAYBOOK.md` di `.opencode/docs/`
  repo ini, sehingga default validasi = `rtk go mod download && rtk go mod tidy && rtk go test ./...`
  per service `pjp` (lihat AGENTS.md root + service `pjp/Makefile` ketika perlu).
- `pjp` adalah service Gin (di-confirmed). Pattern validasi = `cd pjp && rtk go test ./service/live_monitoring/... ./repository/live_monitoring/... ./model/... ./data/response/...`.
- Repo tidak punya HTTP/integration test untuk live monitoring; validasi curl manual dengan
  Bearer token lokal/env terpisah (lihat .opencode/docs/SECURITY.md).

## 6. Open Questions / Decisions Needed

1. **Opsi A atau B untuk principal?** (lihat §2). Default = A (mobile.visits.file_url sesuai Jira & doc).
2. Jika A: ada kemungkinan beberapa `mobile.visits` row per (emp_code, outlet_code, date).
   Apakah implementasi boleh memilih "any non-empty file_url" (mirip distributor yang ambil row_number=1)
   atau harus "visit terbaru"? Default = pilih row_number=1 (visit terbaru) untuk konsistensi
   dengan distributor.
