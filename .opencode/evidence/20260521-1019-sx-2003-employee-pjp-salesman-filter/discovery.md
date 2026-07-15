# Discovery — SX-2003 Filter Employee role "salesman" pada `GET /master/v1/employee-pjp`

## Files diperiksa

- `master/controller/employee_controller.go` — handler `ListPJP` di route `app.Group("/v1/employee-pjp", middleware.JWTProtected())` line 52, set `dataFilter.CustId` dan `dataFilter.ParentCustId` dari JWT lewat `c.Locals` (line 152-153). Mengembalikan `data`, `total`, `lastPage`.
- `master/entity/employee.go` line 216-234 — `EmployeePJPQueryFilter` (CustId, ParentCustId, Page, Limit, Query, Sort, IsActive []int, FilterCustId *string, DistributorId *int) dan `EmployeePJPResponse` (EmpId, EmpCode, EmpName).
- `master/model/employee.go` line 120-125 — `EmployeePJP` minimal (`emp_id`, `emp_code`, `emp_name`).
- `master/service/employee_service.go` line 788-803 — `ListPJP` adapter ke repository, no business logic.
- `master/repository/employee_repository.go` line 999-1091 — `FindAllForPJP`. Membangun query string concatenation manual; tidak join `sys.m_user/sys.user_roles/sys.m_role`; tidak filter `mr.role_name = 'salesman'`; tidak filter `me.is_active = true` secara default (hanya kalau `IsActive` diisi); cabang scope:
  - `DistributorId > 0` → `me.cust_id = ( SELECT mc.cust_id FROM smc.m_customer mc WHERE mc.distributor_id = %d LIMIT 1 )`.
  - `FilterCustId != ""` → `me.cust_id = '<FilterCustId>'`.
  - default → `me.cust_id = '<JWT cust_id>'`.
- `master/repository/salesman_repository.go` line 85-134 — pola `buildSalesmanCustScopeCondition(distributorIDs, parentCustId, custId)` dengan principal `s.cust_id = '<parent>'` dan distributor map lewat `smc.m_customer.parent_cust_id` + `mc.distributor_id IN (...)`. Reuse-able.
- `master/repository/business_unit_repository.go` line 70-140 — pola `sqlx.In` + `Rebind` + parameterized args pada master service Fiber.
- `master/repository/outlet_repository.go` line 4362-4408 — helper `FindCustIdsByDistributorIds`/`FindCustIdsByParentCustId`.
- `master/repository/salesman_repository_test.go` & `controller/query_filter_parser_test.go` — pola test repository SQL string + test helper builder.
- `mobile/repository/m_user.go` line 60-72 — confirm tabel role: `sys.user_roles`, `sys.m_role`. Schema-prefix `sys.` ada.
- `system/repository/m_menu_repository.go` line 50, 66, 108 — penggunaan `sys.user_roles.role_id`, `sys.user_roles.user_id`, `sys.user_roles.cust_id`. Mendukung asumsi join key kolom tersebut.
- `docs/Monitoring Activity - BE.md` — point Employee menyebut filter cust_id by login type (distributor/principal/principal+distributor) dan `mr.role_name = 'salesman'`. Skenario sesuai SX-2003.
- `.opencode/docs/ARCHITECTURE.md` — Controller → Service → Repository; tenant rules `cust_id`/`parent_cust_id`; schema prefix penting.
- `.opencode/docs/QUALITY.md` — validate per service: `cd master && rtk go mod download && rtk go mod tidy && rtk go test ./...`.
- `master/go.mod` — Fiber v2; pakai `sqlx`, `lib/pq`. Tidak ada framework testing khusus. `DATA-DOG/go-sqlmock` tersedia.

## Pola proyek yang ditemukan

- Repository master sering memakai concat string (vulnerable) untuk dynamic clauses, tetapi pola lebih baru (`business_unit_repository`, `outlet_repository`, `m_distributor_repository`) sudah pakai `sqlx.In` + `Rebind` + args. Plan ini sebaiknya tetap aman:
  - Filter cust_id list (string) wajib lewat `sqlx.In` + `Rebind` agar tidak rentan SQL injection dan konsisten dengan pola baru.
  - `role_name` hardcoded literal `'salesman'`, atau lebih aman `LOWER(mr.role_name) = 'salesman'`.
- Tenant scoping pakai `parent_cust_id` (principal scope) atau mapping via `smc.m_customer`.
- JWT context: controller pasti set `cust_id` dan `parent_cust_id`; angka panjang `cust_id` 6 digit utk principal, > 6 digit utk distributor (sesuai dokumen + token sample).

## Reuse candidates

- `buildSalesmanCustScopeCondition` (salesman_repository) sebagai referensi pattern, tapi tidak diperluas ke skema role join.
- `FindCustIdsByDistributorIds`/`FindCustIdsByParentCustId` jika butuh resolusi cust_id list eksplisit.
- Pattern test `master/repository/salesman_repository_test.go` untuk membuat unit test fungsi builder query baru.
- Pattern parameterized + `sqlx.In` `business_unit_repository.go`.

## Commands & docs dicek

- `rtk docker compose -f docker-compose.yml ps` — output kosong, hanya warning version. Tidak menjalankan service. Tidak diperlukan untuk perubahan code-level pure.
- Doc resmi internal: `docs/Monitoring Activity - BE.md` (point Employee).
- `.opencode/docs/AGENT_ROUTING.md`, `ARCHITECTURE.md`, `QUALITY.md` — patuh layering & validation.

## Constraints

- Layer Controller → Service → Repository.
- Tenant aware: `cust_id`/`parent_cust_id` dari JWT; jangan bocor lintas principal.
- Schema prefix: `mst.m_employee`, `smc.m_customer`, `sys.m_user`, `sys.user_roles`, `sys.m_role`. Confirmed by usages dan dokumen.
- Endpoint masih dipakai FE; default param dari FE: `is_active[]=1`, `limit=9999`, `sort=area_id:asc`. `sort` tidak relevan ke kolom yang diselect, tapi backend tetap menerima string sort generik.
- Tidak boleh mengubah kontrak response. Plan ini hanya menambah filter di query.

## Risks

- Cross-tenant role join: `sys.user_roles`/`sys.m_role` punya `cust_id`. Harus di-AND dengan `cust_id` employee (atau `parent_cust_id`) supaya role lookup sesuai tenant employee.
- Jika `mr.role_name` ditulis dengan casing campur di DB, gunakan `LOWER(mr.role_name) = 'salesman'`.
- Jika satu employee punya multi-role + duplikasi join, query bisa balikin baris ganda. Gunakan `SELECT DISTINCT` atau `EXISTS (...)` subquery untuk filter role agar hasil tetap unik.
- `FindAllForPJP` saat ini concat literal `cust_id`/`distributor_id`. Memperkenalkan list `IN (...)` literal akan tetap aman untuk tipe int, tapi untuk string `cust_id` sebaiknya gunakan `sqlx.In` + `Rebind` untuk menghindari injection lewat FE-controlled params (`cust_id` query). Konsistensi pola baru.
- Endpoint lain yang reuse `FindAllForPJP`: hanya `EmployeeService.ListPJP` (grep `FindAllForPJP` hanya 2 hits di repository + service). Jadi perubahan terbatas hanya untuk endpoint ini.

## Test footprint

- Tidak ada `employee_controller_test.go` atau `employee_repository_test.go` saat ini.
- Pola test builder query string ada (salesman/business_unit). Buat unit test builder string baru untuk:
  - `buildEmployeePJPQuery` (atau builder ekstrak) memastikan join + WHERE role_name + tenant scope.
- `go-sqlmock` tersedia bila ingin test eksekusi.

## Decisions/Assumptions

- TDD via builder unit test (string contains assertions) konsisten dengan repo.
- Tetap pakai `me.is_active = true AND me.is_del = false` sesuai SQL referensi tiket (saat ini `is_del=false` ada, `is_active` belum dipaksa default → ikuti tiket: SQL referensi gunakan `me.is_active = true`. Aman karena FE selalu kirim `is_active[]=1`; default ke `is_active=true` tidak break compat untuk FE existing, tapi ubah default behavior. Jadikan keputusan eksplisit: pertahankan default existing — `is_active` hanya dipaksa via param, tidak di-hardcode default — supaya non-breaking. Filter role tetap diterapkan di semua skenario).
- `LOWER(mr.role_name) = 'salesman'` untuk safety case-insensitive.
- Role join harus tenant-aware: `mu.cust_id = me.cust_id`, `ur.cust_id = mu.cust_id`, `mr.cust_id = ur.cust_id`. Sesuai dokumen tiket.
- Hindari row duplikat: gunakan `EXISTS (SELECT 1 FROM sys.m_user mu JOIN sys.user_roles ur ... JOIN sys.m_role mr ... WHERE mu.emp_id = me.emp_id AND mu.cust_id = me.cust_id AND LOWER(mr.role_name) = 'salesman')` agar tidak menimbulkan cartesian/duplikasi karena multi-role.
