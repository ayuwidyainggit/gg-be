# Open Questions — SX-2034

Task ID: `20260521-1515-sx2034-extra-call-monitoring`

## Q1. Backfill data lama
Banyak baris `pjp_principles.destinations_history` dengan `is_extra_call = true` AND `destination_id IS NULL` (akibat bug `StoreFromListPrinciple`). Apakah dieksekusi backfill (lihat lampiran SQL di plan utama) untuk staging dan production setelah verifikasi?

- Default rekomendasi: ya, staging dulu, lalu production setelah `@quality-gate` review hasilnya.
- Dampak skip: data extra call lama tetap tidak muncul di Monitoring Activity meski fix insert + fix query sudah live.
- Risiko: salah join key (`destination_code` vs field lain) bisa mengubah row ke ID yang salah → mitigasi: dry-run SELECT + transaction.

## Q2. Distributor extra call
Endpoint `POST /scylla-mobile/api/v1/m-outlets/from-list` saat ini hanya membawa `outlet_id`; tidak ada cabang `distributor_id`. Apakah scope SX-2034 termasuk extra call type distributor, atau cukup outlet (sesuai sample Jira)?

- Default rekomendasi: outlet only sekarang. Buat follow-up issue untuk distributor extra call (perlu desain payload + UI).
- Catatan: query monitoring tetap inklusif via `dh.destination_type` sehingga ketika distributor extra call live nanti, FE tidak perlu adjust kontrak besar.

## Q3. Konfirmasi join key backfill
Asumsi: `destinations_history.destination_code` = `m_outlet.outlet_code` (untuk outlet) dan `m_distributor.distributor_code` (untuk distributor). Mohon konfirmasi sebelum eksekusi staging.

- Jika ternyata bisa multi-match (mis. `outlet_code` tidak unik antar `cust_id`), tambahkan filter `cust_id` di UPDATE.
