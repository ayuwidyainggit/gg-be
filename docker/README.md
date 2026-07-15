# Docker Setup untuk Scylla Backend Services

Setup Docker untuk menjalankan semua service Scylla Backend termasuk PostgreSQL dengan Citus extension.

## 🚀 Quick Start

### Development Mode (Hot Reload) - Recommended untuk Development

Untuk development dengan hot reload (code changes langsung terupdate):

```bash
# Start development mode dengan volume mounting
./docker/start-dev.sh

# Atau manual
cd docker
docker-compose -f docker-compose.dev.yml up -d
```

**Keuntungan:**
- ✅ Code changes di local langsung terlihat di container
- ✅ Tidak perlu rebuild image
- ✅ Restart cepat untuk apply changes

Lihat [Development Mode Guide](./DEVELOPMENT.md) untuk detail lengkap.

### Production Mode (Build Image)

Untuk production atau testing dengan compiled binary:

```bash
# 1. Setup Environment (Optional)
cd docker
cp .env.example .env
# Edit .env jika perlu (untuk OBS configuration)

# 2. Start All Services
docker-compose up -d --build

# Atau menggunakan helper script
./docker/start-all.sh
```

### 3. Verify All Containers Running

```bash
docker-compose ps

# Atau
docker ps | grep scylla
```

### 4. Check Service Health

```bash
# System Service
curl http://localhost:9001/ping

# Master Service
curl http://localhost:9002/ping

# Inventory Service
curl http://localhost:9003/ping

# Sales Service
curl http://localhost:9004/ping

# Finance Service
curl http://localhost:9005/ping

# Mobile Service
curl http://localhost:9008/ping

# Cronjob Service
curl http://localhost:9100/ping
```

### 5. Connect ke Database

```bash
# Menggunakan psql (port 54321 untuk menghindari konflik dengan PostgreSQL lokal)
psql -h localhost -p 54321 -U postgres -d scylla_db

# Atau menggunakan Docker
docker exec -it scylla-postgres-citus psql -U postgres -d scylla_db
```

### 6. Verify Citus Extension

```sql
-- Di dalam psql
\dx

-- Akan muncul:
-- citus
-- citus_columnar
```

## 📋 Configuration

### Service Ports

| Service | Port | Health Check |
|---------|------|--------------|
| System | 9001 | http://localhost:9001/ping |
| Master | 9002 | http://localhost:9002/ping |
| Inventory | 9003 | http://localhost:9003/ping |
| Sales | 9004 | http://localhost:9004/ping |
| Finance | 9005 | http://localhost:9005/ping |
| Mobile | 9008 | http://localhost:9008/ping |
| Cronjob | 9100 | http://localhost:9100/ping |
| PostgreSQL | 54321 | psql -h localhost -p 54321 |
| Redis | 6379 | redis-cli -h localhost ping |

### Database Credentials

- **Host**: `postgres-citus` (dalam Docker network) atau `localhost` (dari host)
- **Port**: `5432` (dalam Docker) atau `54321` (dari host)
- **User**: `postgres`
- **Password**: `postgres`
- **Database**: `scylla_db`

### Redis Credentials

- **Host**: `redis` (dalam Docker network) atau `localhost` (dari host)
- **Port**: `6379`
- **Password**: (kosong)

### Volumes

- `postgres_citus_data`: PostgreSQL data persistence
- `redis_data`: Redis data persistence

## 🔧 Commands

### Start All Services

```bash
# Start semua service
docker-compose up -d

# Start dengan build
docker-compose up -d --build

# Start specific service
docker-compose up -d postgres-citus
docker-compose up -d system
```

### Stop All Services

```bash
# Stop semua service
docker-compose down

# Stop dan remove volumes (hapus semua data)
docker-compose down -v
```

### View Logs

```bash
# Logs semua service
docker-compose logs -f

# Logs specific service
docker-compose logs -f system
docker-compose logs -f master
docker-compose logs -f postgres-citus
```

### Restart Services

```bash
# Restart semua service
docker-compose restart

# Restart specific service
docker-compose restart system
```

### Rebuild Services

```bash
# Rebuild semua service
docker-compose build

# Rebuild specific service
docker-compose build system

# Rebuild dan restart
docker-compose up -d --build system
```

### Service Management

```bash
# Start specific service
docker-compose start system

# Stop specific service
docker-compose stop system

# Scale service (jika diperlukan)
docker-compose up -d --scale system=2
```

## 🔗 Connect dari Application

### Untuk Development (dari Host Machine)

Update `.env` file di setiap service untuk connect ke Docker services:

```env
# Database (menggunakan port 54321 untuk menghindari konflik)
DB_HOST=localhost
DB_PORT=54321
DB_USER=postgres
DB_PASS=postgres
DB_NAME=scylla_db

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
```

### Untuk Services dalam Docker Network

Services dalam Docker otomatis menggunakan:
- `DB_HOST=postgres-citus` (service name)
- `REDIS_HOST=redis` (service name)

## 📦 Clone Database ke Docker

Setelah container running, gunakan script clone dengan port 54321:

```bash
# Set environment variable untuk local database
export LOCAL_DB_HOST=localhost
export LOCAL_DB_PORT=54321
export LOCAL_DB_USER=postgres
export LOCAL_DB_NAME=scylla_db

# Clone dari staging
./scripts/clone_staging.sh
```

## 🏗️ Service Architecture

```
┌─────────────────────────────────────────────────┐
│              Docker Network                     │
│                                                 │
│  ┌──────────────┐    ┌──────────────┐          │
│  │ PostgreSQL   │    │    Redis     │          │
│  │ (Citus)      │    │              │          │
│  │ :5432        │    │   :6379      │          │
│  └──────────────┘    └──────────────┘          │
│         │                    │                   │
│         └────────┬───────────┘                   │
│                  │                               │
│  ┌──────────────┼──────────────┐               │
│  │              │               │               │
│  System    Master    Inventory                  │
│  :9001     :9002     :9003                      │
│                                                 │
│  Sales     Finance   Mobile    Cronjob          │
│  :9004     :9005     :9008     :9100            │
│                                                 │
└─────────────────────────────────────────────────┘
         │
         │ Port Mapping
         │
┌────────┴────────────────────────────────────────┐
│              Host Machine                       │
│                                                 │
│  PostgreSQL: localhost:54321                    │
│  Redis:      localhost:6379                     │
│  System:     localhost:9001                     │
│  Master:     localhost:9002                     │
│  Inventory:  localhost:9003                     │
│  Sales:      localhost:9004                     │
│  Finance:    localhost:9005                     │
│  Mobile:     localhost:9008                     │
│  Cronjob:    localhost:9100                     │
└─────────────────────────────────────────────────┘
```

## 🐛 Troubleshooting

### Port Already in Use

Jika port 54321 sudah digunakan:

```bash
# Cek process yang menggunakan port
lsof -i :54321

# Atau ubah port di docker-compose.postgres-citus.yml
ports:
  - "54322:5432"  # Ganti ke port lain
```

### Container Won't Start

```bash
# Check logs
docker-compose -f docker-compose.postgres-citus.yml logs

# Remove dan recreate
docker-compose -f docker-compose.postgres-citus.yml down -v
docker-compose -f docker-compose.postgres-citus.yml up -d
```

### Extension Not Found

Extension Citus akan otomatis di-install saat container pertama kali dibuat. Jika tidak muncul:

```sql
-- Connect ke database
psql -h localhost -p 54321 -U postgres -d scylla_db

-- Create extension manually
CREATE EXTENSION IF NOT EXISTS citus;
CREATE EXTENSION IF NOT EXISTS citus_columnar;
```

## 📝 Environment Variables

### Docker Compose Mode

**TIDAK PERLU** file `.env` di service directory. Environment variables sudah dikonfigurasi di `docker-compose.yml`.

### Local Development Mode

**PERLU** file `.env` jika menjalankan service dengan `go run main.go`:

```bash
# Generate .env untuk semua service (connect ke Docker)
./scripts/generate-env.sh docker

# Atau untuk local PostgreSQL
./scripts/generate-env.sh local
```

Lihat [Environment Variables Setup Guide](./ENV_SETUP.md) untuk detail lengkap.

## 📚 Related Documentation

- [Development Mode](./DEVELOPMENT.md) - Panduan development dengan hot reload
- [Environment Variables Setup](./ENV_SETUP.md) - Panduan setup .env
- [Database Scripts](../scripts/README.md)
- [Database Documentation](../docs/DATABASE.md)

---

**Last Updated**: 2024

