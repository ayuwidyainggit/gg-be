# Open Questions — SX-2054 Answer Frequency

Task ID: `20260525-1738-sx-2054-answer-frequency`

## Pertanyaan material

1. Apakah new write harus **menolak** legacy `Multiple` mulai SX-2054?
   - Rekomendasi: ya, create/update hanya menerima tiga value baru: `One Time`, `Multiple Times, One Day`, `Multiple Times, Different Day`.
   - Alasan: requirement menyebut value setelah perubahan, dan legacy `Multiple` hanya untuk data lama.

2. Apakah SX-2054 harus mengubah behavior mobile submit untuk membedakan `Multiple Times, One Day` vs `Multiple Times, Different Day`?
   - Rekomendasi planner: jangan ubah behavior submit tanpa aturan bisnis eksplisit.
   - Fakta: `mobile/service/survey.go` sekarang selalu blok submit kedua dengan `CheckExistingSubmission(...)`, sehingga semua survey efektif one-submit.
   - Risiko: kalau product menganggap dua value baru sudah harus mengubah duplicate prevention, requirement masih kurang.

3. Apakah DB perlu langsung diberi `CHECK` constraint?
   - Rekomendasi planner: tahap pertama widen-only `VARCHAR(50)` agar new value masuk dan legacy aman; tambah transitional `CHECK` hanya setelah distinct-value audit.
   - Alasan: local migration tidak punya constraint; data aktual mungkin punya nilai liar.

## Asumsi untuk lanjut implementasi bila tidak dijawab

- New write menolak `Multiple`.
- Read path tetap raw dan bisa menampilkan legacy `Multiple`.
- Schema migration pertama hanya widen column ke `VARCHAR(50)`.
- Tidak ada mass migration dari `Multiple` ke value baru.
- Submit/mobile rule hanya diaudit dan didokumentasikan, bukan diubah, kecuali product memberi rule eksplisit.
