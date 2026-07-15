# Database Scripts

Kumpulan script untuk membantu operasi database di Scylla Backend.

## 📋 Daftar Script

### 1. `clone_db.sh` - Clone Database dari Remote ke Local

Script untuk clone database dari remote (berdasarkan konfigurasi di `.env`) ke local PostgreSQL.

**Usage:**
```bash
# Clone dari .env di root atau service pertama yang ditemukan
./scripts/clone_db.sh

# Clone dari .env di service tertentu
./scripts/clone_db.sh finance
./scripts/clone_db.sh master
./scripts/clone_db.sh cronjob
```

**Features:**
- Otomatis mencari file `.env` di root atau service
- Test koneksi sebelum clone
- Opsi untuk membuat backup file
- Auto-create local database jika belum ada
- Colored output untuk feedback yang jelas

**Environment Variables:**
Script akan membaca dari `.env`:
- `DB_HOST` - Remote database host
- `DB_PORT` - Remote database port
- `DB_USER` - Remote database user
- `DB_PASS` - Remote database password
- `DB_NAME` - Remote database name

**Local Database Config:**
Default local database config:
- Host: `localhost`
- Port: `5432`
- User: `postgres`
- Database: `scylla_db`

Bisa di-override dengan environment variables:
```bash
export LOCAL_DB_HOST=localhost
export LOCAL_DB_PORT=5432
export LOCAL_DB_USER=postgres
export LOCAL_DB_NAME=scylla_db_local
./scripts/clone_db.sh
```

**Contoh:**
```bash
# Clone dari finance/.env
./scripts/clone_db.sh finance

# Clone dengan custom local database name
LOCAL_DB_NAME=scylla_db_dev ./scripts/clone_db.sh
```

### 2. `clone_staging.sh` - Clone Database dari Staging (Interactive)

Script interaktif untuk clone database dari staging ke local dengan support multiple database selection.

**Usage:**
```bash
./scripts/clone_staging.sh
```

**Features:**
- Interactive input untuk staging credentials
- List semua database yang tersedia di staging
- Pilih database yang ingin di-clone
- Auto-check dan install Citus extension jika tersedia
- Handle Citus extension errors secara otomatis

**Contoh:**
```bash
./scripts/clone_staging.sh
# Akan muncul prompt untuk:
# - Staging DB credentials
# - Pilih database dari list
# - Nama database lokal
```

### 3. `install_citus.sh` - Install Citus Extension

Script untuk install Citus extension di local PostgreSQL.

**Usage:**
```bash
./scripts/install_citus.sh
```

**Features:**
- Auto-detect OS (macOS, Ubuntu/Debian, RHEL/CentOS)
- Install Citus via package manager
- Create extension di database

### 4. `sync_remote_to_local.sh` - Sync Remote `scylla_citus_dev` ke Local `ggn_scyllax` (Automation-first)

Script non-interaktif untuk sync remote ke local via `pg_dump -Fc` dan `pg_restore`.

**Usage:**
```bash
# Lihat help
./scripts/sync_remote_to_local.sh --help

# Preflight saja (cek koneksi, cek requirement, cek citus source/local)
SOURCE_DB_PASSWORD='<from-secure-env>' LOCAL_DB_PASSWORD='<from-secure-env>' \
  ./scripts/sync_remote_to_local.sh --preflight-only

# Dry run (cek tanpa dump/restore)
SOURCE_DB_PASSWORD='<from-secure-env>' LOCAL_DB_PASSWORD='<from-secure-env>' \
  ./scripts/sync_remote_to_local.sh --dry-run

# Full sync destructive local (wajib --drop --yes)
SOURCE_DB_PASSWORD='<from-secure-env>' LOCAL_DB_PASSWORD='<from-secure-env>' \
  ./scripts/sync_remote_to_local.sh --install-citus --drop --yes
```

**Default Remote Config:**
- Host: `103.28.219.73`
- Port: `25431`
- User: `postgres`
- DB: `scylla_citus_dev`

**Default Local Config:**
- Host: `localhost`
- Port: `5432`
- User: `postgres`
- DB: `ggn_scyllax`

**Required Environment Variables:**
- `SOURCE_DB_PASSWORD`
- `LOCAL_DB_PASSWORD`

**Flags:**
- `--install-citus`
- `--drop`
- `--yes`
- `--preflight-only`
- `--dry-run`
- `--keep-dump`
- `--dump-file /path/to/file.dump`
- `--allow-nonlocal-target`
- `--allow-custom-target`

**Safety Rules:**
- Remote read-only: script hanya pakai `pg_dump` dan SELECT-only `psql` query.
- Tidak ada remote DDL/DML.
- Destructive action hanya local.
- Drop/recreate local butuh `--drop --yes`.
- Destructive mode default hanya boleh ke `localhost` / `127.0.0.1`; host lain perlu `--allow-nonlocal-target`.
- Destructive mode default hanya boleh ke DB `ggn_scyllax`; nama lain perlu `--allow-custom-target`.
- Tolak local DB berbahaya: `postgres`, `template0`, `template1`, atau kosong.
- Default dump file: `/tmp/scylla_citus_dev_YYYYMMDD_HHMMSS.dump`.
- Restore pakai `pg_restore --clean --if-exists --no-owner --no-privileges`.

**Citus Note:**
- Jika source pakai Citus dan local belum siap, jalankan dengan `--install-citus`.
- Script akan panggil `./scripts/install_citus.sh --non-interactive --create-extension ...`.
- PostgreSQL local juga harus preload Citus lewat `shared_preload_libraries = 'citus'`. Jika PostgreSQL di-upgrade/reinstall, cek lagi setting ini dan restart service.

### 5. `restore_db.sh` - Restore Database dari Backup File

Script untuk restore database dari file backup (`.dump`, `.sql`, atau `.sql.gz`) ke local PostgreSQL.

**Usage:**
```bash
# Restore dari backup.dump (default)
./scripts/restore_db.sh

# Restore dari file tertentu
./scripts/restore_db.sh backup.dump scylla_db
./scripts/restore_db.sh backup.sql scylla_db
./scripts/restore_db.sh backup.sql.gz scylla_db
```

**Features:**
- Support multiple backup formats (`.dump`, `.backup`, `.sql`, `.sql.gz`)
- Auto-detect backup file format
- Auto-create database jika belum ada
- Colored output untuk feedback yang jelas

**Contoh:**
```bash
# Restore dari backup.dump yang ada di root
./scripts/restore_db.sh backup.dump scylla_db

# Restore dari compressed SQL
./scripts/restore_db.sh backup.sql.gz scylla_db
```

## 🔧 Prerequisites

Script memerlukan PostgreSQL client tools:
- `pg_dump` - Untuk dump database
- `pg_restore` - Untuk restore dari custom format
- `psql` - Untuk restore dari SQL file
- `createdb` - Untuk create database
- `dropdb` - Untuk drop database

### Citus Extension (Optional)

Jika database staging menggunakan Citus extension, Anda bisa install di local:

```bash
# Install Citus extension
./scripts/install_citus.sh

# Atau manual:
# macOS: brew install citus
# Ubuntu/Debian: See https://docs.citusdata.com/en/stable/installation/
```

**Install PostgreSQL Client Tools:**

**macOS:**
```bash
brew install postgresql
```

**Ubuntu/Debian:**
```bash
sudo apt-get install postgresql-client
```

**CentOS/RHEL:**
```bash
sudo yum install postgresql
```

## 📝 Contoh Penggunaan

### Scenario 1: Clone Database dari Production ke Local

```bash
# 1. Pastikan .env sudah dikonfigurasi dengan production credentials
# 2. Clone database
./scripts/clone_db.sh

# Script akan:
# - Membaca .env
# - Test koneksi ke production
# - Create local database
# - Clone semua data
```

### Scenario 2: Restore dari Backup File

```bash
# 1. Pastikan backup file ada
ls -lh backup.dump

# 2. Restore ke local database
./scripts/restore_db.sh backup.dump scylla_db
```

### Scenario 3: Clone dengan Backup File

```bash
# Clone dan simpan backup file
./scripts/clone_db.sh finance
# Pilih 'y' ketika ditanya "Do you want to create a backup file?"
# Backup akan disimpan sebagai backup_YYYYMMDD_HHMMSS.dump
```

## ⚠️ Troubleshooting

### Error: Connection refused

**Problem:** Tidak bisa connect ke remote database

**Solution:**
- Check apakah remote database accessible
- Verify credentials di `.env`
- Check firewall/security group settings
- Test koneksi manual: `psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME`

### Error: Database already exists

**Problem:** Local database sudah ada

**Solution:**
- Script akan menanyakan apakah ingin drop dan recreate
- Pilih 'y' untuk drop dan recreate
- Atau drop manual: `dropdb scylla_db`

### Error: Permission denied

**Problem:** Tidak punya permission untuk create/drop database

**Solution:**
- Pastikan user PostgreSQL punya permission
- Atau gunakan superuser: `sudo -u postgres ./scripts/clone_db.sh`

### Error: pg_dump not found

**Problem:** PostgreSQL client tools belum terinstall

**Solution:**
- Install PostgreSQL client tools (lihat Prerequisites)
- Pastikan `pg_dump`, `psql`, dll ada di PATH

## 🔒 Security Notes

- **Jangan commit file `.env`** yang berisi production credentials
- **Jangan commit backup files** yang berisi production data
- Gunakan `.gitignore` untuk exclude:
  ```
  .env
  *.dump
  backup_*.dump
  backup_*.sql
  ```

## 📚 Related Documentation

- [Database Documentation](../docs/DATABASE.md)
- [Development Guidelines](../docs/DEVELOPMENT.md)
- [Quick Start Guide](../docs/QUICK_START.md)

---

**Last Updated**: 2024

