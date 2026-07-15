# SX-2516 C1-F6 remediation

## Result
Worker no longer runs empty-scope replace or marks it SUCCESS. `ProcessSecondarySalesImport` now loads staging by `history_id`, rejects empty/incomplete rows, derives distinct date scope, and marks history `FAILED` before any order mutation.

## Schema decision
Chosen fail-closed path. Existing 22-column `import.sales_update_temp` schema cannot reconstruct validated mapped `model.Order`/`model.OrderDetail`: missing `outlet_id`, `salesman_id`, `pro_id`, `wh_id`, and mapped detail fields (`qty1..`, conversion, prices, amount, VAT, promo/discount fields). No migration added. In-memory payload was rejected for worker correctness because current public worker API only accepts `history_id`; process restart would lose payload and persisted staging still could not rebuild records.

## Changes
- Added repository `FindSalesUpdateTempByHistoryID`, filtering `status_insert='SUCCESS'`, ordered by ID.
- Enqueue staging no longer writes blank NOT NULL mapping fields; writes only document identity/date and success status.
- Worker validates history tenant, nonempty document/date, distinct date scope, then fails with explicit schema-gap error. History becomes `FAILED`; never `SUCCESS`.
- No controller validate=true wiring, G/H, smoke, migration, dependency, env, or compose changes.

## Validation
- `cd sales && rtk go mod download && rtk go mod tidy`: PASS
- `cd sales && rtk gofmt -l service/order_service.go repository/order_repository.go`: PASS, no output
- `cd sales && rtk go test ./controller/... ./service/... ./repository/...`: PASS, 362 tests
- `cd sales && rtk go build ./...`: PASS
- Focused new rollback/scope/preservation tests: NOT added; worker now intentionally fails closed before mutation because staging schema cannot reconstruct mapped orders. Existing replacement repository methods remain unchanged and transaction-aware.

## Blocker
C1-F6 cannot ship as runnable replace until approved persistence change adds complete mapped order/detail payload, or worker contract is redesigned to safely retain parsed payload across worker invocation. Current behavior is deliberately non-runnable and cannot report false `SUCCESS`.
