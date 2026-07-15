# Plan SX-2016 — Fix Monitoring Activity Principal salah/tidak lengkap location point

## Goal

Endpoint `GET /scylla-pjp/api/v1/live-monitoring-principal?date=1779278400&status[]=Approved&status[]=Need+Review&emp_id=482` mengembalikan semua destination route principal tanggal 2026-05-20 secara benar: satu destination satu row, activity mobile terhubung ke destination yang tepat, dan titik peta dashboard tidak hilang karena pagination row mentah/cross join.

## Non-goals

- Tidak ubah kontrak mobile `/mobile/visits/arrive` atau `/mobile/visits/leave`.
- Tidak ubah UI dashboard.
- Tidak ubah approval workflow PJP.
- Tidak memperbaiki data master koordinat placeholder `(1,2)` dalam kode.
- Tidak menyimpan token/password/secret ke test, log, atau artifact.

## Scope

Target module: `pjp`.

Perubahan utama:
- Perbaiki JOIN `pjp_principles.outlet_visit_list` di repository principal monitoring agar join berdasarkan `pjp_id`, `date`, dan `outlet_code/destination_code`, bukan hanya `pjp_code`.
- Isi `leave_at`, `arrive_longitude`, dan `arrive_latitude` untuk principal seperti distributor, dari OVL principal.
- Ubah pagination principal agar paging di level employee, bukan row mentah destination, supaya satu employee tidak kehilangan destination karena `LIMIT 10`.
- Tambah regression tests untuk transform/pagination dan repo SQL behavior sesuai kemampuan test existing.

Perubahan hardening opsional:
- Resolve child `cust_id` di principal service memakai `GetChildCustIDs` seperti distributor. DB case Princessa tidak butuh ini (`C26002` sama), tapi berguna untuk principal lain. Implement hanya jika tidak memperbesar scope/interface pain.

## Requirements

- Untuk `emp_id=482`, response harus berisi destination `BMI260003` dan `BMI260004` pada route 7010, bukan hanya duplikasi `BMI260005`.
- Join OVL harus tidak membuat cross product.
- `arrive_at/leave_at` harus menempel pada destination yang sama dengan `ovl.outlet_code = d.destination_code`.
- `arrive_longitude/arrive_latitude` untuk `BMI260003` harus bisa berisi mobile coordinate `-122.084000 / 37.421998` jika field response dipakai dashboard/QA.
- Default `limit=10` tidak boleh memotong destination list employee 482 menjadi 10 duplikasi dari destination pertama.
- `status[]=Approved&status[]=Need Review` tetap memfilter `pjp.approval_status`.
- Date conversion tetap `epochToDateString` dengan `Asia/Jakarta`.

## Acceptance Criteria

- Exact staging request return `data != null` untuk `emp_id=482` dan `route_data[].destination_data` berisi 5 destination unik route 7010.
- Response mengandung `destination_code=BMI260003` dengan `arrive_at=1779236763189`, `leave_at=1779237085474`, `arrive_latitude≈37.421998`, `arrive_longitude≈-122.084000`.
- Response mengandung `destination_code=BMI260004` dengan `arrive_at=1779268286992`.
- Tidak ada duplikasi 10 row `BMI260005` akibat cross product + SQL limit.
- User non-principal/distributor monitoring tidak regresi.

## Existing Patterns/Reuse

- Reuse distributor pattern for per-employee pagination: first get scoped employee IDs, slice page, then fetch all detail rows for paged employees with `limit=0`.
- Reuse distributor response fields `ArriveLongitude`, `ArriveLatitude`, `LeaveAt` in `response.LiveMonitoringDestinationData`.
- Reuse `buildLiveMonitoringDayRange`/timezone concepts only if needed; principal OVL already filters date via `ovl.date = DATE(rpp.date)`.
- Reuse Go service tests in `pjp/service/live_monitoring/get_distributor_service_test.go` as stub style.

## Constraints

- Layer Controller → Service → Repository → DB.
- `pjp` is separate Go module; validate from `/pjp`.
- Shell commands stay `rtk`-prefixed.
- Do not commit/store secrets/JWT/DB password.
- DB validation must be read-only unless explicitly requested.

## Risks

- Master destination coordinates for route 7010 mostly placeholder `(1,2)`; code fix cannot make those map points geographically correct.
- Dashboard may plot only destination `longitude/latitude`; after query fix it will show 5 destinations, but 4 may appear near `(1,2)` until master data fixed.
- If frontend expects old duplicate structure, UI behavior may change; this is desired but should be QA-smoked.
- Child `cust_id` hardening might change principal scope for other users; use distributor pattern with fallback to reduce risk.

## Decisions/Assumptions

- DB validation overrides earlier hypothesis: Princessa login `cust_id=C26002`, `parent_cust_id=C26002`; salesman `MS9990/emp_id=482` also `cust_id=C26002`. Child scope mismatch is not SX-2016 root cause for this data.
- Root cause: repository principal query cross joins destinations with all OVL rows via `pjp_code`, then `LIMIT 10` truncates to duplicate first destination.
- Dashboard currently reads destination coordinates; add arrive coordinates for completeness/parity but not as primary map contract unless FE later confirms.
- Open question: should data team fix `pjp_principles.destinations` coordinates currently set to `(1,2)` for `BMI260003`, `BMI260004`, `BMI260005`, `BMI260015`?

## TDD/Test Plan

TDD required: ya. Ini bug production query/tenant activity mapping.

Existing test patterns:
- `get_distributor_service_test.go` has stub repo, per-employee pagination test, and transform tests.
- Principal currently lacks enough tests around duplicate/cross-join behavior.

First failing/regression tests:

1. `TestTransformPrincipalRows_DoesNotDuplicateCrossJoinedDestinationsAfterRepoFix`
   - Input rows should model 5 distinct destination rows with OVL matched per destination.
   - Assert output has 5 `destination_data` entries and includes `BMI260003` with arrive/leave + arrive coordinate.
   - This guards transform after adding `LeaveAt/ArriveLongitude/ArriveLatitude`.

2. `TestGetPrincipalMonitoring_PaginatesEmployeeScopeBeforeDetailRows`
   - Stub repo returns scoped employee IDs `[482]` or count `1`; service should fetch rows for employee 482 with `limit=0` detail query (or equivalent no SQL detail limit).
   - Assert destination list is not truncated by default `limit=10`.
   - If adding `GetPrincipalEmployeeIDs` to repo interface, mirror distributor design.

3. Repository SQL test if sqlmock/pattern exists; if not, rely on service tests + DB smoke evidence:
   - Verify JOIN condition changes to `ovl.pjp_id = pjp.id`, `ovl.date = DATE(rpp.date)`, `ovl.outlet_code = d.destination_code`.

Green step:
- Update repository SELECT and JOIN.
- Add fields to `model.LiveMonitoringPrincipalRow`: `LeaveAt`, `ArriveLongitude`, `ArriveLatitude`.
- Update `transformPrincipalRows` to set `LeaveAt`, `ArriveLongitude`, `ArriveLatitude`.
- Update principal service pagination to page employees before fetching all detail rows, or remove row-level limit if only one employee page but preserve paging metadata.

Refactor step:
- Consider shared helper for pagination/scope only after tests pass; avoid broad abstraction.

Edge cases:
- OVL row absent → destination still appears with nil activity fields.
- Multiple OVL historical rows same pjp_code different date → no longer joined due date condition.
- OVL pjp_code duplicates across days → no longer cross-product.
- Empty emp filter → employee paging still works.

Commands:
- `rtk go test ./service/live_monitoring -run 'Test(TransformPrincipalRows|GetPrincipalMonitoring)'`
- `rtk go test ./service/live_monitoring`
- `rtk go test ./...`

## Implementation Steps

1. Add principal regression tests for OVL fields + no truncation.
2. Update model fields in `pjp/model/live_monitoring.go`.
3. Update repository principal SELECT/JOIN in `pjp/repository/live_monitoring/get_principal_repository.go`:
   - join OVL by `pjp_id`, `date`, `outlet_code = destination_code`.
   - select `ovl.leave_at`, `arrive_longitude`, `arrive_latitude`.
4. Update `transformPrincipalRows` fields.
5. Fix principal pagination:
   - preferred: add repository method `GetPrincipalEmployeeIDs` mirroring distributor and page employees before detail rows.
   - acceptable smaller fix: do not apply SQL `Limit/Offset` to raw detail rows; keep count/paging by employee. Use only if product accepts full detail rows for page employees.
6. Run tests.
7. Smoke endpoint staging with same account/token after deploy.
8. `@quality-gate` final review.

## Expected Files to Change

Primary:
- `pjp/repository/live_monitoring/get_principal_repository.go`
- `pjp/model/live_monitoring.go`
- `pjp/service/live_monitoring/get_principal_service.go`
- `pjp/service/live_monitoring/get_principal_service_test.go`

Possible if interface pagination method added:
- `pjp/repository/live_monitoring/live_monitoring_repository.go`
- existing test stubs in `pjp/service/live_monitoring/*_test.go`

No planned changes:
- mobile arrive/leave handlers.
- dashboard frontend.
- migrations.
- package files.

## Agent/Tool Routing

- `@orchestrator`: execute plan and integrate validation.
- `@fixer`: code tests + repository/service fix.
- `@explorer`: DB smoke and trace if endpoint still wrong.
- `@quality-gate`: final review; tenant/query regression risk.
- `@architect`: not needed.
- External docs/GitHub/web/browser: not needed; local + DB evidence sufficient.

## Execution-ready Worklist / Handoff Contract

`start_with`: `SX2016-01`

| id | action | depends_on | owner/lane | validation/check | exit criteria | status | blocker | requires_user_decision |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| SX2016-01 | Add regression tests for principal destination activity fields and no row truncation. | none | `@fixer` | `rtk go test ./service/live_monitoring -run 'Test(TransformPrincipalRows|GetPrincipalMonitoring)'` from `pjp` | Tests fail before fix or document if direct green. | ready | none | no |
| SX2016-02 | Add `LeaveAt`, `ArriveLongitude`, `ArriveLatitude` to principal row model and transform. | SX2016-01 | `@fixer` | same targeted command | Principal response can carry OVL activity coordinate fields. | ready | none | no |
| SX2016-03 | Fix principal repository OVL JOIN to match pjp/date/destination. | SX2016-02 | `@fixer` | targeted service tests + manual SQL evidence if available | No cross product; rows model one destination ↔ one OVL. | ready | none | no |
| SX2016-04 | Fix principal pagination to avoid raw row truncation. | SX2016-03 | `@fixer` | `rtk go test ./service/live_monitoring` | Default limit does not drop destination rows for employee page. | ready | none | no |
| SX2016-05 | Run full PJP tests. | SX2016-04 | `@fixer` | `rtk go test ./...` from `pjp` | Pass or unrelated failures documented. | ready | none | no |
| SX2016-06 | Staging smoke after deploy with Princessa login. | SX2016-05 | `@orchestrator` | curl exact endpoint; inspect `BMI260003`, `BMI260004`, unique destination count | Endpoint returns 5 unique destinations and mobile activity fields. | blocked | Need deployed fixed build. | yes |
| SX2016-07 | Quality gate review. | SX2016-05 | `@quality-gate` | review diff/test/smoke evidence | No blocking tenant/query regression. | ready | none | no |

## Validation Commands

From repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

From `pjp` directory:

```bash
rtk go test ./service/live_monitoring -run 'Test(TransformPrincipalRows|GetPrincipalMonitoring)'
rtk go test ./service/live_monitoring
rtk go test ./...
```

Staging smoke after fixed deploy, with fresh token only:

```bash
curl "https://best.scyllax.online/scylla-pjp/api/v1/live-monitoring-principal?date=1779278400&status[]=Approved&status[]=Need+Review&emp_id=482" \
  -H "Authorization: Bearer <token>" \
  -H "Accept: application/json"
```

Expected response checks:
- `data[0].emp_id == 482`
- destination codes include `162612`, `BMI260003`, `BMI260004`, `BMI260005`, `BMI260015`
- no 10x duplicate `BMI260005`
- `BMI260003.arrive_at == 1779236763189`
- `BMI260003.leave_at == 1779237085474`
- `BMI260003.arrive_longitude == -122.084000`
- `BMI260003.arrive_latitude == 37.421998`

## Evidence Requirements

Kept evidence:
- `.opencode/evidence/20260520-1702-sx-2016-monitoring-principal-empty/discovery.md`

Implementation evidence needed:
- test output targeted + full `pjp` module.
- diff summary.
- staging smoke after deploy.

Research gate:
- Local project discovery: done.
- DB read-only validation: done.
- Official docs/context7: skipped; no external library behavior.
- GitHub/web: skipped; not relevant.
- Browser screenshot: skipped; backend endpoint sufficient. UI QA can run after endpoint fixed.

## Done Criteria

- Principal repository no longer cross-joins OVL rows.
- Principal response includes all unique route destinations for employee page.
- Principal response includes arrive/leave and actual mobile coordinate fields for matched destination.
- `rtk go test ./service/live_monitoring` and `rtk go test ./...` pass or unrelated failures documented.
- Staging endpoint after deploy returns expected `BMI260003`/`BMI260004` data.
- `@quality-gate` passes.

## Final Planning Summary

Artifacts created/updated:
- `.opencode/plans/20260520-1702-sx-2016-monitoring-principal-empty.md` — source of truth implementation handoff.
- `.opencode/evidence/20260520-1702-sx-2016-monitoring-principal-empty/discovery.md` — kept because it contains DB/endpoint proof and exact root cause.

Key decisions:
- Root cause is query JOIN + row-level pagination, not mobile persistence.
- Fix by joining OVL on `pjp_id + date + outlet_code/destination_code` and by avoiding SQL row-limit truncation of destination details.

Assumptions resolved:
- Account and DB validation completed via user-provided credentials; secrets not persisted.
- Data mobile exists in `pjp_principles.outlet_visit_list`.
- Endpoint currently returns data but wrong/duplicated destination rows.

Remaining open questions:
- Whether data team should fix placeholder destination coordinates `(1,2)`.
- Whether FE wants actual arrive marker via `arrive_longitude/arrive_latitude` or master destination marker only.

Readiness:
- Ready for implementation via `@fixer`, starting `SX2016-01`.
- Staging smoke is blocked until fixed build deployed.

Cleanup:
- No draft artifacts created. Evidence kept intentionally.
