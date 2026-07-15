# SX-1989 Discovery — Manage Prices tidak ter-update untuk produk principal lama

Tanggal discovery: 2026-05-20 13:43 WIB
Service: `master/` (Fiber, modul Go terpisah)

## File yang diperiksa

- `master/controller/product_controller.go` (handler `GET /master/v1/products`, switch `mode=lookup_dist_price`).
- `master/service/product_service.go` (`LookupDistPrice` di sekitar baris 4163-4187, store/dist propagation di 4243-4296, 4350-4385).
- `master/repository/product_repository.go`:
  - `FindAllByCustIdLookupDistPrice` (1194-1324) — dipakai saat tidak ada `distributor_id` JWT, hanya filter `p.distributor_id IS NULL` dan tidak baca `mst.m_transaction_price` sama sekali.
  - `FindAllByDistributorLookupDistPrice` (1326-1573) — dipakai saat user distributor (kasus SX-1989). Pakai CTE `priced_products` yang menggabungkan `mst.m_transaction_price` via `pricing_lookup_pro_id = COALESCE(NULLIF(p.parent_pro_id, 0), p.pro_id)` (1354, 1521-1545).
  - `resolveDistPriceGroupID` (141-167).
- `master/controller/m_price_controller.go` (route `/v1/prices`, RMQ subscriber `processPublishManagePriceMessage`).
- `master/service/m_price_service.go`:
  - `Store/Update/Publish/PublishByRMQ` (128-324).
  - `prepareCreateRequest`/`prepareUpdateRequest` (461-505) — ambil snapshot via `FindOneProductSnapshotByProID`.
  - `applyPublishedProductPrices` (525-535) — untuk principal panggil `UpdatePrincipalAssignedProductPrices(price.ProID, …)`; untuk distributor panggil `UpdateDistributorProductPrices(custID, distID, price.ProID, …)`.
  - `syncTransactionPrices` (537-567) — insert/update `mst.m_transaction_price` dengan `pro_id = price.ProID`.
  - `setupTransactionPriceData/Update` (1302-1347) — pasang `Source = 10`, `Coverage = price.Coverage`, `DistributorID`, dst.
- `master/repository/m_price_repository.go`:
  - `FindOneProductSnapshotByProID/ByCode` (464-526) — proyeksikan `parent_pro_id` (kolom ada).
  - `UpdatePrincipalAssignedProductPrices` (604-645) — update `mst.m_product` di mana `parent_pro_id = ?`.
  - `UpdateDistributorProductPrices` (647-682) — update `mst.m_product` per `cust_id`, `pro_id`, opsional `distributor_id`.
  - `PublishByRMQ` (684-722).
- `master/model/m_price.go`, `master/model/m_transaction_price.go`, `master/entity/m_price.go`, `master/entity/product.go` — struct request/response/DTO.
- `master/repository/m_price_repository_test.go`, `master/service/m_price_service_test.go` — pola test pakai `DATA-DOG/go-sqlmock` + stub repository.
- Linked code untuk konteks:
  - `mobile/repository/product.go` — pakai `manage_minimum_price`, bukan jalur ini.
  - `inventory/service/wh_trf_service.go` (191) — konsumen `mode=lookup_dist_price` lain.
  - `sales/repository/promotion_repository.go` (724-735) — pakai pola CTE serupa.
  - `master/service/m_sp_price_service.go`, `master/service/dist_price_service.go` — jalur lain yang tulis `m_transaction_price`.

## Pola proyek yang dipakai

- Layering Controller → Service → Repository (lihat `m_price_*` dan `product_*`).
- Multi-tenant: `cust_id` = distributor child; `parent_cust_id` = principal; sebagian repo memakai placeholder string langsung yang sudah lama jadi pola di repo ini, jadi perubahan harus ikut pola sambil hati-hati ke escaping (sudah ada di kode).
- Transaksi/RMQ: `Store` Manage Prices → enqueue / immediate publish via `RMQ_MANAGE_PRICE_CREATE_EVENT` → `PublishByRMQ` → `applyPublishedProductPrices` + `syncTransactionPrices` + `MPriceRepository.PublishByRMQ`.
- Test pattern wajib: `sqlmock` + regex `(?s)…` untuk repo, stub manual untuk service.

## Reuse candidates

- `MPriceProductSnapshot.ParentProID` sudah dibawa di snapshot, tinggal dipakai untuk routing target update + insert `m_transaction_price`.
- `setupTransactionPriceData` / `setupTransactionPriceDataUpdate` sudah satu titik untuk override `pro_id` agar lookup match dengan jalur baca.
- `UpdatePrincipalAssignedProductPrices` sudah pakai `parent_pro_id`, jadi update `mst.m_product` aman; yang ompong adalah update `m_transaction_price`.

## Akar masalah

### Hipotesis awal dari pembacaan kode

1. Lookup `mode=lookup_dist_price` (jalur distributor) **membaca** harga dari `mst.m_transaction_price` dengan kunci:
   - `mtp_mg_pr.pro_id = paged_products.pricing_lookup_pro_id`,
   - `pricing_lookup_pro_id = COALESCE(NULLIF(p.parent_pro_id, 0), p.pro_id)`,
   - `mtp_mg_pr.cust_id = parent_cust_id`,
   - `distributor_id = (CASE WHEN coverage = 'N' THEN 0 ELSE jwtDistributorID END)` atau `price_group_reff = distPriceGroupID`,
   - `start_date <= order_date`.
2. Saat principal melakukan **Manage Prices update**:
   - `request.ProID` = `pro_id` produk **principal** (parent), `request.CustID = parent_cust_id`.
   - `syncTransactionPrices` selalu pakai `price.ProID` (principal pro_id) sebagai `pro_id` baris baru di `mst.m_transaction_price`.
3. Ini sempat mengarah ke hipotesis mismatch parent-child legacy (`parent_pro_id` kosong/null) untuk produk principal lama.

### Root cause final setelah validasi staging DB

Hipotesis di atas **tidak cukup** untuk kasus JY1-005. Audit staging menunjukkan:
- child row distributor 102 untuk `JY1-005` **sudah punya** `parent_pro_id = 8429`;
- `mst.m_transaction_price` child generic juga **sudah ada**:
  - `cust_id = C260020001`
  - `pro_id = 8457`
  - `start_date = 2026-05-16`
  - `purch_price1 = 50000`
  - `sell_price1 = 150000`
- parent generic row lama juga ada:
  - `cust_id = C26002`
  - `pro_id = 8429`
  - `start_date = 2026-05-13`
  - `purch_price1 = 60000`
  - `sell_price1 = 1000000`

**Bug aktual:** `FindAllByDistributorLookupDistPrice` hanya membaca:
- child outlet-specific row (`mtp`) dengan `outlet_id = current outlet`, lalu
- parent generic row (`mtp_mg_pr`, `mtp_mg_pr_sell`) via `pricing_lookup_pro_id`.

Query **tidak membaca child generic row** (`cust_id = child`, `pro_id = child`, `COALESCE(outlet_id,0)=0`, `start_date <= order_date`).

Akibat untuk `order_date = 2026-05-16`:
- child generic row yang benar (`50000 / 150000`, start_date 2026-05-16) diabaikan,
- parent old row (`60000 / 1000000`, start_date 2026-05-13) dipilih,
- lookup mengembalikan harga lama walau Manage Prices distributor/child sudah punya row lebih baru.

### Kesimpulan final

Root cause utama adalah **precedence query salah** pada lookup distributor:
- source kandidat generic hanya parent,
- child generic tidak ikut dipertimbangkan,
- sehingga latest applicable price di child tenant kalah oleh parent old price.

Fix final: `UNION ALL` child generic + parent generic, lalu `ORDER BY start_date DESC, scope_priority ASC LIMIT 1`.
Child diberi `scope_priority = 0`, parent `scope_priority = 1`, jadi child menang jika tanggal sama, dan latest date menang jika berbeda.

Backfill `parent_pro_id` tetap berguna untuk 3 row legacy lain yang fixable, tetapi **bukan** root cause utama kasus JY1-005.

## Konstrain & risiko

- Repo ini punya kebijakan keras: write logic harus di service layer, repo cuma data access. Perubahan boleh tambah method repo, tapi logika routing harus di service.
- Banyak query pakai string concat. Hindari menambah parameter baru tanpa rebind kalau repo terkait sudah pakai placeholder.
- Backfill data `parent_pro_id` perlu dipertimbangkan; kalau dilakukan, harus reversible dan mempertimbangkan tenant.
- `mst.m_transaction_price` punya jalur lain (`m_sp_price_service`, `dist_price_service`) yang juga insert. Hindari ubah skema atau makna kolom.
- RMQ subscriber `processPublishManagePriceMessage` jalan goroutine; perubahan harus ramah idempotensi.
- TMS/PJP Gin pakai jalur berbeda; perubahan ini cukup di `master/` saja.

## Pertanyaan yang masih terbuka

1. Untuk produk principal lama tanpa child row, perilaku yang diinginkan saat distributor melihat lookup: membuat child on-the-fly (mahal) ATAU lookup memang harus jatuh ke baris principal (level 0)?
2. Ada toleransi backfill data? Atau fix harus runtime-only?
3. Apakah `dist_price_grp_id` token user sengaja `0`? Param query mengirim `dist_price_group_id=2`. Repo override pakai `m_distributor.dist_price_grp_id`, perlu konfirmasi distributor 102 punya `dist_price_grp_id` valid (bukan 0).
