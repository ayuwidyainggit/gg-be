# Discovery SX-2253 Available Stock Same `pro_id` Per Row

Task id: `20260619-0831-sx-2253-available-stock-row`

## File diperiksa

- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `sales/service/order_service.go`
- `sales/service/order_stock_helper.go`
- `sales/service/order_stock_helper_test.go`
- `sales/service/order_service_test.go`
- `sales/repository/order_repository.go`
- `sales/entity/order_detail.go`
- `sales/model/order_detail.go`

## Pola proyek

- Service Sales adalah module Go sendiri di `sales/`.
- Layer wajib: Controller → Service → Repository → DB.
- Validasi harus jalan dari `sales/`.
- Command repo-local memakai prefix `rtk`.
- `DetailV2` ada di `sales/service/order_service.go`.
- Test existing untuk `DetailV2` ada di `sales/service/order_service_test.go`.
- Stock display helper ada di `sales/service/order_stock_helper.go`.

## Temuan teknis

- Response detail memakai `qty1_stok`, `qty2_stok`, `qty3_stok` sebagai available/displayed stock breakdown.
- `DetailV2` mengambil warehouse stock via `FindWarehouseStockByWhIdAndProIds(custID, whID, proIds)`.
- Repository method `FindWarehouseStockByWhIdAndProIds` return `map[pro_id] -> warehouse qty`; ini stock lookup product-level, bukan order qty aggregation.
- `DetailV2` menghitung display stock tiga kali:
  - sales detail lines: `sales/service/order_service.go:2955-2964`
  - purchase detail lines: `sales/service/order_service.go:3063-3072`
  - final detail lines: `sales/service/order_service.go:3225-3234`
- Hitungan sekarang memakai `computeDisplayedAvailableStockBreakdown(whStockQty, current row qty..., !useWarehouseCurrentOnly, convUnit2, convUnit3)`.
- Existing test `TestDetailV2_SameSKURewardDoesNotContaminateNormalRow` menjaga row normal/reward same SKU tidak tercampur promo flag, tetapi belum assert `qty*_stok` row-level.
- Existing tests `TestDetailV2_Cancelled_UsesWarehouseCurrentOnlyForDisplayedStock` dan `TestDetailV2_NonCancelled_KeepsExistingDisplayedStockBehavior` menjaga behavior stock display cancel/non-cancel.

## T1 trace `qty*_stok` write path

- `sales/service/order_service.go:5555-5565` — update partial detail snapshot; service writes `updates["qty1_stok"]`, `updates["qty2_stok"]`, `updates["qty3_stok"]` from `canonicalAPIStockBreakdown(currentStock, conv2, conv3)`.
- `sales/repository/order_repository.go:75` — repository interface declaration for `RefreshOrderDetailStock`.
- `sales/repository/order_repository.go:1043-1050` — repository method implementation; direct DB write of passed `qty*_stok` values.
- `RefreshOrderDetailStock` has no call site in `sales/` from `rg -n 'RefreshOrderDetailStock' sales/`.
- Conclusion: no evidence to modify `sales/repository/order_repository.go`. Current bug more likely in `DetailV2` response mapping or upstream row data, not this unused repository writer path.

## T2-T4 tests: row-level stock already correct

- `TestDetailV2_SameProIDNormalAndRewardRow_ComputesStockIndependently` (added at `sales/service/order_service_test.go:1326-1511`): `pro_id=8435`, wh=100, normal row qty1=10, reward row qty1=5, conv2=10, conv3=5.
  - Actual: normal `Qty1Stok=2, Qty2Stok=1, Qty3Stok=0`; reward `Qty1Stok=2, Qty2Stok=0, Qty3Stok=5`. Rows independent.
- `TestDetailV2_SameProIDTwoNormalRows_ComputesStockIndependently` (added at `sales/service/order_service_test.go:1513-1606`): two `item_type=1` rows same `pro_id=8435`, qty1=10 and qty1=5, wh=100.
  - Actual: row A `Qty1Stok=2, Qty2Stok=1, Qty3Stok=0`; row B `Qty1Stok=2, Qty2Stok=0, Qty3Stok=5`. Rows independent.
- `TestDetailV2_DifferentProID_StockByProduct` (added at `sales/service/order_service_test.go:1608-1668`): row A `pro_id=8435` wh=100 qty1=10, row B `pro_id=8436` wh=50 qty1=5.
  - Actual: row A `Qty1Stok=2, Qty2Stok=1, Qty3Stok=0`; row B `Qty1Stok=1, Qty2Stok=0, Qty3Stok=5`. Stock by product.
- All three tests PASS in current code. `DetailV2` mapping in `sales/service/order_service.go:2898-3072` already uses per-row `computeDisplayedAvailableStockBreakdown` with `whStockQty` from `warehouseStockMap[pro_id]` and the current row's own `Qty1/Qty2/Qty3`. No aggregation by `pro_id` in this path. The plan's `computeDisplayedAvailableStockBreakdown` invariant is intact.

## T5 — no code change required

- `DetailV2` per-row loop at `sales/service/order_service.go:2898-3072` already row-keyed. No fix needed in `sales/service/order_service.go`.
- `sales/repository/order_repository.go` not modified. `RefreshOrderDetailStock` is unused inside `sales/`.

## T6 — full service suite

- `rtk go test ./service -count=1` → 195 passed, 0 failed. No regression.

## T7 — diff boundary check

- Not a git repo (`git diff --stat` warns not a git repo). Equivalent: only files changed are
  - `sales/service/order_service_test.go` (T2/T3/T4 tests added; L1326-1668)
  - `.opencode/evidence/20260619-0831-sx-2253-available-stock-row/discovery.md` (T1/T2-T4 evidence)
- No production code changed. Diff boundary respected.

## T8 — final summary

- Status: PASS_WITH_RISKS. New tests pass; existing `TestDetailV2_*` and `TestComputeDisplayedAvailableStockBreakdown*` all pass; 195 service tests pass; no production code modified.
- Risk: plan's bug premise (per-row `qty*_stok` uses aggregated qty by `pro_id` in `DetailV2` mapping) is not reproducible in the current `DetailV2` code path. If the bug exists in another path not covered (e.g., an older mobile sync, an external system, or pre-`DetailV2` data with stale `qty*_stok` columns), this worklist does not address it. Executor recommendation: route to `@oracle` for review of plan premise, or re-open with a failing production repro beyond `DetailV2` mapping before further code edits.
- Open question: where is the user-reported bug actually observable? Need a failing production repro or specific call path that exhibits aggregated `qty*_stok` per row.

## Reuse candidates

- Pakai `computeDisplayedAvailableStockBreakdown` untuk invariant row-level.
- Tambah test di `order_service_test.go`, dekat test same SKU existing atau stock display tests.
- Reuse mock `mockOrderRepositoryDetailV2`.

## Risiko

- Jika bug root berada di persisted `qty*_stok` snapshot bukan `DetailV2`, executor perlu trace `RefreshOrderDetailStock` call sites sebelum fix final.
- Jangan ubah `FindWarehouseStockByWhIdAndProIds` menjadi row-keyed; warehouse stock memang product-level.
- Jangan ubah promo aggregation helpers kecuali bukti menunjukkan field display stock memakai aggregate promo qty.
- Jangan salin kredensial Jira/staging ke repo, log, test, evidence.

## Source strategy

- Local repo evidence dipakai.
- Jira URL tidak di-fetch karena akses auth dan berisi kredensial test; prompt user cukup sebagai source-approved issue summary.
- External docs tidak diperlukan; fix Go service/test lokal.
- Browser/staging validation direncanakan sebagai optional manual evidence oleh executor dengan env/token lokal, bukan artifact commit.
