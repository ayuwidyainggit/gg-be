# Discovery — secondary-sales endpoint dev → qa

Tanggal: 2026-05-18 (Asia/Jakarta).
Repo (Go module): `sales/` di dalam monorepo `scylla-be`.
Worktree planner aktif: `/Users/ujang/Projects/Geekgarden/scylla-be/sales` di branch `dev` (HEAD `b7eafd9`).
Remote `origin/qa` HEAD: `2ec0152`.
Remote `origin/dev` HEAD: `b7eafd9`.
URL produksi yang ditanyakan: `https://best.scyllax.online/sales/v1/reports/secondary-sales` → service `sales`.

## Pemetaan endpoint

Route terdaftar di `sales/controller/report_controller.go`:

- `POST /v1/reports/secondary-sales` → `ReportController.SecondarySales`
- `GET /v1/reports/secondary-sales/sum-date` → `SecondaryReportSalesSumMonth`
- `GET /v1/reports/secondary-sales/group` → `SecondaryReportSalesGroup`
- `GET /v1/reports/secondary-sales/trend-sales` → `SecondaryReportSalesTrendSales`
- `POST /v1/extract/secondary-sales` → `SecondarySalesDashboardExtract`

Lapisan dependency (Controller → Service → Repository):

- Controller: `sales/controller/report_controller.go`
- Service: `sales/service/report_service.go` (`PublishSecondarySalesReport`, `SubscribeSecondarySalesReport`, dsb.)
- Repository: `sales/repository/report_repository.go` (`SecondarySales`, `SecondarySalesUnion`, `SecondarySalesUnionPagination`, helper `buildSecondarySalesUnionQuery`, `secondarySalesProductSelect`, `secondarySalesProductJoins`, `CountSecondarySalesReportByDate`)
- Entity/model: `sales/entity/report.go`, `sales/model/report.go`
- Constant/env yang dipakai endpoint: `sales/pkg/constant/constant.go`, `sales/pkg/config/env/env.go`
- Bootstrap: `sales/main.go`
- Test: `sales/service/report_service_test.go`, `sales/repository/report_repository_test.go`

## Diff stat `origin/qa..origin/dev` untuk file terkait endpoint

```
controller/report_controller.go       |  16 +/-
entity/report.go                      |   8 +/-
model/report.go                       |  40 +/-
pkg/config/env/env.go                 |  17 +
pkg/constant/constant.go              |   2 +/-
repository/report_repository.go       | 117 +/-
service/report_service.go             |  23 +/-
main.go                               |  10 +
repository/report_repository_test.go  |  30 +
service/report_service_test.go        |  55 +
```

## Commit yang relevan dengan endpoint (urut paling lama → terbaru, di `dev`, belum ada di `qa`)

- `c278c29` fix(report): activity sales geotag & visit join for list/export
- `107f5bf` fix activity report view
- `c8461cd` take out geotag_status_desc from export salesman activity
- `6355e7a` fix(returns): resolve product/unit join for distributor-scoped master data
- `d53e4c3` fix(repository): use fallback product and supplier fields in report queries
- `bb92526` fix(repository): prefer customer product fields in report queries
- `8586507` fix(service): build report object filenames from report id
- `e8c8662` fix(repository): count secondary sales by report name prefix
- `7bdd374` fix(report): handle secondary sales export failures
- `8fad42c` fix(repository): correct supplier lookup in secondary sales query
- `9adb94d` fix(secondary-sales): include canvas invoice rows in report
- `31ba66d` fix(secondary-sales): add COALESCE to gross_sales calculation to prevent NULL propagation
- `54a63b8` fix(secondary-sales): restore report export completion

Catatan tambahan:

- `qa` punya commit `45beac3 chore(qa): restore sales service from dev backup and add GET alias for print proforma invoice` yang tidak ada di `dev`. Artinya `qa` bukan ancestor murni dari `dev`; merge cherry-pick langsung tanpa merge commit dari `qa` perlu hati-hati.
- Beberapa commit di atas juga menyentuh helper bersama (`activity-report` dan supplier/product join) sehingga cherry-pick parsial endpoint saja akan hampir pasti memicu konflik di `repository/report_repository.go` dan `service/report_service.go`.

## Worktree yang sudah ada (penting untuk konflik nama)

`git worktree list`:

```
/Users/ujang/Projects/Geekgarden/scylla-be/sales                                 b7eafd9 [dev]
/Users/ujang/Projects/Geekgarden/scylla-be-restore-worktrees-20260505/sales      5aaf8d3 [demo-05052026]
/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260505-1300/sales-source  3e534e4 (detached HEAD)
/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260505-1300/sales-target  a92c328 (detached HEAD)
/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260513/sales              4c5a016 [demo-13052026]
/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260518/sales              1e7586f [demo-18052026]
```

Branch `demo-18052026` sudah ada (HEAD `1e7586f`), termasuk worktree-nya. Membuat ulang worktree `demo-18052026` baru akan ditolak Git. Lihat open questions untuk penamaan.

## Konstrain repo & arsitektur

- Service `sales` adalah Go module sendiri (`sales/go.mod`); Fiber stack, layering Controller → Service → Repository → DB harus dijaga.
- Kebijakan repo (`AGENTS.md`): perintah shell pakai prefix `rtk`; `rtk go mod download && rtk go mod tidy && rtk go test ./...` adalah baseline validasi service `sales`.
- Tidak ada migration baru untuk endpoint ini; perubahan murni di kode + helper SQL inline.
- Tidak ada perubahan kontrak HTTP request/response: signature `POST /v1/reports/secondary-sales` tidak berubah; perubahan bersifat behavior fix (NULL safe gross_sales, sertakan invoice canvas, supplier lookup benar, report name prefix benar, build object filename dari report id, kegagalan publish ditangani, validasi env OBS di startup).

## Risiko utama

- Konflik merge di `repository/report_repository.go` dan `service/report_service.go` saat sinkronisasi sebagian (banyak commit menyentuh file yang sama).
- `main.go` menambah `env.ValidateRequired(...)` untuk OBS_HUAWEI_*; jika env target QA belum lengkap, service gagal boot.
- Helper `effectiveGeotagStatus`, `geotagStatusDescFrom`, `activityReportLocationActual` dipakai bersama oleh activity-report; bawa parsial bisa membuat compile error kalau tidak ikut.
- Struktur worktree existing menempatkan worktree per-tanggal di sibling folder root (`scylla-be-worktrees-YYYYMMDD/sales`), bukan di dalam `scylla-be`. Plan harus mengikuti pola itu.

## Hasil eksekusi (2026-05-18 21:26 Asia/Jakarta)

### Worktree dibuat

```
git worktree add -b demo-18052026-2026 /Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260518-2026/sales origin/qa
```

Branch `demo-18052026-2026` dibuat dari `origin/qa` (HEAD `2ec0152`).
Path: `/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260518-2026/sales`

### Patch diterapkan

Bounded checkout dari `origin/dev` untuk 10 file scope:

```
git checkout origin/dev -- controller/report_controller.go service/report_service.go repository/report_repository.go entity/report.go model/report.go pkg/config/env/env.go pkg/constant/constant.go main.go repository/report_repository_test.go service/report_service_test.go
```

Tidak ada konflik. Tidak ada file di luar scope yang tersentuh.

### Test results

```
rtk go test ./repository -run 'TestBuildSecondarySalesUnionQuery'
→ 3 passed

rtk go test ./service -run 'TestPublishSecondarySalesReportMarksFailedWhenRabbitMQPublishFails|TestSubscribeSecondarySalesReportUploadsWorkbookWithExpectedValues'
→ 2 passed

rtk go test ./...
→ 150 passed in 22 packages
```

### Quality gate

Verdict: `PASS_WITH_RISKS`
Blocker: tidak ada.
Risks (non-blocking):
- MEDIUM: `main.go` `env.ValidateRequired` akan panic saat boot jika `OBS_HUAWEI_*` belum ada di env QA.
- MEDIUM: `godotenv.Load` panic jika `.env` file tidak ada (pre-existing, bukan regresi baru).
- LOW: `publishExportMessage` sekarang synchronous; RabbitMQ latency/failure langsung surfacing ke user.
- LOW: `msg.Ack(false)` menggantikan `msg.Ack(true)` — behavior fix yang benar tapi perlu dicatat di release notes.
- LOW: trailing space inconsistency di `model/report.go` field `ActualLongitude`/`ActualLatitude` (kosmetik).
- INFO: `ValidateRequired` belum punya unit test sendiri.

### Commit lokal

```
0266527 fix(sales): port secondary-sales endpoint fixes to QA base
```

Branch `demo-18052026-2026` ahead 1 dari `origin/qa`.
Working tree clean.
