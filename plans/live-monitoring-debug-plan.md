# Live Monitoring Fix Plan

## Scope

Issue yang akan diperbaiki ada pada endpoint berikut:

- [`GetDistributorMonitoring`](../pjp/controller/live_monitoring/get_distributor_controller.go:33)
- [`GetMonitoringDetail`](../pjp/controller/live_monitoring/get_detail_controller.go:28)

## Diagnosis Summary

### 1. List Monitoring

Akar masalah paling mungkin:

- Query repository sudah mengambil [`ovl.leave_at`](../pjp/repository/live_monitoring/get_distributor_repository.go) **belum** dilakukan di implementasi saat ini.
- Model row [`LiveMonitoringDistributorRow`](../pjp/model/live_monitoring.go:30) belum punya field `leave_at`.
- Response DTO [`LiveMonitoringDestinationData`](../pjp/data/response/live_monitoring_response.go:42) belum punya field `leave_at`.
- Mapper [`transformDistributorRows()`](../pjp/service/live_monitoring/get_distributor_service.go:87) belum meneruskan nilai tersebut ke response.

Dampak:

- Data `leave_at` ada di database, tetapi hilang di layer repository → model → response.

### 2. Detail Monitoring

Akar masalah paling mungkin:

- Query agregasi distributor di [`GetVisitInformationDistributor()`](../pjp/repository/live_monitoring/get_detail_repository.go:85) masih memakai pola lama berbasis join [`pjp.route_outlet`](../pjp/repository/live_monitoring/get_detail_repository.go:106) dan hitungan `start` / `finish`.
- Formula existing belum memisahkan:
  - planned non extra call
  - extra call
  - on going berdasarkan `arrive_at is not null and leave_at is null`
  - visited berdasarkan `arrive_at is not null and leave_at is not null`
  - skipped berdasarkan `skip_at is not null and leave_at is null`
- Perhitungan extra call di service [`getDistributorVisitInfo()`](../pjp/service/live_monitoring/get_detail_service.go:183) masih memakai `totalVisits - matched`, sehingga tidak sesuai definisi issue doc.

Dampak:

- Nilai `planned`, `extra_call`, `on_going`, `visited`, `skipped` tidak sesuai data nyata di `pjp.outlet_visit_list`.

## File Change Plan

### A. Files yang akan diubah untuk List Monitoring

1. [`pjp/repository/live_monitoring/get_distributor_repository.go`](../pjp/repository/live_monitoring/get_distributor_repository.go)
   - Tambahkan select `ovl.leave_at`
   - Pastikan alias konsisten dengan struct model

2. [`pjp/model/live_monitoring.go`](../pjp/model/live_monitoring.go:30)
   - Tambahkan field `LeaveAt *int64`

3. [`pjp/data/response/live_monitoring_response.go`](../pjp/data/response/live_monitoring_response.go:42)
   - Tambahkan field `LeaveAt *int64` pada response destination distributor

4. [`pjp/service/live_monitoring/get_distributor_service.go`](../pjp/service/live_monitoring/get_distributor_service.go:150)
   - Mapping `row.LeaveAt` ke response `leave_at`

### B. Files yang akan diubah untuk Detail Monitoring

1. [`pjp/model/live_monitoring.go`](../pjp/model/live_monitoring.go:61)
   - Ubah atau perluas struct [`VisitInformationRow`](../pjp/model/live_monitoring.go:61) agar bisa menampung kolom hasil agregasi baru, minimal `Plan`, `ExtraCall`, `OnGoing`, `Visited`, `TotalSkip`

2. [`pjp/repository/live_monitoring/get_detail_repository.go`](../pjp/repository/live_monitoring/get_detail_repository.go)
   - Refactor [`GetVisitInformationDistributor()`](../pjp/repository/live_monitoring/get_detail_repository.go:85) agar mengikuti query perbaikan issue doc
   - Kemungkinan tidak lagi bergantung pada [`CountTotalVisitsDistributor()`](../pjp/repository/live_monitoring/get_detail_repository.go:130) untuk menghitung extra call distributor
   - Jika perlu, sederhanakan atau hapus dependensi formula `matched`

3. [`pjp/service/live_monitoring/get_detail_service.go`](../pjp/service/live_monitoring/get_detail_service.go:183)
   - Sesuaikan [`getDistributorVisitInfo()`](../pjp/service/live_monitoring/get_detail_service.go:183) supaya memakai nilai `extra_call` langsung dari repository, bukan `totalVisits - matched`
   - Review apakah jalur principal perlu dibiarkan seperti sekarang atau disamakan polanya tanpa mengubah behavior principal yang existing

### C. Files opsional untuk validasi dokumentasi dan kontrak API

4. [`pjp/docs/swagger.json`](../pjp/docs/swagger.json)
   - Update hanya bila Swagger perlu disinkronkan secara eksplisit
   - Tidak wajib untuk fix runtime bila proyek belum regenerate docs saat ini

## Implementation Sequence

1. Tambahkan kontrak data `leave_at` pada model dan response
2. Perbarui query distributor list agar select `leave_at`
3. Perbarui mapper list response
4. Refactor query detail distributor agar sesuai formula issue
5. Sesuaikan service detail distributor agar tidak menghitung extra call dengan formula lama
6. Jalankan verifikasi curl untuk dua endpoint
7. Cocokkan hasil dengan data SQL remote yang sudah tervalidasi

## Expected Result After Fix

### List Monitoring

Untuk endpoint [`/api/v1/live-monitoring-distributor`](../pjp/router/router.go:95):

- `destination_data` akan mengandung `leave_at`
- Kasus BOEDY 8 akan menampilkan `leave_at = 1773820841188`
- Kasus BOEDY 6 akan menampilkan `leave_at = 1773821305595`

### Detail Monitoring

Untuk endpoint [`/api/v1/monitoring_locations/details`](../pjp/router/router.go:95):

nilai visit information distributor untuk kasus emp `360`, date `2026-03-18`, distributor `67` diharapkan menjadi:

- `planned = 2`
- `extra_call = 2`
- `on_going = 1`
- `visited = 2`
- `skipped = 0`

## Notes

- Tidak ada indikasi perlu dependency baru.
- Perubahan tetap berada di layer yang sesuai: repository untuk query, model/response untuk kontrak data, service untuk orchestration.
- KiloCode skill reusable yang relevan hanya memberi panduan umum backend; tidak ada utility KiloCode spesifik yang langsung menyelesaikan query live monitoring ini, jadi perubahan dilakukan dengan extend kode existing.
