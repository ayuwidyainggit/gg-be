# Development Mode dengan Docker

Panduan untuk development dengan Docker menggunakan volume mounting untuk hot reload.

## 🚀 Quick Start

### Start Development Mode

```bash
# Start semua service dengan hot reload
./docker/start-dev.sh

# Atau manual
cd docker
docker-compose -f docker-compose.dev.yml up -d
```

## 🔥 Hot Reload dengan Air

### Cara Kerja

1. **Volume Mounting**: Source code di local di-mount ke container
2. **Air Hot Reload**: Menggunakan Air untuk auto-rebuild dan restart saat code berubah
3. **File Watching**: Air memantau perubahan file `.go` dan otomatis rebuild

### Setup Air (First Time)

```bash
# Setup Air configuration untuk semua service
./docker/setup-air.sh
```

Script ini akan membuat file `.air.toml` di setiap service directory.

### Perubahan Code

Saat Anda edit code di local:
1. File di local langsung terlihat di container (via volume mount)
2. **Air otomatis detect perubahan** dan rebuild
3. Service **otomatis restart** dengan code baru
4. Tidak perlu manual restart!

### Manual Restart Service (Jika Perlu)

```bash
# Restart specific service
docker-compose -f docker-compose.dev.yml restart system

# Restart semua service
docker-compose -f docker-compose.dev.yml restart
```

## 📝 Development Workflow

### 1. Setup Air (First Time Only)

```bash
# Setup Air configuration untuk semua service
./docker/setup-air.sh
```

### 2. Start Development Environment

```bash
./docker/start-dev.sh
```

### 3. Edit Code di Local

Edit code di direktori service (contoh: `system/controller/user_controller.go`)

**Air akan otomatis:**
- Detect perubahan file
- Rebuild binary
- Restart service
- Tidak perlu manual restart!

### 4. Test Changes

```bash
# Test endpoint
curl http://localhost:9001/ping
```

### 5. View Logs (Optional)

```bash
# View logs untuk melihat Air rebuild process
docker-compose -f docker-compose.dev.yml logs -f system
```

## 🔍 View Logs

```bash
# Logs semua service
docker-compose -f docker-compose.dev.yml logs -f

# Logs specific service
docker-compose -f docker-compose.dev.yml logs -f system

# Logs dengan tail
docker-compose -f docker-compose.dev.yml logs -f --tail=100 system
```

## 🛠️ Air Configuration

### Customize Air Behavior

Edit file `.air.toml` di service directory untuk customize:

```toml
# Contoh: Exclude directory tertentu
[build]
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]

# Contoh: Include file extension tertentu
  include_ext = ["go", "tpl", "tmpl", "html"]

# Contoh: Delay sebelum rebuild (ms)
  delay = 1000
```

### Air Logs

Air akan menampilkan log saat rebuild:
- Build errors akan ditampilkan di console
- Log file: `tmp/build-errors.log` (di dalam container)

### Troubleshooting Air

Jika Air tidak auto-reload:

```bash
# Check Air process
docker exec scylla-system-dev ps aux | grep air

# Check logs
docker-compose -f docker-compose.dev.yml logs -f system

# Restart service
docker-compose -f docker-compose.dev.yml restart system
```

#### Container Restart Terus

Jika container restart terus, kemungkinan:
1. **Air tidak terinstall**: Check logs untuk error install
2. **Go version mismatch**: Air v1.49.0 digunakan untuk kompatibilitas dengan Go 1.23.5
3. **Build error**: Check logs untuk compile errors

```bash
# Check error logs
docker-compose -f docker-compose.dev.yml logs system | grep -i error

# Check jika Air binary ada
docker exec scylla-system-dev ls -la /root/go/bin/air
```

## 📊 Perbandingan Mode

| Feature | Production Mode | Development Mode |
|---------|----------------|------------------|
| Dockerfile | Multi-stage build | Go image langsung |
| Binary | Compiled binary | Air hot reload |
| Volume Mount | ❌ Tidak ada | ✅ Source code mounted |
| Hot Reload | ❌ Tidak ada | ✅ **Auto restart dengan Air** |
| Build Time | Lambat (build image) | Cepat (incremental build) |
| Restart Speed | Lambat (rebuild) | **Otomatis saat code berubah** |

## 🎯 Best Practices

### 1. Development Workflow

```bash
# 1. Start development environment
./docker/start-dev.sh

# 2. Edit code di local
# (edit file di service directory)

# 3. Restart service
docker-compose -f docker-compose.dev.yml restart [service]

# 4. Check logs
docker-compose -f docker-compose.dev.yml logs -f [service]
```

### 2. Debugging

```bash
# Attach ke container untuk debugging
docker exec -it scylla-system-dev sh

# Di dalam container
cd /app
go run main.go
```

### 3. Testing Changes

```bash
# Test endpoint setelah perubahan
curl http://localhost:9001/ping

# Test dengan verbose
curl -v http://localhost:9001/v1/users
```

## ⚠️ Important Notes

1. **Volume Mounting**: Source code di-mount sebagai read-write, perubahan di local langsung terlihat
2. **Go Modules**: Go modules di-cache di volume terpisah untuk performa
3. **Restart Required**: Perlu restart container untuk load perubahan (tidak auto-reload)
4. **File Permissions**: Pastikan file permissions correct untuk volume mounting

## 🔧 Troubleshooting

### Service Tidak Restart

```bash
# Check container status
docker-compose -f docker-compose.dev.yml ps

# Check logs untuk error
docker-compose -f docker-compose.dev.yml logs system
```

### Code Changes Tidak Terlihat

```bash
# Verify volume mount
docker exec scylla-system-dev ls -la /app

# Check file timestamp
docker exec scylla-system-dev stat /app/main.go
```

### Go Module Issues

```bash
# Rebuild go modules
docker-compose -f docker-compose.dev.yml exec system go mod download
docker-compose -f docker-compose.dev.yml restart system
```

## 📚 Related Documentation

- [Docker Setup](./README.md)
- [Environment Variables](./ENV_SETUP.md)
- [Development Guidelines](../docs/DEVELOPMENT.md)

---

**Last Updated**: 2024

