# Discovery — SX-2445 / SX-2448 / SX-2452 Survey Level Target Distributor

Ticket terkait: SX-2047 (parent enhance Manage Survey - Survey List) → SX-2445, SX-2448, SX-2452.
Service target: `master` (modul `master/`, Fiber + sqlx + Postgres, default port `9002`, env `master/.env`, baseline `rtk go mod download && rtk go mod tidy && rtk go test ./...` per `.opencode/docs/QUALITY.md`).
Rencana dieksekusi oleh `@fixer` (TDD) setelah plan ini disetujui; planner hanya menulis artefak.

## Repo state saat ini

- `master/entity/survey.go`:
  - `CreateSurveyBody` / `UpdateSurveyBody` tidak punya field `level_target`, `cust_id` payload, atau `target_distributor_id`.
  - `CustId` pada body bertipe `string` dan diisi dari token (`c.Locals("cust_id")`) — bukan scope owner numerik dari FE.
  - `SurveyDetailResponse` punya `DistributorId`, `AreaId`, `BusinessUnits` (extended dengan `BusinessUnitName`, `Name`, `Type` dari SX-1789). Belum ada `LevelTarget`, `TargetCustId`, atau `TargetDistributor`.
- `master/model/survey.go`:
  - `model.Survey` merepresentasikan `mst.m_survey`; tidak ada `LevelTarget` atau `TargetCustId` di level header.
- `master/model/survey_area.go`:
  - Kolom `survey_area_id`, `survey_id`, `distributor_id`, `area_id`, `is_del` + joined `area_name`/`distributor_name`. Tidak ada `target_cust_id`.
- `master/migration/mst.survey/`:
  - `001_create_tables.sql` membuat `m_survey_area` dengan `distributor_id` (awalnya nullable, lalu diset NOT NULL di 002).
  - `002_add_distributor_and_salesman.sql` backfill `distributor_id = 0` untuk principal sentinel.
  - `003_add_survey_salesman_guardrails.sql` dan `004_alter_answer_frequency_length.sql` ada, tidak relevan dengan schema baru.
- `master/service/survey_service.go`:
  - `Store` / `Update` sudah:
    - parse date → `ErrSurveyInvalidDateFormat`,
    - cek overlap title → `ErrSurveyTitleConflict`,
    - normalisasi BU via `normalizeBusinessUnitSelection` (0 = principal, >0 = distributor riil),
    - resolve area via `FindSurveyAreasByDistributorIds` (gabungan payload + lookup),
    - resolve salesman via `resolveSurveyCustIds` + `resolveSalesmanCustIds`,
    - `buildSurveyAreas` menulis baris `model.SurveyArea` (DistributorId, AreaId) — belum ada `TargetCustId`.
  - `Detail` membentuk `BusinessUnits` dengan sentinel `Principal` ketika `DistributorId == 0`, plus `area`/`outlet`/`salesman` di `target_survey`.
  - `Update` sudah replace-on-update untuk `m_survey_area` dan `m_survey_salesman` lewat `DeleteAreasBySurveyId` + `StoreAreas` (pola yang sama bisa dipakai untuk `m_survey_distributor` baru).
- `master/repository/survey_repository.go`:
  - `FindOneById` join `mst.m_distributor` via `MIN(distributor_id)` di `m_survey_area`. Tidak expose distributor list multi-row.
  - `FindAreasBySurveyId` left join `mst.m_area` dan `mst.m_distributor`; tidak expose `target_cust_id`.
  - `StoreAreas`, `DeleteAreasBySurveyId`, `FindAreasBySurveyId` ada dan idempotent (pakai `is_del`).
  - Tidak ada method `StoreSurveyDistributors`, `DeleteSurveyDistributorsBySurveyId`, `FindSurveyDistributorsBySurveyId`.
- `master/controller/survey_controller.go`:
  - Route `POST /v1/survey`, `GET /v1/survey/:survey_id`, `PUT /v1/survey/:survey_id` (prefix `/v1` di-mount lewat `main.go`).
  - `Create` / `Update` parsing body via `json.Unmarshal` (untuk dukung `FlexibleIntArray`).
  - `Deactivate` tidak relevan dengan scope ticket.
- `master/pkg/validation/validation.go`:
  - Custom validator `answer_frequency` terdaftar. Untuk `level_target` perlu tambahkan validator `level_target` baru (mirip `oneof=Salesman Outlet Distributor`) dan registrasikan translation EN/ID.
- `master/main.go`:
  - `surveyRepository := repository.NewSurveyRepository(postgreDB)` dan `surveyService := service.NewSurveyService(txManagerRepository, surveyRepository, salesmanRepository, ...)` di line 135 dan 252. Implementasi baru akan menambah wire-up dependensi distributor lookup bila perlu.

## Plan lama terkait (reused as pattern)

- `.opencode/plans/20260504-1034-sx-1789-survey-business-units.md` → pola normalize BU 0/positive, replace-on-update, sentinel `Principal` di response.
- `.opencode/plans/20260504-0846-sx-1906-survey-principal-only.md` → cara principal-only `[0]` menyimpan `m_survey_area` dengan `distributor_id=0` + payload area.
- `.opencode/plans/20260504-2058-sx-1915-salesman-business-unit.md` → referensi pola validasi cust scope untuk multi distributor.
- `.opencode/plans/20260504-2058-sx-1915-salesman-business-unit.md` (refs) → `buildSalesTeamCustScopeCondition` untuk scope `0` + child distributor via `smc.m_customer`.
- Pola test Red/Green/Refactor dengan `surveyRepositoryRedStub` (interface `SurveyRepository`) dan `salesmanRepositoryStub` di `master/service/survey_service_test.go`.

## Gap yang harus diisi plan

1. **Schema baru menurut DOCX**:
   - Tambah kolom `mst.m_survey.level_target` (doc menyebut enum `Outlet | Distributor | Salesman`; repo lokal belum punya field ini).
   - Tambah kolom `mst.m_survey_area.target_cust_id VARCHAR(10) NULL` (doc `Create_Survey_Database.docx` line 116-119; untuk principal scope field ini berisi `cust_id` business unit, untuk distributor scope `NULL`).
   - Tabel baru `mst.m_survey_distributor` (doc `Create_Survey_Database.docx` line 171-198; doc menyebut kolom `id`, `cust_id`, `survey_id`, `distributor_id`).
2. **Payload FE vs tenant token**:
   - DOCX BE menyebut field body **`target_cust_id`** (bukan `cust_id`) bertipe `varchar(10)` dan mandatory saat business unit principal (`Enhance_Create_Survey_BE.extracted.md` line 760-765, 1700-1705).
   - Field ini tidak sama dengan `request.CustId` string tenant dari token, tetapi secara domain masih berupa `cust_id` business unit. Service tetap harus memvalidasi scope tenant.
3. **Perilaku principal vs distributor menurut DOCX**:
   - Business unit **Distributor** → `target_cust_id = NULL`, `m_survey_area` berisi `distributor_id` + `area_id` dari distributor, dan `m_survey_distributor` hanya dipakai jika `level_target = Distributor`.
   - Business unit **Principal** → `target_cust_id` berisi `cust_id` principal, dan contoh doc menunjukkan `m_survey_area.distributor_id = NULL`, `area_id = NULL`, `target_cust_id = 'C22001'` untuk principal-only (`Enhance_Create_Survey_BE.extracted.md` line 1322-1334, 1457-1469, 1588-1600).
4. **Service/repo/response extension**:
   - Repository `SurveyRepository`:
     - `Store` / `Update` harus menulis `level_target` dan `target_cust_id` (di area).
     - Tambah `StoreSurveyDistributors`, `DeleteSurveyDistributorsBySurveyId`, `FindSurveyDistributorsBySurveyId`.
     - `FindAreasBySurveyId` harus expose `target_cust_id` (kolom `target_cust_id` di select).
   - Service:
     - `Store` / `Update` normalisasi `level_target` (default fallback behavior lihat Assumption A2) dan `target_distributor_id` (filter `0` dan negative).
     - `Detail` tambahkan `LevelTarget`, `TargetCustId` (header), dan `TargetDistributor` list.
   - Entity:
     - Tambah `LevelTarget string` + `TargetCustId string` + `TargetDistributor []SurveyDistributorResponse` di `SurveyDetailResponse`.
     - Tambah field payload baru di `CreateSurveyBody` / `UpdateSurveyBody` (lihat Implementation Steps §3.1).
   - Validator: tambah `level_target` rule dengan translation EN/ID.
4. **Migration numbering**: lanjut `005_add_level_target_and_target_cust_id.sql` (atau file 005 khusus `m_survey_distributor`). Pattern: schema `mst`, prefix `m_`, snake_case, idempotent (`IF NOT EXISTS` / `ADD COLUMN IF NOT EXISTS`).
5. **Acceptance test**:
   - Tambah test untuk 6 kombinasi BU × level_target (salesman/outlet/distributor + distributor/principal) pada `survey_service_test.go` dan `survey_controller_test.go` (post-valid).
   - Test edit round-trip (salesman → outlet) + cek tidak ada duplikat `m_survey_area` / `m_survey_distributor`.
   - Test detail response shape (level_target, target_cust_id, target_distributor).
6. **Handoff / sumber daya**:
   - Schema doc resmi (Google Doc) tidak tersedia di repo. Plan harus menyatakan ini sebagai Assumption A1 dan merekomendasikan cek manual saat eksekusi.
   - Tidak ada dependensi library baru yang perlu dipasang (`go.mod` cukup).

## Confirmed vs Assumed Audit

| Material claim | Status | Source |
| --- | --- | --- |
| `master` service pakai Fiber + sqlx + Postgres | confirmed_repo | `master/go.mod`, `master/main.go:135,252`, `.opencode/docs/SERVICE_MATRIX.md` |
| `POST /v1/survey` dan `PUT /v1/survey/:survey_id` ada | confirmed_repo | `master/controller/survey_controller.go:30-35` |
| `Store` / `Update` sudah replace-on-update untuk `m_survey_area` | confirmed_repo | `master/service/survey_service.go:600-610`, `master/service/survey_service_test.go:698` |
| Sentinel `0` di `distributor_id` diperlakukan sebagai Principal | confirmed_repo | `master/service/survey_service.go:147-163`, `master/migration/mst.survey/002_add_distributor_and_salesman.sql` |
| `m_survey_area` saat ini punya kolom `distributor_id NOT NULL`, `area_id NOT NULL` | confirmed_repo | `master/migration/mst.survey/001_create_tables.sql:27-33`, `002_add_distributor_and_salesman.sql:1-9` |
| Tidak ada kolom `level_target` / `target_cust_id` / `m_survey_distributor` | confirmed_repo | `rg "level_target|target_distributor|target_cust_id|m_survey_distributor" master/` returns 0 hit |
| `cust_id` dari token bertipe `string` (`C1001`) di semua request survey existing | confirmed_repo | `master/controller/survey_controller.go:139,143,208,212` |
| DOCX memakai field body `target_cust_id` bertipe `varchar(10)` dan mandatory saat business unit principal | confirmed_docs | `Enhance_Create_Survey_BE.extracted.md:760-765`, `1700-1705` |
| DOCX detail response menambah `level_target`, `business_unit[].target_cust_id`, `target_cust_name`, dan `target_distributor` | confirmed_docs | `Enhance_Create_Survey_BE.extracted.md:254-257`, `273-296`, `379-388` |
| DOCX `m_survey_distributor` berisi `cust_id`, `survey_id`, `distributor_id` | confirmed_docs | `Create_Survey_Database.extracted.md:171-198` |
| DOCX contoh principal-only menaruh `target_cust_id` di `m_survey_area` dan `distributor_id/area_id = NULL` | confirmed_docs | `Enhance_Create_Survey_BE.extracted.md:1322-1334`, `1457-1469`, `1588-1600` |
| `m_survey_area.distributor_id` punya FK ke `m_distributor` di staging | assumption | doc schema lokal tidak menunjukkan FK eksplisit; plan merekomendasikan cek staging sebelum deploy migration |
| 6 kombinasi BU × level_target valid (sesuai tabel di tiket) | user_confirmed | prompt SX-2445/2448/2452 |
| Ticket summary user menyebut field payload `cust_id`, tetapi DOCX BE menyebut `target_cust_id`; implementasi harus memihak DOCX atau support alias sementara | assumption | prompt user vs DOCX berbeda nama field |
| `target_distributor_id` biasanya kosong untuk principal scope | user_confirmed | prompt SX-2445 (Case D–E–F) + DOCX examples principal (`target_distributor_id: []`) |
| Validator `level_target` akan didaftarkan EN/ID | assumption | repo belum punya rule ini; default aman: tambahkan dengan translation sama seperti `answer_frequency` |
| `m_survey.cust_id` di repo lokal adalah `VARCHAR(10)` tetapi doc database menyebut `int(8)` | assumption | repo vs doc divergen; implementasi harus mengikuti schema repo aktual sambil memenuhi contract baru di field tambahan |

## Implikasi & risiko

- **Resiko schema**: bila staging punya FK `m_survey_area.distributor_id` ke `m_distributor`, persist `0` akan gagal. Mitigasi: cek staging manual sebelum deploy; existing plan SX-1789 sudah membahas ini.
- **Resiko kontrak**: rename/extend `SurveyDetailResponse` bisa membuat FE existing bergantung pada field tertentu. Mitigasi: tambahkan field baru di top level, tidak menghapus field lama.
- **Resiko determinisme**: `Store` saat ini menentukan `principal`/`distributor` berdasarkan payload `DistributorId` saja. Plan akan menambah `level_target` sebagai penentu utama untuk menyimpan `m_survey_distributor`. Pastikan helper normalization diuji.
- **Resiko transaksional**: replace `m_survey_area` + insert `m_survey_distributor` + replace `m_survey_salesman` harus tetap dalam satu `txManager.WithinTransaction`. Pola existing sudah aman; tidak ada perubahan `txManager` yang diperlukan.

## Cleanup artifact sebelumnya

- Plan lama `20260504-*` tetap di `.opencode/plans/` sebagai histori. Plan baru ini (`20260707-SX-2445-2448-2452-...`) akan menjadi sumber tunggal untuk eksekusi 3 tiket ini; cleanup tidak dilakukan oleh planner (lihat Final Planning Summary §Cleanup).
