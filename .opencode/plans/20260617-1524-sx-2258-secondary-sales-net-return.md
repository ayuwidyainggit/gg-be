# Goal

Memperbaiki kalkulasi Backend Secondary Sales untuk `SX-2258` agar metrik order-based dikurangi return:

- `Number of Product Sold` / response field `qty` = total qty order - total qty return.
- `Discount and Promo` / response field `total_discount_promo` = discount+promo order - discount+promo return.
- Filter order dan return konsisten untuk `cust_id`, date range, outlet, salesman, product, dan `o.data_status IN (6,7)`.

Readiness: `ready-for-implementation`.

# Non-goals

- Tidak mengubah PPN rule tanpa konfirmasi Product/QA.
- Tidak mengubah FE formatter currency/number.
- Tidak mengubah export Excel detail kecuali ditemukan bahwa FE memakai export path untuk metrik yang sama.
- Tidak memasukkan kredensial Jira/staging ke repo, log, plan, atau test.
- Tidak mengubah schema DB/migration.
- Tidak mengubah modul selain target runtime tanpa bukti deployment butuh mirror patch.

# Scope

Scope utama:

- Modul `sales/`.
- Dashboard summary endpoint:
  - `GET /v1/reports/secondary-sales/sum-date`
  - `sales/controller/report_controller.go`
  - `sales/service/report_service.go`
  - `sales/repository/report_repository.go`
- Test target:
  - `sales/repository/report_repository_test.go`
  - `sales/service/report_service_test.go` bila response mapping/filter propagation perlu coverage.

Scope audit tambahan:

- `GET /v1/reports/secondary-sales/trend-sales` karena query trend punya bug `discount_promo order + return` juga.
- `GET /v1/reports/secondary-sales/group` untuk memastikan net sales existing tidak ikut rusak.
- `pjp-sales/` hanya mirror patch bila deployment/staging ternyata memakai module itu, atau branch policy repo mengharuskan parity `sales/` dan `pjp-sales/`.

# Requirements

1. `qty` pada summary card harus net:
   - `COALESCE(order_summary.qty, 0) - COALESCE(return_summary.qty_return, 0)`.
2. `total_discount_promo` pada summary card harus net:
   - `COALESCE(order_summary.discount_promo, 0) - COALESCE(return_summary.discount_promo, 0)`.
3. Order discount/promo mengikuti referensi QA:
   - `COALESCE(od.disc_value_final, 0)`
   - `COALESCE(od.promo_final1, 0) + ... + COALESCE(od.promo_final5, 0)`
4. Return discount/promo mengikuti referensi QA:
   - `COALESCE(rd.disc_value, 0) + COALESCE(rd.promo_value, 0)`.
5. Order qty mengikuti referensi QA:
   - `(qty3 * conv2 * conv3) + (qty2 * conv2) + qty1` dengan `COALESCE`.
6. Return qty mengikuti referensi QA:
   - `(qty3 * conv2 * conv3) + (qty2 * conv2) + qty1` dengan `COALESCE`.
7. Return harus join invoice order:
   - `r.return_no = rd.return_no`
   - `r.cust_id = rd.cust_id`
   - `o.invoice_no = r.invoice_no`
   - `o.cust_id = r.cust_id`
8. Filter order side:
   - `o.cust_id IN ?`
   - `o.data_status IN (6,7)`
   - `o.invoice_date` date range
   - `o.outlet_id`
   - `o.salesman_id`
   - `od.pro_id`
9. Filter return side:
   - `rd.cust_id IN ?`
   - `o.data_status IN (6,7)`
   - `o.invoice_date` date range
   - `r.outlet_id`
   - `r.salesman_id`
   - `rd.product_id`
10. Existing `month/year` behavior tetap jalan. Bila `from/to` query tersedia, gunakan `from/to`; jika tidak, derive date range dari `month/year`.

# Acceptance Criteria

1. Dataset QA menghasilkan:
   - `Number of Product Sold = 134`
   - `Discount and Promo = 1.238.740`
2. `total_discount_promo` tidak lagi memakai `order + return`.
3. `qty` tidak lagi hanya `order qty`, tetapi `order qty - return qty`.
4. Date/outlet/salesman/product filters diterapkan simetris pada order dan return.
5. `Return Value`, `Return Rate`, `Gross Sales`, `Net Sales exclude PPN`, `Net Sales include PPN`, dan `Total PPN` tidak berubah tanpa bukti QA/Product.
6. Test otomatis minimal mencakup:
   - order tanpa return
   - order dengan return sebagian
   - filter product/outlet/salesman/date
   - null discount/promo/vat/qty
7. Targeted tests di `sales/` pass.
8. Jika staging/local DB tersedia, response endpoint diverifikasi terhadap SQL QA.

# Existing Patterns/Reuse

Reuse:

- `sales/controller/report_controller.go:521` route `SecondaryReportSalesSumMonth` sudah parse query dan auth context.
- `sales/service/report_service.go:1374` sudah resolve authorized `cust_id`/`parent_cust_id`.
- `sales/repository/report_repository.go:1135` sudah punya CTE `order_summary` dan `return_summary`; ini titik fix utama.
- `sales/repository/report_repository_test.go` sudah punya dry-run SQL tests untuk summary, trend, group.
- `sales/service/report_service_test.go` sudah punya mock repository untuk service-level coverage.

Jangan reimplement:

- Auth/scope resolver `resolveSecondaryDashboardCustIDs`.
- Controller response wrapper.
- Group query net-sales union, kecuali filter propagation memang jadi scope baru.

Repo-local evidence cukup. Official docs/GitHub/web skipped karena bug SQL internal dan library behavior tidak version-sensitive.

# Constraints

- Ikuti layering: Controller -> Service -> Repository -> DB.
- Semua query tetap parameterized.
- Jangan hardcode QA `cust_id`, date, outlet, salesman, product, atau expected numbers ke production code.
- Jangan copy kredensial staging/Jira.
- `rtk` wajib untuk shell workflow repo ini.
- Validate dari `sales/`, bukan root Go module.
- `PROJECT_STACK.md`, `PROJECT_COMMANDS.md`, `FRAMEWORK_PLAYBOOK.md` tidak ada; pakai `.opencode/docs/QUALITY.md` dan `SERVICE_MATRIX.md`.

# Risks

1. **Field mismatch promo**
   - Current summary memakai `promo_value_final`; QA reference memakai `promo_final1..5`.
   - Mitigasi: Red test harus mengunci formula QA atau staging compare harus membuktikan field canonical.
2. **PPN ambiguity**
   - Jira reference total PPN punya formula berbeda dari existing; user minta jangan ubah PPN.
   - Mitigasi: jangan ubah `total_ppn` kecuali Product/QA minta.
3. **Filter contract gap**
   - Current `SecondarySalesReportDashboardSumPayload` belum punya `from/to`, outlet, salesman, product.
   - Mitigasi: tambah optional query fields backward-compatible; default tetap `month/year`.
4. **Mirror module drift**
   - `pjp-sales/` punya file serupa.
   - User decision: staging target untuk SX-2258 adalah `sales`; jangan mirror `pjp-sales/` untuk task ini kecuali instruksi baru muncul.
5. **Double-count return**
   - Return join invoice bisa double bila invoice tidak unik per cust.
   - Mitigasi: verify with SQL counts in staging/local; do not add joins beyond QA reference.
6. **Tests only assert SQL text**
   - Dry-run SQL tests tidak membuktikan numeric result.
   - Mitigasi: tambah pure arithmetic regression test dan staging/local DB verification bila tersedia.

# Decisions/Assumptions

## Decisions

- Primary source of truth implementation: `.opencode/plans/20260617-1524-sx-2258-secondary-sales-net-return.md`.
- Target module dikunci oleh user: `sales`.
- Target utama: `sales/repository/report_repository.go` summary query.
- Trend `total_discount_promo` juga wajib ikut fix di PR yang sama.
- Formula promo/discount mengikuti docs/Jira reference: order `disc_value_final + promo_final1..5`, return `disc_value + promo_value`.
- `qty` alias tetap dipakai untuk response compatibility, tetapi nilainya menjadi net sold qty.
- `qty_return` tetap return qty bruto untuk display `Number of Product Return`.
- `total_discount_promo` menjadi net discount/promo.
- `total_ppn` tidak diubah.
- Add optional filters to summary endpoint in backward-compatible way.

## Assumptions / Open Questions

- User answered: staging target module adalah `sales`.
- User answered: trend `total_discount_promo` perlu ikut fix.
- User answered: promo/discount formula ikuti docs/Jira reference, bukan asumsi `promo_value_final`.
- Assumption: FE dashboard `Number of Product Sold` dan `Discount and Promo` membaca `GET /v1/reports/secondary-sales/sum-date`.
- Assumption: Existing response field `qty` adalah FE label `Number of Product Sold`.
- Assumption: `month/year` tetap wajib bagi current FE; optional `from/to` dapat override date range bila FE mengirim date range.
- Open validation: staging credentials/dataset QA harus digunakan oleh engineer dengan mekanisme aman internal, bukan dari plan.

Question gate: tidak bertanya sekarang karena user sudah memberi formula dan acceptance. Ambiguity filter/date ditangani lewat backward-compatible assumption.

# Execution Source of Truth

Precedence saat implementasi:

1. Instruksi user terbaru untuk `SX-2258`.
2. Security rule: jangan expose/copy secrets, jangan pakai remote DB default tanpa izin.
3. Non-negotiable Implementation Invariants di plan ini.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria dan Done Criteria.
6. Implementation Steps.
7. Follow-up/recommendation.

Jika konflik terjadi, executor wajib ikuti sumber lebih tinggi dan catat konflik di evidence.

# Non-negotiable Implementation Invariants

- `qty` summary card harus net sold qty: order qty - return qty.
- `qty_return` tetap return qty, bukan net qty.
- `total_discount_promo` harus order discount/promo - return discount/promo.
- Order discount/promo harus memakai `COALESCE(od.disc_value_final, 0) + COALESCE(od.promo_final1, 0) + COALESCE(od.promo_final2, 0) + COALESCE(od.promo_final3, 0) + COALESCE(od.promo_final4, 0) + COALESCE(od.promo_final5, 0)`.
- Return discount/promo harus memakai `COALESCE(rd.disc_value, 0) + COALESCE(rd.promo_value, 0)`.
- Return date filter untuk SX-2258 summary harus memakai `o.invoice_date`, bukan `r.return_date`.
- Return outlet/salesman/product filters harus memakai `r.outlet_id`, `r.salesman_id`, `rd.product_id`.
- Order status filter tetap `o.data_status IN (6,7)` pada order dan return branch.
- Return join invoice wajib ada: `JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id`.
- PPN formula jangan diubah tanpa explicit Product/QA decision.
- Test harus fail pada formula lama `discount_promo + return` dan `os.qty AS qty`.
- Source edits hanya oleh implementation lane, bukan planner.

# Do Not / Reject If

Reject implementation jika:

- `total_discount_promo` masih mengandung `+ rs.discount_promo` di summary SX-2258 path.
- `qty` masih `os.qty AS qty` tanpa subtract return.
- Return branch difilter dengan `r.return_date` untuk summary SX-2258.
- Product filter return memakai `od.pro_id` atau order side alias yang salah.
- Outlet/salesman return memakai `o.outlet_id` / `o.salesman_id` padahal QA reference minta `r.outlet_id` / `r.salesman_id`.
- Query string concat user input langsung.
- PPN diubah untuk mengejar angka tanpa QA confirmation.
- Staging credentials ditulis ke test/plan/log.
- Tests hanya update expected text tanpa Red/Green reason.

# Diff Boundary

Allowed source files:

- `sales/entity/report.go`
- `sales/controller/report_controller.go` bila query parsing perlu adjustment.
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`
- `sales/service/report_service_test.go`
- `sales/controller/so_controller_test.go` atau controller report test bila compile break dari interface change.

Conditional allowed:

- `pjp-sales/...` tidak termasuk scope SX-2258 karena user sudah mengunci target ke `sales`; hanya ubah bila ada instruksi baru eksplisit.

Allowed evidence files:

- `.opencode/evidence/20260617-1524-sx-2258-secondary-sales-net-return/**`

Out-of-boundary changes must be reverted or justified in verification evidence before final quality gate.

# TDD/Test Plan

TDD required: yes.

Reason: production report math, money/qty metrics, regression from QA failed staging.

## Existing test patterns

- Dry-run SQL assert in `sales/repository/report_repository_test.go`.
- Service-level mock in `sales/service/report_service_test.go`.

## First failing tests

1. Update/add repository dry-run test for `SecondarySalesReportSumReportByMonth`:
   - expect `COALESCE(os.qty, 0) - COALESCE(rs.qty_return, 0) AS qty`
   - reject `os.qty AS qty`
   - expect `COALESCE(os.discount_promo, 0) - COALESCE(rs.discount_promo, 0) AS total_discount_promo`
   - reject `(os.discount_promo + rs.discount_promo)` and `COALESCE(os.discount_promo, 0) + COALESCE(rs.discount_promo, 0)`
   - expect return invoice join.
2. Add filter SQL test for summary:
   - input has `from/to`, `outlet_ids`, `salesman_ids`, `pro_ids`
   - expect order side `o.invoice_date`, `o.outlet_id`, `o.salesman_id`, `od.pro_id`
   - expect return side `o.invoice_date`, `r.outlet_id`, `r.salesman_id`, `rd.product_id`
3. Add arithmetic regression test:
   - order qty `150`, return qty `16`, expect `134`
   - order discount/promo `1_500_000`, return discount/promo `261_260`, expect `1_238_740`
   - null promo parts treated as zero.
4. If trend endpoint remains in scope, update trend test to reject plus formula for `total_discount_promo`.

## Green step

- Refactor repository summary function to accept filter payload or new internal filter struct.
- Compute date range from `from/to` when present, otherwise from `month/year`.
- Replace summary select formulas.
- Add optional filters to SQL with parameter binding.
- Use individual `COALESCE(od.promo_finalN, 0)` for order promo if staging/query evidence agrees with QA reference.

## Refactor step

- Extract small helper for date range/filter building if raw query becomes hard to read.
- Keep alias names stable for response mapping.
- Keep tests readable and specific to business formula.

## Edge cases

- Order without return -> net qty equals order qty; net discount equals order discount.
- Partial return -> subtract return qty/discount.
- Full return -> net can become zero.
- Null qty/conv -> `COALESCE` gives safe zero/one.
- Null individual promo fields -> treated as zero.
- Empty filters -> no extra filter clauses beyond cust/date/status.
- Multi-cust parent request -> `cust_id IN ?` still works.

# Implementation Steps

1. Confirm runtime path.
   - Use code evidence: `GET /v1/reports/secondary-sales/sum-date` -> service -> repository summary.
2. Add/adjust tests first.
   - Repository dry-run SQL tests for net qty, net discount promo, filters.
   - Arithmetic test for QA-style numbers.
3. Extend summary filter input.
   - Add optional fields to `SecondarySalesReportDashboardSumPayload`:
     - `From *int64 query:"from"`
     - `To *int64 query:"to"`
     - `OutletIDs []int64 query:"outlet_ids"`
     - `SalesmanIDs []int64 query:"salesman_ids"`
     - `ProIDs []int64 query:"pro_ids"`
   - Keep existing `Month`, `Year`, `CustID`, `CustIDs`.
4. Change repository interface/function.
   - Prefer passing a small filter struct or existing payload from service to repository.
   - Preserve `custIDs` resolved by service.
5. Build date range.
   - If `From` and `To` exist: convert via existing `str.UnixTimestampToUtcTime`.
   - Else: `dateFrom = first day month/year`, `dateTo = next month`.
6. Update `SecondarySalesReportSumReportByMonth` SQL.
   - Replace `total_discount_promo` plus with subtract.
   - Replace `os.qty AS qty` with net qty.
   - Keep `rs.qty_return AS qty_return`.
   - Keep `total_ppn` unchanged.
   - Apply optional filters symmetrically.
7. Align order discount promo source.
   - If implementing QA formula directly, use:
     - `COALESCE(od.disc_value_final, 0) + COALESCE(od.promo_final1, 0) + ... + COALESCE(od.promo_final5, 0)`.
   - If keeping `promo_value_final`, executor must add evidence showing it equals QA promo sum for target dataset and get Product/QA acceptance.
8. Audit trend sales.
   - Change trend `total_discount_promo` from plus to subtract if FE uses same metric in trend.
   - Keep this in same PR only if test confirms no unrelated behavior break.
9. Run tests.
10. If DB access available, compare local/staging response to QA SQL.
11. Run quality-gate review.

# Expected Files to Change

Primary:

- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`
- `sales/entity/report.go`
- `sales/service/report_service.go`

Possible:

- `sales/service/report_service_test.go`
- `sales/controller/so_controller_test.go` if mock interface compile breaks.
- `pjp-sales/repository/report_repository.go` and tests only if mirror target confirmed.

# Agent/Tool Routing

- `@artifact-planner`: plan/evidence only.
- `@fixer` or backend implementation lane: source edits + tests.
- `@oracle`: optional review if formula/promo field ambiguity remains.
- `@quality-gate`: required final signoff because money/qty report logic and staging QA failed.
- Staging DB verification: engineer/internal safe DB access only.

# Executor Handoff Prompt

Copy ke `@orchestrator` / implementation lane:

```text
Implement SX-2258 using `.opencode/plans/20260617-1524-sx-2258-secondary-sales-net-return.md` as source of truth.

Scope: module `sales/`, Secondary Sales dashboard summary. Fix `qty` to order qty minus return qty and `total_discount_promo` to order discount/promo minus return discount/promo. Add backward-compatible optional filters for `from/to`, `outlet_ids`, `salesman_ids`, `pro_ids` if needed by summary endpoint. Preserve auth cust scope, `o.data_status IN (6,7)`, return invoice join, and PPN behavior.

must_preserve:
- Controller -> Service -> Repository layering.
- `qty_return`, `net_sales_return`, `return_rate` semantics.
- `total_ppn` unchanged unless Product/QA explicitly approves.
- parameterized SQL only.
- no secrets in files/logs.

do_not_touch:
- migrations/schema.
- FE code.
- unrelated report modules.
- `pjp-sales/` unless deployment target requires mirror patch.

validation:
- From repo root: `rtk docker compose -f docker-compose.yml ps`.
- From `sales/`: `rtk go test ./repository -run 'SecondarySalesReport(SumReportByMonth|TrendSales)'`.
- From `sales/`: `rtk go test ./service -run SecondarySalesReportSumReportByMonth` if service tests changed.
- From `sales/`: `rtk go test ./...` before final.
- If staging/local DB available, compare endpoint response against QA SQL and verify `qty=134`, `total_discount_promo=1238740`.

return evidence:
- changed files.
- before/after formula.
- test commands output.
- staging/local response or reason not run.
- any unresolved promo field decision.
```

# Execution-ready Worklist / Handoff Contract

`start_with`: `T1`

## T1 — Add failing SQL regression tests

- `depends_on`: none
- `owner/lane`: `@fixer`
- action: Update `sales/repository/report_repository_test.go` to fail on old summary formulas and verify return invoice/filter aliases.
- validation: `rtk go test ./repository -run TestSecondarySalesReportSumReportByMonth`
- exit criteria: tests fail before production change for `+ rs.discount_promo` and `os.qty AS qty`.
- blocking status: `ready`
- blocker reason: none
- requires_user_decision: no
- must_preserve: dry-run test pattern; no DB dependency.
- do_not_touch: production query until Red state captured.
- evidence_update: record failing test names/output in task evidence.
- exit_verification: failing output shows expected old-formula mismatch.

## T2 — Add arithmetic regression test

- `depends_on`: T1
- `owner/lane`: `@fixer`
- action: Add pure test for `150 - 16 = 134` and `1_500_000 - 261_260 = 1_238_740`, including null promo pieces as zero.
- validation: `rtk go test ./repository -run SX2258`
- exit criteria: test exists and documents QA formula.
- blocking status: `ready`
- blocker reason: none
- requires_user_decision: no
- must_preserve: no hardcoded QA cust/date in production code.
- do_not_touch: staging credentials.
- evidence_update: record test name.
- exit_verification: test passes after formula helper/query fix or fails only before expected production change.

## T3 — Extend summary filter contract backward-compatible

- `depends_on`: T1
- `owner/lane`: `@fixer`
- action: Add optional query fields to summary payload and pass resolved filters from service to repository.
- validation: `rtk go test ./service -run SecondarySalesReportSumReportByMonth`
- exit criteria: existing month/year requests still work; optional filters reach repository test/mocks.
- blocking status: `ready`
- blocker reason: none
- requires_user_decision: no
- must_preserve: auth cust resolution.
- do_not_touch: export request DTO unless needed separately.
- evidence_update: note payload fields added.
- exit_verification: service tests compile/pass.

## T4 — Fix summary query formulas and filters

- `depends_on`: T1, T2, T3
- `owner/lane`: `@fixer`
- action: Modify `SecondarySalesReportSumReportByMonth` SQL to net qty and net discount promo, plus optional filters.
- validation: `rtk go test ./repository -run TestSecondarySalesReportSumReportByMonth`
- exit criteria: repository tests pass and reject old formula.
- blocking status: `ready`
- blocker reason: none
- requires_user_decision: no
- must_preserve: `total_ppn`, `net_sales_return`, `return_rate` existing semantics.
- do_not_touch: unrelated report queries.
- evidence_update: include before/after SQL formula snippets.
- exit_verification: SQL text contains net formulas and correct aliases.

## T5 — Fix trend discount formula

- `depends_on`: T4
- `owner/lane`: `@fixer`
- action: Change `buildSecondarySalesReportTrendSalesSQL` `total_discount_promo` from order+return to order-return; user confirmed trend must be fixed in same task.
- validation: `rtk go test ./repository -run TestSecondarySalesReportTrendSalesSQLUsesSourceTablesAndNetSalesFormula`
- exit criteria: trend test updated and rejects plus formula.
- blocking status: `ready`
- blocker reason: none
- requires_user_decision: no
- must_preserve: yearly trend month list and `net_sales` formula.
- do_not_touch: group query unless tests reveal break.
- evidence_update: record trend formula before/after.
- exit_verification: trend tests pass and SQL contains subtract formula.

## T6 — Full local validation

- `depends_on`: T4, T5
- `owner/lane`: `@fixer`
- action: Run targeted then full sales tests.
- validation: `rtk go test ./repository -run 'SecondarySalesReport(SumReportByMonth|TrendSales)' && rtk go test ./service -run SecondarySalesReport && rtk go test ./...`
- exit criteria: tests pass or failures documented as unrelated with evidence.
- blocking status: `ready`
- blocker reason: none
- requires_user_decision: no
- must_preserve: no skipped tests without reason.
- do_not_touch: `.env`, secrets, lockfiles unless Go tooling legitimately changes none.
- evidence_update: paste command outputs summary.
- exit_verification: full test result captured.

## T7 — Staging/local data verification

- `depends_on`: T6
- `owner/lane`: implementation engineer with safe DB/API access
- action: Compare endpoint response to QA SQL for Jira dataset.
- validation: endpoint/API response and SQL result show `qty=134`, `total_discount_promo=1238740`.
- exit criteria: QA expected numbers match or discrepancy documented with row-level diff.
- blocking status: `blocked`
- blocker reason: planner lacks staging credentials/access; engineer must use safe internal access.
- requires_user_decision: no
- must_preserve: no credentials in repo/evidence.
- do_not_touch: Jira secrets.
- evidence_update: store sanitized response/sample under `.opencode/evidence/20260617-1524-sx-2258-secondary-sales-net-return/`.
- exit_verification: sanitized evidence path listed.

## T8 — Quality gate

- `depends_on`: T6; T7 if DB access available
- `owner/lane`: `@quality-gate`
- action: Review formula, filters, tests, scope, secrets, and regression risk.
- validation: quality-gate report.
- exit criteria: pass or remediation tasks created.
- blocking status: `ready`
- blocker reason: none
- requires_user_decision: no
- must_preserve: evidence over assertion.
- do_not_touch: source files during review.
- evidence_update: final quality summary.
- exit_verification: quality gate result recorded.

# Validation Commands

Repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

Module `sales/`:

```bash
rtk go test ./repository -run TestSecondarySalesReportSumReportByMonth
rtk go test ./repository -run TestSecondarySalesReportTrendSalesSQLUsesSourceTablesAndNetSalesFormula
rtk go test ./service -run SecondarySalesReportSumReportByMonth
rtk go test ./...
```

Optional DB/API verification jika safe access tersedia:

```bash
# Run QA reference SQL with safe internal DB access; do not store credentials.
# Call GET /v1/reports/secondary-sales/sum-date with same sanitized params.
```

# Evidence Requirements

Already created/kept:

- `.opencode/evidence/20260617-1524-sx-2258-secondary-sales-net-return/discovery.md`
- `.opencode/evidence/20260617-1524-sx-2258-secondary-sales-net-return/index.json`

Implementation must add sanitized evidence:

- failing Red test output
- passing Green test output
- full `rtk go test ./...` summary
- before/after formula snippets
- staging/local response sample if available
- note that `pjp-sales/` mirror was not changed because user confirmed target `sales`

Source strategy:

- Used: repo-local docs, code, tests, prior Secondary Sales plan, compose status.
- Skipped official docs/context7: no external library behavior needed.
- Skipped GitHub/web: issue is internal SQL/report logic.
- Skipped browser/screenshots: backend numeric defect, no UI parity task.
- Staging SQL skipped by planner because credentials must remain internal; executor should run if access exists.

# Done Criteria

- Code formulas match SX-2258 requirements.
- Tests prove old formulas fail and new formulas pass.
- Optional filters do not break existing `month/year` dashboard calls.
- No PPN behavior changed without explicit evidence.
- No secrets or staging creds stored.
- `sales/` tests pass.
- Staging/local verification completed or blocked reason documented.
- `@quality-gate` final signoff completed.

# Final Planning Summary

Artifacts created:

- `.opencode/plans/20260617-1524-sx-2258-secondary-sales-net-return.md`
- `.opencode/evidence/20260617-1524-sx-2258-secondary-sales-net-return/discovery.md`
- `.opencode/evidence/20260617-1524-sx-2258-secondary-sales-net-return/index.json`

Key decisions:

- Fix primary summary endpoint in `sales/`.
- User confirmed target module `sales`; `pjp-sales` out of scope.
- Net `qty` and `total_discount_promo` in repository SQL.
- Trend `total_discount_promo` wajib ikut fixed.
- Promo/discount formula mengikuti docs/Jira reference: `disc_value_final + promo_final1..5` untuk order, `disc_value + promo_value` untuk return.
- Add backward-compatible optional filters if summary endpoint must support QA filters.
- Preserve PPN.
- Keep evidence because implementation replay needs discovered route/query/test paths.

Assumptions:

- FE label maps `Number of Product Sold` to response `qty`.
- FE label maps `Discount and Promo` to response `total_discount_promo`.

Open questions:

- Tidak ada question produk/scope yang masih blocking implementasi core.
- Staging exact-number validation tetap menunggu safe DB/API access.

Cleanup:

- No draft artifacts created.
- Evidence kept intentionally for implementation replay.

Plan Quality Gate: `PASS_FOR_SLICE` for source implementation; staging exact-number validation remains blocked until safe DB/API access exists.
