# Scylla Live Monitoring API - Postman Collection

Collection ini berisi endpoint untuk testing Live Monitoring setelah penyesuaian dengan dokumentasi `Monitoring Activity - BE.md`.

## 📁 Files

1. **Scylla-Live-Monitoring.postman_collection.json** - Collection lengkap dengan request examples
2. **Scylla-Local-Environment.postman_environment.json** - Environment variables untuk local development

## 🚀 Cara Import ke Postman

### Step 1: Import Collection
1. Buka Postman
2. Klik **Import** button (pojok kiri atas)
3. Pilih file `Scylla-Live-Monitoring.postman_collection.json`
4. Klik **Import**

### Step 2: Import Environment
1. Klik **Environments** tab (pojok kiri)
2. Klik **Import** button
3. Pilih file `Scylla-Local-Environment.postman_environment.json`
4. Klik **Import**
5. Aktifkan environment dengan klik dropdown di pojok kanan atas Postman dan pilih "Scylla Local Environment"

## 📋 Endpoint yang Tersedia

### 1. Login (Get Token)
- **URL**: `POST http://localhost:9001/v1/users/login`
- **Body**:
  ```json
  {
    "email": "princ@idetama.id",
    "password": "admin"
  }
  ```
- **Response**: Token JWT yang digunakan untuk endpoint lainnya

### 2. Live Monitoring - Principal
- **URL**: `GET http://localhost:9010/api/v1/live-monitoring-principal`
- **Query Parameters**:
  - `date` (required): Epoch timestamp (contoh: 1769904000 = 2026-02-01)
  - `status[]` (required): Array status filter (contoh: Approved)
  - `emp_id` (optional): Comma-separated employee IDs filter (contoh: 221,279,264)
  - `region_id` (optional): Region ID filter
  - `area_id` (optional): Area ID filter
  - `distributor_id` (optional): Distributor ID filter

### 3. Live Monitoring - Distributor
- **URL**: `GET http://localhost:9010/api/v1/live-monitoring-distributor`
- **Query Parameters**: Sama seperti Principal

### 4. Monitoring Detail
- **URL**: `GET http://localhost:9010/api/v1/monitoring_locations/details`
- **Query Parameters**:
  - `emp_id` (required): Employee ID
  - `date` (required): Date in YYYY-MM-DD format (contoh: 2026-02-01)
  - `distributor_id` (optional): Distributor ID

## 🔧 Perubahan yang Sudah Diimplementasikan

### Principal Monitoring
- ✅ Response menggunakan `pjp_code` (bukan `pjp_id`)
- ✅ Field `skip_reason` ditambahkan
- ✅ LEFT JOIN untuk distributor (data tetap muncul walau distributor kosong)
- ✅ JOIN via `pjp_code` untuk visit list
- ✅ Filter area/region/distributor menjadi optional

### Distributor Monitoring
- ✅ Field `destination_type: "Outlet"` ditambahkan (hardcoded)
- ✅ JOIN via `route_code` untuk visit list
- ✅ Filter area/region/distributor menjadi optional

### Detail Monitoring
- ✅ Field `status` ditambahkan pada data Shipment

## 🧪 Testing dengan Data Dummy

Untuk testing dengan data lengkap, pastikan data dummy sudah diinsert ke database:

```sql
-- Data Principal
INSERT INTO pjp_principles.permanent_journey_plans (id, pjp_code, salesman_id, approval_status, cust_id) 
VALUES (8888, 8888, 9999, 'Approved', 'C22001');

-- Data Distributor
INSERT INTO pjp.permanent_journey_plans (id, pjp_code, salesman_id, approval_status, cust_id) 
VALUES (7777, 7777, 9999, 'Approved', 'C22001');

-- Lengkapi dengan route, destination, dan visit list (lihat query lengkap di repository)
```

## ⚠️ Troubleshooting

### Response "No Data"
Jika endpoint mengembalikan "No Data", kemungkinan penyebab:
1. Data dummy belum lengkap (cek tabel destinations, outlet_visit_list)
2. Filter tidak match dengan data (cek date format, status, cust_id)
3. JWT token expired (lakukan login ulang)

### HTTP 404 Not Found
Pastikan service PJP sudah berjalan:
```bash
docker-compose -f docker/docker-compose.dev.yml ps pjp
```

### HTTP 401 Unauthorized
Token JWT expired atau invalid. Lakukan login ulang di endpoint "1. Login".

## 📚 Response Structure

### Principal/Distributor Response
```json
{
  "message": "Success",
  "data": [
    {
      "emp_id": 9999,
      "emp_code": "SLS9999",
      "emp_name": "Salesman Test",
      "distributor_id": 999,
      "area_id": 1,
      "region_id": 1,
      "pjp_data": [
        {
          "pjp_code": 8888,
          "approval_status": "Approved",
          "route_data": [
            {
              "route_code": "8888",
              "route_name": "Route Test",
              "destination_data": [
                {
                  "destination_id": 1001,
                  "destination_code": "DEST001",
                  "destination_type": "Outlet",
                  "destination_name": "Toko Test",
                  "longitude": -6.17511,
                  "latitude": 106.865036,
                  "arrive_at": null,
                  "start": 1769904000,
                  "finish": 1769907600,
                  "skip_at": null,
                  "skip_reason": "Toko Tutup"
                }
              ]
            }
          ]
        }
      ]
    }
  ],
  "paging": {
    "total_record": 1,
    "page_current": 1,
    "page_limit": 10,
    "page_total": 1
  },
  "request_id": "xxx"
}
```

### Detail Response
```json
{
  "message": "Success",
  "data": [
    {
      "visit_information": {
        "activity_date": "2026-02-01",
        "company_name": "PT Sukses Makmur",
        "company_code": "SM001",
        "level": "Distributor",
        "emp_id": 9999,
        "emp_code": "SLS9999",
        "emp_name": "Salesman Test",
        "activity_time": "2026-02-01 08:00:00",
        "planned": 1,
        "on_going": 0,
        "extra_call": 0,
        "visited": 1,
        "skipped": 0
      },
      "sales": [...],
      "return": [...],
      "collection": [...],
      "expense": [...],
      "shipment": [
        {
          "shipment_no": "SHP001",
          "status": "Delivered",
          "shipment_data": [...]
        }
      ]
    }
  ],
  "request_id": "xxx"
}
```

## 📞 Support
Jika ada masalah atau pertanyaan, silakan hubungi tim development.