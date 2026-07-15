Discovery SX-2124 — Business Unit cust_name

Task id: `20260602-1111-sx-2124-business-unit-cust-name`
Waktu: `2026-06-02 11:11 Asia/Jakarta`

File diperiksa:

- `master/controller/business_unit_controller.go`
- `master/service/business_unit_service.go`
- `master/repository/business_unit_repository.go`
- `master/entity/business_unit.go`
- `master/model/business_unit.go`
- `master/service/business_unit_service_test.go`
- `master/repository/business_unit_repository_test.go`
- `master/repository/employee_scope_repository.go`
- `master/repository/employee_repository.go`
- `.opencode/docs/ARCHITECTURE.md`

Pola ditemukan:

- Endpoint `GET /master/v1/business-unit` diroute oleh `BusinessUnitController.Route` ke `/v1/business-unit` di service `master`.
- Controller mengambil konteks JWT: `cust_id`, `parent_cust_id`, `user_name`, `employee_id`, `distributor_id`.
- Service `GetBusinessUnit` selalu memanggil `FindUserByUsername(dataFilter.UserName)`.
- Principal path: `DistributorId == nil || *DistributorId == 0`.
- Principal path mengambil scope lewat `FindEmployeeDropdownScope(EmployeeId, CustId)` lalu `NormalizeScopeSet`.
- Principal response sekarang mengisi `UserFullname: userInfo.UserFullname` dari `sys.m_user.user_fullname`.
- Repository `FindUserByUsername` query: `SELECT user_id, user_fullname FROM sys.m_user WHERE user_name = $1 AND is_del = false`.
- DTO `BusinessUnitPrincipalResponse` belum punya `cust_name`.
- Query distributor SX-2079 sudah memakai `DISTINCT`, mapping joins untuk `specific`, `parent_cust_id` fallback, dan `IN (?)` untuk multi-value `region_id` / `area_id`.
- `FindEmployeeDropdownScope` saat ini hanya select scope dari `mst.m_employee`, belum join `smc.m_customer`.

Kandidat reuse:

- Pakai interface `BusinessUnitRepository` sebagai tempat method customer lookup baru agar service tetap Controller → Service → Repository.
- Pakai pola `repo.Get(&model, query, args...)` dari `FindUserByUsername`.
- Pakai model baru kecil `CustomerInfo` atau extend `UserInfo` tidak direkomendasikan karena sumber data beda.
- Bisa alternatif extend `FindEmployeeDropdownScope` join `smc.m_customer`, tapi interface dipakai region/area juga; perubahan model employee berisiko lebih luas.

Constraint:

- Jaga service layer contract dari `.opencode/docs/ARCHITECTURE.md`.
- Jaga tenant rules: `cust_id` untuk current principal, `parent_cust_id` untuk parent-company dropdown distributor scope.
- Jangan ubah SX-2079 query scope kecuali test membuktikan perlu.
- Jangan commit token/password/header auth staging.

Risiko:

- `cust_name` field baru bisa memengaruhi FE hanya bila FE strict schema, kecil untuk JSON additive field.
- Missing customer row perlu perilaku eksplisit. Karena controller sudah map `sql.ErrNoRows` ke 404, default aman adalah bubble-up error dari repository; jangan fallback ke user fullname tanpa keputusan produk.
- Distributor path masih pakai `userInfo.UserFullname`; jangan ubah kecuali ada bug sama.

Research gate:

- Local project discovery: dilakukan, wajib untuk rencana implementasi.
- Official docs/context7: tidak diperlukan; isu memakai Go/sqlx internal repo, bukan API library baru.
- GitHub: tidak diperlukan; tidak bergantung upstream.
- Web search: tidak diperlukan; fakta cukup dari Jira prompt dan repo lokal.
- Browser/screenshot: tidak diperlukan; BE API defect, bukan visual parity.
