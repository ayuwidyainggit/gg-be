# Verification

Date: 2026-06-18
Task: 1:1 port `POST /sales/v1/reports/secondary-sales` from `dev` to `bugfix/SX-2209-qa`
Service: `sales`
Branch: `bugfix/SX-2209-qa`

## Source strategy

Literal port from `dev` using direct file checkout for report files:
- `controller/report_controller.go`
- `entity/report.go`
- `service/report_service.go`
- `repository/report_repository.go`
- `model/report.go`
- `pkg/constant/constant.go`
- `service/report_service_test.go`
- `repository/report_repository_test.go`

Command used:
```bash
git checkout dev -- controller/report_controller.go entity/report.go service/report_service.go repository/report_repository.go model/report.go pkg/constant/constant.go service/report_service_test.go repository/report_repository_test.go
```

## Scope-preserving deviations

Needed compile-only branch reconciliation outside requested endpoint logic:
- `controller/so_controller_test.go`
  - updated mock `SecondarySalesReportTrendSales` signature to match dev service interface (`[]string`)
- `repository/report_repository.go`
  - adapted activity-report helper calls to existing `activity_report_query.go` signatures on `bugfix/SX-2209-qa`
  - added `buildActivitySalesmanGroupSalesSQL` and `buildActivitySalesmanGroupReturnSQL` so repo compiles with dev-ported activity report methods

Reason:
- direct dev port introduced compile mismatch against branch-local activity report query helpers
- no secondary-sales endpoint behavior changed by these reconciliations

## Validation

Commands run:
```bash
rtk docker compose -f docker-compose.yml ps
rtk go test ./controller -run 'TestSecondary|TestActivity'
rtk go test ./service -run 'Test.*SecondarySales|Test.*Activity'
rtk go test ./repository -run 'Test.*SecondarySales|Test.*Activity'
rtk go test ./...
```

Results:
- `rtk docker compose -f docker-compose.yml ps` -> warning only about obsolete compose `version` key; command ran
- `rtk go test ./controller -run 'TestSecondary|TestActivity'` -> pass (`10 passed in 1 packages`)
- `rtk go test ./service -run 'Test.*SecondarySales|Test.*Activity'` -> pass (`28 passed in 1 packages`)
- `rtk go test ./repository -run 'Test.*SecondarySales|Test.*Activity'` -> pass (`31 passed in 1 packages`)
- `rtk go test ./...` -> pass (`222 passed in 22 packages`)

## Notes

Expected dev behavior now present for endpoint path:
- request body `cust_id` accepts string or array via dev raw-body parsing path
- auth-owned fields still sourced from JWT locals, not body spoofing
- publish/export payload flow follows dev implementation for this endpoint
