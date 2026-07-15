# Repo Documentation Guide

Folder `docs/` berisi panduan tambahan (quick start, API, database, development notes).  
Untuk workflow agent/operator dan aturan operasional repo, **system of record ada di `.opencode/docs/`**.

## Source of Truth Order

Jika ada konflik antar dokumen, gunakan urutan berikut:

1. Runtime/module-local evidence: `docker-compose.yml`, module `go.mod`, env files, module `Makefile`
2. `.opencode/docs/` (terutama `index.md`, `AGENT_ROUTING.md`, `ARCHITECTURE.md`, `SERVICE_MATRIX.md`, `QUALITY.md`)
3. Root [`README.md`](../README.md)
4. Dokumen tambahan di folder `docs/`

## Recommended Reading

Mulai dari:

- [../README.md](../README.md)
- [../.opencode/docs/index.md](../.opencode/docs/index.md)

Lalu gunakan dokumen `docs/` sesuai kebutuhan:

- [QUICK_START.md](./QUICK_START.md)
- [DEVELOPMENT.md](./DEVELOPMENT.md)
- [API_STRUCTURE.md](./API_STRUCTURE.md)
- [DATABASE.md](./DATABASE.md)

## Notes

- Beberapa service README masih template/stale; cek status audit di `.opencode/docs/SERVICE_MATRIX.md`.
- Saat update dokumentasi, prioritaskan perubahan yang bisa diverifikasi dari file repo saat ini.
