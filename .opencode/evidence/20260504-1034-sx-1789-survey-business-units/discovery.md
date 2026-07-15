# Discovery — SX-1789 Survey Business Units

## Files inspected

- `master/entity/survey.go`
- `master/service/survey_service.go`
- `master/repository/survey_repository.go`
- `master/model/survey.go`
- `master/model/survey_area.go`
- `master/model/survey_salesman.go`
- `master/service/survey_service_test.go`
- `master/migration/mst.survey/001_create_tables.sql`
- `master/migration/mst.survey/002_add_distributor_and_salesman.sql`
- `master/entity/business_unit.go`
- `master/repository/sales_team_repository.go`
- `master/repository/outlet_repository.go`

## Commands/docs checked

- `rtk docker compose -f docker-compose.yml ps`
  - `master`, `system`, dan `redis` sedang berjalan.
- Local discovery via `Glob`, `Grep`, dan `Read`.
- Official docs/GitHub/web/browser tidak diperlukan karena defect berada pada logic lokal Go + SQL existing, bukan perilaku library eksternal atau UI visual.

## Project patterns found

- Modul target adalah `master/`, module Go terpisah dengan service/repository/test sendiri.
- Alur survey mengikuti `controller -> service -> repository -> DB`.
- `Store` dan `Update` survey sudah memakai `txManager.WithinTransaction`.
- `distributor_id` request memakai `entity.FlexibleIntArray`, mendukung scalar atau array.
- Perbaikan sebelumnya memakai:
  - `normalizePositiveInts([]int(request.DistributorId))`, sehingga `0` selalu dibuang.
  - `FindSurveyAreasByDistributorIds(distributorIds)` untuk resolve `(distributor_id, area_id)` dari `mst.m_distributor`.
  - `buildSurveyAreas(...)` dedupe berdasarkan pasangan `[2]int{distributor_id, area_id}`.
- Detail survey membangun `DistributorId`, `AreaId`, dan `BusinessUnits` dari `FindAreasBySurveyId`.
- `SurveyBusinessUnit` saat ini hanya berisi `distributor_id` dan `area_id`; belum ada nama/label/type.

## Reuse candidates

- Reuse `FlexibleIntArray` untuk contract `distributor_id` array.
- Extend `normalizePositiveInts` dengan helper baru yang tidak membuang sinyal Principal, misalnya `normalizeSurveyBusinessUnits`.
- Reuse `buildSurveyAreas` untuk distributor `> 0`, lalu tambahkan row Principal secara eksplisit.
- Reuse transaction create/update yang sudah ada.
- Reuse sales team principal handling pattern: `buildSalesTeamCustScopeCondition` sudah memperlakukan `distributor_id = 0` sebagai principal/global scope dan `>0` sebagai distributor scope.
- Reuse schema `mst.m_survey_area` untuk row Principal jika tidak ada FK ke `mst.m_distributor`; migration lokal menunjukkan `distributor_id INT NOT NULL` tanpa FK, dan migration `002` pernah mengisi `distributor_id = 0` untuk row lama.

## Constraints

- Jangan mengubah typo API existing `efective_date_start` / `efective_date_end`.
- Jangan breaking field response lama; jika menambah field nama/type, tetap pertahankan `distributor_id` dan `area_id`.
- `mst.m_survey_area.area_id` adalah `INT NOT NULL`; Principal tidak cocok memakai `NULL` tanpa migration schema.
- `FindOneById` masih memakai subquery `MIN(distributor_id)` lalu join ke `mst.m_distributor`; jika ada row Principal `0`, join distributor name tidak akan ditemukan. Detail utama perlu membaca business units dari `FindAreasBySurveyId`, bukan mengandalkan scalar distributor pada `FindOneById`.
- `buildSurveySalesmen` saat ini menyimpan semua salesman dengan satu `defaultCustId` dari `resolveSurveyCustIds`; ini berisiko untuk multi distributor/Principal karena join detail salesman memakai `ss.cust_id`.

## Risks

- Jika environment staging memiliki FK tersembunyi pada `mst.m_survey_area.distributor_id`, row `0` bisa gagal; perlu validasi schema staging sebelum deploy.
- Jika FE hanya membaca `distributor_name` scalar top-level, response tambahan `business_units[].name` perlu dikomunikasikan; namun field lama tetap dijaga.
- Employee/salesman bisa terlihat hilang jika `cust_id` mapping disimpan dengan satu distributor cust_id sementara salesman milik cust_id lain.
- Endpoint outlet/sales-team punya beberapa path/filter; perlu test spesifik untuk path yang dipakai Manage Survey.

## Research gate decision

- Local project discovery: diperlukan dan sudah dilakukan.
- Official docs/context7: tidak diperlukan; tidak ada dependency/library baru.
- GitHub: tidak diperlukan; tidak bergantung upstream repo/issue.
- Brave/web search: tidak diperlukan; semua fakta berasal dari Jira prompt dan kode lokal.
- Browser/screenshot: tidak diperlukan; ini defect backend API, bukan visual parity.
