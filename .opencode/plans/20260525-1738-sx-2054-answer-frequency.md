# Plan — SX-2054 Answer Frequency Survey

Task ID: `20260525-1738-sx-2054-answer-frequency`

Primary source of truth: `.opencode/plans/20260525-1738-sx-2054-answer-frequency.md`

## Goal

Membuat backend survey siap menerima, menyimpan, dan membaca `answer_frequency` baru untuk `mst.m_survey` pada `POST /master/v1/survey` dan `PUT /master/v1/survey/:survey_id`, tanpa merusak row legacy bernilai `Multiple`.

## Non-goals

- Tidak melakukan mass migration dari `Multiple` ke salah satu value baru.
- Tidak mengubah label UI `Same day` / `Different Day`; backend memakai value API exact.
- Tidak mengarang rule submit baru untuk `One Day` vs `Different Day` tanpa requirement product.
- Tidak menaruh Authorization, cookie, token, credential, atau cURL sensitif ke source/test/artifact.
- Tidak memperbaiki risiko sort injection survey list dalam task ini kecuali diminta terpisah.

## Scope

Module utama: `master`.

Endpoint utama:

- `POST /master/v1/survey`
- `PUT /master/v1/survey/:survey_id`

Area wajib:

- schema `mst.m_survey.answer_frequency`
- request validation create/update
- persistence create/update
- read path list/detail/report yang expose `answer_frequency`
- audit downstream `mobile` yang baca `mst.m_survey.answer_frequency`
- test controller/service sesuai pattern repo

## Requirements

1. New write hanya menerima value exact:
   - `One Time`
   - `Multiple Times, One Day`
   - `Multiple Times, Different Day`
2. Legacy value `Multiple` tetap bisa dibaca pada list/detail/report/mobile selama ada row lama.
3. New write tidak boleh normalize dua value baru menjadi `Multiple`.
4. Schema harus cukup panjang untuk value baru terpanjang.
5. Create dan update menyimpan string exact dari payload valid.
6. Read path mengembalikan string exact dari DB.
7. Semua literal `Multiple` yang relevan di test/code diaudit; happy-path write baru tidak boleh memakai `Multiple`.
8. Mobile submit behavior diaudit. Jika perlu behavior beda untuk `One Day` vs `Different Day`, wajib keputusan product sebelum implementasi.

## Acceptance Criteria

- `mst.m_survey.answer_frequency` bisa menyimpan `Multiple Times, One Day` dan `Multiple Times, Different Day`.
- `POST /master/v1/survey` menerima dan menyimpan `One Time`.
- `POST /master/v1/survey` menerima dan menyimpan `Multiple Times, One Day`.
- `POST /master/v1/survey` menerima dan menyimpan `Multiple Times, Different Day`.
- `PUT /master/v1/survey/:survey_id` menerima dan menyimpan tiga value baru.
- `Multiple` ditolak untuk new create/update, kecuali user/product memutuskan temporary write compatibility.
- List/detail/report/mobile read path tetap return raw DB value, termasuk legacy `Multiple`.
- Tidak ada mapping paksa value baru menjadi `Multiple`.
- Test create/update validation dan exact persistence lewat.
- Legacy strategy tertulis di final summary implementasi.
- Tidak ada credential/token baru di source/test/commit.

## Existing Patterns/Reuse

- `master` memakai Fiber dan layer `Controller → Service → Repository → DB`.
- `Store()` dan `Update()` sudah memakai transaction via `txManager.WithinTransaction(...)`.
- `survey_repository.go` sudah pass-through `AnswerFrequency` ke `INSERT`/`UPDATE`.
- `List()` dan `Detail()` di service sudah pass-through raw `AnswerFrequency` dari repository.
- `survey_report_repository.go` dan `mobile/repository/survey.go` juga select raw `s.answer_frequency` / `ms.answer_frequency`.
- Test pattern sudah ada di `master/controller/survey_controller_test.go` dan `master/service/survey_service_test.go`.
- Migration survey memakai raw SQL berurutan di `master/migration/mst.survey/`.
- KiloCode/project utility khusus tidak ditemukan untuk enum survey; gunakan pattern repo sendiri.

## Constraints

- `answer_frequency` saat ini `VARCHAR(20)`; value baru panjang 23 dan 29 karakter.
- Tidak ada `CHECK` constraint lokal untuk `mst.m_survey.answer_frequency`.
- Data aktual bisa berisi value lain karena DB lama permisif; audit distinct value perlu sebelum constraint keras.
- Validation package memakai `go-playground/validator/v10`.
- Inline `oneof` untuk value berkoma dan spasi rapuh; rencana pakai custom validation tag.
- Perintah shell repo pakai `rtk` sesuai AGENTS lokal.
- Validasi dilakukan dari module `master/`, bukan repo root.

## Risks

1. Release app tanpa migration widen akan gagal insert/update karena `VARCHAR(20)`.
2. Menambah `CHECK` langsung tanpa data audit bisa memecah row existing yang punya value liar.
3. Menolak `Multiple` untuk new write bisa memecah FE lama bila FE belum deploy bersamaan.
4. Mengizinkan `Multiple` untuk new write bisa melanggar requirement SX-2054.
5. Mobile submit sekarang memblokir submit kedua untuk semua survey; jika product berharap `Multiple Times` mengubah duplicate prevention, scope masih kurang.
6. Controller/service tests existing memakai `Multiple` sebagai happy path; update harus hati-hati agar test lama tetap bermakna atau diubah menjadi invalid/legacy-read case.
7. `FindAllByCustId` punya raw `ORDER BY` dari `filter.Sort`; di luar scope tapi perlu dicatat bila touched.

## Decisions/Assumptions

- **Keputusan:** value API/persistence exact adalah `One Time`, `Multiple Times, One Day`, `Multiple Times, Different Day`.
- **Keputusan:** read path tetap permissive dan raw untuk legacy `Multiple`.
- **Keputusan:** migration aman pertama adalah widen `answer_frequency` ke `VARCHAR(50) NOT NULL`.
- **Keputusan:** jangan bulk-update row `Multiple`.
- **Keputusan:** jangan pasang final `CHECK` yang menolak `Multiple` pada task ini.
- **Asumsi:** new create/update menolak `Multiple` karena requirement menyebut value baru untuk write.
- **Asumsi:** `mobile` behavior submit tidak diubah sampai ada rule product eksplisit.
- **Open question:** apakah product mau temporary write compatibility untuk `Multiple` selama FE rollout?
- **Open question:** apakah `Multiple Times, One Day` dan `Multiple Times, Different Day` harus mengubah duplicate-prevention mobile sekarang atau task lanjutan?

## TDD/Test Plan

- **TDD required:** ya.
- **Alasan:** task menyentuh kontrak API, validation, DB schema, persistence, dan legacy compatibility.
- **Existing test patterns:** gunakan controller tests di `master/controller/survey_controller_test.go`; service tests dan stubs di `master/service/survey_service_test.go`; bila repository test tersedia, pakai `sqlmock`, tapi saat discovery tidak ada repository test khusus `survey_repository.go`.

### Red step

1. Tambah controller validation test create yang membuktikan `Multiple Times, One Day` dan `Multiple Times, Different Day` awalnya ditolak oleh tag lama.
2. Tambah controller validation test update yang membuktikan dua value baru awalnya ditolak.
3. Tambah test yang memastikan `Multiple` ditolak untuk create/update baru.
4. Tambah service test/stub capture agar `Store()` dan `Update()` menyimpan exact value, bukan normalize.
5. Tambah read-path regression service `Detail()` dan `List()` untuk memastikan legacy `Multiple` tetap keluar raw.
6. Tambah migration-review check manual: nilai baru lebih panjang dari `VARCHAR(20)` sehingga butuh migration.

### Green step

1. Tambah constants/helper terpusat, contoh lokasi yang disarankan: `master/entity/survey.go` atau file baru kecil `master/entity/survey_answer_frequency.go`.
2. Define constants:
   - `AnswerFrequencyOneTime = "One Time"`
   - `AnswerFrequencyMultipleTimesOneDay = "Multiple Times, One Day"`
   - `AnswerFrequencyMultipleTimesDifferentDay = "Multiple Times, Different Day"`
   - `AnswerFrequencyLegacyMultiple = "Multiple"` hanya untuk read/legacy/test.
3. Define helper:
   - `IsValidSurveyAnswerFrequencyForWrite(value string) bool`
   - opsional `IsLegacySurveyAnswerFrequency(value string) bool`
   - opsional `IsMultipleSurveyAnswerFrequency(value string) bool` yang mengenali dua value baru dan legacy `Multiple` bila dipakai untuk backward-compatible read logic.
4. Replace inline tag `oneof=Multiple 'One Time'` dengan custom validation tag, misalnya `answer_frequency`, pada create/update request.
5. Register custom validation di `master/pkg/validation/validation.go` memakai `RegisterValidation`, sesuai documented `go-playground/validator/v10` pattern.
6. Tambah custom translation singkat untuk `answer_frequency` bila error default tidak cukup jelas.
7. Tambah migration baru di `master/migration/mst.survey/004_alter_answer_frequency_length.sql`:
   ```sql
   ALTER TABLE mst.m_survey
   ALTER COLUMN answer_frequency TYPE VARCHAR(50);
   ```
8. Jangan tambah data update massal.
9. Jangan tambah strict `CHECK` sebelum distinct-value audit.
10. Update tests agar happy path baru memakai tiga value requirement, bukan `Multiple`.

### Refactor step

- Hilangkan duplikasi value literal di tests sebanyak mungkin dengan constants bila import cycle tidak terjadi.
- Keep read mapping pass-through; jangan buat mapper baru yang mengubah string.
- Bila ada helper `IsMultipleSurveyAnswerFrequency`, pakai hanya untuk grouping behavior yang benar-benar ada; jangan pakai untuk persistence normalization.

### Edge cases

- Empty `answer_frequency` tetap rejected by `required`.
- Unknown typo/case mismatch rejected.
- Legacy `Multiple` rejected pada create/update baru.
- Legacy `Multiple` tetap bisa dibaca pada detail/list/report/mobile.
- New values dengan koma/spasi harus diterima exact.
- `Multiple Times, Different Day` harus lolos DB length.

### Commands

Dari `master/`:

```bash
rtk go test ./controller -run 'TestSurveyController_.*AnswerFrequency|TestSurveyController_Create|TestSurveyController_Update' -v
rtk go test ./service -run 'TestSurveyService_.*AnswerFrequency|TestSurveyService_Store|TestSurveyService_Update|TestSurveyService_Detail|TestSurveyService_List' -v
rtk go test ./repository -run 'TestSurveyReport' -v
rtk go test ./...
```

Dari repo root untuk runtime evidence:

```bash
rtk docker compose -f docker-compose.yml ps
rtk docker compose -f docker-compose.yml up -d
```

## Implementation Steps

1. Jalankan preflight grep lagi di `master` dan `mobile` untuk memastikan tidak ada usage baru:
   ```bash
   rg -n "answer_frequency|AnswerFrequency|\bMultiple\b|One Time|Multiple Times" master mobile
   ```
2. Tambah Red tests untuk create/update validation, exact persistence, dan legacy read compatibility.
3. Tambah constants/helper answer frequency di layer yang tidak membuat import cycle.
4. Register custom validator `answer_frequency` di `master/pkg/validation/validation.go`.
5. Ganti tag pada `CreateSurveyBody.AnswerFrequency` dan `UpdateSurveyBody.AnswerFrequency` menjadi `validate:"required,answer_frequency"`.
6. Tambah migration widen column ke `VARCHAR(50)` di `master/migration/mst.survey/004_alter_answer_frequency_length.sql`.
7. Pastikan `Store()` dan `Update()` tetap pass-through `request.AnswerFrequency` exact.
8. Audit dan update tests yang masih memakai `Multiple` sebagai valid write; ubah ke value baru atau jadikan invalid-write/legacy-read test.
9. Pastikan list/detail/report/mobile tetap raw pass-through; jangan ubah query read kecuali test menunjukkan normalisasi tersembunyi.
10. Jalankan focused tests, lalu full `master` tests.
11. Jika DB tersedia, jalankan preflight distinct-value audit dan smoke create/update/read manual.
12. Catat mobile duplicate-submission gap sebagai known impact bila belum ada rule product.

## Expected Files to Change

Wajib:

- `master/entity/survey.go`
- `master/pkg/validation/validation.go`
- `master/migration/mst.survey/004_alter_answer_frequency_length.sql`
- `master/controller/survey_controller_test.go`
- `master/service/survey_service_test.go`

Opsional bila helper dipisah:

- `master/entity/survey_answer_frequency.go`
- `master/pkg/validation/validation_survey.go`

Kemungkinan tidak berubah karena sudah pass-through:

- `master/service/survey_service.go`
- `master/repository/survey_repository.go`
- `master/repository/survey_report_repository.go`
- `mobile/repository/survey.go`
- `mobile/service/survey.go`

## Agent/Tool Routing

- `@orchestrator`: jalankan handoff, koordinasi implementation dan validation.
- `@fixer`: implementasi perubahan code/migration/tests di `master`.
- `@explorer`: discovery tambahan bila grep menemukan usage baru atau docs internal tersebar.
- `@oracle`: review bila akan menambah DB `CHECK` atau mengubah mobile submit semantics.
- `@quality-gate`: final signoff karena task menyentuh API contract, schema, legacy data, dan mobile impact.
- `@librarian/context7`: sudah dipakai untuk verifikasi pola `go-playground/validator/v10` custom validation.

## Execution-ready Worklist / Handoff Contract

`start_with`: `SX2054-01`

| Task ID | Action | depends_on | Owner/lane | Validation/check | Exit criteria | Status | Blocker | requires_user_decision |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| SX2054-01 | Re-run grep audit untuk `answer_frequency`, `Multiple`, `One Time`, `Multiple Times` di `master mobile` | none | `@explorer` atau `@fixer` | `rg -n "answer_frequency|AnswerFrequency|\bMultiple\b|One Time|Multiple Times" master mobile` | Semua titik relevan dicatat; tidak ada usage tersembunyi sebelum edit | ready | none | no |
| SX2054-02 | Tulis Red tests create/update validation untuk tiga value baru dan reject `Multiple` | SX2054-01 | `@fixer` | `rtk go test ./controller -run 'AnswerFrequency|Create|Update' -v` dari `master/` harus gagal sebelum fix | Tests gagal karena validator lama | ready | none | no |
| SX2054-03 | Tulis Red tests exact persistence dan legacy read raw | SX2054-01 | `@fixer` | `rtk go test ./service -run 'AnswerFrequency|Detail|List|Store|Update' -v` harus gagal untuk new contract bila belum fix | Tests menunjukkan expected new contract | ready | none | no |
| SX2054-04 | Tambah constants/helper answer frequency tanpa import cycle | SX2054-02,SX2054-03 | `@fixer` | `rtk go test ./entity ./pkg/validation` bila package valid | Constants/helper tersedia dan tests compile | ready | none | no |
| SX2054-05 | Register custom validator `answer_frequency` dan update create/update tags | SX2054-04 | `@fixer` | Controller tests pass for allowed values, fail for `Multiple`/unknown | API boundary menerima exact new values dan reject invalid | ready | none | no |
| SX2054-06 | Tambah migration widen `mst.m_survey.answer_frequency` ke `VARCHAR(50)` | SX2054-01 | `@fixer` | Static SQL review; DB smoke bila tersedia | Migration ada, no DML, no mass conversion | ready | none | no |
| SX2054-07 | Update old happy-path tests yang memakai `Multiple` sebagai valid write | SX2054-05 | `@fixer` | `rtk go test ./controller ./service -v` | Tests selaras new write contract | ready | none | no |
| SX2054-08 | Audit read paths list/detail/report/mobile tetap raw pass-through | SX2054-05,SX2054-06 | `@fixer` | grep + targeted tests/manual inspection | Tidak ada normalisasi ke `Multiple`; legacy read aman | ready | none | no |
| SX2054-09 | Jalankan focused tests dan full `master` tests | SX2054-07,SX2054-08 | `@fixer` | `rtk go test ./...` dari `master/` | Semua test target lewat atau failure terdokumentasi bukan dari SX-2054 | ready | none | no |
| SX2054-10 | DB preflight distinct-value audit dan manual smoke create/update/read bila DB tersedia | SX2054-06,SX2054-09 | `@orchestrator`/`@fixer` | SQL distinct + API/DB smoke | Evidence manual berisi stored exact values dan legacy read status | ready | none jika DB tersedia | no |
| SX2054-11 | Putuskan mobile submit semantics untuk `One Day` vs `Different Day` | SX2054-01 | `@architect`/`@oracle` + product | Product rule tertulis | Behavior rule jelas atau task lanjutan dibuat | blocked | Requirement eksplisit belum ada | yes |
| SX2054-12 | Final quality gate | SX2054-09,SX2054-10 | `@quality-gate` | Review diff, tests, migration, evidence | Gate pass atau issue tercatat | ready | none setelah validation | no |

## Validation Commands

Dari repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

Dari `master/`:

```bash
rtk go mod download && rtk go mod tidy
rtk go test ./controller -run 'TestSurveyController_.*AnswerFrequency|TestSurveyController_Create|TestSurveyController_Update' -v
rtk go test ./service -run 'TestSurveyService_.*AnswerFrequency|TestSurveyService_Store|TestSurveyService_Update|TestSurveyService_Detail|TestSurveyService_List' -v
rtk go test ./repository -run 'TestSurveyReport' -v
rtk go test ./...
```

DB preflight bila environment siap:

```sql
SELECT answer_frequency, COUNT(*)
FROM mst.m_survey
GROUP BY 1
ORDER BY 1;
```

Manual API smoke bila service dan auth dev tersedia:

1. Create survey dengan `answer_frequency = "One Time"`.
2. Create survey dengan `answer_frequency = "Multiple Times, One Day"`.
3. Create survey dengan `answer_frequency = "Multiple Times, Different Day"`.
4. Update existing survey ke tiap value baru.
5. Baca list/detail dan pastikan output sama dengan DB.
6. Baca row legacy `Multiple` bila ada dan pastikan tidak error.

## Evidence Requirements

- Test output focused dan full `master`.
- Migration file path dan SQL content.
- Grep audit setelah implementasi untuk `answer_frequency` dan literal `Multiple`.
- DB distinct-value audit bila DB tersedia.
- Manual create/update/read smoke evidence bila auth/runtime tersedia.
- Catatan mobile submit gap bila belum ada keputusan product.

## Done Criteria

- Migration widen tersedia dan aman.
- Create/update validation menerima tiga value baru exact.
- Create/update validation tidak menerima `Multiple` untuk new write, kecuali ada keputusan eksplisit berbeda.
- Persistence menyimpan value exact.
- Read path tidak normalize dan tetap legacy-compatible.
- Tests relevan pass.
- Manual verification tercatat atau blocker runtime jelas.
- Open mobile semantics ditutup oleh keputusan product atau tercatat sebagai follow-up.
- Quality gate pass.

## Research Gate

- **Local project discovery:** dilakukan; wajib karena task non-trivial dan menyentuh schema/API.
- **Official docs/context7:** dilakukan untuk `go-playground/validator/v10` custom validation karena inline `oneof` dengan value berkoma/spasi rapuh.
- **GitHub:** tidak diperlukan; perubahan berdasar repo lokal dan library docs cukup.
- **Brave/web search:** tidak diperlukan; source of truth user/Jira sudah diberikan di prompt, tidak perlu fakta eksternal.
- **Browser/screenshot:** tidak diperlukan; task backend API/schema, bukan UI visual.

## Final Planning Summary

Artifacts dibuat dan dipakai:

- Primary plan: `.opencode/plans/20260525-1738-sx-2054-answer-frequency.md`
- Discovery evidence kept: `.opencode/evidence/20260525-1738-sx-2054-answer-frequency/discovery.md`
- Open questions kept: `.opencode/draft/20260525-1738-sx-2054-answer-frequency/open-questions.md`

Key decisions:

- Widen `answer_frequency` ke `VARCHAR(50)`.
- New writes menerima tiga value baru exact.
- Legacy `Multiple` tetap raw-readable, tidak dimigrasi massal.
- Read path tidak dinormalisasi.
- Mobile submit semantics tidak diubah tanpa product rule.

Assumptions:

- FE sudah atau akan mengirim exact API values dari comment BE/API.
- New write boleh menolak `Multiple`.
- DB actual dapat diaudit sebelum constraint tambahan.

Remaining open questions:

- Apakah `Multiple` perlu temporary write compatibility selama FE rollout?
- Apakah `Multiple Times, One Day` dan `Multiple Times, Different Day` harus mengubah duplicate-prevention mobile di task ini?

Readiness:

- Implementasi create/update/schema siap dijalankan mulai `SX2054-01`.
- Mobile submit behavior task `SX2054-11` blocked sampai product decision.

Cleanup performed:

- Draft/evidence tidak dihapus karena masih operasional: discovery evidence membantu implementasi, open questions masih unresolved.
