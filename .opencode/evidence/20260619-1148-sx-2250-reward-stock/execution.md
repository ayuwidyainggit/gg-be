# Execution Evidence — SX-2250

Task id: `20260619-1148-sx-2250-reward-stock`
Jira: SX-2250
Related: SX-2253, SX-2090, SX-2129
Target service: `sales/`

## Summary

Feedback terbaru menunjukkan stock movement sudah benar, tetapi available stock harus mengikuti Inventory Stock Report. Root cause ditemukan: Sales `DetailV2` memakai `inv.warehouse_stock.qty` snapshot, sementara Inventory report memakai ledger `inv.stock` (`SUM(qty_in)-SUM(qty_out)`). Production logic sekarang diubah agar `FindWarehouseStockByWhIdAndProIds` memakai basis ledger yang sama dengan Inventory report. Note: tidak ada `stock_date` cutoff di kode saat ini; query memakai seluruh ledger cumulative per `cust_id`/`wh_id`/`pro_id`. Tanggal cutoff untuk historical order menjadi follow-up jika FE/PM minta.

## Tests

Baseline + new tests run from `sales/` with `rtk`:

- `rtk go test ./service -run 'DetailV2|Stock' -count=1` -> 28 passed
- `rtk go test ./service -run 'TestDetailV2_SX2250' -count=1 -v` -> 2 passed
- `rtk go test ./service -run 'TestDetailV2_SX2250|TestDetailV2_SameSKURewardDoesNotContaminateNormalRow|TestDetailV2_Cancelled_UsesWarehouseCurrentOnlyForDisplayedStock|TestDetailV2_NonCancelled_KeepsExistingDisplayedStockBehavior' -count=1` -> 5 passed
- `rtk go test ./service -run 'Stock' -count=1` -> 20 passed
- `rtk go test ./service -count=1` -> 197 passed

New tests added (in `sales/service/order_service_test.go`):

- `TestDetailV2_SX2250_SameProIDNormalAndReward_ComputesStockPerRow` — same `pro_id`, normal row item_type=1 + reward row item_type=2, qty 0 0 1, warehouse current 0, conv 5/1, assert each row in `details.normal` and `details_final.normal` returns `Qty1Stok=0, Qty2Stok=0, Qty3Stok=1`; flags preserved.
- `TestDetailV2_SX2250_SameProIDDifferentQty_ComputesStockPerRow` — same `pro_id`, row A qty 1 -> stock 0 0 3, row B qty 2 -> stock 0 0 4, no aggregation.

## Production logic changes

Changed `sales/repository/order_repository.go::FindWarehouseStockByWhIdAndProIds`:

- Before: read `inv.warehouse_stock.qty` snapshot.
- After: read `inv.stock` ledger with `COALESCE(SUM(st.qty_in),0) - COALESCE(SUM(st.qty_out),0)` grouped by `pro_id`.
- No date filter is applied in this patch; it mirrors cumulative stock basis required by current FE feedback. If historical `ro_date` parity is required later, add explicit `stock_date <= ro_date` scope and tests.
- Reason: Inventory Stock Report (`inventory/repository/stock_repository.go:425-432`) uses the same ledger formula, and FE validation expects Detail Sales Order available stock to match Inventory report.

`DetailV2` stock display remains row-scoped and now uses priority `final > sales > purchase` for `QTY_ORDER` via `stockDisplayQtyByPriority` in `sales/service/order_service.go:2476-2496`:

- sales normal: `sales/service/order_service.go:2946-2982`
- purchase normal: `sales/service/order_service.go:3073-3091`
- final normal: `sales/service/order_service.go:3236-3254`
- promo row dimove ke `.Normal` lewat `movePromoDetailsToNormal` (`sales/service/order_service.go:2467-2474`, dipanggil `2997`/`3257`).
- helper `sales/service/order_stock_helper.go` tetap canonical source-of-truth, signature tidak diubah.

## On-customer stock movement trace

`buildRewardProductStockDeltasFromModels` di `syncRewardProductState` sudah menginklusi reward product ke stock delta, tidak ada exclusion ditemukan. Tidak ada edit, hanya verifikasi kode. T1-T4 done.

## Runtime verification via local stack

Service hidup di docker:

- `scylla-system` port 9001
- `scylla-sales` port 9004
- DB Postgres `ggn_scyllax` di host `127.0.0.1:5432`

### DB query

Order lokal `SO2606190002` punya 2 row pro_id 10752 (sama dengan Jira). Periksa detail:

```sql
select od.order_detail_id, od.pro_id, p.pro_code, od.item_type, od.qty, od.qty1, od.qty2, od.qty3, od.qty1_stok, od.qty2_stok, od.qty3_stok, od.is_product_promotion_so, od.is_product_promotion_final, p.conv_unit2, p.conv_unit3
from sls.order_detail od join mst.m_product p on p.pro_id=od.pro_id
where od.ro_no='SO2606190002' and od.cust_id='C260020001'
order by od.item_type, od.order_detail_id;
```

Hasil:

| order_detail_id | pro_id | pro_code | item_type | qty | qty1_stok (db) | is_product_promotion_so | is_product_promotion_final |
|---|---|---|---|---|---|---|---|
| 7474 | 10752 | AF-006 | 1 | 1 | 4 | false | false |
| 7477 | 10752 | AF-006 | 2 | 1 | (null) | true | true |

Stok ledger current (basis Inventory Stock Report):

```sql
select pro_id, sum(qty_in)-sum(qty_out) as qty
from inv.stock
where cust_id='C260020001' and wh_id=350 and pro_id=10752 and stock_date <= '2026-06-19'
group by pro_id;
```

Hasil: `qty = 0` (AF-006 kosong sesuai catatan FE terbaru). Catatan: `inv.warehouse_stock.qty=8` di local, tapi snapshot itu tidak dipakai lagi setelah patch.

Conv: `conv_unit2=5`, `conv_unit3=1`.

### API call

Login via `POST http://127.0.0.1:9001/v1/users/login` (user lokal, tidak ditulis ke file). Token dipakai `Authorization: Bearer ...` ke `GET http://127.0.0.1:9004/v2/orders/SO2606190002`. Response `details.normal` dan `details_final.normal` untuk 7474/7477:

```json
{"order_detail_id":7474, "pro_id":10752, "item_type":1, "qty1":1, "qty2":0, "qty3":0, "qty1_final":1, "qty2_final":0, "qty3_final":0, "qty1_stok":0, "qty2_stok":0, "qty3_stok":1, "is_product_promotion_so":false, "is_product_promotion_final":false}
{"order_detail_id":7477, "pro_id":10752, "item_type":2, "qty1":1, "qty2":0, "qty3":0, "qty1_final":1, "qty2_final":0, "qty3_final":0, "qty1_stok":0, "qty2_stok":0, "qty3_stok":1, "is_product_promotion_so":true, "is_product_promotion_final":true}
```

Math (unit-basis check):

- API call site passes current row `Qty3/Qty2/Qty1` to `computeDisplayedAvailableStockBreakdown` in that order; helper signature labels them `qtyLarge, qtyMedium, qtySmall`.
- `toTotalSmallFromAPIUnits` converts back to base units: `total = small + medium*conv2 + large*conv2*conv3`.
- For this row, `Qty3=0, Qty2=0, Qty1=1`, conv `5/1`. Note: `Qty1` here means "smallest unit" per the conversion semantics in `sales/pkg/conversion/quantity.go` (largest-count stored as Qty3, remainder as Qty1).
- `toTotalSmallFromAPIUnits(0, 0, 1, 5, 1)` = `0 + 0 + 1 = 1` base unit.
- `computeDisplayedAvailableStockBreakdown(0, 0, 0, 1, true, 5, 1)` -> `totalSmall = 0 + 1 = 1` (warehouse current 0 dari ledger setelah patch).
- `canonicalAPIStockBreakdown(1, 5, 1)` -> `Qty3 = 1/5 = 0`, rem `1`. `Qty2 = 1/5 = 0`, rem `1`. `Qty1 = 1`. Returns `Qty1=converted.Qty3=0, Qty2=0, Qty3=converted.Qty1=1` -> API `qty1_stok=0, qty2_stok=0, qty3_stok=1`. Cocok dengan API dan sesuai expected FE.

Stok DB kolom `sls.order_detail.qty1_stok=4` adalah snapshot lama (4 base units saat order dibuat) dan tidak dipakai di display; helper hitung ulang per request dari `FindWarehouseStockByWhIdAndProIds` (sekarang basis ledger `inv.stock`).

`details.normal` dan `details_final.normal` punya nilai identik dan flags benar. Row normal+reward sama `pro_id` tidak saling kontaminasi (tiap row qty=1 -> API `0 0 1`).

## Runtime verification — SO2606190003 cutoff date feedback

Feedback latest: smallest `104 + qty 1 = 105`, available stock should be `21 0 0` for AF-006.

Local DB check:

```sql
select pro_id, sum(qty_in)-sum(qty_out) as qty
from inv.stock
where cust_id='C260020001' and wh_id=350 and pro_id=10752 and stock_date <= '2026-06-19'
group by pro_id;
```

Hasil: `qty = 104`.

Runtime `GET /v2/orders/SO2606190003` after patch:

```text
details.normal:
7484 AF-006 item_type=1 qty*_stok = 21 0 0, qty*_final = 1 0 0
7487 AF-006 item_type=2 qty*_stok = 21 0 0, qty*_final = 1 0 0

details_final.normal:
7484 AF-006 item_type=1 qty*_stok = 21 0 0, qty*_final = 1 0 0
7487 AF-006 item_type=2 qty*_stok = 21 0 0, qty*_final = 1 0 0
```

Math: `stock_small=104`, `qty_order=1`, available `105`. Conv `5/1` -> 21 largest, 0 middle, 0 smallest. Cocok dengan expected FE.

## Diff boundary check

- Changed: `sales/repository/order_repository.go` (stock source switched to `inv.stock` ledger), `sales/service/order_service_test.go` (2 test baru), `.opencode/evidence/20260619-1148-sx-2250-reward-stock/execution.md` (this file).
- Not touched: `.env`, `docker-compose.yml`, migrations, other services, schema files.

## No-secret note

- Login lokal menggunakan credential user test (tidak ditulis ke source/test/log).
- Tidak ada token, password, atau curl auth dari Jira yang disalin ke repo.
- Token runtime disimpan sementara di `/tmp/sx2250_token` (file temp sistem), tidak di-commit.
- Verifikasi runtime di luar repo working tree, tidak masuk git diff.

## Sanitized PM/FE note received after implementation

Catatan dari Widya / BE YOGIE (sanitized; bearer token/curl auth intentionally ignored and not copied):

- Product: AF-006 (`pro_id=10752`).
- Initial warehouse stock: `0 0 4`; initial on-customer: `1 0 3`.
- After process: warehouse stock `0 0 2`; on-customer order `2 0 0`.
- `conv_unit2=5`, `conv_unit3=1`.
- SO includes normal product AF-006 qty `0 0 1` and reward product AF-006 qty `0 0 1`, so stock movement by 2 pcs is expected and already correct.
- Detail Sales Order expected:
  - `order_detail_id=7474`: `is_product_promotion_so=false`, `is_product_promotion_final=false`, qty `0 0 1`, display stock `0 0 3` from `0 0 2 + 0 0 1`.
  - `order_detail_id=7477`: `is_product_promotion_so=true`, `is_product_promotion_final=true`, qty `0 0 1`, display stock `0 0 3` from `0 0 2 + 0 0 1`.
- Required behavior: `details.normal` and `details_final.normal` must distinguish normal row and reward row even when `pro_id` is same; no `pro_id` aggregate contamination.

## Sanitized PM/FE followup (Widya / BE YOGIE)

Feedback: stock movement sudah benar, available stock harus `0 0 1` untuk AF-006 dengan qty order `0 0 1` karena stock kosong. Validasinya dengan Inventory Stock Report `/inventory/v1/stocks/report`.

- Bearer token/curl auth dari evidence tidak dipakai/dicopy ke source/test/evidence.
- Verifikasi: `inv.stock` ledger `SUM(qty_in)-SUM(qty_out)` untuk AF-006 = `0` (sama dengan Inventory report).

## Sanitized PM/FE followup total smallest (Widya / Yogie)

Feedback: untuk `SO2606190003` AF-006, Inventory total smallest = `104`, qty order = `1`, expected available stock = `21 0 0`. Stock source Sales saat ini sudah pakai `inv.stock` ledger sehingga basis sama dengan Inventory report. Validasi:

- DB local: `inv.stock` SUM(qty_in) - SUM(qty_out) untuk AF-006 wh_id=350 = `104`.
- API lokal: `GET /v2/orders/SO2606190003` -> `details.normal` dan `details_final.normal` untuk 7484 (item_type=1, flags false) dan 7487 (item_type=2, flags true) = `21 0 0`.
- Math: `104 + 1 = 105` smallest, conv `5/1` -> 21 large, 0 medium, 0 small. Cocok.
- Test `TestDetailV2_SX2250_TotalSmallest104PlusOrderQty_ReturnsLarge21` mengunci case ini agar tidak regresi.

## Priority rule for QTY_ORDER on stock display

Catatan terbaru FE/PM (Widya / Yogie) meminta `QTY_ORDER` mengikuti priority:

1. `qty*_final` (jika sudah terisi, terdeteksi via pointer non-nil, termasuk `0,0,0` eksplisit)
2. `qty*` (sales, fallback jika `qty*_final` tidak ada)
3. `qty_po*` (purchase, fallback terakhir)

`WAREHOUSE_STOCK` memakai `inv.stock` ledger (basis Inventory Stock Report). `Available Stock = WAREHOUSE_STOCK + QTY_ORDER` lalu di-convert ke Large/Medium/Small via helper yang sama.

Implementasi: helper `stockDisplayQtyByPriority` di `sales/service/order_service.go:2476-2496` memilih qty per row dengan priority `qty*_final > qty* > qty_po*`, deteksi "sudah terisi" memakai non-nil pointer (bukan non-zero), dipakai di tiga call site `computeDisplayedAvailableStockBreakdown` (sales `2973-2982`, purchase `3082-3091`, final `3245-3254`).

Test:

- `TestDetailV2_SX2250_StockDisplayUsesFinalQtyWhenPresent`: `qty1_final=1`, `qty1=5` (sales) → `details.normal` dan `details_final.normal` mendapat `qty*_stok = 0,0,1` (priority final menang).
- `TestDetailV2_SX2250_StockDisplayFallsBackToSalesWhenFinalUnset`: `Qty1Final/Qty2Final/Qty3Final` nil, `qty1=1` → row hanya muncul di `details.normal` (final tab skip karena active qty 0) dan stock `0,0,1` (fallback ke sales).
- `TestDetailV2_SX2250_TotalSmallest104PlusOrderQty_ReturnsLarge21`: warehouse smallest `104` + qty order `1` = `105`, conv `5/1` → `qty*_stok = 21,0,0` for normal and reward rows in `details.normal` and `details_final.normal`. Mengunci feedback FE untuk `SO2606190003` (AF-006).


## Risks/blockers

- Source switched dari snapshot table ke ledger, sehingga `qty*_stok` di response sekarang bisa beda dari nilai `sls.order_detail.qty*_stok` yang sudah terlanjur tersimpan di DB (kolom itu menjadi stale untuk konsumer yang baca langsung; helper display tidak membacanya).
- `FindWarehouseStockByWhIdAndProIds` sekarang tidak membaca `qty_on_order`/`qty_on_shipping`/`qty_bs`/`qty_exp` di `inv.warehouse_stock`; jika ada business case yang butuh include on-order, butuh diskusi dengan FE apakah `qty_order` dari ledger juga harus di-include dalam display.
- Stock movement path di-trace di kode, tidak diuji runtime proses SO. Tidak ditemukan exclusion reward row.
- Cancelled-order stock date basis: current code pakai `ro.RoDate` sebagai cutoff untuk semua status termasuk CANCELLED. FE/PM belum konfirmasi apakah cancelled SO harus pakai `ro_date` cutoff atau current-date stock. Quality gate mencatat `PASS_WITH_RISKS` pada item ini; perlu konfirmasi product rule dan/atau dedicated test.

