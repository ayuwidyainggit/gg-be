# Environment Variables Setup Guide

Panduan setup environment variables untuk berbagai skenario penggunaan.

## 📋 Skenario Penggunaan

### 1. **Menggunakan Docker Compose (Semua Service di Docker)**

**TIDAK PERLU** update `.env` di masing-masing service.

Environment variables sudah dikonfigurasi di `docker/docker-compose.yml`. Service dalam Docker akan menggunakan environment variables dari docker-compose, bukan dari file `.env`.

**Cara kerja:**
- Docker Compose meng-inject environment variables langsung ke container
- Service membaca dari `os.Getenv()` yang sudah di-set oleh Docker
- File `.env` di service directory **tidak digunakan** saat running di Docker

### 2. **Development Lokal (go run) - Connect ke Docker Services**

**PERLU** update `.env` di masing-masing service.

Jika Anda menjalankan service secara lokal dengan `go run main.go` tapi ingin connect ke PostgreSQL/Redis di Docker:

```env
# Database (connect ke Docker PostgreSQL)
DB_HOST=localhost
DB_PORT=54321          # Port Docker PostgreSQL
DB_USER=postgres
DB_PASS=postgres
DB_NAME=scylla_db
DB_DEBUG=true
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10
DB_MAX_LIFETIME_CONNECTIONS=3600

# Redis (connect ke Docker Redis)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080       # Atau port sesuai service

# OBS (jika diperlukan)
OBS_HUAWEI_AK=your_access_key
OBS_HUAWEI_SK=your_secret_key
OBS_HUAWEI_ENDPOINT=https://obs.ap-southeast-1.myhuaweicloud.com
OBS_HUAWEI_BUCKET=your_bucket_name
```

### 3. **Development Lokal (go run) - Connect ke Local PostgreSQL**

**PERLU** update `.env` dengan konfigurasi local:

```env
# Database (local PostgreSQL)
DB_HOST=localhost
DB_PORT=5432           # Port PostgreSQL lokal
DB_USER=postgres
DB_PASS=your_local_password
DB_NAME=scylla_db
DB_DEBUG=true
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10
DB_MAX_LIFETIME_CONNECTIONS=3600

# Redis (local Redis)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

## 🔄 Perbandingan Setup

| Skenario | .env File | Environment Variables |
|----------|-----------|----------------------|
| **Docker Compose** | ❌ Tidak perlu | ✅ Di docker-compose.yml |
| **Local Dev → Docker DB** | ✅ Perlu (port 54321) | ❌ Tidak perlu |
| **Local Dev → Local DB** | ✅ Perlu (port 5432) | ❌ Tidak perlu |

## 📝 Template .env untuk Development Lokal

Buat file `.env` di setiap service directory dengan template berikut:

### Template untuk Connect ke Docker Services

```env
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database Configuration (Docker PostgreSQL)
DB_HOST=localhost
DB_PORT=54321
DB_USER=postgres
DB_PASS=postgres
DB_NAME=scylla_db
DB_DEBUG=true
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10
DB_MAX_LIFETIME_CONNECTIONS=3600

# Redis Configuration (Docker Redis)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Huawei OBS Configuration (optional)
OBS_HUAWEI_AK=your_access_key
OBS_HUAWEI_SK=your_secret_key
OBS_HUAWEI_ENDPOINT=https://obs.ap-southeast-1.myhuaweicloud.com
OBS_HUAWEI_BUCKET=your_bucket_name

# Service Specific
# Cronjob
INTERVAL_RELOAD_JOBS_IN_SECOND=30
```

### Template untuk Connect ke Local PostgreSQL

```env
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database Configuration (Local PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=your_local_password
DB_NAME=scylla_db
DB_DEBUG=true
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10
DB_MAX_LIFETIME_CONNECTIONS=3600

# Redis Configuration (Local Redis)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Huawei OBS Configuration (optional)
OBS_HUAWEI_AK=your_access_key
OBS_HUAWEI_SK=your_secret_key
OBS_HUAWEI_ENDPOINT=https://obs.ap-southeast-1.myhuaweicloud.com
OBS_HUAWEI_BUCKET=your_bucket_name
```

## 🎯 Quick Setup Script

Buat script untuk generate `.env` file untuk semua service:

```bash
# Generate .env untuk semua service (connect ke Docker)
for service in cronjob finance inventory master mobile sales system; do
  cat > "$service/.env" <<EOF
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=54321
DB_USER=postgres
DB_PASS=postgres
DB_NAME=scylla_db
DB_DEBUG=true
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10
DB_MAX_LIFETIME_CONNECTIONS=3600
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
EOF
done
```

## ⚠️ Important Notes

1. **Docker Compose**: Environment variables di docker-compose.yml **override** file `.env`
2. **Local Development**: Service membaca `.env` file dari direktori service-nya
3. **Port Conflict**: 
   - Docker PostgreSQL: port `54321` (untuk menghindari konflik)
   - Local PostgreSQL: port `5432` (default)
4. **Service Port**: Setiap service memiliki port berbeda (9001, 9002, dll) - pastikan tidak konflik

## 🔍 Verify Environment Variables

### Check dari dalam Docker Container

```bash
# Check environment variables di container
docker exec scylla-system env | grep DB_

# Check specific service
docker exec scylla-master env | grep SERVER_PORT
```

### Check dari Local Development

```bash
# Service akan load .env otomatis saat start
cd system
go run main.go
# Environment variables akan di-load dari system/.env
```

## 📚 Related Documentation

- [Docker Setup](./README.md)
- [Development Guidelines](../docs/DEVELOPMENT.md)
- [Database Scripts](../scripts/README.md)

---

**Last Updated**: 2024

