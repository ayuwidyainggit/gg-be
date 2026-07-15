# Plan — SX-2182 Secondary Sales Multiselect BE

Task id: `20260608-1534-sx-2182-secondary-sales-multiselect`
Readiness: `ready-for-implementation`
Plan Quality Gate: `PASS`
Source of truth: `.opencode/plans/20260608-1534-sx-2182-secondary-sales-multiselect.md`

## Goal

Implementasi BE untuk SX-2182 agar Secondary Sales Report dan endpoint master Business Unit mendukung multi-select Business Unit/Region/Area/Distributor secara backward-compatible, aman tenant-scope, dan tetap mengikuti pola repo `Controller -> Service -> Repository -> DB`.

## Non-goals

- Tidak mengubah modul `pjp-sales/` kecuali user meminta parity terpisah.
- Tidak mengubah kontrak route endpoint.
- Tidak mengubah mekanisme RMQ/export async selain payload/filter `cust_id` yang perlu backward-compatible.
- Tidak membuat FE behavior, UI checkbox, atau Select All.
- Tidak memakai remote dev DB default atau menambah kredensial/env.
- Tidak mengubah formula bisnis lain di luar requirement `sum-date`, `group`, `trend-sales`, export, dan business-unit.

## Scope

Endpoint target:

- `GET /master/v1/business-unit`
- `POST /sales/v1/reports/secondary-sales`
- `GET /sales/v1/reports/secondary-sales/trend-sales`
- `GET /sales/v1/reports/secondary-sales/sum-date`
- `GET /sales/v1/reports/secondary-sales/group`

Module target:

- `master/` untuk business-unit `region_id`/`area_id` query parsing dan repository filter.
- `sales/` untuk Secondary Sales export/dashboard DTO, auth scope resolver, repository filters, response fields, dan tests.

## Requirements

1. `POST /sales/v1/reports/secondary-sales` menerima `cust_id` lama sebagai string dan baru sebagai array string.
2. Query export order dan return mengambil semua selected authorized `cust_id` memakai parameter binding slice, bukan string interpolation.
3. `GET /master/v1/business-unit` menerima `region_id` dan `area_id` single, comma-separated, comma with spaces, repeated, dan bracket syntax bila existing client memakai `[]`.
4. `trend-sales`, `sum-date`, dan `group` menerima `cust_id` single/multi dari query; `trend-sales` tetap boleh membaca GET body jika code existing sudah support, tetapi query param adalah jalur utama.
5. Missing/empty `cust_id` memakai default auth cust/effective allowed scope sesuai behavior lama; `cust_id: []` tidak boleh menjadi `IN ()` atau `ANY('{}')` yang menghasilkan no rows tanpa sengaja.
6. Principal hanya boleh request auth cust atau child cust di bawah `parent_cust_id`.
7. Distributor user tidak boleh request sibling/foreign cust; reject request, jangan silent ignore.
8. `sum-date` memakai `year` filter dan response fields existing yang sudah mengarah ke docs tetap dipertahankan.
9. `group` menambahkan `code` pada response item dan memperbaiki `name` product/category bila data source tersedia.
10. Single-select behavior lama harus tetap sama.

## Acceptance Criteria

- `cust_id: "C260020001"` tetap sukses untuk export dan query efektif `[]string{"C260020001"}`.
- `cust_id: ["C260020001"]` sukses dengan response export shape existing.
- `cust_id: ["C260020001", "C260020002"]` sukses hanya bila user berhak atas semua cust tersebut.
- Unauthorized requested cust mengembalikan 403 untuk endpoint sales dan tidak menjalankan query report/export.
- `region_id=80`, `region_id=80,90`, `region_id=80,%2090`, dan equivalent `area_id` tidak error dan menghasilkan filter multi yang benar.
- Invalid numeric list seperti `region_id=80,abc` menghasilkan 400 atau error validation eksplisit; jangan silent skip token invalid untuk behavior baru.
- Export multi-cust mencakup data semua selected cust pada order dan return branch.
- `trend-sales` menghasilkan 12 bulan dan multi-cust sama dengan agregasi manual beberapa single-cust.
- `sum-date` multi-cust sama dengan agregasi manual beberapa single-cust dan tetap aman divide-by-zero untuk `return_rate`.
- `group` response item punya `id`, `code`, `name`, `net_sales`.
- Query list memakai `IN ?`/slice binding GORM atau `sqlx.In`; tidak ada list ids yang disusun dari raw request string.
- Unit/repository/controller tests menutup backward compatibility dan multi payload/query.

## Existing Patterns/Reuse

- Reuse `master/controller/business_unit_controller.go` helper concept `normalizeIntArrayQuery`; sekarang sudah support comma, whitespace, repeated, `region_id[]`, dan `region_id`.
- Reuse `master/repository/business_unit_repository.go` `sqlx.In` pattern untuk `md.region_id IN (?)` dan `md.area_id IN (?)`.
- Reuse `sales/service/report_service.go` `ErrUnauthorizedCustID` dan HTTP 403 mapping existing di controller.
- Extend `resolveSecondaryDashboardCustID` menjadi resolver multi-cust di service layer, jangan pindahkan auth decision ke repository/controller.
- Reuse GORM `IN ?` slice binding pattern yang sudah dipakai untuk `DistributorIDs`, `SalesmanIDs`, `OutletIDs`, dan `ProIDs` di export query builder.
- Reuse `sales/service/report_service_test.go` mock repository hook pattern.
- Prior plan `.opencode/plans/20260519-1250-secondary-sales-dashboard-filters.md` sudah menyelesaikan sebagian single `year`/`cust_id` dashboard behavior; implementasi SX-2182 harus memperluas ke multi, bukan menghapus hardening lama.

## Constraints

- Shell workflow di repo ini harus `rtk`-prefixed.
- Validasi dijalankan dari module folder target: `master/` dan `sales/`.
- Wajib jaga tenant rules dari `.opencode/docs/ARCHITECTURE.md`: transactional `cust_id`, parent scope lewat `parent_cust_id`.
- Planner tidak mengedit source; implementasi dilakukan oleh `@orchestrator`/`@fixer` setelah plan.
- Jangan membuat query raw dengan interpolasi list dari request.
- Jangan menyentuh `.env`, secrets, dump, atau infra credential.

## Risks

- Broken access control jika multi `cust_id` tidak di-intersect/reject terhadap allowed scope.
- Async RMQ payload bisa gagal decode jika `cust_id` berubah tipe tanpa custom compatibility.
- Empty slice bisa memfilter semua data menjadi kosong jika tidak dinormalisasi menjadi default scope.
- Group `id` bisa tabrakan lintas cust bila dim table id tidak global; perlu verify atau group by code/name/cust-aware key.
- Business-unit helper existing silently skips invalid number; requirement SX-2182 lebih aman bila invalid token jadi 400.
- `trend-sales` sekarang hanya order fact, bukan return; bila docs mengharapkan order+return, catat sebagai follow-up bila tidak ada source requirement eksplisit.

## Decisions/Assumptions

Keputusan plan:

- Untuk sales `cust_id` unauthorized, reject seluruh request dengan 403, bukan ignore sebagian selected cust.
- Untuk business-unit invalid numeric token, ubah menjadi 400 supaya client tahu filter salah.
- Untuk empty/missing `cust_id`, pakai behavior lama: auth cust sebagai default. Jangan otomatis pakai semua child cust kecuali product owner mengonfirmasi Select All/missing memang berarti semua allowed child.
- Untuk principal multi child, validasi semua requested cust dengan parent scope; efisienkan dengan query count/list multi daripada loop N kali bila memungkinkan.
- Untuk query sales, gunakan `IN ?` GORM slice binding karena pattern itu sudah ada di repository.

Assumptions / open questions:

- Assumption slice-safe: `authCustID == parentCustID` menandai principal, `authCustID != parentCustID` menandai distributor, sesuai hardening existing.
- Assumption slice-safe: `cust_id: []` dari FE berarti tidak memilih filter khusus dan fallback ke default auth cust, bukan semua child cust.
- Open non-blocking: bila PO ingin `cust_id` missing/empty untuk principal berarti semua allowed child cust, service resolver harus disesuaikan sebelum implementasi final.
- Open non-blocking: group dim ids perlu diverifikasi apakah global unik; bila tidak, group multi-cust harus tetap menghasilkan agregasi yang benar dengan grouping yang tidak mencampur entitas berbeda.

## Execution Source of Truth

Precedence untuk executor:

1. Instruksi user terbaru yang eksplisit.
2. Safety/security/permission rules repo dan tenant isolation.
3. Non-negotiable Implementation Invariants di plan ini.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria dan Done Criteria.
6. Implementation Steps.
7. Follow-up/rekomendasi non-blocking.

Jika ada konflik, ikuti sumber yang lebih tinggi dan catat konflik serta alasan di evidence implementasi.

## Non-negotiable Implementation Invariants

- Planner posture: jangan menganggap file plan sebagai implementasi; source edits hanya boleh dilakukan lane implementasi setelah plan ini.
- Tenant invariant: tidak boleh ada selected `cust_id` yang tidak authorized masuk query report/export.
- Backward compatibility invariant: string `cust_id` lama, single query `cust_id`, single `region_id`, dan single `area_id` tetap berjalan.
- Binding invariant: list ids harus lewat parameter binding (`IN ?`, `sqlx.In`, atau equivalent), bukan concat raw query dari input.
- Default invariant: empty optional filter arrays berarti no optional filter; empty/missing `cust_id` sales fallback ke auth cust kecuali keputusan user baru mengubahnya.
- Layer invariant: controller parse/validate request, service resolve authorization/effective filters, repository query DB saja.
- Async invariant: RMQ payload lama dan payload baru sama-sama bisa diproses oleh subscriber.
- Evidence invariant: implementer wajib menyimpan hasil tests dan manual query/cURL penting di summary/evidence.

## Do Not / Reject If

- Reject if implementation menghapus `CustID`/`ParentCustID` auth-only protection dan membiarkan body spoof auth owner.
- Reject if unauthorized selected `cust_id` hanya diabaikan diam-diam.
- Reject if query memakai `WHERE cust_id = 'C1,C2'` atau interpolasi string list.
- Reject if `cust_id: []` menghasilkan SQL invalid atau no rows tanpa keputusan eksplisit.
- Reject if perubahan memecah report list owner; `report.list.cust_id` harus tetap auth owner supaya principal melihat export job-nya.
- Reject if `group` menambahkan `code` tapi menghilangkan `name`/`net_sales` existing.
- Reject if implementer mengubah `pjp-sales/` tanpa justifikasi parity yang diminta user.
- Reject if tests hanya happy path single value dan tidak menutup multi/backward compatibility.

## Diff Boundary

Allowed source groups:

- `master/controller/business_unit_controller.go`
- `master/entity/business_unit.go` bila perlu tipe/error DTO
- `master/repository/business_unit_repository.go`
- `master/service/business_unit_service.go` hanya bila scope validation perlu
- `master/controller/*business_unit*_test.go`
- `master/repository/*business_unit*_test.go`
- `master/service/*business_unit*_test.go`
- `sales/entity/report.go`
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/model/report.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`
- `sales/controller/report_controller_test.go` bila dibutuhkan untuk body/query parsing tests

Evidence/report exceptions:

- `.opencode/evidence/20260608-1534-sx-2182-secondary-sales-multiselect/**`
- `.opencode/plans/20260608-1534-sx-2182-secondary-sales-multiselect.md`

Out-of-boundary changes harus direvert atau dijustifikasi di verification evidence sebelum final quality gate.

## TDD/Test Plan

TDD required: ya.

Reason:

- Ini perubahan production API behavior, export async payload, DB query, dan tenant isolation.

Existing test patterns:

- `master/controller/business_unit_controller_test.go` untuk parsing query controller.
- `master/repository/business_unit_repository_test.go` untuk SQL string builder.
- `sales/service/report_service_test.go` untuk service auth resolver, RMQ publish payload, export subscriber behavior.
- Tambahkan `sales/repository/report_repository_test.go` untuk SQL builder/report query checks bila belum ada.

First failing/regression tests:

1. `master` controller: `region_id=80,90&area_id=89,90` dan `region_id=80,%2090&area_id=89,%2090` menghasilkan `[]int{80,90}` dan `[]int{89,90}`.
2. `master` controller: `region_id=80,abc` menghasilkan 400 atau explicit validation error.
3. `sales` controller/body parser: export body menerima `cust_id` string lama dan array baru.
4. `sales` service: principal multi child valid menghasilkan effective `CustIDs []string{"C1","C2"}` dan report owner tetap auth cust.
5. `sales` service: distributor request sibling multi menghasilkan `ErrUnauthorizedCustID` dan tidak publish/query.
6. `sales` subscriber: payload lama single `cust_id` dan payload baru `cust_ids`/array sama-sama diteruskan ke repository dengan effective slice benar.
7. `sales` repository export: order dan return where clause memakai `IN ?` untuk cust ids dan params memuat slice cust ids.
8. `sales` repository sum/group/trend: all dashboard queries memakai `IN ?` cust ids dan `dt."year" = ?` saat applicable.
9. `sales` group response: maps `Code` ke JSON `code` untuk outlet/salesman/product_category/product.

Green step:

- Tambahkan normalizer/custom type list input, multi resolver auth, DTO fields, repository interface signature updates, dan query binding sampai tests lulus.

Refactor step:

- Extract helper kecil bila duplikasi tinggi:
  - `normalizeStringListInput`
  - `normalizeQueryStringList`
  - `resolveSecondaryDashboardCustIDs`
  - `validateBusinessUnitIntArrayQuery`
- Pastikan naming tidak ambigu: `RequestedCustID` untuk raw legacy single, `RequestedCustIDs`/`CustIDs` untuk normalized effective list.

Edge cases:

- `cust_id` missing, empty string, empty array.
- `cust_id` comma query with spaces.
- Duplicate cust ids deduped preserving stable order.
- Principal request auth cust plus child cust.
- Principal request one valid and one invalid child: reject all.
- Distributor request auth cust only: allow.
- Optional filters `distributor_ids/outlet_ids/salesman_ids/pro_ids` empty arrays: no filter.
- Business-unit repeated query values and bracket/non-bracket syntax.

Commands:

Dari repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

Dari `master/`:

```bash
rtk go test ./controller -run BusinessUnit
rtk go test ./repository -run BusinessUnit
rtk go test ./service -run BusinessUnit
rtk go test ./...
```

Dari `sales/`:

```bash
rtk go test ./service -run 'SecondarySalesReport|PublishSecondarySalesReport|SubscribeSecondarySalesReport'
rtk go test ./repository -run 'SecondarySales'
rtk go test ./controller -run 'SecondarySales'
rtk go test ./...
```

## Implementation Steps

1. **Harden business-unit query parsing**
   - Update `normalizeIntArrayQuery` agar bisa mengembalikan error untuk invalid numeric token.
   - Controller harus map invalid token ke 400.
   - Pertahankan support key `region_id[]`, `region_id`, `area_id[]`, `area_id`.

2. **Add/extend list input types for sales**
   - Ubah export DTO private agar `cust_id` bisa decode string atau array string.
   - Pertimbangkan custom type dengan `UnmarshalJSON` agar `cust_id: "C1"`, `cust_id: ["C1"]`, dan `cust_id: []` aman.
   - Untuk GET query dashboard, parse `cust_id` dari raw query args supaya comma-separated dan repeated bisa dinormalisasi.

3. **Extend sales entity/request model**
   - Tambahkan normalized `RequestedCustIDs []string` atau `CustIDs []string` di `SecondarySalesReportQueryFilter`.
   - Jaga legacy `RequestedCustID string` bila dibutuhkan untuk payload lama/backward compatibility.
   - Tambahkan `Code string json:"code"` di `SecondarySalesReportGroupResp` dan model scan `SecondarySalesReportGroup`.

4. **Implement multi auth resolver in service**
   - Tambahkan `resolveSecondaryDashboardCustIDs(authCustID, parentCustID string, requested []string) ([]string, error)`.
   - Empty requested -> `[]string{authCustID}`.
   - Requested same as auth -> allow.
   - Distributor with requested not only auth -> `ErrUnauthorizedCustID`.
   - Principal requested child list -> validate all belong to `parentCustID`; prefer repository method multi count/list.

5. **Update export publish/subscriber flow**
   - Publish: normalize body `cust_id`, resolve effective cust ids, store `report.list.cust_id` as auth owner, serialize effective `CustIDs` for RMQ.
   - Subscriber: support old payload (`RequestedCustID`/`CustID`) and new payload (`CustIDs`/array). Never revert multi to single.

6. **Update sales repository signatures and queries**
   - Change export builder from single `CustID` to `CustIDs` slice.
   - Replace `od.cust_id = ?` and `rd.cust_id = ?` with `od.cust_id IN ?` and `rd.cust_id IN ?` or equivalent GORM slice binding.
   - Update sum/group/trend repository methods from `custID string` to `custIDs []string`.
   - For group query select `id, code, name, net_sales` and group by `id, code, name`; verify source column names from dim tables.

7. **Update dashboard service/controller flow**
   - `sum-date`, `group`, `trend-sales` normalize `cust_id` query/body to slice and call multi resolver.
   - Keep existing `year` validation/fallback.
   - Map `Code` in group response.

8. **Add tests and run validation**
   - Add Red tests first, implement Green, then refactor helper naming.
   - Run targeted and full module tests.
   - If runtime API is available, run manual cURL cases from user prompt and compare single vs multi aggregation.

## Expected Files to Change

Likely:

- `master/controller/business_unit_controller.go`
- `master/controller/business_unit_controller_test.go`
- `master/repository/business_unit_repository_test.go`
- `sales/entity/report.go`
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/model/report.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`

Possible:

- `master/entity/business_unit.go`
- `master/service/business_unit_service.go`
- `sales/controller/report_controller_test.go`

## Agent/Tool Routing

- `@artifact-planner`: artifact plan dan evidence saja.
- `@orchestrator`: jalankan implementasi berdasarkan plan, pecah `master` dan `sales` work bila perlu.
- `@fixer`: bounded source edits dan tests di `master/` dan `sales/`.
- `@explorer`: optional follow-up bila implementer perlu mencari dim table `code` column atau hidden tests.
- `@quality-gate`: wajib final review karena tenant/auth/report export sensitive.

## Executor Handoff Prompt

Copyable prompt untuk implementation lane:

```text
Implement SX-2182 using .opencode/plans/20260608-1534-sx-2182-secondary-sales-multiselect.md as source of truth.

Scope:
- master GET /v1/business-unit supports region_id/area_id single, comma, comma+space, repeated, bracket/non-bracket; invalid numeric token must return explicit 400.
- sales Secondary Sales export/dashboard supports cust_id string legacy and multi string array/comma query.
- sales endpoints: POST /v1/reports/secondary-sales, GET /secondary-sales/trend-sales, /sum-date, /group.

must_preserve:
- Controller -> Service -> Repository -> DB layering.
- report.list.cust_id remains auth owner for export jobs.
- unauthorized selected cust_id rejects request; do not silently ignore.
- empty/missing cust_id falls back to auth cust unless user gives newer instruction.
- list filters use parameter binding, not raw interpolation.
- RMQ subscriber handles old single-cust payload and new multi-cust payload.

do_not_touch:
- Do not edit pjp-sales, .env, secrets, infra credentials, unrelated services, or docs outside requested evidence unless justified.

validation:
- root: rtk docker compose -f docker-compose.yml ps
- master: rtk go test ./controller -run BusinessUnit; rtk go test ./repository -run BusinessUnit; rtk go test ./...
- sales: rtk go test ./service -run 'SecondarySalesReport|PublishSecondarySalesReport|SubscribeSecondarySalesReport'; rtk go test ./repository -run 'SecondarySales'; rtk go test ./...

return/evidence:
- Changed files.
- Tests run and outputs.
- Manual cURL/runtime checks if available.
- Explicit note on unauthorized cust behavior, empty cust behavior, query binding evidence, and group code/name source.
- Any out-of-boundary changes or unresolved assumptions.
```

## Execution-ready Worklist / Handoff Contract

`start_with`: `T1`

| Task | Action | depends_on | Owner/lane | Validation/check | Exit criteria | Status | requires_user_decision | must_preserve | do_not_touch | evidence_update | exit_verification |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T1 | Add failing tests for `master` business-unit non-bracket comma, whitespace, and invalid token parsing | none | `@fixer` | `rtk go test ./controller -run BusinessUnit` from `master/` | Tests fail before implementation or document already passing cases | ready | no | existing bracket/repeated support | unrelated master endpoints | note test names | failing/passing test output |
| T2 | Implement strict multi int query parsing for business-unit | T1 | `@fixer` | `rtk go test ./controller -run BusinessUnit` | valid forms parse; invalid token 400 | ready | no | single value compatibility | repository logic unless needed | parser behavior note | passing test output |
| T3 | Verify/adjust business-unit repository multi filter tests | T2 | `@fixer` | `rtk go test ./repository -run BusinessUnit` | `IN (?)` + `sqlx.In` remains used | ready | no | parameterized binding | raw SQL interpolation | SQL assertion note | passing test output |
| T4 | Add sales tests for export `cust_id` string/array/empty and unauthorized multi | none | `@fixer` | `rtk go test ./service -run 'PublishSecondarySalesReport|SubscribeSecondarySalesReport'` | Tests capture legacy + new payload behavior | ready | no | auth owner and ErrUnauthorizedCustID | repository SQL first | test names | failing/passing test output |
| T5 | Implement sales list input DTO/custom JSON/query normalizers | T4 | `@fixer` | targeted sales service/controller tests | string, array, comma query normalize with dedupe | ready | no | legacy payload shape | unrelated report DTOs | normalizer cases | passing test output |
| T6 | Implement multi `cust_id` resolver and repository scope validation | T5 | `@fixer` | `rtk go test ./service -run SecondarySalesReport` | effective cust ids resolved safely; unauthorized rejected | ready | no | no silent ignore | controller direct DB calls | auth behavior note | passing test output |
| T7 | Update export RMQ publish/subscriber and union query builder for multi cust | T6 | `@fixer` | `rtk go test ./service -run 'PublishSecondarySalesReport|SubscribeSecondarySalesReport'`; `rtk go test ./repository -run SecondarySales` | old and new payloads work; export SQL uses bound slice | ready | no | report owner auth cust | Excel column shape | RMQ payload + SQL evidence | passing test output |
| T8 | Update `sum-date`, `group`, and `trend-sales` repository/service/controller paths for multi cust | T6 | `@fixer` | `rtk go test ./service -run SecondarySalesReport`; `rtk go test ./repository -run SecondarySales` | dashboard queries aggregate selected cust ids | ready | no | year filter and fallback behavior | unrelated activity report | query/aggregation note | passing test output |
| T9 | Add `code` to group model/response and verify source columns | T8 | `@fixer` with optional `@explorer` | group repository/service tests | response includes `id/code/name/net_sales`; name not empty when source exists | ready | no | existing fields remain | response envelope | dim code source note | passing test output |
| T10 | Run full module validation and prepare quality-gate evidence | T1-T9 | `@orchestrator` then `@quality-gate` | `rtk go test ./...` in `master/` and `sales/` | all relevant tests pass or blockers documented | ready | no | evidence over assertion | source outside boundary | final validation summary | quality-gate handoff |

## Validation Commands

Root:

```bash
rtk docker compose -f docker-compose.yml ps
```

Master module:

```bash
rtk go test ./controller -run BusinessUnit
rtk go test ./repository -run BusinessUnit
rtk go test ./service -run BusinessUnit
rtk go test ./...
```

Sales module:

```bash
rtk go test ./service -run 'SecondarySalesReport|PublishSecondarySalesReport|SubscribeSecondarySalesReport'
rtk go test ./repository -run 'SecondarySales'
rtk go test ./controller -run 'SecondarySales'
rtk go test ./...
```

Manual smoke if runtime/token available:

```bash
curl '/master/v1/business-unit?is_active=1&region_id=80,90&area_id=89,90&q=&page=1&limit=99'
curl '/sales/v1/reports/secondary-sales/trend-sales?year=2026&cust_id=C260020001,C260020002'
curl '/sales/v1/reports/secondary-sales/sum-date?month=5&year=2026&cust_id=C260020001,C260020002'
curl '/sales/v1/reports/secondary-sales/group?month=4&year=2026&cust_id=C260020001,C260020002&group_by=outlet'
```

## Evidence Requirements

Implementation evidence must include:

- Changed files list.
- Test commands and outputs for `master` and `sales`.
- SQL/builder evidence showing list filters use binding (`IN ?`, `sqlx.In`, or equivalent).
- Auth behavior evidence for principal valid multi child and distributor unauthorized sibling.
- RMQ payload compatibility note for old single and new multi payload.
- Group `code` source mapping evidence for each `group_by`.
- Manual runtime/cURL evidence if credentials/runtime are available; otherwise note skipped reason.

Source strategy used for this plan:

- Local project discovery: used.
- Jira/issue details from user prompt: used as product reference.
- Official docs/context7: skipped because implementation depends on repo-local Go/GORM/sqlx patterns and Jira requirements, not external API behavior.
- GitHub/web search: skipped because no upstream behavior needed.
- Browser/runtime: skipped at planning phase because no implementation or runtime token was available.

Kept evidence artifact:

- `.opencode/evidence/20260608-1534-sx-2182-secondary-sales-multiselect/discovery.md`
- `.opencode/evidence/20260608-1534-sx-2182-secondary-sales-multiselect/index.json`

## Done Criteria

- All acceptance criteria pass.
- `master/` and `sales/` targeted tests pass.
- `master/` and `sales/` full `rtk go test ./...` pass or any unrelated pre-existing failures are documented with evidence.
- No source changes outside diff boundary unless justified.
- No secrets/env/dumps touched.
- Final `@quality-gate` signoff completed because this touches tenant/auth/report data.

## Final Planning Summary

Artifacts created:

- Primary plan: `.opencode/plans/20260608-1534-sx-2182-secondary-sales-multiselect.md`
- Evidence: `.opencode/evidence/20260608-1534-sx-2182-secondary-sales-multiselect/discovery.md`
- Evidence manifest: `.opencode/evidence/20260608-1534-sx-2182-secondary-sales-multiselect/index.json`

Artifacts consulted:

- Repo docs: `.opencode/docs/index.md`, `AGENT_ROUTING.md`, `ARCHITECTURE.md`, `QUALITY.md`.
- Source files in `sales/` and `master/` listed in discovery evidence.
- Prior plans for Secondary Sales export/dashboard hardening.

Key decisions:

- Reject unauthorized multi `cust_id` rather than silently filtering.
- Keep missing/empty `cust_id` fallback to auth cust for backward compatibility.
- Use repo-local binding patterns: GORM `IN ?` and `sqlx.In`.
- Keep evidence folder because it remains operationally useful for implementation handoff.

Assumptions:

- Principal/distributor detection remains `authCustID == parentCustID` vs `authCustID != parentCustID`.
- Empty selected `cust_id` means default auth cust unless user clarifies otherwise.

Open questions:

- Non-blocking: PO may decide principal empty `cust_id` means all allowed child cust. If so, resolver default must change before final implementation.
- Non-blocking: verify dim `id` uniqueness for group multi-cust during implementation.

Cleanup performed:

- No draft artifacts were created.
- Evidence artifacts were kept intentionally for handoff; none deleted as stale.

Readiness:

- `ready-for-implementation` / `PASS` for bounded implementation using this plan.
