# Verification — SX-2246 Add Purchase Details

Task ID: `20260619-1754-sx-2246-add-purchase-details`
Date: 2026-06-19

## Scope executed

Allowed files changed:

- `sales/entity/edit_order_enhance.go`
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`

Temporary validation file created and removed:

- `sales/service/sx2246_localdb_validate_test.go` with build tag `sx2246_validate`

No controller/repository/model/schema changes.

## Implementation notes

- `entity.AddPurchaseOrderDetail` now includes:
  - `UnitId1 string json:"unit_id1,omitempty"`
  - `UnitId2 string json:"unit_id2,omitempty"`
  - `UnitId3 string json:"unit_id3,omitempty"`
  - `IsProductPromotionPo *bool json:"is_product_promotion_po,omitempty"`
- `createOrderDetailFromPurchaseOrder` now:
  - preserves raw request qty into `original_qty_po1/2/3`.
  - resolves stock with `StockRepository.GetCurrentStock` when repository is available.
  - converts stock with `canonicalAPIStockBreakdown`.
  - caps `qty_po1/2/3`, `qty1/2/3`, `qty1_final/2_final/3_final` per UOM level.
  - writes `qty1_stok/2_stok/3_stok` at insert.
  - uses payload `unit_id1/2/3` when non-empty, otherwise product master fallback.
  - writes payload `is_product_promotion_po`, default `false`.
- `UpdateEnhance` add purchase loop now returns explicit errors when `wh_id` or `ro_date` is missing before add insert.
- No destructive dedupe added.
- Existing `add_purchase_details` alias remained in `normalizeEnhancePromoFlags`.

## TDD / unit validation

Command:

```bash
cd /Users/ujang/Projects/Geekgarden/scylla-be/sales
rtk go test ./service -run 'TestCreateOrderDetailFromPurchaseOrder|TestUpdateEnhance.*AddPurchase' -count=1
```

Result:

```text
Go test: 7 passed in 1 packages
```

Command:

```bash
cd /Users/ujang/Projects/Geekgarden/scylla-be/sales
rtk go test ./... -count=1
```

Result:

```text
Go test: 290 passed in 22 packages
```

Command:

```bash
cd /Users/ujang/Projects/Geekgarden/scylla-be/sales
rtk go vet ./entity/... ./service/...
```

Result:

```text
Go vet: No issues found
```

Gofmt check:

```bash
rtk gofmt -l entity/edit_order_enhance.go service/order_service.go service/order_service_test.go
```

Result: no output.

## Local database validation (`ggn_scyllax`)

DB access used:

```text
host=127.0.0.1 port=5432 user=postgres password=postgres dbname=ggn_scyllax sslmode=disable
```

Schema compatibility verified:

```bash
PGPASSWORD=postgres psql -h 127.0.0.1 -p 5432 -U postgres -d ggn_scyllax -c "\d sls.order_detail"
```

Relevant columns exist:

- `qty1`, `qty2`, `qty3`
- `qty1_final`, `qty2_final`, `qty3_final`
- `qty_po1`, `qty_po2`, `qty_po3`
- `original_qty_po1`, `original_qty_po2`, `original_qty_po3`
- `qty1_stok`, `qty2_stok`, `qty3_stok`
- `unit_id1`, `unit_id2`, `unit_id3`
- `is_product_promotion_po`

Sample local products verified:

```sql
SELECT pro_id, pro_name, conv_unit2, conv_unit3, unit_id1, unit_id2, unit_id3, vat, sell_price1, sell_price2, sell_price3
FROM mst.m_product
WHERE pro_id IN (10813, 8436)
ORDER BY pro_id;
```

Observed:

```text
8436  Jersey Medan Chief        conv_unit2=5 conv_unit3=1 unit_id1=PCS unit_id2=CTN unit_id3=CTN
10813 Jersey Manchester United  conv_unit2=5 conv_unit3=1 unit_id1=PCS unit_id2=CTN unit_id3=CTN
```

Temporary DB validation test:

- Created build-tagged test `service/sx2246_localdb_validate_test.go`.
- Seeded temporary order `SX2606240001` in `sls."order"` with `data_status=1` (`Need Review`) under `cust_id=C260020001`, `wh_id=350`.
- Seeded deterministic `inv.warehouse_stock` rows for products `10813` and `8436`.
- Called `createOrderDetailFromPurchaseOrder` against real repositories and DB.
- Queried `sls.order_detail` back and asserted:
  - raw `original_qty_po*` persisted.
  - capped `qty_po*` values persisted.
  - payload `unit_id*` persisted.
  - payload `is_product_promotion_po` persisted.
  - stock snapshot fields set.
- Cleanup removed all temporary `sls.order`, `sls.order_detail`, and `inv.warehouse_stock` rows.
- Temporary validation file removed after run.

Command:

```bash
cd /Users/ujang/Projects/Geekgarden/scylla-be/sales
SX2246_DSN='host=127.0.0.1 port=5432 user=postgres password=postgres dbname=ggn_scyllax sslmode=disable' rtk go test -tags=sx2246_validate ./service -run TestSX2246_ValidateAgainstLocalDB -count=1 -v
```

Result:

```text
Go test: 3 passed in 1 packages
```

Cleanup verification:

```sql
SELECT COUNT(*) AS leftover_order_details FROM sls.order_detail WHERE ro_no = 'SX2606240001';
SELECT COUNT(*) AS leftover_order FROM sls."order" WHERE ro_no = 'SX2606240001';
SELECT COUNT(*) AS leftover_stock FROM inv.warehouse_stock WHERE cust_id = 'C260020001' AND wh_id = 350 AND pro_id IN (10813,8436);
```

Result:

```text
leftover_order_details = 0
leftover_order = 0
leftover_stock = 0
```

## Plan compliance checkpoint

- Execution-ready tasks T1-T6 completed.
- T7 quality-gate: `PASS_WITH_RISKS` (see Quality Gate Remediation section below).
- Diff boundary respected for production code.
- Temporary DB validation file was removed before final tests.
- No secrets copied into source files. DB password appears only in verification command text for local `postgres/postgres`, matching repo-local runtime guidance.
- No `.env` files modified.
- No git commit created because workspace root is not a git repository.

## Quality Gate Remediation

Gate verdict: `PASS_WITH_RISKS`. Required-before-PASS item completed.

Item R1 (required_before_PASS, `requires_user_decision: no`):

- Finding: real-DB rollback proof on partial add failure was missing.
- Action: added build-tagged test `sales/service/sx2246_rollback_proof_test.go` using real `WithinTransaction` and a wrapping repo that returns a simulated `FindProductByID` error for the second add, then asserted zero leftover rows.
- Command:

  ```bash
  cd /Users/ujang/Projects/Geekgarden/scylla-be/sales
  SX2246_DSN='host=127.0.0.1 port=5432 user=postgres password=postgres dbname=ggn_scyllax sslmode=disable' rtk go test -tags=sx2246_validate ./service -run TestSX2246_RollsBackOnSecondInsertFailure_RealDB -count=1 -v
  ```

- Result: `Go test: 1 passed in 1 packages`.
- Post-test cleanup query: zero leftover rows for order `SX2606240002`, order header, and stock seeds.
- Temp file removed after run; final test suite re-run: 7/7 targeted pass, 290/290 full sales pass.

Non-blocking follow-ups (deferred):

- F1: stock cap interpretation for mixed-UOM allocation needs FE/product sign-off (`requires_user_decision: yes`).
- F2: duplicate resend / idempotency policy for `add_purchase_details` needs product sign-off (`requires_user_decision: yes`).

## Final status

`PASS_WITH_RISKS` with R1 closed. Implementation matches plan, tests + real DB validation pass, no regression in `purchase_order` update path (covered by 290-test full suite). Diff boundary clean. Two product decisions (F1, F2) deferred to owner.

## Risks / follow-up

- Duplicate resend policy remains undefined; implementation does not delete or dedupe rows.
- Cap algorithm follows plan decision: per-level cap after `canonicalAPIStockBreakdown`. If product owner expects optimal total-smallest-unit allocation across UOM levels, that is a separate requirement.
