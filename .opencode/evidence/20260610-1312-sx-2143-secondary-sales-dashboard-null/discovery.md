# Discovery Evidence — SX-2143 Secondary Sales Dashboard Data Null

Task ID: `20260610-1312-sx-2143-secondary-sales-dashboard-null`
Tanggal: 2026-06-10 Asia/Jakarta
Mode: Maintenance Stability Mode, artifact-only planning.

## Sumber yang digunakan

- Prompt user untuk Jira SX-2143, SX-2201, dan ringkasan BE docs.
- Repo-local harness:
  - `AGENTS.md`
  - `.opencode/docs/index.md`
  - `.opencode/docs/ARCHITECTURE.md`
  - `.opencode/docs/QUALITY.md`
  - `.opencode/docs/SECURITY.md`
  - `.opencode/docs/AGENT_ROUTING.md`
- Source lokal modul `sales`:
  - `sales/controller/report_controller.go`
  - `sales/service/report_service.go`
  - `sales/entity/report.go`
  - `sales/model/report.go`
  - `sales/repository/report_repository.go`
  - `sales/service/report_service_test.go`
  - `sales/repository/report_repository_test.go`
  - `sales/controller/so_controller_test.go`
  - `sales/go.mod`
- Tidak mengambil Jira/Google Docs langsung karena kemungkinan butuh kredensial dan user sudah memberi potongan relevan. Klaim dokumen eksternal dibatasi ke isi prompt.

## Runtime dan status lokal

- `rtk docker compose -f docker-compose.yml ps` dari repo root berhasil; service `scylla-sales` status `Up`, port `9004` terbuka.
- Peringatan command: `rtk` menolak filter tidak dipercaya di `.rtk/filters.toml`; ini tidak memblokir command.
- `git status --short` di `sales` kosong saat discovery.
- Targeted validation baseline: `rtk go test ./service ./repository ./controller` dari `sales` lulus: `Go test: 205 passed in 3 packages`.

## Pola arsitektur yang relevan

- Repo rule: perubahan harus menjaga Controller → Service → Repository → DB.
- Multi-tenant/cust rule: jaga filter `cust_id`; principal hanya boleh akses child via scope check; distributor default ke `cust_id` login.
- Shell command repo ini memakai prefix `rtk`.
- Service `sales` adalah Go module mandiri; validasi dari direktori `sales`.

## Temuan kode endpoint

### Route dan controller

- `sales/controller/report_controller.go:67-72` mendaftarkan route:
  - `GET /v1/reports/secondary-sales/sum-date` → `SecondaryReportSalesSumMonth`
  - `GET /v1/reports/secondary-sales/group` → `SecondaryReportSalesGroup`
- `SecondaryReportSalesSumMonth` di `report_controller.go:397-431` memakai `QueryParser` ke `entity.SecondarySalesReportDashboardSumPayload`, validasi struct, lalu meneruskan `authCustID` dan `parentCustID` dari JWT locals ke service.
- `SecondaryReportSalesGroup` di `report_controller.go:481-515` memakai pola serupa untuk `entity.SecondarySalesReportDashboardGroupPayload`.

### Payload dan response entity

- `sales/entity/report.go:214-218`:
  - `Month int query:"month" validate:"required,gte=1,lte=12"`
  - `Year *int query:"year" validate:"omitempty,gte=2000,lte=9999"`
  - `CustID string query:"cust_id" validate:"omitempty"`
- `sales/entity/report.go:225-230` group payload sama, plus `GroupBy string query:"group_by"`.
- `sales/entity/report.go:231-245` response `SumReportByMonthModelResp` sudah memiliki numeric non-pointer untuk `total_gross_sale`, `total_discount_promo`, `total_ppn`, `net_sales_exc_ppn`, `net_sales`, counts, `qty`, `qty_return`, `return_rate`, `net_sales_return`; `last_update` pointer.
- `sales/entity/report.go:247-251` `SecondarySalesReportGroupResp` saat discovery hanya punya `id`, `name`, `net_sales`; tidak ada `code` meski plan SX-2172 lama menyebut code.

### Service

- `resolveSecondaryDashboardYear` di `sales/service/report_service.go:227-233` memakai `time.Now().Year()` ketika `year == nil`.
- `resolveSecondaryDashboardCustID` di `report_service.go:1274-1292`:
  - default ke `authCustID` bila `requestedCustID` kosong atau sama auth.
  - jika user bukan principal (`authCustID != parentCustID`), requested sibling/child ditolak dengan `ErrUnauthorizedCustID`.
  - principal mengecek `ExistsCustomerInParentScope(requestedCustID, parentCustID)` sebelum boleh memakai child.
- `SecondarySalesReportSumReportByMonth` di `report_service.go:1294-1332`:
  - resolve effective cust id.
  - resolve effective year.
  - panggil repository order summary dan return summary dengan `(effectiveCustID, month, effectiveYear)`.
  - mapping response numeric dari `sumReportModel`.
  - return rate aman saat qty order `0`.
  - catatan: `sumReportReturnModel` hanya dipakai untuk `LastUpdate`; qty/net return sudah dihitung di query summary utama.
- `SecondarySalesReportGroupSales` di `report_service.go:1335-1372`:
  - resolve effective cust id + year.
  - branch `outlet`, `salesman`, `product_category`, default product.
  - mapping response saat discovery tidak mengisi `code` karena response entity tidak punya field code.

### Repository

- `sales/repository/report_repository.go:47-52` interface dashboard secondary sales sudah menerima `year int` untuk summary dan semua group branch.
- `SecondarySalesReportSumReportByMonth` di `report_repository.go:1099-1145`:
  - memakai raw SQL CTE `order_summary` dan `return_summary` dari `report.fact_orders` + `report.fact_returns`.
  - filter order: `fo.cust_id = ? AND dt.month = ? AND dt."year" = ?`.
  - filter return: `fr.cust_id = ? AND dt.month = ? AND dt."year" = ?`.
  - aggregate memakai `COALESCE(SUM(...), 0)`.
  - final menghitung order minus return untuk gross, ppn, net sales exc/inc; discount_promo order + return.
  - `GREATEST(os.last_update, rs.last_update)` dapat menjadi `NULL` bila salah satu sisi null di PostgreSQL; perlu verifikasi/normalisasi jika ingin last_update tetap maksimum non-null.
- `SecondarySalesReportReturnSumReportByMonth` di `report_repository.go:1150-1158` sudah filter `dt."year"` dan COALESCE numeric.
- `buildSecondarySalesReportGroupQuery` di `report_repository.go:1162-1208`:
  - order branch dan return branch dipadukan via `UNION ALL`.
  - setiap branch filter `cust_id`, `dt.month`, `dt."year"`.
  - return branch memakai `fr.net_sales_exclude_ppn * -1`.
  - final `COALESCE(SUM(net_sales), 0)`.
  - saat discovery select tidak menyertakan `code`.

## Test yang sudah ada

- `sales/service/report_service_test.go:413-453`: explicit year + child cust principal diteruskan ke repository.
- `sales/service/report_service_test.go:455-479`: unauthorized distributor sibling ditolak.
- `sales/service/report_service_test.go:481-513`: missing year fallback ke `time.Now().Year()` dan cust fallback ke auth cust.
- `sales/service/report_service_test.go:515-536`: return rate aman ketika qty order `0`.
- `sales/service/report_service_test.go:538-577`: PPN dan net sales exc/inc mapping.
- `sales/service/report_service_test.go:579-665`: group semua branch memakai fallback year.
- `sales/repository/report_repository_test.go:243-260`: summary SQL memakai quoted year filter dan expected vars order.
- `sales/repository/report_repository_test.go:262-279`: return summary SQL memakai quoted year filter.
- `sales/repository/report_repository_test.go:457-518`: group queries memakai quoted year filter untuk order/return dan include return subtraction.
- `sales/repository/report_repository_test.go:520-547`: summary SQL combine order and return fact.
- `sales/controller/so_controller_test.go:138-210`: controller returns 403 untuk unauthorized cust di sum/group.

## Gap dan risiko yang masih perlu divalidasi implementor

1. Banyak perubahan SX-2143 tampak sudah sebagian/seluruhnya ada di working tree saat discovery: year handling, fallback year, COALESCE, CTE order+return, dan tests baseline sudah ada. Implementor perlu membandingkan dengan branch target/Jira branch aktual sebelum mengklaim fixed.
2. Tidak ada controller tests eksplisit untuk parsing `month=6&year=2026&cust_id=C260020001` sukses, missing `year` sukses, invalid `month/year` 400. Test forbidden ada, tetapi parsing sukses perlu ditambah agar acceptance lebih kuat.
3. `group` response saat discovery tidak punya `code`, bertentangan dengan acceptance user yang menyebut group endpoints tetap return `code`, `name`, `net_sales` sesuai enhancement sebelumnya. Ini perlu keputusan implementasi: tambahkan `Code` ke `model.SecondarySalesReportGroup` dan `entity.SecondarySalesReportGroupResp`, serta SQL alias `code`, atau konfirmasi bahwa branch lokal belum memuat SX-2172 enhancement.
4. `GREATEST(os.last_update, rs.last_update)` berpotensi menghasilkan null jika salah satu sisi null. Untuk numeric fields aman karena non-pointer + COALESCE, tetapi `last_update` boleh `null` sesuai expected minimal.
5. Query target dari prompt memakai source live table `sls.order`/`sls.order_detail`/`sls.return`/`sls.return_det`, sedangkan kode saat discovery memakai reporting facts `report.fact_orders`/`report.fact_returns`. Karena extract pipeline sudah mengisi facts dari source order/return, pilihan paling kecil-risiko adalah mempertahankan facts kecuali smoke DB membuktikan facts belum terisi untuk Juni 2026.
6. Untuk fallback missing year, `time.Now().Year()` pada tanggal discovery 2026 akan memenuhi Jira staging Juni 2026. Risiko: bila issue direplay di tahun lain, request `?month=6` tanpa year akan mengarah tahun saat runtime, bukan tahun data target. Ini sesuai instruksi user sebagai opsi fallback, tetapi FE tetap harus mengirim `&year=2026`.
7. Manual API test butuh valid staging token dari secure source; tidak boleh hardcode token di artifact atau source.

## Reuse candidates

- Reuse existing `resolveSecondaryDashboardYear`, `resolveSecondaryDashboardCustID`, raw SQL CTE summary, group query builder, dry-run repository test helpers, and service mock tests.
- Extend existing controller tests in `sales/controller/so_controller_test.go` instead of adding a new test file unless readability requires.
- Extend `model.SecondarySalesReportGroup`/`entity.SecondarySalesReportGroupResp` only if `code` must be restored/added.

## Source strategy keputusan

- Local project discovery: used and sufficient for implementation plan.
- Official docs/context7: skipped because no version-sensitive Go/GORM behavior is central; GORM dry-run patterns already exist locally.
- GitHub: skipped because no upstream repo behavior is needed.
- Web search: skipped; Jira/Google Docs are credentialed and user supplied relevant excerpts.
- Browser/screenshot: skipped; backend API bug, no visual parity work.
