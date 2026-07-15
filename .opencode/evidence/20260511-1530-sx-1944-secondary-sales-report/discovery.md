## Discovery Evidence — SX-1944 Secondary Sales Report

### Files inspected
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/entity/report.go`
- `sales/model/report.go`
- `sales/service/report_service_test.go`
- `sales/pkg/str/generator.go`
- `SecondarySales-080526-003 (2).xlsx`
- `AGENTS.md`

### Commands and checks run
- `rtk docker compose -f docker-compose.yml ps`
  - Result: required root services are up, including `sales`.
- Timestamp verification for QA payload:
  - `1778086800` -> `2026-05-06T17:00:00Z` -> `2026-05-07 00:00:00 Asia/Jakarta`
  - `1778173199` -> `2026-05-07T16:59:59Z` -> `2026-05-07 23:59:59 Asia/Jakarta`

### Project patterns found
- Route for target endpoint is in `sales/controller/report_controller.go`:
  - `POST /v1/reports/secondary-sales`
- Controller injects auth context into request filter:
  - `cust_id`
  - `parent_cust_id`
  - `user_fullname`
- Export is asynchronous via RMQ:
  - request publishes job in `PublishSecondarySalesReport`
  - worker consumes in `processSecondarySalesExportMessage`
  - workbook is generated in `SubscribeSecondarySalesReport`
- Export source already reuses repository query through `ReportRepository.SecondarySalesUnion(...)`.
- Shared SQL builder already exists:
  - `buildSecondarySalesUnionQuery(dataFilter, withPagination)`
  - both paginated and export variants call the same builder
- Date filter pattern in this module uses `str.UnixTimestampToUtcTime(...)` and SQL `BETWEEN ? AND ?`.

### Reuse candidates
- Reuse existing controller/service/repository flow for secondary sales.
- Reuse `buildSecondarySalesUnionQuery(...)` as the single source of truth rather than introducing a new separate query path.
- Reuse existing service export tests in `sales/service/report_service_test.go` for workbook/output consistency.
- Reuse `str.UnixTimestampToUtcTime(...)` because the QA timestamps already map correctly to the intended Asia/Jakarta day boundary when converted from Unix seconds.

### Root-cause evidence from existing code
- Current query already moved away from old `SecondarySales(...)` and uses `sls.order_detail` + `sls."order"` with `invoice_date`, `invoice_no is not null`, and `data_status in (6,7)`.
- The highest-risk defect area in current implementation is product fallback join:
  - `LEFT JOIN mst.m_product pp ON pp.cust_id = ? AND (pp.pro_id = detail_product_id OR (cp.pro_code IS NOT NULL AND pp.pro_code = cp.pro_code))`
- That `OR` join can duplicate rows if parent product data is non-unique for `pro_code` or if both matching routes resolve to multiple rows.
- Because the query uses `UNION ALL` and no outer deduplication, any duplicated join row becomes duplicated export/report output and inflates totals.
- Existing return branch keeps `ppn` positive while gross/net values are negative. User instructed to assume existing behavior unless stronger business evidence appears.

### Spreadsheet evidence
- File `SecondarySales-080526-003 (2).xlsx` contains sheet `Report`.
- Relevant output columns exist:
  - `GrossSales`
  - `NetSalesExcPPN`
  - `PPN`
  - `NetSalesIncPPN`
- The sheet also contains a total row, so query duplication would directly affect exported totals.

### Constraints
- Scope confirmed: `sales/` module only.
- Do not commit Jira/staging credentials or tokens.
- Follow repo architecture: Controller -> Service -> Repository -> DB.
- Keep export and report response aligned through the same repository query path.
- Planner mode only: no implementation edits outside `.opencode/`.

### Risks
- Duplicate row inflation from parent/child product fallback join.
- Silent regression if only export path is fixed but paginated/report path diverges.
- Date boundary regression if query changes away from current Unix-to-UTC conversion without proof.
- Product filter behavior ambiguity if fallback product mapping changes selected identifier semantics.
- Return `PPN` sign mismatch if business rule differs from existing implementation.

### Research gate outcome
- Local project discovery: required and completed.
- Official docs / context7: not required; issue is internal SQL/report behavior, not library-version-sensitive.
- GitHub/upstream research: not required; logic is repo-local.
- Brave/web search: not required; no external current facts needed.
- Browser/screenshot evidence: not required for backend report bug.
