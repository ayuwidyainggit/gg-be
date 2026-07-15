# Plan — SX-2224 API Secondary Sales Group

Task ID: `20260612-0916-sx-2224-secondary-sales-group-code`
Readiness: `ready-for-implementation`
Quality Gate: `PASS`
Primary source of truth: this file.

## Goal
Implementasi final SX-2224 untuk `GET /sales/v1/reports/secondary-sales/group`: pastikan 4 variasi `group_by` mengembalikan `code`, `product` dan `product_category` memakai nama master yang tidak kosong ketika data master tersedia, dan branch `outlet` mengikuti dokumen SX-2224 sebagai grouping outlet murni.

## Non-goals
- Tidak mengubah route, auth middleware, envelope response, atau endpoint report lain.
- Tidak membuat migrasi DB atau mengubah schema.
- Tidak mengganti sumber data dashboard ke `report.fact_orders`; query saat ini sudah memakai source tables `sls.order` dan `sls.return`.
- Tidak membuat fallback sintetis `Product <id>` kecuali nanti ada keputusan user baru; DB lokal membuktikan product `10733` punya `pro_name` master.
- Tidak menyentuh module `pjp-sales` kecuali user eksplisit meminta sinkronisasi terpisah.

## Scope
Target utama:
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`

Target opsional hanya bila test memperlihatkan perlu:
- `sales/service/report_service_test.go`
- `sales/client_test.http`

Tidak diharapkan berubah:
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/model/report.go`
- `sales/entity/report.go`

## Requirements
- `group_by=outlet` harus mengikuti SX-2224:
  - `id = outlet_id`
  - `code = mst.m_outlet.outlet_code`
  - `name = mst.m_outlet.outlet_name`
  - Grouping tidak boleh lagi memakai `salesman_id` atau `emp_name > outlet_name`.
- `group_by=salesman` tetap:
  - `id = salesman_id`
  - `code = mst.m_employee.emp_code`
  - `name = mst.m_employee.emp_name`
- `group_by=product_category` tetap memakai kategori produk, bukan salesman:
  - `id = mst.m_product_cat.pcat_id` fallback `0`
  - `code = mst.m_product_cat.pcat_code`
  - `name = mst.m_product_cat.pcat_name`
- `group_by=product` tetap memakai master product:
  - `id = mst.m_product.pro_id` fallback source row product id
  - `code = mst.m_product.pro_code`
  - `name = mst.m_product.pro_name`
- Semua branch harus mempertahankan `cust_id IN ?`, date range dari `month/year`, return subtraction, dan `ORDER BY net_sales DESC`.
- Response tetap `{ message, data, request_id }`, item data tetap `id`, `code`, `name`, `net_sales`.

## Acceptance Criteria
- `outlet` response berisi `id` outlet, `code` outlet, `name` outlet tanpa nama salesman terkonkatenasi.
- `salesman` response berisi `code` salesman dan nama salesman.
- `product_category` response berisi `code`/`name` kategori dari `mst.m_product_cat` saat master tersedia.
- `product` response berisi `code`/`name` product dari `mst.m_product` saat master tersedia.
- Product `10733` pada DB lokal `ggn_scyllax` menghasilkan sumber master `code=JY1-002`, `name=Jersey Manchester City FC`; jangan pakai placeholder `Product 10733` untuk data ini.
- Filter `month/year/cust_id` konsisten untuk order dan return.
- Authorization `cust_id` existing tetap aman; principal boleh child scope valid, distributor hanya auth cust id.
- Targeted tests lulus dari direktori `sales`.

## Existing Patterns/Reuse
- Reuse `buildSecondarySalesReportGroupQuery(groupBy string)` sebagai pusat perubahan SQL.
- Reuse model `model.SecondarySalesReportGroup` dan DTO `entity.SecondarySalesReportGroupResp`; keduanya sudah punya `Code`.
- Reuse service mapping `Code: r.Code`; tidak perlu ubah service jika SQL alias tetap `code`.
- Reuse dry-run SQL helpers di `sales/repository/report_repository_test.go`: `newReportRepoDryRunDB`, `latestRecordedQuery`, `assertSecondarySalesGroupDateVars`.
- Reuse existing tests untuk product/category master fallback; ubah test outlet yang sekarang masih mencerminkan SX-2172.

## Constraints
- Jalankan command dari service dir `sales` dan pakai prefix `rtk`.
- Tenant join master harus row-level: `mo.cust_id = o.cust_id` / `rd.cust_id`, `mp.cust_id = od.cust_id` / `rd.cust_id`, employee via salesman customer context.
- `mst.m_product_cat` terkonfirmasi punya `pcat_id`, `pcat_code`, `pcat_name`; tidak terlihat butuh `cust_id` filter pada category.
- Compose saat discovery tidak running; DB lokal bisa diakses langsung via Postgres host lokal.
- Repo sudah punya tracked infra credentials; jangan copy, memperluas, atau commit secret/env.

## Risks
- Mengubah `outlet` dari SX-2172 mapping ke SX-2224 mapping dapat mengubah jumlah row dan interpretasi UI dari “Sales by Customer” menjadi outlet murni. Ini disengaja berdasarkan jawaban user.
- Test lama `TestSecondarySalesReportGroupOutletUsesSalesmanOutletDisplayMapping` akan gagal dan harus diganti/diubah agar sesuai SX-2224.
- Jika transaksi punya product/category tanpa master, `name` masih bisa kosong; requirement “tidak kosong” hanya dapat dijamin ketika master tersedia. Jangan ubah `LEFT JOIN` menjadi `JOIN` karena akan membuang transaksi lama.
- PJP mirror module memiliki file serupa, tapi scope task adalah `sales`; perubahan ke `pjp-sales` berisiko tanpa instruksi eksplisit.

## Decisions/Assumptions
- Keputusan user: branch `outlet` harus mengikuti dokumen SX-2224, bukan mapping SX-2172.
- Hasil cek DB: `mst.m_product.pro_id=10733` punya `pro_code=JY1-002`, `pro_name=Jersey Manchester City FC`; contoh `Product 10733` dianggap placeholder, bukan requirement fallback sintetis.
- Source of truth untuk product/category adalah master `mst.m_product` dan `mst.m_product_cat`.
- `report.dim_*` tidak diperlukan untuk slice ini karena query existing sudah memakai source transaction tables + master joins dan tests sudah mengunci pola itu.
- Official docs/GitHub/web diskip karena tidak ada perilaku library eksternal yang material.

## Execution Source of Truth
Prioritas executor:
1. Instruksi eksplisit terbaru dari user.
2. Safety/security/tenant rules repo.
3. Non-negotiable Implementation Invariants di plan ini.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria dan Done Criteria.
6. Implementation Steps.
7. Follow-up non-blocking.

Jika konflik, ikuti prioritas lebih tinggi dan catat konflik di evidence verifikasi.

## Non-negotiable Implementation Invariants
- Planner artifact-only; plan ini belum mengubah source.
- Kontrak Controller → Service → Repository → DB harus tetap dipertahankan.
- Final SQL alias wajib tetap `id`, `code`, `name`, `net_sales`.
- `outlet` wajib outlet murni: `o.outlet_id`/`r.outlet_id`, `mo.outlet_code`, `mo.outlet_name`.
- `product` dan `product_category` wajib `LEFT JOIN` ke master agar transaksi tidak hilang saat master kosong.
- Date filter wajib tetap range dari `month/year`; jangan downgrade ke `month` saja.
- `cust_id IN ?` dan row-level master joins tidak boleh hilang.
- Return rows wajib subtract net sales dengan nilai negatif.
- Jangan menyentuh env, migrations, package files, lockfiles, atau endpoint unrelated.

## Do Not / Reject If
- Reject jika branch `outlet` masih memakai `salesman_id` sebagai `id` atau `CONCAT_WS(...emp_name..., ...outlet_name...)` sebagai `name`.
- Reject jika branch `product_category` mengambil code/name salesman.
- Reject jika `product` memakai `Product <id>` untuk data yang punya `pro_name` master.
- Reject jika `LEFT JOIN mst.m_product` / `mst.m_product_cat` diganti menjadi `JOIN` tanpa bukti aman.
- Reject jika filter `year` hilang atau menjadi default current year tanpa param request.
- Reject jika authorization `cust_id` dilewati atau query menerima cust id arbitrary di service.
- Reject jika response envelope/field JSON berubah di luar penambahan/keberadaan `code` yang sudah ada.

## Diff Boundary
Allowed source changes:
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`
- `sales/service/report_service_test.go` hanya bila perlu update expectation branch/service.
- `sales/client_test.http` opsional untuk manual smoke examples.

Allowed evidence changes:
- `.opencode/evidence/20260612-0916-sx-2224-secondary-sales-group-code/**`
- `.opencode/plans/20260612-0916-sx-2224-secondary-sales-group-code.md`

Any out-of-boundary change must be reverted or justified in verification evidence before quality gate.

## TDD/Test Plan
TDD required: yes, karena ini API/query behavior yang tenant-sensitive.

Existing test patterns:
- Repository dry-run SQL tests in `sales/repository/report_repository_test.go`.
- Service mock branch tests in `sales/service/report_service_test.go`.

Red step:
- Ubah/rename `TestSecondarySalesReportGroupOutletUsesSalesmanOutletDisplayMapping` menjadi expectation SX-2224, misalnya `TestSecondarySalesReportGroupOutletUsesOutletDisplayMapping`.
- Assert SQL contains:
  - `o.outlet_id AS id`
  - `COALESCE(mo.outlet_code, '') AS code`
  - `COALESCE(mo.outlet_name, '') AS name`
  - `LEFT JOIN mst.m_outlet mo ON mo.outlet_id = o.outlet_id AND mo.cust_id = o.cust_id`
  - `r.outlet_id AS id`
  - `LEFT JOIN mst.m_outlet mo ON mo.outlet_id = r.outlet_id AND mo.cust_id = rd.cust_id`
  - `GROUP BY id, code, name`
- Assert SQL does not contain outlet branch `CONCAT_WS(' > '` or `o.salesman_id AS id` / `r.salesman_id AS id`.

Green step:
- Update `buildSecondarySalesReportGroupQuery` outlet select/join fragments only.
- Keep salesman/product/product_category fragments unless tests expose mismatch.

Refactor step:
- Jika SQL map tetap jelas, jangan refactor besar.
- Bila perlu, rename test only; jangan ubah API contract.

Edge cases:
- Data kosong untuk month/year tetap return empty data.
- Multi-cust principal request tetap memakai `cust_id IN ?`.
- Distributor request ke cust_id lain tetap unauthorized.
- Product/category master kosong tidak boleh drop rows.
- Return rows untuk outlet/product/category sama harus subtract dari net sales.

Commands:
```bash
rtk go test ./repository -run 'TestSecondarySalesReportGroup'
rtk go test ./service -run 'TestSecondarySalesReportGroupSales'
rtk go test ./controller -run 'TestSecondaryReportSalesGroup'
rtk go test ./...
```

Optional DB smoke if runtime service is available:
```bash
rtk docker compose -f docker-compose.yml ps
curl '<base>/sales/v1/reports/secondary-sales/group?month=6&year=2026&cust_id=C260020001&group_by=outlet' -H 'Authorization: Bearer <token>'
curl '<base>/sales/v1/reports/secondary-sales/group?month=6&year=2026&cust_id=C260020001&group_by=salesman' -H 'Authorization: Bearer <token>'
curl '<base>/sales/v1/reports/secondary-sales/group?month=6&year=2026&cust_id=C260020001&group_by=product_category' -H 'Authorization: Bearer <token>'
curl '<base>/sales/v1/reports/secondary-sales/group?month=6&year=2026&cust_id=C260020001&group_by=product' -H 'Authorization: Bearer <token>'
```

## Implementation Steps
1. Update repository dry-run outlet test to SX-2224 outlet murni expectation.
2. Run targeted repository test to confirm Red if source still uses SX-2172 mapping.
3. Edit `buildSecondarySalesReportGroupQuery` outlet branch:
   - Order select:
     - `o.outlet_id AS id`
     - `COALESCE(mo.outlet_code, '') AS code`
     - `COALESCE(mo.outlet_name, '') AS name`
   - Return select:
     - `r.outlet_id AS id`
     - `COALESCE(mo.outlet_code, '') AS code`
     - `COALESCE(mo.outlet_name, '') AS name`
   - Keep `LEFT JOIN mst.m_outlet` for order and return.
   - Remove salesman/employee joins from outlet branch unless another test proves they are still needed elsewhere.
4. Run repository group tests.
5. Run service/controller targeted tests; update only test expectations if existing tests assert old outlet id/name.
6. Run full `rtk go test ./...` from `sales`; record unrelated failures if any.
7. If local runtime/token available, smoke 4 `group_by` responses and record actual sample output in evidence.
8. Run `@quality-gate` final review before claiming done because this touches report SQL + tenant-sensitive filters.

## Expected Files to Change
Likely:
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`

Possibly:
- `sales/service/report_service_test.go`
- `sales/controller/so_controller_test.go` only if outlet expected output currently assumes old mapping.
- `sales/client_test.http` optional only.

## Agent/Tool Routing
- `@orchestrator`: coordinate implementation and validation.
- `@fixer`: bounded code/test edits.
- `@explorer`: only if executor needs more schema/pattern discovery.
- `@quality-gate`: final signoff required.
- `@librarian`: not needed unless external docs unexpectedly become material.

## Executor Handoff Prompt
```text
Implement SX-2224 using `.opencode/plans/20260612-0916-sx-2224-secondary-sales-group-code.md` as the source of truth. Scope is the `sales` module only. Most code/name support already exists from SX-2172; the required source change is to make `group_by=outlet` follow SX-2224 outlet-murni mapping: `id=o/r.outlet_id`, `code=mo.outlet_code`, `name=mo.outlet_name`, not salesman id or `emp_name > outlet_name`. Preserve salesman/product_category/product master mappings, `cust_id IN ?`, date range from month/year, return subtraction, and JSON aliases `id/code/name/net_sales`. Use TDD: update repository dry-run SQL test first, then edit `buildSecondarySalesReportGroupQuery`. Do not touch env, migrations, package files, lockfiles, unrelated endpoints, or `pjp-sales`. Validate from `sales` with `rtk go test ./repository -run 'TestSecondarySalesReportGroup'`, `rtk go test ./service -run 'TestSecondarySalesReportGroupSales'`, `rtk go test ./controller -run 'TestSecondaryReportSalesGroup'`, and `rtk go test ./...`. Return changed files, test outputs, SQL behavior summary, and any unrelated failures with evidence.
```

## Execution-ready Worklist / Handoff Contract
`start_with`: `T1`

### T1 — Update outlet regression test
- `depends_on`: none
- `owner/lane`: `@fixer`
- `action`: Replace old SX-2172 outlet SQL expectation with SX-2224 outlet-murni expectation in repository dry-run test.
- `validation`: `rtk go test ./repository -run 'TestSecondarySalesReportGroupOutlet'`
- `exit_criteria`: Test asserts `outlet_id`, `outlet_code`, `outlet_name` and rejects `salesman_id` / `CONCAT_WS` for outlet branch.
- `blocking_status`: ready
- `requires_user_decision`: no
- `must_preserve`: existing date range and return subtraction assertions.
- `do_not_touch`: source SQL until test expectation is updated.
- `evidence_update`: record test name and intended Red/Green behavior.
- `exit_verification`: failing or passing targeted output with explanation if source was changed in same iteration.

### T2 — Change outlet SQL mapping
- `depends_on`: T1
- `owner/lane`: `@fixer`
- `action`: Update outlet entries in `orderSelect`, `returnSelect`, `orderJoin`, `returnJoin` inside `buildSecondarySalesReportGroupQuery`.
- `validation`: `rtk go test ./repository -run 'TestSecondarySalesReportGroup'`
- `exit_criteria`: Repository group tests pass; SQL uses outlet id/code/name and still includes source table/date/cust filters.
- `blocking_status`: ready
- `requires_user_decision`: no
- `must_preserve`: product/category/salesman mappings and return subtraction.
- `do_not_touch`: service/controller unless tests require expectation updates.
- `evidence_update`: capture SQL fragments or test output.
- `exit_verification`: targeted repository test pass.

### T3 — Verify service/controller contract
- `depends_on`: T2
- `owner/lane`: `@fixer`
- `action`: Run targeted service/controller tests and update only tests that assume old outlet mapping.
- `validation`: `rtk go test ./service -run 'TestSecondarySalesReportGroupSales' && rtk go test ./controller -run 'TestSecondaryReportSalesGroup'`
- `exit_criteria`: Tests pass or unrelated pre-existing failures are documented.
- `blocking_status`: ready
- `requires_user_decision`: no
- `must_preserve`: response envelope and auth behavior.
- `do_not_touch`: controller implementation unless actual regression proves necessary.
- `evidence_update`: record changed tests and outputs.
- `exit_verification`: targeted service/controller test output.

### T4 — Full validation and optional runtime smoke
- `depends_on`: T3
- `owner/lane`: `@fixer`
- `action`: Run full sales tests and optional curl smoke for 4 group_by variants if runtime/token available.
- `validation`: `rtk go test ./...`
- `exit_criteria`: Full tests pass, or unrelated failures documented with exact command output; optional smoke confirms `code`/`name` fields.
- `blocking_status`: ready
- `requires_user_decision`: no
- `must_preserve`: no extra files outside diff boundary.
- `do_not_touch`: secrets/env.
- `evidence_update`: write verification evidence under `.opencode/evidence/20260612-0916-sx-2224-secondary-sales-group-code/verification.md` if implementation proceeds.
- `exit_verification`: test outputs and optional before/after JSON snippets.

### T5 — Final quality gate
- `depends_on`: T4
- `owner/lane`: `@quality-gate`
- `action`: Review diff, tests, tenant filters, response contract, and evidence.
- `validation`: quality-gate checklist + changed file review.
- `exit_criteria`: PASS or explicit remediation list.
- `blocking_status`: ready
- `requires_user_decision`: no
- `must_preserve`: all invariants in this plan.
- `do_not_touch`: implementation edits by reviewer.
- `evidence_update`: final review result.
- `exit_verification`: quality-gate summary.

## Validation Commands
Run from `/Users/ujang/Projects/Geekgarden/scylla-be/sales`:
```bash
rtk go test ./repository -run 'TestSecondarySalesReportGroup'
rtk go test ./service -run 'TestSecondarySalesReportGroupSales'
rtk go test ./controller -run 'TestSecondaryReportSalesGroup'
rtk go test ./...
```

Optional DB check already used during planning:
```sql
SELECT pro_id, pro_code, pro_name, cust_id, pcat_id FROM mst.m_product WHERE pro_id = 10733;
SELECT pcat_id, pcat_code, pcat_name FROM mst.m_product_cat WHERE pcat_id = 77;
```

## Evidence Requirements
- Keep `.opencode/evidence/20260612-0916-sx-2224-secondary-sales-group-code/discovery.md` because it contains repo/DB evidence and the user decision gate outcome.
- During implementation, add `verification.md` with command outputs and optional curl responses.
- Evidence must state whether smoke used live runtime or only dry-run/unit tests.
- Source strategy used: repo-local code/tests + DB local `ggn_scyllax`; external docs/web skipped as irrelevant to local SQL mapping.

## Done Criteria
- Outlet branch uses outlet id/code/name.
- Salesman/product_category/product branch still satisfy SX-2224 code/name requirements.
- Tests listed in Validation Commands are run and results recorded.
- No out-of-boundary source changes remain.
- Quality gate passes or remediation is completed.

## Final Planning Summary
Artifacts created:
- Primary plan: `.opencode/plans/20260612-0916-sx-2224-secondary-sales-group-code.md`
- Evidence kept: `.opencode/evidence/20260612-0916-sx-2224-secondary-sales-group-code/discovery.md`

Key decisions:
- User selected SX-2224 outlet-murni mapping over existing SX-2172 outlet salesman mapping.
- Product `10733` was checked in local DB and has real master name; no synthetic `Product <id>` fallback is planned.
- Product/category source remains master `mst.*`.

Open questions:
- None blocking for implementation. Optional future question: whether `pjp-sales` should mirror this endpoint behavior separately.

Cleanup:
- No draft artifacts were created.
- Evidence is intentionally kept because DB findings and decision gate are operationally useful for executor and reviewer.

Readiness: `ready-for-implementation`.
