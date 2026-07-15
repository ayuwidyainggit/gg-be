# Open Questions — SX-2038

Task ID: `20260522-1101-sx-2038-monitoring-detail-null`

Material yang masih perlu konfirmasi user/QA sebelum implementasi finalize:

- **Q1 — Mapping `emp_id`**: di curl QA `emp_id=484`, di komentar dev.fe `salesman_id=482`. Apakah test user `princessa@gmail.com` benar-benar terhubung ke `m_salesman.emp_id = 482`, atau `484`? Default asumsi: query pakai parameter apapun yang dikirim FE (tanpa lookup tambahan, karena `pjp.salesman_id == m_salesman.emp_id` di skema saat ini); QA harus reproduce dengan ID yang benar.
- **Q2 — Status PJP yang valid**: query detail principal saat ini hanya `approval_status = 'Approved'`. Jalur distributor & query referensi dev.fe pakai `IN ('Approved','Need Review')`. Konfirmasi ke product/dev.fe: status mana yang harus ditampilkan di Live Monitoring Detail Principal?
- **Q3 — Sumber data plan principal**: apakah Live Monitoring Detail Principal harus menghitung extra-call (sumber `pjp_principles.destinations_history` dengan `is_extra_call=true`) sebagai bagian dari `planned/visited/skipped/extra_call`, atau cukup destinations reguler? Jalur list principal sudah memisahkan dua loader (SX-2034); detail diharapkan paritas.
- **Q4 — Akses DB staging untuk validasi**: planner tidak punya akses langsung ke staging DB. Siapa yang menjalankan query debug Step 1-6 (Yoggie / QA / DBA)? Tanpa hasil itu, Hipotesis 1-3 tidak bisa difinalisasi sebelum implementasi; default rekomendasi: `@fixer` jalankan saat implementasi sambil menambah unit test untuk semua skenario.
- **Q5 — Definisi sukses Acceptance Criteria**: apakah AC SX-2038 mengharuskan `data` jadi array of object (dengan `visit_information`, `sales`, `return`, dst) atau cukup `data` non-null? Default asumsi mengikuti kontrak Postman `Scylla-Live-Monitoring-Complete` (array of object).

Default operasional kalau user tidak menjawab dalam waktu dekat:
- Hipotesis utama yang di-fix duluan: Q2 (status filter) + Q3 (extra-call source) — risiko rendah, paritas dengan jalur lain.
- Q1 ditangani via verifikasi QA di staging (bukan code change).
