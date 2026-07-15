# Plan — SX-1879 Export Issue Data

## Goal

Perbaiki backend export order Excel untuk `SX-1879` agar endpoint `GET /sales/v1/download` menghasilkan sheet `Purchase Order`, `Sales Order`, dan `Final Order` yang null-safe, memakai mapping financial field yang benar, dan menampilkan amount dengan separator ribuan Indonesia.

## Non-goals

- Tidak mengubah flow FE, UI `Sales > Order List`, atau endpoint contract selain isi file export/report yang dihasilkan.
- Tidak menambahkan credential/token staging ke repo.
- Tidak melakukan rewrite besar pada repository query atau arsitektur async report kecuali terbukti diperlukan.
- Tidak mengubah schema database atau migration.

## Scope

- Module: `sales`.
- Endpoint/controller: `sales/controller/so_controller.go` untuk tracing dan manual verification.
- Service/export: `sales/service/so_service.go`.
- Tests: `sales/service/so_service_test.go`, dan test helper baru bila perlu tetap di package `service`.
- Repository query: `sales/repository/so_repository.go` hanya jika formula membutuhkan kolom tambahan yang sudah tidak tersedia di model saat ini.

## Requirements

- Export tetap berhasil walaupun field nullable di data PO/SO/Final bernilai `NULL`.
- Amount nullable pada kolom currency/export financial harus fallback ke `0`, bukan blank/error.
- Text/date nullable tetap fallback ke empty string.
- Kolom financial untuk Sales Order dan Final Order harus sesuai QA untuk fixture setara `SO2604290009`:
  - `Discount`: `0`
  - `Net Sales (Exc PPN)`: `12.000.000`
  - `PPN`: `1.100.000`
  - `Gross`: `12.000.000`
- Amount/currency export harus memakai Indonesian thousands separator, contoh `11000000 -> 11.000.000`.
- Perubahan harus minimal, teruji, dan tidak menurunkan tenant filter atau order export lain.

## Acceptance Criteria

1. `GET /sales/v1/download?start_date=1777420800&end_date=1777593599&salesman_id[]=421` membuat report ready di environment yang valid tanpa hardcoded token.
2. Sheet `Purchase Order` tidak gagal saat data null; field teks kosong menjadi `""`, amount null menjadi `0`.
3. Sheet `Sales Order` dan `Final Order` memiliki kolom `Discount`, `Net Sales (ExcPPN)`, `PPN`, dan `Gross` sesuai expected untuk fixture QA/setara.
4. Amount export tampil dengan titik separator ribuan Indonesia.
5. Automated tests untuk formatter dan mapping/export row lulus.
6. `go test ./...` di module `sales` lulus atau blocker dicatat dengan error spesifik.

## Existing Patterns/Reuse

- Reuse flow existing: controller → service → repository → DB.
- Reuse `generateDownloadSalesOrderExcel` dan empat sheet writer di `sales/service/so_service.go`.
- Reuse `excelize` (`github.com/xuri/excelize/v2 v2.9.1`) yang sudah ada di `sales/go.mod`.
- Reuse `so_service_test.go` mocks (`mockSoRepository`, `mockReportRepository`) dan helper pembacaan Excel base64.
- Tidak ditemukan helper formatter amount Indonesia existing di module `sales`; buat helper kecil lokal di `so_service.go` atau file helper service baru bila perlu.

## Constraints

- Repo adalah multi-module Go; jalankan command dari `sales/` untuk test module.
- Instruksi repo meminta `rtk` untuk command; mandatory compose check sudah dijalankan.
- `sales` container tidak sedang up saat discovery; manual endpoint verification mungkin perlu `rtk docker compose -f docker-compose.yml up -d sales` atau `rtk go run main.go` dari `sales/` dengan env valid.
- Endpoint async: response awal berisi metadata report; file Excel tersimpan sebagai base64 di report list/history, sehingga verifikasi harus mengambil report output setelah status `Ready`.
- Jangan log/commit credential atau data sensitif.

## Risks

- Formula financial mungkin memiliki definisi domain yang tidak terdokumentasi. Discovery menunjukkan kemungkinan bug mapping, tetapi implementation harus mempertahankan perubahan minimal dan mengunci expected QA dengan test.
- Existing tests saat ini mengharapkan mapping lama yang mencurigakan (`DiscValueFinal` sebagai `Promotion`, `VatValueFinal` sebagai `Discount`); test tersebut harus diubah sesuai behavior baru, bukan dihapus tanpa pengganti.
- Export sebagai string berformat `11.000.000` memenuhi QA visual, tetapi dapat mengubah tipe cell untuk konsumen Excel. Jika tipe numeric penting, gunakan numeric cell + custom number format dan validasi tampilan Excel. Karena QA meminta titik separator pada value, opsi paling aman untuk defect ini adalah helper formatted string khusus amount columns, kecuali product/QA meminta numeric cell type.
- QA expected `PPN = 1.100.000` dan `Gross = 12.000.000` tidak konsisten dengan `gross = net + ppn`; kemungkinan `Gross` di export berarti amount sebelum PPN. Jangan rename header; sesuaikan mapping export sesuai expected.

## Decisions/Assumptions

- Interaction level: Assumption-first. Requirements cukup jelas untuk membuat plan; tidak perlu question gate sebelum implementasi, tetapi beberapa asumsi harus divalidasi saat coding/manual test.
- Asumsi formula minimal berdasarkan field existing:
  - `grossSales` tetap dihitung dari `qty * sell_price` sesuai sheet masing-masing.
  - `discount` harus berasal dari `disc_value_final` dengan fallback `0`.
  - `ppn` harus berasal dari stored `vat_value_final` dengan fallback `0`, bukan hasil rekalkulasi dari persen `vat`.
  - `netSales` dan `gross` untuk expected QA harus tidak dikurangi `vat_value_final`; kemungkinan `grossSales - discount/promotion` atau langsung `grossSales` bila tidak ada diskon/promosi.
  - `promotion` saat ini belum punya field khusus selain `disc_value_final`; jangan isi promotion dengan VAT. Jika tidak ada sumber promosi terpisah, pertahankan `0` atau sumber existing yang terbukti benar setelah inspect domain/order detail.
- Open question rendah-risiko untuk implementer/QA: apakah Excel amount harus string berformat atau numeric cell dengan display format? Plan merekomendasikan string berformat untuk memenuhi expected QA eksplisit, dengan catatan risiko tipe cell.

## TDD/Test Plan

### TDD Required

Ya. Ini bug production logic dan output contract export, sehingga TDD/regression test wajib.

### Reason

- Mencegah regresi mapping financial PO/SO/Final.
- Mencegah nilai null kembali menjadi blank/error.
- Mencegah separator ribuan hilang di Excel output.

### Existing Test Patterns

- `sales/service/so_service_test.go` sudah memakai package `service`, mock repository, `excelize.NewFile`, dan `excelize.OpenReader` untuk verifikasi output.
- Tambahkan test di file yang sama atau file `so_download_export_test.go` dalam package `service`.

### First Failing/Regression Tests (Red)

1. `TestFormatDownloadAmount_UsesIndonesianThousandsSeparator`
   - Cases: `nil -> "0"`, `0 -> "0"`, `11000000 -> "11.000.000"`, `12000000 -> "12.000.000"`, `1100000 -> "1.100.000"`.
   - Di Go, karena tidak ada `undefined`, gunakan `nil *float64` dan value pointer/non-pointer helper sesuai desain.
2. `TestMapSoToEntity_MapsFinancialFieldsForSX1879`
   - Input: `Qty1=1`, `SellPrice1=12000000`, `DiscValueFinal=0`, `VatValueFinal=1100000`, `Vat=11`.
   - Expected row: `Discount=0`, `NetSales=12000000`, `Vat=1100000`, `Gross=12000000`.
3. `TestMapFinalToEntity_MapsFinancialFieldsForSX1879`
   - Fixture serupa memakai `Qty1Final` dan `SellPriceFinal1`.
4. `TestCreateSalesOrderSheet_FormatsAmountColumnsWithIndonesianSeparator`
   - Build sheet dengan row amount values, read cells `Z4:AE4` atau exact amount columns, assert strings include `12.000.000`, `1.100.000`, `0`.
5. `TestCreatePurchaseOrderSheet_DefaultsNullableAmountColumnsToZero`
   - Row with nil amount pointers should produce `0` for amount cells, not blank.

### Green Step

- Implement helper amount formatting and apply it only to amount/currency columns in PO/SO/Final sheet writers:
  - Selling prices: columns `Q:V`.
  - Financial amount columns: `Z:AE` (`GrossSales`, `Promotion`, `Discount`, `Net Sales (ExcPPN)`, `PPN`, `Gross`).
  - Quantity columns remain numeric/no thousands format unless QA requests otherwise.
- Correct mapper logic for PO/SO/Final financial fields:
  - Stop using `VatValueFinal` as `Discount`.
  - Use `VatValueFinal` as `Vat` output amount.
  - Use `DiscValueFinal` as discount only if existing domain confirms it is discount; otherwise inspect nearby order/list/export behavior before finalizing.
  - Make net/gross calculations match expected QA fixture.
- Update old tests that encoded wrong behavior to assert new mapping.

### Refactor Step

- Remove duplication carefully after tests pass:
  - Extract small helper such as `formatDownloadAmountFromPointer(*float64) string` and maybe `calculateDownloadFinancials(...)`.
  - Keep changes local and readable; do not introduce new dependency.
- Ensure helpers have clear names and avoid broad abstraction across unrelated services.

### Edge Cases

- Nil `DiscValueFinal`, `VatValueFinal`, `Vat`, quantity, and price pointers.
- Zero and negative adjustments, if domain allows returns/discount corrections.
- Decimal values: decide whether to round/no decimals. QA wants no decimal fraction; use zero fraction unless business requires otherwise.
- Multiple item rows for same SO; ensure line-level values do not double-count if repository values are already line-level totals.
- Missing PO number filtering currently excludes blank PO rows; ensure this is still intended for Purchase Order tab.

### Commands

Run from `sales/`:

```bash
rtk go test ./service -run 'TestFormatDownloadAmount|TestMap.*SX1879|TestCreate.*Amount|TestCreatePurchaseOrderSheet_DefaultsNullableAmountColumnsToZero'
rtk go test ./...
```

If `rtk` is unavailable in executor context, run equivalent direct `go test` and record why.

## Implementation Steps

1. Inspect nearby domain usage of `disc_value_final`, `vat_value_final`, `vat`, and order totals in `sales/model/order_detail.go`, `sales/model/order.go`, and service detail/list mapping to confirm formula semantics.
2. Add Red tests listed above. Run targeted tests and confirm failure.
3. Add amount formatter helper in `sales/service/so_service.go` near existing dereference helpers, or in a small `sales/service/so_download_formatter.go` if it improves readability.
   - Prefer no new package dependency.
   - Implement Indonesian grouping deterministically, e.g. integer rounding + inserting `.` every three digits, or use `golang.org/x/text/message` if already available indirectly and acceptable. Simpler local helper is preferred.
4. Correct financial mapping in `mapPoToEntity`, `mapSoToEntity`, and `mapFinalToEntity`.
   - Replace `discount = VatValueFinal` with correct field.
   - Replace computed `vat = netSales * vat%` with stored `VatValueFinal` if present.
   - Adjust `netSales`/`gross` to satisfy QA and domain semantics.
5. Apply formatter in `createPurchaseOrderSheet`, `createSalesOrderSheet`, and `createFinalOrderSheet` for amount columns.
   - Ensure nil amount output is `"0"`.
   - Keep non-amount text/date columns as empty string fallback.
6. Update existing tests that assert old mapping (`TestMapPoToEntity_UsesDiscountAsPromotionAndKeepsPriceOrder`, `TestMapSoToEntity_KeepsUnitPriceOrderAndUsesDiscountAsPromotion`, `TestMapFinalToEntity_KeepsUnitPriceOrderAndUsesDiscountAsPromotion`) to new expected labels/values.
7. Run targeted tests, then full `rtk go test ./...` in `sales/`.
8. If env/token available, perform manual staging/local verification:
   - Start sales service if needed.
   - Call endpoint with QA query using environment-provided authorization.
   - Wait for report `Ready` and inspect generated Excel/report history.
   - Verify `SO2604290009` or fixture-equivalent row values in Sales Order and Final Order sheets.

## Expected Files to Change

- `sales/service/so_service.go`
  - Formatter helper.
  - Financial mapping correction.
  - Sheet amount formatting and null fallback.
- `sales/service/so_service_test.go`
  - New regression tests and update tests that reflect old mapping.
- Optional only if needed:
  - `sales/service/so_download_formatter.go`
  - `sales/service/so_download_formatter_test.go`
  - `sales/repository/so_repository.go` and `sales/model/so_download.go` if formula requires additional selected columns after domain inspection.

## Agent/Tool Routing

- Implementation: route to `@fixer` / `opencode-fixer` for bounded code edits and TDD.
- Architecture/formula review if uncertainty persists: route to `@oracle` after local domain inspection, especially if `gross`, `net`, `promotion`, and `discount` semantics conflict.
- Security/privacy reviewer is not required unless staging tokens/logs or sensitive export data are handled beyond local env.
- Release engineer is not required for code-only defect unless deployment/migration concerns emerge.

## Validation Commands

From repo root before runtime verification if services needed:

```bash
rtk docker compose -f docker-compose.yml ps
rtk docker compose -f docker-compose.yml up -d sales
```

From `sales/`:

```bash
rtk go test ./service -run 'TestFormatDownloadAmount|TestMap.*SX1879|TestCreate.*Amount|TestCreatePurchaseOrderSheet_DefaultsNullableAmountColumnsToZero'
rtk go test ./...
```

Manual endpoint verification, using environment-provided token only:

```bash
rtk curl -H "Authorization: Bearer $SCYLLA_STAGING_TOKEN" \
  "$SCYLLA_STAGING_SALES_BASE_URL/v1/download?start_date=1777420800&end_date=1777593599&salesman_id[]=421"
```

Then inspect generated report file from the existing download-history/report retrieval flow.

## Evidence Requirements

- Automated test output for targeted tests and full `go test ./...`.
- If manual staging/local verification is possible:
  - Request URL without token value.
  - Report ID/name and status transition to `Ready`.
  - Excel inspection notes for sheets and columns; include values but avoid sensitive customer/token data.
- If staging data is inaccessible, provide fixture-based test evidence and note blocker.
- Official docs/context7 not required unless implementing numeric cell style with custom Excel format; if used, record doc source.
- GitHub/web/browser not required for this local backend defect.

## Done Criteria

- Plan implemented with minimal source changes in `sales` module.
- Regression tests demonstrate null safety, formatter behavior, and corrected financial mapping.
- Full `sales` tests pass or known unrelated failures are documented.
- Manual verification evidence is attached when credentials/env are available.
- PR description references `SX-1879` and summarizes root cause, changes, and verification.

## Final Planning Summary

- Artifacts created:
  - Primary source of truth: `.opencode/plans/20260504-2141-sx-1879-export-issue-data.md`.
  - Discovery evidence kept: `.opencode/evidence/20260504-2141-sx-1879-export-issue-data/discovery.md` because it contains inspected files, root-cause candidates, and research-gate notes useful for implementer.
- Key decisions:
  - Use TDD/regression-first implementation.
  - Keep changes localized in `sales/service/so_service.go` and `sales/service/so_service_test.go` unless additional domain columns are required.
  - Prefer formatted string amount output for QA-visible Indonesian separator, unless implementation validation requires numeric cell display formatting.
- Questions asked: none; requirements were sufficient for a concrete plan.
- Assumptions:
  - `vat_value_final` is the stored PPN amount and should not be exported as discount.
  - `disc_value_final` is the discount amount unless further code inspection proves it is promotion.
  - `Gross` expected by QA represents the pre-PPN sales amount for this export.
- Remaining open questions:
  - Whether downstream consumers require numeric Excel cell type instead of formatted string.
  - Whether `Promotion` has a separate source field not currently selected by repository.
- Readiness: ready for bounded implementation by `@fixer` with TDD.
- Cleanup performed: no draft artifacts were created. Evidence was intentionally kept for handoff.
