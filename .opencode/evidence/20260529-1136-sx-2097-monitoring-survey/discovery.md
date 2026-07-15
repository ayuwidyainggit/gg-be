# Discovery SX-2097 Monitoring Survey

## Files inspected
- `AGENTS.md`
- `.opencode/docs/index.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `docs/Monitoring Activity - BE.md`
- `pjp/router/live_monitoring.go`
- `pjp/controller/live_monitoring/get_detail_controller.go`
- `pjp/service/live_monitoring/get_detail_service.go`
- `pjp/service/live_monitoring/get_detail_service_test.go`
- `pjp/repository/live_monitoring/live_monitoring_repository.go`
- `pjp/repository/live_monitoring/get_detail_repository.go`
- `pjp/repository/live_monitoring/get_detail_repository_test.go`
- `pjp/data/request/live_monitoring_request.go`
- `pjp/data/response/live_monitoring_response.go`
- `pjp/model/live_monitoring.go`

## Project patterns found
- Service target: `pjp`, Gin-based Go module.
- Endpoint route: `GET /v1/monitoring_locations/details` in `pjp/router/live_monitoring.go` maps to `controller.GetMonitoringDetail`.
- Controller parses `LiveMonitoringDetailRequest`, gets `customerID` and `userID`, calls `service.GetMonitoringDetail`, wraps non-nil result as one item array.
- Service flow in `pjp/service/live_monitoring/get_detail_service.go`:
  1. determine principal/distributor from `distributor_id` nil check,
  2. fetch visit info,
  3. return nil when no visit info,
  4. resolve `salesmanCustID` using `GetSalesmanCustID`,
  5. fetch `sales`, `return`, `expense`, `shipment`, `collection`,
  6. assemble `response.LiveMonitoringDetailData`.
- Repository interface is centralized in `pjp/repository/live_monitoring/live_monitoring_repository.go`.
- Detail queries live in `pjp/repository/live_monitoring/get_detail_repository.go`.
- Existing data sections use GORM query builder or raw SQL with bound `?` args, not string interpolation.
- Current response DTO lacks `survey_data` in `LiveMonitoringDetailData`.
- Current model lacks `SurveyDataRow`.
- Current service test uses `detailRepoStub`; any repository interface method addition requires updating this stub plus principal/distributor stubs if they implement full interface.
- Repository unit tests use `sqlmock` in `pjp/repository/live_monitoring/get_detail_repository_test.go`.

## Reuse candidates
- Reuse service assembler pattern from `Sales`, `Return`, `Expense`, and `Collection`.
- Reuse tenant isolation approach: after visit info, resolve `salesmanCustID`, then pass `targetCustIDs := []string{salesmanCustID}` to data-section queries.
- Reuse GORM query builder style from `GetSales`, `GetReturns`, and `GetShipments` for survey query.
- Reuse `sqlmock` repository test pattern from `TestGetCollections_AllocatesCollectionPerInvoicePerOutlet`.
- Reuse `detailRepoStub` service test pattern to assert request date and emp ID are passed into new repository method.

## Docs checked
- `docs/Monitoring Activity - BE.md` lines 664-722 confirm `survey_data` addition and query source.
- Docs table has swapped descriptions for `outlet_code` and `outlet_name`; SQL/query intent says `outlet_code = mo.outlet_code`, `outlet_name = mo.outlet_name`.
- Existing doc response at lines 780-782 lacks `survey_data`, matching Jira failure.

## Constraints
- Keep Controller → Service → Repository layering.
- Keep existing response fields unchanged.
- `survey_data` must always be JSON array when detail data exists; no `null`, no omitted field.
- Query must use bound params.
- Filter by request `emp_id`, request `date`, and exact `status = 'Submitted'` unless data evidence forces product decision.
- Date compare must use date-cast semantics like existing `DATE(...) = ?`; for `answer_date` use `DATE(sa.answer_date) = ?` or equivalent bound expression.
- Tenant safety likely needs `sa.cust_id` and/or `mo.cust_id` filter if columns exist and existing section pattern requires customer isolation. Must confirm DB columns or code schema before final implementation.
- Do not add `leave_at` filter because reference query does not include it.

## Risks
- `mst.survey_answer` may have `cust_id`; omitting tenant filter could leak cross-customer data when emp IDs collide.
- `COUNT` scan type into `int` may need explicit model type compatible with GORM/Postgres count.
- `status` casing may differ in actual data; changing to `ILIKE` without confirmation could break intended contract.
- Adding repository interface method requires updating all service test stubs implementing `LiveMonitoringRepository`.
- Manual staging verification needs real token and data; not available in prompt.

## Research gate result
- Local project discovery: done and required.
- Official docs/context7: not needed; standard Go/GORM/Postgres patterns already present locally.
- GitHub: not needed; no upstream behavior dependency.
- Brave/web search: not needed; Jira/docs and repo enough.
- Browser/screenshot: not needed; BE endpoint change only.

## Suggested first regression test
`TestGetMonitoringDetail_IncludesSurveyData` in `pjp/service/live_monitoring/get_detail_service_test.go`, expecting `survey_data` with one row and verifying repo received `custIDs`, `date`, `empID`.
