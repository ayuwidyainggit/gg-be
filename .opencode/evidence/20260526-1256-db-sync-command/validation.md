# Validation Evidence: Database Sync Command

Task ID: `20260526-1256-db-sync-command`
Date: `2026-05-26`

## Files changed

Repo files:
- `scripts/sync_remote_to_local.sh`
- `scripts/README.md`

Local machine config outside repo:
- `/opt/homebrew/var/postgresql@18/postgresql.conf`

Repo docs/scripts not changed during this task:
- `scripts/install_citus.sh`

## Commands run

### Syntax and help

```bash
bash -n "scripts/sync_remote_to_local.sh"
bash -n "scripts/install_citus.sh"
bash "scripts/sync_remote_to_local.sh" --help
bash "scripts/install_citus.sh" --help
```

Outcome:
- Passed.

### Missing secret failure

```bash
env -u SOURCE_DB_PASSWORD -u LOCAL_DB_PASSWORD bash "scripts/sync_remote_to_local.sh" --dry-run
```

Outcome:
- Failed as expected with: `SOURCE_DB_PASSWORD is required`
- No password literal printed.

### Runtime baseline

```bash
rtk docker compose -f "docker-compose.yml" ps
```

Outcome:
- Compose command ran.
- No services active at validation time.

### Preflight before local Citus preload fix

```bash
SOURCE_DB_PASSWORD='<redacted>' LOCAL_DB_PASSWORD='<redacted>' bash "scripts/sync_remote_to_local.sh" --preflight-only
```

Outcome:
- Remote preflight succeeded.
- Local preflight succeeded.
- Source reported Citus usage.
- Failed safely with: `Local Citus not available. Re-run with --install-citus`

### Local Citus diagnosis

```bash
psql --version
PGPASSWORD='<redacted>' psql -h localhost -p 5432 -U postgres -d postgres -At -c "show config_file; show data_directory; show shared_preload_libraries;"
PGPASSWORD='<redacted>' psql -h localhost -p 5432 -U postgres -d postgres -At -c "select name, default_version, installed_version from pg_available_extensions where name in ('citus','citus_columnar');"
brew info citus
```

Outcome:
- Local server: PostgreSQL 18.3.
- Citus extension files available locally.
- Failure root cause: `shared_preload_libraries` did not include `citus`.

### Local PostgreSQL config remediation

Change applied outside repo:
- In `/opt/homebrew/var/postgresql@18/postgresql.conf`, enabled:
  - `shared_preload_libraries = 'citus'`

Then restarted service:

```bash
brew services restart postgresql@18
```

Verification:

```bash
PGPASSWORD='<redacted>' psql -h localhost -p 5432 -U postgres -d postgres -At -c "show shared_preload_libraries;"
```

Outcome:
- Returned `citus`.

### Preflight after Citus preload fix

```bash
SOURCE_DB_PASSWORD='<redacted>' LOCAL_DB_PASSWORD='<redacted>' bash "scripts/sync_remote_to_local.sh" --preflight-only
```

Outcome:
- Passed.

### Full sync execution

```bash
SOURCE_DB_PASSWORD='<redacted>' LOCAL_DB_PASSWORD='<redacted>' bash "scripts/sync_remote_to_local.sh" --install-citus --drop --yes
```

Outcome:
- Local DB `ggn_scyllax` dropped and recreated.
- Local `citus` and `citus_columnar` extensions created in target DB.
- Remote dump created via `pg_dump -Fc` to `/tmp/scylla_citus_dev_20260526_141405.dump`.
- Dump restored via `pg_restore --clean --if-exists --no-owner --no-privileges`.
- Post-restore validation in script returned:
  - `ggn_scyllax|postgres`
  - `484` non-system tables from `information_schema.tables`
  - extensions: `citus`, `citus_columnar`, `pgcrypto`, `plpgsql`, `uuid-ossp`
- Dump file removed after success.

### Manual post-checks

```bash
PGPASSWORD='<redacted>' psql -h localhost -p 5432 -U postgres -d ggn_scyllax -c "select current_database(), current_user;"
PGPASSWORD='<redacted>' psql -h localhost -p 5432 -U postgres -d ggn_scyllax -c "select schemaname, count(*) from pg_tables where schemaname not in ('pg_catalog','information_schema') group by schemaname order by schemaname;"
PGPASSWORD='<redacted>' psql -h localhost -p 5432 -U postgres -d ggn_scyllax -At -c "select extname from pg_extension order by extname;"
```

Outcome:
- Current DB/user: `ggn_scyllax` / `postgres`
- Non-system schema counts:
  - `acf=61`
  - `columnar_internal=4`
  - `import=45`
  - `inv=58`
  - `mobile=2`
  - `mst=153`
  - `orlin=9`
  - `picklist=3`
  - `pjp=9`
  - `pjp_principles=9`
  - `promo=12`
  - `public=2`
  - `report=8`
  - `sls=45`
  - `smc=8`
  - `sys=22`
  - `tms=3`
- Extensions:
  - `citus`
  - `citus_columnar`
  - `pgcrypto`
  - `plpgsql`
  - `uuid-ossp`

### Extra destructive guard validation

```bash
SOURCE_DB_PASSWORD='redacted' LOCAL_DB_PASSWORD='redacted' LOCAL_DB_HOST='192.0.2.10' bash "scripts/sync_remote_to_local.sh" --drop --yes --dry-run
```

Outcome:
- Failed as expected with: `Refuse destructive non-local target host '192.0.2.10'. Use --allow-nonlocal-target to override`

## Remote safety confirmation

Remote command surface used during execution:
- `psql ... -c "SELECT ..."` for preflight and extension detection only
- `pg_dump -Fc` for dump creation

Not used on remote:
- `dropdb`
- `createdb`
- `pg_restore`
- `CREATE/ALTER/DROP/INSERT/UPDATE/DELETE/TRUNCATE`

Conclusion:
- Remote database remained read-only from script behavior.

## Remaining risks / notes

- Sync depends on local PostgreSQL keeping `shared_preload_libraries = 'citus'`.
- Local machine config was changed outside repo to satisfy Citus runtime requirement.
- Future local PostgreSQL upgrades/reinstalls may require reapplying Citus preload.
- Destructive override flags now exist for non-local host / custom target; use only when intentionally targeting non-default local destination.
