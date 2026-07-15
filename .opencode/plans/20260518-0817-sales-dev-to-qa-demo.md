# Plan: bandingkan dan ambil perubahan endpoint `GET /v2/orders/:ro_no` dari `dev` ke `qa`

## Goal

Memastikan apakah code path untuk endpoint `https://be.scyllax.online/sales/v2/orders/SO2605170001` berbeda antara branch `dev` dan `qa`, lalu menyiapkan rencana implementasi lengkap untuk mengambil perubahan relevan dari `dev` ke branch `qa` memakai worktree baru bernama `demo-18052026`.

## Non-goals

- Tidak melakukan implementasi, merge, cherry-pick, atau push branch.
- Tidak menjalankan migration produksi.
- Tidak memverifikasi response endpoint live terhadap environment remote.
- Tidak merencanakan sinkronisasi seluruh `dev` ke `qa` tanpa batas scope.

## Scope

Scope utama:
- `sales/controller/order_controller.go`
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`
- `sales/service/order_stock_helper.go`
- `sales/service/order_stock_helper_test.go`
- `sales/repository/order_repository.go`

Scope pendukung yang perlu diputuskan sebelum eksekusi:
- Apakah ikut bawa perubahan `service/so_service.go`
- Apakah perlu bawa commit pendukung lain yang menjadi dependency perilaku `DetailV2`
- Apakah strategi terbaik `merge dev -> qa`, `cherry-pick minimal`, atau `checkout file` selektif

## Requirements

- Bandingkan branch `dev` dan `qa` pada code path endpoint `GET /v2/orders/:ro_no`.
- Identifikasi perubahan logika vs perubahan kosmetik.
- Siapkan worktree baru dari `qa` bernama `demo-18052026`.
- Rencana harus aman untuk service `sales` dan mengikuti validasi repo lokal:
  - `rtk go mod download && rtk go mod tidy`
  - `rtk go test ./...`
  - targeted test untuk `DetailV2`
- Rencana harus mempertimbangkan migration yang tertinggal di `qa`.

## Acceptance Criteria

- Terdapat kesimpulan jelas apakah endpoint target berbeda antara `dev` dan `qa`.
- Terdapat daftar file yang wajib dibawa untuk perubahan endpoint.
- Terdapat keputusan strategi promosi code yang paling aman.
- Terdapat langkah kerja eksekusi yang atomic dan siap dijalankan dari worktree `demo-18052026`.
- Terdapat daftar validasi test dan risk check pasca-promosi.

## Existing Patterns/Reuse

Reuse dulu, jangan bikin ulang.

Pattern yang ditemukan:
- Route endpoint sudah ada di `controller/order_controller.go:56-58` dan memanggil `OrderService.DetailV2`.
- Logika utama response endpoint ada di `service/order_service.go:2742+`.
- Dev sudah punya helper reusable `service/order_stock_helper.go` untuk normalisasi stok canonical L/M/S.
- Dev sudah punya regression test reusable di `service/order_stock_helper_test.go` dan `service/order_service_test.go`.
- `repository/order_repository.go` untuk endpoint ini tidak punya perubahan logika; beda hanya whitespace.

Kesimpulan reuse:
- Untuk endpoint target, reuse terbaik ialah mengambil helper dan test yang sudah ada di `dev`, bukan rewrite lokal di `qa`.
- Tidak ada utility repo lain yang lebih cocok daripada helper `order_stock_helper.go`.

## Constraints

- Repo root bukan git repo; git root relevan ada di `/Users/ujang/Projects/Geekgarden/scylla-be/sales`.
- User minta worktree baru bernama `demo-18052026`.
- Branch `demo-18052026` lokal/remote sudah tampak ada varian `demo-13052026`, `demo-05052026`; nama `demo-18052026` harus dicek tabrakan dulu saat eksekusi.
- Banyak commit di `dev` belum ada di `qa`; ambil semua tanpa scope akan berisiko tinggi.
- Planner ini hanya menulis artefak. Tidak implementasi.

## Risks

1. Perubahan nilai `qty1_stok/qty2_stok/qty3_stok`
   - `dev` menormalkan stok ke representasi canonical L/M/S.
   - `qa` masih pakai kalkulasi inline yang menghasilkan mapping berbeda.
   - Efek: payload endpoint target berubah walau data order sama.

2. Dependency commit tersembunyi
   - Commit `3085b2c` fokus stok detail.
   - Tetapi branch `dev` juga punya banyak commit lain di `DetailV2` path, mis. promo snapshot, unit fallback, VAT, filter zero-qty, status, `opr_type`.
   - Jika keluhan user sebenarnya bukan stok, cherry-pick satu commit bisa kurang.

3. Migration drift
   - `dev` punya banyak file migration yang belum ada di `qa`.
   - Jika perubahan yang diambil menyentuh field baru, query/runtime di `qa` bisa gagal tanpa migration.

4. Scope creep dari `so_service.go`
   - Ada perubahan terpisah di `service/so_service.go`, tidak langsung di endpoint target.
   - Jangan ikut terbawa tanpa kebutuhan bisnis jelas.

## Decisions/Assumptions

Keputusan sementara:
- Endpoint target memang **berbeda** antara `dev` dan `qa`.
- Perbedaan paling langsung dan terisolasi untuk endpoint target ada pada normalisasi `qty*_stok` di `DetailV2`.
- `repository/order_repository.go` tidak perlu dibawa untuk kasus ini karena beda hanya formatting.
- Worktree baru harus dibuat dari branch `qa`, bukan dari `dev`.

Assumptions:
- Tujuan user: ambil perubahan perilaku endpoint `GET /v2/orders/:ro_no` dari `dev` ke `qa`, bukan sinkronisasi seluruh fitur `sales`.
- Ticket ini fokus pada mismatch response endpoint target, terutama field stok yang ditampilkan.
- Jika saat implementasi ditemukan mismatch juga pada promo snapshot/VAT/order status, scope harus dinaikkan ke mini-bundle commit terkait `DetailV2`.

Open question material:
- Apakah user ingin ambil **minimal fix** untuk endpoint ini saja, atau **semua perubahan `DetailV2` yang sudah ada di dev**?

## TDD/Test Plan

TDD required: **Ya**.

Alasan:
- Task menyentuh production logic endpoint API.
- Perubahan output field stok mudah regress.
- Dev sudah menyediakan pola regression test yang bisa direuse.

Existing test patterns:
- `service/order_service_test.go` punya test `DetailV2`.
- `service/order_stock_helper_test.go` di dev mengisolasi helper conversion logic.

First failing/regression test:
- Jalankan test `DetailV2` yang memverifikasi:
  - cancelled order memakai warehouse current only
  - non-cancelled order memakai canonical L/M/S stock breakdown
  - final detail filter zero effective qty tetap benar
- Jika helper belum ada di `qa`, test helper dari dev akan fail dulu. Itu Red pertama yang diharapkan.

Green step:
- Bawa helper `order_stock_helper.go` + integrasi `order_service.go` + test terkait dari `dev`.

Refactor step:
- Pastikan tidak ada duplicate inline conversion tersisa di `DetailV2`, `Update`, `UpdateEnhance`.
- Jangan sentuh file di luar kebutuhan minimal tanpa bukti.

Edge cases:
- `convUnit2 <= 0`
- `convUnit3 <= 0`
- stok nol
- cancelled order
- non-cancelled order dengan existing qty pada order
- final details dengan deleted rows / effective qty nol

Commands:
- `rtk go mod download && rtk go mod tidy`
- `rtk go test ./service/...`
- `rtk go test ./...`
- bila perlu targeted: `rtk go test ./service -run 'TestDetailV2|TestCanonicalAPIStockBreakdown|TestComputeDisplayedAvailableStockBreakdown|TestStore_DoesNotPersistStockSnapshotDuringInitialCreate'`

## Implementation Steps

1. Siapkan worktree baru dari `qa` dengan branch kerja `demo-18052026`.
2. Verifikasi nama worktree/branch belum bentrok lokal/remote.
3. Di worktree baru, bandingkan `qa` terhadap commit kandidat `3085b2c` dan file terkait.
4. Jalankan Red state di `qa` dengan test `DetailV2` yang relevan.
5. Ambil perubahan minimal dari `dev` untuk endpoint target:
   - `service/order_stock_helper.go`
   - `service/order_stock_helper_test.go`
   - hunk relevan di `service/order_service.go`
   - hunk/test relevan di `service/order_service_test.go`
6. Jangan bawa `repository/order_repository.go` kecuali ada conflict resolution yang membutuhkan.
7. Jalankan targeted tests.
8. Jalankan `rtk go test ./...` di module `sales`.
9. Jika gagal karena dependency commit lain, evaluasi perluasan scope ke bundle `DetailV2` terkait, urutan prioritas:
   - `5534509` cancelled stock behavior
   - `0252e97` no-cust zero-qty filter
   - `25c2419` assertion tightening
   - commit `DetailV2` support lain yang terbukti dibutuhkan
10. Setelah hijau, lakukan smoke diff akhir terhadap `qa` untuk memastikan hanya file minimal yang berubah.
11. Siapkan PR/merge request dari `demo-18052026` ke `qa` setelah review.

## Expected Files to Change

Minimal target:
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`
- `sales/service/order_stock_helper.go`
- `sales/service/order_stock_helper_test.go`

Biasanya tidak perlu berubah untuk kasus ini:
- `sales/controller/order_controller.go`
- `sales/repository/order_repository.go`

Conditional only if evidence demands:
- file migration terkait `sls.order` / `sls.order_detail`
- file pendukung `DetailV2` lain bila tests/runtime membuktikan dependency

## Agent/Tool Routing

- `@artifact-planner` — plan ini, source of truth artifact
- `@explorer` — kalau saat eksekusi perlu cari dependency commit tambahan atau file impact lebih luas
- `@fixer` — implementasi bounded di worktree `demo-18052026`
- `@quality-gate` — final signoff karena menyentuh response API production
- `@oracle` — opsional bila saat eksekusi muncul tradeoff antara cherry-pick minimal vs merge subset yang lebih besar

Tool decisions:
- Local project discovery: **dipakai**
- Official docs/context7: **tidak perlu**, karena issue murni codebase lokal Go/business logic
- GitHub/API: **tidak perlu**, repo lokal cukup
- Brave/web search: **tidak perlu**
- Browser evidence: **tidak perlu** untuk task backend ini

## Execution-ready Worklist / Handoff Contract

`start_with: T01`

| Task ID | Action | depends_on | owner/lane | Validation | Exit criteria | status | blocker reason | requires_user_decision |
|---|---|---|---|---|---|---|---|---|
| T01 | Verifikasi `sales` git root, branch `qa`, dan ketersediaan nama worktree `demo-18052026` | none | `@fixer` | `git rev-parse --show-toplevel`, `git branch --all --list '*demo-18052026*'`, `git worktree list` | Jelas branch/worktree target aman dipakai | ready | none | no |
| T02 | Buat worktree baru dari `qa` untuk branch `demo-18052026` | T01 | `@fixer` | `git worktree add ...` lalu `git status` | Worktree baru aktif dan bersih | ready | none | no |
| T03 | Jalankan baseline test Red pada `qa` worktree untuk area `DetailV2` | T02 | `@fixer` | `rtk go test ./service -run 'TestDetailV2|TestCanonicalAPIStockBreakdown|TestComputeDisplayedAvailableStockBreakdown|TestStore_DoesNotPersistStockSnapshotDuringInitialCreate'` | Bukti awal state `qa` dan failure/reachability test didapat | ready | none | no |
| T04 | Ambil perubahan minimal dari `dev` untuk helper stock breakdown dan integrasi `order_service` | T03 | `@fixer` | `git diff --stat qa..HEAD`, compile/test targeted | Empat file minimal ter-update tanpa file liar | ready | none | no |
| T05 | Jalankan targeted tests untuk helper dan `DetailV2` | T04 | `@fixer` | command test targeted di atas | Semua test target hijau | ready | none | no |
| T06 | Jalankan full module validation `sales` | T05 | `@fixer` | `rtk go mod download && rtk go mod tidy && rtk go test ./...` | Module `sales` hijau | ready | none | no |
| T07 | Review diff akhir dan cek dependency tambahan/migration drift | T06 | `@quality-gate` | `git diff --name-only qa...HEAD`, review hasil test | Dipastikan scope aman, tidak ada migration wajib yang tertinggal untuk change ini | ready | none | no |
| T08 | Jika T07 menemukan dependency tambahan, perluas scope secara minimal dan ulang T05-T07 | T07 | `@fixer` | targeted diff + test | Scope tambahan terdokumentasi dan hijau | blocked | Hanya jalan bila T07 menemukan blocker | no |
| T09 | Siapkan branch untuk PR ke `qa` | T07 | `@fixer` | `git status`, ringkasan diff | Branch siap direview | ready | none | no |
| T10 | Final signoff untuk merge ke `qa` | T09 | `@quality-gate` | review evidence + test output | Rekomendasi merge/no-merge jelas | ready | none | no |

## Validation Commands

Jalankan dari worktree module `sales`:

```bash
git rev-parse --show-toplevel
git status
git worktree list
git branch --all --list '*demo-18052026*'
rtk go mod download && rtk go mod tidy
rtk go test ./service -run 'TestDetailV2|TestCanonicalAPIStockBreakdown|TestComputeDisplayedAvailableStockBreakdown|TestStore_DoesNotPersistStockSnapshotDuringInitialCreate'
rtk go test ./...
git diff --name-only qa...HEAD
git diff --stat qa...HEAD
```

Worktree creation candidate:

```bash
git worktree add ../sales-demo-18052026 -b demo-18052026 qa
```

Jika branch `demo-18052026` sudah ada:

```bash
git worktree add ../sales-demo-18052026 demo-18052026
```

## Evidence Requirements

Wajib simpan bukti berikut saat implementasi:
- hasil `git worktree list`
- hasil `git diff --name-only qa...HEAD`
- output targeted test `DetailV2`
- output `rtk go test ./...`
- jika scope meluas, alasan dan commit/file tambahan yang ikut dibawa
- catatan apakah migration diperlukan atau tidak untuk scope final

## Done Criteria

Selesai jika semua kondisi ini benar:
- worktree `demo-18052026` dibuat dari `qa`
- perubahan endpoint target dari `dev` sudah diambil secara minimal/terkontrol
- test target dan full module `sales` hijau
- diff akhir bersih dan dapat dijelaskan
- `@quality-gate` menyatakan aman untuk merge ke `qa`

## Final Planning Summary

Artifacts dibuat:
- Primary plan: `.opencode/plans/20260518-0817-sales-dev-to-qa-demo.md`
- Evidence kept: `.opencode/evidence/20260518-0817-sales-dev-to-qa-demo/discovery.md`

Key decisions:
- Ada perubahan nyata antara `dev` dan `qa` untuk endpoint target.
- Perubahan paling relevan dan terisolasi untuk endpoint target berpusat di commit `3085b2c`.
- `repository/order_repository.go` tidak perlu dibawa untuk kasus ini.
- Strategi awal paling aman: worktree dari `qa` + ambil minimal fix + test ketat.

Assumptions:
- User ingin promosi scoped fix, bukan full sync `dev` ke `qa`.

Remaining open questions:
- Ambil minimal fix endpoint saja, atau bundle semua perubahan `DetailV2` yang sudah ada di `dev`?

Readiness:
- Siap dieksekusi oleh `@orchestrator`/`@fixer` tanpa replanning, untuk jalur minimal-fix.

Cleanup performed:
- Tidak ada draft artifact dibuat. Evidence discovery disimpan karena masih berguna untuk eksekusi.
