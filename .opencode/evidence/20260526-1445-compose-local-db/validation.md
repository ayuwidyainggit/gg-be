# Validation Evidence: Compose Runtime Switched to Local DB

Task ID: `20260526-1445-compose-local-db`
Date: `2026-05-26`

## Files changed

- `docker-compose.yml`
- `AGENTS.md`
- `.opencode/docs/ARCHITECTURE.md`

## Intent

- Semua service compose harus pakai database local host, bukan remote dev DB.
- Agent/harness guidance harus menganggap local Postgres sebagai runtime default.

## Config changes

### docker-compose.yml

Semua service compose yang sebelumnya memakai remote dev DB diubah ke baseline local berikut:

- Host: `host.docker.internal`
- Port: `5432`
- User: `postgres`
- Password: `postgres`
- DB: `ggn_scyllax`

Perubahan diterapkan pada:
- `system`
- `master`
- `inventory`
- `sales`
- `finance`
- `mobile`
- `pjp` (`POSTGRES_*` dan `DB_*`)
- `cronjob`
- `tms` (`POSTGRES_*` dan `DB_*`)

### AGENTS.md

Ditambah rule runtime eksplisit:
- compose harus target local Postgres `host.docker.internal:5432`, `postgres/postgres`, DB `ggn_scyllax`
- jangan balik ke remote dev DB default

### .opencode/docs/ARCHITECTURE.md

Ditambah baseline compose DB runtime:
- `host.docker.internal:5432`
- `postgres/postgres`
- `ggn_scyllax`

## Commands run

### Restart stack

```bash
rtk docker compose -f "docker-compose.yml" ps
rtk docker compose -f "docker-compose.yml" up -d
```

Outcome:
- Compose stack restarted with updated config.
- Warning only: compose `version` attribute obsolete.

### Running containers

```bash
docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'
```

Outcome:
- Containers up:
  - `scylla-system`
  - `scylla-master`
  - `scylla-inventory`
  - `scylla-sales`
  - `scylla-finance`
  - `scylla-mobile`
  - `scylla-pjp`
  - `scylla-cronjob`
  - `scylla-tms`
  - `scylla-rabbitmq`
  - `scylla-redis`

### Env verification inside containers

```bash
docker inspect scylla-system --format '{{range .Config.Env}}{{println .}}{{end}}'
docker inspect scylla-tms --format '{{range .Config.Env}}{{println .}}{{end}}'
docker inspect scylla-pjp --format '{{range .Config.Env}}{{println .}}{{end}}'
```

Outcome:
- `scylla-system` has:
  - `DB_HOST=host.docker.internal`
  - `DB_PORT=5432`
  - `DB_USER=postgres`
  - `DB_PASS=postgres`
  - `DB_NAME=ggn_scyllax`
- `scylla-tms` has:
  - `POSTGRES_HOST=host.docker.internal`
  - `POSTGRES_PORT=5432`
  - `POSTGRES_USER=postgres`
  - `POSTGRES_PASSWORD=postgres`
  - `POSTGRES_DB=ggn_scyllax`
  - plus matching `DB_*`
- `scylla-pjp` has:
  - `POSTGRES_HOST=host.docker.internal`
  - `POSTGRES_PORT=5432`
  - `POSTGRES_USER=postgres`
  - `POSTGRES_PASSWORD=postgres`
  - `POSTGRES_DB=ggn_scyllax`
  - plus matching `DB_*`

### Reachability from containers to local DB

```bash
docker exec scylla-system sh -lc 'nc -zvw3 host.docker.internal 5432'
docker exec scylla-tms sh -lc 'nc -zvw3 host.docker.internal 5432'
docker exec scylla-pjp sh -lc 'nc -zvw3 host.docker.internal 5432'
```

Outcome:
- All returned `host.docker.internal (...:5432) open`

### Conclusive DB login checks from containers

```bash
docker exec scylla-system sh -lc 'command -v psql >/dev/null 2>&1 || apk add --no-cache postgresql-client >/dev/null; PGPASSWORD="$DB_PASS" psql "host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASS dbname=$DB_NAME application_name=scylla-system-compose-check" -At -c "select current_database(), current_user"'
docker exec scylla-tms sh -lc 'command -v psql >/dev/null 2>&1 || apk add --no-cache postgresql-client >/dev/null; PGPASSWORD="$POSTGRES_PASSWORD" psql "host=$POSTGRES_HOST port=$POSTGRES_PORT user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB application_name=scylla-tms-compose-check" -At -c "select current_database(), current_user"'
docker exec scylla-pjp sh -lc 'command -v psql >/dev/null 2>&1 || apk add --no-cache postgresql-client >/dev/null; PGPASSWORD="$POSTGRES_PASSWORD" psql "host=$POSTGRES_HOST port=$POSTGRES_PORT user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB application_name=scylla-pjp-compose-check" -At -c "select current_database(), current_user"'
```

Outcome:
- Each container connected successfully and returned:
  - `ggn_scyllax|postgres`
- This is conclusive proof containers are using local DB target and credentials.

### App runtime health

```bash
curl -i "http://127.0.0.1:9001/ping"
curl -i "http://127.0.0.1:9002/ping"
```

Outcome:
- `system` returned `HTTP/1.1 200 OK` and body `It works`
- `master` returned `HTTP/1.1 200 OK` and body `It works`

### App logs

```bash
docker logs --tail 120 scylla-system
docker logs --tail 120 scylla-master
```

Outcome:
- Services built and reached running state.
- `system` served `/ping` request successfully.

### Local DB activity checks

```bash
PGPASSWORD='postgres' psql -h localhost -p 5432 -U postgres -d postgres -c "select client_addr, usename, datname, application_name, state, count(*) from pg_stat_activity where datname='ggn_scyllax' group by 1,2,3,4,5 order by count(*) desc;"
```

Outcome:
- Existing local activity visible on `ggn_scyllax`.
- Short-lived tagged `psql` sessions from containers were not captured in `pg_stat_activity` snapshot because they disconnected immediately after query, but direct successful logins above already prove DB target correctness.

## Conclusion

- Docker compose config now points to local DB, not remote dev DB.
- Agent and harness repo guidance now enforce local DB baseline.
- Docker stack restarted successfully.
- Verified running containers carry local DB env vars.
- Verified containers can reach `host.docker.internal:5432`.
- Verified containers can authenticate to `ggn_scyllax` as `postgres` using compose env vars.
- Verified key services are up and serving requests.

## Remaining notes

- `host.docker.internal` assumption is valid for Docker Desktop on macOS. Linux daemon setups may need explicit host-gateway mapping.
- Compose emits non-blocking warning that `version` key is obsolete.
