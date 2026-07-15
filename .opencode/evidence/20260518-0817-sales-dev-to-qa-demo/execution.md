# Execution evidence: demo-18052026

Tanggal: 2026-05-18
Service: `sales`
Worktree: `/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260518/sales`
Branch: `demo-18052026`
Commit hasil: `1e7586f fix(order): normalize available stock breakdown for sales order detail`

## Worktree

```text
/Users/ujang/Projects/Geekgarden/scylla-be/sales                                 8c995a5 [dev]
/Users/ujang/Projects/Geekgarden/scylla-be-restore-worktrees-20260505/sales      5aaf8d3 [demo-05052026]
/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260505-1300/sales-source  3e534e4 (detached HEAD)
/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260505-1300/sales-target  a92c328 (detached HEAD)
/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260513/sales              4c5a016 [demo-13052026]
/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260518/sales              1e7586f [demo-18052026]
```

## File dibawa dari dev ke qa scope

- `go.mod`
- `service/order_service.go`
- `service/order_service_test.go`
- `service/order_stock_helper.go`
- `service/order_stock_helper_test.go`

## Diff terhadap qa

```text
go.mod
service/order_service.go
service/order_service_test.go
service/order_stock_helper.go
service/order_stock_helper_test.go
```

Stat:

```text
 go.mod                             |   2 +-
 service/order_service.go           |  85 +++++++++-----------------
 service/order_service_test.go      | 119 ++++++++++++++++++++++++++++++++-----
 service/order_stock_helper.go      |  73 +++++++++++++++++++++++
 service/order_stock_helper_test.go |  45 ++++++++++++++
 5 files changed, 249 insertions(+), 75 deletions(-)
```

## Validation

### Targeted tests

Command:

```bash
rtk go mod download && rtk go mod tidy && rtk go test ./service -run 'TestDetailV2|TestCanonicalAPIStockBreakdown|TestComputeDisplayedAvailableStockBreakdown|TestStore_DoesNotPersistStockSnapshotDuringInitialCreate'
```

Hasil:

```text
Go test: 16 passed in 1 packages
```

### Full module tests

Command:

```bash
rtk go test ./...
```

Hasil final:

```text
Go test: 146 passed in 22 packages
```

## Issue tambahan yang ditemukan saat eksekusi

- Full test pertama gagal pada `TestPromoV2Only_CSVAcceptanceCoverageForReferenceSOs`.
- Akar masalah bukan logika endpoint target, tapi helper test `readCSVReferenceFixture` tidak mengenali layout worktree baru di luar folder `scylla-be` utama.
- Solusi: tambah path candidates ke `service/order_service_test.go` agar fixture CSV tetap terbaca pada layout worktree `/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260518/sales`.
- Karena perubahan ini ada di file test saja dan dibutuhkan agar full module validation hijau, perubahan dianggap in-scope untuk validasi.

## Migration impact

- Tidak ada file migration yang dibawa.
- Scope final tidak menambah field DB baru.
- Tidak ada migration wajib untuk backport ini.

## Notes

- `go.mod` berubah karena `rtk go mod tidy` memindahkan `github.com/jackc/pgx/v5` dari indirect ke direct, konsisten dengan import yang kini dipakai langsung.
- `repository/order_repository.go` tidak dibawa karena beda hanya formatting, tidak ada perubahan logika.
