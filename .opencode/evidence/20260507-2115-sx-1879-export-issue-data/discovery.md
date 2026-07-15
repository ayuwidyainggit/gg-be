# Discovery — SX-1879 Export Issue Data

## Files Inspected

- `AGENTS.md` — repo constraints: multi-module Go, `sales` service on port `9004`, strict controller → service → repository layering, `rtk` command expectation, no secrets in source/logs.
- `sales/controller/so_controller.go` — `Route` registers `GET /v1/download` under JWT protected `/v1`; with service prefix this matches `/sales/v1/download`.
- `sales/service/so_service.go` — export generation, sheet writers, amount formatter, mapper functions for PO/SO/Final/QTY Summary.
- `sales/repository/so_repository.go` — `FindDownloadDataPo`, `FindDownloadDataSo`, `FindDownloadDataFinal`, and selected source columns.
- `sales/model/so_download.go` — nullable pointer fields for selected order detail export columns.
- `sales/entity/so_download.go` — export row DTO fields for workbook sheets.
- `sales/service/so_service_test.go` — existing tests already contain SX-1879 regression coverage for formatter, SO/PO/Final mapping, sheet formatting, and nullable PO amounts.
- `.opencode/plans/20260504-2141-sx-1879-export-issue-data.md` — prior SX-1879 plan; many requested fixes appear already implemented in current code.
- `.opencode/plans/20260507-2105-sx-1878-1879-po-final-export.md` — newer combined plan; notes that current SX-1879 work may mainly be verification/hardening.

## Commands / Docs Checked

- `rtk docker compose -f docker-compose.yml ps` from repo root.
  - Result: `system`, `master`, `sales`, `finance` are up; `sales` is up on `9004`.
  - RTK emitted warning about untrusted `.rtk/filters.toml`; output still returned.
- Local code search via glob/grep/read tools.
- No official docs/context7 needed: this is project-local Go/export logic, and `excelize` usage is already established in project tests.
- No GitHub/web/browser evidence needed for planning: defect depends on local backend code plus Jira evidence supplied by user. Staging token is intentionally not used or stored.

## Project Patterns Found

- Endpoint flow: controller calls `SoService.Download`; async generator `generateDownloadSalesOrderExcel` fetches PO/SO/Final/QTY data via repository, builds Excel with `excelize`, stores base64 content in report repository with status `Ready`.
- Sheet writers are in `sales/service/so_service.go`:
  - `createPurchaseOrderSheet`
  - `createSalesOrderSheet`
  - `createFinalOrderSheet`
  - `createQtySummarySheet`
- Existing amount helper:
  - `formatDownloadAmount(*float64) string` returns `0` for nil and applies Indonesian thousands separator by inserting `.`.
  - `formatDownloadAmountValue(float64) string` rounds to integer and supports negative values.
- Current mappers already map:
  - `Discount` from `DiscValueFinal` with fallback `0`.
  - `Vat` from `VatValueFinal` with fallback `0`.
  - `NetSales` as `grossSales - promotion - discount`, with `promotion = 0`.
  - `Gross` as `grossSales`.
- Current SO expected fixture (`Qty1=1`, `SellPrice1=12000000`, `DiscValueFinal=0`, `VatValueFinal=1100000`) already maps to `Discount=0`, `NetSales=12000000`, `Vat=1100000`, `Gross=12000000`.
- Sheet amount columns `Q:V` and `Z:AE` already call `formatDownloadAmount` for PO/SO/Final. Quantity columns still use `derefFloat64`, returning blank for nil.
- Existing tests in `so_service_test.go` include:
  - `TestFormatDownloadAmount_UsesIndonesianThousandsSeparator`
  - `TestMapPoToEntity_SX1879MapsFinancialFieldsCorrectly`
  - `TestMapSoToEntity_SX1879MapsFinancialFieldsCorrectly`
  - `TestMapFinalToEntity_SX1879MapsFinancialFieldsCorrectly`
  - `TestCreateSalesOrderSheet_FormatsAmountColumnsWithIndonesianSeparator`
  - `TestCreatePurchaseOrderSheet_DefaultsNullableAmountColumnsToZero`

## Reuse Candidates

- Reuse existing `formatDownloadAmount` for all amount/currency columns.
- Reuse existing `so_service_test.go` helpers (`ptrFloat64`, `excelize.NewFile`, cell assertions, base64 workbook readers).
- Reuse existing mapper functions if tests pass; add only targeted hardening rather than new helpers.
- Reuse Final Order formula as comparison baseline for Sales Order if staging evidence still fails.

## Constraints

- Do not hardcode/copy Jira/staging token; manual reproduction must use env such as `$SCYLLA_STAGING_TOKEN`.
- Root has no `go.mod`; test commands must run from `sales/`.
- Existing source appears to include fixes from prior plan, so implementation may be verification-first and minimal.
- Endpoint is JWT-protected and async report based; direct curl may not immediately return `.xlsx` bytes if local route follows report creation flow.

## Risks

- QA evidence says latest retest still fails, but current local code already contains many expected fixes; staging may be behind current branch, or remaining bug may be data-specific repository source/row filtering not covered by fixture tests.
- Existing `filterDownloadDataPoWithPONumber` drops PO rows with blank `po_no`; if QA's "null data PO" row has blank `po_no`, Purchase Order sheet may omit it rather than safely exporting fallback `order_no`.
- Current nullable amount fallback covers amount columns, but nullable quantity columns still become blank. This is likely acceptable unless QA expects numeric nullable qty to become `0`.
- Current `formatDownloadAmountValue` manually rounds; values with fractions may need domain confirmation if decimals should be preserved or rounded.
- Manual staging reproduction requires valid credentials/token and may expose sensitive data; do not store output artifacts with secrets.
