# Scylla Backend Services

Scylla Backend adalah monorepo backend ERP berbasis Go dengan banyak service/module untuk master data, sales, inventory, finance, mobile, workflow PJP, dan cronjob. Untuk workflow agent/operator, repo ini memakai `.opencode/docs/` sebagai system of record.

## Overview

Domain utama yang dicakup:

- **System** — user management, config, menu, notification
- **Master** — master data produk, outlet, employee, pricing, discount
- **Sales** — order, invoice, return, promotion, approval
- **Inventory** — stock, warehouse, transfer, adjustment
- **Finance** — AP/AR, payment, cheque, deposit, VAT
- **Mobile** — API untuk aplikasi mobile
- **TMS** — transport/task management module terpisah di repo
- **PJP family** — `pjp`, `pjp-principle`, `pjp-sales`
- **Cronjob** — scheduler untuk task otomatis

## System of Record

Mulai dari dokumen berikut:

- `.opencode/docs/index.md`
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/SERVICE_MATRIX.md`
- `.opencode/docs/QUALITY.md`

Jika README ini bertentangan dengan `docker-compose.yml`, `go.mod`, `Makefile`, env file module, atau `.opencode/docs/`, utamakan sumber yang paling dekat dengan module/runtime yang sedang disentuh.

## Architecture

Repo ini adalah **multi-module Go monorepo**:

- tidak ada root `go.mod`
- setiap service/module punya `go.mod` sendiri
- mayoritas service utama bersifat Fiber-oriented
- `pjp` adalah service Gin-based utama yang berbeda dari service Fiber utama

Layering yang harus dipertahankan:

```text
Controller → Service → Repository → DB
```

Aturan penting:

- controller tidak memanggil repository langsung
- repository tidak memuat business logic
- write operation harus berjalan lewat transaction service layer
- query tenant-aware wajib menjaga `cust_id` / `parent_cust_id` sesuai konteks

## Services and Modules

### Compose-managed default stack

Root `docker-compose.yml` mengelola service berikut:

| Module | Runtime style | Default port |
| --- | --- | --- |
| `system` | Fiber-oriented | `9001` |
| `master` | Fiber-oriented | `9002` |
| `inventory` | Fiber-oriented | `9003` |
| `sales` | Fiber-oriented | `9004` |
| `finance` | Fiber-oriented | `9005` |
| `tms` | Fiber-oriented | `9006` |
| `mobile` | Fiber-oriented | `9008` |
| `pjp` | Gin | `9010` |
| `cronjob` | Fiber-oriented | `9100` |
| `redis` | infra dependency | `6379` |

### Extra modules in repo

Module berikut ada di repo tetapi tidak menjadi bagian dari default root compose stack:

- `pjp-principle`
- `pjp-sales`

Untuk matrix yang lebih detail soal env, Makefile, migration style, dan status authority README per module, lihat `.opencode/docs/SERVICE_MATRIX.md`.

## Prerequisites

- **Go** — ikuti versi pada `go.mod` module target; versi berbeda antar module
- **PostgreSQL** — untuk akses DB lokal/remote dan restore scripts
- **Redis** — dipakai beberapa service
- **Docker** — untuk default local stack
- **RTK** — workflow shell repo ini memakai prefix `rtk`
- PostgreSQL client tools (`pg_dump`, `pg_restore`, `psql`, `createdb`, `dropdb`) bila memakai script database

## Setup

### 1. Clone repository

```bash
git clone <repository-url>
cd scylla-be
```

### 2. Check runtime stack

Jalankan dari root repo:

```bash
rtk docker compose -f docker-compose.yml ps
```

Kalau service yang dibutuhkan belum jalan:

```bash
rtk docker compose -f docker-compose.yml up -d
```

### 3. Setup environment files

Convention env file berbeda per module:

- mayoritas service memakai `.env`
- `mobile` punya `.env` dan `.env.example`
- `tms` punya beberapa varian `.env.*`
- `pjp` dan `pjp-principle` punya `.env`, sementara workflow Makefile/module docs juga mereferensikan `development.env`

Selalu cek file env yang benar-benar dipakai module target sebelum run/test/migrate.

### 4. Database setup

Script database tersedia di `scripts/`.

Contoh:

```bash
./scripts/clone_db.sh finance
./scripts/restore_db.sh backup.dump scylla_db
```

Detail lengkap ada di `scripts/README.md`.

## Development Workflow

Jalankan dari direktori module target:

```bash
cd <module>
rtk go mod download
rtk go mod tidy
rtk go test ./...
rtk go run main.go
```

Catatan:

- banyak compose-managed service memakai `air` dalam flow compose/dev-reload
- untuk module yang punya `Makefile`, gunakan target module-local bila tersedia
- untuk perubahan non-trivial, cek `.opencode/docs/QUALITY.md` dan `.opencode/docs/ARCHITECTURE.md`

## Validation

Baseline validation:

```bash
rtk go mod download && rtk go mod tidy
rtk go test ./...
rtk go test ./path/to/pkg -run TestName
```

Untuk cheat sheet per service/module, lihat `.opencode/docs/QUALITY.md`.

## Database and Migrations

Workflow migration tidak seragam antar module:

- `tms` memakai `migrations/` dan `rtk make migrateUp`
- `pjp` memakai `database/migrate/` dan `rtk make migrateUp`
- `pjp-principle` memakai `database/migrate/` dan `rtk make migrateUp`
- module lain tidak punya workflow migration lokal yang terdokumentasi jelas di root docs

Jangan anggap semua service punya `Makefile` atau migration command yang sama.

## Additional Docs

- `.opencode/docs/index.md` — pintu masuk docs workflow agent/operator
- `scripts/README.md` — clone/restore DB scripts
- `docs/README.md` — docs tambahan repo
- `.opencode/docs/ARCHITECTURE.md` — aturan arsitektur lintas module
- `.opencode/docs/SERVICE_MATRIX.md` — inventory service/module + audit authority README
- `pjp/README.md` dan `pjp-principle/README.md` — konteks tambahan workflow PJP family

Catatan penting:

- beberapa README service masih template/stale
- bila ragu, utamakan `docker-compose.yml`, `go.mod`, `Makefile`, env file module, dan `.opencode/docs/`
