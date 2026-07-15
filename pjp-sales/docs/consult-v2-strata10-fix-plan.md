# Consult V2 – strata10 (Strata Quantity, Fixed Value Per Order) – Analisis & Plan

## 1. Request (client_test.http)

```json
POST /v2/promotions/consult
{
  "order_date": "2025-12-17",
  "outlet_id": 1404,
  "salesman_id": 359,
  "wh_id": 63,
  "details": [
    { "pro_id": 710, "qty1": 10, "qty2": 0, "qty3": 0, "gross_value": 6000000 },
    { "pro_id": 534, "qty1": 5,  "qty2": 0, "qty3": 0, "gross_value": 5000000 },
    { "pro_id": 673, "qty1": 5,  "qty2": 0, "qty3": 0, "gross_value": 5400000 }
  ]
}
```

- Total order gross (dari request): 6M + 5M + 5.4M = **16.400.000**
- Promo: **strata10** – Strata Syarat Quantity & Reward Value PerOrder Not Sequential (Quantity), `strata_reward` = 2000, `fixed_value`, `per_order`

---

## 2. Perbandingan Actual vs Expected

| Aspek | Actual | Expected |
|-------|--------|----------|
| **products_eligible** | `[710]` (1 produk) | `[710, 534, 673]` (3 produk) |
| **total_gross_value** | 16.400.000 | 4.300.000 |
| **reward_value** | Satu entry (710): gross 6.000.000, promo1 2000, net 5.998.000 | Tiga entry: masing-masing promo1 = 667, net = gross − 667 |
| **reward_value[].gross_value** | 6.000.000 (satu produk) | 1.000.000, 1.500.000, 1.800.000 |
| **strata_rule_uom** | "smallest" | (tidak ada di expected; opsional) |

Expected menyebut: *"strata_reward(2000) / product_eligible(3) → jumlah product_eligible"* → **promo1 per produk = 2000/3 ≈ 667**.

---

## 3. Root cause

### 3.1 Hanya satu produk eligible (710)

- **Phase 5 & 6** di `ConsultV2`: `validatedPromoProductGroups[promoID]` diisi hanya dari **product criteria**.
- Jika strata10 di DB hanya punya **satu** product criteria (misalnya pro_id 710, mandatory), maka hanya 710 yang masuk ke `validatedPromoProductGroups["strata10"]`.
- Produk 534 dan 673 tidak ada di product criteria → tidak pernah ditambahkan → hanya 710 yang dapat reward (2000 penuh).

### 3.2 total_gross_value

- **Saat ini** (baris ~2590): `response.TotalGrossValue = orderTotalGross` → selalu **total seluruh order** (16.400.000).
- **Expected**: 4.300.000 = 1.000.000 + 1.500.000 + 1.800.000 = **jumlah gross_value dari reward_value** (hanya produk eligible untuk promo ini).

Jadi expected menganggap `total_gross_value` = **total gross untuk produk eligible promo ini**, bukan total order.

### 3.3 Pembagian reward per_order fixed_value

- Untuk **per_order** + **fixed_value**, kode saat ini (sekitar 2867–2887) sudah benar:  
  `promo1 = Round(strataReward / productCount)`.  
  Jika `productCount == 3`, maka promo1 = 667 per produk.  
- Masalah bukan di rumus, tapi di **productCount**: karena hanya 710 yang eligible, productCount = 1 → promo1 = 2000 untuk satu produk.

### 3.4 reward_value[].gross_value (1M, 1.5M, 1.8M)

- Expected memakai gross per produk **bukan** dari request (6M, 5M, 5.4M), tapi 1M, 1.5M, 1.8M.
- Kemungkinan: **gross_value di response** = nilai yang dihitung dari master (misalnya unit price × qty), atau dari sumber lain. Perlu konfirmasi requirement: apakah response harus pakai `detail.GrossValue` dari request atau “gross dari sistem” (misalnya harga × qty).

---

## 4. Plan perbaikan ConsultV2

### 4.1 Memperbaiki products_eligible (semua 3 produk ikut)

**Opsi A – Data (disarankan jika aturan promo memang “semua produk order”):**

- Pastikan promo **strata10**:
  - **Tanpa** product criteria (semua produk order eligible), atau
  - Product criteria memuat **710, 534, 673** (mandatory atau optional sesuai aturan bisnis).

Dengan ini, Phase 5 & 6 akan mengisi `validatedPromoProductGroups["strata10"]` dengan ketiga produk, dan tidak perlu ubah alur strata.

**Opsi B – Code (jika aturan bisnis: strata quantity = semua detail order ikut reward):**

- Untuk promo **tipe strata** dengan **quantity rule**:
  - Setelah Phase 7b (validasi strata), jika strata match:
    - Hitung `ruleValue` dari **semua** `req.Details` (bukan hanya dari `validatedPromoProductGroups[promoID]`).
    - Jika `ruleValue` masih masuk range strata yang sama:
      - Set `validatedPromoProductGroups[promoID]` = semua `req.Details` untuk promo ini, dan
      - Update `subTotalValidatedPromoProductGroups[promoID]` = total gross dari semua detail.
- Risiko: mengabaikan product criteria untuk strata; hanya lakukan jika requirement jelas “strata quantity = seluruh order”.

Rekomendasi: cek dulu data strata10. Jika product criteria sengaja hanya 710, pakai Opsi A (tambah 534 & 673 di criteria). Jika memang harus “semua produk order”, pakai Opsi B.

### 4.2 total_gross_value sesuai expected

- **Saat ini:** `response.TotalGrossValue = orderTotalGross` (total order).
- **Expected:** `total_gross_value` = 4.300.000 = jumlah gross produk **eligible untuk promo ini**.

**Perubahan kode:**

- Untuk response **strata** (misalnya saat `useStrataNonSequentialFixedValue` atau path strata lainnya), set:
  - `response.TotalGrossValue = totalGrossValueForPromo`  
  dengan `totalGrossValueForPromo` = jumlah `GrossValue` dari `validatedPromoProductGroups[promoID]` (sudah dihitung di baris ~2727–2730).
- Tetap gunakan `orderTotalGross` untuk response yang memang dimaksud “total order” (jika ada path lain yang butuh itu).

Ini membuat `total_gross_value` = total gross **eligible** untuk promo tersebut, konsisten dengan expected 4.300.000 (setelah eligible = 3 produk dan jika gross_value per produk pakai nilai yang sama dengan expected).

### 4.3 Reward per produk (promo1 = 667)

- Tidak perlu ubah rumus: `promo1 = Round(strata_reward / productCount)` sudah benar.
- Cukup pastikan **productCount = 3** (dari 4.1), maka promo1 = 667 per produk.

### 4.4 reward_value[].gross_value (1M, 1.5M, 1.8M)

- Saat ini kode pakai `rv.GrossValue = float64(detail.GrossValue)` dari request (6M, 5M, 5.4M).
- Expected pakai 1M, 1.5M, 1.8M. Itu bisa berarti:
  - **Sumber lain:** misalnya unit price dari master × qty (per UOM), atau dari service price/ConversionWithPrice.
- **Langkah:**
  1. Konfirmasi ke product/business: apakah `reward_value[].gross_value` harus sama dengan request, atau harus “gross dari sistem” (harga × qty)?
  2. Jika harus “gross dari sistem”: tambah step (di ConsultV2 atau helper) untuk hitung gross per produk dari harga × qty (misalnya pakai existing conversion/pricing service), lalu pakai nilai itu untuk `rv.GrossValue` dan untuk `total_gross_value` (sum of eligible).

---

## 5. Ringkasan aksi

| # | Aksi | Lokasi / catatan |
|---|------|-------------------|
| 1 | Pastikan strata10 punya product criteria kosong atau berisi 710, 534, 673 | Data / DB |
| 2 | (Opsional) Untuk strata quantity, expand eligible = semua req.Details jika strata match | Phase 7b, `promotion_service.go` |
| 3 | Untuk response strata: set `response.TotalGrossValue = totalGrossValueForPromo` | Sekitar baris 2589–2590, bedakan path strata vs slab/order total |
| 4 | Konfirmasi sumber `reward_value[].gross_value`; jika dari sistem, tambah hitungan gross dari harga × qty | Phase 8, blok `useStrataNonSequentialFixedValue` |

Dengan 1 (atau 2) + 3, hasil consult akan punya `products_eligible = [710, 534, 673]`, `total_gross_value = 4.300.000` (jika gross per produk sesuai expected), dan `promo1 = 667` per produk. Poin 4 hanya perlu jika requirement memang gross_value response bukan dari request.
