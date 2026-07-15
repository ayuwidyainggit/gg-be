# Plan — SX-1903 Survey Report Attachment Links

## Goal

Perbaiki export `DownloadSurveyReport-*` agar kolom `Attachment 1`, `Attachment 2`, dan `Attachment 3` berisi link attachment yang valid, bukan `file_key`/filename mentah, dan cell Excel dibuat sebagai hyperlink ke public OBS URL sesuai pola storage existing.

## Non-goals

- Tidak mengubah flow menu frontend `Report → Survey list` atau `Download History`.
- Tidak mengubah struktur header/report non-attachment.
- Tidak membuat storage private, signed URL, atau endpoint download authorized baru dalam scope bugfix ini.
- Tidak memigrasi data `mst.survey_answer_files`.
- Tidak memasukkan credential admin/Jira/SOP ke source, test, fixture, atau artifact.

## Scope

- Modul target: `master`.
- Flow utama:
  - `master/controller/survey_report_controller.go` tetap memanggil `SurveyReportService.Export`.
  - `master/service/survey_report_service.go` melakukan enrichment `Attachment1/2/3` menjadi URL dan menulis hyperlink Excel.
  - `master/repository/survey_report_repository.go` tetap mengambil `file_key` dengan tenant scope.
  - `master/adapter/obs-huawei.go` diperluas agar bisa resolve object key menjadi public URL berdasarkan `FileBaseUrl` existing.
  - `master/main.go` meng-inject resolver ke `NewSurveyReportService`.
- Test difokuskan pada resolver, workbook hyperlink/cell value, service export mapping, dan query regression.

## Requirements

- Attachment kosong/null tetap menghasilkan cell kosong tanpa error.
- Attachment object key seperti `d7o3c10estqgiqkueq50.jpg` atau `folder/file.jpg` harus menjadi URL `https://<bucket>.<endpoint>/<encoded-key>`.
- Attachment yang sudah berupa full URL hanya boleh dipertahankan jika host sama dengan configured OBS public host.
- Attachment full URL external harus tidak diekspor sebagai link external; pilih perilaku aman berupa cell kosong atau error terkontrol sesuai pola service yang paling mudah dites. Rekomendasi: cell kosong untuk mencegah satu data buruk menggagalkan seluruh export, plus test behavior.
- Excel cell attachment harus memiliki hyperlink target URL valid. Display text boleh filename/object key atau URL; keputusan user: gunakan hyperlink. Rekomendasi implementasi: display filename/key agar readable, hyperlink target full URL.
- Jika Excelize hyperlink gagal, export harus mengembalikan error dan `report_list` di-update `FAILED` seperti error generation lainnya.
- Tidak double-prefix untuk URL valid.
- Tidak hardcode bucket/domain; gunakan `OBS_HUAWEI_BUCKET` dan `OBS_HUAWEI_ENDPOINT` melalui adapter yang sudah diinisialisasi.

## Acceptance Criteria

1. Export survey report untuk filter evidence menghasilkan hyperlink attachment valid.
2. Kolom attachment tidak lagi hanya berisi filename mentah tanpa link.
3. Link mengarah ke public OBS URL sesuai configured bucket/endpoint dan host tervalidasi.
4. Data non-attachment, header, filter, urutan row, dan Download History tetap normal.
5. Export tanpa attachment tetap berhasil.
6. Workbook `.xlsx` tetap valid dan cell attachment memiliki hyperlink relationship.
7. Automated test ditambahkan/diupdate untuk attachment link generation dan edge cases.

## Existing Patterns/Reuse

- Reuse pattern upload OBS di `master/adapter/obs-huawei.go`:
  - `FileBaseUrl = https://<bucket>.<endpoint>`.
  - `UploadFile` mengembalikan `fmt.Sprintf("%v/%v", o.FileBaseUrl, key)`.
- Reuse `excelize/v2` yang sudah dipakai di `generateSurveyReportExcel`.
- Reuse `SurveyReportExportRow.Attachment1/2/3` sebagai intermediate object key dari repository.
- Reuse `survey_report_service_test.go` dan tambahkan inspeksi workbook, bukan hanya cek base64 signature.
- Reuse `survey_report_repository_test.go` untuk assertion query string terkait tenant scope dan attachment ordering.

## Constraints

- Repository ini multi-module Go; jalankan command validasi dari `master/`.
- Project `AGENTS.md` meminta shell command memakai `rtk`; ikuti pola repo saat implementasi/validasi.
- Jangan expose atau menyalin secret dari compose/workflow/env.
- Strict layer flow tetap dijaga: repository query DB, service mapping/orchestration, adapter storage URL.
- Query attachment harus tetap memakai `saf.cust_id = sa.cust_id`.
- TDD wajib karena ini bugfix logic/export behavior.

## Risks

- Public OBS URL bersifat long-lived dan bisa dibagikan; ini sesuai keputusan user untuk bugfix minimal, tetapi bukan kontrol akses kuat untuk attachment sensitif.
- Jika bucket/ACL berubah private di masa depan, link public OBS tidak akan cukup; perlu authorized download endpoint atau URL generation saat user membuka file.
- Signed URL tidak dipilih karena workbook tersimpan di Download History dan URL bisa expired.
- Jika data lama menyimpan URL external, resolver harus memblokir agar export tidak menjadi vektor phishing.
- Jika `file_key` memiliki karakter khusus, resolver harus encode path segment dengan benar tanpa mengubah `/` sebagai separator folder.

## Decisions/Assumptions

- Pertanyaan sudah diajukan dan dijawab:
  - Tipe link: `Public OBS URL`.
  - Format Excel: `Excel hyperlink`.
  - Full URL lama: `Validate same host`.
- Asumsi implementasi:
  - Object di OBS tetap public-readable sesuai `obs.AclPublicRead` existing.
  - `file_key` adalah object key yang cukup untuk membangun URL deterministic.
  - Attachment export tetap maksimal 3 kolom sesuai header existing.
  - Display text hyperlink boleh menggunakan basename/object key; target harus full URL.
- Open questions tersisa:
  - Perlu konfirmasi schema apakah `mst.survey_answer_files` punya `is_del`. Jika ada, implementation agent sebaiknya menambah filter `saf.is_del = false` setelah memastikan tidak memutus data valid.

## TDD/Test Plan

- TDD required: ya, karena ini bugfix behavior export dan file URL security-sensitive.
- Reason: bug saat ini dapat dibuktikan dengan workbook yang hanya berisi `file_key` tanpa hyperlink.
- Existing test patterns:
  - `master/service/survey_report_service_test.go` untuk helper/service export workbook.
  - `master/repository/survey_report_repository_test.go` untuk query string behavior.
  - `github.com/DATA-DOG/go-sqlmock` sudah tersedia untuk repository tests.
- First failing/regression test:
  1. Tambahkan test di `survey_report_service_test.go` yang membuat row dengan `Attachment1: "survey/C26004/file a.jpg"`, memakai fake resolver/base URL, generate workbook, decode base64, buka dengan Excelize, lalu assert:
     - cell `M2` display text bukan kosong,
     - hyperlink `M2` target adalah `https://<bucket>.<endpoint>/survey/C26004/file%20a.jpg`,
     - `N2` kosong tanpa hyperlink saat attachment kosong,
     - `O2` sesuai attachment ketiga jika ada.
  2. Current implementation akan gagal karena `generateSurveyReportExcel` tidak membuat hyperlink dan hanya menulis raw string.
- Green step:
  - Tambahkan resolver URL dan transform/hyperlink writing sampai test workbook lulus.
- Refactor step:
  - Pisahkan helper kecil seperti `resolveSurveyReportAttachmentLink` dan `setSurveyReportAttachmentCell` agar edge cases mudah dites.
  - Pastikan service constructor tetap jelas dan dependency sempit.
- Edge cases:
  - Empty string → empty cell, no hyperlink.
  - Normal key → encoded OBS URL.
  - Leading slash → tidak double slash.
  - Subfolder path → `/` tetap separator, segment lain ter-encode.
  - Already full URL same OBS host → preserve/normalize target, display basename/key.
  - Full URL external host → no external hyperlink emitted.
  - Multiple attachments → mapped to `Attachment 1..3` in order.
- Commands:
  - Dari `master/`: `rtk go test ./service -run 'TestGenerateSurveyReportExcel|TestSurveyReportAttachment'`
  - Dari `master/`: `rtk go test ./repository -run TestFindExportRows`
  - Dari `master/`: `rtk go test ./...`

## Implementation Steps

1. **Red test workbook hyperlink**
   - Update/add test in `master/service/survey_report_service_test.go` that decodes generated base64 workbook and verifies hyperlink target/cell values for attachment columns.
   - Use a fake URL resolver or helper-level test depending on chosen function signature.

2. **Add narrow URL resolver contract**
   - Define a small interface in service or adapter boundary, e.g. `type ObjectURLResolver interface { ResolveObjectURL(key string) string }` or `(string, bool)` if invalid URL needs explicit handling.
   - Prefer `(string, bool)` or `(string, error)` internally to handle external URL rejection deterministically.

3. **Extend OBS adapter**
   - In `master/adapter/obs-huawei.go`, implement resolver using existing `FileBaseUrl`.
   - Normalize:
     - trim whitespace,
     - return empty for empty key,
     - remove leading `/`,
     - if input is full URL, allow only when host matches `FileBaseUrl` host,
     - encode path segments while preserving `/`.
   - Avoid network calls to OBS.

4. **Inject resolver into survey report service**
   - Change `NewSurveyReportService(surveyReportRepo, reportListRepo)` to include resolver dependency.
   - Update `master/main.go` to pass `obsAdapter`.
   - Update affected tests/fakes for constructor changes.

5. **Enrich export rows or Excel writer**
   - Recommended structure:
     - `Export` fetches rows.
     - `generateSurveyReportExcel(rows, resolver)` or a pre-step transforms attachment key into a struct containing `Display` and `URL`.
   - Since hyperlinks require both display and target, prefer a small internal export-cell model rather than overwriting raw `Attachment1` with URL only.

6. **Write hyperlink with Excelize**
   - For each attachment cell `M/N/O`:
     - set display text, e.g. original `file_key` basename or original string,
     - call Excelize hyperlink API for external link target,
     - optionally apply hyperlink style (blue + underline) if supported and consistent.
   - Keep non-attachment fields unchanged.

7. **Add resolver security tests**
   - Test empty, leading slash, special chars, same-host full URL, external full URL.

8. **Repository query regression**
   - Add/adjust test to assert export attachment subqueries keep `saf.cust_id = sa.cust_id` and `ORDER BY saf.survey_answer_files`.
   - If schema confirms `is_del`, add filter and test it.

9. **Manual/integration validation**
   - Run export with Jira filter if environment data/credentials are available through SOP.
   - Download from Download History and verify the `Survey Report` sheet attachment cell displays as hyperlink and opens public OBS URL.

## Expected Files to Change

- `master/adapter/obs-huawei.go`
  - Add URL resolver method and path/host normalization helper.
- `master/service/survey_report_service.go`
  - Inject resolver and generate attachment hyperlinks.
- `master/main.go`
  - Pass `obsAdapter` into `NewSurveyReportService`.
- `master/service/survey_report_service_test.go`
  - Add workbook/hyperlink tests and resolver fake.
- `master/repository/survey_report_repository_test.go`
  - Add export query regression assertions.
- Possibly `master/repository/survey_report_repository.go`
  - Only if schema confirms deleted-file filtering is needed.

## Agent/Tool Routing

- Implementation: route to build/fixer implementation agent with TDD.
- Security/architecture escalation: route to oracle if product changes from public OBS URL to authorized endpoint/private bucket.
- Documentation/library lookup: librarian only if Excelize hyperlink API behavior is unclear for current version.
- No UI/designer/browser routing needed for backend export plan; browser evidence not relevant.

## Validation Commands

Run from `master/` unless stated otherwise:

```bash
rtk go test ./service -run 'TestGenerateSurveyReportExcel|TestSurveyReportAttachment|TestResolve'
rtk go test ./repository -run 'TestBuildSurveyReport|TestFindExportRows'
rtk go test ./...
```

Optional manual check when authorized credentials/data are available:

```bash
rtk docker compose -f ../docker-compose.yml ps
```

Then trigger `/v1/survey-report/export` through normal authenticated app flow with Jira filter and verify file from Download History.

## Evidence Requirements

- Local project discovery: required and completed. See `.opencode/evidence/20260503-1258-sx-1903-survey-attachment-links/discovery.md`.
- Official docs/context7: not required for plan; Excelize is already in project and only standard hyperlink/read APIs are needed. Implementation may consult docs if API signature uncertainty occurs.
- GitHub: not required; plan does not depend on upstream issue/source behavior.
- Brave/web search: not required; no current external facts needed.
- Browser/screenshot capture: not required; backend Excel export issue, not visual UI parity.
- Manual evidence after implementation should include:
  - sample output cell display and hyperlink target,
  - test command outputs,
  - note whether Jira evidence filter could be verified in environment.

## Done Criteria

- Tests listed above pass.
- Generated workbook attachment cells have hyperlink relationships with target public OBS URL.
- Empty attachment cells remain empty and have no hyperlink.
- Existing Download History download still returns valid `.xlsx`.
- No credentials or secrets added.
- Root cause and changed files reported in implementation summary.

## Final Planning Summary

- Artifacts created/consulted:
  - Created and kept `.opencode/evidence/20260503-1258-sx-1903-survey-attachment-links/discovery.md` because it contains operational discovery evidence for implementation handoff.
  - Created this primary plan `.opencode/plans/20260503-1258-sx-1903-survey-attachment-links.md` as the source of truth.
- Key decisions:
  - Use public OBS URL generated from existing OBS adapter configuration.
  - Write Excel hyperlink cells for attachments.
  - Validate already-full URLs against configured OBS host; do not emit arbitrary external URLs.
- Assumptions:
  - Existing public-read OBS object policy remains valid for this bugfix.
  - Maximum three attachment columns remains unchanged.
- Remaining open questions:
  - Confirm whether `mst.survey_answer_files` has `is_del`; if yes, consider adding deleted-file filtering.
- Readiness:
  - Ready for TDD implementation by build/fixer agent.
- Cleanup performed:
  - No draft artifacts were created.
  - Discovery evidence was intentionally kept for handoff; no stale draft/evidence cleanup needed.
