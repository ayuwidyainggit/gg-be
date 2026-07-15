# Discovery — SX-2224 [BE] API Secondary Sales Group (code + name fix)

Task ID: `20260612-0916-sx-2224-secondary-sales-group-code`
Tanggal discovery: 2026-06-12 (Asia/Jakarta)

## Endpoint & peta kode aktual
- Route: `sales/controller/report_controller.go:97` → `reportRouteV1.Get("/secondary-sales/group", controller.SecondaryReportSalesGroup)`
- Service: `sales/service/report_service.go:1417` `SecondarySalesReportGroupSales`, switch `group_by` di `:1434-1443`
- Repository fungsi: `sales/repository/report_repository.go:1352-1390` (outlet/salesman/product_category/product)
- Query builder terpusat: `buildSecondarySalesReportGroupQuery(groupBy string)` di `report_repository.go` (select/join maps `:1255-1326`, template SQL `:1328-1349`)
- Model: `sales/model/report.go:436` `SecondarySalesReportGroup{ID, Code, Name, NetSales}`
- Response DTO: `sales/entity/report.go:263` `SecondarySalesReportGroupResp{ID, Code json:"code", Name, NetSales}`
- Tes regresi SQL: `sales/repository/report_repository_test.go:571-740`
- Tes branch service: `sales/service/report_service_test.go:100-128`

## Temuan kunci: scope inti SX-2224 SUDAH ada (via SX-2172)
Plan SX-2172 `.opencode/plans/20260609-1442-sx-2172-secondary-sales-dashboard.md` sudah:
- Menambah field `code` ke model + DTO + mapping service (`Code: r.Code` di `report_service.go:1453`).
- Mengisi `code`/`name` per branch dari master `mst.*` (bukan `report.dim_*`):
  - outlet: `COALESCE(mo.outlet_code,'') AS code`, name komposit `CONCAT_WS(' > ', e.emp_name, mo.outlet_name)`, `id = o.salesman_id`
  - salesman: `e.emp_code` / `e.emp_name`, `id = o.salesman_id`
  - product_category: `mpc.pcat_code` / `mpc.pcat_name`, `id = mpc.pcat_id`
  - product: `mp.pro_code` / `mp.pro_name`, `id = mp.pro_id` fallback `od.pro_id`
- Filter periode konsisten via rentang `o.invoice_date >= dateFrom AND < dateTo` (dateFrom dari `year`+`month`) di order & return branch untuk semua group_by → param `year` sudah efektif diterapkan.
- Authorization `cust_id` principal vs distributor di `resolveSecondaryDashboardCustIDs` (`report_service.go:1331-1361`), error `ErrUnauthorizedCustID`.

## Divergensi SX-2224 doc vs implementasi SX-2172 (perlu perubahan)
1. Branch `outlet`:
   - Doc SX-2224 expected: grouping outlet murni → `id = outlet_id`, `code = outlet_code`, `name = outlet_name` (contoh `id:1730, code:BM001, name:"Toko tosca"`).
   - Implementasi sekarang (SX-2172 "Sales by Customer"): `id = salesman_id`, `name = "emp_name > outlet_name"`.
   - Keputusan user (question gate 2026-06-12): **ikuti dokumen SX-2224** → ubah ke grouping outlet murni.
2. Branch `product` fallback nama:
   - Doc menunjukkan `name: "Product 10733"` (terlihat fallback sintetis dari id).
   - Cek DB lokal `ggn_scyllax`: `mst.m_product` `pro_id=10733` → `pro_code=JY1-002`, `pro_name="Jersey Manchester City FC"`, `cust_id=C260020001`, `pcat_id=77`.
   - Kategori `pcat_id=77` → `pcat_code=02`, `pcat_name="Jersey"`.
   - Kesimpulan: `"Product 10733"` adalah placeholder ilustratif; data master nyata punya nama. Tidak perlu fallback sintetis. Pertahankan master-priority `COALESCE(NULLIF(mp.pro_name,''),'')`.

## Bukti skema master (DB lokal ggn_scyllax)
- `mst.m_outlet(outlet_id, outlet_code, outlet_name, cust_id)` — contoh `1723 | BMI260004 | Toko biru | C260020001`.
- `mst.m_product_cat(pcat_id, pcat_code, pcat_name)` — `77 | 02 | Jersey`.
- `mst.m_employee(emp_id, emp_code varchar(10), emp_name varchar(150))`, PK `(cust_id, emp_id)`.
- `mst.m_product(pro_id, pro_code, pro_name, cust_id, pcat_id)`.

## Constraints repo
- Shell prefix `rtk` wajib di repo ini.
- Compose tidak running saat discovery (`rtk docker compose ps` kosong); query master dilakukan langsung ke Postgres lokal `localhost:5432 / postgres / ggn_scyllax`.
- Kontrak layer: Controller → Service → Repository → DB. Tenant join harus row-level (`mst.*.cust_id = sls.*.cust_id`).
- Validasi per service dir `sales`.

## Source strategy
- Repo-local evidence (kode + tes + plan SX-2172): utama.
- DB lokal `ggn_scyllax`: dipakai untuk validasi product/category/outlet master.
- Jira SX-2224 + doc `docs/Secondary Sales Report_BE.md`: konteks requirement.
- Official docs/GitHub/web: diskip, bug murni query SQL lokal + mapping internal, tidak bergantung library eksternal.
