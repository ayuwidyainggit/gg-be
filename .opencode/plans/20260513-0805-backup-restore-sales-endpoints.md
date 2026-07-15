# Goal

Menyiapkan rencana implementasi untuk membackup logic endpoint sales yang terkait dengan `PATCH /v1/orders/enhance/:ro_no` dan `POST /v1/validate-order`, lalu merestore logic dari branch `dev` ke branch worktree `demo-13052026` yang berbasis `qa` pada repo `sales`.

# Non-goals

- Tidak mengubah source code implementation pada tahap planner ini.
- Tidak melakukan restore untuk service selain `sales`.
- Tidak memverifikasi perilaku production endpoint live di luar fakta repo lokal, karena user sudah mengklarifikasi method target.
- Tidak membuat commit, merge request, atau push remote pada tahap planner ini.

# Scope

- Backup source code terkait endpoint:
  - `PATCH /v1/orders/enhance/:ro_no`
  - `POST /v1/validate-order`
- Identifikasi file controller, service, entity, dan file pendukung yang terlibat.
- Pembuatan branch/worktree baru `demo-13052026` dari branch `qa` pada repo `sales`.
- Restore dengan cara menyalin logic/file relevan dari branch `dev` ke branch demo berbasis `qa`.
- Validasi minimal bahwa code hasil restore dapat diinspeksi dan, bila memungkinkan, lulus test/kompilasi yang relevan di module `sales`.

# Requirements

1. Semua command shell dijalankan dengan prefix `rtk` sesuai kebijakan lokal repo.
2. Operasi git harus dijalankan dari directory `sales/`, bukan root repo, karena root bukan git repo.
3. Backup harus mencakup seluruh source yang mempengaruhi dua endpoint target, bukan hanya satu file route.
4. Branch target harus bernama `demo-13052026` dan dibuat dari `qa`.
5. Restore harus memakai source dari branch `dev`.
6. Implementasi restore harus meminimalkan perubahan di luar area endpoint terkait.
7. Hasil akhir harus menyisakan working tree target yang jelas untuk direview sebelum commit/push.

# Acceptance Criteria

1. Tersedia salinan backup yang bisa dipakai untuk rollback atau audit terhadap logic endpoint terkait dari branch `dev`.
2. Branch/worktree `demo-13052026` berhasil dibuat dari `qa` tanpa mengganggu worktree lain.
3. Logic endpoint `PATCH /v1/orders/enhance/:ro_no` pada branch target telah dipulihkan dari `dev`.
4. Logic endpoint `POST /v1/validate-order` pada branch target telah dipulihkan dari `dev`.
5. File-file yang berubah di branch target terbatas pada area endpoint terkait atau dependensi langsungnya.
6. Ada hasil verifikasi minimal berupa diff inspection dan command validasi yang menunjukkan restore tidak rusak secara struktural.

# Existing Patterns/Reuse

- Reuse utama adalah implementation existing pada branch `dev`; ini harus menjadi sumber kebenaran restore, bukan penulisan ulang logic.
- Route existing yang sudah teridentifikasi:
  - `sales/controller/order_controller.go` → `roRouteV1.Patch("/enhance/:ro_no", controller.UpdateEnhance)`
  - `sales/controller/validate_order_controller.go` → `invoiceRouteV1.Post("/", controller.ValidateOrder)`
- Service existing yang sudah teridentifikasi:
  - `sales/service/order_service.go` → `UpdateEnhance(...)`
  - `sales/service/validate_order_service.go` → `ValidateOrder(...)`
- Entity request yang langsung terkait enhance:
  - `sales/entity/edit_order_enhance.go`
- Test pattern reusable untuk controller/helper ada di `sales/controller/so_controller_test.go`; bila perlu, pola test ringan di controller package dapat diikuti untuk coverage tambahan.

# Constraints

- Planner hanya boleh menulis artefak `.opencode/`.
- Root `scylla-be` bukan git repo; repo git yang relevan untuk task ini adalah `sales/`.
- Ada worktree existing lain di mesin lokal, sehingga implementasi harus menghindari path/branch collision.
- Repo mengandung secrets di beberapa file menurut AGENTS lokal; implementasi tidak boleh menyalin/mengekspos secrets.
- Tidak boleh mengandalkan asumsi bahwa dua endpoint hanya menyentuh controller; service/entity/repository terkait harus ikut ditinjau saat restore.

# Risks

- Restore berbasis copy file penuh dari `dev` dapat membawa perubahan tidak terkait bila file yang sama mengandung fitur lain.
- Restore berbasis hunk/diff parsial lebih aman terhadap scope, tetapi lebih rawan miss dependency antar file.
- Branch `qa` bisa tertinggal jauh dari `dev`, sehingga restore dapat memicu konflik API struct atau method signature.
- Branch/worktree `demo-13052026` mungkin sudah ada saat implementasi dijalankan dan perlu penanganan non-destruktif.

# Decisions/Assumptions

- Diputuskan berdasarkan jawaban user bahwa endpoint enhance yang diproses adalah **PATCH saja**, bukan POST.
- Diputuskan bahwa restore berarti **copy file/logic dari `dev` ke branch demo berbasis `qa`**, bukan cherry-pick commit sebagai pendekatan utama.
- Diputuskan nama target branch/worktree adalah `demo-13052026`.
- Asumsi kerja: file paling relevan akan mencakup controller/service/entity yang sudah teridentifikasi, namun implementer tetap harus memverifikasi diff dependensi tambahan sebelum finalisasi.
- Tidak ada open question yang tersisa untuk memulai implementasi terarah.

# TDD/Test Plan

## Apakah TDD wajib?

Tidak wajib secara penuh untuk task ini, karena inti pekerjaan adalah **restore/mem-porting logic existing** antar branch, bukan mendesain behavior baru. Namun verification-first tetap wajib.

## Alasan exemption

- Source of truth behavior sudah ada di branch `dev`.
- Tujuan utama adalah parity logic antara `dev` dan branch demo berbasis `qa`.

## Existing test patterns

- `sales/controller/so_controller_test.go` menunjukkan pattern unit test ringan pada package `controller` menggunakan `fiber.New()` dan `httptest`.
- Perlu cek apakah sudah ada test lain yang langsung menyentuh order/validate-order sebelum menambah test baru.

## First failing/regression test

- Jika ditemukan gap test yang murah ditambahkan, prioritas pertama adalah regression test ringan untuk route/controller helper atau request parsing yang terkait endpoint restore.
- Bila belum feasible, ganti Red step dengan **diff-based regression baseline**:
  1. ambil diff implementation endpoint target antara `qa` dan `dev`,
  2. restore ke branch demo,
  3. pastikan diff branch demo terhadap `dev` untuk area endpoint menjadi nol atau sesuai scope yang disepakati.

## Green step

- Restore logic/file relevan dari `dev` ke branch `demo-13052026`.
- Jalankan verifikasi minimal:
  - `rtk git diff -- sales/controller/order_controller.go sales/controller/validate_order_controller.go sales/service/order_service.go sales/service/validate_order_service.go sales/entity/edit_order_enhance.go`
  - `rtk go test ./...` di `sales/` bila waktu/dependency memungkinkan.

## Refactor step

- Refactor hanya bila diperlukan untuk menyelesaikan incompatibility kecil antara `qa` dan logic dari `dev` tanpa mengubah behavior endpoint target.
- Hindari refactor non-esensial.

## Edge cases

- Branch target atau worktree target sudah ada.
- File target di `qa` punya signature berbeda sehingga copy langsung gagal compile.
- Ada dependency logic endpoint di file repository/model lain yang awalnya tidak terdeteksi.
- Test full module terlalu berat atau gagal karena faktor eksternal non-task; jika begitu, fallback ke targeted build/test dan diff inspection.

## Commands

- `rtk git status`
- `rtk git branch --all`
- `rtk git worktree list`
- `rtk git diff qa..dev -- <relevant-files>`
- `rtk go test ./...`

# Implementation Steps

1. Dari directory `sales/`, verifikasi working tree bersih dan cek ulang branch/worktree existing.
2. Buat backup artefak source dari branch `dev` untuk file/logic endpoint target. Bentuk backup yang direkomendasikan:
   - patch file: `rtk git diff qa..dev -- <relevant-files> > ...` pada lokasi kerja aman di luar source target, atau
   - salinan file referensi dari `dev` untuk area endpoint.
3. Bandingkan `dev` vs `qa` untuk file kandidat berikut:
   - `controller/order_controller.go`
   - `controller/validate_order_controller.go`
   - `service/order_service.go`
   - `service/validate_order_service.go`
   - `entity/edit_order_enhance.go`
   - tambah file lain bila diff menunjukkan dependency langsung.
4. Buat worktree baru dari `qa` untuk branch `demo-13052026` pada path yang belum dipakai, misalnya sibling folder terdedikasi.
5. Di worktree target, restore logic dari `dev` dengan pendekatan **scope-minimized copy**:
   - utamakan copy hunk atau blok logic endpoint terkait,
   - gunakan copy file penuh hanya bila file memang spesifik ke endpoint atau perubahan lain di file itu memang harus ikut agar compile/behavior konsisten.
6. Jalankan inspection diff pada worktree target untuk memastikan perubahan terbatas pada endpoint terkait.
7. Jalankan validasi teknis bertahap:
   - targeted diff check,
   - targeted test bila ada package yang relevan,
   - `rtk go test ./...` di `sales/` bila realistis.
8. Dokumentasikan file yang berubah dan hasil validasi untuk handoff implementasi/review.

# Expected Files to Change

Perkiraan minimal file yang mungkin berubah saat implementasi:

- `sales/controller/order_controller.go`
- `sales/controller/validate_order_controller.go`
- `sales/service/order_service.go`
- `sales/service/validate_order_service.go`
- `sales/entity/edit_order_enhance.go`

Kemungkinan tambahan bila ada dependency langsung dari branch `dev`:

- repository/model/helper yang dipanggil khusus oleh logic enhance/validate-order
- test file baru/terkait yang dibutuhkan untuk regression verification

# Agent/Tool Routing

- Planner/artifact writer: sudah selesai di artefak ini.
- Implementasi: agent implementer/fixer atau operator manual yang boleh edit source pada repo `sales`.
- Tool utama implementasi:
  - shell/git via `rtk`
  - file diff/inspection lokal
- Tool tambahan hanya bila diperlukan:
  - explorer untuk menelusuri dependency tambahan endpoint
  - oracle bila muncul konflik arsitektur atau pilihan restore yang berisiko

# Validation Commands

Jalankan dari `sales/` kecuali disebut lain.

1. Status dan worktree
   - `rtk git status`
   - `rtk git branch --all`
   - `rtk git worktree list`
2. Baseline diff source vs target
   - `rtk git diff qa..dev -- controller/order_controller.go controller/validate_order_controller.go service/order_service.go service/validate_order_service.go entity/edit_order_enhance.go`
3. Setelah restore pada worktree target
   - `rtk git status`
   - `rtk git diff -- controller/order_controller.go controller/validate_order_controller.go service/order_service.go service/validate_order_service.go entity/edit_order_enhance.go`
4. Validasi Go
   - `rtk go test ./...`

# Evidence Requirements

- Simpan bukti branch asal (`dev`) dan target (`qa` → `demo-13052026`).
- Simpan daftar file final yang berubah.
- Simpan diff sebelum dan sesudah restore untuk area endpoint.
- Simpan output singkat hasil validasi (`rtk git status`, `rtk git worktree list`, dan test bila dijalankan).
- Jika ada file tambahan di luar kandidat awal, catat alasan dependensinya.

# Done Criteria

- Backup logic endpoint dari `dev` tersedia dan dapat dirujuk.
- Worktree `demo-13052026` berhasil dibuat dari `qa`.
- Logic dua endpoint target telah dipulihkan dari `dev` ke branch demo.
- Scope perubahan sudah diinspeksi dan dinilai terbatas/relevan.
- Validasi minimal telah dijalankan dan hasilnya terdokumentasi.
- Perubahan siap diteruskan ke tahap implementasi/review/commit oleh agent atau operator yang berwenang.

# Final Planning Summary

- Pertanyaan material sudah diajukan dan dijawab user:
  - endpoint enhance diperlakukan sebagai `PATCH` saja,
  - restore berarti copy logic/file dari `dev` ke branch demo berbasis `qa`,
  - nama target adalah `demo-13052026`.
- Artefak yang dibuat:
  - source of truth plan: `.opencode/plans/20260513-0805-backup-restore-sales-endpoints.md`
  - evidence discovery: `.opencode/evidence/20260513-0805-backup-restore-sales-endpoints/discovery.md`
- Keputusan utama:
  - operasikan git dari `sales/`, bukan root,
  - gunakan implementation di `dev` sebagai sumber restore,
  - minimalkan scope dengan restore hanya pada logic endpoint terkait dan dependency langsungnya.
- Asumsi utama:
  - file kandidat yang sudah diidentifikasi mencakup mayoritas area perubahan, walau implementer tetap harus memverifikasi dependency tambahan.
- Open questions: tidak ada blocker tersisa.
- Readiness for implementation: **siap** untuk tahap implementasi bounded pada repo `sales`.
- Cleanup performed:
  - tidak ada draft file yang dibuat, sehingga tidak ada cleanup draft.
  - evidence discovery dipertahankan karena masih berguna sebagai acuan implementasi dan audit scope.
