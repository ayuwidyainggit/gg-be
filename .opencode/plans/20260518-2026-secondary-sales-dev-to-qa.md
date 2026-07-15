# Plan — sinkronisasi `secondary-sales` dari `dev` ke `qa`

Task ID: `20260518-2026-secondary-sales-dev-to-qa`
Tanggal: 2026-05-18 20:26 Asia/Jakarta
Service target: `sales`
Endpoint target: `POST /v1/reports/secondary-sales`
Branch sumber: `origin/dev` (`b7eafd9`)
Branch target: `origin/qa` (`2ec0152`)
Primary source of truth: `.opencode/plans/20260518-2026-secondary-sales-dev-to-qa.md`
Evidence: `.opencode/evidence/20260518-2026-secondary-sales-dev-to-qa/discovery.md`
Open questions: `.opencode/draft/20260518-2026-secondary-sales-dev-to-qa/open-questions.md`

## Goal

Buat rencana lengkap untuk membawa perubahan kode endpoint `POST /v1/reports/secondary-sales` dari branch `dev` ke branch `qa` memakai worktree baru bernama sesuai pola `demo-DDMMYYYY`, tanpa mengubah source code saat fase planning.

## Non-goals

- Tidak melakukan merge/cherry-pick/source edit dalam fase plan ini.
- Tidak push branch remote.
- Tidak buka MR/PR.
- Tidak mengubah env, secrets, package, lockfile, test, atau kode produksi.
- Tidak mengubah branch/worktree existing yang sudah ada tanpa keputusan user.

## Scope

Cakupan kode endpoint berdasarkan route dan dependency langsung:

- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/entity/report.go`
- `sales/model/report.go`
- `sales/pkg/config/env/env.go`
- `sales/pkg/constant/constant.go`
- `sales/main.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`

Endpoint terkait yang terpengaruh karena helper/dashboard report sama:

- `GET /v1/reports/secondary-sales/sum-date`
- `GET /v1/reports/secondary-sales/group`
- `GET /v1/reports/secondary-sales/trend-sales`
- `POST /v1/extract/secondary-sales`

## Requirements

- Worktree baru harus berbasis `qa`, lalu menerima perubahan dari `dev` untuk endpoint `secondary-sales`.
- Layering Controller → Service → Repository → DB harus tetap utuh.
- Tenant filter `cust_id` dan `parent_cust_id` harus dipertahankan.
- Perintah validasi service `sales` memakai prefix `rtk` sesuai repo policy.
- Tidak boleh commit secrets atau `.env`.
- Jika nama worktree/branch `demo-18052026` sudah ada, implementor harus minta keputusan user sebelum tindakan destruktif.
- Sinkronisasi harus membawa test yang relevan dari `dev` agar behavior fixes punya bukti.

## Acceptance Criteria

- Worktree baru ada dan berada di branch demo yang disetujui user, berbasis `origin/qa`.
- Diff final terhadap `origin/qa` hanya memuat perubahan yang diperlukan untuk `secondary-sales` dan dependency compile/test langsung, kecuali user memilih merge penuh `dev`.
- Endpoint `POST /v1/reports/secondary-sales` memakai behavior `dev` untuk:
  - publish export message gagal → report ditandai `FILE_STATUS_FAILED`.
  - export completion kembali benar.
  - filename/object name dibangun dari `report_id` saat perlu.
  - sequence report dihitung dari prefix `SecondarySales-<date>-%`.
  - invoice canvas ikut masuk report.
  - `gross_sales` aman dari NULL via `COALESCE`.
  - product/supplier fallback memakai data customer/distributor sesuai helper `dev`.
- `rtk go test ./...` di directory `sales` lulus.
- Targeted tests untuk `secondary-sales` lulus:
  - `rtk go test ./repository -run 'TestBuildSecondarySalesUnionQuery'`
  - `rtk go test ./service -run 'TestPublishSecondarySalesReportMarksFailedWhenRabbitMQPublishFails|TestSubscribeSecondarySalesReportUploadsWorkbookWithExpectedValues'`
- Compile tidak rusak karena perubahan helper bersama.
- QA env punya `OBS_HUAWEI_AK`, `OBS_HUAWEI_SK`, `OBS_HUAWEI_ENDPOINT`, `OBS_HUAWEI_BUCKET`, atau validasi startup disesuaikan lewat keputusan rilis.

## Existing Patterns/Reuse

- Reuse pola route di `ReportController.Route`: `app.Group('/v1/reports', middleware.JWTProtected())` dan handler controller existing.
- Reuse service pattern: controller hanya parse request dan panggil `ReportService`; business logic tetap di service.
- Reuse repository SQL builder existing `buildSecondarySalesUnionQuery`, `secondarySalesProductSelect`, dan `secondarySalesProductJoins` dari `dev`, bukan tulis ulang.
- Reuse tests dari `dev` untuk regression coverage:
  - `TestBuildSecondarySalesUnionQueryCoalescesGrossSalesOperands`
  - `TestPublishSecondarySalesReportMarksFailedWhenRabbitMQPublishFails`
  - test workbook existing di `service/report_service_test.go`
- Reuse repo validation policy dari `.opencode/docs/QUALITY.md`.
- Tidak ditemukan utilitas KiloCode/project lain yang menggantikan logic endpoint; perubahan terbaik adalah reuse/port kode `dev` yang sudah ada.

## Constraints

- Repo root `/Users/ujang/Projects/Geekgarden/scylla-be` bukan git repo; `sales/` adalah git repo/module sendiri.
- Branch aktif saat discovery: `sales` di `dev`.
- `origin/dev` dan `origin/qa` berhasil di-fetch.
- Branch/worktree `demo-18052026` sudah ada:
  - branch `demo-18052026`
  - worktree `/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260518/sales`
- Git tidak boleh checkout branch yang sama ke dua worktree sekaligus.
- Menghapus/reset worktree existing bersifat destruktif; butuh keputusan user.
- Perubahan `main.go` menambah validasi env OBS; risiko boot QA bila env belum lengkap.
- Beberapa commit `dev` juga menyentuh activity-report helper; cherry-pick terlalu parsial berisiko compile error.

## Risks

- Konflik di `repository/report_repository.go` tinggi karena banyak commit endpoint dan activity-report mengubah helper/SQL yang berdekatan.
- Konflik di `service/report_service.go` sedang karena helper publish report dan export flow dipakai bersama.
- `qa` punya commit sendiri `45beac3 chore(qa): restore sales service from dev backup and add GET alias for print proforma invoice`; merge harus menjaga perubahan QA-only ini.
- `env.ValidateRequired` untuk OBS dapat membuat QA gagal startup kalau konfigurasi belum siap.
- Cherry-pick commit lama satu per satu dapat membawa perubahan luas di luar endpoint karena commit historis bercampur.
- Squash patch terbatas lebih terkendali, tapi perlu tes compile penuh untuk memastikan dependency tak tertinggal.

## Decisions/Assumptions

Keputusan planner:

- Rekomendasi strategi: squash patch terbatas dari `origin/dev` ke branch demo berbasis `origin/qa`, mencakup file endpoint dan tests terkait.
- Rekomendasi fallback: merge penuh `origin/dev` hanya jika patch terbatas gagal compile/test karena dependency tersembunyi yang terlalu luas.
- Rekomendasi branch/worktree karena konflik nama: jangan overwrite `demo-18052026`; buat nama aman seperti `demo-18052026-2026` kecuali user minta pakai existing.

Asumsi:

- User ingin tanggal hari ini `18-05-2026`, sehingga pola `demo-DDMMYYYY` berarti `demo-18052026`.
- User ingin plan dulu, bukan implementasi langsung.
- Branch remote `origin/dev` dan `origin/qa` adalah sumber kebenaran untuk compare.

Open questions material:

1. `demo-18052026` sudah ada. Pilih pakai existing, buat suffix baru, atau hapus/recreate?
2. Pilih strategi bounded squash patch, cherry-pick commit, atau merge penuh `dev`?
3. Perlu push remote + MR ke `qa`, atau lokal saja?
4. Env QA sudah punya `OBS_HUAWEI_*`?

Detail open questions ada di `.opencode/draft/20260518-2026-secondary-sales-dev-to-qa/open-questions.md`.

## TDD/Test Plan

TDD required: ya.

Reason:

- Ini perubahan behavior endpoint/report export, SQL report, file export, RabbitMQ failure handling, dan env startup; regression risk tinggi.

Existing test patterns:

- Unit test SQL builder di `sales/repository/report_repository_test.go`.
- Unit test service/report export di `sales/service/report_service_test.go` dengan mock repository/adapter dan `excelize`.
- Mock config baru dari `dev`: `mockConfigEnv` untuk `REPORT_DELAY_SECONDS`.

Red step:

- Di branch demo berbasis `qa`, jalankan targeted tests yang ada di `dev` setelah patch test diterapkan tapi sebelum production patch bila memungkinkan.
- Test yang harus gagal di QA sebelum perubahan behavior:
  - `TestBuildSecondarySalesUnionQueryCoalescesGrossSalesOperands`: SQL QA belum memaksa `COALESCE` untuk semua operand `gross_sales`.
  - `TestPublishSecondarySalesReportMarksFailedWhenRabbitMQPublishFails`: QA belum menandai report failed saat RabbitMQ publish error.

Green step:

- Terapkan patch bounded dari `dev` untuk file scope.
- Jalankan targeted tests hingga lulus.
- Jalankan `rtk go test ./...` dari directory `sales`.

Refactor step:

- Rapikan konflik tanpa mengubah kontrak endpoint.
- Pastikan helper bersama tidak duplikatif dan tidak memindahkan business logic ke controller/repository.
- Pastikan no broad unrelated changes jika strategi bounded patch dipilih.

Edge cases:

- `qty*_final` atau `sell_price*` NULL tidak membuat `gross_sales` NULL.
- Return rows dengan `qty*` atau `sell_price*` NULL tetap hitung 0-safe.
- Product/customer supplier fallback tetap memakai tenant/scope benar.
- Report publish gagal langsung update status failed.
- OBS env missing memberi error startup jelas atau dikonfirmasi siap di QA.
- Canvas invoice rows ikut union query dan tidak merusak pagination/count.

Commands:

```bash
rtk go test ./repository -run 'TestBuildSecondarySalesUnionQuery'
rtk go test ./service -run 'TestPublishSecondarySalesReportMarksFailedWhenRabbitMQPublishFails|TestSubscribeSecondarySalesReportUploadsWorkbookWithExpectedValues'
rtk go test ./...
```

## Implementation Steps

1. Minta keputusan user untuk konflik `demo-18052026` dan strategi sync.
2. Fetch branch terbaru:

```bash
git fetch origin dev qa
```

3. Buat worktree baru berbasis `origin/qa` dengan nama yang disetujui.

Jika user setuju nama aman rekomendasi:

```bash
git worktree add -b demo-18052026-2026 ../scylla-be-worktrees-20260518-2026/sales origin/qa
```

Jika user minta exact `demo-18052026`, gunakan worktree/branch existing, jangan buat baru:

```bash
git worktree list
git status --short --branch
```

Lalu minta instruksi reset/reuse sebelum lanjut.

4. Di worktree demo, terapkan patch bounded dari `origin/dev` untuk file scope.

Opsi bounded patch:

```bash
git checkout origin/dev -- controller/report_controller.go service/report_service.go repository/report_repository.go entity/report.go model/report.go pkg/config/env/env.go pkg/constant/constant.go main.go repository/report_repository_test.go service/report_service_test.go
```

5. Review diff terhadap `origin/qa`:

```bash
git diff --stat origin/qa...HEAD
git diff --name-status origin/qa...HEAD
git diff origin/qa...HEAD -- controller/report_controller.go service/report_service.go repository/report_repository.go entity/report.go model/report.go pkg/config/env/env.go pkg/constant/constant.go main.go repository/report_repository_test.go service/report_service_test.go
```

6. Kalau patch bounded compile gagal karena dependency non-endpoint, identifikasi missing symbols dan ambil dependency minimum dari `origin/dev`, bukan langsung semua file.

7. Jalankan targeted tests.

8. Jalankan full `sales` tests.

9. Jika test lulus, commit di branch demo dengan pesan:

```text
sync(secondary-sales): bring dev report export fixes to qa
```

10. Minta keputusan user sebelum push/MR.

## Expected Files to Change

Target bounded files:

- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/entity/report.go`
- `sales/model/report.go`
- `sales/pkg/config/env/env.go`
- `sales/pkg/constant/constant.go`
- `sales/main.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`

No expected changes:

- `go.mod`
- `go.sum`
- migrations
- `.env`
- docker/compose files
- docs outside `.opencode/`

## Agent/Tool Routing

- `@artifact-planner`: plan artifact owner, selesai di file ini.
- `@orchestrator`: eksekusi handoff setelah user menjawab open questions.
- `@explorer`: bantu discovery ulang jika konflik muncul atau missing dependency tidak jelas.
- `@fixer`: bounded implementation di worktree demo, resolve conflicts, run tests.
- `@oracle`: review jika merge strategy berubah ke full merge `dev` atau konflik menyentuh architecture/tenant SQL.
- `@quality-gate`: final signoff setelah tests dan diff review.
- `@librarian`: tidak diperlukan; tidak ada library/API eksternal baru.
- Browser/Playwright: tidak diperlukan; ini backend endpoint/report export, bukan UI flow.

## Execution-ready Worklist / Handoff Contract

`start_with`: `W1`

### W1 — putuskan konflik nama worktree

- action: Minta user memilih cara menangani branch/worktree `demo-18052026` yang sudah ada.
- depends_on: `none`
- owner/lane: `@orchestrator`
- validation: `git worktree list` dan `git branch --list demo-18052026 demo-18052026-*`
- exit criteria: User memilih salah satu: reuse existing, create suffix new, atau delete/recreate.
- blocking status: `blocked`
- blocker reason: Branch/worktree `demo-18052026` sudah ada; tindakan destruktif tidak boleh diasumsikan.
- requires_user_decision: `yes`

### W2 — putuskan strategi sync

- action: Pilih `bounded squash patch` (recommended), `cherry-pick commits`, atau `merge full dev`.
- depends_on: `W1`
- owner/lane: `@orchestrator`
- validation: keputusan user dicatat di handoff/commit notes.
- exit criteria: Strategi sync final jelas.
- blocking status: `blocked`
- blocker reason: Strategi memengaruhi blast radius dan risiko konflik.
- requires_user_decision: `yes`

### W3 — buat worktree demo dari `origin/qa`

- action: Jalankan `git fetch origin dev qa`, lalu `git worktree add -b <approved-branch> <approved-path>/sales origin/qa`.
- depends_on: `W1`, `W2`
- owner/lane: `@fixer`
- validation: `git -C <worktree>/sales status --short --branch` menunjukkan branch demo clean berbasis `qa`.
- exit criteria: Worktree baru/approved siap dan clean.
- blocking status: `ready` setelah W1/W2 selesai.
- blocker reason: `none` setelah keputusan user.
- requires_user_decision: `no`

### W4 — terapkan patch dari `dev`

- action: Ambil file bounded dari `origin/dev` ke worktree demo atau jalankan strategi final yang dipilih user.
- depends_on: `W3`
- owner/lane: `@fixer`
- validation: `git diff --name-status origin/qa...HEAD` hanya menunjukkan file sesuai scope atau scope yang disetujui.
- exit criteria: Patch applied, no unresolved conflict, diff scope jelas.
- blocking status: `ready` setelah W3.
- blocker reason: `none` setelah worktree clean.
- requires_user_decision: `no`

### W5 — resolve dependency/compile gap minimum

- action: Jika compile/test gagal karena dependency dari `dev` belum terbawa, ambil dependency minimum dan catat alasannya.
- depends_on: `W4`
- owner/lane: `@fixer` dengan `@explorer` bila perlu.
- validation: `rtk go test ./repository -run 'TestBuildSecondarySalesUnionQuery'` dan `rtk go test ./service -run 'TestPublishSecondarySalesReportMarksFailedWhenRabbitMQPublishFails|TestSubscribeSecondarySalesReportUploadsWorkbookWithExpectedValues'`.
- exit criteria: Targeted tests compile dan lulus.
- blocking status: `ready` setelah W4.
- blocker reason: `none`; escalate ke `@oracle` kalau dependency meluas ke tenant SQL architecture.
- requires_user_decision: `no`

### W6 — validasi penuh service `sales`

- action: Jalankan full test suite `sales`.
- depends_on: `W5`
- owner/lane: `@fixer`
- validation: `rtk go test ./...`
- exit criteria: Semua test lulus atau failure non-related didokumentasikan dengan bukti.
- blocking status: `ready` setelah W5.
- blocker reason: `none`
- requires_user_decision: `no`

### W7 — review env OBS QA

- action: Verifikasi kesiapan env `OBS_HUAWEI_AK`, `OBS_HUAWEI_SK`, `OBS_HUAWEI_ENDPOINT`, `OBS_HUAWEI_BUCKET` untuk QA.
- depends_on: `W4`
- owner/lane: `@orchestrator` + user/ops
- validation: Konfirmasi env QA atau config deployment menyatakan keys ada.
- exit criteria: Risiko startup env jelas: ready atau ditunda dengan mitigation.
- blocking status: `blocked`
- blocker reason: Planner tidak punya akses ke env QA; perubahan `main.go` dapat memblokir boot.
- requires_user_decision: `yes`

### W8 — commit lokal

- action: Commit perubahan setelah W6 lulus dan W7 punya keputusan.
- depends_on: `W6`, `W7`
- owner/lane: `@fixer`
- validation: `git status --short --branch` clean setelah commit; `git log -1 --oneline` menampilkan commit sync.
- exit criteria: Commit lokal siap review.
- blocking status: `ready` setelah W6/W7.
- blocker reason: `none` setelah validations.
- requires_user_decision: `no`

### W9 — quality gate

- action: Review final diff, scope, tests, tenant/query risk, env risk.
- depends_on: `W8`
- owner/lane: `@quality-gate`
- validation: evidence test output + diff summary + env decision tersedia.
- exit criteria: Signoff pass atau daftar blocker final.
- blocking status: `ready` setelah W8.
- blocker reason: `none`
- requires_user_decision: `no`

### W10 — push/MR ops opsional

- action: Push branch demo dan buka MR ke `qa` hanya jika user minta.
- depends_on: `W9`
- owner/lane: `@orchestrator`
- validation: remote branch/MR URL tersedia.
- exit criteria: MR siap review, atau lokal-only selesai.
- blocking status: `blocked`
- blocker reason: User belum minta push/MR.
- requires_user_decision: `yes`

## Validation Commands

Dari repo root parent, cek runtime compose bila perlu:

```bash
rtk docker compose -f docker-compose.yml ps
```

Dari directory worktree service `sales`:

```bash
rtk go mod download
rtk go test ./repository -run 'TestBuildSecondarySalesUnionQuery'
rtk go test ./service -run 'TestPublishSecondarySalesReportMarksFailedWhenRabbitMQPublishFails|TestSubscribeSecondarySalesReportUploadsWorkbookWithExpectedValues'
rtk go test ./...
```

Diff/scope checks:

```bash
git status --short --branch
git diff --stat origin/qa...HEAD
git diff --name-status origin/qa...HEAD
git log --oneline origin/qa..HEAD
```

Optional smoke if service can run in QA-like env:

```bash
rtk docker compose -f docker-compose.yml up -d sales
```

Then call endpoint with valid JWT and payload from `sales/client_test.http` after env/credentials are ready.

## Evidence Requirements

Must keep evidence for implementation handoff:

- Discovery artifact: `.opencode/evidence/20260518-2026-secondary-sales-dev-to-qa/discovery.md`
- Open questions: `.opencode/draft/20260518-2026-secondary-sales-dev-to-qa/open-questions.md`

Implementation evidence to collect later:

- `git worktree list` output after creating worktree.
- `git status --short --branch` before/after patch.
- `git diff --stat origin/qa...HEAD` and `git diff --name-status origin/qa...HEAD`.
- Targeted test command outputs.
- Full `rtk go test ./...` output.
- Env OBS readiness note.
- Final quality-gate notes.

Research Gate decision:

- Local project discovery: required and done.
- Official docs/context7: not needed; no library/API behavior change beyond existing Go code.
- GitHub: not needed; private Git remote/local branch compare sufficient.
- Brave/web search: not needed; no external/current facts needed.
- Browser/screenshot: not needed; backend endpoint/report export, no UI parity.

## Done Criteria

Planning done when:

- Primary plan exists at `.opencode/plans/20260518-2026-secondary-sales-dev-to-qa.md`.
- Discovery evidence exists.
- Open questions documented because branch/worktree name collision blocks safe implementation.
- Worklist is actionable for `@orchestrator` after user answers W1/W2/W7/W10 decisions.

Implementation done later when:

- Approved worktree exists.
- Approved sync strategy applied.
- Tests pass.
- Env risk resolved.
- Commit ready on demo branch.
- Quality gate passes.

## Final Planning Summary

Artifacts created:

- `.opencode/plans/20260518-2026-secondary-sales-dev-to-qa.md` — source of truth implementation handoff.
- `.opencode/evidence/20260518-2026-secondary-sales-dev-to-qa/discovery.md` — kept because implementation needs exact branch SHAs, diff scope, existing worktrees, and risk notes.
- `.opencode/draft/20260518-2026-secondary-sales-dev-to-qa/open-questions.md` — kept because unresolved user decisions block safe execution.

Key findings:

- Ada changes nyata `dev` vs `qa` untuk endpoint `POST /v1/reports/secondary-sales`.
- Diff bounded menyentuh controller/service/repository/entity/model/config/constant/main plus tests.
- Perubahan utama: RabbitMQ publish failure handling, OBS required env validation, report filename/report count fixes, SQL product/supplier fallback, canvas invoice inclusion, NULL-safe `gross_sales`, dan dashboard sum COALESCE.
- Branch/worktree `demo-18052026` sudah ada; tidak aman membuat ulang tanpa keputusan user.

Key decisions:

- Planner merekomendasikan bounded squash patch dari `origin/dev` ke branch baru berbasis `origin/qa`.
- Planner merekomendasikan nama aman `demo-18052026-2026` jika user setuju karena `demo-18052026` sudah ada.

Assumptions:

- User ingin plan dulu.
- `DDMMYYYY` memakai tanggal lokal 18 Mei 2026.
- `origin/dev` dan `origin/qa` adalah source of truth.

Remaining open questions:

- Cara menangani `demo-18052026` yang sudah ada.
- Strategi sync final.
- Push/MR atau lokal-only.
- Kesiapan env OBS di QA.

Readiness:

- Plan siap untuk implementation setelah user menjawab open questions.
- Worklist W1 adalah blocker pertama.

Cleanup performed:

- Tidak ada draft/evidence yang dihapus karena keduanya masih operasional untuk handoff.
