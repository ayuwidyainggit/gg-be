# Validation Evidence: Air Cache Reuse in Compose

Task ID: `20260526-1506-air-cache-compose`
Date: `2026-05-26`

## Files changed

- `docker-compose.yml`

## Intent

- Semua service Go yang jalan via `air` di compose harus reuse cache.
- Jangan install/download `air` dari nol tiap restart kalau versi tetap `v1.52.3`.
- Reuse cache Go modules dan build cache lintas restart container.

## Config changes

### Air command strategy

Semua app service yang pakai `air` diubah dari pola selalu-install:

```sh
go install github.com/air-verse/air@v1.52.3 && /go/bin/air -c .air.toml
```

menjadi pola conditional install:

```sh
export PATH=/go/bin:/usr/local/go/bin:$PATH
[ -x /go/bin/air ] && /go/bin/air -v 2>/dev/null | grep -Eq '(^|[[:space:]])v1\.52\.3($|[[:space:]])' || /usr/local/go/bin/go install github.com/air-verse/air@v1.52.3
exec /go/bin/air -c .air.toml
```

Artinya:
- kalau `/go/bin/air` ada dan versinya tepat `v1.52.3` -> langsung pakai cache
- kalau belum ada / versi beda -> install sekali lalu reuse

### Named volumes added

Semua service app Go di compose sekarang mount cache bersama:

- `go_mod_cache:/go/pkg/mod`
- `go_build_cache:/root/.cache/go-build`
- `go_bin_cache:/go/bin`

Volume dideklarasikan di top-level compose.

## Services covered

- `system`
- `master`
- `inventory`
- `sales`
- `finance`
- `mobile`
- `pjp`
- `cronjob`
- `tms`

## Commands run

### Recreate stack with updated compose

```bash
rtk docker compose -f "docker-compose.yml" up -d
```

Outcome:
- Stack recreated successfully.
- Non-blocking warning remains: compose `version` field obsolete.

### Verify cache volumes exist

```bash
docker volume ls --format '{{.Name}}' | rg 'go_(mod|build|bin)_cache|scylla'
```

Outcome:
- Found:
  - `scylla_go_mod_cache`
  - `scylla_go_build_cache`
  - `scylla_go_bin_cache`

### Verify mounts inside running container

```bash
docker inspect scylla-system --format '{{range .Mounts}}{{println .Type .Name .Destination}}{{end}}'
```

Outcome:
- `volume scylla_go_mod_cache /go/pkg/mod`
- `volume scylla_go_build_cache /root/.cache/go-build`
- `volume scylla_go_bin_cache /go/bin`
- source bind mount tetap `/app`

### First warm start behavior

```bash
docker logs --tail 80 scylla-system
```

Outcome:
- First warm boot still showed `go: downloading ...` lines.
- This expected, because cache volume was empty and had to be populated once.

### Verify Air binary cached

```bash
docker exec scylla-system sh -lc 'ls -l /go/bin; /go/bin/air -v'
```

Outcome:
- `/go/bin/air` exists.
- Version output:
  - `v1.52.3`

### Verify cache sizes

```bash
docker exec scylla-system sh -lc 'du -sh /go/pkg/mod /root/.cache/go-build /go/bin'
```

Outcome:
- `/go/pkg/mod` = `1.1G`
- `/root/.cache/go-build` = `620.3M`
- `/go/bin` = `12.8M`

This confirms caches are populated and persisted in mounted volumes.

### Warm restart proof: system

```bash
docker restart scylla-system
docker logs --since 40s scylla-system
```

Outcome:
- Restart log showed:
  - `air` banner
  - file watching
  - `building...`
- No repeated `go: downloading github.com/air-verse/air v1.52.3`
- No repeated long dependency download list for `air` itself
- This proves cached `air` binary reused after warm install.

### Warm restart proof: pjp

```bash
docker restart scylla-pjp
docker logs --since 60s scylla-pjp
docker exec scylla-pjp sh -lc '/go/bin/air -v'
```

Outcome:
- Restart log showed `air` banner, file watching, `building...`, `running...`
- App connected to database and booted routes.
- No repeated `go: downloading github.com/air-verse/air v1.52.3`
- `air -v` returned `v1.52.3`

### Warm restart proof: tms

```bash
docker restart scylla-tms
docker logs --since 40s scylla-tms
docker exec scylla-tms sh -lc '/go/bin/air -v'
```

Outcome:
- Restart log showed `air` banner and `building...` without repeated Air download lines.
- `air -v` returned `v1.52.3`
- App process itself panicked on missing `.env` after build (`panic: open .env: no such file or directory`).
- This panic is app/runtime config issue for `tms`, not Air cache regression.

## Fixes applied during task

### Compose variable interpolation issue

Initial command version used shell variable in compose command and triggered compose interpolation warning for `AIR_VERSION`.

Resolved by removing compose-time variable dependence and pinning literal version in command.

### PATH issue

Initial attempt relied on `go` being in runtime shell PATH and failed inside container.

Resolved by:
- prepending `PATH=/go/bin:/usr/local/go/bin:$PATH`
- calling `/usr/local/go/bin/go` explicitly for install path

## Conclusion

- Semua service compose yang pakai `air` sekarang support cache reuse.
- `air` tidak perlu install ulang dari nol tiap restart jika versi tetap `v1.52.3`.
- Go module cache, Go build cache, dan Go binary cache sekarang persisted via named volumes.
- First boot tetap warm cache sekali. Restart berikutnya reuse cache.

## Remaining notes

- First start setelah volume baru memang tetap download/build sekali.
- Kalau versi `air` diganti di compose, install ulang akan terjadi sekali lalu cache baru dipakai.
- `tms` punya masalah runtime `.env` sendiri yang terpisah dari task cache ini.
- Non-blocking compose warning masih ada untuk field `version` obsolete.
