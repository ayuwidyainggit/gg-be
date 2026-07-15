# Validation SX-2079

Module: `master`
Date: 2026-05-26

## Commands

```bash
rtk go test ./service -run 'TestNormalizeDropdownScope|TestRegionService|TestAreaService|TestBusinessUnitService'
# Go test: 10 passed in 1 packages

rtk go test ./repository -run 'TestBuildRegionListQuery|TestBuildAreaListQuery|TestBuildFindDistributorsByCustIDQuery'
# Go test: 8 passed in 1 packages

rtk go test ./controller -run 'Test.*BusinessUnit|Test.*Region|Test.*Area'
# Go test: 4 passed in 1 packages

rtk go test ./...
# Go test: 317 passed in 23 packages
```

## Coverage evidence

- Scope normalization helper tested for `Specific`, `SPESIFIC`, `selected`, `ALL`, empty, unknown.
- Region repository query builder tested for:
  - specific region mapping join
  - all scope cust filter path
- Area repository query builder tested for:
  - specific area mapping join
  - region mapping fallback path
  - `ma` alias on search clause
- Business-unit repository query builder tested for:
  - all/all/all parent_cust_id fallback
  - specific distributor mapping join
  - combined region+area mapping join
- Controller tests confirm token locals `employee_id`/`distributor_id` and array query parsing for:
  - business-unit
  - regions
  - areas
- Service tests confirm:
  - principal scope load for region/area/business-unit
  - missing `employee_id` fails principal path
  - non-principal region path skips scope lookup
  - non-principal area path skips scope lookup
  - non-principal business-unit distributor path does not trigger employee scope lookup
- Additional repository parity test confirms non-principal region query still scopes by expected parent tenant argument.

## Manual/runtime smoke

- Not executed.
- Reason: no safe local token fixture in repo; docs contain real-looking tokens and were intentionally not reused.
