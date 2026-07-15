# Plan — SX-1965 Save Survey Select All

Task ID: `20260513-1611-sx-1965-survey-select-all-save`

Primary source of truth: `.opencode/plans/20260513-1611-sx-1965-survey-select-all-save.md`

## Goal

Memperbaiki backend `master` pada endpoint `POST /master/v1/survey` agar flow create survey dengan selection all tetap aman dan sukses untuk salesman yang valid dalam scope principal/distributor, sambil mengembalikan error yang jelas bila ada salesman invalid atau inactive.

## Non-goals

- Tidak mengubah kontrak FE selain memperkaya response error invalid salesman bila diperlukan.
- Tidak mengubah schema database atau migrasi.
- Tidak mengubah rule akses di endpoint list salesman/outlet di luar kebutuhan langsung create survey.
- Tidak menaruh token, JWT, atau kredensial sensitif ke source code, fixture, atau log permanen.

## Scope

- Module: `master`
- Endpoint utama: `POST /master/v1/survey`
- Endpoint pendamping untuk konsistensi: `PUT /master/v1/survey/:survey_id`
- Area kode utama:
  - `master/controller/survey_controller.go`
  - `master/service/survey_service.go`
  - `master/repository/survey_repository.go`
  - `master/repository/salesman_repository.go`
  - `master/entity/survey.go`
  - `master/service/survey_service_test.go`
  - opsional `master/controller/survey_controller_test.go`

## Requirements

1. Principal dapat membuat survey untuk:
   - salesman milik principal itu sendiri; dan
   - salesman milik distributor child di bawah principal tersebut.
2. Distributor hanya dapat membuat survey untuk salesman milik distributor login.
3. `distributor_id` yang mengandung `0` harus dinormalisasi aman sebagai sentinel select-all/principal-scope, bukan distributor riil.
4. `0` tidak boleh disimpan sebagai distributor aktual yang di-resolve dari FE, kecuali existing sentinel internal `m_survey_area.distributor_id = 0` yang memang sudah dipakai modul survey untuk principal scope area.
5. Validasi salesman harus dapat mengidentifikasi salesman yang invalid atau inactive secara detail.
6. Create/update survey tetap atomic dalam transaction.
7. Error DB/ORM perlu dibungkus lebih baik dan sementara dapat ditrace via logging terstruktur minimal.
8. Himpunan salesman yang diterima pada create survey harus selaras dengan himpunan salesman yang dikembalikan endpoint `GET /master/v1/salesman` untuk filter FE yang sama.

## Acceptance Criteria

- Principal dapat save survey saat memilih Area All, Business Unit All, Sales Team All, dan Salesman All, selama semua salesman valid dalam scope principal.
- Payload dengan `emp_id: [450,435,415,421,458,459,466]` sukses jika seluruh salesman valid menurut rule bisnis final; bila tidak, response wajib menjelaskan salesman invalid secara spesifik.
- Salesman distributor child principal diterima sebagai target valid.
- Distributor user tetap ditolak bila mencoba memilih salesman distributor lain.
- `distributor_id: [0, ...]` tidak menyebabkan lookup/insert distributor aktual `0` dan tidak mematahkan save.
- Untuk scope FE seperti `sales_team_id=82,81,80,78,77,66,65` dan `distributor_id=0,102,103,119`, salesman yang muncul dari endpoint list salesman menjadi acuan kandidat valid pada create survey.
- Jika ada salesman invalid/inactive/out-of-scope, response memuat `invalid_emp_id` dan `invalid_salesman`.
- Tidak ada partial insert ketika salah satu insert target gagal.
- Regression test mencakup principal scope, distributor scope, sentinel `0`, invalid salesman, dan rollback.

## Existing Patterns/Reuse

- `normalizeBusinessUnitSelection()` sudah melakukan dedupe distributor positif dan menandai keberadaan sentinel `0`.
- `resolveSurveyCustIds()` sudah merepresentasikan rule child distributor melalui `FindCustIdsByDistributorIds(parentCustId, distributorIds)`.
- `surveyServiceImpl.Store()` dan `Update()` sudah memakai `txManager.WithinTransaction(...)`.
- `survey_service_test.go` sudah punya banyak regression test untuk principal-only, mixed `0 + distributor`, child distributor, dan rollback.
- Existing survey area behavior sudah menggunakan sentinel internal `distributor_id = 0` untuk principal area mapping; ini perlu dipertahankan bila masih relevan dengan detail response.

## Constraints

- Wajib mengikuti arsitektur `Controller → Service → Repository → DB`.
- Write tetap berada di service-layer transaction.
- Validation logic jangan bergantung pada FE mengirim data bersih; backend harus tetap aman.
- Validasi tenant harus menghormati `cust_id` dan `parent_cust_id`.
- Validasi dan evidence runtime utama difokuskan ke local `master` service, sesuai keputusan user.
- Perintah validasi mengikuti repo guidance untuk module `master`.

## Risks

1. Data DB aktual menunjukkan `emp_id 435 / Erling Braut Caraka` inactive, sehingga payload QA mungkin memang berisi data invalid.
2. Klarifikasi terbaru menunjukkan create survey seharusnya menerima data yang lolos dari endpoint list salesman. Jika list endpoint masih mengembalikan inactive salesman, ada mismatch kontrak antar endpoint yang harus diputuskan dengan evidence.
3. Bila response body diubah dari message-only menjadi message + detail invalid, ada risiko compatibility untuk consumer lain.
4. Logging debug yang terlalu verbose bisa bocorkan context sensitif bila tidak dibatasi.
5. Jika service saat ini tidak mengecek `is_active` secara eksplisit di validasi salesman, implementer perlu memastikan apakah query/validator existing sudah atau belum menolak inactive secara benar.

## Decisions/Assumptions

- **Keputusan:** backend harus mengembalikan error actionable saat salesman invalid, bukan `ErrSurveySalesmanNotFound` generik saja.
- **Keputusan:** sentinel `0` diperlakukan sebagai bagian dari scope principal/select-all, tetapi tidak boleh diteruskan ke distributor lookup positif.
- **Keputusan:** implementasi harus menjaga behavior mixed principal + child distributor seperti test existing `DistributorId: {67,0,68}`.
- **Keputusan baru:** source-of-truth praktis untuk kandidat salesman valid adalah hasil endpoint `GET /master/v1/salesman` dengan filter FE yang sama. Validator create tidak boleh lebih ketat tanpa alasan bisnis yang terbukti.
- **Asumsi revisi:** bila endpoint list salesman mengembalikan salesman inactive untuk scope tersebut, implementer harus membuktikan apakah create perlu mengikuti list atau justru list yang perlu dibenahi; jangan diam-diam mempertahankan mismatch.
- **Asumsi:** local DB dev `scylla_citus_dev` cukup representatif untuk root-cause verification awal.
- **Open question tersisa:** apakah product menghendaki create mengikuti penuh hasil endpoint list salesman, termasuk bila list mengandung inactive, atau create tetap boleh menolak inactive dengan error detail. Ini perlu dipastikan lewat evidence runtime sebelum final signoff.

## TDD/Test Plan

- **TDD required:** Ya.
- **Alasan:** ini bug fix production logic, menyentuh access scope, transaction safety, dan error contract.
- **Existing test patterns:** gunakan `surveyRepositoryRedStub`, `salesmanRepositoryStub`, dan `transactionManagerStub` di `master/service/survey_service_test.go`.

### Red step

Tambahkan test baru terlebih dahulu:

1. `TestSurveyService_Store_ShouldAllowPrincipalOwnedAndChildDistributorSalesmen_WhenDistributorSelectionContainsZero`
   - request principal dengan `DistributorId: {0,102,103,119}` atau versi ringkas `{0,102}`
   - `findCustIdsResult` mengembalikan child cust IDs
   - stub salesman valid per cust:
     - principal-owned salesman valid pada `custId` principal
     - child distributor salesman valid pada cust child
   - expect success dan `storeSalesmenInput` menyimpan cust_id per-employee yang benar.

2. `TestSurveyService_Store_ShouldReturnDetailedInvalidSalesmen_WhenAnyEmpIdIsOutOfScopeOrInactive`
   - request mixed emp_id berisi satu salesman invalid/inactive
   - expect typed error baru yang memuat daftar `invalid_emp_id` dan `invalid_salesman`
   - expect tidak ada insert salesman dan transaction gagal.

3. `TestSurveyService_Store_ShouldIgnoreZeroForDistributorLookup_WhenSelectAllPayloadContainsZero`
   - memastikan `findAreasInput` dan `findCustIdsInput` hanya memuat distributor positif.

4. `TestSurveyService_Store_ShouldRollback_WhenStoreSalesmenFails_AfterDetailedValidationPasses`
   - pertahankan regression rollback existing.

5. `TestSurveyService_Update_ShouldMirrorCreateValidationRules_ForPrincipalAndDistributorScope`
   - optional tetapi direkomendasikan agar create/update konsisten.

6. Tambahkan atau siapkan evidence test/integration yang membandingkan source scope:
   - hasil resolver create survey vs hasil endpoint/filter salesman untuk scope yang sama.
   - bila tidak feasible sebagai automated test penuh, minimal sebagai manual evidence/query comparison.

Jika controller test mudah ditambah, tambahkan satu test untuk memastikan typed validation error diterjemahkan menjadi HTTP `400` dengan payload detail.

### Green step

Implementasi minimal yang direncanakan:

1. Tambah typed error/domain error baru, misalnya `SurveyInvalidSalesmenError`:
   - `Message`
   - `InvalidEmpID []int`
   - `InvalidSalesman []string`
2. Refactor validasi salesman dari `resolveSalesmanCustIds(...)` agar tidak berhenti di error generik pertama.
3. Validation flow baru perlu:
   - normalisasi unique emp ids;
   - resolve semua cust candidates dari principal/distributor scope;
   - cek setiap salesman terhadap scope customer yang valid;
   - ambil detail nama salesman dan status aktif bila perlu melalui repository/helper baru atau query existing yang diperluas;
   - kumpulkan invalid list, lalu return typed error bila ada.
4. Pastikan principal request dapat memvalidasi dua jenis salesman sekaligus:
   - `cust_id = principal`
   - `cust_id IN child distributor customers`
5. Pertahankan logic `normalizeBusinessUnitSelection()` sehingga `0` tidak ikut dalam lookup distributor positif.
6. Bungkus error DB insert dengan context yang lebih jelas dan tambahkan logging sementara yang tidak mencatat token.

### Refactor step

- Ekstrak helper khusus validasi target salesman agar `Store()` dan `Update()` tidak menduplikasi logic.
- Bila perlu, buat repository helper read-only ringan untuk mengambil identitas invalid salesman tanpa mencampur business logic ke repository.
- Rapikan typed error handling di controller agar response contract tetap konsisten.

### Edge cases

- `distributor_id: [0]` principal-only.
- `distributor_id: [0, 102, 103, 119]` mixed principal + child distributor.
- principal-owned salesman aktif.
- principal-owned salesman inactive.
- child distributor salesman aktif.
- distributor user mengirim salesman milik distributor lain.
- `emp_id` tidak ada di DB.
- duplicate `emp_id` di payload.
- `StoreSalesmen` gagal setelah header survey berhasil dibuat.

### Commands

Jalankan dari `master/`:

```bash
rtk go mod download && rtk go mod tidy
rtk go test ./service -run 'TestSurveyService_Store|TestSurveyService_Update'
rtk go test ./controller -run 'TestSurveyController_Create'
rtk go test ./...
```

## Implementation Steps

1. Reproduce root cause dengan unit test baru di `master/service/survey_service_test.go` berdasarkan payload campuran principal + child distributor.
2. Tambahkan typed error untuk invalid salesman agar tidak lagi mengandalkan `ErrSurveySalesmanNotFound` generik.
3. Audit `resolveSalesmanCustIds()`:
   - cek apakah cukup menggunakan `FindOneByEmpIdAndCustId(...)` existing;
   - bila belum cukup untuk nama/status, tambahkan helper repository read-only yang bisa mengambil detail invalid salesman.
4. Bandingkan logic scope create survey dengan logic scope `GET /master/v1/salesman`.
   - target utamanya adalah menghilangkan mismatch acceptance antara picker FE dan validator create.
5. Implementasikan validator target salesman terpusat yang:
   - menerima principal `custId`, `parentCustId`, resolved distributor scope, dan payload `emp_id`;
   - memetakan salesman valid ke `cust_id` yang benar;
   - mengumpulkan invalid IDs dan names;
   - menandai inactive sebagai invalid bila rule final tetap demikian.
6. Pastikan `distributor_id` yang mengandung `0`:
   - tidak ikut `FindCustIdsByDistributorIds(...)` sebagai ID lookup;
   - tidak ikut `FindSurveyAreasByDistributorIds(...)` sebagai distributor positif.
7. Tambahkan logging sementara terstruktur di service sekitar:
   - resolved auth context (`cust_id`, `parent_cust_id`, `user_id`)
   - input normalized `distributor_id`, `area_id`, `emp_id`
   - hasil resolved cust scope
   - daftar invalid salesman
   - raw DB/ORM error sebelum dibungkus
   Logging harus bebas token dan mudah dihapus setelah verifikasi.
8. Update controller create/update agar dapat mengenali typed invalid-salesman error dan mengembalikan body actionable.
9. Pertahankan transaction semantics existing; jika salah satu insert gagal, seluruh transaksi rollback.
10. Jalankan focused tests, lalu full `master` tests.
11. Lakukan retest local runtime utama memakai payload QA yang disanitasi dan token dev valid.
12. Tambahkan evidence perbandingan dengan hasil `GET /master/v1/salesman` untuk filter yang sama; jika ada mismatch, dokumentasikan sebagai root cause final.
13. Jika local runtime lolos dan akses tersedia, retest staging sebagai langkah lanjutan, tetapi bukan blocker utama untuk plan ini.

## Expected Files to Change

- `master/service/survey_service.go`
- `master/service/survey_service_test.go`
- `master/controller/survey_controller.go`
- `master/entity/survey.go` bila typed response error perlu struct baru
- opsional `master/repository/salesman_repository.go`
- opsional `master/controller/survey_controller_test.go`

## Agent/Tool Routing

- Implementasi: `@fixer`
- Discovery tambahan bila perlu: `@explorer`
- Review arsitektur/risk bila validator melebar ke rule tenant lain: `@oracle`
- Final signoff setelah implementasi karena menyentuh tenant scoping dan error contract: `@quality-gate`

## Validation Commands

Pre-check dari repo root:

```bash
rtk docker compose -f docker-compose.yml ps
rtk docker compose -f docker-compose.yml up -d
```

Targeted validation dari `master/`:

```bash
rtk go mod download && rtk go mod tidy
rtk go test ./service -run 'TestSurveyService_Store|TestSurveyService_Update'
rtk go test ./controller -run 'TestSurveyController_Create'
rtk go test ./...
```

Optional local smoke check setelah service hidup:

```bash
curl http://localhost:9002/ping
```

Authenticated local retest yang direncanakan:

```bash
curl 'http://localhost:9002/master/v1/survey' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <DEV_TOKEN>' \
  --data-raw '{
    "survey_title":"Testing May 12",
    "efective_date_start":"2026-05-12",
    "efective_date_end":"2026-05-13",
    "answer_frequency":"Multiple",
    "response_type":"Optional",
    "target_type":"Specific",
    "distributor_id":[0,102,103,119],
    "area_id":[91,88],
    "outlet_id":[],
    "survey_template_id":53,
    "emp_id":[450,435,415,421,458,459,466]
  }'
```

## Evidence Requirements

- Wajib ada evidence discovery lokal di `.opencode/evidence/20260513-1611-sx-1965-survey-select-all-save/discovery.md`.
- Wajib ada output test Red awal untuk skenario invalid/selection-all.
- Wajib ada output Green test targeted setelah fix.
- Wajib ada hasil `rtk go test ./...` atau penjelasan blocker yang terukur.
- Wajib ada hasil local retest payload QA atau payload turunan yang membuktikan:
  - sukses bila semua data valid; atau
  - error detail bila ada salesman invalid/inactive.
- Wajib ada evidence yang menunjukkan apakah `GET /master/v1/salesman` untuk filter FE yang sama memang mengembalikan salesman payload QA tersebut.
- Wajib ada catatan log terstruktur sementara yang menunjukkan penyebab root cause tanpa token sensitif.
- Official docs/context7, GitHub, web search, dan browser evidence tidak diperlukan untuk task ini; keputusan didasarkan pada code lokal dan data master DB.

## Done Criteria

- Tersedia regression test baru untuk principal + child distributor scope pada create survey.
- Tersedia regression test untuk sentinel `0` yang aman.
- Tersedia regression test untuk invalid/inactive salesman dengan error detail.
- Controller mengembalikan response actionable untuk invalid salesman.
- Transaction rollback tetap terjaga.
- Local retest memberi evidence yang cukup untuk menjelaskan apakah payload QA gagal karena bug scope, karena inactive salesman, atau kombinasi keduanya.
- Ringkasan root cause untuk PR/Jira bisa ditulis singkat berdasarkan hasil implementasi, misalnya:
  - validasi target salesman pada create survey tidak selaras dengan source scope dari endpoint list salesman; dan/atau
  - payload QA mengandung salesman inactive `emp_id 435`, tetapi backend hanya memberi error generik sehingga akar invalid tidak terlihat.

## Final Planning Summary

- Artifacts created:
  - `.opencode/plans/20260513-1611-sx-1965-survey-select-all-save.md`
  - `.opencode/evidence/20260513-1611-sx-1965-survey-select-all-save/discovery.md`
- Key decisions:
  - fokus validation utama di local `master` runtime dan test otomatis;
  - tambah typed invalid-salesman error dengan `invalid_emp_id` dan `invalid_salesman`;
  - pertahankan normalisasi sentinel `0` agar tidak dipakai sebagai distributor lookup positif;
  - jadikan hasil endpoint list salesman dengan filter FE yang sama sebagai acuan kompatibilitas validator create.
- Assumptions:
  - create dan list salesman idealnya harus konsisten untuk scope FE yang sama;
  - DB dev cukup representatif untuk root-cause awal.
- Open questions:
  - belum ada blocker implementasi; pertanyaan bisnis tersisa adalah bagaimana menangani kasus ketika list endpoint masih mengembalikan salesman inactive.
- Readiness:
  - siap diimplementasikan dengan TDD.
- Cleanup performed:
  - tidak membuat draft tambahan karena requirement sudah cukup jelas;
  - evidence discovery dipertahankan karena masih operasional untuk implementer.
