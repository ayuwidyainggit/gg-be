# SX-2516 — Discovery Evidence

## Scope

Rencana perubahan `sales` untuk replace secondary-sales import. Tidak ada source edit pada tahap planning.

## Files inspected

- `AGENTS.md`
- `.opencode/docs/PROJECT_STACK.md`
- `.opencode/docs/PROJECT_COMMANDS.md`
- `.opencode/docs/FRAMEWORK_PLAYBOOK.md`
- `.opencode/docs/PROJECT_DETECTED_TOOLS.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/MCP.md`
- `sales/controller/order_controller.go:39-70,886-970`
- `sales/controller/order_controller_test.go:650-746`
- `sales/entity/order.go:57-104`
- `sales/model/order.go`
- `sales/model/order_detail.go:7-109`
- `sales/service/order_service.go:370-489,6543-7269`
- `sales/repository/dbtransaction.go:1-56`
- `sales/go.mod`
- `master/entity/history.go:5-21`
- `master/repository/outlet_repository.go:2720-2800`
- `master/service/outlet_service.go:5289-5451`
- `master/service/product_service.go:1892-1941`

## Confirmed patterns and reuse

1. Sales route exists: `POST /v1/orders/export-template/import`, handled by `OrderController.ImportFromUrl`; API gateway-visible prefix `/sales/v1` is user-confirmed endpoint contract.
2. Current handler accepts `{url, validate}` and uses login `cust_id`, `parent_cust_id`, `user_id`; `validate` currently parses as string via `isTruthy`.
3. Existing exported template has exact 16 headers in `orderImportHeaders`; `ExportTemplate` generates XLSX only. `parseImportOrders` validates same header and parses via existing `excelize/v2`.
4. Existing parser resolves outlet and salesman with `cust_id`, distributor product mapping with product code/name, parent product, distributor unit slots, parent unit conversions/prices, then builds imported `CreateOrderBody` with `is_sales_mapping=true`, `data_status=6`, `data_source=3`.
5. Existing sync `ImportOrders` validates then calls `Store` once per invoice; each `Store` opens its own transaction. This does not satisfy scope-level replace atomicity.
6. `sls.order.is_sales_mapping` exists in model and prior raw SQL migration. No new mapping flag migration required.
7. `sales/model/order_detail.go` contains all requested unit/price/amount/detail fields.
8. Sales has transaction helper `WithinTransaction(ctx, fn)` injecting a GORM transaction into context. New repository writes must use this context path.
9. `import.import_history` is an existing shared table pattern in master. `import.sales_update_temp` does not exist in repository.
10. `sales/migration/` contains raw SQL, including `sales/migration/sls.order/`; no generated migration tool exists.
11. `sales/go.mod` already includes `github.com/xuri/excelize/v2 v2.9.1`, `github.com/rs/xid`, `github.com/go-co-op/gocron/v2`, and `github.com/streadway/amqp`; no dependency addition planned.
12. Sales tests use `net/http/httptest`, mocks over service interfaces, and table-driven Go tests. Existing controller test SX-2475 validates `validate=false` does not call import.

## Source strategy

- Repo-local: used. Sufficient for current endpoint, models, parser, transaction, migration, and tests.
- User Jira/doc text: used as behavior authority, endpoint conflict resolved by user confirmation.
- Official docs/Context7: skipped. No version-sensitive new API/dependency. Existing `excelize`, Fiber, GORM, gocron reused unchanged.
- GitHub/web/browser: skipped. No upstream or visual dependency.
- `sequential-thinking`: used for scope/risk framing; planner skill required it.
- Read-only advisors: `@explorer` and `@architect`; findings incorporated.

## Constraints and risks

- `import.import_history` schema columns/status semantics must be inspected in target local DB before migration; repo model comes from master and may be incomplete.
- User selected existing `import.import_history` and `import.sales_update_temp` flow “as is”; plan must not invent separate history table, status GET endpoint, OBS persistence, or watchdog.
- Async job may be non-durable across process restart because user rejected added persistence/watchdog. Claim is bounded: accepted request creates/holds history+staging before background replace work starts; no retry guarantee claimed.
- Concurrent imports with same `cust_id` and date scopes need transaction-level serialization. Use PostgreSQL transaction advisory locks keyed by `(cust_id, normalized document_date)` before delete/insert; no new library.
- `ro_no` stays equal to `DocumentNo`; delete matching mapping header/detail first then recreate inside one transaction. Any failing delete/insert rolls all back.
- Do not use customer input as tenant scope. Always derive `cust_id`/user ID from JWT locals.

## Confirmed vs Assumed Audit

| Claim | Class | Evidence |
|---|---|---|
| Route is mounted as `/v1/orders/export-template/import` in sales | confirmed_repo | `sales/controller/order_controller.go:43-56` |
| Gateway contract is `/sales/v1/orders/export-template/import` | user_confirmed | question gate 2026-07-15 |
| `url` is final request field | user_confirmed | question gate 2026-07-15 |
| `false=validate`, `true=write` | user_confirmed | question gate 2026-07-15 |
| `is_sales_mapping` exists | confirmed_repo | `sales/model/order.go:16`; `order_service.go:395-397` |
| Parser header contract is 16-column `orderImportHeaders` | confirmed_repo | `sales/service/order_service.go:6548-6553` |
| `import.sales_update_temp` absent in repo | confirmed_repo | explorer grep report; no model/migration match |
| Existing `import.import_history` can store desired statuses unchanged | assumption | User requested as-is; local DB DDL must verify capacity before source edit |
| Background job survives restart/retries | unverified | User chose as-is/no persistence/watchdog; plan does not claim this |
| Existing DDL has cascading foreign key header/detail | unverified | must inspect `\d+ sls.order_detail` before migration and use explicit detail-first delete regardless |
| Workload maximum is 4200 rows per request | user_confirmed | question gate calculation: 20 invoices/day × 10 products × 3 UOM × 7 days |
| CSV/XLS/XLSX support exists | partial | task says all; current parser/export uses `excelize.OpenReader`, which must be tested with actual CSV/XLS fixtures before claim |
