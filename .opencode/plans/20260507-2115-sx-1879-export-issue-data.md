# Plan — SX-1879 Export Issue Data

## Goal

Pastikan backend `sales` untuk `GET /sales/v1/download` menghasilkan export order `.xlsx` yang benar untuk Jira `SX-1879`: Purchase Order null-safe, Sales Order financial columns sesuai expected QA, Final Order tidak regress, dan amount tetap memakai separator ribuan Indonesia.

## Non-goals

- Tidak mengubah FE/menu `Sales > Order List`.
- Tidak hardcode/copy token Jira, credential staging, atau data sensitif ke source, test, log, atau artifact.
- Tidak mengubah schema/migration database kecuali implementer menemukan bukti kolom sumber memang tidak tersedia.
- Tidak rewrite arsitektur async report/export.
- Tidak mengubah Final Order behavior yang sudah PASS kecuali hanya untuk deduplikasi helper yang terbukti netral oleh test.

## Scope

- Module target: `sales`.
- Endpoint routing: `sales/controller/so_controller.go` (`GET /v1/download`, protected by JWT; exposed as `/sales/v1/download` behind gateway/service prefix).
- Export service: `sales/service/so_service.go`.
- Export repository/model: `sales/repository/so_repository.go`, `sales/model/so_download.go` only if remaining bug requires source column/filter changes.
- Tests: `sales/service/so_service_test.go`, optional new `sales/service/so_download_export_test.go` if file split improves readability.

## Requirements

- Filter QA to verify: `start_date=1777420800`, `end_date=1777593599`, `salesman_id[]=421`.
- Purchase Order export must not fail or emit literal `null`, `undefined`, or `NaN` for nullable data.
- Nullable text/date fields should export as empty string.
- Nullable amount fields should export as formatted `0`.
- Sales Order row(s) for QA scenario must export:
  - `Discount = 0`
  - `Net Sales (Exc PPN) = 12.000.000`
  - `PPN = 1.100.000`
  - `Gross = 12.000.000`
- Final Order must remain aligned with its passing expected behavior.
- Amount columns must use Indonesian thousands separator (`.`), e.g. `11.000.000`, `12.000.000`, `1.100.000`.
- Automated regression tests must cover formatter, Sales Order financial mapping, Purchase Order null data, and Final Order non-regression.

## Acceptance Criteria

1. Targeted automated tests for SX-1879 pass in `sales` module.
2. Full `rtk go test ./...` from `sales/` passes, or any blocker is recorded with exact error.
3. Existing or added tests prove `formatDownloadAmount(nil) == "0"`, `11000000 -> "11.000.000"`, `12000000 -> "12.000.000"`, and `1100000 -> "1.100.000"`.
4. Tests prove Sales Order export mapping/sheet cells output `Discount=0`, `Net Sales (ExcPPN)=12.000.000`, `PPN=1.100.000`, `Gross=12.000.000`.
5. Tests prove Purchase Order amount cells tolerate nil pointers and output `0` instead of blank/invalid values for amount columns.
6. Tests prove Final Order financial mapping remains correct for the same expected QA fixture.
7. If valid staging/local auth is available, manual reproduction verifies the QA filter without storing token.

## Existing Patterns/Reuse

- Reuse current export flow in `generateDownloadSalesOrderExcel` and sheet writer functions.
- Reuse existing helper `formatDownloadAmount` / `formatDownloadAmountValue`; current code already applies Indonesian thousands separator and nil-to-`0` behavior.
- Reuse existing mapper structure for PO/SO/Final and compare SO with Final where formula should match.
- Reuse existing tests in `sales/service/so_service_test.go`; discovery shows several SX-1879 tests already exist and may already satisfy most requested coverage.
- Reuse `excelize` dependency already present in `sales/go.mod` and test patterns that read workbook cell values.
- Tidak ditemukan kebutuhan helper/KiloCode baru; pilih Reuse > Extend > Create.

## Constraints

- Repo multi-module Go; commands for this fix run from `sales/`.
- Repo instruction requires `rtk` prefix for shell commands in this project. Compose status was checked with `rtk docker compose -f docker-compose.yml ps`.
- Staging curl needs a valid token from env/local auth only; do not paste token into commands recorded in artifacts.
- Current local code appears already contains prior SX-1879 fixes; implementation should start with verification to avoid redundant or risky changes.
- Endpoint uses JWT and async report generation; validating actual `.xlsx` may require existing report retrieval/download-history flow, not only immediate `GET /v1/download` response.

## Risks

- QA latest failure may be from staging branch not containing current local fixes, or from a data-specific edge not covered by existing fixtures.
- `filterDownloadDataPoWithPONumber` currently excludes rows whose `po_no` is blank. If QA's null-data PO case has empty `po_no` but a valid `order_no`, this may be the remaining Purchase Order issue and needs a regression test before change.
- Current amount columns are formatted as strings. This matches QA visible expected output, but Excel consumers needing numeric cells could be affected if behavior changes further.
- The meaning of `Gross` in QA expected equals pre-PPN gross/net (`12.000.000`), not `net + ppn`; do not reinterpret without product confirmation.
- Quantity nils currently export blank through `derefFloat64`; only change to `0` if QA evidence confirms qty cells are part of the null-data failure.

## Decisions/Assumptions

- Interaction level: Assumption-first. Requirements and local patterns are sufficient for a concrete implementation/verification plan; no blocking question gate needed.
- Question Gate: not asked. The only material unknown is whether PO rows with blank `po_no` but valid `order_no` should appear; plan treats this as a likely implementation hypothesis to prove by test/evidence, not a blocker.
- Financial formula assumption, matching current code and QA expected:
  - `grossSales = sum(qty * sell_price)` for the sheet's relevant qty/price fields.
  - `promotion = 0` unless a project field/source is proven.
  - `discount = disc_value_final` fallback `0`.
  - `netSales = grossSales - promotion - discount`.
  - `ppn = vat_value_final` fallback `0`.
  - `gross = grossSales` for this export contract.
- If tests already pass, implementation should focus on missing edge tests and staging/local verification rather than changing formula.

## TDD/Test Plan

### TDD Required

Ya. Ini bug output contract backend dan export data production.

### Reason

- Mencegah regresi decimal separator yang sudah PASS.
- Mengunci Sales Order formula agar sama dengan expected QA dan Final Order baseline.
- Mengunci null-safe Purchase Order export agar tidak menghasilkan invalid workbook/cell values.

### Existing Test Patterns

- `sales/service/so_service_test.go` already has package-level helper `ptrFloat64` and `excelize.NewFile`/`GetCellValue` based assertions.
- Existing SX-1879 tests already cover formatter, mapper, Sales Order sheet amount formatting, and Purchase Order amount nil fallback.

### First Failing/Regression Test

Run existing tests first as the Red/verification gate:

```bash
rtk go test ./service -run 'TestFormatDownloadAmount|TestMap.*SX1879|TestCreateSalesOrderSheet_FormatsAmountColumnsWithIndonesianSeparator|TestCreatePurchaseOrderSheet_DefaultsNullableAmountColumnsToZero|TestFilterDownloadDataPo'
```

If they pass, add the smallest missing regression test before any code change:

1. `TestFilterDownloadDataPoWithPONumber_AllowsOrderNoFallbackForNullDataRows` if discovery/reproduction shows PO null-data rows are omitted due blank `po_no`.
2. `TestCreateFinalOrderSheet_FormatsSX1879AmountColumns` if Final Order sheet-level formatting lacks direct coverage.
3. `TestCreatePurchaseOrderSheet_NullTextAndAmountCellsDoNotEmitInvalidLiterals` if QA evidence shows literal `null`/`NaN` in cells beyond amount columns.

### Green Step

- If existing tests fail: fix only the failing mapper/formatter/sheet logic in `so_service.go`.
- If PO blank `po_no` omission is proven: change PO filtering to keep valid rows when either `po_no` or `order_no` is present, while display uses `resolveDownloadPONumber` fallback.
- If Final Order sheet-level coverage is missing: add tests only; no code change if behavior already matches.
- If nullable text/date issue appears: centralize text fallback with existing empty-string pattern; avoid outputting `%!s(<nil>)`, `null`, or `NaN`.

### Refactor Step

- Keep helper changes local and small.
- Consider extracting a shared financial calculation helper only if the same bug requires changing more than one mapper; otherwise leave current explicit PO/SO/Final mapping for readability.
- Do not introduce new dependencies for formatting.

### Edge Cases

- `nil` amount pointers.
- Zero and large integer amounts.
- Fractional amounts requiring rounding to no decimals.
- Negative amounts if return/adjustment rows exist.
- Blank `po_no` with valid `order_no`.
- Missing invoice/outlet/supplier/product fields from left joins.
- Multiple rows for the same SO; verify line-level export does not double-count header-level totals.

### Commands

From `sales/`:

```bash
rtk go test ./service -run 'TestFormatDownloadAmount|TestMap.*SX1879|TestCreateSalesOrderSheet_FormatsAmountColumnsWithIndonesianSeparator|TestCreatePurchaseOrderSheet_DefaultsNullableAmountColumnsToZero|TestFilterDownloadDataPo'
rtk go test ./...
```

## Implementation Steps

1. Run targeted tests above to determine whether current code already satisfies requested unit coverage.
2. Inspect any failing assertions and confirm whether failure is formatter, mapper formula, sheet cell formatting, or PO row filtering.
3. If formatter fails, fix `formatDownloadAmount`/`formatDownloadAmountValue` to return `0` for nil/non-finite and apply `.` thousands separator.
4. If Sales Order financial mapping fails, align `mapSoToEntity` with current Final Order formula and QA expected values: `DiscValueFinal` as Discount, `VatValueFinal` as PPN, gross before PPN.
5. If Purchase Order null-data issue is row omission, update `filterDownloadDataPoWithPONumber`/helper naming to keep rows with either valid `po_no` or valid `order_no`; add tests for blank `po_no` fallback.
6. If Purchase Order null-data issue is invalid cell output, add safe fallback tests and adjust sheet writer/mappers for the affected columns only.
7. Add/adjust Final Order non-regression sheet-level test if not already covered.
8. Run targeted tests and full `rtk go test ./...`.
9. Optional manual verification when valid auth is available:
   - Use `$SCYLLA_STAGING_TOKEN` and configured base URL env.
   - Call QA filter without printing token.
   - Retrieve generated report/export through existing flow and inspect sheets/cells.
10. Summarize root cause from actual failing evidence, not assumptions.

## Expected Files to Change

- Likely tests only if current code passes:
  - `sales/service/so_service_test.go`
- Possible source files if tests/reproduction expose remaining gap:
  - `sales/service/so_service.go`
  - `sales/repository/so_repository.go` only if selected fields/filter query are proven wrong.
  - `sales/model/so_download.go` only if repository must select additional columns.

## Agent/Tool Routing

- Implementation/TDD: route to `@fixer` / `opencode-fixer` for bounded edits and tests.
- Deep formula ambiguity: route to `@oracle` only if SO vs Final amount semantics conflict with product evidence.
- Security/privacy reviewer is not required for code-only fix, but must be considered if anyone attempts to store staging token/exported sensitive workbook.
- Release engineer is not required unless deployment/report infrastructure issues emerge.

## Validation Commands

From repo root if runtime verification needs services:

```bash
rtk docker compose -f docker-compose.yml ps
rtk docker compose -f docker-compose.yml up -d sales
```

From `sales/`:

```bash
rtk go test ./service -run 'TestFormatDownloadAmount|TestMap.*SX1879|TestCreateSalesOrderSheet_FormatsAmountColumnsWithIndonesianSeparator|TestCreatePurchaseOrderSheet_DefaultsNullableAmountColumnsToZero|TestFilterDownloadDataPo'
rtk go test ./...
```

Manual staging verification, with token only from env:

```bash
rtk curl -H "Authorization: Bearer $SCYLLA_STAGING_TOKEN" \
  "$SCYLLA_STAGING_SALES_BASE_URL/v1/download?start_date=1777420800&end_date=1777593599&salesman_id[]=421"
```

Do not paste the token or commit the downloaded workbook unless sanitized and explicitly approved.

## Evidence Requirements

- Test output for targeted SX-1879 tests.
- Full `rtk go test ./...` output or exact blocker.
- If source changes are made: concise diff summary of changed mapper/formatter/filter logic.
- If manual verification is performed: URL without token, report id/status, sheet names checked, cell values for Discount/Net Sales/PPN/Gross, and confirmation token was not stored.
- Keep no staging `.xlsx` artifact in repo unless sanitized and intentionally approved.

## Done Criteria

- Automated tests cover null PO amount handling, SO financial expected values, Final Order non-regression, and Indonesian amount formatting.
- Implementation changes, if any, are minimal and localized.
- No secrets/tokens added to source, tests, fixtures, logs, or `.opencode` artifacts.
- Reported root cause distinguishes whether bug was missing deployment/current branch mismatch, PO row filtering, mapper formula, or formatter/null safety.
- Primary plan remains the implementation source of truth.

## Final Planning Summary

- Artifacts created:
  - `.opencode/plans/20260507-2115-sx-1879-export-issue-data.md` — primary source of truth for implementation.
  - `.opencode/evidence/20260507-2115-sx-1879-export-issue-data/discovery.md` — kept because it records inspected files, current-code evidence, and key risk that local code may already include prior fixes.
- Draft artifacts: none created; no stale draft cleanup needed.
- Key decisions:
  - Start implementation with targeted test verification because current local code already contains formatter and SX-1879 mapping tests.
  - Treat PO blank `po_no` filtering as the leading remaining hypothesis for Purchase Order null-data failure if existing tests pass.
  - Do not use or store Jira/staging token; manual verification uses env only.
- Assumptions:
  - QA expected `Gross` means pre-PPN gross for export, matching current code.
  - String-formatted amount cells are acceptable because QA explicitly validates visible values with `.` separators.
- Open questions:
  - Whether PO rows with blank `po_no` but valid `order_no` should be included. This should be answered by test/reproduction evidence during implementation.
- Readiness: ready for bounded TDD implementation/verification by `@fixer`; no user answer required before starting.
- Cleanup performed: none required; evidence is intentionally kept for implementation handoff.
