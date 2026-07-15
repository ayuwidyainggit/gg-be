# Evidensi Discovery — SX-1958

## Ringkasan Temuan
- Endpoint list outlet berada di service `master` pada route `GET /v1/outlets`; prefix `/master` sangat mungkin ditambahkan di gateway/proxy.
- Query `is_active` **bukan** default backend pada endpoint ini. Filter hanya diterapkan bila request mengirim `is_active` dan field `OutletQueryFilter.IsActive` terisi.
- Flow PJP yang mengambil outlet by salesman/PJP di service `pjp` saat ini memanggil endpoint master dengan `outlet_id=<ids>&limit=9999` tanpa `include_inactive` dan tanpa `verification_status`.
- Karena user evidence menunjukkan request aktual Add New Order mengandung `is_active=1&verification_status=1&outlet_id=...`, ada kemungkinan caller lain (FE/gateway/service lain) masih menambahkan `is_active=1`. Perubahan backend yang paling aman adalah menambah flag eksplisit untuk bypass filter aktif, bukan mengubah contract global `outlet_id`.

## File yang Diinspeksi
- `master/controller/outlet_controller.go`
- `master/entity/outlet.go`
- `master/repository/outlet_repository.go`
- `master/service/outlet_service.go`
- `master/controller/outlet_controller_test.go`
- `master/controller/query_filter_parser_test.go`
- `pjp/service/third_party/get_outlet_by_sales_codes_service.go`
- `pjp/service/third_party/get_outlet_picklist_by_sales_codes_service.go`
- `pjp/service/third_party/get_outlet_service.go`
- `pjp/controller/third_party/get_outlet_by_sales_codes_controller.go`
- `pjp/router/third_party.go`
- `pjp/service/visit_service.go` (hasil grep)

## Pattern / Reuse yang Ditemukan
- Parsing query list integer reuse helper existing:
  - `parseIntSliceQuery(...)`
  - `parseIntSliceQueryAllowZero(...)`
- `OutletQueryFilter` sudah menjadi contract tunggal untuk controller → service → repository list outlet.
- Repository sudah memiliki pola filter opsional via pointer/int:
  - `if dataFilter.IsActive != nil { ... }`
- PJP already has dedicated caller services for salesman/PJP outlet retrieval:
  - `GetOutletBySalesCodes`
  - `GetOutletPicklistBySalesCodes`
- Existing test pattern tersedia untuk parser helper di controller tests.

## Trace Parameter

### Route / Controller
- `master/controller/outlet_controller.go:33-42`
  - route group `/v1/outlets`
  - list handler `controller.List`
- `master/controller/outlet_controller.go:111-233`
  - `c.QueryParser(&dataFilter)` mengisi field query sederhana seperti `page`, `limit`, `sort`, `is_active`
  - `verification_status`, `outlet_id`, `ot_class_id`, `ot_type_id`, `ot_grp_id`, `distributor_id`, `outlet_status` diparse eksplisit dari query args
  - default pagination: `page=1`, `limit=10`

### Entity Contract
- `master/entity/outlet.go:804-826`
  - `IsActive *int     \`query:"is_active"\``
  - `VerificationStatus []int`
  - `OutletID []int`
  - belum ada field `include_inactive`

### Service
- `master/service/outlet_service.go:1227-1232`
  - list service hanya meneruskan filter ke repository setelah scope distributor/customer resolved

### Repository / SQL Builder
- `master/repository/outlet_repository.go:926-934`
  - hanya menambah `AND o.is_active = true/false` jika `dataFilter.IsActive != nil`
- `master/repository/outlet_repository.go:942-946`
  - `outlet_id`, `verification_status`, `ot_class_id`, `ot_grp_id`, `ot_type_id` ditambahkan dengan helper filter IN

## Caller / Aggregator yang Relevan
- `pjp/controller/third_party/get_outlet_by_sales_codes_controller.go:27-69`
  - endpoint PJP `GET /outlets/salesman`
- `pjp/router/third_party.go:20-21`
  - route `/outlets/salesman` dan `/outlets-picklist/salesman`
- `pjp/service/third_party/get_outlet_by_sales_codes_service.go:62-64`
  - membangun URL `.../v1/outlets?outlet_id=%s&limit=9999`
- `pjp/service/third_party/get_outlet_picklist_by_sales_codes_service.go:62-64`
  - membangun URL yang sama

## Keputusan Awal Kontrak
- Kontrak endpoint umum `/v1/outlets` tetap backward-compatible.
- Tambahkan flag eksplisit `include_inactive=1` di backend master.
- Saat `include_inactive=1`, backend harus **mengabaikan** filter `is_active`, termasuk bila caller masih mengirim `is_active=1`.
- Flow PJP/salesman caller di service `pjp` perlu ikut mengirim `include_inactive=1` agar intent eksplisit dan aman terhadap regression.

## Kandidat File Perubahan Implementasi
- `master/entity/outlet.go`
- `master/controller/outlet_controller.go`
- `master/repository/outlet_repository.go`
- `master/controller/outlet_controller_test.go`
- `pjp/service/third_party/get_outlet_by_sales_codes_service.go`
- `pjp/service/third_party/get_outlet_picklist_by_sales_codes_service.go`
- Mungkin tambahan test di `pjp` bila test pattern caller tersedia.

## TDD / Test Reuse
- Existing parser tests:
  - `master/controller/query_filter_parser_test.go`
  - `master/controller/outlet_controller_test.go`
- Untuk repository/service list outlet belum terlihat stub khusus pada discovery cepat; kemungkinan lebih murah membuat unit test controller/repository helper daripada integration DB test penuh bila harness test tidak tersedia.

## Constraint
- Repo root bukan git repo menurut environment metadata, jadi verifikasi git style/history lokal tidak tersedia via default safety flow.
- Instruksi repo mewajibkan penggunaan prefix `rtk`, tetapi tool shell policy global melarang prefix tersebut. Pada planning artifact, command examples mengikuti konvensi repo (`rtk ...`) agar implementer konsisten.
- `.opencode/docs/AGENT_ROUTING.md` dan `.opencode/docs/SKILLS.md` tidak ditemukan di repo ini, sehingga discovery mengandalkan `AGENTS.md` dan codebase executable facts.
- Data staging/Jira/cURL sensitif tidak boleh dihardcode ke code/test.

## Risk
- Jika hanya caller PJP diubah, flow lain yang masih memanggil `/master/v1/outlets?...is_active=1...` akan tetap salah untuk use case PJP Add New Order.
- Jika contract diubah menjadi “selalu abaikan `is_active` saat `outlet_id` ada”, screen/list lain yang memang butuh active-only by selected IDs bisa regression.
- Belum ada bukti lokal yang memetakan secara pasti request web Add New Order ke route PJP vs caller lain; perlu manual verification staging setelah implementasi.

## Sumber / Command yang Dipakai
- Docker service check via `docker compose -f docker-compose.yml ps`
- Read/Grep lokal pada file-file master dan pjp di atas
- Read-only explorer subagent untuk trace awal route/controller/service/repository
