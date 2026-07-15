# SX-2124 — Business Unit pakai `cust_name`

Task id: `20260602-1111-sx-2124-business-unit-cust-name`

## Goal

Perbaiki `GET /master/v1/business-unit` pada principal path agar dropdown display memakai `smc.m_customer.cust_name` dari current principal `cust_id`, bukan `sys.m_user.user_fullname` atau JWT user name.

## Non-goals

- Tidak ubah endpoint selain `GET /master/v1/business-unit`.
- Tidak ubah SX-2079 scope enforcement untuk region, area, dan distributor mapping.
- Tidak ubah distributor-user path kecuali test menunjukkan regresi langsung.
- Tidak hardcode `C26002` atau `PT. Madura Sejahtera`.
- Tidak memakai token/password/header auth dari Jira.

## Scope

Target module: `master`.

File utama kemungkinan berubah:

- `master/entity/business_unit.go`
- `master/model/business_unit.go`
- `master/repository/business_unit_repository.go`
- `master/service/business_unit_service.go`
- `master/service/business_unit_service_test.go`
- `master/repository/business_unit_repository_test.go` bila query helper ditambah test.

## Requirements

- Principal response `data.user_fullname` harus berisi customer name dari `smc.m_customer.cust_name` berdasarkan current `cust_id`.
- Tambahkan `data.cust_name` pada `BusinessUnitPrincipalResponse` bila aman, additive dan backward-compatible.
- `data.user_id` tetap dari `sys.m_user` agar contract lama tidak rusak.
- `distributor_data[].distributor_name` tetap nama distributor/business unit.
- Query params tetap: `q`, `page`, `limit`, `sort`, `is_active`, `region_id`, `area_id`, `distributor_id`.
- Principal scope SX-2079 tetap: `specific` mapping, all fallback, explicit filter intersection-safe.
- Missing customer row mengikuti perilaku repository sekarang: error naik ke controller; `sql.ErrNoRows` jadi `404 RECORD_NOT_FOUND`.

## Acceptance Criteria

- Principal user `cust_id=C26002`, `user_fullname=Princessa Ahsani Taqwim` menerima `data.user_fullname = "PT. Madura Sejahtera"`.
- Principal user `cust_id=C26002`, `user_fullname=Agung Citra` menerima `data.user_fullname = "PT. Madura Sejahtera"`.
- Bila `cust_name` field ditambah, `data.cust_name = "PT. Madura Sejahtera"`.
- Request tanpa `region_id` / `area_id` dan request dengan `region_id=84&area_id=96` sama-sama mempertahankan customer display.
- Distributor list tetap terfilter sesuai scope SX-2079 dan tidak duplikat.
- Distributor-user response tetap pakai perilaku lama.
- Test service membuktikan principal tidak lagi memakai `userInfo.UserFullname` sebagai display.
- Tidak ada secret/token/password masuk fixture, log, atau code.

## Existing Patterns/Reuse

- Reuse `BusinessUnitRepository` untuk data access; service tidak query DB langsung.
- Reuse `repo.Get(&model, query, args...)` seperti `FindUserByUsername`.
- Reuse `BusinessUnitQueryFilter.CustId` dari JWT sebagai lookup key principal.
- Reuse current `FindEmployeeDropdownScope` dan `NormalizeScopeSet`; jangan gabungkan customer lookup ke employee scope dulu karena region/area services memakai interface sama.
- Reuse test stub `businessUnitRepositoryStub`, tambahkan method customer lookup.

Tidak ditemukan util KiloCode/project yang langsung menyelesaikan lookup customer name untuk endpoint ini. Pola terdekat ada di repository lain yang select `cust_id, cust_name, parent_cust_id` dari `smc.m_customer`.

## Constraints

- Ikuti Controller → Service → Repository → DB.
- Tenant rule: principal display lookup pakai `dataFilter.CustId`; distributor list all-scope fallback tetap pakai `ParentCustId` bila ada.
- `master` adalah Go module mandiri; jalankan command dari `master/`.
- Repo guidance lokal minta command shell pakai `rtk`.
- Jangan ubah source file di luar planning mode sampai handoff ke `@orchestrator` / `@fixer`.

## Risks

- Field `cust_name` additive biasanya aman, tapi FE strict parser bisa sensitif. Mitigasi: tetap isi `user_fullname` dengan customer name.
- Missing `smc.m_customer` row bisa menyebabkan 404, bukan fallback user name. Ini memperlihatkan bad master data; sejalan dengan source-of-truth requirement.
- Menambah lookup terpisah menambah satu query per principal request. Risiko kecil; endpoint dropdown already DB-bound. Bisa optimasi nanti dengan join scope bila perlu.
- `sort` saat ini interpolasi raw column/direction; bukan bagian SX-2124, tapi jangan perluas risiko. Catat untuk security follow-up bila quality gate minta.

## Decisions/Assumptions

- Keputusan: pakai `dataFilter.CustId` untuk principal customer lookup.
- Keputusan: `data.user_fullname` tetap ada dan nilainya diganti ke `cust_name` untuk principal path.
- Keputusan: tambahkan `cust_name` ke `BusinessUnitPrincipalResponse` sebagai field additive bila implementer setuju dengan backward-compatible schema.
- Asumsi: distributor path tetap memakai `sys.m_user.user_fullname` karena issue hanya principal dropdown top-level.
- Asumsi: no-row customer harus error, bukan fallback ke user name.
- Pertanyaan tidak ditanyakan karena prompt sudah memberi expected behavior, ID source preference, dan backward-compatible response direction.

## TDD/Test Plan

TDD wajib karena ini bug produksi API behavior.

Existing test patterns:

- `master/service/business_unit_service_test.go` memakai stub repo + scope repo.
- `master/repository/business_unit_repository_test.go` menguji query builder SX-2079.

Red step:

- Tambah test service: principal `userInfo.UserFullname = "Princessa Ahsani Taqwim"`, repository customer name `"PT. Madura Sejahtera"`, expected `resp.UserFullname == "PT. Madura Sejahtera"` dan bila field ada `resp.CustName == "PT. Madura Sejahtera"`.
- Tambah test service P2 dengan `userInfo.UserFullname = "Agung Citra"`, same `cust_id`, same expected customer name.
- Tambah assertion bahwa lookup customer dipanggil dengan `C26002`.
- Tambah test missing customer row: repository returns `sql.ErrNoRows`, service returns error.

Green step:

- Tambah model `CustomerInfo` atau method return string.
- Tambah interface method `FindCustomerNameByCustId(custId string) (string, error)`.
- Implement query `SELECT cust_name FROM smc.m_customer WHERE cust_id = $1 LIMIT 1`.
- Pada principal path, setelah scope lookup atau sebelum response, ambil customer name lalu set `UserFullname` ke customer name.
- Tambah `CustName string \`json:"cust_name"\`` ke `BusinessUnitPrincipalResponse` bila dipakai.

Refactor step:

- Pastikan mapping distributor tetap tidak berubah.
- Hindari duplicate empty string normalization logic.
- Pastikan test names jelas merujuk SX-2124.

Edge cases:

- `cust_name` empty string: response empty string bila DB berisi empty, jangan fallback user fullname.
- `cust_id` missing: customer lookup gagal; service error.
- Distributor user: customer lookup tidak dipanggil.
- Principal with `region_id` / `area_id`: customer display tetap sama, filter tetap diteruskan.

Commands:

```bash
rtk go test ./service -run 'TestBusinessUnitService_GetBusinessUnit'
rtk go test ./repository -run 'TestBuildFindDistributorsByCustIDQuery'
rtk go test ./...
```

Jalankan dari `master/`.

## Implementation Steps

1. Dari `master/`, jalankan targeted failing test setelah menambah test SX-2124.
2. Update `master/model/business_unit.go` bila butuh `CustomerInfo`.
3. Update `BusinessUnitRepository` interface dengan `FindCustomerNameByCustId`.
4. Implement `FindCustomerNameByCustId` di `business_unit_repository.go` memakai `smc.m_customer`.
5. Update service principal path agar mengambil customer name dan memakai itu untuk `UserFullname`.
6. Tambah `CustName` field di `BusinessUnitPrincipalResponse` bila additive field diterima.
7. Update test stub dan service tests.
8. Jalankan targeted tests, lalu full `rtk go test ./...` di `master/`.
9. Manual smoke local/staging dengan token fresh, tidak disimpan di repo.
10. Quality gate review untuk SX-2079 regression dan secret hygiene.

## Expected Files to Change

- `master/entity/business_unit.go`
- `master/model/business_unit.go`
- `master/repository/business_unit_repository.go`
- `master/service/business_unit_service.go`
- `master/service/business_unit_service_test.go`

Opsional:

- `master/repository/business_unit_repository_test.go`
- `master/controller/business_unit_controller_test.go` bila controller contract perlu JSON assertion.

## Agent/Tool Routing

- `@orchestrator`: jalankan handoff dan integrasi.
- `@fixer`: implementasi bounded + tests Red → Green → Refactor.
- `@quality-gate`: final review karena bug staging, tenant/data-source behavior, dan SX-2079 regression risk.
- `@explorer`: opsional bila implementer butuh discovery tambahan.
- `@architect`: tidak diperlukan; perubahan kecil dan contract jelas.
- `@librarian`/Context7/GitHub/Web: tidak diperlukan; tidak ada dependency external baru.

## Execution-ready Worklist / Handoff Contract

`start_with`: `T1`

| id | action | depends_on | owner/lane | validation | exit criteria | status | requires_user_decision |
|---|---|---|---|---|---|---|---|
| T1 | Tambah failing service tests SX-2124 untuk principal customer display, dua user fullname berbeda, same `cust_id` | none | `@fixer` | `rtk go test ./service -run 'TestBusinessUnitService_GetBusinessUnit'` | Test gagal karena response masih pakai `sys.m_user.user_fullname` | ready | no |
| T2 | Tambah repository contract `FindCustomerNameByCustId` dan stub support | T1 | `@fixer` | compile targeted test | Stub mencatat `cust_id` lookup dan bisa return customer name/error | ready | no |
| T3 | Implement query customer name dari `smc.m_customer` | T2 | `@fixer` | `rtk go test ./repository -run 'TestBuildFindDistributorsByCustIDQuery'` plus compile | Repository punya method kecil, parameterized, no hardcode | ready | no |
| T4 | Update principal response mapping agar `user_fullname` dan opsional `cust_name` memakai customer name | T3 | `@fixer` | `rtk go test ./service -run 'TestBusinessUnitService_GetBusinessUnit'` | Tests P1/P2/missing-customer/distributor-no-lookup pass | ready | no |
| T5 | Pastikan SX-2079 query builder tetap tidak berubah/regress | T4 | `@fixer` | `rtk go test ./repository -run 'TestBuildFindDistributorsByCustIDQuery'` | Scope mapping, `DISTINCT`, `IN (?)`, parent fallback tests pass | ready | no |
| T6 | Jalankan full master tests | T5 | `@fixer` | `rtk go test ./...` | Semua test `master` pass atau failure unrelated terdokumentasi | ready | no |
| T7 | Smoke API local/staging dengan token fresh | T6 | `@orchestrator` | `curl` endpoint dengan dan tanpa `region_id`/`area_id` | `data.user_fullname` menjadi customer name; distributor rows tetap benar | ready | no |
| T8 | Final quality gate | T7 | `@quality-gate` | Review diff + test evidence | No secret leak, no SX-2079 regression, acceptance criteria met | ready | no |

## Validation Commands

Dari `master/`:

```bash
rtk go test ./service -run 'TestBusinessUnitService_GetBusinessUnit'
rtk go test ./repository -run 'TestBuildFindDistributorsByCustIDQuery'
rtk go test ./...
```

Manual smoke, token fresh dari env/local/staging:

```bash
curl -H "Authorization: Bearer <TOKEN_PRINCIPAL>" \
  "https://best.scyllax.online/master/v1/business-unit?is_active=1&q=&page=1&limit=99"

curl -H "Authorization: Bearer <TOKEN_PRINCIPAL>" \
  "https://best.scyllax.online/master/v1/business-unit?is_active=1&region_id=84&area_id=96&q=&page=1&limit=99"
```

DB smoke bila akses tersedia:

```sql
SELECT cust_id, cust_name
FROM smc.m_customer
WHERE cust_id = 'C26002';
```

## Evidence Requirements

- Test output targeted service.
- Test output query builder repository.
- Full `rtk go test ./...` output atau daftar failure unrelated.
- Manual curl response sanitized tanpa token/header secret.
- DB check result sanitized bila dilakukan.
- Diff review memastikan `sys.m_user.user_fullname` tidak lagi jadi display source principal.

## Done Criteria

- Primary acceptance criteria pass.
- Tests pass or unrelated failures documented.
- No secrets committed.
- SX-2079 regression tests pass.
- `@quality-gate` approve atau issues resolved.

## Final Planning Summary

Artefak dibuat:

- `.opencode/plans/20260602-1111-sx-2124-business-unit-cust-name.md`
- `.opencode/evidence/20260602-1111-sx-2124-business-unit-cust-name/discovery.md`
- `.opencode/evidence/20260602-1111-sx-2124-business-unit-cust-name/index.json`

Keputusan kunci:

- Principal display lookup pakai `smc.m_customer.cust_name` by current `cust_id`.
- `user_fullname` tetap ada sebagai backward-compatible FE display field, nilainya diganti ke customer name pada principal path.
- `cust_name` additive field direkomendasikan.
- Distributor path tidak disentuh.

Asumsi:

- Missing customer row error, bukan fallback ke user fullname.
- `cust_id` adalah source ID benar untuk principal; `parent_cust_id` tetap untuk distributor scope fallback.

Open questions:

- Tidak ada blocker. FE confirmation untuk `cust_name` additive bagus, tapi tidak wajib karena `user_fullname` tetap dipertahankan.

Readiness:

- Siap implementasi oleh `@orchestrator` → `@fixer` tanpa replanning.

Cleanup:

- Tidak ada draft dibuat.
- Evidence discovery dipertahankan karena operasional berguna untuk implementer dan quality gate.
