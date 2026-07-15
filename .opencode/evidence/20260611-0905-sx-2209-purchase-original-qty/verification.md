# Verification SX-2209

Task ID: `20260611-0905-sx-2209-purchase-original-qty`
Tanggal: 2026-06-11 Asia/Jakarta

## Changed files

- `sales/service/order_service.go`
- `sales/service/order_service_test.go`

## Implementation summary

- Menambahkan helper `hasPurchaseDisplayQty(detail model.OrderDetailRead) bool` untuk rule display Purchase Order:
  - `qty_po1/2/3 > 0`, atau
  - `original_qty_po1/2/3 > 0`.
- Menambahkan helper `shouldIncludePurchaseDetailRow(detail model.OrderDetailRead) bool` dengan guard `ItemType == 2` tetap excluded dan fallback existing Sales Order tetap dipertahankan.
- Mengubah append `response.PurchaseDetails.Normal` di `DetailV2` agar memakai `shouldIncludePurchaseDetailRow(detail)`.
- Tidak mengubah `activeQtyForTab` global.
- Tidak mengubah nilai `qty_po*` dari `original_qty_po*`.
- Tidak mengubah stock mutation/process order/promo/VAT/amount formula.

## Regression test

Ditambahkan:

- `TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero`

Coverage:

- `order_detail_id = 7273`
- `pro_id = 10813`
- `qty_po1/2/3 = 0`
- `original_qty_po3 = 3`
- Sales/final qty nil/zero supaya row tidak lolos dari tab lain.
- Expected: satu row muncul di `PurchaseDetails.Normal`, `qty_po3` tetap zero, `original_qty_po3` tetap 3.

## Go validation

Commands dijalankan dari `sales` module:

```bash
rtk go test ./service -run 'TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero'
rtk go test ./service -run 'TestDetailV2_PurchaseDetails'
rtk go test ./service -run 'TestDetailV2'
rtk go test ./...
```

Results:

- `TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero`: `1 passed in 1 packages`
- `TestDetailV2_PurchaseDetails`: `2 passed in 1 packages`
- `TestDetailV2`: `11 passed in 1 packages`
- Full sales module: `245 passed in 22 packages`

Red evidence dari implementer:

- Test baru awalnya gagal sebelum fix dengan failure: `expected one purchase detail row from original_qty_po values only`.

## Runtime preflight

Command dari repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

Result:

- Compose services aktif, termasuk `scylla-sales` pada port `9004`.
- `rtk` memberi warning untrusted project filters, tetapi command tetap berjalan tanpa filter.

## Local DB validation

Database connection:

```bash
PGPASSWORD='postgres' psql -h localhost -p 5432 -U postgres -d ggn_scyllax -v ON_ERROR_STOP=1 -c "SELECT current_database() AS db, current_user AS user_name;"
```

Result:

```text
 db          | user_name
-------------+-----------
 ggn_scyllax | postgres
```

Schema check:

```sql
SELECT column_name
FROM information_schema.columns
WHERE table_schema='sls'
  AND table_name='order_detail'
  AND column_name IN ('order_detail_id','ro_no','order_no','pro_id','original_qty_po1','original_qty_po2','original_qty_po3','qty_po1','qty_po2','qty_po3')
ORDER BY column_name;
```

Result confirms local DB has:

- `order_detail_id`
- `ro_no`
- `pro_id`
- `original_qty_po1/2/3`
- `qty_po1/2/3`

Sample row validation:

```sql
SELECT od.order_detail_id,
       od.ro_no,
       od.pro_id,
       od.original_qty_po1,
       od.original_qty_po2,
       od.original_qty_po3,
       od.qty_po1,
       od.qty_po2,
       od.qty_po3,
       CASE
         WHEN COALESCE(od.qty_po1,0) > 0
           OR COALESCE(od.qty_po2,0) > 0
           OR COALESCE(od.qty_po3,0) > 0
           OR COALESCE(od.original_qty_po1,0) > 0
           OR COALESCE(od.original_qty_po2,0) > 0
           OR COALESCE(od.original_qty_po3,0) > 0
         THEN true ELSE false
       END AS should_display_purchase
FROM sls.order_detail od
WHERE od.ro_no = 'SO2606100013'
ORDER BY od.order_detail_id;
```

Result:

```text
 order_detail_id |    ro_no     | pro_id | original_qty_po1 | original_qty_po2 | original_qty_po3 | qty_po1 | qty_po2 | qty_po3 | should_display_purchase
-----------------+--------------+--------+------------------+------------------+------------------+---------+---------+---------+-------------------------
            7273 | SO2606100013 |  10813 |                0 |                0 |                3 |       0 |       0 |       0 | t
            7274 | SO2606100013 |  10743 |                0 |                0 |                2 |       0 |       0 |       2 | t
            7275 | SO2606100013 |   8438 |                2 |                0 |                1 |       2 |       0 |       1 | t
```

Conclusion:

- Local DB `ggn_scyllax` contains exact SX-2209 sample row.
- Row `7273` has `qty_po1/2/3 = 0` and `original_qty_po3 = 3`.
- New display predicate marks row `7273` eligible for `purchase_details.normal[]`.

## API smoke

Authenticated local Docker smoke was run against direct local services. Token was used only in-memory and was not written to files.

Login endpoint:

```text
POST http://localhost:9001/v1/users/login
status=200
token_present=true
cust_id=C260020001
parent_cust_id=C26002
user_id=141
```

Sales endpoint:

```text
GET http://localhost:9004/v2/orders/SO2606100013
status=200
purchase_normal_count=3
row_7273_present=true
```

Sanitized row evidence:

```json
{
  "order_detail_id": 7273,
  "pro_id": 10813,
  "qty_po1": 0,
  "qty_po2": 0,
  "qty_po3": 0,
  "original_qty_po1": 0,
  "original_qty_po2": 0,
  "original_qty_po3": 3
}
```

## Diff boundary

Changed files are within allowed boundary:

- `sales/service/order_service.go`
- `sales/service/order_service_test.go`

Out-of-boundary source changes: none.

Mechanical evidence for changed files:

```text
993078ea10d83b70393983022040fd6d2a5291a02a33d22cfde4f7ea01985f5b  sales/service/order_service.go
2561c5219f4ad791335f7a6a55cc10965fde9f2aabf00390b5ed682f2bf066c3  sales/service/order_service_test.go
```

## Git status

Repo root is not detected as a git repository by `git rev-parse --show-toplevel` / `git status --short`:

```text
fatal: not a git repository (or any of the parent directories): .git
```

No local commit was created.

## Known diagnostics

Tooling previously reported LSP/editor diagnostics in existing files around `OrderType` and report repository signatures. These did not block `rtk go test ./...`, which passed.

## Plan compliance checkpoint

- `qty_po*` not overwritten from `original_qty_po*`: pass.
- `original_qty_po*` only used for purchase row display decision: pass.
- No stock mutation/process order change: pass.
- No promo/VAT/amount formula change: pass.
- `details.normal[]` and `details_final.normal[]` semantics unchanged: pass.
- Diff boundary respected: pass.
- Local DB `ggn_scyllax` sample validated: pass.
- Tests passed: pass.

## Quality gate

Final quality gate result: `PASS`.

Notes:

- Remediation from initial `PASS_WITH_RISKS` was completed by adding mechanical evidence for non-git workspace.
- Residual low risk: no token-based live API smoke; covered by source review, targeted regression, full module tests, and local DB sample evidence.
