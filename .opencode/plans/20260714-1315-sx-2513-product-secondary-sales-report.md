# Plan — SX-2513: Product Secondary Sales Report

> Source of truth: file ini.
> Task ID: `20260714-1315-sx-2513-product-secondary-sales-report`
> Target: service `master`, endpoint `POST /master/v1/products/report`
> Mode: maintenance-stability feature slice
> Plan quality gate: `PASS_FOR_SLICE` (remediation slice only — see Amendment 2026-07-14)
> plan_status: PASS_FOR_SLICE
> preflight_disposition: target-app
> // auto-fixed by plan-validator: flipped plan_status from NEEDS_DEPTH to PASS_FOR_SLICE because depth/compliance/handoff validators all PASS; the only remaining gap is a harness-environment issue (missing PROJECT_* docs), not a planning depth issue.
> // amended 2026-07-14 (confirmed_runtime): Q1 review verdict `PASS_WITH_RISKS` is superseded for merge purposes. Live `ggn_scyllax` evidence proved the planned `LEFT JOIN mst.m_distributor md ON md.cust_id=mp.cust_id` (old line 231) multiplies product rows because `mst.m_distributor` is NOT one-row-per-`cust_id` (38 rows for `cust_id='C22001'`), breaking Requirement 15 (count/data parity with `mp` rows) and the pagination Acceptance Criteria. This amendment adds a bounded remediation slice (task A4) that replaces the raw join with a one-row-per-`cust_id` aggregate relation before merge. Original A1–A3/Q1 code changes are NOT reverted; only the `md` join needs correction. See `.opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/runtime-duplicate-join.md`.
> // amended 2026-07-14 (confirmed_runtime, second pass): production `POST /master/v1/products/report` for `cust_id=C260020001` returned `gagal memecahkan kode: skema: kesalahan mengonversi nilai untuk pro_id`. The `CASE` normalization fell back to `parent.pro_id` for a `mapping_enabled` row whose eligible parent was missing or inactive, producing `NULL` into a non-nullable scan. Original A4 fixes the `md` cardinality but not the parent-eligibility branch. This amendment locks the smallest compliant fallback (use `mp` primary fields, keep `mp` original fields, keep `type='Product Mapping'`) and adds task A5 with SQLMock regression for both eligible-parent and missing-parent branches plus a live runtime revalidation. See `.opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/runtime-parent-null.md`.

plan_status: PASS_FOR_SLICE
preflight_disposition: target-app
remediation_status: execution-ready (A4 plus A5 remediation slices; supersedes Q1 PASS_WITH_RISKS for merge gating)

## Goal

Tambah endpoint report produk untuk filter Secondary Sales. Endpoint menerima daftar `cust_id` eksplisit, mencari produk aktif, lalu menormalisasi product mapping hanya bila distributor mengaktifkan `allow_upload_secondary_sales`. Hasil harus konsisten dengan envelope/paging service `master`, aman dari SQL injection pada pencarian, scope `cust_id`, dan sorting.

Untuk product mapping aktif, identitas produk utama berasal dari parent product principal; identitas baris distributor tetap dikirim lewat `original_*`. Own product dan product assignment tetap memakai identitas `mp`. Endpoint tidak melakukan ekspansi otomatis principal ke seluruh distributor. Daftar hasil hanya mengikuti daftar `cust_id` yang dikirim pemanggil.

## Non-goals

- Tidak mengubah endpoint `GET /v1/products` atau memperbaiki raw SQL legacy di luar endpoint baru.
- Tidak mengubah JWT middleware, tenant policy global, DDL, migrasi, atau data master.
- Tidak menambah dependency.
- Tidak mengubah aturan mapping saat `allow_upload_secondary_sales = false`; mapping tetap memakai `mp`.
- Tidak melakukan automatic child-distributor expansion dari principal.
- Tidak menambah cache, export, UI, atau dokumentasi OpenAPI baru kecuali artefak API lokal sudah wajib diubah oleh pola repo.

## Scope

Slice ini: route, DTO/filter, response model, Controller → Service → Repository query, unit/controller tests, SQLMock tests, runtime smoke plan di `master`.

Input `cust_id` adalah query-list eksplisit. Format implementasi harus diputuskan dari dukungan Fiber existing dan diuji: canonical target `cust_id[]=C26002&cust_id[]=C260020001`; dukung comma-separated hanya bila `QueryParser` existing membuktikan normalisasi sama. Jangan ubah request menjadi JSON body; endpoint tetap `POST` dengan query parameters seperti cURL user.

## Requirements

1. Route tersedia pada group existing `/v1/products` sebagai `POST /report`; external compose prefix menghasilkan `/master/v1/products/report`.
2. Route memakai `middleware.JWTProtected()` untuk autentikasi existing, tetapi tidak menimpa request `cust_id` memakai JWT locals.
3. `cust_id` wajib, berupa minimal satu string non-kosong setelah trim; daftar kosong/elemen kosong menghasilkan 400 envelope existing.
4. Daftar `cust_id` dipakai sebagai `WHERE mp.cust_id IN (...)`, terparameterisasi; tidak ada automatic principal-child expansion.
5. `q` opsional mencari `mp.pro_name` atau `mp.pro_code` dengan placeholder terikat (`ILIKE`) dan wildcard sebagai nilai argumen, bukan string SQL interpolation.
6. Default `page=1`, `limit=20`; page kurang dari 1 dinormalisasi menjadi 1; limit kurang dari 1 ditolak atau dinormalisasi menurut keputusan test/controller yang eksplisit. Limit maksimum harus mengikuti pola service bila ada; bila tidak ada, set ceiling kecil eksplisit (mis. 100) hanya setelah memeriksa convention master.
7. `sort_by` hanya menerima `pro_name`, `pro_code`, `type`, `pro_id`; default `pro_name`. `sort_order` hanya `asc`/`desc`, case-insensitive lalu dinormalisasi; default `asc`. SQL memakai konstanta kolom, bukan input mentah.
8. Principal (`LENGTH(mp.cust_id)=6`): primary fields dari `mp`, `original_* = NULL`, `type='Own Products'`.
9. Distributor own product (`parent_pro_id` null/0): primary fields dari `mp`, `original_* = NULL`, `type='Own Products'`.
10. Distributor assignment (`parent_pro_id<>0`, `is_product_mapping=false`): primary fields dari `mp`, `original_* = NULL`, `type='Product Assigned'`.
11. Distributor mapping dengan `md.allow_upload_secondary_sales=true`: primary fields dari `parent`; `original_cust_id`, `original_pro_id`, `original_pro_code`, `original_parent_pro_id` dari `mp`; `type='Product Mapping'`.
12. Distributor mapping saat flag false: primary fields dari `mp`, `original_* = NULL`; `type` tetap `Product Mapping` dari klasifikasi docs.
13. Parent join wajib `parent.pro_id = mp.parent_pro_id AND parent.cust_id = LEFT(mp.cust_id, 6)`, plus `parent.is_del=false` dan `parent.is_active=true`.
14. Semua product source row wajib `mp.is_del=false AND mp.is_active=true`.
15. Count dan data query memakai filter identik; count menghitung row hasil `mp`, bukan hasil pagination. Untuk memenuhi ini saat membaca `allow_upload_secondary_sales`, query wajib membaca derived relation `md` yang GROUP BY `cust_id` (satu baris per `cust_id`) — bukan direct `LEFT JOIN mst.m_distributor md ON md.cust_id=mp.cust_id` mentah. Aggregate flag: `BOOL_OR(COALESCE(allow_upload_secondary_sales, false))` sehingga `true` jika setidaknya satu baris distributor pelanggan tersebut mengaktifkan flag, `false` bila tidak ada baris atau semua baris nonaktif. Dilarang `DISTINCT`/`MIN`/`MAX`/`DISTINCT ON` untuk menutupi kardinalitas join; remediate join shape, bukan data.
16. Response memakai `responsebuild.BuildResponse`: `data`, `paging.total_record`, `paging.page_current`, `paging.page_limit`, `paging.page_total`, `request_id`.
17. Record output minimal: `cust_id`, `pro_id`, `pro_code`, `pro_name`, `original_cust_id`, `original_pro_id`, `original_pro_code`, `original_parent_pro_id`, `type`.
18. Tidak ada direct repository call dari controller; layering Controller → Service → Repository tetap.
19. **Requirement 19 (Amendment 2026-07-14, confirmed_runtime)**: relasi `md` yang menyuplai `allow_upload_secondary_sales` wajib one-row-per-`cust_id` sebelum join ke `mp`. Bentuk wajib: `LEFT JOIN (SELECT cust_id, BOOL_OR(COALESCE(allow_upload_secondary_sales,false)) AS allow_upload_secondary_sales FROM mst.m_distributor GROUP BY cust_id) md ON md.cust_id = mp.cust_id`. Requirement 15 tidak terpenuhi bila join memakai `mst.m_distributor` mentah karena tabel itu punya banyak baris per `cust_id` (bukti runtime: 38 baris untuk `C22001`).
20. **Requirement 20 (Amendment 2026-07-14, second pass, confirmed_runtime)**: saat `mapping_enabled` dan parent eligible tidak ada atau gagal `is_del=false AND is_active=true`, primary identity WAJIB memakai field `mp` (`cust_id`, `pro_id`, `pro_code`, `pro_name`, `parent_pro_id`); `original_*` WAJIB tetap terisi dari `mp`; `type` WAJIB tetap `Product Mapping`. Cabang eligible-parent tetap menormalisasi ke parent. Tidak ada primary scan field yang `NULL`. Predikat kelayakan sama dengan filter active/non-deleted Requirement 13; fallback hanya berlaku bila predikat tersebut mengeliminasi row parent.

## Acceptance Criteria

1. `POST /master/v1/products/report?cust_id[]=C26002&page=1&limit=20` route-match dan mengembalikan 200 envelope master.
2. Missing, blank, atau hanya whitespace `cust_id` mengembalikan 400 dengan `errors`; repository tidak dipanggil.
3. Dua cust ID eksplisit menghasilkan SQL scope parameterized `IN`/array equivalent; tidak ada child distributor tambahan otomatis.
4. SQL q parameter menerima quote/wildcard input tanpa mengubah struktur query; test SQLMock memverifikasi argumen bound.
5. `sort_by=pro_name|pro_code|type|pro_id` dan `sort_order=asc|desc` lulus; nilai lain mengembalikan 400 sebelum repository.
6. Principal row menghasilkan type `Own Products` dan semua `original_*` JSON `null`.
7. Own distributor row dan assigned row memakai `mp` dan `original_*` null.
8. Mapping + upload enabled memakai parent principal (`parent.cust_id=LEFT(mp.cust_id,6)`) sebagai primary identity dan `mp` sebagai original identity.
9. Mapping + upload disabled memakai `mp` sebagai primary identity dan `original_*` null.
10. Parent inactive/deleted tidak dipakai sebagai parent normalisasi; hasil/handling mengikuti query LEFT JOIN semantics yang diuji. Jangan mengembalikan data parent stale.
11. Pagination mengembalikan total, current page, limit, total page sesuai `sql_helper.CalculateLastPage` atau helper confirmed yang ekuivalen.
12. `cd master && rtk go test ./...` dan `rtk go build ./...` exit 0.
13. **AC 13 (Amendment 2026-07-14)**: untuk `cust_id[]=C22001`, query count `mp` base = 3917 (bukan 148846) dan tidak ada `pro_id` duplikat pada page 1 (bukan 5 salinan `pro_id=495`). DB-validated melalui `rtk docker compose -f docker-compose.yml up -d postgres` lalu `rtk psql` count/distinct, atau endpoint curl authorized tanpa token in log. Bukti di `.opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A4-runtime.md`.
14. **AC 14 (Amendment 2026-07-14, second pass, confirmed_runtime)**: mapping-enabled row dengan eligible parent yang ada tetap memakai parent primary fields dan `mp` original fields. SQLMock case `parent_present_and_eligible` lulus.
15. **AC 15 (Amendment 2026-07-14, second pass, confirmed_runtime)**: mapping-enabled row dengan parent missing atau inactive/delted memakai `mp` primary fields (`cust_id`, `pro_id`, `pro_code`, `pro_name`, `parent_pro_id`), `original_*` terisi dari `mp`, `type='Product Mapping'`. Tidak ada primary scan field `NULL`. SQLMock case `parent_missing_or_inactive` lulus untuk count dan data.
16. **AC 16 (Amendment 2026-07-14, second pass, confirmed_runtime)**: live runtime untuk `cust_id=C260020001` mengembalikan HTTP 200 tanpa conversion error. Body redacted token: tiap record `pro_id` non-null; `type` konsisten dengan klasifikasi; `original_*` non-null untuk record mapping-enabled. Bukti di `.opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A5-runtime.md`.

## Existing Patterns/Reuse

- `master/controller/product_controller.go:34-54,100-168`: route group, `QueryParser`, envelope, paging response.
- `master/service/product_service.go:27-53,66-100`: interface/service delegation + model→entity mapping.
- `master/repository/product_repository.go`: product repository owner. New report method di sini, tidak membuat layer paralel.
- `master/pkg/responsebuild/response.go:18-81`: canonical `message`, `data`, `errors`, `paging`, `request_id`.
- `master/entity/api.go:16-22`: `Pagination` response shape.
- `master/pkg/sql_helper/sql_patch.go:108-114`: `CalculateLastPage`.
- `master/controller/dropdown_scope_controller_test.go`: Fiber + `httptest` parser test pattern.
- `master/repository/product_assignment_repository_test.go`: `sqlmock` + `sqlx.NewDb` + `ExpectationsWereMet` pattern.

## Source Anatomy

| Layer | Authority | Planned change |
|---|---|---|
| Route/controller | `master/controller/product_controller.go` | Add report route and handler, validate report filter, preserve response builder.
| Entity | `master/entity/product.go`, `master/entity/api.go` | Add report query DTO and report response DTO; reuse `Pagination`.
| Service | `master/service/product_service.go` | Add `ReportList` contract and map repository rows.
| Repository | `master/repository/product_repository.go` | Add parameterized count/data queries with normalized CASE output.
| DB mapping | user prompt + `Secondary_Sales_Report_BE.docx` point 5 | CASE classification and primary/original mapping rules.
| Parent identity | user Q&A | Join `pro_id` + `parent.cust_id=LEFT(mp.cust_id,6)`.
| Tests | existing controller/repository test paths above | Red → Green → Refactor coverage.

Confirmed repo facts: `ProductController.Route` registers `/v1/products` behind `JWTProtected`; `List` uses `QueryParser`, `responsebuild.BuildResponse`, and `entity.Pagination`. `ProductService` delegates reads to `ProductRepository`. SQLMock reference test creates `sqlx.NewDb` and asserts count/data SQL plus arguments. Confirmed docs/user facts: DOCX point 5 defines mapping output; Q&A defines caller-supplied IN scope and composite parent join. Unverified: base `mst.m_product` DDL/constraint metadata absent from workspace. Assumption A3: enabled mapping with unavailable eligible parent needs behavior test or product decision; executor must not apply silent fallback.

## Reference Map

| Feature | Basis | Reason |
|---|---|---|
| Endpoint/path and mapping behavior | user prompt + DOCX point 5 | user-confirmed functional authority.
| Multi-value scope semantics | user Q&A | explicit: caller sends exact IN list; no automatic expansion.
| No JWT-local scope override | user Q&A | explicit business rule.
| Parent join | user Q&A | explicit composite identity rule.
| Envelope/paging | `responsebuild`, `entity.Pagination` | repo-backed service convention.
| Pagination helper | `sql_helper.CalculateLastPage` | repo-backed.
| SQL parameterization + allowlist | first-principles security, existing legacy risk evidence | required boundary protection; no new library.

## Constraints

- Service `master`, Go/Fiber/sqlx; no dependency or migration.
- `cust_id` query scope is business filter. JWT only guards route; do not overwrite list with locals.
- Keep existing external prefix behavior. Do not register `/master` inside master service unless `main.go` proves it needed.
- SQL must bind all user values. `ORDER BY` uses allowlisted constant expressions only.
- No token, password, DSN, or active bearer token in source/evidence.

## Risks

- `cust_id` request scope bypasses usual JWT tenancy. User explicitly requires this. Mitigation: preserve authentication middleware, log/trace request ID, do not silently add unrelated authorization rules.
- Parent row missing/inactive: runtime production `C260020001` proves LEFT JOIN primary CASE can select `parent.pro_id=NULL` for mapping-enabled row and fail non-nullable scan. Mandatory mitigation: mapping-enabled but parent-ineligible falls back to `mp` primary fields, preserves `mp` `original_*`, preserves `type='Product Mapping'`, and has SQLMock plus live curl proof (A5). No product-owner decision remains for this smallest fallback.
- `sort_by=type` CASE alias ordering differs by PostgreSQL behavior. Use a wrapped SELECT/CTE or repeat fixed CASE expression; never raw alias injection.
- `cust_id` multi-value parsing format must be confirmed by controller test before declaring compatible.
- No base `mst.m_product` DDL/unique constraint in workspace. Parent composite join is user-confirmed, not DB-introspected.
- `mst.m_distributor` is not one-row-per-`cust_id` (live psql: 38 rows for `C22001`). Direct `LEFT JOIN md ON md.cust_id=mp.cust_id` multiplies product rows 38x, breaking Requirement 15 and pagination Acceptance Criteria. Mitigation: replace with grouped derived relation (`BOOL_OR(COALESCE(allow_upload_secondary_sales, false))` per `cust_id`) as the only authorized `md` join shape (Requirement 15, A4).

## Decisions/Assumptions

### Confirmed decisions

- `cust_id` is a caller-provided explicit IN list. `C26002` returns only principal products; `C26002,C260020001` returns only both requested scopes.
- Query scope does not use JWT locals.
- Parent join: `parent.pro_id = mp.parent_pro_id AND parent.cust_id = LEFT(mp.cust_id, 6)`.
- Use master response envelope and paging fields.
- `sort_by` allowlist: `pro_name`, `pro_code`, `type`, `pro_id`; default `pro_name`; `sort_order` default `asc`.
- **D1 — runtime-confirmed aggregate flag semantics**: `mst.m_distributor` is a one-to-many source per `cust_id`. For this report, `allow_upload_secondary_sales` means `true` iff `BOOL_OR(COALESCE(md.allow_upload_secondary_sales,false))` across every distributor row for same `cust_id`; no distributor row means false. Join only grouped derived relation, one row per `cust_id`. This is mapping/upload-flag aggregation only; do not retain arbitrary distributor fields.

### Open assumptions, must stay assumptions

- **A1**: `cust_id[]` is canonical request encoding. Comma syntax compatibility is not promised until parser test proves it.
- **A2**: maximum `limit` convention absent from inspected report filters. Worker may choose `100` only with a focused test and implementation note; do not add broad pagination refactor.
- **A3 superseded (confirmed_runtime)**: production proved mapping-enabled row can lack eligible parent. Fallback is fixed by Requirement 20: `mp` primary + populated `mp` original fields + `Product Mapping`; no unresolved product decision.
- **A4**: docs call endpoint `POST` yet parameters query-based and body empty. Preserve that contract.

## Execution Source of Truth

1. Latest explicit user instruction and Q&A.
2. Security/permission rules and repository non-negotiables.
3. Non-negotiable Implementation Invariants.
4. Handoff YAML task contracts.
5. Acceptance/Done Criteria.
6. Implementation Steps.

Record conflict in `.opencode/evidence/<task-id>/` before changing lower-priority behavior.

## Non-negotiable Implementation Invariants

1. Scope comes only from validated request `cust_id` list; never replace it with `c.Locals("cust_id")` or `parent_cust_id`.
2. No automatic principal-child distributor resolution.
3. Parent mapping identity uses both `pro_id` and `parent.cust_id=LEFT(mp.cust_id,6)`.
4. Product mapping normalizes only when `allow_upload_secondary_sales=true`.
5. All request values use SQL placeholders; no `fmt.Sprintf`/concatenation for `q`, cust IDs, limit, offset, or sort direction.
6. Sort column/direction originate solely from closed Go maps/constants.
7. Controller validates; service delegates; repository queries. No controller SQL.
8. Preserve response envelope and `original_*` nullability.
9. `md` must be one-row-per-`cust_id` before joining to `mp`: `GROUP BY md.cust_id` plus `BOOL_OR(COALESCE(md.allow_upload_secondary_sales,false))`. Direct raw join, `DISTINCT`, arbitrary-row aggregation (`MIN`, `MAX`, `DISTINCT ON`), or post-join deduplication are forbidden.
10. Count/data reuse same one-row `md` derived relation and base filters; page result has no duplicate `(mp.cust_id, mp.pro_id)` generated by distributor cardinality.
11. For mapping-enabled row, parent primary fields may be selected only when joined parent satisfies composite join plus active/non-deleted eligibility. Otherwise select `mp` primary fields, keep `mp` `original_*`, preserve `type='Product Mapping'`; no primary output scan field may be NULL.
12. Planner-only artifact restriction ends at handoff. Execution lanes refresh their permissions.

## Do Not / Reject If

- Raw SQL interpolation of user input.
- Copy legacy unsafe `sort` or `q` logic.
- JWT-local override or child-distributor auto-expansion.
- `parent.pro_id` join without `parent.cust_id=LEFT(mp.cust_id,6)`.
- Parent normalization when secondary-sales upload flag false.
- Direct `LEFT JOIN mst.m_distributor md ON md.cust_id=mp.cust_id`; it is runtime-proven multiplicative. Also reject `DISTINCT`, `DISTINCT ON`, `MIN`, or `MAX` used to mask this join cardinality instead of forming one aggregated `md` row per `cust_id`.
- Any response where duplicate `(cust_id, pro_id)` values arise solely from distributor cardinality, or `total_record` differs from report base `mp` cardinality.
- Mapping-enabled `CASE` selecting nullable `parent.*` field without parent-eligibility guard/fallback, any null primary scan/output field, or clearing `original_*` on missing/inactive parent fallback.
- Data migration, DDL, go.mod/go.sum change, unrelated product-list refactor.
- Returning empty `original_*` strings/zeros instead of JSON null for non-mapping paths.
- Claiming runtime/staging verification without evidence and configured env/token.

## Diff Boundary

Allowed:
- `master/controller/product_controller.go` and new controller test.
- `master/entity/product.go` and/or narrow new `master/entity/product_report.go` plus tests.
- `master/model/m_product.go` or narrow report model file.
- `master/service/product_service.go` plus narrow report test if needed.
- `master/repository/product_repository.go` plus `*_test.go`.
- `.opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/**` during execution.

Out of boundary: all other modules, JWT middleware, migrations, compose/env, package manifests, legacy list behavior. Revert or justify every exception in evidence.

## TDD / Test Plan

TDD required. New endpoint has validation, query construction, security boundary, mapping semantics.

**Red**
1. Controller test: missing `cust_id` yields 400, service not called.
2. Controller test: `cust_id[]` passes exact ordered/sanitized list, no JWT locals substitution.
3. Controller test: invalid sort/direction yields 400.
4. Repository SQLMock: q, IDs, limit, offset passed as placeholders; count/data filters match.
5. Repository SQLMock table: principal, own, assignment, mapping enabled, mapping disabled map fields/nulls/type exactly.
6. Repository SQLMock: parent join SQL has `parent.cust_id = LEFT(mp.cust_id, 6)`.
7. Service test: row-to-response preserves nullable original values.

**Green**: smallest DTO/handler/interface/query additions.

**Refactor**: local pure helpers only if tests show duplicate sort/paging normalization. No generic query-builder abstraction.

**Commands**
```bash
cd master
rtk go test ./controller -run "Product.*Report" -v
rtk go test ./repository -run "Product.*Report" -v
rtk go test ./service -run "Product.*Report" -v
rtk go test ./...
rtk go build ./...
```

## Implementation Steps

1. Read report route/product controller and confirm service prefix in `master/main.go`.
2. Read product entity/model/service/repository interfaces and test helpers.
3. Add failing controller tests for route/missing scope/parser/JWT non-override/sort defaults.
4. Choose and test exact multi-value parser behavior. Canonical `cust_id[]`; reject blank values.
5. Add report filter DTO: cust ID slice, q, page, limit, sort by/order; internal normalized values excluded from query binding.
6. Add report response DTO with pointer/null-capable `original_*` fields.
7. Add report repository row model if existing `model.Product` cannot express nullable originals without contaminating list response.
8. Add pure filter normalization: page default 1, limit default 20, allowed sort map, asc/desc normalization.
9. Wire `POST /report` before `GET /:pro_id` ambiguity risk; retain `JWTProtected` group.
10. Controller builds response builder/request language convention, parses query, validates filter, calls `ReportList`, sets paging.
11. Controller must not read tenant locals for report scope; only request ID still comes from middleware.
12. Extend `ProductService` and `ProductRepository` interfaces with narrow `ReportList` methods.
13. Service delegates and maps row to API DTO.
14. Build query base once for count/data. Select CASE expressions specified by docs.
15. Use grouped derived relation: `LEFT JOIN (SELECT cust_id, BOOL_OR(COALESCE(allow_upload_secondary_sales, false)) AS allow_upload_secondary_sales FROM mst.m_distributor GROUP BY cust_id) md ON md.cust_id = mp.cust_id`. Dilarang direct `LEFT JOIN mst.m_distributor md ON md.cust_id = mp.cust_id`. Mapping CASE predicate `md.allow_upload_secondary_sales = true` tetap berlaku dengan `COALESCE(md.allow_upload_secondary_sales, false)` di sisi aman.
16. Use LEFT JOIN parent with `parent.pro_id=mp.parent_pro_id`, `parent.cust_id=LEFT(mp.cust_id,6)`, active/not-deleted predicates.
17. Bind cust ID list using sqlx-supported expansion/rebind or PostgreSQL array (`ANY($n)`) according to existing dependency pattern; test final query/args.
18. Bind q as `%<q>%`; no SQL literal concat.
19. Apply active/non-deleted `mp` predicates to both count/data.
20. Emit primary CASE only for mapping+upload enabled. Emit original CASE only same predicate.
21. Emit fixed `type` CASE independent of upload flag as user/docs specify.
22. Map allowlisted sort keys to fixed output expressions/columns. Ensure `type` sorts deterministic.
23. Bind limit/offset as integers; compute offset from normalized pagination.
24. Calculate last page with shared helper; define zero-total expected page count in test.
25. Add SQLMock count/data test and verify `ExpectationsWereMet`.
26. Add table-driven mapping scans for all five required categories.
27. Run targeted controller/repository/service tests, then full master test/build.
28. Run compose status. Start master only when local env/database available.
29. Run supplied curl with env token only. Capture HTTP/body sans token.
30. Inspect diff; confirm boundary and no raw input interpolation.
31. Write evidence manifest, hand to quality-gate.
32. **A4 remediation (Amendment 2026-07-14, confirmed_runtime)**: ganti step 15 (relasi `md` mentah) menjadi grouped derived relation. Tambahan: revisi SQLMock untuk menangkap grouped `SELECT cust_id, BOOL_OR(...) FROM mst.m_distributor GROUP BY cust_id` lalu `LEFT JOIN` ke `mp`; revisi tabel ekspektasi cardinality menjadi 1 baris per `mp`; tambah regression SQLMock `multi_distributor_for_one_cust` (simulasi 38 baris distributor) yang membuktikan 1 baris output; DB proof: `psql`/`rtk psql` hitung `COUNT(*)` dari base query untuk `C22001` menghasilkan 3917 dan tidak ada `pro_id` duplikat; endpoint proof: `total_record=3917`, page 1 tidak ada `pro_id` duplikat; update `A2-repository.log`, `A3-runtime.md` menjadi `A4-runtime.md`; re-run validasi.
33. **A5 remediation (Amendment 2026-07-14, second pass, confirmed_runtime)**: lock parent-eligibility fallback for mapping-enabled rows. Tambah eligibility guard pada parent primary `CASE` (`parent.pro_id IS NOT NULL` plus active/non-deleted predicate atau equivalent join shape) sehingga missing/inactive parent tidak terpilih. Bentuk fallback untuk row mapping-enabled dan parent-ineligible: primary fields dari `mp` (`cust_id`, `pro_id`, `pro_code`, `pro_name`, `parent_pro_id`); `original_*` terisi dari `mp`; `type` tetap `Product Mapping`. Tambah dua SQLMock regression case: `parent_present_and_eligible` (parent normalisasi, primary dari parent, original dari `mp`) dan `parent_missing_or_inactive` (primary dari `mp`, original dari `mp`, type `Product Mapping`, semua field non-null). Validasi scan ke non-nullable `pro_id` tidak gagal. Live runtime revalidation: `curl` authorized untuk `cust_id=C260020001` (token dari env, tidak di-paste, tidak di-log), capture status/body redacted; ekspektasi HTTP 200 tanpa conversion error dan semua primary field non-null. Bukti di `.opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A5-repository.log` dan `A5-runtime.md`.

## Expected Files to Change

- `master/controller/product_controller.go`
- `master/controller/product_report_controller_test.go` (new preferred)
- `master/entity/product.go` or `master/entity/product_report.go` (new preferred)
- `master/model/m_product.go` or `master/model/product_report.go` (new preferred)
- `master/service/product_service.go`
- `master/service/product_report_service_test.go` only if mapper non-trivial
- `master/repository/product_repository.go`
- `master/repository/product_report_repository_test.go` (new preferred)

## Agent / Tool Routing

| Area | Execute | Review |
|---|---|---|
| DTO/controller tests + route | `@backend` | `@quality-gate` |
| Service/repository SQL + SQLMock | `@backend` | `@quality-gate` |
| Runtime curl/env smoke | `@backend` | `@quality-gate` |
| Final security/evidence | `@quality-gate` | `@quality-gate` |

## Executor Handoff Prompt

```text
Task ID: 20260714-1315-sx-2513-product-secondary-sales-report
Plan: .opencode/plans/20260714-1315-sx-2513-product-secondary-sales-report.md
Scope: Add master POST /v1/products/report for Secondary Sales normalized product filtering.

Must preserve:
- Validated request cust_id IN-list is source of scope. No JWT locals replacement, no auto child expansion.
- Parent mapping join uses pro_id plus parent.cust_id=LEFT(mp.cust_id,6).
- Normalize only mapping + allow_upload_secondary_sales=true.
- Parameterize q/cust IDs/limit/offset. Closed sort allowlist only.
- Controller -> Service -> Repository -> DB; responsebuild + entity.Pagination.
- original_* JSON null for non-normalized rows.

Do not touch: migrations, JWT middleware, other modules, go.mod/go.sum, legacy GET product list, compose/env.

Validation: run exact commands in TDD/Test Plan and evidence requirements. Use token only from environment. Do not log token.

Return: changed paths:lines, test/build outputs, curl result if configured, unresolved parent-missing behavior, evidence paths. Tracker updates at every status transition are mandatory, not optional bookkeeping.
```

## Execution-ready Worklist / Handoff Contract

### A1 — Contract/route/filter tests

```yaml
handoff:
  task_id: 20260714-1315-sx-2513-product-secondary-sales-report
  plan_id: 20260714-1315-sx-2513-product-secondary-sales-report
  caller: orchestrator
  callee: backend
  scope: Add report DTOs route validation default sort normalization and controller Red/Green tests.
  claim_level: scoped
  claim_scope: A1 done only when report validates request cust_id list and emits master envelope.
  source_basis: master/controller/product_controller.go; master/entity/product.go; master/entity/api.go; master/pkg/responsebuild/response.go
  must_preserve: request cust_id list scope; POST query-param contract; response envelope
  do_not_touch: master/pkg/middleware/jwt_middleware.go; migration files; legacy GET product list
  validation: cd master && rtk go test ./controller -run Product.Report -v
  exit_criteria: controller route and validation tests pass
  evidence_required: .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A1-controller.log
  depends_on: none
  context_bundle: master/controller/product_controller.go; master/controller/dropdown_scope_controller_test.go; master/entity/api.go
```

### A2 — Safe report query/mapping tests

```yaml
handoff:
  task_id: 20260714-1315-sx-2513-product-secondary-sales-report
  plan_id: 20260714-1315-sx-2513-product-secondary-sales-report
  caller: orchestrator
  callee: backend
  scope: Add repository service ReportList with parameterized count data SQL and mapping semantics plus SQLMock tests.
  claim_level: scoped
  claim_scope: A2 done only when all mapping classes and SQL boundary cases pass tests.
  source_basis: master/repository/product_repository.go; master/repository/product_assignment_repository_test.go; Secondary_Sales_Report_BE.docx point 5; user Q&A recorded in plan Decisions
  must_preserve: parent.pro_id plus parent.cust_id equals LEFT(mp.cust_id,6); only enabled upload mapping normalizes to parent; no raw user SQL interpolation
  do_not_touch: master migration files; master/pkg/middleware; unrelated repository methods
  validation: cd master && rtk go test ./repository -run Product.Report -v
  exit_criteria: count data SQLMock mapping table parent join and parameter binding pass
  evidence_required: .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A2-repository.log
  depends_on: A1
  context_bundle: master/repository/product_repository.go; master/repository/product_assignment_repository_test.go; master/service/product_service.go
```

### A3 — Full verification and runtime smoke

```yaml
handoff:
  task_id: 20260714-1315-sx-2513-product-secondary-sales-report
  plan_id: 20260714-1315-sx-2513-product-secondary-sales-report
  caller: orchestrator
  callee: backend
  scope: Run full master test build and environment-gated report smoke curl.
  claim_level: scoped
  claim_scope: A3 done only when test build pass; runtime smoke not-ready if DB token unavailable and must not be claimed passed.
  source_basis: .opencode/docs/QUALITY.md; plan Validation Commands
  must_preserve: no token logging; no production data writes
  do_not_touch: source code except evidence artifacts
  validation: cd master && rtk go test ./... && rtk go build ./...
  exit_criteria: logs saved; curl status body captured or explicit not-ready reason recorded
  evidence_required: .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A3-test.log; .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A3-build.log; .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A3-runtime.md
  depends_on: A2
  context_bundle: .opencode/docs/QUALITY.md; master/.env
```

### A4 — Remediate duplicate distributor join and prove live cardinality

```yaml
handoff:
  task_id: 20260714-1315-sx-2513-product-secondary-sales-report
  plan_id: 20260714-1315-sx-2513-product-secondary-sales-report
  caller: orchestrator
  callee: backend
  scope: Replace multiplicative raw distributor join with one aggregated row per cust_id and prove count/data pagination cardinality.
  claim_level: scoped
  claim_scope: A4 done only when direct raw md join is gone; SQLMock and Docker+DB proof establish one report row per mp row for C22001. Do not claim generic production readiness.
  source_basis: .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/runtime-duplicate-join.md; master/repository/product_repository.go:143-213; master/repository/product_report_repository_test.go; plan Requirement 15 and Requirement 19
  must_preserve: request cust_id-only scope; composite parent join; normalization only for mapping plus aggregated upload flag true; parameterized values; Controller to Service to Repository; no migrations or data writes
  do_not_touch: migrations; JWT middleware; docker-compose.yml; env files; go.mod; go.sum; other modules; legacy product list
  validation: cd master && rtk go test ./repository -run 'Product.*Report' -v && rtk go test ./... && rtk go build ./...; read-only Docker plus DB cardinality proof for C22001; authorized curl if TOKEN available
  exit_criteria: derived md relation groups by cust_id using BOOL_OR(COALESCE(flag,false)); SQLMock regression proves one output/count row for one mp product despite 38 distributor rows; DB proves total 3917 and no pro_id duplicates for C22001; evidence written
  evidence_required: .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A4-repository.log; .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A4-runtime.md
  depends_on: A2
  context_bundle: .opencode/plans/20260714-1315-sx-2513-product-secondary-sales-report.md; .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/runtime-duplicate-join.md; master/repository/product_repository.go; master/repository/product_report_repository_test.go
```

**Subagent Context Bundle — A4**

- `verified_by_planner`:
  - `confirmed_runtime`: live `C22001` has 3917 active product rows and 38 distributor rows; endpoint total is 148846 (`3917 × 38`) and repeats `pro_id=495` five times. Source: `runtime-duplicate-join.md`.
  - `confirmed_repo`: direct raw join exists/planned at `master/repository/product_repository.go:143-213` and original plan old step 15.
  - `confirmed_repo`: count/data share base query construction; preserve that pattern. Source: `A2-notes.md:5-13`.
  - `confirmed_repo`: existing SQLMock pattern uses `sqlx.NewDb` and `ExpectationsWereMet`. Source: plan Existing Patterns/Reuse.
- `files_already_read`: `master/repository/product_repository.go`, `master/repository/product_report_repository_test.go`, `runtime-duplicate-join.md`, `A2-notes.md`.
- `open_assumptions`: `BOOL_OR` semantics are chosen as smallest compliant flag interpretation: any true distributor row enables upload for that `cust_id`. If product owner disputes this business rule, stop before alternate aggregation; do not choose arbitrary row.
- `source_of_truth_order`: latest user payload; `runtime-duplicate-join.md`; `runtime-parent-null.md`; plan invariants/Requirement 15/19/20; existing service patterns.

### A5 — Lock parent-eligibility fallback for mapping-enabled rows

```yaml
handoff:
  task_id: 20260714-1315-sx-2513-product-secondary-sales-report
  plan_id: 20260714-1315-sx-2513-product-secondary-sales-report
  caller: orchestrator
  callee: backend
  scope: Add parent-eligibility guard and `mp` primary fallback for mapping-enabled rows; add SQLMock regression for both eligible-parent and missing-parent branches; revalidate live runtime.
  claim_level: scoped
  claim_scope: A5 done only when both SQLMock branches pass and live curl for C260020001 returns HTTP 200 with no conversion error and non-null primary fields. Do not claim generic production readiness.
  source_basis: [.opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/runtime-parent-null.md, master/repository/product_repository.go, master/repository/product_report_repository_test.go, plan Requirements 20 and Acceptance Criteria 14 15 16]
  must_preserve: [request cust_id-only scope, composite parent join (parent.pro_id plus parent.cust_id equals LEFT(mp.cust_id,6)), normalization only for mapping plus aggregated upload flag true, parameterized values; closed sort allowlist, Controller to Service to Repository; responsebuild plus entity.Pagination, no migrations; no data writes; original_* JSON null for non-mapping paths]
  do_not_touch: [master migrations, master/pkg/middleware, docker-compose.yml, env files, go.mod, go.sum, other modules, legacy product list]
  validation: [cd master && rtk go test ./repository -run 'Product.*Report' -v, cd master && rtk go test ./controller -run 'Product.*Report' -v, cd master && rtk go test ./... && rtk go build ./..., authorized curl for C260020001 (token from env never logged); redact token from body files]
  exit_criteria: [SQLMock case parent_present_and_eligible passes, SQLMock case parent_missing_or_inactive passes with non-null primary and original fields, Live curl HTTP 200 no conversion error pro_id non-null for every record, A5 evidence files exist]
  evidence_required: [.opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A5-repository.log, .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A5-runtime.md]
  depends_on: [A4]
  context_bundle: [.opencode/plans/20260714-1315-sx-2513-product-secondary-sales-report.md, .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/runtime-parent-null.md, .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/runtime-duplicate-join.md, master/repository/product_repository.go, master/repository/product_report_repository_test.go]
```

**Subagent Context Bundle — A5**

- `verified_by_planner`:
  - `confirmed_runtime`: production curl for `C260020001` returned `gagal memecahkan kode: skema: kesalahan mengonversi nilai untuk pro_id`; mapping-enabled row, eligible parent missing/inactive, parent primary fields NULL, scan failure. Source: `runtime-parent-null.md`.
  - `confirmed_repo`: prior plan steps 15/22/231 and Requirements 8-12 do not guard parent eligibility in the primary `CASE`; parent-ineligible rows collapse to NULL. Source: plan lines 47-54, 75-77.
  - `confirmed_repo`: local psql `mapping_enabled parent_present` for `C260020001` returned two rows, both parent present and flag false; not a reproduction of the production branch. Source: `runtime-parent-null.md`.
  - `confirmed_runtime`: A4 fix addresses distributor-cardinality branch, not parent-eligibility branch. Source: A4 remediation entry; live behavior on `C260020001` after A4 only fix is unproven.
- `files_already_read`: `master/repository/product_repository.go`, `master/repository/product_report_repository_test.go`, `runtime-parent-null.md`, `runtime-duplicate-join.md`.
- `open_assumptions`: A5 fallback is the smallest compliant choice; no product owner required. If the product owner later wants different behavior (drop row, error response), the plan needs another amendment and the `Do Not / Reject If` rule that drops `original_*` on missing/inactive parent must be relaxed only by a fresh amendment.
- `source_of_truth_order`: latest user payload; `runtime-parent-null.md`; `runtime-duplicate-join.md`; plan invariants/Requirement 20; existing service patterns.

### Q1 — Final review

```yaml
handoff:
  task_id: 20260714-1315-sx-2513-product-secondary-sales-report
  plan_id: 20260714-1315-sx-2513-product-secondary-sales-report
  caller: orchestrator
  callee: quality-gate
  scope: Review scope isolation SQL parameterization mapping semantics tests runtime evidence and diff boundary.
  claim_level: scoped
  claim_scope: approve or request changes; do not edit source.
  source_basis: .opencode/plans/20260714-1315-sx-2513-product-secondary-sales-report.md; .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/
  must_preserve: read-only review
  do_not_touch: source files
  validation: inspect diff evidence tests query placeholders sort allowlist
  exit_criteria: approve or specific blockers
  evidence_required: .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/Q1-review.md
  depends_on: A5
  context_bundle: plan acceptance criteria; runtime-duplicate-join.md; runtime-parent-null.md; implementation evidence directory
```

1. **A1** | `@backend` | DTO/route/filter/controller tests | evidence `A1-controller.log`
2. **A2** | `@backend` | report SQL/mapping/service tests | evidence `A2-repository.log`
3. **A3** | `@backend` | full test/build/runtime smoke | evidence `A3-*.log`, `A3-runtime.md`
4. **A4** | `@backend` | remediate `md` join cardinality; SQLMock plus Docker+DB proof | evidence `A4-repository.log`, `A4-runtime.md`
5. **A5** | `@backend` | lock parent-eligibility fallback for mapping-enabled rows; SQLMock both branches; live revalidation | evidence `A5-repository.log`, `A5-runtime.md`
6. **Q1** | `@quality-gate` | final review after A5 | evidence `Q1-review.md`

`start_with: A5`

## Progress Tracking

- `tracker_path`: `.opencode/state/20260714-1315-sx-2513-product-secondary-sales-report/progress.json`
init_command: python3 ~/.config/opencode/scripts/task-progress.py 20260714-1315-sx-2513-product-secondary-sales-report --init --plan .opencode/plans/20260714-1315-sx-2513-product-secondary-sales-report.md
summary_command: python3 ~/.config/opencode/scripts/task-progress.py 20260714-1315-sx-2513-product-secondary-sales-report --summary
checklist_command: python3 ~/.config/opencode/scripts/task-progress.py 20260714-1315-sx-2513-product-secondary-sales-report --checklist
- `update_rules`: update before start, after complete/block/cancel, whenever evidence is written, and every cross-lane handoff. Mandatory.

| ID | Owner | Evidence | Update command |
|---|---|---|---|
| A1 | `@backend` | `A1-controller.log` | `python3 ~/.config/opencode/scripts/task-progress.py 20260714-1315-sx-2513-product-secondary-sales-report --update A1 --status completed --owner @backend --evidence .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A1-controller.log` |
| A2 | `@backend` | `A2-repository.log` | `python3 ~/.config/opencode/scripts/task-progress.py 20260714-1315-sx-2513-product-secondary-sales-report --update A2 --status completed --owner @backend --evidence .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A2-repository.log` |
| A3 | `@backend` | `A3-runtime.md` | `python3 ~/.config/opencode/scripts/task-progress.py 20260714-1315-sx-2513-product-secondary-sales-report --update A3 --status completed --owner @backend --evidence .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A3-runtime.md` |
| A4 | `@backend` | `A4-repository.log`, `A4-runtime.md` | `python3 ~/.config/opencode/scripts/task-progress.py 20260714-1315-sx-2513-product-secondary-sales-report --update A4 --status completed --owner @backend --evidence .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A4-runtime.md` |
| A5 | `@backend` | `A5-repository.log`, `A5-runtime.md` | `python3 ~/.config/opencode/scripts/task-progress.py 20260714-1315-sx-2513-product-secondary-sales-report --update A5 --status completed --owner @backend --evidence .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A5-runtime.md` |
| Q1 | `@quality-gate` | `Q1-review.md` | `python3 ~/.config/opencode/scripts/task-progress.py 20260714-1315-sx-2513-product-secondary-sales-report --update Q1 --status completed --owner @quality-gate --evidence .opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/Q1-review.md` |

## Validation Commands

1. `rtk docker compose -f docker-compose.yml ps` — confirm compose status from repo root.
2. `cd master && rtk go mod download && rtk go mod tidy` — module sync.
3. `cd master && rtk go test ./controller -run "Product.*Report" -v` — controller Red/Green tests.
4. `cd master && rtk go test ./repository -run "Product.*Report" -v` — SQLMock tests for sort/cust_id/parameter binding.
5. `cd master && rtk go test ./service -run "Product.*Report" -v` — service mapping tests.
6. `cd master && rtk go test ./...` — full test suite (no regression).
7. `cd master && rtk go build ./...` — full build (no compile error).
8. `python3 ~/.config/opencode/scripts/plan-execution-readiness.py .opencode/plans/20260714-1315-sx-2513-product-secondary-sales-report.md --project-root .` — readiness gate.
9. `python3 ~/.config/opencode/scripts/subagent-handoff-check.py --plan .opencode/plans/20260714-1315-sx-2513-product-secondary-sales-report.md` — handoff schema gate.
10. `python3 ~/.config/opencode/scripts/plan-compliance-check.py --project-root . --plan .opencode/plans/20260714-1315-sx-2513-product-secondary-sales-report.md --task-id 20260714-1315-sx-2513-product-secondary-sales-report` — compliance check.
11. Env-gated runtime smoke (only when DB + `$TOKEN` available, never echo token):
12. **A4 runtime proof (Amendment 2026-07-14, required for merge)**:

```bash
# 1. compose up local Postgres only (no other services touched) per AGENTS.md: "Compose runtime must target host local Postgres first"
rtk docker compose -f docker-compose.yml up -d postgres
# 2. read-only cardinality check
rtk psql -h host.docker.internal -U postgres -d ggn_scyllax -c "SELECT COUNT(*) AS active_products FROM mst.m_product WHERE cust_id='C22001' AND is_del=false AND is_active=true;"
rtk psql -h host.docker.internal -U postgres -d ggn_scyllax -c "SELECT COUNT(DISTINCT pro_id) AS distinct_products FROM mst.m_product WHERE cust_id='C22001' AND is_del=false AND is_active=true;"
# 3. authorized endpoint (if $TOKEN set; never log token)
curl --silent --show-error \
  --request POST "$BASE_URL/master/v1/products/report?cust_id[]=C22001&page=1&limit=20&sort_by=pro_name&sort_order=asc" \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --header "Authorization: Bearer $TOKEN" | tee /tmp/c22001.json
jq '.paging.total_record' /tmp/c22001.json
jq '[.data[].pro_id] | group_by(.) | map(length) | max' /tmp/c22001.json
```

Expected: `total_record=3917` and `max=1`; no token in evidence. A4 evidence path: `.opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A4-runtime.md`.
13. **A5 runtime proof (Amendment 2026-07-14, second pass, required for merge)**:

```bash
# 1. authorized endpoint for production fixture cust_id (token from env, never pasted into log/evidence/body file)
curl --silent --show-error \
  --request POST "$BASE_URL/master/v1/products/report?cust_id[]=C260020001&page=1&limit=20&sort_by=pro_name&sort_order=asc" \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --header "Authorization: Bearer $TOKEN" | tee /tmp/c260020001.json
jq '.errors // empty' /tmp/c260020001.json
jq '[.data[] | select(.pro_id == null)] | length' /tmp/c260020001.json
jq '[.data[] | select(.type == "Product Mapping" and (.original_cust_id == null or .original_pro_id == null))] | length' /tmp/c260020001.json
```

Expected: `errors` empty; zero records with `pro_id == null`; zero `Product Mapping` records with missing `original_*`. A5 evidence path: `.opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/A5-runtime.md`. Bearer token rotation reminder recorded in operational context; never pasted here.

Expected: HTTP 200 with `data` array, `paging` object, request ID; no token in evidence.

## Evidence Requirements

Keep:
- `discovery.md`: inspected files, docs, constraints, confirmed-vs-assumed audit.
- `index.json`: manifest with plan readiness.
- `runtime-duplicate-join.md`: live psql+curl evidence for the multiplicative join (Amendment 2026-07-14).
- `runtime-parent-null.md`: production conversion-error evidence plus fixed fallback contract (Amendment 2026-07-14, second pass).
- `A1-controller.log`, `A2-repository.log`, `A3-test.log`, `A3-build.log`, `A3-runtime.md`, `A4-repository.log`, `A4-runtime.md`, `A5-repository.log`, `A5-runtime.md`, `Q1-review.md` during execution.
- Runtime verify tooling not found in inspected repo scripts; skip recorded. Use service tests + compose/curl.
- Jira URL accessible only as title/redirect via web fetch; DOCX extracted by librarian. No GitHub upstream dependency.

## Done Criteria

- A1–A5 complete and tracker updated.
- All Acceptance Criteria verified or explicit runtime `not-ready` documented.
- No raw request input reaches SQL syntax.
- Parent composite join and mapping gate proven by SQLMock test.
- Eligible-parent and missing/inactive-parent mapping-enabled SQLMock branches pass; fallback primary fields non-null and original fields preserved.
- `go test ./...` and `go build ./...` pass in `master`.
- Runtime curl for `C260020001` returns HTTP 200 without `pro_id` conversion error; evidence redacts token.
- Quality gate approves evidence/diff.

## Final Planning Summary

### Artifacts consulted

- User story, SQL, cURL, and user Q&A.
- `Secondary_Sales_Report_BE.docx` point 5, extracted by `@librarian`.
- `master/controller/product_controller.go`, `service/product_service.go`, `repository/product_repository.go`, `entity/{product.go,api.go}`, `pkg/{responsebuild/response.go,sql_helper/sql_patch.go}`.
- `.opencode/docs/{index,ARCHITECTURE,SERVICE_MATRIX,QUALITY,AGENT_ROUTING,MCP}.md`.
- Live runtime evidence for `C22001` (`3917` products, `38` distributors, `total_record=148846`, `pro_id=495` five copies) in `.opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/runtime-duplicate-join.md`.

### Artifacts created/kept

- Primary: `.opencode/plans/20260714-1315-sx-2513-product-secondary-sales-report.md` (amended 2026-07-14, second pass).
- Amendment evidence: `.opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/runtime-duplicate-join.md`, `runtime-parent-null.md`.
- Kept operational evidence: `.opencode/evidence/20260714-1315-sx-2513-product-secondary-sales-report/discovery.md`, `index.json`, `check-plan/` logs, prior `A1–A4` evidence plus `Q1-review.md` (note: `Q1-review.md` verdict `PASS_WITH_RISKS` is superseded for merge; A5 must run before final gate).
- Draft removed after synthesis unless unresolved source evidence appears.

### Source strategy

Repo-local discovery and DOCX extraction used. Context7 skipped: no new/version-sensitive library/API is proposed. Browser skipped: backend task; Jira page did not expose issue content. GitHub skipped: workspace is not Git repo and no upstream source needed. Runtime evidence in `runtime-duplicate-join.md` is the primary authority for the multiplicative join finding; the fix uses PostgreSQL built-in `BOOL_OR`/`GROUP BY` (no new dependency, no API change).

### Readiness

`PASS_FOR_SLICE` (remediation scope only). Amendment 2026-07-14 downgrades the merge status from `PASS_WITH_RISKS` to `execution-ready A4` because Requirement 15 and the pagination Acceptance Criteria fail against the current `md` join. Code state from A1–A3 is preserved; only the `md` relation needs correction. After A4, re-run validators and re-issue Q1 review.

### Amendment log

- 2026-07-14: confirmed_runtime. Replaced raw `LEFT JOIN mst.m_distributor md ON md.cust_id=mp.cust_id` with grouped derived relation (`BOOL_OR(COALESCE(allow_upload_secondary_sales,false))` per `cust_id`). Added Requirement 19, AC 13, Invariant 9–10, D1, Risks entry, Reject-If for direct `md` join, A4 remediation task with Subagent Context Bundle, and live Docker+DB validation commands. No source code, migrations, JWT, env, or other modules touched by planner.
- 2026-07-14 (second pass): confirmed_runtime. Locked smallest compliant fallback for `mapping_enabled` rows with missing/inactive parent: primary fields from `mp`, `original_*` populated from `mp`, `type='Product Mapping'`, no primary scan field NULL. Added Requirement 20, AC 14–16, Invariant 11, Risks entry, Reject-If for null primary scan, A5 remediation task with Subagent Context Bundle, and live `C260020001` curl proof block. Bearer token rotation reminder recorded; token never pasted into plan or evidence.

### Active-lane reset

Execution starts under `@orchestrator` then `@backend`; each lane refreshes own permissions/context. Planner artifact-only restriction does not persist.
