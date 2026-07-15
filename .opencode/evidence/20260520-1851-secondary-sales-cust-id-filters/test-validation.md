# Test Validation — Task 120/121 Secondary Sales `cust_id` Filters

Task id: `20260520-1851-secondary-sales-cust-id-filters`
Tanggal: `2026-05-20`
Module: `sales`

## Commands

```bash
cd /Users/ujang/Projects/Geekgarden/scylla-be/sales && rtk go test ./controller/... ./service/... -count=1
cd /Users/ujang/Projects/Geekgarden/scylla-be/sales && rtk go test ./... -count=1
```

## Output

```text
Go test: 165 passed in 2 packages
Go test: 185 passed in 22 packages
```

## Notes

- Full module suite passed after:
  - export `cust_id` request-body support
  - trend-sales GET body `cust_id` support
  - scope enforcement via `resolveSecondaryDashboardCustID`
  - 403 mapping for unauthorized `cust_id`
  - `alphanum,max=20` validation on request `cust_id`
  - explicit 400 controller tests for invalid `cust_id` format on export and trend endpoints
  - DTO-only body binding comment to prevent direct `BodyParser` into `entity.SecondarySalesReportQueryFilter`
- No failing package remained.
