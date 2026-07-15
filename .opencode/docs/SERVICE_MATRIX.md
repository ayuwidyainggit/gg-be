# Service Matrix & README Authority

Repo-local matrix untuk service/module. Dokumen ini memisahkan detail matrix dari `ARCHITECTURE.md` agar architecture tetap ringkas.

## Authority model

- **authoritative**: README module bisa diandalkan untuk setup/runtime utama.
- **advisory**: README berguna sebagian, tapi tetap harus divalidasi dengan compose/env/Makefile.
- **stale-template**: README masih template/generic dan tidak layak jadi source utama.
- **missing**: README module tidak ditemukan.

## Service/module matrix (evidence-driven)

| Module | In root compose | Runtime style | Default port | Env files present | Makefile | Migration location/style | README authority |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `system` | yes | Fiber-oriented | `9001` | `.env` | no | no module-local migration dir found | stale-template |
| `master` | yes | Fiber-oriented | `9002` | `.env` | no | no module-local migration dir found | stale-template |
| `inventory` | yes | Fiber-oriented | `9003` | `.env` | no | no module-local migration dir found | stale-template |
| `sales` | yes | Fiber-oriented | `9004` | `.env` | no | no module-local migration dir found | stale-template |
| `finance` | yes | Fiber-oriented | `9005` | `.env` | no | no module-local migration dir found | stale-template |
| `tms` | yes | Fiber-oriented | `9006` | `.env.example`, `.env.demo`, `.env.live`, `.env.staging` | yes | `tms/migrations/` (SQL files + `migrate.go`) | missing |
| `mobile` | yes | Fiber-oriented | `9008` | `.env`, `.env.example` | no | no module-local migration dir found | stale-template |
| `pjp` | yes | Gin | `9010` | `.env` | yes | `pjp/database/migrate/` (`*.up.sql` / `*.down.sql`) | advisory |
| `cronjob` | yes | Fiber-oriented | `9100` | `.env` | no | no module-local migration dir found | missing |
| `pjp-principle` | no | PJP-family module | n/a in root compose | `.env` | yes | `pjp-principle/database/migrate/` | advisory |
| `pjp-sales` | no | Sales-family module | n/a in root compose | `.env` | no | no module-local migration dir found | stale-template |

## README authority audit notes

- `system`, `master`, `inventory`, `sales`, `finance`, `mobile`, `pjp-sales`: README berisi template GitLab generic, klasifikasi **stale-template**.
- `pjp`, `pjp-principle`: README berisi langkah dev/migrate nyata, tetapi port/command masih perlu verifikasi silang terhadap env/compose/Makefile, klasifikasi **advisory**.
- `tms`, `cronjob`: `README.md` module tidak ditemukan, klasifikasi **missing**.

Last reviewed: 2026-05-13.
