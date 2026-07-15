# Plan — SX-1578 Business Unit Dropdown Filter

## Goal

Memperbaiki backend `master` agar dropdown Target Survey hanya mengembalikan **Salesman**, **Sales Team**, dan **Outlet** sesuai Business Unit/distributor yang dipilih, dengan dukungan CSV `distributor_id` dan filter outlet approved/valid lewat parameter eksplisit.

## Non-goals

- Tidak mengubah kontrak response DTO.
- Tidak mengubah FE, staging credential, token, atau cookie.
- Tidak meng-hardcode distributor id dari Jira.
- Tidak menjadikan status outlet survey sebagai default global jika parameter status tidak dikirim.
- Tidak melakukan implementasi/source edit dalam mode artifact planner ini.

## Scope

Modul target: `master/`.

Endpoint/service route lokal yang relevan:

- `GET /v1/salesman` (dipublish sebagai `/master/v1/salesman`)
- `GET /v1/sales-teams` (dipublish sebagai `/master/v1/sales-teams`)
- `GET /v1/outlets` (dipublish sebagai `/master/v1/outlets`)

File area yang kemungkinan berubah:

- Controller parser dan tests.
- Entity filter untuk outlet status multi-value.
- Repository query builder salesmen/sales team/outlet dan tests.
- Service outlet jika perlu handling `0`/principal scope pada distributor resolution.

## Requirements

- `distributor_id` menerima single, comma-separated, repeated param, dan `distributor_id[]`.
- `distributor_id=0,67,68` harus diproses sebagai gabungan principal + distributor sesuai domain existing, bukan string tunggal.
- `GET /v1/salesman` memfilter data berdasarkan Business Unit/distributor tanpa merusak `q`, `sales_team_id`, `is_active`, pagination, dan sort.
- `GET /v1/sales-teams` memfilter berdasarkan Business Unit/distributor dan mengaktifkan logic existing untuk `0` principal scope.
- `GET /v1/outlets` memfilter berdasarkan Business Unit/distributor, `verification_status`, dan multi-value outlet status `1,5,6,7` jika parameter dikirim.
- `ot_type_id`, `ot_grp_id`, `ot_class_id`, `outlet_id`, `q`, `sort`, `page`, `limit`, dan `is_active` tetap berfungsi dan digabung dengan `AND`.
- Jika `distributor_id` kosong/tidak dikirim, behavior existing tetap berlaku.
- Filter request tidak boleh memperluas akses di luar tenant/user scope existing.

## Acceptance Criteria

- `/master/v1/salesman?...&distributor_id=102` hanya mengembalikan salesman di distributor/business unit tersebut.
- `/master/v1/salesman?...&distributor_id=0,67,68` menangani principal + distributor sesuai keputusan implementasi yang aman.
- `/master/v1/sales-teams?...&distributor_id=0,67,68` mengembalikan sales team dalam principal/distributor scope yang sesuai.
- `/master/v1/outlets?...&distributor_id=102&verification_status=1&outlet_status=1,5,6,7` hanya mengembalikan outlet distributor terkait dengan `verification_status=1` dan `outlet_status IN (1,5,6,7)`.
- Optional filters tetap bekerja bersamaan dengan filter distributor/status.
- Parser CSV menolak nilai invalid seperti `abc` dengan `400 Bad Request`.
- Pagination total/count mengikuti hasil filtered query.
- Tidak ada data business unit lain bocor di dropdown Target Survey.

## Existing Patterns/Reuse

Discovery lokal menemukan beberapa reuse penting:

- `parseCSVIntValues` dan `parseIntSliceQuery` di `master/controller/query_filter_parser.go` sudah mendukung repeated dan comma-separated integer query.
- `SalesmanQueryFilter.DistributorID []int` sudah ada dan repository salesman sudah memanggil `buildSalesmanCustScopeCondition`.
- `GeneralQueryFilter.DistributorIDs []int` sudah ada untuk sales team.
- `buildSalesTeamCustScopeCondition` sudah punya logic `0` sebagai principal scope dan distributor id positif sebagai `mc.distributor_id IN (...)`.
- `OutletQueryFilter.DistributorID []int`, `VerificationStatus []int`, `OtClassID`, `OtTypeID`, `OtGrpID`, dan `ResolvedCustIdsForDistributor` sudah tersedia.
- `appendIntInFilter` sudah bisa membuat `IN (...)` untuk outlet filters.
- Tests existing sudah ada untuk parser dan query helper sehingga TDD bisa langsung menambah regression tests.

Tidak ditemukan utilitas KiloCode/project lain yang lebih spesifik dari parser/query helper existing untuk kasus ini; rencana mengutamakan extend/reuse helper yang sudah ada.

## Constraints

- Ikuti arsitektur Controller → Service → Repository → DB.
- Repository tetap bertanggung jawab query, service tetap menjaga orchestration/scope, controller parse request.
- Tetap pakai tenant scope `cust_id`/`parent_cust_id` dari JWT/header sesuai pattern existing.
- Jangan menambah secret, token, atau credential.
- Karena repo `master` banyak memakai SQL string concatenation, perubahan harus sekecil mungkin dan tidak memperkenalkan input mentah baru.
- Perintah validasi dijalankan dari `master/` dan mengikuti project-local instruction dengan prefix `rtk`.

## Risks

- **Ambiguitas `0` principal**: sales team sudah mendukung `0`, tetapi parser saat ini membuangnya dan salesman/outlet belum konsisten. Implementasi harus menyamakan behavior tanpa membuka semua customer tenant sembarangan.
- **Global parser impact**: mengubah `parseIntSliceQuery` agar tidak selalu membuang `0` bisa memengaruhi endpoint lain. Lebih aman tambah opsi helper atau wrapper khusus `distributor_id`.
- **Outlet status single-value**: `OutletQueryFilter.OutletStatus *int` tidak cukup untuk `1,5,6,7`; perlu field baru atau parse manual agar backward compatible.
- **Pagination multi-cust outlet**: path `listOutletReadMultiCustForCitus` harus tetap menggunakan filter status/verification/distributor yang sama.
- **SQL injection legacy**: existing `sales_team_id`, `sort`, dan `q` masih raw string concat. Jangan memperburuk risiko; jika disentuh, minimal validasi CSV numeric untuk `sales_team_id` dapat dipertimbangkan sebagai hardening terpisah.

## Decisions/Assumptions

- **Decision**: Backend akan mendukung parameter eksplisit; flow survey harus mengirim `verification_status=1&outlet_status=1,5,6,7` agar tidak mengubah default global `/v1/outlets` untuk modul lain.
- **Decision**: `0` harus dipertahankan untuk `distributor_id` pada endpoint yang membutuhkan principal scope, bukan dibuang oleh parser umum.
- **Decision**: Filter distributor direpresentasikan sebagai scope customer (`cust_id`) melalui `smc.m_customer.distributor_id`/`parent_cust_id`, mengikuti pattern existing.
- **Assumption**: Business Unit yang dipilih FE mengirim daftar distributor id yang sudah dalam scope user; BE tetap harus membatasi via `parent_cust_id`/`cust_id` existing sehingga id luar tenant tidak menghasilkan data.
- **Assumption**: Untuk outlet, status valid survey dikirim oleh FE sebagai parameter eksplisit; BE tidak perlu parameter baru `source=survey_target` kecuali product meminta default khusus.
- **Open Question**: Jika `distributor_id=0` pada salesman/outlet, apakah harus berarti semua principal + child distributor dalam `parent_cust_id`, atau hanya principal cust row? Implementasi direkomendasikan mengikuti sales team: `0` = principal scope aman dalam `parent_cust_id`.

Question gate: tidak perlu menunggu jawaban sebelum implementasi karena requirement Jira sudah menyatakan contoh `0,67,68`, dan repository sales team sudah memberi pola domain untuk `0`.

## TDD/Test Plan

TDD wajib karena ini defect pada API/query behavior.

### Existing test patterns

- `master/controller/query_filter_parser_test.go`
- `master/controller/salesman_controller_test.go`
- `master/controller/outlet_controller_test.go`
- `master/repository/salesman_repository_test.go`
- `master/repository/sales_team_repository_test.go`
- `master/repository/outlet_repository_test.go`

### Red step: first failing/regression tests

Tambahkan tests yang gagal sebelum implementasi:

1. Parser distributor yang mempertahankan `0` untuk `distributor_id=0,67,68` menghasilkan `[]int{0,67,68}` pada endpoint sales team/salesman/outlet yang membutuhkan principal scope.
2. Salesman cust scope dengan `[]int{0,67,68}` menghasilkan kombinasi principal scope + distributor scope yang dibatasi `parent_cust_id`.
3. Outlet query helper mendukung `outlet_status=1,5,6,7` menjadi `o.outlet_status IN (1,5,6,7)`.
4. Outlet filter menggabungkan `verification_status=1`, `outlet_status=1,5,6,7`, `ot_type_id`, `ot_grp_id`, `ot_class_id`, dan `distributor_id` dengan `AND`.
5. Empty `distributor_id` tetap tidak menambah filter distributor baru.
6. Invalid CSV seperti `distributor_id=102,abc` atau `outlet_status=1,abc` menghasilkan error.

### Green step

- Perbarui parser/wrapper agar `0` dapat dipertahankan khusus distributor principal scope.
- Tambahkan field multi-value outlet status, misalnya `OutletStatusIDs []int`, tanpa menghapus `OutletStatus *int` lama.
- Perbarui repository helper untuk salesman/outlet agar mendukung `0` principal scope secara aman.
- Terapkan `appendIntInFilter` atau helper baru pada outlet list utama untuk `OutletStatusIDs`.

### Refactor step

- Hilangkan duplikasi parsing distributor antar controller jika aman.
- Pastikan nama helper jelas, misalnya `parseDistributorIDQueryPreservePrincipal` atau opsi `allowZero`.
- Pertahankan backward compatibility dengan `OutletStatus *int` single value.

### Edge cases

- `distributor_id=null`, `undefined`, kosong.
- `distributor_id=0` saja.
- `distributor_id=0,67,67,68` deduplicate.
- Distributor id luar `parent_cust_id` menghasilkan empty list.
- `outlet_status=0` atau kosong tidak memfilter status.
- Repeated query params: `outlet_status=1,5&outlet_status=6,7`.

### Commands

Dari `master/`:

```bash
rtk go test ./controller -run 'TestParse|Test.*Distributor|Test.*Outlet'
rtk go test ./repository -run 'TestBuildSalesman|TestBuildSalesTeam|TestAppendInt|TestBuildOutlet'
rtk go test ./...
```

## Implementation Steps

1. **Parser strategy**
   - Jangan ubah behavior global secara sembrono.
   - Tambahkan helper baru atau parameter opsi untuk parse integer query yang bisa mempertahankan `0` khusus `distributor_id` principal scope.
   - Update `salesman_controller.go`, `sales_team_controller.go`, dan `outlet_controller.go` agar menggunakan helper yang sesuai.

2. **Sales Team**
   - Pastikan `GET /v1/sales-teams` meneruskan `DistributorIDs` termasuk `0`.
   - Reuse `buildSalesTeamCustScopeCondition`; tests existing untuk `0` harus benar-benar terpakai dari controller parser.

3. **Salesman**
   - Extend `buildSalesmanCustScopeCondition` agar mendukung `0` mirip sales team:
     - `0` → principal scope dalam `parent_cust_id` yang aman.
     - positive ids → `s.cust_id IN (SELECT mc.cust_id ... mc.distributor_id IN (...))`.
     - kombinasi → OR dalam scope parent.
   - Pertahankan fallback `s.cust_id = '<custId>'` saat no distributor filter.
   - Pastikan list dan lookup sama-sama memakai condition ini.

4. **Outlet status multi-value**
   - Tambahkan field baru di `OutletQueryFilter`, misalnya `OutletStatusIDs []int` dengan `query:"-"` atau parse manual dari query args.
   - Di `OutletController.List`, parse `outlet_status` dan `outlet_status[]` sebagai CSV/repeated values.
   - Jika nilai hanya satu, tetap bisa mengisi field baru; repository dapat memprioritaskan `OutletStatusIDs` lalu fallback `OutletStatus *int`.

5. **Outlet distributor principal scope**
   - Tentukan helper resolve distributor ids agar `0` tidak dikirim ke `FindCustIdsByDistributorIds` sebagai distributor real.
   - Rekomendasi: jika `0` ada, include `parentCustId` dan/atau seluruh child cust dalam `parentCustId` sesuai behavior product yang disepakati; jika mengikuti sales team, `0` mewakili principal scope. Untuk dropdown survey dengan `0,67,68`, hasil final harus union aman dari principal scope + distributor ids 67/68.
   - Jangan mengembalikan data di luar `parentCustId`.

6. **Outlet repository filter**
   - Terapkan `o.verification_status IN (...)` sudah ada; pastikan tests menutupinya.
   - Terapkan `o.outlet_status IN (...)` untuk multi-status.
   - Pastikan `ot_class_id`, `ot_type_id`, `ot_grp_id`, `q`, `is_active`, sort, page, limit tetap digabung.

7. **Hardening ringan**
   - Jika sempat, validasi `sales_team_id` CSV numeric sebelum repository raw IN, atau rencanakan sebagai follow-up jika terlalu besar.
   - Jangan mengubah response shape.

8. **Run tests dan API checks**
   - Jalankan unit tests terarah, lalu `rtk go test ./...`.
   - Dengan auth lokal/staging valid, ambil evidence untuk minimal 2 business unit/distributor berbeda.

## Expected Files to Change

Kemungkinan file implementasi:

- `master/controller/query_filter_parser.go`
- `master/controller/salesman_controller.go`
- `master/controller/sales_team_controller.go`
- `master/controller/outlet_controller.go`
- `master/entity/outlet.go`
- `master/repository/salesman_repository.go`
- `master/repository/outlet_repository.go`
- Mungkin `master/service/outlet_service.go` untuk resolution `0` principal scope.

Kemungkinan file test:

- `master/controller/query_filter_parser_test.go`
- `master/controller/salesman_controller_test.go`
- `master/controller/outlet_controller_test.go`
- `master/repository/salesman_repository_test.go`
- `master/repository/sales_team_repository_test.go`
- `master/repository/outlet_repository_test.go`

## Agent/Tool Routing

- Implementasi: gunakan `@fixer` atau build agent dengan TDD Red → Green → Refactor.
- Codebase discovery tambahan: `@explorer` jika perlu mencari pemakaian endpoint lain.
- Review arsitektur/security: `@oracle` jika ada keraguan soal makna `0` principal scope atau tenant authorization.
- Tidak perlu `@designer`, browser, visual asset, atau external docs untuk backend defect ini.

## Validation Commands

Dari repo root, sesuai instruksi project sebelum code work:

```bash
rtk docker compose -f docker-compose.yml ps
```

Jika service perlu dinyalakan:

```bash
rtk docker compose -f docker-compose.yml up -d
```

Dari `master/`:

```bash
rtk go test ./controller -run 'TestParse|Test.*Distributor|Test.*Outlet'
rtk go test ./repository -run 'TestBuildSalesman|TestBuildSalesTeam|TestAppendInt|TestBuildOutlet'
rtk go test ./...
```

API smoke dengan auth valid:

```bash
# Login dengan credential QA yang diberikan user secara out-of-band.
# Jangan simpan email/password/token plaintext ke repo, artifact, shell history permanen, atau final evidence.
# Simpan token hanya sebagai env var sementara, misalnya:
export BASE_URL="<STAGING_OR_LOCAL_BASE_URL>"
export TOKEN="<TOKEN_FROM_LOGIN>"

rtk curl "<BASE_URL>/master/v1/salesman?page=1&limit=70&sort=sales_name:asc&q=&distributor_id=102" \
  -H "Authorization: Bearer ${TOKEN}"

rtk curl "<BASE_URL>/master/v1/outlets?page=1&limit=70&sort=outlet_id:desc&is_active=1&verification_status=1&outlet_status=1,5,6,7&distributor_id=102&q=" \
  -H "Authorization: Bearer ${TOKEN}"

rtk curl "<BASE_URL>/master/v1/sales-teams?page=1&limit=70&distributor_id=0,67,68" \
  -H "Authorization: Bearer ${TOKEN}"
```

Catatan: gunakan token lokal/staging yang valid dari login QA yang diberikan user, bukan token dari Jira/evidence.

Validasi database boleh dilakukan ke database QA/dev yang diberikan user secara out-of-band. Jangan simpan connection string dengan password plaintext di repo/artifact/evidence. Gunakan env var sementara atau `.pgpass` lokal yang tidak dicommit:

```bash
export PGHOST="<DB_HOST_FROM_USER>"
export PGPORT="<DB_PORT_FROM_USER>"
export PGUSER="<DB_USER_FROM_USER>"
export PGPASSWORD="<DB_PASSWORD_FROM_USER>"
export PGDATABASE="<DB_NAME_FROM_USER>"
export PGSSLMODE="disable"

# Contoh validasi bentuk query, sesuaikan distributor ids dan parent/cust scope dari response/login:
rtk psql -c "SELECT distributor_id, cust_id, parent_cust_id FROM smc.m_customer WHERE distributor_id IN (102,67,68) ORDER BY distributor_id;"
rtk psql -c "SELECT o.cust_id, COUNT(*) FROM mst.m_outlet o WHERE o.verification_status = 1 AND o.outlet_status IN (1,5,6,7) GROUP BY o.cust_id ORDER BY o.cust_id;"
```

## Evidence Requirements

Saat implementasi selesai, simpan/rangkum evidence berikut di final developer summary atau artifact implementation evidence:

- Output unit tests dan `rtk go test ./...`.
- Request/response sanitized untuk minimal 2 distributor/business unit berbeda.
- Bukti semua outlet response memenuhi `verification_status=1` dan `outlet_status IN (1,5,6,7)` saat parameter dikirim.
- Bukti salesman/sales team/outlet berubah ketika Business Unit A diganti ke Business Unit B.
- Catatan jika `0` principal scope menghasilkan behavior khusus.

Research gate:

- Local project discovery: dilakukan dan cukup untuk plan.
- Official docs/context7: tidak diperlukan; perubahan bergantung pada Go/Fiber/sqlx pattern lokal yang sudah ada.
- GitHub: tidak diperlukan; tidak bergantung upstream repo.
- Brave/web search: tidak diperlukan; fakta eksternal tidak dibutuhkan.
- Browser/screenshot: tidak diperlukan untuk backend API plan; UI QA cukup dengan API evidence dan FE manual verification.

## Done Criteria

- Tests parser dan repository helper lulus.
- `rtk go test ./...` di `master/` lulus atau ada blocker terdokumentasi.
- API mendukung `distributor_id` CSV/repeated/array style.
- `0` principal scope diproses sesuai keputusan dan tetap tenant-safe.
- Outlet multi-status dan verification status berjalan bersamaan dengan filter distributor dan optional filters.
- Tidak ada response contract change.
- QA summary berisi endpoint berubah, query param didukung, contoh request, dan caveat.

## Final Planning Summary

- Artifact utama dibuat: `.opencode/plans/20260501-1200-sx-1578-bu-dropdown-filter.md`.
- Discovery dibuat sementara di `.opencode/evidence/20260501-1200-sx-1578-bu-dropdown-filter/discovery.md`, lalu dikonsolidasikan ke plan ini dan dibersihkan agar tidak menjadi konteks stale.
- Keputusan kunci: dukung filter eksplisit, jangan hardcode status survey sebagai default global, pertahankan `0` untuk principal scope khusus distributor, dan reuse helper/query pattern existing.
- Asumsi kunci: FE akan mengirim `verification_status=1&outlet_status=1,5,6,7` untuk flow Target Survey; backend hanya memastikan parameter tersebut benar-benar bekerja.
- Open question tidak memblokir: makna detail `0` untuk salesman/outlet direkomendasikan mengikuti sales team (`0` sebagai principal scope dalam `parent_cust_id`) dan harus diverifikasi saat implementation review.
- Readiness: siap untuk implementasi TDD oleh fixer/build agent.
- Cleanup: evidence discovery sudah diringkas di plan dan akan dihapus setelah plan ditulis.
