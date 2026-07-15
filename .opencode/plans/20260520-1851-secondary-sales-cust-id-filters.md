# Plan — Task 120/121 Secondary Sales `cust_id` Filters

Task id: `20260520-1851-secondary-sales-cust-id-filters`
Tanggal: `2026-05-20`
Target service: `sales`
Primary source of truth: `.opencode/plans/20260520-1851-secondary-sales-cust-id-filters.md`

## Goal

Tambahkan filter business-unit `cust_id` ke dua endpoint Secondary Sales Report di service `sales`:

- Task 120 — `POST /sales/v1/reports/secondary-sales` (Export). Body baru menerima `cust_id` opsional.
- Task 121 — `GET /sales/v1/reports/secondary-sales/trend-sales` (Trend Sales). Body JSON baru menerima `cust_id` opsional.

`cust_id` dipakai untuk memfilter data pada `sls.order/sls.order_detail` (export) dan `report.fact_orders` (trend), dengan scope rule:

- Empty atau equal auth → fallback ke auth `cust_id`.
- Principal user (`authCustID == parentCustID`) boleh kirim child `cust_id` dibawah `parent_cust_id`.
- Distributor user (`authCustID != parentCustID`) tidak boleh kirim `cust_id` selain miliknya.

## Non-goals

- Tidak mengubah route, method, atau response shape.
- Tidak mengubah `year` behavior pada Trend Sales (tetap `required` sesuai validator existing dan jawaban user "ikuti existing").
- Tidak mengubah field `report.list.cust_id`: tetap auth user yang login (sesuai keputusan user — supaya principal yang export untuk child distributor tetap melihat baris di `GET /v1/reports`).
- Tidak mengubah parent product LATERAL join di `buildSecondarySalesUnionQuery` (tetap pakai `dataFilter.ParentCustID`, bukan effective cust).
- Tidak menambah filter `region`/`area`/`distributor` (di luar dua task ini).
- Tidak meng-update `docs/Secondary Sales Report_BE.md` kecuali user minta.

## Scope

Endpoint:

- `POST /v1/reports/secondary-sales` (Task 120, Export).
- `GET /v1/reports/secondary-sales/trend-sales` (Task 121, Trend Sales).

Layer target:

- `sales/entity/report.go`
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go` (penyesuaian filter Union supaya pakai `RequestedCustID`/effective cust untuk WHERE; tidak ubah parent product join)
- `sales/service/report_service_test.go`
- `sales/controller/so_controller_test.go`
- `sales/repository/report_repository_test.go`

## Docs Reference Alignment

- `docs/Secondary Sales Report_BE.md` lines 22-50 (Trend Sales) meminta tambah body `cust_id` (filter `report.fact_orders.cust_id`). User memilih bind via JSON body.
- `docs/Secondary Sales Report_BE.md` lines 119-141 (Export Secondary Sales) meminta tambah body `cust_id`.
- Keputusan user:
  - Trend Sales: `cust_id` di body JSON, `year` tetap `required` mengikuti existing.
  - Export: `report.list.cust_id` = auth user.

## Requirements

Common:

- `cust_id` opsional di body. Kosong → fallback ke auth `cust_id`.
- Principal (auth == parent) boleh kirim child `cust_id`. Repo cek via `ExistsCustomerInParentScope`.
- Distributor (auth != parent) hanya boleh kirim `cust_id` yang sama dengan auth atau kosong.
- Cust scope unauthorized → `403 Forbidden` dengan body response standar.

Task 120 (Export):

- Field body baru: `cust_id` opsional.
- Effective cust dipakai sebagai filter WHERE pada query Union (`od.cust_id`, `rd.cust_id`).
- `report.list.cust_id` tetap = auth user.
- `Subscribe...` (RMQ consumer) memakai effective cust dari payload yang dipublish.

Task 121 (Trend Sales):

- Bind body JSON baru berisi `cust_id` opsional. `year` tetap dari query param `?year=...` (existing, `validate:"required"`).
- Effective cust dipakai pada `repository.SecondarySalesReportTrendSales`.
- Response tetap 12 baris bulan dengan zero-fill (tidak diubah).

## Acceptance Criteria

1. Export menerima body `cust_id` dan memfilter row export sesuai effective cust.
2. Export tetap menulis `report.list.cust_id = auth user` agar principal melihat baris export di `GET /v1/reports`.
3. Trend Sales menerima body JSON `{"cust_id": "..."}` dan memfilter `report.fact_orders.cust_id` sesuai effective cust.
4. Trend Sales tetap `year` required dari query param.
5. Distributor user yang mengirim `cust_id` selain miliknya → `403 Forbidden`, tidak query data.
6. Principal user yang mengirim child cust valid → query memakai child cust.
7. Empty `cust_id` di body → query memakai auth cust.
8. Invalid format `cust_id` (mis. spasi, panjang berlebihan) → `400 Bad Request`.
9. Response shape tidak berubah.
10. Unit dan integration-style test menutup: empty fallback, equal auth, principal child valid, principal child invalid, distributor sibling.

## Existing Patterns/Reuse

- `service.resolveSecondaryDashboardCustID(authCustID, parentCustID, requestedCustID)` (sales/service/report_service.go:1163).
- `repository.ExistsCustomerInParentScope` (sales/repository/report_repository.go:68).
- `service.ErrUnauthorizedCustID` + mapping ke `fiber.StatusForbidden` (lihat `SecondaryReportSalesSumMonth` line 357-360).
- Pattern dual-bind body+query Fiber: `c.QueryParser(&request)` lalu `c.BodyParser(&request)` pada struct yang sama (struct dengan tag `query:` dan `json:`).
- `SecondarySalesReportQueryFilter` sudah punya `CustID/ParentCustID` sebagai auth/scope. Tambahkan field baru `RequestedCustID string \`json:"cust_id"\`` supaya tidak bertabrakan dengan auth `CustID` yang di-overwrite handler.
- RMQ payload `dataFilter` di-serialize via `structs.StructToJson(dataFilter)` (line 348 service); cukup pastikan field baru di-export agar otomatis ikut serialized.
- Test pattern: mock repository hooks (`mockReportRepositoryForService`) dan controller mock service (`mockReportServiceForController`).

## Constraints

- Layering Controller→Service→Repository ketat. Controller tidak query DB.
- Tag JSON `cust_id` sudah dipakai oleh `SecondarySalesReportQueryFilter.CustID` (auth). Untuk hindari konflik, body baru dipakai field terpisah dengan dummy alias atau pakai DTO baru. Plan memilih: pisahkan ke DTO controller-only `SecondarySalesExportRequest` lalu copy ke `SecondarySalesReportQueryFilter.RequestedCustID` (baru).
- Validasi `rtk go test ./...` dijalankan dari `cd sales`.
- Tetap pakai `rtk` prefix untuk command shell sesuai `AGENTS.md` repo.
- Jangan commit `.env`/credentials.
- Validasi DB dev wajib read-only dan memakai secret lewat environment variable, bukan hardcoded di file, command history, test fixture, atau artefak. Koneksi target: host `103.28.219.73`, port `25431`, db `scylla_citus_dev`, user `postgres`, `sslmode=disable`; password disimpan out-of-band.

## Risks

- Cross-tenant leak bila scope check terlewat. Mitigasi: sentralisasi via `resolveSecondaryDashboardCustID` di service.
- Bila body Trend Sales tidak terkirim (banyak HTTP client/CDN drop body di GET), `cust_id` jadi empty dan fallback ke auth — perilaku ini aman tapi perlu QA note di bagian Validation.
- Union query saat ini juga menulis `dataFilter.CustID` ke parent product LATERAL join (`pp.cust_id`). Karena `pp.cust_id` adalah filter parent-product (bukan transactional), tetap pakai auth/parent — tidak diganti effective cust.
- `report.list.cust_id` tetap auth — implikasi: row hasil export dari principal untuk child tetap muncul di list principal, tidak muncul di list child. Sudah dikonfirmasi user.
- Penambahan field `RequestedCustID` ke `SecondarySalesReportQueryFilter` ikut ter-serialize ke RMQ payload; pastikan consumer subscribe (`SubscribeSecondarySalesReport`) memakai field yang sama saat memanggil repository Union.

## Decisions/Assumptions

Decisions (dari user):

- Trend Sales `cust_id` → JSON body.
- Trend Sales `year` → tetap required (existing).
- Owner `report.list` row → auth user (untuk Export Task 120).

Assumptions (low-risk):

- Validasi `cust_id` panjang reasonable: `omitempty,alphanum,max=20` selaras pattern Task 118/119.
- `c.BodyParser` di Fiber v2 untuk GET tetap mem-parse body bila Content-Type `application/json` dan body dikirim. Jika body kosong → fallback ke auth.
- Field baru `RequestedCustID` ditambahkan ke `SecondarySalesReportQueryFilter` dan tetap nullable (string kosong = unset).

Open questions: tidak ada (semua dijawab di question gate).

## TDD/Test Plan

TDD required: ya. Alasan: ini menyangkut tenant isolation, API behavior, dan production query.

Existing test patterns:

- `sales/controller/so_controller_test.go` menggunakan `httptest.NewRequest` + Fiber `app.Test`.
- `sales/service/report_service_test.go` memakai mock repository hook.
- `sales/repository/report_repository_test.go` memakai SQL string assertion via `buildSecondarySalesUnionQuery`.

First failing/regression tests:

1. **Service Export — principal child cust valid**: `PublishSecondarySalesReport` dengan `RequestedCustID = "CHILD-1"`, auth = parent. Repository `ExistsCustomerInParentScope` di-mock true. Assert: filter union dipanggil dengan effective cust `CHILD-1`, tetapi `reportList.CustID` yang disimpan = auth.
2. **Service Export — distributor sibling cust**: auth dist, request sibling cust → return `ErrUnauthorizedCustID`, tidak panggil `StoreReportList`.
3. **Service Export — empty fallback**: `RequestedCustID = ""`, scope DB tidak dipanggil, effective cust = auth.
4. **Subscribe Export query**: `SubscribeSecondarySalesReport` memakai effective cust di `dataFilter.CustID` saat memanggil `SecondarySalesUnion`. Verify via assertion field di mock repository.
5. **Service Trend — principal child valid**: `SecondarySalesReportTrendSales(authCustID, parentCustID, year, requestedCustID)` dipanggil dengan child cust valid → repo dipanggil dengan child.
6. **Service Trend — distributor sibling**: auth dist, request sibling → `ErrUnauthorizedCustID`.
7. **Service Trend — empty fallback**: cust kosong, effective = auth.
8. **Controller Export**: POST body `{..., "cust_id": "CHILD-1"}` → service menerima `RequestedCustID` benar dan auth/parent/exportBy diisi.
9. **Controller Trend — body cust_id**: GET `?year=2026` dengan body JSON `{"cust_id": "SIBLING-1"}` untuk distributor → 403.
10. **Controller Trend — query year required**: GET tanpa `year` → 400/422 (perilaku validator existing tetap).
11. **Repository union**: `buildSecondarySalesUnionQuery` ketika `dataFilter.CustID` adalah effective cust tetap membentuk SQL/parameter sesuai existing assertion.

Green step: kembangkan struktur DTO + service helper + controller routing sampai semua test lulus.

Refactor step: kalau ada duplikasi binding cust di service, tarik ke helper (`resolveSecondaryDashboardCustID` sudah cukup; tidak perlu helper baru).

Edge cases:

- `cust_id` kosong/whitespace.
- `cust_id` panjang berlebihan / invalid charset → ditolak validator.
- Body kosong di Trend Sales (Content-Length 0) → equivalent empty cust.
- RMQ payload effect: `RequestedCustID` ikut ter-publish dan `Subscribe...` memakai field benar.

Commands:

```bash
rtk go mod download && rtk go mod tidy
rtk go test ./entity/... ./controller/... ./service/... ./repository/...
rtk go test ./service -run 'TestSecondarySalesReport(Trend|Publish)|TestPublishSecondarySalesReport|TestSubscribeSecondarySalesReport'
rtk go test ./controller -run 'TestSecondarySales|TestTrendSales'
rtk go test ./repository -run 'TestBuildSecondarySalesUnionQuery'
rtk go test ./...
```

## Implementation Steps

1. **Entity** — tambah field `RequestedCustID` ke `SecondarySalesReportQueryFilter` dan struct payload Trend Sales:

   ```go
   // sales/entity/report.go
   type SecondarySalesReportQueryFilter struct {
       CustID         string `json:"-"` // auth dari controller, JSON tag dimatikan agar tidak terisi dari body
       ParentCustID   string `json:"-"`
       RequestedCustID string `json:"cust_id" validate:"omitempty,alphanum,max=20"`
       From           *int64 `json:"from" validate:"required_with=To,omitempty,gte=1000000000"`
       To             *int64 `json:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
       // ...field existing tetap...
   }

   type SecondarySalesReportTrensSalesSumPayload struct {
       Year   int    `query:"year" validate:"required"`
       CustID string `json:"cust_id" validate:"omitempty,alphanum,max=20"`
   }
   ```

   Catatan: ubah JSON tag `CustID/ParentCustID` ke `json:"-"` supaya request body tidak bisa men-spoof auth. Jalankan `grep` di service untuk memastikan tidak ada caller yang serialize fitur tersebut sebagai JSON ke client.

2. **Service helper** — sudah ada `resolveSecondaryDashboardCustID` dan `ErrUnauthorizedCustID`. Reuse.

3. **Service Export `PublishSecondarySalesReport`**:
   - Di awal fungsi, panggil `effectiveCustID, err := service.resolveSecondaryDashboardCustID(dataFilter.CustID, dataFilter.ParentCustID, dataFilter.RequestedCustID)`. Jika err, return.
   - Set `dataFilter.RequestedCustID = effectiveCustID` (atau `dataFilter.CustID = effectiveCustID` setelah backup auth ke variabel `authCustID`).
   - Tetap `reportList.CustID = authCustID` (auth, BUKAN effective). Simpan `authCustID` di variabel sebelum dataFilter dimodifikasi.
   - Kirim `dataFilter` (sudah mengandung effective cust) ke RMQ.

4. **Service Subscribe `SubscribeSecondarySalesReport`**:
   - Tidak ada perubahan logic; karena `dataFilter.CustID` sekarang effective dari publish, repository union otomatis filter sesuai effective cust.
   - Tambah safety: jika `RequestedCustID` ada, paksa `dataFilter.CustID = RequestedCustID` di awal Subscribe untuk konsumen lama yang publish belum effective.

5. **Service Trend Sales**:
   - Ubah signature `SecondarySalesReportTrendSales(authCustID, parentCustID string, year int, requestedCustID string)` (atau pass payload struct). Update interface `ReportService`.
   - Panggil `resolveSecondaryDashboardCustID`. Jika err, return.
   - Lanjut panggil `repository.SecondarySalesReportTrendSales(effectiveCustID, year)`.

6. **Controller Export `SecondarySales`**:
   - Tetap `c.BodyParser(&request)` ke `entity.SecondarySalesReportQueryFilter`. Karena JSON tag `CustID/ParentCustID` di-set `json:"-"`, body tidak bisa men-spoof auth.
   - Setelah BodyParser:
     - simpan `authCustID := c.Locals("cust_id").(string)` dan `parentCustID := c.Locals("parent_cust_id").(string)`,
     - `request.CustID = authCustID; request.ParentCustID = parentCustID`,
     - `request.ExportBy = c.Locals("user_fullname").(string)`.
   - Tambah validator: `controller.validator.ValidateStruct(request, headerAcceptLang)` jika belum ada (saat ini handler ini belum memvalidasi struct; tambahkan untuk tangani format `cust_id`).
   - Service yang return `ErrUnauthorizedCustID` di-mapping ke `fiber.StatusForbidden`.

7. **Controller Trend Sales `SecondaryReportSalesTrendSales`**:
   - Bind: `c.QueryParser(&request)` lalu `c.BodyParser(&request)` (abaikan error body kalau Content-Length 0). Validasi struct setelahnya.
   - Panggil service: `SecondarySalesReportTrendSales(authCustID, parentCustID, request.Year, request.CustID)`.
   - Map `ErrUnauthorizedCustID` → 403.

8. **Repository**:
   - `SecondarySalesUnion`/`SecondarySalesUnionPagination` tidak butuh ubah signature — mereka memakai `dataFilter.CustID` untuk filter transactional dan `dataFilter.ParentCustID` untuk product join.
   - `SecondarySalesReportTrendSales(custID, year)` tidak ubah signature. Service yang resolve cust.

9. **Tests**:
   - Service Export: tambah test publish dengan principal child + distributor sibling + empty fallback. Verifikasi store payload (`reportList.CustID = auth`) dan publish payload (`dataFilter.CustID = effective`).
   - Service Trend Sales: parameter signature baru, test scope.
   - Controller Export: test 403 dengan dist user request sibling.
   - Controller Trend Sales: test body JSON `{"cust_id":"..."}` di GET; pastikan dual-bind bekerja.
   - Repository: pastikan tes existing `TestBuildSecondarySalesUnionQuery*` masih hijau (tidak ada perubahan kontrak).

10. **Validation**:
    - `cd sales && rtk go mod tidy && rtk go test ./...`
    - Manual smoke (post-deploy): scenario di Done Criteria.

## Expected Files to Change

- `sales/entity/report.go`
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/service/report_service_test.go`
- `sales/controller/so_controller_test.go`
- `sales/repository/report_repository_test.go` (selaras kontrak existing; verifikasi tidak regress)
- (Opsional) `sales/repository/report_repository.go` — hanya bila penambahan helper diperlukan; default: tidak.

Tidak diubah:

- `docs/Secondary Sales Report_BE.md` (tunggu permintaan eksplisit).
- File compose, Makefile, env, atau lockfile.

## Agent/Tool Routing

- Implementasi: `@fixer` (bounded code edits + test).
- Review keamanan tenant + final signoff: `@quality-gate`.
- Discovery tambahan jika ada pola serupa: `@explorer`.
- Validasi library/Fiber GET-with-body bila ragu: `@librarian`.

## Execution-ready Worklist / Handoff Contract

Format:

`id | depends_on | action | owner | validation | exit | blocking | requires_user_decision`

```
T01 | none | Tambahkan field `RequestedCustID` ke `SecondarySalesReportQueryFilter` dan ubah JSON tag `CustID/ParentCustID` ke `json:"-"`; tambahkan `CustID` ke `SecondarySalesReportTrensSalesSumPayload` (`json:"cust_id"`). | @fixer | `cd sales && rtk go build ./...` | Compile clean | ready | no
T02 | T01 | Update interface `ReportService.SecondarySalesReportTrendSales` agar terima `(authCustID, parentCustID string, year int, requestedCustID string)`; update implementasi service untuk panggil `resolveSecondaryDashboardCustID` lalu repo. Update mock di test untuk match signature baru. | @fixer | `cd sales && rtk go build ./...` | Compile clean, mock konsisten | ready | no
T03 | T01 | Modifikasi `PublishSecondarySalesReport` — panggil resolver, set `dataFilter.CustID = effective`, tetap `reportList.CustID = authCustID`. | @fixer | `cd sales && rtk go test ./service -run TestPublishSecondarySalesReport` | Test untuk principal child, distributor sibling, empty fallback hijau | ready | no
T04 | T03 | Tambah safety di `SubscribeSecondarySalesReport`: bila `RequestedCustID != ""` dan beda dengan `CustID`, paksa `dataFilter.CustID = RequestedCustID` (idempotent). | @fixer | `cd sales && rtk go test ./service -run TestSubscribeSecondarySalesReport` | Subscribe pakai effective cust untuk Union | ready | no
T05 | T01 | Update controller `SecondarySales`: tambah `ValidateStruct`, mapping `ErrUnauthorizedCustID` → 403. | @fixer | `cd sales && rtk go test ./controller -run TestSecondarySales` | 403 saat distributor request sibling, 200 saat principal child valid | ready | no
T06 | T02 | Update controller `SecondaryReportSalesTrendSales`: dual-bind query+body, panggil service signature baru, mapping 403. | @fixer | `cd sales && rtk go test ./controller -run TestTrendSales` | Body JSON `{"cust_id":"..."}` di GET ter-bind benar; 403 untuk dist sibling | ready | no
T07 | T03,T04,T05,T06 | Tambah/ubah test sesuai TDD/Test Plan. | @fixer | `cd sales && rtk go test ./...` | Semua tes hijau | ready | no
T08 | T07 | Jalankan validasi DB read-only terhadap dev DB untuk memverifikasi scope customer, trend data, dan owner `report.list`; credential harus via env var/secret out-of-band. | @fixer | `psql` read-only SELECT dengan `PGPASSWORD`/secret env; jangan tulis password ke artefak | Bukti query menunjukkan child scope valid dan `report.fact_orders.cust_id` filter bekerja | ready | no
T09 | T08 | Final review tenant isolation + signoff. | @quality-gate | `cd sales && rtk go test ./...` plus review evidence DB redacted | Tidak ada cross-tenant path; field `report.list.cust_id` tetap auth | ready | no
```

`start_with`: `T01`.

## Validation Commands

Dari `cd sales`:

```bash
rtk go mod download && rtk go mod tidy
rtk go test ./entity/... ./controller/... ./service/... ./repository/... -count=1
rtk go test ./service -run 'TestPublishSecondarySalesReport|TestSubscribeSecondarySalesReport|TestSecondarySalesReportTrendSales' -count=1
rtk go test ./controller -run 'TestSecondarySales|TestTrendSales' -count=1
```

Read-only DB validation (gunakan secret/env out-of-band; jangan commit password):

```bash
PGPASSWORD='<redacted>' psql -h 103.28.219.73 -p 25431 -U postgres -d scylla_citus_dev -At -c "select current_database(), current_user"

PGPASSWORD='<redacted>' psql -h 103.28.219.73 -p 25431 -U postgres -d scylla_citus_dev -At -c "select cust_id, parent_cust_id, is_del, is_active from smc.m_customer where parent_cust_id='C26002' order by cust_id limit 20"

PGPASSWORD='<redacted>' psql -h 103.28.219.73 -p 25431 -U postgres -d scylla_citus_dev -At -c "select m.month, coalesce(sum(fo.gross_sale),0)::bigint as total_gross_sale, coalesce(sum(fo.discount + fo.special_discount),0)::bigint as total_discount_promo, coalesce(sum(fo.net_sales_exclude_ppn),0)::bigint as net_sales from (select generate_series(1,12) as month) m left join report.dim_dates dt on dt.month=m.month and dt.year=2026 left join report.fact_orders fo on fo.date_id=dt.id and fo.cust_id='C260020001' group by m.month order by m.month"
```

Optional integration smoke (manual setelah deploy ke staging):

```bash
# Distributor user (auth dist), tidak kirim cust_id → fallback ke auth
curl -X POST 'https://staging.scyllax.online/sales/v1/reports/secondary-sales' \
  -H 'Authorization: Bearer <DIST_TOKEN>' \
  -H 'Content-Type: application/json' \
  --data '{"from":1777568400,"to":1779123599,"distributor_ids":[],"outlet_ids":[],"salesman_ids":[],"pro_ids":[]}'

# Distributor user kirim sibling cust_id → harus 403
curl -X POST 'https://staging.scyllax.online/sales/v1/reports/secondary-sales' \
  -H 'Authorization: Bearer <DIST_TOKEN>' \
  -H 'Content-Type: application/json' \
  --data '{"from":1777568400,"to":1779123599,"cust_id":"OTHER-DIST"}'

# Principal kirim child cust_id valid → 200, data sesuai child
curl -X POST 'https://staging.scyllax.online/sales/v1/reports/secondary-sales' \
  -H 'Authorization: Bearer <PRINCIPAL_TOKEN>' \
  -H 'Content-Type: application/json' \
  --data '{"from":1777568400,"to":1779123599,"cust_id":"C260020001"}'

# Trend Sales — distributor user dengan body sibling → 403
curl -X GET 'https://staging.scyllax.online/sales/v1/reports/secondary-sales/trend-sales?year=2026' \
  -H 'Authorization: Bearer <DIST_TOKEN>' \
  -H 'Content-Type: application/json' \
  --data '{"cust_id":"OTHER-DIST"}'
```

Catatan QA:

- Beberapa proxy/CDN bisa drop body GET; bila Trend Sales tampak abai pada `cust_id`, periksa apakah body sampai ke server (Fiber log atau Apache/Nginx access log).

## Evidence Requirements

- Hasil `rtk go test ./...` (sales).
- Snapshot output staging untuk skenario fallback, principal child valid, dan distributor sibling 403.
- Konfirmasi `report.list.cust_id` di DB = auth user untuk export yang dilakukan principal terhadap child distributor.

## Done Criteria

- Semua tes hijau di `sales`.
- Manual smoke staging memenuhi acceptance criteria 1-8.
- Tidak ada perubahan response shape atau route.
- `report.list` tetap muncul di `GET /v1/reports` untuk pemilik auth.
- `@quality-gate` signoff untuk perubahan tenant isolation.

## Final Planning Summary

- Artefak dibuat:
  - `.opencode/plans/20260520-1851-secondary-sales-cust-id-filters.md` (primary).
  - `.opencode/evidence/20260520-1851-secondary-sales-cust-id-filters/discovery.md`.
  - `.opencode/evidence/20260520-1851-secondary-sales-cust-id-filters/db-validation.md`.
- Question gate: 3 pertanyaan, semua dijawab user.
- Keputusan kunci:
  - Trend Sales `cust_id` di JSON body, `year` tetap required.
  - Export `report.list.cust_id` = auth user.
- Asumsi tersisa: validator `alphanum,max=20` cocok untuk `cust_id`. Jika project punya pola validasi cust berbeda, sesuaikan saat eksekusi.
- Open questions: tidak ada.
- Kesiapan implementasi: `@fixer` boleh mulai dari T01.
- Cleanup: draft tidak dibuat. Evidence `discovery.md` dan `db-validation.md` dipertahankan karena referensi line-number dan hasil DB read-only berguna saat implementasi dan QA (akan direvisit saat `@quality-gate` signoff). Hapus setelah Done.
