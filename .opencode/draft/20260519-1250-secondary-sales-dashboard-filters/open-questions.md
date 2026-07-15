# Open Questions — Task 118/119 Secondary Sales Dashboard Filters

Task id: `20260519-1250-secondary-sales-dashboard-filters`
Status: blokir finalisasi plan sampai keputusan diberikan.

## 1. Kontrak `year` (backward compatibility)

Pertanyaan: bagaimana behavior `year` pada `GET /sales/v1/reports/secondary-sales/sum-date` dan `/secondary-sales/group`?

Opsi:
- A. Optional + fallback `year = current year` bila kosong. FE Task 118/119 wajib kirim `year`. Aman untuk client lama tapi current year jadi default tersembunyi.
- B. Optional + tanpa fallback. Bila `year` kosong, query tetap tanpa `dt."year"` (perilaku lama, multi-year campur). Paling backward-compatible tapi tidak menutup gap multi-year.
- C. Required. `year` wajib; client lama tanpa `year` akan 400. Semantik paling bersih, tapi breaking.

Rekomendasi awal: A bila FE sudah siap kirim `year`; B bila ada client lama yang tidak bisa diubah cepat.

## 2. Sumber `cust_id` (security)

Pertanyaan: apakah `cust_id` query boleh meng-override `cust_id` auth pada endpoint ini?

Opsi:
- A. `cust_id` selalu dari `c.Locals("cust_id")`. Query param `cust_id` diabaikan. Aman dari cross-tenant leak. Tetap sesuai pola controller report existing.
- B. `cust_id` query mengoverride auth, tapi divalidasi memakai scope helper existing.
  - Blocker: di `sales` path tidak ditemukan scope validator BU yang reusable untuk dashboard report; perlu keputusan data model/authorization sebelum aman.
- C. `cust_id` query mengoverride tanpa validasi.
  - Tidak direkomendasikan. Risiko cross-tenant.

Rekomendasi awal: A sampai sumber scope BU tersedia.

## 3. Definisi “business unit” pada konteks ini

Pertanyaan: apa yang dimaksud “business unit” pada Task 118/119?
- Apakah `cust_id` distributor child dari `parent_cust_id` user login (principal melihat distributor di bawahnya)?
- Atau dimensi BU baru di luar `cust_id` yang belum ada di `report.fact_*` saat ini?

Bila jawabannya BU = child `cust_id` di bawah `parent_cust_id` user login, scope check bisa dibangun memakai pattern parent-child seperti di `master/repository/outlet_repository.go` (`smc.m_customer.parent_cust_id`).

Bila BU dimensi baru, ini bukan perubahan kecil dan harus naik ke `@architect` lebih dulu sebelum implementasi.

## 4. Behavior bila `cust_id` query kosong

Asumsi default: fallback ke `c.Locals("cust_id")`. Konfirmasi diperlukan agar tidak menabrak harapan FE.

## 5. Validasi `month`/`year`

Saat ini payload tidak punya `validate:"required"`. Konfirmasi: apakah `month` (1..12) dan `year` (>=2000) perlu validasi range untuk menutup input invalid?
