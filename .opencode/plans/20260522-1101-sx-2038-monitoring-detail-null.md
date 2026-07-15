# Plan SX-2038 — Live Monitoring Detail Null Data

Task ID: `20260522-1101-sx-2038-monitoring-detail-null`
Jira: `SX-2038` — `[Defect][BE] Live Monitoring - View Details display with null data`
Module: `pjp` / Live Monitoring Detail
Endpoint: `GET /scylla-pjp/api/v1/monitoring_locations/details`
Source of truth: file ini.
Evidence: `.opencode/evidence/20260522-1101-sx-2038-monitoring-detail-null/discovery.md`.
Open questions: `.opencode/draft/20260522-1101-sx-2038-monitoring-detail-null/open-questions.md`.

## Goal

Perbaiki endpoint `GET /monitoring_locations/details` agar request seperti `emp_id=484&date=2026-05-22` atau `emp_id=482&date=2026-05-22` tidak menghasilkan `data: null` ketika data PJP valid ada. Response harus berisi `visit_information` dengan summary `planned`, `on_going`, `extra_call`, `visited`, `skipped`, plus section detail existing (`sales`, `return`, `collection`, `expense`, `shipment`) tanpa breaking change endpoint monitoring lain.

## Non-goals

- Tidak mengubah skema DB.
- Tidak mengubah route, auth, atau format parameter endpoint.
- Tidak menjalankan backfill SX-2034 di plan ini.
- Tidak mengubah list endpoint `live-monitoring-principal` / `live-monitoring-distributor` kecuali compile impact dari interface/test.
- Tidak mengganti kontrak response final sebelum FE/product setuju.

## Scope

- Audit dan fix jalur Principal detail: `pjp/service/live_monitoring/get_detail_service.go` dan `pjp/repository/live_monitoring/get_detail_repository.go`.
- Tambah test unit/regression untuk Principal detail supaya `Need Review`, rows `destinations_history`, dan mismatch data tidak kembali jadi `nil` diam-diam.
- Pertahankan jalur Distributor detail yang sudah memakai `IN ('Approved','Need Review')` dan counter terpisah.
- Jalankan query debug staging untuk membuktikan root cause aktual: `emp_id=484/482`, PJP status, `destinations_history`, `outlet_visit_list`, dan hasil query dev.fe.
- Review response no-data: default tetap mengikuti kontrak existing (`data: null`, `message: No Data`) kecuali FE menyetujui perubahan ke `data: []`.

## Requirements

1. Endpoint detail principal wajib menerima `emp_id` dari query dan mencari PJP sesuai domain existing: `pjp.salesman_id == mst.m_salesman.emp_id`.
2. Query detail principal wajib menampilkan PJP dengan `approval_status IN ('Approved', 'Need Review')`, selaras query dev.fe dan jalur distributor.
3. Query detail principal wajib tidak drop semua data hanya karena data regular route/destination tidak lengkap, selama `destinations_history` pada tanggal tersebut ada.
4. Extra-call principal dari `pjp_principles.destinations_history` wajib dihitung sebagai `extra_call`, dengan status `on_going`, `visited`, `skipped` dari `outlet_visit_list` bila match.
5. `cust_id` scope wajib tetap ketat: parent principal → `GetChildCustIDs` → filter `m_salesman.cust_id IN childCustIDs` / `pjp.cust_id` sesuai pola repo.
6. `date` tetap string `YYYY-MM-DD`; jangan ubah ke epoch untuk endpoint detail.
7. Kalau data benar-benar tidak ada, endpoint tetap informatif dan tidak crash.

## Acceptance Criteria

- [ ] `GET /scylla-pjp/api/v1/monitoring_locations/details?emp_id=484&date=2026-05-22` dengan token `princessa@gmail.com` mengembalikan `data` non-null bila staging DB memiliki PJP valid untuk `484`.
- [ ] `GET /scylla-pjp/api/v1/monitoring_locations/details?emp_id=482&date=2026-05-22` diverifikasi juga; bila `482` adalah ID valid, response berisi data.
- [ ] `visit_information` berisi nilai benar untuk `planned`, `on_going`, `extra_call`, `visited`, `skipped`.
- [ ] PJP `Need Review` tidak lagi hilang dari detail principal.
- [ ] Extra-call dari `pjp_principles.destinations_history` tidak hilang karena `destination_id` NULL atau tidak ada row regular `destinations`.
- [ ] Endpoint distributor detail tetap hijau lewat test existing.
- [ ] Tidak ada regresi pada `live-monitoring-principal` dan `live-monitoring-distributor`.

## Existing Patterns/Reuse

- Reuse `GetChildCustIDs` untuk tenant scope principal.
- Reuse pola status distributor detail: `approval_status IN ('Approved', 'Need Review')`.
- Reuse pola list principal SX-2034: loader regular + loader extra-call dari `destinations_history`.
- Reuse `model.VisitInformationRow` untuk hasil count detail.
- Reuse stub test `detailRepoStub` di `get_detail_service_test.go`; tambahkan field/method minimal agar Principal path bisa dites.
- Reuse `GetSalesmanCustID` untuk filter sales/returns/expenses/shipment setelah `visitInfo` ditemukan.
- Tidak ada util KiloCode/project lain yang menyelesaikan mapping ini langsung; perbaikan harus di repo `pjp`.

## Constraints

- Repo multi-module Go; target module `pjp`.
- Layer wajib Controller → Service → Repository → DB.
- Query harus tetap tenant-safe; jangan longgarkan `cust_id` demi membuat data muncul.
- Repo menggunakan GORM, `Take()` untuk single-row not-found, dan `Find()` untuk slice.
- `AGENTS.md` repo mewajibkan `rtk` prefix untuk shell workflow.
- Planner tidak boleh edit source; implementasi dilakukan setelah plan oleh `@orchestrator`/`@fixer`.

## Risks

- Status filter berubah dari `Approved` saja ke `Approved/Need Review`; jika FE sebelumnya sengaja hanya detail approved, behaviour berubah. Mitigasi: dev.fe query dan jalur distributor sudah memakai dua status.
- Switching sumber principal detail ke `destinations_history` bisa double-count regular route jika digabung tanpa guard. Mitigasi: agregasi jelas: planned = `is_extra_call=false`, extra_call = `is_extra_call=true`.
- `emp_id=484` vs `482` mungkin data QA salah, bukan bug code. Mitigasi: staging DB query wajib sebelum/selama implementasi.
- `destination_id` NULL tetap bisa membuat join ke `mst.m_outlet` gagal jika query butuh outlet metadata. Mitigasi: pakai data denormalisasi di `destinations_history` dan `LEFT JOIN` hanya untuk enrichment.
- Query agregat baru bisa berbeda dari kontrak lama `GetVisitInformationPrincipal` berbasis `destinations`. Mitigasi: unit test + API smoke sebelum/after.

## Decisions/Assumptions

Decisions:

- Fix utama bukan lookup `emp_id → salesman_id` baru. Codebase saat ini memperlakukan `pjp.salesman_id` sebagai `mst.m_salesman.emp_id`; mapping tambahan berisiko salah kecuali DB membuktikan sebaliknya.
- Detail Principal harus pakai `approval_status IN ('Approved', 'Need Review')`.
- Untuk principal detail, rekomendasi query agregat baru berbasis `pjp_principles.destinations_history dh` sebagai sumber utama karena issue terkait extra-call dan query dev.fe memakai tabel itu.
- `mst.m_outlet` di query detail harus `LEFT JOIN` bila dipakai, bukan `JOIN`, supaya `dh.destination_id NULL` tidak drop semua rows.
- Response no-data tidak diubah di implementasi pertama; perubahan `data: []` butuh kontrak FE/product.

Assumptions / Open Questions:

- A1: `date=2026-05-22` sudah benar format `YYYY-MM-DD` untuk endpoint detail.
- A2: `cust_id=C26002` dari token principal punya child customer yang mencakup salesman target.
- A3: `Need Review` harus visible di detail karena query dev.fe dan jalur distributor menganggapnya valid.
- Q1: ID valid untuk reproduce adalah `482`, `484`, atau keduanya? Lihat `.opencode/draft/.../open-questions.md`.
- Q2: Product/FE setuju response no-data tetap `null` atau mau kontrak baru `[]`? Default: tidak ubah kontrak.

## TDD/Test Plan

TDD wajib karena ini bug logic/query dan endpoint behaviour.

Existing test patterns:

- `pjp/service/live_monitoring/get_detail_service_test.go` punya stub repo dan test distributor/expense.
- `pjp/service/live_monitoring/get_principal_service_test.go` punya pola principal monitoring + extra-call split.
- Repo-level SQL unit belum tampak; fokus awal service tests + compile + staging DB smoke.

Red step:

1. Tambah test `TestGetMonitoringDetail_PrincipalAllowsNeedReview` di `get_detail_service_test.go`.
   - Stub `GetVisitInformationPrincipal` harus menerima path principal dan return row walau status logical `Need Review` disimulasikan lewat repo method baru/hasil stub.
   - Test awal gagal karena current repo filter hanya `Approved` tidak bisa dibuktikan lewat service stub saja; jika ingin murni Red, tambah repository test dengan `sqlmock` tidak tersedia. Alternatif: buat test service untuk method baru `GetVisitInformationPrincipalFromHistory` setelah interface ditambah; awal compile fail.
2. Tambah test `TestGetMonitoringDetail_PrincipalIncludesExtraCallSummary`.
   - Stub row: `Plan=1`, `ExtraCall=1`, `OnGoing=0`, `Visited=1`, `TotalSkip=0`.
   - Ekspektasi `result.VisitInformation.ExtraCall == 1` dan tidak `nil`.
3. Tambah test `TestGetMonitoringDetail_PrincipalNoDataStillNil`.
   - Stub no rows → result `nil`, err `nil` tetap.
4. Tambah test untuk `EmpID` passthrough: repo menerima `empID=484` dan tidak mengubah ke `salesman_id` lain tanpa resolver.

Green step:

- Implement repo query/counter principal detail baru berbasis `destinations_history` + status `IN`.
- Service `getPrincipalVisitInfo` memanggil method baru atau method lama yang sudah diganti internal query-nya.
- Sesuaikan stub interface di tests.

Refactor step:

- Ekstrak status slice constant/helper bila ada duplikasi: `[]string{constant.ApprovalStatusApproved, "Need Review"}`.
- Jaga query tetap pendek dengan helper builder untuk principal detail scope kalau perlu.
- Jangan refactor controller response di sprint fix ini kecuali FE setuju.

Edge cases:

- `emp_id` tidak punya PJP pada tanggal tersebut → tetap `data: null`, `message: No Data`.
- PJP `Need Review` → muncul.
- `dh.destination_id NULL`, `dh.destination_code` ada → tetap dihitung.
- `outlet_visit_list` tidak ada → planned/extra_call tetap dihitung, visit counters 0.
- `skip_at` dan `leave_at` sama-sama ada → ikuti definisi dev.fe (`skip_at` menghitung skipped; visited menghitung arrive+leave).

Commands:

- `rtk go test ./service/live_monitoring -run TestGetMonitoringDetail`
- `rtk go test ./repository/live_monitoring ./service/live_monitoring`
- `rtk go test ./...`

## Implementation Steps

1. Jalankan DB debug staging dari prompt untuk `emp_id IN (484,482)`, `date='2026-05-22'`, `cust_id='C26002'`; simpan hasil ringkas ke evidence implementasi.
2. Tambah method repo principal detail history, misal:
   - `GetVisitInformationPrincipalFromHistory(ctx, tx, custIDs, date, empID, statuses) (*model.VisitInformationRow, error)`.
3. Query method baru:

```sql
SELECT
  pjp.salesman_id AS emp_id,
  me.emp_code,
  me.emp_name,
  SUM(CASE WHEN dh.is_extra_call = false THEN 1 ELSE 0 END) AS plan,
  SUM(CASE WHEN ovl.arrive_at IS NOT NULL AND ovl.leave_at IS NULL AND ovl.skip_at IS NULL THEN 1 ELSE 0 END) AS on_going,
  SUM(CASE WHEN dh.is_extra_call = true THEN 1 ELSE 0 END) AS extra_call,
  SUM(CASE WHEN ovl.arrive_at IS NOT NULL AND ovl.leave_at IS NOT NULL THEN 1 ELSE 0 END) AS visited,
  SUM(CASE WHEN ovl.skip_at IS NOT NULL THEN 1 ELSE 0 END) AS total_skip,
  COUNT(ovl.outlet_code) AS matched
FROM pjp_principles.destinations_history dh
JOIN pjp_principles.permanent_journey_plans pjp
  ON pjp.id = dh.pjp_id
 AND pjp.cust_id = dh.cust_id
JOIN mst.m_salesman ms
  ON ms.emp_id = pjp.salesman_id
LEFT JOIN mst.m_employee me
  ON me.emp_id = ms.emp_id
LEFT JOIN mst.m_outlet mo
  ON mo.outlet_id = dh.destination_id
LEFT JOIN pjp_principles.outlet_visit_list ovl
  ON ovl.pjp_id = dh.pjp_id
 AND DATE(ovl.date) = DATE(dh.date)
 AND (
      ovl.outlet_id = dh.destination_id
      OR ovl.outlet_code = dh.destination_code
 )
WHERE pjp.salesman_id IN (SELECT emp_id FROM mst.m_salesman WHERE cust_id IN (?))
  AND pjp.salesman_id = ?
  AND DATE(dh.date) = ?
  AND pjp.approval_status IN ('Approved', 'Need Review')
GROUP BY pjp.salesman_id, me.emp_code, me.emp_name;
```

4. Kalau DB membuktikan `destinations_history` tidak bisa jadi sumber untuk planned regular, fallback minimal: ubah method existing `GetVisitInformationPrincipal` status filter ke `IN`, ubah `JOIN mst.m_distributor` ke `LEFT JOIN`, lalu tambah counter extra-call dari `destinations_history` dan merge di service.
5. Update interface `LiveMonitoringRepository` dan stub tests.
6. Update `getPrincipalVisitInfo` agar memanggil method baru dan mengisi `VisitInformationData` dari `VisitInformationRow.ExtraCall`, bukan hitung `CountTotalVisitsPrincipal - Matched` bila path baru sudah akurat.
7. Pertahankan `GetActivityTime`, `GetUserFullname`, dan downstream `GetSales/GetReturns/GetExpenses/GetShipments` behaviour.
8. Tambah unit tests Red → Green → Refactor.
9. Jalankan validation commands.
10. Smoke staging dengan token QA dan simpan response before/after.
11. Minta `@quality-gate` review tenant scope, query status, no-data response, dan regression.

## Expected Files to Change

- `pjp/repository/live_monitoring/live_monitoring_repository.go` — tambah method principal detail history (atau update signature bila memilih replace).
- `pjp/repository/live_monitoring/get_detail_repository.go` — fix query principal detail (`IN` status, `LEFT JOIN`, history source/counter extra-call).
- `pjp/service/live_monitoring/get_detail_service.go` — panggil repo method baru, map `ExtraCall` langsung.
- `pjp/service/live_monitoring/get_detail_service_test.go` — stub + test principal detail regression.
- Opsional: `pjp/constant/pjp_constant.go` — tambah constant `ApprovalStatusNeedReview` atau helper statuses jika project setuju.
- Opsional: `pjp/repository/live_monitoring/get_detail_repository_test.go` — hanya bila menambah SQL-level test feasible tanpa dependency baru.

## Agent/Tool Routing

- `@orchestrator`: mulai eksekusi dari plan ini, pecah ke lane tepat, integrasi final.
- `@fixer`: implementasi bounded di module `pjp`, TDD, unit tests, staging smoke bila akses ada.
- `@oracle`: review query design bila DB evidence menunjukkan perlu ganti sumber data besar dari `destinations` ke `destinations_history`.
- `@quality-gate`: final signoff untuk tenant scope, security, API contract, regression.
- `@librarian`: tidak diperlukan; tidak ada external/library behaviour baru.
- `@architect`: tidak diperlukan kecuali keputusan response contract `null` vs `[]` berubah jadi API-wide contract.

## Execution-ready Worklist / Handoff Contract

start_with: `T1`

```yaml
- id: T1
  action: Jalankan query debug staging untuk emp_id 484 dan 482 pada date 2026-05-22; catat mapping m_salesman, PJP status, destinations_history, outlet_visit_list, query dev.fe dengan JOIN vs LEFT JOIN.
  depends_on: none
  owner: @fixer
  validation: hasil SQL tersimpan di evidence implementasi; minimal mencatat rows/count dan status PJP.
  exit: Root cause data-level terkonfirmasi atau dicatat belum ada akses DB.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T2
  action: Tambah failing regression tests untuk principal detail: Need Review visible, extra_call summary, no-data tetap nil, EmpID passthrough.
  depends_on: T1
  owner: @fixer
  validation: rtk go test ./service/live_monitoring -run TestGetMonitoringDetail_Principal (minimal satu test merah sebelum code fix)
  exit: Test mengekspresikan bug SX-2038.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T3
  action: Tambah/ubah repo principal detail query agar memakai status IN Approved/Need Review dan tidak INNER JOIN ke mst.m_outlet untuk sumber detail.
  depends_on: T2
  owner: @fixer
  validation: rtk go test ./repository/live_monitoring ./service/live_monitoring
  exit: Query compile, tests relevant hijau.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T4
  action: Pakai destinations_history sebagai sumber agregat principal detail atau tambahkan counter extra-call terpisah; map ExtraCall dari hasil query.
  depends_on: T3
  owner: @fixer
  validation: rtk go test ./service/live_monitoring -run TestGetMonitoringDetail_PrincipalIncludesExtraCallSummary
  exit: summary planned/on_going/extra_call/visited/skipped sesuai test.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T5
  action: Refactor status slice/helper dan query builder kecil bila perlu; jangan ubah response contract.
  depends_on: T4
  owner: @fixer
  validation: rtk go test ./service/live_monitoring ./repository/live_monitoring
  exit: Code lebih ringkas, behaviour sama.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T6
  action: Jalankan full test module pjp.
  depends_on: T5
  owner: @fixer
  validation: rtk go test ./...
  exit: Semua test hijau atau failure unrelated terdokumentasi dengan bukti.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T7
  action: Smoke staging endpoint untuk emp_id=484 dan emp_id=482 dengan date=2026-05-22; capture before/after JSON dan request_id.
  depends_on: T6
  owner: @fixer
  validation: cURL endpoint mengembalikan data non-null untuk ID yang punya data; summary benar.
  exit: Evidence response tersimpan; kalau data staging memang tidak ada, hasil DB T1 menjelaskan no-data valid.
  status: ready
  blocker: butuh token staging valid dan akses jaringan
  requires_user_decision: no

- id: T8
  action: Review final dengan @quality-gate untuk tenant filter, query safety, API contract, dan monitoring regression.
  depends_on: T7
  owner: @quality-gate
  validation: checklist QUALITY.md terpenuhi; diff review tidak menemukan blocker.
  exit: Signoff siap PR/deploy.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T9
  action: Jika FE/product meminta `data: []` alih-alih `null`, buat plan/issue terpisah untuk API contract change.
  depends_on: T8
  owner: @orchestrator
  validation: keputusan tertulis FE/product.
  exit: Tidak mencampur contract change ke bugfix SX-2038 tanpa persetujuan.
  status: blocked
  blocker: menunggu keputusan FE/product; default tidak dikerjakan.
  requires_user_decision: yes
```

## Validation Commands

Dari repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

Dari `pjp/`:

```bash
rtk go test ./service/live_monitoring -run TestGetMonitoringDetail
rtk go test ./repository/live_monitoring ./service/live_monitoring
rtk go test ./...
```

Smoke staging:

```bash
curl -G 'https://best.scyllax.online/scylla-pjp/api/v1/monitoring_locations/details' \
  --data-urlencode 'emp_id=484' \
  --data-urlencode 'date=2026-05-22' \
  -H 'Authorization: Bearer <token>' \
  -H 'Accept: application/json, text/plain, */*'

curl -G 'https://best.scyllax.online/scylla-pjp/api/v1/monitoring_locations/details' \
  --data-urlencode 'emp_id=482' \
  --data-urlencode 'date=2026-05-22' \
  -H 'Authorization: Bearer <token>' \
  -H 'Accept: application/json, text/plain, */*'
```

DB verification minimum:

```sql
SELECT emp_id, salesman_id, sales_name, salesman_code, cust_id
FROM mst.m_salesman
WHERE emp_id IN (484, 482);

SELECT id, pjp_code, salesman_id, salesman_code, cust_id, approval_status
FROM pjp_principles.permanent_journey_plans
WHERE salesman_id IN (484, 482)
  AND approval_status IN ('Approved', 'Need Review');

SELECT dh.pjp_id, dh.destination_id, dh.destination_code, dh.destination_type, dh.is_extra_call, dh.cust_id, dh."date"
FROM pjp_principles.destinations_history dh
JOIN pjp_principles.permanent_journey_plans pjp ON pjp.id = dh.pjp_id
WHERE pjp.salesman_id IN (484, 482)
  AND dh."date"::date = '2026-05-22';
```

## Evidence Requirements

- Staging DB result untuk `emp_id=484` dan `482`.
- Before response JSON yang menunjukkan `data:null`.
- After response JSON yang menunjukkan `data:[...]` untuk ID valid.
- Test output `rtk go test ./...` dari `pjp/`.
- Query diff/rationale bila memilih fallback selain `destinations_history`.
- Catatan `No Data` valid bila DB memang tidak punya PJP pada tanggal/emp_id target.
- MCP/source gate: local discovery wajib dan sudah dilakukan. Official docs/context7 tidak diperlukan karena tidak ada library/API baru. GitHub/brave tidak diperlukan karena bug internal. Browser evidence tidak relevan karena BE-only.

## Done Criteria

- Semua worklist `T1`–`T8` selesai atau blocker terdokumentasi.
- Unit/regression tests baru hijau.
- Smoke staging membuktikan data non-null untuk `emp_id` yang benar.
- Summary counts sesuai query DB.
- Tenant scope tidak dilonggarkan.
- `@quality-gate` signoff.
- Open question `data:null` vs `[]` tidak menghalangi bugfix inti; jika ingin ubah contract, follow-up terpisah.

## Final Planning Summary

- Primary plan path: `.opencode/plans/20260522-1101-sx-2038-monitoring-detail-null.md` — source of truth implementation.
- Evidence created/kept: `.opencode/evidence/20260522-1101-sx-2038-monitoring-detail-null/discovery.md` — dipertahankan karena berisi file inspection, smoking gun, root-cause ranking, constraints.
- Draft kept: `.opencode/draft/20260522-1101-sx-2038-monitoring-detail-null/open-questions.md` — dipertahankan karena ada pertanyaan material (`emp_id 484 vs 482`, status, response contract, akses DB).
- Key decisions: tidak tambah lookup `emp_id → salesman_id` dulu; fix status `Approved/Need Review`; pakai `destinations_history`/LEFT JOIN untuk principal detail agar extra-call dan `destination_id NULL` tidak drop semua rows.
- Assumptions: `pjp.salesman_id == m_salesman.emp_id`; status `Need Review` valid untuk detail; response no-data tidak diubah tanpa FE/product.
- Questions asked: belum dikirim sebagai blocking chat; ditulis di draft. Plan tetap implementation-ready dengan asumsi default.
- Cleanup: tidak menghapus evidence/draft karena masih operationally useful untuk implementer dan QA.
- Readiness: siap untuk `@orchestrator` mulai dari `T1`; `T7` butuh token staging valid, `T9` blocked by FE/product decision.
