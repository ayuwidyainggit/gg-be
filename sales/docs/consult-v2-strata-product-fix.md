# ConsultV2 Strata Product Reward – Analisis & Plan

## Request (client_test.http)
- **Promo:** strata06 (Strata Syarat Value & Reward Lebih 1 Jenis Product Not Sequential - Quantity)
- **Order:** total_gross_value = 20.400.000, 3 produk eligible (710, 534, 673)
- **Strata:** 5 tingkat, strata_reward = [1, 2, 3, 4, 5] (quantity "largest"), total 15 unit

## Perbedaan Actual vs Expected

### Reward product 1 (pro_id 534)
| Field   | Actual                         | Expected                       |
|---------|--------------------------------|--------------------------------|
| qty3    | 2                              | 2                              |
| gross_value | 135000                     | 135000                         |
| promo1  | 135000                         | 135000                         |
| promo2  | 270000                         | **135000**                     |
| promo3  | 405000                         | **0**                          |
| promo4  | 540000                         | **0**                          |
| promo5  | 675000                         | **0**                          |

### Reward product 2 (pro_id 484)
| Field   | Actual                         | Expected                       |
|---------|--------------------------------|--------------------------------|
| qty3    | 13                             | 13                             |
| gross_value | 1440000                    | 1440000                        |
| promo1  | 1440000                        | **0**                          |
| promo2  | 2880000                        | **1440000**                    |
| promo3  | 4320000                        | 4320000                        |
| promo4  | 5760000                        | 5760000                        |
| promo5  | 7200000                        | 7200000                        |

## Root cause

**Actual:** Setiap reward product dihitung sebagai `promo_i = gross_value * strata_reward[i]`, seolah seluruh qty produk itu memenuhi semua strata dengan penuh (1, 2, 3, 4, 5 unit per strata).

**Expected:** Unit reward dialokasikan **secara berurutan per strata** (strata 1 dulu, lalu 2, 3, 4, 5) dan **antar produk berurutan** (produk pertama dipakai dulu, habis lalu produk berikutnya).

- Strata butuh: 1 + 2 + 3 + 4 + 5 = 15 unit (largest).
- Pro 534 stock 2: 1 unit → strata 1 (promo1), 1 unit → strata 2 (promo2) → sisa demand strata 2 = 1, strata 3,4,5 masih penuh.
- Pro 484 stock 13: 0 → strata 1, 1 → strata 2, 3 → strata 3, 4 → strata 4, 5 → strata 5.
- Jadi per produk: `promo_i = (jumlah unit produk ini yang dialokasikan ke strata i) × unit_price`.

## Plan perbaikan

1. **Alokasi per strata (sequential):**
   - Demand per strata: `remainingDemand[i] = strataListNonSeq[i].RewardValue` (1, 2, 3, 4, 5).
   - Untuk setiap reward product (urutan list):
     - Stock = stok produk dalam reward UOM.
     - Untuk strata 1..5: `take = min(stock, remainingDemand[i])`, lalu:
       - `allocation[stratum] = take`
       - `remainingDemand[i] -= take`, `stock -= take`.

2. **Response per reward product:**
   - `promo_i = allocation[i] × unit_price` (unit_price = `gross_value` dari ConversionWithPrice untuk 1 unit = SellPrice1/2/3).
   - Tetap pakai `ConversionWithPrice` dengan qty total yang diambil dari produk ini; `gross_value` di response tetap unit price (SellPrice sudah per unit).

3. **Lokasi kode:** Di blok `rewardType == model.RewardTypeProduct`, cabang `useStrataNonSequentialProduct`, saat "Sufficient stock":
   - Sebelum loop `for _, reward := range rewards`, inisialisasi `remainingDemand[5]` dari `strataListNonSeq`.
   - Di dalam loop, setelah hitung `qtyReward` dan sebelum `ConversionWithPrice`, hitung `allocation[5]` untuk produk ini dan update `remainingDemand`.
   - Saat set `rewardProduct.Promo1..Promo5`, pakai `allocation[i] * unitPrice` (bukan `GrossValue * strata_reward[i]`).

## Field response lain

- `strata_rule_uom`, `strata_per_scope`: di expected bisa kosong `""`; bisa ditambahkan di response strata non-sequential jika ada di entity.
- `strata_id` / `strata_desc`: bisa beda ID (env/data) asal struktur sama; tidak wajib diubah untuk logic ini.
