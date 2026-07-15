# Discovery Evidence — SX-2182 Secondary Sales Multiselect

Task id: `20260608-1534-sx-2182-secondary-sales-multiselect`
Tanggal: `2026-06-08`
Mode: Maintenance Stability Mode

## Files inspected

Repo docs:
- `AGENTS.md`
- `.opencode/docs/index.md`
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`

Sales module:
- `sales/entity/report.go`
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/model/report.go`
- `sales/service/report_service_test.go`

Master module:
- `master/controller/business_unit_controller.go`
- `master/entity/business_unit.go`
- `master/service/business_unit_service.go`
- `master/repository/business_unit_repository.go`
- `master/controller/business_unit_controller_test.go`
- `master/repository/business_unit_repository_test.go`

Prior planning artifacts consulted:
- `.opencode/plans/20260511-1530-sx-1944-secondary-sales-report.md`
- `.opencode/plans/20260519-1250-secondary-sales-dashboard-filters.md`

## Project patterns found

- Repo adalah multi-module Go monorepo. Validasi dijalankan per target service folder, terutama `master/` dan `sales/` untuk issue ini.
- Repo-local policy mewajibkan `rtk` prefix untuk shell workflow di repo ini.
- Arsitektur wajib `Controller -> Service -> Repository -> DB`.
- Tenant rules: `cust_id` untuk transaksi distributor, `parent_cust_id` untuk scope parent/principal.
- `sales` memakai Fiber controller, service resolver, GORM repository, dan async RMQ export.
- `master` business-unit memakai Fiber controller, service scope employee, repository `sqlx`, dan `sqlx.In` untuk list query.

## Existing implementation evidence

### `GET /master/v1/business-unit`

- Route ada di `master/controller/business_unit_controller.go:31-34`.
- Controller sudah override `RegionId` dan `AreaId` memakai `normalizeIntArrayQuery(c.Context().QueryArgs(), "region_id[]", "region_id")` dan equivalent untuk `area_id` di lines 55-56.
- Helper `normalizeIntArrayQuery` di lines 105-135 sudah support repeated args, comma-separated values, whitespace trim, dedupe, dan empty token skip.
- Catatan risiko: invalid numeric token sekarang di-skip diam-diam (`strconv.Atoi` error -> `continue`), bukan 400.
- Repository `master/repository/business_unit_repository.go:135-143` sudah memakai `md.region_id IN (?)` dan `md.area_id IN (?)` dengan `sqlx.In`, bukan string interpolation.
- Tests existing:
  - `master/controller/business_unit_controller_test.go` menguji comma-separated `region_id[]=1,2,3`, repeated `region_id[]=1&region_id[]=2`.
  - Belum terlihat test eksplisit untuk non-bracket query `region_id=80,90` dan whitespace `region_id=80, 90`, meski helper menerima key `region_id`.
  - `master/repository/business_unit_repository_test.go` menguji `IN` untuk multi region/area.

### `POST /sales/v1/reports/secondary-sales`

- Route ada di `sales/controller/report_controller.go:67-72`, `Post("/secondary-sales", controller.SecondarySales)`.
- Export request body DTO private `secondarySalesExportBody` masih `RequestedCustID string` di `sales/controller/report_controller.go:38-49`.
- Entity `SecondarySalesReportQueryFilter` masih punya `RequestedCustID string` di `sales/entity/report.go:33-51`.
- Service `PublishSecondarySalesReport` memakai `resolveSecondaryDashboardCustID` single-cust lalu overwrite `dataFilter.CustID = effectiveCustID` di `sales/service/report_service.go:345-355`.
- RMQ payload memakai serialized `SecondarySalesReportQueryFilter`; subscriber punya fallback single `RequestedCustID` ke `CustID` di `sales/service/report_service.go:431-436`.
- Repository builder `buildSecondarySalesUnionQuery` masih `od.cust_id = ?` dan `rd.cust_id = ?` di `sales/repository/report_repository.go:147-152`.
- Optional filters sudah memakai slice binding `IN ?` untuk distributor/salesman/outlet/product di `sales/repository/report_repository.go:166-193`.
- Service tests sudah menutup single child cust, unauthorized sibling, fallback auth cust, dan RMQ payload single cust di `sales/service/report_service_test.go`.

### Dashboard endpoints

- `sum-date`, `group`, dan `trend-sales` routes ada di `sales/controller/report_controller.go:70-72`.
- `sum-date` DTO sudah punya `Month`, optional `Year *int`, dan `CustID string` query di `sales/entity/report.go:214-218`.
- `group` DTO sudah punya `Month`, optional `Year *int`, `CustID string`, dan `GroupBy` di `sales/entity/report.go:225-230`.
- `trend-sales` DTO typo name `SecondarySalesReportTrensSalesSumPayload`, field `CustID string` bertag `json:"cust_id"`, bukan `query:"cust_id"`, tapi Fiber QueryParser dapat mengandalkan field name; perlu verifikasi/ubah agar eksplisit.
- Controller trend-sales sudah mencoba QueryParser dan optional BodyParser untuk GET body di `sales/controller/report_controller.go:434-468`.
- Service resolver `resolveSecondaryDashboardCustID` single-cust di `sales/service/report_service.go:1274-1292`:
  - empty atau same as auth -> auth cust
  - distributor login (`authCustID != parentCustID`) tidak boleh request sibling
  - principal perlu `ExistsCustomerInParentScope(requestedCustID, parentCustID)`
- `sum-date`, `group`, dan `trend-sales` service path semua masih menerima satu effective cust string.
- Repository dashboard queries masih filter single cust:
  - sum `fo.cust_id = ?`, `fr.cust_id = ?` di `sales/repository/report_repository.go:1115` dan `1127`.
  - group builder `fo.cust_id = ?`, `fr.cust_id = ?` di `sales/repository/report_repository.go:1194` dan `1202`.
  - trend `fo.cust_id = ?` di `sales/repository/report_repository.go:1264`.
- Group response entity `SecondarySalesReportGroupResp` belum punya `code` field di `sales/entity/report.go:247-251`.
- Model `SecondarySalesReportGroup` perlu dicek/diubah agar scan `code` bisa masuk.

## Reuse candidates

- Reuse `normalizeIntArrayQuery` concept di `master/controller/business_unit_controller.go` untuk `region_id`/`area_id`; tambahkan strict error path jika requirement 400 dipilih.
- Reuse `sqlx.In` pattern di `master/repository/business_unit_repository.go` untuk master query.
- Reuse `IN ?` slice binding GORM yang sudah dipakai di `buildSecondarySalesUnionQuery` untuk `cust_ids` di `sales/repository/report_repository.go`, bukan raw `ANY(:cust_ids)` string interpolation.
- Reuse existing `ErrUnauthorizedCustID` mapping ke HTTP 403 di controller sales.
- Extend `resolveSecondaryDashboardCustID` menjadi multi resolver, bukan buat auth check baru di controller.
- Extend `sales/service/report_service_test.go` mock hook pattern.

## Constraints

- Planner hanya menulis artifact di `.opencode/`; implementasi source harus dilakukan lane implementasi setelah plan.
- Scope source implementation lintas dua module: `master/` dan `sales/`.
- Jangan mengubah `pjp-sales/` kecuali user secara eksplisit meminta parity.
- Jaga async export RMQ payload backward compatibility.
- Jangan silent allow unauthorized `cust_id`; rekomendasi reject 403 untuk seluruh request jika ada requested cust yang tidak allowed.

## Risks

- Broken access control karena multi `cust_id` memperbesar blast radius data leak.
- `cust_id: []` atau empty query bisa berubah dari default auth cust menjadi no rows kalau langsung dipakai `IN ?` dengan empty slice.
- RMQ payload lama dan baru harus sama-sama bisa diproses oleh subscriber.
- `cust_id` body `string | []string` tidak bisa langsung di-bind ke `string` Go; perlu custom unmarshal/type atau raw body parsing.
- Group `id` mungkin tabrakan antar `cust_id` bila dim tables tidak global; implementer harus verify schema atau group by `(id, code, name)` dan catat evidence.
- Business-unit invalid numeric token saat ini di-skip; acceptance merekomendasikan 400, jadi implementer perlu memutuskan sesuai plan, bukan membiarkan silent skip.

## Source strategy

Digunakan:
- Repo-local docs dan source code.
- Jira/requirement detail dari prompt user sebagai product reference.

Diskip dengan alasan:
- Official docs/context7: tidak dibutuhkan karena behavior utama berasal dari requirement Jira dan pattern Go/GORM/sqlx lokal.
- GitHub/web search: tidak dibutuhkan karena tidak bergantung upstream behavior.
- Browser/runtime smoke: belum dilakukan karena planner tidak mengimplementasi dan tidak ada kredensial/API token runtime disediakan.

## Open questions

Tidak ada pertanyaan blocking untuk membuat plan. Policy authorization yang dipilih di plan: reject seluruh request dengan 403/400 bila ada requested `cust_id` unauthorized, karena lebih aman daripada silent ignore.
