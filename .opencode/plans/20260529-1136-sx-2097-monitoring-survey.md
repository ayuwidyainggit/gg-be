# SX-2097 Monitoring Detail Survey Section Plan

## Goal
Tambahkan `survey_data` ke response `GET /scylla-pjp/v1/monitoring_locations/details` agar FE bisa render survey yang sudah di-submit salesman pada detail live monitoring.

## Non-goals
- Tidak ubah kontrak section existing: `visit_information`, `sales`, `return`, `collection`, `expense`, `shipment`.
- Tidak tambah filter `leave_at` karena query referensi Jira/docs tidak memakainya.
- Tidak ubah status matching menjadi case-insensitive tanpa bukti data atau keputusan PO/QA.
- Tidak ubah route, auth, atau request parameter.

## Scope
- Module target: `pjp`.
- Endpoint target: `GET /v1/monitoring_locations/details`.
- Tambah model row, response DTO, repository method, service mapping, dan test.
- Survey source: `mst.survey_answer` join `mst.m_survey` join `mst.m_outlet`.

## Requirements
- `survey_data` muncul sejajar dengan section lain di setiap item `data`.
- Empty state survey harus `[]`, bukan `null`, saat detail data ada.
- Setiap survey row berisi:
  - `submission`
  - `survey_title`
  - `outlet_code`
  - `outlet_name`
- Filter wajib: request `emp_id`, request `date`, `sa.status = 'Submitted'`.
- Date compare pakai date semantics: `DATE(sa.answer_date) = ?` atau ekuivalen bound param.
- SQL/query wajib pakai parameter binding.
- Tenant/customer isolation wajib mengikuti pola existing bila schema survey punya `cust_id`.

## Acceptance Criteria
- Endpoint `GET /scylla-pjp/v1/monitoring_locations/details` return `survey_data` di `data[0]`.
- `survey_data` grouped by `survey_title + outlet_code + outlet_name`.
- `submission` berisi count `survey_answer_id`.
- `outlet_code` map ke `mst.m_outlet.outlet_code`.
- `outlet_name` map ke `mst.m_outlet.outlet_name`.
- Hanya status exact `Submitted` dihitung.
- Empty survey return `survey_data: []`.
- Existing response tetap backward compatible.
- Tests ditambah/diupdate dan lulus.

## Existing Patterns/Reuse
- Route sudah ada di `pjp/router/live_monitoring.go`.
- Controller sudah ada di `pjp/controller/live_monitoring/get_detail_controller.go`.
- Service flow sudah ada di `pjp/service/live_monitoring/get_detail_service.go`.
- Repository interface di `pjp/repository/live_monitoring/live_monitoring_repository.go`.
- Query detail section di `pjp/repository/live_monitoring/get_detail_repository.go`.
- DTO response di `pjp/data/response/live_monitoring_response.go`.
- Raw row model di `pjp/model/live_monitoring.go`.
- Service tests memakai `detailRepoStub` di `pjp/service/live_monitoring/get_detail_service_test.go`.
- Repository tests memakai `sqlmock` di `pjp/repository/live_monitoring/get_detail_repository_test.go`.
- Reuse lebih baik daripada buat helper baru: ikuti pola `GetSales`, `GetReturns`, `GetShipments`, dan `GetCollections`.

## Constraints
- Layering wajib: Controller → Service → Repository → DB.
- `pjp` adalah Go/Gin module; validasi dari direktori `pjp`.
- Repo policy: gunakan `rtk` untuk shell workflow di repo ini.
- Jangan commit secrets, token, `.env`, atau copy token dari docs.
- Jangan tambah filter bisnis yang tidak ada di Jira/docs.
- Kalau tenant column tersedia, jangan buka data lintas customer.

## Risks
- `mst.survey_answer` schema belum dikonfirmasi di local discovery; perlu validasi column sebelum final query tenant filter.
- Kalau `survey_answer` tidak punya `cust_id`, isolation hanya lewat `emp_id` dan join outlet; perlu catat evidence.
- `COUNT` Postgres bisa scan ke integer berbeda; model perlu tipe aman, ideal `int64` atau `int` sesuai scan test.
- Menambah method repository interface memaksa update semua test stubs.
- Manual staging test butuh token/data yang tidak tersedia di prompt.

## Decisions/Assumptions
- Pertanyaan tidak diajukan; requirement cukup jelas untuk rencana implementasi.
- Asumsi: `survey_data` hanya ditampilkan bila detail monitoring punya `visit_information`; jika service return no-data existing, response tetap `data: null`.
- Asumsi: `submission` singular adalah kontrak response, meski SQL awal docs menyebut `submissions`.
- Keputusan: pakai exact `sa.status = 'Submitted'`.
- Keputusan: tidak pakai `leave_at` filter.
- Keputusan: implement tenant filter bila schema mendukung, dengan prioritas `sa.cust_id IN ?` dan join `mo.cust_id = sa.cust_id`; jika tidak tersedia, dokumentasikan hasil validasi dan gunakan query referensi.

## TDD/Test Plan
- TDD required: ya, karena production logic dan API response contract berubah.
- Existing test patterns:
  - Service unit: `pjp/service/live_monitoring/get_detail_service_test.go` dengan `detailRepoStub`.
  - Repository unit: `pjp/repository/live_monitoring/get_detail_repository_test.go` dengan `sqlmock`.
- Red step pertama:
  - Tambah field `surveys []model.SurveyDataRow` dan capture args di `detailRepoStub`.
  - Tambah `TestGetMonitoringDetail_IncludesSurveyData` yang gagal karena `LiveMonitoringDetailData` belum punya `SurveyData` dan service belum call repo.
- Green step:
  - Tambah `SurveyDataRow`, `SurveyData`, `GetSubmittedSurveyData`, service mapping, dan response field.
- Refactor step:
  - Rapikan mapping loop bila perlu, tanpa ubah behavior existing.
- Edge cases:
  - No survey rows return `survey_data: []`.
  - Only `Submitted` counted by repository query.
  - Date and emp ID passed unchanged from request.
  - `outlet_code` and `outlet_name` not swapped.
  - Existing sections remain empty/non-empty exactly as before.
- Commands:
  - `rtk go test ./service/live_monitoring -run 'TestGetMonitoringDetail_(IncludesSurveyData|NoSurveyReturnsEmptyList|SurveyUsesRequestDateAndEmpID)'`
  - `rtk go test ./repository/live_monitoring -run TestGetSubmittedSurveyData`
  - `rtk go test ./...`

## Implementation Steps
1. Update `pjp/model/live_monitoring.go`:
   - add `SurveyDataRow` with GORM columns `submission`, `survey_title`, `outlet_code`, `outlet_name`.
2. Update `pjp/data/response/live_monitoring_response.go`:
   - add `SurveyData []SurveyData 'json:"survey_data"'` to `LiveMonitoringDetailData`.
   - add `SurveyData` struct with JSON keys required by FE.
3. Update `pjp/repository/live_monitoring/live_monitoring_repository.go`:
   - add `GetSubmittedSurveyData(ctx context.Context, tx *gorm.DB, custIDs []string, date string, empID int) ([]model.SurveyDataRow, error)`.
4. Update repository implementation in `pjp/repository/live_monitoring/get_detail_repository.go`:
   - query `mst.survey_answer sa` join `mst.m_survey ms` and `mst.m_outlet mo`.
   - use bound args for `custIDs`, `date`, `empID`, `Submitted`.
   - group and order by `ms.survey_title`, `mo.outlet_code`, `mo.outlet_name`.
   - include tenant filter if schema supports `sa.cust_id`; otherwise match reference query and record evidence.
5. Update `pjp/service/live_monitoring/get_detail_service.go`:
   - call repository after `targetCustIDs` exists.
   - map rows to `[]response.SurveyData` with `make(..., 0, len(rows))` to preserve empty array.
   - set `SurveyData: surveyData` in result.
6. Update service test stubs:
   - add new interface method to `detailRepoStub`, `principalRepoStub`, `distributorRepoStub` if compile needs it.
7. Add service tests:
   - includes non-empty survey data.
   - no survey returns empty list.
   - repo receives `targetCustIDs`, request `date`, request `empID`.
8. Add repository test:
   - `TestGetSubmittedSurveyData_ReturnsSubmittedSurveyGroupedByTitleAndOutlet` using `sqlmock`.
   - assert query includes `status = ?`, `DATE(sa.answer_date) = ?`, group/order fields, and returns mapped rows.
9. Run targeted tests, then full `pjp` module tests.
10. Optional manual DB/API validation when credentials/data available.

## Expected Files to Change
- `pjp/model/live_monitoring.go`
- `pjp/data/response/live_monitoring_response.go`
- `pjp/repository/live_monitoring/live_monitoring_repository.go`
- `pjp/repository/live_monitoring/get_detail_repository.go`
- `pjp/service/live_monitoring/get_detail_service.go`
- `pjp/service/live_monitoring/get_detail_service_test.go`
- `pjp/repository/live_monitoring/get_detail_repository_test.go`
- Potentially `pjp/service/live_monitoring/get_principal_service_test.go`
- Potentially `pjp/service/live_monitoring/get_distributor_service_test.go`

## Agent/Tool Routing
- `@orchestrator`: start implementation handoff and coordinate validation.
- `@fixer`: code/test implementation.
- `@explorer`: only if schema or compile stubs unclear during implementation.
- `@quality-gate`: final review because API contract and tenant isolation risk changed.
- No `@designer`, browser, GitHub, or external docs needed.

## Execution-ready Worklist / Handoff Contract
`start_with`: `SX2097-01`

| Task | Action | depends_on | Owner | Validation | Exit Criteria | Status | requires_user_decision |
|---|---|---|---|---|---|---|---|
| SX2097-01 | Run targeted grep/read if implementation context stale; confirm files listed above still current. | none | `@explorer` or `@fixer` | `rg "GetMonitoringDetail|LiveMonitoringDetailData|GetShipments|GetCollections" pjp` | Current files confirmed. | ready | no |
| SX2097-02 | Add failing service tests for `survey_data` mapping, empty array, and arg propagation. | SX2097-01 | `@fixer` | `rtk go test ./service/live_monitoring -run TestGetMonitoringDetail_.*Survey` | Tests fail for missing implementation, not syntax-only unrelated issue. | ready | no |
| SX2097-03 | Add `SurveyDataRow` model and `SurveyData` response DTO; add `SurveyData` field to `LiveMonitoringDetailData`. | SX2097-02 | `@fixer` | `rtk go test ./service/live_monitoring -run TestGetMonitoringDetail_.*Survey` | Compile reaches missing repo method/service call phase or tests progress. | ready | no |
| SX2097-04 | Add `GetSubmittedSurveyData` to repository interface and update all stubs to compile. | SX2097-03 | `@fixer` | `rtk go test ./service/live_monitoring -run TestGetMonitoringDetail_.*Survey` | Interface compile errors resolved. | ready | no |
| SX2097-05 | Implement repository query with bound params, grouping, ordering, and tenant filter if schema supports it. | SX2097-04 | `@fixer` | `rtk go test ./repository/live_monitoring -run TestGetSubmittedSurveyData` | Query maps `submission`, `survey_title`, `outlet_code`, `outlet_name`; SQL mock passes. | ready | no |
| SX2097-06 | Wire service call and map survey rows into response with empty array default. | SX2097-05 | `@fixer` | `rtk go test ./service/live_monitoring -run TestGetMonitoringDetail_.*Survey` | Service tests pass; `survey_data` always non-nil when detail exists. | ready | no |
| SX2097-07 | Run broader live monitoring tests. | SX2097-06 | `@fixer` | `rtk go test ./service/live_monitoring ./repository/live_monitoring` | Target packages pass. | ready | no |
| SX2097-08 | Run full pjp tests. | SX2097-07 | `@fixer` | `rtk go test ./...` | Full module passes or unrelated failures documented with evidence. | ready | no |
| SX2097-09 | Validate local DB schema and survey data in `ggn_scyllax`. | SX2097-08 | `@fixer` | `psql -h localhost -p 5432 -U postgres -d ggn_scyllax` with safe local password env | Schema columns exist; SQL reference returns rows or verified empty state for chosen `emp_id/date`. | blocked: needs local DB running and password env | yes |
| SX2097-10 | Validate Docker service with direct cURL against local `pjp` container. | SX2097-09 | `@fixer` | `curl -sS http://localhost:9010/api/v1/monitoring_locations/details?...` | Response includes `data[0].survey_data` array matching DB SQL count/grouping. | blocked: needs `AUTH_TOKEN` and local data | yes |
| SX2097-11 | Final quality review for API contract, tenant safety, DB evidence, Docker cURL evidence, and tests. | SX2097-08 | `@quality-gate` | Review diff and validation output | Approved or specific fixes listed. | ready | no |

## Validation Commands
Run tests from `pjp/`:

```bash
rtk go test ./service/live_monitoring -run 'TestGetMonitoringDetail_.*Survey'
rtk go test ./repository/live_monitoring -run TestGetSubmittedSurveyData
rtk go test ./service/live_monitoring ./repository/live_monitoring
rtk go test ./...
```

Run Docker service from repo root:

```bash
rtk docker compose -f docker-compose.yml ps pjp
rtk docker compose -f docker-compose.yml up -d pjp
rtk docker compose -f docker-compose.yml logs --no-color --since=2m pjp
```

Direct service health check, no gateway prefix:

```bash
curl -sS 'http://localhost:9010/api/v1/health'
```

Set safe local variables before API validation. Do not paste token into plan, shell history screenshots, or Jira:

```bash
export AUTH_TOKEN='<local-or-staging-token>'
export EMP_ID='210'
export ACTIVITY_DATE='2026-05-28'
export DISTRIBUTOR_ID='67'
```

Validate endpoint through Docker-published PJP service:

```bash
curl -sS --location -g "http://localhost:9010/api/v1/monitoring_locations/details?emp_id=${EMP_ID}&date=${ACTIVITY_DATE}&distributor_id=${DISTRIBUTOR_ID}" \
  --header 'Accept: application/json' \
  --header "Authorization: Bearer ${AUTH_TOKEN}" \
  -o /tmp/sx-2097-monitoring-detail.json

jq '.message, (.data[0].survey_data // "missing")' /tmp/sx-2097-monitoring-detail.json
jq -e '.data[0] | has("survey_data") and (.survey_data | type == "array")' /tmp/sx-2097-monitoring-detail.json
```

Validate local DB target is `ggn_scyllax` on local Postgres:

```bash
export PGPASSWORD='<local-postgres-password>'
psql -h localhost -p 5432 -U postgres -d ggn_scyllax -v ON_ERROR_STOP=1 -c 'SELECT current_database(), current_user, inet_server_addr(), inet_server_port();'
```

Validate survey schema columns in local `ggn_scyllax`:

```bash
psql -h localhost -p 5432 -U postgres -d ggn_scyllax -v ON_ERROR_STOP=1 -c "
SELECT table_schema, table_name, column_name, data_type
FROM information_schema.columns
WHERE table_schema = 'mst'
  AND table_name IN ('survey_answer', 'm_survey', 'm_outlet')
  AND column_name IN ('survey_answer_id', 'survey_id', 'outlet_id', 'emp_id', 'answer_date', 'status', 'survey_title', 'outlet_code', 'outlet_name', 'cust_id')
ORDER BY table_name, ordinal_position;
"
```

Validate local DB source data for selected `EMP_ID` and `ACTIVITY_DATE`:

```bash
psql -h localhost -p 5432 -U postgres -d ggn_scyllax -v ON_ERROR_STOP=1 \
  -v emp_id="${EMP_ID}" -v activity_date="${ACTIVITY_DATE}" -c "
SELECT
    sa.emp_id,
    sa.answer_date::date AS answer_date,
    sa.status,
    COUNT(*) AS total
FROM mst.survey_answer sa
WHERE sa.emp_id = :'emp_id'::int
  AND sa.answer_date::date = :'activity_date'::date
GROUP BY sa.emp_id, sa.answer_date::date, sa.status
ORDER BY sa.status;
"
```

Validate local DB expected `survey_data` rows:

```bash
psql -h localhost -p 5432 -U postgres -d ggn_scyllax -v ON_ERROR_STOP=1 \
  -v emp_id="${EMP_ID}" -v activity_date="${ACTIVITY_DATE}" -c "
SELECT
    COUNT(sa.survey_answer_id) AS submission,
    ms.survey_title,
    mo.outlet_code,
    mo.outlet_name
FROM mst.survey_answer sa
JOIN mst.m_survey ms
    ON ms.survey_id = sa.survey_id
JOIN mst.m_outlet mo
    ON mo.outlet_id = sa.outlet_id
WHERE sa.answer_date::date = :'activity_date'::date
  AND sa.emp_id = :'emp_id'::int
  AND sa.status = 'Submitted'
GROUP BY ms.survey_title, mo.outlet_code, mo.outlet_name
ORDER BY ms.survey_title ASC, mo.outlet_code ASC;
"
```

Compare API response with DB result:

```bash
jq -r '.data[0].survey_data[]? | [.submission, .survey_title, .outlet_code, .outlet_name] | @tsv' /tmp/sx-2097-monitoring-detail.json
```

If DB query returns no rows, accepted API result is:

```bash
jq -e '.data[0].survey_data == []' /tmp/sx-2097-monitoring-detail.json
```

## Evidence Requirements
- Test output for targeted service tests.
- Test output for repository test.
- Full `rtk go test ./...` output or documented unrelated failures.
- Docker runtime evidence: `rtk docker compose -f docker-compose.yml ps pjp`, recent `pjp` logs, and `GET /api/v1/health` result.
- Local DB evidence from `ggn_scyllax`: current database check, schema column check, status count query, expected `survey_data` query.
- cURL evidence against `http://localhost:9010/api/v1/monitoring_locations/details` with token redacted and response snippet showing `survey_data` array.
- SQL/schema evidence for tenant filter decision if DB available.
- Never store `AUTH_TOKEN`, `.env`, raw secrets, or full Authorization header in artifacts/Jira.

## Done Criteria
- `survey_data` exists in success detail response.
- Empty survey is `[]`.
- Query uses bound params.
- Tenant filter decision documented and safe.
- Tests cover service mapping and repository query.
- Existing sections remain unchanged.
- `@quality-gate` signs off or all review blockers fixed.

## Final Planning Summary
- Artifacts created:
  - `.opencode/evidence/20260529-1136-sx-2097-monitoring-survey/discovery.md`
  - `.opencode/evidence/20260529-1136-sx-2097-monitoring-survey/index.json`
  - `.opencode/plans/20260529-1136-sx-2097-monitoring-survey.md`
- Key decisions:
  - Use `submission` singular in response.
  - Map `outlet_code` from `mo.outlet_code` and `outlet_name` from `mo.outlet_name` despite docs typo.
  - Use exact `Submitted` status.
  - Do not add `leave_at` filter.
  - Apply tenant filter if schema supports it; otherwise document schema evidence.
  - Validate Docker service directly via `http://localhost:9010/api/v1`, because root compose publishes `pjp` on port `9010` and router group is `/api/v1`.
  - Validate database against local Postgres DB `ggn_scyllax` before comparing API response.
- Questions asked: none.
- Open questions: none blocking implementation; Docker cURL validation needs `AUTH_TOKEN`; DB validation needs local Postgres access.
- Readiness: ready for `@orchestrator` → `@fixer` implementation, then Docker/API/DB validation and `@quality-gate`.
- Cleanup: no durable draft needed; evidence kept because it is useful for implementation handoff and audit.
