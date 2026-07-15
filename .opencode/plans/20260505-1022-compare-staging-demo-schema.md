# Plan — Compare Schema Staging vs Demo dan Update Demo via Migration

Task ID: `20260505-1022-compare-staging-demo-schema`

Status: `ready-for-implementation-with-guardrails`

## Goal

Membuat alur kerja aman untuk:

1. Membandingkan database **staging** dan **demo** pada level schema/object, termasuk schema, table, column, type, default, constraint, index, sequence, view/materialized view, function/procedure, trigger, extension, grant bila relevan, dan metadata Citus bila aktif.
2. Menghasilkan file migration SQL untuk membuat **database demo mengikuti schema staging saat ini**.
3. Memastikan **tidak ada data staging yang dipindahkan ke demo** dan perubahan hanya dilakukan lewat file migration yang direview.

## Non-goals

- Tidak melakukan dump/restore data dari staging ke demo.
- Tidak menjalankan `INSERT`, `UPDATE`, `DELETE`, `COPY`, `TRUNCATE`, `CREATE TABLE AS SELECT`, `SELECT INTO`, atau cross-database data movement dari staging ke demo.
- Tidak menjalankan destructive DDL seperti `DROP COLUMN`, `DROP TABLE`, `DROP FUNCTION`, atau operasi drop/rename lain yang menghapus object demo.
- Tidak langsung apply migration ke demo sebelum ada hasil diff, review manual, backup demo, dan dry-run.
- Tidak menyimpan credential staging/demo di repo, artifact plan, migration, log, atau command history.
- Tidak melakukan scan/mapping detail ke masing-masing repo/service; plan fokus langsung ke database.

## Scope

### Included

- Schema-only inventory dan diff dari staging sebagai source-of-truth terhadap demo sebagai target.
- Object PostgreSQL/Citus yang perlu dicakup:
  - schemas
  - tables dan columns
  - data types, nullable, defaults
  - primary keys, unique constraints, check constraints, foreign keys
  - indexes, termasuk partial/expression indexes
  - sequences dan sequence ownership
  - enum/custom types
  - views/materialized views
  - functions/procedures
  - triggers
  - extensions
  - Citus cluster-level metadata dikecualikan sesuai keputusan user; cukup perlakukan target sebagai PostgreSQL schema/object biasa untuk kebutuhan migration ini.
- Menentukan satu lokasi migration operasional yang disepakati untuk hasil database diff.
- Static guardrail untuk memastikan migration DDL-only.
- Dry-run di clone/temporary demo database sebelum apply remote demo.
- Post-migration schema compare ulang demo vs staging.

### Excluded

- Data row compare antar environment.
- Data sync, backfill data bisnis, seed/reference data, atau correction data demo kecuali disetujui terpisah.
- Perubahan aplikasi Go/API kecuali implementasi diff membuktikan code perlu update untuk schema baru.
- Credential rotation langsung; hanya direncanakan sebagai rekomendasi bila credential sudah dibagikan luas.

## Requirements

1. Gunakan environment variable ephemeral untuk DSN, bukan credential hardcoded.
2. Staging connection idealnya memakai user read-only metadata/schema-only.
3. Demo connection untuk apply migration memakai user DDL-limited, bukan superuser bila memungkinkan.
4. Semua dump/diff harus `schema-only`.
5. Migration harus berbentuk SQL files, tetapi tidak perlu scan/mapping seluruh repo/module; fokus utama adalah database diff dan satu lokasi migration operasional yang disepakati untuk apply ke demo.
6. Migration harus idempotent sejauh aman dan praktis.
7. Risky DDL harus ditandai dan direview; destructive DDL dilarang untuk scope ini.
8. Backup demo wajib sebelum apply ke remote demo.
9. Dry-run wajib di clone/temporary demo database.
10. Hasil akhir dianggap valid hanya bila schema demo setelah migration match staging untuk scope yang disepakati dan data demo tidak berubah secara tidak sengaja.

## Acceptance Criteria

- Ada schema-only dump staging dan demo yang disimpan di lokasi temporary aman di luar repo atau di artifact internal yang tidak dicommit.
- Ada laporan diff yang mengelompokkan perubahan per schema/object dan per risk class: non-destructive, risky, dan unsupported/destructive.
- Ada daftar migration files yang akan dibuat/diubah, dengan scope database object yang jelas; tidak perlu mapping detail ke masing-masing repo/service.
- SQL migration tidak mengandung DML/data movement yang diblokir.
- SQL migration sudah direview untuk function/view/trigger dan dipastikan tidak mengandung destructive DDL.
- Migration berhasil di-apply ke clone/temporary demo database.
- Schema clone-demo setelah migration dibandingkan ulang terhadap staging dan gap yang tersisa terdokumentasi.
- Backup remote demo tersedia sebelum apply final.
- Post-apply demo validation menunjukkan schema demo match staging untuk scope yang disepakati.
- Database smoke checks pass; service tests/smoke checks hanya bila dijalankan opsional.

## Existing Patterns/Reuse

- Repo adalah monorepo Go multi-module tanpa root `go.mod`; module ditemukan di `inventory`, `master`, `pjp`, `finance`, `system`, `sales`, `pjp-principle`, `tms`, `pjp-sales`, `cronjob`, dan `mobile`.
- Migration SQL tersebar:
  - `*/migration/**.sql` untuk service seperti `master`, `mobile`, `finance`, `inventory`, `sales`, dan `pjp-sales`.
  - `pjp/database/migrate/*.sql`.
  - `pjp-principle/database/migrate/*.sql`.
  - `tms/migrations/*.sql`.
- Workflow repo seperti `golang-migrate` di beberapa service diketahui ada, tetapi tidak menjadi fokus implementasi karena user meminta langsung ke database.
- Reuse manual SQL migration conventions untuk service lain yang hanya punya `migration/` folder.
- Reuse PostgreSQL client tools yang sudah disebut di `scripts/README.md`: `pg_dump`, `pg_restore`, `psql`, `createdb`, `dropdb`.
- Jangan memakai script clone data untuk pekerjaan ini karena requirement mengecualikan data.
- User mengonfirmasi tidak perlu scan masing-masing repo/service; reuse repo migration patterns cukup sebagai referensi awal, tetapi implementasi fokus langsung ke database.

## Constraints

- Credential remote sudah diberikan oleh user di chat, tetapi **tidak boleh disalin ke file plan, migration, command output yang disimpan, atau commit**.
- Repo local instructions dan global OpenCode instructions konflik soal prefix `rtk`; untuk eksekusi OpenCode ikuti instruksi global yang lebih tinggi: command tidak memakai `rtk`. Jika operator internal memakai standar repo di luar OpenCode, mereka dapat menyesuaikan.
- Database memakai PostgreSQL/Citus naming dan schema prefix penting: `inv.`, `mst.`, `acf.`, `sys.`, `smc.`, `report.`, `sls.`, `pjp`, `pjp_principles`, dan schema lain yang ditemukan dari metadata.
- Demo mungkin punya data yang tidak sama dengan staging; constraint baru dapat gagal karena data demo, sehingga perlu precheck constraint compatibility tanpa membaca/memindahkan data staging.
- User mengonfirmasi tidak perlu sampai Citus cluster-level checks; plan tetap menjaga DDL PostgreSQL aman tanpa merencanakan metadata/distribution changes Citus.

## Risks

- Diff otomatis dapat menghasilkan `DROP`/rename/type-change yang merusak data demo; output seperti itu harus ditolak atau dicatat sebagai unsupported drift, bukan dieksekusi.
- Constraint baru seperti `NOT NULL`, `UNIQUE`, `FOREIGN KEY`, atau tighter `CHECK` bisa gagal pada data demo.
- `CREATE INDEX CONCURRENTLY` tidak bisa berjalan di dalam transaction block; migration runner perlu disesuaikan bila memakai statement tersebut.
- Function/view/trigger body bisa membawa logic sensitif atau bergantung pada object/data yang tidak ada di demo.
- Schema diff yang dibuat saat staging masih berubah bisa stale sebelum apply.
- Citus cluster-level drift tidak menjadi target pekerjaan ini sesuai keputusan user, sehingga sisa perbedaan metadata Citus tidak dianggap blocker.
- Plaintext credential sudah ada di beberapa repo files menurut instruksi lokal; jangan menambah atau mengekspos secret baru.

## Decisions/Assumptions

### Decisions

- Staging adalah source-of-truth schema; demo adalah target schema.
- Data tidak dibandingkan dan tidak dimigrasikan.
- Migration dibuat incremental dari hasil schema diff, bukan restore schema dump mentah ke demo.
- Perubahan destructive tidak dibuat dan tidak dieksekusi dalam scope ini; catat sebagai unsupported drift bila ditemukan.
- Hasil diff perlu direview manual sebelum menjadi migration final.
- Tidak ada schema/object aplikasi yang dikecualikan dari sync, selain schema system PostgreSQL dan Citus cluster-level metadata.
- Migration tidak boleh mengandung query `DELETE`; selain itu seluruh DML/data movement tetap diblokir untuk menjaga requirement tanpa data migration.
- Migration tidak boleh mengandung destructive DDL, termasuk `DROP COLUMN`, `DROP TABLE`, `DROP FUNCTION`, atau drop/rename object lain yang menghapus/menghilangkan object demo.
- Tidak perlu melakukan Citus cluster-level comparison/update.
- Tidak perlu scan masing-masing repo/service; implementasi langsung fokus ke database schema diff dan migration SQL.

### Assumptions / Open Questions

- Scope database adalah seluruh schema aplikasi pada database yang diberikan, kecuali schema system PostgreSQL (`pg_catalog`, `information_schema`, `pg_toast`) dan Citus cluster-level metadata yang tidak menjadi target pekerjaan ini.
- Diasumsikan demo boleh diubah agar mengikuti staging, tetapi data demo harus dipertahankan.
- Terjawab oleh user: tidak ada schema/object aplikasi yang dikecualikan.
- Terjawab oleh user: jangan ada query `DELETE`. Untuk keamanan, plan juga tetap melarang `INSERT`, `UPDATE`, `TRUNCATE`, `COPY`, `MERGE`, `CREATE TABLE AS SELECT`, dan `SELECT INTO`.
- Terjawab oleh user: tidak perlu sampai Citus cluster.
- Terjawab oleh user: destructive DDL seperti `DROP COLUMN`, `DROP TABLE`, dan `DROP FUNCTION` tidak boleh digunakan. Jika demo punya object tambahan yang tidak ada di staging, gap tersebut dicatat sebagai intentional/unsupported residual drift dan tidak diselesaikan lewat migration ini.
- Terjawab oleh user: tidak perlu scan masing-masing repo/service, langsung fokus ke database.

## TDD/Test Plan

### Apakah TDD wajib?

TDD aplikasi tidak wajib untuk fase compare/migration schema karena perubahan utama adalah DDL dan operasi release database. Namun validation berbasis Red → Green tetap wajib untuk migration safety.

### Reason

Tidak ada production logic Go yang diminta untuk diubah. Risiko utama ada pada schema diff, SQL migration, data preservation, dan Citus/release safety.

### Existing Test Patterns

- Module/service tests tidak menjadi validation utama untuk scope ini; validasi utama adalah schema diff, migration dry-run, dan database smoke checks.
- Service Fiber punya health `GET /ping`; `pjp` Gin-based perlu verifikasi endpoint health aktual.

### Red Step

1. Ambil schema-only dump staging dan demo.
2. Jalankan schema compare untuk membuktikan ada gap demo vs staging.
3. Simpan laporan gap awal yang menjadi baseline Red.

### First Failing/Regression Test

- Schema diff awal harus menunjukkan object/definition demo yang belum sama dengan staging.
- Static SQL scan harus fail bila migration mengandung DML/data movement terlarang.

### Green Step

1. Apply migration ke clone/temporary demo.
2. Dump schema clone-demo setelah migration.
3. Compare ulang clone-demo vs staging.
4. Hasil Green bila gap untuk scope yang disepakati hilang atau sisa gap terdokumentasi sebagai intentionally excluded.

### Refactor Step

- Rapikan migration menjadi per scope database, idempotent, dan mudah rollback/restore.
- Pisahkan risky non-destructive DDL ke file terpisah bila perlu; destructive DDL tetap dilarang dan hanya dicatat sebagai residual drift.
- Tambahkan komentar SQL singkat hanya untuk operasi non-obvious.

### Edge Cases

- `ALTER COLUMN SET NOT NULL` gagal karena data demo memiliki null.
- `ADD UNIQUE` gagal karena duplicate demo data.
- `ADD FOREIGN KEY` gagal karena orphan rows demo.
- `CREATE INDEX CONCURRENTLY` gagal bila runner membungkus transaction.
- Function bergantung extension/schema yang belum ada.
- Sequence ownership/default berbeda meski table/column sama.
- Citus cluster-level metadata tidak divalidasi dalam scope ini.

### Commands

Gunakan env variable ephemeral; contoh nama saja, tanpa credential:

```bash
export STAGING_DATABASE_URL='postgres://<user>:<password>@<staging-host>:<port>/<db>?sslmode=disable'
export DEMO_DATABASE_URL='postgres://<user>:<password>@<demo-host>:<port>/<db>?sslmode=disable'
```

Schema-only dump:

```bash
pg_dump --schema-only --no-owner --no-privileges "$STAGING_DATABASE_URL" > /tmp/scylla_staging_schema.sql
pg_dump --schema-only --no-owner --no-privileges "$DEMO_DATABASE_URL" > /tmp/scylla_demo_schema.sql
diff -u /tmp/scylla_demo_schema.sql /tmp/scylla_staging_schema.sql > /tmp/scylla_schema.diff || true
```

Object inventory checks:

```bash
psql "$STAGING_DATABASE_URL" -v ON_ERROR_STOP=1 -c "SELECT nspname FROM pg_namespace WHERE nspname NOT LIKE 'pg_%' AND nspname <> 'information_schema' ORDER BY 1;"
psql "$DEMO_DATABASE_URL" -v ON_ERROR_STOP=1 -c "SELECT nspname FROM pg_namespace WHERE nspname NOT LIKE 'pg_%' AND nspname <> 'information_schema' ORDER BY 1;"
```

Migration DML guardrail example:

```bash
grep -RInE '\b(INSERT|UPDATE|DELETE|TRUNCATE|COPY|MERGE)\b|CREATE[[:space:]]+TABLE[[:space:]].*AS[[:space:]]+SELECT|SELECT[[:space:]].*INTO' <migration-files>
```

Service tests per affected module:

```bash
go test ./...
```

## Implementation Steps

### Phase 0 — Preparation and Guardrails

1. Buat branch kerja khusus untuk migration schema sync.
2. Pastikan `pg_dump`, `psql`, `pg_restore`, dan schema diff tool tersedia. Pilih tool utama:
   - baseline wajib: `pg_dump --schema-only` + `diff` + catalog queries;
   - lebih baik untuk SQL generation: `migra`, `apgdiff`, `psqldef`, atau Atlas, tetapi hasilnya tetap manual review.
3. Set DSN staging/demo via env variable ephemeral.
4. Verifikasi target connection tanpa mencetak password:
   - `SELECT current_database(), current_user, inet_server_addr(), inet_server_port();`
   - cek marker hostname/db untuk memastikan staging dan demo tidak tertukar.
5. Tentukan schema allowlist dari metadata aplikasi, lalu exclude schema system.

### Phase 1 — Schema Inventory and Baseline Diff

1. Ambil schema-only dump staging dan demo dengan `--schema-only --no-owner --no-privileges`.
2. Ambil catalog inventory terstruktur untuk:
   - table/column/type/default/nullability;
   - constraints;
   - indexes;
   - sequences;
   - views/materialized views;
   - functions/procedures/triggers;
   - extensions;
   - Citus metadata bila extension aktif.
3. Normalize output agar diff stabil, misalnya sort object inventory query by schema/name.
4. Generate diff awal.
5. Klasifikasikan gap:
   - non-destructive: `CREATE SCHEMA`, `CREATE TABLE`, nullable `ADD COLUMN`, `CREATE INDEX`, `CREATE FUNCTION`, `CREATE VIEW`.
   - risky: `ALTER TYPE`, `SET NOT NULL`, `ADD UNIQUE`, `ADD FK`, default change, function replacement besar.
   - unsupported/destructive: `DROP`, rename yang menghilangkan object lama, incompatible type change, `DROP COLUMN`, `DROP TABLE`, `DROP FUNCTION`; jangan buat migration untuk kategori ini.

### Phase 2 — Tentukan Lokasi Migration Operasional

1. Tidak perlu scan atau mapping masing-masing repo/service.
2. Gunakan hasil database diff sebagai source utama untuk menyusun migration.
3. Tempatkan migration di satu lokasi operasional yang disepakati untuk pekerjaan database sync ini, misalnya folder migration service yang diminta operator atau folder migration khusus database sync.
4. Jika tim sudah punya lokasi standar untuk migration database lintas schema, gunakan lokasi tersebut.
5. Jika belum ada lokasi standar, buat satu migration SQL terpisah berdasarkan scope/risk database, bukan berdasarkan ownership repo.

### Phase 3 — Draft Migration SQL

1. Buat migration SQL berdasarkan hasil diff database dengan timestamp baru.
2. Pecah file hanya bila dibutuhkan oleh risiko/urutan dependency, misalnya `001_create_missing_objects.sql`, `002_alter_existing_objects.sql`, dan `003_replace_functions_views.sql`.
3. Jangan membuat file per service hanya untuk mengikuti struktur repo jika user meminta fokus database saja.
4. Gunakan idempotency aman:
   - `CREATE SCHEMA IF NOT EXISTS`;
   - `CREATE TABLE IF NOT EXISTS`;
   - `ALTER TABLE ... ADD COLUMN IF NOT EXISTS`;
   - `CREATE INDEX IF NOT EXISTS` bila tidak memakai concurrently;
   - `CREATE OR REPLACE FUNCTION` hanya setelah review body.
5. Untuk `CREATE INDEX CONCURRENTLY`, pisahkan dari transaction-based migration jika runner tidak mendukung.
6. Untuk risky constraints, tambahkan precheck query terhadap demo data sebelum apply final.
7. Jangan memasukkan DML/data movement.
8. Jangan memasukkan destructive DDL. Object demo tambahan yang tidak ada di staging hanya dilaporkan sebagai residual drift.

### Phase 4 — Static Review and Safety Checks

1. Scan migration untuk blocklist DML/data movement.
2. Tolak setiap `DROP`, rename yang menghapus object lama, incompatible `ALTER TYPE`, `DROP COLUMN`, `DROP TABLE`, `DROP FUNCTION`, atau destructive DDL lain dari migration generated.
3. Review constraint tightening dan default/type changes yang masih non-destructive untuk risiko terhadap data demo.
4. Review function/view/trigger body untuk dependency dan security concern.
5. Pastikan migration tidak mencetak atau menyimpan DSN.
6. Pastikan tidak ada Citus cluster-level DDL seperti distribution/colocation metadata changes karena user menyatakan tidak perlu sampai Citus cluster.

### Phase 5 — Dry-run di Clone/Temporary Demo

1. Buat backup demo sebelum dry-run.
2. Restore backup demo ke database temporary/clone yang isolated.
3. Apply migration ke temporary database langsung via `psql -v ON_ERROR_STOP=1 -f <file>` sesuai urutan dependency.
4. Jangan jalankan scan/test per repo/service kecuali setelah database migration terbukti memerlukan perubahan aplikasi.
5. Dump schema temporary setelah migration.
6. Compare temporary schema vs staging schema.
7. Catat sisa diff yang legitimate atau intentionally excluded.
8. Jalankan smoke check database dan schema compare ulang; module/service tests bersifat opsional dan hanya bila ada indikasi perubahan aplikasi terdampak.

### Phase 6 — Approval Gate

Migration ke remote demo hanya boleh lanjut bila:

- backup/restore test sukses;
- dry-run sukses;
- schema temporary match staging untuk scope yang disepakati;
- static DML scan bersih;
- risky non-destructive DDL sudah direview, tidak ada query `DELETE`, dan tidak ada destructive DDL;
- rollback tersedia;
- maintenance window/owner approval tersedia bila perubahan berisiko.

### Phase 7 — Apply ke Demo

1. Freeze schema staging sementara atau set cut-off timestamp agar diff tidak berubah.
2. Ambil backup remote demo tepat sebelum apply final.
3. Verifikasi target demo connection lagi.
4. Apply migration batch per file/scope database, mulai dari dependencies schema/type/function yang dibutuhkan table lain.
5. Monitor migration output dan DB/app logs.
6. Stop jika ada error; jangan lanjut batch berikutnya sampai root cause jelas.

### Phase 8 — Post-validation

1. Dump schema demo setelah apply.
2. Compare demo-after vs staging schema.
3. Jalankan catalog checks untuk object penting.
4. Jalankan database smoke checks; service tests/smoke checks opsional bila operator ingin validasi aplikasi tambahan.
5. Verifikasi row count demo tidak berubah untuk table yang seharusnya DDL-only bila memungkinkan:
   - ambil count sample sebelum dan sesudah untuk table critical demo;
   - tidak membandingkan dengan staging.
6. Dokumentasikan hasil apply, migration IDs, checksum, duration, sisa diff, dan rollback notes.

## Expected Files to Change

Source code Go tidak diharapkan berubah kecuali ditemukan incompatibility setelah migration review.

Expected migration file location:

- Satu lokasi operasional yang disepakati untuk database sync, misalnya `<chosen-migration-folder>/<timestamp>_sync_demo_schema_to_staging.sql`.
- Jika perlu dipisah karena urutan/risk, gunakan beberapa file di lokasi yang sama, misalnya:
  - `<chosen-migration-folder>/<timestamp>_001_create_missing_objects.sql`
  - `<chosen-migration-folder>/<timestamp>_002_alter_existing_objects.sql`
  - `<chosen-migration-folder>/<timestamp>_003_replace_functions_views.sql`

Tidak perlu membuat/memilih file berdasarkan repo/service ownership.

Optional operational docs/artifacts yang boleh dibuat saat implementasi bila diperlukan:

- sanitized schema diff summary di `.opencode/evidence/<new-implementation-task-id>/` atau internal release artifact yang tidak dicommit.
- rollback SQL notes untuk service manual migration.

## Agent/Tool Routing

- `@explorer`: tidak wajib untuk fase implementasi ini karena user meminta fokus langsung ke database, bukan scan masing-masing repo.
- `@release-engineer`: review release plan, backup, dry-run, migration ordering, rollback, dan post-validation.
- `@security-privacy-reviewer`: review secret handling, DDL-only guardrail, least privilege, audit, dan data exposure risks.
- `@oracle`: review kompleks bila diff menghasilkan banyak unsupported/destructive drift, risky non-destructive DDL, atau tradeoff schema tidak jelas.
- `@fixer`: hanya setelah plan ini, untuk membuat migration files dan menjalankan validation secara bounded.
- `@quality-gate`: final conformance review sebelum commit/apply ke demo.

## Validation Commands

Gunakan command tanpa credential literal. Set DSN lewat env var ephemeral.

### Connectivity and Target Verification

```bash
psql "$STAGING_DATABASE_URL" -v ON_ERROR_STOP=1 -c "SELECT current_database(), current_user, inet_server_addr(), inet_server_port();"
psql "$DEMO_DATABASE_URL" -v ON_ERROR_STOP=1 -c "SELECT current_database(), current_user, inet_server_addr(), inet_server_port();"
```

### Schema-only Dumps

```bash
pg_dump --schema-only --no-owner --no-privileges "$STAGING_DATABASE_URL" > /tmp/scylla_staging_schema.sql
pg_dump --schema-only --no-owner --no-privileges "$DEMO_DATABASE_URL" > /tmp/scylla_demo_schema.sql
diff -u /tmp/scylla_demo_schema.sql /tmp/scylla_staging_schema.sql > /tmp/scylla_schema.diff || true
```

### Demo Backup Before Final Apply

```bash
pg_dump -Fc --no-owner --no-privileges "$DEMO_DATABASE_URL" > /tmp/scylla_demo_before_$(date +%Y%m%d_%H%M%S).dump
```

### DML/Data Movement Guardrail

```bash
grep -RInE '\b(INSERT|UPDATE|DELETE|TRUNCATE|COPY|MERGE)\b|CREATE[[:space:]]+TABLE[[:space:]].*AS[[:space:]]+SELECT|SELECT[[:space:]].*INTO' <migration-files>
```

Expected result: no matches. Any match must be reviewed and justified; for this task it should usually block.

### Destructive DDL Guardrail

```bash
grep -RInE '\bDROP[[:space:]]+(COLUMN|TABLE|FUNCTION|SCHEMA|VIEW|MATERIALIZED[[:space:]]+VIEW|SEQUENCE|TYPE|TRIGGER|INDEX)\b|\bALTER[[:space:]]+TABLE\b.*\bDROP\b|\bDROP\b|\bALTER[[:space:]]+(TABLE|FUNCTION|VIEW|TYPE)\b.*\bRENAME\b' <migration-files>
```

Expected result: no matches. Any match blocks this migration scope and must be reported as residual/unsupported drift.

### Citus Metadata Check

Tidak wajib untuk scope ini karena user menyatakan tidak perlu sampai Citus cluster. Jika implementer tetap ingin memastikan extension tidak mengganggu migration, cukup cek keberadaan extension tanpa menjadikan metadata Citus sebagai acceptance blocker.

```bash
psql "$STAGING_DATABASE_URL" -v ON_ERROR_STOP=1 -c "SELECT * FROM pg_extension WHERE extname = 'citus';"
psql "$DEMO_DATABASE_URL" -v ON_ERROR_STOP=1 -c "SELECT * FROM pg_extension WHERE extname = 'citus';"
```

### Apply Migration Examples

Apply SQL migration langsung ke temporary/demo database sesuai urutan file:

```bash
psql "$TEMP_DEMO_DATABASE_URL" -v ON_ERROR_STOP=1 -f <migration-file.sql>
```

### Module Tests

Opsional. Tidak perlu scan/test masing-masing repo untuk scope plan ini kecuali database migration menunjukkan aplikasi terdampak. Jika dibutuhkan, run dari module terkait:

```bash
go test ./...
```

### Health/Smoke Checks

When services are running:

```bash
curl -f http://localhost:9001/ping
curl -f http://localhost:9002/ping
curl -f http://localhost:9003/ping
curl -f http://localhost:9004/ping
curl -f http://localhost:9005/ping
curl -f http://localhost:9006/ping
curl -f http://localhost:9008/ping
```

For `pjp` on port `9010`, verify actual Gin route before assuming `/ping`.

## Evidence Requirements

Implementation must produce or retain sanitized evidence for:

- schema allowlist/exclusion list;
- staging schema dump command metadata, not credentials;
- demo schema dump command metadata, not credentials;
- diff summary per object/risk class, including unsupported/destructive residual drift that will not be migrated;
- migration file list and checksum;
- DML guardrail scan result;
- Citus extension presence check only if implementer wants additional context; Citus cluster-level metadata comparison is not required.
- dry-run apply result on clone/temporary demo;
- schema compare result after dry-run;
- backup path/checksum for demo before final apply;
- final post-apply schema compare;
- database smoke results; service test/smoke results hanya bila dijalankan opsional.

Evidence should avoid raw credential and avoid committing full schema dumps unless team explicitly approves because function/view definitions may expose internal logic.

## Done Criteria

- Plan ini digunakan sebagai source of truth untuk implementation.
- Hasil schema diff sudah tersedia dan direview.
- Migration SQL dibuat di lokasi operasional yang disepakati untuk database sync.
- Migration lulus static DML guardrail.
- Migration lulus destructive DDL guardrail.
- Migration lulus dry-run di clone/temporary demo.
- Backup demo tersedia sebelum apply final.
- Demo berhasil di-update via migration file, bukan direct manual ad-hoc DDL tanpa file.
- Post-migration schema demo match staging untuk non-destructive scope yang disepakati; residual demo-only objects may remain documented because destructive DDL is prohibited.
- Tidak ada data staging dipindahkan ke demo.
- Rollback/restore path terdokumentasi.

## Final Planning Summary

- Artifact utama dibuat: `.opencode/plans/20260505-1022-compare-staging-demo-schema.md`.
- Evidence discovery dibuat dan dipertahankan: `.opencode/evidence/20260505-1022-compare-staging-demo-schema/discovery.md` sebagai catatan konteks awal dan advisory release/security; implementasi tetap fokus database, bukan scan repo/service.
- Draft artifact tidak dibuat karena semua durable findings langsung dikonsolidasikan ke plan utama.
- Key decisions: schema-only diff, no data migration, incremental migration files, dry-run before demo, no destructive DDL, credentials via env vars only.
- Pertanyaan yang sudah dijawab user: tidak ada schema/object aplikasi yang dikecualikan; tidak boleh ada query `DELETE`; tidak boleh ada destructive DDL seperti `DROP COLUMN`, `DROP TABLE`, atau `DROP FUNCTION`; tidak perlu sampai Citus cluster; tidak perlu scan masing-masing repo/service, fokus langsung ke database.
- Pertanyaan belum terjawab: maintenance/backup procedure resmi demo.
- Readiness: siap untuk fase implementation discovery/diff oleh `@fixer` atau engineer, tetapi belum siap apply ke demo sampai guardrails, diff review, backup, dry-run, dan review risky non-destructive DDL terpenuhi.
- Cleanup performed: tidak ada draft stale; evidence discovery dipertahankan secara sengaja.
