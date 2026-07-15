# Plan — Task 118/119 Secondary Sales Dashboard Filters

Task id: `20260519-1250-secondary-sales-dashboard-filters`
Tanggal: `2026-05-19`
Target service: `sales`
Primary source of truth: `.opencode/plans/20260519-1250-secondary-sales-dashboard-filters.md`

## Goal

Perbaiki `GET /sales/v1/reports/secondary-sales/sum-date` dan `GET /sales/v1/reports/secondary-sales/group` agar mendukung filter `year` dan `cust_id` query secara aman, tanpa mengubah response shape existing.

## Non-goals

- Tidak mengubah route endpoint.
- Tidak mengubah response JSON contract.
- Tidak menambah request body untuk endpoint `GET`.
- Tidak mengubah extraction report atau trend sales di luar kebutuhan compile/test.
- Tidak mengimplementasikan filter `region`/`area` pada Task 118/119 karena bagian endpoint target di docs hanya meminta `year` dan `cust_id`.
- Tidak mengubah `Export Secondary Sales` pada dokumen kecuali user memperluas scope task.
- Tidak membuka akses lintas `cust_id` tanpa scope check.
- Tidak menambah dependency baru bila test bisa memakai pattern repo existing.

## Scope

Endpoint target:

- `GET /sales/v1/reports/secondary-sales/sum-date`
- `GET /sales/v1/reports/secondary-sales/group`

Docs reference target:

- `docs/Secondary Sales Report_BE.md`, bagian `SUM Date` lines 59-93.
- `docs/Secondary Sales Report_BE.md`, bagian `Secondary Sales Group` lines 96-123.

Layer target:

- `sales/entity/report.go`
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`

## Docs Reference Alignment

- Dokumen `SUM Date` meminta tambah `year` dari `report.dim_dates` dan `cust_id` untuk business unit ke `report.fact_orders.cust_id`.
- Dokumen `Secondary Sales Group` meminta tambah `year` dan `cust_id` untuk business unit ke `report.fact_orders.cust_id`.
- Dokumen menulis `Request Body`, tetapi endpoint adalah `GET` dan code existing memakai `QueryParser`; plan tetap memakai query param supaya selaras HTTP behavior existing.
- Response contract pada dokumen tetap `message`, `data`, `request_id`; plan menjaga response shape existing.
- Dokumen global menyebut filter wilayah `Region`, `Area`, `Distributor` untuk Homepage dan Modal Export. Untuk Task 118/119, bagian endpoint hanya merinci `year` dan `cust_id`; `region`/`area`/export endpoint tetap out-of-scope kecuali scope task diperluas.
- `cust_id` query pada plan dipakai sebagai selector business unit/distributor, tetapi wajib scope check terhadap child `cust_id` di bawah `parent_cust_id` user.

## Requirements

- Payload query support `month`, `year`, `cust_id` untuk sum-date.
- Payload query support `month`, `year`, `cust_id`, `group_by` untuk group.
- `month` wajib dan valid `1..12`.
- `year` opsional; bila kosong fallback ke `time.Now().Year()`.
- `year` bila dikirim harus valid, minimal `>= 2000` dan realistis `<= 9999`.
- `cust_id` opsional; bila kosong fallback ke auth `cust_id`.
- `cust_id` query hanya boleh:
  - sama dengan auth `cust_id`; atau
  - child `cust_id` di bawah `parent_cust_id` user login, hanya bila user login principal (`authCustID == parentCustID`).
- Distributor login (`authCustID != parentCustID`) tidak boleh query sibling `cust_id`.
- Repository selalu filter `dt."year" = ?` memakai effective year.
- Semua branch `group_by` memakai filter `cust_id + month + year` konsisten.
- `return_rate` tidak boleh divide-by-zero.

## Acceptance Criteria

1. `sum-date` menerima `month`, `year`, dan `cust_id` query.
2. `group` menerima `month`, `year`, `cust_id`, dan `group_by` query.
3. Jika `year` dikirim, query filter `report.dim_dates` dengan `dt."year" = ?`.
4. Jika `year` kosong, query memakai current year sebagai fallback.
5. Jika `cust_id` kosong, query memakai auth `cust_id`.
6. Jika `cust_id` query sama dengan auth `cust_id`, query memakai auth `cust_id`.
7. Jika principal user mengirim child `cust_id` valid di bawah `parent_cust_id`, query memakai child `cust_id` tersebut.
8. Jika distributor user mengirim sibling/foreign `cust_id`, endpoint mengembalikan 403 dan tidak query report data.
9. Semua branch `group_by` (`outlet`, `salesman`, `product_category`, default `product`) punya filter sama.
10. Response shape existing tetap sama.
11. Unit/repository tests menutup year, cust scope, semua branch group, fallback year, fallback cust, dan divide-by-zero.

## Existing Patterns/Reuse

- Controller report existing memakai `c.QueryParser(&request)` dan `responsebuild.BuildResponse`.
- Tenant locals tersedia dari middleware:
  - `c.Locals("cust_id")`
  - `c.Locals("parent_cust_id")`
- Repo docs menyatakan transaksi distributor pakai `cust_id`; parent master pakai `parent_cust_id`.
- `docs/Secondary Sales Report_BE.md` menjadi referensi requirement lokal untuk Task 118/119.
- `sales/model/companies.go` sudah punya `model.SmcMCustomer` untuk `smc.m_customer`.
- `sales/repository/hierarchy_approval_repository.go` punya pattern scope parent via `smc.m_customer.parent_cust_id`.
- `SecondarySalesReportTrendSales` sudah memakai `dt."year"`, jadi gunakan quoted column untuk konsistensi.
- Test repository bisa reuse `gorm.Config{DryRun: true}` pattern dari `stock_repository_cancel_test.go`.
- Test service bisa reuse mock repository hook pattern di `report_service_test.go`.

## Constraints

- Project local command harus `rtk`-prefixed sesuai `AGENTS.md` repo.
- Validasi harus dijalankan dari folder `sales`.
- Jangan commit secrets, `.env`, atau backup dumps.
- Jangan melewati layer: controller tidak boleh langsung query DB.
- Perubahan auth/tenant wajib final signoff `@quality-gate`.
- Dokumen lokal `docs/Secondary Sales Report_BE.md` sudah dicek untuk alignment; jangan edit dokumen ini kecuali user meminta update docs.
- Tidak ditemukan Swagger/Postman target untuk endpoint ini di `sales`; update docs hanya bila implementer menemukan file API docs relevan saat eksekusi.

## Risks

- Query `cust_id` tanpa scope check menyebabkan cross-tenant leak.
- `year` fallback current year mengubah behavior request lama: dulu multi-year, setelah fix current-year only.
- Pointer `Year *int` perlu validasi benar supaya `year=0` tidak diam-diam fallback.
- Jika service signature berubah, semua compile callers dan tests harus ikut update.
- `return_rate` existing rawan `Inf` bila order qty `0` dan return qty `>0`; harus diperbaiki.

## Decisions/Assumptions

- Pertanyaan diajukan dan dijawab.
- Keputusan user:
  - `year` opsional dengan fallback current year.
  - `cust_id` query override boleh, tetapi wajib scope check.
  - Business unit = child `cust_id` di bawah `parent_cust_id` user.
  - Tambah validasi range `month` dan `year`.
- Asumsi implementasi:
  - Principal user ditandai oleh `authCustID == parentCustID`.
  - Distributor user ditandai oleh `authCustID != parentCustID`.
  - Scope child dicek via `smc.m_customer` memakai `cust_id = requestedCustID AND parent_cust_id = parentCustID`.
  - Jika `smc.m_customer.is_del` tersedia, tambahkan `AND is_del = false`; `is_active = true` hanya jika product owner mengonfirmasi inactive harus ditolak.
  - Error unauthorized cust map ke HTTP `403 Forbidden`.

## TDD/Test Plan

TDD required: ya.

Reason:

- Ini perubahan production logic, API behavior, DB query, dan tenant isolation.

Existing test patterns:

- `sales/service/report_service_test.go` memakai mock repository dengan function hooks.
- `sales/repository/report_repository_test.go` memakai SQL string/assert param helper.
- `sales/repository/stock_repository_cancel_test.go` memakai GORM DryRun untuk SQL inspection.

First failing/regression tests:

1. Service test: `cust_id` query child valid untuk principal harus meneruskan child `cust_id`, explicit `year`, dan `month` ke repository.
2. Service test: distributor login dengan requested sibling `cust_id` harus return `ErrUnauthorizedCustID` dan tidak memanggil report repository.
3. Repository test: sum orders/returns SQL harus mengandung `dt."year" = ?` dan vars berisi `custID, month, year`.
4. Repository test: semua group branch SQL harus mengandung `dt."year" = ?`.
5. Service test: `sum-date` dengan order qty `0` dan return qty `>0` menghasilkan `return_rate = 0`, bukan `Inf`.

Green step:

- Tambahkan payload fields, service resolver, repository scope method, repository year params, dan controller validation/error mapping sampai tests lulus.

Refactor step:

- Extract small helpers bila duplikasi tinggi:
  - `resolveSecondaryDashboardYear(year *int) int`
  - `resolveSecondaryDashboardCustID(authCustID, parentCustID, requestedCustID string) (string, error)` sebagai method service karena butuh repository scope.
  - query builder helpers hanya bila dibutuhkan agar repository DryRun tests stabil.

Edge cases:

- Missing `year` → current year.
- `year=1999` → 400.
- `month=0` atau `month=13` → 400.
- Empty `cust_id` → auth `cust_id`.
- Requested `cust_id` equals auth → allow tanpa scope DB call.
- Principal requested child under same parent → allow.
- Distributor requested sibling under same parent → deny 403.
- Unknown requested `cust_id` → deny 403.
- `group_by` unknown/empty → default product tetap.

Commands:

```bash
rtk go test ./service -run 'TestSecondarySalesReport(SumReportByMonth|GroupSales|Resolve|Dashboard)'
rtk go test ./repository -run 'TestSecondarySalesReport(SumReportByMonth|ReturnSumReportByMonth|Group)'
rtk go test ./...
```

## Implementation Steps

1. Update payload structs di `sales/entity/report.go`:

```go
type SecondarySalesReportDashboardSumPayload struct {
	Month  int  `query:"month" validate:"required,gte=1,lte=12"`
	Year   *int `query:"year" validate:"omitempty,gte=2000,lte=9999"`
	CustID string `query:"cust_id" validate:"omitempty,alphanum"`
}

type SecondarySalesReportDashboardGroupPayload struct {
	Month   int    `query:"month" validate:"required,gte=1,lte=12"`
	Year    *int   `query:"year" validate:"omitempty,gte=2000,lte=9999"`
	CustID  string `query:"cust_id" validate:"omitempty,alphanum"`
	GroupBy string `query:"group_by" validate:"omitempty,oneof=outlet salesman product_category product"`
}
```

Catatan: jika validator `oneof` dengan empty string bermasalah, pakai `omitempty,oneof=...` atau validasi manual. Jangan ubah default product behavior.

2. Update controller `SecondaryReportSalesSumMonth` dan `SecondaryReportSalesGroup`:

- Setelah `QueryParser`, panggil `controller.validator.ValidateStruct(request, headerAcceptLang)`.
- Jika invalid, return `400` dengan pattern response existing.
- Ambil `authCustID := c.Locals("cust_id").(string)` dan `parentCustID := c.Locals("parent_cust_id").(string)`.
- Panggil service dengan signature baru: `SecondarySalesReportSumReportByMonth(authCustID, parentCustID, request)` dan `SecondarySalesReportGroupSales(authCustID, parentCustID, request)`.
- Jika `errors.Is(err, service.ErrUnauthorizedCustID)`, return `403`.
- Error lain tetap `400` seperti existing.

3. Update service interface dan implementation:

- Signature:

```go
SecondarySalesReportSumReportByMonth(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardSumPayload) (...)
SecondarySalesReportGroupSales(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardGroupPayload) (...)
```

- Tambah exported sentinel error di `service/report_service.go`:

```go
var ErrUnauthorizedCustID = errors.New("cust_id is outside authorized scope")
```

- Tambah helper effective year:

```go
func resolveSecondaryDashboardYear(year *int) int {
	if year == nil {
		return time.Now().Year()
	}
	return *year
}
```

- Tambah resolver cust:

```go
func (service *reportServiceImpl) resolveSecondaryDashboardCustID(authCustID, parentCustID, requestedCustID string) (string, error) {
	if requestedCustID == "" || requestedCustID == authCustID {
		return authCustID, nil
	}
	if authCustID != parentCustID {
		return "", ErrUnauthorizedCustID
	}
	allowed, err := service.ReportRepository.ExistsCustomerInParentScope(requestedCustID, parentCustID)
	if err != nil {
		return "", err
	}
	if !allowed {
		return "", ErrUnauthorizedCustID
	}
	return requestedCustID, nil
}
```

4. Update repository interface:

```go
ExistsCustomerInParentScope(custID string, parentCustID string) (bool, error)
SecondarySalesReportSumReportByMonth(custID string, month int, year int) (...)
SecondarySalesReportReturnSumReportByMonth(custID string, month int, year int) (...)
SecondarySalesReportGroupOutlet(custID string, month int, year int) (...)
SecondarySalesReportGroupSalesman(custID string, month int, year int) (...)
SecondarySalesReportProductCategory(custID string, month int, year int) (...)
SecondarySalesReportProduct(custID string, month int, year int) (...)
```

5. Implement `ExistsCustomerInParentScope` di `report_repository.go` memakai `model.SmcMCustomer`:

```go
func (repository *RepositoryReportImpl) ExistsCustomerInParentScope(custID string, parentCustID string) (bool, error) {
	var count int64
	err := repository.Model(&model.SmcMCustomer{}).
		Where("cust_id = ? AND parent_cust_id = ? AND is_del = false", custID, parentCustID).
		Count(&count).Error
	return count > 0, err
}
```

Jika DB lokal membuktikan `is_del` nullable/bermasalah, fallback aman: hilangkan `is_del = false` sesuai existing hierarchy query, tetapi catat di final implementation summary.

6. Update repository report queries:

- Orders sum:

```go
Where("report.fact_orders.cust_id = ? AND dt.month = ? AND dt.\"year\" = ?", custID, month, year)
```

- Returns sum:

```go
Where("report.fact_returns.cust_id = ? AND dt.month = ? AND dt.\"year\" = ?", custID, month, year)
```

- Semua group branch:

```go
Where("report.fact_orders.cust_id = ? AND dt.month = ? AND dt.\"year\" = ?", custID, month, year)
```

- Pertimbangkan `COALESCE(SUM(...), 0)` untuk numeric sums agar response tetap `0`, bukan null scan issue.

7. Update service calls:

- `sum-date` harus memakai effective cust dan effective year untuk orders + returns.
- `group` switch harus memakai effective cust dan effective year untuk semua branch.
- Fix `return_rate`:

```go
if sumReportModel.Qty > 0 {
	data.ReturnRate = (float64(sumReportReturnModel.Qty) / float64(sumReportModel.Qty)) * 100
}
```

8. Tambah/update tests:

- Mock repository di `report_service_test.go` perlu hook untuk:
  - `ExistsCustomerInParentScope`
  - `SecondarySalesReportSumReportByMonth`
  - `SecondarySalesReportReturnSumReportByMonth`
  - group branch methods.
- Repository SQL tests boleh pakai query builder helpers atau GORM DryRun. Jika existing direct methods sulit diinspect, extract private query builder methods dan test builder.

## Expected Files to Change

- `sales/entity/report.go`
- `sales/controller/report_controller.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/service/report_service_test.go`
- `sales/repository/report_repository_test.go`

Possible only if discovered during implementation:

- API docs/Postman/Swagger files, jika endpoint ini terdokumentasi di repo.

## Agent/Tool Routing

- `@orchestrator`: route execution and integration.
- `@fixer`: implement bounded code + tests in `sales` module.
- `@explorer`: optional follow-up only if compile reveals hidden callers or auth patterns.
- `@quality-gate`: mandatory final signoff karena tenant isolation/security-sensitive.
- `@architect`: not needed unless business asks BU model beyond child `cust_id` under `parent_cust_id`.

## Execution-ready Worklist / Handoff Contract

`start_with`: `T01`

| Task | Action | depends_on | owner/lane | validation/check | exit criteria | status | blocker | requires_user_decision |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T01 | Add failing service tests for effective `cust_id`, year fallback, unauthorized scope, and return_rate zero guard. | none | `@fixer` | `rtk go test ./service -run 'TestSecondarySalesReport'` should fail before code | Tests fail for missing implementation, not compile syntax mistakes. | ready | none | no |
| T02 | Add failing repository SQL tests for sum orders, sum returns, and all group branches requiring `dt."year" = ?`. | T01 | `@fixer` | `rtk go test ./repository -run 'TestSecondarySalesReport'` should fail before query update | Tests verify SQL fragments and vars. | ready | none | no |
| T03 | Update payload structs with `year`, `cust_id`, and validation tags. | T02 | `@fixer` | `rtk go test ./service ./repository` compile | Structs parse query fields and preserve response structs. | ready | none | no |
| T04 | Update controller validation, parent/auth locals, service calls, and 403 mapping for unauthorized cust. | T03 | `@fixer` | `rtk go test ./...` compile target | Both handlers parse/validate query and pass `authCustID`, `parentCustID`, request. | ready | none | no |
| T05 | Add service resolver for effective year and effective scoped cust. | T04 | `@fixer` | `rtk go test ./service -run 'TestSecondarySalesReport'` | Empty cust fallback, child cust allow, sibling deny, missing year fallback work. | ready | none | no |
| T06 | Add repository `ExistsCustomerInParentScope` and update interface/mocks. | T05 | `@fixer` | `rtk go test ./service` | Scope DB method compiles and service mock covers it. | ready | none | no |
| T07 | Update repository signatures and queries to require `year` for orders, returns, and group branches. | T06 | `@fixer` | `rtk go test ./repository -run 'TestSecondarySalesReport'` | SQL includes `dt."year" = ?` in every target query. | ready | none | no |
| T08 | Fix return_rate divide-by-zero and preserve summary merge/last_update logic. | T07 | `@fixer` | `rtk go test ./service -run 'TestSecondarySalesReportSumReportByMonth'` | `Qty=0` never produces `Inf`; last_update logic unchanged. | ready | none | no |
| T09 | Run full sales validation. | T08 | `@fixer` | `rtk go mod download && rtk go test ./...` from `sales` | All relevant tests pass. | ready | none | no |
| T10 | Final security/quality review. | T09 | `@quality-gate` | Review diff + test evidence | No tenant leak; scope and year behavior documented in final summary. | ready | none | no |

## Validation Commands

Run from `/Users/ujang/Projects/Geekgarden/scylla-be/sales`:

```bash
rtk go test ./service -run 'TestSecondarySalesReport(SumReportByMonth|GroupSales|Resolve|Dashboard)'
rtk go test ./repository -run 'TestSecondarySalesReport(SumReportByMonth|ReturnSumReportByMonth|Group)'
rtk go mod download && rtk go test ./...
```

Optional runtime smoke after deploy/local service is up:

```bash
curl 'https://<host>/sales/v1/reports/secondary-sales/sum-date?month=5&year=2026&cust_id=C2600200001' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer <VALID_TOKEN>'

curl 'https://<host>/sales/v1/reports/secondary-sales/group?month=5&year=2026&cust_id=C2600200001&group_by=outlet' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer <VALID_TOKEN>'
```

Manual DB check, if DB access allowed:

```sql
SELECT SUM(fo.net_sales_exclude_ppn) AS net_sales
FROM report.fact_orders fo
JOIN report.dim_dates dt ON fo.date_id = dt.id
WHERE fo.cust_id = 'C2600200001'
  AND dt.month = 5
  AND dt."year" = 2026;
```

## Evidence Requirements

Implementation final summary must include:

- Files changed.
- Root cause final.
- Before vs after query fragment for orders, returns, and group.
- Service scope decision: requested `cust_id` allowed only for auth cust or principal child scope.
- Validation commands and outputs.
- Any skipped test and reason.
- Runtime curl evidence if environment/token available; otherwise state blocker.
- Quality-gate result because auth/tenant isolation changed.

## Done Criteria

- Tests pass in `sales` module.
- No arbitrary `cust_id` query access.
- `dt."year" = ?` present in orders summary, returns summary, and every group branch.
- `year` fallback current year documented.
- `month`, `year`, `cust_id`, `group_by` validation documented.
- Response shape unchanged.
- `return_rate` safe when order qty zero.
- `@quality-gate` signoff completed.

## Final Planning Summary

Artifacts created:

- `.opencode/evidence/20260519-1250-secondary-sales-dashboard-filters/discovery.md`
- `.opencode/evidence/20260519-1250-secondary-sales-dashboard-filters/docs-reference-check.md`
- `.opencode/draft/20260519-1250-secondary-sales-dashboard-filters/open-questions.md`
- `.opencode/plans/20260519-1250-secondary-sales-dashboard-filters.md`

Artifacts consulted:

- `docs/Secondary Sales Report_BE.md`.
- Repo docs: `.opencode/docs/index.md`, `ARCHITECTURE.md`, `QUALITY.md`, `AGENT_ROUTING.md`, `SECURITY.md`, `PROMPT_GATES.md`.
- Source files listed in discovery evidence.
- `@explorer` and `@architect` read-only findings.

Key decisions:

- `year` optional, fallback current year.
- `cust_id` query override allowed only with scope check.
- BU = child `cust_id` under `parent_cust_id` user.
- Validate `month`/`year` range.

Open questions:

- None blocking after user decisions.
- Minor implementation choice remains: include `is_active = true` in scope check only if business confirms inactive child must be denied. Default plan only uses `is_del = false`.

Readiness:

- Ready for implementation by `@orchestrator` → `@fixer` with `@quality-gate` final signoff.

Cleanup performed:

- Draft/evidence still kept because implementation handoff may need raw discovery and resolved question trail for tenant-scope audit.
