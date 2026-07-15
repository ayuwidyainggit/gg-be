# 📮 Postman Collection - Outlet Management Enhanced

Collection ini berisi semua endpoint untuk testing enhancement fitur Outlet Management sesuai dokumentasi `docs/Add Outlet Manage Outlet_BE.md`.

## 📁 Files

- **Outlet_Management_Enhanced.postman_collection.json** - Postman Collection dengan semua endpoint
- **Outlet_Management_Enhanced.postman_environment.json** - Environment variables untuk local development

## 🚀 Cara Import ke Postman

### 1. Import Collection
1. Buka Postman
2. Klik **Import** di pojok kiri atas
3. Pilih file `Outlet_Management_Enhanced.postman_collection.json`
4. Klik **Import**

### 2. Import Environment
1. Klik icon gear (⚙️) di pojok kanan atas
2. Klik **Import**
3. Pilih file `Outlet_Management_Enhanced.postman_environment.json`
4. Klik **Import**
5. Pilih environment "Outlet Management - Local Development" dari dropdown

## 📋 Endpoint yang Tersedia

### 🔐 Authentication
1. **Login** - `POST /v1/users/login`
   - Otomatis menyimpan access_token ke environment variable
   - Credential default: `charles@gmail.com` / `admin`

### 🏪 Outlet Management
2. **Upload File** - `POST /v1/files/uploads`
   - Upload foto outlet
   - Response berisi URL untuk field `file_url`

3. **Create Outlet (With Photo)** - `POST /v1/m-outlets`
   - Create outlet dengan field `file_url`
   - Field `file_url` berisi URL hasil upload

4. **Create Outlet (Without Photo)** - `POST /v1/m-outlets`
   - Create outlet tanpa `file_url`
   - Test backward compatibility

5. **Get Outlet List (All Status)** - `GET /v1/outlet-list`
   - List outlet dengan multiple status (0=pending, 1=active)
   - Support pagination, search, sorting

6. **Get Outlet List (Pending Only)** - `GET /v1/outlet-list`
   - Filter hanya outlet pending (status=0)

7. **Get Outlet List (Active Only)** - `GET /v1/outlet-list`
   - Filter hanya outlet active (status=1)

8. **Get Outlet List (With Search)** - `GET /v1/outlet-list`
   - Search by outlet_name atau outlet_code

9. **Delete Outlet** - `DELETE /v1/outlet-list/:outlet_id`
   - Soft delete outlet
   - Update `is_del=true`, `deleted_by`, `deleted_at`

## 🔧 Environment Variables

Collection menggunakan environment variables berikut:

| Variable | Default Value | Description |
|----------|--------------|-------------|
| `base_url` | `http://localhost:9008` | Base URL service mobile |
| `access_token` | (auto-set) | JWT token dari login |
| `test_email` | `charles@gmail.com` | Email untuk testing |
| `test_password` | `admin` | Password untuk testing |
| `user_id` | `120` | User ID untuk testing |

## 📝 Cara Menggunakan

### Step 1: Login
1. Jalankan request **Login** di folder Authentication
2. Token akan otomatis tersimpan di environment variable `access_token`
3. Semua request berikutnya akan otomatis menggunakan token ini

### Step 2: Upload File (Optional)
1. Jalankan request **Upload File**
2. Pilih file foto outlet
3. Copy URL dari response
4. Gunakan URL tersebut di field `file_url` saat create outlet

### Step 3: Create Outlet
1. Jalankan request **Create Outlet (With Photo)** atau **Create Outlet (Without Photo)**
2. Edit body request sesuai kebutuhan
3. Response akan berisi pesan "Berhasil Ditambahkan"

### Step 4: Get Outlet List
1. Jalankan salah satu request **Get Outlet List**
2. Sesuaikan query parameters:
   - `outlet_status`: Filter by status (bisa multiple)
   - `q`: Search keyword
   - `page`: Halaman
   - `limit`: Jumlah data per halaman
   - `sort`: Sorting (format: `field:asc` atau `field:desc`)

### Step 5: Delete Outlet
1. Jalankan request **Delete Outlet**
2. Ganti `:outlet_id` dengan ID outlet yang ingin dihapus
3. Outlet akan di-soft delete (tidak muncul di list)

## 🎯 Query Parameters Outlet List

### outlet_status (Array)
- `0` = Pending
- `1` = Active
- Bisa multiple: `?outlet_status=0&outlet_status=1`

### q (String)
- Search by `outlet_name` atau `outlet_code`
- Example: `?q=toko`

### page (Integer)
- Default: `1`
- Example: `?page=2`

### limit (Integer)
- Default: `9999`
- Example: `?limit=10`

### sort (String)
- Format: `field:direction`
- Default: `outlet_code:asc`
- Example: `?sort=created_date:desc`

## ✅ Testing Checklist

- [ ] Login berhasil dan token tersimpan
- [ ] Upload file berhasil dan dapat URL
- [ ] Create outlet dengan file_url berhasil
- [ ] Create outlet tanpa file_url berhasil (backward compatibility)
- [ ] Get list outlet dengan multiple status
- [ ] Get list outlet pending only
- [ ] Get list outlet active only
- [ ] Search outlet by name/code
- [ ] Delete outlet berhasil
- [ ] Outlet yang dihapus tidak muncul di list

## 🐛 Troubleshooting

### Token Expired
- Jalankan ulang request **Login** untuk mendapatkan token baru

### 401 Unauthorized
- Pastikan environment "Outlet Management - Local Development" sudah dipilih
- Pastikan sudah login dan token tersimpan

### Service Not Running
- Pastikan Docker container `scylla-mobile-dev` sudah running
- Check dengan: `docker ps | grep mobile`
- Service harus running di port 9008

## 📚 Reference

- Dokumentasi: `docs/Add Outlet Manage Outlet_BE.md`
- Implementation Summary: `docs/IMPLEMENTATION_SUMMARY.md`
- Testing Guide: `plans/testing-guide-curl.md`

## 🎉 Features

✅ Auto-save access token setelah login  
✅ Pre-configured request bodies  
✅ Multiple test scenarios  
✅ Query parameter examples  
✅ Environment variables untuk easy configuration  
✅ Backward compatibility testing  

---

**Created for:** Scylla Mobile Service - Outlet Management Enhancement  
**Last Updated:** 2025-01-08