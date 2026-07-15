# Plan — SX-1915 Salesman Business Unit Filter

Task ID: `20260504-2058-sx-1915-salesman-business-unit`
Tanggal: 2026-05-04 Asia/Jakarta
Primary source of truth: `.opencode/plans/20260504-2058-sx-1915-salesman-business-unit.md`

## Goal

Memperbaiki endpoint `GET /master/v1/salesman` agar dropdown Salesman pada flow Manage Survey dapat mengembalikan salesman principal dan distributor secara benar saat FE mengirim `distributor_id=0,<distributor_id>`, khususnya kasus SX-1915 `distributor_id=0,120`, `sales_team_id=79,76,73,70`, dan pencarian `rahmat`.

## Non-goals

- Tidak mengubah kontrak FE atau mengganti parameter `distributor_id=0`.
- Tidak melakukan refactor besar query raw SQL ke parameterized query builder di luar kebutuhan bugfix.
- Tidak mengubah business rule survey creation, sales team, distributor, warehouse, atau Master Salesman di luar list salesman.
- Tidak melakukan perubahan migrasi database.
- Tidak menggunakan kredensial Jira kecuali implementer perlu reproduksi manual di staging dan diizinkan oleh workflow tim.

## Scope

Module target: `master`.

Endpoint target:

- External/proxy: `GET /master/v1/salesman`
- Internal service route: `GET /v1/salesman` di `master/controller/salesman_controller.go`

Area kode utama:

- Parsing query `distributor_id` dan `sales_team_id`.
- Helper repository salesman untuk cust scope principal/distributor.
- Unit test helper scope dan parser.
- Validasi manual/API untuk principal-only, distributor-only, mixed principal+distributor, pagination/total, dan no duplicate.

## Requirements

1. `distributor_id=0` harus berarti include salesman milik principal/current parent customer scope.
2. `distributor_id=<non-zero>` harus tetap berarti include salesman dari distributor child customer terkait, dibatasi `parent_cust_id`.
3. `distributor_id=0,120` harus menghasilkan union scope principal + distributor `120`.
4. `sales_team_id` harus tetap berlaku untuk semua branch principal dan distributor.
5. `q`/search harus tetap berlaku untuk semua branch principal dan distributor.
6. Status aktif dan soft delete existing harus tetap dihormati.
7. Tidak boleh ada salesman dari tenant/business unit yang tidak dipilih.
8. Response dan total pagination harus tidak duplicate.

## Acceptance Criteria

1. `/master/v1/salesman?...&distributor_id=0&q=rahmat&sales_team_id=79&limit=9999` mengembalikan salesman principal yang cocok, termasuk `Rahmat / EMP0012` jika data tersedia dan sales team cocok.
2. `/master/v1/salesman?...&distributor_id=120&sales_team_id=79&limit=9999` mengembalikan salesman distributor `120` saja.
3. `/master/v1/salesman?...&distributor_id=0,120&q=rahmat&sales_team_id=79,76,73,70&limit=9999` mengembalikan union principal + distributor dan menemukan `Rahmat / EMP0012`.
4. Filter `sales_team_id=79,76,73,70` berlaku untuk semua hasil.
5. Search `q` berlaku untuk semua hasil sesuai behavior existing.
6. Salesman business unit lain tidak ikut muncul.
7. Tidak ada duplicate salesman pada response.
8. `total_record` dan `page_total` sesuai dataset setelah filter final.
9. Regression scenario lolos: no distributor filter/default, principal-only, distributor-only, multiple distributors, mixed principal+distributor.

## Existing Patterns/Reuse

Discovery menunjukkan pola reusable sudah ada:

- `master/controller/query_filter_parser.go`
  - `parseIntSliceQueryAllowZero` sudah mempertahankan nilai `0` untuk `distributor_id`.
  - Test parser di `master/controller/salesman_controller_test.go` sudah memastikan `0` tidak dibuang.
- `master/repository/sales_team_repository.go`
  - `buildSalesTeamCustScopeCondition` sudah memodelkan sentinel `0` sebagai principal scope dan id non-zero sebagai distributor child scope, lalu menggabungkan branch dengan `OR` dalam grouping.
- `master/repository/sales_team_repository_test.go`
  - Menyediakan pola test untuk principal-only dan mixed principal+distributor.

Strategi reuse: adaptasi pola `buildSalesTeamCustScopeCondition` ke `buildSalesmanCustScopeCondition` dengan alias `s`.

## Constraints

- Tenant scoping wajib menggunakan `cust_id`/`parent_cust_id`; jangan membuka data lintas tenant.
- Untuk distributor non-zero, tetap map via `smc.m_customer` dengan `mc.parent_cust_id = scopeParentCustID` dan `mc.distributor_id IN (...)`.
- Untuk principal sentinel `0`, gunakan `scopeParentCustID` sebagai cust scope principal. Jika `parentCustId` kosong, fallback ke `custId`, mengikuti pola `sales_team_repository.go`.
- Kondisi principal/distributor harus digroup sebagai `( principal_condition OR distributor_condition )` sebelum ditambah `AND s.is_del = false`, `AND s.sales_team_id IN (...)`, `AND s.sales_name ILIKE ...`, dan filter lain.
- Perintah shell dalam repo ini mengikuti instruksi repo lokal dengan `rtk` prefix.
- Jangan menambah atau mengekspos secrets dari file compose/workflow.

## Risks

- Helper repository saat ini memakai SQL string concatenation; perubahan harus kecil dan tidak memperluas input surface. Karena `distributor_id` diparse sebagai int dan dedupe, risiko injeksi untuk id rendah, tetapi `SalesTeamId`, `Sort`, dan `Query` existing tetap raw string dan di luar scope utama.
- Test existing `TestBuildSalesmanCustScopeCondition_IgnoresZeroWhenDistributorFilterHasValidValues` mengunci perilaku lama yang salah; test ini harus diubah, bukan dipertahankan.
- `FindAllByCustId` dan `FindAllByCustIdLookup` menggunakan helper yang sama; perbaikan helper memengaruhi kedua mode. Ini diinginkan, tetapi perlu regression test.
- Count saat ini `COUNT(*)`; jika join menghasilkan duplicate row, acceptance no duplicate/total benar mungkin memerlukan `DISTINCT` tambahan.
- Service `List` memanggil `FindDetailById(row.EmpId, custId)` memakai `custId` request, bukan `row.CustId`; jika hasil mixed mencakup salesman child distributor, detail mungkin gagal atau kosong. Implementer harus memverifikasi behavior aktual; jika error muncul pada mixed list, rencanakan fix kecil untuk memakai `row.CustId` pada detail lookup.

## Decisions/Assumptions

- Keputusan: implementasi dimulai dari helper `buildSalesmanCustScopeCondition`, bukan controller, karena parser sudah benar dan akar issue berada pada scope condition yang membuang `0` saat ada distributor non-zero.
- Keputusan: principal scope untuk sentinel `0` adalah `s.cust_id = scopeParentCustID`, mengikuti pola sales team dan konteks user principal `cust_id = parent_cust_id = C26004`.
- Keputusan: distributor non-zero tetap memakai subquery `smc.m_customer` agar tenant dibatasi oleh `parent_cust_id`.
- Asumsi: `Rahmat / EMP0012` tersimpan sebagai salesman pada `mst.m_salesman` dengan `s.cust_id` principal `C26004` dan `s.sales_team_id` salah satu dari `79,76,73,70`.
- Asumsi: existing search by `s.sales_name ILIKE '%q%'` cukup untuk `q=rahmat`; jika acceptance juga memerlukan search by code, implementer perlu menambahkan/menyesuaikan dengan behavior existing di endpoint lain setelah konfirmasi.
- Open questions: tidak ada pertanyaan material yang perlu ditahan sebelum implementasi; requirement dan pola lokal sudah cukup jelas.

## TDD/Test Plan

TDD required: Ya.

Alasan: perubahan menyentuh filter tenant/business unit yang security-sensitive dan berpotensi regression pada query list production.

Existing test patterns:

- `master/repository/salesman_repository_test.go` untuk helper scope condition.
- `master/controller/salesman_controller_test.go` untuk parser query `distributor_id`.
- `master/repository/sales_team_repository_test.go` sebagai referensi perilaku sentinel `0` yang benar.

### Red Step — first failing/regression test

Tambahkan/ubah test di `master/repository/salesman_repository_test.go`:

1. Ubah test lama `TestBuildSalesmanCustScopeCondition_IgnoresZeroWhenDistributorFilterHasValidValues` menjadi ekspektasi baru, misalnya `TestBuildSalesmanCustScopeCondition_WithPrincipalAndDistributorScope`.
   - Input: `[]int{0, 67, 67, 68}`, `parentCustId="C22001"`, `custId="C220010001"`.
   - Expected contains: `s.cust_id = 'C22001'`.
   - Expected contains: `mc.distributor_id IN (67,68)`.
   - Expected contains: ` OR `.
   - Expected starts/contains grouped condition `( ... )` agar precedence aman.
2. Tambahkan test principal-only:
   - Input: `[]int{0}`.
   - Expected contains `s.cust_id = 'C22001'`.
   - Expected not contains `mc.distributor_id IN`.
3. Pastikan distributor-only tetap:
   - Input: `[]int{10,20}`.
   - Expected contains distributor subquery dan tidak contains principal `s.cust_id = 'C22001'` jika tidak ada sentinel `0`.
4. Tambahkan guard negative/duplicate:
   - Input: `[]int{0, 120, -1, 120}`.
   - Expected includes principal and `IN (120)`, tidak include `-1`.

Jika ingin lebih kuat tanpa DB integration, tambahkan helper test untuk membangun `qWhere` terpisah hanya bila implementer mengekstrak builder query; jangan refactor besar bila tidak perlu.

### Green Step

Implementasi minimal:

1. Update `buildSalesmanCustScopeCondition` agar:
   - Mendeteksi `includePrincipalScope` saat ada `distributorID == 0`.
   - Mengabaikan negative id.
   - Deduplicate id.
   - Membuat `conditions`:
     - Principal: `s.cust_id = '<scopeParentCustID>'`.
     - Distributor: `s.cust_id IN (SELECT mc.cust_id FROM smc.m_customer mc WHERE mc.parent_cust_id = '<scopeParentCustID>' AND mc.distributor_id IN (...))`.
   - Return `( <condition1> OR <condition2> )` saat ada conditions.
   - Fallback ke `s.cust_id = '<custId>'` saat no distributor filter atau semua id invalid/non-positive tanpa principal sentinel.
2. Jalankan unit test repository dan controller terkait.
3. Jika test/manual menunjukkan duplicate, tambahkan `DISTINCT`/`COUNT(DISTINCT ...)` dengan scope minimal dan test bila memungkinkan.

### Refactor Step

- Hilangkan duplikasi logika dengan tetap menjaga readability helper.
- Pertimbangkan membuat helper kecil untuk dedupe distributor IDs hanya jika membuat kode lebih jelas.
- Jangan refactor query raw SQL luas ke `sqlx.In`/bind args dalam task ini kecuali blocker test/security muncul.

### Edge Cases

- `distributor_id` kosong/default: tetap `s.cust_id = custId`.
- `distributor_id=0`: principal-only.
- `distributor_id=120`: distributor-only.
- `distributor_id=0,120`: mixed principal + distributor.
- `distributor_id=0,120,120,-1`: dedupe dan abaikan negative.
- `parentCustId` kosong: fallback ke `custId`.
- `sales_team_id` tetap diterapkan setelah grouped scope.
- `q` kosong dan non-kosong.
- `is_active=1` dan `is_active=2` jika endpoint menggunakan filter tersebut.

### Commands

Jalankan dari module `master`:

```bash
rtk go test ./controller -run 'TestParseSalesmanDistributor'
rtk go test ./repository -run 'TestBuildSalesmanCustScopeCondition|TestFindAllByCustId_AppliesDistributorFilterToQuery'
rtk go test ./repository ./controller
rtk go test ./...
```

Jika `go test ./...` terlalu luas/terblokir dependency environment, minimal laporkan blocker dan hasil test targeted di atas.

## Implementation Steps

1. Buka `master/repository/salesman_repository_test.go`.
2. Ubah test yang saat ini mengharapkan `0` diabaikan ketika ada id valid menjadi test yang mengharapkan union principal + distributor.
3. Tambahkan test principal-only untuk `[]int{0}` dan guard dedupe/negative jika belum tercakup.
4. Jalankan Red test repository targeted dan pastikan gagal karena helper lama tidak memasukkan principal scope.
5. Update `buildSalesmanCustScopeCondition` di `master/repository/salesman_repository.go` dengan pola `buildSalesTeamCustScopeCondition`, menggunakan alias `s`.
6. Jalankan Green test targeted.
7. Review query string hasil helper untuk memastikan grouping `( ... OR ... )` benar.
8. Jalankan test controller parser untuk memastikan parsing `0` tetap aman.
9. Jalankan broader module test sesuai kemampuan environment.
10. Jika manual API tersedia, panggil endpoint lokal/staging dengan token valid untuk mixed `distributor_id=0,120&q=rahmat` dan verifikasi response mengandung `Rahmat / EMP0012`.
11. Dokumentasikan root cause final, file berubah, test, dan bukti response.

## Expected Files to Change

Kemungkinan file implementasi:

- `master/repository/salesman_repository.go`
- `master/repository/salesman_repository_test.go`

Opsional jika ditemukan blocker tambahan:

- `master/service/salesman_service.go` jika detail lookup mixed perlu memakai `row.CustId`.
- `master/controller/salesman_controller_test.go` hanya jika parser regression perlu ditambah; parser saat discovery sudah benar.

Tidak diharapkan berubah:

- File migrasi database.
- File FE.
- Kontrak API.
- `go.mod`/`go.sum`, kecuali test environment memerlukan tidy yang harus dijustifikasi terpisah.

## Agent/Tool Routing

- Implementasi: route ke `@fixer` / `opencode-fixer` dengan plan ini sebagai source of truth.
- Review arsitektur/security opsional: `@oracle` bila implementer menemukan ambiguity pada tenant scoping atau duplicate count.
- Security/privacy review opsional: `@security-privacy-reviewer` bila perubahan meluas ke auth/tenant isolation di luar helper.
- Quality gate: gunakan `@quality-gate` setelah implementasi dan test karena bugfix menyentuh tenant/business-unit filtering.
- Tidak perlu `@designer`, browser visual, document specialist, atau visual asset generator.

## Validation Commands

Pre-check dari repo root bila belum dilakukan:

```bash
rtk docker compose -f docker-compose.yml ps
```

Targeted tests dari `master`:

```bash
rtk go test ./repository -run 'TestBuildSalesmanCustScopeCondition|TestFindAllByCustId_AppliesDistributorFilterToQuery'
rtk go test ./controller -run 'TestParseSalesmanDistributor'
rtk go test ./repository ./controller
```

Full regression bila environment siap:

```bash
rtk go test ./...
```

Manual API setelah deploy/local token valid:

```bash
curl '<BASE_URL>/master/v1/salesman?page=1&sort=sales_name:asc&sales_team_id=79,76,73,70&distributor_id=0,120&q=rahmat&limit=9999' \
  -H 'Authorization: Bearer <TOKEN>' \
  -H 'Accept: application/json'
```

Expected manual evidence:

- Response data mengandung `sales_name` `Rahmat` dan/atau `emp_code` `EMP0012` sesuai field response aktual.
- Tidak ada duplicate `emp_id`.
- Paging total sesuai jumlah unique hasil.

## Evidence Requirements

Implementation handoff harus menghasilkan evidence berikut:

- Red test failure untuk helper scope sebelum fix, atau minimal catatan test lama yang diubah karena mengunci behavior bug.
- Green test output targeted repository/controller.
- Jika memungkinkan, output `rtk go test ./...` atau blocker detail.
- Manual API response snippet lokal/staging yang menunjukkan `Rahmat / EMP0012` untuk request Jira.
- Catatan apakah duplicate/total count diverifikasi.

Sumber yang digunakan dalam planning:

- Local project discovery: digunakan dan dirangkum di `.opencode/evidence/20260504-2058-sx-1915-salesman-business-unit/discovery.md`.
- Official docs/context7: tidak digunakan karena tidak diperlukan.
- GitHub/web: tidak digunakan karena tidak bergantung upstream/current external facts.
- Browser: tidak digunakan karena ini bugfix backend/API, bukan UI visual parity.

## Done Criteria

- Test helper scope salesman mencakup principal-only, distributor-only, mixed principal+distributor, dedupe, dan invalid/negative guard.
- Parser test tetap memastikan `distributor_id=0` dipertahankan.
- `buildSalesmanCustScopeCondition` mengembalikan grouped OR condition untuk mixed scope.
- Test targeted lulus.
- Manual/API verification menunjukkan kasus SX-1915 menemukan `Rahmat / EMP0012` jika data staging tersedia.
- Summary implementasi mencantumkan root cause final, file berubah, test dijalankan, dan evidence endpoint.
- `@quality-gate` pass atau pass with risks tanpa blocker sebelum final/commit.

## Final Planning Summary

Artifacts created:

- Primary plan: `.opencode/plans/20260504-2058-sx-1915-salesman-business-unit.md`
- Discovery evidence: `.opencode/evidence/20260504-2058-sx-1915-salesman-business-unit/discovery.md`

Key decisions:

- Source of bug berada pada `buildSalesmanCustScopeCondition` yang membuang sentinel `0` saat ada distributor non-zero.
- Reuse pola `buildSalesTeamCustScopeCondition` untuk union principal + distributor.
- Tidak mengubah kontrak FE `distributor_id=0,120`.
- TDD wajib karena perubahan menyentuh tenant/business-unit filtering.

Assumptions:

- Principal sentinel `0` harus dipetakan ke `parentCustId`/principal cust scope.
- Existing search by `sales_name` cukup untuk `q=rahmat`, kecuali implementer menemukan behavior Master Salesman juga mencari by code.

Questions:

- Tidak ada pertanyaan material yang ditanyakan; requirement dari Jira dan pola lokal cukup jelas. Pertanyaan minor dicatat sebagai asumsi.

Readiness:

- Siap untuk implementasi oleh `@fixer` berdasarkan plan ini.

Cleanup performed:

- Tidak ada draft artifact yang dibuat.
- Evidence discovery tetap disimpan karena operasional berguna untuk implementer memahami file dan pola yang ditemukan.
