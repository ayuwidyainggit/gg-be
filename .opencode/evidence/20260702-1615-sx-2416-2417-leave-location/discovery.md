# Discovery — SX-2416 & SX-2417 (Leave Location in Live Monitoring)

Task: tambah `leave_longitude` dan `leave_latitude` di response `GET /scylla-pjp/v1/live-monitoring-distributor` dan `GET /scylla-pjp/v1/live-monitoring-principal`.

Mode: maintenance-stability (shortest diff, tanpa refactor besar).

## 1. Confirmed vs Assumed Audit

| # | Claim | Source | Level |
|---|---|---|---|
| 1 | Endpoint distributor ada di service `pjp` pada path `GET /api/v1/live-monitoring-distributor` | `pjp/router/live_monitoring.go:13` | confirmed_repo |
| 2 | Endpoint principal ada di service `pjp` pada path `GET /api/v1/live-monitoring-principal` | `pjp/router/live_monitoring.go:12` | confirmed_repo |
| 3 | Controller distributor hanya pass-through ke service | `pjp/controller/live_monitoring/get_distributor_controller.go:63` | confirmed_repo |
| 4 | Controller principal hanya pass-through ke service | `pjp/controller/live_monitoring/get_principal_controller.go:63` | confirmed_repo |
| 5 | Service distributor mengambil row utama dari `GetDistributorMonitoring`, lalu enrich arrival dari `GetDistributorLatestVisitCoordinates` | `pjp/service/live_monitoring/get_distributor_service.go:88-127` | confirmed_repo |
| 6 | Query distributor utama sudah `LEFT JOIN pjp.outlet_visit_list ovl` | `pjp/repository/live_monitoring/get_distributor_repository.go:141-143` | confirmed_repo |
| 7 | Query distributor utama sudah select `ovl.arrive_at`, `ovl.leave_at`, `ovl.start`, `ovl.finish`, `ovl.skip_at`, `ovl.skip_reason` | `pjp/repository/live_monitoring/get_distributor_repository.go:131-139` | confirmed_repo |
| 8 | `arrive_longitude/arrive_latitude` distributor bukan dari `ovl`, tetapi dari `mobile.visits` via enrichment | `pjp/repository/live_monitoring/get_distributor_repository.go:181-223` + `pjp/service/live_monitoring/get_distributor_service.go:150-160` | confirmed_repo |
| 9 | Query principal regular sudah `LEFT JOIN pjp_principles.outlet_visit_list ovl` | `pjp/repository/live_monitoring/get_principal_repository.go:132` | confirmed_repo |
| 10 | Query principal extra call juga sudah `LEFT JOIN pjp_principles.outlet_visit_list ovl` | `pjp/repository/live_monitoring/get_principal_extra_call_repository.go:75` | confirmed_repo |
| 11 | Query principal regular sudah select `ovl.arrive_at`, `ovl.leave_at`, `ovl.longitude`, `ovl.latitude` | `pjp/repository/live_monitoring/get_principal_repository.go:106-123` | confirmed_repo |
| 12 | Query principal extra call juga sudah select `ovl.arrive_at`, `ovl.leave_at`, `ovl.longitude`, `ovl.latitude` | `pjp/repository/live_monitoring/get_principal_extra_call_repository.go:50-67` | confirmed_repo |
| 13 | Kolom DB `pjp.outlet_visit_list.leave_longitude` dan `leave_latitude` ada | `pjp/model/outlet_visit_list.go:46-47` | confirmed_repo |
| 14 | Kolom DB `pjp_principles.outlet_visit_list.leave_longitude` dan `leave_latitude` ada | `pjp/model/outlet_visit_list_principle.go:46-47` | confirmed_repo |
| 15 | Request DTO leave visit sudah memakai `*string` untuk `leave_longitude` dan `leave_latitude` | `pjp/data/request/visit_request.go:69-77` | confirmed_repo |
| 16 | Response DTO live monitoring saat ini belum punya `leave_longitude` / `leave_latitude` | `pjp/data/response/live_monitoring_response.go:50-68` | confirmed_repo |
| 17 | Row model live monitoring saat ini belum punya `leave_longitude` / `leave_latitude` | `pjp/model/live_monitoring.go:4-82` | confirmed_repo |
| 18 | Mapping distributor ke `destination_data[]` dilakukan di `transformDistributorRows` | `pjp/service/live_monitoring/get_distributor_service.go:417-433` | confirmed_repo |
| 19 | Mapping principal ke `destination_data[]` dilakukan di `transformPrincipalRows` | `pjp/service/live_monitoring/get_principal_service.go:237-254` | confirmed_repo |
| 20 | Referensi prompt user `monitoring_activity_be_doc.txt` / `SX-2361...` / `SX-2034...` / `SX-2038...` tidak ditemukan di repo lokal | `glob` pada root repo menghasilkan no file | confirmed_repo |
| 21 | FE nullability untuk field baru disetujui `null` | jawaban user via `question` tool | user_confirmed |
| 22 | FE contract Google Docs tidak diverifikasi dari repo lokal | URL ada di prompt user, tidak diakses saat discovery | unverified |

## 2. Ringkasan Teknis

Perubahan ini lebih kecil dari patch SX-2421/SX-2422 (`file_url`).

- **Distributor**: tidak perlu ubah jalur `mobile.visits` untuk arrival. Cukup tambahkan select:
  - `NULLIF(ovl.leave_longitude, '') AS leave_longitude`
  - `NULLIF(ovl.leave_latitude, '') AS leave_latitude`
- **Principal regular**: cukup tambahkan dua kolom yang sama dari `pjp_principles.outlet_visit_list ovl`.
- **Principal extra call**: juga perlu tambahkan dua kolom yang sama dari `ovl`, karena endpoint principal menggabungkan regular route + extra call.
- Karena kolom DB bertipe varchar dan user memilih `null` bila kosong, tipe response yang paling aman adalah `*string` dengan `json:"leave_longitude"` dan `json:"leave_latitude"` tanpa `omitempty`.

## 3. Path Map

### Distributor
1. `pjp/router/live_monitoring.go`
2. `pjp/controller/live_monitoring/get_distributor_controller.go`
3. `pjp/service/live_monitoring/get_distributor_service.go`
4. `pjp/repository/live_monitoring/get_distributor_repository.go`
5. `pjp/model/live_monitoring.go`
6. `pjp/data/response/live_monitoring_response.go`

### Principal
1. `pjp/router/live_monitoring.go`
2. `pjp/controller/live_monitoring/get_principal_controller.go`
3. `pjp/service/live_monitoring/get_principal_service.go`
4. `pjp/repository/live_monitoring/get_principal_repository.go`
5. `pjp/repository/live_monitoring/get_principal_extra_call_repository.go`
6. `pjp/model/live_monitoring.go`
7. `pjp/data/response/live_monitoring_response.go`

## 4. Existing Reuse Candidates

- Reuse pattern `*string` nullable dari `FileURL *string` pada live monitoring response/model:
  - `pjp/model/live_monitoring.go:27,67,81`
  - `pjp/data/response/live_monitoring_response.go:63`
- Reuse pola test dari patch sebelumnya:
  - `pjp/service/live_monitoring/get_distributor_service_test.go` sudah punya assertion `file_url`
  - `pjp/service/live_monitoring/get_principal_service_test.go` sudah punya assertion `file_url`
- Reuse diff boundary dan validator dari plan SX-2421/SX-2422.

## 5. Constraint Repo / Harness

- Tidak ada `PROJECT_STACK.md`, `PROJECT_COMMANDS.md`, `FRAMEWORK_PLAYBOOK.md`, atau `PROJECT_DETECTED_TOOLS.md` di `.opencode/docs/`; validasi mengikuti `QUALITY.md` + `AGENTS.md` repo.
- Service target adalah `pjp` (Gin). Jangan menyentuh `pjp-principle/`; endpoint live monitoring ada di service `pjp`.
- Validasi standar per repo:
  - `rtk go mod download && rtk go mod tidy`
  - `rtk go test ./...`
- Untuk task ini, validasi cukup difokuskan ke package live monitoring + build service `pjp`, lalu optional curl manual tanpa menyimpan token.

## 6. Open Questions

Tidak ada blocker tersisa.

- Nullability field baru sudah diputuskan user: `null`.
- Tipe field mengikuti schema DB + request DTO existing: `*string`.
- Google Docs FE contract belum diverifikasi, tetapi prompt user + schema DB sudah cukup kuat untuk slice ini.
