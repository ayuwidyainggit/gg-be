# Open Questions SX-2085

Task ID: `20260528-1440-sx-2085-monitoring-collection-paid`

## Tidak blocking untuk rencana awal

1. FE butuh top-level total seperti `collection.total_paid`, atau cukup sum dari existing `collection[].collection_total`?
   - Default plan: pertahankan contract existing `collection` array, isi `collection_total` per outlet.
   - Jika FE minta total eksplisit, tambah field non-breaking terpisah perlu koordinasi.

2. Staging access tersedia untuk SQL manual dan before/after API capture?
   - Default plan: implementer menjalankan saat eksekusi; blocker hanya untuk QA evidence, bukan untuk code path lokal.

3. Jika satu deposit punya invoice lintas outlet, apakah payment dibagi per outlet atau seluruh payment dihitung ke tiap outlet?
   - Default plan: validasi grain. Bila ada lintas outlet, stop untuk keputusan business karena SQL reference belum menjelaskan split allocation.
