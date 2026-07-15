# Plan — SX-2445 / SX-2448 / SX-2452 Survey Level Target Distributor (BE)

## Goal

Memperluas service `master` (endpoint `POST/PUT/GET /master/v1/survey`) untuk menerima field `level_target`, `target_cust_id` (cust_id business unit, sesuai DOCX `Enhance_Create_Survey_BE.docx`), dan `target_distributor_id` sesuai tiket SX-2445 (create), SX-2448 (edit), dan SX-2452 (detail), dengan skema baru `mst.m_survey.level_target`, `mst.m_survey_area.target_cust_id`, dan tabel baru `mst.m_survey_distributor` agar target Distributor (Sprint 13) ter-handle secara eksplisit, tetap transactional, dan backward-compatible untuk payload lama.

## Non-goals

- Tidak mengubah kontrak field existing (`survey_title`, `efective_date_*`, `answer_frequency`, `response_type`, `target_type`, `distributor_id`, `area_id`, `outlet_id`, `emp_id`, `survey_template_id`).
- Tidak menyentuh service lain (`sales`, `tms`, `pjp`, dll).
- Tidak menulis ulang test lama; cukup extend test baru dan minor update stub `SurveyRepository` di `master/service/survey_service_test.go` untuk field schema baru.
- Tidak menambah library baru; semuanya pakai `database/sql`, `sqlx`, `go-playground/validator/v10` yang sudah ada di `master/go.mod`.
- Tidak melakukan backfill historis `level_target`/`target_cust_id` untuk baris survey lama; biarkan `null` (lihat Out of Scope).

## Scope

- Modul: `master/`.
- Endpoint: `POST /v1/survey`, `PUT /v1/survey/:survey_id`, `GET /v1/survey/:survey_id` (di-mount `/master/v1/...` lewat `master/main.go:461`).
- Layer target:
  - `master/entity/survey.go` — extend request/response struct + validator rule.
  - `master/model/survey.go` dan `master/model/survey_area.go` — tambah kolom model.
  - `master/model/survey_distributor.go` — model baru.
  - `master/repository/survey_repository.go` — extend Store/Update/FindAreas, tambah distributor mapping.
  - `master/service/survey_service.go` — extend Store/Update/Detail.
  - `master/controller/survey_controller.go` — minimal (body sudah auto-bind JSON).
  - `master/pkg/validation/validation.go` — tambah rule `level_target`.
  - `master/migration/mst.survey/005_add_level_target_and_target_cust_id.sql` — schema baru.
  - `master/migration/mst.survey/006_create_m_survey_distributor.sql` — tabel baru.
  - Test: `master/service/survey_service_test.go` + `master/controller/survey_controller_test.go`.
- DB target: schema `mst`, DB `ggn_scyllax` sesuai baseline compose (lihat `.opencode/docs/ARCHITECTURE.md` & `master/.env`).
- Auth context: gunakan `c.Locals("cust_id")` & `c.Locals("parent_cust_id")` (string tenant) seperti pola existing; payload FE `target_cust_id` (`varchar(10)`, DOCX) terpisah dari `request.CustId` (lihat Non-negotiable Invariant #3).

## Requirements

1. Tambah `level_target` di payload `CreateSurveyBody` & `UpdateSurveyBody` (validator `level_target` di EN/ID, value: `Salesman` | `Outlet` | `Distributor` — disusun dari DOCX `Enhance_Create_Survey_BE.extracted.md:742-748`).
2. Tambah `target_cust_id` payload (`varchar(10)`) untuk cust_id business unit — disusun dari DOCX. Petakan ke:
   - `m_survey_area.target_cust_id` (doc `Create_Survey_Database.extracted.md:116-119`),
   - `m_survey_distributor.cust_id` (doc `Create_Survey_Database.extracted.md:181-184`).
   Validasi cross-check dengan tenant token (lihat Non-negotiable Invariant #3).
3. Tambah `target_distributor_id` array — difilter `0`/negatif sebelum query, lalu di-store ke `m_survey_distributor` bila kombinasi BU × level_target membutuhkannya (lihat §Implementation Steps 3.2).
4. Kolom baru `mst.m_survey.level_target VARCHAR(20) NULL` — ditulis saat create/update; backfill `null` untuk baris lama.
5. Kolom baru `mst.m_survey_area.target_cust_id VARCHAR(10) NULL` — diisi sesuai `target_cust_id` payload FE; untuk BU principal menjadi `cust_id` business unit principal, untuk BU distributor bernilai `NULL` (doc `Enhance_Create_Survey_BE.extracted.md:1322-1334`, `1457-1469`, `1588-1600`).
6. Tabel baru `mst.m_survey_distributor` (sesuai doc `Create_Survey_Database.extracted.md:171-198`) berisi `cust_id`, `survey_id`, `distributor_id` (+ timestamp/audit opsional) dengan unique `(survey_id, distributor_id) WHERE is_del = false`.
7. Response detail (`GET`) expose:
   - `level_target` (string),
   - `business_unit[].target_cust_id` (string) dan `target_cust_name` (string, join `smc.m_customer`),
   - `target_distributor` (array, minimal: `id`/`distributor_id`, `distributor_id`, `distributor_code`, `distributor_name`).
8. Validation rules (disusun dari DOCX body field table):
   - `level_target` wajib (Yes) — payload create/edit harus menyertainya.
   - `target_cust_id` mandatory saat business unit principal, NULL saat distributor.
   - `outlet_id` mandatory saat `level_target = Outlet`.
   - `emp_id` mandatory saat `level_target = Salesman`.
   - `target_distributor_id` mandatory saat `level_target = Distributor`.
   - Filter `0`/negatif dari `target_distributor_id` sebelum insert/lookup.
   - `survey_id` harus ada; rule status existing tetap (lihat Assumption A3).
9. `POST` / `PUT` harus idempotent untuk `m_survey_distributor`: delete-then-insert di-update flow.
10. Coverage test:
    - 6 test kombinasi (`salesman|outlet|distributor` × `distributor|principal`) untuk create,
    - 1 test edit round-trip (salesman → outlet) yang assert tidak ada duplikat row,
    - 1 test detail shape (`level_target`, `target_cust_id`, `target_distributor`),
    - 2 test controller create/update yang validasi parsing field baru.

## Acceptance Criteria

Lihat `## Acceptance Criteria` di bawah untuk 8 case (A–G dan H). Setiap case punya expected: HTTP success (atau error untuk test validasi), rows `m_survey_area` & `m_survey_distributor`, dan field response detail.

## Existing Patterns/Reuse

- `master/service/survey_service.go:147-176` — `normalizeBusinessUnitSelection`, `normalizeUniqueInts`, `normalizeBusinessUnitIds` bisa dipakai ulang; tambah helper `normalizeTargetDistributorIds` untuk filter `0`/negatif dengan aturan yang sama.
- `master/service/survey_service.go:442-525` (Store) dan `:528-642` (Update) — pola transaction + replace-on-update `m_survey_area`; replikasi pola ini untuk `m_survey_distributor`.
- `master/service/survey_service.go:652-687` (`buildSurveyAreas`) — extend agar menghasilkan `model.SurveyArea` dengan `TargetCustId` terisi.
- `master/service/survey_service.go:234-422` (`Detail`) — extend di blok `for _, a := range areas` untuk set `target_cust_id` di response business unit / area.
- `master/repository/survey_repository.go:276-303` (`StoreAreas`, `DeleteAreasBySurveyId`, `FindAreasBySurveyId`) — replikasi untuk distributor.
- `master/repository/survey_repository.go:122-156` (`FindAllByCustId`, `FindOneById`) — tambahkan `level_target` & `target_cust_id` di select list.
- `master/pkg/validation/validation.go:59` — tambahkan `vc.RegisterValidation("level_target", levelTarget)` dan translations.
- Test patterns:
  - `master/service/survey_service_test.go:16-180` — `surveyRepositoryRedStub` + `salesmanRepositoryStub`; tambah field `storeDistributorInput`, `deleteDistributorCalled`, `findDistributorsResult`, dst.
  - `master/controller/survey_controller_test.go:162-256` — pola parsing body dengan field tambahan (lihat `TestSurveyController_Create_ShouldParseDistributorAndEmpArrays`).
- Migration patterns: `master/migration/mst.survey/001_create_tables.sql`, `002_add_distributor_and_salesman.sql` — pakai `IF NOT EXISTS`, `ADD COLUMN IF NOT EXISTS`, backfill aman.

## Constraints

- Ikuti flow Controller → Service → Repository → DB.
- Write dalam transaction service-layer (`txManager.WithinTransaction`).
- Jangan hardcode `distributor_id=120`, `area_id=82`, dst dari tiket; semua resolve via repository.
- Tidak menambah dependensi baru di `go.mod`/`go.sum`.
- Multi-tenant: scope `smc.m_customer.cust_id` (string) tetap dipakai untuk resolve `distributor_id → cust_id` saat `target_cust_id` payload kosong; `target_cust_id` divalidasi berada di tenant tree (lihat Non-negotiable Invariant #3).
- Migration idempotent: aman dijalankan ulang; tidak menghapus data.
- Migration `m_survey_area.distributor_id` di staging mungkin tidak punya FK (lihat Assumption A1). Jangan tambahkan FK baru.

## Risks

- **Schema drift**: staging bisa punya `m_survey_area.distributor_id` dengan FK ke `m_distributor(distributor_id)`. INSERT `0` ditolak. Mitigasi: lihat Non-negotiable Invariant #2; Plan §Verification gate untuk cek sebelum deploy.
- **Idempotency replace**: pattern existing `DeleteAreasBySurveyId` (set `is_del=true`) + `StoreAreas` harus konsisten untuk `m_survey_distributor`. Mitigasi: tambahkan helper di repository yang menulis `is_del=true` lalu `INSERT`; test edit round-trip memverifikasi count.
- **Detail performance**: tambah `FindSurveyDistributorsBySurveyId` menambah 1 query di detail; tetap oke untuk volume existing. Mitigasi: jika N row > 100k butuh agregasi, parking.
- **FE contract drift**: nama field response (`target_cust_id` di level survey vs per-area) masih open per tiket SX-2452. Mitigasi: expose `level_target`, `target_cust_id`, dan `target_distributor` di top-level response, sehingga FE bisa memilih representasi. Field per-area (`target_cust_id` di `business_units` atau `target_survey.area`) juga ditambahkan agar backward compatibility aman.
- **Backward compatibility**: payload existing tidak punya `level_target`/`target_distributor_id`. Service harus accept payload tanpa field ini (default fallback ke BU existing, no insert distributor baru). Mitigasi: lihat Assumption A2.

## Decisions/Assumptions

### Decisions

- **D1**: Pakai migration terpisah (file `005_…` dan `006_…`) daripada satu file besar agar rollback dan review lebih mudah.
- **D2**: Tabel `m_survey_distributor` punya `cust_id` (string tenant, `varchar(10)`) sesuai konvensi `m_survey_salesman.cust_id`; sesuai doc BE/Database, field `cust_id` adalah cust_id business unit yang akan ditampilkan di mobile sebagai task (`Create_Survey_Database.extracted.md:117-119`). Field `target_cust_id` di `m_survey_area` diseragamkan `varchar(10)`.
- **D3**: Validator `level_target` EN/ID didaftarkan setelah `answer_frequency` agar reuse pola (lihat `master/pkg/validation/validation.go:59`). Rule DOCX: wajib diisi `Outlet | Distributor | Salesman`.
- **D4**: Field body FE adalah `target_cust_id` (varchar(10)) sesuai DOCX BE. Ticket summary menyebut `cust_id`, tetapi DOCX lebih otoritatif untuk payload body. Implementer mengikuti DOCX; jika FE sudah terlanjut mengirim `cust_id` di payload (alias), mapper sementara dapat dipakai tapi tidak menjadi default. Confirm ke FE (Widya) sebelum eksekusi.
- **D5**: Untuk BU principal, `m_survey_area` row di-insert dengan `distributor_id = NULL`, `area_id = NULL`, dan `target_cust_id` = cust_id principal (DOCX line 1322-1334, 1457-1469, 1588-1600). Ini **berbeda** dengan pola existing service yang menyimpan `distributor_id = 0` (sentinel) untuk principal area (lihat Master §Conflict dengan pola existing). Untuk backward compat, implementer boleh menyimpan kedua representasi (`0` dan NULL) bila query existing masih bergantung; default rekomendasi: ikuti DOCX (NULL), dan update semua query/service yang membaca `m_survey_area.distributor_id` agar `0` (sentinel lama) dipetakan ke NULL.
- **D6**: Detail response expose `level_target`, dan `target_distributor` di top-level `SurveyDetailResponse`; field `target_cust_id`/`target_cust_name` di-embed ke `business_unit[]` (lihat DOCX `Enhance_Create_Survey_BE.extracted.md:273-296`, `412-509`, `510-612`, `613-709`). Field existing `business_units`/`salesman`/`outlet` tetap dipertahankan.
- **D7**: `target_distributor_id` dari FE di-trim (drop `0` dan negatif) lalu di-store ke `m_survey_distributor` bila `level_target = Distributor` (DOCX). Untuk `level_target = Salesman | Outlet`, DOCX tidak menunjukkan insert ke `m_survey_distributor` (lihat DOCX impact DB), sehingga implementer TIDAK insert ke `m_survey_distributor` di level non-Distributor. Pengecualian: bila FE mengirim `target_distributor_id` non-empty saat `level_target` ≠ Distributor, treat sebagai no-op (abaikan, jangan error) — ini sesuai spec DOCX bahwa field itu hanya mandatory di level Distributor.

### Assumptions

- **A1**: Schema lokal `mst.m_survey_area.distributor_id` TIDAK memiliki FK ke `m_distributor(distributor_id)`. Jika staging punya FK, INSERT sentinel `0` (pola lama repo) akan gagal. Jika worker memilih representasi `NULL` sesuai DOCX, risiko FK menurun tetapi query existing yang mengandalkan sentinel `0` harus direview.
- **A2**: Untuk payload legacy yang tidak mengirim `level_target`, service tetap lanjut dengan `level_target = ""` (null di DB) demi backward compatibility. Namun untuk payload enhance baru, `level_target` wajib diisi (DOCX body table menyatakan `Yes`).
- **A3**: Aturan edit existing (status survey) tidak berubah; survey yang sudah di-deactivate tetap tidak bisa diedit. Implementasi baru mengikuti pola existing `surveyServiceImpl.Update`.
- **A4**: Field body FE yang benar adalah `target_cust_id` (`varchar(10)`), bukan `cust_id` integer. Ticket summary user memakai nama singkat `cust_id`; plan memihak DOCX. Jika FE production masih mengirim `cust_id`, worker boleh menambah alias parser sementara dengan evidence dan catatan deprecation.
- **A5**: `target_distributor_id` hanya bermakna untuk `level_target = Distributor` sesuai DOCX. Untuk `Salesman`/`Outlet`, jika field ini ikut terkirim maka abaikan tanpa error dan jangan insert `m_survey_distributor`.
- **A6**: Migration dijalankan manual (`rtk make migrateUp` tidak tersedia di `master`; `master/migration/mst.survey/` adalah folder SQL). Operator menjalankan file `005` lalu `006` lewat CLI/SQL client.
- **A7**: FE koordinasikan nama final field response `target_distributor` vs `distributor` (DOCX daftar atribut menyebut `target_distributor`, contoh JSON level Distributor menulis `distributor`). Plan memihak `target_distributor` sebagai nama contract baru, sambil menjaga backward compatibility bila perlu alias sementara.

## Source Strategy (used vs skipped)

- **Used**:
  - Repo lokal: `master/entity/survey.go`, `master/service/survey_service.go`, `master/repository/survey_repository.go`, `master/controller/survey_controller.go`, `master/migration/mst.survey/*`, `master/service/survey_service_test.go`, `master/controller/survey_controller_test.go`, `master/pkg/validation/validation.go`, `master/main.go`, `master/.env`, `.opencode/docs/{ARCHITECTURE,SERVICE_MATRIX,QUALITY,SECURITY,PROMPT_GATES}.md`, plan lama `20260504-1034-sx-1789-…`, `20260504-0846-sx-1906-…`, `20260504-2058-sx-1915-…`.
  - DOCX lokal yang diekstrak ke evidence:
    - `.opencode/evidence/20260707-sx-2445-2448-2452-survey-level-target-distributor/Enhance_Create_Survey_BE.extracted.md`
    - `.opencode/evidence/20260707-sx-2445-2448-2452-survey-level-target-distributor/Create_Survey_Database.extracted.md`
  - Existing migration numbering (`001`–`004`) dan pattern `IF NOT EXISTS`.
- **Skipped with reason**:
  - **Context7 / GitHub / web search**: tidak perlu, tidak ada library/API baru; semua pattern ada di repo lokal + DOCX lokal.
  - **Browser/visual**: tidak relevan, backend-only.
  - **Library docs (Fiber, sqlx, validator)**: pola existing service sudah terdokumentasi; tidak ambil risiko rekomendasi yang berbeda dari implementasi.
  - **Google Doc remote link**: tidak perlu fetch karena dokumen `.docx` yang sama tersedia lokal di repo dan sudah diekstrak ke evidence.

## Source Anatomy

- **Auth scope & tenant string**:
  - `master/controller/survey_controller.go:139-145` (Create) dan `:208-214` (Update) — set `request.CustId` (string), `request.ParentCustId`, dan `request.CreatedBy/UpdatedBy` dari locals.
  - `master/pkg/middleware/*.go` (sudah dikonfirmasi di `master/main.go`) — `JWTProtected` mengisi locals.
- **Repository contract**:
  - `master/repository/survey_repository.go:14-40` (interface) — pastikan extend.
  - `master/service/survey_service_test.go:16-180` (`surveyRepositoryRedStub`) — implementer harus menambah stub method untuk distributor.
- **Transaction layer**:
  - `master/service/survey_service.go:487-525` (Store) dan `:592-641` (Update) — `txManager.WithinTransaction(ctx, func(txCtx) { tx := repository.GetTxFromContext(txCtx); ... })`.
  - `master/repository/transaction_repository.go` (sudah ada) — `GetTxFromContext`.
- **Migration folder**:
  - `master/migration/mst.survey/001_create_tables.sql:27-33` — pola CREATE TABLE IF NOT EXISTS.
  - `master/migration/mst.survey/002_add_distributor_and_salesman.sql:1-9` — pola `ADD COLUMN IF NOT EXISTS` + backfill.
- **Validator registry**:
  - `master/pkg/validation/validation.go:59, 107-117, 157-167` — pola `RegisterValidation` + translation EN/ID.
- **Existing test layout**:
  - `master/service/survey_service_test.go:313-528` — pola `surveyRepositoryRedStub` + test case nama `TestSurveyService_Store_…`.
  - `master/controller/survey_controller_test.go:116-260` — pola `httptest.NewRequest("POST", "/v1/survey", strings.NewReader(body))` + `app.Test`.

## Reference Map

- **SX-2445 (create)**: repo-backed — replikasi pola `Store` di `master/service/survey_service.go:424-526`; extend payload `entity.CreateSurveyBody` (`master/entity/survey.go:123-138`); extend SQL `Store` & `StoreAreas` (`master/repository/survey_repository.go:205-217, 276-285`).
- **SX-2448 (update)**: repo-backed — replikasi pola `Update` di `master/service/survey_service.go:528-642`; replace-on-update `m_survey_distributor` mengikuti pola `DeleteAreasBySurveyId` + `StoreAreas` di `master/service/survey_service.go:600-610`.
- **SX-2452 (detail)**: repo-backed — extend `Detail` di `master/service/survey_service.go:234-422`; tambah query distributor list di repository `master/repository/survey_repository.go:293-303`.
- **Migration**: docs-backed + repo-backed — pola `001_create_tables.sql` & `002_add_distributor_and_salesman.sql` (lokal); Google Doc link di tiket hanya sebagai referensi final nama kolom; plan mengikuti rekomendasi user di prompt (lihat kolom yang akan ditambah).
- **Test**: repo-backed — replikasi test `TestSurveyService_Store_ShouldPersistPrincipalBusinessUnitAndDistributorMappings` di `master/service/survey_service_test.go:526-572`; pola JSON parsing di `TestSurveyController_Create_ShouldParseDistributorAndEmpArrays` di `master/controller/survey_controller_test.go:162-213`.

## Execution Source of Truth

Precedence (highest to lowest):

1. Latest explicit user instruction in SX-2445/2448/2452 ticket context.
2. `SECURITY.md` dan `AGENTS.md` repo-local (tenant isolation, env key tidak dicommit).
3. Non-negotiable Implementation Invariants (§ bawah).
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria dan Done Criteria.
6. Implementation Steps.
7. TDD/Test Plan.

Konflik dicatat di evidence/verification (lihat Evidence Requirements §1).

## Non-negotiable Implementation Invariants

1. **Planner tidak mengedit source**. Eksekusi di bawah `@fixer`/orchestrator lane, bukan planner.
2. **Schema gate**: sebelum deploy, operator harus memverifikasi staging `mst.m_survey_area.distributor_id` tidak punya FK ke `m_distributor`. Jika ada, implementasi perlu patch tambahan (tidak di slice ini).
3. **Tenant scope**: `target_cust_id` payload FE (`varchar(10)`) divalidasi tipe (non-empty string) dan diserialisasi; service TIDAK mempercayainya sebagai ganti `c.Locals("cust_id")` (string tenant). Cross-check penuh dengan tree tenant (parent_cust_id) di luar slice ini (lihat A4).
4. **Write in transaction**: replace `m_survey_area` + replace `m_survey_distributor` + replace `m_survey_salesman` selalu dalam satu `txManager.WithinTransaction`. Tidak ada partial commit.
5. **Filter `0`/negatif**: `target_distributor_id` di-filter SEBELUM query DB; `distributor_id=0` existing sentinel tetap dipakai untuk `m_survey_area` jika ada di payload (lihat pola `normalizeBusinessUnitSelection`).
6. **Backward compatibility response**: tidak ada field lama di `SurveyDetailResponse` yang dihapus atau diganti nama.
7. **Validator wajib**: `level_target` harus lewat validator EN/ID; jika tidak dikirim pada create, `LevelTarget` di header disimpan `""` (null) dan tidak error (lihat A2) — tapi TIDAK boleh menerima nilai di luar `Salesman`/`Outlet`/`Distributor` jika dikirim.
8. **Idempotency replace**: delete-then-insert `m_survey_distributor` di-update flow; tidak insert duplicate `(survey_id, distributor_id)`.
9. **Test wajib Red→Green**: setiap perubahan Store/Update/Detail diikuti test gagal dulu lalu pass; coverage 6 kombinasi BU × level_target.
10. **Planner handoff**: eksekusi dilakukan di lane berikutnya; planner read-only tidak berlaku di execution lane.

## Do Not / Reject If

- **Do not** insert `m_survey_distributor` dengan `distributor_id <= 0`. Filter sebelum insert.
- **Do not** insert `m_survey_area` baru tanpa `target_cust_id` (kecuali baris yang sudah ada pre-migration, biarkan `null`).
- **Do not** ubah `m_survey.cust_id` (string tenant) — tetap diisi dari token.
- **Do not** return error hanya karena `len(area_id) != len(distributor_id)` (sesuai carry-over SX-1789).
- **Do not** rename `distributor_id`/`area_id`/`outlet_id`/`emp_id`/`survey_template_id` di payload atau response.
- **Do not** gunakan library baru di `go.mod`.
- **Do not** hardcode `distributor_id=120`, `area_id=82`, atau `survey_template_id=59` di service/repo.
- **Reject if** QA report menunjukkan duplicate row `m_survey_area` atau `m_survey_distributor` setelah edit round-trip.
- **Reject if** `level_target` dikirim dengan nilai di luar `Salesman`/`Outlet`/`Distributor` dan service menerima tanpa error.
- **Reject if** migration membuat `m_survey_area.target_cust_id` NOT NULL tanpa backfill aman.
- **Reject if** response detail `GET` menghilangkan field lama (e.g. `business_units`).
- **Reject if** test coverage 6 kombinasi BU × level_target tidak lengkap.

## Diff Boundary

- **Allowed files**:
  - `master/entity/survey.go`, `master/model/survey.go`, `master/model/survey_area.go`, `master/model/survey_distributor.go` (new), `master/repository/survey_repository.go`, `master/service/survey_service.go`, `master/controller/survey_controller.go` (jika perlu), `master/pkg/validation/validation.go`, `master/service/survey_service_test.go`, `master/controller/survey_controller_test.go`, `master/migration/mst.survey/005_add_level_target_and_target_cust_id.sql` (new), `master/migration/mst.survey/006_create_m_survey_distributor.sql` (new).
  - `.opencode/draft/<task-id>/migration.md` (operator notes, optional; tidak masuk commit).
- **Forbidden**:
  - `go.mod` / `go.sum` (kecuali `go mod tidy` yang diizinkan dengan justifikasi eksplisit di evidence).
  - File di service lain (`sales/`, `tms/`, `pjp/`, dst).
  - File konfigurasi runtime (`docker-compose.yml`, `master/.env`).
  - File dokumentasi di luar `.opencode/`.
  - Script `scripts/` di root repo.
- **Out-of-bound change remediations**: revert via `git checkout`; catat di evidence verification.

## TDD/Test Plan

TDD required: **ya**. Mengubah 3 endpoint + 1 tabel baru + 1 kolom baru + 6 kombinasi behavior = production risk tinggi.

### Existing test patterns
- `master/service/survey_service_test.go:16-180` — `surveyRepositoryRedStub` implements `repository.SurveyRepository`. Tambah field:
  - `storeDistributorInput []model.SurveyDistributor`
  - `deleteDistributorCalled bool`
  - `findDistributorsResult []model.SurveyDistributor`
  - `updateLevelTargetCalled bool`
  - `updateTargetCustIdInAreasCalled bool`
- `master/controller/survey_controller_test.go:1-580` — `surveyServiceControllerStub`; extend untuk menerima field baru.

### First failing tests (Red)

1. `TestSurveyService_Store_ShouldInsertLevelTarget_AndTargetCustId_OnSurveyHeader` — payload `level_target="Salesman"`, `cust_id=120`. Expect `Store` (header) menerima `LevelTarget="Salesman"`, `TargetCustId=120`. Saat ini stub `Store` tidak menyimpan field ini; test akan gagal sebelum extension.
2. `TestSurveyService_Store_ShouldInsertTargetDistributorsForDistributorLevel_AndMarkTargetCustIdOnArea` — payload `level_target="Distributor"`, `target_distributor_id=[120]`, `cust_id=120`. Expect `storeDistributorInput` berisi 1 row `(survey_id, 120, cust_id="120")` dan `storeAreasInput` rows berisi `TargetCustId=120`.
3. `TestSurveyService_Store_ShouldFilterZeroAndNegativeTargetDistributorId_BeforeInsert` — payload `target_distributor_id=[0, 120, -1, 120]`. Expect hanya 1 row distributor.
4. `TestSurveyService_Store_ShouldRejectInvalidLevelTarget_OnlyForKnownValues` — payload `level_target="Manager"`. Expect validator/service return error. (Validator EN/ID akan diuji terpisah di test validator.)
5. `TestSurveyService_Store_ShouldNotInsertSurveyDistributor_WhenLevelTargetIsSalesmanAndPrincipalOnly` — payload `level_target="Salesman"`, `distributor_id=[0]`, `target_distributor_id=[]`. Expect `storeDistributorInput` kosong (principal-only case).
6. `TestSurveyService_Store_ShouldPreserveTargetDistributorForDistributorBu_WhenLevelTargetIsOutletOrSalesman` — payload `level_target="Outlet"`, `distributor_id=[120]`, `target_distributor_id=[120]`. Expect `storeDistributorInput` 1 row, `storeAreasInput` ≥ 1 row.
7. `TestSurveyService_Update_ShouldReplaceSurveyDistributors_NoDuplicates` — call Update 2× dengan payload berbeda, assert delete-then-insert dan tidak ada row duplicate (idempotency).
8. `TestSurveyService_Detail_ShouldExposeLevelTarget_TargetCustId_AndTargetDistributorList` — stub `findDetailSurvey{LevelTarget: "Outlet"}`, `findDistributorsResult=[{DistributorId:120, CustId:"120", DistributorName:"PT Makmur"}]`. Expect response top-level memuat `level_target`, `target_cust_id`, dan `target_distributor[0].distributor_id=120`.
9. `TestSurveyService_Store_ShouldMaintainPrincipalAndMultiDistributorBackwardCompatibility` — replikasi test SX-1789/SX-1906 dengan field baru `target_distributor_id=[]`; expect tidak ada regression.
10. `TestSurveyController_Create_ShouldParseLevelTargetAndTargetDistributor` — JSON body berisi `level_target`, `cust_id`, `target_distributor_id`. Expect `serviceStub.lastCreate` memuat field-field ini.

### Green step
- Extend entity, model, repository, service, controller, validator secara inkremental.
- Tambah file migration `005` & `006`.
- Implementasi repository `StoreSurveyDistributors`, `DeleteSurveyDistributorsBySurveyId`, `FindSurveyDistributorsBySurveyId`.
- Implementasi service normalisasi `level_target` (trim/lower) + validasi inclusion set; integrasi ke `Store`/`Update`/`Detail`.

### Refactor step
- Ekstrak helper `resolveTargetDistributorScope` di `master/service/survey_service.go` agar `Store`/`Update` DRY.
- Ekstrak helper `mapSurveyAreaTargetCustId` agar `buildSurveyAreas` dan `Detail` konsisten.

### Edge cases
- `level_target=""` (tidak dikirim) → default null, behavior existing.
- `cust_id=0` (tidak dikirim) → `target_cust_id` di area = `null` (kecuali lookup distributor memberi cust string).
- `target_distributor_id=[0]` → difilter; no insert.
- `target_distributor_id=[]` + `level_target="Distributor"` → no row di `m_survey_distributor` (valid edge; FE mengirim tanpa target aktual).
- Edit: `level_target` berubah → delete-then-insert distributor mapping.

### Commands (validasi per service)

```bash
# dari repo root
rtk docker compose -f docker-compose.yml ps

# dari master/
cd master
rtk go mod download && rtk go mod tidy
rtk go test ./service -run 'TestSurveyService_(Store|Update|Detail)'
rtk go test ./controller -run 'TestSurveyController_(Create|Update|Detail)'
rtk go test ./repository -run TestSurvey
rtk go test ./...
```

## Implementation Steps

1. **Migration schema** (file baru):
   1.1. `master/migration/mst.survey/005_add_level_target_and_target_cust_id.sql`:
   - `ALTER TABLE mst.m_survey ADD COLUMN IF NOT EXISTS level_target VARCHAR(20) NULL;`
   - `ALTER TABLE mst.m_survey_area ADD COLUMN IF NOT EXISTS target_cust_id VARCHAR(10) NULL;`
   - `CREATE INDEX IF NOT EXISTS idx_m_survey_level_target ON mst.m_survey(level_target);`
   1.2. `master/migration/mst.survey/006_create_m_survey_distributor.sql`:
   - `CREATE TABLE IF NOT EXISTS mst.m_survey_distributor (... seperti requirements §6 ...)`.
   - `CREATE UNIQUE INDEX IF NOT EXISTS uniq_m_survey_distributor_survey_distributor ON mst.m_survey_distributor(survey_id, distributor_id) WHERE is_del = false;` (Postgres partial unique index).
2. **Model**:
   2.1. `master/model/survey.go` — tambah field `LevelTarget *string`, `TargetCustId *string` (header; opsional karena backfill null). DB tags mengikuti `db:"level_target"`, `db:"target_cust_id"`.
   2.2. `master/model/survey_area.go` — tambah `TargetCustId *string` `db:"target_cust_id"`.
   2.3. `master/model/survey_distributor.go` (baru) — struct `SurveyDistributor`:
   ```go
   type SurveyDistributor struct {
       SurveyDistributorId int     `db:"survey_distributor_id"`
       SurveyId            int     `db:"survey_id"`
       DistributorId       int     `db:"distributor_id"`
       CustId              string  `db:"cust_id"`
       IsDel               bool    `db:"is_del"`
       CreatedAt           *time.Time `db:"created_at"`
       CreatedBy           *int64     `db:"created_by"`
       UpdatedAt           *time.Time `db:"updated_at"`
       UpdatedBy           *int64     `db:"updated_by"`
       // joined
       DistributorName *string `db:"distributor_name" json:"distributor_name,omitempty"`
   }
   ```
3. **Entity** (`master/entity/survey.go`):
   3.1. `CreateSurveyBody` & `UpdateSurveyBody` (extend, bukan rename field existing):
   - `LevelTarget string` dengan json tag `"level_target"` dan validator `validate:"required,level_target"` (DOCX menyatakan field ini required pada create & edit).
   - `TargetCustId string` dengan json tag `"target_cust_id"` (varchar 10 sesuai DOCX). Validator: minimal length 1 ketika payload tidak null; tidak `required` agar legacy caller tanpa field ini tetap diterima (lihat A2).
   - `TargetDistributorId FlexibleIntArray` dengan json tag `"target_distributor_id"`.
   - `CustId string` (string tenant dari token, existing — tidak diubah) dan `ParentCustId string`, `CreatedBy int64`/`UpdatedBy int64` tetap.
   3.2. `SurveyDetailResponse` (extend):
   - `LevelTarget string` dengan json tag `"level_target"`.
   - `TargetDistributor []SurveyDistributorResponse` dengan json tag `"target_distributor"`.
   - `BusinessUnits []SurveyBusinessUnit` sudah ada; extend `SurveyBusinessUnit` dengan:
     - `TargetCustId string` (`json:"target_cust_id,omitempty"`).
     - `TargetCustName string` (`json:"target_cust_name,omitempty"`) — join ke `smc.m_customer.cust_name` saat `target_cust_id` non-null.
   3.3. `SurveyDistributorResponse` (struct baru): `MSurveyDistributorId int` `json:"id"`, `DistributorId int` `json:"distributor_id"`, `DistributorCode string` `json:"distributor_code,omitempty"`, `DistributorName string` `json:"distributor_name"`. Catatan: DOCX contoh JSON level Distributor menulis field `distributor_code`/`distributor_name`, bukan `cust_id`. Implementasi mengikuti join ke `mst.m_distributor` (sudah ada di `m_survey_area` repo).
   3.4. Validator tag: `validate:"required,level_target"` di Create/Update untuk `LevelTarget`. Tambah rule `level_target` di `master/pkg/validation/validation.go:59` (fungsi `levelTarget(fl validator.FieldLevel) bool` mengembalikan true jika `Salesman|Outlet|Distributor`). Tambah translation EN/ID di blok yang sama dengan `answer_frequency`.
4. **Repository** (`master/repository/survey_repository.go`):
   4.1. Extend interface `SurveyRepository`:
   - `Store(tx, survey)` — tambahkan parameter `levelTarget` dan set di SQL `INSERT` (kolom `level_target`). `model.Survey` di-extend dengan `LevelTarget *string` (DB tag `db:"level_target"`). `target_cust_id` TIDAK disimpan di header `m_survey` (hanya di `m_survey_area`).
   - `Update(tx, surveyId, custId, survey)` — tambahkan `level_target = $X` di SQL.
   - `StoreAreas(tx, areas)` — extend dengan `model.SurveyArea.TargetCustId *string` (DB tag `db:"target_cust_id"`). SQL include `target_cust_id` (sesuai D5, representasi principal: `distributor_id=0` & `area_id=0` dan `target_cust_id=cust_principal` — atau NULL, lihat Conflict note di bawah).
   - `FindAreasBySurveyId` — select `target_cust_id` + left join `smc.m_customer` untuk `cust_name` (target_cust_name).
   - `FindOneById` — tambahkan `level_target` di select.
   - Tambah:
     - `StoreSurveyDistributors(tx *sqlx.Tx, distributors []model.SurveyDistributor) error`
     - `DeleteSurveyDistributorsBySurveyId(tx *sqlx.Tx, surveyId int) error`
     - `FindSurveyDistributorsBySurveyId(surveyId int) ([]model.SurveyDistributor, error)` — left join `mst.m_distributor` untuk `distributor_code`/`distributor_name`.
   4.2. Implementasi: pola `StoreAreas` (loop `tx.Exec`) + `DeleteAreasBySurveyId` (set `is_del=true`) + `FindAreasBySurveyId` (LEFT JOIN).
   4.3. **Conflict dengan pola existing (kritikal)**: pola existing `m_survey_area` untuk principal menggunakan `distributor_id = 0` sentinel. DOCX ingin `distributor_id = NULL` + `target_cust_id` non-null. Keputusan: implementer tetap menulis row di `m_survey_area` untuk BU principal dengan `distributor_id = 0` (backward compat) dan `target_cust_id` = cust_id principal. TIDAK menulis `area_id`/insert distributor sentinel ke `m_survey_distributor`. Untuk BU distributor, `target_cust_id` di area = `NULL`; distribusi area mengikuti pola existing (`FindSurveyAreasByDistributorIds`). Flag: `find_areas` & response tetap harus expose principal via `business_unit.type = "principal"` (seperti pola SX-1789).
5. **Service** (`master/service/survey_service.go`):
   5.1. Tambah helper:
   - `normalizeTargetDistributorIds(values []int) []int` — drop `0` dan negatif, dedupe. (Sesuai D7/Asumption A5: hanya digunakan bila `level_target = Distributor`.)
   - `isPrincipalScope(targetCustId string, distributorIdZeroPresent bool) bool` — deteksi BU principal: `target_cust_id` non-empty ATAU `distributor_id` hanya berisi `0`.
   - `buildSurveyDistributors(surveyId int, targetDistIds []int, custIdPayload string) []model.SurveyDistributor` — hanya dipanggil saat `level_target = Distributor`; setiap row berisi `survey_id`, `distributor_id`, `cust_id` (cust header FE atau default request.CustId).
   5.2. `Store`:
   - Validasi `level_target` (service-side guard meskipun validator sudah ada, untuk defensive; return `ErrSurveyInvalidLevelTarget` bila di luar set).
   - `buildSurveyAreas` di-extend:
     - Jika principal scope: row tunggal dengan `DistributorId = 0` (sentinel) + `AreaId = 0` (atau `NULL` jika worker memilih representasi DOCX — konfirmasi dengan operator) + `TargetCustId = request.TargetCustId`.
     - Jika distributor scope: row mengikuti pola existing `FindSurveyAreasByDistributorIds` + `payloadAreaIds`, `TargetCustId = NULL`.
   - `buildSurveyDistributors` dipanggil hanya jika `request.LevelTarget == "Distributor"`. Insert ke `m_survey_distributor` di dalam `txManager.WithinTransaction`.
   - `m_survey_distributor.cust_id` = `request.CustId` (string tenant dari token) — lihat catatan D2: ini berbeda dengan `target_cust_id` payload. `cust_id` di tabel distributor merepresentasikan owner header survey (yang akan dipakai mobile untuk task assignment).
   - Bila `target_distributor_id` kosong dan `level_target="Distributor"` → no row di `m_survey_distributor`, no error (DOCX: field mandatory secara spec, tapi implementasi defensif: empty list diizinkan).
   5.3. `Update`:
   - Sama dengan `Store`; tambah `DeleteSurveyDistributorsBySurveyId(tx, surveyId)` sebelum `StoreSurveyDistributors` untuk idempotency.
   - `DeleteAreasBySurveyId` tetap dilakukan (replace).
   5.4. `Detail`:
   - Set `LevelTarget` di top-level response.
   - Tambah `FindSurveyDistributorsBySurveyId` call; map ke `[]SurveyDistributorResponse` (pakai `MSurveyDistributorId` → `id` di response sesuai DOCX).
   - Set `TargetCustId` per `business_unit` (extend entity struct) — dari `m_survey_area.target_cust_id`. Tambahkan join ke `smc.m_customer` untuk `target_cust_name`.
   - Pertahankan `target_survey.area`/`outlet`/`salesman` agar tidak hilang (DOCX response list masih menyebut `target_survey`).
6. **Controller** (`master/controller/survey_controller.go`):
   - Tidak ada perubahan routing. Body parsing via `json.Unmarshal` otomatis handle field baru (lihat pola `Create`/`Update` existing).
   - Tambah branch `errors.Is(err, service.ErrSurveyInvalidLevelTarget)` untuk return 400.
   - Tambah branch `errors.Is(err, service.ErrSurveyLevelTargetRequired)` (jika ada custom error baru) untuk return 400.
7. **Test**:
   7.1. `master/service/survey_service_test.go`:
   - Extend `surveyRepositoryRedStub` dengan field distributor (`storeDistributorInput []model.SurveyDistributor`, `deleteDistributorCalled bool`, `findDistributorsResult []model.SurveyDistributor`, dst.).
   - Tambah 8 test di Red step (lihat TDD/Test Plan).
   7.2. `master/controller/survey_controller_test.go`:
   - Extend `surveyServiceControllerStub` agar menerima field payload baru.
   - Tambah 2 test: parse body create + parse body update dengan field baru.
8. **Validator** (`master/pkg/validation/validation.go`):
   - Tambah `vc.RegisterValidation("level_target", levelTarget)` dan translation EN/ID.
   - Fungsi: `func levelTarget(fl validator.FieldLevel) bool` — string harus `Salesman`/`Outlet`/`Distributor` (case-sensitive sesuai input FE per DOCX).
9. **Migration apply** (operator):
   - Jalankan `005_…` lalu `006_…` di staging.
   - Validasi via query `information_schema` (lihat Evidence Requirements §2).

## Acceptance Criteria

Lihat 8 test case (A–G + H) di bawah. Setiap case adalah test berdiri sendiri di `master/service/survey_service_test.go` dan/atau `master/controller/survey_controller_test.go`.

### Case A — `Salesman` + `Distributor` (level_target=Salesman, BU=Distributor)

Request body:

```json
{
  "survey_title": "t1 salesman-dist",
  "efective_date_start": "2026-07-10",
  "efective_date_end":   "2026-07-20",
  "answer_frequency": "One Time",
  "response_type": "Mandatory",
  "target_type": "Specific",
  "distributor_id": [120],
  "area_id": [82],
  "outlet_id": [],
  "emp_id": [370],
  "survey_template_id": 59,
  "level_target": "Salesman",
  "target_cust_id": null,
  "target_distributor_id": []
}
```

Expected:
- `model.Survey.LevelTarget == "Salesman"`.
- `storeAreasInput` berisi row distributor-scope (`DistributorId=120`, `TargetCustId=nil`).
- `storeSalesmenInput` berisi salesman yang dipilih.
- `storeDistributorInput` **kosong** (DOCX impact DB untuk Salesman+Distributor tidak menulis `m_survey_distributor`).

### Case B — `Outlet` + `Distributor`

Payload sama, `level_target: "Outlet"`, `outlet_id: [1234]`, `emp_id: []`. Expected: insert `m_survey_outlet` + row `m_survey_area` distributor-scope (`TargetCustId=nil`) + `storeDistributorInput` kosong.

### Case C — `Distributor` + `Distributor`

Payload: `level_target: "Distributor"`, `target_distributor_id: [120]`, tanpa `outlet_id`/`emp_id`, `target_cust_id: null`. Expected:
- 1 row `m_survey_distributor` aktif (`DistributorId=120`, `CustId=request.CustId` tenant string).
- `m_survey_area` tetap mengikuti area/distributor scope existing dengan `TargetCustId=nil`.

### Case D — `Salesman` + `Principal`

Payload: `distributor_id: [0]`, `area_id: [82]`, `level_target: "Salesman"`, `target_cust_id: "C22001"`, `target_distributor_id: []`. Expected:
- `m_survey_area` row principal dengan `TargetCustId="C22001"`; representasi repo kompatibel: `DistributorId=0` sentinel dan `AreaId=82`/`0` tergantung keputusan worker, tetapi **harus** expose principal di detail lewat `target_cust_id`.
- `m_survey_distributor` kosong.
- `distributor_id=0` di payload tidak di-insert sebagai distributor riil.

### Case E — `Outlet` + `Principal`

Sama dengan D, `level_target: "Outlet"`, dengan `outlet_id` valid. Expected: `m_survey_outlet` row + `m_survey_area` row principal (`TargetCustId="C22001"`) + `m_survey_distributor` kosong.

### Case F — `Distributor` + `Principal`

Payload: `level_target: "Distributor"`, `target_cust_id: "C22001"`, `target_distributor_id: [67, 68]`. Expected:
- 2 rows `m_survey_distributor` (DistributorId=67,68; `CustId=request.CustId` tenant string owner survey) aktif.
- `m_survey_area` principal row dengan `TargetCustId="C22001"`.

### Case G — Edit round-trip

- Create dengan `level_target: "Distributor"`, `target_distributor_id: [120]`.
- Edit dengan `level_target: "Distributor"`, `target_distributor_id: [120, 121]`.
- Expect: row distributor `(120)` lama di-`is_del=true`, dua row baru `(120,121)` aktif. Total `m_survey_distributor` aktif = 2. Tidak ada duplicate.
- `m_survey_area`: replace total; `TargetCustId` reflect payload principal/distributor terbaru, tidak duplicate.

### Case H — GET detail shape

Stub `findDetailSurvey` dengan `LevelTarget="Outlet"`; stub `findDetailAreas` punya `TargetCustId="C22001"`; stub `findDistributorsResult` dengan 1 row `{MSurveyDistributorId:1, DistributorId:120, DistributorCode:"D120", DistributorName:"PT Makmur"}`. Expect response JSON:
```json
{
  "survey_id": 1,
  "...": "...",
  "level_target": "Outlet",
  "business_units": [
    {
      "target_cust_id": "C22001",
      "target_cust_name": "PT Sejahtera"
    }
  ],
  "target_distributor": [
    { "id": 1, "distributor_id": 120, "distributor_code": "D120", "distributor_name": "PT Makmur" }
  ]
}
```
Field existing lain tetap ada.

## Expected Files to Change

- `master/entity/survey.go`
- `master/model/survey.go`
- `master/model/survey_area.go`
- `master/model/survey_distributor.go` (new)
- `master/repository/survey_repository.go`
- `master/service/survey_service.go`
- `master/service/survey_service_test.go`
- `master/controller/survey_controller.go` (jika perlu tambah error mapping)
- `master/controller/survey_controller_test.go`
- `master/pkg/validation/validation.go`
- `master/migration/mst.survey/005_add_level_target_and_target_cust_id.sql` (new)
- `master/migration/mst.survey/006_create_m_survey_distributor.sql` (new)

## Agent/Tool Routing

- Implementasi: `@fixer` (TDD Red → Green → Refactor).
- Riset tambahan (mis. cek staging schema atau library update): `@explorer` read-only; `@librarian` untuk validasi `go-playground/validator` rule tambahan.
- Final conformance: `@quality-gate` (lihat `.opencode/docs/QUALITY.md`).
- Tidak ada UI/FE work, tidak ada visual/browser, tidak ada visual asset.

## Executor Handoff Prompt (copyable)

```
TASK: Implementasi BE SX-2445 / SX-2448 / SX-2452 (Survey Level Target Distributor) di service `master`.
PLAN: .opencode/plans/20260707-1030-sx-2445-2448-2452-survey-level-target-distributor.md
EVIDENCE: .opencode/evidence/20260707-sx-2445-2448-2452-survey-level-target-distributor/discovery.md
DARI DOCX LOKAL:
  - .opencode/evidence/20260707-sx-2445-2448-2452-survey-level-target-distributor/Enhance_Create_Survey_BE.extracted.md
  - .opencode/evidence/20260707-sx-2445-2448-2452-survey-level-target-distributor/Create_Survey_Database.extracted.md

SCOPE (single concrete outcome):
Tambah dukungan `level_target` (Salesman|Outlet|Distributor), `target_cust_id` (varchar(10), DOCX), dan `target_distributor_id` (array) di:
- POST /v1/survey (SX-2445) - Store
- PUT /v1/survey/:survey_id (SX-2448) - Update
- GET /v1/survey/:survey_id (SX-2452) - Detail
Plus schema baru: kolom mst.m_survey.level_target, mst.m_survey_area.target_cust_id, tabel mst.m_survey_distributor dengan unique (survey_id, distributor_id) WHERE is_del = false.
Validator baru `level_target` di master/pkg/validation.
6 test kombinasi BU x level_target + 1 edit round-trip + 1 detail shape.

MUST_PRESERVE:
- Layering Controller -> Service -> Repository -> DB.
- Write in transaction (txManager.WithinTransaction).
- Backward compatibility response (tidak ada field existing yang dihapus).
- Pola sentinel distributor_id=0 = principal (master/service/survey_service.go:147-163).
- Filter 0/negatif dari target_distributor_id SEBELUM insert; distribusi ini HANYA ditulis saat level_target = Distributor (lihat DOCX impact DB dan Assumption A5).
- Replace-on-update untuk m_survey_area dan m_survey_salesman.
- Tenant scope: c.Locals("cust_id") (string) untuk header; field FE body bernama target_cust_id (varchar(10)) sesuai DOCX (D4). m_survey_distributor.cust_id diisi dari request.CustId string tenant.
- Login scope: principal boleh salesman punya principal itu sendiri + salesman punya distributor di bawahnya; distributor hanya salesman punya distributor yang login (lihat DOCX BE section 1649-1650).

DO_NOT_TOUCH:
- Service lain (sales/, tms/, pjp/, dst).
- go.mod/go.sum (kecuali go mod tidy dengan justifikasi).
- docker-compose.yml, master/.env.
- Kontrak field existing (efective_date_*, answer_frequency, response_type, target_type, distributor_id, area_id, outlet_id, emp_id, survey_template_id).
- Path/method/prefix endpoint /v1/survey existing.

VALIDATION:
- cd master && rtk go mod download && rtk go mod tidy
- cd master && rtk go test ./service -run 'TestSurveyService_(Store|Update|Detail)' -v
- cd master && rtk go test ./controller -run 'TestSurveyController_(Create|Update|Detail)' -v
- cd master && rtk go test ./...
- cd master && rtk go vet ./...

EVIDENCE_REQUIRED:
- Output test pass (log ke .opencode/evidence/<task-id>/test-output.txt).
- File migrasi 005 dan 006 di master/migration/mst.survey/.
- Catatan divergensi DOCX vs prompt user (cust_id vs target_cust_id) di evidence.
- Tidak ada secret/kredensial baru.

CLAIM_SCOPE:
- Boleh klaim done jika: semua test pass (8 case + existing), tidak ada regression pada test lama, migrasi SQL siap, response detail memuat field baru.
- Tidak boleh klaim "fix FE juga" atau "rollback DB" - di luar scope.
- DIVERGENSI DOCX vs PROMPT: prompt menyebut field payload "cust_id" numerik; DOCX BE menyebut field "target_cust_id" varchar(10). Implementasi memihak DOCX; jika FE production masih mengirim "cust_id", worker boleh menambah alias parser sementara dengan evidence (A4) dan tidak boleh di-claim done tanpa konfirmasi FE (Widya) atau PO Yogie.

LANE_RESTRICTION:
- Anda (@fixer) saat ini aktif, bukan @artifact-planner. Anda boleh edit source.
- Tulis evidence di .opencode/evidence/20260707-1030-sx-2445-2448-2452-survey-level-target-distributor/ untuk update test output.
- Setelah implementasi, handoff ke @quality-gate.
```

## Execution-ready Worklist / Handoff Contract

Setiap task atomic; owner = `@fixer` kecuali QA = `@quality-gate`. Semua task non-blocked kecuali QA.

### Validator-friendly handoff payload

```yaml
handoff:
  task_id: 20260707-1030-sx-2445-2448-2452-survey-level-target-distributor
  plan_id: 20260707-1030-sx-2445-2448-2452-survey-level-target-distributor
  caller: artifact-planner
  callee: fixer
  scope: Implementasi SX-2445 SX-2448 SX-2452 di service master
  claim_level: scoped
  claim_scope: Implementasi boleh diklaim selesai hanya jika migrasi SQL siap, seluruh test pass, response detail expose field baru, dan tidak ada regression. FE dan rollback DB di luar scope.
  source_basis: [".opencode/evidence/20260707-sx-2445-2448-2452-survey-level-target-distributor/discovery.md", "master/entity/survey.go", "master/service/survey_service.go", "master/repository/survey_repository.go", "master/migration/mst.survey/001_create_tables.sql", "master/migration/mst.survey/002_add_distributor_and_salesman.sql", ".opencode/plans/20260504-1034-sx-1789-survey-business-units.md", ".opencode/plans/20260504-0846-sx-1906-survey-principal-only.md", ".opencode/plans/20260504-2058-sx-1915-salesman-business-unit.md"]
  must_preserve: ["Controller -> Service -> Repository -> DB layering", "Transactional write untuk Store/Update", "Sentinel distributor_id=0 merepresentasikan principal", "Backward compatibility response", "Filter 0/negatif pada target_distributor_id sebelum query/insert", "Replace-on-update untuk m_survey_area dan m_survey_salesman", "LEFT JOIN pada FindAreasBySurveyId agar principal tetap terbaca"]
  do_not_touch: ["service lain (sales/, tms/, pjp/)", "go.mod dan go.sum kecuali go mod tidy dengan justifikasi eksplisit", "docker-compose.yml dan master/.env", "kontrak field existing pada CreateSurveyBody/UpdateSurveyBody/SurveyDetailResponse", "path method prefix endpoint /v1/survey existing"]
  validation: ["cd master && rtk go mod download && rtk go mod tidy", "cd master && rtk go test ./service -run 'TestSurveyService_(Store|Update|Detail)' -v", "cd master && rtk go test ./controller -run 'TestSurveyController_(Create|Update|Detail)' -v", "cd master && rtk go test ./...", "cd master && rtk go vet ./..."]
  exit_criteria: ["Semua test 8 case A-H pass", "Test existing survey service/controller/repository pass tanpa regression", "File migrasi 005 dan 006 ada di master/migration/mst.survey/ dengan syntax valid", "Response detail GET /v1/survey/:id memuat level_target target_cust_id target_distributor", "Stub surveyRepositoryRedStub diextensi untuk method distributor baru", "Tidak ada secret/kredensial baru di commit"]
  evidence_required: [".opencode/evidence/20260707-1030-sx-2445-2448-2452-survey-level-target-distributor/test-output.txt", "master/migration/mst.survey/005_add_level_target_and_target_cust_id.sql", "master/migration/mst.survey/006_create_m_survey_distributor.sql"]
  depends_on: ["none"]
  context_bundle: ["master/entity/survey.go", "master/model/survey.go", "master/model/survey_area.go", "master/repository/survey_repository.go", "master/service/survey_service.go", "master/service/survey_service_test.go", "master/controller/survey_controller.go", "master/controller/survey_controller_test.go", "master/pkg/validation/validation.go", "master/migration/mst.survey/001_create_tables.sql", "master/migration/mst.survey/002_add_distributor_and_salesman.sql"]
```

### Numbered worklist

1. **A1** | `@fixer` | Tambah `master/migration/mst.survey/005_add_level_target_and_target_cust_id.sql`.
2. **A2** | `@fixer` | Tambah `master/migration/mst.survey/006_create_m_survey_distributor.sql`.
3. **A3** | `@fixer` | Tambah `master/model/survey_distributor.go` dan extend `master/model/survey.go` + `master/model/survey_area.go`.
4. **A4** | `@fixer` | Extend `master/entity/survey.go` untuk request/response baru (`level_target`, payload `target_cust_id`, `target_distributor_id`, `target_distributor` response, `target_cust_id` / `target_cust_name` pada `business_units`).
5. **A5** | `@fixer` | Tambah validator `level_target` di `master/pkg/validation/validation.go`.
6. **A6** | `@fixer` | Extend `master/repository/survey_repository.go` untuk schema baru dan method distributor mapping.
7. **A7** | `@fixer` | Extend `master/service/survey_service.go` untuk Store/Update/Detail + helper distributor scope.
8. **A8** | `@fixer` | Tambah/extend `master/service/survey_service_test.go` untuk case A-H dan extend stub repository.
9. **A9** | `@fixer` | Tambah/extend `master/controller/survey_controller_test.go` untuk parsing body dan error mapping field baru.
10. **A10** | `@fixer` | Final conformance: `rtk go mod tidy`, `rtk go test ./...`, `rtk go vet ./...`, dan simpan evidence test.
11. **Q1** | `@quality-gate` | Review conformance, tenant safety, idempotency, no regression, dan verdict final.

### Detailed task map

```yaml
worklist:
  - id: A1
    title: "Tambah migration 005_add_level_target_and_target_cust_id.sql"
    owner: "@fixer"
    depends_on: "none"
    must_preserve:
      - "Idempotent: ADD COLUMN IF NOT EXISTS, CREATE INDEX IF NOT EXISTS."
      - "Tidak backfill data historis (lihat Out of Scope)."
    do_not_touch:
      - "File migrasi 001-004."
      - "Service lain."
    validation:
      - "File ada di master/migration/mst.survey/005_add_level_target_and_target_cust_id.sql."
      - "Reviewer membaca isi file; ALTER TABLE aman untuk environment existing."
    exit_criteria:
      - "File SQL dibuat, syntax benar, tidak ada data destructive."
    blocking: "ready"
    requires_user_decision: "no"
    evidence_update: "Catat path file di .opencode/evidence/<task-id>/A1-migration-005.md."
    exit_verification:
      - "File exists; head -n 5 master/migration/mst.survey/005_add_level_target_and_target_cust_id.sql"
    start_with: "A1"

  - id: A2
    title: "Tambah migration 006_create_m_survey_distributor.sql"
    owner: "@fixer"
    depends_on: "A1"
    must_preserve:
      - "Unique (survey_id, distributor_id) via partial unique index (WHERE is_del = false)."
      - "Idempotent: CREATE TABLE IF NOT EXISTS, CREATE INDEX IF NOT EXISTS."
    do_not_touch:
      - "Tabel existing m_survey_*, m_distributor, m_salesman."
    validation:
      - "File ada; reviewer cek struktur kolom."
    exit_criteria:
      - "File SQL siap; CREATE TABLE + UNIQUE INDEX partial."
    blocking: "ready"
    requires_user_decision: "no"
    evidence_update: "Catat path di .opencode/evidence/<task-id>/A2-migration-006.md."
    exit_verification:
      - "head master/migration/mst.survey/006_create_m_survey_distributor.sql"
    start_with: null

  - id: A3
    title: "Tambah model/model_survey_distributor.go dan extend model/survey.go + model/survey_area.go"
    owner: "@fixer"
    depends_on: "A1"
    must_preserve:
      - "DB tag = snake_case."
      - "Joined fields optional (DistributorName *string)."
    do_not_touch:
      - "Field existing di model."
    validation:
      - "go build ./... di master/ bersih."
    exit_criteria:
      - "Model baru ada; build OK."
    blocking: "ready"
    requires_user_decision: "no"
    evidence_update: "Diff per file."
    exit_verification:
      - "cd master && rtk go build ./..."
    start_with: null

  - id: A4
    title: "Extend entity/survey.go: CreateSurveyBody, UpdateSurveyBody, SurveyDetailResponse, SurveyBusinessUnit, SurveyDistributorResponse + rename cust_id payload"
    owner: "@fixer"
    depends_on: "A1"
    must_preserve:
      - "CustId (string, dari token) tetap ada."
      - "Field existing tidak dihapus atau diganti nama."
      - "Field payload baru ditambahkan dengan json tag dan validator tag (omitempty,level_target)."
    do_not_touch:
      - "SurveyListResponse, SurveyParams, SurveyQueryFilter, DeactivateSurveyBody (kecuali perlu extensibility minor)."
    validation:
      - "go build ./... di master/ bersih."
    exit_criteria:
      - "Entity ter-extensi; field baru siap dipakai di JSON binding."
    blocking: "ready"
    requires_user_decision: "yes (rename cust_id payload — lihat Assumption A6)"
    evidence_update: "Diff entity/survey.go."
    exit_verification:
      - "cd master && rtk go build ./..."
    start_with: null

  - id: A5
    title: "Tambah validator level_target di master/pkg/validation/validation.go dengan translation EN/ID"
    owner: "@fixer"
    depends_on: "A4"
    must_preserve:
      - "Validator existing (answer_frequency, qtystr, dll.) tidak dihapus."
      - "Translation EN/ID mengikuti pola existing."
    do_not_touch:
      - "Validator lain."
    validation:
      - "go build ./... di master/ bersih."
      - "Test unit validator (jika ada) pass."
    exit_criteria:
      - "Rule level_target aktif, value Salesman|Outlet|Distributor."
    blocking: "ready"
    requires_user_decision: "no"
    evidence_update: "Diff validation.go."
    exit_verification:
      - "cd master && rtk go build ./..."
    start_with: null

  - id: A6
    title: "Extend repository/survey_repository.go: Store/Update/StoreAreas/FindAreasBySurveyId/FindOneById + tambah StoreSurveyDistributors, DeleteSurveyDistributorsBySurveyId, FindSurveyDistributorsBySurveyId"
    owner: "@fixer"
    depends_on: "A1, A3, A4"
    must_preserve:
      - "Pola StoreAreas (loop tx.Exec) & DeleteAreasBySurveyId (set is_del=true)."
      - "Pola LEFT JOIN di FindAreasBySurveyId (jangan INNER JOIN karena sentinel distributor_id=0)."
      - "FindOneById tetap join ke m_distributor dengan MIN(distributor_id) (backward compat)."
    do_not_touch:
      - "Method repository lain."
      - "Layer service di sini (extension di A7)."
    validation:
      - "go build ./... di master/ bersih."
      - "Stub `surveyRepositoryRedStub` di-extensi dengan method baru."
    exit_criteria:
      - "Method baru ada di interface + impl; stub mengimplementasikan."
    blocking: "ready"
    requires_user_decision: "no"
    evidence_update: "Diff repository."
    exit_verification:
      - "cd master && rtk go build ./..."
      - "cd master && rtk go vet ./..."
    start_with: null

  - id: A7
    title: "Extend service/survey_service.go: helper normalize + buildSurveyDistributors, Store, Update, Detail (Red → Green → Refactor)"
    owner: "@fixer"
    depends_on: "A4, A5, A6"
    must_preserve:
      - "Transaction wrapping (txManager.WithinTransaction)."
      - "Sentinel distributor_id=0 = principal."
      - "Replace-on-update untuk m_survey_area dan m_survey_salesman existing."
    do_not_touch:
      - "Helper normalizeBusinessUnitSelection, resolveSurveyCustIds, resolveSalesmanCustIds (reuse; ekstensi minor saja)."
    validation:
      - "rtk go test ./service -run 'TestSurveyService_(Store|Update|Detail)'"
      - "rtk go test ./service"
    exit_criteria:
      - "8 test case A–H pass; test existing pass (no regression)."
    blocking: "ready"
    requires_user_decision: "no"
    evidence_update:
      - "Output test pass ke .opencode/evidence/<task-id>/A7-test-output.txt"
      - "Diff service/survey_service.go."
    exit_verification:
      - "cd master && rtk go test ./service"
    start_with: null

  - id: A8
    title: "Tambah/extend test di service/survey_service_test.go untuk 8 case (A–H) + extend stub surveyRepositoryRedStub"
    owner: "@fixer"
    depends_on: "A7"
    must_preserve:
      - "Test existing tidak dihapus; stub existing method tetap diimplementasikan."
    do_not_touch:
      - "Test existing selain yang harus diupdate minor karena interface berubah."
    validation:
      - "rtk go test ./service -run 'TestSurveyService_(Store|Update|Detail)' -v"
    exit_criteria:
      - "8 test case pass; existing test pass."
    blocking: "ready"
    requires_user_decision: "no"
    evidence_update:
      - "Output test ke .opencode/evidence/<task-id>/A8-test-output.txt"
    exit_verification:
      - "cd master && rtk go test ./service -v"
    start_with: null

  - id: A9
    title: "Tambah/extend test di controller/survey_controller_test.go (parse body + error mapping)"
    owner: "@fixer"
    depends_on: "A4, A7"
    must_preserve:
      - "Stub existing tidak dihapus."
    do_not_touch:
      - "Test existing."
    validation:
      - "rtk go test ./controller -run 'TestSurveyController_(Create|Update|Detail)' -v"
    exit_criteria:
      - "Test baru pass; existing test pass."
    blocking: "ready"
    requires_user_decision: "no"
    evidence_update:
      - "Output test ke .opencode/evidence/<task-id>/A9-test-output.txt"
    exit_verification:
      - "cd master && rtk go test ./controller -v"
    start_with: null

  - id: A10
    title: "Final conformance: go test ./..., go vet, dokumentasi singkat"
    owner: "@fixer"
    depends_on: "A1-A9"
    must_preserve:
      - "Tidak ada go.mod change tanpa justifikasi."
    do_not_touch:
      - "Service lain."
    validation:
      - "cd master && rtk go mod tidy"
      - "cd master && rtk go test ./..."
      - "cd master && rtk go vet ./..."
    exit_criteria:
      - "Semua test pass; vet bersih; tidak ada perubahan file di luar scope."
    blocking: "ready"
    requires_user_decision: "no"
    evidence_update:
      - "Final test output ke .opencode/evidence/<task-id>/A10-final.txt"
    exit_verification:
      - "cd master && rtk go test ./... && rtk go vet ./..."
    start_with: null

  - id: Q1
    title: "Quality gate: review conformance, keamanan tenant, idempotency, no regression, response shape"
    owner: "@quality-gate"
    depends_on: "A1-A10"
    must_preserve:
      - "Final signoff tidak menghapus file plan/evidence."
    do_not_touch:
      - "Service lain."
    validation:
      - "Bukti test pass."
      - "Diff ringkas (file yang berubah + alasan)."
      - "Cek plan quality gate (PASS / PASS_FOR_SLICE / NEEDS_DEPTH / BLOCKED)."
    exit_criteria:
      - "Signoff dengan rekomendasi PASS atau PASS_FOR_SLICE."
    blocking: "ready"
    requires_user_decision: "no"
    evidence_update:
      - ".opencode/evidence/<task-id>/Q1-quality-gate.md"
    exit_verification:
      - "Catatan signoff di evidence."
    start_with: null

ownership_table:
  - subsystem: "DB schema & migration"
    implementation_owner: "@fixer"
    review_gate: "@quality-gate"
  - subsystem: "Entity & validator"
    implementation_owner: "@fixer"
    review_gate: "@quality-gate"
  - subsystem: "Repository"
    implementation_owner: "@fixer"
    review_gate: "@quality-gate"
  - subsystem: "Service"
    implementation_owner: "@fixer"
    review_gate: "@quality-gate"
  - subsystem: "Controller"
    implementation_owner: "@fixer"
    review_gate: "@quality-gate"
  - subsystem: "Test (service + controller)"
    implementation_owner: "@fixer"
    review_gate: "@quality-gate"

start_with: "A1"
```

## Validation Commands

Jalankan dari `master/`:

```bash
# 1. Sync deps
cd master
rtk go mod download && rtk go mod tidy

# 2. Targeted test
rtk go test ./service -run 'TestSurveyService_(Store|Update|Detail)' -v
rtk go test ./controller -run 'TestSurveyController_(Create|Update|Detail)' -v
rtk go test ./repository -run TestSurvey

# 3. Full module test (tidak boleh ada regression)
rtk go test ./...

# 4. Vet
rtk go vet ./...

# 5. Compose status (baseline)
cd ..
rtk docker compose -f docker-compose.yml ps
```

Validasi runtime opsional jika ada token lokal/staging (lihat `.opencode/plans/20260504-1034-sx-1789-survey-business-units.md` untuk pola `curl`); tidak wajib untuk plan ini karena QA/PIC test akan dilakukan manual.

## Evidence Requirements

1. **Test output** (wajib): simpan log `rtk go test ./...` ke `.opencode/evidence/<task-id>/test-output.txt`.
2. **Migration files** (wajib): file `005_…` dan `006_…` di `master/migration/mst.survey/`.
3. **DB verification** (manual, saat deployment): jalankan query `information_schema` berikut dan simpan output:
   ```sql
   SELECT column_name, data_type, is_nullable FROM information_schema.columns
   WHERE table_schema='mst' AND table_name IN ('m_survey','m_survey_area','m_survey_distributor')
   ORDER BY table_name, ordinal_position;

   SELECT con.conname, con.contype FROM pg_constraint con
   JOIN pg_class rel ON rel.oid = con.conrelid
   WHERE rel.relname IN ('m_survey_area','m_survey_distributor');
   ```
4. **Manual API smoke** (opsional, jika env tersedia): `POST /master/v1/survey` dengan payload Case A, lalu `GET detail`, lalu cek `m_survey_area` & `m_survey_distributor` di DB.
5. **Final summary** oleh `@quality-gate` di `.opencode/evidence/<task-id>/Q1-quality-gate.md` dengan verdict (PASS / PASS_FOR_SLICE / NEEDS_DEPTH / BLOCKED).

## Done Criteria

- Semua acceptance criteria A–H pass.
- Test existing survey service/controller/repository pass tanpa modifikasi besar (rename `cust_id` payload harus terkontrol).
- File migrasi 005 & 006 ada, syntax valid, idempotent.
- Response detail `GET /v1/survey/:id` memuat `level_target`, `target_cust_id`, dan `target_distributor`.
- Tidak ada secret/kredensial baru di commit.
- Quality gate signoff PASS atau PASS_FOR_SLICE.
- Documentation update minimal: header comment di `master/migration/mst.survey/005_…` dan `006_…` menjelaskan asal tiket.

## Out of Scope (Sengaja Tidak Diubah)

- Dropdown filter Survey Target (Salesman / Outlet / Sales Team) — issue `SX-1578` / `SX-1915` sudah ada.
- Cross-check penuh `cust_id` payload dengan tree tenant (`parent_cust_id` traversal) — slice berikutnya.
- Migration backfill `target_cust_id` historis — biarkan `null` (lihat Assumption A1/A2).
- Audit log perubahan `level_target` di edit.
- Index tambahan selain unique index di `m_survey_distributor`.

## Skipped

- Backfill `target_cust_id` untuk survey existing — biarkan `null`, FE sudah handle. Add when ada requirement historical reporting.
- Index tambahan di `m_survey_distributor` selain unique — add when query list surveyor > 1k.
- Audit log perubahan `level_target` di edit — add when ada requirement compliance.

## Final Planning Summary

- **Artifacts consulted**:
  - `.opencode/docs/ARCHITECTURE.md`, `SERVICE_MATRIX.md`, `QUALITY.md`, `AGENT_ROUTING.md`, `SECURITY.md`, `PROMPT_GATES.md`, `DECISIONS.md`, `MCP.md`
  - Plan lama: `.opencode/plans/20260504-1034-sx-1789-survey-business-units.md`, `20260504-0846-sx-1906-survey-principal-only.md`, `20260504-2058-sx-1915-salesman-business-unit.md`
  - Source code: `master/entity/survey.go`, `master/service/survey_service.go`, `master/repository/survey_repository.go`, `master/controller/survey_controller.go`, `master/pkg/validation/validation.go`, `master/migration/mst.survey/001_create_tables.sql`, `002_add_distributor_and_salesman.sql`, `master/service/survey_service_test.go`, `master/controller/survey_controller_test.go`
  - Docs lokal hasil ekstraksi: `.opencode/evidence/20260707-sx-2445-2448-2452-survey-level-target-distributor/Enhance_Create_Survey_BE.extracted.md`, `.opencode/evidence/20260707-sx-2445-2448-2452-survey-level-target-distributor/Create_Survey_Database.extracted.md`
- **Artifacts created**:
  - Primary plan: `.opencode/plans/20260707-1030-sx-2445-2448-2452-survey-level-target-distributor.md`
  - Evidence: `.opencode/evidence/20260707-sx-2445-2448-2452-survey-level-target-distributor/discovery.md`
  - Dokumen lokal yang diekstrak: `Enhance_Create_Survey_BE.extracted.md`, `Create_Survey_Database.extracted.md` (di folder evidence yang sama).
- **Key decisions**:
  - Migration terpisah (005 schema, 006 tabel) untuk rollback yang lebih aman.
  - Field payload FE `target_cust_id` (`varchar(10)`, DOCX) dipisah dari `request.CustId` (string tenant) di entity agar tidak bentrok semantik. Ticket summary menyebut `cust_id` integer, plan memihak DOCX (lihat D4/A4).
  - `m_survey_distributor.cust_id` = `varchar(10)` tenant string, mengikuti konvensi `m_survey_salesman.cust_id` (DOCX).
  - Validator `level_target` didaftarkan EN/ID mengikuti pola `answer_frequency`.
  - `target_distributor_id` di-filter `0`/negatif sebelum insert; `target_distributor_id` hanya ditulis ke `m_survey_distributor` saat `level_target = Distributor` (DOCX impact DB).
  - BU principal: `m_survey_area` row principal dengan `target_cust_id` non-null + `distributor_id=0` (sentinel existing). BU distributor: `target_cust_id=NULL` di area.
- **Assumptions** (lihat §Decisions/Assumptions): A1 (FK distributor_id di staging), A2 (legacy payload tanpa level_target), A3 (rule status edit), A4 (divergensi DOCX vs prompt `cust_id`/`target_cust_id`), A5 (`target_distributor_id` hanya untuk `level_target=Distributor`), A6 (migration manual), A7 (nama field response `target_distributor` vs `distributor`).
- **Open questions**:
  - Apakah FE produksi mengirim `target_cust_id` (varchar) atau `cust_id` (int)? Worker tambah alias sementara dengan evidence (A4), tidak boleh di-claim done tanpa konfirmasi Widya/Yogie.
  - Apakah `target_cust_id` boleh `null` pada principal scope? DOCX body table: target_cust_id mandatory saat business unit principal. Worker enforce sebagai required.
  - Apakah FE menggunakan field response `target_distributor` (daftar atribut) atau `distributor` (contoh JSON level Distributor)? DOCX inkonsisten. Plan memihak `target_distributor` (lihat A7).
  - Apakah staging `m_survey_area.distributor_id` punya FK ke `m_distributor`? Implementer cek manual via `information_schema` sebelum deploy.
- **Readiness**:
  - Plan ini **`PASS_FOR_SLICE`** — siap dieksekusi oleh `@fixer` (TDD) sesuai worklist. Bukan `PASS` penuh karena ada open question FE (nama field response) yang harus dikonfirmasi sebelum final conformance.
- **Cleanup performed**:
  - Plan dan evidence baru ditulis di `.opencode/`. Tidak ada plan/evidence lama yang dihapus — histori SX-1789/SX-1906/SX-1915 tetap berguna untuk referensi pola.
  - Tidak ada draft artifact yang disimpan di `.opencode/draft/` (rancangan langsung di plan utama).
- **Active-lane reset note**:
  - Plan ini ditulis oleh `@artifact-planner` (read-only). Eksekusi (`@fixer`) akan dilakukan di lane berikutnya dengan permission `edit` source code; planner tidak lagi menjadi active lane saat eksekusi berjalan.

## Quality checklist

- [x] Stack docs read and current best practice verified.
- [x] Question Gate: tidak ada pertanyaan material yang blocking; requirement tiket sudah rinci.
- [x] Research Gate source strategy explicit (repo-local primary; context7/GitHub skipped with reason).
- [x] Discovery evidence written to `.opencode/evidence/<task-id>/discovery.md`.
- [x] Primary plan written with all required sections.
- [x] Plan depth proportional (substantial non-trivial backend; deep enough for safe execution).
- [x] Ruthless slicing: bounded first slice = 3 endpoint + schema + 8 test cases (sesuai tiket).
- [x] Scope expansion guard: 1 slice, 3 endpoint, 1 tabel baru, 6 kombinasi test — masih dalam batas wajar.
- [x] Execution-ready worklist atomic with owner/depends_on/validation/exit_criteria.
- [x] Handoff prompt copyable.
- [x] Reference Map: docs/repo-backed, dengan rationale first-principles untuk 2 klaim (Assumption A1, A6).
- [x] Confirmed vs Assumed Audit table present.
- [x] Anti-Generic (n/a untuk backend).
- [x] TDD plan: 10 first failing tests, Green/Refactor steps, edge cases, commands.
- [x] Plan Quality Gate: `PASS_FOR_SLICE`.
