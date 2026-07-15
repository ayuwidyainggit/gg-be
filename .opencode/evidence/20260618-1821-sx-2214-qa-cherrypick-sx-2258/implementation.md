# Implementation Evidence — SX-2214 + SX-2258 ke branch `bugfix/SX-2214-qa`

Task ID: `20260618-1821-sx-2214-qa-cherrypick-sx-2258`
Mode: Maintenance Stability Mode (git workflow)
Tanggal: 2026-06-18 Asia/Jakarta
Branch hasil: `sales/bugfix/SX-2214-qa` (base = `qa`)
Push: tidak (per user decision).

## Ringkas

Source commits dari `sales/dev` di-apply ke branch baru `bugfix/SX-2214-qa` (base = `qa`):

| Source SHA | Subject | Status |
|---|---|---|
| `fd8817e` | fix(invoice): use final order totals | applied as new commit `0f61997 fix(invoice): use final order totals` (conflict resolution: outlet-status block di-drop) |
| `4ebacfe` | fix(secondary-sales): subtract returns from sold metrics | skip (equivalent in qa via `ca2e5d7 feat(report): port secondary-sales from dev` dll) |
| `4d258f9` | fix(secondary-sales): drop return-status filter from summary | applied as new commit `64e7076 fix(secondary-sales): drop return-status filter from summary` (conflict resolution: keep dev test block) |

Branch tip saat ini (after squash per user decision):

```text
64e7076 fix(secondary-sales): drop return-status filter from summary
0f61997 fix(invoice): use final order totals
78a0c9c (qa) Merge branch 'bugfix/SX-2154-qa' into 'qa'
```

Note: cherry-pick selalu membuat commit baru (SHA baru), jadi source SHA tidak muncul literal. Yang masuk history branch adalah 2 commit hasil cherry-pick (`0f61997` + `64e7076`).

## Worklog

- T1: `git -C sales fetch --all --prune` — exit 0. Catatan: banyak remote branch `bugfix/SX-XXXX-*` sudah dihapus di server, tapi `dev` dan `qa` masih ada dan reachable.
- T2: `git -C sales checkout qa && git -C sales pull --ff-only` — `qa` up to date. Tip `qa` = `78a0c9c`.
- T3: `git -C sales checkout -b bugfix/SX-2214-qa qa` — branch baru dibuat.
- T4: `git -C sales cherry-pick fd8817e` — conflict 1 file (`service/invoice_service.go`).
- T5: Conflict SX-2214 resolved dengan take dev block (memang call `UpdateOutletStatusFromPreDormantIfSet`). Commit `dc5dda6` dibuat.
- T6: `git -C sales cherry-pick 4ebacfe` — empty patch. `git cherry HEAD 4ebacfe` menunjukkan perubahan sudah equivalent di `qa` via commit `ca2e5d7` dll. Skip dengan `--skip`. (Per user decision: terima skip.)
- T7: `git -C sales cherry-pick 4d258f9` — conflict 1 file (`controller/so_controller_test.go`).
- T8: Conflict SX-2258 resolved dengan menghapus marker `<<<<<<< HEAD` / `=======` / `>>>>>>> 4d258f9`. Dev-side test block (test functions untuk `TestSecondaryReportSalesSumMonth*` dan `TestSecondaryReportSalesGroup*`) tetap dimasukkan. Commit `4093c68` dibuat.
- T9: `rtk go mod download && rtk go mod tidy` di `sales/` — sukses.
- T10: `rtk go test ./service -run TestInvoice` — **fail build** karena `service.InvoiceRepository.UpdateOutletStatusFromPreDormantIfSet undefined`. Conflict resolution SX-2214 mempertahankan call ke method yang hanya ada di dev.
- T10 (remedy): hapus blok outlet-status dari `service/invoice_service.go` (9 baris, di luar scope SX-2214 fix; qa tidak punya method). Commit `d0a17e0` dibuat untuk resolution.
- T10 (re-run): `rtk go test ./service -run TestInvoice` — **8 passed**.
- T11: `rtk go test ./repository -run TestSecondarySalesReport` — **11 passed**.
- T12: `rtk go test ./service -run TestReport` — sebelumnya build fail, setelah T10 remedy **No tests found** (tidak ada test function dengan prefix `TestReport` di package `service`; secondary sales service tests pakai prefix `TestSecondarySales*` yang sudah termasuk di T13).
- T13: `rtk go test ./...` — **236 passed in 22 packages**.

### Quality gate round 1

- `@quality-gate` return `NEEDS_FIX` dengan blocker:
  1. cherry-pick set tidak literal (`4ebacfe` skip, `fd8817e` jadi `dc5dda6`, `4d258f9` jadi `4093c68`, plus extra `d0a17e0`).
  2. evidence file log terpisah tidak dibuat (cuma `implementation.md`).
  3. `model/invoice.go` 82-line change di luar plan's expected file list.
  4. `service/report_service.go` (non-test) modified tanpa plan list.

### User decisions (question gate)

- `4ebacfe` skip: terima, evidence jelaskan.
- `d0a17e0` resolution: squash ke SX-2214 commit.
- Evidence: single file (`implementation.md` + `quality-gate.md`).
- Diff scope: update plan scope sesuai hasil.

### Post-decision rewrite

- `git -C sales reset --soft qa` — undo 3 commit cherry-pick (staged).
- Re-stage SX-2214 files (model + repository + service) → commit `0f61997 fix(invoice): use final order totals`. (outlet-status block sudah absent dari working tree dari remedy sebelumnya.)
- Re-stage SX-2258 files (controller/so_controller_test.go + report_repository + report_service_test) → commit `64e7076 fix(secondary-sales): drop return-status filter from summary`.
- Re-run T10–T13:
  - `rtk go test ./service -run TestInvoice`: 8 passed.
  - `rtk go test ./repository -run TestSecondarySalesReport`: 11 passed.
  - `rtk go test ./...`: 236 passed in 22 packages.

## Conflict resolution notes

### SX-2214 (`fd8817e`)

File: `sales/service/invoice_service.go`. Region: `BulkUpdate` transaction block setelah `Update(...)`. Konflik karena qa tidak punya `UpdateOutletStatusFromPreDormantIfSet` call.

Resolusi: cherry-pick ambil dev block (call method baru). Hasil: build fail karena `repository.InvoiceRepository` interface di qa belum expose method.

Remedy: hapus block outlet-status dari `invoice_service.go` (9 baris). Method ini terkait fitur dev `order_type / taking-order flow` (lihat `4e389fe feat(order): port order_type/taking-order flow to match dev`) yang di luar scope SX-2214 fix. Setelah squash, blok ini tidak ada di final commit `0f61997`.

`mockInvoiceRepositoryFinal.UpdateOutletStatusFromPreDormantIfSet` di `service/invoice_service_test.go` tetap ada sebagai unused method — Go interface duck-typing izinkan extra method pada mock.

### SX-2258 #2 (`4d258f9`)

File: `sales/controller/so_controller_test.go`. Region: antara `TestSecondaryReportSalesGroupReturns*` dan `TestSecondarySalesExportReturnsForbiddenForDistributorSiblingCust`. Konflik karena qa tidak punya test block dev (test untuk `TestSecondaryReportSalesSumMonth*` dan `TestSecondaryReportSalesGroup*`).

Resolusi: ambil dev block (test functions baru untuk secondary report endpoint). Marker `<<<<<<< HEAD` / `=======` / `>>>>>>> 4d258f9` dihapus. Tidak ada logika bisnis yang berubah.

## Diff stat `qa..bugfix/SX-2214-qa`

```text
 controller/so_controller_test.go     | 247 +++++++++++++++++++++
 model/invoice.go                     |  82 ++++---
 model/invoice_detail.go              |  11 +
 repository/invoice_repository.go     |   4 +-
 repository/report_repository.go      |   2 +-
 repository/report_repository_test.go |   2 +-
 service/invoice_amount.go            |  62 ++++++
 service/invoice_service.go           | 104 +++++++---
 service/invoice_service_test.go      | 417 +++++++++++++++++++++++++++++++++++
 service/report_service_test.go       |  58 +++--
 10 files changed, 905 insertions(+), 75 deletions(-)
```

File sesuai scope (per user decision, scope di-update):
- SX-2214: `model/invoice.go`, `model/invoice_detail.go`, `repository/invoice_repository.go`, `service/invoice_amount.go` (new), `service/invoice_service.go`, `service/invoice_service_test.go` (new).
- SX-2258: `repository/report_repository.go`, `repository/report_repository_test.go`, `service/report_service_test.go`, `controller/so_controller_test.go` (test block baru).

Tidak ada file di luar scope (tidak ada `go.mod` / `go.sum` / migration / entity order_type / open_api).

## Validation output

| Step | Command | Result |
|---|---|---|
| T9 | `rtk go mod download && rtk go mod tidy` (di `sales/`) | exit 0 |
| T10 | `rtk go test ./service -run TestInvoice` | 8 passed in 1 packages |
| T11 | `rtk go test ./repository -run TestSecondarySalesReport` | 11 passed in 1 packages |
| T12 | `rtk go test ./service -run TestReport` | No tests found (regex tidak match; secondary sales service tests pakai prefix `TestSecondarySales*`) |
| T13 | `rtk go test ./...` | 236 passed in 22 packages |

## Pre-existing fail

Tidak ada baseline test fail di qa yang tercatat. T10 dan T12 fail murni karena cherry-pick SX-2214 menarik dependency yang tidak ada di qa; setelah T10 remedy, semua hijau.

## Risk residual

- `service/invoice_service_test.go` punya mock method `UpdateOutletStatusFromPreDormantIfSet` yang unreferenced (Go duck-typing izinkan). Cleanup follow-up bisa hapus method mock setelah branch stabil.
- Branch history local sudah di-squash (no force-push, no amend, no `-i` interactive conflict — pakai `reset --soft` + recreate commit). User dapat push kapanpun.

## Done criteria

- [x] `bugfix/SX-2214-qa` branch exists di `sales/`, base = `qa`.
- [x] Source commit SX-2214 (`fd8817e`) dan SX-2258 #2 (`4d258f9`) di-apply sebagai commit baru di branch (`0f61997`, `64e7076`). SX-2258 #1 (`4ebacfe`) skip karena equivalent di qa.
- [x] `diff --stat qa..bugfix/SX-2214-qa` hanya berisi file invoice + secondary-sales scope.
- [x] `rtk go test ./...` hijau (236/236).
- [x] Tidak ada secret / `.env` / file sync DB ter-commit.
- [x] Tidak ada push otomatis.

## Push

Tidak push. User push manual via:

```bash
git -C sales push -u origin bugfix/SX-2214-qa
```

## Plan compliance

Plan file updated post-eksekusi:
- `Expected Files to Change` section: ditambah `model/invoice.go` ke scope.
- `Evidence Requirements` section: disederhanakan jadi single evidence file (per user decision).

Plan tetap source of truth; deviasi dicatat di sini dan di plan.
