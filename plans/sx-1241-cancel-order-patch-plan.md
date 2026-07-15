# SX-1241 Cancel Order Stock Consistency Patch Plan Finalized

## Context
Defect scope:
- On Cust Order not updated correctly after cancel
- Available Stock becomes incorrect after cancel
- possible distributor or warehouse stock summary inconsistency

Confirmed code path:
- Main cancel path: `sales/service/order_service.go` via `BulkUpdateStatus`
- Reverse stock path: `sales/repository/stock_repository.go` via `CancelSalesStockUpdates`
- Detail order endpoint for UI `Final Order` tab: `GET /sales/v2/orders/:ro_no`
- Detail endpoint controller path: `sales/controller/order_controller.go` via `DetailV2`
- Detail endpoint service path: `sales/service/order_service.go` via `DetailV2`
- Current cancel implementation updates:
  - `inv.stock`
  - `inv.warehouse_stock`
- Current implementation does not yet show explicit synchronization for projection or summary layers consumed by UI
- Current reversal basis is derived from historical `SO` ledger rows, not validated against latest detail or final state
- There are still multiple cancel paths and they are not unified yet
- Discovery update for UI `Available Stock`:
  - response field comes from `data.details_final.normal[].qty1_stok`, `qty2_stok`, `qty3_stok`
  - these fields are not primarily read from persisted `sls.order_detail.qty*_stok`
  - `DetailV2` recomputes them from warehouse stock source plus line quantity displayed in tab
  - warehouse stock source is fetched through `FindWarehouseStockByWhIdAndProIds`

---

## 1. Review Kekuatan Plan

Bagian plan yang sudah kuat dan aman untuk dijalankan minggu ini:

- **Tetap mempertahankan orchestration existing** di `BulkUpdateStatus`.
  Ini tepat untuk patch minggu ini karena blast radius kecil dan tidak memaksa rewrite cancel engine.

- **Tetap mempertahankan reverse ledger existing** di `CancelSalesStockUpdates`.
  Ini keputusan yang aman karena perilaku reverse yang sudah berjalan untuk `inv.stock` dan `inv.warehouse_stock` tidak dirombak besar-besaran.

- **Membatasi reconcile hanya ke affected key `(cust_id, wh_id, pro_id)`**.
  Ini penting untuk production safety, karena menghindari recompute tenant-wide atau warehouse-wide yang mahal dan berisiko side effect.

- **Mengubah urutan transaksi menjadi status update sebelum reconcile**.
  Ini keputusan yang benar untuk patch ini, terutama jika projection UI bergantung pada `order.data_status`.

- **Menambahkan sanity guard** sebelum cancel melanjutkan flow.
  Ini penting untuk mencegah silent corruption atau false success.

- **Mewajibkan discovery source-of-truth UI sebelum implementasi reconcile**.
  Ini sangat benar; tanpa ini tim bisa menulis ke table yang salah dan defect tetap terlihat.

- **Membedakan patch cepat vs proper fix**.
  Ini baik dari sisi delivery dan governance: patch minggu ini fokus self-healing summary, proper fix belakangan fokus unifikasi cancel flow dan validasi basis.

Kesimpulan reviewer: fondasi plan sudah **cukup aman**, **cukup implementable**, dan **cukup minim regression** untuk patch minggu ini, selama keputusan final di bawah dijalankan tanpa ambigu.

---

## 2. Final Design Decisions

Keputusan desain final paling aman untuk patch minggu ini:

### 2.1 Definisi `active reserving detail`
Patch-week definition:
- `item_type = 1`
- detail normal reservable saja
- bukan promo item
- bukan reward item
- quantity reserving dianggap aktif bila `qty_final > 0`
- bila `qty_final` tidak reliable pada jalur tertentu, fallback ke `qty > 0`

Rule ini **hanya untuk sanity guard**, bukan untuk mengganti basis reverse yang tetap berasal dari ledger historis.

### 2.2 Promo dan reward item
- Promo item **tidak** dianggap `active reserving detail`
- Reward item **tidak** dianggap `active reserving detail`
- Promo atau reward item **tidak boleh** menjadi alasan fail-fast `basis kosong`

### 2.3 Delta-based vs rebuild-based contract
Field yang **tetap delta-based**:
- `inv.stock`
- `inv.warehouse_stock.qty`
- `inv.warehouse_stock.qty_on_order`

Field yang **harus authoritative rebuild**:
- projection UI-facing `On Cust Order`
- projection UI-facing `Available Stock`
- distributor summary **hanya jika** discovery membuktikan summary itu bagian dari failing read path

### 2.4 Owner final `qty_on_order`
Keputusan final patch-week:
- `inv.warehouse_stock.qty_on_order` tetap dimiliki oleh existing reverse flow di `CancelSalesStockUpdates`
- Rebuild step **tidak boleh overwrite** `qty_on_order`
- Tujuan utama patch adalah memperbaiki layer projection yang dibaca UI, bukan menambahkan second writer ke field warehouse summary yang sudah punya existing owner

### 2.5 Distributor summary default rule
- `RebuildDistributorStockSummary` default **no-op / not called**
- Method ini hanya boleh dipanggil bila discovery source-of-truth membuktikan distributor summary memang bagian dari endpoint atau screen yang menunjukkan defect
- Jika discovery belum selesai atau hasilnya negatif, distributor summary **tidak disentuh**

### 2.6 Kapan reconcile dijalankan
Reconcile dijalankan dalam urutan berikut:
1. reverse stock existing
2. update order status ke cancelled
3. rebuild projection UI-facing yang confirmed
4. commit

### 2.7 Kapan fail-fast vs self-heal
**Fail-fast** bila:
- order tidak ditemukan
- `data_status` `nil`
- transition ke cancel invalid
- `basis kosong` dan ada `active reserving detail`
- `basis ada` tetapi `affected key` kosong
- reconcile gagal

**Self-heal** bila:
- sebagian `ref_det_id` sudah memiliki reverse row
- repeated cancel atau retry butuh memastikan projection sinkron
- projection mismatch masih bisa diperbaiki dengan rebuild authoritative per affected key

---

## 3. Final Ownership Contract

| Komponen data | Single writer or owner | Update mechanism | Authoritative source | Boleh ditulis method lain? |
|---|---|---|---|---|
| `inv.stock` | `CancelSalesStockUpdates` | Delta append-only reverse rows | Historical stock ledger | Tidak |
| `inv.warehouse_stock.qty` | `CancelSalesStockUpdates` | Delta upsert existing logic | Reverse cancel mutation | Tidak |
| `inv.warehouse_stock.qty_on_order` | `CancelSalesStockUpdates` | Delta existing logic | Reverse cancel mutation against reservation ledger basis | Tidak |
| On Cust Order UI projection | `RebuildOnCustOrderProjection` atau method projection khusus yang dikunci setelah discovery | Rebuild authoritative | Confirmed UI source path | Tidak |
| Available Stock UI projection | `RebuildAvailableStockProjection` | Rebuild authoritative | Confirmed UI formula or projection source | Tidak |
| Distributor stock summary | `RebuildDistributorStockSummary` hanya jika discovery confirmed | Rebuild authoritative | Confirmed distributor summary source | Tidak |

### Contract rules
- Satu field hanya boleh punya satu owner pada patch ini
- Rebuild method tidak boleh overwrite field yang sudah dimiliki delta writer existing
- `qty_on_order` dianggap **locked ownership** pada `CancelSalesStockUpdates`
- Projection rebuild hanya menyentuh target UI-facing yang sudah confirmed dari discovery

---

## 4. Final Rule untuk `qty_on_order`

Keputusan final:
- Patch minggu ini **tetap membiarkan `CancelSalesStockUpdates` menjadi owner final `inv.warehouse_stock.qty_on_order`**
- Rebuild step **jangan menyentuh** `qty_on_order`
- Pengecualian tidak berlaku untuk patch minggu ini; bila discovery membuktikan UI membaca field ini langsung, itu tetap tidak mengubah ownership patch-week kecuali ada review ulang eksplisit

Alasan production-safe:
- field ini sudah punya existing writer yang aktif
- menambah overwrite kedua di rebuild akan menciptakan dual-writer risk
- dual-writer pada field reservation summary adalah sumber regresi paling berbahaya pada patch inventory cepat
- lebih aman menjaga `qty_on_order` tetap single-writer, lalu memperbaiki projection UI pada layer projection yang memang dibaca UI

---

## 5. Final Rule untuk Distributor Summary

Rule final:
- `RebuildDistributorStockSummary` **wajib dipanggil hanya jika** discovery source-of-truth membuktikan distributor summary dipakai oleh endpoint atau screen yang menunjukkan defect
- `RebuildDistributorStockSummary` **tidak boleh dipanggil** bila:
  - discovery belum selesai
  - discovery menunjukkan UI tidak membaca distributor summary
  - summary hanya bersifat reporting downstream dan bukan bagian read path defect
- Untuk patch minggu ini, method itu **default tidak dipanggil sama sekali** sampai discovery selesai dan confirmed

Rule untuk mencegah scope creep:
- tambahkan flag keputusan di plan: `distributor_summary_reconcile = false` secara default
- flag hanya boleh berubah ke `true` setelah discovery table terisi dan direview
- tidak boleh ada generic reconcile semua summary
- bila distributor summary tidak confirmed, patch selesai tanpa menyentuh area itu

---

## 6. Ready-for-Coding Checklist

Checklist final sebelum implementasi dimulai:

- [ ] Discovery source UI untuk `On Cust Order` selesai sampai endpoint, controller, service, repository, dan table atau view level
- [ ] Discovery source UI untuk `Available Stock` selesai sampai endpoint, controller, service, repository, dan table atau view level
- [ ] Diputuskan apakah distributor summary bagian dari failing read path
- [ ] Definisi `active reserving detail` ditulis eksplisit di plan sesuai keputusan final
- [ ] Disepakati bahwa promo dan reward item dikecualikan dari `active reserving detail`
- [ ] Ownership `inv.warehouse_stock.qty_on_order` dikunci final pada `CancelSalesStockUpdates`
- [ ] Field ownership semua komponen tidak overlap
- [ ] Affected key dedup memakai deterministic sorting
- [ ] Rollback semantics tertulis: reconcile failure = full transaction rollback
- [ ] Semua repository methods yang dipakai cancel dan reconcile menerima `ctx` dan memakai `repo.model(ctx)` pattern
- [ ] Reconcile target hanya mencakup projection yang confirmed dibaca UI
- [ ] Test matrix minimal untuk unit, integration, regression, contract, rollback sudah disetujui developer dan reviewer

---

## 6A. Discovery Update Locked for `Available Stock`

Hasil discovery terbaru yang sudah dianggap locked untuk patch minggu ini:

### Endpoint dan flow backend
- UI detail order page untuk tab `Final Order` dibackup oleh endpoint `GET /sales/v2/orders/:ro_no`
- Route ada di `sales/controller/order_controller.go` pada `Route`
- Handler controller ada di `DetailV2`
- Service pembentuk response ada di `sales/service/order_service.go` pada `DetailV2`
- DTO response yang mengandung `details_final` ada di `sales/entity/order.go` pada `OrderResponse`

### Source field `details_final.normal[].qty*_stok`
- Detail order dasar dibaca dari `sls.order_detail` lewat repository `FindDetail`
- Warehouse stock current dibaca lewat repository `FindWarehouseStockByWhIdAndProIds`
- `DetailV2` lalu menghitung ulang `qty1_stok`, `qty2_stok`, `qty3_stok`
- Formula efektif endpoint saat ini:
  - `available_stock_display = warehouse_stock_current + qty_line_yang_sedang_ditampilkan`
- Untuk tab `Final Order`, quantity line yang dipakai adalah quantity final row yang telah dikonversi ke unit display

### Implikasi arsitektural
- Source authoritative untuk `Available Stock` pada detail page saat ini adalah warehouse stock query plus service-side recomputation
- Persisted field `sls.order_detail.qty1_stok`, `qty2_stok`, `qty3_stok` bukan source-of-truth utama untuk endpoint detail v2
- Karena itu patch minggu ini tidak boleh berasumsi bahwa update persisted `qty*_stok` saja akan memperbaiki UI
- Patch harus divalidasi terhadap contract endpoint `DetailV2`, bukan hanya ledger internal atau kolom snapshot detail

### Risiko bug yang paling mungkin setelah discovery ini
- reverse cancel tidak mengembalikan warehouse stock source ke nilai benar
- atau formula `warehouse stock + final line qty` tetap membuat response salah ketika order sudah `CANCELLED`
- atau cancel basis historis tidak sinkron dengan qty final terbaru setelah edit qty lalu print proforma

## 7. Discovery Output Template

Gunakan template ini dan wajib diisi sebelum implementasi reconcile:

| UI metric | Screen or component | Endpoint | Controller | Service | Repository | Table or view | Depends on order status Y or N | Depends on ledger Y or N | Needs reconcile in this patch Y or N | Notes |
|---|---|---|---|---|---|---|---|---|---|---|
| On Cust Order |  |  |  |  |  |  |  |  |  |  |
| Available Stock | Detail Order page Final Order tab | `GET /sales/v2/orders/:ro_no` | `DetailV2` | `DetailV2` | `FindDetail` + `FindWarehouseStockByWhIdAndProIds` | `sls.order_detail` + warehouse stock table | Y | Y | Y | Response `qty*_stok` dihitung ulang dari warehouse stock current + qty line final yang sedang ditampilkan |
| Distributor Stock Summary |  |  |  |  |  |  |  |  |  |  |
| Warehouse Stock Summary | Detail Order page Final Order tab indirect source | `GET /sales/v2/orders/:ro_no` | `DetailV2` | `DetailV2` | `FindWarehouseStockByWhIdAndProIds` | warehouse stock table | N | Y | Y | Bukan field UI langsung, tetapi source dasar untuk `Available Stock` |

### Wajib terisi sebelum coding
- endpoint nyata yang dipakai UI
- query atau source table yang dipakai
- apakah metric bergantung pada `data_status`
- apakah metric perlu reconcile pada patch minggu ini

---

## 8. Final Strengthened Patch Flow

### Final patch flow recommendation

1. **Load order header**
   - apa yang dibaca:
     - order by `ro_no`, `cust_id`
   - apa yang divalidasi:
     - order ditemukan
     - `data_status` tidak `nil`
     - data dasar seperti `ro_date`, `cust_id`, current status tersedia
   - apa yang diupdate:
     - belum ada update
   - kenapa step ini ada:
     - semua keputusan cancel bergantung pada state order sekarang
   - kondisi error atau rollback:
     - order tidak ditemukan → error
     - `data_status nil` → error

2. **Validate cancel transition**
   - apa yang dibaca:
     - current `data_status`
   - apa yang divalidasi:
     - jika sudah cancelled → idempotent success
     - jika transition invalid → fail-fast
   - apa yang diupdate:
     - belum ada update
   - kenapa step ini ada:
     - mencegah illegal transition dan duplicate reverse
   - kondisi error atau rollback:
     - invalid transition → rollback

3. **Load order details untuk sanity scope**
   - apa yang dibaca:
     - reservable detail order
   - apa yang divalidasi:
     - tentukan apakah ada `active reserving detail`
   - apa yang diupdate:
     - belum ada update
   - kenapa step ini ada:
     - dipakai untuk mendeteksi state abnormal saat basis kosong
   - kondisi error atau rollback:
     - read detail gagal → rollback

4. **Load cancel basis dari `GetCancelStockBasis`**
   - apa yang dibaca:
     - basis ledger historis `SO`
   - apa yang divalidasi:
     - basis boleh kosong hanya bila tidak ada `active reserving detail`
   - apa yang diupdate:
     - belum ada update
   - kenapa step ini ada:
     - reverse patch-week tetap menggunakan source existing yang paling aman
   - kondisi error atau rollback:
     - query basis gagal → rollback

5. **Run sanity guards**
   - apa yang dibaca:
     - detail order + basis
   - apa yang divalidasi:
     - `basis kosong + active reserving detail ada` → fail-fast
     - `basis ada tetapi row valid nol` → fail-fast
   - apa yang diupdate:
     - belum ada update
   - kenapa step ini ada:
     - mencegah silent inconsistent cancel
   - kondisi error atau rollback:
     - guard gagal → rollback total

6. **Build affected keys**
   - apa yang dibaca:
     - basis rows
   - apa yang divalidasi:
     - dedup berdasarkan `(cust_id, wh_id, pro_id)`
     - sorting deterministic stabil
     - hasil tidak boleh kosong bila basis ada
   - apa yang diupdate:
     - belum ada update
   - kenapa step ini ada:
     - reconcile harus minimal scope, stabil, dan terprediksi
   - kondisi error atau rollback:
     - affected key kosong → rollback

7. **Reverse stock dengan `CancelSalesStockUpdates`**
   - apa yang dibaca:
     - basis rows
     - existing reverse refs
   - apa yang divalidasi:
     - row yang sudah reversed tidak ditulis ulang
     - partial reverse existing boleh self-heal
   - apa yang diupdate:
     - append reverse rows ke `inv.stock`
     - update delta `inv.warehouse_stock.qty`
     - update delta `inv.warehouse_stock.qty_on_order`
   - kenapa step ini ada:
     - existing reverse engine paling aman dipertahankan untuk patch minggu ini
   - kondisi error atau rollback:
     - reverse write gagal → rollback total

8. **Update order status ke cancelled**
   - apa yang dibaca:
     - target order header
   - apa yang divalidasi:
     - update hanya untuk order target dan tenant benar
   - apa yang diupdate:
     - `data_status = CANCELLED`
   - kenapa step ini ada:
     - projection reconcile harus membaca final business state
   - kondisi error atau rollback:
     - update status gagal → rollback total

9. **Run projection reconcile untuk confirmed UI-facing targets**
   - apa yang dibaca:
      - affected keys
      - source-of-truth projection yang sudah confirmed
      - state order terbaru dalam transaction yang sama
   - apa yang divalidasi:
      - hanya reconcile target yang confirmed dari discovery yang dipanggil
      - distributor summary tetap skip bila belum confirmed
      - tidak ada method yang overwrite `qty_on_order`
      - untuk `Available Stock`, validasi harus mengacu ke source yang dipakai `DetailV2`, yaitu warehouse stock current dan formula response detail
   - apa yang diupdate:
      - rebuild authoritative `On Cust Order`
      - source yang memengaruhi `Available Stock` pada endpoint detail
      - optionally rebuild distributor summary hanya bila confirmed
   - kenapa step ini ada:
      - defect paling mungkin berada pada projection layer yang dibaca UI
- kondisi error atau rollback:
     - satu reconcile method gagal → rollback total

10. **Commit transaction**
   - apa yang dibaca:
     - state transaction
   - apa yang divalidasi:
     - seluruh step sebelumnya sukses
   - apa yang diupdate:
     - commit transaction
   - kenapa step ini ada:
     - menjamin ledger, status, dan projection committed secara atomik
   - kondisi error atau rollback:
     - commit gagal → return error
     - tidak boleh ada swallow error sebelumnya

### Aturan wajib selama flow
- semua repository call harus pakai `txCtx`
- jangan ada swallow error di reconcile step
- jangan ada field yang ditulis oleh delta logic dan rebuild logic sekaligus tanpa kontrak eksplisit
- affected key processing order harus stabil

---

## 9. Final Test Matrix

### A. Unit tests

#### `TestBulkUpdateStatus_CancelOrder_ReconcileAfterStatusUpdate`
- tujuan: memastikan reconcile dijalankan setelah status update
- setup: mock repository dan transaction, cancel request valid
- expected result: urutan call reverse → update status → reconcile
- risk yang dicegah: projection membaca status lama

#### `TestBuildAffectedCancelKeys_DedupAndSortStable`
- tujuan: memastikan affected key dedup dan sorting stabil
- setup: basis rows acak dengan duplicate SKU
- expected result: unique sorted keys
- risk yang dicegah: flaky behavior dan flaky tests

#### `TestCancelOrder_FailFastWhenBasisEmptyButActiveReservingExists`
- tujuan: memastikan guard abnormal state bekerja
- setup: detail aktif ada, basis kosong
- expected result: error eksplisit
- risk yang dicegah: silent corruption

#### `TestCancelOrder_SelfHealWhenPartialReverseAlreadyExists`
- tujuan: memastikan partial reverse existing tetap aman
- setup: sebagian `ref_det_id` sudah reversed
- expected result: hanya remaining rows yang ditulis, reconcile tetap jalan
- risk yang dicegah: retry inconsistency

#### `TestCancelOrder_IdempotentWhenAlreadyCancelled`
- tujuan: repeated request aman
- setup: order status sudah cancelled
- expected result: no duplicate reverse, no error fatal
- risk yang dicegah: double reverse

#### `TestCancelOrder_SkipPromoRewardInActiveReservingGuard`
- tujuan: promo dan reward item tidak memicu guard salah
- setup: hanya promo atau reward detail aktif
- expected result: tidak dianggap `active reserving detail`
- risk yang dicegah: false fail-fast

### B. Integration tests

#### `Integration_CancelOrder_NormalReservationFlow`
- tujuan: baseline happy path cancel
- setup: SO reserving normal
- expected result: ledger reverse benar, warehouse stock benar, projection benar
- risk yang dicegah: patch merusak flow utama

#### `Integration_CancelOrder_AfterEditQty`
- tujuan: cover repro utama
- setup: edit qty lalu cancel
- expected result: final stock state sinkron
- risk yang dicegah: stale reservation mismatch

#### `Integration_CancelOrder_AfterEditPrice`
- tujuan: memastikan price-only edit tidak merusak stock state
- setup: edit price lalu cancel
- expected result: quantity-related projection tetap benar
- risk yang dicegah: state drift setelah non-qty edit

#### `Integration_CancelOrder_AfterPrintProforma`
- tujuan: cover state transition proforma
- setup: print proforma lalu cancel
- expected result: projection tetap benar
- risk yang dicegah: jalur state-specific defect

#### `Integration_CancelOrder_MultiOrdersSameSKU`
- tujuan: cancel satu order tidak merusak reservation order lain
- setup: dua SO dengan SKU sama, satu dibatalkan
- expected result: hanya efek order target yang hilang
- risk yang dicegah: over-release reservation

#### `Integration_CancelOrder_ReverseRowsAlreadyExistPartially`
- tujuan: memastikan self-healing partial reverse
- setup: seed sebagian reverse row sebelum cancel
- expected result: hasil akhir tetap konsisten setelah commit
- risk yang dicegah: interrupted transaction scenarios

### C. Regression test SX-1241

#### `Regression_SX1241_EditQty_PrintProforma_FinalOrder_Cancel`
- tujuan: mirror reproduksi ticket seakurat mungkin
- setup:
  1. edit qty
  2. print proforma
  3. masuk final order
  4. cancel
- expected result:
  - reverse rows benar
  - `inv.warehouse_stock` sinkron
  - `On Cust Order` sinkron
  - `Available Stock` sinkron
- risk yang dicegah: defect re-open

#### `Regression_SX1241_EditPrice_PrintProforma_FinalOrder_Cancel`
- tujuan: cover varian price edit
- setup:
  1. edit price
  2. print proforma
  3. masuk final order
  4. cancel
- expected result: projection stock tetap benar
- risk yang dicegah: hidden coupling qty vs price edit path

### D. Contract tests untuk UI source-of-truth

#### `Contract_OnCustOrder_UIEndpoint_MatchesDocumentedSource`
- tujuan: mengunci source-of-truth yang ditemukan saat discovery
- setup: panggil endpoint UI dan query source documented path
- expected result: nilai sama
- risk yang dicegah: backend patch source yang salah

#### `Contract_AvailableStock_UIEndpoint_MatchesDocumentedSource`
- tujuan: mengunci formula atau table source available stock
- setup: endpoint UI + query authoritative source
- expected result: nilai sama
- risk yang dicegah: projection mismatch tidak terdeteksi

#### `Contract_DetailV2_FinalOrder_QtyStok_MatchesWarehouseStockPlusDisplayedFinalQty`
- tujuan: mengunci contract aktual endpoint `DetailV2` untuk `details_final.normal[].qty1_stok`, `qty2_stok`, `qty3_stok`
- setup: seed warehouse stock, seed final order line qty, panggil `GET /sales/v2/orders/:ro_no`
- expected result: `qty*_stok` sama dengan warehouse stock current yang sudah dikonversi + final line qty yang sedang ditampilkan
- risk yang dicegah: patch memperbaiki ledger tetapi contract UI detail tetap salah

#### `Contract_DetailV2_FinalOrder_QtyStok_ChangesAfterCancelAccordingToSource`
- tujuan: memastikan sesudah cancel, response detail order ikut berubah sesuai source-of-truth yang didokumentasikan
- setup: seed reserving order, ambil response before, cancel order, ambil response after
- expected result: `details_final.normal[].qty*_stok` berubah konsisten terhadap warehouse stock source dan keputusan bisnis final untuk order cancelled
- risk yang dicegah: source warehouse sudah berubah tetapi endpoint detail masih menampilkan nilai stale atau formula bisnis salah

#### `Contract_DistributorSummary_OnlyIfConfirmedInDiscovery`
- tujuan: memastikan summary hanya jadi bagian patch jika memang confirmed
- setup: test hanya aktif bila discovery flag menyatakan summary bagian dari read path
- expected result: tidak ada distributor summary assertion jika discovery tidak confirm
- risk yang dicegah: scope creep tersembunyi

### E. Rollback or transaction tests

#### `Integration_CancelOrder_RollbackWhenProjectionReconcileFails`
- tujuan: memastikan atomicity penuh
- setup: paksa reconcile gagal setelah reverse dan status update step
- expected result: reverse rows tidak committed, status tidak berubah, projection tidak berubah
- risk yang dicegah: partial success corruption

#### `Integration_CancelOrder_AllReadsAndWritesUseSameTxCtx`
- tujuan: memastikan reconcile membaca state pada transaction yang sama
- setup: tx-aware assertions pada repository layer
- expected result: semua method menerima `txCtx`
- risk yang dicegah: read-your-own-write inconsistency

#### `Integration_CancelOrder_CommitOnlyAfterAllReconcilePass`
- tujuan: memastikan commit hanya terjadi di akhir
- setup: semua reconcile pass lalu verifikasi persisted state
- expected result: state konsisten atomik
- risk yang dicegah: premature commit

---

## 10. Missing Test Cases yang Sering Terlewat

Edge cases yang tetap wajib dipantau:

- **mixed item types**
  - detail normal dan promo atau reward bercampur

- **promo or reward item present**
  - reward tidak boleh dianggap reservation utama

- **zero qty detail**
  - detail ada tetapi qty nol

- **stale summary with already reversed rows**
  - projection sudah salah sebelum retry cancel

- **multi warehouse possibility**
  - basis row bisa berisi lebih dari satu `wh_id`

- **nil field or null status**
  - `data_status == nil`
  - `wh_id == nil`
  - `ro_date == nil`

- **deterministic ordering issue from map dedup**
  - hasil rebuild harus stabil

- **rollback behavior when reconcile fails**
  - reverse + status update tidak boleh committed jika reconcile gagal

- **all rows already reversed but order not yet cancelled**
  - patch harus tetap bisa menyelesaikan status dan projection tanpa duplicate reverse

- **basis row exists but detail has become inactive**
  - patch-week handling harus documented dan deterministic

---

## 10A. Final Decision for `Available Stock` After Cancel

### Rule final yang dikunci untuk patch minggu ini
- Saat order berstatus `CANCELLED`, formula display `Available Stock` **tidak boleh lagi menambahkan displayed line qty order tersebut**
- Dengan kata lain, untuk order `CANCELLED`, expected UI value adalah:
  - `Available Stock = warehouse stock current saja`
- Rule ini berlaku untuk semua tab yang menampilkan line dari order yang sama:
  - `Sales Order`
  - `Final Order`
  - `Purchase Order`
- Bila endpoint detail tetap mengembalikan historical line rows untuk tujuan audit atau visibility, qty line tersebut **boleh tetap tampil**, tetapi **tidak boleh ikut menambah `qty*_stok` display**

### Alasan principal-level
- Setelah cancel, order tidak lagi memegang reservation yang valid secara bisnis
- Reverse flow `CancelSalesStockUpdates` sudah menjadi single owner untuk mengembalikan stock reservation dan `qty_on_order`
- Jika endpoint detail masih menghitung `warehouse stock current + displayed line qty`, maka UI akan melakukan double-count terhadap stock yang sebenarnya sudah dilepas saat cancel
- Itu membuat response tidak lagi merepresentasikan availability real-time, melainkan availability semu yang masih menganggap order cancelled sebagai reservasi aktif
- Rule ini paling production-safe karena selaras dengan definisi inventory pasca cancel: order cancelled tidak punya hak reserve, maka line qty cancelled tidak boleh memengaruhi `Available Stock`

### Implementasi rule di patch
- `DetailV2` tetap boleh memakai formula existing untuk order aktif atau non-cancelled sesuai tab behavior saat ini
- Tetapi ketika header `data_status = CANCELLED`, service harus branch ke rule khusus:
  - `qty1_stok`, `qty2_stok`, `qty3_stok` dibentuk dari warehouse stock current yang sudah dikonversi
  - tanpa tambahan qty row order
- Rule ini harus diperlakukan sebagai contract UI, bukan sekadar side effect implementasi

### Expected UI value after cancel
- `Final Order` tab: `details_final.normal[].qty*_stok = warehouse stock current converted`
- `Sales Order` tab: `details.normal[].qty*_stok = warehouse stock current converted`
- `Purchase Order` tab: jika field ini ikut diproyeksikan atau diwariskan oleh endpoint detail, nilainya juga harus mengikuti warehouse stock current converted tanpa tambahan qty cancelled order

### Catatan scope control
- Patch minggu ini **tidak** mengubah rule untuk order aktif
- Patch minggu ini **hanya** mengunci exception business rule untuk status `CANCELLED`
- Ini menjaga blast radius kecil dan langsung menutup mismatch utama pada UI setelah cancel

---

## 10B. Focused Discovery Plan Khusus `On Cust Order`

Tujuan section ini adalah **mengunci source-of-truth UI untuk metric `On Cust Order`** sebelum coding dimulai.
Section ini tidak membahas redesign patch lain; fokus hanya pada menemukan read path UI yang benar-benar dipakai QA.

### 10B.1 Candidate UI Locations paling mungkin

Urutan kandidat berikut disusun dari evidence code yang paling kuat saat ini.

#### Kandidat 1. Distributor stock page atau stock per warehouse page
Alasan relevan:
- repository inventory sudah secara eksplisit membaca `qty_on_order` dari `inv.warehouse_stock`
- query paling kuat ada pada `inventory/repository/warehouse_stock_repository.go` di method `FindAllByCustId`
- field response `qty_on_order` sudah tersedia pada DTO `DistributorStockList`
- ini cocok dengan UI tabel stock yang biasanya menampilkan `Qty`, `On Order`, `On Shipping`

Evidence code yang sudah ada:
- `inventory/repository/warehouse_stock_repository.go` method `FindAllByCustId`
- select field mencakup `whs.qty_on_order`
- `inventory/entity/warehouse_stock.go` type `DistributorStockList`

#### Kandidat 2. Warehouse stock detail atau product stock popup per warehouse
Alasan relevan:
- flow inventory memiliki endpoint product list pada warehouse stock controller
- meski query product list saat ini lebih condong menampilkan `qty`, screen jenis popup sering memanggil endpoint tambahan atau memakai DTO turunan yang masih bersumber dari warehouse stock
- bila QA melihat metric saat membuka stock popup dari order flow, kandidat ini harus dicek lebih dulu

Evidence code yang sudah ada:
- `inventory/controller/warehouse_stock_controller.go` method `ProductList`
- `inventory/service/warehouse_stock_service.go` method `ProductList`
- `inventory/repository/warehouse_stock_repository.go` method `ProductList`

#### Kandidat 3. Product lookup modal saat create atau edit order
Alasan relevan:
- label bisnis seperti `On Cust Order` sering dipakai untuk membantu salesman atau admin melihat reservation saat memilih produk
- lookup mode dapat memakai source berbeda dari stock list biasa
- kandidat ini penting bila QA melihat metric bukan di inventory page melainkan di flow order

Evidence code yang sudah ada:
- trace price or lookup repository perlu diprioritaskan ke query lookup yang sudah mengandung `qty_on_order`

#### Kandidat 4. Distributor price lookup atau product lookup berbasis pricing
Alasan relevan:
- repository master sudah memiliki query lookup yang menghitung `qty_on_order` lewat agregasi `inv.wh_stock`
- jika UI order memakai lookup ini, maka source metric bukan direct table warehouse stock, melainkan summary or view `inv.wh_stock`

Evidence code yang sudah ada:
- `master/repository/dist_price_repository.go` method `FindAllByCustIdLookupMode`
- query menghitung `SUM(COALESCE(wh.qty_on_order, 0))` dari `inv.wh_stock`

#### Kandidat 5. Stock summary component atau dashboard inventory widget
Alasan relevan:
- bila QA melihat metric di layar ringkasan stock, read path bisa datang dari summary component, bukan dari order detail
- kandidat ini lebih lemah dibanding warehouse stock page, tetapi tetap harus dicatat sebagai fallback bila network trace tidak mengarah ke endpoint list biasa

#### Kandidat 6. Order detail atau final order page
Alasan relevan:
- defect ticket berasal dari flow cancel order
- namun sampai saat ini belum ada bukti bahwa `On Cust Order` dibaca dari endpoint detail order seperti `DetailV2`
- karena itu kandidat ini relevan secara bisnis, tetapi belum kuat secara teknis

### 10B.2 Discovery strategy paling efisien

Urutan langkah wajib agar discovery minim trial-and-error:

1. **Kunci screen exact terlebih dahulu**
   - jangan mulai dari grep backend tanpa tahu screen exact
   - output: nama menu, page, tab, komponen, dan screenshot marker

2. **Tangkap network request saat metric tampil**
   - buka screen yang memuat `On Cust Order`
   - filter response JSON yang memuat nilai metric tersebut
   - output: endpoint utama dan endpoint pendukung

3. **Trace route ke controller**
   - cari route registration dari endpoint hasil network trace
   - output: controller file dan handler exact

4. **Trace controller ke service**
   - pastikan service hanya pass-through atau ada enrichment field
   - output: service method dan helper mapping yang membentuk metric

5. **Trace service ke repository dan query**
   - di sini source metric biasanya terkunci
   - output: repository method, select clause, join, subquery, dan table atau view yang dipakai

6. **Klasifikasikan source table atau view**
   - salah satu kategori berikut harus dipilih:
     - `inv.warehouse_stock.qty_on_order`
     - `inv.wh_stock.qty_on_order`
     - aggregate order aktif
     - summary distributor
     - projection atau view lain

7. **Cek dependency terhadap status order**
   - jika read path memakai stock precomputed, dependency status kemungkinan ada di write path
   - jika read path memakai aggregate order query, dependency status ada langsung di query read
   - output: daftar status yang dihitung atau dikecualikan

8. **Cek dimensionality metric**
   - jika endpoint butuh `wh_id`, metric cenderung warehouse-scoped
   - jika query meng-`SUM` lintas warehouse hanya dengan `cust_id`, metric cenderung distributor-scoped
   - output: warehouse-scoped atau distributor-scoped

9. **Putuskan patch relevance**
   - metric patch-relevant bila memang muncul pada screen defect atau endpoint validasi pasca cancel
   - output: `needs_reconcile = yes or no`

### 10B.3 Exact developer checklist

Checklist operasional yang harus dijalankan developer:

#### A. UI dan network
- cari label `On Cust Order` di frontend atau API mapping
- catat screen exact tempat QA melihat metric
- capture request network sebelum cancel dan sesudah cancel
- simpan satu sample response JSON yang memuat metric target

#### B. Keyword search backend
Jalankan pencarian dengan keyword berikut:
- `qty_on_order`
- `qtyOnOrder`
- `on order`
- `on cust order`
- `on_cust_order`
- `on_customer_order`
- `customer_order`
- `reserved_stock`
- `reservation`
- `available_stock`
- `warehouse_stock`
- `wh_stock`
- `qty_order1`
- `qty_inc_on_order1`
- `total_qty_inc_on_order`

#### C. Prioritas file yang harus dicek lebih dulu
- `inventory/repository/warehouse_stock_repository.go`
- `inventory/service/warehouse_stock_service.go`
- `inventory/controller/warehouse_stock_controller.go`
- `master/repository/dist_price_repository.go`
- `inventory/repository/stock_repository.go`
- `sales/repository/stock_repository.go`

#### D. Checklist trace code
- dari endpoint hasil network trace, cari route exact
- catat controller handler
- catat service method
- catat repository method
- screenshot atau copy select clause yang memuat source metric
- tentukan field exact yang mengisi metric UI

#### E. Checklist validasi bisnis
- apakah metric berasal dari stock precomputed atau aggregate order aktif
- apakah query read memiliki filter `data_status`
- apakah order `CANCELLED` masih ikut dihitung
- apakah metric warehouse-scoped atau distributor-scoped
- apakah metric dibaca dari summary distributor atau tidak

#### F. Checklist validasi patch relevance
- apakah screen ini benar screen yang dipakai QA pada defect
- apakah source sudah sinkron setelah reverse cancel
- apakah masih perlu reconcile layer lain setelah source ketemu

### 10B.4 Search keywords dan heuristik

#### Keyword utama
- `qty_on_order`
- `qtyOnOrder`
- `on_order`
- `on cust order`
- `on_cust_order`
- `on_customer_order`
- `customer_order`
- `reserved_stock`
- `reservation`
- `available_stock`
- `warehouse_stock`
- `wh_stock`

#### DTO atau response alternatif yang patut dicurigai
- `DistributorStockList`
- `WarehouseStock`
- `ProductWarehouseList`
- `StockReport`
- `DistPriceLookup`
- `qty_order1`
- `qty_order2`
- `qty_order3`
- `qty_inc_on_order1`
- `qty_inc_on_order2`
- `qty_inc_on_order3`
- `total_qty_inc_on_order`

#### Nama kolom atau istilah alternatif di DB atau code
- `qty_on_order`
- `qty_order`
- `on_order_qty`
- `customer_order_qty`
- `reserved_qty`
- `allocation_qty`
- `booked_stock`

#### Heuristik cepat untuk narrowing
- bila field muncul bersama `qty` dan `qty_on_shipping`, source paling mungkin adalah stock table atau stock view
- bila query memakai `SUM(wh.qty_on_order)`, source kemungkinan distributor-wide atau lookup lintas warehouse
- bila field select langsung `whs.qty_on_order` dengan `wh_id`, source kemungkinan warehouse-scoped
- bila UI memakai `qty_order1..3` atau `qty_inc_on_order1..3`, kemungkinan UI membaca endpoint report atau endpoint dengan conversion layer, bukan raw stock field langsung

### 10B.5 Decision tree menentukan source-of-truth

```text
Apakah UI membaca field precomputed stock?
├─ Ya, dari inv.warehouse_stock.qty_on_order
│  ├─ Implikasi: bug utama paling mungkin ada di write path cancel atau reverse
│  ├─ Patch area: reverse writer, idempotency, tx boundary
│  └─ Validasi: row warehouse target harus berubah setelah cancel
├─ Ya, dari inv.wh_stock.qty_on_order
│  ├─ Implikasi: ada projection atau view layer di atas warehouse stock
│  ├─ Patch area: sync source projection atau pastikan view source ikut berubah
│  └─ Validasi: bandingkan inv.wh_stock dengan inv.warehouse_stock
├─ Tidak, dari aggregate active order
│  ├─ Implikasi: bug ada pada filter status order di query read
│  ├─ Patch area: query read side dan exclusion status cancelled
│  └─ Validasi: compare active order aggregate dengan metric UI
├─ Tidak, dari computed projection lain
│  ├─ Implikasi: query atau mapper membentuk metric sendiri
│  ├─ Patch area: service or repository projection logic
│  └─ Validasi: cocokkan formula projection dengan state pasca cancel
└─ Ya, dari distributor summary
   ├─ Implikasi: warehouse reverse bisa benar tetapi summary tetap stale
   ├─ Patch area: summary reconcile
   └─ Validasi: warehouse source benar, summary source masih salah
```

### 10B.6 Evidence template wajib diisi

| UI screen/component | Endpoint | Controller | Service | Repository | Table or view | Exact field name | Depends on order status Y or N | Depends on warehouse Y or N | Depends on distributor summary Y or N | Patch relevant Y or N | Notes |
|---|---|---|---|---|---|---|---|---|---|---|---|
| `On Cust Order` |  |  |  |  |  |  |  |  |  |  |  |

Kolom ini wajib diisi lengkap sebelum coding:
- UI screen/component
- endpoint
- controller
- service
- repository
- table or view
- exact field name
- depends on order status Y or N
- depends on warehouse Y or N
- depends on distributor summary Y or N
- patch relevant Y or N
- notes

### 10B.7 SQL validation template

#### Jika source adalah `inv.warehouse_stock`
```sql
SELECT
  cust_id,
  wh_id,
  pro_id,
  qty,
  qty_on_order,
  qty_on_shipping,
  updated_at
FROM inv.warehouse_stock
WHERE cust_id = :cust_id
  AND wh_id = :wh_id
  AND pro_id = :pro_id;
```

#### Jika source adalah `inv.wh_stock`
```sql
SELECT
  cust_id,
  wh_id,
  pro_id,
  qty,
  qty_on_order,
  qty_on_shipping
FROM inv.wh_stock
WHERE cust_id = :cust_id
  AND (:wh_id IS NULL OR wh_id = :wh_id)
  AND pro_id = :pro_id;
```

#### Jika source adalah aggregate order aktif
```sql
SELECT
  oh.cust_id,
  oh.wh_id,
  od.pro_id,
  SUM(COALESCE(od.qty, 0)) AS active_order_qty
FROM sales.order_header oh
JOIN sales.order_detail od ON od.ro_no = oh.ro_no
WHERE oh.cust_id = :cust_id
  AND oh.wh_id = :wh_id
  AND od.pro_id = :pro_id
  AND oh.data_status IN (:active_statuses)
  AND oh.data_status <> :cancelled_status
GROUP BY oh.cust_id, oh.wh_id, od.pro_id;
```

#### Jika source adalah distributor stock summary
```sql
SELECT
  cust_id,
  distributor_id,
  pro_id,
  qty_on_order,
  updated_at
FROM <schema>.<distributor_stock_summary>
WHERE cust_id = :cust_id
  AND distributor_id = :distributor_id
  AND pro_id = :pro_id;
```

#### Query pembanding source read vs expected active orders
```sql
SELECT
  ws.cust_id,
  ws.wh_id,
  ws.pro_id,
  ws.qty_on_order AS source_qty_on_order,
  COALESCE(ord.active_qty, 0) AS expected_active_qty
FROM inv.warehouse_stock ws
LEFT JOIN (
  SELECT
    oh.cust_id,
    oh.wh_id,
    od.pro_id,
    SUM(COALESCE(od.qty, 0)) AS active_qty
  FROM sales.order_header oh
  JOIN sales.order_detail od ON od.ro_no = oh.ro_no
  WHERE oh.cust_id = :cust_id
    AND oh.data_status IN (:active_statuses)
    AND oh.data_status <> :cancelled_status
  GROUP BY oh.cust_id, oh.wh_id, od.pro_id
) ord
  ON ord.cust_id = ws.cust_id
 AND ord.wh_id = ws.wh_id
 AND ord.pro_id = ws.pro_id
WHERE ws.cust_id = :cust_id
  AND ws.wh_id = :wh_id
  AND ws.pro_id = :pro_id;
```

### 10B.8 Acceptance criteria untuk discovery `On Cust Order`

Discovery dianggap selesai hanya jika semua poin berikut terpenuhi:
- screen exact tempat QA melihat metric sudah diketahui
- endpoint exact yang mengisi metric sudah diketahui
- controller, service, dan repository read path sudah diketahui
- source table atau view final sudah diketahui
- exact field name backend sudah diketahui
- sudah diketahui apakah metric warehouse-scoped atau distributor-scoped
- sudah diketahui apakah dependency terhadap `data_status` terjadi di read path atau hanya di write path
- sudah diketahui apakah order `CANCELLED` masih ikut dihitung atau tidak
- sudah ada keputusan apakah source ini patch-relevant minggu ini
- sudah tersedia minimal satu query SQL validasi yang bisa dipakai membandingkan actual vs expected

### 10B.9 Output final yang wajib diserahkan developer

Developer wajib menyerahkan hasil discovery berikut:
- mapping table lengkap untuk `On Cust Order`
- route trace lengkap: endpoint → controller → service → repository
- source SQL atau select clause exact
- klasifikasi source-of-truth:
  - `inv.warehouse_stock`
  - `inv.wh_stock`
  - aggregate active order
  - distributor summary
  - projection lain
- keputusan patch relevance
- rekomendasi test case endpoint source aktual

### 10B.10 Evidence backend paling kuat saat ini

Sebelum network trace FE dikonfirmasi, evidence backend terkuat saat ini adalah:
- candidate source pertama: `inventory/repository/warehouse_stock_repository.go` method `FindAllByCustId` yang membaca `whs.qty_on_order` dari `inv.warehouse_stock`
- candidate source kedua: `master/repository/dist_price_repository.go` method `FindAllByCustIdLookupMode` yang membaca agregat `SUM(wh.qty_on_order)` dari `inv.wh_stock`

Kesimpulan sementara:
- discovery harus dipusatkan dulu ke apakah QA melihat metric pada stock list atau lookup yang membaca `qty_on_order` precomputed
- jangan mengasumsikan `On Cust Order` berasal dari detail order endpoint sampai network trace membuktikannya

---

## 10C. Final Decision Gate untuk Distributor Summary

### Rule final
Distributor summary **hanya dianggap relevan** bila ada bukti bahwa metric pada screen defect membaca langsung atau tidak langsung dari distributor-level summary tersebut.

### Bukti minimal untuk memasukkan distributor summary ke patch
Semua bukti berikut harus ada:
- screen atau endpoint yang terdampak memang memakai distributor summary
- controller ke service ke repository trace menunjukkan query menuju distributor summary table atau view
- ada bukti bahwa setelah cancel, distributor summary source tetap stale walaupun reverse flow warehouse sudah benar
- ada hubungan jelas bahwa mismatch UI tidak bisa dijelaskan hanya oleh `DetailV2` source atau warehouse stock source

### Kapan harus tetap skip
Distributor summary wajib tetap di-skip bila salah satu kondisi berikut terjadi:
- metric UI yang gagal tidak membaca distributor summary
- trace code berhenti di warehouse-level source atau order-level projection
- distributor summary hanya dipakai untuk reporting downstream dan bukan live UI path
- belum ada bukti stale state pada distributor summary setelah cancel
- memasukkan distributor summary hanya berdasarkan dugaan tanpa endpoint atau query evidence

### Keputusan patch-week
- default value tetap: `distributor_summary_reconcile = false`
- perubahan ke `true` membutuhkan bukti minimal lengkap seperti di atas
- jika bukti belum lengkap saat coding dimulai, patch tetap jalan tanpa distributor summary

---

## 10D. Acceptance Criteria Final

### Available Stock after cancel
- setelah cancel sukses, endpoint detail order yang dipakai UI harus mengembalikan `qty1_stok`, `qty2_stok`, `qty3_stok` yang merepresentasikan warehouse stock current yang sudah dikonversi
- untuk order dengan `data_status = CANCELLED`, line qty order tersebut tidak boleh lagi menambah nilai `Available Stock`
- nilai pada tab `Final Order` harus sinkron dengan source warehouse stock yang sama transactionally setelah reverse cancel selesai
- bila tab lain memakai source field yang sama dari endpoint detail, hasilnya juga harus mengikuti rule cancelled yang sama

### On Cust Order after cancel
- setelah cancel sukses, metric `On Cust Order` pada screen atau endpoint yang dipakai QA harus turun atau kembali sesuai state order aktif yang tersisa
- order yang sudah `CANCELLED` tidak boleh lagi berkontribusi ke metric `On Cust Order`
- bila source metric adalah summary atau projection, nilai summary harus sinkron dengan source-of-truth pasca reverse cancel dalam boundary transaksi yang sama atau dengan mekanisme reconcile yang approved
- jika discovery membuktikan UI source berbeda dari asumsi awal, patch hanya dianggap selesai bila source aktual tersebut ikut tervalidasi

---

## 10E. Tambahan Test Cases Penutup Blocker

### `DetailV2` behavior when order is `CANCELLED`
#### `Contract_DetailV2_CancelledOrder_AvailableStockUsesWarehouseCurrentOnly`
- tujuan: mengunci rule final bahwa order cancelled tidak menambah line qty ke `Available Stock`
- setup: seed order cancelled dengan detail final qty non-zero dan warehouse stock current known value
- expected result: `details_final.normal[].qty*_stok` sama dengan warehouse stock current converted saja
- risk yang dicegah: double-count stock setelah cancel

#### `Contract_DetailV2_NonCancelledOrder_AvailableStockKeepsExistingBehavior`
- tujuan: memastikan patch tidak mengubah behavior order aktif
- setup: seed order processed atau active dengan detail qty known value
- expected result: formula existing tetap berjalan untuk non-cancelled path
- risk yang dicegah: regression pada screen order aktif

### `On Cust Order` endpoint or source after cancel
#### `Contract_OnCustOrder_Source_DoesNotCountCancelledOrder`
- tujuan: mengunci bahwa source `On Cust Order` tidak lagi memasukkan cancelled order
- setup: seed beberapa order untuk SKU sama, cancel satu order, query endpoint atau source aktual
- expected result: metric hanya menghitung order aktif yang tersisa
- risk yang dicegah: stale reservation visibility setelah cancel

#### `Integration_CancelOrder_OnCustOrder_UIPath_ReflectsPostCancelState`
- tujuan: memastikan endpoint UI yang benar-benar dipakai QA ikut berubah setelah cancel
- setup: ambil response sebelum cancel, lakukan cancel, ambil response sesudah cancel
- expected result: field `On Cust Order` berubah sesuai source-of-truth yang didokumentasikan
- risk yang dicegah: reverse flow benar tetapi UI path tidak sinkron

### guard jika UI source berbeda dengan asumsi awal
#### `Guard_DiscoveryFailsIfOnCustOrderSourceDiffersFromPlan`
- tujuan: mencegah dev mengimplementasikan patch ke source yang salah
- setup: discovery menemukan endpoint atau table berbeda dari asumsi awal plan
- expected result: implementasi reconcile untuk `On Cust Order` tidak dilanjutkan sebelum discovery table dan acceptance criteria di-update
- risk yang dicegah: false fix pada summary atau table yang tidak dipakai UI

#### `Guard_DistributorSummaryRemainsSkippedWithoutEvidence`
- tujuan: mencegah scope creep saat source `On Cust Order` belum mengarah ke distributor summary
- setup: discovery tidak menemukan evidence distributor summary di UI path
- expected result: patch tetap skip distributor summary
- risk yang dicegah: perubahan tidak perlu pada area berisiko tinggi

---

## 11. Reviewer Verdict

### Apakah plan ini sudah ready for coding?
**Secara desain patch, plan ini sudah implementation-ready.**
Satu-satunya dependency yang masih pending sebelum coding dimulai adalah **discovery source-of-truth UI untuk `On Cust Order`** dan konfirmasi akhir apakah distributor summary relevan atau tetap bisa di-skip.

Keputusan final untuk `Available Stock` saat order `CANCELLED` **sudah terkunci** di [`## 10A. Final Decision for `Available Stock` After Cancel`](plans/sx-1241-cancel-order-patch-plan.md).
Karena itu, `Available Stock` **bukan lagi blocker**.

### Blocker terakhir
Blocker yang tersisa hanya dua hal berikut:

1. **Source-of-truth `On Cust Order` belum terkunci sampai read path UI level**
   Yang belum diketahui secara final:
   - screen atau component exact yang dipakai QA
   - endpoint exact yang mengisi metric
   - controller, service, dan repository read path yang benar
   - source table atau view final
   - apakah metric warehouse-scoped atau distributor-scoped
   - apakah dependency terhadap order status ada di read path atau hanya di write path

2. **Relevansi distributor summary belum confirmed**
   Yang belum diketahui secara final:
   - apakah UI source `On Cust Order` membaca distributor summary secara langsung atau tidak langsung
   - apakah mismatch pasca cancel masih terjadi walaupun source warehouse atau source utama sudah benar

### Evidence minimum agar discovery `On Cust Order` dianggap selesai
Discovery `On Cust Order` baru dianggap selesai bila minimal evidence berikut sudah tersedia:
- screen exact tempat QA melihat metric sudah diketahui
- endpoint exact yang dipakai UI sudah diketahui
- trace lengkap `controller → service → repository` sudah diketahui
- source table atau view final sudah diketahui
- exact field name backend yang mengisi metric sudah diketahui
- sudah diketahui apakah metric warehouse-scoped atau distributor-scoped
- sudah diketahui apakah cancelled order masih ikut dihitung atau tidak
- sudah ada keputusan eksplisit `patch relevant = yes or no`

### Kapan patch boleh mulai coding meskipun distributor summary tetap skip
Patch **boleh mulai coding** segera setelah:
- discovery `On Cust Order` memenuhi evidence minimum di atas
- hasil discovery menunjukkan UI source tidak membaca distributor summary
  **atau**
- tidak ada evidence endpoint-to-query bahwa distributor summary adalah bagian dari failing read path

Dengan kata lain:
- `distributor_summary_reconcile` boleh tetap `false`
- patch tetap sah untuk diimplementasikan selama source aktual `On Cust Order` sudah terkunci dan distributor summary tidak terbukti relevan

### Jika blocker discovery selesai
Maka plan ini **ready for coding** dengan guardrail wajib berikut:
- tetap gunakan `BulkUpdateStatus` sebagai orchestration patch-week
- tetap gunakan `CancelSalesStockUpdates` sebagai single owner untuk reverse ledger, `inv.warehouse_stock.qty`, dan `inv.warehouse_stock.qty_on_order`
- jangan biarkan reconcile overwrite `qty_on_order`
- distributor summary default skip sampai ada evidence yang mewajibkan inclusion
- semua reconcile harus scoped by deterministic affected key
- semua error reconcile harus rollback entire transaction
- semua repository methods harus konsisten memakai `txCtx`

### Statement final
**Plan SX-1241 sudah implementation-ready pending `On Cust Order` discovery only.**
Jika discovery mengunci source UI `On Cust Order` dan tidak menemukan evidence kuat untuk distributor summary, maka coding patch dapat langsung dimulai tanpa redesign tambahan.

### Kesimpulan principal reviewer
- arah patch sudah benar
- scope patch sudah cukup terkontrol
- ownership `qty_on_order` sudah dianggap final dan terkunci
- keputusan final `Available Stock` pasca cancel sudah locked dan bukan blocker lagi
- blocker nyata sebelum coding hanya tersisa penguncian source-of-truth `On Cust Order` dan konfirmasi apakah distributor summary relevan atau tetap skip
- setelah discovery itu selesai, plan ini layak langsung masuk tahap implementasi dengan regression risk yang terkontrol
