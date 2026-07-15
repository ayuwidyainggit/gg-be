# Plan — Backup Code Endpoint dari `dev/staging` sebelum Replace ke `demo/qa`

Task ID: `20260505-1240-backup-endpoint-code`
Tanggal: 2026-05-05 12:40 Asia/Jakarta
Status: **ready untuk eksekusi backup dan validasi replace dengan worktree**

## Goal

Membuat rencana backup lengkap untuk Go code backend endpoint yang disebutkan user dari branch sumber `dev/staging`, agar aman sebelum code di branch `demo/qa` diganti. Backup harus mencakup file `main`/wiring, route/handler/controller, service, repository, entity/model/helper, dan test terkait per module Go (`master`, `inventory`, `sales`, `finance`, `pjp`) tanpa mencampur perubahan lokal yang tidak relevan atau secret. Migration SQL tidak ikut replace karena database sudah siap.

## Non-goals

- Tidak melakukan replace, merge, cherry-pick, checkout destruktif, atau commit implementasi dalam tahap plan ini.
- Tidak mengubah source code, test, migration, config, `.env`, atau lockfile.
- Tidak mengekspos isi secret dari file config/env.
- Tidak membuat backup database/data produksi; scope replace hanya Go code endpoint, bukan migration SQL.
- Tidak mendesain ulang endpoint/API contract.

## Scope

### Module yang masuk scope

1. `master`
   - Manage Survey - Survey
   - Manage Survey - Template
   - Monitoring field option `business-unit`
   - Sales Target - Distributor
   - Master - Distributor
   - Report - Survey List
   - Master - Approval - Outlet Profile
   - Common/field options: areas, regions, distributors, employees, emp-groups, business-unit, sales-teams, salesman, outlet-types, outlet-groups, outlet-classes, outlets, suppliers, warehouses bila endpoint berada di master.
2. `inventory`
   - Stock Disposal
   - Stock Disposal Products
   - Goods Receipt reference route/file terkait.
3. `sales`
   - Order List
   - Add New Order
   - Download Sales Order / Generate Invoice
   - Proforma Invoice
   - Print Invoice
4. `finance`
   - Setup Parameter Web - Expense
   - AR - Collection
5. `pjp`
   - Monitoring Activity: live monitoring principal/distributor/location details.

### Endpoint inventory dari user

#### 1. Manage Survey - Survey (`master`)
- `GET master/v1/survey`
- `POST master/v1/survey`
- `GET master/v1/survey/{id}`
- `PUT master/v1/survey/{id}`
- `PATCH master/v1/survey/{id}/deactivate`
- Field options: `master/v1/areas`, `master/v1/survey_template`, `master/v1/business-unit`, `master/v1/sales-teams`, `master/v1/salesman`, `master/v1/outlet-types`, `master/v1/outlet-groups`, `master/v1/outlet-classes`, `master/v1/outlets`

#### 2. Manage Survey - Template (`master`)
- `GET/POST master/v1/survey_template`
- `GET/PUT/DELETE master/v1/survey_template/{id}`

#### 3. Inventory - Stock Disposal (`inventory`, plus master options)
- `GET/POST inventory/v1/stock-disposal`
- `GET inventory/v1/stock-disposal/{id}`
- `GET inventory/v1/stock-disposal/products`
- `GET master/v1/suppliers`
- `GET master/v1/warehouses`
- `GET inventory/v1/goods-reipt/{id}` dari request user perlu diverifikasi terhadap code karena discovery menemukan route `goods-receipts`.

#### 4. Monitoring Activity (`pjp`, plus master option)
- `GET scylla-pjp/api/v1/live-monitoring-principal`
- `GET scylla-pjp/api/v1/live-monitoring-distributor`
- `GET scylla-pjp/api/v1/monitoring_locations/details`
- `GET master/v1/business-unit`

#### 5-9. Sales (`sales`)
- `GET sales/v1/orders`
- `GET sales/v1/orders/{id}`
- `GET sales/v2/orders/{id}`
- `POST sales/v1/orders`
- `PATCH sales/v1/orders/{id}`
- `GET sales/v1/validate`
- `GET sales/v1/validate-detail`
- `GET sales/v1/stock`
- `GET sales/v1/min-price`
- `GET sales/v1/conversion-product`
- `POST sales/v1/generate-invoice`
- `POST sales/v1/print_proforma_invoice`
- `PATCH sales/v1/invoices/print/{id}`

#### 10. Setup Parameter Web - Expense (`finance`)
- `GET finance/v1/expense`
- `GET finance/v1/expense/{id}`
- `POST finance/v1/expense`
- `PATCH finance/v1/expense/{id}`
- `POST finance/v1/expense/upload`
- Catatan discovery: controller yang ditemukan memakai `/v1/expense-type`; perlu verifikasi apakah `/v1/expense` berada di file lain/branch lain atau nama request adalah alias produk.

#### 11-12. Master (`master`)
- Sales Target Distributor: `GET/POST master/v1/sales-target-distributor`, `GET/PUT master/v1/sales-target-distributor/{id}`
- Distributor: `GET/POST master/v1/distributors`, `GET/PUT/DELETE master/v1/distributors/{id}`
- Options: `areas`, `regions`, `distributors`

#### 13. AR - Collection (`finance`)
- `GET/POST finance/v1/account-receivables/collection`
- `GET/PATCH/DELETE finance/v1/account-receivables/collection/{id}`
- `POST finance/v1/account-receivables/collection/print`
- `GET finance/v1/account-receivables/collection/invoice`

#### 14-15. Report / Approval (`master`)
- `GET master/v1/survey-report`
- `GET master/v1/survey-report/{id}`
- `POST master/v1/survey-report/export`
- `PATCH master/v1/outlet-list/approval`

#### Common options (`master`)
- `GET master/v1/employees`
- `GET master/v1/emp-groups`

## Requirements

1. Backup dilakukan per module/repo karena tidak ada root `go.mod` dan root bukan git repo.
2. Backup dan restore harus selalu diawali sinkronisasi remote terbaru: `git fetch --all --prune`, lalu `git pull --ff-only` di worktree source/target sebelum membuat bundle, patch, atau apply restore.
3. Backup harus dimulai dengan memastikan working tree bersih atau perubahan lokal diselamatkan tanpa mencampur ke backup.
4. Backup harus menyimpan:
   - branch/source commit SHA,
   - target branch commit SHA,
   - daftar file path terkait endpoint/fungsi,
   - patch `--binary`,
   - bundle source branch untuk recovery immutable,
   - metadata bundle per endpoint/fungsi dari rantai aktual `main`, route/handler/controller, service, repository, entity/model/helper, dan test terkait bila ada.
5. Backup hanya mencakup Go code endpoint dan file pendukung Go yang diperlukan. Migration SQL tidak ikut replace karena database sudah siap menurut user.
6. Backup tidak boleh menyertakan `.env`, token, credentials, file temporary, migration SQL, atau perubahan lokal tidak relevan.
7. Endpoint path harus diverifikasi dari route actual, bukan hanya asumsi dari daftar user.
8. Rencana eksekusi harus aman terhadap branch yang berbeda antar module (`dev` vs `staging`, `demo` vs `demo/28-03-2026` vs `qa`).
9. Validasi setelah backup harus memastikan patch bisa diaplikasikan secara kering (`git apply --check`) ke target branch/worktree tanpa langsung mengganti target.
10. Build/test gate wajib: sebelum dinyatakan aman, worktree target yang menerima patch endpoint harus lolos minimal `rtk go test ./...`; bila test penuh gagal karena env/dependency, jalankan compile-only `rtk go test ./... -run '^$'` dan klasifikasikan kegagalan test penuh.

## Acceptance Criteria

- Ada satu folder backup per endpoint/fungsi di bawah module yang memuat minimal:
  - `metadata.md`
  - `file-list.txt`
  - `source-to-target.patch`
  - referensi `source-code.bundle` module
  - `source-status.txt`
  - `target-status.txt`
  - `sha.txt`
- Semua endpoint dalam scope dipetakan ke rantai file `main`/route/handler/controller/service/repository/entity/model/helper/test atau diberi status `needs-verification` dengan alasan.
- Dirty working tree tidak hilang dan tidak ikut backup kecuali user eksplisit menyetujui.
- Tidak ada `.env`/secret yang masuk backup patch atau commit.
- Patch backup lolos `git apply --check` terhadap branch target yang dipilih, atau konflik dicatat lengkap.
- Worktree hasil apply lolos `rtk go test ./...`, atau kegagalan diklasifikasikan sebagai baseline existing/env-dependent dan dibuktikan dengan test baseline sebelum patch.
- Plan ini dipakai sebagai sumber kebenaran sebelum tindakan replace.

## Existing Patterns/Reuse

### Pola repo yang harus digunakan

- Operasi git dilakukan di masing-masing folder module: `master`, `inventory`, `sales`, `finance`, `pjp`.
- Fiber modules memakai route group `/v1/...` di controller, service/repository/domain files mengikuti naming `{domain}_controller.go`, `{domain}_service.go`, `{domain}_repository.go`.
- `pjp` memakai Gin router, contoh `pjp/router/live_monitoring.go`.
- Gunakan route/controller existing sebagai dasar pathspec backup; jangan membuat daftar file manual dari nol tanpa discovery lanjutan.

### Reuse candidates

- Reuse `git diff --binary`, `git format-patch`, `git bundle`, `git worktree` untuk backup aman.
- Reuse existing test files yang ditemukan:
  - `master/service/survey_service_test.go`, `master/controller/survey_controller_test.go`, `master/controller/survey_template_controller_test.go`, `master/controller/survey_report_controller_test.go`, `master/service/sales_target_distributor_service_test.go`, dll.
  - `inventory/service/stock_disposal_service_test.go`
  - `sales/service/order_service_test.go`, `sales/service/invoice_service_concurrency_test.go`
  - `finance/*_test.go` untuk AR/expense bila tersedia.
  - `pjp/data/request/live_monitoring_request_test.go`

## Constraints

- Root compose discovery menunjukkan yang sedang `Up`: `master`, `system`, `redis`; service `inventory`, `sales`, `finance`, `pjp` tidak sedang up.
- Working tree discovery:
  - `inventory`: modified `go.mod`, `go.sum`.
  - `finance`: modified `.env`; jangan stage/backup/print isi.
  - `pjp`: modified `.air.toml`, `.gitignore`.
- Branch discovery:
  - `master`, `inventory`, `sales`, `finance` saat ini di `dev`.
  - `pjp` saat ini di `staging`.
  - Branch demo tidak konsisten antar module: ada `demo`, `demo/28-03-2026`, `demo/28-03-2026-qa`; target harus dipilih eksplisit.
- Instruksi proyek meminta command dengan `rtk`; gunakan `rtk` untuk command repo ini.
- Jangan gunakan perintah destruktif seperti `git reset --hard`, `git clean`, force push, atau overwrite branch tanpa persetujuan eksplisit.

## Risks

1. **Branch target salah**: `demo/qa` tidak seragam antar module.
2. **Backup tidak lengkap**: route actual berbeda dengan endpoint list, terutama `goods-receipt` vs `goods-receipts`, `expense` vs `expense-type`, dan sales validate/stock/min-price/conversion route.
3. **Secret leakage**: `.env` di `finance` modified dan repo mengandung plaintext credentials di beberapa file.
4. **Local changes overwritten**: beberapa module dirty sebelum replace.
5. **Schema drift**: migration terkait survey, stock disposal, expense, sales order/invoice, monitoring indexes sengaja tidak ikut replace karena user menyatakan database sudah siap; jika Go code mengasumsikan schema yang belum ada di target DB, build bisa aman tetapi runtime tetap berisiko.
6. **Per-module branch divergence**: dev/staging mungkin tidak punya commit yang sama antar module dan target.

## Decisions/Assumptions

### Decisions

- Backup diperlakukan sebagai operasi **per module**, bukan root repo.
- Source branch sudah diputuskan user:
  - `master`, `inventory`, `sales`, `finance`: `origin/dev`.
  - `pjp`: `origin/staging`.
- Target branch sudah diputuskan user:
  - `master`: `origin/qa`
  - `inventory`: `origin/qa`
  - `sales`: `origin/qa`
  - `finance`: `origin/qa`
  - `pjp`: `origin/demo`
- Backup harus berupa bundle per endpoint/fungsi. Setiap endpoint/fungsi mengambil rantai file aktual dari `main`/wiring, route/handler/controller, service/use-case, repository/query, entity/request/response/model, helper/util yang dipanggil langsung, dan test terkait bila ada.
- Go code endpoint saja yang masuk patch replace; migration SQL dikecualikan karena database sudah siap.
- Gunakan `git worktree` saja untuk validasi dan replace agar tidak mengganggu working tree kotor; jangan stash otomatis.
- Pull/fast-forward wajib dilakukan di worktree source sebelum backup dan di worktree target sebelum restore/apply. Jika `git pull --ff-only` gagal karena branch diverged, hentikan proses; jangan merge otomatis.
- Backup artifact operasional sebaiknya ditempatkan di folder eksternal yang tidak ikut commit, misalnya `../scylla-be-endpoint-backup-YYYYMMDD-HHMM/` atau path yang user tentukan. Jangan simpan patch besar/secret-risk di repo kecuali user meminta.

### Assumptions / Open Questions

Sudah dijawab user:

1. Source branch: `dev` untuk `master/inventory/sales/finance`, `staging` untuk `pjp`.
2. Target: `qa` untuk `master/inventory/sales/finance`, `demo` untuk `pjp`.
3. Format backup: bundle setiap endpoint/fungsi dari rantai `main`, handler/controller, service, repository, dan file pendukung Go terkait.
4. Scope file: Go code endpoint saja; database/migration tidak ikut.
5. Working tree handling: wajib pakai worktree, tidak stash.
6. Pull dulu sebelum backup dan sebelum restore/apply; gunakan fast-forward only agar tidak membuat merge commit otomatis.

Tidak ada open question branch target yang tersisa berdasarkan jawaban user. Eksekusi tetap harus memverifikasi branch remote tersebut tersedia sebelum membuat worktree.

## TDD/Test Plan

### Apakah TDD diperlukan?

TDD tidak wajib untuk tahap backup murni karena tidak ada perubahan production code. Namun validasi backup/replace wajib dilakukan dengan pendekatan regresi: sebelum replace, identifikasi test endpoint/domain yang sudah ada; setelah patch diterapkan di worktree target, jalankan test terkait untuk memastikan code hasil replace tetap buildable.

### Existing test patterns

- `go test ./...` per module.
- Test domain tersedia di beberapa module sesuai discovery:
  - `master`: survey, survey report, sales target distributor, distributor.
  - `inventory`: stock disposal service.
  - `sales`: order service dan invoice concurrency.
  - `finance`: AR/expense repository/service/controller test sebagian tersedia.
  - `pjp`: request test live monitoring.

### First failing/regression test

Untuk eksekusi replace nanti, Red step adalah menjalankan test di target branch sebelum patch:

```bash
rtk go test ./...
```

per module target. Catat baseline failure yang sudah ada agar tidak dianggap regresi dari replace.

### Green step

Setelah patch source endpoint diterapkan ke worktree target, jalankan:

```bash
rtk go test ./...
```

Untuk service yang bisa dijalankan, lanjut smoke test route dengan JWT/test token yang valid bila tersedia.

### Refactor step

Tidak ada refactor dalam backup. Jika patch konflik, buat rencana resolusi konflik terpisah dan jangan langsung edit source tanpa approval.

### Edge cases

- Route berubah nama antara branch source dan target.
- Migration SQL berbeda antara source dan target tetapi sengaja tidak ikut replace; catat sebagai runtime risk bila Go code butuh schema baru.
- `go.mod`/`go.sum` berbeda antar branch.
- Test gagal karena service remote DB/env, bukan karena patch.
- Endpoint upload/print/export membutuhkan file system atau external dependency.

### Commands

Gunakan command berikut per module di worktree validasi, bukan working tree utama jika masih dirty:

```bash
rtk git status --short --branch
rtk go test ./...
rtk go test ./... -run '^$'
```

`rtk go test ./... -run '^$'` digunakan sebagai compile-only fallback bila test penuh membutuhkan DB/env. Jika compile-only gagal setelah patch, replace harus diblokir karena build tidak aman.

## Implementation Steps

> Bagian ini adalah langkah eksekusi yang harus dilakukan oleh agen implementasi/operator setelah branch target dipastikan. Jangan dilakukan dalam planning mode.

### 1. Persiapan dan freeze scope

1. Gunakan source yang sudah diputuskan:
   - `master`: `origin/dev`
   - `inventory`: `origin/dev`
   - `sales`: `origin/dev`
   - `finance`: `origin/dev`
   - `pjp`: `origin/staging`
2. Gunakan target yang sudah diputuskan:
   - `master`: `origin/qa`
   - `inventory`: `origin/qa`
   - `sales`: `origin/qa`
   - `finance`: `origin/qa`
   - `pjp`: `origin/demo`
3. Buat nama backup run, contoh `endpoint-backup-20260505-1240`.
4. Tentukan folder backup eksternal, contoh `../scylla-be-endpoint-backup-20260505-1240/`.
5. Untuk tiap module (`master`, `inventory`, `sales`, `finance`, `pjp`), jalankan:

```bash
rtk git status --short --branch
rtk git fetch --all --prune
```

6. Jika ada dirty files, jangan checkout di folder utama dan jangan stash otomatis. Semua operasi source/target dilakukan lewat `git worktree`.

7. Buat worktree source dan target dari remote branch yang sudah diputuskan, lalu pull fast-forward sebelum backup/restore:

```bash
rtk git worktree add ../_wt-<module>-source-20260505 <source-branch>
rtk git worktree add ../_wt-<module>-target-20260505 <target-branch>
```

Di masing-masing worktree source dan target:

```bash
rtk git pull --ff-only
rtk git status --short --branch
rtk git rev-parse HEAD
```

Jika `rtk git pull --ff-only` gagal karena branch diverged atau konflik, hentikan proses untuk module tersebut dan minta keputusan. Jangan melakukan merge otomatis.

### 2. Buat metadata backup per module

Untuk setiap module, simpan metadata:

```bash
rtk git rev-parse <source-branch>
rtk git rev-parse <target-branch>
rtk git rev-parse HEAD
rtk git status --short --branch
rtk git log --oneline --decorate -20
```

Output disimpan sebagai:

- `<backup-root>/<module>/sha.txt`
- `<backup-root>/<module>/source-status.txt`
- `<backup-root>/<module>/target-status.txt`
- `<backup-root>/<module>/recent-log.txt`

### 3. Lengkapi pathspec file per endpoint/fungsi

Gunakan discovery awal, lalu perluas dengan route grep. Bundle harus dibuat per endpoint/fungsi, bukan hanya per module. Untuk setiap endpoint, file list harus diturunkan dari call chain aktual:

1. `main.go` atau wiring/dependency injection yang mendaftarkan controller/router.
2. route/handler/controller endpoint.
3. service/use-case yang dipanggil.
4. repository/query layer yang dipanggil.
5. entity/request/response/model yang dipakai langsung.
6. helper/util internal yang dipanggil langsung oleh endpoint.
7. test terkait bila ada.

File Go shared oleh banyak endpoint boleh masuk beberapa bundle endpoint jika memang dipakai oleh chain tersebut. Migration SQL tidak masuk bundle/patch replace.

Pathspec awal untuk discovery per module:

#### `master`

```text
controller/survey_controller.go
service/survey_service.go
repository/survey_repository.go
entity/survey.go
model/survey.go
model/survey_area.go
model/survey_salesman.go
model/survey_outlet.go
model/survey_detail.go
controller/survey_template_controller.go
service/survey_template_service.go
repository/survey_template_repository.go
entity/survey_template.go
model/survey_template.go
controller/survey_report_controller.go
service/survey_report_service.go
repository/survey_report_repository.go
entity/survey_report.go
model/survey_report.go
controller/sales_target_distributor_controller.go
service/sales_target_distributor_service.go
repository/sales_target_distributor_repository.go
entity/sales_target_distributor.go
model/sales_target_distributor.go
controller/m_distributor_controller.go
service/m_distributor_service.go
repository/m_distributor_repository.go
entity/m_distributor.go
model/m_distributor.go
controller/outlet_controller.go
service/outlet_service.go
repository/outlet_repository.go
entity/outlet.go
model/outlet_cr.go
controller/business_unit_controller.go
controller/sales_team_controller.go
controller/outlet_type_controller.go
controller/outlet_group_controller.go
controller/outlet_class_controller.go
controller/employee_controller.go
controller/emp_group_controller.go
```

Tambahkan service/repository/entity/model untuk option endpoints bila diff source-target menunjukkan perubahan.

#### `inventory`

```text
controller/stock_disposal_controller.go
service/stock_disposal_service.go
repository/stock_disposal_repository.go
entity/stock_disposal.go
model/stock_disposal.go
model/stock_disposal_detail.go
controller/gr_controller.go
controller/gr_branch_controller.go
```

#### `sales`

```text
controller/order_controller.go
service/order_service.go
service/validate_order_service.go
repository/order_repository.go
repository/validate_order_repository.go
entity/order.go
entity/order_detail.go
entity/validate_order.go
entity/edit_order_enhance.go
model/order.go
model/order_detail.go
model/validate_order.go
controller/invoice_controller.go
service/invoice_service.go
repository/invoice_repository.go
entity/invoice.go
entity/invoice_detail.go
model/invoice.go
model/invoice_detail.go
```

Catatan: route exact `validate`, `validate-detail`, `stock`, `min-price`, `conversion-product`, `generate-invoice`, dan `invoices/print` harus diverifikasi dengan grep lanjutan karena discovery exact string belum lengkap.

#### `finance`

```text
controller/expense_controller.go
service/expense_service.go
repository/expense_repository.go
entity/expense.go
model/expense.go
model/expense_type.go
model/expense_file.go
controller/*ar*
service/ar_service.go
repository/ar_repository.go
entity/ar*.go
model/*collection*.go
```

Catatan: perlu map controller AR collection exact dan route expense exact karena discovery menunjukkan `expense-type`, bukan `expense`.

#### `pjp`

```text
router/live_monitoring.go
controller/live_monitoring/live_monitoring_controller.go
service/live_monitoring/live_monitoring_service.go
repository/live_monitoring/live_monitoring_repository.go
data/request/live_monitoring_request.go
data/response/live_monitoring_response.go
model/live_monitoring.go
```

### 4. Buat bundle dan patch per endpoint/fungsi

Untuk tiap endpoint/fungsi, buat folder seperti:

```text
<backup-root>/<module>/<endpoint-slug>/
```

Isi minimal:

```text
metadata.md
file-list.txt
source-to-target.patch
apply-check.txt
test-build.txt
```

`git bundle` tidak bisa membatasi hanya path tertentu. Karena itu gunakan dua artefak:

1. `<backup-root>/<module>/<module>-source.bundle`: bundle commit source branch untuk recovery immutable.
2. `<backup-root>/<module>/<endpoint-slug>/source-to-target.patch`: patch scoped path endpoint/fungsi untuk replace aman.

Untuk tiap module, bundle source branch dibuat sekali:

```bash
rtk git bundle create <backup-root>/<module>/<module>-source.bundle <source-branch>
```

Untuk tiap endpoint/fungsi, buat patch scoped:

```bash
rtk git diff --binary <target-branch>...<source-branch> -- <endpoint-pathspecs> > <backup-root>/<module>/<endpoint-slug>/source-to-target.patch
```

Jika shell redirection tidak ingin dipakai langsung, operator dapat menjalankan command tanpa `rtk` filter hanya di lingkungan aman atau memakai script internal yang tidak mencetak secret. Pastikan patch tidak memuat `.env`.

### 5. Validasi patch terhadap target dengan worktree

Untuk tiap module:

```bash
rtk git worktree add ../_wt-<module>-target-20260505 <target-branch>
```

Di worktree target:

```bash
rtk go test ./...
```

Untuk bundle endpoint, validasi patch endpoint satu per satu:

```bash
rtk git apply --check <backup-root>/<module>/<endpoint-slug>/source-to-target.patch
```

Jika `git apply --check` gagal, catat konflik di `<backup-root>/<module>/<endpoint-slug>/conflicts.md` dan jangan lanjut replace otomatis untuk endpoint tersebut.

### 6. Build/test safety gate

Setelah semua patch endpoint yang akan direplace bisa di-apply ke worktree target, jalankan build/test gate per module:

```bash
rtk go test ./...
rtk go test ./... -run '^$'
```

`rtk go test ./...` adalah gate utama. `rtk go test ./... -run '^$'` adalah compile-only fallback untuk membedakan build error dari test integration yang membutuhkan DB/env. Jangan klaim aman bila compile-only gagal. Klasifikasi kegagalan:

1. **Baseline existing**: gagal juga sebelum apply patch di target worktree.
2. **Env/dependency**: membutuhkan DB/service/env yang tidak tersedia; catat dependency dan endpoint terdampak.
3. **Patch regression**: gagal hanya setelah apply patch; blok replace sampai diperbaiki.

### 7. Replace branch target hanya setelah backup tervalidasi

Setelah user menyetujui dan backup tervalidasi, implementer dapat memilih strategi:

1. **Patch apply scoped**: apply patch hanya path endpoint ke branch target.
2. **Cherry-pick commit terpilih**: bila endpoint source punya commit bersih dan tidak membawa perubahan lain.
3. **File restore scoped**: restore file pathspec dari source ke target, lalu test.

Strategi rekomendasi: patch apply scoped karena backup sudah membatasi path endpoint.

## Expected Files to Change

Tidak ada source file yang berubah pada tahap plan ini.

Jika eksekusi backup/replace dilakukan nanti, file yang mungkin berubah di target branch adalah file dalam pathspec module pada `Implementation Steps`, terutama:

Expected files dibatasi ke Go code dan test terkait; migration SQL tidak ikut replace:

- `master/main.go`, `master/controller|service|repository|entity|model|pkg` yang dipakai endpoint.
- `inventory/main.go`, `inventory/controller|service|repository|entity|model|pkg` yang dipakai endpoint.
- `sales/main.go`, `sales/controller|service|repository|entity|model|pkg` yang dipakai endpoint.
- `finance/main.go`, `finance/controller|service|repository|entity|model|pkg` yang dipakai endpoint.
- `pjp/main.go`, `pjp/router|controller|service|repository|data|model|pkg` yang dipakai endpoint.

## Agent/Tool Routing

- `@artifact-planner`: sudah digunakan untuk plan dan artefak `.opencode`.
- `@explorer`: direkomendasikan bila implementer perlu discovery route lanjutan yang luas sebelum patch final.
- `@oracle`: opsional bila muncul konflik branch besar atau keputusan strategi replace berisiko.
- `@security-privacy-reviewer`: direkomendasikan sebelum backup/commit jika patch berpotensi menyentuh config, upload, auth, payment/finance, atau secret.
- `@release-engineer`: direkomendasikan bila replace akan dipush ke environment `demo/qa` dengan CI/CD, deployment, rollback, atau migration.
- `@fixer`: hanya setelah plan disetujui untuk menjalankan backup/replace scoped dan validasi, bukan dalam planning mode.

## Validation Commands

### Preflight per module

```bash
rtk git status --short --branch
rtk git fetch --all --prune
rtk git pull --ff-only
rtk git branch --list "dev*" "staging" "demo*" "qa"
rtk git branch -r --list "origin/dev*" "origin/staging" "origin/demo*" "origin/qa"
```

### Backup integrity

```bash
rtk git rev-parse <source-branch>
rtk git rev-parse <target-branch>
rtk git rev-parse HEAD
rtk git diff --name-status <target-branch>...<source-branch> -- <pathspecs>
rtk git apply --check <backup-root>/<module>/<endpoint-slug>/source-to-target.patch
```

### Test/build per module

```bash
rtk go test ./...
rtk go test ./... -run '^$'
```

`rtk go test ./...` wajib dicoba. `rtk go test ./... -run '^$'` wajib dicoba bila test penuh gagal untuk memastikan compile/build tetap aman.

### Optional smoke tests

Jika service dan JWT/test token tersedia, lakukan smoke test endpoint penting:

- `GET /ping` untuk Fiber services yang berjalan.
- `GET /v1/survey`, `GET /v1/survey_template`, `GET /v1/stock-disposal`, `GET /v1/orders`, `GET /v1/account-receivables/collection`, `GET /api/v1/live-monitoring-principal` sesuai base URL module.

## Evidence Requirements

Sebelum replace dianggap siap, simpan evidence berikut di folder backup eksternal:

- Source/target SHA per module setelah `git pull --ff-only`.
- Status working tree sebelum backup.
- File list final per endpoint/fungsi.
- Bundle source branch per module dan patch scoped per endpoint/fungsi.
- Hasil `git apply --check` per endpoint/fungsi.
- Hasil `go test ./...` dan/atau compile-only `go test ./... -run '^$'` per module atau alasan tidak bisa dijalankan.
- Catatan route mismatch dan resolusinya:
  - `inventory/v1/goods-receipt/{id}` vs route code `goods-receipts`.
  - `finance/v1/expense` vs route code `expense-type`.
  - sales validate/stock/min-price/conversion/generate-invoice route exact.
- Catatan secret scan manual: pastikan `.env`, credentials, token tidak ada di backup patch.
- Catatan bahwa migration SQL tidak masuk patch/bundle endpoint replace.

### Research Gate

- Local project discovery: **sudah dilakukan** dan wajib dilanjutkan sebelum eksekusi final untuk route yang masih mismatch.
- Official docs/context7: **tidak diperlukan** karena operasi memakai Git/Go standard dan pola lokal lebih menentukan.
- GitHub upstream: **tidak diperlukan** kecuali branch remote tidak lengkap atau perlu PR/Actions status.
- Brave/web search: **tidak diperlukan** karena tidak bergantung informasi eksternal.
- Browser/screenshot: **tidak diperlukan** karena task backend backup code, bukan UI/visual.

## Done Criteria

Backup/replace workflow boleh dianggap selesai bila:

1. Source branch dan target branch exact per module sudah diverifikasi tersedia di remote/local.
2. Source dan target worktree sudah `git pull --ff-only` sebelum backup dan sebelum restore/apply.
3. Dirty working tree ditangani aman tanpa kehilangan perubahan lokal.
4. Backup metadata + bundle source + patch scoped dibuat per endpoint/fungsi.
5. Patch endpoint/fungsi lolos `git apply --check` di worktree target atau konflik terdokumentasi.
6. Test/build baseline dan test/build setelah apply dijalankan atau blocker dicatat.
7. Tidak ada secret/local `.env` dalam backup/commit.
8. User menerima laporan backup path dan daftar konflik/risiko sebelum replace.

## Final Planning Summary

### Artifacts created/consulted

- Primary plan: `.opencode/plans/20260505-1240-backup-endpoint-code.md`
- Discovery evidence dibuat dan kemudian dikonsolidasikan: `.opencode/evidence/20260505-1240-backup-endpoint-code/discovery.md`

### Key decisions

- Backup harus per module/repo mandiri.
- Source sudah dipastikan user: `origin/dev` untuk `master/inventory/sales/finance`, `origin/staging` untuk `pjp`.
- Gunakan worktree saja untuk menghindari gangguan pada dirty working tree; tidak stash.
- Pull dulu dengan `rtk git pull --ff-only` di worktree source sebelum backup dan di worktree target sebelum restore/apply.
- Backup harus per endpoint/fungsi dari chain `main`/route/handler/controller/service/repository/entity/model/helper/test terkait.
- Patch backup harus scoped by endpoint pathspec dan tervalidasi dengan `git apply --check`.
- Build/test gate wajib sebelum dinyatakan aman: `rtk go test ./...` atau compile-only `rtk go test ./... -run '^$'` bila test integration membutuhkan env.
- Migration SQL tidak ikut karena database sudah siap.

### Assumptions

- User ingin backup Go code endpoint sebelum replace, bukan backup DB/data.
- Target sudah final: `qa` untuk `master/inventory/sales/finance`, `demo` untuk `pjp`.

### Remaining open questions

Tidak ada open question branch yang tersisa. Satu-satunya gate sebelum replace adalah verifikasi remote branch tersedia dan build/test gate berhasil di worktree patched.

### Readiness for implementation

Plan siap dipakai sebagai handoff. Eksekusi backup dari source dapat dimulai dengan worktree dan bundle/patch scoped per endpoint. Replace ke target aman dimulai setelah branch remote diverifikasi tersedia dan build/test gate di worktree patched lolos atau failure terklasifikasi sebagai baseline/env-dependent.

### Cleanup performed

Evidence discovery dipertahankan karena masih operasional untuk eksekusi berikutnya: `.opencode/evidence/20260505-1240-backup-endpoint-code/discovery.md` berisi mapping route awal, status branch, dirty files, dan risiko yang belum semuanya terselesaikan.
