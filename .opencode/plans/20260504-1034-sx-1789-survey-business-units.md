# Plan — SX-1789 Survey Business Units Principal, Multi Distributor, dan Edit Salesman

## Goal

Memperbaiki backend `master` agar create/update/detail survey konsisten mendukung Business Units berupa Principal (`distributor_id = 0`), multi distributor (`> 0`), dan kombinasi keduanya, serta memastikan edit kedua tidak menghapus salesman/employee yang masih valid.

## Non-goals

- Tidak mengubah contract typo existing `efective_date_start` dan `efective_date_end`.
- Tidak mengubah UI/FE.
- Tidak menyimpan atau menambahkan secret/token dari Jira/evidence.
- Tidak melakukan redesign schema besar kecuali validasi staging membuktikan row Principal di `mst.m_survey_area` tidak aman.

## Scope

- Modul: `master/`.
- Flow utama: `POST /master/v1/survey`, `GET /master/v1/survey/{survey_id}`, `PUT/PATCH /master/v1/survey/{survey_id}` sesuai routing existing.
- Area kode target: entity survey, service survey, repository survey, tests survey, dan bila diperlukan filter pendukung `sales-teams`/`outlets` yang dipakai Manage Survey.

## Requirements

- `distributor_id = 0` harus dipertahankan sebagai pilihan Principal, bukan dibuang total.
- Distributor `> 0` tetap di-resolve ke area dari `mst.m_distributor.area_id`.
- Dedupe mapping distributor memakai pasangan `(distributor_id, area_id)`.
- Detail/edit response harus mengandung data yang cukup agar FE merender `Principal` dan nama distributor.
- Update survey harus replace mapping salesman berdasarkan `request.emp_id` unik yang valid, di dalam transaction.
- Kombinasi `0,67,68` tidak boleh membuat query filter pendukung kosong atau collapse ke satu distributor.

## Acceptance Criteria

1. Create survey dengan `distributor_id: [67, 68]` berhasil.
2. Detail multi distributor mengembalikan semua Business Units dengan nama benar.
3. Create survey dengan `distributor_id: [0]` berhasil jika Principal valid.
4. Detail Principal-only menampilkan data setara `Business Units = Principal`.
5. Create survey dengan `distributor_id: [0, 67, 68]` berhasil bila kombinasi valid.
6. Detail kombinasi mengembalikan Principal dan semua distributor.
7. Edit dari 2 business units menjadi 4 business units berhasil.
8. Edit attempt kedua tidak menghapus salesman/employee yang masih ada di `emp_id` dan valid.
9. Detail/edit response tidak collapse ke satu distributor.
10. Dedupe mapping menggunakan `(distributor_id, area_id)`, bukan hanya `area_id`.
11. Single distributor tidak regression.
12. Unit/integration test pass.

## Existing Patterns/Reuse

- Reuse `entity.FlexibleIntArray` untuk request/response ID array.
- Extend pola `normalizePositiveInts` dengan helper baru yang memisahkan Principal dan distributor positif.
- Reuse `FindSurveyAreasByDistributorIds` dan `buildSurveyAreas` untuk distributor `> 0`.
- Reuse transaction existing pada `Store` dan `Update`.
- Reuse pola `sales_team_repository.go` yang sudah memperlakukan `0` sebagai principal/global scope pada filter distributor.
- Reuse `mst.m_survey_area` sebagai storage selected business unit jika validasi DB membuktikan tidak ada FK ke `mst.m_distributor`; migration lokal mendukung `distributor_id = 0` historis.

## Constraints

- `mst.m_survey_area.area_id` saat ini `NOT NULL`, jadi Principal tidak bisa memakai `NULL` tanpa migration.
- `mst.m_survey_area.distributor_id` pada migration lokal tidak memiliki FK, tetapi staging harus dicek.
- `FindOneById` scalar distributor memakai `MIN(distributor_id)` dan join `mst.m_distributor`; row `0` tidak bisa menghasilkan nama Principal melalui join biasa.
- Response lama perlu tetap ada: `distributor_id`, `area_id`, `business_units[].distributor_id`, `business_units[].area_id`.

## Risks

- Staging schema bisa berbeda dari migration lokal; jika ada FK pada `distributor_id`, strategi row `0` harus diganti ke kolom/table baru.
- FE mungkin membaca field berbeda untuk nama Business Unit; tambahkan field kompatibel seperti `business_unit_name`, `name`, dan `type` tanpa menghapus field lama.
- `buildSurveySalesmen` menyimpan satu `cust_id` untuk semua salesman; untuk multi distributor bisa menyebabkan detail salesman hilang jika join nama memakai cust_id yang tidak sesuai.
- Soft delete lalu insert ulang bisa membuat duplicate aktif jika filter `is_del` atau transaction gagal; test perlu menangkap idempotency update kedua.

## Decisions/Assumptions

- Pertanyaan tidak diajukan karena prompt memberi acceptance criteria lengkap dan discovery lokal cukup untuk membuat rencana implementasi.
- Asumsi utama: Principal valid sebagai selected Business Unit dan boleh dipersist di `mst.m_survey_area` dengan `distributor_id = 0` serta `area_id` dari payload yang sudah dipilih, karena schema lokal tidak memiliki FK dan migration `002` pernah mengisi `0`.
- Jika staging memiliki FK/constraint yang menolak `0`, ubah keputusan menjadi migration baru untuk menyimpan flag/relasi Principal secara eksplisit, misalnya `mst.m_survey_business_unit` atau kolom/flag yang disetujui tim DB.
- Untuk Principal-only, `resolveSurveyCustIds` tetap harus memvalidasi salesman terhadap principal/current cust scope, bukan mengirim lookup distributor `0` ke `smc.m_customer`.
- Untuk response, `area_id` Principal dapat memakai area payload/mapping aktif karena contract existing `area_id` non-null; FE harus memakai `type/name` untuk label Principal.

## TDD/Test Plan

- TDD required: ya, karena ini bug produksi pada API behavior, mapping persistensi, response detail, dan update transaction.
- Existing test patterns:
  - `master/service/survey_service_test.go` memakai `surveyRepositoryRedStub`, `salesmanRepositoryStub`, `transactionManagerStub`.
  - Test existing sudah ada untuk multi distributor dan zero distributor yang saat ini mengharapkan `0` diabaikan; test ini perlu diubah sesuai requirement baru.
  - `master/controller/survey_controller_test.go` sudah memvalidasi response detail basic.
- First failing/regression tests:
  1. `TestSurveyService_Store_ShouldPersistPrincipalBusinessUnit` dengan `DistributorId: {0}`, `AreaId: {82}`; assert `StoreAreas` berisi row `{DistributorId: 0, AreaId: 82}` dan tidak memanggil lookup distributor untuk `0`.
  2. `TestSurveyService_Detail_ShouldReturnPrincipalBusinessUnitName` dengan `FindAreasBySurveyId` mengembalikan row Principal; assert `DistributorId` mengandung `0`, `BusinessUnits` mengandung `name/business_unit_name = Principal`, `type = principal`.
  3. `TestSurveyService_Store_ShouldPersistPrincipalAndDistributors` dengan `{0,67,68}`; assert row Principal + `{67,82}` + `{68,82}` tanpa duplicate.
  4. `TestSurveyService_Update_ShouldKeepSalesmenOnSecondEdit` dengan update payload sama dua kali; assert delete+store menulis semua unique `emp_id` yang dikirim.
- Green step:
  - Implement helper normalisasi Business Unit.
  - Extend model/entity response Business Unit dengan fields nama/type kompatibel.
  - Persist row Principal atau fallback strategi schema bila staging menolak row `0`.
  - Fix salesman mapping agar unique request `emp_id` disimpan ulang lengkap dan detail tidak hilang karena cust_id salah.
- Refactor step:
  - Hindari duplikasi logic Store/Update dengan helper reusable: normalize, build survey areas, build business unit response, build survey salesmen.
  - Pastikan helper kecil, deterministic, dan mudah dites.
- Edge cases:
  - `distributor_id: []`, `area_id: []` tetap kompatibel general survey.
  - `distributor_id: [0]` dengan `area_id: []` perlu keputusan: jika Principal wajib area, return validation error; jika tidak, gunakan sentinel area sesuai schema. Karena QA payload Principal menyertakan area, plan awal validasi area wajib untuk Principal mapping.
  - Duplicate input `[0,0,67,67]` tidak menghasilkan duplicate mapping.
  - Distributor berbeda dengan area sama tetap muncul dua item.
  - Negative ID diabaikan/ditolak sesuai pattern existing; jangan lookup negatif.
- Commands:
  - `rtk go test ./service -run 'TestSurveyService_(Store|Detail|Update)'`
  - `rtk go test ./controller -run TestSurveyController`
  - `rtk go test ./...`

## Implementation Steps

1. Tambahkan/ubah entity response Business Unit.
   - Extend `entity.SurveyBusinessUnit` dengan field opsional yang tetap backward compatible:
     - `BusinessUnitName string json:"business_unit_name,omitempty"`
     - `Name string json:"name,omitempty"`
     - `Type string json:"type,omitempty"`
   - Pertahankan `DistributorId` dan `AreaId`.
2. Tambahkan helper normalisasi selected business units di `survey_service.go`.
   - Return `HasPrincipal bool` dan `DistributorIDs []int`.
   - Dedupe `0` dan distributor positif secara stabil.
   - Jangan kirim `0` ke `FindSurveyAreasByDistributorIds` atau `FindCustIdsByDistributorIds`.
3. Ubah `Store` dan `Update`.
   - Gunakan helper baru.
   - Build mapping distributor positif via `FindSurveyAreasByDistributorIds` + `buildSurveyAreas`.
   - Jika `HasPrincipal`, append mapping Principal ke `surveyAreas` dengan `DistributorId: 0` dan area dari `request.AreaId` yang sudah dinormalisasi.
   - Dedupe final mapping berdasarkan `(distributor_id, area_id)`.
   - Jika Principal dipilih tanpa area dan schema tidak mendukung `area_id NULL`, return validation error yang jelas atau ikuti convention existing jika ditemukan.
4. Perbaiki detail response.
   - Saat loop `FindAreasBySurveyId`, jika `DistributorId == 0`, append `0` ke `response.DistributorId` dan `SurveyBusinessUnit{DistributorId: 0, AreaId: a.AreaId, BusinessUnitName: "Principal", Name: "Principal", Type: "principal"}`.
   - Untuk distributor `>0`, isi `Type: "distributor"` dan nama dari query repository.
   - Perlu extend `model.SurveyArea` joined fields dengan `DistributorName *string` dan query `FindAreasBySurveyId` join ke `mst.m_distributor`.
   - Hindari penggunaan scalar `FindOneById` untuk menentukan daftar Business Units.
5. Perbaiki query `FindAreasBySurveyId`.
   - Join `mst.m_area` tetap.
   - Tambah `LEFT JOIN mst.m_distributor md ON md.distributor_id = sa.distributor_id AND md.is_del = false`.
   - Select `md.distributor_name` sebagai joined field untuk item distributor.
   - Jangan inner join agar row Principal `0` tetap keluar.
6. Perbaiki salesman/employee mapping update.
   - Buat helper `normalizeUniqueInts` untuk `emp_id` agar delete-then-insert menulis semua `emp_id` unik.
   - Review `resolveSurveyCustIds`: jika Principal dipilih, cust scope harus mencakup principal/current parent scope, bukan hanya child cust dari distributor positif.
   - Jika multi distributor menghasilkan beberapa cust_id, jangan simpan semua salesman dengan satu default cust_id bila detail lookup membutuhkan cust_id aktual. Tambahkan repository lookup `FindSalesmanCustIdsByEmpIds(parentCustId/custScope, empIds)` atau ubah `buildSurveySalesmen` agar memilih cust_id valid per salesman.
   - Minimal green untuk defect: delete+insert ulang semua `request.emp_id` unik dalam transaction, dengan cust_id yang membuat join detail berhasil.
7. Review filter pendukung.
   - `sales_team_repository.go` sudah punya pattern `0` sebagai principal/global scope; tambahkan/pertahankan test untuk `0,67,68`.
   - Cek path outlet yang dipakai endpoint Manage Survey; update helper filter agar `0` tidak menjadi `mc.distributor_id IN (0,67,68)` murni yang mempersempit hasil. Pola yang disarankan: `0` menambahkan principal/global scope, ID positif menambahkan filter distributor.
8. Update tests.
   - Ubah test lama `ShouldIgnoreZeroDistributor...` menjadi test baru yang memastikan `0` tidak dilookup sebagai distributor tetapi tetap dipersist/dikembalikan sebagai Principal.
   - Tambah tests detail names dan update idempotency.
   - Tambah controller response test bila contract JSON business_units berubah.
9. Manual verification staging/local.
   - Jalankan create/detail Principal-only, Principal+distributor, multi distributor, update dua kali.
   - Jalankan query DB mapping survey area dan survey salesman sesuai nama tabel aktual.

## Expected Files to Change

- `master/entity/survey.go`
- `master/model/survey_area.go`
- `master/service/survey_service.go`
- `master/repository/survey_repository.go`
- `master/service/survey_service_test.go`
- `master/controller/survey_controller_test.go`
- Kemungkinan: `master/repository/sales_team_repository_test.go`
- Kemungkinan: `master/repository/outlet_repository.go` dan test terkait outlet filter, bila path Manage Survey memakai helper yang belum mendukung `0,67,68`.
- Kemungkinan migration baru di `master/migration/mst.survey/` hanya jika staging menolak strategi row `distributor_id = 0`.

## Agent/Tool Routing

- Implementasi: gunakan `@fixer` atau build agent dengan TDD Red → Green → Refactor.
- Review risiko arsitektur/schema bila staging punya FK berbeda: gunakan `@oracle`.
- Discovery tambahan kode: gunakan `@explorer` bila perlu mencari path outlet/sales-team spesifik.
- Jangan gunakan designer/browser visual karena ini backend API defect.

## Validation Commands

Jalankan dari `master/`:

```bash
rtk go test ./service -run 'TestSurveyService_(Store|Detail|Update)'
rtk go test ./controller -run TestSurveyController
rtk go test ./repository -run 'Test(SalesTeam|Outlet|Survey)'
rtk go test ./...
```

Validasi service sebelum code work dari repo root sudah dilakukan:

```bash
rtk docker compose -f docker-compose.yml ps
```

## Evidence Requirements

- Automated:
  - Output `rtk go test ./...` dari `master/`.
  - Output test spesifik survey service/controller.
- Manual API:
  - `POST /master/v1/survey` Principal-only lalu `GET detail` menampilkan Principal.
  - `POST /master/v1/survey` Principal+distributor lalu `GET detail` menampilkan semua Business Units.
  - `PUT/PATCH` update dua kali dengan `emp_id` sama lalu detail tetap menampilkan semua salesman.
- DB:
  - `SELECT survey_id, distributor_id, area_id, is_del FROM mst.m_survey_area WHERE survey_id = :survey_id ORDER BY distributor_id, area_id;`
  - `SELECT survey_id, salesman_id, is_del FROM mst.m_survey_salesman WHERE survey_id = :survey_id ORDER BY salesman_id;`
- Schema gate sebelum deploy:
  - Pastikan `mst.m_survey_area.distributor_id` tidak punya FK yang menolak `0`; bila ada, gunakan strategi migration/flag alternatif.

## Done Criteria

- Semua acceptance criteria terpenuhi.
- Tests survey service/controller/repository pass.
- Manual API dan DB evidence tersedia.
- Response lama tetap kompatibel dan response baru cukup untuk FE render `Principal`.
- Tidak ada secret/token ditambahkan ke repo.

## Final Planning Summary

- Artifacts created/consulted:
  - Primary plan: `.opencode/plans/20260504-1034-sx-1789-survey-business-units.md`
  - Evidence kept: `.opencode/evidence/20260504-1034-sx-1789-survey-business-units/discovery.md` karena berguna sebagai handoff implementasi dan audit discovery.
- Key decisions:
  - Pisahkan selected Business Unit menjadi `HasPrincipal` dan distributor positif.
  - Jangan lookup `0` ke `mst.m_distributor`, tetapi jangan hilangkan dari selected Business Units.
  - Strategi awal persist Principal: row `mst.m_survey_area` dengan `distributor_id = 0` dan `area_id` dari payload, dengan schema gate staging.
  - Detail response harus inject nama/type Principal dan nama distributor dari join distributor.
  - Update salesman harus berbasis `request.emp_id` unik dan transaction-safe.
- Assumptions:
  - Principal merupakan pilihan valid untuk survey.
  - Area tetap wajib/tersedia untuk Principal sesuai payload QA karena schema `area_id NOT NULL`.
  - Staging schema sejalan dengan migration lokal atau akan divalidasi sebelum deploy.
- Open questions:
  - Tidak ada blocker untuk rencana; satu schema gate tersisa: apakah staging punya FK pada `mst.m_survey_area.distributor_id`.
- Readiness:
  - Siap untuk implementasi TDD oleh build/fixer agent.
- Cleanup performed:
  - Draft artifact tidak diperlukan; durable findings dikonsolidasikan ke primary plan dan discovery evidence.
