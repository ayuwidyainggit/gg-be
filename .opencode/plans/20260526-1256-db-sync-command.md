# Plan: Shell Command Sync Remote DB ke Local

Task ID: `20260526-1256-db-sync-command`

## Goal

Buat shell command/script aman untuk sync database remote `scylla_citus_dev` di `103.28.219.73:25431` ke database local `ggn_scyllax`.

Plan ini mengutamakan reuse pola script repo dan mencegah secret tertulis ulang di file/script/artifact.

## Non-goals

- Tidak menjalankan sync database saat planning.
- Tidak mengubah data remote.
- Tidak membuat migration aplikasi.
- Tidak mengubah source Go service.
- Tidak menyimpan password remote/local di repo.
- Tidak menjamin restore Citus berhasil jika local PostgreSQL belum punya extension Citus kompatibel.

## Scope

Masuk scope:

- Shell script non-interaktif/opsional interaktif untuk sync remote → local.
- Default source: `scylla_citus_dev` di `103.28.219.73:25431` user `postgres`.
- Default target: `ggn_scyllax` di `localhost:5432` user `postgres`.
- Mode dump custom format: `pg_dump -Fc` lalu `pg_restore`.
- Opsi aman untuk drop/recreate target local.
- Cek prerequisite PostgreSQL client tools.
- Cek koneksi source dan target.
- Cek/beri instruksi untuk Citus extension.
- Validasi hasil restore minimal schema/table count.
- Dokumentasi usage singkat di `scripts/README.md` bila implementasi disetujui.

Keluar scope:

- Backup incremental/differential.
- Sinkronisasi dua arah.
- Masking/anonymization data.
- Scheduler/cron.
- Restore selective schema/table kecuali diminta lanjutan.

## Requirements

- Script harus bisa dipanggil dari repo root.
- Script harus gagal cepat bila `pg_dump`, `pg_restore`, `psql`, `createdb`, atau `dropdb` tidak ada.
- Source credential harus lewat environment variable, contoh `SOURCE_DB_PASSWORD`, bukan hardcoded.
- Local credential harus lewat env automation: `LOCAL_DB_PASSWORD`, `LOCAL_DB_HOST`, `LOCAL_DB_PORT`, `LOCAL_DB_USER`, `LOCAL_DB_NAME`.
- Default local db name harus `ggn_scyllax`.
- Local target boleh drop/recreate karena user sudah menyetujui, tapi script tetap harus membatasi destructive action hanya ke `LOCAL_DB_NAME=ggn_scyllax` default atau target eksplisit.
- Mode utama harus automation/non-interaktif, contoh `--drop --yes --install-citus`.
- Script harus memasang/menyiapkan Citus local bila missing melalui helper existing `scripts/install_citus.sh` atau instruksi runnable yang aman.
- Remote database wajib read-only dari sisi command: hanya `pg_dump` dan `psql` SELECT/preflight; tidak boleh ada `psql` DDL/DML ke remote.
- Dump file default disimpan di path temp atau `./tmp/` yang tidak ikut commit; bila membuat folder repo, harus masuk `.gitignore` dulu atau pakai `/tmp`.
- Restore harus memakai `--clean --if-exists --no-owner --no-privileges`.
- Script harus menampilkan source/target tanpa password.
- Script harus memberi pesan jelas untuk Citus jika install otomatis gagal.
- Script tidak boleh memfilter error restore sampai exit code hilang; error harus memengaruhi status akhir.

## Acceptance Criteria

- Ada shell command/script siap pakai untuk sync remote `scylla_citus_dev` ke local `ggn_scyllax`.
- Password remote tidak muncul di file hasil implementasi, plan, log contoh, atau command docs final.
- `bash -n scripts/sync_remote_to_local.sh` sukses.
- Script dengan env kurang lengkap gagal dengan pesan jelas.
- Script dapat melakukan dry validation koneksi source dan target tanpa restore penuh.
- Script dapat membuat database local bila belum ada.
- Script hanya drop/recreate `ggn_scyllax` jika user memberi flag eksplisit atau menjawab konfirmasi.
- Restore memakai custom dump + `pg_restore`, bukan pipe raw yang mencampur stderr/stdout.
- Validasi pasca-restore menampilkan minimal jumlah schema/table dan `current_database()`.
- README scripts memuat usage tanpa password literal.

## Existing Patterns/Reuse

- Reuse pola dari `scripts/clone_db.sh`: color output, requirement checks, env loading, connection test, create local db.
- Reuse pola dari `scripts/restore_db.sh`: support `pg_restore --clean --if-exists --no-owner --no-acl` dan post-restore stats.
- Reuse pola dari `scripts/clone_staging.sh`: Citus availability check dan prompt safety.
- Reuse `scripts/install_citus.sh` sebagai instruksi saat local Citus belum tersedia.
- Jangan reuse direct pipe branch `clone_db.sh` sebagai default karena stdout/stderr filtering bisa merusak SQL stream dan menutupi exit code.

Tidak ada KiloCode/project utility lain yang lebih tepat daripada script existing di `scripts/` untuk kebutuhan ini.

## Constraints

- Remote password sudah diberikan user di chat, tapi plan/implementasi harus memperlakukannya sebagai secret dan tidak menuliskannya lagi.
- Repo sudah punya plaintext credential di `docker-compose.yml`; jangan memperluas exposure.
- Local `ggn_scyllax` mungkin berisi data penting; drop wajib eksplisit.
- Source bernama `scylla_citus_dev`; local mungkin butuh `citus` dan `citus_columnar`.
- Citus restore penuh bisa butuh local PostgreSQL versi/extension kompatibel.
- Runtime repo memakai `rtk docker compose -f docker-compose.yml ps` sebagai check awal bila perlu validasi environment.

## Risks

- Data local hilang jika `--drop` salah target.
- Remote dump besar memakan waktu dan disk.
- Restore gagal bila local role/extension/schema tidak cocok.
- Distributed Citus metadata bisa gagal restore di single-node local tanpa setup Citus.
- Secret bocor ke shell history jika user menjalankan inline `SOURCE_DB_PASSWORD='...' command`.
- Existing scripts punya pola filter error Citus; implementasi baru harus menjaga exit code valid.

Mitigasi:

- Require `--yes` bersama `--drop` untuk non-interaktif destructive mode, atau prompt typed confirmation `ggn_scyllax`.
- Gunakan env exported dari shell session atau `.pgpass` lokal, bukan inline command di history.
- Print password redacted saja.
- Simpan dump di `/tmp/scylla_citus_dev_<timestamp>.dump` default.
- Jalankan preflight Citus dan tampilkan instruksi `./scripts/install_citus.sh`.

## Decisions/Assumptions

- Interaction level: assumption-first dengan blocker kecil untuk credential dan destructive behavior.
- Keputusan teknis: buat script baru `scripts/sync_remote_to_local.sh`, bukan memodifikasi `clone_db.sh`, supaya flow khusus `scylla_citus_dev` → `ggn_scyllax` jelas dan tidak mengubah perilaku script lama.
- Source default:
  - `SOURCE_DB_HOST=103.28.219.73`
  - `SOURCE_DB_PORT=25431`
  - `SOURCE_DB_USER=postgres`
  - `SOURCE_DB_NAME=scylla_citus_dev`
- Target default:
  - `LOCAL_DB_HOST=localhost`
  - `LOCAL_DB_PORT=5432`
  - `LOCAL_DB_USER=postgres`
  - `LOCAL_DB_NAME=ggn_scyllax`
- Password source/local wajib dari env; automation tidak boleh meminta prompt.
- Local PostgreSQL final dari user: host `localhost`, port `5432`, user `postgres`; password lewat `LOCAL_DB_PASSWORD` saat execution, tidak ditulis di file.
- Target `ggn_scyllax` boleh di-drop/recreate.
- Citus local belum terinstall; implementasi harus memasang/menyiapkan Citus sebelum restore atau stop dengan instruksi jelas bila OS/package manager tidak mendukung auto-install.
- Script mode utama: automation penuh/non-interaktif.
- Remote DB guard: command ke remote hanya boleh `pg_dump` dan `psql` SELECT/preflight; tidak boleh `drop`, `create`, `restore`, DDL, atau DML ke remote.

Open questions untuk implementasi: tidak ada yang memblokir. Credential source tetap harus diberikan out-of-band melalui `SOURCE_DB_PASSWORD` saat execution, bukan ditulis di file.

## TDD/Test Plan

TDD required: ya, karena shell script mengubah database local dan destructive risk tinggi.

Reason:

- Behavior parsing flag, env validation, destructive guard, dan command construction harus aman sebelum restore sungguhan.

Existing test patterns:

- Tidak ditemukan test shell khusus di repo.
- Existing validation repo untuk scripts memakai command manual.
- Gunakan minimal shell validation tanpa dependency baru: `bash -n`, mode `--help`, mode missing env, dan preflight dry-run.

First failing/regression test:

- Jalankan script tanpa `SOURCE_DB_PASSWORD` dalam mode non-interaktif; expected exit non-zero dan pesan `SOURCE_DB_PASSWORD is required` tanpa menampilkan password.

Green step:

- Implement env/flag parser sampai test missing `SOURCE_DB_PASSWORD`, missing `LOCAL_DB_PASSWORD`, `--help`, dan dry-run pass.
- Tambah preflight command dengan real `psql` hanya ketika credential tersedia; remote preflight harus SELECT-only.

Refactor step:

- Pisah fungsi: `print_*`, `require_commands`, `load_config`, `confirm_destructive`, `test_connection`, `check_citus`, `dump_remote`, `prepare_local_db`, `restore_dump`, `validate_restore`, `cleanup_dump`.
- Pastikan setiap fungsi return code dicek.

Edge cases:

- Source password kosong.
- Local password lewat `LOCAL_DB_PASSWORD`, termasuk nilai default lokal yang user punya.
- Local database sudah ada tanpa `--drop`.
- Local database tidak ada.
- `pg_dump` gagal koneksi.
- `pg_restore` mengembalikan error.
- Citus missing.
- Dump file sudah ada.
- User cancel confirmation.
- Disk temp tidak writable.

Commands:

```bash
bash -n scripts/sync_remote_to_local.sh
scripts/sync_remote_to_local.sh --help
env -u SOURCE_DB_PASSWORD scripts/sync_remote_to_local.sh --dry-run
SOURCE_DB_PASSWORD='<from-secure-env>' scripts/sync_remote_to_local.sh --dry-run
```

Integration validation, hanya setelah user set credential:

```bash
rtk docker compose -f docker-compose.yml ps
SOURCE_DB_PASSWORD='<from-secure-env>' scripts/sync_remote_to_local.sh --preflight-only
SOURCE_DB_PASSWORD='<from-secure-env>' scripts/sync_remote_to_local.sh --drop --yes
psql -h localhost -p 5432 -U postgres -d ggn_scyllax -c "select current_database(), current_user;"
```

## Implementation Steps

1. Buat `scripts/sync_remote_to_local.sh` dengan strict mode `set -Eeuo pipefail`.
2. Tambah `usage()` yang menjelaskan env dan flags tanpa password literal.
3. Implement defaults source/target dan env override.
4. Implement command requirement check untuk `pg_dump`, `pg_restore`, `psql`, `createdb`, `dropdb`.
5. Implement secret input fallback via `read -rsp` bila `SOURCE_DB_PASSWORD` kosong dan mode interaktif.
6. Implement local password optional via `LOCAL_DB_PASSWORD`; jika kosong biarkan client auth default.
7. Implement `--dry-run` / `--preflight-only` untuk validasi tanpa dump/restore.
8. Implement source connection test dengan SELECT-only: `select current_database(), current_user, version();`.
9. Implement local connection test ke db `postgres` memakai `LOCAL_DB_PASSWORD`.
10. Implement Citus install/check:
    - Query source `pg_extension` untuk `citus` dan `citus_columnar` memakai SELECT-only.
    - Query local `pg_available_extensions` dan target `pg_extension` bila db ada.
    - Bila source pakai Citus dan local belum siap, jalankan helper `scripts/install_citus.sh` dalam mode automation bila memungkinkan; jika helper masih interaktif, implementor boleh extend helper dengan flag non-interaktif atau menambahkan instruksi install non-interaktif di script baru.
    - Jika auto-install gagal, stop sebelum restore dengan pesan jelas.
11. Implement destructive guard automation:
    - Jika target exists dan `--drop --yes`, terminate local connections lalu `dropdb` + `createdb` hanya pada local host/port.
    - Jika target exists tanpa `--drop --yes`, fail fast agar automation tidak ambigu.
    - Validasi target destructive: tolak jika `LOCAL_DB_NAME` kosong, `postgres`, `template0`, `template1`, atau bukan target yang user set eksplisit.
12. Implement dump remote ke custom format file:
    - Default `/tmp/scylla_citus_dev_YYYYMMDD_HHMMSS.dump`.
    - `pg_dump -Fc --no-owner --no-privileges`.
13. Implement restore:
    - `pg_restore --clean --if-exists --no-owner --no-privileges -d ggn_scyllax <dump>`.
    - Jangan pipe stderr ke parser yang menghilangkan exit code.
14. Implement post-restore validation:
    - `select current_database(), current_user;`
    - schema/table count dari `information_schema.tables` atau `pg_tables` untuk semua schema non-system.
    - extension list.
15. Implement cleanup dump default, dengan flag `--keep-dump` untuk menyimpan file.
16. Update `scripts/README.md` dengan usage, env, examples, Citus note, destructive warning.
17. Validasi syntax, dry-run, preflight, lalu integration restore bila credential dan izin tersedia.
18. Jalankan `@quality-gate` final karena menyentuh ops script, destructive local data, dan secret handling.

## Expected Files to Change

Implementasi nanti kemungkinan mengubah:

- `scripts/sync_remote_to_local.sh` baru.
- `scripts/README.md` update usage.
- Opsional `.gitignore` jika dump dir repo-local dipilih; recommended tidak perlu karena default `/tmp`.

Tidak perlu ubah:

- Source Go service.
- `docker-compose.yml`.
- `.env` files.
- Migration files.
- Lockfiles.

## Agent/Tool Routing

- `@orchestrator`: jalankan handoff dan integrasi tugas.
- `@fixer`: implement shell script + README update + local validations.
- `@explorer`: bila implementor butuh discovery tambahan terhadap existing scripts.
- `@oracle`: opsional review destructive guard/secret posture jika script makin kompleks.
- `@quality-gate`: wajib final review karena destructive local DB + secret handling.
- `@artifact-planner`: selesai setelah plan ini; tidak melakukan implementasi source.

Research gate:

- Local project discovery: dilakukan, wajib untuk reuse existing scripts.
- Official docs/context7: tidak diperlukan; PostgreSQL CLI behavior standar dan repo sudah punya pola.
- GitHub: tidak diperlukan; tidak tergantung upstream repo/source issue.
- Brave/web search: tidak diperlukan; tidak ada fakta eksternal/current yang menentukan plan.
- Browser/screenshot: tidak relevan.

## Execution-ready Worklist / Handoff Contract

`start_with`: `T01`

| Task | depends_on | owner/lane | action | validation/check | exit criteria | status | blocker | requires_user_decision |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T01 | none | @fixer | Inspect existing `scripts/*.sh` and confirm no newer sync helper exists | `ls scripts` plus read target scripts | reuse decision confirmed | ready | none | no |
| T02 | T01 | @fixer | Create `scripts/sync_remote_to_local.sh` with usage, config defaults, strict mode, and print helpers | `bash -n scripts/sync_remote_to_local.sh` | syntax valid and `--help` works | ready | none | no |
| T03 | T02 | @fixer | Add env/flag parser and secret-safe automation validation | `env -u SOURCE_DB_PASSWORD scripts/sync_remote_to_local.sh --dry-run` and `env -u LOCAL_DB_PASSWORD scripts/sync_remote_to_local.sh --dry-run` | missing secrets fail safely and no password printed | ready | none | no |
| T04 | T03 | @fixer | Add command prerequisite checks | temporarily run `PATH=/nonexistent scripts/sync_remote_to_local.sh --dry-run` if safe, or review function | clear missing command errors | ready | none | no |
| T05 | T03 | @fixer | Add source/local connection preflight and `--preflight-only`; remote checks SELECT-only | `SOURCE_DB_PASSWORD='<from-secure-env>' LOCAL_DB_PASSWORD='<from-secure-env>' scripts/sync_remote_to_local.sh --preflight-only` | source/local connection checks pass or clear failure shown; no remote writes possible | blocked | needs source/local password at execution | no |
| T06 | T03 | @fixer | Add Citus detection plus automation install/prepare path | run with `--dry-run`; if local available run extension queries and installer path | local Citus ready or script stops before restore with clear install failure | ready | auto-install may depend on OS/package manager | no |
| T07 | T03 | @fixer | Add safe local DB create/drop/recreate guard for automation | dry-run/review plus throwaway db test if needed | no destructive action unless `--drop --yes`; destructive target cannot be system DB | ready | none | no |
| T08 | T05,T06,T07 | @fixer | Add `pg_dump -Fc` to temp file and `pg_restore` restore flow | `bash -n`; integration only with credentials | commands preserve exit code; dump path safe; remote only uses `pg_dump` | blocked | needs credentials to run full sync | no |
| T09 | T08 | @fixer | Add post-restore validation and dump cleanup/keep flag | real restore validation query | table/schema stats printed; dump removed unless `--keep-dump` | blocked | depends on real restore | no |
| T10 | T02,T09 | @fixer | Update `scripts/README.md` with usage and warnings | read README section | examples contain no literal password | ready | none | no |
| T11 | T10 | @quality-gate | Review destructive guard, secret handling, validation evidence | review changed files and commands output | signoff or required fixes listed | ready | none | no |

## Validation Commands

Syntax and help:

```bash
bash -n scripts/sync_remote_to_local.sh
scripts/sync_remote_to_local.sh --help
```

Secret-safe failure:

```bash
env -u SOURCE_DB_PASSWORD scripts/sync_remote_to_local.sh --dry-run
```

Repo runtime baseline if local services matter:

```bash
rtk docker compose -f docker-compose.yml ps
```

Preflight dengan credential dari env, tanpa menulis password literal di repo:

```bash
export SOURCE_DB_PASSWORD='<from-secure-env>'
export LOCAL_DB_PASSWORD='<from-secure-env>'
scripts/sync_remote_to_local.sh --preflight-only
```

Full local sync automation. Ini hanya boleh menulis ke local `ggn_scyllax`; remote hanya `pg_dump`/SELECT:

```bash
export SOURCE_DB_PASSWORD='<from-secure-env>'
export LOCAL_DB_PASSWORD='<from-secure-env>'
scripts/sync_remote_to_local.sh --install-citus --drop --yes
```

Manual post-check:

```bash
psql -h localhost -p 5432 -U postgres -d ggn_scyllax -c "select current_database(), current_user;"
psql -h localhost -p 5432 -U postgres -d ggn_scyllax -c "select schemaname, count(*) from pg_tables where schemaname not in ('pg_catalog','information_schema') group by schemaname order by schemaname;"
```

## Evidence Requirements

Implementation evidence harus mencatat:

- Changed files.
- `bash -n` result.
- `--help` output summary.
- Missing secret failure output, redacted.
- Preflight connection result, redacted.
- Citus check result.
- If full restore run: dump path, duration, restore exit code, post-restore table/schema counts.
- Confirmation that password tidak ditulis ke repo/log artifact.
- `@quality-gate` result.

Kept planning evidence:

- `.opencode/evidence/20260526-1256-db-sync-command/discovery.md`
- `.opencode/evidence/20260526-1256-db-sync-command/index.json`

## Done Criteria

- `scripts/sync_remote_to_local.sh` exists and executable.
- Script default target is `ggn_scyllax`.
- Source config default matches remote dev DB except password.
- Password only accepted via env for automation.
- Destructive local drop diizinkan untuk `ggn_scyllax` karena user sudah approve, tapi automation tetap harus memakai `--drop --yes` agar tidak salah target.
- `bash -n` passes.
- Dry-run/preflight behaviors pass.
- README usage added without secret.
- Real sync either completed successfully with evidence, or clearly marked blocked by credential/local Citus/local DB access.
- `@quality-gate` signs off or blockers resolved.

## Final Planning Summary

Artifacts created:

- Primary plan: `.opencode/plans/20260526-1256-db-sync-command.md`
- Discovery evidence: `.opencode/evidence/20260526-1256-db-sync-command/discovery.md`
- Evidence manifest: `.opencode/evidence/20260526-1256-db-sync-command/index.json`

Key decisions:

- Buat script baru `scripts/sync_remote_to_local.sh`.
- Gunakan custom-format dump + `pg_restore` sebagai default.
- Jangan hardcode password; gunakan env automation.
- Local target default `ggn_scyllax`.
- User sudah approve local drop/recreate, tetap gated oleh `--drop --yes` agar automation eksplisit.
- Citus install/prepare wajib karena local belum punya Citus.
- Remote write dilarang; script hanya boleh memakai `pg_dump` dan SELECT-only `psql` ke remote.

Assumptions:

- Local PostgreSQL ada di `localhost:5432`.
- Local user default `postgres`.
- Source dan local password diberikan lewat env saat execution.
- Full restore boleh overwrite local `ggn_scyllax`.

Open questions:

- Tidak ada open question yang memblokir implementasi. Credential tetap harus diberikan out-of-band saat run.

Readiness:

- Plan siap untuk implementasi oleh `@orchestrator`/`@fixer`.
- Tugas real preflight/restore hanya blocked oleh ketersediaan credential runtime dan kemampuan install Citus di mesin local.

Cleanup performed:

- Draft tidak digunakan.
- Evidence discovery tetap disimpan karena berguna untuk implementor dan audit plan.
