# Discovery — SX-2214 Final Invoice Total dari Final Order

Task ID: `20260617-1236-sx-2214-final-invoice-total`  
Mode: Maintenance Stability Mode  
Tanggal: 2026-06-17 Asia/Jakarta

## Files inspected

Repo docs:

- `AGENTS.md`
- `.opencode/docs/index.md`
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `.opencode/docs/SECURITY.md`
- `.opencode/docs/PROJECT_STACK.md` — tidak ada.
- `.opencode/docs/PROJECT_COMMANDS.md` — tidak ada.

Invoice BE path:

- `sales/controller/invoice_controller.go`
- `sales/service/invoice_service.go`
- `sales/service/invoice_amount.go`
- `sales/service/invoice_service_test.go`
- `sales/repository/invoice_repository.go`
- `sales/model/invoice.go`
- `sales/model/invoice_detail.go`
- `sales/entity/invoice.go`

Reference docs eksternal lokal dari prompt:

- `/Volumes/External/Downloads/BrowserOS/prompts/sx-2131-docs/db_doc.txt`
- `/Volumes/External/Downloads/BrowserOS/prompts/sx-2131-docs/be_doc.txt`
- `/Volumes/External/Downloads/BrowserOS/prompts/sx-2131-docs/fe_doc.txt`

Existing artifact terkait:

- `.opencode/plans/20260615-1530-sx-2214-invoice-final-order.md`
- `.opencode/evidence/20260615-1530-sx-2214-invoice-final-order/discovery.md`

## Project patterns found

- Repo multi-module Go; target module `sales` punya `sales/go.mod`.
- Repo bukan git worktree lokal: `git status --short` gagal dengan `fatal: not a git repository (or any of the parent directories): .git`.
- Validasi repo docs: jalan dari module target dengan `rtk go test ./...` atau targeted `rtk go test ./service -run TestName`.
- Route invoice utama ada di `sales/controller/invoice_controller.go`:
  - `GET /v1/invoices` → `InvoiceService.List`
  - `GET /v1/invoices/details` → `InvoiceService.Details`
  - `GET /v1/invoices/:ro_no` → `InvoiceService.Detail`
  - `POST /v1/invoices/` → `InvoiceService.BulkUpdate`
  - `PATCH /v1/invoices/print/:invoice_no` → print flag only
- Tidak ada generator PDF/download final invoice eksplisit ditemukan di `sales`; klaim PDF fixed harus diverifikasi lewat source consumer atau API trace.

## Root cause evidence

Root cause lama, dari plan/evidence 20260615 dan file path saat ini:

1. Detail invoice lama memakai Sales Order fields dari `sls.order_detail`:
   - `qty1`, `qty2`, `qty3`
   - `sell_price1`, `sell_price2`, `sell_price3`
   - `amount`, `disc_value`, `vat_value`
2. List invoice lama memakai header stale dari `sls.order`:
   - `sub_total`, `disc_value`, `promo_value`, `vat_value`, `total`
3. Generate invoice lama tidak recompute final header total saat `InvoiceService.BulkUpdate`.

Current local code sudah punya fix candidate:

- `sales/model/invoice_detail.go` `InvoiceDetRead` sudah punya:
  - `Qty1Final`, `Qty2Final`, `Qty3Final`
  - `Qty4Final`, `Qty5Final`
  - `SellPriceFinal1`, `SellPriceFinal2`, `SellPriceFinal3`
  - `AmountFinal`
  - `DiscValueFinal`
  - `PromoFinal1..5`
  - `VatValueFinal`
- `sales/service/invoice_amount.go` sudah punya shared helper:
  - `calculateInvoiceFinalLineAmount(detail model.InvoiceDetRead)`
  - `calculateInvoiceFinalTotals(details []model.InvoiceDetRead)`
  - Formula helper memakai final qty/final price, nil-safe promo final per field, discount final, VAT final.
- `sales/service/invoice_service.go` sudah memakai helper di:
  - `Detail`: lines 56-61, totals response dari final details.
  - `List`: lines 120-130, list totals dari final details.
  - `Details`: lines 173-183, details totals dari final details.
  - `BulkUpdate`: lines 279-294, generated invoice writes `SubTotalFinal`, `PromoValueFinal`, `DiscValueFinal`, `VatValueFinal`, `TotalFinal`.
  - `mapInvoiceFinalDetailResponse`: lines 359-381, response qty/price/amount/net memakai final fields.
- `sales/service/invoice_service_test.go` sudah punya regression coverage:
  - `TestInvoiceFinalLineAmountUsesFinalOrderFieldsAndNullPromo`
  - `TestInvoiceFinalLineAmountSubtractsAllFinalPromoFields`
  - `TestInvoiceDetailUsesFinalOrderFieldsAndHeaderTotals`
  - `TestInvoiceListUsesFinalOrderTotals`
  - `TestInvoiceDetailsUsesFinalOrderTotals`
  - `TestInvoiceBulkUpdatePersistsFinalHeaderTotals`

## Validation run

Command run:

```bash
rtk go test ./service -run 'TestInvoice'
```

Working directory:

```text
/Users/ujang/Projects/Geekgarden/scylla-be/sales
```

Result:

```text
Go test: 8 passed in 1 packages
```

## Remaining verification gaps

- Full `rtk go test ./...` belum dijalankan; perlu executor/quality gate karena scope bisa besar.
- Manual staging API belum dijalankan; token secure/local env diperlukan dan tidak boleh ditulis ke artifact.
- PDF/download final invoice path belum ditemukan di `sales`; perlu trace consumer before claim.
- Karena repo lokal bukan git worktree, diff terhadap baseline tidak tersedia. Executor perlu memastikan perubahan ada di target branch sebenarnya.

## Source strategy

Used:

- Repo-local docs and code discovery.
- Existing SX-2214 plan/evidence in `.opencode/`.
- User-provided Jira summary, formulas, sample evidence, guardrails.
- Local extracted docs from BrowserOS paths.

Skipped:

- Context7/official docs: tidak perlu; defect business logic repo-local.
- GitHub/web search: tidak perlu; private repo/local issue.
- Browser capture: tidak perlu untuk planner; manual API/PDF verification planned by executor with secure token.

## Reuse candidates

- Keep `sales/service/invoice_amount.go` helper as single calculation source.
- Keep current invoice service mapping pattern instead of duplicating SQL aggregates.
- Reuse existing tests in `sales/service/invoice_service_test.go` for regression.
- Reuse 20260615 SX-2214 plan as historical root-cause artifact, but current plan is latest source of truth.

## Constraints and risks

- No hardcoded `SO2606100015`, `C260020001`, token, password, bearer.
- Preserve Purchase Order, Sales Order, Proforma behavior.
- Preserve service-layer transaction in `InvoiceService.BulkUpdate`.
- Maintain tenant filter `cust_id` in repository queries.
- Currency uses existing `float64` style; avoid new dependency unless project decides broader money refactor.
- Do not claim PDF fixed until PDF route/consumer confirmed uses fixed API data.
