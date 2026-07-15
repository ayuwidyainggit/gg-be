# Discovery SX-2085 — Monitoring Activity Collection Paid

Task ID: `20260528-1440-sx-2085-monitoring-collection-paid`

## File yang diinspeksi

- `docs/Monitoring Activity - BE.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/plans/20260522-1101-sx-2038-monitoring-detail-null.md`
- `pjp/router/live_monitoring.go`
- `pjp/controller/live_monitoring/get_detail_controller.go`
- `pjp/data/request/live_monitoring_request.go`
- `pjp/data/response/live_monitoring_response.go`
- `pjp/model/live_monitoring.go`
- `pjp/repository/live_monitoring/live_monitoring_repository.go`
- `pjp/repository/live_monitoring/get_detail_repository.go`
- `pjp/service/live_monitoring/get_detail_service.go`
- `pjp/service/live_monitoring/get_detail_service_test.go`

## Endpoint dan kontrak existing

- Endpoint detail: `GET /scylla-pjp/api/v1/monitoring_locations/details`.
- Query params detail: `emp_id` int, `distributor_id` optional, `date` string `YYYY-MM-DD`.
- Controller membungkus `data` sebagai array berisi satu object saat sukses.
- Response detail sudah punya field `collection` sebagai `[]CollectionData`.
- `CollectionData` punya `outlet_id`, `outlet_code`, `outlet_name`, `collection_total`, semua pointer.
- Service saat ini hardcode `collection := []response.CollectionData{}` dan `collection_summary` selalu `none`.

## Pola project yang ditemukan

- Module target: `pjp`, Gin-based.
- Layer wajib: Controller → Service → Repository → DB.
- Repository pakai GORM query builder dengan `Table`, `Select`, `Joins`, `Where`, `Group`, `Order`, `Find`.
- Detail section lain sudah punya pola repository row → response mapping:
  - `GetSales` → `[]model.SalesRow` → `[]response.SalesData`
  - `GetReturns` → `[]model.ReturnRow` → `[]response.ReturnData`
  - `GetExpenses` → `[]model.ExpenseRow` → `[]response.ExpenseData`
  - `GetShipments` → `[]model.ShipmentRow` → `[]response.ShipmentData`
- `GetMonitoringDetail` sudah resolve `salesmanCustID` dari `GetSalesmanCustID`, lalu memakai `targetCustIDs := []string{salesmanCustID}` untuk transactional sections.

## Reuse candidates

- Tambah `GetCollections(ctx, tx, custIDs, date, empID)` sejajar `GetSales` / `GetReturns` / `GetShipments`.
- Tambah `model.CollectionRow` sejajar `SalesRow` dan `ReturnRow`.
- Reuse response contract existing `response.CollectionData.CollectionTotal`.
- Reuse `buildVisitSummary(len(collection))` untuk `collection_summary` count/status.
- Reuse `targetCustIDs` dari salesman cust_id agar tenant filter tetap ketat.

## Root cause indikatif

- Field `collection` sudah ada di response, tetapi service mengisi array kosong permanen.
- Repository `pjp` belum punya query collection detail dari `acf.deposit` / `acf.deposit_payment`.
- Search dalam `pjp` tidak menemukan pemakaian `deposit_payment`, `deposit_detail`, `acf.deposit`, atau `payment_amount`.

## Query dan grain risk

- Jira reference memakai join `acf.deposit` → `acf.deposit_detail` → `sls.order` → `acf.deposit_payment` lalu `SUM(dp.payment_amount)` per outlet.
- Risiko duplikasi ada bila satu `deposit_no,cust_id` punya banyak `deposit_detail` dan banyak row `deposit_payment`.
- Plan harus pakai agregasi aman:
  - aggregate `acf.deposit_payment` per `deposit_no,cust_id` dulu,
  - ambil `deposit_outlet` distinct per `deposit_no,cust_id,outlet_id`,
  - sum `COALESCE(payment_amount,0)` per outlet.

## Commands/docs dicek

- `rtk docker compose -f docker-compose.yml ps` dijalankan; output tool hanya menampilkan warning compose `version` obsolete, status service tidak terlihat karena output terpotong/aneh.
- Local docs dicek; `docs/Monitoring Activity - BE.md` line 582+ mendokumentasikan endpoint detail dan line 645+ mendokumentasikan enhancement collection.
- `.opencode/docs/ARCHITECTURE.md` dicek untuk layering, tenant, schema constraints.
- External docs/context7/GitHub/brave/browser tidak diperlukan; bug internal query/response di repo lokal.

## Constraints

- Jangan ubah source saat planner aktif.
- Jangan hardcode `2026-05-28` atau `421`.
- Date dari request detail adalah string `YYYY-MM-DD`, bukan epoch.
- Filter valid wajib `d.deposit_date = req.Date`, `d.emp_id = req.EmpID`, `d.collection_no IS NOT NULL`.
- Amount wajib dari `acf.deposit_payment.payment_amount`, bukan `sls.order.total`.
- Tenant filter wajib pakai `cust_id IN targetCustIDs` di tabel transaksi relevan.

## Risiko

- Staging DB/token belum tersedia di sesi ini; SQL manual dan before/after API evidence harus dilakukan saat implementasi.
- `acf.deposit_payment` grain belum dibuktikan di DB; query aman wajib untuk cegah duplicate sum.
- Existing response tidak punya top-level `collection.total_paid`; menambah object bisa breaking. Rekomendasi awal: isi array existing `collection[].collection_total`; FE bisa sum. Top-level total hanya kalau FE sepakat.
- `outlet_id` dalam `deposit_outlet` bisa muncul lebih dari satu untuk satu deposit bila invoice lintas outlet; plan perlu validate data grain dan catat bila business rule harus split amount.
