# Discovery — Task 120/121 Secondary Sales `cust_id` Filters

Task id: `20260520-1851-secondary-sales-cust-id-filters`
Tanggal: `2026-05-20`
Target service: `sales`

## Files inspected

- `sales/controller/report_controller.go`
  - Route export POST: `/v1/reports/secondary-sales` → `controller.SecondarySales` (line 41).
  - Route trend sales GET: `/v1/reports/secondary-sales/trend-sales` → `controller.SecondaryReportSalesTrendSales` (line 45).
  - Handler `SecondarySales` (line 136-172) memakai `c.BodyParser(&request)` ke `entity.SecondarySalesReportQueryFilter`, lalu meng-overwrite `request.CustID = c.Locals("cust_id")` dan `request.ParentCustID = c.Locals("parent_cust_id")`. Implikasi: bila ditambahkan field `cust_id` ke body, harus dipindah ke field lain dulu sebelum overwrite, atau hentikan overwrite tanpa scope check.
  - Handler `SecondaryReportSalesTrendSales` (line 367-390) memakai `c.QueryParser(&request)` ke `entity.SecondarySalesReportTrensSalesSumPayload`, lalu memanggil service dengan `c.Locals("cust_id")` langsung sebagai `custID`. Belum ada field `CustID` di payload.
- `sales/entity/report.go`
  - `SecondarySalesReportQueryFilter` (line 33-48) sudah punya field `CustID string \`json:"cust_id"\`` dan `ParentCustID string \`json:"parent_cust_id"\``. Kalau body `cust_id` user mau dipakai sebagai filter override, harus pakai field baru terpisah supaya tidak bertabrakan dengan auth `cust_id` yang ditulis di handler.
  - `SecondarySalesReportTrensSalesSumPayload` (line 209-211) hanya berisi `Year int \`query:"year" validate:"required"\``. Belum ada `CustID`.
  - `SecondarySalesReportDashboardSumPayload` (line 203-207) dan `SecondarySalesReportDashboardGroupPayload` (line 213-218) sudah punya `CustID string \`query:"cust_id"\``, hasil Task 118/119 (lihat plan `.opencode/plans/20260519-1250-secondary-sales-dashboard-filters.md`).
- `sales/service/report_service.go`
  - Helper `resolveSecondaryDashboardCustID(authCustID, parentCustID, requestedCustID string) (string, error)` (line 1163-1181) sudah ada dan reusable. Logikanya: empty atau equal auth → auth; principal (`auth == parent`) cek scope `ExistsCustomerInParentScope`; selain itu `ErrUnauthorizedCustID`.
  - Helper `resolveSecondaryDashboardYear(year *int) int` (line 187) sudah ada untuk fallback current year.
  - `PublishSecondarySalesReport(dataFilter entity.SecondarySalesReportQueryFilter)` (line 305-356) menulis `reportList.CustID = dataFilter.CustID` (line 319). Karena `dataFilter.CustID` di-overwrite dari `c.Locals("cust_id")` di controller, ini selalu auth `cust_id` — sesuai keputusan user (owner report.list = auth user).
  - `SubscribeSecondarySalesReport(dataFilter)` (line 374-466) memanggil `service.ReportRepository.SecondarySalesUnion(dataFilter)` lewat dataFilter di-publish lewat RabbitMQ (line 348 `structs.StructToJson(dataFilter)`), jadi field baru harus ada di struct yang sama supaya ikut serialized.
  - `SecondarySalesReportTrendSales(custID string, year int)` (line 1262-1278) menerima `custID` raw — belum ada scope check terhadap principal/distributor.
- `sales/repository/report_repository.go`
  - `buildSecondarySalesUnionQuery(dataFilter entity.SecondarySalesReportQueryFilter, withPagination bool)` (line 136) memakai `dataFilter.CustID` di:
    - `whereOrder := "od.cust_id = ? AND o.data_status IN (6,7)"` (line 147)
    - `whereReturn := "rd.cust_id = ? AND o.data_status IN (6,7)"` (line 148)
    - `paramsOrder := []interface{}{dataFilter.CustID}` (line 151)
    - `paramsReturn := []interface{}{dataFilter.CustID}` (line 152)
    - `allParams = append(allParams, dataFilter.CustID, dataFilter.ParentCustID)` (line 317) untuk LATERAL join `mst.m_product` parent.
  - `ExistsCustomerInParentScope(custID, parentCustID)` (line 68-75) sudah dipakai untuk Task 118/119 scope check.
  - `SecondarySalesReportTrendSales(custID string, year int)` (line 1148-1183) memakai `custID` langsung di SQL `LEFT JOIN report.fact_orders fo ON fo.date_id = dt.id AND fo.cust_id = ?`.

## Project patterns found

- Tenant: `cust_id` = data per-distributor; `parent_cust_id` = parent principal. Principal user terdeteksi `authCustID == parentCustID` (lihat `resolveSecondaryDashboardCustID`).
- Validasi memakai `validate:"..."` tag dan `c.Locals("cust_id")`/`c.Locals("parent_cust_id")` dari middleware JWT.
- Error scope unauthorized memakai `ErrUnauthorizedCustID` lalu di-mapping ke `fiber.StatusForbidden` di controller (lihat `SecondaryReportSalesSumMonth` line 357-360).
- Handler test pattern memakai `httptest.NewRequest` + `app.Test(req)` (`sales/controller/so_controller_test.go`).
- Service test pattern: mock repository hooks (`mockReportRepositoryForService` di `sales/service/report_service_test.go`).
- Repository test pattern: SQL string assertion via `buildSecondarySalesUnionQuery` (`sales/repository/report_repository_test.go`).

## Reuse candidates

- `resolveSecondaryDashboardCustID` di service untuk Task 120 dan 121.
- `ExistsCustomerInParentScope` repository untuk scope check.
- `ErrUnauthorizedCustID` error mapping → 403 di controller.
- Pola binding query+body: bisa dual-bind (`c.QueryParser` + `c.BodyParser`) di Task 121 tanpa duplikasi struct, atau dua struct dipisah. Plan memilih dual-bind dengan struct tunggal yang punya tag `query:"year"` dan `json:"cust_id"`.

## Commands / docs checked

- `docs/Secondary Sales Report_BE.md` lines 22-50 (Trend Sales), lines 119-141 (Export Secondary Sales).
- `.opencode/plans/20260519-1250-secondary-sales-dashboard-filters.md` (Task 118/119) untuk reference pattern.
- `rtk docker compose -f docker-compose.yml ps` di repo root → compose definition valid (warning `version` obsolete diabaikan, sesuai catatan `.opencode/docs/QUALITY.md`).

## Constraints

- `sales` adalah service Fiber dengan layering Controller→Service→Repository→DB ketat.
- Validasi `rtk go test ./...` dijalankan dari `cd sales`.
- Subscribe export pakai RabbitMQ JSON serialize; field baru WAJIB ada di struct yang dikirim ke RMQ.
- Tidak ada Swagger generator aktif untuk service `sales`; doc API hanya di `docs/Secondary Sales Report_BE.md`.

## Risks

- Cross-tenant leak bila `cust_id` body/query dipakai tanpa scope check.
- `report.list.cust_id` harus tetap auth user (keputusan user) supaya endpoint `GET /v1/reports` tetap menampilkan baris hasil export ke principal yang membuat.
- Field `CustID` di `SecondarySalesReportQueryFilter` saat ini meng-encode "auth cust_id" karena di-overwrite di handler. Menambahkan `cust_id` body baru harus di field terpisah (mis. `RequestedCustID`) supaya tidak bertabrakan dengan auth.
- Trend Sales endpoint method GET tapi mengikut docs harus pakai JSON body. Beberapa client/proxy/CDN bisa drop body di GET. Plan tetap meng-bind body sesuai keputusan user; berikan note untuk QA.
