# Discovery Evidence — SX-1879 Export Issue Data

## Files Inspected

- `sales/controller/so_controller.go`
  - Route `GET /v1/download` is registered in `SoController.Route`, mounted under the sales service base path by deployment/API gateway.
  - `Download` parses `start_date`, `end_date`, and both `salesman_id` / `salesman_id[]`, validates max 31-day range, injects `cust_id`, `parent_cust_id`, and `user_fullname`, then calls `SoService.Download`.
- `sales/service/so_service.go`
  - `Download` creates an async report entry with `FILE_STATUS_PROCESSING`, then calls `generateDownloadSalesOrderExcel` in a goroutine.
  - `generateDownloadSalesOrderExcel` queries four datasets and writes sheets: `Purchase Order`, `Sales Order`, `Final Order`, `QTY Summary`.
  - The sheet writers currently use `derefFloat64` for amount columns, returning raw `float64` values or empty string for nil.
  - `mapPoToEntity`, `mapSoToEntity`, and `mapFinalToEntity` compute `GrossSales`, `Promotion`, `Discount`, `NetSales`, `Vat`, and `Gross` in service layer.
- `sales/repository/so_repository.go`
  - `FindDownloadDataPo`, `FindDownloadDataSo`, `FindDownloadDataFinal`, and `FindDownloadQtySummary` query `sls.order_detail` joined to `sls.order` and master tables.
  - Queries include `sls.order_detail.cust_id = ?`, `item_type = 1`, date range, and optional salesman filter.
  - Product/supplier fields already use several `COALESCE` fallbacks, but many order/customer/salesman/unit/amount fields remain nullable by model.
- `sales/model/so_download.go`
  - Download model fields are mostly pointer types for nullable columns.
  - `VatValueFinal`, `DiscValueFinal`, and `Vat` are available for all three financial sheets.
- `sales/entity/so_download.go`
  - Export row entities keep amount fields as `*float64`.
- `sales/service/so_service_test.go`
  - Existing tests cover async Excel generation headers, date range, PO number filtering, basic mapper field order, and QTY summary null default.
  - Some tests currently encode the suspicious behavior: `DiscValueFinal` maps to `Promotion`, while `VatValueFinal` maps to `Discount`.

## Project Patterns Found

- Module is `sales` with `go 1.23.0` and `github.com/xuri/excelize/v2 v2.9.1`.
- Strict layer flow is followed: controller → service → repository → DB.
- Existing export code centralizes Excel creation inside `sales/service/so_service.go` rather than a dedicated export helper.
- Null-safe string/date mapping is mostly explicit in mapper functions.
- Numeric null handling is inconsistent: QTY Summary PO quantity uses `derefFloat64Zero`, while amount columns use `derefFloat64` and return blank for nil.
- No existing reusable Indonesian amount formatter was found in the sales module.

## Reuse Candidates

- Reuse `excelize` already present for cell writing and optional cell style/number format.
- Reuse current mapper functions and add small helpers near existing dereference helpers rather than adding a new layer.
- Reuse existing `so_service_test.go` mocks and Excel open/read helpers for regression tests.
- Extend tests around mapper calculations and Excel cell values instead of introducing new framework/dependency.

## Commands / Docs Checked

- Required repository command run first: `rtk docker compose -f docker-compose.yml ps`.
  - Only `master`, `redis`, and `system` were up; `sales` was not running. For planning, static discovery was sufficient.
- Local file discovery via `Glob`, `Grep`, and `Read`.
- Official docs were not required for `excelize` because current code already uses stable `SetCellValue`, `NewStyle`, and `SetCellStyle`; if implementation chooses numeric cell styles, verify exact `excelize.NumFmt` / custom number format behavior before coding.
- GitHub and web search were not required; this is a local defect in project-specific export logic.
- Browser/screenshot evidence was not required; this is backend Excel output, not UI parity.

## Likely Root Cause Candidates

1. Financial field mapping appears swapped or mislabeled:
   - Current code sets `promotion = DiscValueFinal` and `discount = VatValueFinal` in PO/SO/Final mappers.
   - QA expects `Discount = 0`, `PPN = 1.100.000`, and `Gross = 12.000.000` for a case where using VAT value as discount would reduce net sales and then recompute PPN incorrectly.
2. Current code recomputes VAT as `netSales * vatPercent / 100` instead of using the stored `vat_value_final` value already selected by repository.
   - QA expected `PPN = 1.100.000` while `Net Sales (Exc PPN) = 12.000.000`; this is not a simple 11% of 12.000.000.
3. Current `gross := netSales + vat` conflicts with QA expected `Gross = 12.000.000` when `PPN = 1.100.000`, suggesting export `Gross` is expected to represent sales amount before PPN, or existing domain uses stored gross/base amount differently.
4. Amount columns are written as raw numeric cells without Indonesian thousands separator display; `GetCellValue`/download may show `11000000` instead of `11.000.000`.
5. Purchase Order data with null financial fields may become blank where business expects `0`, and nullable non-financial fields need continued safe fallbacks.

## Constraints

- Do not hardcode staging credentials, Jira tokens, or sensitive data.
- Preserve tenant filtering with `cust_id` and parent lookup behavior.
- Avoid broad formula changes without tests documenting expected behavior.
- Endpoint creates async report data; manual verification must inspect report output/history or decoded `FileBase64`, not expect direct file bytes from initial `GET /v1/download` call.
- Repository root has AGENTS instructions requiring `rtk`, while global OpenCode instructions say not to prefix commands. The repository-local instruction is more specific for this repo; `rtk` was used for the mandatory compose check.

## Risks

- Existing tests currently assert the old/suspicious field mapping, so implementation must update tests deliberately.
- QA expected values may be row-level or order-level; local code currently writes detail rows. Need verify whether duplicate item rows should each carry line-level values or order-level totals.
- If Excel cells are exported as formatted strings, downstream consumers may lose numeric type; if numeric cells with custom format are used, QA tooling may still read raw values depending on how it inspects the file.
- Staging data for `SO2604290009` may not be accessible in local environment; create fixture-equivalent automated tests and add manual staging verification steps.
