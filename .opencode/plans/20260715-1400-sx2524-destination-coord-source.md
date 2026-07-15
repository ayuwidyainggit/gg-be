# SX-2524 — destination_data lon/lat source: `mst.m_outlet`, not `route_outlet_history`

plan_status: PASS_FOR_SLICE
preflight_disposition: `target-app`


## Goal
`GET /scylla-pjp/api/v1/live-monitoring-distributor` mengembalikan `destination_data[].longitude` dan `.latitude` dari master outlet (`mst.m_outlet.longitude/latitude`) — bukan dari snapshot baris `pjp.route_outlet_history` per-hari. Berlaku untuk `pjp_data` maupun `extra_call_data`. Tidak ada regresi di `arrive_*`, `leave_*`, `start`, `finish`, `skip_at`.

## Non-goals
- Tidak menyentuh endpoint `/scylla-pjp/v1/live-monitoring-principal` (skema `pjp_principles`, berbeda).
- Tidak mengubah FE marker rendering; hanya kontrak response BE.
- Tidak migrasi schema; tidak menulis ulang tabel; tidak menyentuh data historis.
- Tidak mengubah token/JWT, tenancy, atau response envelope.

## Scope
- Modul `pjp` (Gin, GORM, port 9010).
- Satu file SQL/query: `pjp/repository/live_monitoring/get_distributor_repository.go` (fungsi `GetDistributorMonitoring`).
- Optional: tambah test regresi.
- Verifikasi runtime: cURL ke Staging (env `1.0.0 (153)`) sesuai tiket.

## Requirements
1. `destination_data[].longitude` = `mst.m_outlet.longitude` untuk outlet `destination_id`.
2. `destination_data[].latitude` = `mst.m_outlet.latitude` untuk outlet `destination_id`.
3. Berlaku untuk blok `pjp_data` dan `extra_call_data` (satu query, satu transform, satu DTO).
4. `arrive_longitude/latitude` tetap diisi dari `mobile.visits` lewat `GetDistributorLatestVisitCoordinates`; `0` kalau belum ada visit.
5. `leave_longitude/latitude` tetap dari `pjp.outlet_visit_list` (alias `ovl`).
6. `start`, `finish`, `skip_at` tidak berubah.
7. SQL harus `LEFT JOIN` ke `mst.m_outlet`; `COALESCE(CAST(... AS FLOAT), 0)` tetap dipakai untuk handle varchar NULL/kosong.
8. Jika outlet tidak ditemukan di `mst.m_outlet` (data drift), fallback ke `0` — bukan ke `roh` (mengikuti pola COALESCE; tidak menambah perilaku baru).

## Acceptance Criteria
- [ ] Kompilasi: `cd pjp && rtk go build ./...` lulus tanpa warning baru.
- [ ] Unit: `cd pjp && rtk go test ./... -run Distributor -v` lulus; minimal satu test baru menegaskan transform `destination_data.longitude/latitude` = nilai input row (tidak ditimpa/service tidak source ulang dari `roh` setelah enrichment).
- [ ] SQL diff: di `pjp/repository/live_monitoring/get_distributor_repository.go` line 129-130, `roh.longitude`/`roh.latitude` diganti `mst.m_outlet.longitude`/`mst.m_outlet.latitude`, dan `LEFT JOIN mst.m_outlet mo ON mo.outlet_id = roh.outlet_id` ditambahkan (atau gunakan alias existing bila sudah ada).
- [ ] Runtime Staging: cURL `GET /scylla-pjp/api/v1/live-monitoring-distributor?date=1784030400&status=Approved&status=Need+Review&emp_id=466&distributor_id=102` mengembalikan `destination_data[].longitude=106.82096445685606` dan `.latitude=-6.248543092691333` untuk `Fahmi Garage` (atau nilai master `mst.m_outlet` saat itu), bukan lagi nilai history. Ticket menyebut posisi sebagai `(latitude, longitude)`; contract JSON harus tetap `(longitude, latitude)`.
- [ ] Tidak ada regresi: `arrive_*`, `leave_*`, `start`, `finish`, `skip_at` sama dengan sebelum fix (cek field-shape manual via JSON diff).
- [ ] Tidak ada perubahan response shape atau field lain.

## Existing Patterns/Reuse
- Pola `COALESCE(CAST(... AS FLOAT), 0)` sudah dipakai di query `GetDistributorMonitoring` (line 129-130 untuk `roh`, line 133-134 untuk `arrive_*`) — reuse pola yang sama untuk switch ke `mst.m_outlet`. Tidak bikin helper SQL baru.
- Pola `NULLIF(<column>, '')` + `CAST AS FLOAT` + `COALESCE(..., 0)` sudah dipakai di `GetDistributorLatestVisitCoordinates` (line 169) dan di kolom `leave_*` (line 135-136) untuk handle `varchar` kosong. Pakai pola yang sama untuk `mo.longitude`/`mo.latitude` karena `mst.m_outlet` bertipe `varchar(125)`.
- `distributorRepoStub` di `pjp/service/live_monitoring/get_distributor_service_test.go:15-188` sudah implement `LiveMonitoringRepository`. Test baru cukup menambahkan field assertion; tidak perlu instantiate stub baru atau library mocking tambahan.
- Pola test service `TestTransformDistributorRows_AssignsLeaveLocation` (line 537) dan `TestTransformDistributorRows_NilLeaveLocation` (line 570) adalah template yang harus diikuti untuk test `TestTransformDistributorRows_UsesRowLongitudeLatitudeAsDestinationCoordinates`. Pola yang sama: stub inject canned `[]model.LiveMonitoringDistributorRow`, panggil `transformDistributorRows`, assert terhadap field destination_data.
- Pola controller test `TestGetDistributorMonitoring_ReturnsMonitoringPayloadWithoutStaleData` (`pjp/controller/live_monitoring/get_distributor_controller_test.go:44`) tetap dipakai untuk smoke HTTP jika regression test ditambah di layer controller.
- Service `transformDistributorRows` line 398 sudah split `pjp_data` vs `extra_call_data` lewat `IsExtraCall`; tidak perlu sentuh untuk fix ini. Patch tidak boleh menambah percabangan transform baru.
- Tidak butuh migrasi / bukan domain generator; tidak ada generator path untuk fix ini, hanya edit source GORM + test Go. Tidak perlu library baru.
- `rtk go test` / `rtk go build` di `pjp/` adalah command resmi per `.opencode/docs/PROJECT_COMMANDS.md`; tidak perlu tooling tambahan.

## Source Anatomy
| Subsystem | File | Lines | Peran |
|---|---|---|---|
| Router | `pjp/router/live_monitoring.go` | 13 | register `GET /live-monitoring-distributor` |
| Controller | `pjp/controller/live_monitoring/get_distributor_controller.go` | 33 | bind query, panggil service, return response |
| Service | `pjp/service/live_monitoring/get_distributor_service.go` | 24 | orchestrate query + enrichment + transform |
| Repository (main) | `pjp/repository/live_monitoring/get_distributor_repository.go` | 107 | query SQL yang akan di-patch (line 129-130) |
| Repository (visits) | `pjp/repository/live_monitoring/get_distributor_repository.go` | 169 | `GetDistributorLatestVisitCoordinates` — di luar scope |
| Model | `pjp/model/live_monitoring.go` | 38-77 | `LiveMonitoringDistributorRow.Longitude/Latitude` |
| Response DTO | `pjp/data/response/live_monitoring_response.go` | 51-70 | `LiveMonitoringDestinationData.Longitude/Latitude` |
| Test stub | `pjp/service/live_monitoring/get_distributor_service_test.go` | 15-188 | `distributorRepoStub` |
| Tests existing | file yang sama | 190-598 | pola `TestGetDistributorMonitoring_*`, `TestTransformDistributorRows_*` |

Query utama mengambil `roh` sebagai base row untuk route harian; `ovl` hanya menyumbang koordinat pulang. `LiveMonitoringDistributorRow` membawa dua koordinat destination tunggal menuju transform. Service memperkaya data visit dan metadata tanpa mengubah koordinat destination. `transformDistributorRows` kemudian membangun route normal atau extra call berdasarkan `IsExtraCall`; keduanya memakai field `Longitude`/`Latitude` dari row yang sama. Karena itu patch tepat ada di SQL projection dan join master outlet, bukan controller/service/DTO. Repository dan test source sudah diinspeksi; exact join alias harus dicek ulang sebelum patch untuk mencegah alias duplikat.

## Reference Map
- **Ticket-focused comment (SX-2524#17118)**: source of truth untuk arah fix — `repo-backed` (`user_confirmed`).
- **Doc "Monitoring Activity - BE" (GDoc `19NxHFOBpqT5XO2X6RmLyBC8JRANkbHnwoGzRg6ldK-Y`)**: menjelaskan intended source = `mst.m_outlet` — `docs-backed` (`user_confirmed`).
- **SQL source di `pjp/repository/live_monitoring/get_distributor_repository.go:129-130`**: konfirmasi bug di level kode — `confirmed_repo`.
- **cURL blok di tiket**: `runtime-backed` smoke test pasca-fix.
- Skip reference UI/templates (bukan UI work).

## Constraints
- `pjp` adalah Gin, bukan Fiber — jangan pakai `fiber.Ctx` / Fiber middleware.
- GORM dipakai di module ini — pakai pattern `db.Raw(...).Scan(...)` atau builder yang sudah ada; tidak migrasi ke sqlx.
- Tenant/scope rules per `.opencode/docs/ARCHITECTURE.md` berlaku; query tidak boleh bocor antar distributor.
- Tidak boleh menambah library baru.
- Tidak boleh menulis/menyalin `.env` atau JWT ke file plan/evidence.

## Risks
- `mst.m_outlet` kemungkinan sudah ada di join lain (mis. untuk `arrive_longitude` versi lama atau `extra_call_data` di query yang sama). Cek dulu apakah join duplikat dibutuhkan.
- Kalau `roh.outlet_id` orphan (tidak ada di `mst.m_outlet`), output jadi `0` — lebih baik dari output `roh` yang salah, tapi harus dikonfirmasi FE/QA acceptable.
- Cast `varchar(125)` ke `FLOAT` bisa gagal di value non-numeric; `COALESCE` di `FLOAT` mungkin tidak menangkap error cast → cek `TryCast`/`NULLIF` pattern atau fallback eksplisit. Pola existing `NULLIF(..., '')` dari `arrive_*`/`leave_*` lebih aman; copy pola itu.
- Fix terbalik: hanya patch distributor; jangan touch principal — mereka query beda, beda tabel referensi.

## Decisions/Assumptions
- **Decision**: pakai `LEFT JOIN mst.m_outlet mo ON mo.outlet_id = roh.outlet_id` (bukan `JOIN`) — agar row tetap muncul walau master outlet hilang, dengan `longitude/latitude` jadi `0` lewat `COALESCE`. Alasan: FE tidak akan crash kalau ada satu outlet tanpa master; dan mencegah kehilangan data visit yang masih punya `arrive_*`/`leave_*` valid.
- **Decision**: pakai `NULLIF(mo.longitude, '')` lalu `CAST(... AS FLOAT)` lalu `COALESCE(..., 0)` — match pola `arrive_*`/`leave_*` existing, lebih aman dari `CAST` polos varchar.
- **Decision**: tidak pisah `pjp_data` vs `extra_call_data` jadi 2 query — sama seperti sekarang, satu query sudah cukup karena split terjadi di service.
- **Assumption**: `mst.m_outlet.outlet_id` = `roh.outlet_id` adalah kunci join yang benar (sudah dipakai di query doc/GDoc dengan `roh.outlet_id = mo.outlet_id`).
- **Assumption**: FE consumer tidak cache value ini lintas hari — perubahan source aman dilakukan tanpa versioning.
- **Assumption**: tidak ada FE logic yang bergantung pada `destination_data.longitude` = nilai history (lihat dengan FE saat re-test, bukan asumsi planner).

## Execution Source of Truth
Urutan preseden untuk implementasi:
1. Aman: PG schema + tenant read-only, `.opencode/docs/ARCHITECTURE.md`.
2. Kontrak: ticket SX-2524 + focused comment #17118 + GDoc `19NxHFOBpqT5XO2X6RmLyBC8JRANkbHnwoGzRg6ldK-Y`.
3. Invariants: bagian "Non-negotiable Implementation Invariants" di bawah.
4. Worklist: bagian "Execution-ready Worklist / Handoff Contract".
5. Acceptance + Done Criteria di bawah.
6. Implementation Steps.
7. Follow-up.

## Non-negotiable Implementation Invariants
- HANYA patch `pjp/repository/live_monitoring/get_distributor_repository.go` di `GetDistributorMonitoring`. Method lain (`GetDistributorLatestVisitCoordinates`, dst.) tidak boleh berubah.
- `arrive_longitude/latitude` masih `0` di SQL utama dan di-enrich dari `mobile.visits`. Jangan isi dari `mst.m_outlet`.
- `leave_longitude/leave_latitude` masih dari `ovl`. Jangan isi dari `mst.m_outlet`.
- `pjp_data` dan `extra_call_data` kembali dari satu query, satu transform, satu DTO.
- Tidak mengubah response shape, field name, atau urutan field.
- Tidak menyentuh `live-monitoring-principal` atau service/service-module lain.
- Tidak menulis/menyalin `.env`, JWT, atau secret ke file plan/evidence/commit.
- Tiap perubahan source file harus dibarengi `task-progress.py --update` (lihat Progress Tracking).
- Worker (fixer) hanya edit source + test di modul `pjp/`. Tidak edit file plan/evidence/selainnya.

## Do Not / Reject If
- Reject: ganti `destination_data.longitude/latitude` ke `arrive_longitude/latitude` (sumber berbeda, semantics berbeda).
- Reject: hard-code lat/lng outlet spesifik di query.
- Reject: tambah kolom baru ke DTO; field harus tetap `longitude`/`latitude`.
- Reject: migrasi schema (column rename, new column, copy) untuk fix ini.
- Reject: ubah endpoint `live-monitoring-principal`.
- Reject: "fix" dengan menambahkan feature flag/konfigurasi runtime — terlalu besar untuk satu defect.
- Reject: skip sqlmock/regression test atau cURL smoke, lalu klaim selesai.
- Reject: commit `.env` atau secret.

## Diff Boundary
- Allowed:
  - `pjp/repository/live_monitoring/get_distributor_repository.go` (patch line 129-130 + tambah join jika belum ada).
  - `pjp/repository/live_monitoring/get_distributor_repository_test.go` jika sudah ada (cek dulu); kalau belum, tambah test di `pjp/service/live_monitoring/get_distributor_service_test.go` atau controller test.
  - `.opencode/evidence/20260715-1400-sx2524-destination-coord-source/` (test logs, curl output, sql diff).
- Forbidden:
  - `pjp/router/**`
  - `pjp/controller/live_monitoring/**` (kecuali tambah test)
  - `pjp/service/live_monitoring/**` source (kecuali tambah test)
  - `pjp/model/live_monitoring.go`
  - `pjp/data/response/live_monitoring_response.go`
  - Semua service lain (`pjp-principle`, `master`, `mobile`, dll.)
  - `.env`, `*.key`, `*.pem`, compose secrets.
- Out-of-scope change harus di-revert atau justifikasi di evidence sebelum `@quality-gate`.

## Grounding Contract
- `repo-backed` claim: `destination_data.longitude/latitude` saat ini dari `roh.longitude/roh.latitude` (line 129-130 `pjp/repository/live_monitoring/get_distributor_repository.go`).
- `repo-backed` claim: `arrive_*` di-enrich dari `mobile.visits` (line 169), `leave_*` dari `ovl` (line 135-136).
- `repo-backed` claim: `pjp_data` vs `extra_call_data` split lewat `is_extra_call` di `transformDistributorRows` (line 398).
- `user_confirmed` claim: ticket SX-2524 + focused comment #17118 arah fix; Principal endpoint di luar scope; Staging creds.
- `docs-backed` claim: `mst.m_outlet.longitude/latitude` adalah `varchar(125)` (GDoc `19NxHFOBpqT5XO2X6RmLyBC8JRANkbHnwoGzRg6ldK-Y`).
- `runtime-backed` claim: cURL Staging + JWT dipakai sebagai smoke test pasca-fix.
- Lihat `Source Anatomy` (repo files) + `Reference Map` (klasifikasi per claim) + `Evidence Requirements` (paths). Plan ini tidak bergantung pada klaim tanpa bukti; klaim yang belum konfirmasi terdaftar di `Decisions / Assumptions`.

## TDD / Test Plan
- **TDD required**: iya. Defect memiliki field-level assertion yang stabil.
- **Existing test pattern**:
  - Service: `distributorRepoStub` inject canned `[]model.LiveMonitoringDistributorRow` → service transform → assert response.
  - Controller: `controllerServiceStub` inject canned response → HTTP call → assert JSON.
- **First failing test (Red)**:
  - Test service-level: `TestTransformDistributorRows_UsesRowLongitudeLatitudeAsDestinationCoordinates` —
    inject 2 row (1 `IsExtraCall=false`, 1 `IsExtraCall=true`) dengan `Longitude=106.82096445685606, Latitude=-6.248543092691333` dan `DestinationID` valid.
    Panggil `transformDistributorRows`.
    Assert: di `pjp_data[].route_data[].destination_data[]` ada entry dengan `longitude == 106.82096445685606 && latitude == -6.248543092691333`.
    Assert: di `extra_call_data[].route_data[].destination_data[]` ada entry identik.
  - Catatan: ini test kontrak, bukan test SQL. SQL level akan divalidasi lewat cURL runtime.
- **Green**: patch query (lihat Implementation Steps). Test harusnya tetap pass karena model `Longitude/Latitude` di-read langsung dari row; perubahan di repository akan menghasilkan row dengan nilai baru.
- **Refactor**: tidak ada duplikasi yang perlu di-refactor; satu baris SQL.
- **Edge cases**:
  - `roh.outlet_id` orphan → `mo.longitude = NULL` → `COALESCE(..., 0)` → `destination_data.longitude = 0` (verifikasi tidak panic).
  - `mo.longitude` berisi string non-numeric → `CAST AS FLOAT` error → wrap dengan `NULLIF` + COALESCE; pattern lihat `arrive_longitude`.
- **Exempt**: tidak perlu test visual / browser; backend-only.
- **Run command**:
  - `cd pjp && rtk go test ./service/live_monitoring/ -run TestTransformDistributorRows -v`
  - `cd pjp && rtk go test ./controller/live_monitoring/ -run TestGetDistributor -v`
  - `cd pjp && rtk go test ./... -run Distributor -v`
  - `cd pjp && rtk go build ./...`

## Implementation Steps
1. Buka `pjp/repository/live_monitoring/get_distributor_repository.go`. Verifikasi line 107-160 (fungsi `GetDistributorMonitoring`). Identifikasi alias `roh`, `mo`, `pjp`, `ovl` dan join existing.
2. Cek apakah sudah ada `JOIN mst.m_outlet mo` (atau alias setara) untuk `arrive_*`/lokasi lain. Jika belum, tambahkan `LEFT JOIN mst.m_outlet mo ON mo.outlet_id = roh.outlet_id` ke chain join.
3. Ganti ekspresi:
   - `COALESCE(CAST(roh.longitude AS FLOAT), 0)` → `COALESCE(CAST(NULLIF(mo.longitude, '') AS FLOAT), 0)`
   - `COALESCE(CAST(roh.latitude AS FLOAT), 0)` → `COALESCE(CAST(NULLIF(mo.latitude, '') AS FLOAT), 0)`
   - Pertahankan `0` untuk `arrive_longitude/latitude` di main query.
   - Pertahankan `NULLIF(ovl.leave_longitude, '')` / `NULLIF(ovl.leave_latitude, '')`.
4. Tambah test baru di `pjp/service/live_monitoring/get_distributor_service_test.go` sesuai "First failing test" di TDD Plan. Pattern: lihat `TestTransformDistributorRows_AssignsLeaveLocation` line 537.
5. Jalankan test service + controller: harus pass.
6. Jalankan `cd pjp && rtk go build ./...` — harus lulus tanpa warning baru.
7. (Opsional tapi direkomendasikan) Build image / restart `pjp` di compose lokal:
   - `rtk docker compose -f docker-compose.yml up -d pjp`
   - Smoke lokal (jika ada data) atau langsung ke Staging.
8. Curl Staging pakai JWT dari tiket. Pipe ke `jq '.data[].pjp_data[].route_data[].destination_data[] | select(.destination_name | test("Fahmi")) | {longitude, latitude}'` — ekspektasi `-6.248543092691333` / `106.82096445685606` (atau nilai master saat itu).
9. Bandingkan JSON response `arrive_*`/`leave_*`/`start`/`finish`/`skip_at` sebelum vs sesudah dengan `diff` atau simpan sebelum/after ke `.opencode/evidence/<task-id>/response-before.json` dan `response-after.json`.
10. Catat ringkasan diff + test log + curl output ke `.opencode/evidence/20260715-1400-sx2524-destination-coord-source/result.md`.

## Expected Files to Change
- `pjp/repository/live_monitoring/get_distributor_repository.go` (patch ~2 baris + 1 join)
- `pjp/service/live_monitoring/get_distributor_service_test.go` (tambah 1 test ~40 baris)
- `.opencode/evidence/20260715-1400-sx2524-destination-coord-source/result.md` (baru)

## Agent / Tool Routing
| Area | Owner |
|---|---|
| Source edit (repository + test) | `@fixer` |
| Review signoff (BE correctness + regression risk) | `@quality-gate` |
| FE re-test di Staging | `@frontend` (assignee `dev.fe` original) |
| Runtime smoke cURL Staging | `@fixer` saat fix; atau FE saat re-test |
| Plan owner | `@artifact-planner` (selesai setelah plan ini) |

## Executor Handoff Prompt
Task: SX-2524 BE fix — destination_data lon/lat source switch.
Scope: pjp module only; patch GetDistributorMonitoring to source destination_data.{longitude,latitude} from mst.m_outlet (via LEFT JOIN) instead of pjp.route_outlet_history. Apply to pjp_data and extra_call_data (single query). Do NOT touch live-monitoring-principal, mobile, master, or any other service.

Must preserve: arrive_longitude/latitude still filled from mobile.visits; leave_longitude/latitude still from pjp.outlet_visit_list (ovl); start/finish/skip_at untouched; response field names/shape unchanged; split by is_extra_call unchanged; LEFT JOIN retains orphan rows with coord = 0.

Do not touch: pjp/router/**; pjp/controller/live_monitoring/** except tests; pjp/service/live_monitoring/** source except tests; pjp/model/live_monitoring.go; pjp/data/response/live_monitoring_response.go; other modules; .env; *.key; *.pem; JWT values.

Validate: `cd pjp && rtk go build ./...`; `cd pjp && rtk go test ./... -run Distributor -v`; Staging curl with ticket parameters; for Fahmi Garage assert JSON longitude=106.82096445685606 and latitude=-6.248543092691333, or current master values.

Exit: all validation passes; one regression test `TestTransformDistributorRows_UsesRowLongitudeLatitudeAsDestinationCoordinates` passes; Staging curl returns master coords; arrive_*/leave_*/start/finish/skip_at unchanged; evidence exists under `.opencode/evidence/20260715-1400-sx2524-destination-coord-source/`.

Claim limit: may claim BE distributor fix verified. May not claim FE marker visual verification or principal endpoint review.
Return: changed files, test log path, curl output path, before/after diff note.

## Execution-ready Worklist / Handoff Contract

```yaml
handoff:
  task_id: 20260715-1400-sx2524-destination-coord-source
  plan_id: 20260715-1400-sx2524-destination-coord-source
  caller: orchestrator
  callee: fixer
  scope: switch destination_data.{longitude,latitude} source in pjp GetDistributorMonitoring from roh to mst.m_outlet (LEFT JOIN)
  claim_level: scoped
  claim_scope: "May claim BE distributor fix merged, master coordinates returned, regression test passes. Must not claim FE visual verification, principal review, or performance benchmark."
  source_basis: ["pjp/repository/live_monitoring/get_distributor_repository.go:107-160", "pjp/model/live_monitoring.go:38-77", "pjp/data/response/live_monitoring_response.go:51-70", "pjp/service/live_monitoring/get_distributor_service_test.go:537-598", "ticket SX-2524 focused comment #17118", "GDoc 19NxHFOBpqT5XO2X6RmLyBC8JRANkbHnwoGzRg6ldK-Y"]
  must_preserve: ["arrive_longitude/latitude source mobile.visits", "leave_longitude/latitude source outlet_visit_list", "response field names and shape", "pjp_data vs extra_call_data single-query split", "LEFT JOIN not INNER"]
  do_not_touch: ["pjp/router/**", "pjp/controller/live_monitoring/** except new tests", "pjp/service/live_monitoring/** source except new tests", "pjp/model/live_monitoring.go", "pjp/data/response/live_monitoring_response.go", "pjp-principle master mobile or any other module", ".env *.key *.pem JWT secrets"]
  validation: ["cd pjp && rtk go build ./...", "cd pjp && rtk go test ./... -run Distributor -v", "curl Staging endpoint with ticket JWT and date=1784030400"]
  exit_criteria: ["all validation commands pass", "new regression test added and passes", "curl on Staging shows mst.m_outlet coords", "no regression on arrive/leave/start/finish/skip", "evidence saved"]
  evidence_required: [".opencode/evidence/20260715-1400-sx2524-destination-coord-source/result.md", ".opencode/evidence/20260715-1400-sx2524-destination-coord-source/curl-after.json", ".opencode/evidence/20260715-1400-sx2524-destination-coord-source/test.log"]
  depends_on: ["none"]
  context_bundle: ["pjp/repository/live_monitoring/get_distributor_repository.go", "pjp/model/live_monitoring.go", "pjp/data/response/live_monitoring_response.go", "pjp/service/live_monitoring/get_distributor_service_test.go", "ticket SX-2524 focused comment #17118", "GDoc 19NxHFOBpqT5XO2X6RmLyBC8JRANkbHnwoGzRg6ldK-Y"]
```

Context bundle detail (for worker, outside strict handoff schema):
- `verified_by_planner` confirmed_repo: destination_data.longitude/latitude currently from roh.longitude/roh.latitude (pjp/repository/live_monitoring/get_distributor_repository.go:129-130); arrive_* enriched from mobile.visits (line 133-134 + 169); leave_* from outlet_visit_list (line 135-136); pjp_data and extra_call_data share query, split by is_extra_call in transformDistributorRows (pjp/service/live_monitoring/get_distributor_service.go:398); route registered at pjp/router/live_monitoring.go:13
- `verified_by_planner` confirmed_docs: mst.m_outlet.longitude/latitude is varchar(125) per GDoc 19NxHFOBpqT5XO2X6RmLyBC8JRANkbHnwoGzRg6ldK-Y
- `verified_by_planner` user_confirmed: Staging creds adminbm@gmail.com/Admin123 env 1.0.0 (153); Principal endpoint out of scope
- `open_assumptions` (worker must not turn into fact): mst.m_outlet.outlet_id = roh.outlet_id join key; FE consumer does not cache destination_data across dates; FE accepts coord=0 for orphan outlets
- `source_of_truth_order` for this task: ticket focused comment #17118 -> GDoc 19NxHFOBpqT5XO2X6RmLyBC8JRANkbHnwoGzRg6ldK-Y -> repo code (current state) -> runtime cURL on Staging (post-fix evidence)

## Progress Tracking
- **tracker_path**: `.opencode/state/20260715-1400-sx2524-destination-coord-source/progress.json`
- init_command: `python3 ~/.config/opencode/scripts/task-progress.py 20260715-1400-sx2524-destination-coord-source --init --plan .opencode/plans/20260715-1400-sx2524-destination-coord-source.md`
- summary_command: `python3 ~/.config/opencode/scripts/task-progress.py 20260715-1400-sx2524-destination-coord-source --summary`
- checklist_command: `python3 ~/.config/opencode/scripts/task-progress.py 20260715-1400-sx2524-destination-coord-source --checklist`
- update_rules: tracker update wajib setiap `pending` -> `in_progress` -> `completed`/`blocked`/`cancelled`; setiap evidence ditulis -> `--evidence <path>`; setiap cross-lane handoff juga update status.
- task_map (one row per worklist id; update command listed per row):

| worklist id | owner | evidence path | update command |
|---|---|---|---|
| B1 | `@fixer` | `.opencode/evidence/20260715-1400-sx2524-destination-coord-source/result.md` | `python3 ~/.config/opencode/scripts/task-progress.py 20260715-1400-sx2524-destination-coord-source --update B1 --status in_progress --owner fixer` lalu `--status completed --evidence .opencode/evidence/20260715-1400-sx2524-destination-coord-source/result.md` |
| B2 | `@fixer` | `.opencode/evidence/20260715-1400-sx2524-destination-coord-source/test.log` | `python3 ~/.config/opencode/scripts/task-progress.py 20260715-1400-sx2524-destination-coord-source --update B2 --status in_progress --owner fixer` lalu `--status completed --evidence .opencode/evidence/20260715-1400-sx2524-destination-coord-source/test.log` |
| B3 | `@fixer` | `.opencode/evidence/20260715-1400-sx2524-destination-coord-source/curl-after.json` | `python3 ~/.config/opencode/scripts/task-progress.py 20260715-1400-sx2524-destination-coord-source --update B3 --status in_progress --owner fixer` lalu `--status completed --evidence .opencode/evidence/20260715-1400-sx2524-destination-coord-source/curl-after.json` |
| Q1 | `@quality-gate` | `.opencode/evidence/20260715-1400-sx2524-destination-coord-source/quality-gate.md` | `python3 ~/.config/opencode/scripts/task-progress.py 20260715-1400-sx2524-destination-coord-source --update Q1 --status in_progress --owner quality-gate` lalu `--status completed --evidence .opencode/evidence/20260715-1400-sx2524-destination-coord-source/quality-gate.md` |

Worklist:
1. **B1** | `@fixer` | Patch `pjp/repository/live_monitoring/get_distributor_repository.go` (join + 2 kolom) dan tambah regression test di `pjp/service/live_monitoring/get_distributor_service_test.go`. depends_on: none.
2. **B2** | `@fixer` | Jalankan `rtk go build` + `rtk go test -run Distributor -v` di `pjp/`; simpan log ke evidence. depends_on: B1.
3. **B3** | `@fixer` | Curl Staging dengan JWT + payload dari tiket; simpan response, bandingkan `arrive_*/leave_*/start/finish/skip_at` vs pre-fix. depends_on: B2.
4. **Q1** | `@quality-gate` | Signoff: SQL diff + test log + curl diff + invariants terpenuhi. depends_on: B3.

start_with: B1.

## Validation Commands
1. `cd pjp && rtk go build ./...` — exit 0; tidak ada warning baru.
2. `cd pjp && rtk go test ./service/live_monitoring/ -run TestTransformDistributorRows -v` — termasuk test regresi baru, semua PASS.
3. `cd pjp && rtk go test ./controller/live_monitoring/ -run TestGetDistributor -v` — semua PASS.
4. `cd pjp && rtk go test ./... -run Distributor -v` — semua PASS, tidak ada regresi di test lain.
5. `cd pjp && rtk go test ./...` — full module test, baseline regression.
6. Staging smoke (JWT dari tiket; tidak log secret value):
   ```bash
   TOKEN='<jwt-from-ticket>'
   curl 'https://best.scyllax.online/scylla-pjp/api/v1/live-monitoring-distributor?date=1784030400&status%5B%5D=Approved&status%5B%5D=Need+Review&emp_id=466&distributor_id=102' \
     -H "Authorization: Bearer $TOKEN" \
     -H 'Origin: https://staging.scyllax.online' \
     -H 'Referer: https://staging.scyllax.online/' \
     -H 'Accept: application/json, text/plain, */*' \
     -o .opencode/evidence/20260715-1400-sx2524-destination-coord-source/curl-after.json
   ```
7. Assert koordinat outlet (Python stdlib, tidak butuh jq):
   ```bash
   python3 -c "
   import json
   d=json.load(open('.opencode/evidence/20260715-1400-sx2524-destination-coord-source/curl-after.json'))
   def walk(n,p=''):
     if isinstance(n,dict):
       if 'fahmi' in n.get('destination_name','').lower():
         print(p,'lon=',n.get('longitude'),'lat=',n.get('latitude'))
       [walk(v,p+'/'+k) for k,v in n.items()]
     elif isinstance(n,list):
       [walk(x,p+f'[{i}]') for i,x in enumerate(n)]
   walk(d)
   "
   ```
8. SQL sanity di Staging DB (eksekusi lewat kredensial user; tidak simpan password di evidence):
   ```sql
   SELECT roh.outlet_id, roh.longitude AS hist_long, roh.latitude AS hist_lat,
          mo.longitude AS mst_long, mo.latitude AS mst_lat
   FROM pjp.route_outlet_history roh
   JOIN mst.m_outlet mo ON mo.outlet_id = roh.outlet_id
   WHERE roh.outlet_id IN (SELECT outlet_id FROM mst.m_outlet WHERE outlet_name ILIKE '%Fahmi Garage%')
     AND roh."date" = DATE '2026-07-14'
   ORDER BY roh.outlet_id;
   ```
   Ekspektasi: `hist_long/lat != mst_long/lat` (konfirmasi original bug) dan response JSON `longitude=106.82096445685606`, `latitude=-6.248543092691333` per ticket.

Expected output:
- Build: exit 0, no new warnings.
- Unit: all pass, new `TestTransformDistributorRows_UsesRowLongitudeLatitudeAsDestinationCoordinates` PASS.
- Curl Staging: prints a row with longitude/latitude equal to `mst.m_outlet` value (target JSON `longitude=106.82096445685606`, `latitude=-6.248543092691333` per ticket).
- Quick SQL sanity (run on Staging DB after fix, do not store creds in evidence):
  ```sql
  SELECT roh.outlet_id,
         roh.longitude AS hist_long, roh.latitude AS hist_lat,
         mo.longitude  AS mst_long,  mo.latitude  AS mst_lat
  FROM pjp.route_outlet_history roh
  JOIN mst.m_outlet mo ON mo.outlet_id = roh.outlet_id
  WHERE roh.outlet_id IN (SELECT outlet_id FROM mst.m_outlet WHERE outlet_name ILIKE '%Fahmi Garage%')
    AND roh."date" = DATE '2026-07-14'
  ORDER BY roh.outlet_id;
  ```
  Hist != Mst is the original bug; fix means response matches Mst.

## Evidence Requirements
- `index.json` (sudah ditulis di langkah discovery) — task-scoped evidence manifest.
- `result.md` — diff ringkas (file + line + before/after SQL), test log ringkas, link ke evidence lain.
- `curl-after.json` — response Staging setelah fix.
- `curl-before.json` — response Staging sebelum fix (opsional tapi direkomendasikan; ambil sebelum patch).
- `test.log` — output `rtk go test -v` untuk `Distributor` scope.
- `sql-diff.txt` — diff SQL (before/after) yang sudah di-ekstrak dari `get_distributor_repository.go` (manual, bukan dari git).

## Done Criteria
- Kode ter-patch: 2 ekspresi di line 129-130 repository, plus join `mst.m_outlet` jika belum ada.
- Regression test baru PASS.
- `rtk go build ./...` dan `rtk go test ./... -run Distributor` PASS tanpa regression.
- Curl Staging mengembalikan `destination_data[].longitude/latitude` = `mst.m_outlet` value untuk outlet yang terdampak.
- `arrive_*/leave_*/start/finish/skip_at` sama dengan sebelum fix.
- Evidence disimpan di `.opencode/evidence/20260715-1400-sx2524-destination-coord-source/`.
- Tracker progress di-update untuk B1/B2/B3/Q1.
- FE re-test di Staging sudah dikonfirmasi lewat koordinator `dev.fe` (di luar lane ini; tidak blocker BE signoff).

## Final Planning Summary
- **Mode**: maintenance-stability. Defect terisolasi, satu query, satu arah fix.
- **Artifacts created**:
  - `.opencode/plans/20260715-1400-sx2524-destination-coord-source.md` (this file)
  - `.opencode/evidence/20260715-1400-sx2524-destination-coord-source/index.json`
  - `.opencode/draft/20260715-1400-sx2524-destination-coord-source/` (kosong, placeholder; boleh dihapus saat final)
- **Key decisions**:
  - LEFT JOIN, bukan INNER (preserve orphan rows dengan coord 0).
  - `NULLIF(mo.longitude, '')` + `CAST AS FLOAT` + `COALESCE(..., 0)` (match pola `arrive_*`/`leave_*`).
  - Single query untuk `pjp_data` + `extra_call_data` (sama seperti sekarang).
  - Principal endpoint di luar scope.
- **Assumptions**:
  - Join key `roh.outlet_id = mo.outlet_id` benar.
  - FE tidak cache value lintas hari; boleh berubah source tanpa versioning.
  - FE menerima coord 0 untuk orphan outlet.
- **Open questions**: tidak ada yang memblokir. FE re-test pasca-deploy menjadi gate berikutnya, di luar lane planner.
- **Readiness**: `PASS_FOR_SLICE` (defect kecil, single-slice, terdefinisi jelas).
- **Readiness script fallback**: `python3 scripts/plan-execution-readiness.py` tidak ditemukan di repo (cek path). Fallback: gunakan output `validate-plan-depth.py` + `plan-compliance-check.py` + `subagent-handoff-check.py` sebagai readiness proxy. Catat di `result.md`.
- **Cleanup performed**: tidak ada draft evidence yang perlu dihapus; `index.json` adalah manifest resmi.
- **Source strategy**: repo-local evidence (kode) + ticket focused comment (user_confirmed) + GDoc (docs) + cURL Staging (runtime). Skip external web research — tidak perlu.
- **Active-lane reset**: eksekusi harus dilakukan oleh lane berikutnya (`@orchestrator` -> `@fixer`/quality gate). Pembatasan read-only planner tidak berlaku setelah plan di-hand off.
