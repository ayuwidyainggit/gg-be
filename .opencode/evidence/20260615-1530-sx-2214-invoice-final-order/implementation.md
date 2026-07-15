# Implementation Evidence — SX-2214

Task ID: `20260615-1530-sx-2214-invoice-final-order`
Date: 2026-06-15 Asia/Jakarta

## Changed files

- `sales/model/invoice.go`
- `sales/model/invoice_detail.go`
- `sales/repository/invoice_repository.go`
- `sales/service/invoice_amount.go`
- `sales/service/invoice_service.go`
- `sales/service/invoice_service_test.go`

## Root cause confirmed

- `sales/service/invoice_service.go` previously automapped `InvoiceDetRead` directly to invoice responses, so invoice detail/list detail used Sales Order fields:
  - `qty1`, `qty2`, `qty3`
  - `sell_price1`, `sell_price2`, `sell_price3`
  - `amount`
  - `disc_value`
  - `vat_value`
- `sales/repository/invoice_repository.go` invoice list selected stale header fields such as `sls.order.total`.
- `InvoiceService.BulkUpdate` generated invoice number/status without recomputing final invoice totals from final order detail fields.

## Implementation summary

- Added final invoice amount helper in `sales/service/invoice_amount.go`.
- Formula source:
  - `qty1_final`, `qty2_final`, `qty3_final`
  - `sell_price_final1`, `sell_price_final2`, `sell_price_final3`
  - `promo_final1..5`
  - `disc_value_final`
  - `vat_value_final`
- Helper does not use `amount` or `amount_final` as canonical source.
- Promo sum is nil-safe per field.
- Applied helper to:
  - `InvoiceService.Detail`
  - `InvoiceService.List`
  - `InvoiceService.Details`
  - `InvoiceService.BulkUpdate`
- Detail response now maps final qty/price/amount/discount/VAT/net into invoice-context response fields.
- List/detail response header money fields now reflect final totals.
- `BulkUpdate` persists final header columns:
  - `sub_total_final`
  - `promo_value_final`
  - `disc_value_final`
  - `vat_value_final`
  - `total_final`
- `BulkUpdate` clears stale non-final request monetary fields before update payload.
- `FindAllByInvoiceNombersAndCustId` now aliases `sls.order.ro_no AS order_no` so `Details()` can fetch detail rows by order number.

## Tests added

`sales/service/invoice_service_test.go` covers:

- final line amount uses final fields and null promo fields sum to zero.
- all `promo_final1..5` fields subtract from net.
- `Detail()` uses final order fields and final header totals.
- `List()` uses final order totals over stale header totals.
- `Details()` uses final order totals.
- `BulkUpdate()` persists final header totals and ignores stale request money fields.

## Validation run

Targeted:

```bash
rtk go test ./service -run 'TestInvoice.*Final|TestInvoiceBulkUpdate'
```

Result:

```text
Go test: 8 passed in 1 packages
```

Full sales suite:

```bash
rtk go test ./...
```

Result:

```text
Go test: 258 passed, 3 failed, 1 skipped in 22 packages
```

Failures are unrelated external/flaky translator tests:

- `pkg/texttranslator/translator_test.go:69` expected text mismatch.
- `TestTranslator_Translate/test_Conight/go-googletrans_version` hit `429 Too Many Requests`.
- `TestTranslator_Translate` panicked from expected status `200`, got `429`.

Full output path from runner:

```text
~/Library/Application Support/rtk/tee/1781514245_go_test.log
```

## Revalidation 2026-06-17

Targeted ulang:

```bash
rtk go test ./service -run 'TestInvoice.*Final|TestInvoiceBulkUpdate'
```

Result:

```text
Go test: 8 passed in 1 packages
```

Full sales suite ulang:

```bash
rtk go test ./...
```

Result:

```text
Go test: 262 passed in 22 packages
```

## PDF/download trace

Repo-local trace found no final invoice PDF/download generator in this repo.

Evidence:

- `sales/controller/invoice_controller.go` has JSON routes only:
  - `GET /v1/invoices`
  - `GET /v1/invoices/details`
  - `GET /v1/invoices/:ro_no`
  - `POST /v1/invoices/`
  - `PATCH /v1/invoices/print/:invoice_no`
- `Print()` only updates `is_printed`, `printed_by`, `printed_at` and returns JSON message.
- No `application/pdf`, PDF library, file send, or final invoice download route found repo-wide.
- TMS services consume `/v1/invoices` JSON data, but no local PDF renderer proof.

Claim limit:

- Fixed invoice API source data for final-order totals.
- Cannot truthfully claim direct PDF renderer fixed in this repo.
- If PDF/download consumes fixed `/v1/invoices` APIs, downstream PDF should receive corrected values after deploy.
- If PDF renderer is in FE/BFF/external service, trace/fix there remains follow-up.

## Quality gate

Initial quality gate: `PASS_WITH_RISKS` because full suite had unrelated translator test failures.

Final revalidation quality gate on 2026-06-17: `PASS`.

Evidence:

- targeted `rtk go test ./service -run 'TestInvoice.*Final|TestInvoiceBulkUpdate'` => 8 passed
- full `rtk go test ./...` => 262 passed in 22 packages
- `sales` git status clean
- final formula source, nil-safe promo, invoice list/detail override, transaction-safe final header persistence, and PDF claim limit reviewed

No blocker found in SX-2214 code path.

## Data repair / migration note

- No migration added; local models and plan evidence show final columns already exist.
- Existing generated invoices may need regeneration or DB backfill of `sls.order.*_final` header fields if already generated before this fix and stored final header totals are stale.
- Safer repair strategy: recompute from `sls.order_detail` final fields with the same formula and update only affected invoice rows; do not hardcode Jira sample in repair script.

## Local DB/API verification

Local runtime:

```text
scylla-system: up on 9001
scylla-sales: up on 9004
local DB: ggn_scyllax on localhost:5432
```

DB connection check:

```sql
SELECT current_database() AS db, current_user AS usr;
```

Result:

```text
db = ggn_scyllax
usr = postgres
```

Sample header in local DB for `SO2606100015`:

```text
cust_id = C260020001
invoice_no = INV2606100001
data_status = 6
sub_total = 3.600.000
sub_total_final = 18.600.000
vat_value = 360.000
vat_value_final = 1.860.000
total = 3.960.000
total_final = 20.460.000
```

Final formula DB query result from `sls.order_detail`:

```text
gross = 18.600.000
promo = 0
discount = 0
vat = 1.860.000
expected_invoice_total = 20.460.000
```

Local login/API smoke:

```text
POST http://localhost:9001/v1/users/login => 200
access token present = true
cust_id = C260020001
parent_cust_id = C26002
```

No token/password stored in artifact.

Local invoice API validation:

```text
GET http://localhost:9004/v1/invoices?q=SO2606100015&is_invoice=true&page=1&limit=10 => 200
rows = 1
order_no = SO2606100015
sub_total = 18.600.000
vat_value = 1.860.000
total = 20.460.000
```

Invoice list details:

```text
TP-012 qty = 1,0,1 price = 1.500.000,15.000.000,15.000.000 amount = 16.500.000 vat = 1.650.000 net = 18.150.000
TP-013 qty = 4,0,1 price = 150.000,1.500.000,1.500.000 amount = 2.100.000 vat = 210.000 net = 2.310.000
```

Invoice detail API:

```text
GET http://localhost:9004/v1/invoices/SO2606100015 => 200
sub_total = 18.600.000
vat_value = 1.860.000
total = 20.460.000
normal_details = 2
```

Order V2 Final Order control:

```text
GET http://localhost:9004/v2/orders/SO2606100015 => 200
sub_total_final = 18.600.000
vat_value_final = 1.860.000
total_final = 20.460.000
final_details = 2
```

Conclusion:

- Local DB final formula matches expected Jira evidence.
- Local Invoice List API matches Final Order total.
- Local Invoice Detail API line values use final qty/price and match Final Order amounts.
- Local `GET /v2/orders/{ro_no}` Final Order remains aligned.
