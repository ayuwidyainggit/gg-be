# Plan SX-2034 — Fix Extra Call di Monitoring Activity (Principal)

Task ID: `20260521-1515-sx2034-extra-call-monitoring`
Sprint: SX Sprint 13. Reporter: product.management. Env: staging & demo.
Source of truth: file ini. Evidence: `.opencode/evidence/20260521-1515-sx2034-extra-call-monitoring/discovery.md`.

## Goal

Outlet Extra Call yang dibuat dari mobile harus muncul di endpoint `GET /scylla-pjp/api/v1/live-monitoring-principal` dan terlihat di peta Monitoring Activity Principal beserta route line-nya, untuk PJP `Approved`/`Need Review` pada tanggal terkait. Bug `destination_id` NULL pada `pjp_principles.destinations_history` saat create extra call ditutup, dan query principal monitoring memuat extra call (paritas dengan jalur distributor).

## Non-goals

- Tidak mengubah skema DB.
- Tidak mengubah kontrak request endpoint `live-monitoring-principal` (hanya menambah `extra_call_data` payload dalam struktur yang sudah ada).
- Tidak menyentuh jalur distributor monitoring (sudah benar).
- Tidak mengubah jalur create PJP awal (`SubmitPjpPrincipal`) yang sudah benar.
- Tidak refactor besar service `mobile.MOutletService`.

## Scope

- `mobile/repository/m_outlet.go`: perbaiki `StoreFromListPrinciple` agar `destination_id` diisi dengan `outlet.OutletId` (untuk outlet); siapkan branch distributor jika diperlukan (lihat Open Question).
- `mobile/service/m_outlet.go`: tambah validasi outlet ditemukan + (opsional) propagate destination type yang tepat.
- `pjp/repository/live_monitoring/get_principal_repository.go`: tambahkan loader baris extra call dari `pjp_principles.destinations_history` (filter `is_extra_call = true`) dengan resolve outlet/distributor lewat `destination_id`+`destination_type`, plus join `outlet_visit_list` untuk `start/finish/skip`.
- `pjp/model/live_monitoring.go`: tambah `IsExtraCall bool` + `Address` opsional di `LiveMonitoringPrincipalRow`.
- `pjp/service/live_monitoring/get_principal_service.go::transformPrincipalRows`: split rows ke `RouteData`/`ExtraCallData` (mengikuti pola distributor).
- Optional data fix: script SQL idempoten untuk backfill `destination_id` pada `destinations_history` di mana `is_extra_call = true` AND `destination_id IS NULL`.

## Requirements

1. Saat `POST /scylla-mobile/api/v1/m-outlets/from-list` (extra call) berjalan untuk schema principal, baris baru di `pjp_principles.destinations_history` wajib memiliki `destination_id` non-NULL sesuai `destination_type`.
2. `GET /scylla-pjp/api/v1/live-monitoring-principal?date=...&status[]=Approved&status[]=Need+Review&emp_id=482` mengembalikan `pjp_data[].extra_call_data[]` yang terisi untuk PJP yang punya extra call pada tanggal tsb.
3. Setiap entri extra call di response memuat: `route_code`, `route_name`, dan `destination_data[]` lengkap (id, code, type, name, longitude, latitude, address kalau ada di model, plus `start`/`finish`/`skip_at`/`arrive_*`).
4. PJP normal (`is_extra_call = false`/tidak ada di `destinations_history`) tetap muncul di `route_data` seperti semula tanpa duplikasi atau hilang.
5. Backfill data lama dijalankan terkontrol di staging dulu, dan tidak menyentuh data dengan `destination_id` yang sudah terisi.

## Acceptance Criteria

- [ ] Buat extra call baru via mobile → row baru di `pjp_principles.destinations_history` punya `destination_id = m_outlet.outlet_id` (atau `m_distributor.distributor_id` saat type distributor) dan `is_extra_call = true`.
- [ ] Hit `GET /scylla-pjp/api/v1/live-monitoring-principal?date=1779364800&status[]=Approved&status[]=Need+Review&emp_id=482` di staging → ada minimal 1 entri `extra_call_data` non-empty pada PJP `toko akbar` salesman 482, dengan outlet `toko principal`.
- [ ] PJP non-extra-call tetap utuh (verifikasi via response sebelum/sesudah pada salesman lain di hari yang sama).
- [ ] Snapshot SQL pada `destinations_history` menunjukkan tidak ada NULL baru pada `destination_id` untuk row `is_extra_call = true` setelah deploy.
- [ ] Unit test baru lulus: split principal rows menghasilkan `RouteData` vs `ExtraCallData` sesuai `IsExtraCall`.

## Existing Patterns / Reuse

- Reuse pola distributor: `pjp/repository/live_monitoring/get_distributor_repository.go` (`roh.is_extra_call`), `get_distributor_service.go` (`extraRouteMap`/`pjp.ExtraCallData`). Pattern langsung dipadankan ke principal.
- Reuse pola JOIN outlet/distributor lewat `destinations_history`: `mobile/repository/pjp_principal.go` baris 148-176, 210-211.
- Reuse pola insert benar: `mobile/service/pjp_principal.go::SubmitPjpPrincipal` (`DestinationId: outlet.OutletId` / `distributor.DistributorId`).
- Reuse `transformPrincipalRows` yang ada — tidak buat path transform baru, hanya perpanjang.

## Constraints

- Layer wajib Controller→Service→Repository→DB; transactional write tetap di service.
- `pjp` service Fiber-based; ikuti gaya GORM existing (`tx.WithContext(ctx).Table(...).Joins(...)`).
- Tx-context extraction harus dihormati saat memodifikasi repo write.
- Tenant isolation: query principal sudah scope via `cust_id` parent + child cust_ids → tetap dipakai. Jangan kendor.
- `salesman_id` Jira (`482`) = `pjp.salesman_id` = `emp_id`; jangan tukar.
- Hindari N+1 saat menambah loader extra call: lakukan single-query LEFT JOIN, gabung di service.

## Risks

- Mengubah base query principal beresiko regresi PJP normal. Mitigasi: tambah loader extra call sebagai query terpisah dan merge di service, biarkan query existing tidak diubah.
- Backfill SQL bisa salah join key. Mitigasi: dry-run `SELECT` count dulu, lalu UPDATE pada staging dengan transaction + LIMIT/CTE checked.
- Distributor-typed extra call tidak ada di endpoint `from-list` saat ini (entity hanya bawa `outlet_id`). Mitigasi: cabang ditangani di repo via `destination_type`-aware update; sampai FE/BE sepakat distributor extra call, jalur principal di-fix untuk outlet dulu, query monitoring tetap kompatibel jika type distributor muncul.
- Data lama dengan `destination_id` salah (misal terisi `OldPjpId` non-nil) bisa lolos backfill. Mitigasi: backfill hanya target rows yang `destination_id IS NULL` AND `is_extra_call = true`; rows non-NULL salah perlu remediation manual terpisah.

## Decisions / Assumptions

Decisions:
- Repair di sumber: ganti bind value kolom `destination_id` di `StoreFromListPrinciple` dari `outlet.OldPjpId` → `outlet.OutletId`. Pengembangan minimal, sesuai bug.
- Untuk monitoring, tambah loader query `GetPrincipalExtraCalls` yang membaca `pjp_principles.destinations_history` dengan `is_extra_call = true` + LEFT JOIN `m_outlet`/`m_distributor` + LEFT JOIN `outlet_visit_list ovl` (cocokkan `pjp_id`, `date`, `outlet_id`, `is_extra_call=true`).
- Service `GetPrincipalMonitoring` memanggil dua loader (existing + extra call), lalu `transformPrincipalRows` menerima flag `IsExtraCall` per row dan memetakan ke `RouteData` (false) atau `ExtraCallData` (true).

Assumptions / Open Questions (lihat juga `.opencode/draft/.../open-questions.md`):
- A1: extra call principal saat ini hanya outlet (sesuai `entity.ExtraCallOutlet.OutletIDs`). Distributor extra call dianggap belum live untuk principal; query monitoring tetap inklusif bila type distributor muncul.
- A2: `outlet_visit_list` untuk principal extra call sudah mengisi `outlet_id` (dikonfirmasi dari `StoreFromListOutletVisitListPrinciple`). JOIN ke ovl pakai `ovl.outlet_id = dh.destination_id AND ovl.pjp_id = dh.pjp_id AND ovl.date = dh.date AND ovl.is_extra_call = true`.
- A3: `destinations_history.date` adalah `timestamp/date` dengan timezone serupa input mobile. Filter tanggal monitoring memakai `DATE(dh.date) = ?`.
- Q1: backfill data lama → jalankan? (Default rekomendasi: ya, untuk staging dulu, kemudian production setelah verifikasi.)
- Q2: distributor extra call → masuk scope sekarang atau follow-up issue?

## TDD / Test Plan

TDD wajib (perubahan logika DB write + query + transform).

Existing tests:
- `pjp/service/live_monitoring/get_principal_service_test.go` (stub repo + transform).
- `pjp/service/live_monitoring/get_distributor_service_test.go` (pola split route vs extra call).

Red (gagal dulu):
1. `pjp/service/live_monitoring/get_principal_service_test.go`: tambah test baru `TestGetPrincipalMonitoring_PopulatesExtraCallData` — stub repo mengembalikan rows mix `IsExtraCall=true/false`, ekspektasi: `result[0].PjpData[0].ExtraCallData` punya 1 entri dengan outlet name/coord, dan `RouteData` tidak kemasukan row tsb.
2. (Opsional) repo-level test integratif jika harness ada; jika tidak, validasi via Postman + DB query.
3. `mobile/service` (atau repo) — test bahwa setelah `StoreFromList` untuk skema principal, value bind ke `destination_id` adalah `OutletId`. Jika repo-level test belum ada, tambahkan unit test ringan via tabel mock atau verifikasi argumen via wrapper. Jika tidak feasible cepat, validasi manual via DB cek + test integrasi end-to-end.

Green:
- Implementasi loader extra call + split di transform.
- Ganti bind value `$11` di `StoreFromListPrinciple`.

Refactor:
- Ekstrak helper `splitPrincipalRowsByExtraCall` (mirror pattern distributor) supaya `transformPrincipalRows` tetap pendek.
- Pastikan struktur `LiveMonitoringPrincipalRow` tidak duplikasi tag.

Edge cases:
- Salesman tidak punya extra call → `extra_call_data: []` (bukan null).
- Extra call ada tapi `outlet_visit_list` belum (outlet baru saja dibuat, belum visited) → tetap muncul dengan koordinat dari `dh`, `start/finish` null.
- Outlet tidak ditemukan (`mst.m_outlet` row missing) → row di-skip, log warning ringan.
- Tanggal di luar rentang attendance/business day → tidak menambahkan extra call.
- `is_extra_call = true` tapi `destination_id` masih NULL (data lama tanpa backfill) → tidak muncul; backfill diharapkan menyelesaikan.

Validation commands:
- `rtk go mod download && rtk go mod tidy` (per service).
- `rtk go test ./service/live_monitoring/... ./repository/live_monitoring/...` di dalam `pjp/`.
- `rtk go test ./service/... ./repository/...` di dalam `mobile/` untuk path yang tersentuh.
- Manual: cURL ke endpoint dengan payload Jira, verifikasi response.

## Implementation Steps

1. Tambah field `IsExtraCall bool` (gorm:"column:is_extra_call") + (opsional) `DestinationAddress string` ke `pjp/model/live_monitoring.go::LiveMonitoringPrincipalRow`.
2. Tambah method baru di `pjp/repository/live_monitoring/get_principal_repository.go` (atau file baru `get_principal_extra_call_repository.go`):
   - `GetPrincipalExtraCalls(ctx, tx, custIDs, date, regionID, areaID, distributorID, empIDs, statuses) ([]model.LiveMonitoringPrincipalRow, error)`
   - Query: `FROM pjp_principles.destinations_history dh JOIN pjp_principles.permanent_journey_plans pjp ON pjp.id = dh.pjp_id AND dh.cust_id = pjp.cust_id JOIN pjp_principles.routes r ON r.route_code = dh.route_code LEFT JOIN mst.m_outlet mo ON mo.outlet_id = dh.destination_id AND dh.destination_type = 'outlet' LEFT JOIN mst.m_distributor md ON md.distributor_id = dh.destination_id AND dh.destination_type = 'distributor' JOIN mst.m_salesman ms ON ms.emp_id = pjp.salesman_id JOIN mst.m_employee me ON me.emp_id = ms.emp_id JOIN smc.m_customer mc ON ms.cust_id = mc.cust_id LEFT JOIN mst.m_distributor mdcust ON mdcust.distributor_id = mc.distributor_id LEFT JOIN pjp_principles.outlet_visit_list ovl ON ovl.outlet_id = dh.destination_id AND ovl.pjp_id = dh.pjp_id AND ovl.date = DATE(dh.date) AND ovl.is_extra_call = true`
   - SELECT field-list paritas dengan loader existing + `true AS is_extra_call`, gunakan COALESCE supaya outlet name/coord fallback ke `dh.*` saat distributor.
   - WHERE: `dh.is_extra_call = true`, `pjp.cust_id IN (...)`, `DATE(dh.date) = ?`, `pjp.approval_status IN ?`, plus filter region/area/distributor/empIDs sama.
3. Tambah method ke interface `LiveMonitoringRepository` (`pjp/repository/live_monitoring/live_monitoring_repository.go`) + adjust stubs (`*_test.go`).
4. Update `pjp/service/live_monitoring/get_principal_service.go::GetPrincipalMonitoring`:
   - Setelah `GetPrincipalMonitoring` rows existing, panggil `GetPrincipalExtraCalls` untuk paged employees.
   - Gabungkan rows lalu tandai sumbernya (`IsExtraCall`).
   - `transformPrincipalRows` direfactor: rows non-extra-call → `RouteData[]`; rows extra-call → `ExtraCallData[]`.
5. Edit `mobile/repository/m_outlet.go::StoreFromListPrinciple`:
   - Ganti bind value `$11` dari `outlet.OldPjpId` ke `outlet.OutletId`.
   - Tambah guard: kalau `outlet.OutletId == 0` → return error eksplisit (jangan biarkan NULL/0).
   - Catatan: kolom `old_pjp_id` ($21) sudah benar pakai `outlet.OldPjpId`, tidak perlu dipindah.
6. (Opsional) Tambah validasi di `mobile/service/m_outlet.go::StoreFromList`: jika `len(outlets) != len(request.OutletIDs)` → return error agar partial insert tidak terjadi.
7. (Opsional) Tambah skenario branching distributor di `StoreFromListPrinciple` jika scope diperluas (lihat Open Question Q2). Jika tidak, dokumentasikan.
8. Backfill SQL (lihat `Validation Commands` & file SQL terlampir di catatan plan): jalankan di staging → verifikasi count → production.
9. Jalankan unit test `pjp` + smoke API.
10. Lewatkan ke `@quality-gate` untuk security/regresi sign-off.

## Expected Files to Change

- `pjp/model/live_monitoring.go` (tambah `IsExtraCall`).
- `pjp/repository/live_monitoring/live_monitoring_repository.go` (interface + stub).
- `pjp/repository/live_monitoring/get_principal_repository.go` (atau file baru `get_principal_extra_call_repository.go`).
- `pjp/service/live_monitoring/get_principal_service.go` (panggil loader + split transform).
- `pjp/service/live_monitoring/get_principal_service_test.go` (test baru).
- `pjp/service/live_monitoring/get_detail_service_test.go` & `get_distributor_service_test.go` (sinkronisasi stub interface bila perlu).
- `mobile/repository/m_outlet.go` (`StoreFromListPrinciple` fix bind).
- `mobile/service/m_outlet.go` (validasi outlet length, opsional).
- (Opsional) `scripts/sql/sx2034-backfill-destinations-history.sql` (jika repo punya tempat untuk SQL ad-hoc).

## Agent / Tool Routing

- `@orchestrator`: integrasi & routing.
- `@fixer`: implementasi semua perubahan kode + test (Red-Green-Refactor).
- `@oracle`: review query design (regresi PJP normal, performance JOIN baru).
- `@quality-gate`: security/regresi sign-off, validasi tenant scope, tinjau backfill SQL sebelum staging.
- `@librarian`: tidak diperlukan (no external API doc lookup).
- `@architect`: tidak diperlukan (perubahan tidak mengubah platform boundary).

## Execution-ready Worklist / Handoff Contract

start_with: `T1`
Semua task `ready` kecuali yang di-mark `blocked`.

```
- id: T1
  action: Tambah field IsExtraCall (dan optional DestinationAddress) ke LiveMonitoringPrincipalRow di pjp/model/live_monitoring.go
  depends_on: none
  owner: @fixer
  validation: rtk go build ./... (dalam pjp/)
  exit: Struct field tersedia, kompil sukses
  status: ready
  requires_user_decision: no
- id: T2
  action: Tambah method GetPrincipalExtraCalls ke interface + implementasi repo (file baru get_principal_extra_call_repository.go) + sinkronkan stub di test files
  depends_on: T1
  owner: @fixer
  validation: rtk go vet ./... + rtk go test ./repository/live_monitoring/... -run Compile (compile-only ok)
  exit: Method baru tersedia di interface dan terpasang di liveMonitoringRepository
  status: ready
  requires_user_decision: no
- id: T3
  action: Tambah unit test merah TestGetPrincipalMonitoring_PopulatesExtraCallData di get_principal_service_test.go (stub mengembalikan rows extra call)
  depends_on: T2
  owner: @fixer
  validation: rtk go test ./service/live_monitoring/... -run TestGetPrincipalMonitoring_PopulatesExtraCallData (harus FAIL dulu)
  exit: Test gagal sebagaimana diharapkan (Red)
  status: ready
  requires_user_decision: no
- id: T4
  action: Update get_principal_service.go untuk memanggil GetPrincipalExtraCalls dan split transformPrincipalRows ke RouteData vs ExtraCallData
  depends_on: T3
  owner: @fixer
  validation: rtk go test ./service/live_monitoring/... (semua hijau)
  exit: Test T3 hijau, regresi test lain tetap hijau (Green)
  status: ready
  requires_user_decision: no
- id: T5
  action: Refactor helper splitPrincipalRowsByExtraCall (atau equivalent) supaya transform tetap ringkas, tanpa mengubah behaviour
  depends_on: T4
  owner: @fixer
  validation: rtk go test ./service/live_monitoring/...
  exit: Hijau, struktur lebih bersih
  status: ready
  requires_user_decision: no
- id: T6
  action: Fix mobile/repository/m_outlet.go::StoreFromListPrinciple — ganti bind $11 menjadi outlet.OutletId; tambah guard outlet.OutletId != 0
  depends_on: none
  owner: @fixer
  validation: rtk go build ./... (mobile/) + manual cURL POST /m-outlets/from-list di staging → SELECT destination_id terisi
  exit: Insert baru selalu set destination_id = m_outlet.outlet_id
  status: ready
  requires_user_decision: no
- id: T7
  action: Tambah validasi length outlets di mobile/service/m_outlet.go::StoreFromList (opsional tapi disarankan) — return error jika len(outlets) != len(request.OutletIDs)
  depends_on: T6
  owner: @fixer
  validation: rtk go test ./service/... (mobile/) atau manual integrasi
  exit: Partial insert tidak mungkin terjadi
  status: ready
  requires_user_decision: no
- id: T8
  action: Susun script backfill .sql idempoten untuk destinations_history rows is_extra_call=true AND destination_id IS NULL (lihat catatan SQL di plan)
  depends_on: T6
  owner: @fixer
  validation: SELECT count sebelum/sesudah di staging; UPDATE harus 0 setelah idempotent rerun
  exit: Skrip review-ready
  status: ready
  requires_user_decision: yes
- id: T9
  action: Jalankan backfill di staging (tx, dry-run SELECT dulu) dan verifikasi via query debug Jira pada salesman_id=482, date=2026-05-21
  depends_on: T8
  owner: @quality-gate
  validation: query Jira menunjukkan destination_id terisi; cURL endpoint live-monitoring-principal mengembalikan extra_call_data terisi
  exit: Bukti screenshot/log respons disimpan ke .opencode/evidence/
  status: blocked
  blocker: Menunggu jawaban Q1 (apakah jalankan backfill di staging) dan akses DB
  requires_user_decision: yes
- id: T10
  action: Smoke API end-to-end di staging dengan emp_id=482 + status Approved/Need Review + date=1779364800; capture before/after JSON
  depends_on: T4, T6, T9
  owner: @quality-gate
  validation: extra_call_data berisi outlet "toko principal"; PJP "toko akbar" tetap ada di route_data
  exit: Bukti respons disimpan
  status: ready
  requires_user_decision: no
- id: T11
  action: Final review oleh @quality-gate (regresi PJP normal, security/tenant, perfQA join baru) dan sign-off
  depends_on: T10
  owner: @quality-gate
  validation: Checklist QUALITY.md terpenuhi
  exit: Sign-off & PR siap merge
  status: ready
  requires_user_decision: no
```

## Validation Commands

- `rtk docker compose -f docker-compose.yml ps` (pastikan service up).
- `rtk go mod download && rtk go mod tidy` di `pjp/` dan `mobile/`.
- `rtk go test ./service/live_monitoring/... ./repository/live_monitoring/...` di `pjp/`.
- `rtk go test ./service/... ./repository/...` di `mobile/` (lingkup tersentuh).
- cURL staging:
  - `curl -G "https://best.scyllax.online/scylla-pjp/api/v1/live-monitoring-principal" --data-urlencode "date=1779364800" --data-urlencode "status[]=Approved" --data-urlencode "status[]=Need Review" --data-urlencode "emp_id=482" -H "Authorization: Bearer <token>"`.
- DB verifikasi (gunakan query debug Jira, salesman_id=482, date=2026-05-21).

## Evidence Requirements

- Sebelum implementasi: snapshot SELECT count `is_extra_call=true AND destination_id IS NULL` di `pjp_principles.destinations_history` (staging).
- Setelah backfill (jika dijalankan): SELECT count sama harus 0; query debug Jira menampilkan baris extra call dengan `destination_id` terisi.
- Setelah implementasi: simpan response JSON `before` dan `after` endpoint ke `.opencode/evidence/<task-id>/api-response-before.json` & `.../api-response-after.json`.
- Catat command yang dijalankan + viewport tidak relevan (BE only).

## Done Criteria

- Semua tasks T1-T11 `completed` (T9 jika user setuju backfill).
- Tests baru hijau, regresi tidak red.
- Bukti API menampilkan `extra_call_data` terisi untuk salesman 482 pada 2026-05-21.
- `@quality-gate` sign-off.

## Final Planning Summary

- Plan path: `.opencode/plans/20260521-1515-sx2034-extra-call-monitoring.md` (file ini, source of truth).
- Evidence path: `.opencode/evidence/20260521-1515-sx2034-extra-call-monitoring/discovery.md` (dipertahankan karena query debug Jira & smoking gun mapping diperlukan untuk implementasi & QG).
- Open questions: lihat `.opencode/draft/20260521-1515-sx2034-extra-call-monitoring/open-questions.md` (Q1 backfill, Q2 distributor extra call). Plan utama tetap actionable dengan asumsi default.
- Decisions: dual-loader pattern (existing + extra call) di repo principal; fix bind kolom `destination_id`; reuse pola distributor untuk split response.
- MCP/tools used: hanya local code search; tidak panggil context7/brave/GitHub karena root cause sepenuhnya internal.
- Cleanup: draft & evidence dipertahankan karena masih relevan untuk implementasi & verifikasi DB; akan di-GC setelah QG sign-off.
- Readiness: Implementation-ready. Hanya T9 menunggu jawaban Q1 dari user (boleh dijalankan secara default di staging dengan persetujuan tertulis).

## Lampiran — Backfill SQL (review dulu sebelum jalan)

```sql
-- DRY RUN: hitung kandidat
SELECT count(*)
FROM pjp_principles.destinations_history dh
WHERE dh.is_extra_call = true
  AND dh.destination_id IS NULL;

-- Backfill outlet (kolom kunci: destination_code = m_outlet.outlet_code, samakan cust_id jika perlu)
BEGIN;
UPDATE pjp_principles.destinations_history dh
SET destination_id = mo.outlet_id
FROM mst.m_outlet mo
WHERE dh.is_extra_call = true
  AND dh.destination_type = 'outlet'
  AND dh.destination_id IS NULL
  AND dh.destination_code = mo.outlet_code;
-- Verifikasi:
-- SELECT count(*) FROM pjp_principles.destinations_history WHERE is_extra_call=true AND destination_type='outlet' AND destination_id IS NULL;
COMMIT;

-- Backfill distributor (jika scope diperluas)
BEGIN;
UPDATE pjp_principles.destinations_history dh
SET destination_id = md.distributor_id
FROM mst.m_distributor md
WHERE dh.is_extra_call = true
  AND dh.destination_type = 'distributor'
  AND dh.destination_id IS NULL
  AND dh.destination_code = md.distributor_code;
COMMIT;
```

> ⚠️ Validasi join key di staging dulu (kemungkinan `destination_code` cocok dengan `outlet_code`/`distributor_code`; konfirmasi sebelum dijalankan production).
