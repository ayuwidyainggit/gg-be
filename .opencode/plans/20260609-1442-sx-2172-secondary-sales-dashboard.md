# Plan — SX-2172 Secondary Sales Dashboard

Task ID: `20260609-1442-sx-2172-secondary-sales-dashboard`
Readiness: `ready-for-implementation`
Quality Gate: `PASS_FOR_SLICE`
Primary source of truth: this file.

## Goal
Perbaiki `GET /v1/reports/secondary-sales/group` untuk issue SX-2172 dengan perubahan minimal pada modul `sales`, terutama query dashboard group agar label dan grouping `outlet`, `salesman`, `product_category`, dan `product` sesuai expected Jira.

## Non-goals
- Tidak merombak struktur endpoint, payload, atau response JSON.
- Tidak mengubah controller route.
- Tidak mengubah proses extract/report fact selain jika test membuktikan bug berasal dari data extraction.
- Tidak melakukan migrasi DB.
- Tidak memperbaiki LSP error existing yang tidak terkait langsung dengan SX-2172, kecuali menghalangi test target.

## Scope
Target utama:
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`

Target opsional hanya jika dibutuhkan oleh mapping:
- `sales/service/report_service.go`
- `sales/service/report_service_test.go`
- `sales/model/report.go`
- `sales/entity/report.go`
- `sales/client_test.http`

## Requirements
- `group_by=outlet` / Sales by Customer:
  - Mapping yang dipilih user: `id = salesman_id`, `code = outlet_code`, `name = salesman_name + " > " + outlet_name`.
  - Query harus group by kombinasi salesman dan outlet agar `net_sales` tidak mencampur salesman berbeda pada outlet yang sama.
- `group_by=salesman` / Sales by Salesman:
  - `id = salesman_id`, `code = salesman_code`.
  - `name` boleh tetap `salesman_name` agar backward compatible.
- `group_by=product_category` / Sales by Category:
  - Master table harus jadi sumber prioritas: `mst.m_product` via `pro_id + cust_id`, lalu `mst.m_product_cat` via `pcat_id`.
  - Fallback ke report dim hanya ketika master kosong/tidak valid.
  - Handle `report.dim_products.category_id` kosong, `0`, atau null tanpa kehilangan kategori master.
- `group_by=product` / Sales by Product:
  - Master table `mst.m_product` harus jadi sumber prioritas untuk product id/code/name.
  - Fallback ke `report.dim_products` hanya ketika master kosong/tidak valid.
  - Product name tidak boleh kosong jika master/dim punya data.
- `fact_returns` tetap mengurangi `net_sales` dengan `fr.net_sales_exclude_ppn * -1`.
- Filter `cust_id IN ?`, `dt.month = ?`, dan `dt."year" = ?` tetap ada untuk order dan return branch.

## Acceptance Criteria
- `Sales by Customer` menghasilkan field yang mendukung tampilan `Salesman ID > Outlet Code - Salesman Name > Outlet Name` melalui `id=salesman_id`, `code=outlet_code`, `name=salesman_name + " > " + outlet_name`.
- `Sales by Salesman` menghasilkan `id=salesman_id`, `code=salesman_code`, `name=salesman_name`.
- `Sales by Category` memakai master category dari `mst.m_product` + `mst.m_product_cat` ketika report dim category kosong/0/null.
- `Sales by Product` menghasilkan `id=product_id`, `code=product_code`, `name=product_name` dari master lebih dulu.
- Output response tetap memakai `id`, `code`, `name`, `net_sales`.
- Branch outlet dan salesman tidak kehilangan return subtraction, month/year filter, dan multi-cust `IN ?` binding.
- `rtk go test ./repository -run 'TestSecondarySalesReportGroup'` lulus dari direktori `sales`.
- `rtk go test ./service -run 'TestSecondarySalesReportGroupSales'` lulus dari direktori `sales`.
- Jika memungkinkan, `rtk go test ./...` lulus dari direktori `sales`, atau kegagalan unrelated dicatat dengan bukti.

## Existing Patterns/Reuse
- Reuse `buildSecondarySalesReportGroupQuery(groupBy string)` sebagai pusat perubahan SQL.
- Reuse raw SQL + `repository.Raw(query, custIDs, month, year, custIDs, month, year).Find(&results)`.
- Reuse model `model.SecondarySalesReportGroup` karena alias SQL masih bisa tetap `id`, `code`, `name`, `net_sales`.
- Reuse response `entity.SecondarySalesReportGroupResp`; tidak perlu contract baru.
- Reuse dry-run SQL test helpers di `sales/repository/report_repository_test.go`: `newReportRepoDryRunDB`, `latestRecordedQuery`.
- Reuse service branch tests di `sales/service/report_service_test.go` untuk memastikan branch selection tetap aman.

## Constraints
- Repo-local rule: validasi dari target service directory `sales`.
- Shell command repo ini memakai prefix `rtk`.
- Tenant rule: jangan hilangkan `cust_id`; master product join untuk dashboard group harus memakai row-level fact cust: `mp.cust_id = fo.cust_id` dan `mp.cust_id = fr.cust_id`.
- Schema prefixes penting: `report.`, `mst.`.
- Response field terbatas pada `id`, `code`, `name`, `net_sales`.
- Jangan copy atau menambah secret/env.

## Risks
- `mst.m_product_cat` mungkin tidak punya kolom `cust_id`; ikuti pola existing extract yang join `mst.m_product_cat prdcat ON prdcat.pcat_id = COALESCE(pp.pcat_id, cp.pcat_id)` kecuali schema lokal membuktikan sebaliknya.
- Mengubah outlet grouping dari outlet-only menjadi salesman+outlet dapat mengubah jumlah row. Ini disengaja berdasarkan jawaban user untuk mapping `ID=salesman, code=outlet`.
- Jika ada transaksi lama tanpa master product dan tanpa dim product valid, label bisa tetap kosong; query harus tetap mengembalikan row dengan fallback `id=fo.pro_id/fr.pro_id`, `code=''`, `name=''` atau label aman, bukan drop row karena inner join.
- Current LSP diagnostics menunjukkan error existing pada `queryActivitySalesReportRows` argument mismatch di file yang sama; executor harus bedakan jika test gagal karena isu unrelated.

## Decisions/Assumptions
- Keputusan user: untuk `Sales by Customer`, gunakan `id=salesman_id`, `code=outlet_code`, `name=salesman_name + " > " + outlet_name`.
- Master product/category adalah authoritative source untuk `product` dan `product_category`.
- Fallback dim report tetap diperlukan agar data lama tidak hilang ketika master tidak lengkap.
- Tidak perlu perubahan controller/entity jika SQL alias tetap sama.
- Source strategy: repo-local evidence + konteks Jira dari prompt. Official docs, GitHub, dan web search diskip karena bug adalah query SQL lokal dan tidak bergantung library/API eksternal.

## Execution Source of Truth
Urutan prioritas executor:
1. Instruksi eksplisit terbaru dari user.
2. Safety/security/tenant rules repo.
3. Non-negotiable Implementation Invariants di plan ini.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria dan Done Criteria.
6. Implementation Steps.
7. Rekomendasi/follow-up non-blocking.

Jika ada konflik, ikuti sumber dengan prioritas lebih tinggi dan catat konflik di evidence verifikasi.

## Non-negotiable Implementation Invariants
- Planner artifact-only: jangan menganggap plan ini sudah mengubah source.
- Controller → Service → Repository → DB harus tetap dipertahankan.
- Query repository tidak boleh menghapus `cust_id IN ?`, `month`, atau quoted year filter.
- `fact_returns` harus tetap dikurangi dari net sales.
- Master product/category harus prioritas untuk `product` dan `product_category`.
- Fallback dim harus mempertahankan row lama; jangan ubah `LEFT JOIN` master/dim menjadi `JOIN` yang bisa menghilangkan transaksi.
- Alias final harus tetap `id`, `code`, `name`, `net_sales`.
- Jangan mengubah package, lockfile, env, migration, atau unrelated report endpoints.

## Do Not / Reject If
- Reject jika `product_category` masih hanya join `report.dim_products` → `report.dim_product_categories` tanpa master fallback.
- Reject jika `product` masih hanya join `report.dim_products`.
- Reject jika join master memakai hardcoded cust id atau auth cust tunggal untuk multi-cust fact rows.
- Reject jika `fact_returns` tidak lagi dikurangi.
- Reject jika outlet/salesman branch kehilangan `dt."year"` filter.
- Reject jika response contract berubah menambah field wajib baru tanpa user approval.
- Reject jika source changes melebar ke migration/config/env.

## Diff Boundary
Allowed source changes:
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`
- `sales/service/report_service.go` hanya jika mapping tidak bisa selesai via SQL alias.
- `sales/service/report_service_test.go` hanya untuk branch/mapping regression.
- `sales/model/report.go` dan `sales/entity/report.go` hanya jika sangat diperlukan; expected tidak memerlukan ini.
- `sales/client_test.http` opsional untuk manual request examples.

Allowed planning/evidence changes:
- `.opencode/plans/20260609-1442-sx-2172-secondary-sales-dashboard.md`
- `.opencode/evidence/20260609-1442-sx-2172-secondary-sales-dashboard/**`

Any out-of-boundary change must be reverted or justified in verification evidence before final quality gate.

## TDD/Test Plan
TDD required: yes, karena ini bugfix query/API behavior.

Existing test patterns:
- Repository dry-run SQL tests in `sales/repository/report_repository_test.go`.
- Service mock branch tests in `sales/service/report_service_test.go`.

Red step:
- Tambah/ubah repository tests yang gagal pada query existing:
  - `TestSecondarySalesReportGroupProductCategoryUsesMasterCategoryFallback`
  - `TestSecondarySalesReportGroupProductUsesMasterProductFallback`
  - `TestSecondarySalesReportGroupOutletUsesSalesmanOutletDisplayMapping`
- Test harus assert SQL contains fragments untuk `mst.m_product`, `mst.m_product_cat`, `COALESCE`, `NULLIF`, row-level `fo.cust_id`/`fr.cust_id`, dan alias final `id`, `code`, `name`.

Green step:
- Update `buildSecondarySalesReportGroupQuery` maps untuk `outlet`, `product_category`, dan `product`.
- Pertahankan params order: `custIDs, month, year, custIDs, month, year`.

Refactor step:
- Jika map SQL menjadi terlalu sulit dibaca, ekstrak helper kecil untuk select/join fragments, tetapi jangan ubah behavior.
- Hindari restructure besar.

Edge cases:
- `dim_products.category_id = 0` atau null.
- Master product ada, dim product kosong.
- Master product kosong, dim product ada.
- Product name master kosong string; fallback ke dim name.
- Return rows punya product/category sama dan harus subtract.
- Multi-cust request dengan `custIDs` lebih dari satu.

Commands:
```bash
rtk go test ./repository -run 'TestSecondarySalesReportGroup'
rtk go test ./service -run 'TestSecondarySalesReportGroupSales'
rtk go test ./...
```

## Implementation Steps
1. Tambahkan regression tests repository untuk group SQL.
2. Update outlet branch di `buildSecondarySalesReportGroupQuery`:
   - order select: `fo.salesman_id AS id`, outlet code sebagai `code`, concatenated salesman/outlet label sebagai `name`, net sales order.
   - return select: `fr.salesman_id AS id`, outlet code sebagai `code`, concatenated salesman/outlet label sebagai `name`, return net sales negatif.
   - join outlet dim dan salesman dim pada order/return branch.
   - group final by `id`, `code`, `name`.
3. Pastikan salesman branch tetap `dsls.id`, `dsls.code`, `dsls.name`.
4. Update product category branch:
   - Use `LEFT JOIN report.dim_products dprd ON fo.pro_id = dprd.id`.
   - Use `LEFT JOIN report.dim_product_categories dprdctr ON NULLIF(dprd.category_id, 0) = dprdctr.id`.
   - Use `LEFT JOIN mst.m_product mp ON mp.pro_id = fo.pro_id AND mp.cust_id = fo.cust_id`.
   - Use `LEFT JOIN mst.m_product_cat mpc ON mpc.pcat_id = mp.pcat_id` unless schema requires `cust_id` filter.
   - Select `COALESCE(NULLIF(mpc.pcat_id, 0), NULLIF(dprdctr.id, 0), 0) AS id`, `COALESCE(NULLIF(mpc.pcat_code, ''), dprdctr.code, '') AS code`, `COALESCE(NULLIF(mpc.pcat_name, ''), dprdctr.name, '') AS name`.
   - Mirror same logic for return branch with `fr`.
5. Update product branch:
   - Use `LEFT JOIN mst.m_product mp ON mp.pro_id = fo.pro_id AND mp.cust_id = fo.cust_id` plus `LEFT JOIN report.dim_products dprd ON fo.pro_id = dprd.id`.
   - Select `COALESCE(NULLIF(mp.pro_id, 0), NULLIF(dprd.id, 0), fo.pro_id) AS id`, `COALESCE(NULLIF(mp.pro_code, ''), dprd.code, '') AS code`, `COALESCE(NULLIF(mp.pro_name, ''), dprd.name, '') AS name`.
   - Mirror same logic for return branch with `fr`.
6. Run targeted tests.
7. If targeted tests pass, run `rtk go test ./...`; record unrelated failures if existing LSP/activity-report mismatch blocks all tests.
8. Optional manual smoke with runtime DB if available:
   - `rtk docker compose -f docker-compose.yml ps` from repo root.
   - Hit four sample URLs and compare JSON fields.

## Expected Files to Change
Likely:
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`

Possibly:
- `sales/service/report_service_test.go`
- `sales/client_test.http`

Not expected:
- `sales/entity/report.go`
- `sales/model/report.go`
- `sales/service/report_service.go`

## Agent/Tool Routing
- `@orchestrator`: coordinate implementation and evidence.
- `@fixer`: source/test edits for bounded backend bugfix.
- `@quality-gate`: final signoff because this touches report SQL and tenant-sensitive data.
- `@explorer`: only if executor needs more schema/pattern discovery.
- `@librarian`: not needed unless external DB/library behavior unexpectedly becomes central.

## Executor Handoff Prompt
Copyable prompt:

```text
Implement SX-2172 in the `sales` module using `.opencode/plans/20260609-1442-sx-2172-secondary-sales-dashboard.md` as source of truth. Scope: fix `GET /v1/reports/secondary-sales/group` group queries for outlet/customer, salesman, product_category, and product. Must preserve response fields `id`, `code`, `name`, `net_sales`, tenant filters, `dt."year"`, and return subtraction. Use TDD: add failing repository dry-run SQL tests first, then update `buildSecondarySalesReportGroupQuery`. Do not touch env, migrations, package files, or unrelated endpoints. Validate with `rtk go test ./repository -run 'TestSecondarySalesReportGroup'`, `rtk go test ./service -run 'TestSecondarySalesReportGroupSales'`, then `rtk go test ./...` from `sales`. Return changed files, commands run, SQL behavior summary, and any unrelated failures with evidence.
```

## Execution-ready Worklist / Handoff Contract
`start_with`: `T1`

### T1 — Add regression SQL tests
- `depends_on`: none
- `owner/lane`: `@fixer`
- `action`: Add repository dry-run tests for product category master fallback, product master fallback, and outlet salesman+outlet display mapping.
- `validation`: `rtk go test ./repository -run 'TestSecondarySalesReportGroup'`
- `exit_criteria`: Tests fail before SQL change or clearly assert missing fragments in existing query.
- `blocking_status`: ready
- `requires_user_decision`: no
- `must_preserve`: existing quoted year and return subtraction tests.
- `do_not_touch`: controller, env, package files.
- `evidence_update`: record test names and expected SQL fragments.
- `exit_verification`: show failing output or explain if same commit includes green result.

### T2 — Update group query SQL
- `depends_on`: T1
- `owner/lane`: `@fixer`
- `action`: Modify `buildSecondarySalesReportGroupQuery` select/join maps to implement outlet mapping and master-priority product/category fallback.
- `validation`: `rtk go test ./repository -run 'TestSecondarySalesReportGroup'`
- `exit_criteria`: Repository group tests pass and params remain `custIDs, month, year, custIDs, month, year`.
- `blocking_status`: ready
- `requires_user_decision`: no
- `must_preserve`: response aliases and row-level `cust_id` joins.
- `do_not_touch`: service/controller unless required by tests.
- `evidence_update`: capture relevant SQL fragments from tests.
- `exit_verification`: targeted repository test output.

### T3 — Verify service mapping
- `depends_on`: T2
- `owner/lane`: `@fixer`
- `action`: Run service branch tests; update only if mapping/branch test needs explicit outlet expectation.
- `validation`: `rtk go test ./service -run 'TestSecondarySalesReportGroupSales'`
- `exit_criteria`: Service branch test passes and maps `Code` non-empty for all branches.
- `blocking_status`: ready
- `requires_user_decision`: no
- `must_preserve`: service branch selection and fallback-to-product behavior.
- `do_not_touch`: repository SQL beyond T2 unless service test reveals mismatch.
- `evidence_update`: record service test output.
- `exit_verification`: targeted service test output.

### T4 — Full module validation and manual smoke
- `depends_on`: T3
- `owner/lane`: `@fixer`
- `action`: Run broader tests and optional endpoint smoke when runtime DB is available.
- `validation`: `rtk go test ./...`; optional four sample GET requests.
- `exit_criteria`: Full tests pass or unrelated failures are documented with exact output and reason.
- `blocking_status`: ready
- `requires_user_decision`: no
- `must_preserve`: no remote DB defaults; use local compose posture if runtime smoke is attempted.
- `do_not_touch`: compose/env secrets.
- `evidence_update`: command outputs and smoke JSON samples if run.
- `exit_verification`: test output summary.

### T5 — Quality gate
- `depends_on`: T4
- `owner/lane`: `@quality-gate`
- `action`: Review diff boundary, SQL tenant safety, tests, and acceptance criteria.
- `validation`: inspect diff and evidence.
- `exit_criteria`: PASS or explicit blockers.
- `blocking_status`: ready
- `requires_user_decision`: no
- `must_preserve`: all invariants in plan.
- `do_not_touch`: source edits by reviewer.
- `evidence_update`: final signoff summary.
- `exit_verification`: quality gate result.

## Validation Commands
From repo root for runtime status if needed:
```bash
rtk docker compose -f docker-compose.yml ps
```

From `sales` directory:
```bash
rtk go test ./repository -run 'TestSecondarySalesReportGroup'
rtk go test ./service -run 'TestSecondarySalesReportGroupSales'
rtk go test ./...
```

Optional smoke requests after service is running:
```http
GET /sales/v1/reports/secondary-sales/group?month=4&cust_id=C260020001&group_by=outlet
GET /sales/v1/reports/secondary-sales/group?month=4&cust_id=C260020001&group_by=salesman
GET /sales/v1/reports/secondary-sales/group?month=4&cust_id=C260020001&group_by=product_category
GET /sales/v1/reports/secondary-sales/group?month=4&cust_id=C260020001&group_by=product
```

## Evidence Requirements
Implementation evidence must include:
- Changed files list.
- Test commands and outputs.
- SQL fragment summary proving master-priority fallback.
- Note whether manual runtime smoke was run; if skipped, state why.
- Note any unrelated existing failures, especially current LSP diagnostics around `queryActivitySalesReportRows` if they affect `rtk go test ./...`.

Research gate decision:
- Local project discovery: used and sufficient.
- Official docs/context7: skipped; no version-sensitive library/API behavior.
- GitHub: skipped; upstream behavior not needed.
- Web search: skipped; Jira context provided expected behavior.
- Browser/screenshot: not applicable for backend API bugfix.

## Done Criteria
- SQL uses master product/category fallback for `product_category` and `product`.
- Outlet mapping follows user decision.
- Targeted repository and service tests pass.
- Full `sales` module test status is known and documented.
- No out-of-boundary changes remain.
- Final quality gate has enough evidence to assess tenant/query safety.

## Final Planning Summary
Artifacts created:
- `.opencode/plans/20260609-1442-sx-2172-secondary-sales-dashboard.md`
- `.opencode/evidence/20260609-1442-sx-2172-secondary-sales-dashboard/discovery.md`
- `.opencode/evidence/20260609-1442-sx-2172-secondary-sales-dashboard/index.json`

Artifacts kept:
- Evidence directory kept because it contains discovery notes and source strategy useful for executor replay.

Key decisions:
- `product` and `product_category` must prioritize master tables and fallback to report dim.
- `outlet` display mapping uses user-selected `id=salesman_id`, `code=outlet_code`, `name=salesman_name + " > " + outlet_name`.
- No model/entity response expansion planned.

Questions:
- Question gate was asked and answered for outlet mapping.
- No remaining blocking open questions.

Readiness:
- `ready-for-implementation` for bounded backend fix.
- `PASS_FOR_SLICE` because runtime DB smoke may depend on local data availability, but code/test work is executable without replanning.

Cleanup performed:
- No stale drafts were created, so no draft cleanup needed.
