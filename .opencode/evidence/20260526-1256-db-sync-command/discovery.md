# Discovery Evidence: Database Sync Command

Task ID: `20260526-1256-db-sync-command`

## Files inspected

- `scripts/README.md`
- `scripts/clone_db.sh`
- `scripts/restore_db.sh`
- `scripts/clone_staging.sh`
- `scripts/install_citus.sh`
- `docker-compose.yml`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `.opencode/docs/AGENT_ROUTING.md`

## Project patterns found

- Repo sudah punya helper database di `scripts/`:
  - `clone_db.sh` untuk clone remote ke local dari `.env`.
  - `restore_db.sh` untuk restore backup file.
  - `clone_staging.sh` untuk flow interaktif staging.
  - `install_citus.sh` untuk Citus local.
- `scripts/README.md` sudah dokumentasikan PostgreSQL client tools: `pg_dump`, `pg_restore`, `psql`, `createdb`, `dropdb`.
- `docker-compose.yml` mengandung remote DB dev sama dengan input user: host `103.28.219.73`, port `25431`, user `postgres`, db `scylla_citus_dev`. Password tidak disalin ke plan; gunakan secret/env.
- Repo guidance minta runtime check diawali `rtk docker compose -f docker-compose.yml ps` bila menyentuh runtime repo.
- Repo punya Citus concern: `install_citus.sh` dan clone scripts sudah punya handling `citus` / `citus_columnar`.

## Reuse candidates

- Reuse pola helper output dan prerequisite dari `scripts/clone_db.sh`, `restore_db.sh`, `clone_staging.sh`.
- Extend atau buat shell command baru di `scripts/sync_remote_to_local.sh` bila perlu flow non-interaktif khusus `scylla_citus_dev` → `ggn_scyllax`.
- Prefer custom-format dump (`pg_dump -Fc` + `pg_restore`) daripada SQL pipe langsung karena lebih aman untuk restore besar, retry, dan validasi.

## Constraints

- Jangan hardcode password baru di file, command history, docs, atau artifact. Gunakan `SOURCE_DB_PASSWORD` / `PGPASSWORD` dari env atau `.pgpass` lokal.
- User sudah mengizinkan target local `ggn_scyllax` drop/recreate; tetap perlu flag eksplisit `--drop --yes` agar automation tidak salah target.
- Citus source bisa punya extension/object yang local belum punya. User menyatakan local belum install Citus; script harus install/prepare local Citus atau stop sebelum restore bila gagal.
- Restore harus memakai `--no-owner --no-privileges` agar tidak bergantung role remote.
- Local PostgreSQL target: `localhost:5432`, user `postgres`, db `ggn_scyllax`; password tetap lewat env runtime dan tidak ditulis ke artifact.
- Remote database tidak boleh berubah: remote command hanya `pg_dump` dan `psql` SELECT/preflight.

## Risks

- Data local `ggn_scyllax` hilang bila drop/recreate dijalankan.
- Remote dump besar bisa lama dan butuh disk cukup bila memakai dump file.
- Citus distributed table metadata bisa gagal restore bila local bukan Citus coordinator yang cocok.
- Existing `clone_db.sh` direct pipe path berisiko karena stdout/stderr difilter dan diteruskan ke `psql`; plan harus mengutamakan custom dump path atau script baru.
- Remote dev credential sudah pernah muncul di repo tracked infra; jangan memperluas secret exposure.

## Commands/docs checked

- Tidak menjalankan koneksi DB karena user meminta plan, bukan execution.
- Local discovery memakai read-only file inspection.
- Official docs/context7 tidak diperlukan: fitur memakai PostgreSQL CLI standar yang sudah ada di repo scripts.
- GitHub/web/browser tidak diperlukan: tidak ada dependency upstream baru, reference UI, atau fakta eksternal.
