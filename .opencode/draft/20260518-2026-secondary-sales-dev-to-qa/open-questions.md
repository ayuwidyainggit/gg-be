# Open questions — sinkronisasi secondary-sales dev → qa

1. Branch `demo-18052026` sudah ada (HEAD `1e7586f`) beserta worktree `scylla-be-worktrees-20260518/sales`. Pilihan:
   a) Pakai branch `demo-18052026` existing, reset/perbarui dari `qa` lalu apply perubahan dev → qa di sana.
   b) Buat branch baru dengan suffix waktu, misal `demo-18052026-2026` (HHMM) untuk menghindari tabrakan.
   c) Hapus worktree+branch `demo-18052026` lama, buat ulang dari `qa`.
   Rekomendasi default: (b) `demo-18052026-2026` dengan worktree baru `scylla-be-worktrees-20260518-2026/sales`. Aman, tidak destruktif.

2. Strategi sinkronisasi yang diinginkan:
   a) Cherry-pick hanya commit yang menyentuh endpoint `secondary-sales` (13 commit di discovery), risiko konflik tinggi karena helper bersama.
   b) Merge `origin/dev` ke branch baru (atas dasar `qa`), bawa semua perubahan dev. Paling aman secara compile/test, tapi cakupan lebih luas dari sekadar endpoint ini.
   c) Squash patch khusus berisi diff `origin/qa..origin/dev` terbatas pada file: `sales/controller/report_controller.go`, `sales/service/report_service.go`, `sales/repository/report_repository.go`, `sales/entity/report.go`, `sales/model/report.go`, `sales/pkg/config/env/env.go`, `sales/pkg/constant/constant.go`, `sales/main.go`, dan dua test file terkait — sebagai satu commit “sync(secondary-sales): bring dev fixes to qa”.
   Rekomendasi default: (c) untuk bounded blast radius; jatuh ke (b) kalau test gagal karena ketergantungan helper non-endpoint.

3. Apakah perlu push branch ke remote (`origin/demo-18052026-2026`) dan buka MR ke `qa`, atau cukup branch lokal dulu untuk demo?

4. Env QA: apakah `OBS_HUAWEI_AK/SK/ENDPOINT/BUCKET` sudah ter-set di environment QA target? Tanpa itu service gagal boot karena `env.ValidateRequired` baru dari commit `54a63b8`.
