# Discovery Evidence — Compare Staging dan Demo Schema

Task ID: `20260505-1022-compare-staging-demo-schema`

## Files/Area yang Diinspeksi

- Root repo `/Users/ujang/Projects/Geekgarden/scylla-be` berisi multi-module Go services tanpa root `go.mod`.
- `AGENTS.md` lokal: aturan repo menyebut root compose, service ports, migration workflow untuk `tms`, `pjp`, dan `pjp-principle`, serta peringatan ada plaintext credentials di repo.
- `docker-compose.yml` dicek via `docker compose -f docker-compose.yml ps`; service yang aktif saat discovery hanya `system`, `master`, dan `redis`.
- `**/go.mod`: modul ditemukan di `inventory`, `master`, `pjp`, `finance`, `system`, `sales`, `pjp-principle`, `tms`, `pjp-sales`, `cronjob`, dan `mobile`.
- Migration folders/patterns:
  - `*/migration/**.sql` untuk beberapa service seperti `master`, `mobile`, `finance`, `inventory`, `sales`, `pjp-sales`.
  - `pjp/database/migrate/*.sql` dengan format `.up.sql`/`.down.sql` dan beberapa single `.sql`.
  - `pjp-principle/database/migrate/*.sql` dengan format versioned `.up.sql`/`.down.sql`.
  - `tms/migrations/*.sql` dengan campuran `.up.sql`, `.down.sql`, dan single `.sql`.
- `tms/Makefile` memakai `migrate -path migrations -database $(DATABASE_URL) -verbose up/down/force/drop`.
- `pjp/Makefile` memakai `migrate -path database/migrate -database $(DATABASE_URL) -verbose up/down/force/drop`.
- `pjp/utils/migrate.go` dan `tms/migrations/migrate.go` berisi helper `AutoMigrate` GORM untuk table awal, tetapi repo juga memakai SQL migrations.
- `docs/DATABASE.md` menyebut manual migration memakai SQL files di `migration/`, dengan best practice version control, reversible, testing, dan backup.
- `scripts/README.md` menyebut kebutuhan PostgreSQL client tools `pg_dump`, `pg_restore`, `psql`, `createdb`, `dropdb`, dan catatan Citus extension.

## Project Patterns Found

- Tidak ada satu migration runner root untuk seluruh database; migration tersebar per service/module.
- Banyak migration manual SQL di folder service, tidak semuanya mengikuti `golang-migrate` pair `.up/.down`.
- Object DB memakai schema-prefixed names seperti `mst.`, `inv.`, `acf.`, `sys.`, `sls.`, `pjp`, `pjp_principles`, `report`, dan lain-lain.
- SQL migration lokal sering memakai idempotency parsial seperti `CREATE TABLE IF NOT EXISTS`, `CREATE INDEX IF NOT EXISTS`, `ALTER TABLE ... ADD COLUMN IF NOT EXISTS`, tetapi tidak konsisten di semua file.
- `pjp` dan `tms` punya workflow `make migrateUp`; service lain tampak menyimpan SQL manual di `migration/` tanpa Makefile migration runner yang seragam.

## Reuse Candidates

- Reuse existing service migration folders untuk menempatkan file migration sesuai ownership object.
- Reuse `golang-migrate` workflow untuk `pjp`, `pjp-principle`, dan `tms` bila perubahan menyentuh schema mereka.
- Reuse PostgreSQL tools yang sudah didokumentasikan (`pg_dump`, `psql`, `pg_restore`) untuk schema-only dump, validation, dan backup.
- Reuse existing script documentation sebagai panduan prerequisite, tetapi jangan memakai script clone data untuk pekerjaan ini karena user meminta tanpa data.

## Commands/Docs Checked

- `docker compose -f docker-compose.yml ps` untuk status service lokal.
- `glob **/go.mod` untuk module inventory.
- `glob **/*migrat*`, `glob **/*.sql`, `glob **/Makefile` untuk migration inventory.
- `grep` untuk pola `pg_dump`, `migrate -path`, `DATABASE_URL`, dan SQL DDL.
- `docs/DATABASE.md` dan `scripts/README.md` untuk migration/backup conventions.

## Advisory Sources

- `release-engineer` advisor merekomendasikan schema-only diff, backup demo, dry-run ke clone demo, klasifikasi destructive/risky/non-destructive, Citus-aware review, dan gate sebelum apply.
- `security-privacy-reviewer` advisor merekomendasikan secret handling via env ephemeral, credential terpisah staging read-only vs demo DDL-only, DDL-only guardrails, audit logging, blocklist DML/data movement, dan target DB verification.

## Constraints

- User memberikan credential remote staging dan demo; plan tidak boleh menyimpan credential mentah di artifact final.
- Scope user eksplisit: compare schema/table/function/dll, **kecuali data**; update demo sesuai staging saat ini **kecuali data**; update memakai file migration.
- Repo punya instruksi konflik tentang `rtk`: instruksi global OpenCode melarang prefix `rtk`, sementara `AGENTS.md` lokal meminta `rtk`. Discovery menjalankan command tanpa `rtk` mengikuti global OpenCode yang lebih tinggi.
- Root folder bukan git repo menurut environment, tetapi submodule service kemungkinan punya `.git` masing-masing; plan tidak akan melakukan commit.
- Remote database kemungkinan Citus/PostgreSQL; Citus metadata perlu diperiksa sebelum final SQL generated.

## Risks

- Schema diff tool dapat menghasilkan destructive DDL yang berbahaya untuk demo data.
- Function/view definition bisa membawa logic atau default yang sensitif; diff artifact perlu diperlakukan sebagai sensitive internal artifact.
- Migration lint harus memastikan tidak ada `INSERT`, `UPDATE`, `DELETE`, `COPY`, `CREATE TABLE AS SELECT`, atau cross-database data movement.
- Constraint baru dapat gagal bila data demo melanggar constraint meski schema staging valid.
- `DROP`/rename/type-change harus membutuhkan approval eksplisit atau dipisah sebagai release khusus.
- Tanpa freeze schema staging, diff bisa stale sebelum migration dijalankan.
