# Discovery: sales dev → qa (endpoint GET /v2/orders/:ro_no)

**Task ID:** 20260518-0817-sales-dev-to-qa-demo  
**Tanggal:** 2026-05-18  
**Service:** `sales` (Go module di `/Users/ujang/Projects/Geekgarden/scylla-be/sales`)

## File yang Diinspeksi

| File | Keterangan |
|---|---|
| `controller/order_controller.go` | Route `GET /v2/orders/:ro_no` → `DetailV2` handler |
| `service/order_service.go` | Implementasi `DetailV2`, `Update`, `UpdateEnhance` |
| `service/order_stock_helper.go` | Helper `canonicalAPIStockBreakdown`, `computeDisplayedAvailableStockBreakdown`, `applyStockBreakdownToPointers` |
| `service/order_stock_helper_test.go` | Unit test untuk helper di atas |
| `service/order_service_test.go` | Regression test `DetailV2` |
| `repository/order_repository.go` | Query `FindByNo`, `FindByNoNoCustID`, `FindOutletByID`, dll |

## Perbedaan dev vs qa (arah: dev lebih baru dari qa)

### 1. `service/order_stock_helper.go` — **FILE BARU di dev, tidak ada di qa**
- Memperkenalkan tipe `apiStockBreakdown` dan fungsi:
  - `safeConvUnits`
  - `toTotalSmallFromAPIUnits`
  - `canonicalAPIStockBreakdown` — konversi stok total-small → L/M/S canonical
  - `computeDisplayedAvailableStockBreakdown` — gabungkan stok gudang + order qty
  - `applyStockBreakdownToPointers` — assign ke pointer field response

### 2. `service/order_service_test.go` — **FILE BARU di dev, tidak ada di qa**
- Unit test untuk semua fungsi di `order_stock_helper.go`

### 3. `service/order_service.go` — **BERBEDA (dev lebih baru)**
- `DetailV2` (sales tab & final tab): refactor kalkulasi `qty1_stok/qty2_stok/qty3_stok`
  - **qa (lama):** inline `conversion.Qty.ConvToQtyConversion()` + addisi manual jika `!useWarehouseCurrentOnly`
  - **dev (baru):** pakai `computeDisplayedAvailableStockBreakdown` + `applyStockBreakdownToPointers`
- `Update`: refactor stock breakdown pakai `canonicalAPIStockBreakdown` + `applyStockBreakdownToPointers`
- `UpdateEnhance`: sama, refactor ke `canonicalAPIStockBreakdown`
- `TestStore_DoesNotPersistStockSnapshotDuringInitialCreate` — test baru di dev
- `readCSVReferenceFixture` — path candidates diperluas di dev untuk worktree support

### 4. `service/order_service_test.go` — **BERBEDA (dev lebih baru)**
- Nilai ekspektasi `qty1_stok/qty3_stok` berubah di beberapa test karena logika konversi yang dikoreksi
- Test `TestStore_DoesNotPersistStockSnapshotDuringInitialCreate` hanya ada di dev

### 5. `repository/order_repository.go` — **BERBEDA (whitespace only)**
- Hanya perbedaan trailing space di SQL string. Tidak ada perubahan logika.

### 6. `so_service.go` — **BERBEDA**
- `hasValidDownloadPONumber`: di dev hanya cek `poNo`, di qa cek `poNo` OR `orderNo`
- `filterDownloadDataPoWithPONumber`: mengikuti perubahan di atas

## Commit Kunci di dev yang Belum Ada di qa

| Commit | Pesan |
|---|---|
| `3085b2c` | fix: normalize available stock breakdown for sales order detail |
| `055478a` | feat(order): determine sales order status from validation rules |
| `10abb3f` | fix(service): centralize float comparison tolerance |
| `f9491a6` | fix(order): batch resolve promo item unit fallback |
| `7f8f97f` | fix(service): batch load reward product masters for snapshots |
| `6320e27` | fix(order): batch load final order details |
| `3fc3a0f` | fix(order): include opr_type in order response payloads |
| `8ecce0a` | feat(sales): support edit-order enhance flow on PATCH /v1/orders/:ro_no |
| ... | (total ~60+ commit di dev belum di qa) |

## Merge Base

`548a909ca24e357aedcb0289dc86a6a664333d5f`

## Risiko

- Perubahan logika `qty_stok` di `DetailV2` mengubah nilai response endpoint `GET /v2/orders/:ro_no` — **breaking bagi consumer yang bergantung pada nilai lama**
- `so_service.go` mengubah filter PO number — perlu verifikasi apakah fitur download SO terpengaruh
- Banyak commit di dev belum di qa; cherry-pick selektif berisiko konflik; merge/rebase lebih aman
- Migration SQL baru di dev belum ada di qa — harus dijalankan setelah merge
