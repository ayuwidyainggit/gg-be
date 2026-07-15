# A3 — Full verification and runtime smoke

## Test / build

`confirmed_runtime` dari `master` modul:

```text
rtk go test ./...     -> 408 passed in 23 packages
rtk go build ./...    -> Success
```

Tidak ada regresi. Acceptance Criteria #12 (full test + build exit 0) terpenuhi untuk layer kompilasi/test.

## Runtime smoke

Status: `not-ready`. Alasan:

- `rtk docker compose -f docker-compose.yml ps` dari run sebelumnya: tidak ada service yang berjalan.
- `$TOKEN` bearer tidak dikonfigurasi pada environment run ini.
- Plan Constraints/Do Not/Reject If: klaim runtime/staging verification tanpa evidence dan env/token configured dilarang.
- A3 remediation tambahan (menjalankan `rtk docker compose -f docker-compose.yml up -d` dan curl smoke) adalah tindakan multi-service yang membutuhkan approval eksplisit user sebelum lanjut, karena:
  1. compose up -d menarik image, membuat kontainer untuk semua service (system, master, inventory, sales, finance, tms, mobile, pjp, cronjob, redis, rabbitmq).
  2. Plan Constraints only allows local Postgres target first; compose membawa seluruh stack.
  3. `clone_db.sh`/restore dan `install_citus.sh` di `scripts/` butuh approval terpisah.

## Diff Boundary check

Perubahan A1+A2+remediation hanya menyentuh:
- `master/entity/product.go`
- `master/service/product_service.go`
- `master/controller/product_controller.go`
- `master/controller/product_report_controller_test.go` (new)
- `master/repository/product_repository.go`
- `master/repository/product_report_repository_test.go` (new)

Tidak menyentuh: migrations, JWT middleware, package manifest, compose/env, legacy GET product list, modul lain. Boundary aman.

## Open assumption

A3 (plan assumption A3): enabled mapping dengan parent yang tidak eligible (missing/inactive) — query LEFT JOIN natural behavior dapat menghasilkan null primary fields. Belum ada product-owner decision; flag ke Q1 review.

## Rekomendasi Q1 review

- Approve A1+A2+A3-test/build.
- Minta product owner memutuskan A3 mapping-parent-missing behavior sebelum release.
- Minta approval terpisah untuk compose up + curl smoke jika env/token siap.
