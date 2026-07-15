# Discovery — SX-2214 Total Invoice List Tidak Sama dengan Final Order

Task ID: `20260615-1530-sx-2214-invoice-final-order`
Mode: Maintenance Stability Mode

## Files inspected

Repo docs:

- `AGENTS.md`
- `.opencode/docs/index.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `.opencode/docs/PROJECT_STACK.md` — tidak ada.
- `.opencode/docs/PROJECT_COMMANDS.md` — tidak ada.

Sales invoice path:

- `sales/controller/invoice_controller.go`
- `sales/service/invoice_service.go`
- `sales/repository/invoice_repository.go`
- `sales/model/invoice.go`
- `sales/model/invoice_detail.go`
- `sales/entity/invoice.go`
- `sales/entity/invoice_detail.go`
- `sales/service/invoice_service_concurrency_test.go`

Order/final-order path:

- `sales/model/order.go`
- `sales/model/order_detail.go`
- `sales/service/order_service.go`
- `sales/repository/order_repository.go`

## Commands/docs checked

- `glob("*/go.mod")` confirmed `sales/go.mod` as target module.
- `glob("**/*invoice*")` under `sales` found invoice controller/service/repository/model/entity and concurrency tests.
- `grep` for invoice/download/pdf/final fields under `sales`.
- Repo quality docs require validation in target module with `rtk go test ./...`; runtime starts from root with `rtk docker compose -f docker-compose.yml ps`.

## Project patterns found

- Repo is multi-module Go monorepo. Target module is `sales`.
- Layering required: Controller → Service → Repository → DB.
- Tenant filter required: `cust_id` on transactional reads/writes.
- Main invoice route is `sales/controller/invoice_controller.go`:
  - `GET /v1/invoices` → `InvoiceService.List`
  - `GET /v1/invoices/details` → `InvoiceService.Details`
  - `GET /v1/invoices/:ro_no` → `InvoiceService.Detail`
  - `POST /v1/invoices/` → `InvoiceService.BulkUpdate` generate/update invoice info
  - `PATCH /v1/invoices/print/:invoice_no` → print flag only
- No explicit PDF/download generator code found in `sales` by local grep for `pdf`, `download`, and invoice terms. Likely PDF/download consumer uses invoice list/detail API response or another repo/service not present here. Implementation must verify with FE/BFF or route inventory before claiming PDF fixed.
- Proforma print path is separate in `OrderService.PrintProformaInvoice`, not final invoice PDF path.

## Root cause candidates from local evidence

### Candidate 1 — Invoice detail response maps Sales Order fields

`RepositoryInvoiceImpl.FindDetail` selects raw `sls.order_detail.*` plus product data:

```go
func (repository *RepositoryInvoiceImpl) FindDetail(roNo string, custId string) (details []model.InvoiceDetRead, err error) {
    err = repository.Select(`sls.order_detail.*, ...`).
        Joins("LEFT JOIN mst.m_product p on p.pro_id = sls.order_detail.pro_id").
        Where("sls.order_detail.ro_no = ? AND sls.order_detail.cust_id=?", roNo, custId).
        Find(&details).Error
}
```

But `model.InvoiceDetRead` and `entity.InvoiceDetResponse` expose invoice-facing fields as Sales Order fields:

- `Qty1`, `Qty2`, `Qty3` map `qty1`, `qty2`, `qty3`.
- `SellPrice1`, `SellPrice2`, `SellPrice3` map `sell_price1`, `sell_price2`, `sell_price3`.
- `Amount` maps `amount`.
- `DiscValue` maps `disc_value`.
- `VatValue` maps `vat_value`.
- `NetValue`, `PriceIncludePpn`, `PriceExcludePpn` exist but no final-field calculation found in invoice service.

`InvoiceService.Detail`, `List`, and `Details` automap these raw fields. This matches Jira note: invoice download shows Sales Order quantity.

### Candidate 2 — Invoice list header total maps stale Sales Order header fields

`RepositoryInvoiceImpl.FindAllByCustId` selects invoice list totals from `sls.order` non-final columns:

- `sls.order.sub_total`
- `sls.order.disc_value`
- `sls.order.promo_value`
- `sls.order.vat_value`
- `sls.order.total`

It also selects `promo_value_final` and `promo_bg_value_final`, but not `sub_total_final`, `disc_value_final`, `vat_value_final`, or `total_final`. `model.InvoiceList` also lacks `SubTotalFinal`, `DiscValueFinal`, `VatValueFinal`, `TotalFinal`, so `entity.InvoiceListResponse.Total` remains stale Sales Order total.

### Candidate 3 — Invoice generation only sets invoice number/status, not amount fields

`InvoiceService.BulkUpdate` sets invoice date/status and generated invoice number, then `InvoiceRepository.Update` writes `model.Invoice` derived from request. No recompute of final invoice amount found inside `BulkUpdate`.

Existing `OrderService` final-order update logic already maintains final headers:

- `recomputePromoStateForTab(..., promoSnapshotTabFinalOrder, ...)` writes:
  - `SubTotalFinal`
  - `DiscValueFinal`
  - `VatValueFinal`
  - `TotalFinal`
  - `PromoValueFinal`
- fallback final recompute around `order_service.go:5837-5863` uses final qty/price and final values, but it does not include promo in `totalFinal` in that branch; executor must inspect whether branch can apply to SX-2214 flow.

## Reuse candidates

- `getValueOrDefault` in `sales/service/order_service.go` pattern for nil-safe numeric access.
- `calculateVatValue` in `sales/service/order_service.go` for VAT formula pattern when VAT must be recalculated, though SX-2214 requires using persisted `vat_value_final` as source unless DB inconsistency found.
- `recomputePromoStateForTab` in `sales/service/order_service.go` as source of Final Order header convention.
- `model.OrderDetailRead` already has final fields:
  - `Qty1Final`, `Qty2Final`, `Qty3Final`
  - `SellPriceFinal1`, `SellPriceFinal2`, `SellPriceFinal3`
  - `PromoFinal1..5`
  - `DiscValueFinal`
  - `VatValueFinal`
  - `AmountFinal`
- `model.OrderList` already has final headers:
  - `SubTotalFinal`
  - `PromoValueFinal`
  - `PromoBgValueFinal`
  - `DiscValueFinal`
  - `VatValueFinal`
  - `TotalFinal`
- Existing test pattern for invoice service in `sales/service/invoice_service_concurrency_test.go` uses mocks and `NewInvoiceService`.
- Existing SQL-string regression style exists in `sales/repository/report_repository_test.go`.

## Constraints

- Use `rtk` prefix for shell validation in this repo.
- Validate in `sales` module.
- No token, bearer, Jira credential, staging password, `.env`, or live customer secret in source/test/artifacts.
- Do not hardcode `SO2606100015` or `C260020001` in production code.
- Do not change Purchase Order, Sales Order, or Proforma behavior unless direct dependency proven.
- Preserve service-layer transactions and repository tx-context extraction for writes.
- Currency currently uses `float64` across invoice/order models. Plan should prefer existing project type for minimum diff, unless repo already has decimal helper not found in discovery.

## Risks

- Final invoice PDF/download code may live outside `sales` or outside repository; local BE can only fix API source unless PDF route is found during executor discovery.
- `AmountFinal` may be stale or semantically net/gross depending path. Do not use it as sole source unless DB evidence proves it matches Final Order.
- `promo_value_final` vs `promo_final1..5` convention may differ between header and line. Jira formula says line promo must sum `promo_final1..5`; use that for invoice line/footer source.
- `InvoiceService.BulkUpdate` currently updates stock using `Qty*Final` but `UnitPrice: detail.SellPrice1`; if unit price affects stock valuation, executor should inspect, but stock valuation change is out of scope unless directly linked.
- `InvoiceRepository.FindAllByCustId` sort accepts raw `dataFilter.Sort`; existing risk out of SX-2214 scope.
- If existing global tests fail, executor must isolate targeted tests and document pre-existing failures.

## Source strategy

Used:

- Repo-local docs and code discovery.
- User-provided Jira summary, sample evidence, formulas, and guardrails.

Skipped:

- Official docs/context7: not needed; Go/GORM behavior not version-sensitive for plan.
- GitHub/web search: not needed; defect is repo-local business logic.
- Browser/screenshot capture: not possible/needed in planner; manual verification planned with local/staging API using secure token outside artifacts.

## Open questions

No blocking question. User supplied enough requirements. Remaining uncertainty is implementation discovery item: where final invoice PDF/download is generated if not in `sales`.
