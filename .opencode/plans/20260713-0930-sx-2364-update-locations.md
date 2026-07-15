# Plan: SX-2364 [BE] Get Location Timeline (`/v1/monitoring_locations/update-locations`)

- task_id: `20260713-0930-sx-2364-update-locations`
- module: `pjp` (Gin + GORM, `scyllax-pjp` Go 1.21)
- lane: `@artifact-planner` → handoff to `@orchestrator` and `@fixer` for implementation
- mode: maintenance-stability (single bounded endpoint slice)

## Goal

Tambahkan endpoint baru `GET /v1/monitoring_locations/update-locations` di service `pjp` yang mengembalikan timeline lokasi satu karyawan per hari. Sumber timeline: `mobile.attendances` (clock-in/clock-out), `sys.user_location` (gps), dan branch outlet-visit (`pjp.outlet_visit_list` untuk distributor; `pjp_principles.outlet_visit_list` untuk principal). Endpoint dipakai FE untuk polling 5-menit pada halaman Location Monitoring, sehingga memperbaiki data titik koordinat harian yang sebelumnya kosong.

## Non-goals

- Tidak ada publish event ke message broker; FE polling manual per 5 menit.
- Tidak ada WebSocket hub, tidak ada consumer `web-monitoring-location`, tidak ada publisher `location.updated`. Sudah didefer ke plan `20260630-1441-sx-2361-latest-location.md` (sibling scope).
- Tidak memperbaiki bug coordinate-reset pada list-monitoring distributor (DI-2 di `.opencode/plans/20260630-1441-sx-2361-latest-location.md`). Ticket terpisah.
- Tidak menambah field baru ke response `LiveMonitoringData`/`LiveMonitoringPaging`/`LiveMonitoringDetailData`. Endpoint ini berdiri sendiri.
- Tidak menyentuh Swagger docs di luar endpoint baru ini.
- Tidak membuat publisher/consumer baru; tidak menambah dependensi.

## Scope

In-scope:
- Service: `pjp`
- 1 route baru, 1 request DTO, 1 service method, 1 repo method, 1 response struct.
- 1 file `*_test.go` (controller, sqlmock untuk repository) mengikuti pola `pjp/controller/live_monitoring/get_distributor_controller_test.go` dan `pjp/repository/live_monitoring/get_detail_repository_test.go`.
- Resolver principal vs distributor via 1 query `mst.m_employee` dengan scope tenant dari JWT.

Out of scope (next slice):
- Tambah branch di `pjp_principles.permanent_journey_plans` untuk history-based timeline.
- Tambah field `accuracy` di `sys.user_location` jika belum ada.
- Background job pre-aggregating daily timeline ke cache.
- Bug fix list-monitoring distributor (DI-2).

## Requirements

1. Route: `GET /v1/monitoring_locations/update-locations` di bawah router yang sudah ada. Tambahkan entry di `pjp/router/live_monitoring.go`.
2. Auth: middleware `pjp/middleware/jwt.go` sudah aktif untuk group route yang sama. Tidak ada perubahan middleware.
3. Query parameters:
   - `emp_id` (int, required, `binding:"required"`). Merujuk `pjp.salesman_id`.
   - `date` (string, optional, format `YYYY-MM-DD` Asia/Jakarta). Default ke hari ini (`Asia/Jakarta`) jika kosong.
4. Response body mengikuti contoh pada dokumen Point 8, tanpa pagination:
   - `message` (string, `"Success"` atau `"No Data"`).
   - `data.timeline` (seluruh timeline karyawan untuk satu hari, terurut `recorded_at ASC`, `sequence` mulai `1`).
   - `request_id` (UUID v4 string).
5. Item timeline:
   - `sequence` (int, urut ASC).
   - `type` (string, salah satu dari `clock_in`, `clock_out`, `arrive`, `leave`, `gps`).
   - `latitude` (float64, `0` jika sumber kosong).
   - `longitude` (float64, `0` jika sumber kosong).
   - `destination_id` (*int64, nullable).
   - `destination_type` (*string, nullable; hanya berisi untuk outlet principal karena `pjp_principles.outlet_visit_list` punya `destination_type`).
   - `destination_name` (*string, nullable).
   - `recorded_at` (RFC3339 string).
6. Resolver principal vs distributor:
   - Query `SELECT cust_id FROM mst.m_employee WHERE emp_id = $1 AND cust_id IN (SELECT cust_id FROM smc.m_customer WHERE cust_id = $2 OR parent_cust_id = $2) LIMIT 1`. Parameter kedua adalah JWT `custID` (currentCustomerId). Parameter pertama adalah `emp_id` request.
   - Jika baris tidak ada → `404 Not Found` (jangan bocorkan).
   - Jika `cust_id` panjangnya `> 6` → distributor; selain itu → principal.
   - Definisi "panjang > 6" mengikuti jawaban user: `length(cust_id) > 6`. Hardcode tidak diperbolehkan, gunakan `LENGTH(cust_id) > 6` di query.
7. Union query timeline:
   - **Sumber A — `mobile.attendances`**: filter `emp_code = (SELECT emp_code FROM mst.m_employee WHERE emp_id = $empID LIMIT 1)`; range tanggal. `type = 1` → `clock_in`; `type != 1` → `clock_out`. Latitude/longitude dari kolom attendance, `created_at` jadi `recorded_at`.
   - **Sumber B — `sys.user_location`**: filter `emp_id = $empID`, `created_at::date = $date`. `type = 'gps'`, `recorded_at = created_at`.
   - **Sumber C — outlet-visit** (branch):
     - Distributor: `pjp.outlet_visit_list` join `pjp.permanent_journey_plans ON pjp_id`. Cross join lateral `arrive_at/leave_at` jadi dua row, koordinat dari `latitude/longitude` dan `leave_latitude/leave_longitude`. `type` `arrive`/`leave`. `destination_id = outlet_id`. `destination_type` selalu `null` (kolom ini hanya di principal, lihat `pjp/model/outlet_visit_list_principle.go`).
     - Principal: `pjp_principles.outlet_visit_list` join `pjp_principles.permanent_journey_plans ON pjp_id`. Cross join lateral seperti di atas, `destination_type` diisi dari kolom `destination_type` (kolom ini di distributor juga ada di tabel, lihat `pjp/model/outlet_visit_list.go:48`, jadi distributor dapat di-set null by code). `destination_name` join ke `mst.m_outlet outlet_name WHERE outlet_id`.
   - Urutan: union semua sumber lalu `ORDER BY recorded_at ASC, type ASC, destination_id ASC` (tie-breaker untuk stabilitas FE).
8. Tenant scope:
   - Sumber A & C difilter `cust_id IN (SELECT cust_id FROM smc.m_customer WHERE cust_id = $jwtCust OR parent_cust_id = $jwtCust)`. Pola ini sudah dipakai di `pjp/repository/pjp/get_destination_details_repository.go:64`.
   - Sumber B (`sys.user_location`) tidak punya `cust_id` di PK, scope via `emp_id` saja (tidak ada filter tambahan; lihat `mobile/model/user_location.go:7` yang punya `cust_id`—tetap filter `ul.cust_id IN (...)` untuk konsistensi tenant).
9. Output:
    - SQL mengembalikan seluruh row union untuk satu `emp_id` dan satu hari, tanpa `LIMIT`, `OFFSET`, atau query `count`.
10. Service interface baru:
     - `GetUpdateLocations(ctx, req request.UpdateLocationsRequest, custID string) (response.UpdateLocationsResponse, error)`.
     - Mengembalikan struct dengan `Timeline []TimelineItem`.
11. Repository interface baru:
     - `GetUpdateLocations(ctx, tx, empID int, date string, jwtCust string, branch string) ([]model.UpdateLocationRow, error)`.
     - `GetEmployeeRole(ctx, tx, empID int, jwtCust string) (string, error)` (returns `cust_id` string, error if not found).

## Acceptance Criteria

- AC1: Hit `GET /v1/monitoring_locations/update-locations?emp_id=479&date=2026-07-08` dengan JWT valid. Response `200` dengan body sesuai kontrak Point 8 (Message, Data.timeline, RequestID). Timeline urut `recorded_at ASC`, tanpa `paging`.
- AC2: Hit tanpa `emp_id` → response `400 BAD_REQUEST` dengan body error Gin `gin.H{message, request_id}` (mengikuti pola `get_distributor_controller.go:55-58`).
- AC3: Hit tanpa `Authorization` header → `401 Unauthorized` dari middleware (sudah ada di `pjp/middleware/jwt.go`).
- AC4: `emp_id` valid tapi `mst.m_employee` tidak punya baris untuk tenant caller → `404 Not Found` (jangan bocorkan ke tenant lain).
- AC5: `emp_id` valid dan role resolver menandai `cust_id` length > 6 → query outlet branch ke `pjp.outlet_visit_list`. Test: `emp_id` 479 dengan `cust_id` distributor, timeline berisi `arrive`/`leave` dari `pjp.outlet_visit_list` saja.
- AC6: `emp_id` valid dan `cust_id` length ≤ 6 → query outlet branch ke `pjp_principles.outlet_visit_list`. Timeline `destination_type` terisi.
- AC7: `date` kosong → server default ke hari ini `Asia/Jakarta` (sudah ada `loadJakartaLocation` di `service/live_monitoring/live_monitoring_service.go:83`).
- AC8: Tidak ada data untuk emp+date → `200` dengan `data.timeline` array kosong, `message = "No Data"`.
- AC9: Service `go test ./service/live_monitoring/...` dan `go test ./repository/live_monitoring/...` dan `go test ./controller/live_monitoring/...` lulus.
- AC10: Lint/format (`gofmt -l pjp/...`) tanpa diff.
- AC11: Build bersih (`go build ./...` di service `pjp`).

## Existing Patterns/Reuse

- Router: `pjp/router/live_monitoring.go`. Tambah 1 entry.
- Controller pattern: `pjp/controller/live_monitoring/get_distributor_controller.go`. Copy struktur handler + `gin.H{...}` payload, ganti endpoint dan service call.
- Controller interface: `pjp/controller/live_monitoring/live_monitoring_controller.go:10`. Tambah method `GetUpdateLocations(ctx *gin.Context)`.
- Service pattern: `pjp/service/live_monitoring/live_monitoring_service.go`. Extend interface + struct impl.
- Service timezone: `loadJakartaLocation` dan `epochToDateString` (reused).
- Request DTO: `pjp/data/request/live_monitoring_request.go`. Tambah `UpdateLocationsRequest` struct. Pakai `binding:"required"` untuk `emp_id`. Untuk `date` opsional dengan format YYYY-MM-DD, gunakan validator tag `omitempty,datetime=2006-01-02` (lihat library validator).
- Response DTO: `pjp/data/response/live_monitoring_response.go`. Tambah `UpdateLocationsResponse`, `TimelineItem`; jangan pakai `LiveMonitoringPaging`.
- Repository interface: `pjp/repository/live_monitoring/live_monitoring_repository.go`. Tambah 2 method.
- Repository test pattern: `pjp/repository/live_monitoring/get_detail_repository_test.go` (sqlmock, `regexp.QuoteMeta` untuk SQL match).
- Controller test pattern: `pjp/controller/live_monitoring/get_distributor_controller_test.go` (gin recorder, `controllerServiceStub`).
- Tenant scope pattern: `pjp/repository/pjp/get_destination_details_repository.go:64` (`cust_id IN (SELECT cust_id FROM smc.m_customer WHERE cust_id = $1 OR parent_cust_id = $1)`).
- Cross-schema reference: `sys.*`, `mobile.*`, `pjp_principles.*`, `pjp.*`, `mst.*`, `smc.*` sudah dipakai di repo existing. Tidak ada barrier.
- `mst.m_employee`: dipakai di banyak service lain (`sls.*`); pola join sudah lazim.
- `sys.user_location` model: `mobile/model/user_location.go` (sama DB schema, dipanggil dari `pjp`).
- HTTP error convention: `400` untuk binding, `401` untuk JWT, `404` untuk not found (lihat `exception/error_handler.go`). Controller inline `gin.H{...}` (bukan `exception.ErrorHandler`) mengikuti pola `get_distributor_controller.go`.

## Source Anatomy

| Subsystem | Authority inspected | Concrete rule for implementation |
|---|---|---|
| API contract | `Monitoring_Activity_BE.docx`, Point 8 (`textutil` extract lines 2254-2507) + user instruction 2026-07-13 | Route, required `emp_id`, source union, response key names, timeline types, branch distinction; pagination dihapus oleh user. |
| Router | `pjp/router/live_monitoring.go:10-15` | Register only `router.GET("/monitoring_locations/update-locations", controller.GetUpdateLocations)`. |
| Controller | `pjp/controller/live_monitoring/get_distributor_controller.go:33-97` | Get JWT cust ID, create UUID request ID, 30-second context, bind query, return inline `gin.H`. |
| Request/response | `pjp/data/request/live_monitoring_request.go:8-47`; `pjp/data/response/live_monitoring_response.go:72-78` | Reuse form binding; buat DTO response timeline baru tanpa `LiveMonitoringPaging`; jangan mutasi DTO existing. |
| Service / time | `pjp/service/live_monitoring/live_monitoring_service.go:15-103` | Reuse `jakartaLocation`; default/parse business date di service. |
| Tenant access | `pjp/repository/pjp/get_destination_details_repository.go:57-64`; `pjp/repository/live_monitoring/get_detail_repository.go:574-585` | Scope employee and each source to JWT customer plus direct child customer rows. |
| Location source: GPS | `mobile/model/user_location.go:5-16` | Query `sys.user_location` by `emp_id`, tenant `cust_id`, and business-day range; no accuracy field. |
| Location source: distributor visit | `pjp/model/outlet_visit_list.go:7-63` | Convert `arrive_at` and `leave_at` to two rows using stored arrive/leave coordinates. |
| Location source: principal visit | `pjp/model/outlet_visit_list_principle.go:7-61`; Point 8 principal SQL | Use `pjp_principles.outlet_visit_list` branch. Raw select is required because model does not expose `destination_type`. |
| Auth/error | `pjp/middleware/jwt.go`; `pjp/exception/error_handler.go:62-117` | JWT owns 401. Controller query bind owns 400. No global error-handler change. |
| Test | `pjp/controller/live_monitoring/get_distributor_controller_test.go`; `pjp/repository/live_monitoring/get_detail_repository_test.go`; `pjp/go.mod` | Reuse Gin recorder and installed `github.com/DATA-DOG/go-sqlmock v1.5.2`. |

## Reference Map

| Feature | Basis | Why sufficient |
|---|---|---|
| Route and JSON contract | `docs-backed` | Point 8 directly names exact endpoint, query fields, source types, and response example. |
| Attendance timeline rows | `docs-backed` + `repo-backed` | Point 8 SQL maps `mobile.attendances.type` to clock events; PJP already queries this schema in `get_distributor_repository.go`. |
| GPS timeline rows | `docs-backed` + `repo-backed` | Point 8 SQL names `sys.user_location`; model confirms fields and table name. |
| Visit timeline rows | `docs-backed` + `repo-backed` | Point 8 supplies lateral-values mapping; two visit models confirm times and coordinates. |
| Principal/distributor decision | `user_confirmed` | User chose `mst.m_employee.cust_id` length branch. |
| Full-day timeline and ASC order | `user_confirmed` | User menghapus pagination; endpoint mengembalikan seluruh timeline satu hari dalam urutan kronologis. |
| 400 missing emp_id | `user_confirmed` + `repo-backed` | User selected 400; controller binding convention supports it. |
| Tenant scope | `repo-backed` | Existing PJP query pattern scopes self+child customer before data reads. |
| No WebSocket/event changes | `repo-backed` | Sibling plan SX-2361 owns deferred realtime transport scope. |

- `docs-backed`: Point 8 menentukan contract endpoint dan 3 source timeline.
- `repo-backed`: PJP Gin/GORM controller, tenant scope, dan testing pattern adalah authority implementation.
- `user_confirmed`: role rule, tanpa pagination, chronological order, dan error 400 adalah keputusan final user.

### Confirmed vs Assumed Audit

| Material claim | Status | Evidence / rationale |
|---|---|---|
| Endpoint path and Point 8 source union | confirmed_docs | `Monitoring_Activity_BE.docx` extract lines 2254-2507. |
| `sys.user_location` field set | confirmed_repo | `mobile/model/user_location.go:5-16`. |
| Distributor/principal outlet tables and columns | confirmed_repo | `pjp/model/outlet_visit_list*.go`. |
| PJP controller binding returns 400 | confirmed_repo | `get_distributor_controller.go:53-59`. |
| JWT handles missing/invalid token with 401 | confirmed_repo | `pjp/middleware/jwt.go`. |
| Go test dependency sqlmock exists | confirmed_repo | `pjp/go.mod`. |
| Role branch by `mst.m_employee.cust_id` length > 6 | user_confirmed | Question Gate answer, 2026-07-13. |
| No pagination, full-day timeline | user_confirmed | User revisi 2026-07-13: hapus `page`/`limit`, kembalikan seluruh timeline satu hari. |
| Timeline order is `recorded_at ASC` | user_confirmed | Question Gate answer, 2026-07-13. |
| Missing `emp_id` returns 400 | user_confirmed | Question Gate answer, 2026-07-13. |
| Empty date defaults today Jakarta | assumption | Existing live-monitoring time convention; Point 8 does not specify default. |
| Employee outside scope returns 404 generic | assumption | Security-safe behavior; must not be converted into a stated product fact. |
| Distributor `destination_type` is null | assumption | Point 8 distributor SQL selects null although table model has field. |
| Destination name join uses `mst.m_outlet` | assumption | Existing live-monitoring query precedent; verify exact join before coding. |

## Constraints

- C1: Hanya edit file di bawah `pjp/`. Service lain (`mobile/`, `sls/`, `master/`, dll) tidak boleh disentuh.
- C2: Tidak menambah dependency Go baru. Pakai `gin`, `gorm`, `google/uuid`, `DATA-DOG/go-sqlmock` yang sudah ada.
- C3: Tidak menambah tabel DB atau migration.
- C4: Branch identifier principal vs distributor **tidak** boleh di-hardcode ke nilai `cust_id` tertentu; gunakan rule `LENGTH(cust_id) > 6` di query resolver.
- C5: Response JSON tidak boleh bocor `cust_id` internal atau nilai filter tenant.
- C6: File baru mengikuti Go naming convention: snake_case untuk file, CamelCase untuk export.
- C7: SQL parameterized. Tidak ada string concatenation untuk input user (raw query di repo existing sudah pakai `?` placeholder).
- C8: Test coverage minimum: 1 happy path (principal), 1 happy path (distributor), 1 validation failure (no `emp_id`), 1 not-found (emp di luar tenant).
- C9: Tidak ada breaking change pada route yang sudah ada. Hanya tambah 1 entry router.

## Risks

- R1: Definisi `cust_id` length > 6 mungkin rapuh. Field `mst.m_employee.cust_id` adalah string; distributor `cust_id` panjangnya 6 (`C22001XXXXX`) di dokumentasi monitoring. Mitigasi: rule hanya di resolver, dicatat di `Decisions/Assumptions`. Lihat `mobile/pkg/middleware/jwt_middleware.go` contoh token `C220010001` (10 digit) untuk distributor. Panjang tidak konsisten.
- R2: Performance union 3 sources per employee per day. Volume harian rendah (1 emp × 1 day, max ~30 timeline points). Mitigasi: index hint tidak perlu karena tabel sudah dipakai route lain dengan pola serupa.
- R3: SQL kompleks (UNION ALL 3 cabang). Jika Postgres behavior union berbeda dari docs, hasil bisa kosong. Mitigasi: sqlmock test verifikasi regex SQL, dan integration test dengan DB dev (`best.scyllax.online` jika ada env).
- R4: `sys.user_location` mungkin tidak punya index `(emp_id, created_at)`. Mitigasi: fallback ke `ORDER BY created_at` + filter tanggal; jika lambat, defer ke next slice dengan index.
- R5: Kolom `destination_type` di `pjp_principles.outlet_visit_list` belum ada di Go model `OutletVisitListPrinciple`. Mitigasi: query SELECT eksplisit, tidak scan ke model.
- R6: Tenant scope untuk `sys.user_location` (`cust_id IN (jwt_cust OR parent_cust_id)`) mungkin mengecualikan user yang valid. Pola ini konsisten dengan route lain, risiko rendah.

## Decisions/Assumptions

- D1: Role resolver cek `mst.m_employee` dengan `cust_id` dari `smc.m_customer` (parent_cust_id) mengikuti `pjp/repository/pjp/get_destination_details_repository.go:64`. **(user_confirmed)**
- D2: Tidak ada pagination. Endpoint mengembalikan seluruh timeline satu hari (user revisi 2026-07-13, override jawaban sebelumnya). **(user_confirmed)**
- D3: Urutan timeline `recorded_at ASC`. `(user_confirmed)`. Tie-breaker `type ASC, destination_id ASC` ditambahkan planner untuk deterministik FE render.
- D4: Validasi `emp_id` hilang → `400 BAD_REQUEST` mengikuti konvensi PJP. **(user_confirmed)**. Trace: `exception/error_handler.go:110`.
- D5: `date` kosong → default hari ini `Asia/Jakarta`. **assumption**, label `repo_local_evidence` (pola `live_monitoring`).
- D6: Resolver 404 generik saat emp di luar tenant. **assumption**, label `repo_local_evidence` (pola FE lain jika ada) — perlu cek frontend expectation, fallback ke `200 empty timeline` jika FE tidak siap handle `404`.
- D7: Resolver pakai `LENGTH(cust_id) > 6` di SQL, bukan `emp_id > 6 digit` dari dokumen Point 8. **(user_confirmed)**.
- D8: Distributor query `destination_type` selalu `null` (kolom `destination_type` di `pjp.outlet_visit_list` ada tapi per-docs Point 8 distributor tidak set). **assumption**, label `repo_local_evidence` (cek `pjp/model/outlet_visit_list.go:48`).
- D9: `destination_name` join ke `mst.m_outlet.outlet_name`. **assumption**, label `repo_local_evidence` (lihat `pjp/repository/live_monitoring/get_distributor_repository.go:622`).
- D10: `accuracy` tidak ada di kontrak Point 8 → tidak di-serve. **assumption**, label `user_confirmed` (docs Point 8 tidak sebut).

## Execution Source of Truth

Urutan precedence untuk implementer:

1. Dokumen `Monitoring_Activity_BE.docx` Point 8 (kontak kebenaran tertinggi untuk kontrak API).
2. Security/permission rules (`exception/error_handler.go`, `pjp/middleware/jwt.go`).
3. Non-negotiable Implementation Invariants (di bawah).
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria.
6. Implementation Steps.
7. Sibling plan `20260630-1441-sx-2361-latest-location.md` (referensi WS scope, BUKAN untuk outlet 8 endpoint ini).

Konflik yang muncul saat implementasi harus direkam di `.opencode/evidence/20260713-0930-sx-2364-update-locations/verification.md` oleh executor, dan dirujuk ke sini sebelum finalisasi.

## Non-negotiable Implementation Invariants

- NII1: Route path exact `/monitoring_locations/update-locations` (tanpa `/v1` prefix; prefix ditambahkan di router group).
- NII2: Endpoint di bawah middleware JWT yang sudah ada (group route yang sudah di-mount di `pjp/router/router.go`).
- NII3: Output JSON key case mengikuti contoh Point 8 (`snake_case`). Go struct pakai tag `json:"snake_case"`.
- NII4: SQL parameterized (`?` placeholder) untuk semua input user.
- NII5: Tenant scope wajib. Jangan pernah membaca row dari `cust_id` di luar JWT scope. Untuk `sys.user_location` scope via `ul.cust_id IN (...)` walaupun pola `mobile.visits` di repo lain tidak melakukan ini, lakukan untuk konsistensi security posture.
- NII6: Resolver principal vs distributor via `LENGTH(cust_id) > 6` SQL, bukan Go-side heuristic.
- NII7: Tidak ada panic-recover swallow. Service return error, controller map ke 500 (mengikuti pola `get_distributor_controller.go:64-69`).
- NII8: Tidak ada goroutine leak. `context.WithTimeout` 30 detik seperti `get_distributor_controller.go:49`.
- NII9: Test wajib pakai sqlmock untuk repository. Tidak ada test yang bergantung ke DB live.
- NII10: Bug fix list-monitoring distributor (DI-2 di plan SX-2361) **tidak** termasuk scope. Jangan di-bundle.

## Do Not / Reject If

- Jangan pakai `?` query pada `mobile.attendances` dengan string raw. Pakai `tx.WithContext(ctx).Raw(sql, args...).Scan(&rows)` seperti `get_distributor_repository.go:480`.
- Jangan terima `emp_id` dari query body. Hanya query string.
- Jangan expose `cust_id` di response. Field internal.
- Jangan panggil `db.Raw("... " + empID)`. Pakai parameterized.
- Jangan rubah response `LiveMonitoringData`, `LiveMonitoringPaging`, atau `LiveMonitoringDetailData`. Buat struct baru dan jangan embed `LiveMonitoringPaging` di endpoint baru.
- Jangan buat publisher/consumer baru. Tidak relevan.
- Jangan buat WebSocket. Tidak relevan.
- Jangan bundle dengan bug fix `current_coordinate` distributor. Itu SX-2361.
- Jangan hardcode nilai `cust_id` panjang 6 atau 7 di kode. SQL only.
- Jangan tambah file di luar `pjp/`.
- Reject impl jika `service.GetUpdateLocations` lebih dari ~120 baris (sign of over-engineering atau business logic bocor ke controller). Reject jika controller lebih dari ~80 baris.

## Diff Boundary

- Allowed groups:
  - `pjp/router/live_monitoring.go` (1 line added).
  - `pjp/controller/live_monitoring/live_monitoring_controller.go` (1 method signature added to interface).
  - `pjp/controller/live_monitoring/get_update_locations_controller.go` (new).
  - `pjp/controller/live_monitoring/get_update_locations_controller_test.go` (new).
  - `pjp/service/live_monitoring/live_monitoring_service.go` (1 method signature added to interface).
  - `pjp/service/live_monitoring/get_update_locations_service.go` (new).
  - `pjp/repository/live_monitoring/live_monitoring_repository.go` (2 method signatures added).
  - `pjp/repository/live_monitoring/get_update_locations_repository.go` (new).
  - `pjp/repository/live_monitoring/get_update_locations_repository_test.go` (new).
  - `pjp/data/request/live_monitoring_request.go` (1 struct added).
  - `pjp/data/response/live_monitoring_response.go` (2 struct added).
- Disallowed:
  - `pjp/exception/error_handler.go` (no change).
  - `pjp/middleware/jwt.go` (no change).
  - `pjp/router/router.go` (no change — group router sudah ada).
  - `pjp/main.go` (no change — service & controller di-wire di tempat lain; cek `pjp/router/router.go` lebih dulu, jika tidak di-wire di sana, tambahkan wiring di router, BUKAN main).
  - `pjp/docs/docs.go` (auto-generated, jangan edit tangan; biarkan `swag init` regenerate atau skip).
  - File di luar `pjp/`.

- Generated report exceptions: tidak ada.

## TDD / Test Plan

- TDD wajib. Red → Green → Refactor.
- T1 — Repository test (`get_update_locations_repository_test.go`):
   - Red: tulis `TestGetUpdateLocations_Distributor_Success` dengan sqlmock expectation resolver role dan 1 select union tanpa `count`, `LIMIT`, atau `OFFSET`. Pastikan gagal karena method belum ada.
  - Green: implement `GetUpdateLocations` dan `GetEmployeeRole`.
  - Refactor: pecah helper `buildUnionSQL(branch, empID, date, jwtCust)` jika >80 baris.
- T2 — Repository test `TestGetUpdateLocations_Principal_Branch` (memverifikasi union pakai `pjp_principles.outlet_visit_list`).
- T3 — Repository test `TestGetEmployeeRole_NotFound` (returns sql.ErrNoRows atau wrapped error).
- T4 — Controller test (`get_update_locations_controller_test.go`):
  - `TestGetUpdateLocations_Happy`: stub service return 5 timeline items, assert 200 + JSON shape.
  - `TestGetUpdateLocations_MissingEmpID`: stub bind error, assert 400.
  - `TestGetUpdateLocations_NotFound`: stub service return `ErrNotFound`, assert 404 (perlu helper atau constant `NotFoundError`).
  - `TestGetUpdateLocations_NoData`: stub service return empty timeline, assert 200 + message "No Data".
- Validation commands: `cd pjp && go mod download && go mod tidy && gofmt -l . && go vet ./... && go test ./controller/live_monitoring/... ./service/live_monitoring/... ./repository/live_monitoring/...`.

## Implementation Steps

1. Tambah `UpdateLocationsRequest` di `pjp/data/request/live_monitoring_request.go`:
   - Fields: `EmpID int form:"emp_id" binding:"required"`, `Date string form:"date"` (format check via service). Tidak ada `Page` atau `Limit`.
2. Tambah `UpdateLocationsResponse` dan `TimelineItem` di `pjp/data/response/live_monitoring_response.go`:
   - `UpdateLocationsResponse` memuat `Timeline []TimelineItem`, `Message string`, dan `RequestID string`; tidak embed `LiveMonitoringPaging`.
   - `TimelineItem` fields sesuai AC5/AC6/AC7 dengan pointer untuk nullable.
3. Extend `LiveMonitoringRepository` interface di `pjp/repository/live_monitoring/live_monitoring_repository.go` dengan `GetUpdateLocations` dan `GetEmployeeRole`.
4. Buat `pjp/repository/live_monitoring/get_update_locations_repository.go`:
   - `GetEmployeeRole(ctx, tx, empID, jwtCust)` — query `SELECT cust_id FROM mst.m_employee WHERE emp_id = $1 AND cust_id IN (SELECT cust_id FROM smc.m_customer WHERE cust_id = $2 OR parent_cust_id = $2) LIMIT 1`.
   - `GetUpdateLocations(ctx, tx, empID, date, jwtCust, branch)` — branch enum `"principal"` atau `"distributor"`. Bangun UNION ALL 3 sumber dan kembalikan seluruh row urut kronologis; tanpa `count`, `LIMIT`, atau `OFFSET`.
   - Scan ke `[]model.UpdateLocationRow` (struct baru di `pjp/model` atau di file repo).
5. Tambah model `UpdateLocationRow` di `pjp/model/update_location.go` (jika belum ada).
6. Extend `LiveMonitoringService` interface dengan `GetUpdateLocations(ctx, req, custID) (response.UpdateLocationsResponse, error)`.
7. Buat `pjp/service/live_monitoring/get_update_locations_service.go`:
   - Validate `date` format (regex `^\d{4}-\d{2}-\d{2}$`); invalid → return error yang controller map ke 400.
   - Default `date` ke `time.Now().In(jakartaLocation).Format("2006-01-02")` jika kosong.
   - Call `repository.GetEmployeeRole`. Not found → wrap dengan sentinel `ErrNotFound` (lihat `exception/error_handler.go:121`).
   - Determine branch dari `LENGTH(cust_id) > 6` (sudah di-resolve di SQL, service tinggal switch dari string length Go). Atau: tambah kolom `Branch string` di return `GetEmployeeRole`. Pilih kedua agar logika role terpusat di repo.
   - Call `repository.GetUpdateLocations`. Map rows → `[]response.TimelineItem`. Assign `sequence` mulai 1.
   - Tidak ada perhitungan paging.
8. Buat `pjp/controller/live_monitoring/get_update_locations_controller.go`:
   - Pattern copy dari `get_distributor_controller.go`.
   - JWT check via `helper.GetCurrentCustomerId` (401).
   - Bind query ke `request.UpdateLocationsRequest` (400).
   - Call service (500 on error).
   - Build `gin.H{message, data, request_id}`. Tidak ada key `paging`.
9. Extend `LiveMonitoringController` interface dengan `GetUpdateLocations(ctx *gin.Context)`.
10. Tambah route di `pjp/router/live_monitoring.go`:
    - `router.GET("/monitoring_locations/update-locations", controller.GetUpdateLocations)`.
11. Tulis test repository (T1, T2, T3) dengan sqlmock.
12. Tulis test controller (T4) dengan gin recorder + `controllerServiceStub`.
13. Jalankan `gofmt -l pjp/`, `go vet ./...`, `go test ./...` di service `pjp`. Pastikan nol diff dan zero fail.
14. Buat evidence:
    - `.opencode/evidence/20260713-0930-sx-2364-update-locations/route-diff.txt`.
    - `.opencode/evidence/20260713-0930-sx-2364-update-locations/sql-query.txt` (SQL yang dipakai).
    - `.opencode/evidence/20260713-0930-sx-2364-update-locations/test-output.txt`.
    - `.opencode/evidence/20260713-0930-sx-2364-update-locations/curl-trace.txt` (mock curl dengan token redacted).
    - `.opencode/evidence/20260713-0930-sx-2364-update-locations/index.json` (manifest).
    - `.opencode/evidence/20260713-0930-sx-2364-update-locations/discovery.md` (ringkasan bukti repo + refer Point 8).

## Expected Files to Change

- Modified: `pjp/router/live_monitoring.go` (+1 line).
- Modified: `pjp/controller/live_monitoring/live_monitoring_controller.go` (+1 method di interface).
- Modified: `pjp/service/live_monitoring/live_monitoring_service.go` (+1 method di interface).
- Modified: `pjp/repository/live_monitoring/live_monitoring_repository.go` (+2 method di interface).
- Modified: `pjp/data/request/live_monitoring_request.go` (+1 struct).
- Modified: `pjp/data/response/live_monitoring_response.go` (+2 struct).
- New: `pjp/controller/live_monitoring/get_update_locations_controller.go`.
- New: `pjp/controller/live_monitoring/get_update_locations_controller_test.go`.
- New: `pjp/service/live_monitoring/get_update_locations_service.go`.
- New: `pjp/repository/live_monitoring/get_update_locations_repository.go`.
- New: `pjp/repository/live_monitoring/get_update_locations_repository_test.go`.
- New: `pjp/model/update_location.go`.

12 file changes, 11 added, 6 modified. Total LOC di file baru ~400-600 baris (perkiraan; aktual mengikuti pola existing).

## Agent / Tool Routing

- Implementation: `@fixer` (Go backend bounded work).
- Review: `@quality-gate` (security + DB access correctness + JSON shape).
- Architecture question: `@architect` (jika struktur `union query` di-decompose ke CTE).
- Docs help: tidak ada eksternal lib baru, tidak perlu `@librarian`.
- Browser evidence: tidak UI.

## Executor Handoff Prompt

Copy-paste untuk `@orchestrator`:

~~~text
Task: SX-2364 [BE] Get Location Timeline
Plan: .opencode/plans/20260713-0930-sx-2364-update-locations.md
Module: pjp
Source-of-truth: docs/Monitoring_Activity_BE.docx Point 8 (BARU, bukan list-monitoring)

Scope:
- Tambah route GET /v1/monitoring_locations/update-locations di service pjp.
- Timeline items: clock_in, clock_out, gps, arrive, leave.
- 3 UNION ALL sources: mobile.attendances, sys.user_location, branch outlet_visit_list.
- Branch: distributor = pjp.outlet_visit_list; principal = pjp_principles.outlet_visit_list.
- Resolver principal vs distributor via LENGTH(mst.m_employee.cust_id) > 6.
- Tanpa pagination. Kembalikan seluruh timeline satu hari, ordered by `recorded_at ASC`.

Must preserve:
- Tenant scope JWT. cust_id selalu di-restrict ke smc.m_customer self+parent.
- SQL parameterized (?). Tidak ada raw concat.
- Field JSON contract persis contoh Point 8 (snake_case).
- Service GetUpdateLocations tidak lebih dari ~120 LOC.
- 404 generik saat emp di luar tenant. Jangan bocor.

Do not touch:
- File di luar pjp/.
- pjp/exception/error_handler.go, pjp/middleware/jwt.go, pjp/router/router.go.
- pjp/docs/docs.go (auto-generated, skip).
- Bug fix list-monitoring distributor current_coordinate (DI-2 di SX-2361). TIKET TERPISAH.

Validation:
- cd pjp && go mod download && go mod tidy
- gofmt -l pjp/ (expect kosong)
- go vet ./...
- go test ./controller/live_monitoring/... ./service/live_monitoring/... ./repository/live_monitoring/...
- go build ./...

Evidence (wajib di .opencode/evidence/20260713-0930-sx-2364-update-locations/):
- route-diff.txt
- sql-query.txt (union SQL yang dipakai)
- test-output.txt
- curl-trace.txt (token redacted, contoh request+response shape)
- index.json (manifest)
- discovery.md (sumber repo + Point 8 reference)

Return:
- File yang berubah (path list).
- Output test pass/fail.
- Bukti curl shape.
- Task tracker update via `python3 ~/.config/opencode/scripts/task-progress.py 20260713-0930-sx-2364-update-locations --update <TASK_ID> --status <status> --owner @fixer --evidence <path>`.

Claim level: scoped. Claim scope boleh: "endpoint diimplement sesuai plan, semua test pass, evidence disimpan". Claim scope TIDAK boleh: "AC diterima final oleh FE" atau "di-deploy ke production".
~~~

## Execution tracker worklist

1. **W1** | `@fixer` | Add request and response DTOs.
2. **W2** | `@fixer` | Add model and repository interface methods.
3. **W3** | `@fixer` | Implement tenant-scoped update-locations repository query.
4. **W4** | `@fixer` | Add repository sqlmock tests.
5. **W5** | `@fixer` | Implement service date default, resolver, branch, mapping.
6. **W6** | `@fixer` | Add controller request binding and HTTP mapping.
7. **W7** | `@fixer` | Add controller Gin recorder tests.
8. **W8** | `@fixer` | Register protected route.
9. **W9** | `@fixer` | Run formatting, vet, build, and tests.
10. **W10** | `@fixer` | Save validation evidence and finalize manifest.

## Execution-ready Worklist / Handoff Contract

~~~yaml
handoff:
  task_id: 20260713-0930-sx-2364-update-locations
  plan_id: 20260713-0930-sx-2364-update-locations
  caller: orchestrator
  callee: fixer
  scope: Implementasi 1 endpoint + DTO + service + repo + 2 test file di service pjp
  claim_level: scoped
  claim_scope:
    may_claim:
      - "endpoint diimplement sesuai plan"
      - "semua test pass (sqlmock + gin recorder)"
      - "evidence disimpan ke .opencode/evidence/<task-id>/"
    may_not_claim:
      - "AC FE validated end-to-end"
      - "deployed to production"
      - "fix bug list-monitoring distributor (DI-2)"
  source_basis:
    - .opencode/plans/20260713-0930-sx-2364-update-locations.md (this file)
    - docs/Monitoring_Activity_BE.docx Point 8 (line 2254-2507)
    - pjp/router/live_monitoring.go
    - pjp/service/live_monitoring/live_monitoring_service.go
    - pjp/repository/live_monitoring/live_monitoring_repository.go
    - pjp/data/request/live_monitoring_request.go
    - pjp/data/response/live_monitoring_response.go
    - pjp/controller/live_monitoring/get_distributor_controller.go
    - pjp/middleware/jwt.go
    - pjp/exception/error_handler.go
  must_preserve:
    - JWT middleware path
    - Tenant scope pattern cust_id IN (jwt OR parent)
    - snake_case JSON contract
    - parameterized SQL
    - bounded service LOC
  do_not_touch:
    - pjp/exception/error_handler.go
    - pjp/middleware/jwt.go
    - pjp/router/router.go
    - pjp/main.go
    - pjp/docs/docs.go
    - file di luar pjp/
    - bug DI-2 list-monitoring distributor
  validation:
    - "cd pjp && gofmt -l . | wc -l | grep '^0$'"
    - "cd pjp && go vet ./..."
    - "cd pjp && go test ./controller/live_monitoring/... ./service/live_monitoring/... ./repository/live_monitoring/... -count=1"
    - "cd pjp && go build ./..."
  exit_criteria:
    - "gofmt -l . | wc -l = 0"
    - "go vet clean"
    - "go test zero fail"
    - "go build clean"
    - "5+ evidence files ada di .opencode/evidence/20260713-0930-sx-2364-update-locations/"
    - "index.json valid (cek via jq .index.json)"
  evidence_required:
    - .opencode/evidence/20260713-0930-sx-2364-update-locations/route-diff.txt
    - .opencode/evidence/20260713-0930-sx-2364-update-locations/sql-query.txt
    - .opencode/evidence/20260713-0930-sx-2364-update-locations/test-output.txt
    - .opencode/evidence/20260713-0930-sx-2364-update-locations/curl-trace.txt
    - .opencode/evidence/20260713-0930-sx-2364-update-locations/index.json
    - .opencode/evidence/20260713-0930-sx-2364-update-locations/discovery.md
  depends_on: none
  context_bundle:
    verified_by_planner:
      - id: docs-point8-timeline
        fact: "Point 8 mendefinisikan endpoint /v1/monitoring_locations/update-locations dengan 3 UNION ALL sources (mobile.attendances clock_in/out, sys.user_location gps, outlet_visit_list arrive/leave) branched principal vs distributor"
        source: "docs/Monitoring_Activity_BE.docx lines 2254-2507"
        level: confirmed_docs
      - id: schema-user-location
        fact: "sys.user_location punya kolom id, cust_id, emp_id, latitude, longitude, created_at, updated_at (string lat/lng)"
        source: "mobile/model/user_location.go:5-13"
        level: confirmed_repo
      - id: schema-outlet-distributor
        fact: "pjp.outlet_visit_list punya arrive_at, leave_at, latitude, longitude, leave_latitude, leave_longitude, outlet_id; kolom destination_type ada di model tapi docs Point 8 distributor tidak set"
        source: "pjp/model/outlet_visit_list.go:7-60"
        level: confirmed_repo
      - id: schema-outlet-principal
        fact: "pjp_principles.outlet_visit_list tidak punya kolom destination_type di Go model, tapi kolom ada di DB per docs Point 8"
        source: "pjp/model/outlet_visit_list_principle.go:7-58 + docs/Monitoring_Activity_BE.docx query principal"
        level: confirmed_repo
      - id: jwt-401-existing
        fact: "pjp/middleware/jwt.go sudah handle token missing/invalid dengan 401 response"
        source: "pjp/middleware/jwt.go"
        level: confirmed_repo
      - id: controller-convention
        fact: "controller PJP pakai gin.H inline response, BUKAN exception.ErrorHandler; 400 untuk binding, 500 untuk service error, 200 dengan MsgNoData saat data kosong"
        source: "pjp/controller/live_monitoring/get_distributor_controller.go:55-78"
        level: confirmed_repo
      - id: tenant-scope-pattern
        fact: "smc.m_customer.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE cust_id = $1 OR parent_cust_id = $1) adalah pola tenant scope yang dipakai"
        source: "pjp/repository/pjp/get_destination_details_repository.go:64"
        level: confirmed_repo
      - id: no-paging-decision
        fact: "User revisi 2026-07-13: pagination dihapus, endpoint mengembalikan seluruh timeline satu hari tanpa `paging` object dan tanpa `count`/`LIMIT`/`OFFSET`"
        source: "user instruction 2026-07-13 (chat transcript + plan revisi)"
        level: user_confirmed
    files_already_read:
      - pjp/router/live_monitoring.go
      - pjp/router/router.go
      - pjp/data/request/live_monitoring_request.go
      - pjp/data/response/live_monitoring_response.go
      - pjp/controller/live_monitoring/live_monitoring_controller.go
      - pjp/controller/live_monitoring/get_distributor_controller.go
      - pjp/controller/live_monitoring/get_distributor_controller_test.go
      - pjp/service/live_monitoring/live_monitoring_service.go
      - pjp/service/live_monitoring/get_distributor_service.go
      - pjp/repository/live_monitoring/live_monitoring_repository.go
      - pjp/repository/live_monitoring/get_distributor_repository.go
      - pjp/repository/live_monitoring/get_detail_repository.go
      - pjp/repository/live_monitoring/get_detail_repository_test.go
      - pjp/repository/pjp/get_destination_details_repository.go
      - pjp/middleware/jwt.go
      - pjp/exception/error_handler.go
      - pjp/model/outlet_visit_list.go
      - pjp/model/outlet_visit_list_principle.go
      - mobile/model/user_location.go
      - pjp/go.mod
    open_assumptions:
      - D5: date kosong default ke hari ini Asia/Jakarta. Boleh diubah jika FE tidak handle.
      - D6: 404 generik saat emp di luar tenant. Bisa diubah ke 200 empty timeline jika FE design.
      - D8: distributor destination_type selalu null. Bisa di-set jika FE butuh string "Outlet".
      - D9: destination_name join ke mst.m_outlet. Mungkin perlu LEFT JOIN jika outlet di-soft-delete.
    source_of_truth_order:
      - 1. docs/Monitoring_Activity_BE.docx Point 8
      - 2. exception/error_handler.go + jwt middleware security rules
      - 3. Non-negotiable Implementation Invariants (this plan)
      - 4. Execution-ready Worklist (this plan)
      - 5. Acceptance Criteria (this plan)
      - 6. Implementation Steps (this plan)
```

### Worklist tasks (atomic, ordered)

```yaml
worklist:
  - id: W1
    action: Tambah UpdateLocationsRequest + TimelineItem + UpdateLocationsResponse DTO
    depends_on: none
    owner: @fixer
    validation: "cd pjp && gofmt -l ./data/"
    exit_criteria: "2 file modified, 1 file modified (response)"
    blocking: ready
    must_preserve: "JSON tag snake_case"
    do_not_touch: "LiveMonitoringData, LiveMonitoringPaging existing; jangan embed paging di response endpoint baru"
    evidence_update: "evidence/20260713-0930-sx-2364-update-locations/dto-diff.txt"
    exit_verification: "gofmt clean + go build clean"
    start_with: W1
  - id: W2
    action: Extend LiveMonitoringRepository interface + tambah model UpdateLocationRow
    depends_on: W1
    owner: @fixer
    validation: "go build ./repository/live_monitoring/..."
    exit_criteria: "interface + model + file baru"
    blocking: ready
    must_preserve: "interface contract existing"
    do_not_touch: "other repo method"
    evidence_update: "evidence/.../repo-interface-diff.txt"
    exit_verification: "go build clean"
  - id: W3
    action: Implement GetEmployeeRole + GetUpdateLocations repository
    depends_on: W2
    owner: @fixer
    validation: "go test ./repository/live_monitoring/... -run TestGetEmployeeRole"
    exit_criteria: "3 source query (attendance + user_location + outlet), role resolver query"
    blocking: ready
    must_preserve: "tenant scope, parameterized SQL"
    do_not_touch: "file di luar pjp/repository/live_monitoring/"
    evidence_update: "evidence/.../sql-query.txt"
    exit_verification: "go test pass"
  - id: W4
    action: Tulis test sqlmock untuk GetEmployeeRole + GetUpdateLocations (T1, T2, T3)
    depends_on: W3
    owner: @fixer
    validation: "go test ./repository/live_monitoring/... -run 'TestGet(EmployeeRole|UpdateLocations)'"
    exit_criteria: "3 test function pass, coverage 100% untuk file"
    blocking: ready
    must_preserve: "test pattern dari get_detail_repository_test.go"
    do_not_touch: "test existing"
    evidence_update: "evidence/.../test-output.txt"
    exit_verification: "go test pass"
  - id: W5
    action: Extend LiveMonitoringService interface + implement GetUpdateLocations service
    depends_on: W4
    owner: @fixer
    validation: "go build ./service/live_monitoring/..."
    exit_criteria: "1 method di interface + 1 file service baru + NotFound sentinel"
    blocking: ready
    must_preserve: "LENGTH rule di-resolve di repo, service tinggal branch switch"
    do_not_touch: "service method existing"
    evidence_update: "evidence/.../service-diff.txt"
    exit_verification: "go build clean"
  - id: W6
    action: Tambah handler GetUpdateLocations controller + extend LiveMonitoringController interface
    depends_on: W5
    owner: @fixer
    validation: "go build ./controller/live_monitoring/..."
    exit_criteria: "1 file controller baru + 1 method di interface"
    blocking: ready
    must_preserve: "JWT 401, inline gin.H response"
    do_not_touch: "controller existing"
    evidence_update: "evidence/.../controller-diff.txt"
    exit_verification: "go build clean"
  - id: W7
    action: Tulis test gin recorder untuk GetUpdateLocations (T4)
    depends_on: W6
    owner: @fixer
    validation: "go test ./controller/live_monitoring/... -run TestGetUpdateLocations"
    exit_criteria: "4 test function pass"
    blocking: ready
    must_preserve: "controllerServiceStub pattern"
    do_not_touch: "test existing"
    evidence_update: "evidence/.../test-output.txt"
    exit_verification: "go test pass"
  - id: W8
    action: Wire route di pjp/router/live_monitoring.go
    depends_on: W7
    owner: @fixer
    validation: "go build ./... && grep -n 'update-locations' pjp/router/live_monitoring.go"
    exit_criteria: "1 line added ke router"
    blocking: ready
    must_preserve: "group router existing"
    do_not_touch: "pjp/router/router.go, pjp/main.go"
    evidence_update: "evidence/.../route-diff.txt"
    exit_verification: "go build clean + grep match"
  - id: W9
    action: Final gofmt + vet + test + build di pjp
    depends_on: W8
    owner: @fixer
    validation: "gofmt -l . | wc -l | grep ^0$ ; go vet ./... ; go test ./... -count=1 ; go build ./..."
    exit_criteria: "zero diff, zero fail, clean build"
    blocking: ready
    must_preserve: "all test existing pass"
    do_not_touch: "all"
    evidence_update: "evidence/.../final-validation.txt"
    exit_verification: "all 4 commands exit 0"
  - id: W10
    action: Tulis index.json + discovery.md di evidence/
    depends_on: W9
    owner: @fixer
    validation: "jq . .opencode/evidence/20260713-0930-sx-2364-update-locations/index.json | head -10"
    exit_criteria: "index.json valid, discovery.md referensi docs + repo"
    blocking: ready
    must_preserve: "evidence naming convention"
    do_not_touch: "existing evidence"
    evidence_update: "evidence/.../index.json + discovery.md"
    exit_verification: "jq valid + file readable"
```

### Execution ownership table

| Subsystem | Implementation owner | Review gate |
|---|---|---|
| DTO (request/response) | @fixer | @quality-gate |
| Repository (interface + impl) | @fixer | @quality-gate |
| Service (interface + impl) | @fixer | @quality-gate |
| Controller (interface + impl) | @fixer | @quality-gate |
| Router | @fixer | @quality-gate |
| Test (sqlmock + gin) | @fixer | @quality-gate |
| Evidence | @fixer | @quality-gate |
| Final signoff | n/a | @quality-gate |

## Progress Tracking

- tracker_path: `.opencode/state/20260713-0930-sx-2364-update-locations/progress.json`
- init_command: `python3 ~/.config/opencode/scripts/task-progress.py 20260713-0930-sx-2364-update-locations --init --plan .opencode/plans/20260713-0930-sx-2364-update-locations.md`
- summary_command: `python3 ~/.config/opencode/scripts/task-progress.py 20260713-0930-sx-2364-update-locations --summary`
- checklist_command: `python3 ~/.config/opencode/scripts/task-progress.py 20260713-0930-sx-2364-update-locations --checklist`
- update_rules:
  - W1: status `pending` → `in_progress` saat mulai; `in_progress` → `completed` setelah `gofmt -l ./data/` clean.
  - W2: sama dengan W1 + `go build` clean.
  - W3: `in_progress` → `completed` setelah `go test -run TestGetEmployeeRole` pass.
  - W4: `in_progress` → `completed` setelah `go test -run 'TestGet(EmployeeRole|UpdateLocations)'` pass.
  - W5: `in_progress` → `completed` setelah `go build` clean.
  - W6: sama.
  - W7: `in_progress` → `completed` setelah `go test -run TestGetUpdateLocations` pass.
  - W8: `in_progress` → `completed` setelah `go build` + `grep` match.
  - W9: `in_progress` → `completed` setelah 4 command exit 0.
  - W10: `in_progress` → `completed` setelah `jq . index.json` valid.
  - Update required setiap status transition. Update juga setiap kali evidence file baru ditulis.
  - Update wajib di cross-lane handoff (orchestrator → fixer → quality-gate).
- task_map:
  - W1: `python3 ~/.config/opencode/scripts/task-progress.py 20260713-0930-sx-2364-update-locations --update W1 --status <status> --owner @fixer --evidence .opencode/evidence/20260713-0930-sx-2364-update-locations/dto-diff.txt`
  - W2: `--update W2 --evidence .opencode/evidence/20260713-0930-sx-2364-update-locations/repo-interface-diff.txt`
  - W3: `--update W3 --evidence .opencode/evidence/20260713-0930-sx-2364-update-locations/sql-query.txt`
  - W4: `--update W4 --evidence .opencode/evidence/20260713-0930-sx-2364-update-locations/test-output.txt`
  - W5: `--update W5 --evidence .opencode/evidence/20260713-0930-sx-2364-update-locations/service-diff.txt`
  - W6: `--update W6 --evidence .opencode/evidence/20260713-0930-sx-2364-update-locations/controller-diff.txt`
  - W7: `--update W7 --evidence .opencode/evidence/20260713-0930-sx-2364-update-locations/test-output.txt`
  - W8: `--update W8 --evidence .opencode/evidence/20260713-0930-sx-2364-update-locations/route-diff.txt`
  - W9: `--update W9 --evidence .opencode/evidence/20260713-0930-sx-2364-update-locations/final-validation.txt`
  - W10: `--update W10 --evidence .opencode/evidence/20260713-0930-sx-2364-update-locations/index.json`

Tracker updates at every status transition are mandatory, not optional bookkeeping.

## Validation Commands

1. `cd pjp && go mod download && go mod tidy`
2. `cd pjp && gofmt -l .` (expect empty output)
3. `cd pjp && go vet ./...`
4. `cd pjp && go build ./...`
5. `cd pjp && go test ./controller/live_monitoring/... ./service/live_monitoring/... ./repository/live_monitoring/... -count=1 -v`
6. `cd pjp && go test ./... -count=1`
7. `grep -n 'update-locations' pjp/router/live_monitoring.go` (expect 1 match)
8. `python3 ~/.config/opencode/scripts/subagent-handoff-check.py --plan .opencode/plans/20260713-0930-sx-2364-update-locations.md`
9. `python3 ~/.config/opencode/scripts/plan-compliance-check.py --project-root . --plan .opencode/plans/20260713-0930-sx-2364-update-locations.md --task-id 20260713-0930-sx-2364-update-locations`
10. `python3 ~/.config/opencode/scripts/plan-execution-readiness.py .opencode/plans/20260713-0930-sx-2364-update-locations.md --project-root .`
11. `python3 ~/.config/opencode/scripts/validate-plan-depth.py .opencode/plans/20260713-0930-sx-2364-update-locations.md`
12. `python3 ~/.config/opencode/scripts/task-progress.py 20260713-0930-sx-2364-update-locations --summary`

## Evidence Requirements

- `.opencode/evidence/20260713-0930-sx-2364-update-locations/index.json`: manifest daftar file bukti.
- `.opencode/evidence/20260713-0930-sx-2364-update-locations/discovery.md`: ringkasan bukti repo (paths yang dibaca) + referensi dokumen Point 8.
- `.opencode/evidence/20260713-0930-sx-2364-update-locations/route-diff.txt`: git diff `pjp/router/live_monitoring.go`.
- `.opencode/evidence/20260713-0930-sx-2364-update-locations/sql-query.txt`: full SQL union yang dipakai (3 cabang + count).
- `.opencode/evidence/20260713-0930-sx-2364-update-locations/test-output.txt`: output `go test -v`.
- `.opencode/evidence/20260713-0930-sx-2364-update-locations/curl-trace.txt`: mock curl dengan token redacted, berisi request + expected response shape.
- `.opencode/evidence/20260713-0930-sx-2364-update-locations/final-validation.txt`: output dari validation commands 2-4.
- `.opencode/evidence/20260713-0930-sx-2364-update-locations/dto-diff.txt`, `repo-interface-diff.txt`, `service-diff.txt`, `controller-diff.txt`: file diff kecil (boleh digabung dengan route-diff jika diff kecil).

## Done Criteria

- C1: 12 file changes seperti Expected Files.
- C2: Semua 12 AC pass (validator + integration jika ada DB).
- C3: 9 dari 10 evidence file ada (route-diff, sql-query, test-output, curl-trace, index, discovery, final-validation, dto-diff, controller-diff) sebelum final summary. Sisanya (repo-interface, service) boleh digabung.
- C4: `subagent-handoff-check.py` exit 0.
- C5: `plan-compliance-check.py` exit 0.
- C6: `plan-execution-readiness.py` exit 0.
- C7: `validate-plan-depth.py` exit 0 (PASS_FOR_SLICE acceptable).
- C8: `task-progress.py --summary` report semua task W1-W10 completed.
- C9: `gofmt -l` empty; `go vet` clean; `go test` zero fail; `go build` clean.

## Final Planning Summary

- Artifacts consulted: `docs/Monitoring_Activity_BE.docx` (Point 8), `pjp/router/live_monitoring.go`, `pjp/router/router.go`, `pjp/service/live_monitoring/live_monitoring_service.go`, `pjp/service/live_monitoring/get_distributor_service.go`, `pjp/repository/live_monitoring/live_monitoring_repository.go`, `pjp/repository/live_monitoring/get_distributor_repository.go`, `pjp/repository/live_monitoring/get_detail_repository.go`, `pjp/repository/live_monitoring/get_detail_repository_test.go`, `pjp/repository/pjp/get_destination_details_repository.go`, `pjp/data/request/live_monitoring_request.go`, `pjp/data/response/live_monitoring_response.go`, `pjp/controller/live_monitoring/live_monitoring_controller.go`, `pjp/controller/live_monitoring/get_distributor_controller.go`, `pjp/controller/live_monitoring/get_distributor_controller_test.go`, `pjp/middleware/jwt.go`, `pjp/exception/error_handler.go`, `pjp/model/outlet_visit_list.go`, `pjp/model/outlet_visit_list_principle.go`, `mobile/model/user_location.go`, `pjp/go.mod`, sibling plan `.opencode/plans/20260630-1441-sx-2361-latest-location.md`.
- Artifacts created: this plan, 1 discovery file (placeholder), 1 index.json (placeholder).
- Key decisions:
  1. Resolver principal vs distributor via `LENGTH(mst.m_employee.cust_id) > 6` di SQL, single source of truth.
  2. Tanpa pagination (revisi user 2026-07-13). Endpoint mengembalikan seluruh timeline satu hari per request.
  3. Urutan `recorded_at ASC` (override deskripsi `updated_at desc` di dokumen).
  4. Validasi `400 BAD_REQUEST` (override dokumen awal yang sebut 422).
- Assumptions (label `repo_local_evidence` kecuali disebut): D5 date default hari ini, D6 404 generik, D8 distributor `destination_type` null, D9 destination_name join ke mst.m_outlet, D10 `accuracy` tidak di-serve.
- Open questions: tidak ada (semua dijawab lewat question tool 4 pertanyaan). Asumsi di atas adalah reversible; bisa di-promote ke question jika executor ragu saat implementasi.
- Readiness: PASS_FOR_SLICE (scope dibatasi 1 endpoint; sibling `DI-2` bug fix di SX-2361 tetap terpisah).
- Cleanup: tidak ada draft/evidence yang perlu dihapus dari phase ini. Folder `draft/20260713-0930-sx-2364-update-locations/` berisi catatan pendek (akan ditambahkan setelah plan ini). Folder `state/20260713-0930-sx-2364-update-locations/` akan di-init setelah plan ini final.

### Active-lane reset note

Eksekusi harus dilakukan di bawah lane berikutnya (`@orchestrator` atau `@fixer`). Pembatasan read-only `@artifact-planner` tidak berlaku di lane eksekusi. Worker harus refresh konteks aktif (hak baca/tulis, mode Go service) sebelum implementasi.
