# Plan — SX-1906 Survey Principal-only Business Unit

## Goal

Perbaiki backend `master` agar `POST /master/v1/survey` berhasil membuat survey ketika payload create berisi `distributor_id: [0]` untuk Business Unit principal-only, dengan `area_id` dan `emp_id` tetap diproses.

## Non-goals

- Tidak mengubah kontrak frontend selain menerima kontrak existing `0 = principal`.
- Tidak menghapus validasi distributor/salesman untuk flow distributor normal.
- Tidak mengubah schema database kecuali implementasi menemukan FK/staging schema berbeda dari migration lokal.
- Tidak hardcode `emp_id: 369`, `survey_template_id: 63`, atau data QA tertentu.

## Scope

- Modul: `master`.
- Endpoint: `POST /master/v1/survey`, kemungkinan juga `PUT /master/v1/survey/:survey_id` agar perilaku create/update konsisten.
- Area kode utama:
  - `master/controller/survey_controller.go`
  - `master/entity/survey.go`
  - `master/service/survey_service.go`
  - `master/repository/survey_repository.go` bila perlu query pendukung area principal.
  - `master/service/survey_service_test.go`

## Requirements

- Payload `distributor_id: [0]` harus dikenali sebagai principal-only sentinel, bukan distributor riil.
- `0` tidak boleh dipakai untuk lookup `mst.m_distributor` atau validasi distributor riil.
- `area_id` pada payload principal-only harus tetap menghasilkan target area tersimpan/terlihat sesuai schema existing.
- `emp_id` harus tetap divalidasi dan disimpan sebagai target salesman dalam konteks principal/admin principal (`cust_id`/`parent_cust_id`).
- `outlet_id: []` tetap valid untuk `target_type: "Specific"` bila `emp_id` dipilih.
- Positive distributor IDs tetap mengikuti behavior existing.
- Invalid positive distributor ID tetap gagal dengan error jelas bila flow existing memang memvalidasinya; jangan jadikan semua distributor optional.

## Acceptance Criteria

- Request create survey dengan `distributor_id: [0]`, selected `area_id`, selected `emp_id`, empty `outlet_id`, dan template valid return HTTP success standar (`201`, message `Survey has been successfully created`).
- Survey muncul di list/detail; detail mempertahankan selected area dan salesman.
- Tidak ada query lookup distributor dengan `IN (0)`.
- Tidak ada insert yang memperlakukan `0` sebagai distributor riil. Jika row `mst.m_survey_area.distributor_id = 0` dipakai, itu hanya sebagai sentinel principal-only karena local schema memiliki `NOT NULL` dan migration historis sudah memakai `0`.
- Distributor normal `[123]` tetap sukses.
- Outlet target existing tetap sukses.
- Invalid positive distributor ID tetap menghasilkan validation/error yang jelas.
- Test otomatis ditambahkan untuk principal-only case dan regression normal.

## Existing Patterns/Reuse

- `CreateSurveyBody.DistributorId` sudah memakai `FlexibleIntArray`; tidak perlu ubah parsing JSON.
- `normalizePositiveInts` sudah mengeluarkan `0` dari daftar distributor lookup, sehingga dapat dipertahankan untuk mencegah `FindSurveyAreasByDistributorIds([0])`.
- `resolveSurveyCustIds` sudah return `custId` ketika tidak ada positive distributor IDs; ini cocok untuk principal-only context.
- Transactional write sudah ada di `SurveyService.Store` dan `Update` via `txManager.WithinTransaction`.
- Test doubles di `survey_service_test.go` sudah memadai untuk unit test Red → Green.
- Dokumentasi lokal `docs/Create Survey_BE.md` mencatat `distributor_id` query dengan `0 = principal`, dan migration `002_add_distributor_and_salesman.sql` pernah backfill `distributor_id = 0`, sehingga sentinel `0` punya preseden lokal.

## Constraints

- Ikuti flow Controller → Service → Repository → DB.
- Write tetap dalam transaction service layer.
- Jangan menjalankan test dari repo root karena tidak ada root `go.mod`; jalankan dari `master/`.
- Perintah shell di sesi OpenCode mengikuti instruksi global: tanpa prefix `rtk`, meskipun repo `AGENTS.md` lama meminta `rtk`.
- Jangan mengekspos token/credential dari evidence Jira atau file repo.

## Risks

- Representasi DB untuk principal-only belum eksplisit di spec create survey. Local schema `mst.m_survey_area.distributor_id NOT NULL` mendorong penggunaan sentinel `0` untuk area principal-only.
- Mixed payload `[0, <real_id>]` ambigu. Existing test mengharapkan `0` diabaikan saat ada positive distributor IDs; pertahankan behavior ini sampai domain owner memutuskan lain.
- Staging schema mungkin punya constraint tambahan yang tidak ada di migration lokal. Implementer perlu melihat error aktual/log DB bila local HTTP reproduction tidak sama.

## Decisions/Assumptions

- **Keputusan teknis utama:** Perlakukan `0` sebagai principal-only sentinel hanya ketika tidak ada positive distributor ID dalam request.
- **Create principal-only:** Normalisasi distributor lookup ke empty slice; jangan lookup distributor `0`; validasi salesman terhadap `custId`/`parentCustId`; simpan area target sebagai `m_survey_area` rows dengan `distributor_id = 0` dan selected positive `area_id`, karena schema lokal mewajibkan `distributor_id NOT NULL` dan historical migration memakai `0`.
- **Mixed `[0, 123]`:** Untuk sekarang, abaikan `0` dan proses distributor positif saja, sesuai existing test `TestSurveyService_Store_ShouldIgnoreZeroDistributorAndResolveAreas`.
- **Open question:** Apakah domain owner ingin mixed `[0, <real_id>]` ditolak eksplisit? Jika ya, tambahkan validasi baru dan ubah test existing.

## TDD/Test Plan

- **TDD required:** Ya, karena ini bug fix production logic endpoint create survey.
- **Reason:** Mengubah interpretasi sentinel `distributor_id = 0`, target area, dan target salesman; regression risk tinggi terhadap distributor normal.
- **Existing test patterns:** `master/service/survey_service_test.go` memakai `surveyRepositoryRedStub`, `salesmanRepositoryStub`, dan `transactionManagerStub`.
- **First failing/regression test:** Tambahkan `TestSurveyService_Store_ShouldCreatePrincipalOnlySurveyWithSelectedAreasAndSalesman`:
  - Request `DistributorId: {0}`, `AreaId: {89,86,85}`, `EmpId: {369}`, `TargetType: "Specific"`, valid date/template/cust.
  - Stub salesman `369` valid.
  - Expect no error.
  - Expect `repo.findAreasInput` dan `repo.findCustIdsInput` kosong atau tidak berisi `0`.
  - Expect `repo.storeAreasInput` berisi rows `{DistributorId: 0, AreaId: 89/86/85}`.
  - Expect `repo.storeSalesmenInput` berisi salesman `369` dengan cust principal request.
- **Green step:** Ubah helper build area agar principal-only request dengan selected areas menghasilkan sentinel area rows tanpa distributor lookup; pastikan mixed flow tetap pakai existing positive distributor mapping.
- **Refactor step:** Ekstrak helper kecil bila perlu, misalnya `hasPrincipalOnlyDistributorSentinel(values []int)`, `buildPrincipalSurveyAreas(areaIds []int)`, atau `isPrincipalOnlyRequest(raw, normalized []int)` untuk menjaga `Store`/`Update` mudah dibaca.
- **Edge cases:**
  - `distributor_id: [0]` + `area_id: []` + `emp_id` valid: survey dan salesman tetap tersimpan, tidak ada area rows.
  - `distributor_id: [0]` + invalid `emp_id`: tetap fail `ErrSurveySalesmanNotFound`.
  - `distributor_id: [0, 123]`: existing behavior tetap ignore `0`, proses `123`.
  - `distributor_id: [999999999]`: positive invalid tetap tidak silently menjadi principal-only.
  - Update survey principal-only bila endpoint update menerima payload sejenis.
- **Commands:**
  - `go test ./service -run 'TestSurveyService_Store|TestSurveyService_Update'`
  - `go test ./...`
  - Optional after service up: authenticated `curl`/Postman against local/staging-like endpoint with sanitized payload and valid token.

## Implementation Steps

1. Reproduce unit failure by adding the principal-only test in `master/service/survey_service_test.go` before production changes.
2. Inspect whether `buildSurveyAreas` is used by both `Store` and `Update`; plan to keep behavior consistent.
3. Add helper to detect principal-only sentinel:
   - `rawDistributorIds` contains `0`.
   - `normalizedDistributorIds` length is `0`.
4. Adjust create/update area building:
   - For normal/mixed with positive distributors: keep existing `FindSurveyAreasByDistributorIds(normalizedDistributorIds)` and `buildSurveyAreas` behavior.
   - For principal-only: do not call distributor lookup with `0`; build area rows directly from normalized positive `request.AreaId` using `DistributorId: 0`.
   - If `area_id` empty, return no area rows without `ErrSurveyAreaDistributorRequired`.
5. Ensure salesman cust resolution for principal-only uses `resolveSurveyCustIds` with empty positive distributor IDs and therefore returns request `CustId`.
6. Preserve mixed test behavior: `[67,0,68]` should still lookup `[67,68]` and must not insert `DistributorId: 0`.
7. If positive distributor ID has selected `area_id` but `FindSurveyAreasByDistributorIds` returns no matching area, verify current intended behavior; if currently missing explicit invalid distributor error, consider adding a targeted validation only if tests/domain support it.
8. Run focused tests, then full `master` module tests.
9. If authenticated environment is available, capture actual HTTP status/body for the sanitized Jira-like payload and verify list/detail.

## Expected Files to Change

- `master/service/survey_service.go`
- `master/service/survey_service_test.go`
- Possibly `master/repository/survey_repository.go` only if implementation chooses repository validation/helper for area existence.
- No package/lock/config changes expected.

## Agent/Tool Routing

- Implementation: route to `@fixer` or build agent using `opencode-fixer` / `opencode-build`.
- Research/code search if needed: `@explorer` only.
- Architecture review if deciding mixed sentinel policy or DB representation changes: `@oracle`.
- No UI/designer/browser work required except optional API verification.

## Validation Commands

Run from `master/`:

```bash
go test ./service -run 'TestSurveyService_Store|TestSurveyService_Update'
go test ./...
```

Optional runtime checks:

```bash
docker compose -f ../docker-compose.yml ps
curl http://localhost:9002/ping
```

Authenticated API retest requires a valid non-sanitized token:

```bash
curl 'http://localhost:9002/master/v1/survey' \
  -H 'Accept: application/json, text/plain, */*' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <valid-token>' \
  --data-raw '{"survey_title":"Testing Survey Principal","efective_date_start":"2026-03-05","efective_date_end":"2026-05-05","answer_frequency":"One Time","response_type":"Mandatory","target_type":"Specific","distributor_id":[0],"area_id":[89,86,85,84,83,82,70],"outlet_id":[],"survey_template_id":63,"emp_id":[369]}'
```

## Evidence Requirements

- Required local evidence:
  - Output of failing test before fix.
  - Output of focused passing test after fix.
  - Output of `go test ./...` after fix.
- Required API evidence if token/environment available:
  - HTTP status and response body for principal-only create.
  - Detail/list verification showing area and salesman persisted/resolved.
  - Log/SQL evidence confirming no distributor lookup with `0`.
- Official docs/context7 not needed; behavior is local code/domain-specific.
- GitHub/web search not needed; no upstream dependency decision.
- Browser/screenshot not needed; backend-only defect.

## Done Criteria

- Principal-only survey create test exists and passes.
- Normal distributor create/update tests still pass.
- Mixed `[0, positive]` existing regression still passes or is intentionally changed with domain approval.
- Full `master` tests pass or failures are documented as unrelated with evidence.
- Implementation summary/PR documents root cause: `buildSurveyAreas` required positive distributor IDs when `area_id` was present, so principal-only `[0]` normalized to empty and failed or lost target area.
- QA can retest SX-1906 on staging with same scenario and receive success.

## Final Planning Summary

- Artifacts created:
  - `.opencode/plans/20260504-0846-sx-1906-survey-principal-only.md` — source of truth for implementation.
  - `.opencode/evidence/20260504-0846-sx-1906-survey-principal-only/discovery.md` — kept because it records inspected files, local schema evidence, and risks useful for implementation.
- Key decisions:
  - Treat `0` as principal-only sentinel only when no positive distributor IDs exist.
  - Use `distributor_id = 0` in `mst.m_survey_area` only as sentinel representation for principal-only selected areas, not as distributor lookup/FK.
  - Preserve existing mixed behavior that ignores `0` when positive distributors are present.
- Assumptions/questions:
  - Assumed local schema/migration is authoritative enough to store sentinel `0` in survey area rows.
  - Remaining domain question: whether mixed `[0, <real_id>]` should continue to be allowed/ignored or become validation error.
- Readiness:
  - Ready for implementation with TDD. No user question blocks the fix because existing docs/tests provide a low-risk default for mixed behavior and schema sentinel.
- Cleanup performed:
  - Discovery evidence was not deleted because it remains operationally useful for the implementer.
