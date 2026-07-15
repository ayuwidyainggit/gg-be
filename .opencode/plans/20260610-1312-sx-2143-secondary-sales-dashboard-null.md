# Plan — SX-2143 Secondary Sales Dashboard Data Null

Task ID: `20260610-1312-sx-2143-secondary-sales-dashboard-null`
Readiness: `ready-for-implementation`
Quality Gate: `PASS_FOR_SLICE`
Primary source of truth: this file.

## Goal

Perbaiki dan verifikasi defect SX-2143 pada backend `sales`: dashboard Secondary Sales tidak boleh mengembalikan data numeric `null`/kosong untuk periode Juni 2026 ketika data ada, terutama endpoint:

```http
GET /sales/v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001
```

Plan ini juga mencakup fallback ketika FE belum mengirim `year`, serta regression coverage untuk endpoint `group` yang memakai filter `month`, `year`, `cust_id`, dan `group_by`.

## Non-goals

- Tidak hardcode token, credential, staging base URL rahasia, atau data sensitif dari Jira/Docs.
- Tidak memperbaiki URL salah format seperti `month=6?year2026`; API hanya perlu mendukung format query yang benar `month=6&year=2026`.
- Tidak membuka akses arbitrary `cust_id` di luar scope user.
- Tidak merombak extract pipeline atau membuat migrasi DB kecuali evidence DB membuktikan facts tidak bisa dipakai.
- Tidak mengubah kontrak FE selain menambahkan/menjaga field yang sudah disepakati seperti `code` pada group bila branch target memang membutuhkannya.
- Tidak mengubah modul selain `sales`.

## Scope

Target utama:

- `sales/controller/report_controller.go`
- `sales/entity/report.go`
- `sales/model/report.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/controller/so_controller_test.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`

Target opsional:

- `sales/client_test.http` untuk contoh manual request non-secret.
- `.opencode/evidence/20260610-1312-sx-2143-secondary-sales-dashboard-null/**` untuk evidence implementasi dan smoke results.

## Requirements

- `sum-date` membaca query param `month`, `year`, dan `cust_id`.
- `group` membaca query param `month`, `year`, `cust_id`, dan `group_by`.
- `month` valid `1..12`; invalid harus 400 sesuai pola validator.
- `year` optional, integer 4 digit dalam batas validator saat dikirim; invalid harus 400.
- Jika `year` dikirim, gunakan sebagai filter utama untuk summary dan semua group branch.
- Jika `year` tidak dikirim, fallback BE harus memakai `time.Now().Year()` supaya request existing dari FE tidak menghasilkan numeric `null`. Catatan: ini slice-safe untuk runtime 2026; FE tetap wajib mengirim `&year=2026` untuk determinisme lintas tahun.
- `cust_id` efektif harus mengikuti aturan auth:
  - distributor user default ke `cust_id` login saat `cust_id` tidak dikirim.
  - distributor user tidak boleh mengambil sibling/arbitrary `cust_id`.
  - principal user boleh mengambil child/distributor hanya jika scope check parent-child lolos.
- Summary `sum-date` harus memakai filter `cust_id` dan `month + year` atau date range ekuivalen yang tepat untuk 1–30 Juni 2026.
- Numeric aggregate fields harus number, bukan JSON `null`; gunakan `COALESCE` di SQL dan/atau non-pointer response normalization.
- Response `sum-date` harus menjaga field minimal:
  - `total_gross_sale`
  - `total_discount_promo`
  - `total_ppn`
  - `net_sales_exc_ppn`
  - `net_sales`
  - `total_salesman`
  - `total_outlet`
  - `total_product`
  - `qty`
  - `qty_return`
  - `return_rate`
  - `net_sales_return`
  - `last_update`
- Summary order + return harus konsisten dengan BE docs dari prompt: order valid dikurangi return; PPN dan net sales dihitung null-safe. Implementasi boleh memakai `report.fact_orders`/`report.fact_returns` jika facts terbukti sinkron dengan source tables.
- Group endpoints harus menerapkan filter `cust_id`, `month`, dan `year` pada semua branch:
  - `outlet`
  - `salesman`
  - `product_category`
  - default/`product`
- Group endpoints harus mempertahankan return subtraction (`fact_returns` negatif) dan, jika branch target sudah ada enhancement SX-2172, tetap return `code`, `name`, `net_sales`.

## Acceptance Criteria

- `GET /sales/v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001` mengembalikan dashboard Secondary Sales untuk periode 1–30 Juni 2026 bila data staging/local ada.
- `GET /sales/v1/reports/secondary-sales/sum-date?month=6&cust_id=C260020001` tidak menghasilkan numeric `null`; BE memakai fallback year.
- BE support query param `year` di `sum-date` dan `group` endpoints.
- Query memakai `cust_id = 'C260020001'` dari request/effective auth; tidak mengubahnya menjadi `C2600200001` atau variasi lain.
- Semua numeric aggregate fields di response `sum-date` adalah angka; default `0` jika tidak ada row.
- `return_rate` aman ketika `qty = 0`.
- Group branch `outlet`, `salesman`, `product_category`, dan default/`product` semuanya memakai year filter yang sama dan tetap subtract returns.
- Bila response group di branch target memiliki `code`, field `code` tidak hilang dan tidak null ketika source data ada.
- Regression tests ditambahkan/diupdate untuk controller parsing, service fallback/auth behavior, repository SQL year/COALESCE, dan group branch year filter.
- Tidak ada token/credential hardcoded di code, test, log, artifact, atau summary.

## Existing Patterns/Reuse

- Reuse `SecondaryReportSalesSumMonth` dan `SecondaryReportSalesGroup` sebagai controller parsing + validation entrypoints.
- Reuse DTO:
  - `entity.SecondarySalesReportDashboardSumPayload`
  - `entity.SecondarySalesReportDashboardGroupPayload`
- Reuse `resolveSecondaryDashboardYear(year *int) int` untuk fallback current year.
- Reuse `resolveSecondaryDashboardCustID(authCustID, parentCustID, requestedCustID)` untuk auth scope.
- Reuse repository CTE summary in `SecondarySalesReportSumReportByMonth` when facts are valid.
- Reuse `buildSecondarySalesReportGroupQuery(groupBy string)` for all group branch SQL.
- Reuse dry-run GORM SQL helpers in `sales/repository/report_repository_test.go`:
  - `newReportRepoDryRunDB`
  - `latestRecordedQuery`
- Reuse service mock patterns in `sales/service/report_service_test.go`.
- Reuse controller mock service setup in `sales/controller/so_controller_test.go`.

Discovery found much of the intended SX-2143 behavior already present in the local working tree: optional `year`, current-year fallback, `dt."year"` filters, summary order+return CTE, COALESCE numeric aggregates, auth cust resolution, and baseline tests. Executor must still implement any missing gaps and verify against the actual target branch/runtime before claiming Jira fixed.

## Constraints

- Follow repo-local rules in `AGENTS.md` and `.opencode/docs/*`.
- Validate from `sales` directory, not repo root.
- Use `rtk` prefix for shell workflows in this repo.
- Preserve Controller → Service → Repository → DB layering.
- Preserve tenant and scope rules: no controller-to-repository shortcuts, no arbitrary cust access, no hardcoded cust fallback.
- Schema prefixes matter: `report.`, `sls.`, `mst.`, `smc.`.
- Manual API test requires a valid token from secure local/environment source; never paste it into source/artifacts.
- Jira and Google Docs direct fetch were skipped because likely credentialed; prompt excerpts are treated as supplied reference.

## Risks

- Fallback `time.Now().Year()` is non-deterministic across calendar years. It resolves the current FE gap, but FE must send `&year=2026` for deterministic historical periods.
- Local code already includes changes related to SX-2143/SX-2172. If implementation branch differs, executor must compare and port only missing pieces.
- `sum-date` currently uses reporting facts (`report.fact_orders`/`report.fact_returns`) rather than live `sls.*` CTE from prompt. This is lower-risk if extract is correct, but if staging facts are missing for Juni 2026, executor may need to switch summary to source-table/date-range query or fix extraction.
- `GREATEST(os.last_update, rs.last_update)` may return null if either side is null in PostgreSQL. `last_update: null` is allowed by expected response, but if dashboard needs a non-null value when one side exists, use `COALESCE/GREATEST` pattern carefully.
- Current `SecondarySalesReportGroupResp` in local discovery lacks `code`. Acceptance mentions `code`; this may already be solved in another branch or must be added here.
- Manual DB/API verification depends on local/staging data and token availability.

## Decisions/Assumptions

- Source strategy: repo-local evidence plus user-supplied Jira/BE-doc excerpts. Official docs/context7, GitHub, web search, and browser evidence were intentionally skipped because this is a local Go/GORM SQL defect and external references are credentialed or not material.
- Assumption: `time.Now().Year()` fallback is acceptable as temporary compatibility behavior until SX-2201 FE sends `year`.
- Assumption: reporting facts are the preferred dashboard data source unless DB debug queries prove facts are stale/missing for `C260020001` Juni 2026.
- Assumption: `last_update` may remain `null` when no rows exist; numeric fields must not be null.
- No question gate was needed before writing this plan because the user supplied concrete endpoint, data period, expected mapping, and acceptance criteria. The only remaining decisions are implementation-time branch/runtime facts, covered by worklist tasks.

## Execution Source of Truth

Executor must follow this precedence:

1. Latest explicit user instruction.
2. Safety/security/permission rules, especially no secrets and cust scope.
3. Non-negotiable Implementation Invariants in this plan.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria and Done Criteria.
6. Implementation Steps.
7. Non-blocking notes and recommendations.

If sources conflict, follow the higher-priority source and record the conflict plus resolution in verification evidence.

## Non-negotiable Implementation Invariants

- This plan is artifact-only; source code is not changed yet by `@artifact-planner`.
- Keep Controller → Service → Repository → DB separation.
- `cust_id` must always be resolved through auth-aware service behavior; do not trust request `cust_id` blindly.
- Do not hardcode `C260020001`, `2026`, tokens, staging URLs, or user credentials in production code.
- Query behavior must be parameterized and SQL-injection safe.
- `year` filter must apply consistently to order and return portions of summary/group queries.
- Numeric response fields must be non-pointer or explicitly normalized to zero.
- `return_rate` must not divide by zero.
- Group return rows must subtract from net sales, not add.
- If `code` exists in group contract on target branch, implementation must preserve it through model/entity/service mapping and SQL alias.
- Any generated or debug logging must not expose Authorization headers or tokens and must not be production-noisy.

## Do Not / Reject If

- Reject if `sum-date` filters only `dt.month` without `dt."year"` or equivalent date range.
- Reject if `group` filters only one branch by year while order/return branch differs.
- Reject if missing rows produce JSON numeric `null`.
- Reject if code mutates/normalizes `cust_id` into the wrong ID such as adding an extra zero.
- Reject if distributor can request sibling/arbitrary `cust_id`.
- Reject if implementation hardcodes the Jira sample `C260020001` or year `2026` outside tests/smoke examples.
- Reject if FE contract fields are removed or renamed.
- Reject if implementation logs tokens, raw Authorization headers, or copied credentials.
- Reject if changes touch unrelated modules, env files, package files, lockfiles, migrations, or deployment config without explicit evidence.
- Reject if final claim says staging/Jira fixed without API/DB evidence or clearly labels it unverified.

## Diff Boundary

Allowed source/test changes:

- `sales/controller/report_controller.go`
- `sales/entity/report.go`
- `sales/model/report.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/controller/so_controller_test.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`
- `sales/client_test.http` only for non-secret sample requests.

Allowed evidence changes:

- `.opencode/evidence/20260610-1312-sx-2143-secondary-sales-dashboard-null/**`
- `.opencode/plans/20260610-1312-sx-2143-secondary-sales-dashboard-null.md`

Out-of-boundary changes must be reverted or justified in verification evidence before final quality gate.

## TDD/Test Plan

TDD required: yes. This is a production backend regression involving report SQL, auth-scoped cust filters, API parsing, and numeric contract.

Existing test patterns:

- Controller Fiber tests in `sales/controller/so_controller_test.go`.
- Service mock repository tests in `sales/service/report_service_test.go`.
- Repository dry-run SQL tests in `sales/repository/report_repository_test.go`.

Red step:

1. Add/confirm controller parsing tests:
   - `TestSecondaryReportSalesSumMonthParsesMonthYearCustID`
   - `TestSecondaryReportSalesSumMonthAllowsMissingYear`
   - `TestSecondaryReportSalesSumMonthRejectsInvalidMonthYear`
   - `TestSecondaryReportSalesGroupParsesMonthYearCustIDGroupBy`
2. Add/confirm service tests:
   - explicit `year=2026` forwarded to repository for summary and group.
   - missing year falls back to current year.
   - unauthorized distributor requested cust returns `ErrUnauthorizedCustID` and repository is not called.
   - `return_rate` stays zero when qty order zero.
   - if group `code` exists in target branch, service maps `Code` from model to entity.
3. Add/confirm repository dry-run tests:
   - summary SQL contains `dt."year" = $3` and `dt."year" = $6` or equivalent date range params.
   - summary SQL contains COALESCE for numeric aggregates.
   - summary SQL combines order + return and subtracts return gross/ppn/net sales appropriately.
   - group branch SQL for outlet/salesman/product_category/product contains year filters for both branches and return subtraction.
   - if group code contract exists, group SQL emits `code` alias.

Green step:

- Implement only missing behavior discovered by failing tests.
- Prefer extending existing helpers/functions instead of creating parallel paths.
- Keep repository SQL parameter order covered by tests.

Refactor step:

- If raw SQL becomes hard to read, extract narrowly named helper fragments without changing behavior.
- Remove temporary debug logs or gate them behind local-only debug mechanism before final.
- Keep comments English-only and only if they clarify non-obvious behavior.

Edge cases:

- `month=1` and `month=12` validation passes.
- missing `year` during calendar year 2026 uses `2026`.
- invalid `year=99`, `year=10000`, non-int year returns 400/unprocessable according to parser/validator behavior.
- no matching rows returns zero numeric values, not null.
- order rows exist but return rows do not; numeric order values still appear.
- return rows exist but order rows do not; numeric result remains deterministic and no panic/divide-by-zero.
- `cust_id=C260020001` remains exact.
- unknown `group_by` falls back to product branch as existing service behavior.

Commands:

```bash
rtk go test ./controller -run 'TestSecondaryReportSales'
rtk go test ./service -run 'TestSecondarySalesReport(SumReportByMonth|GroupSales)'
rtk go test ./repository -run 'TestSecondarySalesReport(SumReportByMonth|ReturnSumReportByMonth|Group)'
rtk go test ./service ./repository ./controller
rtk go test ./...
```

## Implementation Steps

1. Re-run baseline from `sales`:
   - `git status --short`
   - `rtk go test ./service ./repository ./controller`
2. Inspect current target branch for existing SX-2143/SX-2172 code because local discovery already shows many fixes present.
3. Add missing controller parsing tests in `sales/controller/so_controller_test.go`.
4. Add missing service tests in `sales/service/report_service_test.go` for explicit year on group if absent, fallback year, unauthorized cust, zero numeric/return-rate mapping, and group code mapping if code contract exists.
5. Add missing repository SQL tests in `sales/repository/report_repository_test.go` for `year`, COALESCE, order+return CTE, group all branches, and group `code` alias if target branch requires it.
6. Update `entity.SecondarySalesReportDashboardSumPayload` and `SecondarySalesReportDashboardGroupPayload` only if target branch lacks optional `Year *int` and validation tags.
7. Update controller only if QueryParser/validation does not already support `year`, `cust_id`, `group_by`, or invalid input behavior.
8. Update service only if target branch lacks:
   - `resolveSecondaryDashboardYear` fallback,
   - auth-aware `resolveSecondaryDashboardCustID`,
   - passing effective year into all repository calls,
   - zero-safe return rate,
   - group `code` mapping when required.
9. Update repository only if target branch lacks:
   - `year` parameter in interface/method signatures,
   - `dt."year"` filters or date range equivalent,
   - COALESCE numeric aggregates,
   - order+return combination for summary,
   - return subtraction for group,
   - `code` alias for group when required.
10. If DB smoke shows reporting facts empty but source `sls.*` has rows for Juni 2026, decide one of two paths and record evidence:
    - fix extraction/fact population if dashboard architecture requires facts,
    - or switch `sum-date` to source-table/date-range query from prompt if docs require live source calculation.
11. Run targeted tests, then full `rtk go test ./...` in `sales`.
12. If a valid secure token is available, perform manual API smoke without recording token:
    - `sum-date?month=6&year=2026&cust_id=C260020001`
    - `sum-date?month=6&cust_id=C260020001`
    - `group?month=6&year=2026&cust_id=C260020001&group_by=outlet`
13. Record changed files, commands, outputs, response sample with numeric fields, and any unverified staging limitations in evidence.
14. Route final review to `@quality-gate` because report SQL + tenant auth are security/data-scope sensitive.

## Expected Files to Change

Likely if branch is missing fixes:

- `sales/entity/report.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/controller/so_controller_test.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`

Possibly:

- `sales/model/report.go` if group `code` must be added/restored.
- `sales/controller/report_controller.go` only if validation/parser behavior must change.
- `sales/client_test.http` for non-secret example curl.

Not expected:

- `go.mod`, `go.sum`, migrations, env files, Docker/compose config, or other services.

## Agent/Tool Routing

- `@orchestrator`: coordinate implementation from this plan and keep evidence coherent.
- `@fixer`: bounded code/test edits in the `sales` module.
- `@explorer`: optional if target branch differs and needs more local discovery.
- `@oracle`: optional if deciding facts vs live `sls.*` query becomes architectural/risk-heavy.
- `@quality-gate`: required final signoff because this touches tenant-scoped reporting SQL and auth behavior.
- `@librarian`: not required unless executor needs credentialed docs extraction or updated external documentation not present in prompt.

## Executor Handoff Prompt

```text
Implement SX-2143 using `.opencode/plans/20260610-1312-sx-2143-secondary-sales-dashboard-null.md` as the source of truth. Work only in the `sales` module unless evidence justifies otherwise. Fix/verify `GET /sales/v1/reports/secondary-sales/sum-date` and `/group` so `month`, optional `year`, and `cust_id` are parsed, `year` is applied to summary and group queries, missing `year` falls back to `time.Now().Year()`, auth-scoped cust resolution is preserved, numeric aggregates are never JSON null, order and return are combined/subtracted correctly, and group branches preserve `code` if the target branch contract includes it. Use TDD: add failing controller/service/repository regression tests first, then implement only missing behavior. Do not hardcode token, credentials, `C260020001`, or year `2026` in production code. Do not touch env, package files, migrations, compose config, or unrelated modules. Validate from `sales` with targeted `rtk go test` commands and full `rtk go test ./...`. If a secure token is available, smoke test the three SX-2143 URLs without recording the token. Return root cause, changed files/functions, year handling, sanitized sample response, tests/commands, and FE note to use `&year=2026` not `?year2026`.
```

## Execution-ready Worklist / Handoff Contract

`start_with`: `T1`

### T1 — Baseline and branch diff discovery

- `depends_on`: none
- `owner/lane`: `@orchestrator` + `@explorer` if needed
- `action`: Verify current branch state, existing SX-2143 code, git status, and baseline tests.
- `validation`: `git status --short`; `rtk go test ./service ./repository ./controller`
- `exit_criteria`: Executor knows which required behaviors are already present vs missing on target branch.
- `blocking_status`: ready
- `blocker_reason`: none
- `requires_user_decision`: no
- `must_preserve`: no source changes during discovery except planned tests later.
- `do_not_touch`: env, secrets, unrelated modules.
- `evidence_update`: record branch/status and baseline test output under task evidence.
- `exit_verification`: baseline command output or explicit blocker.

### T2 — Add controller regression tests

- `depends_on`: T1
- `owner/lane`: `@fixer`
- `action`: Add/extend Fiber controller tests for successful parsing of `month/year/cust_id`, missing year, group `group_by`, and invalid month/year behavior.
- `validation`: `rtk go test ./controller -run 'TestSecondaryReportSales'`
- `exit_criteria`: Tests fail before implementation if parser/validation missing, otherwise pass and document behavior already present.
- `blocking_status`: ready
- `blocker_reason`: none
- `requires_user_decision`: no
- `must_preserve`: response payload structure and 403 unauthorized behavior.
- `do_not_touch`: repository SQL in this task.
- `evidence_update`: list added test names and expected statuses.
- `exit_verification`: targeted controller test output.

### T3 — Add/confirm service regression tests

- `depends_on`: T2
- `owner/lane`: `@fixer`
- `action`: Add/confirm service tests for explicit year, missing-year fallback, cust scope, zero return rate, summary mapping, all group branches, and group code mapping if applicable.
- `validation`: `rtk go test ./service -run 'TestSecondarySalesReport(SumReportByMonth|GroupSales)'`
- `exit_criteria`: Service tests cover effective cust/year behavior and pass after implementation.
- `blocking_status`: ready
- `blocker_reason`: none
- `requires_user_decision`: no
- `must_preserve`: `ErrUnauthorizedCustID` behavior and no repository call for unauthorized requested cust.
- `do_not_touch`: controller routes and repository SQL unless required by later task.
- `evidence_update`: record mock assertions for `custID`, `month`, `year`, and `groupBy`.
- `exit_verification`: targeted service test output.

### T4 — Add/confirm repository SQL regression tests

- `depends_on`: T3
- `owner/lane`: `@fixer`
- `action`: Add/confirm dry-run SQL tests for summary and group year filters, COALESCE numeric aggregates, order+return subtraction, param order, and group code alias if contract requires it.
- `validation`: `rtk go test ./repository -run 'TestSecondarySalesReport(SumReportByMonth|ReturnSumReportByMonth|Group)'`
- `exit_criteria`: SQL tests fail on missing behavior or pass when behavior is already present.
- `blocking_status`: ready
- `blocker_reason`: none
- `requires_user_decision`: no
- `must_preserve`: quoted `dt."year"` or date-range equivalent must be asserted for both order and return branches.
- `do_not_touch`: service/controller in this task.
- `evidence_update`: record SQL fragments asserted.
- `exit_verification`: targeted repository test output.

### T5 — Implement missing backend behavior

- `depends_on`: T2, T3, T4
- `owner/lane`: `@fixer`
- `action`: Update DTO/service/repository/model/entity code only where tests prove gaps exist.
- `validation`: run targeted controller, service, and repository test commands.
- `exit_criteria`: All regression tests from T2–T4 pass.
- `blocking_status`: ready
- `blocker_reason`: none
- `requires_user_decision`: no unless evidence proves facts-vs-live-source table choice changes architecture.
- `must_preserve`: tenant auth, numeric zero normalization, year fallback, return subtraction, no hardcoded secrets/data.
- `do_not_touch`: env, migrations, package files, unrelated endpoints.
- `evidence_update`: list changed functions and before/after behavior.
- `exit_verification`: targeted tests pass.

### T6 — Database/API smoke verification

- `depends_on`: T5
- `owner/lane`: `@orchestrator` or `@fixer`
- `action`: Run safe DB debug SQL and API smoke if local/staging token is securely available.
- `validation`: DB count queries for `report.dim_dates` and `report.fact_orders`; sanitized curl responses for target endpoints.
- `exit_criteria`: Either data exists and API returns non-null numbers, or missing token/data is documented as a verification blocker.
- `blocking_status`: ready
- `blocker_reason`: token/data may be unavailable but should not block code tests.
- `requires_user_decision`: no unless DB facts are empty while source `sls.*` has data and requires architecture choice.
- `must_preserve`: never record Authorization token or credentials.
- `do_not_touch`: production data and destructive SQL.
- `evidence_update`: store sanitized response snippets and DB counts.
- `exit_verification`: evidence file with no secrets.

### T7 — Full validation and quality gate

- `depends_on`: T6
- `owner/lane`: `@quality-gate`
- `action`: Run final tests and review diff boundaries/security/tenant behavior.
- `validation`: `rtk go test ./...` from `sales`; quality-gate review evidence.
- `exit_criteria`: Full tests pass or unrelated failures documented; no out-of-boundary changes; acceptance criteria mapped to evidence.
- `blocking_status`: ready
- `blocker_reason`: none
- `requires_user_decision`: no
- `must_preserve`: no final Jira/staging claim without evidence.
- `do_not_touch`: no new source edits during final review except fixes routed back to `@fixer`.
- `evidence_update`: final command output, diff summary, risk notes.
- `exit_verification`: quality-gate PASS or explicit remediation list.

## Validation Commands

From repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

From `sales`:

```bash
git status --short
rtk go test ./controller -run 'TestSecondaryReportSales'
rtk go test ./service -run 'TestSecondarySalesReport(SumReportByMonth|GroupSales)'
rtk go test ./repository -run 'TestSecondarySalesReport(SumReportByMonth|ReturnSumReportByMonth|Group)'
rtk go test ./service ./repository ./controller
rtk go test ./...
```

Safe DB debug SQL, only against approved local/staging DB:

```sql
SELECT id, year, month, day
FROM report.dim_dates
WHERE year = 2026 AND month = 6
ORDER BY day;

SELECT COUNT(*) AS total_rows,
       SUM(COALESCE(fo.net_sales_exclude_ppn, 0)) AS total_net_sales_exc_ppn,
       SUM(COALESCE(fo.gross_sale, 0)) AS total_gross_sale
FROM report.fact_orders fo
JOIN report.dim_dates dd ON dd.id = fo.date_id
WHERE fo.cust_id = 'C260020001'
  AND dd.year = 2026
  AND dd.month = 6;

SELECT fo.cust_id, COUNT(*)
FROM report.fact_orders fo
JOIN report.dim_dates dd ON dd.id = fo.date_id
WHERE dd.year = 2026
  AND dd.month = 6
  AND fo.cust_id LIKE 'C260020%'
GROUP BY fo.cust_id
ORDER BY fo.cust_id;
```

Manual API smoke with secure token only; do not store token:

```bash
curl '{{base_url}}/sales/v1/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer <VALID_STAGING_TOKEN>'

curl '{{base_url}}/sales/v1/reports/secondary-sales/sum-date?month=6&cust_id=C260020001' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer <VALID_STAGING_TOKEN>'

curl '{{base_url}}/sales/v1/reports/secondary-sales/group?month=6&year=2026&cust_id=C260020001&group_by=outlet' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer <VALID_STAGING_TOKEN>'
```

## Evidence Requirements

Required evidence before final claim:

- Changed files/functions summary.
- Test commands and outputs.
- Repository SQL behavior evidence for year filters and COALESCE.
- Auth/cust scope behavior evidence.
- Sanitized manual response sample or explicit note that token/data unavailable.
- DB debug count for `C260020001` Juni 2026 when DB access is available.
- Confirmation no token/credential was added to code/tests/artifacts.
- Quality-gate result.

Source strategy used for this plan:

- Local project discovery: used.
- User-supplied Jira/BE-doc excerpts: used.
- Official docs/context7: skipped; no unfamiliar/version-sensitive library behavior.
- GitHub: skipped; no upstream source dependency.
- Web search/Jira/Google Docs direct fetch: skipped due credentialed/private references and sufficient user excerpts.
- Browser/screenshot: skipped; backend API bug.

## Done Criteria

- All applicable acceptance criteria are backed by tests or smoke evidence.
- `rtk go test ./service ./repository ./controller` passes.
- `rtk go test ./...` passes, or unrelated failures are documented with evidence and accepted by quality gate.
- Manual API smoke confirms target endpoint returns non-null numeric fields, or lack of token/data is explicitly documented.
- Final response to user includes:
  1. Root cause found.
  2. Files/functions changed.
  3. Year handling when sent and missing.
  4. Sanitized sample response for `month=6&year=2026&cust_id=C260020001`.
  5. Tests run and results.
  6. FE note: use `&year=2026`, not `?year2026`.
- No secrets or credentials introduced.
- `@quality-gate` gives PASS or all remediation items are completed.

## Final Planning Summary

Artifacts created and kept:

- `.opencode/plans/20260610-1312-sx-2143-secondary-sales-dashboard-null.md` — primary implementation source of truth.
- `.opencode/evidence/20260610-1312-sx-2143-secondary-sales-dashboard-null/discovery.md` — kept as operational evidence because local discovery shows many SX-2143 behaviors already present and executor/quality-gate need replayable details.
- `.opencode/evidence/20260610-1312-sx-2143-secondary-sales-dashboard-null/index.json` — evidence manifest.

Artifacts deleted/cleaned:

- No draft artifacts were created, so no stale draft cleanup was needed.

Key decisions:

- Use Maintenance Stability Mode and smallest safe backend fix.
- Prefer existing report facts and current service/repository patterns unless DB evidence proves facts are stale/missing.
- Keep `time.Now().Year()` fallback for FE compatibility while requiring FE to send explicit `year` for deterministic behavior.
- Require TDD/regression-first execution because this is report SQL and tenant-scoped behavior.

Assumptions:

- Prompt excerpts are authoritative enough for Jira/BE-doc requirements.
- Local branch may already contain partial/full SX-2143 fixes; executor must confirm against target branch and runtime.
- `last_update` may be null; numeric fields may not.

Remaining open questions resolved by follow-up check:

- SX-2172 `code` enhancement is present in local code: entity/model include `Code`, service maps it, repository emits `AS code`, and tests assert it.
- Local `ggn_scyllax` facts are not populated for `C260020001` Juni 2026: `report.fact_orders = 0` and `report.fact_returns = 0`; however source `sls.*` has 15 valid orders, 18 order detail rows, and 1 return row for the same cust/date range. Details are kept in `.opencode/evidence/20260610-1312-sx-2143-secondary-sales-dashboard-null/db-check.md`.

Readiness:

- `ready-for-implementation` for a bounded first slice.
- `PASS_FOR_SLICE` because source tests can be implemented now; final staging/Jira closure still depends on secure token and DB/API smoke evidence.
