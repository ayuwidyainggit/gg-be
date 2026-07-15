# Plan — SX-XXXX Survey DistributorId Empty + TargetCustId (Maintenance Compatibility Patch)

Task ID: `20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id`
Source of truth: `.opencode/plans/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id.md`

## Goal

Memastikan `POST /master/v1/survey` (Store) dan `PUT /master/v1/survey/:survey_id` (Update) pada modul `master` menerima payload `distributor_id: []` ketika `target_cust_id` diisi, dengan memperlakukan scope sebagai principal (mirip jalur sentinel `[0]`) namun hanya ketika caller adalah principal tenant. Tujuan akhirnya: backend tetap backward-compatible untuk klien FE yang mengirim empty array, tanpa menambah kontrak FE baru dan tanpa membuka celah tenant scope untuk caller non-principal.

## Non-goals

- Tidak mengubah schema database, migrasi, atau repository query.
- Tidak mengubah kontrak API keluar (response body, status code) selain yang sudah dipakai.
- Tidak memperbaiki bug FE yang mengirim payload kosong jika memang tidak ada intent principal scope; FE tetap harus mengirim `distributor_id: [0]` atau list distributor riil bila tidak ada target_cust_id. Patch ini hanya membuat backend toleran terhadap `[]` + `target_cust_id`.
- Tidak menambah validasi, role, atau RBAC baru.
- Tidak menyentuh service lain di luar `master/service/survey_service.go` dan test-nya.

## Scope

- Module: `master`
- Endpoints: `POST /master/v1/survey`, `PUT /master/v1/survey/:survey_id`
- Files paling mungkin berubah (diff boundary):
  - `master/service/survey_service.go` (tweak pada `Store` dan `Update`, kemungkinan helper kecil)
  - `master/service/survey_service_test.go` (tambah regression test)
- Files yang hampir pasti TIDAK berubah (perlu dicek ulang saat eksekusi):
  - `master/entity/survey.go` (struktur sudah cukup: `DistributorId FlexibleIntArray`, `TargetCustId string`)
  - `master/controller/survey_controller.go` (controller cukup meneruskan; tidak ada logika tambahan yang perlu di sana)
  - `master/repository/*.go` (tidak ada perubahan query)

## Requirements

1. Saat `DistributorId` kosong (`[]`) dan `TargetCustId` non-empty, dan caller adalah principal (`request.CustId == request.ParentCustId` atau token principal), backend memperlakukan scope sebagai principal, sehingga `hasPrincipal = true` untuk logika downstream (`buildSurveyAreas`, `resolveTargetCustIdForAreas`, `resolveSalesmanCustIds`).
2. Perilaku existing yang harus tetap utuh:
   - `DistributorId: [0]` → principal scope (sentinel) tetap jalan seperti sekarang.
   - `DistributorId: [N, M, ...]` (id positif) → distributor scope tetap jalan, `0` (jika ada) tetap menjadi penanda principal dalam mixed list.
   - `DistributorId: []` + `TargetCustId` kosong → general survey path (tanpa principal scope, tanpa area rows) tetap jalan seperti `TestSurveyService_Store_ShouldStayCompatible_WhenEmpIdAndAreaAreEmpty`.
   - `AreaId` non-empty + `DistributorId: []` + `TargetCustId` kosong → tetap `ErrSurveyAreaDistributorRequired` (tidak boleh ter-bypass).
3. `TargetCustId` hanya dianggap sebagai bukti principal scope jika `request.CustId == request.ParentCustId` (i.e. caller adalah principal). Untuk caller distributor child, abaikan `[]` + `TargetCustId` sebagai indikasi principal; biarkan path existing berlaku (umumnya gagal dengan area/distributor rule).
4. `TargetCustId` yang dikirim FE, ketika principal scope aktif, harus tetap tunduk pada `resolveTargetCustIdForAreas` (override hanya pada row principal; row distributor tidak diubah).
5. Tidak ada perubahan transaksi, urutan insert, atau struktur response.
6. Patch minimum: tidak lebih dari satu helper baru dan dua call site di `Store`/`Update`. Tidak ada refactor besar.
7. Tenant scope safety: tidak boleh ada jalur baru yang menulis `distributor_id = 0` ke `m_survey_distributor` (tabel itu dikontrol oleh `buildSurveyDistributors(levelTarget, targetDistributorId)` dan tetap tidak terkait dengan `DistributorId` payload).

## Acceptance Criteria

1. `Store` dengan `DistributorId: []`, `TargetCustId: "C22001"`, `CustId: "C22001"`, `ParentCustId: "C22001"`, `AreaId: [82]` → success; menghasilkan survey area row dengan `DistributorId = 0` dan `TargetCustId = "C22001"`.
2. `Store` dengan `DistributorId: []`, `TargetCustId: ""`, `AreaId: [82]`, principal caller → `ErrSurveyAreaDistributorRequired` (tidak ter-bypass).
3. `Store` dengan `DistributorId: []`, `TargetCustId: ""`, tanpa `AreaId`, principal caller → success sebagai general survey (zero area rows, zero salesman rows bila `emp_id` kosong).
4. `Store` dengan `DistributorId: [0]`, `TargetCustId: "C22001"` → sama hasilnya dengan test existing `TestSurveyService_Store_ShouldCreatePrincipalOnlySurveyWithSelectedAreasAndSalesman` (tidak ada regresi).
5. `Store` dengan `DistributorId: [67]`, `TargetCustId: "C22001"`, distributor child caller (`ParentCustId != CustId`) → principal scope tidak aktif; tetap distributor-scope behavior (test akan meng-cover ini agar tidak ada false-positive dari patch).
6. `Update` dengan payload yang sama dengan Store (kriteria 1) → success; survey area row terganti (DeleteAreasBySurveyId dipanggil lalu StoreAreas menulis ulang).
7. `Update` tidak menambah `distributor_id = 0` ke `m_survey_distributor` (cek melalui `storeDistributorInput` stub tetap kosong untuk `LevelTarget != "Distributor"`).
8. Tests existing di `survey_service_test.go` tetap lulus tanpa modifikasi.

## Existing Patterns/Reuse

- `normalizeBusinessUnitSelection([]int)` di `master/service/survey_service.go` (sekitar baris 148) — dipakai di Store dan Update; letakkan tweak tepat setelahnya atau bungkus dalam helper.
- `isPrincipalScope(distributorIds []int)` — sudah ada untuk deteksi `0`; tidak relevan langsung untuk kasus ini (kasus kita `[]` tanpa `0`), tapi logikanya menunjukkan konvensi repo.
- `resolveTargetCustIdForAreas(...)` (baris 785) — tetap dipakai apa adanya; tweak di Store/Update hanya mengubah input `principalScope`/`hasPrincipal` yang diteruskan ke fungsi ini.
- `buildSurveyAreas(...)` (baris 742) — tidak diubah.
- `surveyRepositoryRedStub` dan `salesmanRepositoryStub` di `master/service/survey_service_test.go` — reuse penuh; cukup tambahkan test baru.
- Pattern mix `DistributorId: {0, 67, 68}` diuji oleh `TestSurveyService_Store_ShouldPersistPrincipalAndResolveAreas` — pertahankan sebagai regression.

## Source Anatomy

### Service layer
- `master/service/survey_service.go:481-598` — `Store()` membentuk `hasPrincipal`, `distributorIds`, `surveyAreas`, `salesmanCustIds`, lalu menjalankan write transaction. Ini titik patch utama. `confirmed_repo`.
- `master/service/survey_service.go:600-732` — `Update()` mengulang flow Store dengan delete-and-replace. Patch wajib mirror di sini. `confirmed_repo`.
- `master/service/survey_service.go:148-164` — `normalizeBusinessUnitSelection()` mendeteksi sentinel `0` dan memisahkan distributor positif. Jangan ubah. `confirmed_repo`.
- `master/service/survey_service.go:742-777` — `buildSurveyAreas()` menetapkan area rows dan error `ErrSurveyAreaDistributorRequired` bila area ada tapi distributor/principal scope tidak ada. `confirmed_repo`.
- `master/service/survey_service.go:785-801` — `resolveTargetCustIdForAreas()` hanya men-stamp `target_cust_id` ke row `DistributorId == 0`. Ini alasan patch cukup mengubah `hasPrincipal`, bukan fungsi ini. `confirmed_repo`.

### Entity / request contract
- `master/entity/survey.go:137-176` — `CreateSurveyBody` dan `UpdateSurveyBody` sudah punya `DistributorId FlexibleIntArray` dan `TargetCustId string`. Tidak perlu field baru. `confirmed_repo`.
- `master/entity/survey.go:11-27` — `FlexibleIntArray.UnmarshalJSON` sudah menerima array dan single int. Empty array tidak error. `confirmed_repo`.

### Controller layer
- `master/controller/survey_controller.go:125-179` — `Create()` hanya unmarshal, inject `cust_id/parent_cust_id/user_id`, validate, lalu panggil service. Tidak perlu patch. `confirmed_repo`.
- `master/controller/survey_controller.go:181-247` — `Update()` sama. Tidak perlu patch. `confirmed_repo`.

### Test harness
- `master/service/survey_service_test.go:16-209` — `surveyRepositoryRedStub`, `salesmanRepositoryStub`, `transactionManagerStub` cukup untuk semua regression baru. `confirmed_repo`.
- `master/service/survey_service_test.go:669-720` — principal-only sentinel `[0]` existing behavior sudah dijaga oleh test. `confirmed_repo`.
- `master/service/survey_service_test.go:942-986` — create/update parity already tested for principal+distributor mixed scope; new patch harus mirror pola ini. `confirmed_repo`.

## Reference Map

- **Feature: infer principal scope from empty `distributor_id` + `target_cust_id`**
  - Basis: `repo-backed`
  - Sumber: `master/service/survey_service.go` helper flow + comments `resolveTargetCustIdForAreas`.
  - Alasan cukup: bug dan fix seluruhnya berada di service request normalization; tidak perlu docs eksternal.
- **Feature: preserve sentinel `[0]` principal behavior**
  - Basis: `repo-backed`
  - Sumber: `normalizeBusinessUnitSelection()`, principal tests existing di `survey_service_test.go`.
  - Alasan cukup: sentinel adalah kontrak internal repo, bukan library behavior.
- **Feature: preserve distributor behavior for positive ids**
  - Basis: `repo-backed`
  - Sumber: `buildSurveyAreas()`, existing tests distributor/mixed scope.
  - Alasan cukup: semua semantics sudah tertutup oleh code + tests lokal.
- **Feature: no controller/entity change**
  - Basis: `repo-backed`
  - Sumber: `master/controller/survey_controller.go`, `master/entity/survey.go`.
  - Alasan cukup: field yang dibutuhkan sudah tersedia; controller tidak memuat branching business logic.

## Confirmed vs Assumed Audit

| Claim | Status | Source |
| --- | --- | --- |
| Store memakai `normalizeBusinessUnitSelection` sebelum area/salesman resolution | confirmed_repo | `master/service/survey_service.go:499-518` |
| Update memakai flow yang sama | confirmed_repo | `master/service/survey_service.go:623-642` |
| Empty `distributorIds` di `resolveSurveyCustIds` menghasilkan `[]string{custId}` | confirmed_repo | `master/service/survey_service.go:79-99` |
| `target_cust_id` hanya dipakai pada row principal (`DistributorId == 0`) | confirmed_repo | `master/service/survey_service.go:785-801` |
| Controller tidak perlu diubah untuk patch ini | confirmed_repo | `master/controller/survey_controller.go:125-247` |
| `FlexibleIntArray` sudah menerima empty array | confirmed_repo | `master/entity/survey.go:11-27` |
| Principal caller dikenali dari `custId == parentCustId` atau parent kosong | assumption | disimpulkan dari `resolveSurveyCustIds` fallback `parentCustId == "" -> custId`; belum diverifikasi runtime |
| FE kadang mengirim `target_cust_id == custId` untuk principal flow | assumption | berasal dari task statement, belum diverifikasi runtime |
| FE bug masih mungkin ada sesudah patch backend | user_confirmed | user request eksplisit: backend backward-compatible walau frontend bug mungkin masih ada |

## Progress Tracking

- `tracker_path`: `.opencode/state/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id/progress.json`
- `init_command`: `python3 ~/.config/opencode/scripts/task-progress.py --project-root . --task 20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id --init --plan .opencode/plans/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id.md`
- `summary_command`: `python3 ~/.config/opencode/scripts/task-progress.py 20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id --summary`
- `checklist_command`: `python3 ~/.config/opencode/scripts/task-progress.py 20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id --checklist`
- `update_rules`:
  - update wajib saat task pindah `pending -> in_progress`
  - update wajib setelah `completed`, `blocked`, atau `cancelled`
  - update wajib setiap evidence file ditulis
  - update wajib di setiap handoff lintas lane
- `task_map`:

| Task ID | Owner | Evidence Path | Update Command |
| --- | --- | --- | --- |
| A1 | `@fixer` | `.opencode/evidence/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id/test.log` | `python3 ~/.config/opencode/scripts/task-progress.py 20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id --update A1 --status completed --owner @fixer --evidence .opencode/evidence/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id/test.log` |
| A2 | `@quality-gate` | `.opencode/evidence/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id/test.log` | `python3 ~/.config/opencode/scripts/task-progress.py 20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id --update A2 --status completed --owner @quality-gate --evidence .opencode/evidence/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id/test.log` |

tracker updates at every status transition are mandatory, not optional bookkeeping.

## Execution ownership table

| Area | Implementation owner | Review gate owner |
| --- | --- | --- |
| Service normalization + principal inference | `@fixer` | `@quality-gate` |
| Regression tests | `@fixer` | `@quality-gate` |
| Final conformance/risk review | `@quality-gate` | `@quality-gate` |

## Existing Patterns/Reuse

- `normalizeBusinessUnitSelection([]int)` di `master/service/survey_service.go` (sekitar baris 148) — dipakai di Store dan Update; letakkan tweak tepat setelahnya atau bungkus dalam helper.
- `isPrincipalScope(distributorIds []int)` — sudah ada untuk deteksi `0`; tidak relevan langsung untuk kasus ini (kasus kita `[]` tanpa `0`), tapi logikanya menunjukkan konvensi repo.
- `resolveTargetCustIdForAreas(...)` (baris 785) — tetap dipakai apa adanya; tweak di Store/Update hanya mengubah input `principalScope`/`hasPrincipal` yang diteruskan ke fungsi ini.
- `buildSurveyAreas(...)` (baris 742) — tidak diubah.
- `surveyRepositoryRedStub` dan `salesmanRepositoryStub` di `master/service/survey_service_test.go` — reuse penuh; cukup tambahkan test baru.
- Pattern mix `DistributorId: {0, 67, 68}` diuji oleh `TestSurveyService_Store_ShouldPersistPrincipalAndResolveAreas` — pertahankan sebagai regression.

## Constraints

- Arsitektur `Controller → Service → Repository → DB` tidak boleh dilanggar; tweak hanya di service layer.
- Transaksi di service (`txManager.WithinTransaction`) tidak boleh berubah; urutan `Store` → `StoreAreas` → `StoreSurveyDistributors` → `StoreSalesmen` → `StoreOutlets` → `StoreDetails` (Store) dan versi Update masing-masing tetap sama.
- `distributor_id = 0` hanya ditulis ke `m_survey_area` untuk principal row (per komentar `resolveTargetCustIdForAreas`); tidak boleh ditulis ke `m_survey_distributor`.
- Validation tenant: principal hanya boleh dianggap principal jika `request.CustId == request.ParentCustId`. Untuk caller distributor child, principal-scope shortcut tidak boleh aktif meskipun `target_cust_id` dikirim.
- Tidak boleh menambah dependency baru.
- Validasi via `rtk go test ./service -run 'TestSurveyService_Store|TestSurveyService_Update'` dari `master/`.

## Risks

1. **False-positive principal scope**: bila `TargetCustId` non-empty diterjemahkan sebagai principal scope untuk caller distributor child, akan ada overwrite row area principal. Mitigasi: guard `request.CustId == request.ParentCustId` (atau pola token principal yang setara) sebelum mengaktifkan `hasPrincipal`.
2. **FE confusion**: FE yang mengirim `distributor_id: []` tanpa `target_cust_id` akan tetap gagal dengan `ErrSurveyAreaDistributorRequired` bila ada `area_id`. Ini bisa dianggap regresi FE, padahal backend hanya menerima `[]` + `target_cust_id` sesuai requirement. Catat di ringkasan bahwa FE harus konsisten: pakai `[0]` atau isi `target_cust_id`.
3. **Konsistensi create vs update**: bila hanya `Store` yang di-tweak, `Update` akan tetap gagal pada payload yang sama. Mitigasi: tweak di kedua fungsi dengan helper yang sama.
4. **Test flake**: stub `surveyRepositoryRedStub` menyimpan `findAreasInput`/`findCustIdsInput`; pastikan test baru tidak bergantung pada urutan call.

## Decisions/Assumptions

- **Decision**: aktivasi `hasPrincipal` ketika `DistributorId` kosong, `TargetCustId` non-empty, dan caller adalah principal. Helper kecil (mis. `func inferPrincipalScope(distributorIds []int, targetCustId, custId, parentCustId string) bool`) ditambahkan di `survey_service.go` agar Store dan Update pakai sumber kebenaran yang sama.
- **Decision**: helper hanya membungkus keputusan `hasPrincipal`; tidak menyentuh `distributorIds` (tetap `[]` untuk lookup, sesuai `normalizeBusinessUnitSelection`).
- **Assumption**: principal caller dikenali dari `custId == parentCustId` (mengikuti konvensi di `resolveSurveyCustIds` baris 84-87). Bila middleware menyuntikkan parent_cust_id berbeda untuk principal, perlu konfirmasi user. Catat di open question.
- **Assumption**: `target_cust_id` di FE, ketika principal, berisi cust_id principal sendiri (lihat test existing yang memakai `CustId: "C22001"` sebagai principal dan `TargetCustId: "C22001"`). Tidak perlu validasi tambahan bahwa `target_cust_id == custId`; biarkan `resolveTargetCustIdForAreas` yang men-stamp `target_cust_id` ke row principal.
- **Open question**: apakah FE kadang mengirim `parent_cust_id` kosong untuk principal (kasus langka). Jika iya, helper perlu fallback `if parentCustId == "" { parentCustId = custId }` (sudah jadi pola di `resolveSurveyCustIds`).

## Execution Source of Truth

Urutan precedence untuk eksekusi:

1. Pertanyaan user terbaru (tidak ada untuk patch ini).
2. Aturan keamanan/tenant repo (`.opencode/docs/ARCHITECTURE.md`).
3. Invariant di section "Non-negotiable Implementation Invariants" plan ini.
4. Acceptance Criteria plan ini.
5. Implementation Steps plan ini.
6. Rekomendasi/follow-up.

## Non-negotiable Implementation Invariants

1. `hasPrincipal` hanya true jika salah satu kondisi terpenuhi:
   a. `normalizeBusinessUnitSelection(...)` mengembalikan `true` (ada `0` di payload), atau
   b. `DistributorId` kosong, `TargetCustId` non-empty, dan `custId == parentCustId` (atau `parentCustId == ""`).
2. `DistributorId` kosong + principal inferred + `AreaId` non-empty + `TargetCustId` kosong → tetap `ErrSurveyAreaDistributorRequired` (guard di `buildSurveyAreas` masih berlaku karena `hasPrincipal = true` akan membuat `len(distributorIds) == 0` tetap gagal jika `!hasPrincipal`; di sini `hasPrincipal` true, tapi area_id harus non-empty dan area lookup tetap kosong → `buildSurveyAreas` akan membuat principal row, bukan error. Test #2 di atas mengunci perilaku ini).
3. `m_survey_distributor` tidak pernah diisi dengan `0`; hanya dikontrol oleh `buildSurveyDistributors(levelTarget, targetDistributorId)`.
4. `resolveSalesmanCustIds` tetap dipanggil dengan `distributorIds = []` ketika `DistributorId = []`; helper internal `resolveSurveyCustIds` sudah benar mengembalikan `[]string{custId}` (yang akan di-merge dengan principal custId oleh `mergeUniqueCustIDs` jika `hasPrincipal`).
5. Behavior `TestSurveyService_Store_ShouldCreatePrincipalOnlySurveyWithSelectedAreasAndSalesman` (DistributorId: `[0]`) tidak boleh berubah.

## Do Not / Reject If

- Jangan ubah `normalizeBusinessUnitSelection` atau `buildSurveyAreas`/`buildSurveyDistributors`/`resolveSalesmanCustIds`/`resolveTargetCustIdForAreas`.
- Jangan tambah field ke `entity.CreateSurveyBody` / `entity.UpdateSurveyBody`.
- Jangan ubah transaksi atau urutan insert.
- Jangan tulis `distributor_id = 0` ke tabel `m_survey_distributor`.
- Jangan aktifkan principal scope untuk caller distributor child hanya karena `target_cust_id` dikirim.
- Jangan refactor besar; patch harus <= 30 baris pada file service (tidak termasuk test).
- Jangan implementasikan perubahan FE; cukup catat sebagai catatan di Final Planning Summary.

## Diff Boundary

- **Allowed files**:
  - `master/service/survey_service.go` (tambah helper + tweak Store/Update)
  - `master/service/survey_service_test.go` (tambah test)
- **Out of boundary** (jangan disentuh, revert jika berubah):
  - `master/controller/survey_controller.go`
  - `master/entity/survey.go`
  - `master/repository/*.go`
  - file lain di repo
- **Evidence path**: `.opencode/evidence/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id/test.log` (output `go test`), opsional `.opencode/evidence/.../manual-trace.md` jika ada retest manual.

## TDD/Test Plan

TDD required (production logic, service-layer).

### Red step — tambahkan test berikut di `master/service/survey_service_test.go`

1. `TestSurveyService_Store_ShouldAcceptEmptyDistributorIdWithTargetCustId_AsPrincipal`
   - `DistributorId: FlexibleIntArray{}` (atau setara `entity.FlexibleIntArray(nil)`)
   - `TargetCustId: "C22001"`, `CustId: "C22001"`, `ParentCustId: "C22001"`
   - `AreaId: []int{82}`, `EmpId: entity.FlexibleIntArray{369}` (salesman valid)
   - Stub `findCustIdsResult: []string{"C22001"}` dan `findAreasResult: nil` (tidak ada lookup karena `distributorIds=[]`)
   - Expect: `err == nil`, `len(repo.storeAreasInput) == 1` dengan `DistributorId = 0` dan `AreaId = 82`; `findAreasInput` kosong.

2. `TestSurveyService_Update_ShouldAcceptEmptyDistributorIdWithTargetCustId_AsPrincipal`
   - Setup `findOneSurvey` existing seperti test Update existing.
   - Payload sama dengan #1, `TargetCustId: "C22001"`.
   - Expect: `err == nil`, `repo.deleteAreasCalled == true`, `len(repo.storeAreasInput) == 1` dengan `DistributorId = 0` dan `AreaId = 82`.

3. `TestSurveyService_Store_ShouldIgnoreEmptyDistributorIdWithTargetCustId_WhenCallerIsDistributorChild`
   - `CustId: "C220010001"`, `ParentCustId: "C22001"`, `TargetCustId: "C22001"`, `DistributorId: FlexibleIntArray{}`, `AreaId: []int{82}`.
   - Expect: error (bisa `ErrSurveyAreaDistributorRequired` atau setara) — patch tidak boleh aktif untuk caller non-principal.

4. `TestSurveyService_Store_ShouldStillRequireDistributorId_WhenAreaIdProvidedWithoutTargetCustId`
   - `DistributorId: FlexibleIntArray{}`, `TargetCustId: ""`, `AreaId: []int{82}`, principal caller.
   - Expect: `ErrSurveyAreaDistributorRequired` (regression guard).

5. (Opsional) `TestSurveyService_Store_ShouldPreservePrincipalSentinelBehavior_WhenDistributorIdContainsZero` — cukup kalau test existing `TestSurveyService_Store_ShouldCreatePrincipalOnlySurveyWithSelectedAreasAndSalesman` masih lulus.

### Green step

Implementasi minimal di `master/service/survey_service.go`:

1. Tambah helper:
   ```go
   func inferPrincipalScopeFromTargetCustId(distributorIds []int, targetCustId, custId, parentCustId string) bool {
       if len(distributorIds) > 0 {
           return false // distributor ids ada; biarkan normalizeBusinessUnitSelection yang menentukan
       }
       if targetCustId == "" {
           return false
       }
       effectiveParent := parentCustId
       if effectiveParent == "" {
           effectiveParent = custId
       }
       return custId != "" && custId == effectiveParent
   }
   ```
2. Di `Store` (baris ~499), setelah `hasPrincipal, distributorIds := normalizeBusinessUnitSelection(...)`, tambahkan:
   ```go
   if !hasPrincipal {
       if inferPrincipalScopeFromTargetCustId(distributorIds, request.TargetCustId, request.CustId, request.ParentCustId) {
           hasPrincipal = true
       }
   }
   ```
3. Tweak yang sama persis di `Update` (baris ~623).
4. Tidak ada perubahan lain.

### Refactor step

- Tidak wajib; helper sudah kecil.
- Jika muncul duplikasi lain, ekstrak helper di langkah Refactor.

### Edge cases (covered by tests + cases existing)

- `DistributorId: []` + `TargetCustId: ""` + area → error.
- `DistributorId: []` + `TargetCustId: ""` + no area + no emp → general survey.
- `DistributorId: []` + `TargetCustId: non-empty` + principal caller → principal scope.
- `DistributorId: []` + `TargetCustId: non-empty` + distributor child caller → bukan principal.
- `DistributorId: [0]` (sentinel existing) → tetap principal scope (test existing lulus).
- `DistributorId: [N, ...]` (positif) → distributor scope, `0` (jika ada) mixed (test existing lulus).
- Update dengan payload di atas → replace area rows.
- `m_survey_distributor` tetap tidak terisi `0` (cek via `storeDistributorInput` stub tetap `nil`/`len==0` saat `level_target != "Distributor"`).

### Commands

Jalankan dari `master/`:

```bash
rtk go mod download && rtk go mod tidy
rtk go test ./service -run 'TestSurveyService_Store|TestSurveyService_Update' -v
rtk go test ./...
```

## Implementation Steps

1. **A1** | `@fixer` | Baca ulang `master/service/survey_service.go` Store dan Update untuk konfirmasi titik injeksi helper.
2. **A1** | `@fixer` | Tambahkan test Red (5 test sesuai TDD/Test Plan) di `master/service/survey_service_test.go`.
3. **A1** | `@fixer` | Jalankan `rtk go test ./service -run 'TestSurveyService_Store|TestSurveyService_Update'` dari `master/`; pastikan gagal (Red).
4. **A1** | `@fixer` | Implementasikan helper `inferPrincipalScopeFromTargetCustId` di `master/service/survey_service.go` (top-level, dekat helper `normalizeBusinessUnitSelection`).
5. **A1** | `@fixer` | Tambah 2 baris tweak di `Store` dan 2 baris tweak di `Update` setelah `normalizeBusinessUnitSelection(...)`.
6. **A1** | `@fixer` | Jalankan ulang test; pastikan Red → Green.
7. **A1** | `@fixer` | Jalankan `rtk go test ./...` dari `master/`; pastikan tidak ada regresi.
8. **A1** | `@fixer` | Simpan jejak eksekusi di `.opencode/evidence/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id/test.log`.
9. **A2** | `@quality-gate` | Review diff boundary, tenant invariants, dan output test; pastikan claim hanya backend backward-compatible.
10. **A2** | `@quality-gate` | Handoff final ke `@orchestrator`; planner tidak implementasi lebih lanjut.

## Worklist

1. **A1** | `@fixer` | Implementasi patch service + tambah regression tests
2. **A2** | `@quality-gate` | Final conformance review untuk tenant rule, diff boundary, dan claim scope

## Expected Files to Change

- `master/service/survey_service.go`
- `master/service/survey_service_test.go`

## Agent/Tool Routing

- **Owner**: `@fixer` (bounded implementation di service + test, multi-file tapi atomik).
- **Review**: `@quality-gate` (material change, menyentuh production logic + tenant rules).
- **Skip**: `@designer`, `@architect`, `@oracle` — tidak perlu untuk maintenance slice sekecil ini.
- **Lane note**: planner tidak implementasi. Setelah plan ini di-PASS, lane berpindah ke `@orchestrator`/`@fixer` dengan permission implementation aktif.

## Executor Handoff Prompt

```text
Task ID: 20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id
Source of truth: .opencode/plans/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id.md
Scope: master service Store/Update — terima `distributor_id: []` saat `target_cust_id` diisi (caller principal).

must_preserve:
  - Normalisasi distributor_id existing (`normalizeBusinessUnitSelection`) untuk path `[0]` dan `[N,...]`.
  - `resolveTargetCustIdForAreas` semantics; row principal di-stamp `target_cust_id`, row distributor tidak.
  - Transaksi dan urutan insert di Store/Update.
  - Tenant rule: principal hanya jika `custId == parentCustId` (atau `parentCustId == ""`).
  - Tests existing di survey_service_test.go (jangan modifikasi kecuali menambah).
  - `m_survey_distributor` tidak pernah diisi `0`; hanya `buildSurveyDistributors(levelTarget, targetDistributorId)` yang mengontrol.

do_not_touch:
  - master/controller/survey_controller.go
  - master/entity/survey.go
  - master/repository/*
  - file lain di luar diff boundary.

validation:
  - rtk go mod download && rtk go mod tidy
  - rtk go test ./service -run 'TestSurveyService_Store|TestSurveyService_Update' -v
  - rtk go test ./...

claim_scope:
  - Boleh klaim "implemented" hanya jika semua 5 test baru di TDD/Test Plan hijau dan tests existing lulus.
  - Jangan klaim "FE fixed". Patch ini hanya membuat backend backward-compatible. Bug FE (mengirim `distributor_id: []` tanpa `target_cust_id` saat area dipilih) tetap harus diperbaiki FE; catat di evidence.

evidence_required:
  - Output `go test` di .opencode/evidence/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id/test.log
  - Diff ringkas (git diff) di evidence yang sama atau inline di handoff.
```

## Execution-ready Worklist / Handoff Contract

```yaml
handoff: compatibility_patch
plan_id: 20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id
task_id: 20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id
caller: orchestrator
callee: fixer
scope: Tambah helper + tweak Store/Update di master/service/survey_service.go; tambah 5 regression test di master/service/survey_service_test.go.
claim_level: scoped
claim_scope: Implementasi + tests hijau; tidak mencakup perbaikan FE.
source_basis: ["master/service/survey_service.go (Store ~481-598, Update ~600-732, helpers 148-218, 742-801)", "master/service/survey_service_test.go (stub patterns + existing principal/0/distributor tests)", "master/entity/survey.go (struct: DistributorId FlexibleIntArray, TargetCustId string, no field changes)", "master/controller/survey_controller.go (no changes; controller forwarding only)", ".opencode/docs/ARCHITECTURE.md (Controller → Service → Repository → DB)"]
must_preserve: ["hasPrincipal dari normalizeBusinessUnitSelection saat ada 0", "behavior [N,...] distributor scope", "behavior DistributorId:[0] principal sentinel", "ErrSurveyAreaDistributorRequired saat area tanpa distributor", "transaksi dan urutan insert", "m_survey_distributor tidak terisi 0"]
do_not_touch: ["master/controller/survey_controller.go", "master/entity/survey.go", "master/repository/*"]
validation: ["rtk go mod download && rtk go mod tidy", "rtk go test ./service -run 'TestSurveyService_Store|TestSurveyService_Update' -v", "rtk go test ./..."]
exit_criteria: ["5 test baru hijau (Store principal empty+target, Update principal empty+target, Store distributor child empty+target tetap non-principal, Store empty+area no target tetap error, existing test [0] masih lulus)", "Tests existing di survey_service_test.go tidak dimodifikasi", "rtk go test ./... hijau"]
evidence_required: [".opencode/evidence/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id/test.log (test output)"]
depends_on: ["none"]
context_bundle: ["master/service/survey_service.go", "master/service/survey_service_test.go", "master/entity/survey.go", "master/controller/survey_controller.go", ".opencode/docs/ARCHITECTURE.md", ".opencode/docs/QUALITY.md"]
```

## Validation Commands

Jalankan dari direktori `master/`:

```bash
rtk go mod download
rtk go mod tidy
rtk go test ./service -run 'TestSurveyService_Store|TestSurveyService_Update' -v
rtk go test ./service
rtk go test ./...
```

Expected:
- 5 test baru hijau.
- Semua test existing tetap hijau.
- Tidak ada error build / vet.

## Evidence Requirements

- `.opencode/evidence/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id/test.log` — output `go test -v`.
- Opsional `.opencode/evidence/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id/manual-trace.md` untuk retest manual dengan curl/Postman.
- Catatan eksplisit di evidence: FE bug (mengirim `distributor_id: []` tanpa `target_cust_id`) TIDAK ditutup oleh patch ini.

## Done Criteria

- 5 test baru di TDD/Test Plan hijau.
- Tests existing di `survey_service_test.go` tidak dimodifikasi dan tetap hijau.
- Diff hanya pada `master/service/survey_service.go` dan `master/service/survey_service_test.go`.
- Tidak ada perubahan pada controller, entity, repository, atau file lain.
- `rtk go test ./...` di `master/` hijau.
- Evidence test tersimpan.
- `@quality-gate` signoff (material change + tenant rule).

## Final Planning Summary

- **Artifacts consulted**:
  - `master/service/survey_service.go` (Store, Update, helper area, principal scope)
  - `master/service/survey_service_test.go` (stub pattern, existing principal/0/distributor tests)
  - `master/entity/survey.go` (struct fields, no changes needed)
  - `master/controller/survey_controller.go` (no changes needed)
  - `.opencode/docs/ARCHITECTURE.md` (layering + tenant rules)
  - `.opencode/docs/QUALITY.md` (validation commands)
  - Prior plan `.opencode/plans/20260513-1611-sx-1965-survey-select-all-save.md` (survey plan shape reference)
- **Artifacts created**:
  - `.opencode/plans/20260710-1400-sx-survey-distributor-id-empty-with-target-cust-id.md` (this plan)
- **Key decisions**:
  - Patch minimum: satu helper + 2 baris tweak di Store dan Update.
  - Principal-scope inference hanya jika `custId == parentCustId` (atau parent kosong) untuk keamanan tenant.
  - Tidak mengubah kontrak API atau FE.
- **Assumptions**:
  - Principal caller dikenali dari `custId == parentCustId` sesuai konvensi `resolveSurveyCustIds`.
  - FE biasanya mengirim `target_cust_id == custId` untuk principal.
- **Open questions**:
  - Apakah middleware FE pernah menyuntikkan `parent_cust_id` non-canonical untuk principal? Belum ada evidence runtime; helper mencakup fallback `parent_cust_id == ""`.
  - FE bug (empty distributor_id tanpa target_cust_id) perlu perbaikan FE terpisah; bukan bagian patch ini.
- **Readiness**: `PASS` (maintenance stability, diff boundary minimal, TDD terdefinisi, command valid tersedia, evidence path siap).
- **Cleanup performed**:
  - Tidak ada draft/evidence yang perlu di-cleanup; plan ditulis langsung sebagai primary.
- **Active-lane reset note**:
  - Eksekusi harus dilakukan di lane `@orchestrator`/`@fixer` berikutnya. Permission implementation aktif di lane itu, bukan di lane planner. Planner hanya menulis plan dan berhenti di sini.
