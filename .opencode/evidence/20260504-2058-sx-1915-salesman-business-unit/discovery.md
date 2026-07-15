# Discovery Evidence — SX-1915 Salesman Business Unit Filter

Task ID: `20260504-2058-sx-1915-salesman-business-unit`
Tanggal: 2026-05-04 Asia/Jakarta

## Files Inspected

- `master/controller/salesman_controller.go`
  - Route `GET /v1/salesman` berada di `SalesmanController.List`.
  - Query `distributor_id` diparse via `parseSalesmanDistributorIDQuery(c.Context().QueryArgs())` lalu masuk ke `dataFilter.DistributorID`.
  - Controller memilih `custId` dari header `cust_id` jika ada; jika tidak, dari JWT locals `cust_id` dan `parent_cust_id`.
- `master/controller/query_filter_parser.go`
  - `parseIntSliceQueryAllowZero` mempertahankan nilai `0` khusus untuk `distributor_id`.
  - Mendukung repeated query, comma-separated, dan `distributor_id[]`.
- `master/controller/salesman_controller_test.go`
  - Sudah ada test bahwa `0` dipertahankan sebagai principal scope saat parsing.
- `master/repository/salesman_repository.go`
  - Scope filter ada di helper `buildSalesmanCustScopeCondition(distributorIDs []int, parentCustId, custId string)`.
  - `FindAllByCustId` dan `FindAllByCustIdLookup` sama-sama memakai helper tersebut.
  - Jika `distributorIDs` berisi nilai valid non-zero, helper saat ini membuang `0` dan hanya menghasilkan `s.cust_id IN (SELECT mc.cust_id ... mc.distributor_id IN (...))`.
  - Ini berarti request `distributor_id=0,120` tidak memasukkan scope principal/current parent customer.
- `master/repository/salesman_repository_test.go`
  - Ada test `TestBuildSalesmanCustScopeCondition_IgnoresZeroWhenDistributorFilterHasValidValues` yang mengunci perilaku bermasalah: principal scope diabaikan saat filter berisi `0` dan id distributor valid.
- `master/repository/sales_team_repository.go`
  - Ada pola reusable `buildSalesTeamCustScopeCondition` yang sudah benar untuk sentinel `0`: `includePrincipalScope` menghasilkan `a.cust_id = parentCustId`, id non-zero menghasilkan subquery distributor, lalu digabung dengan `OR` dalam tanda kurung.
- `master/repository/sales_team_repository_test.go`
  - Ada test pembanding untuk principal-only dan mixed principal+distributor.

## Project Patterns Found

- Repo adalah monorepo Go multi-module; target perbaikan ada di module `master`.
- Route publik staging `/master/v1/salesman` kemungkinan diproksi ke service master route internal `/v1/salesman`.
- Filter tenant/cust scope di repository masih dibangun sebagai SQL string helper, bukan query builder parameterized.
- Pola reusable untuk mixed principal+distributor sudah ada di `sales_team_repository.go` dan sebaiknya diadaptasi ke alias salesman `s`.
- Test yang paling ringan saat ini adalah unit test helper repository dan parser controller; belum terlihat integration test DB untuk endpoint ini.

## Reuse Candidates

- Reuse struktur `buildSalesTeamCustScopeCondition` untuk memperbaiki `buildSalesmanCustScopeCondition`.
- Reuse test style di `sales_team_repository_test.go` untuk menambahkan test principal-only dan mixed pada salesman.
- Reuse parser `parseSalesmanDistributorIDQuery` karena sudah mempertahankan `0`; tidak perlu ubah kontrak FE.

## Commands / Docs Checked

- `rtk docker compose -f docker-compose.yml ps`
  - `master`, `system`, dan `redis` aktif; service lain tidak aktif dari output saat discovery.
- Local search/read via `Glob`, `Grep`, dan `Read` pada file controller/service/repository/test module `master`.
- Official docs/context7 tidak diperlukan karena perbaikan memakai pola SQL/Go lokal yang sudah ada.
- GitHub/web/browser tidak diperlukan untuk rencana ini karena masalah terlokalisasi di codebase lokal dan evidence Jira sudah cukup untuk acceptance behavior.

## Constraints

- Harus menjaga kontrak FE `distributor_id=0,120`.
- `0` adalah sentinel principal/current customer, bukan id distributor literal.
- Query harus tetap membatasi ke `parent_cust_id` untuk distributor child customer agar tidak bocor ke tenant lain.
- Grouping kondisi principal/distributor harus dalam tanda kurung agar `sales_team_id`, `q`, `is_active`, `is_del`, pagination, dan sort berlaku ke semua branch.
- Ada instruksi repo yang meminta `rtk` prefix untuk shell, tetapi instruksi global OpenCode meminta tidak memakai `rtk`. Dalam sesi ini perintah mandatory repo dijalankan dengan `rtk` untuk menghormati instruksi repo lokal yang lebih spesifik.

## Risks

- Helper saat ini melakukan SQL string concatenation; rencana fix tidak memperburuk risiko ini, tetapi tetap bukan refactor total ke parameterized query agar scope SX-1915 tetap kecil.
- `FindAllByCustId` memanggil `FindDetailById(row.EmpId, custId)` di service dengan `custId` login, sehingga jika list mixed mengembalikan salesman distributor child, detail lookup mungkin tidak sesuai tenant. Ini bukan akar issue dropdown jika detail kosong bisa ditoleransi, tetapi perlu diperhatikan saat implementer menjalankan test/manual response.
- Ada `fmt.Println("dataFilter.IsActive:", ...)` di repository yang melanggar style repo, tetapi bukan bagian langsung dari SX-1915; jangan perluas scope kecuali fixer memutuskan cleanup kecil aman.
- Count memakai `COUNT(*)`; jika join menimbulkan duplicate, acceptance meminta unique ids/total benar. Perlu verifikasi apakah join `mst.m_salesman_canvas` bisa menghasilkan lebih dari satu row per salesman. Jika ya, pertimbangkan `COUNT(DISTINCT s.emp_id)` dan `SELECT DISTINCT` atau `DISTINCT ON` secara hati-hati.

## Research Gate Decision

- Local project discovery: required dan sudah dilakukan.
- Official docs/context7: tidak diperlukan; tidak ada API eksternal/version-sensitive yang menentukan solusi.
- GitHub: tidak diperlukan; tidak bergantung upstream.
- Brave/web search: tidak diperlukan; evidence Jira dan local code cukup.
- Browser/screenshot capture: tidak diperlukan untuk rencana BE; manual API evidence cukup.
