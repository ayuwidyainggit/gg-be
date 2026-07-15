# Plan SX-2003 — Filter `employee-pjp` hanya role Salesman

## Goal

Endpoint `GET /master/v1/employee-pjp` hanya mengembalikan employee aktif sesuai filter request dan role `salesman`, untuk dropdown Monitoring Activity.

## Non-goals

- Tidak ubah kontrak response: tetap `emp_id`, `emp_code`, `emp_name`.
- Tidak ubah UI Monitoring Activity.
- Tidak ubah data role/user/employee di DB.
- Tidak ubah endpoint `GET /master/v1/employees`, `employee-lookup`, atau endpoint PJP lain.
- Tidak simpan token/password staging dari ticket ke test, log, atau artifact.

## Scope

Target module: `master`.

Perubahan utama:
- Update query `FindAllForPJP` agar filter employee yang punya role `salesman` lewat tabel `sys.m_user`, `sys.user_roles`, `sys.m_role`.
- Tetap hormati filter existing: `cust_id` JWT/query, `distributor_id`, `is_active[]`, `q`, `sort`, `page`, `limit`.
- Tambah regression tests untuk query builder/filter role.
- Perkuat query builder agar parameterized untuk nilai user-controlled (`cust_id`, `q`, list cust_id) bila refactor dilakukan.

## Requirements

- `GET /master/v1/employee-pjp` harus mengecualikan employee tanpa role `salesman`.
- Filter role harus tenant-aware:
  - `mu.emp_id = me.emp_id`
  - `mu.cust_id = me.cust_id`
  - `ur.user_id = mu.user_id`
  - `ur.cust_id = mu.cust_id`
  - `mr.role_id = ur.role_id`
  - `mr.cust_id = ur.cust_id`
- Filter role case-insensitive: `LOWER(mr.role_name) = 'salesman'`.
- Principal/distributor scope tetap sesuai behavior endpoint saat ini:
  - `distributor_id > 0` → employee dari mapped distributor `cust_id` via `smc.m_customer.distributor_id`.
  - `cust_id` query tidak kosong → employee dari `cust_id` tersebut.
  - default → employee dari JWT `cust_id`.
- Jika implementasi memperluas skenario principal+distributor list sesuai dokumen, lakukan tanpa bocor lintas `parent_cust_id`.
- Pagination total harus menghitung employee unik, bukan row role duplikat.

## Acceptance Criteria

- Endpoint `GET /master/v1/employee-pjp` hanya return employee yang punya role `salesman`.
- Employee non-salesman tidak muncul.
- Filter berlaku untuk distributor, principal-only, dan principal+distributor sesuai request scope yang endpoint dukung.
- `q`, `is_active[]`, `page`, `limit`, `sort`, `cust_id`, `distributor_id` tetap berfungsi.
- Tidak ada duplicate employee akibat employee multi-role.
- Tidak ada regresi endpoint lain yang tidak memakai `FindAllForPJP`.
- `rtk go test ./...` di module `master` pass atau failure unrelated terdokumentasi.

## Existing Patterns/Reuse

- Route ditemukan di `master/controller/employee_controller.go` line 52-53: `/v1/employee-pjp` → `ListPJP`.
- Handler `ListPJP` line 139-179 parse `EmployeePJPQueryFilter`, inject `CustId`/`ParentCustId` dari JWT, panggil service.
- Service `master/service/employee_service.go` line 788-803 hanya mapping model → response; query logic ada di repository.
- Repository target `master/repository/employee_repository.go` line 999-1091 membangun `FindAllForPJP` tanpa role join/filter.
- Reuse pola test repository query builder dari:
  - `master/repository/business_unit_repository_test.go`
  - `master/repository/salesman_repository_test.go`
- Reuse pola safer query building dari `business_unit_repository.go`: `sqlx.In` + `Rebind` + args.
- Reuse tenant scope idea dari `buildSalesmanCustScopeCondition` untuk principal/distributor, tapi jangan copy buta karena alias table beda (`me` bukan `s`).

Tidak ada util KiloCode/project khusus untuk role filter employee yang langsung bisa dipakai. Buat helper query kecil di `employee_repository.go` lebih baik daripada reimplement global abstraction.

## Constraints

- Layering wajib: Controller → Service → Repository → DB.
- Schema prefix wajib: `mst.`, `smc.`, `sys.`.
- Shell validation tetap `rtk`-prefixed.
- Validate dari directory module `master`, bukan root.
- Jangan commit/sebar token staging dari issue.
- Jangan ubah source di luar `master` kecuali test/plan evidence minta eksplisit.
- Repo berisi `.env` dan credential infra tracked; jangan buka/copy/expand.

## Risks

- Query lama memakai string concat; menambah `IN ('...')` manual akan memperbesar risiko injection jika `cust_id` dari query FE tidak diparameterkan.
- Multi-role employee bisa duplicate bila memakai raw `JOIN`. Mitigasi: gunakan `EXISTS` filter role atau `SELECT DISTINCT me.emp_id...` + count distinct.
- `sort=area_id:asc` dari FE bisa mengacu kolom yang tidak ada di select/from query employee minimal; existing behavior mungkin sudah rentan. Jangan ubah behavior besar kecuali test menunjukkan error.
- Case DB `role_name` mungkin `Salesman`/`salesman`; mitigasi dengan `LOWER`.
- Principal+distributor skenario di dokumen memakai `cust_id IN (...)`; endpoint existing hanya punya `DistributorId *int`, bukan array. Perlu hati-hati: jangan mengubah API parsing besar tanpa keputusan.

## Decisions/Assumptions

- Pakai `EXISTS` subquery untuk role filter, bukan raw join di main query, agar employee unik walau multi-role.
- Pakai `LOWER(mr.role_name) = 'salesman'` untuk case-insensitive.
- Pertahankan response tanpa `role_name`, walau SQL referensi select `mr.role_name`; acceptance criteria hanya response employee list.
- Pertahankan default `is_active` behavior existing: aktif hanya difilter bila `is_active[]=1` dikirim. FE Monitoring Activity memang mengirim `is_active[]=1`. Jika product ingin selalu aktif meski param kosong, butuh keputusan terpisah.
- Tidak perlu official docs/context7; ini query SQL internal dan pola repo lokal cukup.
- Tidak perlu GitHub/web/browser; tidak bergantung upstream/external/reference UI.
- Tidak perlu `@architect`; perubahan bounded repository filter tanpa data model/migration.
- Pertanyaan tidak diajukan: requirement cukup jelas untuk plan. Ambiguitas `principal+distributor` dicatat sebagai implementation caveat karena endpoint sekarang `DistributorId *int`, bukan array.

## TDD/Test Plan

TDD required: ya. Ini perubahan query production, tenant scope, dan security-sensitive role filtering.

Existing test patterns:
- `master/repository/business_unit_repository_test.go` test query builder contains + args length.
- `master/repository/salesman_repository_test.go` test helper builder conditions.
- `github.com/DATA-DOG/go-sqlmock` tersedia di `master/go.mod`.

First failing/regression tests:

1. `TestBuildEmployeePJPQuery_AppliesSalesmanRoleExistsFilter`
   - Ekstrak query builder dari `FindAllForPJP` ke helper, misalnya `buildEmployeePJPQuery(dataFilter)`.
   - Assert count/select query berisi `EXISTS`, `sys.m_user`, `sys.user_roles`, `sys.m_role`, `mu.emp_id = me.emp_id`, tenant joins, dan `LOWER(mr.role_name) = 'salesman'`.
   - Assert tidak ada raw `JOIN sys.m_role` di main query jika pakai `EXISTS`.

2. `TestBuildEmployeePJPQuery_DistributorScopeUsesCustomerMapping`
   - Input `DistributorId=67`, `ParentCustId="C26002"`, `CustId="C260020001"`.
   - Assert query scope memakai `smc.m_customer` dan `distributor_id` arg; bila hardening dilakukan, harus constrain `parent_cust_id`.

3. `TestBuildEmployeePJPQuery_FilterCustIDUsesArgsNotConcatenation`
   - Input `FilterCustId="C26002' OR '1'='1"`.
   - Assert query tidak mengandung literal injection string dan args memuat value.

4. Opsional integration-style with `sqlmock`:
   - Expect count query + select query; return one salesman row; assert service response mapped.

Green step:
- Refactor `FindAllForPJP` ke builder helper yang menghasilkan `countQuery`, `countArgs`, `selectQuery`, `selectArgs`, `limit`, `offset`.
- Tambah role filter `EXISTS` dengan `LOWER(mr.role_name) = 'salesman'`.
- Parameterize `cust_id`, search, distributor id, active flag.
- Gunakan `sqlx.In` + `Rebind` hanya jika ada list/slice; untuk scalar pakai `?` lalu `Rebind`.
- Eksekusi query dengan args.

Refactor step:
- Jika helper terlalu kompleks, split kecil:
  - `buildEmployeePJPRoleExistsClause()`
  - `buildEmployeePJPScopeClause(filter)`
  - `buildEmployeePJPSortClause(sort)`
- Jangan abstraksi lintas repository sekarang.

Edge cases:
- Employee punya user tapi role bukan `salesman` → excluded.
- Employee tidak punya `sys.m_user` → excluded.
- Employee punya beberapa roles termasuk `salesman` → included satu kali.
- `role_name='Salesman'` → included karena `LOWER`.
- Empty result → existing controller returns `data` empty or No Data only if repository returns err with total 0; jangan ubah behavior kecuali existing tests minta.

Commands:

```bash
rtk go test ./repository -run 'TestBuildEmployeePJPQuery'
rtk go test ./service -run 'Test.*PJP'
rtk go test ./...
```

## Implementation Steps

1. Tambah tests di `master/repository/employee_repository_test.go` untuk query builder role filter dan scope.
2. Ekstrak query construction dari `FindAllForPJP` menjadi helper testable.
3. Tambah `EXISTS` role filter:
   ```sql
   EXISTS (
     SELECT 1
     FROM sys.m_user mu
     JOIN sys.user_roles ur
       ON ur.user_id = mu.user_id
      AND ur.cust_id = mu.cust_id
     JOIN sys.m_role mr
       ON mr.role_id = ur.role_id
      AND mr.cust_id = ur.cust_id
     WHERE mu.emp_id = me.emp_id
       AND mu.cust_id = me.cust_id
       AND LOWER(mr.role_name) = 'salesman'
   )
   ```
4. Parameterize existing filters:
   - `me.cust_id = ?`
   - distributor subquery `mc.distributor_id = ?`
   - optional `mc.parent_cust_id = ?` if using `ParentCustId` to avoid cross-principal distributor id collision.
   - `q` uses args `%query%`.
5. Ensure count query counts unique employees. If using `EXISTS`, `COUNT(*)` remains okay.
6. Keep select fields `me.emp_id, me.emp_code, me.emp_name` and response unchanged.
7. Run targeted repository tests.
8. Run full `master` tests.
9. Manual smoke after deploy using fresh token only; do not store token.
10. `@quality-gate` review for tenant/security regression.

## Expected Files to Change

Primary:
- `master/repository/employee_repository.go`
- `master/repository/employee_repository_test.go` (new)

Likely unchanged:
- `master/controller/employee_controller.go`
- `master/service/employee_service.go`
- `master/entity/employee.go`
- `master/model/employee.go`

No planned changes:
- migrations
- `go.mod` / `go.sum`
- docs outside `.opencode/`
- PJP services

## Agent/Tool Routing

- `@orchestrator`: run plan, coordinate implementation + validation.
- `@fixer`: implement bounded repository/test changes.
- `@explorer`: only if implementation discovers endpoint reuse/scope conflicts.
- `@quality-gate`: final signoff because tenant/security role filter.
- `@architect`: not needed; no data model/migration/product architecture change.
- `@librarian`/context7/GitHub/web/browser: not needed; local SQL + docs enough.

## Execution-ready Worklist / Handoff Contract

`start_with`: `SX2003-01`

| id | action | depends_on | owner/lane | validation/check | exit criteria | status | blocker | requires_user_decision |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| SX2003-01 | Add repository query-builder tests for role `salesman` filter and scope behavior. | none | `@fixer` | `rtk go test ./repository -run 'TestBuildEmployeePJPQuery'` from `master` | Tests fail before implementation or direct-green documented. | ready | none | no |
| SX2003-02 | Extract `FindAllForPJP` query construction into helper returning queries and args. | SX2003-01 | `@fixer` | same targeted command | Helper testable; existing behavior preserved except planned role filter. | ready | none | no |
| SX2003-03 | Add tenant-aware `EXISTS` role filter for `LOWER(mr.role_name) = 'salesman'`. | SX2003-02 | `@fixer` | `rtk go test ./repository -run 'TestBuildEmployeePJPQuery'` | Query filters role without duplicate employee rows. | ready | none | no |
| SX2003-04 | Parameterize cust/search/distributor filters while preserving endpoint behavior. | SX2003-03 | `@fixer` | targeted repository tests | No user-controlled raw string literals in PJP query; tests cover injection-shaped cust_id. | ready | none | no |
| SX2003-05 | Run full master module tests. | SX2003-04 | `@fixer` | `rtk go test ./...` from `master` | Pass or unrelated failures documented with exact output. | ready | none | no |
| SX2003-06 | Manual smoke endpoint after deploy with fresh distributor/principal tokens. | SX2003-05 | `@orchestrator` | curl endpoint, inspect all rows role via DB or response-backed check | Only salesman employees returned for distributor and principal users. | blocked | Needs deployed build and valid non-persisted staging token. | yes |
| SX2003-07 | Final quality/security review. | SX2003-05 | `@quality-gate` | Review diff, tests, tenant scope, SQL injection risk | No blocking regression or documented accepted risk. | ready | none | no |

## Validation Commands

From repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

From `master` directory:

```bash
rtk go test ./repository -run 'TestBuildEmployeePJPQuery'
rtk go test ./service -run 'Test.*PJP'
rtk go test ./...
```

Manual smoke after fixed build deploy, fresh token only:

```bash
curl "https://best.scyllax.online/master/v1/employee-pjp?q=&page=1&limit=9999&sort=area_id:asc&is_active[]=1&region_id[]=1&cust_id=10&distributor_id=5" \
  -H "Authorization: Bearer <token>" \
  -H "Accept: application/json"
```

Response checks:
- `data` contains only employees that DB join resolves to `LOWER(mr.role_name) = 'salesman'`.
- Non-salesman employee sampled from same `cust_id` absent.
- No duplicate `emp_id`.
- Paging `total_record` equals filtered count.

Optional read-only DB verification query after deploy:

```sql
SELECT me.emp_id, me.emp_code, me.emp_name, mr.role_name
FROM mst.m_employee me
JOIN sys.m_user mu
  ON mu.emp_id = me.emp_id
 AND mu.cust_id = me.cust_id
JOIN sys.user_roles ur
  ON ur.user_id = mu.user_id
 AND ur.cust_id = mu.cust_id
JOIN sys.m_role mr
  ON mr.role_id = ur.role_id
 AND mr.cust_id = ur.cust_id
WHERE me.cust_id IN (<scoped_cust_ids>)
  AND me.is_del = false
  AND me.is_active = true
  AND LOWER(mr.role_name) = 'salesman';
```

## Evidence Requirements

Kept evidence:
- `.opencode/evidence/20260521-1019-sx-2003-employee-pjp-salesman-filter/discovery.md` — local discovery, files inspected, reuse candidates, constraints, risks.

Implementation evidence needed:
- Diff summary for `master/repository/employee_repository.go` and tests.
- Targeted test output.
- Full `master` test output.
- Manual smoke after deploy with fresh token, no token saved.

Research gate:
- Local project discovery: done, required.
- Official docs/context7: skipped, no version-sensitive library behavior.
- GitHub: skipped, no upstream dependency.
- Brave/web: skipped, internal Jira/doc + local repo enough.
- Browser/screenshot: skipped, backend-only endpoint.

## Done Criteria

- Primary plan followed or deviations documented.
- Tests added before/with fix and passing.
- `FindAllForPJP` filters role `salesman` tenant-aware.
- Query avoids duplicate employee rows from multi-role joins.
- Query avoids raw user-controlled string interpolation for modified filters.
- Acceptance criteria verified locally and smoke-ready for staging.
- `@quality-gate` signoff complete for tenant/security risk.

## Final Planning Summary

Artifacts created:
- Primary plan: `.opencode/plans/20260521-1019-sx-2003-employee-pjp-salesman-filter.md`
- Evidence kept: `.opencode/evidence/20260521-1019-sx-2003-employee-pjp-salesman-filter/discovery.md` because implementation should reuse exact file/line findings and risk notes.

Artifacts cleaned:
- No draft artifacts created.
- Evidence kept intentionally; not stale.

Key decisions:
- Use `EXISTS` role filter instead of raw main join to avoid duplicate employee rows.
- Use `LOWER(mr.role_name) = 'salesman'`.
- Keep response unchanged.
- Keep `is_active` behavior tied to query param to avoid accidental endpoint contract change.

Assumptions:
- Existing endpoint scope behavior for `cust_id`/`distributor_id` is acceptable unless FE/product asks array distributor/principal+distributor expansion.
- `sys.user_roles` and `sys.m_role` tenant keys match docs and usage in repo.

Open questions:
- None blocking for implementation.
- Non-blocking: whether API should support multiple distributor ids/principal+all distributors in this endpoint later. Current ticket can be satisfied by current scope + role filter.

Readiness:
- Ready for `@orchestrator` → `@fixer` implementation starting at `SX2003-01`.
